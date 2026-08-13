// Package service 是业务逻辑层，按业务域拆分文件：
// auth.go / user.go / course.go / quiz.go / review.go / dashboard.go / rbac.go 等。
package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"go.uber.org/fx"
	"gorm.io/datatypes"
)

var Module = fx.Provide(
	rbac.NewResolver,
	authutil.NewJWTManager,
	NewAuthService,
	NewUserService,
	NewCourseService,
	NewQuizService,
	NewReviewService,
	NewDashboardService,
	NewSettingsService,
	NewAIService,
	NewKnowledgeService,
	NewRBACService,
	NewCustomerService,
	NewDocumentService,
	NewModelRouter,
	NewAiModelService,
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
