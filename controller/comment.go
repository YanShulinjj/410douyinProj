package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	var CommentsMap map[uint][]Comment
	token := c.Query("token")
	text := c.Query("comment_text")
	video_id, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		fmt.Println("Get video_id faild! ", err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video not found!"})
	}
	// 增加comment
	fmt.Println("Adding comment...")
	//从redis中读取CommentsMap
	CommentsMap = Getrediscomment("CommentsMap", Client)
	newcomment, err := AddComment(token, uint(video_id), text)
	// DemoVideos[VideosBuffer[uint(video_id)]].CommentCount += 1
	// 输出Videos的coment
	CommentsMap[uint(video_id)] = append(CommentsMap[uint(video_id)], newcomment)
	//更新redis中的CommentsMap
	Saverediscomment(" CommentsMap", Client, CommentsMap)
	VideosBuffer := GetVideosBuffer("VideosBuffer", Client)
	DemoVideos[VideosBuffer[uint(video_id)]].CommentCount += 1
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Add comment failed!"})
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

//
func CommentList(c *gin.Context) {
	video_id, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		fmt.Println("Get video_id faild! ", err)
	}
	// fmt.Printf("Getting comment of video: [%d]\n", video_id)
	// comments := GetComments(uint(video_id))
	// 更新视频评论数
	// DemoVideos[VideosBuffer[uint(video_id)]].CommentCount = int64(len(comments))
	//从redis中读取CommentsMap
	var CommentsMap map[uint][]Comment
	CommentsMap = Getrediscomment("CommentsMap", Client)
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsMap[uint(video_id)],
		// CommentList: VideoCommentsMap[video_id],
	})
}
