package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// 注册新用户
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if VerifyEmailFormat(username) == false {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Please enter the email format!"},
		}) //验证邮箱格式是否正确
	} else if _, exist := checkUserName(username); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		// 如果用户已经存在，输出提示用户已经存在
	} else if CheckPasswordLever(password) == false {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "You password is not safe,The password must contain at least six characters of uppercase and lowercase letters, digits and symbols"},
		}) //密码强度必须为字⺟⼤⼩写+数字+符号，6位以上
	} else {
		// 添加新用户
		id, _ := AddUserInfo(username, password)
		// 更新缓存
		fmt.Print(DemoVideos)
		fmt.Println("更新视频！！！")
		initVideos()
		initFavorite()
		fmt.Print(DemoVideos)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    username + "_" + password,
		})
	}
}

// 用户登陆
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + "_" + password
	if user, exist := FindUserInfo(token); exist {
		// 更新视频和喜欢列表缓存
		fmt.Print(DemoVideos)
		fmt.Println("更新视频！！！")
		initVideos()
		initFavorite()
		fmt.Print(DemoVideos)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.ID,
			Token:    token,
		})
	} else if _, exist := checkUserName(username); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password not correct! "},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

// 查找用户
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	// 首先根据token查找到用户ID
	if user, exits := FindUserInfo(token); exits {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		fmt.Println("User doesn't exits!****")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
