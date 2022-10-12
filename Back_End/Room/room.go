package Room

import (
	"EnChat/DB"
	"EnChat/Token"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {
	r.GET("/data/:room_id", func(ctx *gin.Context) {
		room_Id := ctx.Param("room_id")
		res := struct {
			Name      string
			Admin     int
			Join_List []DB.Join_Room
		}{}

		room_data := DB.Room{}
		if DB.GetDB().Where("id = ?", room_Id).First(&room_data).Error != nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Not found Room"})
			return
		}

		res.Name = room_data.Name
		res.Admin = room_data.Admin
		DB.GetDB().Where("room_id = ?", room_Id).Find(&res.Join_List)

		ctx.JSON(200, res)
	})
	r.POST("/create", func(ctx *gin.Context) {
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)
		body := struct {
			Name string `json:"name" binding:"required"`
		}{}
		ctx.BindJSON(&body)

		room_data := DB.Room{}
		D := DB.GetDB().Create(&DB.Room{Admin: Token_data.UserID, Name: body.Name})
		D.Order("created_at desc").Where("name = ? and admin = ?", body.Name, Token_data.UserID).Limit(1).Find(&room_data)
		DB.GetDB().Create(&DB.Join_Room{User_Id: Token_data.UserID, Room_Id: room_data.ID})

		ctx.JSON(200, gin.H{"code": 200})
	})
	r.DELETE("/remove/:room_id", func(ctx *gin.Context) {
		room_Id := ctx.Param("room_id")
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)

		room_data := DB.Room{}

		if DB.GetDB().Where("id = ? and admin = ?", room_Id, Token_data.UserID).First(&room_data).Error != nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Not found Room"})
			return
		}
		DB.GetDB().Where("id = ?", room_Id).Delete(&DB.Room{})
		ctx.JSON(200, gin.H{"code": 200})
	})
}
