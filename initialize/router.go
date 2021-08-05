package initialize

import (
	"crontab/middlewares"
	"crontab/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//中间件必须放在前面
	Router.Use(middlewares.TraceLogger(), middlewares.GinLogger(), middlewares.GinRecovery(true), middlewares.Cors()) //中间件
	apiRouter(Router)
	return Router
}

//api Router
func apiRouter(Router *gin.Engine) {
	ApiGroup := Router.Group("/api/v1/")
	router.UserRouter(ApiGroup) //用户相关
	router.InitRouter(ApiGroup) //公用相关
	router.JobRouter(ApiGroup)  //job相关
}
