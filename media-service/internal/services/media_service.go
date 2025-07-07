package services

import (
	// "log"
	"media-service/internal/pkg/domain_models"
	"media-service/internal/repository"
	// "net/url"
)

type MediaService interface {
	SavePic(params domain_models.SavePicParams) (string, error)
}

type MediaServiceImpl struct {
	mediaRepository repository.MediaRepository
}

func NewMediaService(mediaRepository repository.MediaRepository) MediaService {
	return &MediaServiceImpl {
		mediaRepository: mediaRepository,
	}
}

func (m MediaServiceImpl) SavePic(params domain_models.SavePicParams) (string, error) {
	file := params.File
	generatedUrl, err := m.mediaRepository.SavePic(file)
	if err != nil {
		return "", err
	}
	return generatedUrl.String(), nil
}