package middlewares

import (
	"crontab/response"
	"crontab/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerToken := c.Request.Header.Get("Authorization")
		if headerToken == "" {
			response.Error(c, http.StatusUnauthorized, 401, "未授权", "")
			c.Abort()
			return
		}
		tokenString := strings.Split(headerToken, " ")
		tokenType := tokenString[0]
		token := tokenString[1]
		if tokenType != "Bearer" || token == "" {
			response.Error(c, http.StatusUnauthorized, 401, "未授权", "")
			c.Abort()
			return
		}
		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.ErrTokenExpired {
				response.Error(c, http.StatusUnauthorized, 401, "授权已过期", "")
				c.Abort()
				return
			}
			response.Error(c, http.StatusUnauthorized, 401, "未登陆", "")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}
