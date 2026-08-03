package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cng1985/ai-learning-server/internal/middleware"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"gorm.io/datatypes"
)

type RBACService struct {
	roles *repository.RoleRepo
}

func NewRBACService(roles *repository.RoleRepo) *RBACService {
	return &RBACService{roles: roles}
}

var roleNames = map[string]string{
	"admin": "管理员", "reviewer": "审核员", "operator": "运营",
	"learner": "学员", "guest": "游客",
}

func (s *RBACService) ListPermissions() []rbac.PermissionInfo {
	return rbac.AllPermissions
}

func (s *RBACService) GetRoleName(role string) string {
	return roleNames[role]
}

func (s *RBACService) GetPermissions(role string) []string {
	custom, _ := s.loadCustomMap()
	perms := custom[role]
	if perms == nil {
		perms = rbac.DefaultRolePermissions[role]
	}
	if perms == nil {
		return []string{}
	}
	return perms
}

func (s *RBACService) EnrichUser(user *model.User) model.AuthUser {
	return model.AuthUser{
		User:        *user,
		Permissions: s.GetPermissions(user.Role),
		RoleName:    s.GetRoleName(user.Role),
	}
}

func (s *RBACService) ListRoles() ([]model.RoleInfo, error) {
	custom, _ := s.loadCustomMap()
	roles := []string{"admin", "reviewer", "operator", "learner", "guest"}
	var result []model.RoleInfo
	for _, role := range roles {
		perms := custom[role]
		if perms == nil {
			perms = rbac.DefaultRolePermissions[role]
		}
		result = append(result, model.RoleInfo{
			Role: role, Name: roleNames[role], Permissions: perms,
		})
	}
	return result, nil
}

func (s *RBACService) UpdateRole(role string, permissions []string) error {
	if roleNames[role] == "" {
		return errors.New("角色不存在")
	}
	b, _ := json.Marshal(permissions)
	rp := model.RolePermission{
		Role: role, Permissions: datatypes.JSON(b), UpdatedAt: time.Now().UnixMilli(),
	}
	if err := s.roles.Upsert(&rp); err != nil {
		return err
	}
	custom, _ := s.loadCustomMap()
	middleware.SetCustomPermissions(custom)
	return nil
}

func (s *RBACService) SyncToMiddleware() {
	custom, _ := s.loadCustomMap()
	middleware.SetCustomPermissions(custom)
}

func (s *RBACService) loadCustomMap() (map[string][]string, error) {
	list, err := s.roles.List()
	if err != nil {
		return nil, err
	}
	m := map[string][]string{}
	for _, rp := range list {
		var perms []string
		_ = json.Unmarshal(rp.Permissions, &perms)
		m[rp.Role] = perms
	}
	return m, nil
}
