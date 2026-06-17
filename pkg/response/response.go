package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type R struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{})      { c.JSON(http.StatusOK, R{Code: 0, Message: "success", Data: data}) }
func Created(c *gin.Context, data interface{}) { c.JSON(http.StatusCreated, R{Code: 0, Message: "created", Data: data}) }
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, R{Code: code, Message: msg})
}
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, R{Code: 401, Message: msg})
}
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, R{Code: 403, Message: msg})
}
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, R{Code: 404, Message: msg})
}
func InternalError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, R{Code: 500, Message: msg})
}
