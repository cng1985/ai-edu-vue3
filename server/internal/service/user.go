package service

import (
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
)

// UserService 负责管理端的用户 CRUD。
type UserService struct{ users *repository.UserRepo }

func NewUserService(users *repository.UserRepo) *UserService { return &UserService{users: users} }

func (s *UserService) List(keyword, role, status string, page, pageSize int) (*model.PageResult[model.User], error) {
	page, pageSize = normalizePage(page, pageSize)
	users, total, err := s.users.List(keyword, role, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = []model.User{}
	}
	return &model.PageResult[model.User]{List: users, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *UserService) Get(id string) (*model.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("用户不存在")
	}
	return user, nil
}

func (s *UserService) Create(req model.UserUpsertRequest) (*model.User, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return nil, apperr.BadRequest("用户名和密码必填")
	}
	if req.Role != "" && !rbac.IsValidRole(req.Role) {
		return nil, apperr.BadRequest("非法角色")
	}
	if _, err := s.users.FindByUsername(username); err == nil {
		return nil, apperr.BadRequest("用户名已存在")
	}
	hash, err := authutil.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := model.User{
		ID:           genID("u"),
		Username:     username,
		Nickname:     defaultStr(req.Nickname, username),
		PasswordHash: hash,
		Role:         defaultStr(req.Role, rbac.RoleLearner),
		Status:       defaultStr(req.Status, model.UserStatusActive),
		Avatar:       defaultStr(req.Avatar, strings.ToUpper(username[:1])),
		AvatarColor:  defaultStr(req.AvatarColor, "#6366f1"),
		JoinedAt:     time.Now().UnixMilli(),
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(id string, req model.UserUpsertRequest) (*model.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("用户不存在")
	}
	if req.Role != "" && !rbac.IsValidRole(req.Role) {
		return nil, apperr.BadRequest("非法角色")
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.AvatarColor != "" {
		user.AvatarColor = req.AvatarColor
	}
	if req.Password != "" {
		hash, err := authutil.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}
	if err := s.users.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(id string) error {
	user, err := s.users.FindByID(id)
	if err != nil {
		return apperr.NotFound("用户不存在")
	}
	if user.Role == rbac.RoleAdmin {
		n, _ := s.users.CountByRole(rbac.RoleAdmin)
		if n <= 1 {
			return apperr.BadRequest("不能删除最后一个管理员")
		}
	}
	return s.users.Delete(id)
}

func defaultStr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return page, pageSize
}
