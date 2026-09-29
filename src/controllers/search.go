package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSearchPage(c *gin.Context){
	c.HTML(http.StatusOK, "search.html", gin.H{
		"title": "Login",
	})
}