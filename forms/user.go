package forms

type PasswordLoginRequest struct {
	Password  string `json:"password" form:"password" binding:"required,min=3,max=20"`
	Username  string `json:"username" form:"username" binding:"required"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required,max=5,min=5"`
	CaptchaId string `json:"captcha_id" form:"captcha_id" binding:"required"`
}

type UserListRequest struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"page_size" form:"page_size" binding:"required"`
}
