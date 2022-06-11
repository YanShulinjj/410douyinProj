package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/mylog"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

var UsePublishVideosMap = map[uint][]Video{}

func init() {
	initPublishVideosMap()
}

func initPublishVideosMap() {
	// 对于每一个用户，初始化其点赞列表
	users := GetUsersBriefInfo()
	for _, user := range users {
		videosamples := FindUsersVideos(user.ID)
		for _, videosample := range videosamples {
			video := Video{
				Id:            videosample.ID,
				Author:        user,
				PlayUrl:       videosample.PlayUrl,
				CoverUrl:      videosample.CoverUrl,
				FavoriteCount: videosample.FavoriteCount,
				CommentCount:  videosample.CommentCount,
				IsFavorite:    false,
			}
			UsePublishVideosMap[user.ID] = append([]Video{video}, UsePublishVideosMap[user.ID]...)
		}
	}
	// logger.Printf("Initlize PublishVideosMap[key=user_id, v=videos].\n")
}

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
	mylog.Logger.Printf("Saved the video. [filename = %s].\n", finalName)
	// 使用 ffmpeg 提取视频第1秒处的帧作为封面
	cmd := exec.Command("E:/software/ffmpeg/bin/ffmpeg.exe", "-i", saveFile, "-ss", "00:00:01", filepath.Join("./public/", finalName+".jpg"))
	cmd.Run()
	// update feed
	mylog.Logger.Printf("Capture the video's cover. [filename = %s].\n", finalName)
	vs, err := AddVideo(token, finalName)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	mylog.Logger.Printf("Add the Video to database. [filename = %s].\n", finalName)
	video := Video{
		Id:            vs.ID,
		Author:        user,
		PlayUrl:       vs.PlayUrl,
		CoverUrl:      vs.CoverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	}
	DemoVideos = append([]Video{video}, DemoVideos...)

	// 添加到自己发表的列表
	UsePublishVideosMap[user.ID] = append([]Video{video}, UsePublishVideosMap[user.ID]...)
	mylog.Logger.Printf("Add the Video to aothor's public List. [video_id = %d].\n", vs.ID)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		mylog.Logger.Println("Get Video_id failed!", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get user_id faile !"})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: UsePublishVideosMap[uint(user_id)],
	})
}
