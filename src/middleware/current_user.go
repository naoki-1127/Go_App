// Package middleware は、gin のルーティングに差し込む横断的な処理をまとめる。
package middleware

import (
	"github.com/gin-gonic/gin"

	connect "github.com/hrs-o/docker-go/config"
	"github.com/hrs-o/docker-go/repository"
	"github.com/hrs-o/docker-go/session"
)

// CurrentUser はログイン中のユーザー情報を gin.Context に乗せる。
// _layout.html の header が参照する userName / isAdmin を、
// 画面を描画するハンドラ側で毎回書かなくて済むようにするための共通処理。
// 未ログイン、またはユーザー取得に失敗した場合は何もセットしない
// (= テンプレート側では未ログイン扱いになる)。
func CurrentUser(c *gin.Context) {
	userID, ok := session.UserID(c)
	if !ok {
		c.Next()
		return
	}

	repo := repository.NewPgUserRepository(connect.Pool())
	user, err := repo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.Next()
		return
	}

	c.Set("userName", user.Name)
	c.Set("isAdmin", user.Role == "admin")
	c.Next()
}
