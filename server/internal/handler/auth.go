package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入用户名和密码")
		return
	}
	res, err := h.svc.Login(req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Register(req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "注册成功")
}

func (h *AuthHandler) GuestLogin(c *gin.Context) {
	res, err := h.svc.GuestLogin()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	user, err := h.svc.Me(claims.ID)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user)
}

func (h *AuthHandler) Permissions(c *gin.Context) {
	claims := middleware.GetClaims(c)
	response.OK(c, gin.H{
		"role":        claims.Role,
		"roleName":    h.svc.RoleName(claims.Role),
		"permissions": h.svc.Permissions(claims.Role),
	})
}

func (h *AuthHandler) RefreshPermissions(c *gin.Context) {
	claims := middleware.GetClaims(c)
	user, err := h.svc.RefreshPermissions(claims.ID)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user, "权限已刷新")
}
