package services

import (
	_"log"
	"net/url"
	"notes-service/internal/models"
	"notes-service/internal/repository"
)

type NotesService interface {
	GetTodos(userId uint, urlParams url.Values) ([]models.Todo, error)
	CreateTask(userID uint, title string, priority string, description string,
	 category string, createdAt string, completedAt string) (int, error)
}

type NotesServiceImpl struct {
	notesRepository *repository.NotesRepository
}

func NewNotesService(notesRepository *repository.NotesRepository) NotesService {
	return &NotesServiceImpl{
		notesRepository: notesRepository,
	}
}

func (s *NotesServiceImpl) GetTodos(userId uint, urlParams url.Values) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(userId, urlParams)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *NotesServiceImpl) CreateTask(userID uint, title string, priority string, description string,
	 category string, createdAt string, completedAt string) (int, error) {
		id, err := s.notesRepository.CreateTodo(userID uint, title string, priority string, description string,
category string, createdAt string, completedAt string)
	 if err != nil {
		return nil, err
	 }
	 return id, nil
}