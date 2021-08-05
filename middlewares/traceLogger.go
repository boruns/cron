package middlewares

import (
	"crontab/global"
	"crontab/utils"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

func TraceLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 每个请求生成的请求traceId具有全局唯一性
		u1, _ := uuid.NewV4()
		traceId := u1.String()
		global.TraceId = traceId
		utils.LoggerNewContext(ctx, zap.String("traceId", traceId))
		// 为日志添加请求的地址以及请求参数等信息
		utils.LoggerNewContext(ctx, zap.String("request.method", ctx.Request.Method))
		headersStr, _ := json.Marshal(ctx.Request.Header)
		utils.LoggerNewContext(ctx, zap.String("request.headers", string(headersStr)))
		utils.LoggerNewContext(ctx, zap.String("request.url", ctx.Request.URL.String()))

		// 将请求参数json序列化后添加进日志上下文
		if ctx.Request.Form == nil {
			ctx.Request.ParseMultipartForm(32 << 20)
		}
		form, _ := json.Marshal(ctx.Request.Form)
		utils.LoggerNewContext(ctx, zap.String("request.params", string(form)))

		ctx.Next()
	}
}
