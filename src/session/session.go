// Package session は、ログイン中のユーザーIDを保持するための
// 最小限の署名付きCookieセッションを提供する。
//
// gin-contrib/sessions 等のライブラリは導入せず、標準ライブラリの
// crypto/hmac だけで実装している(セッションIDを保存するだけの用途に対して
// 外部ライブラリが依存関係として重すぎるため)。
package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	cookieName    = "app_session"
	cookieMaxAgeS = 60 * 60 * 24 * 30 // 30日
)

func secret() []byte {
	return []byte(os.Getenv("SESSION_SECRET"))
}

func sign(value string) string {
	mac := hmac.New(sha256.New, secret())
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// SetUserID はログイン済みユーザーIDを、改ざん検知用の署名を付けてcookieに保存する。
func SetUserID(c *gin.Context, userID uuid.UUID) {
	value := userID.String()
	c.SetCookie(cookieName, value+"."+sign(value), cookieMaxAgeS, "/", "", false, true)
}

// UserID はcookieからログイン中のユーザーIDを取り出す。
// 未ログイン、cookieが無い、または署名が不正な場合は ok=false を返す。
func UserID(c *gin.Context) (id uuid.UUID, ok bool) {
	raw, err := c.Cookie(cookieName)
	if err != nil || raw == "" {
		return uuid.Nil, false
	}

	value, sig, found := strings.Cut(raw, ".")
	if !found {
		return uuid.Nil, false
	}
	if !hmac.Equal([]byte(sig), []byte(sign(value))) {
		return uuid.Nil, false
	}

	id, err = uuid.Parse(value)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// Clear はログアウト時にセッションcookieを破棄する。
func Clear(c *gin.Context) {
	c.SetCookie(cookieName, "", -1, "/", "", false, true)
}
