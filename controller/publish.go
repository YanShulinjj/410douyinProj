package controller

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	user, exist := FindUserInfo(token)
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	// fmt.Printf("videos type: %T", data)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 使用 ffmpeg 提取视频第一帧作为封面
	cmd := exec.Command("E:/software/ffmpeg/bin/ffmpeg.exe", "-i", saveFile, filepath.Join("./public/", finalName+".jpg"))
	cmd.Run()
	// update feed
	fmt.Println("update feed videos....")
	AddVideo(token, finalName)
	// DemoVideos = append([]Video{
	// 	{
	// 		PlayUrl:       BaseURL + "static/" + finalName,
	// 		CoverUrl:      BaseURL + "static/" + finalName + ".jpg",
	// 		FavoriteCount: 0,
	// 		CommentCount:  0,
	// 		IsFavorite:    false,
	// 	},
	// }, DemoVideos...)

	// save to file
	// writeVideosFile("./data/videos")

	initVideos()

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
