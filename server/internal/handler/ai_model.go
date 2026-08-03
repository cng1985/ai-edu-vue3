package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AiModelHandler struct{ svc *service.AiModelService }

func NewAiModelHandler(svc *service.AiModelService) *AiModelHandler {
	return &AiModelHandler{svc: svc}
}

func (h *AiModelHandler) Overview(c *gin.Context) {
	overview, err := h.svc.Overview()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, overview)
}

func (h *AiModelHandler) SetDefault(c *gin.Context) {
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
	response.OK(c, nil, "默认虚拟模型已更新")
}

func (h *AiModelHandler) ResolveTest(c *gin.Context) {
	code := c.Query("virtualModel")
	resolved, err := h.svc.ResolveTest(code)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, resolved)
}

// --- CanonicalModel ---

func (h *AiModelHandler) ListCanonicalModels(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListCanonicalModels(c.Query("keyword"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AiModelHandler) CreateCanonicalModel(c *gin.Context) {
	var m model.CanonicalModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateCanonicalModel(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) UpdateCanonicalModel(c *gin.Context) {
	var m model.CanonicalModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateCanonicalModel(c.Param("id"), m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *AiModelHandler) DeleteCanonicalModel(c *gin.Context) {
	if err := h.svc.DeleteCanonicalModel(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- Capability ---

func (h *AiModelHandler) ListCapabilities(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListCapabilities(c.Query("keyword"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AiModelHandler) CreateCapability(c *gin.Context) {
	var m model.Capability
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateCapability(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) UpdateCapability(c *gin.Context) {
	var m model.Capability
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateCapability(c.Param("id"), m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *AiModelHandler) DeleteCapability(c *gin.Context) {
	if err := h.svc.DeleteCapability(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- CapabilityModel ---

func (h *AiModelHandler) ListCapabilityModels(c *gin.Context) {
	list, err := h.svc.ListCapabilityModels(c.Query("canonicalModelId"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, list)
}

func (h *AiModelHandler) CreateCapabilityModel(c *gin.Context) {
	var m model.CapabilityModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateCapabilityModel(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) DeleteCapabilityModel(c *gin.Context) {
	if err := h.svc.DeleteCapabilityModel(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- Provider ---

func (h *AiModelHandler) ListProviders(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListProviders(c.Query("keyword"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AiModelHandler) GetProvider(c *gin.Context) {
	p, err := h.svc.GetProvider(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *AiModelHandler) CreateProvider(c *gin.Context) {
	var m model.Provider
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateProvider(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) UpdateProvider(c *gin.Context) {
	var m model.Provider
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateProvider(c.Param("id"), m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *AiModelHandler) DeleteProvider(c *gin.Context) {
	if err := h.svc.DeleteProvider(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- ProviderModel ---

func (h *AiModelHandler) ListProviderModels(c *gin.Context) {
	list, err := h.svc.ListProviderModels(c.Query("providerId"), c.Query("canonicalModelId"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, list)
}

func (h *AiModelHandler) CreateProviderModel(c *gin.Context) {
	var m model.ProviderModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateProviderModel(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) UpdateProviderModel(c *gin.Context) {
	var m model.ProviderModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateProviderModel(c.Param("id"), m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *AiModelHandler) DeleteProviderModel(c *gin.Context) {
	if err := h.svc.DeleteProviderModel(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- VirtualModel ---

func (h *AiModelHandler) ListVirtualModels(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListVirtualModels(c.Query("keyword"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AiModelHandler) CreateVirtualModel(c *gin.Context) {
	var m model.VirtualModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateVirtualModel(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) UpdateVirtualModel(c *gin.Context) {
	var m model.VirtualModel
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateVirtualModel(c.Param("id"), m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *AiModelHandler) DeleteVirtualModel(c *gin.Context) {
	if err := h.svc.DeleteVirtualModel(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- VirtualModelMapping ---

func (h *AiModelHandler) ListVirtualModelMappings(c *gin.Context) {
	list, err := h.svc.ListVirtualModelMappings(c.Query("virtualModelId"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, list)
}

func (h *AiModelHandler) CreateVirtualModelMapping(c *gin.Context) {
	var m model.VirtualModelMapping
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.CreateVirtualModelMapping(m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *AiModelHandler) UpdateVirtualModelMapping(c *gin.Context) {
	var m model.VirtualModelMapping
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateVirtualModelMapping(c.Param("id"), m)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *AiModelHandler) ListVirtualModelOptions(c *gin.Context) {
	list, err := h.svc.ListVirtualModelOptions()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, list)
}

func (h *AiModelHandler) DeleteVirtualModelMapping(c *gin.Context) {
	if err := h.svc.DeleteVirtualModelMapping(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}
