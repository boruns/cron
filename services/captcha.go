package services

import (
	"crontab/global"
	"time"

	"github.com/mojocn/base64Captcha"
)

type CaptchaService struct {
}

type captchaRedisStore struct {
}

const captchaPrefix = "captcha:"

func (r *captchaRedisStore) Set(id, value string) error {
	key := captchaPrefix + id
	return global.Redis.Set(key, value, time.Minute*10).Err()
}

func (r *captchaRedisStore) Get(id string, clear bool) string {
	key := captchaPrefix + id
	val, err := global.Redis.Get(key).Result()
	if err != nil {
		return ""
	}
	if clear {
		if err := global.Redis.Del(key).Err(); err != nil {
			return ""
		}
	}
	return val
}

func (r *captchaRedisStore) Verify(id, answer string, clear bool) bool {
	v := r.Get(id, clear)
	return v == answer
}

//生成验证码
func (c *CaptchaService) GetCaptcha() (id, b64s string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, &captchaRedisStore{})
	return cp.Generate()
}
