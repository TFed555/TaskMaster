package controllers

import (
	"encoding/json"
	"log"
	"media-service/internal/pkg/domain_models"
	"media-service/internal/pkg/responses"
	"media-service/internal/services"
	"net/http"
	"shared/middleware"
	"strings"
)

type MediaController struct {
	mediaService services.MediaService
	authService middleware.AuthService
}

func NewMediaController(mediaService services.MediaService, 
							authService middleware.AuthService) MediaController {
	return MediaController{
		mediaService: mediaService,
		authService: authService,
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

	// log.Print(newFile.Name, newFile.Link)

	w.WriteHeader(http.StatusAccepted)
	linkMsg := responses.GenerateResponse{
		Name: newFile.Name,
		Link: newFile.Link,
	}

	json.NewEncoder(w).Encode(linkMsg)
}

func (m MediaController) SaveAvatar(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// userID, ok := ctx.Value(middleware.UserIdKey).(uint)
	// if !ok {
	// 	http.Error(w, "Unathorized", http.StatusUnauthorized)
	// }
	log.Print("Header:", r.Header)
	
	cookiesmas := r.Header.Get("Cookie")
	newString := strings.Split(cookiesmas, "; ")
	var refreshToken string
	for _, el := range newString {
		newEl := strings.Split(el, "=")
		if newEl[0] == "refresh_token" {
			refreshToken = newEl[1]
		}
	}
	// log.Print(cookiesmas)
	// log.Print(newString)
	log.Print(refreshToken)
	result, errMsg := m.authService.ValidateToken(refreshToken)

	if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			return
	}

	userID, err := m.authService.ParseUserId(refreshToken)

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
	
	success, newErr := m.mediaService.SaveAvatarPic(params)
	if !success || newErr != nil {
		log.Print(newErr)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := responses.ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: "Can't save pic",
		}
		json.NewEncoder(w).Encode(errMsg)
		log.Print("Не удалось сохранить")
		return
	}

	link, newErr := m.mediaService.GetAvatarPic(params.Name)
	if newErr != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := responses.ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: "Can't get pic",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	linkMsg := responses.LinkResponse{
		Link: link,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(linkMsg)
}
