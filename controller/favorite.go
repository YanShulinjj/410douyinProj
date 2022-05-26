package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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
			vs := VID2Video(vid)
			UserFavoriteListMap[user.Name+"_"+user.Password] = append(UserFavoriteListMap[user.Name+"_"+user.Password], vs)
		}

	}
}

func VID2Video(video_id string) Video {
	fmt.Println(video_id)
	vid, err := strconv.Atoi(video_id)
	if err != nil {
		fmt.Println("*****Get video_id faild! ", err)
		return Video{}
	}
	video := FindVideo(vid)
	DemoVideos[VideosBuffer[uint(vid)]].FavoriteCount += 1
	DemoVideos[VideosBuffer[uint(vid)]].IsFavorite = true
	// 查找作者
	author, exist := FindUserByID(video.UserRefer)
	if !exist {
		author = DemoUser
	}
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

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	video_id_str := c.Query("video_id")
	vid, err := strconv.Atoi(video_id_str)
	if err != nil {
		fmt.Println("Get video_id faild! ", err)
	}
	if user, exist := FindUserInfo(token); exist {

		if DemoVideos[VideosBuffer[uint(vid)]].IsFavorite {
			DemoVideos[VideosBuffer[uint(vid)]].IsFavorite = false
			DemoVideos[VideosBuffer[uint(vid)]].FavoriteCount -= 1
			// 删除
			del_idx := -1
			for _, video := range UserFavoriteListMap[token] {
				del_idx++
				if int(video.Id) == vid {
					break
				}
			}
			// 从喜欢列表中移除
			if del_idx < len(UserFavoriteListMap[token]) {
				UserFavoriteListMap[token] = append(UserFavoriteListMap[token][:del_idx], UserFavoriteListMap[token][del_idx+1:]...)
			}
			//
			// 还没写更新到数据库
			//

			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "从喜欢列表移除！"})
		} else {
			// 将当前的video添加到当前登陆的user 的like里面
			user.LikeVideosID = user.LikeVideosID + "." + video_id_str
			UpdateUser(user)
			// 更新Map
			vs := VID2Video(video_id_str)
			UserFavoriteListMap[token] = append(UserFavoriteListMap[token], vs)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功加入喜欢列表！"})
		}

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: UserFavoriteListMap[token],
	})
}
