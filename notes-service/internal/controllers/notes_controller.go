package controllers

import (
	"encoding/json"
	"log"
	"net/http"
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
	// UserID uint   `json:"userId"`
	Title  string `json:"title"`
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
	log.Print(userID)
	todo, err := n.notesService.GetTodos(userID)
	if err != nil {
		log.Println("D")
	}
	response := TodoResponse{
		// UserID: uint(todo.UserId),
		Title: todo.Title,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}