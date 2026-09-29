package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hrs-o/docker-go/db/sqlc"
	"github.com/stretchr/testify/assert"
)

// fakeUserRepository はDBを使わないインメモリの UserRepository 実装。
// FindOrCreateUser のロジックだけを高速・オフラインでテストするために使う。
type fakeUserRepository struct {
	users []User
}

func (f *fakeUserRepository) FindByProvider(ctx context.Context, provider, providerUserID string) (*User, error) {
	for i := range f.users {
		if f.users[i].Provider == provider && f.users[i].ProviderUserID == providerUserID {
			return &f.users[i], nil
		}
	}
	return nil, ErrUserNotFound
}

func (f *fakeUserRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (*User, error) {
	user := User{
		ID:             uuid.New(),
		Name:           arg.Name,
		Email:          arg.Email,
		Provider:       arg.Provider,
		ProviderUserID: arg.ProviderUserID,
	}
	f.users = append(f.users, user)
	return &user, nil
}

func TestFindOrCreateUser_CreatesNewUserWhenNotExists(t *testing.T) {
	repo := &fakeUserRepository{}

	user, err := FindOrCreateUser(context.Background(), repo, "google", "sub-123", "Taro", "taro@example.com", "member")

	assert.NoError(t, err)
	assert.Equal(t, "taro@example.com", user.Email)
	assert.Equal(t, "google", user.Provider)
	assert.Equal(t, "sub-123", user.ProviderUserID)
	assert.Len(t, repo.users, 1, "新規ユーザーが1件作成されるはず")
}

func TestFindOrCreateUser_ReturnsExistingUserWithoutCreating(t *testing.T) {
	existingID := uuid.New()
	repo := &fakeUserRepository{
		users: []User{
			{ID: existingID, Provider: "google", ProviderUserID: "sub-123", Name: "Taro", Email: "taro@example.com" ,Role: "member"},
		},
	}

	user, err := FindOrCreateUser(context.Background(), repo, "google", "sub-123", "Taro", "taro@example.com", "member")

	assert.NoError(t, err)
	assert.Equal(t, existingID, user.ID)
	assert.Len(t, repo.users, 1, "既存ユーザーがいる場合は新規作成しないはず")
}
