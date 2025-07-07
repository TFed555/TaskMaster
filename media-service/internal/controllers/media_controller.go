package controllers

import (
	"encoding/json"
	"log"
	"media-service/internal/pkg/domain_models"
	"media-service/internal/pkg/responses"
	"media-service/internal/services"
	"net/http"
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

func (m MediaController) SaveAvatar(w http.ResponseWriter, r *http.Request) {
	var req responses.SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	params := domain_models.SavePicParams{
		File: req.File,
	}
	link, err := m.mediaService.SavePic(params)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Print(link)

	w.WriteHeader(http.StatusAccepted)
	linkMsg := responses.SaveResponse{
		Link: link,
	}

	json.NewEncoder(w).Encode(linkMsg)
}
