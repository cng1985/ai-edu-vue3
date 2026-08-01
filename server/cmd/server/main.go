package main

import (
	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/database"
	"github.com/cng1985/ai-learning-server/internal/handler"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/internal/router"
	"github.com/cng1985/ai-learning-server/internal/seed"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		database.Module,
		repository.Module,
		service.Module,
		handler.Module,
		router.Module,
		fx.Invoke(seed.Run),
		fx.Invoke(func(_ *gin.Engine) {}),
	).Run()
}
