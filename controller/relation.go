package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var UserFollowMap = map[uint][]User{}
var UserFollowerMap = map[uint][]User{}
var UserFollowCountMap = map[uint]int64{}
var UserFollowerCountMap = map[uint]int64{}

func init() {
	initMaps()
}

func initMaps() {
	users := GetUsersBriefInfo()
	for _, user := range users {
		// 解析follow
		fuids := strings.Split(user.FollowID, ".")[1:]
		for _, fuid := range fuids {
			fid, err := strconv.Atoi(fuid)
			if err != nil {
				fmt.Println("*****Get video_id faild! ", err)
				continue
			}
			follow, _ := FindUserByID(uint(fid))
			UserFollowMap[user.ID] = append(UserFollowMap[user.ID], follow)
			UserFollowerMap[follow.ID] = append(UserFollowerMap[follow.ID], user)
		}
		UserFollowCountMap[user.ID] = user.FollowCount
		UserFollowerCountMap[user.ID] = user.FollowerCount
	}
}

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// 点击关注
// 将对方用户添加到自身列表
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		fmt.Println("Get user_id faild! ", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get user_id faile !"})
	}
	if user, exist := FindUserInfo(token); exist {
		follow, _ := FindUserByID(uint(to_user_id))
		follow.IsFollow = true
		// 更改缓冲区
		UserFollowMap[user.ID] = append(UserFollowMap[user.ID], follow)
		UserFollowCountMap[user.ID]++

		UserFollowerMap[uint(to_user_id)] = append(UserFollowerMap[uint(to_user_id)], user)
		UserFollowerCountMap[uint(to_user_id)]++
		// 写入数据库
		user.FollowCount++
		user.FollowID += "." + strconv.Itoa(to_user_id)
		UpdateUser(user)
		follow.FollowerCount++
		follow.FollowerID += "." + strconv.Itoa(int(user.ID))
		UpdateUser(follow)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get user_id faile !"})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: UserFollowMap[uint(user_id)],
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		fmt.Println("Get user_id faild! ", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get user_id faile !"})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: UserFollowerMap[uint(user_id)],
	})
}
