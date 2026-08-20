package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserReporitory(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (userrepository *UserRepository) Create(ctx context.Context, username string, passwordHash string) error {

	sql_query := `INSERT INTO users (username, password_hash) VALUES(
	$1, $2
	)
	`
	_, err := userrepository.conn.Exec(ctx, sql_query, username, passwordHash)
	return err

}

func (userepository *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	sql_query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	var u User
	err := userepository.conn.QueryRow(ctx, sql_query, username).Scan(&u.ID, &u.Username, &u.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
