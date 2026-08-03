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
	view, err := h.svc.GetSystemView()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, view)
}

func (h *SettingsHandler) Resolve(c *gin.Context) {
	code := c.Query("virtualModel")
	resolved, err := h.svc.ResolveVirtualModel(code)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, resolved)
}

func (h *SettingsHandler) SetDefaultVirtualModel(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请指定虚拟模型编码")
		return
	}
	if err := h.svc.SetDefaultVirtualModel(req.Code); err != nil {
		failErr(c, err)
		return
	}
	view, err := h.svc.GetSystemView()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, view, "默认虚拟模型已更新")
}

func (h *SettingsHandler) SaveProvider(c *gin.Context) {
	var p model.Provider
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.SaveProvider(p)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "保存成功")
}

func (h *SettingsHandler) UpdateProvider(c *gin.Context) {
	var p model.Provider
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	p.ID = c.Param("id")
	res, err := h.svc.SaveProvider(p)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "保存成功")
}

func (h *SettingsHandler) QuickSetup(c *gin.Context) {
	var req model.AiModelQuickSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	view, err := h.svc.QuickSetup(req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, view, "大模型配置已初始化")
}

func (h *SettingsHandler) ReindexKnowledge(c *gin.Context) {
	status, err := h.svc.ReindexKnowledge(c.Request.Context())
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, status, "知识库索引重建完成")
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
