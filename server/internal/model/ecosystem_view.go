package model

// 本文件为成长生态的接口视图与请求体。

// KnowledgeBrief 知识点摘要。
type KnowledgeBrief struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Domain           string `json:"domain"`
	Difficulty       int    `json:"difficulty"`
	EstimatedMinutes int    `json:"estimatedMinutes"`
}

// KnowledgeStateView 带遗忘衰减的知识掌握视图。
type KnowledgeStateView struct {
	Knowledge        KnowledgeBrief `json:"knowledge"`
	StudyCount       int            `json:"studyCount"`
	PracticeCount    int            `json:"practiceCount"`
	Accuracy         float64        `json:"accuracy"`
	Mastery          float64        `json:"mastery"`
	EffectiveMastery float64        `json:"effectiveMastery"`
	Confidence       float64        `json:"confidence"`
	Status           string         `json:"status"`
	Due              bool           `json:"due"`
	LastStudiedAt    int64          `json:"lastStudiedAt"`
	NextReviewAt     int64          `json:"nextReviewAt"`
}

// SkillStateView 技能状态视图。
type SkillStateView struct {
	Skill          Skill   `json:"skill"`
	Level          int     `json:"level"`
	LevelName      string  `json:"levelName"`
	Score          float64 `json:"score"`
	KnowledgeScore float64 `json:"knowledgeScore"`
	ProjectScore   float64 `json:"projectScore"`
	TaskScore      float64 `json:"taskScore"`
	EvidenceCount  int     `json:"evidenceCount"`
}

// SkillGapItem 单项技能差距。
type SkillGapItem struct {
	SkillID        string  `json:"skillId"`
	SkillName      string  `json:"skillName"`
	CapabilityID   string  `json:"capabilityId,omitempty"`
	CapabilityName string  `json:"capabilityName,omitempty"`
	Required       int     `json:"required"`
	Current        int     `json:"current"`
	Gap            int     `json:"gap"`
	Fit            float64 `json:"fit"`
	Weight         int     `json:"weight"`
}

// CapabilityReadiness 能力达成度。
type CapabilityReadiness struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Weight    int            `json:"weight"`
	Readiness int            `json:"readiness"`
	Skills    []SkillGapItem `json:"skills"`
}

// GapAnalysis 职业目标的能力差距分析。
type GapAnalysis struct {
	Career       *Career               `json:"career"`
	Role         *JobRole              `json:"role"`
	Readiness    int                   `json:"readiness"`
	MetCount     int                   `json:"metCount"`
	TotalCount   int                   `json:"totalCount"`
	Capabilities []CapabilityReadiness `json:"capabilities"`
	Gaps         []SkillGapItem        `json:"gaps"`
}

// KnowledgeRecommendation 推荐学习的知识点。
type KnowledgeRecommendation struct {
	Knowledge      KnowledgeBrief   `json:"knowledge"`
	SkillID        string           `json:"skillId"`
	SkillName      string           `json:"skillName"`
	Mastery        float64          `json:"mastery"`
	Ready          bool             `json:"ready"`
	MissingPrereqs []KnowledgeBrief `json:"missingPrereqs,omitempty"`
	Reason         string           `json:"reason"`
}

// PlanItem 学习计划项。
type PlanItem struct {
	Type    string `json:"type"` // knowledge | project
	RefID   string `json:"refId"`
	SubID   string `json:"subId,omitempty"`
	Title   string `json:"title"`
	Minutes int    `json:"minutes"`
	Reason  string `json:"reason,omitempty"`
}

// PlanWeek 一周的学习计划。
type PlanWeek struct {
	Week  int        `json:"week"`
	Focus string     `json:"focus"`
	Items []PlanItem `json:"items"`
}

// LearningPlan AI 学习规划结果。
type LearningPlan struct {
	Summary      string     `json:"summary"`
	TotalMinutes int        `json:"totalMinutes"`
	Weeks        []PlanWeek `json:"weeks"`
}

// GoalRequest 设置职业目标。
type GoalRequest struct {
	RoleID      string `json:"roleId"`
	WeeklyHours int    `json:"weeklyHours"`
	TargetWeeks int    `json:"targetWeeks"`
	Motivation  string `json:"motivation"`
}

// LearningEventRequest 上报学习事件。
type LearningEventRequest struct {
	KnowledgeID string `json:"knowledgeId"`
	Type        string `json:"type"`
	Correct     int    `json:"correct"`
	Total       int    `json:"total"`
	Minutes     int    `json:"minutes"`
	Source      string `json:"source"`
}

// PracticeAnswerRequest 提交知识点练习答案（服务端判分）。
type PracticeAnswerRequest struct {
	Answers []int `json:"answers"`
}

// PracticeResult 练习判分结果。
type PracticeResult struct {
	Correct int                `json:"correct"`
	Total   int                `json:"total"`
	Results []bool             `json:"results"`
	State   KnowledgeStateView `json:"state"`
}

// ChapterCompleteRequest 课程章节完成 / 测验成绩上报。
type ChapterCompleteRequest struct {
	CourseID  string `json:"courseId"`
	ChapterID string `json:"chapterId"`
	Correct   int    `json:"correct"`
	Total     int    `json:"total"`
}

// GraphNode 知识图谱节点（含个人掌握度）。
type GraphNode struct {
	KnowledgeBrief
	Summary  string   `json:"summary"`
	SkillIDs []string `json:"skillIds"`
	Depth    int      `json:"depth"`
	Mastery  float64  `json:"mastery"`
	Status   string   `json:"status"`
}

// KnowledgeGraph 知识图谱。
type KnowledgeGraph struct {
	Nodes []GraphNode         `json:"nodes"`
	Edges []KnowledgeRelation `json:"edges"`
}

// KnowledgeDetail 知识点详情。
type KnowledgeDetail struct {
	Knowledge KnowledgePoint      `json:"knowledge"`
	Questions []PracticeQuestion  `json:"questions"`
	State     KnowledgeStateView  `json:"state"`
	Skills    []Skill             `json:"skills"`
	Relations []KnowledgeLink     `json:"relations"`
	Chapters  []ChapterRef        `json:"chapters"`
	Resources []KnowledgeResource `json:"resources"`
}

// PracticeQuestion 不含答案的练习题。
type PracticeQuestion struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
}

// KnowledgeLink 知识关系（带方向说明）。
type KnowledgeLink struct {
	Type      string         `json:"type"`
	Direction string         `json:"direction"` // in | out
	Knowledge KnowledgeBrief `json:"knowledge"`
	Mastery   float64        `json:"mastery"`
}

// ChapterRef 引用某知识点的课程章节。
type ChapterRef struct {
	CourseID     string `json:"courseId"`
	CourseTitle  string `json:"courseTitle"`
	ChapterID    string `json:"chapterId"`
	ChapterTitle string `json:"chapterTitle"`
}

// RoleDetail 岗位详情：能力 → 技能要求，以及个人当前等级。
type RoleDetail struct {
	Career *Career          `json:"career"`
	Role   *JobRole         `json:"role"`
	Gap    *GapAnalysis     `json:"gap,omitempty"`
	Skills map[string]Skill `json:"skills"`
}

// MatchView 人才与任务的匹配结果。
type MatchView struct {
	Score     int            `json:"score"`
	Qualified bool           `json:"qualified"`
	MetCount  int            `json:"metCount"`
	Items     []SkillGapItem `json:"items"`
}

// OpportunityView 任务市场列表项。
type OpportunityView struct {
	Opportunity
	Requirements   []SkillRequirement      `json:"requirements"`
	Match          *MatchView              `json:"match,omitempty"`
	ApplicantCount int                     `json:"applicantCount"`
	MyApplication  *OpportunityApplication `json:"myApplication,omitempty"`
}

// OpportunityRequest 发布/编辑 IT 任务。
type OpportunityRequest struct {
	Title        string             `json:"title"`
	Company      string             `json:"company"`
	Description  string             `json:"description"`
	Requirements []SkillRequirement `json:"requirements"`
	Budget       float64            `json:"budget"`
	DurationDays int                `json:"durationDays"`
	Mode         string             `json:"mode"`
	Status       string             `json:"status"`
}

// ApplyRequest 申请任务。
type ApplyRequest struct {
	Message string `json:"message"`
}

// ApplicationDecision 企业处理申请：accept / reject / complete。
type ApplicationDecision struct {
	Action string  `json:"action"`
	Rating int     `json:"rating"`
	Review string  `json:"review"`
	Payout float64 `json:"payout"`
}

// ApplicationView 申请记录（含人才摘要）。
type ApplicationView struct {
	OpportunityApplication
	User        *UserBrief   `json:"user,omitempty"`
	Opportunity *Opportunity `json:"opportunity,omitempty"`
	Match       *MatchView   `json:"match,omitempty"`
}

// CandidateView 任务推荐人才。
type CandidateView struct {
	User   UserBrief        `json:"user"`
	Match  MatchView        `json:"match"`
	Goal   string           `json:"goal"`
	Skills []SkillStateView `json:"skills"`
}

// SubmissionRequest 提交项目任务。
type SubmissionRequest struct {
	Content string `json:"content"`
	RepoURL string `json:"repoUrl"`
}

// SubmissionView 提交记录与评审结果。
type SubmissionView struct {
	TaskSubmission
	TaskTitle    string     `json:"taskTitle"`
	ProjectTitle string     `json:"projectTitle"`
	Evidence     []Evidence `json:"evidence,omitempty"`
}

// ProjectTaskView 项目任务（含要求与个人进度）。
type ProjectTaskView struct {
	ProjectTask
	Requirements []SkillRequirement `json:"requirements"`
	Match        *MatchView         `json:"match,omitempty"`
	BestScore    int                `json:"bestScore"`
	Submitted    bool               `json:"submitted"`
}

// ProjectView 项目视图。
type ProjectView struct {
	Project
	SkillIDs  []string          `json:"skillIds"`
	Tasks     []ProjectTaskView `json:"tasks"`
	Progress  int               `json:"progress"`
	Completed int               `json:"completed"`
}

// ProjectRequest 管理端编辑项目。
type ProjectRequest struct {
	Title       string               `json:"title"`
	Summary     string               `json:"summary"`
	Description string               `json:"description"`
	Domain      string               `json:"domain"`
	Difficulty  int                  `json:"difficulty"`
	Icon        string               `json:"icon"`
	Color       string               `json:"color"`
	SkillIDs    []string             `json:"skillIds"`
	Status      string               `json:"status"`
	Tasks       []ProjectTaskRequest `json:"tasks"`
}

// ProjectTaskRequest 管理端编辑任务。
type ProjectTaskRequest struct {
	ID             string             `json:"id"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	Deliverable    string             `json:"deliverable"`
	Requirements   []SkillRequirement `json:"requirements"`
	EstimatedHours int                `json:"estimatedHours"`
}

// PostRequest 发帖。
type PostRequest struct {
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Tags         []string `json:"tags"`
	KnowledgeIDs []string `json:"knowledgeIds"`
}

// AnswerRequest 回答。
type AnswerRequest struct {
	Content string `json:"content"`
}

// PostDetail 帖子详情。
type PostDetail struct {
	Post      CommunityPost     `json:"post"`
	Answers   []CommunityAnswer `json:"answers"`
	Knowledge []KnowledgeBrief  `json:"knowledge"`
}

// ResourceRequest 发布知识资产。
type ResourceRequest struct {
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	Content      string   `json:"content"`
	URL          string   `json:"url"`
	Tags         []string `json:"tags"`
	KnowledgeIDs []string `json:"knowledgeIds"`
	SkillIDs     []string `json:"skillIds"`
}

// FlywheelStats 个人成长飞轮：学习 → 实践 → 产出 → 收益 → 投资学习 → 更强能力。
type FlywheelStats struct {
	LearningEvents    int64   `json:"learningEvents"`
	LearningMinutes   int     `json:"learningMinutes"`
	MasteredCount     int     `json:"masteredCount"`
	ProjectsPracticed int     `json:"projectsPracticed"`
	Submissions       int     `json:"submissions"`
	EvidenceCount     int     `json:"evidenceCount"`
	TasksCompleted    int     `json:"tasksCompleted"`
	Income            float64 `json:"income"`
	Contributions     int     `json:"contributions"`
	AvgSkillLevel     float64 `json:"avgSkillLevel"`
}

// ActivityDay 每日学习活跃度。
type ActivityDay struct {
	Date    string `json:"date"`
	Events  int    `json:"events"`
	Minutes int    `json:"minutes"`
}

// ContributionStats 社区贡献。
type ContributionStats struct {
	Posts       int64 `json:"posts"`
	Answers     int64 `json:"answers"`
	AnswerLikes int64 `json:"answerLikes"`
	Resources   int64 `json:"resources"`
}

// TalentProfile 人才画像：比传统简历更准确的动态能力档案。
type TalentProfile struct {
	User          UserBrief         `json:"user"`
	JoinedAt      int64             `json:"joinedAt"`
	Goal          *UserGoal         `json:"goal"`
	GoalCareer    *Career           `json:"goalCareer,omitempty"`
	GoalRole      *JobRole          `json:"goalRole,omitempty"`
	Readiness     int               `json:"readiness"`
	Skills        []SkillStateView  `json:"skills"`
	Knowledge     KnowledgeSummary  `json:"knowledge"`
	Projects      []SubmissionView  `json:"projects"`
	Tasks         []ApplicationView `json:"tasks"`
	Evidence      []Evidence        `json:"evidence"`
	Contributions ContributionStats `json:"contributions"`
	Income        []IncomeRecord    `json:"income"`
	TotalIncome   float64           `json:"totalIncome"`
	Flywheel      FlywheelStats     `json:"flywheel"`
	Activity      []ActivityDay     `json:"activity"`
}

// KnowledgeSummary 知识掌握概况。
type KnowledgeSummary struct {
	Total      int             `json:"total"`
	Mastered   int             `json:"mastered"`
	Learning   int             `json:"learning"`
	Due        int             `json:"due"`
	AvgMastery float64         `json:"avgMastery"`
	ByDomain   []DomainMastery `json:"byDomain"`
}

// DomainMastery 领域掌握度。
type DomainMastery struct {
	Domain  string  `json:"domain"`
	Count   int     `json:"count"`
	Mastery float64 `json:"mastery"`
}

// TalentBrief 人才库列表项。
type TalentBrief struct {
	User          UserBrief        `json:"user"`
	Goal          string           `json:"goal"`
	Readiness     int              `json:"readiness"`
	AvgLevel      float64          `json:"avgLevel"`
	TopSkills     []SkillStateView `json:"topSkills"`
	EvidenceCount int              `json:"evidenceCount"`
	TotalIncome   float64          `json:"totalIncome"`
}

// GrowthOverview 成长中心总览。
type GrowthOverview struct {
	Goal            *UserGoal                 `json:"goal"`
	Gap             *GapAnalysis              `json:"gap"`
	Recommendations []KnowledgeRecommendation `json:"recommendations"`
	DueReviews      []KnowledgeStateView      `json:"dueReviews"`
	Flywheel        FlywheelStats             `json:"flywheel"`
	Activity        []ActivityDay             `json:"activity"`
	NextProjects    []ProjectTaskRef          `json:"nextProjects"`
	Opportunities   []OpportunityView         `json:"opportunities"`
	LatestPlan      *PipelineRun              `json:"latestPlan,omitempty"`
}

// ProjectTaskRef 推荐的项目任务。
type ProjectTaskRef struct {
	ProjectID    string    `json:"projectId"`
	ProjectTitle string    `json:"projectTitle"`
	TaskID       string    `json:"taskId"`
	TaskTitle    string    `json:"taskTitle"`
	Match        MatchView `json:"match"`
}

// EcoDashboard 管理端生态看板。
type EcoDashboard struct {
	Careers           int64            `json:"careers"`
	Roles             int64            `json:"roles"`
	Skills            int64            `json:"skills"`
	Knowledge         int64            `json:"knowledge"`
	Relations         int64            `json:"relations"`
	Projects          int64            `json:"projects"`
	Submissions       int64            `json:"submissions"`
	Opportunities     int64            `json:"opportunities"`
	OpenOpportunities int64            `json:"openOpportunities"`
	Applications      int64            `json:"applications"`
	CompletedTasks    int64            `json:"completedTasks"`
	Posts             int64            `json:"posts"`
	Resources         map[string]int64 `json:"resources"`
	Goals             int64            `json:"goals"`
	LearningEvents    int64            `json:"learningEvents"`
	Evidence          int64            `json:"evidence"`
	TotalIncome       float64          `json:"totalIncome"`
	PipelineRuns      int64            `json:"pipelineRuns"`
	LevelDistribution []int            `json:"levelDistribution"`
}

// AgentInfo AI Agent 描述。
type AgentInfo struct {
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Responsibilities []string `json:"responsibilities"`
	Starters         []string `json:"starters"`
	Icon             string   `json:"icon"`
	Color            string   `json:"color"`
}

// AgentChatRequest 与 Agent 对话。
type AgentChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history,omitempty"`
}

// KernelRunRequest 运行 AI 学习内核。
type KernelRunRequest struct {
	Question string `json:"question"`
}

// KernelOutput 学习内核输出。
type KernelOutput struct {
	Gap             *GapAnalysis              `json:"gap"`
	Recommendations []KnowledgeRecommendation `json:"recommendations"`
	Plan            *LearningPlan             `json:"plan"`
	Response        string                    `json:"response"`
	ResponseSource  string                    `json:"responseSource"`
	Evaluation      KernelEvaluation          `json:"evaluation"`
	Profile         ProfileDelta              `json:"profile"`
}

// KernelEvaluation 对本次规划的自评估。
type KernelEvaluation struct {
	Coverage    int      `json:"coverage"`    // 推荐覆盖的差距技能比例
	Feasibility int      `json:"feasibility"` // 计划周期与目标周期的匹配度
	Notes       []string `json:"notes"`
}

// ProfileDelta 画像更新结果。
type ProfileDelta struct {
	SkillsUpdated int `json:"skillsUpdated"`
	Readiness     int `json:"readiness"`
}
