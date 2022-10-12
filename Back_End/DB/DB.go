package DB

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID           int    `gorm:"unique"`
	Name         string `gorm:"unique"`
	Email        string `gorm:"unique"`
	Password     string `gorm:"not_null"`
	Profile_Link string
}
type Friend struct {
	Friend_1 int  `gorm:"not_null"`
	Friend_2 int  `gorm:"not_null"`
	user_1   User `gorm:"foreignKey:friend_1"`
	user_2   User `gorm:"foreignKey:friend_2"`
}
type Room struct {
	ID        int    `gorm:"unique"`
	Admin     int    `gorm:"not_null"`
	Name      string `gorm:"not_null"`
	CreatedAt time.Time
	user      User `gorm:"foreignKey:admin"`
}
type Join_Room struct {
	ID      int  `gorm:"unique"`
	Room_Id int  `gorm:"not_null"`
	User_Id int  `gorm:"not_null"`
	user    User `gorm:"foreignKey:user_id"`
	room    Room `gorm:"foreignKey:Room_Id"`
}
type Msg struct {
	ID        int       `gorm:"unique"`
	User_Name string    `gorm:"not_null"`
	Join_id   int       `gorm:"not_null"`
	Data      string    `gorm:"not_null"`
	join_Room Join_Room `gorm:"foreignKey:join_id"`
}

var dB *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Ok! Env load")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", os.Getenv("DB_User_Name"), os.Getenv("DB_User_Password"), os.Getenv("DB_Createress"), os.Getenv("DB_Database_Name"))
	//dsn := "root:skymkey08@tcp(127.0.0.1:3306)/EnChat"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if err := db.AutoMigrate(&Msg{}, &User{}, &Friend{}, &Room{}, &Join_Room{}); err != nil {
		fmt.Println(err)
	}
	dB = db
	fmt.Println("DB rady")
}
func GetDB() *gorm.DB {
	return dB
}
