// Package service 是业务逻辑层，按业务域拆分文件：
// 基础域 auth / user / course / quiz / review / rbac / settings / ai_model / knowledge；
// 成长生态 catalog / growth / practice / market / community / talent / kernel / agent。
package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync/atomic"
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
	NewModelRouter,
	NewAiModelService,
	NewCatalogService,
	NewGrowthService,
	NewProjectService,
	NewMarketService,
	NewCommunityService,
	NewKernelService,
	NewAgentService,
	NewTalentService,
)

var idSeq = atomic.Uint64{}

func init() { idSeq.Store(uint64(rand.Intn(1000))) }

func genID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixMilli(), idSeq.Add(1))
}

func countQuestions(q datatypes.JSON) int {
	var questions []model.Question
	if err := json.Unmarshal(q, &questions); err != nil {
		return 0
	}
	return len(questions)
}

func toJSON(v interface{}) datatypes.JSON {
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON("null")
	}
	return datatypes.JSON(b)
}

func parseStrings(raw datatypes.JSON) []string {
	var out []string
	_ = json.Unmarshal(raw, &out)
	return nonNil(out)
}

func parseRequirements(raw datatypes.JSON) []model.SkillRequirement {
	var out []model.SkillRequirement
	_ = json.Unmarshal(raw, &out)
	return nonNil(out)
}

func userBrief(u model.User) model.UserBrief {
	return model.UserBrief{ID: u.ID, Nickname: u.Nickname, Avatar: u.Avatar, AvatarColor: u.AvatarColor, Role: u.Role}
}
