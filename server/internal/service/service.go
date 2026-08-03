package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"go.uber.org/fx"
	"gorm.io/datatypes"
)

type AuthService struct {
	users *repository.UserRepo
	jwt   *authutil.JWTManager
	rbac  *RBACService
}

type UserService struct{ users *repository.UserRepo }
type CourseService struct {
	courses *repository.CourseRepo
}
type QuizService struct{ quizzes *repository.QuizRepo }
type ReviewService struct {
	reviews *repository.ReviewRepo
	courses *repository.CourseRepo
	quizzes *repository.QuizRepo
}
type DashboardService struct {
	users   *repository.UserRepo
	courses *repository.CourseRepo
	quizzes *repository.QuizRepo
	reviews *repository.ReviewRepo
}

func NewAuthService(users *repository.UserRepo, jwt *authutil.JWTManager, rbac *RBACService) *AuthService {
	return &AuthService{users: users, jwt: jwt, rbac: rbac}
}
func NewUserService(users *repository.UserRepo) *UserService { return &UserService{users: users} }
func NewCourseService(courses *repository.CourseRepo) *CourseService {
	return &CourseService{courses: courses}
}
func NewQuizService(quizzes *repository.QuizRepo) *QuizService { return &QuizService{quizzes: quizzes} }
func NewReviewService(reviews *repository.ReviewRepo, courses *repository.CourseRepo, quizzes *repository.QuizRepo) *ReviewService {
	return &ReviewService{reviews: reviews, courses: courses, quizzes: quizzes}
}
func NewDashboardService(users *repository.UserRepo, courses *repository.CourseRepo, quizzes *repository.QuizRepo, reviews *repository.ReviewRepo) *DashboardService {
	return &DashboardService{users: users, courses: courses, quizzes: quizzes, reviews: reviews}
}

var Module = fx.Provide(
	authutil.NewJWTManager,
	NewAuthService,
	NewUserService,
	NewCourseService,
	NewQuizService,
	NewReviewService,
	NewDashboardService,
	NewAIService,
	NewRBACService,
)

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d_%04d", prefix, time.Now().UnixMilli(), rand.Intn(10000))
}

func countQuestions(q datatypes.JSON) int {
	var questions []model.Question
	if err := json.Unmarshal(q, &questions); err != nil {
		return 0
	}
	return len(questions)
}

// --- Auth ---

func (s *AuthService) Login(username, password, portal string) (*model.LoginResponse, error) {
	user, err := s.users.FindByUsername(strings.TrimSpace(username))
	if err != nil || !authutil.VerifyPassword(password, user.PasswordHash) {
		return nil, errors.New("用户名或密码错误")
	}
	if user.Status == "disabled" {
		return nil, errors.New("账号已被禁用")
	}
	if portal == "admin" && !rbac.IsAdminRole(user.Role) {
		return nil, errors.New("无权限访问管理后台")
	}
	token, err := s.jwt.Sign(model.Claims{ID: user.ID, Username: user.Username, Role: user.Role})
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: s.rbac.EnrichUser(user)}, nil
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	nickname := strings.TrimSpace(req.Nickname)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`).MatchString(username) {
		return nil, errors.New("用户名需为 3~20 位字母、数字、下划线或短横线")
	}
	if nickname == "" {
		return nil, errors.New("请填写昵称")
	}
	if len(nickname) > 16 {
		return nil, errors.New("昵称最长 16 个字符")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("密码至少 6 位")
	}
	if _, err := s.users.FindByUsername(username); err == nil {
		return nil, errors.New("该用户名已被注册")
	}
	hash, err := authutil.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	avatar := req.Avatar
	if avatar == "" {
		avatar = strings.ToUpper(username[:1])
	}
	color := req.AvatarColor
	if color == "" {
		color = "#6366f1"
	}
	user := model.User{
		ID: genID("u"), Username: username, Nickname: nickname,
		PasswordHash: hash, Role: "learner", Status: "active",
		Avatar: avatar, AvatarColor: color, JoinedAt: time.Now().UnixMilli(),
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}
	token, err := s.jwt.Sign(model.Claims{ID: user.ID, Username: user.Username, Role: user.Role})
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: s.rbac.EnrichUser(&user)}, nil
}

func (s *AuthService) GuestLogin() (*model.LoginResponse, error) {
	stamp := time.Now().UnixMilli()
	id := fmt.Sprintf("guest_%d", stamp)
	username := fmt.Sprintf("guest_%s", genID("g")[2:6])
	user := model.User{
		ID: id, Username: username, Nickname: "游客",
		Role: "guest", Status: "active",
		Avatar: "访", AvatarColor: "#94a3b8", JoinedAt: stamp,
	}
	token, err := s.jwt.Sign(model.Claims{ID: user.ID, Username: user.Username, Role: "guest"})
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: s.rbac.EnrichUser(&user)}, nil
}

func (s *AuthService) Me(id string) (*model.AuthUser, error) {
	if strings.HasPrefix(id, "guest_") {
		guest := &model.User{
			ID: id, Username: "guest", Nickname: "游客",
			Role: "guest", Status: "active",
			Avatar: "访", AvatarColor: "#94a3b8",
		}
		authUser := s.rbac.EnrichUser(guest)
		return &authUser, nil
	}
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	authUser := s.rbac.EnrichUser(user)
	return &authUser, nil
}

func (s *AuthService) Permissions(role string) []string {
	return s.rbac.GetPermissions(role)
}

func (s *AuthService) RoleName(role string) string {
	return s.rbac.GetRoleName(role)
}

// --- User ---

func (s *UserService) List(keyword, role, status string, page, pageSize int) (*model.PageResult[model.User], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	users, total, err := s.users.List(keyword, role, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = []model.User{}
	}
	return &model.PageResult[model.User]{List: users, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *UserService) Get(id string) (*model.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (s *UserService) Create(req model.User, password string) (*model.User, error) {
	if req.Username == "" || password == "" {
		return nil, errors.New("用户名和密码必填")
	}
	if _, err := s.users.FindByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}
	hash, err := authutil.HashPassword(password)
	if err != nil {
		return nil, err
	}
	if req.Role == "" {
		req.Role = "learner"
	}
	if req.Status == "" {
		req.Status = "active"
	}
	if req.Avatar == "" {
		req.Avatar = strings.ToUpper(req.Username[:1])
	}
	if req.AvatarColor == "" {
		req.AvatarColor = "#6366f1"
	}
	if req.Nickname == "" {
		req.Nickname = req.Username
	}
	user := model.User{
		ID: genID("u"), Username: strings.TrimSpace(req.Username), Nickname: req.Nickname,
		PasswordHash: hash, Role: req.Role, Status: req.Status,
		Avatar: req.Avatar, AvatarColor: req.AvatarColor, JoinedAt: time.Now().UnixMilli(),
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(id string, req model.User, password string) (*model.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.AvatarColor != "" {
		user.AvatarColor = req.AvatarColor
	}
	if password != "" {
		hash, err := authutil.HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}
	if err := s.users.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(id string) error {
	user, err := s.users.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user.Role == "admin" {
		n, _ := s.users.CountByRole("admin")
		if n <= 1 {
			return errors.New("不能删除最后一个管理员")
		}
	}
	return s.users.Delete(id)
}

// --- Course ---

func (s *CourseService) List(keyword, status string, page, pageSize int) (*model.PageResult[model.Course], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
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
		return nil, errors.New("课程不存在")
	}
	return course, nil
}

func (s *CourseService) Create(course model.Course) (*model.Course, error) {
	if course.Title == "" {
		return nil, errors.New("课程标题必填")
	}
	if course.ID == "" {
		course.ID = genID("course")
	}
	if s.courses.Exists(course.ID) {
		return nil, errors.New("课程 ID 已存在")
	}
	now := time.Now().UnixMilli()
	if course.Level == "" {
		course.Level = "入门"
	}
	if course.Icon == "" {
		course.Icon = "📚"
	}
	if course.Accent == "" {
		course.Accent = "#6366f1"
	}
	if course.EstimatedMinutes == 0 {
		course.EstimatedMinutes = 60
	}
	if course.Status == "" {
		course.Status = "draft"
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
		return nil, errors.New("课程不存在")
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
		return errors.New("课程不存在")
	}
	return s.courses.Delete(id)
}

func (s *CourseService) AddChapter(courseID string, ch model.Chapter) (*model.Chapter, error) {
	if _, err := s.courses.FindByID(courseID); err != nil {
		return nil, errors.New("课程不存在")
	}
	if ch.Title == "" {
		return nil, errors.New("章节标题必填")
	}
	if ch.ID == "" {
		ch.ID = genID("ch")
	}
	if s.courses.ChapterExists(courseID, ch.ID) {
		return nil, errors.New("章节 ID 已存在")
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
		return nil, errors.New("章节不存在")
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
		return errors.New("章节不存在")
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
		return nil, errors.New("课程不存在或未发布")
	}
	return course, nil
}

// --- Quiz ---

func (s *QuizService) List(keyword, courseID, status string, page, pageSize int) (*model.PageResult[model.Quiz], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
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
		return nil, errors.New("测验不存在")
	}
	quiz.QuestionCount = countQuestions(quiz.Questions)
	return quiz, nil
}

func (s *QuizService) Create(quiz model.Quiz) (*model.Quiz, error) {
	if quiz.Title == "" || quiz.CourseID == "" {
		return nil, errors.New("标题和关联课程必填")
	}
	if quiz.ID == "" {
		quiz.ID = genID("quiz")
	}
	if s.quizzes.Exists(quiz.ID) {
		return nil, errors.New("测验 ID 已存在")
	}
	now := time.Now().UnixMilli()
	if quiz.Status == "" {
		quiz.Status = "draft"
	}
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
		return nil, errors.New("测验不存在")
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
		return errors.New("测验不存在")
	}
	return s.quizzes.Delete(id)
}

// --- Review ---

func (s *ReviewService) List(status, reviewType string, page, pageSize int) (*model.PageResult[model.Review], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
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
		return nil, errors.New("审核记录不存在")
	}
	return review, nil
}

func (s *ReviewService) Approve(id, reviewerID, comment string) (*model.Review, error) {
	review, err := s.reviews.FindByID(id)
	if err != nil {
		return nil, errors.New("审核记录不存在")
	}
	if review.Status != "pending" {
		return nil, errors.New("该记录已处理")
	}
	review.Status = "approved"
	review.ReviewerID = reviewerID
	review.ReviewedAt = time.Now().UnixMilli()
	review.Comment = comment

	if review.Type == "chapter" {
		ch, err := s.courses.FindChapter(review.CourseID, review.TargetID)
		if err == nil {
			ch.Content = review.Content
			ch.Status = "published"
			ch.UpdatedAt = time.Now().UnixMilli()
			_ = s.courses.UpdateChapter(ch)
		}
	} else if review.Type == "quiz" {
		var question model.Question
		if err := json.Unmarshal([]byte(review.Content), &question); err == nil && question.Text != "" {
			quiz, err := s.quizzes.FindByID(review.TargetID)
			if err == nil {
				var questions []model.Question
				_ = json.Unmarshal(quiz.Questions, &questions)
				questions = append(questions, question)
				b, _ := json.Marshal(questions)
				quiz.Questions = datatypes.JSON(b)
				quiz.UpdatedAt = time.Now().UnixMilli()
				_ = s.quizzes.Update(quiz)
			}
		}
	}

	if err := s.reviews.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) Reject(id, reviewerID, comment string) (*model.Review, error) {
	review, err := s.reviews.FindByID(id)
	if err != nil {
		return nil, errors.New("审核记录不存在")
	}
	if review.Status != "pending" {
		return nil, errors.New("该记录已处理")
	}
	review.Status = "rejected"
	review.ReviewerID = reviewerID
	review.ReviewedAt = time.Now().UnixMilli()
	if comment == "" {
		comment = "内容不符合发布标准"
	}
	review.Comment = comment
	if err := s.reviews.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

// --- Dashboard ---

func (s *DashboardService) Stats() (*model.DashboardStats, error) {
	totalUsers, _ := s.users.Total()
	learners, _ := s.users.CountByRole("learner")
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
