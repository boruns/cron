package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code int, msg interface{}, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Error(c *gin.Context, httpCode, code int, msg interface{}, data interface{}) {
	c.JSON(httpCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
