package main

import (
	"EnChat/Friend"
	"EnChat/Join"
	"EnChat/Room"
	"EnChat/Token"
	"EnChat/User"
	"EnChat/socket"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func Server() *gin.Engine {
	r := gin.Default()

	r.Use(GinMiddleware("*"))
	r.GET("/ws", socket.Socket)

	user_r := r.Group("user")
	User.Route(user_r)

	r.Use(Token.AuthMiddleWare)

	friend_r := r.Group("friend")
	Friend.Route(friend_r)

	room_r := r.Group("room")
	Room.Route(room_r)

	join_r := r.Group("join")
	Join.Route(join_r)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"code": 200})
	})
	return r
}

func main() {
	r := Server()
	r.Run(":1009")
}
