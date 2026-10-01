package service

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"gorm.io/datatypes"
)

// Catalog 是职业-技能-知识体系的内存快照，供计算引擎一次性使用。
type Catalog struct {
	Careers         map[string]model.Career
	Roles           map[string]model.JobRole
	Skills          map[string]model.Skill
	SkillList       []model.Skill
	Knowledge       map[string]model.KnowledgePoint
	KnowledgeList   []model.KnowledgePoint
	SkillKnowledge  map[string][]model.SkillKnowledge // skillID → 知识点
	KnowledgeSkills map[string][]string               // knowledgeID → 技能
	Prereqs         map[string][]string               // knowledgeID → 前置知识
	Relations       []model.KnowledgeRelation
}

// Brief 返回知识点摘要，未知 ID 时仅填充 ID。
func (c *Catalog) Brief(id string) model.KnowledgeBrief {
	k, ok := c.Knowledge[id]
	if !ok {
		return model.KnowledgeBrief{ID: id, Name: id}
	}
	return knowledgeBrief(k)
}

func (c *Catalog) SkillName(id string) string {
	if s, ok := c.Skills[id]; ok {
		return s.Name
	}
	return id
}

func knowledgeBrief(k model.KnowledgePoint) model.KnowledgeBrief {
	return model.KnowledgeBrief{ID: k.ID, Name: k.Name, Domain: k.Domain, Difficulty: k.Difficulty, EstimatedMinutes: k.EstimatedMinutes}
}

// CatalogService 职业体系、技能体系与知识图谱。
type CatalogService struct {
	repo    *repository.CatalogRepo
	courses *repository.CourseRepo
}

func NewCatalogService(repo *repository.CatalogRepo, courses *repository.CourseRepo) *CatalogService {
	return &CatalogService{repo: repo, courses: courses}
}

// Load 读取完整目录快照。
func (s *CatalogService) Load() (*Catalog, error) {
	careers, err := s.repo.ListCareers()
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.ListRoles()
	if err != nil {
		return nil, err
	}
	skills, err := s.repo.ListSkills()
	if err != nil {
		return nil, err
	}
	knowledge, err := s.repo.ListKnowledge(false)
	if err != nil {
		return nil, err
	}
	links, err := s.repo.ListSkillKnowledge()
	if err != nil {
		return nil, err
	}
	relations, err := s.repo.ListRelations()
	if err != nil {
		return nil, err
	}
	c := &Catalog{
		Careers: map[string]model.Career{}, Roles: map[string]model.JobRole{},
		Skills: map[string]model.Skill{}, SkillList: skills,
		Knowledge: map[string]model.KnowledgePoint{}, KnowledgeList: knowledge,
		SkillKnowledge: map[string][]model.SkillKnowledge{}, KnowledgeSkills: map[string][]string{},
		Prereqs: map[string][]string{}, Relations: relations,
	}
	for _, it := range careers {
		it.Roles = nil
		c.Careers[it.ID] = it
	}
	for _, it := range roles {
		c.Roles[it.ID] = it
	}
	for _, it := range skills {
		c.Skills[it.ID] = it
	}
	for _, it := range knowledge {
		c.Knowledge[it.ID] = it
	}
	for _, l := range links {
		c.SkillKnowledge[l.SkillID] = append(c.SkillKnowledge[l.SkillID], l)
		c.KnowledgeSkills[l.KnowledgeID] = append(c.KnowledgeSkills[l.KnowledgeID], l.SkillID)
	}
	for _, r := range relations {
		if r.Type == model.RelationPrerequisite || r.Type == model.RelationAdvanced {
			c.Prereqs[r.ToID] = append(c.Prereqs[r.ToID], r.FromID)
		}
	}
	return c, nil
}

func (s *CatalogService) ListCareers() ([]model.Career, error) { return s.repo.ListCareers() }

func (s *CatalogService) ListSkills() ([]model.Skill, error) { return s.repo.ListSkills() }

// RoleFull 返回岗位及其能力与技能要求，并填充技能详情。
func (s *CatalogService) RoleFull(roleID string) (*model.Career, *model.JobRole, error) {
	role, err := s.repo.FindRoleFull(roleID)
	if err != nil {
		return nil, nil, apperr.NotFound("岗位不存在")
	}
	career, _ := s.repo.FindCareer(role.CareerID)
	skills, _ := s.repo.ListSkills()
	byID := map[string]model.Skill{}
	for _, sk := range skills {
		byID[sk.ID] = sk
	}
	for i := range role.Capabilities {
		for j := range role.Capabilities[i].Skills {
			if sk, ok := byID[role.Capabilities[i].Skills[j].SkillID]; ok {
				skCopy := sk
				role.Capabilities[i].Skills[j].Skill = &skCopy
			}
		}
	}
	return career, role, nil
}

// KnowledgeFull 读取包含正文与题目的知识点。
func (s *CatalogService) KnowledgeFull(id string) (*model.KnowledgePoint, error) {
	k, err := s.repo.FindKnowledge(id)
	if err != nil {
		return nil, apperr.NotFound("知识点不存在")
	}
	return k, nil
}

// ChapterRefs 返回引用指定知识点的课程章节。
func (s *CatalogService) ChapterRefs(knowledgeID string) []model.ChapterRef {
	maps, _ := s.repo.ListChapterKnowledge("")
	var refs []model.ChapterRef
	courseCache := map[string]*model.Course{}
	for _, m := range maps {
		if m.KnowledgeID != knowledgeID {
			continue
		}
		course, ok := courseCache[m.CourseID]
		if !ok {
			c, err := s.courses.FindByID(m.CourseID)
			if err != nil {
				continue
			}
			course = c
			courseCache[m.CourseID] = c
		}
		ref := model.ChapterRef{CourseID: course.ID, CourseTitle: course.Title, ChapterID: m.ChapterID}
		for _, ch := range course.Chapters {
			if ch.ID == m.ChapterID {
				ref.ChapterTitle = ch.Title
			}
		}
		refs = append(refs, ref)
	}
	return refs
}

// ChapterKnowledge 返回章节映射的知识点 ID；chapterID 为空时返回整门课程的知识点。
func (s *CatalogService) ChapterKnowledge(courseID, chapterID string) []string {
	maps, _ := s.repo.ListChapterKnowledge(courseID)
	seen := map[string]bool{}
	var ids []string
	for _, m := range maps {
		if chapterID != "" && m.ChapterID != chapterID {
			continue
		}
		if !seen[m.KnowledgeID] {
			seen[m.KnowledgeID] = true
			ids = append(ids, m.KnowledgeID)
		}
	}
	return ids
}

func (s *CatalogService) ListChapterMappings(courseID string) ([]model.ChapterKnowledge, error) {
	return s.repo.ListChapterKnowledge(courseID)
}

// ---------- 管理端维护 ----------

func (s *CatalogService) SaveCareer(id string, in model.Career) (*model.Career, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, apperr.BadRequest("请填写职业名称")
	}
	if id != "" {
		old, err := s.repo.FindCareer(id)
		if err != nil {
			return nil, apperr.NotFound("职业不存在")
		}
		in.ID, in.CreatedAt = old.ID, old.CreatedAt
	} else {
		in.ID = slugOr(in.ID, "career")
		in.CreatedAt = time.Now().UnixMilli()
	}
	in.Roles = nil
	return &in, s.repo.Save(&in)
}

func (s *CatalogService) DeleteCareer(id string) error { return s.repo.DeleteCareer(id) }

func (s *CatalogService) SaveRole(id string, in model.JobRole) (*model.JobRole, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" || in.CareerID == "" {
		return nil, apperr.BadRequest("请填写岗位名称并选择所属职业")
	}
	if _, err := s.repo.FindCareer(in.CareerID); err != nil {
		return nil, apperr.BadRequest("所属职业不存在")
	}
	if id != "" {
		in.ID = id
	} else {
		in.ID = slugOr(in.ID, "role")
	}
	in.Capabilities = nil
	return &in, s.repo.Save(&in)
}

func (s *CatalogService) DeleteRole(id string) error { return s.repo.DeleteRole(id) }

// SaveCapability 保存能力及其技能要求。
func (s *CatalogService) SaveCapability(id string, in model.RoleCapability) (*model.RoleCapability, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" || in.RoleID == "" {
		return nil, apperr.BadRequest("请填写能力名称并选择所属岗位")
	}
	if id != "" {
		in.ID = id
	} else {
		in.ID = slugOr(in.ID, "cap")
	}
	if in.Weight <= 0 {
		in.Weight = 1
	}
	skills := in.Skills
	in.Skills = nil
	if err := s.repo.Save(&in); err != nil {
		return nil, err
	}
	items := make([]model.CapabilitySkill, 0, len(skills))
	for _, sk := range skills {
		if sk.SkillID == "" {
			continue
		}
		items = append(items, model.CapabilitySkill{
			ID: genID("cs"), CapabilityID: in.ID, SkillID: sk.SkillID,
			RequiredLevel: clampLevel(sk.RequiredLevel), Weight: maxInt(sk.Weight, 1),
		})
	}
	if err := s.repo.ReplaceCapabilitySkills(in.ID, items); err != nil {
		return nil, err
	}
	in.Skills = items
	return &in, nil
}

func (s *CatalogService) DeleteCapability(id string) error { return s.repo.DeleteCapability(id) }

func (s *CatalogService) SaveSkill(id string, in model.Skill) (*model.Skill, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, apperr.BadRequest("请填写技能名称")
	}
	if id != "" {
		old, err := s.repo.FindSkill(id)
		if err != nil {
			return nil, apperr.NotFound("技能不存在")
		}
		in.ID, in.CreatedAt = old.ID, old.CreatedAt
	} else {
		in.ID = slugOr(in.ID, "skill")
		in.CreatedAt = time.Now().UnixMilli()
	}
	if len(in.Tags) == 0 {
		in.Tags = datatypes.JSON("[]")
	}
	return &in, s.repo.Save(&in)
}

func (s *CatalogService) DeleteSkill(id string) error { return s.repo.DeleteSkill(id) }

// KnowledgeUpsert 管理端编辑知识点。
type KnowledgeUpsert struct {
	model.KnowledgePoint
	Questions []model.Question `json:"questions"`
	SkillIDs  []string         `json:"skillIds"`
}

func (s *CatalogService) SaveKnowledge(id string, in KnowledgeUpsert) (*model.KnowledgePoint, error) {
	k := in.KnowledgePoint
	k.Name = strings.TrimSpace(k.Name)
	if k.Name == "" {
		return nil, apperr.BadRequest("请填写知识点名称")
	}
	now := time.Now().UnixMilli()
	if id != "" {
		old, err := s.repo.FindKnowledge(id)
		if err != nil {
			return nil, apperr.NotFound("知识点不存在")
		}
		k.ID, k.CreatedAt = old.ID, old.CreatedAt
	} else {
		k.ID = slugOr(k.ID, "kp")
		k.CreatedAt = now
	}
	k.UpdatedAt = now
	if k.Difficulty <= 0 {
		k.Difficulty = 1
	}
	if k.EstimatedMinutes <= 0 {
		k.EstimatedMinutes = 20
	}
	qs, _ := json.Marshal(nonNil(in.Questions))
	k.Questions = datatypes.JSON(qs)
	if len(k.Tags) == 0 {
		k.Tags = datatypes.JSON("[]")
	}
	if err := s.repo.Save(&k); err != nil {
		return nil, err
	}
	if in.SkillIDs != nil {
		links := make([]model.SkillKnowledge, 0, len(in.SkillIDs))
		for _, sid := range in.SkillIDs {
			links = append(links, model.SkillKnowledge{ID: genID("sk"), SkillID: sid, KnowledgeID: k.ID, Weight: 1})
		}
		if err := s.repo.ReplaceSkillKnowledge(k.ID, links); err != nil {
			return nil, err
		}
	}
	return &k, nil
}

func (s *CatalogService) DeleteKnowledge(id string) error { return s.repo.DeleteKnowledge(id) }

var validRelations = map[string]bool{
	model.RelationPrerequisite: true, model.RelationRelated: true,
	model.RelationSimilar: true, model.RelationAdvanced: true,
}

func (s *CatalogService) AddRelation(in model.KnowledgeRelation) (*model.KnowledgeRelation, error) {
	if in.FromID == "" || in.ToID == "" || in.FromID == in.ToID {
		return nil, apperr.BadRequest("请选择两个不同的知识点")
	}
	if !validRelations[in.Type] {
		return nil, apperr.BadRequest("关系类型无效")
	}
	if _, err := s.repo.FindKnowledge(in.FromID); err != nil {
		return nil, apperr.BadRequest("起点知识点不存在")
	}
	if _, err := s.repo.FindKnowledge(in.ToID); err != nil {
		return nil, apperr.BadRequest("终点知识点不存在")
	}
	if in.Type == model.RelationPrerequisite || in.Type == model.RelationAdvanced {
		c, err := s.Load()
		if err == nil && reaches(c.Prereqs, in.FromID, in.ToID) {
			return nil, apperr.BadRequest("该前置关系会形成环路")
		}
	}
	in.ID = genID("rel")
	return &in, s.repo.Save(&in)
}

// reaches 判断 from 是否（间接）以 to 为前置，用于阻止环路。
func reaches(prereqs map[string][]string, from, to string) bool {
	seen := map[string]bool{}
	stack := []string{from}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur == to {
			return true
		}
		if seen[cur] {
			continue
		}
		seen[cur] = true
		stack = append(stack, prereqs[cur]...)
	}
	return false
}

func (s *CatalogService) DeleteRelation(id string) error {
	return s.repo.Delete(&model.KnowledgeRelation{}, id)
}

func (s *CatalogService) SetChapterKnowledge(courseID, chapterID string, knowledgeIDs []string) error {
	if !s.courses.ChapterExists(courseID, chapterID) {
		return apperr.NotFound("章节不存在")
	}
	items := make([]model.ChapterKnowledge, 0, len(knowledgeIDs))
	for _, kid := range knowledgeIDs {
		items = append(items, model.ChapterKnowledge{ID: genID("ck"), CourseID: courseID, ChapterID: chapterID, KnowledgeID: kid})
	}
	return s.repo.ReplaceChapterKnowledge(courseID, chapterID, items)
}

// AdminKnowledge 管理端知识点列表（含关联技能）。
type AdminKnowledge struct {
	model.KnowledgePoint
	Questions     []model.Question `json:"questions"`
	SkillIDs      []string         `json:"skillIds"`
	QuestionCount int              `json:"questionCount"`
}

func (s *CatalogService) AdminListKnowledge() ([]AdminKnowledge, error) {
	list, err := s.repo.ListKnowledge(true)
	if err != nil {
		return nil, err
	}
	links, _ := s.repo.ListSkillKnowledge()
	bySkill := map[string][]string{}
	for _, l := range links {
		bySkill[l.KnowledgeID] = append(bySkill[l.KnowledgeID], l.SkillID)
	}
	out := make([]AdminKnowledge, 0, len(list))
	for _, k := range list {
		var qs []model.Question
		_ = json.Unmarshal(k.Questions, &qs)
		k.Questions = nil
		out = append(out, AdminKnowledge{KnowledgePoint: k, Questions: nonNil(qs), SkillIDs: nonNil(bySkill[k.ID]), QuestionCount: len(qs)})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Domain < out[j].Domain })
	return out, nil
}

func (s *CatalogService) ListRelations() ([]model.KnowledgeRelation, error) { return s.repo.ListRelations() }

func slugOr(id, prefix string) string {
	id = strings.TrimSpace(id)
	if id != "" {
		return id
	}
	return genID(prefix)
}

func clampLevel(l int) int {
	if l < 0 {
		return 0
	}
	if l > 5 {
		return 5
	}
	return l
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func nonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
