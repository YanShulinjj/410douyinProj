package controller

import (
	"github.com/RaymondCode/simple-demo/mylog"
	"strconv"
	"strings"
)

var BaseURL = "http://192.168.137.1:8080/"
var DemoVideos = []Video{}
var VideosBuffer = map[uint]int{}
var CommentsMap = map[uint][]Comment{}

func init() {
	initVideos()
}

func initVideos() {
	videos := GetVideos()
	DemoVideos = []Video{}
	for i, video := range videos {
		// 通过user_refer 查找用户信息
		user, exist := FindUserByID(video.UserRefer)
		if !exist {
			// fmt.Println("User not exist!")
			mylog.Logger.Panicf("User: [user_id = %d] is not exist!", video.UserRefer)
		}
		comments := GetComments(video.ID)
		CommentsMap[video.ID] = comments
		vs := Video{
			Id:            video.ID,
			Author:        user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  int64(len(comments)),
			IsFavorite:    false,
		}
		// 查找User的likeID 是否包含此视频
		vids := strings.Split(user.LikeVideosID, ".")[1:]
		for _, vid := range vids {
			if vid == strconv.Itoa(int(video.ID)) {
				vs.IsFavorite = true
				break
			}
		}
		VideosBuffer[video.ID] = i
		DemoVideos = append(DemoVideos, vs)
	}
	// logger.Printf(": Initlize the vidoesList, [len=%d]", len(DemoVideos))
}
