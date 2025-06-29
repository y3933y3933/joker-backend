package httpx

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerErrorResponse(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error(err.Error(), "method", c.Request.Method, "url", c.Request.URL)

	message := "the server encountered a problem and could not process your request"
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}

func BadRequestResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func NotFoundResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func ForbiddenResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
}
