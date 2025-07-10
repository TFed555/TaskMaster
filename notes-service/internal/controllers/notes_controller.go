package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"notes-service/internal/pkg/domain_models"
	"notes-service/internal/pkg/responses"
	"notes-service/internal/services"
	_ "shared/middleware"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type NotesController struct {
	notesService services.NotesService
}

//go:generate mockery --name=NotesService --dir=../services --output=./mocks --case=underscore
func NewNotesController(notesService services.NotesService) NotesController {
	return NotesController{
		notesService: notesService,
	}
}


func (n NotesController) Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test endpoint works!",
	})
}

func (n NotesController) GetTodos(w http.ResponseWriter, r *http.Request) {
	// /api/todos?createdAt=(date YYYY-MM-DD)&filter=(after | before)&offset=(int)&limit=(int)
	// ctx := r.Context()
	// userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	log.Print(url.ParseQuery(r.URL.RawQuery))
	urlParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Print("Не удалось распарсить url")
	}
	todos, err := n.notesService.GetTodos(urlParams)

	// todos, err := n.notesService.GetTodos(userID, urlParams)

	if err != nil {
		log.Println(err)
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Не удалось получить данные",
		})
	}

	for idx, el := range todos {
		log.Println(idx, el)
	}

	masTodos := make([]responses.OneTodoResponse, 0)

	for _, el := range todos {
		todo := responses.OneTodoResponse{
			Title:       el.Title,
			Priority:    el.Priority,
			Category:    el.Category,
			Description: el.Description,
			CreatedAt:   el.CreatedAt,
			CompletedAt: el.CompletedAt,
			ID:          *el.ID,
		}
		masTodos = append(masTodos, todo)
	}

	response := responses.TodoResponse{
		Todos: &masTodos,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (n NotesController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// 	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	var req responses.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("NotesController: %v", req)
	todoBody := domain_models.Todo{
		UserId: req.UserID,
		Title: req.Title,
		Priority: req.Priority,
		Description: req.Description,
		Category: req.Category,
		CreatedAt: req.CreatedAt,
		CompletedAt: req.CompletedAt,
	}

	log.Print("NotesController:", todoBody)
	id, err := n.notesService.CreateTask(todoBody)

	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		// errMsg := ErrorResponse{
		// 	Status:  http.StatusInternalServerError,
		// 	Message: "User does not exists",
		// }
		json.NewEncoder(w).Encode("bad")
		return
	}

	response := responses.CreateResponse{
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n NotesController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// 	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	var req responses.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todoBody := domain_models.Todo{
		ID: &req.ID,
		Title: req.Title,
		Priority: req.Priority,
		Description: req.Description,
		Category: req.Category,
		CreatedAt: req.CreatedAt,
		CompletedAt: req.CompletedAt,
	}

	id, err := n.notesService.UpdateTask(todoBody)

	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		// errMsg := ErrorResponse{
		// 	Status:  http.StatusInternalServerError,
		// 	Message: "User does not exists",
		// }
		json.NewEncoder(w).Encode("bad")
		return
	}

	response := responses.CreateResponse{
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n NotesController) GetArchivedTodos(w http.ResponseWriter, r *http.Request) {
	// /api/archivedtodos?createdAt=(date YYYY-MM-DD)&filter=(after | before)&offset=(int)&limit=(int)
	// ctx := r.Context()
	// userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	log.Print(url.ParseQuery(r.URL.RawQuery))
	urlParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Print("Не удалось распарсить url")
	}
	todos, err := n.notesService.GetArchivedTodos(urlParams)

	// todos, err := n.notesService.GetTodos(userID, urlParams)

	if err != nil {
		log.Println(err)
	}

	for idx, el := range todos {
		log.Println(idx, el)
	}

	// var masTodos []OneTodoResponse
	masTodos := make([]responses.OneTodoResponse, 0)

	for _, el := range todos {
		todo := responses.OneTodoResponse{
			Title:       el.Title,
			Priority:    el.Priority,
			Category:    el.Category,
			Description: el.Description,
			CreatedAt:   el.CreatedAt,
			CompletedAt: el.CompletedAt,
			ID:          *el.ID,
		}
		masTodos = append(masTodos, todo)
	}
	response := responses.TodoResponse{
		Todos: &masTodos,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (n NotesController) ArchiveTodo(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// 	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	taskId := chi.URLParam(r, "id")
	if taskId == "" {
		log.Print("Не удалось получить id задачи")
	}

	id, err := strconv.Atoi(taskId)
	if err != nil {
		log.Print("Не удалось преобразовать id задачи")
	}

	//передавать в модели
	id, err = n.notesService.ArchiveTask(id)

	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		// errMsg := ErrorResponse{
		// 	Status:  http.StatusInternalServerError,
		// 	Message: "User does not exists",
		// }
		json.NewEncoder(w).Encode("bad")
		return
	}

	response := responses.CreateResponse{
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n NotesController) GetOneTodo(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// 	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	taskId := chi.URLParam(r, "id")
	if taskId == "" {
		log.Print("Не удалось получить id задачи")
	}
	id, err := strconv.Atoi(taskId)
	if err != nil {
		log.Print("Не удалось преобразовать id задачи")
	}
	todo, err := n.notesService.GetOneTodo(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Can't get task info", http.StatusBadRequest)
		return
	}

	response := responses.OneTodoResponse{
		ID:          *todo.ID,
		Title:       todo.Title,
		Priority:    todo.Priority,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		CompletedAt: todo.CompletedAt,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n NotesController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "id")
	if taskId == "" {
		log.Print("Не удалось получить id задачи")
	}
	id, err := strconv.Atoi(taskId)
	if err != nil {
		log.Print("Не удалось преобразовать id задачи")
	}
	success, err := n.notesService.DeleteTodo(id)
	if !success || err != nil {
		log.Println(err)
		http.Error(w, "Can't get task info", http.StatusBadRequest)
		return
	}

	response := responses.DeleteResponse{
		Message: "Deleted successfully",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n NotesController) RestoreTodo(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "id")
	if taskId == "" {
		log.Print("Не удалось получить id задачи")
	}
	id, err := strconv.Atoi(taskId)
	if err != nil {
		log.Print("Не удалось преобразовать id задачи")
	}
	newId, err := n.notesService.RestoreTodo(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Can't get task info", http.StatusBadRequest)
		return
	}

	response := responses.CreateResponse{
		ID: newId,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}