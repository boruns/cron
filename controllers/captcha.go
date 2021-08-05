package controllers

import (
	"crontab/response"
	"crontab/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CaptchaController struct {
}

func (c *CaptchaController) GetCaptcha(ctx *gin.Context) {
	captchaService := services.CaptchaService{}
	id, b64s, err := captchaService.GetCaptcha()
	if err != nil {
		zap.S().Errorf("验证码生成错误,:%s ", err.Error())
		response.Error(ctx, http.StatusInternalServerError, 500, "验证码生成错误", "")
		return
	}
	response.Success(ctx, 200, "生成验证码成功", gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
