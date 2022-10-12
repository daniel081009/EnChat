package socket

import (
	"sync"
)

type Temp struct {
	Rooms    map[int]*Room
	Sessions map[string]*Session
	Mutex    *sync.Mutex
}

func (temp *Temp) Init() {
	temp.Rooms = map[int]*Room{}
	temp.Sessions = map[string]*Session{}
	temp.Mutex = &sync.Mutex{}
}

func (temp *Temp) Room_Create(Room_id int) {
	Room := &Room{Room_Id: Room_id}
	Room.init()

	temp.Mutex.Lock()
	defer temp.Mutex.Unlock()
	temp.Rooms[Room.Room_Id] = Room
}

func (temp *Temp) Session_Create(Sess *Session) {
	temp.Mutex.Lock()
	defer temp.Mutex.Unlock()

	Sess.init()
	temp.Sessions[Sess.token] = Sess
}

func (temp *Temp) Session_Delete(token string) {
	temp.Mutex.Lock()
	defer temp.Mutex.Unlock()

	delete(temp.Sessions, token)
}
func (temp *Temp) Room_Delete(Room_Id int) {
	temp.Mutex.Lock()
	defer temp.Mutex.Unlock()

	delete(temp.Rooms, Room_Id)
}
