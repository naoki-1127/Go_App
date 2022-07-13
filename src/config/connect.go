package connect

import (
	"os"
)

func Dbconnect() string {
	DBUSER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("PASSWORD")
	HOST := os.Getenv("HOST")
	PORTN := os.Getenv("PORTN")
	DBNAME := os.Getenv("DBNAME")

	dsn := DBUSER + ":" + PASSWORD + "@tcp(" + HOST + ":" + PORTN + ")/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn
}
