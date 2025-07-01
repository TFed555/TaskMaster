package repository

import (
	"auth-service/internal/models"
	"database/sql"
	_"errors"
	"fmt"

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