package services

import (
	_"log"
	"net/url"
	"notes-service/internal/models"
	"notes-service/internal/repository"
)

type NotesService interface {
	GetTodos(userId uint, urlParams url.Values) (*models.Todo, error)
}

type NotesServiceImpl struct {
	notesRepository *repository.NotesRepository
}

func NewNotesService(notesRepository *repository.NotesRepository) NotesService {
	return &NotesServiceImpl{
		notesRepository: notesRepository,
	}
}

func (s *NotesServiceImpl) GetTodos(userId uint, urlParams url.Values) (*models.Todo, error) {
	createdAt := urlParams.Get("createdAt")
	filter := urlParams.Get("after")
	limit := urlParams.Get("limit")
	offset := urlParams.Get("offset")
	task, err := s.notesRepository.GetTodos(userId, createdAt, filter, limit, offset)
	if err != nil {
		return nil, err
	}
	return task, nil
}