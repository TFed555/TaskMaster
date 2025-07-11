package services

import (
	"log"
	_ "log"
	"net/url"
	"notes-service/internal/models"
	"notes-service/internal/pkg/domain_models"
	"notes-service/internal/repository"
)

type NotesService interface {
	GetTodos(urlParams url.Values) ([]models.Todo, error)
	GetArchivedTodos(urlParams url.Values) ([]models.Todo, error)
	CreateTask(todoBody domain_models.Todo) (int, error)
	UpdateTask(todoBody domain_models.Todo) (int, error)
	ArchiveTask(ID int) (int, error)
	GetOneTodo(taskId int) (*models.Todo, error)
	DeleteTodo(taskId int) (bool, error)
	RestoreTodo(taskId int) (int, error)
	CreateTag(tag domain_models.Tag) (int, error)
}

type NotesServiceImpl struct {
	notesRepository repository.NotesRepository
}

func NewNotesService(notesRepository repository.NotesRepository) NotesService {
	return &NotesServiceImpl{
		notesRepository: notesRepository,
	}
}

func (s NotesServiceImpl) GetTodos(urlParams url.Values) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(urlParams, "todos")
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s NotesServiceImpl) CreateTask(todoBody domain_models.Todo) (int, error) {
	todo := models.Todo{
		Title: todoBody.Title,
		Priority: todoBody.Priority,
		Category: todoBody.Category,
		Description: todoBody.Description,
		CreatedAt: todoBody.CreatedAt,
		CompletedAt: todoBody.CompletedAt,
		UserId: todoBody.UserId,
	}
	log.Print("NotesService:", todo)
	id, err := s.notesRepository.CreateTodo(todo)
	 if err != nil {
		return -1, err
	 }
	 return id, nil
}

func (s NotesServiceImpl) ArchiveTask(ID int) (int, error) {

	id, err := s.notesRepository.ArchiveTodo(ID)
	 if err != nil {
		return -1, err
	 }
	 return id, nil
}

func (s NotesServiceImpl) UpdateTask(todoBody domain_models.Todo) (int, error) {
	todo := models.Todo{
		ID: todoBody.ID,
		Title: todoBody.Title,
		Priority: todoBody.Priority,
		Category: todoBody.Category,
		Description: todoBody.Description,
		CreatedAt: todoBody.CreatedAt,
		CompletedAt: todoBody.CompletedAt,
	}
	id, err := s.notesRepository.UpdateTodo(todo)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s NotesServiceImpl) GetArchivedTodos(urlParams url.Values) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(urlParams, "archived_todos")
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s NotesServiceImpl) GetOneTodo(taskId int) (*models.Todo, error) {
	todo, err := s.notesRepository.GetTodoByID(taskId, "todos")
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s NotesServiceImpl) DeleteTodo(taskId int) (bool, error) {
	result, err := s.notesRepository.DeleteTodo(taskId)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s NotesServiceImpl) RestoreTodo(taskId int) (int, error) {
	id, err := s.notesRepository.RestoreTodo(taskId)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s NotesServiceImpl) CreateTag(tag domain_models.Tag) (int, error) {
	tagName, userId := tag.Name, tag.UserID
	id, err := s.notesRepository.CreateTag(tagName, userId)
	if err != nil {
		return -1, err
	}
	return id, nil
}