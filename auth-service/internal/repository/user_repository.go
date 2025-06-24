package repository

import (
	"gorm.io/gorm"
	"auth-service/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}	

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (error) {
	return r.db.Create(user).Error
	
}
