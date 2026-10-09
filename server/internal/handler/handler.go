// Package handler 是 HTTP 接口层，按业务域拆分文件：
// 基础域 auth / user / course / quiz / review / dashboard / rbac / ai / app / settings / customer / knowledge / ai_model；
// 成长生态 growth / practice / community / kernel / catalog_admin。
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
	Auth         *AuthHandler
	Users        *UserHandler
	Courses      *CourseHandler
	Quizzes      *QuizHandler
	Reviews      *ReviewHandler
	Dashboard    *DashboardHandler
	AI           *AIHandler
	RBAC         *RBACHandler
	App          *AppHandler
	Settings     *SettingsHandler
	Customer     *CustomerHandler
	Knowledge    *KnowledgeHandler
	AiModel      *AiModelHandler
	Growth       *GrowthHandler
	Practice     *PracticeHandler
	Community    *CommunityHandler
	Kernel       *KernelHandler
	CatalogAdmin *CatalogAdminHandler
}

type handlerParams struct {
	fx.In
	Auth         *AuthHandler
	Users        *UserHandler
	Courses      *CourseHandler
	Quizzes      *QuizHandler
	Reviews      *ReviewHandler
	Dashboard    *DashboardHandler
	AI           *AIHandler
	RBAC         *RBACHandler
	App          *AppHandler
	Settings     *SettingsHandler
	Customer     *CustomerHandler
	Knowledge    *KnowledgeHandler
	AiModel      *AiModelHandler
	Growth       *GrowthHandler
	Practice     *PracticeHandler
	Community    *CommunityHandler
	Kernel       *KernelHandler
	CatalogAdmin *CatalogAdminHandler
}

func NewHandlers(p handlerParams) *Handlers {
	return &Handlers{
		Auth: p.Auth, Users: p.Users, Courses: p.Courses, Quizzes: p.Quizzes, Reviews: p.Reviews,
		Dashboard: p.Dashboard, AI: p.AI, RBAC: p.RBAC, App: p.App, Settings: p.Settings,
		Customer: p.Customer, Knowledge: p.Knowledge, AiModel: p.AiModel,
		Growth: p.Growth, Practice: p.Practice, Community: p.Community, Kernel: p.Kernel, CatalogAdmin: p.CatalogAdmin,
	}
}

var Module = fx.Provide(
	NewAuthHandler, NewUserHandler, NewCourseHandler,
	NewQuizHandler, NewReviewHandler, NewDashboardHandler,
	NewAIHandler, NewRBACHandler, NewAppHandler, NewSettingsHandler, NewCustomerHandler,
	NewKnowledgeHandler, NewAiModelHandler,
	NewGrowthHandler, NewPracticeHandler, NewCommunityHandler, NewKernelHandler, NewCatalogAdminHandler,
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
	case "统一模型不存在", "能力标签不存在", "关联不存在", "厂商不存在", "厂商模型不存在", "虚拟模型不存在", "映射不存在":
		code = http.StatusNotFound
	}
	response.Fail(c, code, code, err.Error())
}
