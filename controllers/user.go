package controllers

import (
	"crontab/forms"
	"crontab/response"
	"crontab/services"
	"crontab/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (u *UserController) PasswordLogin(c *gin.Context) {
	request := forms.PasswordLoginRequest{}
	if err := c.ShouldBind(&request); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	userService := services.UserService{}
	user, err := userService.UserLogin(request)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), "")
		return
	}
	response.Success(c, 200, "success", user)
}

// GetUserList 获取用户列表
func (u *UserController) GetUserList(c *gin.Context) {
	// 获取参数
	UserListForm := forms.UserListRequest{}
	if err := c.ShouldBind(&UserListForm); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	userService := services.UserService{}
	total, userList := userService.GetUserList(UserListForm)
	if total == 0 {
		response.Error(c, http.StatusBadRequest, 400, "获取用户列表失败", "")
		return
	}
	response.Success(c, 200, "获取用户列表成功", map[string]interface{}{
		"total":    total,
		"userlist": userList,
	})
}
