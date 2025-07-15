package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
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

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) HandleRegisterUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, user)
}
