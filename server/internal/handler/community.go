package handler

import (
	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/gin-gonic/gin"
)

// CommunityHandler 知识社区与知识资产接口。
type CommunityHandler struct{ svc *service.CommunityService }

func NewCommunityHandler(svc *service.CommunityService) *CommunityHandler {
	return &CommunityHandler{svc: svc}
}

func (h *CommunityHandler) Posts(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListPosts(c.Query("type"), c.Query("keyword"), c.Query("authorId"), page, pageSize)
	reply(c, res, err)
}

func (h *CommunityHandler) Post(c *gin.Context) {
	d, err := h.svc.GetPost(c.Param("id"))
	reply(c, d, err)
}

func (h *CommunityHandler) CreatePost(c *gin.Context) {
	var req model.PostRequest
	if !bindJSON(c, &req) {
		return
	}
	p, err := h.svc.CreatePost(middleware.GetClaims(c).ID, req)
	reply(c, p, err, "发布成功")
}

func (h *CommunityHandler) Answer(c *gin.Context) {
	var req model.AnswerRequest
	if !bindJSON(c, &req) {
		return
	}
	a, err := h.svc.Answer(middleware.GetClaims(c).ID, c.Param("id"), req)
	reply(c, a, err, "回答已发布")
}

func (h *CommunityHandler) LikePost(c *gin.Context) {
	reply(c, nil, h.svc.LikePost(c.Param("id")))
}

func (h *CommunityHandler) LikeAnswer(c *gin.Context) {
	a, err := h.svc.LikeAnswer(c.Param("id"))
	reply(c, a, err)
}

func (h *CommunityHandler) AcceptAnswer(c *gin.Context) {
	reply(c, nil, h.svc.AcceptAnswer(middleware.GetClaims(c).ID, c.Param("id")), "已采纳")
}

func (h *CommunityHandler) Summarize(c *gin.Context) {
	d, err := h.svc.Summarize(c.Request.Context(), middleware.GetClaims(c).ID, c.Param("id"))
	reply(c, d, err, "已沉淀为知识资产")
}

func (h *CommunityHandler) DeletePost(c *gin.Context) {
	reply(c, nil, h.svc.DeletePost(c.Param("id")), "已删除")
}

func (h *CommunityHandler) Resources(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.ListResources(c.Query("type"), c.Query("keyword"), page, pageSize)
	reply(c, res, err)
}

func (h *CommunityHandler) Resource(c *gin.Context) {
	r, err := h.svc.GetResource(c.Param("id"))
	reply(c, r, err)
}

func (h *CommunityHandler) CreateResource(c *gin.Context) {
	var req model.ResourceRequest
	if !bindJSON(c, &req) {
		return
	}
	r, err := h.svc.CreateResource(middleware.GetClaims(c).ID, req)
	reply(c, r, err, "知识资产已发布")
}

func (h *CommunityHandler) DeleteResource(c *gin.Context) {
	reply(c, nil, h.svc.DeleteResource(c.Param("id")), "已删除")
}
