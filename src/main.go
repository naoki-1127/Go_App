package main

import (
	"fmt"
	"io"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	connect "github.com/hrs-o/docker-go/config"
	"github.com/hrs-o/docker-go/db"
	"github.com/hrs-o/docker-go/router"
	"github.com/joho/godotenv"
)

func main() {

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DisableConsoleColor()
	//env 読み込み
	loadEnv()
	//migrate
	if err := db.Migrate(connect.DatabaseURL()); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	//route読み込み
	router.Router()

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}
