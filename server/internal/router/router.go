package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/handler"
	"github.com/cng1985/ai-learning-server/internal/middleware"
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
	rbacSvc *service.RBACService,
	settingsSvc *service.SettingsService,
	lc fx.Lifecycle,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), corsMiddleware())

	rbacSvc.SyncToMiddleware()

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
	authed := r.Group("/api/v1", middleware.Auth(jwt))
	{
		authed.GET("/auth/me", h.Auth.Me)
		authed.GET("/auth/permissions", h.Auth.Permissions)
		authed.POST("/auth/permissions/refresh", h.Auth.RefreshPermissions)

		// AI 服务
		ai := authed.Group("/ai", middleware.RequirePermission(rbac.PermAIChat))
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
			app.GET("/courses", middleware.RequirePermission(rbac.PermCourseRead), h.App.ListCourses)
			app.GET("/courses/:id", middleware.RequirePermission(rbac.PermCourseRead), h.App.GetCourse)
			app.GET("/quizzes/:id", middleware.RequirePermission(rbac.PermQuizRead), h.App.GetQuiz)

			// 客户咨询
			support := app.Group("/support", middleware.RequirePermission(rbac.PermCustomerChat))
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
	r.GET("/api/v1/ws/support", middleware.Auth(jwt), h.Customer.HandleWS)

	// 管理端接口
	admin := r.Group("/api/v1", middleware.Auth(jwt), middleware.RequireAdmin())
	{
		users := admin.Group("/users", middleware.RequirePermission(rbac.PermUserRead))
		{
			users.GET("", h.Users.List)
			users.GET("/:id", h.Users.Get)
			users.POST("", middleware.RequirePermission(rbac.PermUserCreate), h.Users.Create)
			users.PUT("/:id", middleware.RequirePermission(rbac.PermUserUpdate), h.Users.Update)
			users.DELETE("/:id", middleware.RequirePermission(rbac.PermUserDelete), h.Users.Delete)
		}

		courses := admin.Group("/courses", middleware.RequirePermission(rbac.PermCourseRead))
		{
			courses.GET("", h.Courses.List)
			courses.GET("/:id", h.Courses.Get)
			courses.POST("", middleware.RequirePermission(rbac.PermCourseWrite), h.Courses.Create)
			courses.PUT("/:id", middleware.RequirePermission(rbac.PermCourseWrite), h.Courses.Update)
			courses.DELETE("/:id", middleware.RequirePermission(rbac.PermCourseDelete), h.Courses.Delete)
			courses.POST("/:id/chapters", middleware.RequirePermission(rbac.PermCourseWrite), h.Courses.AddChapter)
			courses.PUT("/:id/chapters/:chapterId", middleware.RequirePermission(rbac.PermCourseWrite), h.Courses.UpdateChapter)
			courses.DELETE("/:id/chapters/:chapterId", middleware.RequirePermission(rbac.PermCourseDelete), h.Courses.DeleteChapter)
		}

		quizzes := admin.Group("/quizzes", middleware.RequirePermission(rbac.PermQuizRead))
		{
			quizzes.GET("", h.Quizzes.List)
			quizzes.GET("/:id", h.Quizzes.Get)
			quizzes.POST("", middleware.RequirePermission(rbac.PermQuizWrite), h.Quizzes.Create)
			quizzes.PUT("/:id", middleware.RequirePermission(rbac.PermQuizWrite), h.Quizzes.Update)
			quizzes.DELETE("/:id", middleware.RequirePermission(rbac.PermQuizDelete), h.Quizzes.Delete)
		}

		reviews := admin.Group("/reviews", middleware.RequirePermission(rbac.PermReviewRead))
		{
			reviews.GET("", h.Reviews.List)
			reviews.GET("/:id", h.Reviews.Get)
			reviews.POST("/:id/approve", middleware.RequirePermission(rbac.PermReviewApprove), h.Reviews.Approve)
			reviews.POST("/:id/reject", middleware.RequirePermission(rbac.PermReviewApprove), h.Reviews.Reject)
		}

		admin.GET("/dashboard/stats", middleware.RequirePermission(rbac.PermDashboard), h.Dashboard.Stats)

		// 权限管理
		roles := admin.Group("/roles", middleware.RequirePermission(rbac.PermRoleManage))
		{
			roles.GET("", h.RBAC.ListRoles)
			roles.PUT("/:role", h.RBAC.UpdateRole)
		}
		admin.GET("/permissions", middleware.RequirePermission(rbac.PermRoleManage), h.RBAC.ListPermissions)

		settings := admin.Group("/settings", middleware.RequirePermission(rbac.PermSettingsManage))
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
		aiModels := admin.Group("/ai-models", middleware.RequirePermission(rbac.PermAiModelRead))
		{
			aiModels.GET("/overview", h.AiModel.Overview)
			aiModels.GET("/resolve", h.AiModel.ResolveTest)
			aiModels.PUT("/default", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.SetDefault)

			aiModels.GET("/canonical-models", h.AiModel.ListCanonicalModels)
			aiModels.POST("/canonical-models", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateCanonicalModel)
			aiModels.PUT("/canonical-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.UpdateCanonicalModel)
			aiModels.DELETE("/canonical-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteCanonicalModel)

			aiModels.GET("/capabilities", h.AiModel.ListCapabilities)
			aiModels.POST("/capabilities", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateCapability)
			aiModels.PUT("/capabilities/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.UpdateCapability)
			aiModels.DELETE("/capabilities/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteCapability)

			aiModels.GET("/capability-models", h.AiModel.ListCapabilityModels)
			aiModels.POST("/capability-models", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateCapabilityModel)
			aiModels.DELETE("/capability-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteCapabilityModel)

			aiModels.GET("/providers", h.AiModel.ListProviders)
			aiModels.GET("/providers/:id", h.AiModel.GetProvider)
			aiModels.POST("/providers", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateProvider)
			aiModels.PUT("/providers/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.UpdateProvider)
			aiModels.DELETE("/providers/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteProvider)

			aiModels.GET("/provider-models", h.AiModel.ListProviderModels)
			aiModels.POST("/provider-models", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateProviderModel)
			aiModels.PUT("/provider-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.UpdateProviderModel)
			aiModels.DELETE("/provider-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteProviderModel)

			aiModels.GET("/virtual-models", h.AiModel.ListVirtualModels)
			aiModels.GET("/virtual-models/options", h.AiModel.ListVirtualModelOptions)
			aiModels.POST("/virtual-models", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateVirtualModel)
			aiModels.PUT("/virtual-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.UpdateVirtualModel)
			aiModels.DELETE("/virtual-models/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteVirtualModel)

			aiModels.GET("/virtual-model-mappings", h.AiModel.ListVirtualModelMappings)
			aiModels.POST("/virtual-model-mappings", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.CreateVirtualModelMapping)
			aiModels.PUT("/virtual-model-mappings/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.UpdateVirtualModelMapping)
			aiModels.DELETE("/virtual-model-mappings/:id", middleware.RequirePermission(rbac.PermAiModelManage), h.AiModel.DeleteVirtualModelMapping)
		}

		// 知识库管理
		knowledge := admin.Group("/knowledge", middleware.RequirePermission(rbac.PermKnowledgeRead))
		{
			knowledge.GET("/status", h.Knowledge.Status)
			knowledge.GET("/chunks", h.Knowledge.ListChunks)
			knowledge.GET("/search", h.Knowledge.Search)
			knowledge.POST("/reindex", middleware.RequirePermission(rbac.PermKnowledgeManage), h.Knowledge.Reindex)
		}

		// 客户咨询管理
		customers := admin.Group("/customers", middleware.RequirePermission(rbac.PermCustomerRead))
		{
			customers.GET("/stats", h.Customer.AdminStats)
			customers.GET("/tickets", h.Customer.AdminListTickets)
			customers.GET("/tickets/:id", h.Customer.AdminGetTicket)
			customers.GET("/tickets/:id/messages", h.Customer.AdminListMessages)
			customers.POST("/tickets/:id/reply", middleware.RequirePermission(rbac.PermCustomerReply), h.Customer.AdminReply)
			customers.PUT("/tickets/:id/status", middleware.RequirePermission(rbac.PermCustomerReply), h.Customer.AdminUpdateStatus)
		}

		// 单据管理
		documents := admin.Group("/documents", middleware.RequirePermission(rbac.PermDocumentRead))
		{
			documents.GET("", h.Document.List)
			documents.GET("/export", middleware.RequirePermission(rbac.PermDocumentExport), h.Document.Export)
			documents.GET("/import/template", middleware.RequirePermission(rbac.PermDocumentImport), h.Document.ExportTemplate)
			documents.POST("/import", middleware.RequirePermission(rbac.PermDocumentImport), h.Document.Import)
			documents.GET("/import/:taskId/progress", middleware.RequirePermission(rbac.PermDocumentImport), h.Document.ImportProgress)
			documents.GET("/:id", h.Document.Get)
			documents.POST("", middleware.RequirePermission(rbac.PermDocumentWrite), h.Document.Create)
			documents.PUT("/:id", middleware.RequirePermission(rbac.PermDocumentWrite), h.Document.Update)
			documents.DELETE("/:id", middleware.RequirePermission(rbac.PermDocumentDelete), h.Document.Delete)
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.Port)
			go func() {
				fmt.Printf("🚀 API 服务已启动: http://localhost:%s\n", cfg.Port)
				fmt.Println("   管理端: admin / admin123 | reviewer / review123 | operator / oper123")
				fmt.Println("   学员端: 支持 /auth/register /auth/login /auth/guest")
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
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
