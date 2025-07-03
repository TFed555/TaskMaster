package services

import (
	_"log"
	"net/url"
	"notes-service/internal/models"
	"notes-service/internal/repository"
)

type NotesService interface {
	GetTodos(urlParams url.Values) ([]models.Todo, error)
	GetArchivedTodos(urlParams url.Values) ([]models.Todo, error)
	CreateTask(userID uint, title string, priority string, description string,
	 category string, createdAt string, completedAt string) (int, error)
	UpdateTask(ID int, title string, priority string,
	description string, category string, completedat string) (int, error)
	ArchiveTask(ID int) (int, error)
	GetOneTodo(taskId int) (*models.Todo, error)
}

type NotesServiceImpl struct {
	notesRepository *repository.NotesRepository
}

func NewNotesService(notesRepository *repository.NotesRepository) NotesService {
	return &NotesServiceImpl{
		notesRepository: notesRepository,
	}
}

func (s *NotesServiceImpl) GetTodos(urlParams url.Values) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(urlParams, "todos")
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *NotesServiceImpl) CreateTask(userID uint, title string, priority string, description string,
	 category string, createdAt string, completedAt string) (int, error) {

	id, err := s.notesRepository.CreateTodo(userID, title, priority, description,
			category, createdAt, completedAt)
	 if err != nil {
		return -1, err
	 }
	 return id, nil
}

func (s *NotesServiceImpl) ArchiveTask(ID int) (int, error) {

	id, err := s.notesRepository.ArchiveTodo(ID)
	 if err != nil {
		return -1, err
	 }
	 return id, nil
}

func (s *NotesServiceImpl) UpdateTask(ID int, title string, priority string,
	description string, category string, completedat string) (int, error) {
	id, err := s.notesRepository.UpdateTodo(ID, title, priority, description, category, completedat)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *NotesServiceImpl) GetArchivedTodos(urlParams url.Values) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(urlParams, "archived_todos")
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *NotesServiceImpl) GetOneTodo(taskId int) (*models.Todo, error) {
	todo, err := s.notesRepository.GetTodoByID(taskId)
	if err != nil {
		return nil, err
	}
	return todo, nil
}