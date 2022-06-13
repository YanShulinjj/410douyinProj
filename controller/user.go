package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/mylog"
	"github.com/gin-gonic/gin"
	"net/http"
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
   if  VerifyEmailFormat(username) {
	   // 如果用户已经存在，输出提示不新建用户
	   if _, exist := checkUserName(username); exist {
		   mylog.Logger.Printf("User:[username=%s] registered failed, user already exist!\n", username)
		   c.JSON(http.StatusOK, UserLoginResponse{
			   Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		   })
	   } else {

		   // 添加新用户
		   user_id, _ := AddUserInfo(username, password)
		   // 更新缓存
		   initVideos()
		   initFavorite()
		   mylog.Logger.Printf("User:[user_id=%d] registered successfully!\n", user_id)
		   c.JSON(http.StatusOK, UserLoginResponse{
			   Response: Response{StatusCode: 0},
			   UserId:   user_id,
			   Token:    username + "_" + password,
		   })
	   }
   }else {
	   c.JSON(http.StatusOK, UserLoginResponse{
		   Response: Response{StatusCode: 1, StatusMsg: "User name does not conform to mailbox format"},
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
		initVideos()
		initFavorite()
		mylog.Logger.Printf("User:[user_id=%d] login in successfully!\n", user.ID)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.ID,
			Token:    token,
		})
	} else if _, exist := checkUserName(username); exist {
		mylog.Logger.Printf("User:[user_name=%s] login in failed,Password not correct!\n", username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Password not correct! "},
		})
	} else {
		mylog.Logger.Printf("User:[user_name=%s] login in failed,No such user!\n", username)
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
