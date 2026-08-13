// Package apperr 提供带语义分类的业务错误类型，
// 供 handler 层统一映射为 HTTP 状态码，替代按错误消息字符串匹配的做法。
package apperr

import (
	"errors"
	"net/http"
)

// Kind 表示错误的语义分类。
type Kind int

const (
	KindBadRequest Kind = iota
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
	KindInternal
)

// Error 是携带分类信息的业务错误。
type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Err }

func New(kind Kind, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

func Wrap(kind Kind, message string, err error) *Error {
	return &Error{Kind: kind, Message: message, Err: err}
}

func BadRequest(message string) *Error   { return New(KindBadRequest, message) }
func Unauthorized(message string) *Error { return New(KindUnauthorized, message) }
func Forbidden(message string) *Error    { return New(KindForbidden, message) }
func NotFound(message string) *Error     { return New(KindNotFound, message) }
func Conflict(message string) *Error     { return New(KindConflict, message) }

func Internal(message string, err error) *Error {
	return Wrap(KindInternal, message, err)
}

// HTTPStatus 将错误映射为 HTTP 状态码。
// 第二个返回值表示该错误是否为 *Error 类型（未识别时由调用方决定兜底状态码）。
func HTTPStatus(err error) (int, bool) {
	var e *Error
	if !errors.As(err, &e) {
		return 0, false
	}
	switch e.Kind {
	case KindBadRequest:
		return http.StatusBadRequest, true
	case KindUnauthorized:
		return http.StatusUnauthorized, true
	case KindForbidden:
		return http.StatusForbidden, true
	case KindNotFound:
		return http.StatusNotFound, true
	case KindConflict:
		return http.StatusConflict, true
	default:
		return http.StatusInternalServerError, true
	}
}
