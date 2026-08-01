package rbac

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
	PermAIChat      = "ai:chat"
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
	{Code: PermAIChat, Name: "AI 对话", Group: "AI 服务"},
}

type PermissionInfo struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

// DefaultRolePermissions 默认角色权限映射
var DefaultRolePermissions = map[string][]string{
	"admin": {
		PermUserRead, PermUserCreate, PermUserUpdate, PermUserDelete,
		PermCourseRead, PermCourseWrite, PermCourseDelete,
		PermQuizRead, PermQuizWrite, PermQuizDelete,
		PermReviewRead, PermReviewApprove,
		PermDashboard, PermRoleManage, PermAIChat,
	},
	"reviewer": {
		PermCourseRead, PermQuizRead,
		PermReviewRead, PermReviewApprove,
		PermDashboard, PermAIChat,
	},
	"operator": {
		PermCourseRead, PermCourseWrite,
		PermQuizRead, PermQuizWrite,
		PermDashboard, PermAIChat,
	},
	"learner": {
		PermCourseRead, PermQuizRead, PermAIChat,
	},
	"guest": {
		PermAIChat,
	},
}

func HasPermission(role, perm string, custom map[string][]string) bool {
	perms := custom[role]
	if perms == nil {
		perms = DefaultRolePermissions[role]
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

func IsAdminRole(role string) bool {
	return role == "admin" || role == "reviewer" || role == "operator"
}

func IsAppRole(role string) bool {
	return role == "learner" || role == "guest"
}
