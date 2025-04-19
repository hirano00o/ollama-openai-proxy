package provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	"github.com/hirano00o/ollama-openai-proxy/model"
)

type Provider struct {
	client *openai.Client
	models []string
}

func NewProvider(apiKey string) *Provider {
	return &Provider{
		client: openai.NewClient(apiKey),
		models: []string{},
	}
}

func (p *Provider) Chat(ctx *gin.Context, messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error) {
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	return p.client.CreateChatCompletion(ctx, req)
}

func (p *Provider) ChatStream(ctx *gin.Context, messages []openai.ChatCompletionMessage, model string) (*openai.ChatCompletionStream, error) {
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}
	return p.client.CreateChatCompletionStream(ctx, req)
}

func (p *Provider) GetModels(ctx *gin.Context) ([]model.Model, error) {
	currentTime := time.Now().Format(time.RFC3339)

	modelList, err := p.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	p.models = []string{}

	models := make([]model.Model, 0, len(modelList.Models))
	for _, m := range modelList.Models {
		parts := strings.Split(m.ID, "/")
		name := parts[len(parts)-1]

		p.models = append(p.models, m.ID)

		models = append(models, model.Model{
			Name:       name,
			Model:      name,
			ModifiedAt: currentTime,
			Size:       0,
			Digest:     name,
			Details: model.Details{
				ParentModel:       "dummy",
				Format:            "gguf",
				Family:            "dummy",
				Families:          []string{"dummy"},
				ParameterSize:     "0B",
				QuantizationLevel: "dummy",
			},
		})
	}

	return models, nil
}

func (p *Provider) GetModelDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"license":    "dummy",
		"system":     "dummy",
		"modifiedAt": time.Now().Format(time.RFC3339),
		"details": map[string]interface{}{
			"format":             "gguf",
			"parameter_size":     "0B",
			"quantization_level": "",
		},
		"model_info": map[string]interface{}{
			"architecture":    "dummy",
			"context_length":  0,
			"parameter_count": 0,
		},
	}, nil
}

func (p *Provider) GetFullModelName(ctx *gin.Context, alias string) (string, error) {
	if len(p.models) == 0 {
		if _, err := p.GetModels(ctx); err != nil {
			return "", fmt.Errorf("failed to get models: %w", err)
		}
	}

	for _, name := range p.models {
		if name == alias {
			return name, nil
		}
	}

	for _, name := range p.models {
		if strings.HasPrefix(name, alias) {
			return name, nil
		}
	}

	return alias, nil
}
