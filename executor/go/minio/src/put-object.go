package src

import (
	"log"

	minio "github.com/minio/minio-go"
)

func PutObject(endpoint string, accessKey string, secretKey string, ssl bool, objName string, bucketName string, filePath string, contentType string) error {
	minioClient, err := minio.New(endpoint, accessKey, secretKey, ssl)
	if err != nil {
		log.Fatalln(err)
	}

	//minioClient.SetCustomTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})

	found, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	if !found {
		err = minioClient.MakeBucket(bucketName, "ap-northeast-2")
	}

	_, err = minioClient.FPutObject(bucketName, objName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
