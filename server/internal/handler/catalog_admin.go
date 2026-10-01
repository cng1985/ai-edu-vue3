package handler

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// CatalogAdminHandler 管理端维护职业、岗位、能力、技能与知识图谱。
type CatalogAdminHandler struct {
	svc    *service.CatalogService
	talent *service.TalentService
}

func NewCatalogAdminHandler(svc *service.CatalogService, talent *service.TalentService) *CatalogAdminHandler {
	return &CatalogAdminHandler{svc: svc, talent: talent}
}

func (h *CatalogAdminHandler) Dashboard(c *gin.Context) {
	response.OK(c, h.talent.Dashboard())
}

func (h *CatalogAdminHandler) SaveCareer(c *gin.Context) {
	var in model.Career
	if !bindJSON(c, &in) {
		return
	}
	v, err := h.svc.SaveCareer(c.Param("id"), in)
	reply(c, v, err, "职业已保存")
}

func (h *CatalogAdminHandler) DeleteCareer(c *gin.Context) {
	reply(c, nil, h.svc.DeleteCareer(c.Param("id")), "职业已删除")
}

func (h *CatalogAdminHandler) Role(c *gin.Context) {
	career, role, err := h.svc.RoleFull(c.Param("id"))
	reply(c, gin.H{"career": career, "role": role}, err)
}

func (h *CatalogAdminHandler) SaveRole(c *gin.Context) {
	var in model.JobRole
	if !bindJSON(c, &in) {
		return
	}
	v, err := h.svc.SaveRole(c.Param("id"), in)
	reply(c, v, err, "岗位已保存")
}

func (h *CatalogAdminHandler) DeleteRole(c *gin.Context) {
	reply(c, nil, h.svc.DeleteRole(c.Param("id")), "岗位已删除")
}

func (h *CatalogAdminHandler) SaveCapability(c *gin.Context) {
	var in model.RoleCapability
	if !bindJSON(c, &in) {
		return
	}
	v, err := h.svc.SaveCapability(c.Param("id"), in)
	reply(c, v, err, "能力已保存")
}

func (h *CatalogAdminHandler) DeleteCapability(c *gin.Context) {
	reply(c, nil, h.svc.DeleteCapability(c.Param("id")), "能力已删除")
}

func (h *CatalogAdminHandler) SaveSkill(c *gin.Context) {
	var in model.Skill
	if !bindJSON(c, &in) {
		return
	}
	v, err := h.svc.SaveSkill(c.Param("id"), in)
	reply(c, v, err, "技能已保存")
}

func (h *CatalogAdminHandler) DeleteSkill(c *gin.Context) {
	reply(c, nil, h.svc.DeleteSkill(c.Param("id")), "技能已删除")
}

func (h *CatalogAdminHandler) Knowledge(c *gin.Context) {
	list, err := h.svc.AdminListKnowledge()
	reply(c, list, err)
}

func (h *CatalogAdminHandler) SaveKnowledge(c *gin.Context) {
	var in service.KnowledgeUpsert
	if !bindJSON(c, &in) {
		return
	}
	v, err := h.svc.SaveKnowledge(c.Param("id"), in)
	reply(c, v, err, "知识点已保存")
}

func (h *CatalogAdminHandler) DeleteKnowledge(c *gin.Context) {
	reply(c, nil, h.svc.DeleteKnowledge(c.Param("id")), "知识点已删除")
}

func (h *CatalogAdminHandler) Relations(c *gin.Context) {
	list, err := h.svc.ListRelations()
	reply(c, list, err)
}

func (h *CatalogAdminHandler) AddRelation(c *gin.Context) {
	var in model.KnowledgeRelation
	if !bindJSON(c, &in) {
		return
	}
	v, err := h.svc.AddRelation(in)
	reply(c, v, err, "关系已添加")
}

func (h *CatalogAdminHandler) DeleteRelation(c *gin.Context) {
	reply(c, nil, h.svc.DeleteRelation(c.Param("id")), "关系已删除")
}

func (h *CatalogAdminHandler) ChapterMappings(c *gin.Context) {
	list, err := h.svc.ListChapterMappings(c.Query("courseId"))
	reply(c, list, err)
}

func (h *CatalogAdminHandler) SetChapterKnowledge(c *gin.Context) {
	var in struct {
		KnowledgeIDs []string `json:"knowledgeIds"`
	}
	if !bindJSON(c, &in) {
		return
	}
	reply(c, nil, h.svc.SetChapterKnowledge(c.Param("courseId"), c.Param("chapterId"), in.KnowledgeIDs), "章节知识映射已更新")
}
