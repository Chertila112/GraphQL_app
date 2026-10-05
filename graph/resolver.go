package graph

import (
	"github.com/Chertila112/server/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resolver struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}
