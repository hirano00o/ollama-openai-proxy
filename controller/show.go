package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Router) Show(ctx *gin.Context) {
	var req map[string]string
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	modelName := req["name"]
	if len(modelName) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "model name is required"})
		return
	}

	details, err := a.prv.GetModelDetails()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, details)
}
