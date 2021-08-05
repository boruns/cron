package router

import (
	"crontab/controllers"
	"crontab/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(Router *gin.RouterGroup) {
	userController := controllers.UserController{}
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), userController.GetUserList)
		UserRouter.POST("login", userController.PasswordLogin)
	}
}
