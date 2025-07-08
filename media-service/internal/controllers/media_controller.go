package controllers

import (
	"encoding/json"
	"log"
	"media-service/internal/pkg/domain_models"
	"media-service/internal/pkg/responses"
	"media-service/internal/services"
	"net/http"
	"shared/middleware"
)

type MediaController struct {
	mediaService services.MediaService
}

func NewMediaController(mediaService services.MediaService) MediaController {
	return MediaController{
		mediaService: mediaService,
	}
}

func (m MediaController) Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test endpoint works!",
	})
}

func (m MediaController) Generate(w http.ResponseWriter, r *http.Request) {
	newFile, err := m.mediaService.GenerateLinkPic()
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Print(newFile.Name, newFile.Link)

	w.WriteHeader(http.StatusAccepted)
	linkMsg := responses.GenerateResponse{
		Name: newFile.Name,
		Link: newFile.Link,
	}

	json.NewEncoder(w).Encode(linkMsg)
}

func (m MediaController) SaveAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserIdKey).(uint)
	if !ok {
		http.Error(w, "Unathorized", http.StatusUnauthorized)
	}
	log.Printf("Controller received userID: %v", userID)
	var req responses.SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	params := domain_models.ConfirmSaveParams{
		Name: req.Name,
		UserID: userID,
	}
	success, err := m.mediaService.SaveAvatarPic(params)
	if !success || err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := responses.ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: "Can't save pic",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(success)
}
