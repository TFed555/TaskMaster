package repository

import (
	"log"
	"media-service/internal/config"
	"net/url"

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

func (r MediaRepository) SavePic(file string) (url.URL, error) {
	//TODO: генерировать по айди а не названию
	generatedUrl, err := r.minioClient.GenerateLink(file)
	if err != nil {
		log.Fatal("Can't generate url")
		return url.URL{}, err
	}

	return generatedUrl, nil
}