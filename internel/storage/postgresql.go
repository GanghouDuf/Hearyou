package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
)

func Connect(databaseUrl string, ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {

		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
