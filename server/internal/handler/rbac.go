package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type RBACHandler struct{ svc *service.RBACService }

func NewRBACHandler(svc *service.RBACService) *RBACHandler { return &RBACHandler{svc: svc} }

func (h *RBACHandler) ListPermissions(c *gin.Context) {
	response.OK(c, h.svc.ListPermissions())
}

func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.svc.ListRoles()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, roles)
}

func (h *RBACHandler) UpdateRole(c *gin.Context) {
	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	if err := h.svc.UpdateRole(c.Param("role"), req.Permissions); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "权限已更新")
}
