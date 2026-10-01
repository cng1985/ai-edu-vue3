package repository

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"gorm.io/gorm"
)

// CommunityRepo 知识社区与知识资产。
type CommunityRepo struct{ db *gorm.DB }

func NewCommunityRepo(db *gorm.DB) *CommunityRepo { return &CommunityRepo{db: db} }

func (r *CommunityRepo) ListPosts(postType, keyword, authorID string, page, pageSize int) ([]model.CommunityPost, int64, error) {
	q := r.db.Model(&model.CommunityPost{}).Where("status = ?", "published")
	if postType != "" {
		q = q.Where("type = ?", postType)
	}
	if authorID != "" {
		q = q.Where("author_id = ?", authorID)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR content LIKE ?", kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CommunityPost
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *CommunityRepo) FindPost(id string) (*model.CommunityPost, error) {
	var p model.CommunityPost
	err := r.db.First(&p, "id = ?", id).Error
	return &p, err
}

func (r *CommunityRepo) SavePost(p *model.CommunityPost) error { return r.db.Save(p).Error }

func (r *CommunityRepo) IncrPost(id, column string, delta int) error {
	return r.db.Model(&model.CommunityPost{}).Where("id = ?", id).
		UpdateColumn(column, gorm.Expr(column+" + ?", delta)).Error
}

func (r *CommunityRepo) DeletePost(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CommunityAnswer{}, "post_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.CommunityPost{}, "id = ?", id).Error
	})
}

func (r *CommunityRepo) ListAnswers(postID string) ([]model.CommunityAnswer, error) {
	var list []model.CommunityAnswer
	err := r.db.Where("post_id = ?", postID).Order("accepted DESC, likes DESC, created_at ASC").Find(&list).Error
	return list, err
}

func (r *CommunityRepo) FindAnswer(id string) (*model.CommunityAnswer, error) {
	var a model.CommunityAnswer
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *CommunityRepo) SaveAnswer(a *model.CommunityAnswer) error { return r.db.Save(a).Error }

func (r *CommunityRepo) CountPostsBy(authorID string) int64 {
	var n int64
	r.db.Model(&model.CommunityPost{}).Where("author_id = ?", authorID).Count(&n)
	return n
}

func (r *CommunityRepo) CountAnswersBy(authorID string) (count int64, likes int64) {
	r.db.Model(&model.CommunityAnswer{}).Where("author_id = ?", authorID).Count(&count)
	r.db.Model(&model.CommunityAnswer{}).Where("author_id = ?", authorID).Select("COALESCE(SUM(likes), 0)").Scan(&likes)
	return
}

func (r *CommunityRepo) CountPosts() int64 {
	var n int64
	r.db.Model(&model.CommunityPost{}).Count(&n)
	return n
}

func (r *CommunityRepo) ListResources(resType, keyword string, page, pageSize int) ([]model.KnowledgeResource, int64, error) {
	q := r.db.Model(&model.KnowledgeResource{}).Where("status = ?", "published")
	if resType != "" {
		q = q.Where("type = ?", resType)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR summary LIKE ?", kw, kw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.KnowledgeResource
	err := q.Omit("content").Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *CommunityRepo) FindResource(id string) (*model.KnowledgeResource, error) {
	var res model.KnowledgeResource
	err := r.db.First(&res, "id = ?", id).Error
	return &res, err
}

func (r *CommunityRepo) SaveResource(res *model.KnowledgeResource) error { return r.db.Save(res).Error }

func (r *CommunityRepo) DeleteResource(id string) error {
	return r.db.Delete(&model.KnowledgeResource{}, "id = ?", id).Error
}

func (r *CommunityRepo) CountResourcesBy(authorID string) int64 {
	var n int64
	r.db.Model(&model.KnowledgeResource{}).Where("author_id = ?", authorID).Count(&n)
	return n
}

func (r *CommunityRepo) CountResourcesByType() map[string]int64 {
	type row struct {
		Type  string
		Count int64
	}
	var rows []row
	r.db.Model(&model.KnowledgeResource{}).Where("status = ?", "published").Select("type, COUNT(*) as count").Group("type").Scan(&rows)
	out := map[string]int64{}
	for _, it := range rows {
		out[it.Type] = it.Count
	}
	return out
}

// KernelRepo AI Pipeline 运行记录。
type KernelRepo struct{ db *gorm.DB }

func NewKernelRepo(db *gorm.DB) *KernelRepo { return &KernelRepo{db: db} }

func (r *KernelRepo) Save(run *model.PipelineRun) error { return r.db.Save(run).Error }

func (r *KernelRepo) List(userID string, limit int) ([]model.PipelineRun, error) {
	var list []model.PipelineRun
	q := r.db.Order("created_at DESC")
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *KernelRepo) Latest(userID, pipeline string) (*model.PipelineRun, error) {
	var run model.PipelineRun
	err := r.db.Where("user_id = ? AND pipeline = ? AND status = ?", userID, pipeline, "ok").
		Order("created_at DESC").First(&run).Error
	return &run, err
}

func (r *KernelRepo) Count() int64 {
	var n int64
	r.db.Model(&model.PipelineRun{}).Count(&n)
	return n
}
