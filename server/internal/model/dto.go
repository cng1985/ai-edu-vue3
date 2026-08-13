package model

// 本文件集中存放认证与用户管理相关的请求 DTO，
// 避免 handler 直接绑定数据库实体。

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Portal   string `json:"portal"` // admin=管理后台登录，其他为学员端
}

type RegisterRequest struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	AvatarColor string `json:"avatarColor"`
}

// UserUpsertRequest 管理端创建/更新用户的请求体
type UserUpsertRequest struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	Avatar      string `json:"avatar"`
	AvatarColor string `json:"avatarColor"`
}

// ReviewDecisionRequest 审核通过/驳回请求体
type ReviewDecisionRequest struct {
	Comment string `json:"comment"`
}
