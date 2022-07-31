package connect

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Dbconnect() *gorm.DB {
	DBUSER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("PASSWORD")
	HOST := os.Getenv("HOST")
	PORTN := os.Getenv("PORTN")
	DBNAME := os.Getenv("DBNAME")

	dsn := DBUSER + ":" + PASSWORD + "@tcp(" + HOST + ":" + PORTN + ")/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	} else {
		log.Printf("DB connect success")
	}
	return db
}
