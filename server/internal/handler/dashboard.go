package handler

import (
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct{ svc *service.DashboardService }

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, stats)
}
