package repository

import (
	"database/sql"
	"fmt"
	"notes-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type NotesRepository struct {
	db *sqlx.DB
}

func NewNotesRepository(db *sqlx.DB) *NotesRepository {
	return &NotesRepository{
		db: db,
	}
}

func (n *NotesRepository) GetTodos(userId int) (*models.Todo, error) {
	const op = "repository.user_repository.GetTodos"

	query := (`
		SELECT * FROM notes.users
		WHERE id = $1
		LIMIT 1`)

	var task models.Todo

	err := n.db.QueryRowx(query, userId).StructScan(&task)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("%s: task not found", op)
        }
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return &task, nil
}