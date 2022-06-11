package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/mylog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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
	mylog.Logger.Printf("Initlize: FavoriteMap[key=user_token,v=videos]\n")
}

// 根据video_id 转成 Video对象
func VID2Video(video_id string, init bool) Video {
	fmt.Println(video_id)
	vid, err := strconv.Atoi(video_id)
	if err != nil {
		return Video{}
	}
	video := FindVideo(vid)
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
	token := c.Query("token")
	video_id_str := c.Query("video_id")
	vid, err := strconv.Atoi(video_id_str)
	action_type := c.Query("action_type")
	if err != nil {
		mylog.Logger.Println("Get Video_id failed!", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Get video_id faild! "})
	}
	if user, exist := FindUserInfo(token); exist {
		if action_type == "2" {
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
			if del_idx > -1 && del_idx < len(UserFavoriteListMap[token]) {
				UserFavoriteListMap[token] = append(UserFavoriteListMap[token][:del_idx], UserFavoriteListMap[token][del_idx+1:]...)
				vids := strings.Split(user.LikeVideosID, ".")[1:]
				vids = append(vids[:del_idx], vids[del_idx+1:]...)
				vids_str := "." + strings.Join(vids, ".")
				user.LikeVideosID = vids_str
			}
			mylog.Logger.Printf("User: [token=%s] canceled like [video_id=%d].\n", token, vid)
			// 更新到数据库
			UpdateUser(user)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "从喜欢列表移除！"})
		} else {
			// 将当前的video添加到当前登陆的user 的like里面
			user.LikeVideosID = user.LikeVideosID + "." + video_id_str
			UpdateUser(user)
			// 更新Map
			vs := VID2Video(video_id_str, false)
			UserFavoriteListMap[token] = append(UserFavoriteListMap[token], vs)
			mylog.Logger.Printf("User: [token=%s] liked [video_id=%d].\n", token, vid)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "成功加入喜欢列表！"})
		}
	} else {
		mylog.Logger.Printf("User: [token=%s] is not existent!\n", token)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	fmt.Println("userflist: ", UserFavoriteListMap[token])
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: UserFavoriteListMap[token],
	})
}
