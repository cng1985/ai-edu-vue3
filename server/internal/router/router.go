package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/handler"
	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewEngine(
	cfg *config.Config,
	jwt *authutil.JWTManager,
	h *handler.Handlers,
	lc fx.Lifecycle,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), corsMiddleware())

	r.GET("/api/v1/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "healthy"}, "ok")
	})
	r.POST("/api/v1/auth/login", h.Auth.Login)

	auth := r.Group("/api/v1", middleware.Auth(jwt), middleware.RequireAdmin())
	{
		auth.GET("/auth/me", h.Auth.Me)

		users := auth.Group("/users")
		{
			users.GET("", h.Users.List)
			users.GET("/:id", h.Users.Get)
			users.POST("", middleware.RequireRole("admin"), h.Users.Create)
			users.PUT("/:id", middleware.RequireRole("admin"), h.Users.Update)
			users.DELETE("/:id", middleware.RequireRole("admin"), h.Users.Delete)
		}

		courses := auth.Group("/courses")
		{
			courses.GET("", h.Courses.List)
			courses.GET("/:id", h.Courses.Get)
			courses.POST("", middleware.RequireRole("admin", "operator"), h.Courses.Create)
			courses.PUT("/:id", middleware.RequireRole("admin", "operator"), h.Courses.Update)
			courses.DELETE("/:id", middleware.RequireRole("admin"), h.Courses.Delete)
			courses.POST("/:id/chapters", middleware.RequireRole("admin", "operator"), h.Courses.AddChapter)
			courses.PUT("/:id/chapters/:chapterId", middleware.RequireRole("admin", "operator"), h.Courses.UpdateChapter)
			courses.DELETE("/:id/chapters/:chapterId", middleware.RequireRole("admin"), h.Courses.DeleteChapter)
		}

		quizzes := auth.Group("/quizzes")
		{
			quizzes.GET("", h.Quizzes.List)
			quizzes.GET("/:id", h.Quizzes.Get)
			quizzes.POST("", middleware.RequireRole("admin", "operator"), h.Quizzes.Create)
			quizzes.PUT("/:id", middleware.RequireRole("admin", "operator"), h.Quizzes.Update)
			quizzes.DELETE("/:id", middleware.RequireRole("admin"), h.Quizzes.Delete)
		}

		reviews := auth.Group("/reviews")
		{
			reviews.GET("", h.Reviews.List)
			reviews.GET("/:id", h.Reviews.Get)
			reviews.POST("/:id/approve", middleware.RequireRole("admin", "reviewer"), h.Reviews.Approve)
			reviews.POST("/:id/reject", middleware.RequireRole("admin", "reviewer"), h.Reviews.Reject)
		}

		auth.GET("/dashboard/stats", h.Dashboard.Stats)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.Port)
			go func() {
				fmt.Printf("🚀 API 服务已启动: http://localhost:%s\n", cfg.Port)
				fmt.Println("   管理端登录: admin / admin123")
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
