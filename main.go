package main

import "github.com/gin-gonic/gin"

func main() {
	api := gin.Default()

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	users := api.Group("/users")

	users.POST("/", func(ctx *gin.Context) {
	})

	users.GET("/:id", func(ctx *gin.Context) {
	})

	users.PUT("/:id", func(ctx *gin.Context) {
	})

	users.DELETE("/:id", func(ctx *gin.Context) {
	})

	tasks := api.Group("/tasks")

	tasks.POST("/", func(ctx *gin.Context) {
	})

	tasks.GET("/:id", func(ctx *gin.Context) {
	})

	tasks.GET("?assignedTo={:id}", func(ctx *gin.Context) {
	})

	tasks.PUT("/:id", func(ctx *gin.Context) {
	})

	tasks.DELETE("/:id", func(ctx *gin.Context) {
	})

	auth := api.Group("/auth")

	auth.POST("/login", func(ctx *gin.Context) {
	})

	auth.POST("/logout", func(ctx *gin.Context) {
	})

	api.Run(":8080")
}
