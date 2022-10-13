package socket

import (
	"EnChat/DB"
	"fmt"
)

func HandMsg(msg *Standard, Sess *Session) {
	switch msg.Type {
	case "Get_Join_Room_List":
		Sess.Send_Msg(Standard{Type: msg.Type + "_Res", Data: fmt.Sprint(Sess.Get_Join_Rooms())})
	case "Join_Room":
		temp := struct {
			Room_Id int
			Join_Id int
		}{}
		msg.StringTo(&temp)
		if DB.GetDB().Where("user_id= ? and room_id = ? and id = ?", Sess.userinfo.UserID, temp.Room_Id, temp.Join_Id).First(&DB.Join_Room{}).Error != nil {
			return
		}
		Sess.JoinRoom(temp.Room_Id, temp.Join_Id)
		Sess.Send_Msg(Standard{Type: msg.Type + "_Res", Data: "OK!"})
	}
}
