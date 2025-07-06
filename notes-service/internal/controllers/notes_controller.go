package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	_ "notes-service/internal/models"
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

// вынести в internal/pkg
type TodoResponse struct {
	Todos *[]OneTodoResponse `json:"todos"`
}

type CreateRequest struct {
	Title       string `json:"title"`
	Priority    string `json:"priority"`
	Category    string `json:"category"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	CompletedAt string `json:"completedAt,omitempty"`
	UserID      int    `json:"userID"`
}

type ArchiveRequest struct {
	ID int `json:"id"`
}

type CreateResponse struct {
	ID int `json:"id"`
}

type UpdateRequest struct {
	Title       string `json:"title,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	CompletedAt string `json:"completedAt,omitempty"`
	ID          int    `json:"id"`
}

type OneTodoResponse struct {
	Title       string  `json:"title,omitempty"`
	Priority    string  `json:"priority,omitempty"`
	Category    string  `json:"category,omitempty"`
	Description string  `json:"description,omitempty"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	CompletedAt *string `json:"completedAt,omitempty"`
	ID          int     `json:"id"`
}

type ErrorResponse struct {
	Status  uint   `json:"code"`
	Message string `json:"message"`
}

func (n *NotesController) Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test endpoint works!",
	})
}

func (n *NotesController) GetTodos(w http.ResponseWriter, r *http.Request) {
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
	}

	for idx, el := range todos {
		log.Println(idx, el)
	}

	masTodos := make([]OneTodoResponse, 0)

	for _, el := range todos {
		todo := OneTodoResponse{
			Title:       el.Title,
			Priority:    el.Priority,
			Category:    el.Category,
			Description: el.Description,
			CreatedAt:   el.CreatedAt,
			CompletedAt: el.CompletedAt,
			ID:          el.ID,
		}
		masTodos = append(masTodos, todo)
	}

	response := TodoResponse{
		Todos: &masTodos,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (n *NotesController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// 	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//передавать в модели
	id, err := n.notesService.CreateTask(uint(req.UserID), req.Title, req.Priority, req.Description, req.Category, req.CreatedAt, req.CompletedAt)

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

	response := CreateResponse{
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n *NotesController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// 	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	// log.Print("UserID:", userID)
	// if !ok {
	//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//     return
	// }
	// log.Printf("Controller received userID: %v", userID)

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := n.notesService.UpdateTask(req.ID, req.Title, req.Priority, req.Description, req.Category, req.CompletedAt)

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

	response := CreateResponse{
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n *NotesController) GetArchivedTodos(w http.ResponseWriter, r *http.Request) {
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
	masTodos := make([]OneTodoResponse, 0)

	for _, el := range todos {
		todo := OneTodoResponse{
			Title:       el.Title,
			Priority:    el.Priority,
			Category:    el.Category,
			Description: el.Description,
			CreatedAt:   el.CreatedAt,
			CompletedAt: el.CompletedAt,
			ID:          el.ID,
		}
		masTodos = append(masTodos, todo)
	}
	response := TodoResponse{
		Todos: &masTodos,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (n *NotesController) ArchiveTodo(w http.ResponseWriter, r *http.Request) {
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

	response := CreateResponse{
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (n *NotesController) GetOneTodo(w http.ResponseWriter, r *http.Request) {
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

	response := OneTodoResponse{
		ID:          todo.ID,
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
