package router

import (
	"crontab/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.RouterGroup) {
	captchaController := controllers.CaptchaController{}
	testController := controllers.TestController{}
	BaseRouter := router.Group("base")
	{
		BaseRouter.GET("captcha", captchaController.GetCaptcha)
		BaseRouter.GET("test", testController.Test)
	}
}
