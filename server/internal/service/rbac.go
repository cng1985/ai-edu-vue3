package service

import (
	"encoding/json"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"gorm.io/datatypes"
)

// RBACService 管理角色权限：数据库为持久化来源，Resolver 为运行时唯一读取入口。
// 启动时调用 SyncResolver 灌入数据，UpdateRole 修改后立即重新同步，
// 因此鉴权中间件与业务层读到的权限始终一致。
type RBACService struct {
	roles    *repository.RoleRepo
	resolver *rbac.Resolver
}

func NewRBACService(roles *repository.RoleRepo, resolver *rbac.Resolver) *RBACService {
	return &RBACService{roles: roles, resolver: resolver}
}

func (s *RBACService) Resolver() *rbac.Resolver {
	return s.resolver
}

func (s *RBACService) ListPermissions() []rbac.PermissionInfo {
	return rbac.AllPermissions
}

func (s *RBACService) GetRoleName(role string) string {
	return rbac.RoleNames[role]
}

func (s *RBACService) GetPermissions(role string) []string {
	return s.resolver.Permissions(role)
}

func (s *RBACService) EnrichUser(user *model.User) model.AuthUser {
	return model.AuthUser{
		User:        *user,
		Permissions: s.GetPermissions(user.Role),
		RoleName:    s.GetRoleName(user.Role),
	}
}

func (s *RBACService) ListRoles() ([]model.RoleInfo, error) {
	var result []model.RoleInfo
	for _, role := range rbac.Roles() {
		result = append(result, model.RoleInfo{
			Role:        role,
			Name:        rbac.RoleNames[role],
			Permissions: s.resolver.Permissions(role),
		})
	}
	return result, nil
}

func (s *RBACService) UpdateRole(role string, permissions []string) error {
	if !rbac.IsValidRole(role) {
		return apperr.NotFound("角色不存在")
	}
	b, _ := json.Marshal(permissions)
	rp := model.RolePermission{
		Role: role, Permissions: datatypes.JSON(b), UpdatedAt: time.Now().UnixMilli(),
	}
	if err := s.roles.Upsert(&rp); err != nil {
		return err
	}
	return s.SyncResolver()
}

// SyncResolver 从数据库加载自定义角色权限并同步到 Resolver。
func (s *RBACService) SyncResolver() error {
	list, err := s.roles.List()
	if err != nil {
		return err
	}
	custom := map[string][]string{}
	for _, rp := range list {
		var perms []string
		_ = json.Unmarshal(rp.Permissions, &perms)
		custom[rp.Role] = perms
	}
	s.resolver.SetCustom(custom)
	return nil
}
