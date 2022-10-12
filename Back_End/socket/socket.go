package socket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"EnChat/DB"
	"EnChat/Token"
	"EnChat/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Standard struct {
	Type string
	Data interface{}
}

func (s Standard) ToString() string {
	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
	}

	return string(data)
}

var Server Temp = Temp{}

func init() {
	Server.Init()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		_, err := Token.CheckJWT(r.URL.Query().Get("Token"))
		if err == nil {
			return true
		} else {
			return false
		}
	},
}

func Socket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	utils.ErrHandle(err)

	Token_string := ctx.Request.URL.Query().Get("Token")
	Token_data, _ := Token.CheckJWT(Token_string)

	fmt.Println("New Connection:", Token_data.UserID)
	Server.Session_Create(&Session{token: Token_string, conn: conn, userinfo: Token_data})

	join_room_list := []DB.Join_Room{}
	DB.GetDB().Where("user_id=?", Token_data.UserID).Find(&join_room_list)
	for _, data := range join_room_list {
		fmt.Println("Join Room:", data.Room_Id)
		Server.Sessions[Token_string].JoinRoom(data.Room_Id, data.ID)
	}

	go Server.Sessions[Token_string].Read()

	// for { // 로그용
	// 	fmt.Println(len(Server.Rooms), Server.Rooms)
	// 	for i, data := range Server.Rooms {
	// 		fmt.Print(i, len(data.conn_list), "\n")
	// 	}
	// 	time.Sleep(time.Second * 1)
	// }
}
