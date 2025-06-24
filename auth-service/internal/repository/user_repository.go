package repository

import (
	_"gorm.io/gorm"
	"auth-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}	

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	//вынести в отдельный файл
	query := `
		INSERT INTO users (login, email, password, created_at)
		VALUES (:login, :email, :password, NOW())
		RETURNING id, created_at
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(user, user)	
}


func (r *UserRepository) GetByEmail(email string) *models.User, error {
	query := `
		SELECT FROM users (login, email, password, created_at)
		WHERE email = ?
	`

	stmt, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(user, user)
}
