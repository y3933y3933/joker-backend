package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) HandleRegisterUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	user, err := h.authService.CreateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, user)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrInvalidCredentials):
			httpx.UnAuthorized(c, err)
		default:
			httpx.ServerErrorResponse(c, h.logger, err)
		}
		return
	}

	res := &LoginResponse{
		Token: token,
	}

	httpx.SuccessResponse(c, res)
}
