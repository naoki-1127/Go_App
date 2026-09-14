package connect

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Dbconnect() *gorm.DB {
	DBUSER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("PASSWORD")
	HOST := os.Getenv("HOST")
	PORTN := os.Getenv("PORTN")
	DBNAME := os.Getenv("DBNAME")

	dsn := "host=" + HOST + " user=" + DBUSER + " password=" + PASSWORD + " dbname=" + DBNAME + " port=" + PORTN + " sslmode=disable TimeZone=Asia/Tokyo"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	} else {
		log.Printf("DB connect success")
	}
	return db
}
