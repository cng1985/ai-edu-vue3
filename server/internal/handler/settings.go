package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct{ svc *service.SettingsService }

func NewSettingsHandler(svc *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

func (h *SettingsHandler) Get(c *gin.Context) {
	response.OK(c, h.svc.GetView())
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var req model.SettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	view, err := h.svc.Update(req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, view, "配置已保存")
}
