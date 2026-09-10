package rbac

import "sync"

// Resolver 是运行时的权限解析器，作为鉴权中间件与业务层的唯一权限数据源。
// 自定义角色权限（来自数据库 role_permissions 表）通过 SetCustom 同步进来，
// 未自定义的角色回退到 DefaultRolePermissions。
type Resolver struct {
	mu     sync.RWMutex
	custom map[string][]string
}

func NewResolver() *Resolver {
	return &Resolver{custom: map[string][]string{}}
}

// SetCustom 整体替换自定义角色权限表（写入方为 RBAC 服务）。
func (r *Resolver) SetCustom(perms map[string][]string) {
	copied := make(map[string][]string, len(perms))
	for role, list := range perms {
		copied[role] = append([]string(nil), list...)
	}
	r.mu.Lock()
	r.custom = copied
	r.mu.Unlock()
}

// Permissions 返回角色的有效权限列表（副本，调用方可安全修改）。
func (r *Resolver) Permissions(role string) []string {
	r.mu.RLock()
	perms, ok := r.custom[role]
	r.mu.RUnlock()
	if !ok {
		perms = DefaultRolePermissions[role]
	}
	return append([]string{}, perms...)
}

// Has 判断角色是否拥有指定权限。
func (r *Resolver) Has(role, perm string) bool {
	r.mu.RLock()
	perms, ok := r.custom[role]
	r.mu.RUnlock()
	if !ok {
		perms = DefaultRolePermissions[role]
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
