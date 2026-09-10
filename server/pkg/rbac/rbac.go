// Package rbac 是角色与权限的唯一后端数据源：
// 角色定义、权限码、默认角色权限映射均在此维护。
package rbac

// 角色常量
const (
	RoleAdmin    = "admin"
	RoleReviewer = "reviewer"
	RoleOperator = "operator"
	RoleLearner  = "learner"
	RoleGuest    = "guest"
)

// RoleNames 角色显示名
var RoleNames = map[string]string{
	RoleAdmin:    "管理员",
	RoleReviewer: "审核员",
	RoleOperator: "运营",
	RoleLearner:  "学员",
	RoleGuest:    "游客",
}

// Roles 返回全部角色（固定顺序，供列表展示）
func Roles() []string {
	return []string{RoleAdmin, RoleReviewer, RoleOperator, RoleLearner, RoleGuest}
}

// IsValidRole 判断角色是否为系统已定义角色
func IsValidRole(role string) bool {
	_, ok := RoleNames[role]
	return ok
}

// AdminPortalRoles 允许登录管理后台的角色
var AdminPortalRoles = []string{RoleAdmin, RoleReviewer, RoleOperator}

// IsAdminRole 判断角色是否可访问管理后台
func IsAdminRole(role string) bool {
	for _, r := range AdminPortalRoles {
		if r == role {
			return true
		}
	}
	return false
}

// 权限常量
const (
	PermUserRead    = "user:read"
	PermUserCreate  = "user:create"
	PermUserUpdate  = "user:update"
	PermUserDelete  = "user:delete"
	PermCourseRead  = "course:read"
	PermCourseWrite = "course:write"
	PermCourseDelete = "course:delete"
	PermQuizRead    = "quiz:read"
	PermQuizWrite   = "quiz:write"
	PermQuizDelete  = "quiz:delete"
	PermReviewRead  = "review:read"
	PermReviewApprove = "review:approve"
	PermDashboard   = "dashboard:read"
	PermRoleManage  = "role:manage"
	PermSettingsManage = "settings:manage"
	PermAIChat      = "ai:chat"
	PermCustomerChat  = "customer:chat"
	PermCustomerRead  = "customer:read"
	PermCustomerReply = "customer:reply"
	PermDocumentRead   = "document:read"
	PermDocumentWrite  = "document:write"
	PermDocumentDelete = "document:delete"
	PermDocumentImport = "document:import"
	PermDocumentExport = "document:export"
	PermKnowledgeRead  = "knowledge:read"
	PermKnowledgeManage = "knowledge:manage"
	PermAiModelRead    = "ai_model:read"
	PermAiModelManage  = "ai_model:manage"
)

var AllPermissions = []PermissionInfo{
	{Code: PermUserRead, Name: "查看用户", Group: "用户管理"},
	{Code: PermUserCreate, Name: "创建用户", Group: "用户管理"},
	{Code: PermUserUpdate, Name: "编辑用户", Group: "用户管理"},
	{Code: PermUserDelete, Name: "删除用户", Group: "用户管理"},
	{Code: PermCourseRead, Name: "查看课程", Group: "内容管理"},
	{Code: PermCourseWrite, Name: "编辑课程", Group: "内容管理"},
	{Code: PermCourseDelete, Name: "删除课程", Group: "内容管理"},
	{Code: PermQuizRead, Name: "查看题库", Group: "内容管理"},
	{Code: PermQuizWrite, Name: "编辑题库", Group: "内容管理"},
	{Code: PermQuizDelete, Name: "删除题库", Group: "内容管理"},
	{Code: PermReviewRead, Name: "查看审核", Group: "内容审核"},
	{Code: PermReviewApprove, Name: "审核操作", Group: "内容审核"},
	{Code: PermDashboard, Name: "数据看板", Group: "运营管理"},
	{Code: PermRoleManage, Name: "权限管理", Group: "系统管理"},
	{Code: PermSettingsManage, Name: "系统设置", Group: "系统管理"},
	{Code: PermAIChat, Name: "AI 对话", Group: "AI 服务"},
	{Code: PermCustomerChat, Name: "客户咨询", Group: "客户服务"},
	{Code: PermCustomerRead, Name: "查看客户咨询", Group: "客户服务"},
	{Code: PermCustomerReply, Name: "回复客户咨询", Group: "客户服务"},
	{Code: PermDocumentRead, Name: "查看单据", Group: "单据管理"},
	{Code: PermDocumentWrite, Name: "编辑单据", Group: "单据管理"},
	{Code: PermDocumentDelete, Name: "删除单据", Group: "单据管理"},
	{Code: PermDocumentImport, Name: "导入单据", Group: "单据管理"},
	{Code: PermDocumentExport, Name: "导出单据", Group: "单据管理"},
	{Code: PermKnowledgeRead, Name: "查看知识库", Group: "知识库"},
	{Code: PermKnowledgeManage, Name: "管理知识库", Group: "知识库"},
	{Code: PermAiModelRead, Name: "查看 AI 模型配置", Group: "AI 服务"},
	{Code: PermAiModelManage, Name: "管理 AI 模型配置", Group: "AI 服务"},
}

type PermissionInfo struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

// DefaultRolePermissions 默认角色权限映射
var DefaultRolePermissions = map[string][]string{
	RoleAdmin: {
		PermUserRead, PermUserCreate, PermUserUpdate, PermUserDelete,
		PermCourseRead, PermCourseWrite, PermCourseDelete,
		PermQuizRead, PermQuizWrite, PermQuizDelete,
		PermReviewRead, PermReviewApprove,
		PermDashboard, PermRoleManage, PermSettingsManage, PermAIChat,
		PermCustomerRead, PermCustomerReply,
		PermDocumentRead, PermDocumentWrite, PermDocumentDelete, PermDocumentImport, PermDocumentExport,
		PermKnowledgeRead, PermKnowledgeManage,
		PermAiModelRead, PermAiModelManage,
	},
	RoleReviewer: {
		PermCourseRead, PermQuizRead,
		PermReviewRead, PermReviewApprove,
		PermDashboard, PermAIChat,
		PermKnowledgeRead,
	},
	RoleOperator: {
		PermCourseRead, PermCourseWrite,
		PermQuizRead, PermQuizWrite,
		PermDashboard, PermAIChat,
		PermCustomerRead, PermCustomerReply,
		PermDocumentRead, PermDocumentWrite, PermDocumentImport, PermDocumentExport,
		PermKnowledgeRead,
	},
	RoleLearner: {
		PermCourseRead, PermQuizRead, PermAIChat, PermCustomerChat,
	},
	RoleGuest: {
		PermAIChat,
	},
}
