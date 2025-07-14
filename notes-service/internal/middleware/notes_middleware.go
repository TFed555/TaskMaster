package notes_middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"notes-service/internal/pkg/responses"
	. "notes-service/internal/services"
	"shared/middleware"
	"strconv"
	"strings"

	_ "github.com/go-chi/chi"
)

type NotesMiddleware struct {
	notesService NotesService
}

func NewNotesMiddleware(notesService NotesService) NotesMiddleware {
	return NotesMiddleware{
		notesService: notesService,
	}
}

func (n NotesMiddleware) SetNotesMIddleware(controller http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "POST" {
			controller.ServeHTTP(w, r)
            return
		}
		
		parts := strings.Split(r.URL.Path, "/")
		id := parts[len(parts)-1]
		log.Printf("ID: %s", id)


		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		r.Body = io.NopCloser(&buf)

		var req any
		switch r.Method {
        case "PATCH":
            var updateReq responses.UpdateRequest
			if err := json.NewDecoder(tee).Decode(&updateReq); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			updateReq.ID, _ = strconv.Atoi(id)
			req = updateReq

        case "PUT", "DELETE":
            req = map[string]interface{}{"ID": id}

        default:
            controller.ServeHTTP(w, r)
            return
        }

		ctx := r.Context()
		userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
		log.Print("UserID:", userID)
		if !ok {
	    http.Error(w, "Unauthorized", http.StatusUnauthorized)
	    return
		}
		log.Printf("Controller received userID: %v", userID)
		
		if 	err := n.notesService.Audit(r.Method, userID, req); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
        		return
		}

		controller.ServeHTTP(w, r)
	})
}
