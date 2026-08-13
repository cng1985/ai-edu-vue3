package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"gorm.io/datatypes"
)

// ReviewService 负责内容审核流程；审核通过后将内容发布到课程/题库。
type ReviewService struct {
	reviews   *repository.ReviewRepo
	courses   *repository.CourseRepo
	quizzes   *repository.QuizRepo
	knowledge *KnowledgeService
}

func NewReviewService(reviews *repository.ReviewRepo, courses *repository.CourseRepo, quizzes *repository.QuizRepo, knowledge *KnowledgeService) *ReviewService {
	return &ReviewService{reviews: reviews, courses: courses, quizzes: quizzes, knowledge: knowledge}
}

func (s *ReviewService) List(status, reviewType string, page, pageSize int) (*model.PageResult[model.Review], error) {
	page, pageSize = normalizePage(page, pageSize)
	reviews, total, err := s.reviews.List(status, reviewType, page, pageSize)
	if err != nil {
		return nil, err
	}
	if reviews == nil {
		reviews = []model.Review{}
	}
	return &model.PageResult[model.Review]{List: reviews, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *ReviewService) Get(id string) (*model.Review, error) {
	review, err := s.reviews.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("审核记录不存在")
	}
	return review, nil
}

func (s *ReviewService) Approve(id, reviewerID, comment string) (*model.Review, error) {
	review, err := s.pendingReview(id)
	if err != nil {
		return nil, err
	}
	review.Status = "approved"
	review.ReviewerID = reviewerID
	review.ReviewedAt = time.Now().UnixMilli()
	review.Comment = comment

	switch review.Type {
	case "chapter":
		s.publishChapter(review)
	case "quiz":
		s.appendQuizQuestion(review)
	}

	if err := s.reviews.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) Reject(id, reviewerID, comment string) (*model.Review, error) {
	review, err := s.pendingReview(id)
	if err != nil {
		return nil, err
	}
	review.Status = "rejected"
	review.ReviewerID = reviewerID
	review.ReviewedAt = time.Now().UnixMilli()
	review.Comment = defaultStr(comment, "内容不符合发布标准")
	if err := s.reviews.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) pendingReview(id string) (*model.Review, error) {
	review, err := s.reviews.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("审核记录不存在")
	}
	if review.Status != "pending" {
		return nil, apperr.BadRequest("该记录已处理")
	}
	return review, nil
}

func (s *ReviewService) publishChapter(review *model.Review) {
	ch, err := s.courses.FindChapter(review.CourseID, review.TargetID)
	if err != nil {
		return
	}
	ch.Content = review.Content
	ch.Status = "published"
	ch.UpdatedAt = time.Now().UnixMilli()
	_ = s.courses.UpdateChapter(ch)
	if s.knowledge != nil {
		_ = s.knowledge.IndexChapter(context.Background(), review.CourseID, review.TargetID)
	}
}

func (s *ReviewService) appendQuizQuestion(review *model.Review) {
	var question model.Question
	if err := json.Unmarshal([]byte(review.Content), &question); err != nil || question.Text == "" {
		return
	}
	quiz, err := s.quizzes.FindByID(review.TargetID)
	if err != nil {
		return
	}
	var questions []model.Question
	_ = json.Unmarshal(quiz.Questions, &questions)
	questions = append(questions, question)
	b, _ := json.Marshal(questions)
	quiz.Questions = datatypes.JSON(b)
	quiz.UpdatedAt = time.Now().UnixMilli()
	_ = s.quizzes.Update(quiz)
}
