package utils

import (
	"crontab/global"
	"crontab/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// HandleValidatorError 处理字段校验异常
func HandleValidatorError(c *gin.Context, err error) {
	//如何返回错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		response.Error(c, http.StatusInternalServerError, 500, "字段验证错误", err.Error())
		return
	}
	msg := removeTopStruct(errs.Translate(global.Trans))
	response.Error(c, http.StatusBadRequest, 400, "字段验证错误", msg)
}

//   removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		// 从文本的逗号开始切分   处理后"mobile": "mobile为必填字段"  处理前: "PasswordLoginForm.mobile": "mobile为必填字段"
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

type Func func(f1 validator.FieldLevel) bool

// RegisterValidatorFunc 注册自定义校验tag
func RegisterValidatorFunc(v *validator.Validate, tag string, msgStr string, fn Func) error {
	// 注册tag自定义校验
	if err := v.RegisterValidation(tag, validator.Func(fn)); err != nil {
		return err
	}
	//自定义错误内容
	return v.RegisterTranslation(tag, global.Trans, func(ut ut.Translator) error {
		return ut.Add(tag, "{0}"+msgStr, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}
