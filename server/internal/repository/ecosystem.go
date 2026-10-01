package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

// CatalogRepo 职业、岗位、能力、技能与知识图谱的数据访问。
type CatalogRepo struct{ db *gorm.DB }

func NewCatalogRepo(db *gorm.DB) *CatalogRepo { return &CatalogRepo{db: db} }

func (r *CatalogRepo) DB() *gorm.DB { return r.db }

func (r *CatalogRepo) CountCareers() int64 {
	var n int64
	r.db.Model(&model.Career{}).Count(&n)
	return n
}

func (r *CatalogRepo) ListCareers() ([]model.Career, error) {
	var list []model.Career
	err := r.db.Preload("Roles", func(db *gorm.DB) *gorm.DB { return db.Order("sort ASC") }).
		Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *CatalogRepo) FindCareer(id string) (*model.Career, error) {
	var c model.Career
	err := r.db.First(&c, "id = ?", id).Error
	return &c, err
}

// FindRoleFull 加载岗位及其能力、能力所需技能。
func (r *CatalogRepo) FindRoleFull(id string) (*model.JobRole, error) {
	var role model.JobRole
	err := r.db.Preload("Capabilities", func(db *gorm.DB) *gorm.DB { return db.Order("sort ASC") }).
		Preload("Capabilities.Skills").
		First(&role, "id = ?", id).Error
	return &role, err
}

func (r *CatalogRepo) ListRoles() ([]model.JobRole, error) {
	var list []model.JobRole
	err := r.db.Order("career_id ASC, sort ASC").Find(&list).Error
	return list, err
}

func (r *CatalogRepo) ListSkills() ([]model.Skill, error) {
	var list []model.Skill
	err := r.db.Order("category ASC, name ASC").Find(&list).Error
	return list, err
}

func (r *CatalogRepo) FindSkill(id string) (*model.Skill, error) {
	var s model.Skill
	err := r.db.First(&s, "id = ?", id).Error
	return &s, err
}

func (r *CatalogRepo) ListSkillKnowledge() ([]model.SkillKnowledge, error) {
	var list []model.SkillKnowledge
	err := r.db.Find(&list).Error
	return list, err
}

func (r *CatalogRepo) ListKnowledge(withContent bool) ([]model.KnowledgePoint, error) {
	var list []model.KnowledgePoint
	q := r.db.Order("domain ASC, difficulty ASC, name ASC")
	if !withContent {
		q = q.Omit("content", "questions")
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *CatalogRepo) FindKnowledge(id string) (*model.KnowledgePoint, error) {
	var k model.KnowledgePoint
	err := r.db.First(&k, "id = ?", id).Error
	return &k, err
}

func (r *CatalogRepo) ListRelations() ([]model.KnowledgeRelation, error) {
	var list []model.KnowledgeRelation
	err := r.db.Find(&list).Error
	return list, err
}

func (r *CatalogRepo) ListCapabilitySkills() ([]model.CapabilitySkill, error) {
	var list []model.CapabilitySkill
	err := r.db.Find(&list).Error
	return list, err
}

func (r *CatalogRepo) ListChapterKnowledge(courseID string) ([]model.ChapterKnowledge, error) {
	var list []model.ChapterKnowledge
	q := r.db.Model(&model.ChapterKnowledge{})
	if courseID != "" {
		q = q.Where("course_id = ?", courseID)
	}
	err := q.Find(&list).Error
	return list, err
}

// Save 通用保存（新增或更新）。
func (r *CatalogRepo) Save(v interface{}) error { return r.db.Save(v).Error }

// Delete 按主键删除。
func (r *CatalogRepo) Delete(v interface{}, id string) error {
	return r.db.Delete(v, "id = ?", id).Error
}

// DeleteCareer 级联删除职业下的岗位与能力。
func (r *CatalogRepo) DeleteCareer(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var roleIDs []string
		tx.Model(&model.JobRole{}).Where("career_id = ?", id).Pluck("id", &roleIDs)
		for _, rid := range roleIDs {
			if err := deleteRoleTx(tx, rid); err != nil {
				return err
			}
		}
		return tx.Delete(&model.Career{}, "id = ?", id).Error
	})
}

// DeleteRole 级联删除岗位下的能力与技能要求。
func (r *CatalogRepo) DeleteRole(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error { return deleteRoleTx(tx, id) })
}

func deleteRoleTx(tx *gorm.DB, roleID string) error {
	var capIDs []string
	tx.Model(&model.RoleCapability{}).Where("role_id = ?", roleID).Pluck("id", &capIDs)
	if len(capIDs) > 0 {
		if err := tx.Delete(&model.CapabilitySkill{}, "capability_id IN ?", capIDs).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.RoleCapability{}, "id IN ?", capIDs).Error; err != nil {
			return err
		}
	}
	return tx.Delete(&model.JobRole{}, "id = ?", roleID).Error
}

// DeleteCapability 删除能力及其技能要求。
func (r *CatalogRepo) DeleteCapability(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CapabilitySkill{}, "capability_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.RoleCapability{}, "id = ?", id).Error
	})
}

// DeleteSkill 删除技能及其关联关系。
func (r *CatalogRepo) DeleteSkill(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, m := range []interface{}{&model.CapabilitySkill{}, &model.SkillKnowledge{}, &model.StudentSkillState{}} {
			if err := tx.Delete(m, "skill_id = ?", id).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&model.Skill{}, "id = ?", id).Error
	})
}

// DeleteKnowledge 删除知识点及其图谱边、技能关联、课程映射。
func (r *CatalogRepo) DeleteKnowledge(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.KnowledgeRelation{}, "from_id = ? OR to_id = ?", id, id).Error; err != nil {
			return err
		}
		for _, m := range []interface{}{&model.SkillKnowledge{}, &model.ChapterKnowledge{}} {
			if err := tx.Delete(m, "knowledge_id = ?", id).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&model.KnowledgePoint{}, "id = ?", id).Error
	})
}

// ReplaceSkillKnowledge 覆盖知识点关联的技能列表。
func (r *CatalogRepo) ReplaceSkillKnowledge(knowledgeID string, links []model.SkillKnowledge) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.SkillKnowledge{}, "knowledge_id = ?", knowledgeID).Error; err != nil {
			return err
		}
		if len(links) == 0 {
			return nil
		}
		return tx.Create(&links).Error
	})
}

// ReplaceCapabilitySkills 覆盖能力的技能要求。
func (r *CatalogRepo) ReplaceCapabilitySkills(capabilityID string, items []model.CapabilitySkill) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CapabilitySkill{}, "capability_id = ?", capabilityID).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

// ReplaceChapterKnowledge 覆盖章节映射的知识点。
func (r *CatalogRepo) ReplaceChapterKnowledge(courseID, chapterID string, items []model.ChapterKnowledge) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.ChapterKnowledge{}, "course_id = ? AND chapter_id = ?", courseID, chapterID).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *CatalogRepo) Count(m interface{}) int64 {
	var n int64
	r.db.Model(m).Count(&n)
	return n
}
