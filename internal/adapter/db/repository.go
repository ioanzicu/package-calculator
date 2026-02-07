package db

import (
	dbsqlc "ignis/internal/adapter/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	dbsqlc.Querier
	Close()
}

type repository struct {
	*dbsqlc.Queries
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{
		Queries: dbsqlc.New(pool),
		pool:    pool,
	}
}

func (r *repository) Close() {
	r.pool.Close()
}
