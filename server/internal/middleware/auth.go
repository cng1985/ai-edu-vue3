package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

const ClaimsKey = "claims"

var (
	permMu      sync.RWMutex
	customPerms map[string][]string
)

func SetCustomPermissions(perms map[string][]string) {
	permMu.Lock()
	defer permMu.Unlock()
	customPerms = perms
}

func getCustomPerms() map[string][]string {
	permMu.RLock()
	defer permMu.RUnlock()
	return customPerms
}

func Auth(jwt *authutil.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, 401, "未登录")
			c.Abort()
			return
		}
		claims, err := jwt.Verify(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, 401, "登录已过期")
			c.Abort()
			return
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func OptionalAuth(jwt *authutil.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			if claims, err := jwt.Verify(strings.TrimPrefix(header, "Bearer ")); err == nil {
				c.Set(ClaimsKey, claims)
			}
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet(ClaimsKey).(*model.Claims)
		if !rbac.IsAdminRole(claims.Role) {
			response.Fail(c, http.StatusForbidden, 403, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		claims := c.MustGet(ClaimsKey).(*model.Claims)
		if !allowed[claims.Role] {
			response.Fail(c, http.StatusForbidden, 403, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet(ClaimsKey).(*model.Claims)
		if !rbac.HasPermission(claims.Role, perm, getCustomPerms()) {
			response.Fail(c, http.StatusForbidden, 403, "无权限: "+perm)
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetClaims(c *gin.Context) *model.Claims {
	return c.MustGet(ClaimsKey).(*model.Claims)
}

func HasClaims(c *gin.Context) bool {
	_, ok := c.Get(ClaimsKey)
	return ok
}
