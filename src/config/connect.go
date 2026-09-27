package connect

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

// DatabaseURL は .env の接続情報から postgres:// 形式の接続URLを組み立てる。
// pgx と golang-migrate の両方で使う。
func DatabaseURL() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(os.Getenv("DB_USER"), os.Getenv("PASSWORD")),
		Host:   net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORTN")),
		Path:   os.Getenv("DBNAME"),
	}
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "Asia/Tokyo")
	u.RawQuery = q.Encode()
	return u.String()
}

// Pool はDBのコネクションプールを返す。接続確立は初回のみ行い、以降は
// キャッシュしたプールを使い回す(コントローラー等、リクエストのたびに
// 呼び出される箇所からでも安全に呼べるようにするため)。
func Pool() *pgxpool.Pool {
	poolOnce.Do(func() {
		pool, poolErr = pgxpool.New(context.Background(), DatabaseURL())
		if poolErr != nil {
			return
		}
		log.Printf("DB connect success")
	})
	if poolErr != nil {
		panic(poolErr.Error())
	}
	return pool
}
