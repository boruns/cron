package services

import (
	"crontab/dao"
	"crontab/forms"
	"crontab/utils"
	"errors"
)

type UserService struct {
}

func (u *UserService) UserLogin(request forms.PasswordLoginRequest) (map[string]interface{}, error) {
	//校验验证码
	captchaRedisStore := captchaRedisStore{}
	if !captchaRedisStore.Verify(request.CaptchaId, request.Captcha, true) {
		return nil, errors.New("验证码错误")
	}
	//然后获取用户
	userDao := dao.UserDao{}
	user, ok := userDao.FindUserInfo(request.Username)
	if !ok {
		return nil, errors.New("用户名或者密码错误")
	}
	//判断验证码
	if !utils.ComparePassword(request.Password, user.Password) {
		return nil, errors.New("用户名或者密码错误")
	}
	token, err := utils.LoginCreateToken(user.ID, user.NickName, user.Role)
	if err != nil {
		return nil, errors.New("token生成失败，请稍后再试")
	}

	birthday := ""
	if user.Birthday == nil {
		birthday = ""
	} else {
		birthday = user.Birthday.Format("2006-01-02")
	}
	userMap := map[string]interface{}{
		"id":       user.ID,
		"nickname": user.NickName,
		"head_url": user.HeadUrl,
		"birthday": birthday,
		"address":  user.Address,
		"desc":     user.Desc,
		"gender":   user.Gender,
		"mobile":   user.Mobile,
		"role":     user.Role,
	}
	data := map[string]interface{}{
		"token": token,
		"user":  userMap,
	}
	return data, nil
}

func (u *UserService) GetUserList(request forms.UserListRequest) (int, interface{}) {
	// 获取数据
	userDao := dao.UserDao{}
	total, userlist := userDao.GetUserListDao(request.Page, request.PageSize)
	// 判断
	if (total + len(userlist)) == 0 {
		return 0, nil
	}
	return total, userlist
}
