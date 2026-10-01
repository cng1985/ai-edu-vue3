package model

import "gorm.io/datatypes"

// 本文件定义 AI 教育知识工作生态的核心领域模型，对应完整成长链：
// Career → Role → Capability → Skill → Knowledge → Course → Project → Task → Evidence → Talent Profile

// ---------- 职业体系 ----------

// Career 职业，例如软件工程师、AI 工程师。
type Career struct {
	ID          string    `gorm:"primaryKey;size:64" json:"id"`
	Name        string    `gorm:"size:100" json:"name"`
	Category    string    `gorm:"size:50" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"size:20" json:"icon"`
	Color       string    `gorm:"size:20" json:"color"`
	Demand      string    `gorm:"size:20" json:"demand"`
	SalaryRange string    `gorm:"size:50" json:"salaryRange"`
	Sort        int       `json:"sort"`
	CreatedAt   int64     `json:"createdAt"`
	Roles       []JobRole `gorm:"foreignKey:CareerID" json:"roles,omitempty"`
}

// JobRole 岗位（与 RBAC 角色区分），例如高级 Java 工程师。
type JobRole struct {
	ID           string       `gorm:"primaryKey;size:64" json:"id"`
	CareerID     string       `gorm:"index;size:64" json:"careerId"`
	Name         string       `gorm:"size:100" json:"name"`
	Level        string       `gorm:"size:20" json:"level"`
	Description  string       `gorm:"type:text" json:"description"`
	Sort         int          `json:"sort"`
	Capabilities []RoleCapability `gorm:"foreignKey:RoleID" json:"capabilities,omitempty"`
}

// RoleCapability 能力：能完成什么类型的问题，归属于岗位。
type RoleCapability struct {
	ID          string            `gorm:"primaryKey;size:64" json:"id"`
	RoleID      string            `gorm:"index;size:64" json:"roleId"`
	Name        string            `gorm:"size:100" json:"name"`
	Description string            `gorm:"type:text" json:"description"`
	Weight      int               `json:"weight"`
	Sort        int               `json:"sort"`
	Skills      []CapabilitySkill `gorm:"foreignKey:CapabilityID" json:"skills,omitempty"`
}

// CapabilitySkill 能力所需技能及要求等级。
type CapabilitySkill struct {
	ID            string `gorm:"primaryKey;size:64" json:"id"`
	CapabilityID  string `gorm:"index;size:64" json:"capabilityId"`
	SkillID       string `gorm:"index;size:64" json:"skillId"`
	RequiredLevel int    `json:"requiredLevel"`
	Weight        int    `json:"weight"`
	Skill         *Skill `gorm:"-" json:"skill,omitempty"`
}

// ---------- 技能与知识体系 ----------

// Skill 具体可以执行的能力，等级 L0-L5。
type Skill struct {
	ID          string         `gorm:"primaryKey;size:64" json:"id"`
	Name        string         `gorm:"size:100" json:"name"`
	Category    string         `gorm:"size:50;index" json:"category"`
	Description string         `gorm:"type:text" json:"description"`
	Tags        datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt   int64          `json:"createdAt"`
}

// SkillKnowledge 技能与知识点的多对多关系。
type SkillKnowledge struct {
	ID          string `gorm:"primaryKey;size:64" json:"id"`
	SkillID     string `gorm:"index;size:64" json:"skillId"`
	KnowledgeID string `gorm:"index;size:64" json:"knowledgeId"`
	Weight      int    `json:"weight"`
}

// KnowledgePoint 知识点：人知道什么。
type KnowledgePoint struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`
	Name             string         `gorm:"size:100" json:"name"`
	Domain           string         `gorm:"size:50;index" json:"domain"`
	Summary          string         `gorm:"type:text" json:"summary"`
	Content          string         `gorm:"type:text" json:"content,omitempty"`
	Difficulty       int            `json:"difficulty"`
	EstimatedMinutes int            `json:"estimatedMinutes"`
	Questions        datatypes.JSON `gorm:"type:json" json:"questions,omitempty"`
	Tags             datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt        int64          `json:"createdAt"`
	UpdatedAt        int64          `json:"updatedAt"`
}

// 知识关系类型
const (
	RelationPrerequisite = "prerequisite" // From 是 To 的前置知识
	RelationRelated      = "related"
	RelationSimilar      = "similar"
	RelationAdvanced     = "advanced" // To 是 From 的高级知识
)

// KnowledgeRelation 知识图谱的有向边。
type KnowledgeRelation struct {
	ID     string `gorm:"primaryKey;size:64" json:"id"`
	FromID string `gorm:"index;size:64" json:"fromId"`
	ToID   string `gorm:"index;size:64" json:"toId"`
	Type   string `gorm:"size:20;index" json:"type"`
}

// ChapterKnowledge 课程章节与知识点的映射：课程结构与知识结构分离，知识点可被多个课程复用。
type ChapterKnowledge struct {
	ID          string `gorm:"primaryKey;size:64" json:"id"`
	CourseID    string `gorm:"index;size:64" json:"courseId"`
	ChapterID   string `gorm:"index;size:64" json:"chapterId"`
	KnowledgeID string `gorm:"index;size:64" json:"knowledgeId"`
}

// ---------- 学生学习模型 ----------

// 知识掌握状态
const (
	KnowledgeStatusNew      = "new"
	KnowledgeStatusLearning = "learning"
	KnowledgeStatusMastered = "mastered"
)

// StudentKnowledgeState 学生对某知识点的学习状态。学会不是 Boolean，而是 mastery + confidence。
type StudentKnowledgeState struct {
	ID            string  `gorm:"primaryKey;size:64" json:"id"`
	UserID        string  `gorm:"uniqueIndex:idx_user_knowledge;size:64" json:"userId"`
	KnowledgeID   string  `gorm:"uniqueIndex:idx_user_knowledge;size:64" json:"knowledgeId"`
	StudyCount    int     `json:"studyCount"`
	PracticeCount int     `json:"practiceCount"`
	CorrectCount  int     `json:"correctCount"`
	Accuracy      float64 `json:"accuracy"`
	Mastery       float64 `json:"mastery"`
	Confidence    float64 `json:"confidence"`
	Status        string  `gorm:"size:20" json:"status"`
	LastStudiedAt int64   `json:"lastStudiedAt"`
	NextReviewAt  int64   `json:"nextReviewAt"`
	UpdatedAt     int64   `json:"updatedAt"`
}

// 学习事件类型
const (
	EventStudy    = "study"
	EventPractice = "practice"
	EventReview   = "review"
)

// LearningEvent 学习行为流水，是个人成长数据的原始记录。
type LearningEvent struct {
	ID          string `gorm:"primaryKey;size:64" json:"id"`
	UserID      string `gorm:"index;size:64" json:"userId"`
	KnowledgeID string `gorm:"index;size:64" json:"knowledgeId"`
	Type        string `gorm:"size:20" json:"type"`
	Correct     int    `json:"correct"`
	Total       int    `json:"total"`
	Source      string `gorm:"size:100" json:"source"`
	Minutes     int    `json:"minutes"`
	CreatedAt   int64  `gorm:"index" json:"createdAt"`
}

// StudentSkillState 技能状态 = 知识状态 + 项目证据 + 任务表现。
type StudentSkillState struct {
	ID             string  `gorm:"primaryKey;size:64" json:"id"`
	UserID         string  `gorm:"uniqueIndex:idx_user_skill;size:64" json:"userId"`
	SkillID        string  `gorm:"uniqueIndex:idx_user_skill;size:64" json:"skillId"`
	Level          int     `json:"level"`
	Score          float64 `json:"score"`
	KnowledgeScore float64 `json:"knowledgeScore"`
	ProjectScore   float64 `json:"projectScore"`
	TaskScore      float64 `json:"taskScore"`
	EvidenceCount  int     `json:"evidenceCount"`
	UpdatedAt      int64   `json:"updatedAt"`
}

// UserGoal 用户的职业目标。
type UserGoal struct {
	UserID      string `gorm:"primaryKey;size:64" json:"userId"`
	CareerID    string `gorm:"size:64" json:"careerId"`
	RoleID      string `gorm:"size:64" json:"roleId"`
	WeeklyHours int    `json:"weeklyHours"`
	TargetWeeks int    `json:"targetWeeks"`
	Motivation  string `gorm:"type:text" json:"motivation"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

// ---------- 项目实践体系 ----------

// SkillRequirement 技能要求，例如 Spring Boot L3。
type SkillRequirement struct {
	SkillID string `json:"skillId"`
	Level   int    `json:"level"`
}

// Project 实践项目：将知识转化为能力。
type Project struct {
	ID          string         `gorm:"primaryKey;size:64" json:"id"`
	Title       string         `gorm:"size:200" json:"title"`
	Summary     string         `gorm:"type:text" json:"summary"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Domain      string         `gorm:"size:50" json:"domain"`
	Difficulty  int            `json:"difficulty"`
	Icon        string         `gorm:"size:20" json:"icon"`
	Color       string         `gorm:"size:20" json:"color"`
	AuthorID    string         `gorm:"size:64" json:"authorId"`
	SkillIDs    datatypes.JSON `gorm:"type:json" json:"skillIds"`
	Status      string         `gorm:"size:20;index" json:"status"`
	CreatedAt   int64          `json:"createdAt"`
	Tasks       []ProjectTask  `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// ProjectTask 项目中的任务。
type ProjectTask struct {
	ID             string         `gorm:"primaryKey;size:64" json:"id"`
	ProjectID      string         `gorm:"index;size:64" json:"projectId"`
	Title          string         `gorm:"size:200" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Deliverable    string         `gorm:"type:text" json:"deliverable"`
	Requirements   datatypes.JSON `gorm:"type:json" json:"requirements"`
	EstimatedHours int            `json:"estimatedHours"`
	Sort           int            `json:"sort"`
}

// 提交状态
const (
	SubmissionReviewed = "reviewed"
)

// TaskSubmission 任务执行与提交（Execution）。
type TaskSubmission struct {
	ID         string `gorm:"primaryKey;size:64" json:"id"`
	UserID     string `gorm:"index;size:64" json:"userId"`
	ProjectID  string `gorm:"index;size:64" json:"projectId"`
	TaskID     string `gorm:"index;size:64" json:"taskId"`
	Content    string `gorm:"type:text" json:"content"`
	RepoURL    string `gorm:"size:300" json:"repoUrl"`
	Status     string `gorm:"size:20" json:"status"`
	Score      int    `json:"score"`
	Feedback   string `gorm:"type:text" json:"feedback"`
	Evaluator  string `gorm:"size:20" json:"evaluator"`
	CreatedAt  int64  `json:"createdAt"`
	ReviewedAt int64  `json:"reviewedAt"`
}

// 证据来源
const (
	EvidenceProject     = "project"
	EvidenceOpportunity = "opportunity"
)

// Evidence 能力证明：项目或真实任务产出的可验证证据。
type Evidence struct {
	ID         string `gorm:"primaryKey;size:64" json:"id"`
	UserID     string `gorm:"index;size:64" json:"userId"`
	SkillID    string `gorm:"index;size:64" json:"skillId"`
	SourceType string `gorm:"size:20" json:"sourceType"`
	SourceID   string `gorm:"size:64" json:"sourceId"`
	Title      string `gorm:"size:200" json:"title"`
	Score      int    `json:"score"`
	Detail     string `gorm:"type:text" json:"detail"`
	CreatedAt  int64  `json:"createdAt"`
}

// ---------- IT 任务平台 ----------

// 任务市场状态
const (
	OpportunityOpen       = "open"
	OpportunityInProgress = "in_progress"
	OpportunityClosed     = "closed"
)

// Opportunity 企业发布的真实 IT 任务。
type Opportunity struct {
	ID           string         `gorm:"primaryKey;size:64" json:"id"`
	Title        string         `gorm:"size:200" json:"title"`
	Company      string         `gorm:"size:100" json:"company"`
	PublisherID  string         `gorm:"index;size:64" json:"publisherId"`
	Description  string         `gorm:"type:text" json:"description"`
	Requirements datatypes.JSON `gorm:"type:json" json:"requirements"`
	Budget       float64        `json:"budget"`
	DurationDays int            `json:"durationDays"`
	Mode         string         `gorm:"size:20" json:"mode"`
	Status       string         `gorm:"size:20;index" json:"status"`
	CreatedAt    int64          `json:"createdAt"`
	UpdatedAt    int64          `json:"updatedAt"`
}

// 申请状态
const (
	ApplicationApplied   = "applied"
	ApplicationAccepted  = "accepted"
	ApplicationRejected  = "rejected"
	ApplicationCompleted = "completed"
)

// OpportunityApplication 人才对任务的申请与履约记录。
type OpportunityApplication struct {
	ID            string  `gorm:"primaryKey;size:64" json:"id"`
	OpportunityID string  `gorm:"index;size:64" json:"opportunityId"`
	UserID        string  `gorm:"index;size:64" json:"userId"`
	Message       string  `gorm:"type:text" json:"message"`
	MatchScore    int     `json:"matchScore"`
	Status        string  `gorm:"size:20;index" json:"status"`
	Rating        int     `json:"rating"`
	Review        string  `gorm:"type:text" json:"review"`
	Payout        float64 `json:"payout"`
	CreatedAt     int64   `json:"createdAt"`
	UpdatedAt     int64   `json:"updatedAt"`
}

// IncomeRecord 收入记录：成长飞轮中"收益"环节的数据。
type IncomeRecord struct {
	ID         string  `gorm:"primaryKey;size:64" json:"id"`
	UserID     string  `gorm:"index;size:64" json:"userId"`
	SourceType string  `gorm:"size:20" json:"sourceType"`
	SourceID   string  `gorm:"size:64" json:"sourceId"`
	Title      string  `gorm:"size:200" json:"title"`
	Amount     float64 `json:"amount"`
	CreatedAt  int64   `json:"createdAt"`
}

// ---------- 知识社区 ----------

// 帖子类型
const (
	PostQuestion   = "question"
	PostDiscussion = "discussion"
	PostArticle    = "article"
	PostShare      = "share"
)

// CommunityPost 社区内容：问答、讨论、文章、分享。
type CommunityPost struct {
	ID           string         `gorm:"primaryKey;size:64" json:"id"`
	AuthorID     string         `gorm:"index;size:64" json:"authorId"`
	Type         string         `gorm:"size:20;index" json:"type"`
	Title        string         `gorm:"size:200" json:"title"`
	Content      string         `gorm:"type:text" json:"content"`
	Tags         datatypes.JSON `gorm:"type:json" json:"tags"`
	KnowledgeIDs datatypes.JSON `gorm:"type:json" json:"knowledgeIds"`
	Likes        int            `json:"likes"`
	Views        int            `json:"views"`
	AnswerCount  int            `json:"answerCount"`
	AISummary    string         `gorm:"type:text" json:"aiSummary"`
	ResourceID   string         `gorm:"size:64" json:"resourceId"`
	Status       string         `gorm:"size:20;index" json:"status"`
	CreatedAt    int64          `gorm:"index" json:"createdAt"`
	UpdatedAt    int64          `json:"updatedAt"`
	Author       *UserBrief     `gorm:"-" json:"author,omitempty"`
}

// CommunityAnswer 回答或评论。
type CommunityAnswer struct {
	ID        string     `gorm:"primaryKey;size:64" json:"id"`
	PostID    string     `gorm:"index;size:64" json:"postId"`
	AuthorID  string     `gorm:"index;size:64" json:"authorId"`
	Content   string     `gorm:"type:text" json:"content"`
	Likes     int        `json:"likes"`
	Accepted  bool       `json:"accepted"`
	CreatedAt int64      `json:"createdAt"`
	Author    *UserBrief `gorm:"-" json:"author,omitempty"`
}

// UserBrief 对外展示的用户摘要。
type UserBrief struct {
	ID          string `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	AvatarColor string `json:"avatarColor"`
	Role        string `json:"role"`
}

// ---------- 知识资产 ----------

// 资源类型
var ResourceTypes = []string{"article", "course", "video", "code", "case", "project", "task", "prompt", "sop"}

// KnowledgeResource 统一知识资产：文章、课程、视频、代码、案例、项目、任务、Prompt、SOP。
type KnowledgeResource struct {
	ID           string         `gorm:"primaryKey;size:64" json:"id"`
	Type         string         `gorm:"size:20;index" json:"type"`
	Title        string         `gorm:"size:200" json:"title"`
	Summary      string         `gorm:"type:text" json:"summary"`
	Content      string         `gorm:"type:text" json:"content,omitempty"`
	URL          string         `gorm:"size:300" json:"url"`
	AuthorID     string         `gorm:"index;size:64" json:"authorId"`
	SourceType   string         `gorm:"size:20" json:"sourceType"`
	SourceID     string         `gorm:"size:64" json:"sourceId"`
	KnowledgeIDs datatypes.JSON `gorm:"type:json" json:"knowledgeIds"`
	SkillIDs     datatypes.JSON `gorm:"type:json" json:"skillIds"`
	Tags         datatypes.JSON `gorm:"type:json" json:"tags"`
	Views        int            `json:"views"`
	Likes        int            `json:"likes"`
	Status       string         `gorm:"size:20;index" json:"status"`
	CreatedAt    int64          `gorm:"index" json:"createdAt"`
	Author       *UserBrief     `gorm:"-" json:"author,omitempty"`
}

// ---------- AI 内核 ----------

// PipelineRun AI Pipeline 运行记录（含每个 Stage 的执行轨迹）。
type PipelineRun struct {
	ID         string         `gorm:"primaryKey;size:64" json:"id"`
	UserID     string         `gorm:"index;size:64" json:"userId"`
	Pipeline   string         `gorm:"size:50;index" json:"pipeline"`
	Status     string         `gorm:"size:20" json:"status"`
	Trace      datatypes.JSON `gorm:"type:json" json:"trace"`
	Output     datatypes.JSON `gorm:"type:json" json:"output,omitempty"`
	DurationMs int64          `json:"durationMs"`
	CreatedAt  int64          `gorm:"index" json:"createdAt"`
}
