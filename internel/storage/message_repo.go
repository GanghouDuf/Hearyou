// internal/storage/message_repo.go
package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Message struct {
	ID        int
	AuthorID  int
	RoomID    int
	Author    string
	Payload   string
	CreatedAt time.Time
}

type MessageRepository struct {
	conn *pgx.Conn
}

func NewMessageRepository(conn *pgx.Conn) *MessageRepository {
	return &MessageRepository{conn: conn}
}

func (r *MessageRepository) Save(ctx context.Context, authorID, roomID int, payload string) error {
	sql := `INSERT INTO messages (author_id, room_id, payload) VALUES ($1, $2, $3)`
	_, err := r.conn.Exec(ctx, sql, authorID, roomID, payload)
	return err
}

func (r *MessageRepository) GetHistory(ctx context.Context, roomID, limit int) ([]Message, error) {
	sql := `
		SELECT messages.id, messages.author_id, users.username, messages.room_id, messages.payload, messages.created_at
		FROM messages
		JOIN users ON messages.author_id = users.id
		WHERE messages.room_id = $1
		ORDER BY messages.created_at DESC
		LIMIT $2
	`

	rows, err := r.conn.Query(ctx, sql, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.AuthorID, &m.Author, &m.RoomID, &m.Payload, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
