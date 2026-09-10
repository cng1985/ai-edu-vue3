package handler

import (
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct{ svc *service.UserService }

func NewUserHandler(svc *service.UserService) *UserHandler { return &UserHandler{svc: svc} }

func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("role"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *UserHandler) Get(c *gin.Context) {
	user, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req model.UserUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	user, err := h.svc.Create(req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user, "创建成功")
}

func (h *UserHandler) Update(c *gin.Context) {
	var req model.UserUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	user, err := h.svc.Update(c.Param("id"), req)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user, "更新成功")
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}
