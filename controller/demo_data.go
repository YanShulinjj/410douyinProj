package controller

import "fmt"

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
			fmt.Println("User not exist!")
			user = DemoUser
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
		VideosBuffer[video.ID] = i
		DemoVideos = append(DemoVideos, vs)
	}
	fmt.Println("Num of videos: ", len(DemoVideos))
}

var DemoUser = User{
	ID:       1,
	Name:     "TestUser",
	IsFollow: false,
}

// var VideoCommentsMap = map[int][]Comment{}
