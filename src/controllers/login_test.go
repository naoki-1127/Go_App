package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestLogin は、/login にアクセスするとログインページ(login.html)が
// 表示されることを確認するテスト。「Googleでサインイン」ボタンはこのページから
// 別ルート(/auth/google/login)に遷移する形なので、ここでは200 OKで
// ページが返ることだけを見る。
func TestLogin(t *testing.T) {
	// Login() は c.HTML でテンプレートを描画するため、
	// LoadHTMLGlob 前提のパス(src/ 直下)にカレントディレクトリを合わせる。
	t.Chdir("..")

	r := gin.New()
	r.LoadHTMLGlob("views/*")
	r.GET("/login", Login)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "/auth/google/login", "「Googleでサインイン」ボタンのリンク先を含むはず")
}

// TestGoogleLogin は、/auth/google/login にアクセスすると Google の
// 認可エンドポイントへリダイレクトされることを確認するテスト。
func TestGoogleLogin(t *testing.T) {
	r := gin.New()
	r.GET("/auth/google/login", MemberGoogleLogin)

	req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code, "Googleの認可エンドポイントへリダイレクトするはず")

	location := w.Header().Get("Location")
	assert.Contains(t, location, "accounts.google.com")
	assert.Contains(t, location, "response_type=code")
	assert.Contains(t, location, "scope="+oidc.ScopeOpenID)

	// CSRF対策のstate cookieが発行されていること
	cookies := w.Result().Cookies()
	found := false
	role := "member"
	for _, ck := range cookies {
		if ck.Name == oauthStateCookie {
			found = true
		}
		if ck.Name == role {
			role = "member"
		}
	}
	assert.Equal(t, role, "member")
	assert.True(t, found, "state用のcookieが発行されるはず")
}

// TestOidc は、Google の OIDC discovery エンドポイントから設定を取得し、
// 認可URL・IDトークンVerifierを組み立てられることを確認するテスト。
//
// Google の "/.well-known/openid-configuration" へ実際にHTTPアクセスするため、
// ネットワーク接続が必要。client_id/secret はダミー値で良い(discoveryと
// AuthCodeURLの組み立てには実際の認証情報は不要なため)。
func TestOidc(t *testing.T) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		t.Fatalf("Googleのdiscoveryエンドポイントから設定を取得できなかった: %v", err)
	}

	const (
		clientID     = "dummy-client-id"
		clientSecret = "dummy-client-secret"
		redirectURL  = "http://localhost:3000/auth/google/callback"
	)

	// OpenID Connect 対応の OAuth2 クライアントを設定します。
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// ディスカバリにより OAuth2 エンドポイントが返されます。
		Endpoint: provider.Endpoint(),
		// "openid" は OpenID Connect フローに必要なスコープです。
		Scopes: []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail},
	}

	authURL := oauth2Config.AuthCodeURL("state")
	assert.Contains(t, authURL, "accounts.google.com", "GoogleのOAuth2認可エンドポイントを指しているはず")
	assert.Contains(t, authURL, "client_id="+clientID)
	assert.Contains(t, authURL, "scope=openid")

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	assert.NotNil(t, idTokenVerifier, "IDトークンVerifierが生成されるはず")
}
