package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Router) Tags(ctx *gin.Context) {
	models, err := a.prv.GetModels(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error getting models", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newModels := make([]map[string]interface{}, 0, len(models))
	for _, m := range models {
		newModels = append(newModels, map[string]interface{}{
			"name":        m.Name,
			"model":       m.Model,
			"modified_at": m.ModifiedAt,
			"size":        0,
			"digest":      "dummy",
			"details":     m.Details,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"models": newModels})
}
