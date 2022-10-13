package socket

import (
	"EnChat/Token"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Session struct {
	token      string
	conn       *websocket.Conn
	userinfo   Token.AuthTokenClaims
	Join_Rooms map[int]int
}

func (Sess *Session) init() {
	Sess.Join_Rooms = map[int]int{}
}

func (Sess *Session) close() {
	Sess.conn.Close()
	Server.Session_Delete(Sess.token)

	for key := range Sess.Join_Rooms {
		Sess.LeaveRoom(key)
	}
}

func (Sess *Session) JoinRoom(Room_Id int, Join_Id int) {
	Sess.Join_Rooms[Room_Id] = Join_Id

	if _, Exist := Server.Rooms[Room_Id]; !Exist {
		Server.Room_Create(Room_Id)
	} else if _, Exist := Server.Rooms[Room_Id].conn_list[Join_Id]; Exist {
		return
	}
	Server.Rooms[Room_Id].Create_Conn(Join_Id, Sess.conn)
}

func (Sess *Session) Remove_JoinRoom(Room_Id int) {
	delete(Sess.Join_Rooms, Room_Id)
}

func (Sess *Session) LeaveRoom(Room_Id int) {
	if _, Exist := Server.Rooms[Room_Id]; !Exist {
		return
	}

	if len(Server.Rooms[Room_Id].conn_list) <= 1 {
		fmt.Println("Room Delete")
		Server.Room_Delete(Room_Id)
	} else {
		fmt.Println("Room Conn Delete")
		fmt.Println(len(Server.Rooms[Room_Id].conn_list))
		Server.Rooms[Room_Id].Delete_Conn(Sess.Join_Rooms[Room_Id])
	}

	Sess.Remove_JoinRoom(Room_Id)
}
func (Sess *Session) Get_Join_Rooms() []int {
	rooms := []int{}
	for key := range Sess.Join_Rooms {
		rooms = append(rooms, key)
	}
	return rooms
}

func (Sess *Session) read() {
	defer Sess.close()
	for {
		_, data, err := Sess.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		req := Standard{}
		err = json.Unmarshal(data, &req)
		if err != nil {
			continue
		}

		fmt.Println(req)
		HandMsg(&req, Sess)
	}
}
func (p *Session) Send_Msg(data Standard) {
	err := p.conn.WriteMessage(websocket.TextMessage, []byte(data.ToString()))
	if err != nil {
		fmt.Println(err)
		p.close()
	}
}
