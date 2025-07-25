package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/utils/httpx"
	"github.com/y3933y3933/joker/internal/utils/param"
)

type QuestionHandler struct {
	questionService *service.QuestionService
	logger          *slog.Logger
}

func NewQuestionHandler(logger *slog.Logger, questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		logger:          logger,
		questionService: questionService,
	}
}

type QuestionResponse struct {
	Question store.Question
}

func (h *QuestionHandler) HandleGetPaginatedQuestions(c *gin.Context) {
	params := h.parseQueryParams(c)

	if err := h.questionService.ValidateQuestionParams(params); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}
	result, err := h.questionService.ListQuestions(c.Request.Context(), params)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, result)

}

func (h *QuestionHandler) parseQueryParams(c *gin.Context) service.QuestionQueryParams {
	params := service.QuestionQueryParams{
		Keyword:  c.Query("keyword"),
		Level:    c.Query("level"),
		SortBy:   c.Query("sort_by"),
		Page:     1,
		PageSize: 10,
	}

	params.Page = param.ReadIntQuery(c, "page", 1)
	params.PageSize = param.ReadIntQuery(c, "page_size", 10)

	if params.SortBy == "" {
		params.SortBy = "created_at_desc"
	}

	return params
}

type createQuestionRequest struct {
	Level   string `json:"level" binding:"required,oneof=normal spicy"`
	Content string `json:"content" binding:"required"`
}

func (h *QuestionHandler) HandleCreateQuestion(c *gin.Context) {
	var req createQuestionRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	q, err := h.questionService.CreateQuestion(c.Request.Context(), req.Content, req.Level)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, q)
}

func (h *QuestionHandler) HandleDeleteQuestion(c *gin.Context) {
	id, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	err = h.questionService.DeleteQuestion(c.Request.Context(), id)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, nil)
}

type updateQuestionRequest struct {
	Level   *string `json:"level" binding:"oneof=normal spicy" `
	Content *string `json:"content"`
}

func (h *QuestionHandler) HandleUpdateQuestion(c *gin.Context) {
	id, err := param.ParseIntParam(c, "id")
	if err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	var req updateQuestionRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	q, err := h.questionService.UpdateQuestion(c.Request.Context(), id, req.Content, req.Level)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, q)
}
