// Package handler 是 HTTP 接口层，按业务域拆分文件：
// auth.go / user.go / course.go / quiz.go / review.go / dashboard.go / rbac.go / ai.go / app.go 等。
package handler

import (
	"net/http"
	"strconv"

	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Handlers 聚合全部接口处理器，供路由统一注册。
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
	AiModel   *AiModelHandler
}

func NewHandlers(
	auth *AuthHandler, users *UserHandler, courses *CourseHandler,
	quizzes *QuizHandler, reviews *ReviewHandler, dashboard *DashboardHandler,
	ai *AIHandler, rbac *RBACHandler, app *AppHandler, settings *SettingsHandler,
	customer *CustomerHandler, document *DocumentHandler, knowledge *KnowledgeHandler,
	aiModel *AiModelHandler,
) *Handlers {
	return &Handlers{
		Auth: auth, Users: users, Courses: courses, Quizzes: quizzes,
		Reviews: reviews, Dashboard: dashboard, AI: ai, RBAC: rbac, App: app,
		Settings: settings, Customer: customer, Document: document, Knowledge: knowledge,
		AiModel: aiModel,
	}
}

var Module = fx.Provide(
	NewAuthHandler, NewUserHandler, NewCourseHandler,
	NewQuizHandler, NewReviewHandler, NewDashboardHandler,
	NewAIHandler, NewRBACHandler, NewAppHandler, NewSettingsHandler, NewCustomerHandler,
	NewDocumentHandler, NewKnowledgeHandler,
	NewAiModelHandler,
	NewHandlers,
)

func pageQuery(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	return page, pageSize
}

// failErr 将业务错误映射为 HTTP 响应。
// 优先识别 apperr 类型化错误；对尚未迁移到 apperr 的旧域保留消息匹配兜底。
func failErr(c *gin.Context, err error) {
	if code, ok := apperr.HTTPStatus(err); ok {
		response.Fail(c, code, code, err.Error())
		return
	}
	code := http.StatusBadRequest
	switch err.Error() {
	case "单据不存在", "导入任务不存在或已过期",
		"统一模型不存在", "能力标签不存在", "关联不存在", "厂商不存在", "厂商模型不存在", "虚拟模型不存在", "映射不存在":
		code = http.StatusNotFound
	}
	response.Fail(c, code, code, err.Error())
}
