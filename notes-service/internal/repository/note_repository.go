package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"notes-service/internal/models"
	_"strings"
	_ "time"

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

func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

func (n *NotesRepository) GetTodos(userId uint, urlParams url.Values) ([]models.Todo, error) {
	const op = "repository.user_repository.GetTodos"
	args := []interface{}{userId}
	// parserCreatedAt, err := time.Parse(time.RFC3339, createdAt)
	var sqlFilter string
	if filter := urlParams.Get("filter"); filter != "" {
		if filter == "after" {
			sqlFilter = ">"
		} else {
			sqlFilter = "<"
		}
	} else {
		sqlFilter = "="
	}

	log.Printf("Filter: %s", sqlFilter)

	query := (`
		SELECT id, title, priority, category, description, createdat, completedat, userid FROM notes.todos
		WHERE userid = $1`)

	todos := []models.Todo{}
	counter := 1
	if createdAt := urlParams.Get("createdAt"); createdAt != "" {
		counter ++
		query += fmt.Sprintf(" AND createdat %s $%d::TIMESTAMPTZ ", sqlFilter, counter)
		args = append(args, createdAt)
		log.Printf("Data: %s", createdAt)
	}

	if offset := urlParams.Get("offset"); offset != "" {
		counter ++
		query += fmt.Sprintf(" OFFSET $%d", counter)
		args = append(args, offset)
		log.Printf("Offset: %s", offset)
	}

	if limit := urlParams.Get("limit"); limit != "" {
		counter ++
		query += fmt.Sprintf(" LIMIT $%d ", counter)
		args = append(args, limit)
		log.Printf("Limit: %s", limit)
	}

	err := n.db.Select(&todos, query, args...)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("%s: task not found", op)
        }
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return todos, nil
}

func (n *NotesRepository) CreateTodo(userID uint, title string, priority string, description string,
    category string, createdAt string, completedAt string) (int, error) {
    const op = "repository.user_repository.CreateTodo"

    query := `INSERT INTO notes.todos (userid, title, priority, description, category, createdat`
    values := `VALUES ($1, $2, $3, $4, $5, $6`
    args := []interface{}{userID, title, priority, description, category, createdAt}
    count := 6

    if completedAt != "" {
        query += `, completedat`
        count++

        values += fmt.Sprintf(`, $%d`, count)
        args = append(args, completedAt)
    }
    query +=` )` + values + ` ) RETURNING id`
    var id int
    err := n.db.QueryRow(query, args...).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }
    return id, nil
}
