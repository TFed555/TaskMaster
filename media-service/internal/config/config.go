package config

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	minioClient *minio.Client
}

func NewMinioClient() (MinioClient, error) {

	if err := godotenv.Load("../.env", ".env.local"); err != nil {
		log.Fatal("Can't load .env file")
		return MinioClient{}, err
	}

	endpoint := "127.0.0.1:9000"
	accessKeyID, exists := os.LookupEnv("MINIO_ACCESS_KEY")
	if !exists {
		log.Print("Can't get token")
		return MinioClient{}, errors.New("Can't get token")
	}
	secretAccessKey, exists := os.LookupEnv("MINIO_SECRET_ACCESS_KEY")
		if !exists {
		log.Print("Can't get secret token")
		return MinioClient{}, errors.New("Can't get secret token")
	}

	useSSL := false

	log.Printf("AccessKey: %s, SecretKey: %s", accessKeyID, secretAccessKey)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Region: "us-east-1",
	})

	if err != nil {
		log.Fatalln(err)
		return MinioClient{}, err
	}

	ctx := context.Background()

	bucketName := "new-bucket"
    location := "us-east-1"

    err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
    if err != nil {
                // Check to see if we already own this bucket (which happens if you run this twice)
    exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
    if errBucketExists == nil && exists {
        log.Printf("We already own %s\n", bucketName)
    } else {
            log.Fatalln(err)
    }
    } else {
            log.Printf("Successfully created %s\n", bucketName)
   }

   	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:PutObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::new-bucket/*"],"Sid": ""}]}`

	err = minioClient.SetBucketPolicy(context.Background(), "new-bucket", policy)
	if err != nil {
		log.Println(err)
		return MinioClient{}, err
	}

	return MinioClient{
		minioClient: minioClient,
	}, nil

}

func (m MinioClient) UploadImage(file string) (string, error) {

	log.Printf("%#v\n", m.minioClient)
	ctx := context.Background()

	bucketName := "test-bucket"
        // location := "us-east-1"

    err := m.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
    if err != nil {
                // Check to see if we already own this bucket (which happens if you run this twice)
    exists, errBucketExists := m.minioClient.BucketExists(ctx, bucketName)
    if errBucketExists == nil && exists {
        log.Printf("We already own %s\n", bucketName)
    } else {
            log.Fatalln(err)
    }
    } else {
            log.Printf("Successfully created %s\n", bucketName)
   }

        objectName := "pidor.jpg"
		wd, _ := os.Getwd()
        filePath := filepath.Join(wd, "pidor.jpg")
        contentType := "application/octet-stream"

		log.Print(filePath)
		
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal("Файл не найден: ", err)
		}
		log.Printf("Размер файла: %d байт", fileInfo.Size())

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Ошибка при чтении файла: %v", err)
		}
		log.Printf("Прочитано %d байт", len(data))
		
        // Upload the test file with FPutObject
        info, err := m.minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
        if err != nil {
                log.Fatalln(err)
        }

        log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
		return "Success", nil
}

func (m MinioClient) GenerateLink(file string) (url.URL, error) {
	exists, err := m.minioClient.BucketExists(context.Background(), "new-bucket")
	if err != nil {
		return url.URL{}, errors.New("bucket check failed")
	}
	if !exists {
		return url.URL{}, errors.New("bucket does not exist")
	}

	expiry := time.Second * 24 * 60 * 60

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := 10
	ran_str := make([]byte, length)
	
	for i := 0; i < length; i++ {
		if r.Intn(2) == 0 {
			ran_str[i] = byte(65 + r.Intn(26))
		} else {
			ran_str[i] = byte(97 + r.Intn(26))
		}
	}

	objectname := string(ran_str)
	presignedURL, err := m.minioClient.PresignedPutObject(context.Background(), "new-bucket", objectname, expiry)
	if err != nil {
    	log.Println(err)
    	return url.URL{}, err
	}
	log.Println("Successfully generated presigned URL", presignedURL)

	return *presignedURL, nil
}