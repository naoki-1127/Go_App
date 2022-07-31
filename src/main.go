package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	migrate "github.com/hrs-o/docker-go/models/db"
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
	migrate.Open()
	//route読み込み
	router.Router()

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}
