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

func (h *QuestionHandler) GetPaginatedQuestions(c *gin.Context) {
	params := h.parseQueryParams(c)

	if err := h.questionService.ValidateParams(params); err != nil {
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
