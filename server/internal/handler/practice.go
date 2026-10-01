package handler

import (
	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// PracticeHandler 项目实践与 IT 任务市场接口。
type PracticeHandler struct {
	projects *service.ProjectService
	market   *service.MarketService
	rbac     *service.RBACService
}

func NewPracticeHandler(projects *service.ProjectService, market *service.MarketService, rbacSvc *service.RBACService) *PracticeHandler {
	return &PracticeHandler{projects: projects, market: market, rbac: rbacSvc}
}

func (h *PracticeHandler) canManage(c *gin.Context) bool {
	return h.rbac.Resolver().Has(middleware.GetClaims(c).Role, rbac.PermOpportunityManage)
}

func (h *PracticeHandler) Projects(c *gin.Context) {
	list, err := h.projects.List(middleware.GetClaims(c).ID, "published")
	reply(c, list, err)
}

func (h *PracticeHandler) Project(c *gin.Context) {
	d, err := h.projects.Get(middleware.GetClaims(c).ID, c.Param("id"))
	reply(c, d, err)
}

func (h *PracticeHandler) Submit(c *gin.Context) {
	var req model.SubmissionRequest
	if !bindJSON(c, &req) {
		return
	}
	v, err := h.projects.Submit(c.Request.Context(), middleware.GetClaims(c).ID, c.Param("id"), c.Param("taskId"), req)
	reply(c, v, err, "已提交，Review Agent 评审完成")
}

func (h *PracticeHandler) MySubmissions(c *gin.Context) {
	response.OK(c, h.projects.Submissions(middleware.GetClaims(c).ID, intQuery(c, "limit", 50)))
}

func (h *PracticeHandler) Opportunities(c *gin.Context) {
	claims := middleware.GetClaims(c)
	list, err := h.market.List(claims.ID, claims.Role, c.Query("status"), c.Query("sort"))
	reply(c, list, err)
}

func (h *PracticeHandler) Opportunity(c *gin.Context) {
	claims := middleware.GetClaims(c)
	d, err := h.market.Get(claims.ID, claims.Role, c.Param("id"), h.canManage(c))
	reply(c, d, err)
}

func (h *PracticeHandler) Apply(c *gin.Context) {
	var req model.ApplyRequest
	_ = c.ShouldBindJSON(&req)
	claims := middleware.GetClaims(c)
	a, err := h.market.Apply(claims.ID, claims.Role, c.Param("id"), req)
	reply(c, a, err, "申请已提交")
}

func (h *PracticeHandler) MyApplications(c *gin.Context) {
	response.OK(c, h.market.MyApplications(middleware.GetClaims(c).ID))
}

func (h *PracticeHandler) Published(c *gin.Context) {
	list, err := h.market.Published(middleware.GetClaims(c).ID)
	reply(c, list, err)
}

func (h *PracticeHandler) Publish(c *gin.Context) {
	var req model.OpportunityRequest
	if !bindJSON(c, &req) {
		return
	}
	o, err := h.market.Publish(middleware.GetClaims(c).ID, req)
	reply(c, o, err, "任务已发布")
}

func (h *PracticeHandler) UpdateOpportunity(c *gin.Context) {
	var req model.OpportunityRequest
	if !bindJSON(c, &req) {
		return
	}
	o, err := h.market.Update(middleware.GetClaims(c).ID, c.Param("id"), h.canManage(c), req)
	reply(c, o, err, "任务已更新")
}

func (h *PracticeHandler) DeleteOpportunity(c *gin.Context) {
	err := h.market.Delete(middleware.GetClaims(c).ID, c.Param("id"), h.canManage(c))
	reply(c, nil, err, "任务已删除")
}

func (h *PracticeHandler) Candidates(c *gin.Context) {
	list, err := h.market.Candidates(middleware.GetClaims(c).ID, c.Param("id"), h.canManage(c), intQuery(c, "limit", 10))
	reply(c, list, err)
}

func (h *PracticeHandler) Decide(c *gin.Context) {
	var req model.ApplicationDecision
	if !bindJSON(c, &req) {
		return
	}
	a, err := h.market.Decide(middleware.GetClaims(c).ID, c.Param("id"), h.canManage(c), req)
	reply(c, a, err, "已处理")
}

// ---------- 管理端 ----------

func (h *PracticeHandler) AdminProjects(c *gin.Context) {
	list, err := h.projects.List("", c.Query("status"))
	reply(c, list, err)
}

func (h *PracticeHandler) AdminProject(c *gin.Context) {
	d, err := h.projects.Get("", c.Param("id"))
	reply(c, d, err)
}

func (h *PracticeHandler) SaveProject(c *gin.Context) {
	var req model.ProjectRequest
	if !bindJSON(c, &req) {
		return
	}
	p, err := h.projects.Save(c.Param("id"), middleware.GetClaims(c).ID, req)
	reply(c, p, err, "项目已保存")
}

func (h *PracticeHandler) DeleteProject(c *gin.Context) {
	reply(c, nil, h.projects.Delete(c.Param("id")), "项目已删除")
}

func (h *PracticeHandler) AdminSubmissions(c *gin.Context) {
	response.OK(c, h.projects.AllSubmissions(c.Query("projectId"), intQuery(c, "limit", 100)))
}

func (h *PracticeHandler) AdminOpportunities(c *gin.Context) {
	list, err := h.market.AdminList(c.Query("status"))
	reply(c, list, err)
}
