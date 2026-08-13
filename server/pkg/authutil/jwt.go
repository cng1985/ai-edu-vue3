// Package authutil 提供 JWT 签发/校验与密码哈希工具。
package authutil

import (
	"errors"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidToken = errors.New("invalid token")

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

// tokenClaims 使用标准 RegisteredClaims 承载过期时间，
// 自定义字段保持 id/username/role 的 JSON 键名以兼容旧 token。
type tokenClaims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{secret: []byte(cfg.JWTSecret), ttl: cfg.TokenTTL}
}

func (j *JWTManager) Sign(claims model.Claims) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		UserID:   claims.ID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	return token.SignedString(j.secret)
}

func (j *JWTManager) Verify(tokenStr string) (*model.Claims, error) {
	var claims tokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid || claims.UserID == "" {
		return nil, ErrInvalidToken
	}
	return &model.Claims{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
