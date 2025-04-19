package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func (a Router) Chat(c *gin.Context) {
	var req struct {
		Model    string                         `json:"model"`
		Messages []openai.ChatCompletionMessage `json:"messages"`
		Stream   *bool                          `json:"stream"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	streamRequested := true
	if req.Stream != nil {
		streamRequested = *req.Stream
	}

	if !streamRequested {
		fullModelName, err := a.prv.GetFullModelName(c, req.Model)
		if err != nil {
			slog.ErrorContext(c, "error getting full model name", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		res, err := a.prv.Chat(c, req.Messages, fullModelName)
		if err != nil {
			slog.ErrorContext(c, "failed to get chat res", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(res.Choices) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no res from model"})
			return
		}

		content := ""
		if len(res.Choices) > 0 && len(res.Choices[0].Message.Content) != 0 {
			content = res.Choices[0].Message.Content
		}

		finishReason := "stop"
		if len(res.Choices[0].FinishReason) != 0 {
			finishReason = string(res.Choices[0].FinishReason)
		}

		ollamaResponse := map[string]interface{}{
			"model":      fullModelName,
			"created_at": time.Now().Format(time.RFC3339),
			"message": map[string]string{
				"role":    "assistant",
				"content": content,
			},
			"done":              true,
			"finish_reason":     finishReason,
			"total_duration":    res.Usage.TotalTokens * 10,
			"load_duration":     0,
			"prompt_eval_count": res.Usage.PromptTokens,
			"eval_count":        res.Usage.CompletionTokens,
			"eval_duration":     res.Usage.CompletionTokens * 10,
		}

		c.JSON(http.StatusOK, ollamaResponse)
		return
	}

	slog.InfoContext(c, "requested model", "model", req.Model)
	fullModelName, err := a.prv.GetFullModelName(c, req.Model)
	if err != nil {
		slog.ErrorContext(c, "error getting full model name", "error", err, "model", req.Model)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	slog.InfoContext(c, "using model", "fullModelName", fullModelName)

	stream, err := a.prv.ChatStream(c, req.Messages, fullModelName)
	if err != nil {
		slog.ErrorContext(c, "failed to create stream", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()

	c.Writer.Header().Set("Content-Type", "application/x-ndjson")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	var lastFinishReason string
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			slog.ErrorContext(c, "stream error", "error", err)
			rawMessage, err := json.Marshal(map[string]string{"error": "stream error: " + err.Error()})
			if err != nil {
				slog.ErrorContext(c, "error marshaling stream error", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fmt.Fprintf(c.Writer, "%s\n", string(rawMessage))
			c.Writer.Flush()
			return
		}

		if len(res.Choices) > 0 && len(res.Choices[0].FinishReason) != 0 {
			lastFinishReason = string(res.Choices[0].FinishReason)
		}

		rawMessage, err := json.Marshal(map[string]interface{}{
			"model":      fullModelName,
			"created_at": time.Now().Format(time.RFC3339),
			"message": map[string]string{
				"role":    "assistant",
				"content": res.Choices[0].Delta.Content,
			},
			"done": false,
		})
		if err != nil {
			slog.ErrorContext(c, "error marshaling", "error", err)
			return
		}

		fmt.Fprintf(c.Writer, "%s\n", string(rawMessage))
		c.Writer.Flush()
	}

	if lastFinishReason == "" {
		lastFinishReason = "stop"
	}

	rawMessage, err := json.Marshal(map[string]interface{}{
		"model":      fullModelName,
		"created_at": time.Now().Format(time.RFC3339),
		"message": map[string]string{
			"role":    "assistant",
			"content": "",
		},
		"done":              true,
		"finish_reason":     lastFinishReason,
		"total_duration":    0,
		"load_duration":     0,
		"prompt_eval_count": 0,
		"eval_count":        0,
		"eval_duration":     0,
	})
	if err != nil {
		slog.ErrorContext(c, "error marshaling", "error", err)
		return
	}

	fmt.Fprintf(c.Writer, "%s\n", string(rawMessage))
	c.Writer.Flush()
}
