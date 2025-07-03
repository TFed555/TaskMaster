package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"notes-service/internal/models"
	"strings"
	_ "strings"
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

// func contains(slice []string, item string) bool {
//     for _, v := range slice {
//         if v == item {
//             return true
//         }
//     }
//     return false
// }

func (n *NotesRepository) GetTodos(urlParams url.Values, tableName string) ([]models.Todo, error) {
	const op = "repository.user_repository.GetTodos"
	args := []interface{}{}
	// args := []interface{}{userId}
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

	// query := (`
	// 	SELECT id, title, priority, category, description, createdat, completedat, userid FROM notes.todos
	// 	WHERE userid = $1`)

	query := fmt.Sprintf(`
		SELECT id, title, priority, category, description, createdat, completedat, userid FROM notes.%s `, tableName)

	todos := []models.Todo{}
	counter := 0

	st := "WHERE"

	if userID:=urlParams.Get("userID"); userID != "" {
		counter ++
		st = "AND"
		query += fmt.Sprintf("WHERE userid = $%d ", counter)
		args = append(args, userID)
		log.Printf("UserID: %s", userID)
	}

	if dataToFilter := urlParams.Get("createdAt"); dataToFilter != "" {
		counter ++

		query += fmt.Sprintf(" %s createdat %s $%d::TIMESTAMPTZ ", st, sqlFilter, counter)
		args = append(args, dataToFilter)
		log.Printf("Date: %s", dataToFilter)
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

func (n *NotesRepository) GetTodoByID(ID int) (*models.Todo, error) {
	const op = "repository.notes_repository.GetTodoByID"
	query := `SELECT id, title, priority, category, description, createdat, completedat, userid
			FROM notes.todos WHERE id = $1 LIMIT 1`
	todo := models.Todo{}
	err := n.db.QueryRowx(query, ID).StructScan(&todo)
	if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
	return &todo, nil
}

func (n *NotesRepository) CreateTodo(userID uint, title string, priority string, description string,
    category string, createdAt string, completedAt string) (int, error) {
    const op = "repository.notes_repository.CreateTodo"

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

func (n *NotesRepository) ArchiveTodo(ID int) (int, error) {
    const op = "repository.user_repository.ArchiveTodo"

	todo, err := n.GetTodoByID(ID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

    query := `INSERT INTO notes.archived_todos (userid, title, priority, description, category, createdat`
    values := `VALUES ($1, $2, $3, $4, $5, $6`
    args := []interface{}{todo.UserId, todo.Title, todo.Priority, todo.Description, todo.Category, todo.CreatedAt}
    count := 6

    if *todo.CompletedAt != "" {
        query += `, completedat`
        count++

        values += fmt.Sprintf(`, $%d`, count)
        args = append(args, todo.CompletedAt)
    }
    query +=` )` + values + ` ) RETURNING id`
    var addedId int
    err = n.db.QueryRow(query, args...).Scan(&addedId)
    if err != nil {
        return -1, fmt.Errorf("%s: %w", op, err)
    }
	nextQuery := `DELETE FROM notes.todos WHERE id = $1`
	
	res, err := n.db.Exec(nextQuery, ID)
	if err != nil {
    	return -1, fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
    	return -1, fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

    return addedId, nil
}


func (n *NotesRepository) UpdateTodo(id int, title string, priority string,
	description string, category string, completedat string) (int, error) {
	
    const op = "repository.user_repository.UpdateTodo"

    query := `UPDATE notes.todos SET `
	values := ``
    args := []interface{}{}
    count := 0
	if title != "" {
		count ++
		values += fmt.Sprintf("TITLE=$%d", count)
		args = append(args, title)
	}
	if priority != "" {
		count ++
		values += fmt.Sprintf(" PRIORITY=$%d", count)
		args = append(args, priority)
	}
	if description != "" {
		count ++
		values += fmt.Sprintf(" DESCRIPTION=$%d", count)
		args = append(args, description)
	}
	if category != "" {
		count ++
		values += fmt.Sprintf(" CATEGORY=$%d", count)
		args = append(args, category)
	}
	if completedat != "" {
		count ++
		values += fmt.Sprintf(" COMPLETEDAT=$%d", count)
		args = append(args, completedat)
	}
	values = strings.ReplaceAll(values, " ", ",")
    query += values + ` RETURNING id`
    var dbId int
    err := n.db.QueryRow(query, args...).Scan(&dbId)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }
    return dbId, nil
}
