package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

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

	// 如果用户已经存在，输出提示不新创用户
	if _, exist := FindUserName(username); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 添加新用户
		id, _ := AddUserInfo(username, password)
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
	fmt.Println("Login.....")
	token := username + "_" + password
	if user, exist := FindUserInfo(token); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.ID,
			Token:    token,
		})
	} else if _, exist := FindUserName(username); exist {
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
