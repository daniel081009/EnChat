package main

import (
	"EnChat/socket"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

const (
	Createress = "localhost"
	port       = "1009"
)

func Send_Msg(conn *websocket.Conn, data socket.Standard) {
	if conn == nil {
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(data.ToString()))
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?Token=%s", Createress, port, os.Args[1]), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conn.RemoteAddr().String())
	Send_Msg(conn, socket.Standard{Type: "JoinRoom", Data: 1})
	for {
	}
	defer conn.Close()
}
