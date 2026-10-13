package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hrs-o/docker-go/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	testDB := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(testDB),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("コンテナの作成に失敗: %s", err)
	}
	url, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	if err != nil {
		log.Fatalf("URLの取得に失敗: %s", err)
	}
	err = db.Migrate(url)

	if err != nil {
		log.Fatalf("migrationの実行に失敗: %s", err)
	}

	testPool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("poolに失敗: %s", err)
	}
	code := m.Run()
	testPool.Close()
	os.Exit(code)
}
