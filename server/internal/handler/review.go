package handler

import (
	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct{ svc *service.ReviewService }

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler { return &ReviewHandler{svc: svc} }

func (h *ReviewHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("status"), c.Query("type"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *ReviewHandler) Get(c *gin.Context) {
	review, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, review)
}

func (h *ReviewHandler) Approve(c *gin.Context) {
	var req model.ReviewDecisionRequest
	_ = c.ShouldBindJSON(&req)
	claims := middleware.GetClaims(c)
	review, err := h.svc.Approve(c.Param("id"), claims.ID, req.Comment)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, review, "审核通过")
}

func (h *ReviewHandler) Reject(c *gin.Context) {
	var req model.ReviewDecisionRequest
	_ = c.ShouldBindJSON(&req)
	claims := middleware.GetClaims(c)
	review, err := h.svc.Reject(c.Param("id"), claims.ID, req.Comment)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, review, "已驳回")
}
