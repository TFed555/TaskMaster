package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"notes-service/internal/models"
	"notes-service/internal/services"
	"shared/middleware"
)

type NotesController struct {
	notesService services.NotesService
}

func NewNotesController(notesService services.NotesService) *NotesController {
	return &NotesController{
		notesService: notesService,
	}
}

type TodoResponse struct {
	Todos  []models.Todo  `json:"todos"`
}

type CreateRequest struct {
	Title string `json:"title"`
	Priority string  `json:"priority"`
	Category string `json:"category"`
	Description string `json:"description"`
	CreatedAt	string `json:"createdAt"`
	CompletedAt string	`json:"completedAt,omitempty"`
}

type CreateResponse struct {
	ID int `json:"id"`
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
	ctx := r.Context()
	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	log.Print("UserID:", userID)
	if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
	log.Printf("Controller received userID: %v", userID)

	log.Print(url.ParseQuery(r.URL.RawQuery))
	urlParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Print("Не удалось распарсить url")
	}

	todos, err := n.notesService.GetTodos(userID, urlParams)
	if err != nil {
		log.Println(err)
	}

	// log.Println(todos[0].Title)
	for idx, el := range todos {
		log.Println(idx, el)
	}
	response := TodoResponse{
		Todos: todos,
		// UserID: uint(todo.UserId),
		// Title: todo.Title,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (n *NotesController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
		userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	log.Print("UserID:", userID)
	if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
	log.Printf("Controller received userID: %v", userID)

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := n.notesService.CreateTask(userID, req.Title, req.Priority, req.Description, req.Category, req.CreatedAt, req.CompletedAt)

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

	response := CreateResponse {
		ID: id,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}