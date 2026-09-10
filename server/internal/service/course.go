package service

import (
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
)

// CourseService 负责课程与章节管理。
type CourseService struct {
	courses *repository.CourseRepo
}

func NewCourseService(courses *repository.CourseRepo) *CourseService {
	return &CourseService{courses: courses}
}

func (s *CourseService) List(keyword, status string, page, pageSize int) (*model.PageResult[model.Course], error) {
	page, pageSize = normalizePage(page, pageSize)
	courses, total, err := s.courses.List(keyword, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	if courses == nil {
		courses = []model.Course{}
	}
	return &model.PageResult[model.Course]{List: courses, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *CourseService) Get(id string) (*model.Course, error) {
	course, err := s.courses.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("课程不存在")
	}
	return course, nil
}

func (s *CourseService) Create(course model.Course) (*model.Course, error) {
	if course.Title == "" {
		return nil, apperr.BadRequest("课程标题必填")
	}
	if course.ID == "" {
		course.ID = genID("course")
	}
	if s.courses.Exists(course.ID) {
		return nil, apperr.BadRequest("课程 ID 已存在")
	}
	now := time.Now().UnixMilli()
	course.Level = defaultStr(course.Level, "入门")
	course.Icon = defaultStr(course.Icon, "📚")
	course.Accent = defaultStr(course.Accent, "#6366f1")
	course.Status = defaultStr(course.Status, "draft")
	if course.EstimatedMinutes == 0 {
		course.EstimatedMinutes = 60
	}
	course.CreatedAt = now
	course.UpdatedAt = now
	if err := s.courses.Create(&course); err != nil {
		return nil, err
	}
	return &course, nil
}

func (s *CourseService) Update(id string, req model.Course) (*model.Course, error) {
	course, err := s.courses.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("课程不存在")
	}
	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Level != "" {
		course.Level = req.Level
	}
	if len(req.Tags) > 0 {
		course.Tags = req.Tags
	}
	if req.Icon != "" {
		course.Icon = req.Icon
	}
	if req.Accent != "" {
		course.Accent = req.Accent
	}
	if req.EstimatedMinutes > 0 {
		course.EstimatedMinutes = req.EstimatedMinutes
	}
	if req.Status != "" {
		course.Status = req.Status
	}
	if len(req.Chapters) > 0 {
		course.Chapters = req.Chapters
	}
	course.UpdatedAt = time.Now().UnixMilli()
	if err := s.courses.Update(course); err != nil {
		return nil, err
	}
	return s.courses.FindByID(id)
}

func (s *CourseService) Delete(id string) error {
	if _, err := s.courses.FindByID(id); err != nil {
		return apperr.NotFound("课程不存在")
	}
	return s.courses.Delete(id)
}

func (s *CourseService) AddChapter(courseID string, ch model.Chapter) (*model.Chapter, error) {
	if _, err := s.courses.FindByID(courseID); err != nil {
		return nil, apperr.NotFound("课程不存在")
	}
	if ch.Title == "" {
		return nil, apperr.BadRequest("章节标题必填")
	}
	if ch.ID == "" {
		ch.ID = genID("ch")
	}
	if s.courses.ChapterExists(courseID, ch.ID) {
		return nil, apperr.BadRequest("章节 ID 已存在")
	}
	if ch.Minutes == 0 {
		ch.Minutes = 10
	}
	ch.CourseID = courseID
	ch.Status = "draft"
	ch.UpdatedAt = time.Now().UnixMilli()
	if err := s.courses.AddChapter(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *CourseService) UpdateChapter(courseID, chapterID string, req model.Chapter) (*model.Chapter, error) {
	ch, err := s.courses.FindChapter(courseID, chapterID)
	if err != nil {
		return nil, apperr.NotFound("章节不存在")
	}
	if req.Title != "" {
		ch.Title = req.Title
	}
	if req.Minutes > 0 {
		ch.Minutes = req.Minutes
	}
	if req.Content != "" {
		ch.Content = req.Content
	}
	if req.Status != "" {
		ch.Status = req.Status
	}
	ch.UpdatedAt = time.Now().UnixMilli()
	if err := s.courses.UpdateChapter(ch); err != nil {
		return nil, err
	}
	return ch, nil
}

func (s *CourseService) DeleteChapter(courseID, chapterID string) error {
	if _, err := s.courses.FindChapter(courseID, chapterID); err != nil {
		return apperr.NotFound("章节不存在")
	}
	return s.courses.DeleteChapter(courseID, chapterID)
}

func (s *CourseService) ListPublished() ([]model.Course, error) {
	courses, err := s.courses.ListPublished()
	if err != nil {
		return nil, err
	}
	if courses == nil {
		courses = []model.Course{}
	}
	for i := range courses {
		courses[i].ChapterCount = len(courses[i].Chapters)
		// 列表不返回章节正文
		for j := range courses[i].Chapters {
			courses[i].Chapters[j].Content = ""
		}
	}
	return courses, nil
}

func (s *CourseService) GetPublished(id string) (*model.Course, error) {
	course, err := s.courses.FindPublishedByID(id)
	if err != nil {
		return nil, apperr.NotFound("课程不存在或未发布")
	}
	return course, nil
}
