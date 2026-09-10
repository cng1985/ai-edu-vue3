package apperr

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPStatus(t *testing.T) {
	cases := []struct {
		err  error
		want int
	}{
		{BadRequest("参数错误"), http.StatusBadRequest},
		{Unauthorized("未登录"), http.StatusUnauthorized},
		{Forbidden("无权限"), http.StatusForbidden},
		{NotFound("不存在"), http.StatusNotFound},
		{Conflict("冲突"), http.StatusConflict},
		{Internal("内部错误", errors.New("boom")), http.StatusInternalServerError},
	}
	for _, tc := range cases {
		got, ok := HTTPStatus(tc.err)
		if !ok || got != tc.want {
			t.Errorf("HTTPStatus(%v) = %d, %v; want %d, true", tc.err, got, ok, tc.want)
		}
	}
}

func TestHTTPStatusUnknownError(t *testing.T) {
	if _, ok := HTTPStatus(errors.New("普通错误")); ok {
		t.Error("普通 error 不应被识别为 apperr")
	}
}

func TestWrappedError(t *testing.T) {
	inner := NotFound("用户不存在")
	wrapped := fmt.Errorf("查询失败: %w", inner)
	got, ok := HTTPStatus(wrapped)
	if !ok || got != http.StatusNotFound {
		t.Errorf("包装后的 apperr 应仍可识别, got %d, %v", got, ok)
	}
	if !errors.Is(errors.Unwrap(wrapped), inner) {
		t.Error("Unwrap 链断裂")
	}
}

func TestErrorMessage(t *testing.T) {
	err := Forbidden("账号已被禁用")
	if err.Error() != "账号已被禁用" {
		t.Errorf("Error() = %q", err.Error())
	}
}
