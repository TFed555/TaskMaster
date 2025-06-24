package postgresql

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err:=sql.Open("pgx", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	//вынести в отдельные миграции
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL primary key,
		login nvarchar(64) not null,
		email varchar(255) not null unique,
		password varchar(255) not null unique,
		created_at timestamp with time zone default now())`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec();

	stmt1, err := db.Prepare(`CREATE TABLE IF NOT EXISTS refresh_tokens(
		id SERIAL primary key,
		user_id int REFERENCES ON users(id) on delete cascade,
		token varchar(512) unique not null,
		expires_at TIMESTAMP with time zone not null,
		created_at timestamp with time zone default now()
	)`)
	
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt1.Exec();

	return &Storage{db: db}, nil
	
}