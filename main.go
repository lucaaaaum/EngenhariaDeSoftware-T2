package main

import "github.com/gin-gonic/gin"

func main() {
	api := gin.Default()

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	api.Run(":8080")
}
