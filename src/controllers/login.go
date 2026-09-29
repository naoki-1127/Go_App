package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	connect "github.com/hrs-o/docker-go/config"
	"github.com/hrs-o/docker-go/repository"
	"github.com/hrs-o/docker-go/session"
)

const (
	googleIssuer      = "https://accounts.google.com"
	oauthStateCookie  = "oauth_state"
	oauthStateCookieN = 300 // 5分
)

var (
	googleProvider     *oidc.Provider
	googleProviderOnce sync.Once
	googleProviderErr  error
)

// getGoogleProvider は Google の OIDC discovery エンドポイントから設定を取得する。
// discoveryは初回のみ行い、以降はキャッシュした結果を使い回す。
func getGoogleProvider(ctx context.Context) (*oidc.Provider, error) {
	googleProviderOnce.Do(func() {
		googleProvider, googleProviderErr = oidc.NewProvider(ctx, googleIssuer)
	})
	return googleProvider, googleProviderErr
}

func googleOAuthConfig(provider *oidc.Provider) oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		// "openid" は OpenID Connect フローに必要なスコープ。
		Scopes: []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail},
	}
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Login は /login への GET リクエストに対し、ログインページ(login.html)を表示する。
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// Login は /login への GET リクエストに対し、ログインページ(login.html)を表示する。
func AdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", gin.H{
		"title": "Login",
	})
}

// 利用者ログインページからの動線
func MemberGoogleLogin(c *gin.Context) {
	role := "member"
	GoogleLogin(c,role)
}

// 管理者ログインページからの動線
func AdminGoogleLogin(c *gin.Context) {
	role := "admin"
	GoogleLogin(c,role)
}

// GoogleLogin は「Googleでサインイン」ボタン押下時に呼ばれ、
// Google の認可エンドポイントへリダイレクトする。
func GoogleLogin(c *gin.Context,role string){
	provider, err := getGoogleProvider(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "OIDCプロバイダーの初期化に失敗しました: %v", err)
		return
	}

	state, err := randomState()
	if err != nil {
		c.String(http.StatusInternalServerError, "stateの生成に失敗しました")
		return
	}
	// CSRF対策: コールバック時にこの値と突き合わせる。
	c.SetCookie(oauthStateCookie, state, oauthStateCookieN, "/", "", false, true)
	c.SetCookie("role", role, oauthStateCookieN, "/", "", false, true)

	config := googleOAuthConfig(provider)
	c.Redirect(http.StatusFound, config.AuthCodeURL(state))
}

// GoogleCallback は Google からのリダイレクトを受け取り、
// 認可コードをトークンに交換したうえで ID トークンを検証する。
func GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	stateCookie, err := c.Cookie(oauthStateCookie)
	role, err := c.Cookie("role")
	if err != nil || stateCookie == "" || stateCookie != c.Query("state") {
		c.String(http.StatusBadRequest, "不正なstateです")
		return
	}
	c.SetCookie(oauthStateCookie, "", -1, "/", "", false, true)

	provider, err := getGoogleProvider(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "OIDCプロバイダーの初期化に失敗しました: %v", err)
		return
	}
	config := googleOAuthConfig(provider)

	token, err := config.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.String(http.StatusUnauthorized, "トークン交換に失敗しました: %v", err)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusUnauthorized, "id_tokenが取得できませんでした")
		return
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.String(http.StatusUnauthorized, "IDトークンの検証に失敗しました: %v", err)
		return
	}

	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		c.String(http.StatusInternalServerError, "claimsの取得に失敗しました: %v", err)
		return
	}

	repo := repository.NewPgUserRepository(connect.Pool())
	user, err := repository.FindOrCreateUser(ctx, repo, "google", claims.Sub, claims.Name, claims.Email,role)
	if err != nil {
		c.String(http.StatusInternalServerError, "ユーザーの作成/取得に失敗しました: %v", err)
		return
	}

	session.SetUserID(c, user.ID)
	c.Redirect(http.StatusFound, "/")
}

// Logout はログインセッションを破棄する。
func Logout(c *gin.Context) {
	session.Clear(c)
	c.Redirect(http.StatusFound, "/login")
}
