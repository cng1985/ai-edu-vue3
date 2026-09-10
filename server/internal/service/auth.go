package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

// AuthService 负责登录、注册、游客会话与当前用户信息。
type AuthService struct {
	users *repository.UserRepo
	jwt   *authutil.JWTManager
	rbac  *RBACService
}

func NewAuthService(users *repository.UserRepo, jwt *authutil.JWTManager, rbac *RBACService) *AuthService {
	return &AuthService{users: users, jwt: jwt, rbac: rbac}
}

func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.users.FindByUsername(strings.TrimSpace(req.Username))
	if err != nil || !authutil.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, apperr.Unauthorized("用户名或密码错误")
	}
	if user.Status == model.UserStatusDisabled {
		return nil, apperr.Forbidden("账号已被禁用")
	}
	if req.Portal == "admin" && !rbac.IsAdminRole(user.Role) {
		return nil, apperr.Forbidden("无权限访问管理后台")
	}
	return s.issueSession(user)
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	nickname := strings.TrimSpace(req.Nickname)
	if !usernamePattern.MatchString(username) {
		return nil, apperr.BadRequest("用户名需为 3~20 位字母、数字、下划线或短横线")
	}
	if nickname == "" {
		return nil, apperr.BadRequest("请填写昵称")
	}
	if len(nickname) > 16 {
		return nil, apperr.BadRequest("昵称最长 16 个字符")
	}
	if len(req.Password) < 6 {
		return nil, apperr.BadRequest("密码至少 6 位")
	}
	if _, err := s.users.FindByUsername(username); err == nil {
		return nil, apperr.BadRequest("该用户名已被注册")
	}
	hash, err := authutil.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	avatar := req.Avatar
	if avatar == "" {
		avatar = strings.ToUpper(username[:1])
	}
	color := req.AvatarColor
	if color == "" {
		color = "#6366f1"
	}
	user := model.User{
		ID: genID("u"), Username: username, Nickname: nickname,
		PasswordHash: hash, Role: rbac.RoleLearner, Status: model.UserStatusActive,
		Avatar: avatar, AvatarColor: color, JoinedAt: time.Now().UnixMilli(),
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}
	return s.issueSession(&user)
}

// GuestLogin 签发不落库的游客会话（ID 前缀 guest_，中间件据此跳过数据库校验）。
func (s *AuthService) GuestLogin() (*model.LoginResponse, error) {
	stamp := time.Now().UnixMilli()
	user := guestUser(fmt.Sprintf("guest_%d", stamp))
	user.Username = fmt.Sprintf("guest_%s", genID("g")[2:6])
	user.JoinedAt = stamp
	return s.issueSession(user)
}

func (s *AuthService) Me(id string) (*model.AuthUser, error) {
	if strings.HasPrefix(id, "guest_") {
		authUser := s.rbac.EnrichUser(guestUser(id))
		return &authUser, nil
	}
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("用户不存在")
	}
	authUser := s.rbac.EnrichUser(user)
	return &authUser, nil
}

func (s *AuthService) Permissions(role string) []string {
	return s.rbac.GetPermissions(role)
}

func (s *AuthService) RoleName(role string) string {
	return s.rbac.GetRoleName(role)
}

// RefreshPermissions 从数据库读取最新角色并解析权限（不依赖 JWT 中的 role）
func (s *AuthService) RefreshPermissions(id string) (*model.AuthUser, error) {
	if strings.HasPrefix(id, "guest_") {
		return s.Me(id)
	}
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, apperr.NotFound("用户不存在")
	}
	if user.Status == model.UserStatusDisabled {
		return nil, apperr.Forbidden("账号已被禁用")
	}
	authUser := s.rbac.EnrichUser(user)
	return &authUser, nil
}

func (s *AuthService) issueSession(user *model.User) (*model.LoginResponse, error) {
	token, err := s.jwt.Sign(model.Claims{ID: user.ID, Username: user.Username, Role: user.Role})
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: s.rbac.EnrichUser(user)}, nil
}

func guestUser(id string) *model.User {
	return &model.User{
		ID: id, Username: "guest", Nickname: "游客",
		Role: rbac.RoleGuest, Status: model.UserStatusActive,
		Avatar: "访", AvatarColor: "#94a3b8",
	}
}
