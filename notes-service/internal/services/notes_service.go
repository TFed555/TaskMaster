package services

import (
	_"log"
	"net/url"
	"notes-service/internal/models"
	"notes-service/internal/repository"
)

type NotesService interface {
	GetTodos(userId uint, urlParams url.Values) ([]models.Todo, error)
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