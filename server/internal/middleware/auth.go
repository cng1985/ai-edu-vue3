// Package middleware 提供认证与鉴权中间件。
//
// 鉴权策略：JWT 仅证明"是谁"，角色与账号状态以数据库实时数据为准，
// 因此修改角色或禁用账号后立即生效，无需等待旧 token 过期。
package middleware

import (
	"net/http"
	"strings"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/cng1985/ai-learning-server/pkg/response"
	"github.com/gin-gonic/gin"
)

const ClaimsKey = "claims"

// GuestIDPrefix 游客用户 ID 前缀，游客不落库，跳过数据库校验。
const GuestIDPrefix = "guest_"

// UserStore 提供按 ID 查询用户的能力（由 repository.UserRepo 实现）。
type UserStore interface {
	FindByID(id string) (*model.User, error)
}

// Auth 校验 Bearer JWT，并从数据库加载用户最新的角色与状态写入 Claims。
func Auth(jwt *authutil.JWTManager, users UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			abort(c, http.StatusUnauthorized, "未登录")
			return
		}
		claims, err := jwt.Verify(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			abort(c, http.StatusUnauthorized, "登录已过期")
			return
		}
		if !strings.HasPrefix(claims.ID, GuestIDPrefix) {
			user, err := users.FindByID(claims.ID)
			if err != nil {
				abort(c, http.StatusUnauthorized, "账号不存在或已被删除")
				return
			}
			if user.Status == model.UserStatusDisabled {
				abort(c, http.StatusForbidden, "账号已被禁用")
				return
			}
			claims.Username = user.Username
			claims.Role = user.Role
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

// RequireAdminPortal 要求当前角色可访问管理后台（admin/reviewer/operator）。
func RequireAdminPortal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rbac.IsAdminRole(GetClaims(c).Role) {
			abort(c, http.StatusForbidden, "无权限")
			return
		}
		c.Next()
	}
}

// RequirePermission 要求当前角色拥有指定权限码，权限数据来自 Resolver。
func RequirePermission(resolver *rbac.Resolver, perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !resolver.Has(GetClaims(c).Role, perm) {
			abort(c, http.StatusForbidden, "无权限: "+perm)
			return
		}
		c.Next()
	}
}

func GetClaims(c *gin.Context) *model.Claims {
	return c.MustGet(ClaimsKey).(*model.Claims)
}

func abort(c *gin.Context, code int, message string) {
	response.Fail(c, code, code, message)
	c.Abort()
}
