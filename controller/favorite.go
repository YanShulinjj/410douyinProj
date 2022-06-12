package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// var CommentsMap map[uint][]Comment

// 用于存储用户点赞列表
var UserFavoriteListMap = map[string][]Video{}

func init() {
	initFavorite()
}

func initFavorite() {
	// 对于每一个用户，初始化其点赞列表

	users := GetUsersBriefInfo()
	for _, user := range users {
		// 解析user.LikeVideosID
		vids := strings.Split(user.LikeVideosID, ".")[1:]
		for _, vid := range vids {
			vs := VID2Video(vid, true)
			UserFavoriteListMap[user.Name+"_"+user.Password] = append(UserFavoriteListMap[user.Name+"_"+user.Password], vs)
		}
	}
}

// 根据video_id 转成 Video对象

func VID2Video(video_id string, init bool) Video {
	fmt.Println(video_id)
	vid, err := strconv.Atoi(video_id)
	if err != nil {
		return Video{}
	}
	video := FindVideo(vid)
	// VideosBuffer := GetVideosBuffer("VideosBuffer", Client)

	DemoVideos[VideosBuffer[uint(vid)]].FavoriteCount += 1
	if init {
		DemoVideos[VideosBuffer[uint(vid)]].IsFavorite = false
	} else {
		DemoVideos[VideosBuffer[uint(vid)]].IsFavorite = true
	}

	// 查找作者
	author, exist := FindUserByID(video.UserRefer)
	if !exist {
		panic("User not extis !")
	}

	// 从redis读取
	//	CommentsMap = Getrediscomment("CommentsMap", Client)
	vs := Video{
		Id:            video.ID,
		Author:        author,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  int64(len(CommentsMap[video.ID])),
		IsFavorite:    false,
	}
	return vs
}

// 点击喜欢按钮
// 将当前视频video_id 加入或移除用户喜欢列表

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")           // 当前用户的信息
	video_id_str := c.Query("video_id") // 发布视频用户的信息
	action_type := c.Query("action_type")

	println(video_id_str)

	vid, err := strconv.Atoi(video_id_str)
	if err != nil {
		fmt.Println("Get video_id faild! ", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Get video_id faild! "})
	}
	// 从redis中获取map
	//	VideosBuffer := GetVideosBuffer("VideosBuffer", Client)

	if user, exist := FindUserInfo(token); exist {

		if action_type == "2" {

			DemoVideos[VideosBuffer[uint(vid)]].IsFavorite = false
			DemoVideos[VideosBuffer[uint(vid)]].FavoriteCount -= 1
			fmt.Println("UserID: ", token, "在喜欢VideoID: ", vid)
			// 取消喜欢
			CancelFavorite(user, vid, token)

			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "从喜欢列表移除！"})
		}

		if action_type == "1" {
			// 添加喜欢
			AddFavorite(user, video_id_str, token)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功加入喜欢列表！"})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 取消喜欢

func CancelFavorite(user User, Vid int, token string) {

	for index, video := range UserFavoriteListMap[token] {

		if int(video.Id) == Vid {
			UserFavoriteListMap[token] = append(UserFavoriteListMap[token][:index], UserFavoriteListMap[token][index+1:]...)
			vids := strings.Split(user.LikeVideosID, ".")[1:]
			vids = append(vids[:index], vids[index+1:]...)
			vids_str := "." + strings.Join(vids, ".")
			user.LikeVideosID = vids_str
			UpdateUser(user)
			break

		}
	}

}

// 添加喜欢

func AddFavorite(user User, video_id_str string, token string) {

	if user.LikeVideosID == "." {
		user.LikeVideosID = user.LikeVideosID + video_id_str
	} else {
		user.LikeVideosID = user.LikeVideosID + "." + video_id_str
	}

	// 更新Map
	vs := VID2Video(video_id_str, false)
	UserFavoriteListMap[token] = append(UserFavoriteListMap[token], vs)

	UpdateUser(user)
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")

	fmt.Println(UserFavoriteListMap[token])

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: UserFavoriteListMap[token],
	})
}
