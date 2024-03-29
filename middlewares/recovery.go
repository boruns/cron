package middlewares

import (
	"crontab/utils"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					if brokenPipe {
						utils.LoggerFor(c).Error(c.Request.URL.Path,
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
						// If the connection is dead, we can't write a status to it.
						c.Error(err.(error)) // nolint: errcheck
						c.Abort()
						return
					}

					if stack {
						utils.LoggerFor(c).Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())),
						)
					} else {
						utils.LoggerFor(c).Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
					}
					c.AbortWithStatus(http.StatusInternalServerError)

				}
			}
		}()
		c.Next()
	}
}
