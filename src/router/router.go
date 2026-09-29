package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrs-o/docker-go/controllers"
)

// SetupRouter はルーティング定義だけを行い、*gin.Engine を返す。
// サーバー起動(r.Run)を含まないため、httptest から直接ハンドリングをテストできる。
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/resources", "./resources")
	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	r.GET("/admin/login", controllers.AdminLogin)
	r.GET("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.GET("/auth/google/login", controllers.MemberGoogleLogin)
	r.GET("/admin/auth/google/login", controllers.AdminGoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	r.GET("/reservations",controllers.GetReservationPage)

	r.GET("/search",controllers.GetSearchPage)

	r.GET("/rooms",controllers.GetRoomPage)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

func Router() {
	r := SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
