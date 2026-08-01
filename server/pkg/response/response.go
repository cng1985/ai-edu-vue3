package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}, message ...string) {
	msg := "success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusOK, Body{Code: 0, Message: msg, Data: data})
}

func Fail(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Body{Code: code, Message: message, Data: nil})
}
