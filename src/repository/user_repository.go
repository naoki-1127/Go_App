package repository

import (
	"context"
	"errors"

	"github.com/hrs-o/docker-go/db/sqlc"
	"github.com/jackc/pgx/v5"
)

// ErrUserNotFound はプロバイダー情報に一致するユーザーが存在しないことを表す。
var ErrUserNotFound = errors.New("user not found")

// User は sqlc が identity.users から生成した型。
type User = sqlc.IdentityUser

// UserRepository はユーザーの検索・作成を抽象化する。
// DBを起動しない単体テストが書けるよう、FindOrCreateUser はこのインターフェースに依存する。
type UserRepository interface {
	FindByProvider(ctx context.Context, provider, providerUserID string) (*User, error)
	Create(ctx context.Context, arg sqlc.CreateUserParams) (*User, error)
}

// PgUserRepository は UserRepository の sqlc(pgx)実装。
type PgUserRepository struct {
	q *sqlc.Queries
}

func NewPgUserRepository(db sqlc.DBTX) *PgUserRepository {
	return &PgUserRepository{q: sqlc.New(db)}
}

func (r *PgUserRepository) FindByProvider(ctx context.Context, provider, providerUserID string) (*User, error) {
	user, err := r.q.GetUserByProvider(ctx, sqlc.GetUserByProviderParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PgUserRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (*User, error) {
	user, err := r.q.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindOrCreateUser は provider/providerUserID の組で既存ユーザーを探し、
// いなければ新規作成する。OIDCコールバックのように「同じアカウントで
// 再ログインしても新規作成せず、初回のみ作成する」動作を表す。
func FindOrCreateUser(ctx context.Context, repo UserRepository, provider, providerUserID, name, email string,role string) (*User, error) {
	user, err := repo.FindByProvider(ctx, provider, providerUserID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}
	return repo.Create(ctx, sqlc.CreateUserParams{
		Name:           name,
		Email:          email,
		Provider:       provider,
		ProviderUserID: providerUserID,
		Role: role,
	})
}
