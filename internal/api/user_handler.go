package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

type UserHandler struct {
	userService *service.UserService
	logger      *slog.Logger
}

func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) HandlerGetUserInfo(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	user, err := h.userService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrUserNotFound):
			httpx.BadRequestResponse(c, err)
		default:
			httpx.ServerErrorResponse(c, h.logger, err)
		}
		return
	}

	httpx.SuccessResponse(c, user)

}
