package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func serverErrorResponse(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error(err.Error(), "method", c.Request.Method, "url", c.Request.URL)

	message := "the server encountered a problem and could not process your request"
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}

func successResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}
