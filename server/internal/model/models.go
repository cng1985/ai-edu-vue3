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

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
