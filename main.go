package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	// 启动一个消费者进程
	go func() {
		controller.RunComsumer()
	}()
	r := gin.Default()

	InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
