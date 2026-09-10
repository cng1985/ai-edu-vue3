package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type authFixture struct {
	auth  *AuthService
	users *UserService
	rbac  *RBACService
	repo  *repository.UserRepo
}

func newAuthFixture(t *testing.T) *authFixture {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.RolePermission{}); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}
	userRepo := repository.NewUserRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	rbacSvc := NewRBACService(roleRepo, rbac.NewResolver())
	jwtMgr := authutil.NewJWTManager(&config.Config{JWTSecret: "test-secret", TokenTTL: time.Hour})
	return &authFixture{
		auth:  NewAuthService(userRepo, jwtMgr, rbacSvc),
		users: NewUserService(userRepo),
		rbac:  rbacSvc,
		repo:  userRepo,
	}
}

func (f *authFixture) createUser(t *testing.T, username, password, role, status string) *model.User {
	t.Helper()
	hash, _ := authutil.HashPassword(password)
	user := &model.User{
		ID: "id_" + username, Username: username, Nickname: username,
		PasswordHash: hash, Role: role, Status: status,
		JoinedAt: time.Now().UnixMilli(),
	}
	if err := f.repo.Create(user); err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	return user
}

func assertStatus(t *testing.T, err error, want int) {
	t.Helper()
	if err == nil {
		t.Fatal("期望出错但成功了")
	}
	got, ok := apperr.HTTPStatus(err)
	if !ok {
		t.Fatalf("期望 apperr 类型化错误, got %T: %v", err, err)
	}
	if got != want {
		t.Fatalf("HTTP 状态码 = %d, want %d (err: %v)", got, want, err)
	}
}

func TestLoginSuccess(t *testing.T) {
	f := newAuthFixture(t)
	f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusActive)

	res, err := f.auth.Login(model.LoginRequest{Username: "alice", Password: "secret123"})
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if res.Token == "" {
		t.Error("应返回 token")
	}
	if res.User.Username != "alice" || res.User.RoleName != "学员" {
		t.Errorf("用户信息不正确: %+v", res.User)
	}
	if len(res.User.Permissions) == 0 {
		t.Error("应返回权限列表")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	f := newAuthFixture(t)
	f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusActive)
	_, err := f.auth.Login(model.LoginRequest{Username: "alice", Password: "wrong"})
	assertStatus(t, err, http.StatusUnauthorized)
}

func TestLoginUnknownUser(t *testing.T) {
	f := newAuthFixture(t)
	_, err := f.auth.Login(model.LoginRequest{Username: "nobody", Password: "x"})
	assertStatus(t, err, http.StatusUnauthorized)
}

func TestLoginDisabledUser(t *testing.T) {
	f := newAuthFixture(t)
	f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusDisabled)
	_, err := f.auth.Login(model.LoginRequest{Username: "alice", Password: "secret123"})
	assertStatus(t, err, http.StatusForbidden)
}

// 学员不能从管理端门户登录
func TestLoginAdminPortalRejectsLearner(t *testing.T) {
	f := newAuthFixture(t)
	f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusActive)
	_, err := f.auth.Login(model.LoginRequest{Username: "alice", Password: "secret123", Portal: "admin"})
	assertStatus(t, err, http.StatusForbidden)
}

func TestLoginAdminPortalAllowsOperator(t *testing.T) {
	f := newAuthFixture(t)
	f.createUser(t, "op", "secret123", rbac.RoleOperator, model.UserStatusActive)
	if _, err := f.auth.Login(model.LoginRequest{Username: "op", Password: "secret123", Portal: "admin"}); err != nil {
		t.Fatalf("运营应可登录管理端: %v", err)
	}
}

func TestRegisterValidation(t *testing.T) {
	f := newAuthFixture(t)
	cases := []model.RegisterRequest{
		{Username: "ab", Nickname: "n", Password: "secret123"},         // 用户名过短
		{Username: "alice", Nickname: "", Password: "secret123"},       // 缺昵称
		{Username: "alice", Nickname: "n", Password: "123"},            // 密码过短
		{Username: "含中文", Nickname: "n", Password: "secret123"},        // 非法用户名
	}
	for _, req := range cases {
		if _, err := f.auth.Register(req); err == nil {
			t.Errorf("请求 %+v 应校验失败", req)
		}
	}
}

func TestRegisterDuplicateUsername(t *testing.T) {
	f := newAuthFixture(t)
	f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusActive)
	_, err := f.auth.Register(model.RegisterRequest{Username: "Alice", Nickname: "小A", Password: "secret123"})
	if err == nil {
		t.Fatal("用户名（大小写不敏感）重复应失败")
	}
}

func TestRegisterAssignsLearnerRole(t *testing.T) {
	f := newAuthFixture(t)
	res, err := f.auth.Register(model.RegisterRequest{Username: "bob", Nickname: "小B", Password: "secret123"})
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	if res.User.Role != rbac.RoleLearner {
		t.Errorf("注册用户角色应为 learner, got %s", res.User.Role)
	}
}

func TestGuestLogin(t *testing.T) {
	f := newAuthFixture(t)
	res, err := f.auth.GuestLogin()
	if err != nil {
		t.Fatalf("游客登录失败: %v", err)
	}
	if res.User.Role != rbac.RoleGuest {
		t.Errorf("游客角色错误: %s", res.User.Role)
	}
	// 游客不落库
	if n, _ := f.repo.Total(); n != 0 {
		t.Errorf("游客不应写入数据库, 用户数 = %d", n)
	}
}

func TestUserCreateRejectsInvalidRole(t *testing.T) {
	f := newAuthFixture(t)
	_, err := f.users.Create(model.UserUpsertRequest{Username: "evil", Password: "secret123", Role: "superuser"})
	assertStatus(t, err, http.StatusBadRequest)
}

func TestUserUpdateRejectsInvalidRole(t *testing.T) {
	f := newAuthFixture(t)
	u := f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusActive)
	_, err := f.users.Update(u.ID, model.UserUpsertRequest{Role: "superuser"})
	assertStatus(t, err, http.StatusBadRequest)
}

func TestUserDeleteLastAdminGuard(t *testing.T) {
	f := newAuthFixture(t)
	admin := f.createUser(t, "admin", "secret123", rbac.RoleAdmin, model.UserStatusActive)
	if err := f.users.Delete(admin.ID); err == nil {
		t.Fatal("不应允许删除最后一个管理员")
	}
	f.createUser(t, "admin2", "secret123", rbac.RoleAdmin, model.UserStatusActive)
	if err := f.users.Delete(admin.ID); err != nil {
		t.Fatalf("存在其他管理员时应可删除: %v", err)
	}
}

func TestRBACUpdateRoleSyncsResolver(t *testing.T) {
	f := newAuthFixture(t)
	if !f.rbac.Resolver().Has(rbac.RoleReviewer, rbac.PermReviewApprove) {
		t.Fatal("reviewer 默认应有 review:approve")
	}
	if err := f.rbac.UpdateRole(rbac.RoleReviewer, []string{rbac.PermCourseRead}); err != nil {
		t.Fatalf("更新角色权限失败: %v", err)
	}
	if f.rbac.Resolver().Has(rbac.RoleReviewer, rbac.PermReviewApprove) {
		t.Error("更新后 Resolver 应立即失去 review:approve")
	}
	if !f.rbac.Resolver().Has(rbac.RoleReviewer, rbac.PermCourseRead) {
		t.Error("更新后 Resolver 应有 course:read")
	}
}

func TestRBACUpdateUnknownRole(t *testing.T) {
	f := newAuthFixture(t)
	err := f.rbac.UpdateRole("superuser", []string{rbac.PermCourseRead})
	if err == nil {
		t.Fatal("未知角色应更新失败")
	}
	var e *apperr.Error
	if !errors.As(err, &e) {
		t.Errorf("应返回类型化错误, got %T", err)
	}
}

func TestRefreshPermissionsUsesDBRole(t *testing.T) {
	f := newAuthFixture(t)
	u := f.createUser(t, "alice", "secret123", rbac.RoleLearner, model.UserStatusActive)
	u.Role = rbac.RoleReviewer
	if err := f.repo.Update(u); err != nil {
		t.Fatalf("更新角色失败: %v", err)
	}
	res, err := f.auth.RefreshPermissions(u.ID)
	if err != nil {
		t.Fatalf("刷新权限失败: %v", err)
	}
	if res.Role != rbac.RoleReviewer {
		t.Errorf("应返回数据库最新角色, got %s", res.Role)
	}
}
