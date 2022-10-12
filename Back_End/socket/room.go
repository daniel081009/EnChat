package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	Room_Id   int
	conn_list map[int]*websocket.Conn
	Mutex     sync.Mutex
}

func (room *Room) init() {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.conn_list = map[int]*websocket.Conn{}
}

func (room *Room) Delete_Conn(Join_Id int) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	delete(room.conn_list, Join_Id)
}

func (room *Room) Create_Conn(Join_Id int, conn *websocket.Conn) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	room.conn_list[Join_Id] = conn
}

func (room *Room) Send(data []byte) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	for _, conn := range room.conn_list {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}
