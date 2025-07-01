package controllers

import (
	"encoding/json"
	"net/http"
	"notes-service/internal/services"
)

type NotesController struct {
	notesService services.NotesService
}

func NewNotesController(notesService services.NotesService) *NotesController {
	return &NotesController{
		notesService: notesService,
	}
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
	// todo, err := n.notesService.GetTodos(r.)
}