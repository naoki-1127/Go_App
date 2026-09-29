package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetReservationPage(t *testing.T){
	t.Chdir("..")

	r := gin.New()
	r.LoadHTMLGlob("views/*")
	r.GET("/reservations", GetReservationPage)

	req := httptest.NewRequest(http.MethodGet, "/reservations", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}