package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoomPage(c *gin.Context){
	Render(c, http.StatusOK, "room.html", gin.H{
		"title": "Login",
	})
}