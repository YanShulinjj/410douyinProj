package controller

import (
	"github.com/RaymondCode/simple-demo/mylog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	text := c.Query("comment_text")
	video_id, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		// fmt.Println("Get video_id faild! ", err)
		mylog.Logger.Println("Get Video_id failed!", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video not found!"})
	}
	// 增加comment
	newcomment, err := AddComment(token, uint(video_id), text)

	// 输出Videos的coment
	CommentsMap[uint(video_id)] = append(CommentsMap[uint(video_id)], newcomment)
	DemoVideos[VideosBuffer[uint(video_id)]].CommentCount += 1

	mylog.Logger.Printf("User: [%s] add a comment", token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Add comment failed!"})
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

func CommentList(c *gin.Context) {
	video_id, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		mylog.Logger.Println("Get Video_id failed!", err)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsMap[uint(video_id)],
		// CommentList: VideoCommentsMap[video_id],
	})
}
