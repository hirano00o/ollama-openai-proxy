package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/hirano00o/ollama-openai-proxy/controller"
	"github.com/hirano00o/ollama-openai-proxy/provider"
)

const host = ":11434"

func run() {
	r := gin.Default()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		slog.Error("OPENAI_API_KEY environment variable not set.")
		os.Exit(1)
	}
	p := provider.NewProvider(apiKey)
	c := controller.NewRouter(*p)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Ollama is running")
	})

	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	r.Group("/api").
		GET("/tags", c.Tags).
		POST("/show", c.Show).
		POST("/chat", c.Chat)

	if err := r.Run(host); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}

func main() {
	run()
}
