// Package repository 是数据访问层（DAO），按业务域拆分文件：
// user.go / course.go / quiz.go / review.go / role.go 及其他领域文件。
package repository

import "go.uber.org/fx"

var Module = fx.Provide(
	NewUserRepo,
	NewCourseRepo,
	NewQuizRepo,
	NewReviewRepo,
	NewRoleRepo,
	NewSettingsRepo,
	NewCustomerRepo,
	NewDocumentRepo,
	NewKnowledgeRepo,
	NewAiModelRepo,
)
