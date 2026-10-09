package handler

import (
	"net/http"
	"strconv"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// GrowthHandler 职业体系、知识图谱与个人成长接口。
type GrowthHandler struct {
	catalog   *service.CatalogService
	growth    *service.GrowthService
	market    *service.MarketService
	community *service.CommunityService
	talent    *service.TalentService
}

func NewGrowthHandler(catalog *service.CatalogService, growth *service.GrowthService, market *service.MarketService,
	community *service.CommunityService, talent *service.TalentService) *GrowthHandler {
	return &GrowthHandler{catalog: catalog, growth: growth, market: market, community: community, talent: talent}
}

func bindJSON(c *gin.Context, dst interface{}) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return false
	}
	return true
}

func reply(c *gin.Context, data interface{}, err error, msg ...string) {
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, data, msg...)
}

func intQuery(c *gin.Context, key string, def int) int {
	v, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return def
	}
	return v
}

func (h *GrowthHandler) Careers(c *gin.Context) {
	list, err := h.catalog.ListCareers()
	reply(c, list, err)
}

func (h *GrowthHandler) Skills(c *gin.Context) {
	list, err := h.catalog.ListSkills()
	reply(c, list, err)
}

// Role 岗位详情及当前用户的差距分析。
func (h *GrowthHandler) Role(c *gin.Context) {
	career, role, err := h.catalog.RoleFull(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	skills, _ := h.catalog.ListSkills()
	byID := map[string]model.Skill{}
	for _, s := range skills {
		byID[s.ID] = s
	}
	gap, _ := h.growth.Gap(middleware.GetClaims(c).ID, role.ID)
	response.OK(c, model.RoleDetail{Career: career, Role: role, Gap: gap, Skills: byID})
}

func (h *GrowthHandler) Graph(c *gin.Context) {
	g, err := h.growth.Graph(middleware.GetClaims(c).ID, c.Query("skillId"))
	reply(c, g, err)
}

func (h *GrowthHandler) Knowledge(c *gin.Context) {
	id := c.Param("id")
	d, err := h.growth.KnowledgeDetail(middleware.GetClaims(c).ID, id, h.community.ResourcesForKnowledge(id, 6))
	reply(c, d, err)
}

func (h *GrowthHandler) Overview(c *gin.Context) {
	claims := middleware.GetClaims(c)
	contrib := h.community.Contributions(claims.ID)
	ov, err := h.growth.Overview(claims.ID, h.market.TopMatches(claims.ID, claims.Role, 3), int(contrib.Posts+contrib.Answers+contrib.Resources))
	reply(c, ov, err)
}

func (h *GrowthHandler) Goal(c *gin.Context) {
	response.OK(c, h.growth.Goal(middleware.GetClaims(c).ID))
}

func (h *GrowthHandler) SetGoal(c *gin.Context) {
	var req model.GoalRequest
	if !bindJSON(c, &req) {
		return
	}
	g, err := h.growth.SetGoal(middleware.GetClaims(c).ID, req)
	reply(c, g, err, "职业目标已更新")
}

func (h *GrowthHandler) Gap(c *gin.Context) {
	g, err := h.growth.Gap(middleware.GetClaims(c).ID, c.Query("roleId"))
	reply(c, g, err)
}

func (h *GrowthHandler) SkillStates(c *gin.Context) {
	list, err := h.growth.SkillStates(middleware.GetClaims(c).ID)
	reply(c, list, err)
}

func (h *GrowthHandler) KnowledgeStates(c *gin.Context) {
	list, err := h.growth.KnowledgeStates(middleware.GetClaims(c).ID)
	reply(c, list, err)
}

func (h *GrowthHandler) Recommendations(c *gin.Context) {
	uid := middleware.GetClaims(c).ID
	gap, _ := h.growth.Gap(uid, "")
	list, err := h.growth.Recommend(uid, gap, intQuery(c, "limit", 10))
	reply(c, list, err)
}

func (h *GrowthHandler) DueReviews(c *gin.Context) {
	response.OK(c, h.growth.DueReviews(middleware.GetClaims(c).ID, intQuery(c, "limit", 20)))
}

func (h *GrowthHandler) RecordEvent(c *gin.Context) {
	var req model.LearningEventRequest
	if !bindJSON(c, &req) {
		return
	}
	v, err := h.growth.RecordEvent(middleware.GetClaims(c).ID, req)
	reply(c, v, err)
}

func (h *GrowthHandler) Practice(c *gin.Context) {
	var req model.PracticeAnswerRequest
	if !bindJSON(c, &req) {
		return
	}
	v, err := h.growth.Practice(middleware.GetClaims(c).ID, c.Param("id"), req.Answers)
	reply(c, v, err)
}

func (h *GrowthHandler) CompleteChapter(c *gin.Context) {
	var req model.ChapterCompleteRequest
	if !bindJSON(c, &req) {
		return
	}
	v, err := h.growth.CompleteChapter(middleware.GetClaims(c).ID, req)
	reply(c, v, err)
}

func (h *GrowthHandler) Profile(c *gin.Context) {
	p, err := h.talent.Profile(middleware.GetClaims(c).ID)
	reply(c, p, err)
}

func (h *GrowthHandler) TalentProfile(c *gin.Context) {
	p, err := h.talent.Profile(c.Param("id"))
	reply(c, p, err)
}

func (h *GrowthHandler) Talents(c *gin.Context) {
	list, err := h.talent.List(c.Query("keyword"))
	reply(c, list, err)
}
