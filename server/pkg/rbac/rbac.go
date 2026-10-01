// Package rbac 是角色与权限的唯一后端数据源：
// 角色定义、权限码、默认角色权限映射均在此维护。
package rbac

// 角色常量。生态中的四类参与者：学习者、创作者、企业、平台运营（管理员/审核员/运营）。
const (
	RoleAdmin      = "admin"
	RoleReviewer   = "reviewer"
	RoleOperator   = "operator"
	RoleLearner    = "learner"
	RoleCreator    = "creator"
	RoleEnterprise = "enterprise"
	RoleGuest      = "guest"
)

// RoleNames 角色显示名
var RoleNames = map[string]string{
	RoleAdmin:      "管理员",
	RoleReviewer:   "审核员",
	RoleOperator:   "运营",
	RoleLearner:    "学习者",
	RoleCreator:    "创作者",
	RoleEnterprise: "企业",
	RoleGuest:      "游客",
}

// Roles 返回全部角色（固定顺序，供列表展示）
func Roles() []string {
	return []string{RoleAdmin, RoleReviewer, RoleOperator, RoleLearner, RoleCreator, RoleEnterprise, RoleGuest}
}

// IsValidRole 判断角色是否为系统已定义角色
func IsValidRole(role string) bool {
	_, ok := RoleNames[role]
	return ok
}

// AdminPortalRoles 允许登录管理后台的角色
var AdminPortalRoles = []string{RoleAdmin, RoleReviewer, RoleOperator}

// TalentRoles 会形成人才画像、参与任务匹配的角色
var TalentRoles = []string{RoleLearner, RoleCreator}

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
	PermUserRead       = "user:read"
	PermUserCreate     = "user:create"
	PermUserUpdate     = "user:update"
	PermUserDelete     = "user:delete"
	PermCourseRead     = "course:read"
	PermCourseWrite    = "course:write"
	PermCourseDelete   = "course:delete"
	PermQuizRead       = "quiz:read"
	PermQuizWrite      = "quiz:write"
	PermQuizDelete     = "quiz:delete"
	PermReviewRead     = "review:read"
	PermReviewApprove  = "review:approve"
	PermDashboard      = "dashboard:read"
	PermRoleManage     = "role:manage"
	PermSettingsManage = "settings:manage"
	PermAIChat         = "ai:chat"
	PermCustomerChat   = "customer:chat"
	PermCustomerRead   = "customer:read"
	PermCustomerReply  = "customer:reply"
	PermKnowledgeRead  = "knowledge:read"
	PermKnowledgeManage = "knowledge:manage"
	PermAiModelRead    = "ai_model:read"
	PermAiModelManage  = "ai_model:manage"

	// 成长生态
	PermEcoRead            = "eco:read"            // 浏览职业、技能、知识图谱、项目、任务市场
	PermGrowthWrite        = "growth:write"        // 设定目标、记录学习、提交项目、申请任务
	PermCommunityPost      = "community:post"      // 社区发帖、回答
	PermResourcePublish    = "resource:publish"    // 发布知识资产（文章、Prompt、SOP 等）
	PermOpportunityPublish = "opportunity:publish" // 发布 IT 任务并管理自己发布的任务
	PermTalentRead         = "talent:read"         // 查看人才库与人才画像
	PermEcoManage          = "eco:manage"          // 管理职业/技能/知识图谱/项目
	PermOpportunityManage  = "opportunity:manage"  // 管理全部任务与履约
	PermCommunityManage    = "community:manage"    // 社区与知识资产治理
	PermKernelRead         = "kernel:read"         // 查看 AI 内核运行记录
)

var AllPermissions = []PermissionInfo{
	{Code: PermUserRead, Name: "查看用户", Group: "用户管理"},
	{Code: PermUserCreate, Name: "创建用户", Group: "用户管理"},
	{Code: PermUserUpdate, Name: "编辑用户", Group: "用户管理"},
	{Code: PermUserDelete, Name: "删除用户", Group: "用户管理"},
	{Code: PermEcoRead, Name: "浏览成长生态", Group: "成长生态"},
	{Code: PermGrowthWrite, Name: "记录个人成长", Group: "成长生态"},
	{Code: PermEcoManage, Name: "管理能力与知识体系", Group: "成长生态"},
	{Code: PermCourseRead, Name: "查看课程", Group: "教育体系"},
	{Code: PermCourseWrite, Name: "编辑课程", Group: "教育体系"},
	{Code: PermCourseDelete, Name: "删除课程", Group: "教育体系"},
	{Code: PermQuizRead, Name: "查看题库", Group: "教育体系"},
	{Code: PermQuizWrite, Name: "编辑题库", Group: "教育体系"},
	{Code: PermQuizDelete, Name: "删除题库", Group: "教育体系"},
	{Code: PermOpportunityPublish, Name: "发布 IT 任务", Group: "任务市场"},
	{Code: PermOpportunityManage, Name: "管理任务与履约", Group: "任务市场"},
	{Code: PermTalentRead, Name: "查看人才库", Group: "任务市场"},
	{Code: PermCommunityPost, Name: "社区发帖与回答", Group: "知识社区"},
	{Code: PermResourcePublish, Name: "发布知识资产", Group: "知识社区"},
	{Code: PermCommunityManage, Name: "社区治理", Group: "知识社区"},
	{Code: PermReviewRead, Name: "查看审核", Group: "内容审核"},
	{Code: PermReviewApprove, Name: "审核操作", Group: "内容审核"},
	{Code: PermDashboard, Name: "数据看板", Group: "运营管理"},
	{Code: PermRoleManage, Name: "权限管理", Group: "系统管理"},
	{Code: PermSettingsManage, Name: "系统设置", Group: "系统管理"},
	{Code: PermAIChat, Name: "AI 对话与 Agent", Group: "AI 服务"},
	{Code: PermKernelRead, Name: "查看 AI 内核运行", Group: "AI 服务"},
	{Code: PermAiModelRead, Name: "查看 AI 模型配置", Group: "AI 服务"},
	{Code: PermAiModelManage, Name: "管理 AI 模型配置", Group: "AI 服务"},
	{Code: PermCustomerChat, Name: "客户咨询", Group: "客户服务"},
	{Code: PermCustomerRead, Name: "查看客户咨询", Group: "客户服务"},
	{Code: PermCustomerReply, Name: "回复客户咨询", Group: "客户服务"},
	{Code: PermKnowledgeRead, Name: "查看 RAG 知识库", Group: "知识库"},
	{Code: PermKnowledgeManage, Name: "管理 RAG 知识库", Group: "知识库"},
}

type PermissionInfo struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

// learnerPerms 学习者在生态中的基础权限。
var learnerPerms = []string{
	PermEcoRead, PermGrowthWrite, PermCommunityPost,
	PermCourseRead, PermQuizRead, PermAIChat, PermCustomerChat,
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
		PermKnowledgeRead, PermKnowledgeManage,
		PermAiModelRead, PermAiModelManage,
		PermEcoRead, PermEcoManage, PermOpportunityPublish, PermOpportunityManage,
		PermTalentRead, PermCommunityPost, PermResourcePublish, PermCommunityManage, PermKernelRead,
	},
	RoleReviewer: {
		PermCourseRead, PermQuizRead,
		PermReviewRead, PermReviewApprove,
		PermDashboard, PermAIChat,
		PermKnowledgeRead,
		PermEcoRead, PermCommunityManage,
	},
	RoleOperator: {
		PermCourseRead, PermCourseWrite,
		PermQuizRead, PermQuizWrite,
		PermDashboard, PermAIChat,
		PermCustomerRead, PermCustomerReply,
		PermKnowledgeRead,
		PermEcoRead, PermEcoManage, PermOpportunityManage, PermTalentRead, PermCommunityManage, PermKernelRead,
	},
	RoleLearner: learnerPerms,
	RoleCreator: append(append([]string{}, learnerPerms...), PermResourcePublish),
	RoleEnterprise: {
		PermEcoRead, PermCommunityPost, PermAIChat, PermCustomerChat,
		PermOpportunityPublish, PermTalentRead,
	},
	RoleGuest: {
		PermEcoRead, PermAIChat,
	},
}
