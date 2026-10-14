package repository

import (
	"context"
	"testing"
	"time"

	"github.com/hrs-o/docker-go/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	address := "岐阜県岐阜市1-2-3"
	arg := sqlc.RegistLocationParams{
		Name:     "サンプル",
		Address:  &address,
		Timezone: "Asia/Tokyo",
		OpenTime: pgtype.Time{
			Microseconds: int64(9 * time.Hour / time.Microsecond), // 09:00
			Valid:        true,
		},
		CloseTime: pgtype.Time{
			Microseconds: int64((18*time.Hour + 30*time.Minute) / time.Microsecond), // 18:30
			Valid:        true,
		},
	}
	q := NewPgLocationRepository(testPool)
	location, err := q.Create(context.Background(), arg)
	require.NoError(t, err)
	assert.Equal(t, location.Name, "サンプル")
}
