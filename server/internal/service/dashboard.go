package service

import (
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
)

// DashboardService 汇总运营看板统计数据。
type DashboardService struct {
	users   *repository.UserRepo
	courses *repository.CourseRepo
	quizzes *repository.QuizRepo
	reviews *repository.ReviewRepo
}

func NewDashboardService(users *repository.UserRepo, courses *repository.CourseRepo, quizzes *repository.QuizRepo, reviews *repository.ReviewRepo) *DashboardService {
	return &DashboardService{users: users, courses: courses, quizzes: quizzes, reviews: reviews}
}

func (s *DashboardService) Stats() (*model.DashboardStats, error) {
	totalUsers, _ := s.users.Total()
	learners, _ := s.users.CountByRole(rbac.RoleLearner)
	admins, _ := s.users.CountAdmins()
	active, _ := s.users.CountActive()

	totalCourses, _ := s.courses.Total()
	published, _ := s.courses.CountByStatus("published")
	draft, _ := s.courses.CountByStatus("draft")
	chapters, _ := s.courses.TotalChapters()

	totalQuizzes, _ := s.quizzes.Total()
	quizzes, _, _ := s.quizzes.List("", "", "", 1, 1000)
	questionCount := int64(0)
	for _, q := range quizzes {
		questionCount += int64(countQuestions(q.Questions))
	}

	pending, _ := s.reviews.CountByStatus("pending")
	approved, _ := s.reviews.CountByStatus("approved")
	rejected, _ := s.reviews.CountByStatus("rejected")

	return &model.DashboardStats{
		UserStats:   model.UserStats{Total: totalUsers, Learners: learners, Admins: admins, Active: active},
		CourseStats: model.CourseStats{Total: totalCourses, Published: published, Draft: draft, Chapters: chapters},
		QuizStats:   model.QuizStats{Total: totalQuizzes, Questions: questionCount},
		ReviewStats: model.ReviewStats{Pending: pending, Approved: approved, Rejected: rejected},
	}, nil
}
