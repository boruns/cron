package router

import (
	"crontab/controllers"
	"crontab/middlewares"

	"github.com/gin-gonic/gin"
)

func JobRouter(Router *gin.RouterGroup) {
	jobController := controllers.JobController{}
	JobRouter := Router.Group("job", middlewares.JWTAuth())
	{
		JobRouter.POST("job", jobController.Store)            //保存定时任务
		JobRouter.DELETE("job/:jobId", jobController.Destory) //删除任务
		JobRouter.GET("job", jobController.Index)             //任务列表
		JobRouter.POST("kill/:jobId", jobController.KillJob)  //强杀任务
		JobRouter.PUT("job/:jobId", jobController.UpdateJob)  //更新任务
		JobRouter.GET("job/:jobId", jobController.Show)       //任务详情
		JobRouter.GET("log/:jobId", jobController.Logs)       //job日志
		JobRouter.GET("work", jobController.WorkNode)         //获取工作节点
	}
}
