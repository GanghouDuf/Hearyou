CREATE TABLE messages (
	id SERIAL PRIMARY KEY,
	author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	payload TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_messages_created_at ON messages(created_at);