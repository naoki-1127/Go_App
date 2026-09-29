package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

// レスポンスに乗ったSet-Cookieを、次のリクエストのCookieとして引き継ぐ。
func carryCookies(w *httptest.ResponseRecorder, req *http.Request) {
	for _, ck := range w.Result().Cookies() {
		req.AddCookie(ck)
	}
}

func TestSetUserID_ThenUserID_RoundTrips(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")

	userID := uuid.New()
	c, w := newTestContext()
	SetUserID(c, userID)

	c2, _ := newTestContext()
	carryCookies(w, c2.Request)

	id, ok := UserID(c2)
	assert.True(t, ok)
	assert.Equal(t, userID, id)
}

func TestUserID_NoCookie(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")

	c, _ := newTestContext()
	_, ok := UserID(c)
	assert.False(t, ok)
}

func TestUserID_TamperedValueIsRejected(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")

	c, w := newTestContext()
	SetUserID(c, uuid.New())

	var original string
	for _, ck := range w.Result().Cookies() {
		if ck.Name == cookieName {
			original = ck.Value
		}
	}
	assert.NotEmpty(t, original)

	// 署名はそのままに、ユーザーID部分だけ不正に書き換える。
	_, sig, _ := strings.Cut(original, ".")
	tampered := uuid.New().String() + "." + sig

	c2, _ := newTestContext()
	c2.Request.AddCookie(&http.Cookie{Name: cookieName, Value: tampered})

	_, ok := UserID(c2)
	assert.False(t, ok, "署名が一致しない改ざん済みcookieは拒否されるはず")
}

func TestClear_RemovesCookie(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-secret")

	c, _ := newTestContext()
	Clear(c)

	found := false
	for _, ck := range c.Writer.Header().Values("Set-Cookie") {
		if len(ck) > 0 {
			found = true
		}
	}
	assert.True(t, found, "Clearはcookie削除のSet-Cookieを発行するはず")
}
