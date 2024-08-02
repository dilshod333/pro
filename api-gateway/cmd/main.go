package main

import (
	"api-gateway/api" 
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		api.Conn().GetUserByName(c, name)
	})

	r.GET("/user/id/:id", func(c *gin.Context) {
		id := c.Param("id")
		api.Conn().GetUserId(c, id)
	})

	r.POST("/user/register", func(c *gin.Context) {
		api.Conn().RegisterUser(c)

	})

	r.POST("/user/verify", func(ctx *gin.Context) {
		api.Conn().VeriffyUser(ctx)
	})

	r.POST("/user/login", func(ctx *gin.Context) {
		api.Conn().Loginn(ctx)
	})

	r.POST("/devices", func(ctx *gin.Context) {
		api.Conn().CreateDevicee(ctx)
	})

	r.DELETE("/device/id/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		api.Conn().DeleteByid(ctx, id)
	})

	r.PUT("/device/update", func(ctx *gin.Context) {
		api.Conn().Updateee(ctx)
	})

	r.POST("/device/control", func(ctx *gin.Context) {
		api.Conn().CreateCommand(ctx)
	})

	log.Println("API Gateway running on :7777")
	r.Run(":7777")
}

