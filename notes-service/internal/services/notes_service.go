package services

import (
	"notes-service/internal/models"
	"notes-service/internal/repository"
)

type NotesService interface {
	GetTodos(userId uint) (*models.Todo, error)
}

type NotesServiceImpl struct {
	notesRepository *repository.NotesRepository
}

func NewNotesService(notesRepository *repository.NotesRepository) NotesService {
	return &NotesServiceImpl{
		notesRepository: notesRepository,
	}
}

func (s *NotesServiceImpl) GetTodos(userId uint) (*models.Todo, error) {
	task, err := s.notesRepository.GetTodos(userId)
	if err != nil {
		return nil, err
	}
	return task, nil
}