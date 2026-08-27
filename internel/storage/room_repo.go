package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Room struct {
	ID        int
	Name      string
	CreatedBy int
}

type RoomRepository struct {
	conn *pgx.Conn
}

func NewRoomRepository(conn *pgx.Conn) *RoomRepository {
	return &RoomRepository{conn: conn}
}

// Create создаёт новую комнату и возвращает её ID
func (r *RoomRepository) Create(ctx context.Context, name string, createdBy int) (int, error) {
	sql := `INSERT INTO rooms (name, created_by) VALUES ($1, $2) RETURNING id`

	var id int
	err := r.conn.QueryRow(ctx, sql, name, createdBy).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByName находит комнату по имени
func (r *RoomRepository) GetByName(ctx context.Context, name string) (*Room, error) {
	sql := `SELECT id, name, created_by FROM rooms WHERE name = $1`

	var room Room
	err := r.conn.QueryRow(ctx, sql, name).Scan(&room.ID, &room.Name, &room.CreatedBy)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

// List возвращает все существующие комнаты
func (r *RoomRepository) List(ctx context.Context) ([]Room, error) {
	sql := `SELECT id, name, created_by FROM rooms ORDER BY created_at`

	rows, err := r.conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.CreatedBy); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
