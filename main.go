package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Running this...")
	go func() {
		controller.RunComsumer()
	}()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
