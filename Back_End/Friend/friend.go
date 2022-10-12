package Friend

import (
	"EnChat/DB"
	"EnChat/Token"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {
	r.POST("Create/:friend_Id", func(ctx *gin.Context) {
		friend_Id := ctx.Param("friend_Id")
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)

		friend_data := DB.Friend{}

		if DB.GetDB().Where("friend_1 = ? and friend_2 = ? or friend_1 = ? and friend_2 = ?", Token_data.UserID, friend_Id, friend_Id, Token_data.UserID).First(&friend_data).Error == nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "already"})
			return
		}

		friend_id_int, err := strconv.Atoi(friend_Id)
		if err != nil {
			fmt.Println(err)
			return
		}
		if DB.GetDB().Create(&DB.Friend{Friend_1: Token_data.UserID, Friend_2: friend_id_int}).Error != nil {
			ctx.JSON(200, gin.H{"code": 0})
			return
		}
		ctx.JSON(200, gin.H{"code": 200})
	})

	r.DELETE("remove/:friend_Id", func(ctx *gin.Context) {
		friend_Id := ctx.Param("friend_Id")
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)

		friend_data := DB.Friend{}
		if DB.GetDB().Where("friend_1 = ? and friend_2 = ? or friend_1 = ? and friend_2 = ?", Token_data.UserID, friend_Id, friend_Id, Token_data.UserID).First(&friend_data).Error != nil {
			ctx.JSON(200, gin.H{"code": 0, "msg": "Not found friend"})
			return
		}

		DB.GetDB().Where("friend_1 = ? and friend_2 = ? or friend_1 = ? and friend_2 = ?", Token_data.UserID, friend_Id, friend_Id, Token_data.UserID).Delete(&DB.Friend{})
		ctx.JSON(200, gin.H{"code": 200})
	})
}
