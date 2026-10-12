package controllers

import (
	"github.com/gin-gonic/gin"
)

// Render は c.HTML のラッパー。middleware.CurrentUser が gin.Context に
// セットした userName / isAdmin を、各ページが渡す data に自動で合成してテンプレートへ渡す。
// これにより各ハンドラは「自分はログイン状態を気にしない」ままヘッダーの
// ログイン/ログアウト表示を切り替えられる。
func Render(c *gin.Context, status int, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	if userName, ok := c.Get("userName"); ok {
		data["userName"] = userName
	}
	if isAdmin, ok := c.Get("isAdmin"); ok {
		data["isAdmin"] = isAdmin
	}
	c.HTML(status, name, data)
}
