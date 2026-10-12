package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReservationPage(c *gin.Context){
	Render(c, http.StatusOK, "reservations.html", gin.H{
		"title": "Login",
	})
}