package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/utils/httpx"
)

type AdminHandler struct {
	adminService *service.AdminService
	logger       *slog.Logger
}

func NewAdminHandler(logger *slog.Logger, adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		logger:       logger,
		adminService: adminService,
	}
}

func (h *AdminHandler) HandleDashboardData(c *gin.Context) {
	data, err := h.adminService.GetDashboardData(c.Request.Context())
	if err != nil {
		httpx.ServerErrorResponse(c, h.logger, err)
		return
	}

	httpx.SuccessResponse(c, data)
}
