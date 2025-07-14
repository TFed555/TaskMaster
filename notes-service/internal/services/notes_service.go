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
	GetTodos(urlParams url.Values, userID uint) ([]models.Todo, error)
	GetArchivedTodos(urlParams url.Values, userID uint) ([]models.Todo, error)
	CreateTask(todoBody domain_models.Todo) (int, error)
	UpdateTask(todoBody domain_models.Todo) (int, error)
	ArchiveTask(ID int) (int, error)
	GetOneTodo(taskId int) (*models.Todo, error)
	DeleteTodo(taskId int) (bool, error)
	RestoreTodo(taskId int) (int, error)
	CreateTag(tag domain_models.Tag) (int, error)
	GetTags(userID uint) ([]models.Tag, error)
	UpdateTag(tagBody domain_models.Tag) (int, error)
	DeleteTag(tagId int) (bool, error)
	AddTagToTodo(tagTodo domain_models.TagTodo) (bool, error)
	ReduceTag(tagTodo domain_models.TagTodo) (bool, error)
	Audit(method string, isArchived bool, userID uint, todoId int) (error)
	GetAuditTrail(userID uint) ([]models.HistoryTodo, error)
}

type NotesServiceImpl struct {
	notesRepository repository.NotesRepository
}

func NewNotesService(notesRepository repository.NotesRepository) NotesService {
	return &NotesServiceImpl{
		notesRepository: notesRepository,
	}
}

func (s NotesServiceImpl) GetTodos(urlParams url.Values, userID uint) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(urlParams, userID, "todos")
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

func (s NotesServiceImpl) GetArchivedTodos(urlParams url.Values, userID uint) ([]models.Todo, error) {
	todos, err := s.notesRepository.GetTodos(urlParams, userID, "archived_todos")
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

func (s NotesServiceImpl) GetTags(userID uint) ([]models.Tag, error) {
	tags, err := s.notesRepository.GetTags(userID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (s NotesServiceImpl) UpdateTag(tagBody domain_models.Tag) (int, error) {
	tag := models.Tag{
		ID: *tagBody.ID,
		Name: tagBody.Name,
	}
	id, err := s.notesRepository.UpdateTag(tag)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s NotesServiceImpl) DeleteTag(tagId int) (bool, error) {
	result, err := s.notesRepository.DeleteTag(tagId)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s NotesServiceImpl) AddTagToTodo(tagTodo domain_models.TagTodo) (bool, error) {
	result, err := s.notesRepository.AddTagToTodo(tagTodo.TodoID, tagTodo.TagID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s NotesServiceImpl) ReduceTag(tagTodo domain_models.TagTodo) (bool, error) {
	result, err := s.notesRepository.ReduceTag(tagTodo.TodoID, tagTodo.TagID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s NotesServiceImpl) Audit(method string, isArchived bool, userID uint, todoId int) (error) {

    return s.notesRepository.AuditTodo(todoId, userID, isArchived, method)
}

func (s NotesServiceImpl) GetAuditTrail(userID uint) ([]models.HistoryTodo, error) {
	todos, err := s.notesRepository.GetHistoryTodos(userID)
	if err != nil {
		return []models.HistoryTodo{}, err
	}
	return todos, nil
}