package services

import (
	// "log"
	"media-service/internal/pkg/domain_models"
	"media-service/internal/repository"
	// "net/url"
)

type MediaService interface {
	GenerateLinkPic() (domain_models.File, error)
	SaveAvatarPic(params domain_models.ConfirmSaveParams) (bool, error)
	GetAvatarPic(img_name string) (string, error)
}

type MediaServiceImpl struct {
	mediaRepository repository.MediaRepository
}

func NewMediaService(mediaRepository repository.MediaRepository) MediaService {
	return &MediaServiceImpl {
		mediaRepository: mediaRepository,
	}
}

func (m MediaServiceImpl) GenerateLinkPic() (domain_models.File, error) {
	name, generatedUrl, err := m.mediaRepository.GenerateLink()
	if err != nil {
		return domain_models.File{}, err
	}
	newFile := domain_models.File{
		Name: name,
		Link: generatedUrl.String(),
	}
	return newFile, nil
}

func (m MediaServiceImpl) SaveAvatarPic(params domain_models.ConfirmSaveParams) (bool, error) {
	name, id := params.Name, params.UserID
	success, err := m.mediaRepository.SaveAvatarPic(name, id)
	if err != nil {
		return false, err
	}
	return success>0, nil
}	

func (m MediaServiceImpl) GetAvatarPic(img_name string) (string, error) {
	imgURL, err := m.mediaRepository.GetAvatarPic(img_name)
	if err != nil {
		return "", err
	}
	return imgURL, nil
}