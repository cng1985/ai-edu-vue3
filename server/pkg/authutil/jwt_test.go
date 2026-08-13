package authutil

import (
	"testing"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func newTestManager(ttl time.Duration) *JWTManager {
	return NewJWTManager(&config.Config{JWTSecret: "test-secret", TokenTTL: ttl})
}

func TestSignVerifyRoundtrip(t *testing.T) {
	m := newTestManager(time.Hour)
	token, err := m.Sign(model.Claims{ID: "u_1", Username: "alice", Role: "learner"})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.ID != "u_1" || claims.Username != "alice" || claims.Role != "learner" {
		t.Errorf("claims 不匹配: %+v", claims)
	}
}

func TestVerifyExpiredToken(t *testing.T) {
	m := newTestManager(-time.Minute)
	token, err := m.Sign(model.Claims{ID: "u_1", Username: "alice", Role: "learner"})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if _, err := m.Verify(token); err == nil {
		t.Error("过期 token 应校验失败")
	}
}

func TestVerifyWrongSecret(t *testing.T) {
	m := newTestManager(time.Hour)
	other := NewJWTManager(&config.Config{JWTSecret: "another-secret", TokenTTL: time.Hour})
	token, _ := other.Sign(model.Claims{ID: "u_1", Username: "alice", Role: "learner"})
	if _, err := m.Verify(token); err == nil {
		t.Error("签名不匹配的 token 应校验失败")
	}
}

func TestVerifyMalformedToken(t *testing.T) {
	m := newTestManager(time.Hour)
	for _, tok := range []string{"", "abc", "a.b.c"} {
		if _, err := m.Verify(tok); err == nil {
			t.Errorf("畸形 token %q 应校验失败", tok)
		}
	}
}

// 非字符串类型的 id 字段不应导致 panic（旧实现使用裸类型断言会 panic）
func TestVerifyNonStringClaims(t *testing.T) {
	m := newTestManager(time.Hour)
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  12345,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token, _ := raw.SignedString([]byte("test-secret"))
	if _, err := m.Verify(token); err == nil {
		t.Error("id 非字符串的 token 应校验失败而非 panic")
	}
}

// 无过期时间的 token 必须被拒绝
func TestVerifyTokenWithoutExp(t *testing.T) {
	m := newTestManager(time.Hour)
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": "u_1",
	})
	token, _ := raw.SignedString([]byte("test-secret"))
	if _, err := m.Verify(token); err == nil {
		t.Error("无 exp 的 token 应校验失败")
	}
}

// alg=none 的 token 必须被拒绝
func TestVerifyNoneAlgorithm(t *testing.T) {
	m := newTestManager(time.Hour)
	raw := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"id":  "u_1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token, _ := raw.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if _, err := m.Verify(token); err == nil {
		t.Error("alg=none 的 token 应校验失败")
	}
}

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !VerifyPassword("secret123", hash) {
		t.Error("正确密码应通过校验")
	}
	if VerifyPassword("wrong", hash) {
		t.Error("错误密码不应通过校验")
	}
}
