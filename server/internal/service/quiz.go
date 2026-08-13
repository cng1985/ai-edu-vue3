package service

import (
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
)

// QuizService 负责测验管理。
type QuizService struct{ quizzes *repository.QuizRepo }

func NewQuizService(quizzes *repository.QuizRepo) *QuizService { return &QuizService{quizzes: quizzes} }

func (s *QuizService) List(keyword, courseID, status string, page, pageSize int) (*model.PageResult[model.Quiz], error) {
	page, pageSize = normalizePage(page, pageSize)
	quizzes, total, err := s.quizzes.List(keyword, courseID, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	if quizzes == nil {
		quizzes = []model.Quiz{}
	}
	for i := range quizzes {
		quizzes[i].QuestionCount = countQuestions(quizzes[i].Questions)
	}
	return &model.PageResult[model.Quiz]{List: quizzes, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *QuizService) Get(id string) (*model.Quiz, error) {
	quiz, err := s.quizzes.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("测验不存在")
	}
	quiz.QuestionCount = countQuestions(quiz.Questions)
	return quiz, nil
}

func (s *QuizService) Create(quiz model.Quiz) (*model.Quiz, error) {
	if quiz.Title == "" || quiz.CourseID == "" {
		return nil, apperr.BadRequest("标题和关联课程必填")
	}
	if quiz.ID == "" {
		quiz.ID = genID("quiz")
	}
	if s.quizzes.Exists(quiz.ID) {
		return nil, apperr.BadRequest("测验 ID 已存在")
	}
	now := time.Now().UnixMilli()
	quiz.Status = defaultStr(quiz.Status, "draft")
	quiz.CreatedAt = now
	quiz.UpdatedAt = now
	if err := s.quizzes.Create(&quiz); err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (s *QuizService) Update(id string, req model.Quiz) (*model.Quiz, error) {
	quiz, err := s.quizzes.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("测验不存在")
	}
	if req.Title != "" {
		quiz.Title = req.Title
	}
	if req.Description != "" {
		quiz.Description = req.Description
	}
	if req.CourseID != "" {
		quiz.CourseID = req.CourseID
	}
	if len(req.Questions) > 0 {
		quiz.Questions = req.Questions
	}
	if req.Status != "" {
		quiz.Status = req.Status
	}
	quiz.UpdatedAt = time.Now().UnixMilli()
	if err := s.quizzes.Update(quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *QuizService) Delete(id string) error {
	if _, err := s.quizzes.FindByID(id); err != nil {
		return apperr.NotFound("测验不存在")
	}
	return s.quizzes.Delete(id)
}
