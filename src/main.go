package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	f, _ := os.Create("gin.log")
	err := godotenv.Load(".env")

	gin.DisableConsoleColor()

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	DBUSER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("PASSWORD")
	HOST := os.Getenv("HOST")
	PORTN := os.Getenv("PORTN")
	DBNAME := os.Getenv("DBNAME")

	dsn := DBUSER + ":" + PASSWORD + "@tcp(" + HOST + ":" + PORTN + ")/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Printf(dsn)
	_, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
