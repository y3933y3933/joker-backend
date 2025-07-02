package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

type FeedbackHandler struct {
	feedbackService *service.FeedbackService
	logger          *slog.Logger
}

func NewFeedbackHandler(logger *slog.Logger, feedbackService *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		logger:          logger,
		feedbackService: feedbackService,
	}
}

type createFeedbackRequest struct {
	Type    string `json:"type" binding:"required,oneof=feature issue other" `
	Content string `json:"content" binding:"required"`
}

func (h *FeedbackHandler) HandleCreateFeedback(c *gin.Context) {
	var req createFeedbackRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	err := h.feedbackService.CreateFeedback(c.Request.Context(), &store.Feedback{
		Content: req.Content,
		Type:    req.Type,
	})

	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, nil)
}
