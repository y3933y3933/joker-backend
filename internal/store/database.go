package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/y3933y3933/joker/internal/db/sqlc"
)

func Open(dbUrl string) (*pgxpool.Pool, *sqlc.Queries, error) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("db: open %w", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("db: open %w", err)
	}
	queries := sqlc.New(dbpool)

	fmt.Println("Connected to Database...")
	return dbpool, queries, nil
}
