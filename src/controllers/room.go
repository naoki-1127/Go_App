package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoomPage(c *gin.Context){
	c.HTML(http.StatusOK, "room.html", gin.H{
		"title": "Login",
	})
}