package repository

import (
	"fmt"
	"log"
	"media-service/internal/config"
	"net/url"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MediaRepository struct {
	db *sqlx.DB
	minioClient config.MinioClient
}

func NewMediaRepository(db *sqlx.DB, minioClient config.MinioClient) MediaRepository {
	return MediaRepository{
		db: db,
		minioClient: minioClient,
	}
}

func (r MediaRepository) GenerateLink() (string, url.URL, error) {
	newUUID := uuid.New()
	fmt.Println("Новый UUID версии 4: ", newUUID)
	uuidString := newUUID.String()
	name, generatedUrl, err := r.minioClient.GenerateLink(uuidString)
	if err != nil {
		log.Fatal("Can't generate url")
		return "",url.URL{}, err
	}

	return name, generatedUrl, nil
}


//подумать как изолировать в auth-service
func (r MediaRepository) SaveAvatarPic(name string, id uint) (int, error) {
	const op ="repository.media_repository.SaveAvatarPic"
	
	query := `UPDATE auth.users SET img_path = $1 
				WHERE ID = $2 RETURNING ID`
	var updatedId int
	err := r.db.QueryRow(query, name, id).Scan(&updatedId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return updatedId, nil
}

func (r MediaRepository) GetAvatarPic(img_name string) (string, error) {
	const op = "repository.media_repository.GetAvatarPic"

	imgUrl, err := r.minioClient.GetAvatarPic(img_name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	imgUrlString := imgUrl.String()
	return imgUrlString, nil
}