package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/gin-gonic/gin"
)

type fakeUserStore struct {
	users map[string]*model.User
}

func (s *fakeUserStore) FindByID(id string) (*model.User, error) {
	if u, ok := s.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("record not found")
}

func newTestJWT() *authutil.JWTManager {
	return authutil.NewJWTManager(&config.Config{JWTSecret: "test-secret", TokenTTL: time.Hour})
}

func newRouter(jwt *authutil.JWTManager, store UserStore, extra ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handlers := append([]gin.HandlerFunc{Auth(jwt, store)}, extra...)
	handlers = append(handlers, func(c *gin.Context) {
		claims := GetClaims(c)
		c.JSON(http.StatusOK, gin.H{"id": claims.ID, "role": claims.Role})
	})
	r.GET("/protected", handlers...)
	return r
}

func doRequest(r *gin.Engine, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthMissingToken(t *testing.T) {
	r := newRouter(newTestJWT(), &fakeUserStore{})
	if w := doRequest(r, ""); w.Code != http.StatusUnauthorized {
		t.Errorf("无 token 应返回 401, got %d", w.Code)
	}
}

func TestAuthInvalidToken(t *testing.T) {
	r := newRouter(newTestJWT(), &fakeUserStore{})
	if w := doRequest(r, "not-a-jwt"); w.Code != http.StatusUnauthorized {
		t.Errorf("非法 token 应返回 401, got %d", w.Code)
	}
}

func TestAuthDeletedUser(t *testing.T) {
	jwt := newTestJWT()
	r := newRouter(jwt, &fakeUserStore{users: map[string]*model.User{}})
	token, _ := jwt.Sign(model.Claims{ID: "u_gone", Username: "ghost", Role: "admin"})
	if w := doRequest(r, token); w.Code != http.StatusUnauthorized {
		t.Errorf("已删除用户的 token 应返回 401, got %d", w.Code)
	}
}

func TestAuthDisabledUser(t *testing.T) {
	jwt := newTestJWT()
	store := &fakeUserStore{users: map[string]*model.User{
		"u_1": {ID: "u_1", Username: "alice", Role: "learner", Status: model.UserStatusDisabled},
	}}
	r := newRouter(jwt, store)
	token, _ := jwt.Sign(model.Claims{ID: "u_1", Username: "alice", Role: "learner"})
	if w := doRequest(r, token); w.Code != http.StatusForbidden {
		t.Errorf("禁用账号应返回 403, got %d", w.Code)
	}
}

// 角色变更后旧 token 立即按新角色鉴权（不再信任 JWT 内的 role）
func TestAuthRoleRefreshedFromDB(t *testing.T) {
	jwt := newTestJWT()
	store := &fakeUserStore{users: map[string]*model.User{
		"u_1": {ID: "u_1", Username: "alice", Role: "learner", Status: model.UserStatusActive},
	}}
	r := newRouter(jwt, store, RequireAdminPortal())
	// token 里写的是 admin，但数据库角色已是 learner
	token, _ := jwt.Sign(model.Claims{ID: "u_1", Username: "alice", Role: "admin"})
	if w := doRequest(r, token); w.Code != http.StatusForbidden {
		t.Errorf("数据库角色已降级时应返回 403, got %d", w.Code)
	}
}

func TestAuthGuestSkipsDBLookup(t *testing.T) {
	jwt := newTestJWT()
	r := newRouter(jwt, &fakeUserStore{users: map[string]*model.User{}})
	token, _ := jwt.Sign(model.Claims{ID: "guest_123", Username: "guest", Role: "guest"})
	if w := doRequest(r, token); w.Code != http.StatusOK {
		t.Errorf("游客 token 应跳过数据库校验, got %d", w.Code)
	}
}

func TestRequirePermission(t *testing.T) {
	jwt := newTestJWT()
	store := &fakeUserStore{users: map[string]*model.User{
		"u_1": {ID: "u_1", Username: "alice", Role: rbac.RoleReviewer, Status: model.UserStatusActive},
	}}
	resolver := rbac.NewResolver()
	token, _ := jwt.Sign(model.Claims{ID: "u_1", Username: "alice", Role: rbac.RoleReviewer})

	allowed := newRouter(jwt, store, RequirePermission(resolver, rbac.PermReviewApprove))
	if w := doRequest(allowed, token); w.Code != http.StatusOK {
		t.Errorf("reviewer 默认应有 review:approve, got %d", w.Code)
	}

	denied := newRouter(jwt, store, RequirePermission(resolver, rbac.PermUserDelete))
	if w := doRequest(denied, token); w.Code != http.StatusForbidden {
		t.Errorf("reviewer 不应有 user:delete, got %d", w.Code)
	}

	// 自定义权限收回后立即生效
	resolver.SetCustom(map[string][]string{rbac.RoleReviewer: {rbac.PermCourseRead}})
	if w := doRequest(allowed, token); w.Code != http.StatusForbidden {
		t.Errorf("自定义权限收回后应返回 403, got %d", w.Code)
	}
}
