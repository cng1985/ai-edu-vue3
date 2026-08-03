package handler

import (
	"net/http"
	"strconv"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Handlers struct {
	Auth      *AuthHandler
	Users     *UserHandler
	Courses   *CourseHandler
	Quizzes   *QuizHandler
	Reviews   *ReviewHandler
	Dashboard *DashboardHandler
	AI        *AIHandler
	RBAC      *RBACHandler
	App       *AppHandler
	Settings  *SettingsHandler
	Customer  *CustomerHandler
	Document  *DocumentHandler
	Knowledge *KnowledgeHandler
}

type AuthHandler struct{ svc *service.AuthService }
type UserHandler struct{ svc *service.UserService }
type CourseHandler struct{ svc *service.CourseService }
type QuizHandler struct{ svc *service.QuizService }
type ReviewHandler struct{ svc *service.ReviewService }
type DashboardHandler struct{ svc *service.DashboardService }

func NewHandlers(
	auth *AuthHandler, users *UserHandler, courses *CourseHandler,
	quizzes *QuizHandler, reviews *ReviewHandler, dashboard *DashboardHandler,
	ai *AIHandler, rbac *RBACHandler, app *AppHandler, settings *SettingsHandler,
	customer *CustomerHandler, document *DocumentHandler, knowledge *KnowledgeHandler,
) *Handlers {
	return &Handlers{
		Auth: auth, Users: users, Courses: courses, Quizzes: quizzes,
		Reviews: reviews, Dashboard: dashboard, AI: ai, RBAC: rbac, App: app,
		Settings: settings, Customer: customer, Document: document, Knowledge: knowledge,
	}
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler       { return &AuthHandler{svc: svc} }
func NewUserHandler(svc *service.UserService) *UserHandler       { return &UserHandler{svc: svc} }
func NewCourseHandler(svc *service.CourseService) *CourseHandler { return &CourseHandler{svc: svc} }
func NewQuizHandler(svc *service.QuizService) *QuizHandler       { return &QuizHandler{svc: svc} }
func NewReviewHandler(svc *service.ReviewService) *ReviewHandler { return &ReviewHandler{svc: svc} }
func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

var Module = fx.Provide(
	NewAuthHandler, NewUserHandler, NewCourseHandler,
	NewQuizHandler, NewReviewHandler, NewDashboardHandler,
	NewAIHandler, NewRBACHandler, NewAppHandler, NewSettingsHandler, NewCustomerHandler,
	NewDocumentHandler, NewKnowledgeHandler,
	NewHandlers,
)

func pageQuery(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	return page, pageSize
}

func failErr(c *gin.Context, err error) {
	code := http.StatusBadRequest
	switch err.Error() {
	case "用户名或密码错误":
		code = http.StatusUnauthorized
	case "账号已被禁用", "无权限", "无权限访问管理后台":
		code = http.StatusForbidden
	case "该用户名已被注册":
		code = http.StatusBadRequest
	case "用户不存在", "课程不存在", "课程不存在或未发布", "章节不存在", "测验不存在", "审核记录不存在", "角色不存在", "工单不存在", "单据不存在", "导入任务不存在或已过期":
		code = http.StatusNotFound
	}
	response.Fail(c, code, code, err.Error())
}

// --- Auth ---

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Portal   string `json:"portal"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, 400, "请输入用户名和密码")
		return
	}
	res, err := h.svc.Login(req.Username, req.Password, req.Portal)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	user, err := h.svc.Me(claims.ID)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user)
}

func (h *AuthHandler) Permissions(c *gin.Context) {
	claims := middleware.GetClaims(c)
	response.OK(c, gin.H{
		"role":        claims.Role,
		"roleName":    h.svc.RoleName(claims.Role),
		"permissions": h.svc.Permissions(claims.Role),
	})
}

func (h *AuthHandler) RefreshPermissions(c *gin.Context) {
	claims := middleware.GetClaims(c)
	user, err := h.svc.RefreshPermissions(claims.ID)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user, "权限已刷新")
}

// --- Users ---

func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("role"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *UserHandler) Get(c *gin.Context) {
	user, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		model.User
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	user, err := h.svc.Create(req.User, req.Password)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user, "创建成功")
}

func (h *UserHandler) Update(c *gin.Context) {
	var req struct {
		model.User
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	user, err := h.svc.Update(c.Param("id"), req.User, req.Password)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, user, "更新成功")
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- Courses ---

func (h *CourseHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *CourseHandler) Get(c *gin.Context) {
	course, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, course)
}

func (h *CourseHandler) Create(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Create(course)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *CourseHandler) Update(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Update(c.Param("id"), course)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *CourseHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

func (h *CourseHandler) AddChapter(c *gin.Context) {
	var ch model.Chapter
	if err := c.ShouldBindJSON(&ch); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.AddChapter(c.Param("id"), ch)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "章节创建成功")
}

func (h *CourseHandler) UpdateChapter(c *gin.Context) {
	var ch model.Chapter
	if err := c.ShouldBindJSON(&ch); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.UpdateChapter(c.Param("id"), c.Param("chapterId"), ch)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "章节更新成功")
}

func (h *CourseHandler) DeleteChapter(c *gin.Context) {
	if err := h.svc.DeleteChapter(c.Param("id"), c.Param("chapterId")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "章节删除成功")
}

// --- Quizzes ---

func (h *QuizHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("keyword"), c.Query("courseId"), c.Query("status"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *QuizHandler) Get(c *gin.Context) {
	quiz, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, quiz)
}

func (h *QuizHandler) Create(c *gin.Context) {
	var quiz model.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Create(quiz)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "创建成功")
}

func (h *QuizHandler) Update(c *gin.Context) {
	var quiz model.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "参数错误")
		return
	}
	res, err := h.svc.Update(c.Param("id"), quiz)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res, "更新成功")
}

func (h *QuizHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, nil, "删除成功")
}

// --- Reviews ---

func (h *ReviewHandler) List(c *gin.Context) {
	page, pageSize := pageQuery(c)
	res, err := h.svc.List(c.Query("status"), c.Query("type"), page, pageSize)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *ReviewHandler) Get(c *gin.Context) {
	review, err := h.svc.Get(c.Param("id"))
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, review)
}

func (h *ReviewHandler) Approve(c *gin.Context) {
	var req struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&req)
	claims := middleware.GetClaims(c)
	review, err := h.svc.Approve(c.Param("id"), claims.ID, req.Comment)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, review, "审核通过")
}

func (h *ReviewHandler) Reject(c *gin.Context) {
	var req struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&req)
	claims := middleware.GetClaims(c)
	review, err := h.svc.Reject(c.Param("id"), claims.ID, req.Comment)
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, review, "已驳回")
}

// --- Dashboard ---

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		failErr(c, err)
		return
	}
	response.OK(c, stats)
}
