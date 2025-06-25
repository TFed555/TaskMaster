CREATE TABLE IF NOT EXISTS refresh_tokens(
		id SERIAL primary key,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token varchar(255) not null,
		expires_at TIMESTAMP with time zone not null)
