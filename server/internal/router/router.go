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

		// AI 服务
		ai := authed.Group("/ai", middleware.RequirePermission(rbac.PermAIChat))
		{
			ai.GET("/config", h.AI.Config)
			ai.POST("/chat", h.AI.Chat)
			ai.POST("/chat/stream", h.AI.ChatStream)
		}

		// 学员端接口
		app := authed.Group("/app")
		{
			app.GET("/profile", h.App.Profile)
			app.GET("/courses", middleware.RequirePermission(rbac.PermCourseRead), h.App.ListCourses)
			app.GET("/courses/:id", middleware.RequirePermission(rbac.PermCourseRead), h.App.GetCourse)
			app.GET("/quizzes/:id", middleware.RequirePermission(rbac.PermQuizRead), h.App.GetQuiz)
		}
	}

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
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.Port)
			go func() {
				fmt.Printf("🚀 API 服务已启动: http://localhost:%s\n", cfg.Port)
				fmt.Println("   管理端: admin / admin123 | reviewer / review123 | operator / oper123")
				fmt.Println("   学员端: 支持 /auth/register /auth/login /auth/guest")
				if cfg.LLM.Enabled {
					fmt.Printf("   AI 大模型: 已启用 (%s)\n", cfg.LLM.Model)
				} else {
					fmt.Println("   AI 大模型: 本地模式（设置 LLM_API_KEY 启用）")
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
