package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/httpx"
	"github.com/y3933y3933/joker/internal/utils/param"
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

func (h *FeedbackHandler) HandlerListFeedback(c *gin.Context) {
	params := h.parseQueryParams(c)

	if err := h.feedbackService.ValidateFeedbackParams(params); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	result, err := h.feedbackService.ListFeedback(c.Request.Context(), params)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, result)

}

func (h *FeedbackHandler) parseQueryParams(c *gin.Context) service.FeedbackQueryParams {
	params := service.FeedbackQueryParams{
		Type:         c.Query("type"),
		ReviewStatus: c.Query("reviewStatus"),
		Page:         1,
		PageSize:     10,
	}

	params.Page = param.ReadIntQuery(c, "page", 1)
	params.PageSize = param.ReadIntQuery(c, "page_size", 10)

	return params
}

func (h *FeedbackHandler) HandleGetFeedbackByID(c *gin.Context) {
	id, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	result, err := h.feedbackService.GetFeedbackByID(c.Request.Context(), id)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}
	httpx.SuccessResponse(c, result)

}

func (h *FeedbackHandler) HandleUpdateFeedbackReviewStatus(c *gin.Context) {
	var req struct {
		ReviewStatus string `json:"reviewStatus" binding:"required"`
	}

	id, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Error("Failed to bind request body", slog.Any("error", err))
		httpx.BadRequestResponse(c, err)
		return
	}

	err = h.feedbackService.UpdateFeedbackReviewStatus(c.Request.Context(), id, req.ReviewStatus)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, nil)

}
