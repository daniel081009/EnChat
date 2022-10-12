package User

import (
	"EnChat/DB"
	db "EnChat/DB"
	"EnChat/Token"
	"crypto/sha512"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) {
	r.POST("login", func(ctx *gin.Context) {
		type body struct {
			Email    string `json:"Email" binding:"required"`
			Password string `json:"Password" binding:"required"`
		}
		res_body := body{}
		ctx.BindJSON(&res_body)
		db := DB.GetDB()
		User_Data := DB.User{}
		En := fmt.Sprintf("%x", sha512.Sum512([]byte(res_body.Password)))
		if db.Where("email = ? and password = ?", res_body.Email, En).First(&User_Data).Error != nil {
			ctx.JSON(200, gin.H{"code": 0})
			return
		}
		Token, _ := Token.CreateJWT(User_Data.ID, User_Data.Name, User_Data.Email)
		ctx.SetCookie("Token", Token, 60*60*24, "/", "*", true, true)
		ctx.JSON(200, gin.H{"code": 200})
	})

	r.POST("register", func(ctx *gin.Context) {
		type body struct {
			Name     string `json:"Name" binding:"required"`
			Email    string `json:"Email" binding:"required"`
			Password string `json:"Password" binding:"required"`
		}
		res_body := body{}
		ctx.BindJSON(&res_body)
		db := db.GetDB()
		En := fmt.Sprintf("%x", sha512.Sum512([]byte(res_body.Password)))
		d := db.Create(&DB.User{Name: res_body.Name, Email: res_body.Email, Password: En, Profile_Link: "Def"})
		if d.Error != nil {
			fmt.Println(d.Error)
			ctx.JSON(200, gin.H{
				"code": 0,
				"msg":  "Err sorry",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "OK!",
		})
	})
	r.GET("me", Token.AuthMiddleWare, func(ctx *gin.Context) {
		token, _ := ctx.Cookie("Token")
		Token_data, _ := Token.CheckJWT(token)
		res := struct {
			Id          int
			Name        string
			Email       string
			Friend_list []DB.Friend
			Room_list   []DB.Room
			Join_list   []DB.Join_Room
		}{}
		res.Id = Token_data.UserID
		res.Name = Token_data.Name
		res.Email = Token_data.Email

		D := DB.GetDB()
		D.Where("friend_1 = ? or friend_2= ?", res.Id, res.Id).Find(&res.Friend_list)
		D.Where("admin = ?", res.Id).Find(&res.Room_list)
		D.Where("user_id= ?", res.Id).Find(&res.Join_list)

		ctx.JSON(200, res)
	})
}
