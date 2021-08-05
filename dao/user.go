package dao

import (
	"crontab/global"
	"crontab/models"
)

type UserDao struct {
}

func (u *UserDao) GetUserListDao(page, pageSize int) (int, []interface{}) {
	var users []models.User
	// 分页用户列表数据
	userList := make([]interface{}, 0, len(users))
	//计算偏移量
	offest := (page - 1) * pageSize
	var total int64
	result := global.DB.Model(&models.User{}).Count(&total)
	if total == 0 {
		return 0, userList
	}

	//查询数据
	result.Offset(offest).Limit(pageSize).Find(&users)

	for _, useSingle := range users {
		birthday := ""
		if useSingle.Birthday == nil {
			birthday = ""
		} else {
			// 给未设置生日的初始值
			birthday = useSingle.Birthday.Format("2006-01-02")
		}
		userItemMap := map[string]interface{}{
			"id": useSingle.ID,
			// "password":  useSingle.Password,
			"nick_name": useSingle.NickName,
			"head_url":  useSingle.HeadUrl,
			"birthday":  birthday,
			"address":   useSingle.Address,
			"desc":      useSingle.Desc,
			"gender":    useSingle.Gender,
			"role":      useSingle.Role,
			"mobile":    useSingle.Mobile,
		}
		userList = append(userList, userItemMap)
	}
	return int(total), userList
}

func (u *UserDao) FindUserInfo(username string) (*models.User, bool) {
	var user models.User
	rows := global.DB.Where(&models.User{NickName: username}).Find(&user)
	if rows.RowsAffected < 1 {
		return nil, false
	}
	return &user, true
}
