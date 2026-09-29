package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReservationPage(c *gin.Context){
	c.HTML(http.StatusOK, "reservations.html", gin.H{
		"title": "Login",
	})
}