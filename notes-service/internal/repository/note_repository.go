package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"notes-service/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type NotesRepository struct {
	db *sqlx.DB
}

func NewNotesRepository(db *sqlx.DB) NotesRepository {
	return NotesRepository{
		db: db,
	}
}

func (n NotesRepository) GetTodos(urlParams url.Values, userID uint, tableName string) ([]models.Todo, error) {
	const op = "repository.notes_repository.GetTodos"
	args := []any{}
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


	query := fmt.Sprintf(`
		SELECT id, title, priority, category, description, createdat, completedat, userid FROM notes.%s `, tableName)

	todos := []models.Todo{}
	counter := 0

	st := "WHERE"
	
		counter++
		st = "AND"
		query += fmt.Sprintf("WHERE userid = $%d ", counter)
		args = append(args, int(userID))
		log.Printf("UserID: %s", userID)
	

	if dataToFilter := urlParams.Get("createdAt"); dataToFilter != "" {
		counter++

		query += fmt.Sprintf(" %s createdat %s $%d::TIMESTAMPTZ ", st, sqlFilter, counter)
		args = append(args, dataToFilter)
		log.Printf("Date: %s", dataToFilter)
	}

	if offset := urlParams.Get("offset"); offset != "" {
		counter++
		query += fmt.Sprintf(" OFFSET $%d", counter)
		args = append(args, offset)
		log.Printf("Offset: %s", offset)
	}

	if limit := urlParams.Get("limit"); limit != "" {
		counter++
		query += fmt.Sprintf(" LIMIT $%d ", counter)
		args = append(args, limit)
		log.Printf("Limit: %s", limit)
	}

	err := n.db.Select(&todos, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: todos not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// log.Print(*todos[0].ID)
	for i := range todos {
		el := &todos[i]
		nextQuery := `SELECT id, name FROM notes.tags WHERE id IN (SELECT tagid FROM notes.todo_tags WHERE todoid = $1)`
		tags := []models.Tag{}
		err = n.db.Select(&tags, nextQuery, *el.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		if err == nil {
			el.Tags = append(el.Tags, tags...)
			log.Printf("Mass %v", el)
		}
	}

	return todos, nil
}

func (n NotesRepository) GetTodoByID(ID int, tableName string) (*models.Todo, error) {
	const op = "repository.notes_repository.GetTodoByID"
	query := fmt.Sprintf(`SELECT id, title, priority, category, description, createdat, completedat, userid
			FROM notes.%s WHERE id = $1 LIMIT 1`, tableName)
	todo := models.Todo{}
	err := n.db.QueryRowx(query, ID).StructScan(&todo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: todos not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nextQuery := `SELECT id, name FROM notes.tags WHERE id IN (SELECT tagid FROM notes.todo_tags WHERE todoid = $1)`
	tags := []models.Tag{}
	err = n.db.Select(&tags, nextQuery, *todo.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	if err == nil {
	todo.Tags = append(todo.Tags, tags...)
		log.Printf("Mass %v", todo)
	}

	return &todo, nil
}

func (n NotesRepository) CreateTodo(todo models.Todo) (int, error) {
	const op = "repository.notes_repository.CreateTodo"

	query := `INSERT INTO notes.todos (userid, title, priority, description, category, createdat`
	values := `VALUES ($1, $2, $3, $4, $5, $6`
	args := []any{todo.UserId, todo.Title, todo.Priority, todo.Description, todo.Category, todo.CreatedAt}
	count := 6
	log.Printf("%s, %s", op, todo.CreatedAt)

	if todo.CompletedAt != nil {
		query += `, completedat`
		count++
		log.Printf("%s %s", op, *todo.CompletedAt)
		values += fmt.Sprintf(`, $%d`, count)
		args = append(args, todo.CompletedAt)
	}
	query += ` )` + values + ` ) RETURNING id`
	var id int
	err := n.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (n NotesRepository) ArchiveTodo(ID int) (int, error) {
	const op = "repository.notes_repository.ArchiveTodo"

	todo, err := n.GetTodoByID(ID, "todos")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO notes.archived_todos (userid, title, priority, description, category, createdat`
	values := `VALUES ($1, $2, $3, $4, $5, $6`
	args := []interface{}{todo.UserId, todo.Title, todo.Priority, todo.Description, todo.Category, todo.CreatedAt}
	count := 6

	if todo.CompletedAt != nil {
		query += `, completedat`
		count++

		values += fmt.Sprintf(`, $%d`, count)
		args = append(args, todo.CompletedAt)
	}
	query += ` )` + values + ` ) RETURNING id`
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
	var tagId int
	err = n.db.QueryRow(`SELECT tagId from notes.todo_tags WHERE todoId = $1`, ID).Scan(&tagId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
   		return -1, fmt.Errorf("%s: %w", op, err)
	}
	if err == nil {
		nextQuery = `INSERT INTO notes.archived_todo_tags (todoid, tagid) VALUES ($1, $2)`
		_, err = n.db.Exec(nextQuery, addedId, tagId)
		if err != nil {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
	}
	return addedId, nil
}

func (n NotesRepository) UpdateTodo(todo models.Todo) (int, error) {

	const op = "repository.notes_repository.UpdateTodo"

	query := `UPDATE notes.todos SET `
	values := []string{}
	args := []any{}
	count := 0
	if todo.Title != "" {
		count++
		values = append(values, fmt.Sprintf("TITLE=$%d", count))
		args = append(args, todo.Title)
	}
	if todo.Priority != "" {
		count++
		values = append(values, fmt.Sprintf("PRIORITY=$%d", count))
		args = append(args, todo.Priority)
	}
	if todo.Description != "" {
		count++
		values = append(values, fmt.Sprintf("DESCRIPTION=$%d", count))
		args = append(args, todo.Description)
	}
	if todo.Category != "" {
		count++
		values = append(values, fmt.Sprintf("CATEGORY=$%d", count))
		args = append(args, todo.Category)
	}
	if todo.CompletedAt != nil {
		count++
		values = append(values, fmt.Sprintf(" COMPLETEDAT=NULLIF($%d, '')::date", count))
		args = append(args, todo.CompletedAt)
	}
	query += strings.Join(values, ", ")
	count++
	query += fmt.Sprintf(` WHERE id = $%d RETURNING id`, count)
	args = append(args, todo.ID)
	log.Printf("%s, %d", op, todo.ID)
	var dbId int
	err := n.db.QueryRow(query, args...).Scan(&dbId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return dbId, nil
}

func (n NotesRepository) DeleteTodo(taskId int) (bool, error) {
	const op = "repository.notes_repository.DeleteTodo"

	query := `DELETE FROM notes.archived_todos WHERE id = $1`
	res, err := n.db.Exec(query, taskId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	resultFromArchived, _ := res.RowsAffected()

	var tagId int
	err = n.db.QueryRow(`SELECT tagId from notes.archived_todo_tags WHERE todoId = $1`, taskId).Scan(&tagId)
	
	if err != nil && !errors.Is(err, sql.ErrNoRows){
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		query = `DELETE FROM notes.archived_todo_tags WHERE todoid = $1`
		_, err = n.db.Exec(query, taskId)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
	}

	return resultFromArchived > 0, nil
}

func (n NotesRepository) RestoreTodo(taskId int) (int, error) {
	const op = "repository.notes_repository.RestoreTodo"

	archivedTodo, err := n.GetTodoByID(taskId, "archived_todos")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	query := `INSERT INTO notes.todos (userid, title, priority, description, category, createdat`
	values := `VALUES ($1, $2, $3, $4, $5, $6`
	args := []any{archivedTodo.UserId, archivedTodo.Title,
		archivedTodo.Priority, archivedTodo.Description,
		archivedTodo.Category, archivedTodo.CreatedAt}
	count := 6

	if archivedTodo.CompletedAt != nil {
		query += `, completedat`
		count++

		values += fmt.Sprintf(`, $%d`, count)
		args = append(args, archivedTodo.CompletedAt)
	}
	query += ` )` + values + ` ) RETURNING id`
	var addedId int
	err = n.db.QueryRow(query, args...).Scan(&addedId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	nextQuery := `DELETE FROM notes.archived_todos WHERE id = $1`

	res, err := n.db.Exec(nextQuery, taskId)
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


func (n NotesRepository) CreateTag(tagName string, userId uint) (int, error) {
	const op = "repository.notes_repository.createTag"

	query := `INSERT INTO notes.tags (name, userid)
				VALUES ($1, $2) RETURNING ID`
	var id int
	err := n.db.QueryRow(query, tagName, userId).Scan(&id)
	if err != nil {
		log.Printf("%s, %s", op, err)
		return -1, err
	}

	return id, nil
}

func (n NotesRepository) GetTags(userID uint) ([]models.Tag, error) {
	const op = "repository.notes_repository.getTags"

	query := `SELECT id, name FROM notes.tags WHERE userid = $1`
	tags := []models.Tag{}
	err := n.db.Select(&tags, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: tags not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tags, nil
}

func (n NotesRepository) UpdateTag(tag models.Tag) (int, error) {
	const op = "repository.notes_repository.UpdateTag"

	query := `UPDATE notes.tags SET `
	values := []string{}
	args := []any{}
	count := 0
	if tag.Name != "" {
		count++
		values = append(values, fmt.Sprintf("NAME=NULLIF($%d, '')", count))
		args = append(args, tag.Name)
	}
	query += strings.Join(values, "")
	count++
	query += fmt.Sprintf(` WHERE id = $%d RETURNING id`, count)
	args = append(args, tag.ID)
	log.Printf("%s, %d", op, tag.ID)
	var dbId int
	err := n.db.QueryRow(query, args...).Scan(&dbId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return dbId, nil
}

func (n NotesRepository) DeleteTag(tagId int) (bool, error) {
	const op = "repository.notes_repository.DeleteTag"

	query := `DELETE FROM notes.tags WHERE id = $1`
	res, err := n.db.Exec(query, tagId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	result, _ := res.RowsAffected()
	return result > 0, nil
}

func (n NotesRepository) AddTagToTodo(todoID int, tagID int) (bool, error) {
	const op = "repository.notes_repository.AddTagToTodo"

	query := `SELECT todoid from notes.todo_tags WHERE todoid = $1 and tagid = $2`
	var selectedId int
	err := n.db.QueryRow(query, todoID, tagID).Scan(&selectedId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("%s: %v", op, err)
	}
	if selectedId > 0 {
		return false, fmt.Errorf("This connection already exists")
	}

	nextQuery := `INSERT INTO notes.todo_tags(todoid, tagid) `
	values := `VALUES ($1, $2) RETURNING todoid`
	args := []any{todoID, tagID}

	nextQuery += values
	var insertedId int
	err = n.db.QueryRow(nextQuery, args...).Scan(&insertedId)
	if err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}
	return true, nil
}

func (n NotesRepository) ReduceTag(todoID int, tagID int) (bool, error) {
	const op = "repository.notes_repository.ReduceTag"

	query := `DELETE FROM notes.todo_tags WHERE todoid = $1 AND tagid = $2`
	args := []any{todoID, tagID}

	res, err := n.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}
	return rowsAffected>0, nil
}

func (n NotesRepository) AuditTodo(id int, userId uint, method string) (error) {
	const op = "repository.notes_repository.AuditTodo"
	log.Printf("CALLED FROM REPOSITORY %d", id)
	action := "Deleted"
	idTodo := "archived_todoid"
	if method == "PATCH" {
		action = "Changed"
		idTodo = "todoid"
	}
	query := fmt.Sprintf(`INSERT INTO todos_history (%s, userId, action) `, idTodo)
	values := `VALUES ($1, $2, $3)`
	args := []any{id, userId, action}

	query += values + ` RETURNING id`
	var insertedId int
	err := n.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return err
	}
	log.Print(insertedId)
	return nil
}

func (n NotesRepository) GetHistoryTodos(userID uint) ([]models.HistoryTodo, error) {
	const op = "repository.notes_repository.GetHistoryTodos"
	query := (`SELECT t.title, th.action 
			FROM notes.todos t
			JOIN todos_history th ON t.id = th.todoId
			WHERE t.userid = $1`)
	todos := []models.HistoryTodo{}
	err := n.db.Select(&todos, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: todos not found", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return todos, nil
}