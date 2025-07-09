package repository

import (
	"auth-service/internal/models"
	"auth-service/internal/pkg/domain_models"
	"database/sql"
	_ "errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sqlx.DB
}	

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(user *models.User) (error) {
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


func (r UserRepository) GetByEmail(email string) (models.User, error) {
	const op = "repository.user_repository.GetByEmail"

	query := (`
		SELECT id, login, email, password, created_at, img_path FROM auth.users
		WHERE email = $1
		LIMIT 1`)

	var user models.User

	err := r.db.QueryRowx(query, email).StructScan(&user)
	log.Print(user)

	log.Print(err)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return user, fmt.Errorf("%s: user not found", op)
        }
        return user, fmt.Errorf("%s: %w", op, err)
    }

	return user, nil
}

func (r UserRepository) GetUserByID(userId uint) (*models.User, error) {
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

func (r UserRepository) UpdateUser(user domain_models.User) (int, error) {
	const op = "repository.user_repository.UpdateUser"

	query := `UPDATE auth.users SET `
	values := []string{}
	args := []any{}
	count := 0

	if user.Email != "" {
		count ++
		values = append(values, fmt.Sprintf("EMAIL=$%d", count))
		args = append(args, user.Email)
	}
	if user.Name != "" {
		count ++ 
		values = append(values, fmt.Sprintf("LOGIN=$%d", count))
		args = append(args, user.Name)
	}
	if user.Password != "" {
		hashedPswd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return -1, fmt.Errorf("Can't hash password")
		}
		count ++
		values = append(values, fmt.Sprintf("PASSWORD=$%d", count))
		args = append(args, hashedPswd)
	}

	query += strings.Join(values, ", ")
	count ++
	query += fmt.Sprintf(` WHERE id = $%d RETURNING id`, count)
	args = append(args, user.ID)

	var dbId int
    err := r.db.QueryRow(query, args...).Scan(&dbId)
    if err != nil {
        return -1, fmt.Errorf("%s: %w", op, err)
    }
    return dbId, nil
}

func (r UserRepository) DeleteUser(userID uint) (bool, error) {
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