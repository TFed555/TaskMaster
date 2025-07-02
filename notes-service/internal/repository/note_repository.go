package repository

import (
	"database/sql"
	"fmt"
	"log"
	"notes-service/internal/models"
	_"time"

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

func (n *NotesRepository) GetTodos(userId uint, createdAt string, filter string, limit string, offset string) (*models.Todo, error) {
	const op = "repository.user_repository.GetTodos"

	// var sqlFilter string

	// if filter == "after" {
	// 	sqlFilter = ">"
	// } else {
	// 	sqlFilter = "<"
	// }

	// parserCreatedAt, err := time.Parse(time.RFC3339, createdAt)

	log.Print(createdAt)

	query := (`
		SELECT id, title, priority, category, description, createdat, completedat, userid FROM notes.todos
		WHERE userid = $1 AND createdat > $2::TIMESTAMPTZ`)

	var task models.Todo
	// createdAtFormat, err := time.(createdAt, )

	err := n.db.QueryRowx(query, userId, createdAt).StructScan(&task)

	log.Printf("Returned field of task: %d, %s, %s", task.ID, task.Title, task.Priority)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("%s: task not found", op)
        }
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return &task, nil
}