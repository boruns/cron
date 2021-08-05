package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidateMobile 校验手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	// 利用反射拿到结构体tag含有mobile的key字段
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	return ok
}
