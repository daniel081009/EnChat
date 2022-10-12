package Join

import (
	"EnChat/DB"
	"EnChat/Token"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {
	r.POST("room/:room_id", func(ctx *gin.Context) {
		room_Id := ctx.Param("room_id")
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)

		room_data := DB.Room{}
		if DB.GetDB().Where("id = ?", room_Id).First(&room_data).Error != nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Not found Room"})
			return
		}
		if DB.GetDB().Where("user_id= ? and room_id = ?", Token_data.UserID, room_Id).First(&DB.Join_Room{}).Error == nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "already"})
			return
		}

		room_id_int, _ := strconv.Atoi(room_Id)
		DB.GetDB().Create(&DB.Join_Room{User_Id: Token_data.UserID, Room_Id: room_id_int})

		ctx.JSON(200, gin.H{"code": 200})
	})
	r.DELETE("room/:room_id", func(ctx *gin.Context) {
		room_Id := ctx.Param("room_id")
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)

		room_data := DB.Room{}
		if DB.GetDB().Where("id = ?", room_Id).First(&room_data).Error != nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Not found Room"})
			return
		}
		if DB.GetDB().Where("user_id= ? and room_id = ?", Token_data.UserID, room_Id).First(&DB.Join_Room{}).Error != nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Not found Join"})
			return
		}

		if DB.GetDB().Where("user_id= ? and room_id = ?", Token_data.UserID, room_Id).Delete(&DB.Join_Room{}).Error != nil {
			ctx.JSON(200, gin.H{"code": 0})
			return
		}
		ctx.JSON(200, gin.H{"code": 200})
	})
}
