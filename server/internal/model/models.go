package model

import "gorm.io/datatypes"

type User struct {
	ID           string `gorm:"primaryKey;size:64" json:"id"`
	Username     string `gorm:"uniqueIndex;size:50" json:"username"`
	Nickname     string `gorm:"size:100" json:"nickname"`
	PasswordHash string `gorm:"size:255" json:"-"`
	Role         string `gorm:"size:20;index" json:"role"`
	Status       string `gorm:"size:20;index" json:"status"`
	Avatar       string `gorm:"size:10" json:"avatar"`
	AvatarColor  string `gorm:"size:20" json:"avatarColor"`
	JoinedAt     int64  `json:"joinedAt"`
}

type Course struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`
	Title            string         `gorm:"size:200" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	Level            string         `gorm:"size:20" json:"level"`
	Tags             datatypes.JSON `gorm:"type:json" json:"tags"`
	Icon             string         `gorm:"size:10" json:"icon"`
	Accent           string         `gorm:"size:20" json:"accent"`
	EstimatedMinutes int            `json:"estimatedMinutes"`
	Status           string         `gorm:"size:20;index" json:"status"`
	CreatedAt        int64          `json:"createdAt"`
	UpdatedAt        int64          `json:"updatedAt"`
	Chapters         []Chapter      `gorm:"foreignKey:CourseID" json:"chapters,omitempty"`
	ChapterCount     int            `gorm:"-" json:"chapterCount,omitempty"`
}

type Chapter struct {
	ID        string `gorm:"primaryKey;size:64" json:"id"`
	CourseID  string `gorm:"index;size:64" json:"courseId,omitempty"`
	Title     string `gorm:"size:200" json:"title"`
	Minutes   int    `json:"minutes"`
	Content   string `gorm:"type:text" json:"content,omitempty"`
	Status    string `gorm:"size:20" json:"status"`
	UpdatedAt int64  `json:"updatedAt"`
}

type Quiz struct {
	ID            string         `gorm:"primaryKey;size:64" json:"id"`
	CourseID      string         `gorm:"index;size:64" json:"courseId"`
	Title         string         `gorm:"size:200" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Questions     datatypes.JSON `gorm:"type:json" json:"questions,omitempty"`
	Status        string         `gorm:"size:20;index" json:"status"`
	CreatedAt     int64          `json:"createdAt"`
	UpdatedAt     int64          `json:"updatedAt"`
	QuestionCount int            `gorm:"-" json:"questionCount,omitempty"`
}

type Review struct {
	ID         string `gorm:"primaryKey;size:64" json:"id"`
	Type       string `gorm:"size:20;index" json:"type"`
	CourseID   string `gorm:"size:64" json:"courseId"`
	TargetID   string `gorm:"size:64" json:"targetId"`
	Title      string `gorm:"size:200" json:"title"`
	Content    string `gorm:"type:text" json:"content"`
	Submitter  string `gorm:"size:50" json:"submitter"`
	Status     string `gorm:"size:20;index" json:"status"`
	AIScore    int    `json:"aiScore"`
	AIFeedback string `gorm:"type:text" json:"aiFeedback"`
	ReviewerID string `gorm:"size:64" json:"reviewerId,omitempty"`
	Comment    string `gorm:"type:text" json:"comment,omitempty"`
	ReviewedAt int64  `json:"reviewedAt,omitempty"`
	CreatedAt  int64  `json:"createdAt"`
}

type Question struct {
	Text        string   `json:"text"`
	Options     []string `json:"options"`
	Answer      int      `json:"answer"`
	Explanation string   `json:"explanation"`
}

type PageResult[T any] struct {
	List     []T `json:"list"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type DashboardStats struct {
	UserStats   UserStats   `json:"userStats"`
	CourseStats CourseStats `json:"courseStats"`
	QuizStats   QuizStats   `json:"quizStats"`
	ReviewStats ReviewStats `json:"reviewStats"`
}

type UserStats struct {
	Total    int64 `json:"total"`
	Learners int64 `json:"learners"`
	Admins   int64 `json:"admins"`
	Active   int64 `json:"active"`
}

type CourseStats struct {
	Total     int64 `json:"total"`
	Published int64 `json:"published"`
	Draft     int64 `json:"draft"`
	Chapters  int64 `json:"chapters"`
}

type QuizStats struct {
	Total     int64 `json:"total"`
	Questions int64 `json:"questions"`
}

type ReviewStats struct {
	Pending  int64 `json:"pending"`
	Approved int64 `json:"approved"`
	Rejected int64 `json:"rejected"`
}

type AuthUser struct {
	User
	Permissions []string `json:"permissions"`
	RoleName    string   `json:"roleName"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type RolePermission struct {
	Role        string         `gorm:"primaryKey;size:20" json:"role"`
	Permissions datatypes.JSON `gorm:"type:json" json:"permissions"`
	UpdatedAt   int64          `json:"updatedAt"`
}

type RoleInfo struct {
	Role        string   `json:"role"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type RegisterRequest struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	AvatarColor string `json:"avatarColor"`
}

type AISource struct {
	CourseID     string `json:"courseId"`
	CourseTitle  string `json:"courseTitle"`
	ChapterID    string `json:"chapterId"`
	ChapterTitle string `json:"chapterTitle"`
}

type ChatRequest struct {
	Question string        `json:"question"`
	History  []ChatMessage `json:"history,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResult struct {
	Text    string     `json:"text"`
	Sources []AISource `json:"sources"`
}

type CareerInterviewRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history,omitempty"`
}

type CareerRecommendRequest struct {
	Background  string `json:"background"`
	Interest    string `json:"interest"`
	Experience  string `json:"experience"`
	WeeklyHours int    `json:"weeklyHours"`
}

type CareerRecommendation struct {
	CareerID   string `json:"careerId"`
	Name       string `json:"name"`
	MatchScore int    `json:"matchScore"`
	Reason     string `json:"reason"`
}

type CareerRecommendResult struct {
	Summary         string                 `json:"summary"`
	Recommendations []CareerRecommendation `json:"recommendations"`
}

type GoalDecomposeRequest struct {
	CareerID      string `json:"careerId"`
	CareerName    string `json:"careerName"`
	BaseLevel     string `json:"baseLevel"`
	WeeklyHours   int    `json:"weeklyHours"`
	DurationWeeks int    `json:"durationWeeks"`
}

type LearningStage struct {
	Name          string   `json:"name"`
	DurationWeeks int      `json:"durationWeeks"`
	Objectives    []string `json:"objectives"`
	Topics        []string `json:"topics"`
}

type GoalDecomposeResult struct {
	GoalName    string          `json:"goalName"`
	Difficulty  string          `json:"difficulty"`
	Stages      []LearningStage `json:"stages"`
	Milestones  []string        `json:"milestones"`
	AISummary   string          `json:"aiSummary"`
	Suggestions []string        `json:"suggestions"`
}

type LearningSuggestRequest struct {
	CompetencyProgress []CompetencyProgress `json:"competencyProgress"`
	NextMilestone      string               `json:"nextMilestone"`
	Streak             int                  `json:"streak"`
}

type CompetencyProgress struct {
	Name     string `json:"name"`
	Progress int    `json:"progress"`
}

type LearningSuggestResult struct {
	Suggestions []string `json:"suggestions"`
}

type SystemSetting struct {
	Key       string `gorm:"primaryKey;size:64" json:"key"`
	Value     string `gorm:"type:text" json:"value"`
	UpdatedAt int64  `json:"updatedAt"`
}

type SettingsView struct {
	LLM       LLMSettingsView `json:"llm"`
	UpdatedAt int64           `json:"updatedAt"`
}

type LLMSettingsView struct {
	APIKeyConfigured bool   `json:"apiKeyConfigured"`
	APIKeyMasked     string `json:"apiKeyMasked,omitempty"`
	BaseURL          string `json:"baseUrl"`
	Model            string `json:"model"`
	Enabled          bool   `json:"enabled"`
	Source           string `json:"source"`
}

type SettingsUpdateRequest struct {
	LLM LLMSettingsUpdate `json:"llm"`
}

type LLMSettingsUpdate struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
}

// CustomerTicket 客户咨询工单
type CustomerTicket struct {
	ID            string `gorm:"primaryKey;size:64" json:"id"`
	UserID        string `gorm:"index;size:64" json:"userId"`
	Subject       string `gorm:"size:200" json:"subject"`
	Status        string `gorm:"size:20;index" json:"status"`
	LastMessageAt int64  `json:"lastMessageAt"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
	// 关联展示字段（非持久化）
	UserNickname string `gorm:"-" json:"userNickname,omitempty"`
	UserUsername string `gorm:"-" json:"userUsername,omitempty"`
	LastMessage  string `gorm:"-" json:"lastMessage,omitempty"`
	UnreadCount  int    `gorm:"-" json:"unreadCount,omitempty"`
}

// CustomerMessage 客户咨询消息
type CustomerMessage struct {
	ID         string `gorm:"primaryKey;size:64" json:"id"`
	TicketID   string `gorm:"index;size:64" json:"ticketId"`
	SenderID   string `gorm:"size:64" json:"senderId"`
	SenderRole string `gorm:"size:20" json:"senderRole"`
	Content    string `gorm:"type:text" json:"content"`
	CreatedAt  int64  `json:"createdAt"`
	// 关联展示字段
	SenderNickname string `gorm:"-" json:"senderNickname,omitempty"`
}

type CustomerTicketStats struct {
	Total   int64 `json:"total"`
	Open    int64 `json:"open"`
	Pending int64 `json:"pending"`
	Closed  int64 `json:"closed"`
}

// Document 业务单据
type Document struct {
	ID        string  `gorm:"primaryKey;size:64" json:"id"`
	DocNo     string  `gorm:"uniqueIndex;size:50" json:"docNo"`
	Title     string  `gorm:"size:200" json:"title"`
	Type      string  `gorm:"size:30;index" json:"type"`
	Amount    float64 `json:"amount"`
	Status    string  `gorm:"size:20;index" json:"status"`
	Remark    string  `gorm:"type:text" json:"remark"`
	CreatedBy string  `gorm:"size:50" json:"createdBy"`
	CreatedAt int64   `json:"createdAt"`
	UpdatedAt int64   `json:"updatedAt"`
}

// ImportTaskProgress Excel 导入任务进度
type ImportTaskProgress struct {
	TaskID   string   `json:"taskId"`
	Status   string   `json:"status"` // pending, processing, completed, failed
	Total    int      `json:"total"`
	Current  int      `json:"current"`
	Success  int      `json:"success"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
	Message  string   `json:"message,omitempty"`
}
