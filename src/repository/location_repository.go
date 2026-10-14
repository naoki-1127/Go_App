package repository

import (
	"context"

	"github.com/hrs-o/docker-go/db/sqlc"
)

type Location = sqlc.FacilityLocation

type LocationRepository interface {
	Create(ctx context.Context, arg sqlc.RegistLocationParams) (*Location, error)
	List(ctx context.Context) ([]Location, error)
}

// PgLocationRepository は LocationRepository の sqlc(pgx)実装。
type PgLocationRepository struct {
	q *sqlc.Queries
}

var _ LocationRepository = (*PgLocationRepository)(nil)

func NewPgLocationRepository(db sqlc.DBTX) *PgLocationRepository {
	return &PgLocationRepository{q: sqlc.New(db)}
}

func (r *PgLocationRepository) Create(ctx context.Context, arg sqlc.RegistLocationParams) (*Location, error) {
	location, err := r.q.RegistLocation(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *PgLocationRepository) List(ctx context.Context) ([]Location, error) {
	return r.q.GetLocations(ctx)
}
