package repository

import (
	"auth-service/internal/models"
	"database/sql"
	_ "errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}	

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (error) {
	const op = "repository.user_repository.Create"
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO auth.users (login, email, password, created_at)
		VALUES (:login, :email, :password, NOW())
		RETURNING id, created_at
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	return stmt.Get(user, user)	
}


func (r *UserRepository) GetByEmail(email string) (*models.User, error, bool) {
	const op = "repository.user_repository.GetByEmail"

	query := (`
		SELECT id, login, email, password, created_at FROM auth.users
		WHERE email = $1
		LIMIT 1`)

	var user models.User

	err := r.db.QueryRowx(query, email).StructScan(&user)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("%s: user not found", op), false
        }
        return nil, fmt.Errorf("%s: %w", op, err), false
    }

	return &user, nil, true
}

func (r *UserRepository) GetUserByID(userId uint) (*models.User, error) {
	const op = "repository.user_repository.GetByID"

	query := (`
		SELECT id, login, email, password, created_at FROM auth.users
		WHERE id = $1
		LIMIT 1`)

	var user models.User

	err := r.db.QueryRowx(query, userId).StructScan(&user)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("%s: user not found", op)
        }
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return &user, nil
}

func (r *UserRepository) UpdateUser(userID uint, email string, name string, password string, img_path string) (int, error) {
	const op = "repository.user_repository.UpdateUser"

	query := `UPDATE auth.users SET `
	values := ``
	args := []interface{}{}
	count := 0

	if email != "" {
		count ++
		values += fmt.Sprintf("EMAIL=$%d", count)
		args = append(args, email)
	}
	if name != "" {
		count ++ 
		values += fmt.Sprintf(" LOGIN=$%d", count)
		args = append(args, name)
	}
	if password != "" {
		count ++
		values += fmt.Sprintf(" PASSWORD=$%d", count)
		args = append(args, password)
	}
	if img_path != "" {
		count ++
		values += fmt.Sprintf(" IMG_PATH=$%d", count)
		args = append(args, img_path)
	}
	values = strings.ReplaceAll(values, " ", ",")
	count ++
	query += fmt.Sprintf(values + ` WHERE id = $%d RETURNING id`, count)
	args = append(args, userID)

	var dbId int
    err := r.db.QueryRow(query, args...).Scan(&dbId)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }
    return dbId, nil
}

func (r *UserRepository) DeleteUser(userID uint) (bool, error) {
	const op = "repository.user_repository.DeleteUser"

	query := `DELETE FROM AUTH.USERS WHERE id = $1`

	res, err := r.db.Exec(query, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected, err := res.RowsAffected(); rowsAffected < 0 {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}