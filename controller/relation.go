package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//var UserFollowMap = map[uint][]User{}       //用户关注的人
//var UserFollowerMap = map[uint][]User{}     //该用户的粉丝
//var UserFollowCountMap = map[uint]int64{}   //用户关注人的数量
//var UserFollowerCountMap = map[uint]int64{} //用户粉丝的数量

func init() {
	initMaps()
}

func initMaps() {
	users := GetUsersBriefInfo()
	var UserFollowMap = map[uint][]User{}
	var UserFollowerMap = map[uint][]User{}
	var UserFollowCountMap = map[uint]int64{}
	var UserFollowerCountMap = map[uint]int64{}
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
	SaveUserFollowMap("UserFollowMap", Client, UserFollowMap)
	SaveUserFollowerMap("UserFollowerMap", Client, UserFollowerMap)
	SaveUserFollowCountMap("UserFollowCountMap", Client, UserFollowCountMap)
	SaveUserFollowerCountMap("UserFollowerCountMap", Client, UserFollowerCountMap)
}

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// 点击关注

func RelationAction(c *gin.Context) {

	token := c.Query("token")                              //正在使用中的用户
	to_user_id, err := strconv.Atoi(c.Query("to_user_id")) //发布视频的用户
	action_type := c.Query("action_type")
	UserFollowerMap := GetUserFollowerMap("UserFollowerMap", Client)
	if err != nil {
		fmt.Println("Get user_id faild! ", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get user_id faile !"})
	}
	follow, _ := FindUserByID(uint(to_user_id))
	if user, exist := FindUserInfo(token); exist {
		// 自己关注不了自己
		if follow.Name == user.Name {

			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Can't focus on yourself"})
		} else {
			//如果没有添加关注
			if action_type == "1" {
				sum := 0
				for _, uer := range UserFollowerMap[follow.ID] {
					if user.ID == uer.ID {
						sum++
					}
				}
				if sum == 0 {
					AddConcern(user, follow)
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				} else {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Already concerned"})
				}

			}

			if action_type == "2" {
				CancelConcern(user, follow) //取消关注
				c.JSON(http.StatusOK, Response{StatusCode: 0})
			}

		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Please login"})
	}
}

// 添加关注

func AddConcern(NowUser User, follow User) {
	follow.IsFollow = true
	UserFollowMap := GetUserFollowMap("UserFollowMap", Client)
	UserFollowCountMap := GetUserFollowCountMap("UserFollowCountMap", Client)
	UserFollowerMap := GetUserFollowerMap("UserFollowerMap", Client)
	UserFollowerCountMap := GetUserFollowerCountMap("UserFollowerCountMap", Client)
	// 更改缓冲区
	UserFollowMap[NowUser.ID] = append(UserFollowMap[NowUser.ID], follow)
	UserFollowCountMap[NowUser.ID]++

	UserFollowerMap[follow.ID] = append(UserFollowerMap[follow.ID], NowUser)
	UserFollowerCountMap[follow.ID]++
	SaveUserFollowMap("UserFollowMap", Client, UserFollowMap)
	SaveUserFollowCountMap("UserFollowCountMap", Client, UserFollowCountMap)
	SaveUserFollowerMap("UserFollowerMap", Client, UserFollowerMap)
	SaveUserFollowerCountMap("UserFollowerCountMap", Client, UserFollowerCountMap)
	// 写入数据库
	NowUser.FollowCount = UserFollowCountMap[NowUser.ID]

	if NowUser.FollowID == "." {
		NowUser.FollowID += strconv.Itoa(int(follow.ID))
	} else {

		NowUser.FollowID += "." + strconv.Itoa(int(follow.ID))
	}

	UpdateUser(NowUser)

	follow.FollowerCount = UserFollowerCountMap[follow.ID]
	if follow.FollowerID == "." {
		follow.FollowerID += strconv.Itoa(int(NowUser.ID))
	} else {
		follow.FollowerID += "." + strconv.Itoa(int(NowUser.ID))
	}
	UpdateUser(follow)

}

//取消关注

func CancelConcern(NowUser User, follow User) {

	follow.IsFollow = false

	//从当前登录用户关注列表中移除
	// 更改缓冲区
	UserFollowMap := GetUserFollowMap("UserFollowMap", Client)
	UserFollowCountMap := GetUserFollowCountMap("UserFollowCountMap", Client)

	for index, u := range UserFollowMap[NowUser.ID] {
		if u.ID == follow.ID {
			UserFollowMap[NowUser.ID] = append(UserFollowMap[NowUser.ID][:index], UserFollowMap[NowUser.ID][index+1:]...)
			UserFollowCountMap[NowUser.ID]--
			//  fmt.Println("当前用户的关注数：", UserFollowCountMap[NowUser.ID])
			NowUser.FollowCount = UserFollowCountMap[NowUser.ID]
			uids := strings.Split(NowUser.FollowID, ".")[1:]
			uids = append(uids[:index], uids[index+1:]...)
			uids_str := "." + strings.Join(uids, ".")
			NowUser.FollowID = uids_str
			break
		}
	}
	// 修改对象
	SaveUserFollowMap("UserFollowMap", Client, UserFollowMap)
	SaveUserFollowCountMap("UserFollowCountMap", Client, UserFollowCountMap)
	// 写入数据库
	UpdateUser(NowUser)

	//发布视频用户的粉丝列表里删除
	//修改缓冲区
	UserFollowerMap := GetUserFollowerMap("UserFollowerMap", Client)
	UserFollowerCountMap := GetUserFollowerCountMap("UserFollowerCountMap", Client)
	for index, u := range UserFollowerMap[follow.ID] {

		if u.ID == NowUser.ID {
			UserFollowerMap[follow.ID] = append(UserFollowerMap[follow.ID][:index], UserFollowerMap[follow.ID][index+1:]...)
			UserFollowerCountMap[follow.ID]--
			//	fmt.Println("视频用户的粉丝数：",UserFollowerCountMap[follow.ID])

			follow.FollowerCount = UserFollowerCountMap[follow.ID]
			uids := strings.Split(follow.FollowerID, ".")[1:]
			uids = append(uids[:index], uids[index+1:]...)
			uids_str := "." + strings.Join(uids, ".")

			follow.FollowerID = uids_str
			break
		}
	}
	// 修改对象
	SaveUserFollowerMap("UserFollowerMap", Client, UserFollowerMap)
	SaveUserFollowerCountMap("UserFollowerCountMap", Client, UserFollowerCountMap)
	//写入数据库
	UpdateUser(follow)

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Query("user_id"))
	UserFollowMap := GetUserFollowMap("UserFollowMap", Client)
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
	UserFollowerMap := GetUserFollowMap("UserFollowerMap", Client)
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
