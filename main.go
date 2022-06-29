package main

import (
	"websocket/socket"

	"github.com/gin-gonic/gin"
)

//	"websocket/socket"

func main() {
	router := gin.Default()

	go socket.StartEngine(socket.ActivePipelines)
	go socket.Serversetup()
	router.POST("/pipelines", socket.CreatePipeline)
	router.Run(":8080")

}
