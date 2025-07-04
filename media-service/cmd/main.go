package main


func main() {
		// //temporary, will be in media-service
	// if err := godotenv.Load("../.env", ".env.local"); err != nil {
	// 	log.Fatal("Can't load .env file")
	// }

	// endpoint := "127.0.0.1:9000"
	// accessKeyID, exists := os.LookupEnv("MINIO_ACCESS_KEY")
	// if !exists {
	// 	log.Print("Can't get token")
	// }
	// secretAccessKey, exists := os.LookupEnv("MINIO_SECRET_ACCESS_KEY")
	// 	if !exists {
	// 	log.Print("Can't get secret token")
	// }

	// useSSL := false

	// log.Printf("AccessKey: %s, SecretKey: %s", accessKeyID, secretAccessKey)

	// minioClient, err := minio.New(endpoint, &minio.Options{
	// 	Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	// 	Secure: useSSL,
	// })

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Printf("%#v\n", minioClient)
	// ctx := context.Background()

	//  	bucketName := "test-bucket"
    //     // location := "us-east-1"

    //     err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
    //     if err != nil {
    //             // Check to see if we already own this bucket (which happens if you run this twice)
    //             exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
    //             if errBucketExists == nil && exists {
    //                     log.Printf("We already own %s\n", bucketName)
    //             } else {
    //                     log.Fatalln(err)
    //             }
    //     } else {
    //             log.Printf("Successfully created %s\n", bucketName)
    //     }

    //     // Upload the test file
    //     // Change the value of filePath if the file is in another location
    //     objectName := "pidor.jpg"
	// 	wd, _ := os.Getwd()
    //     filePath := filepath.Join(wd, "pidor.jpg")
    //     contentType := "application/octet-stream"

	// 	log.Print(filePath)
		
	// 	fileInfo, err := os.Stat(filePath)
	// 	if err != nil {
	// 		log.Fatal("Файл не найден: ", err)
	// 	}
	// 	log.Printf("Размер файла: %d байт", fileInfo.Size())

	// 	data, err := os.ReadFile(filePath)
	// 	if err != nil {
	// 		log.Fatalf("Ошибка при чтении файла: %v", err)
	// 	}
	// 	log.Printf("Прочитано %d байт", len(data))
		
    //     // Upload the test file with FPutObject
    //     info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
    //     if err != nil {
    //             log.Fatalln(err)
    //     }

    //     log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}