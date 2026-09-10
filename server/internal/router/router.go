// Package router 负责路由注册与中间件挂载。
//
// 鉴权层级：
//  1. middleware.Auth          —— 校验 JWT 并从数据库刷新角色/状态
//  2. middleware.RequireAdminPortal —— 管理端接口要求后台角色
//  3. middleware.RequirePermission  —— 按权限码做细粒度控制
package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/handler"
	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewEngine(
	cfg *config.Config,
	jwt *authutil.JWTManager,
	h *handler.Handlers,
	users *repository.UserRepo,
	rbacSvc *service.RBACService,
	settingsSvc *service.SettingsService,
	lc fx.Lifecycle,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), corsMiddleware(cfg.CORSOrigins))

	if err := rbacSvc.SyncResolver(); err != nil {
		fmt.Printf("⚠️  加载角色权限失败，将使用默认权限: %v\n", err)
	}

	auth := middleware.Auth(jwt, users)
	perm := func(p string) gin.HandlerFunc {
		return middleware.RequirePermission(rbacSvc.Resolver(), p)
	}

	r.GET("/api/v1/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "healthy"}, "ok")
	})

	// 公开认证接口
	authPublic := r.Group("/api/v1/auth")
	{
		authPublic.POST("/login", h.Auth.Login)
		authPublic.POST("/register", h.Auth.Register)
		authPublic.POST("/guest", h.Auth.GuestLogin)
	}

	// 任意已登录用户
	authed := r.Group("/api/v1", auth)
	{
		authed.GET("/auth/me", h.Auth.Me)
		authed.GET("/auth/permissions", h.Auth.Permissions)
		authed.POST("/auth/permissions/refresh", h.Auth.RefreshPermissions)

		// AI 服务
		ai := authed.Group("/ai", perm(rbac.PermAIChat))
		{
			ai.GET("/config", h.AI.Config)
			ai.POST("/chat", h.AI.Chat)
			ai.POST("/chat/stream", h.AI.ChatStream)
			ai.POST("/career/interview", h.AI.CareerInterview)
			ai.POST("/career/recommend", h.AI.CareerRecommend)
			ai.POST("/goal/decompose", h.AI.GoalDecompose)
			ai.POST("/learning/suggest", h.AI.LearningSuggest)
		}

		// 学员端接口
		app := authed.Group("/app")
		{
			app.GET("/profile", h.App.Profile)
			app.GET("/courses", perm(rbac.PermCourseRead), h.App.ListCourses)
			app.GET("/courses/:id", perm(rbac.PermCourseRead), h.App.GetCourse)
			app.GET("/quizzes/:id", perm(rbac.PermQuizRead), h.App.GetQuiz)

			// 客户咨询
			support := app.Group("/support", perm(rbac.PermCustomerChat))
			{
				support.POST("/tickets", h.Customer.CreateTicket)
				support.GET("/tickets", h.Customer.ListMyTickets)
				support.GET("/tickets/:id", h.Customer.GetTicket)
				support.GET("/tickets/:id/messages", h.Customer.ListMessages)
				support.POST("/tickets/:id/messages", h.Customer.SendMessage)
			}
		}
	}

	// WebSocket 客户咨询（需登录）
	r.GET("/api/v1/ws/support", auth, h.Customer.HandleWS)

	// 管理端接口
	admin := r.Group("/api/v1", auth, middleware.RequireAdminPortal())
	{
		users := admin.Group("/users", perm(rbac.PermUserRead))
		{
			users.GET("", h.Users.List)
			users.GET("/:id", h.Users.Get)
			users.POST("", perm(rbac.PermUserCreate), h.Users.Create)
			users.PUT("/:id", perm(rbac.PermUserUpdate), h.Users.Update)
			users.DELETE("/:id", perm(rbac.PermUserDelete), h.Users.Delete)
		}

		courses := admin.Group("/courses", perm(rbac.PermCourseRead))
		{
			courses.GET("", h.Courses.List)
			courses.GET("/:id", h.Courses.Get)
			courses.POST("", perm(rbac.PermCourseWrite), h.Courses.Create)
			courses.PUT("/:id", perm(rbac.PermCourseWrite), h.Courses.Update)
			courses.DELETE("/:id", perm(rbac.PermCourseDelete), h.Courses.Delete)
			courses.POST("/:id/chapters", perm(rbac.PermCourseWrite), h.Courses.AddChapter)
			courses.PUT("/:id/chapters/:chapterId", perm(rbac.PermCourseWrite), h.Courses.UpdateChapter)
			courses.DELETE("/:id/chapters/:chapterId", perm(rbac.PermCourseDelete), h.Courses.DeleteChapter)
		}

		quizzes := admin.Group("/quizzes", perm(rbac.PermQuizRead))
		{
			quizzes.GET("", h.Quizzes.List)
			quizzes.GET("/:id", h.Quizzes.Get)
			quizzes.POST("", perm(rbac.PermQuizWrite), h.Quizzes.Create)
			quizzes.PUT("/:id", perm(rbac.PermQuizWrite), h.Quizzes.Update)
			quizzes.DELETE("/:id", perm(rbac.PermQuizDelete), h.Quizzes.Delete)
		}

		reviews := admin.Group("/reviews", perm(rbac.PermReviewRead))
		{
			reviews.GET("", h.Reviews.List)
			reviews.GET("/:id", h.Reviews.Get)
			reviews.POST("/:id/approve", perm(rbac.PermReviewApprove), h.Reviews.Approve)
			reviews.POST("/:id/reject", perm(rbac.PermReviewApprove), h.Reviews.Reject)
		}

		admin.GET("/dashboard/stats", perm(rbac.PermDashboard), h.Dashboard.Stats)

		// 权限管理
		roles := admin.Group("/roles", perm(rbac.PermRoleManage))
		{
			roles.GET("", h.RBAC.ListRoles)
			roles.PUT("/:role", h.RBAC.UpdateRole)
		}
		admin.GET("/permissions", perm(rbac.PermRoleManage), h.RBAC.ListPermissions)

		settings := admin.Group("/settings", perm(rbac.PermSettingsManage))
		{
			settings.GET("", h.Settings.Get)
			settings.GET("/resolve", h.Settings.Resolve)
			settings.PUT("/default-virtual-model", h.Settings.SetDefaultVirtualModel)
			settings.POST("/providers", h.Settings.SaveProvider)
			settings.PUT("/providers/:id", h.Settings.UpdateProvider)
			settings.POST("/quick-setup", h.Settings.QuickSetup)
			settings.POST("/knowledge/reindex", h.Settings.ReindexKnowledge)
			settings.PUT("", h.Settings.Update)
		}

		// AI 大模型分层配置
		aiModels := admin.Group("/ai-models", perm(rbac.PermAiModelRead))
		{
			manage := perm(rbac.PermAiModelManage)

			aiModels.GET("/overview", h.AiModel.Overview)
			aiModels.GET("/resolve", h.AiModel.ResolveTest)
			aiModels.PUT("/default", manage, h.AiModel.SetDefault)

			aiModels.GET("/canonical-models", h.AiModel.ListCanonicalModels)
			aiModels.POST("/canonical-models", manage, h.AiModel.CreateCanonicalModel)
			aiModels.PUT("/canonical-models/:id", manage, h.AiModel.UpdateCanonicalModel)
			aiModels.DELETE("/canonical-models/:id", manage, h.AiModel.DeleteCanonicalModel)

			aiModels.GET("/capabilities", h.AiModel.ListCapabilities)
			aiModels.POST("/capabilities", manage, h.AiModel.CreateCapability)
			aiModels.PUT("/capabilities/:id", manage, h.AiModel.UpdateCapability)
			aiModels.DELETE("/capabilities/:id", manage, h.AiModel.DeleteCapability)

			aiModels.GET("/capability-models", h.AiModel.ListCapabilityModels)
			aiModels.POST("/capability-models", manage, h.AiModel.CreateCapabilityModel)
			aiModels.DELETE("/capability-models/:id", manage, h.AiModel.DeleteCapabilityModel)

			aiModels.GET("/providers", h.AiModel.ListProviders)
			aiModels.GET("/providers/:id", h.AiModel.GetProvider)
			aiModels.POST("/providers", manage, h.AiModel.CreateProvider)
			aiModels.PUT("/providers/:id", manage, h.AiModel.UpdateProvider)
			aiModels.DELETE("/providers/:id", manage, h.AiModel.DeleteProvider)

			aiModels.GET("/provider-models", h.AiModel.ListProviderModels)
			aiModels.POST("/provider-models", manage, h.AiModel.CreateProviderModel)
			aiModels.PUT("/provider-models/:id", manage, h.AiModel.UpdateProviderModel)
			aiModels.DELETE("/provider-models/:id", manage, h.AiModel.DeleteProviderModel)

			aiModels.GET("/virtual-models", h.AiModel.ListVirtualModels)
			aiModels.GET("/virtual-models/options", h.AiModel.ListVirtualModelOptions)
			aiModels.POST("/virtual-models", manage, h.AiModel.CreateVirtualModel)
			aiModels.PUT("/virtual-models/:id", manage, h.AiModel.UpdateVirtualModel)
			aiModels.DELETE("/virtual-models/:id", manage, h.AiModel.DeleteVirtualModel)

			aiModels.GET("/virtual-model-mappings", h.AiModel.ListVirtualModelMappings)
			aiModels.POST("/virtual-model-mappings", manage, h.AiModel.CreateVirtualModelMapping)
			aiModels.PUT("/virtual-model-mappings/:id", manage, h.AiModel.UpdateVirtualModelMapping)
			aiModels.DELETE("/virtual-model-mappings/:id", manage, h.AiModel.DeleteVirtualModelMapping)
		}

		// 知识库管理
		knowledge := admin.Group("/knowledge", perm(rbac.PermKnowledgeRead))
		{
			knowledge.GET("/status", h.Knowledge.Status)
			knowledge.GET("/chunks", h.Knowledge.ListChunks)
			knowledge.GET("/search", h.Knowledge.Search)
			knowledge.POST("/reindex", perm(rbac.PermKnowledgeManage), h.Knowledge.Reindex)
		}

		// 客户咨询管理
		customers := admin.Group("/customers", perm(rbac.PermCustomerRead))
		{
			customers.GET("/stats", h.Customer.AdminStats)
			customers.GET("/tickets", h.Customer.AdminListTickets)
			customers.GET("/tickets/:id", h.Customer.AdminGetTicket)
			customers.GET("/tickets/:id/messages", h.Customer.AdminListMessages)
			customers.POST("/tickets/:id/reply", perm(rbac.PermCustomerReply), h.Customer.AdminReply)
			customers.PUT("/tickets/:id/status", perm(rbac.PermCustomerReply), h.Customer.AdminUpdateStatus)
		}

		// 单据管理
		documents := admin.Group("/documents", perm(rbac.PermDocumentRead))
		{
			documents.GET("", h.Document.List)
			documents.GET("/export", perm(rbac.PermDocumentExport), h.Document.Export)
			documents.GET("/import/template", perm(rbac.PermDocumentImport), h.Document.ExportTemplate)
			documents.POST("/import", perm(rbac.PermDocumentImport), h.Document.Import)
			documents.GET("/import/:taskId/progress", perm(rbac.PermDocumentImport), h.Document.ImportProgress)
			documents.GET("/:id", h.Document.Get)
			documents.POST("", perm(rbac.PermDocumentWrite), h.Document.Create)
			documents.PUT("/:id", perm(rbac.PermDocumentWrite), h.Document.Update)
			documents.DELETE("/:id", perm(rbac.PermDocumentDelete), h.Document.Delete)
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.Port)
			go func() {
				fmt.Printf("🚀 API 服务已启动: http://localhost:%s\n", cfg.Port)
				fmt.Println("   演示账号见 README.md「演示账号」章节")
				llmCfg := settingsSvc.LLMConfig()
				if llmCfg.Enabled {
					fmt.Printf("   AI 大模型: 已启用 (%s)\n", llmCfg.Model)
				} else {
					fmt.Println("   AI 大模型: 未就绪（请在管理端「大模型配置」中配置厂商 API Key）")
				}
				if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
	})
	return r
}

// corsMiddleware 按配置的来源白名单处理跨域；origins 为 * 时不限制。
func corsMiddleware(origins string) gin.HandlerFunc {
	allowAll := origins == "" || origins == "*"
	allowed := map[string]bool{}
	if !allowAll {
		for _, o := range strings.Split(origins, ",") {
			allowed[strings.TrimSpace(o)] = true
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		switch {
		case allowAll:
			c.Header("Access-Control-Allow-Origin", "*")
		case allowed[origin]:
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

var Module = fx.Provide(NewEngine)
