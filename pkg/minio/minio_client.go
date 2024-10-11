package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client интерфейс для взаимодействия с Minio
type Client interface {
	InitMinio(endpoint string, user string, pass string, ssl bool, bucket string) error // Метод для инициализации подключения к Minio
	CreateOne(file FileDataType) (string, error)                                        // Метод для создания одного объекта в бакете Minio
	GetOne(objectID string) (string, error)                                             // Метод для получения одного объекта из бакета Minio
}

type minioClient struct {
	mc     *minio.Client
	bucket string
}

func NewMinioClient() Client {
	return &minioClient{}
}

// InitMinio подключается к Minio и создает бакет, если не существует
func (m *minioClient) InitMinio(endpoint string, user string, pass string, ssl bool, bucket string) error {
	m.bucket = bucket
	ctx := context.Background()

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, pass, ""),
		Secure: ssl,
	})
	if err != nil {
		return err
	}

	m.mc = client

	exists, err := m.mc.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		err = m.mc.SetBucketPolicy(ctx, bucket, fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Action": ["s3:GetObject"],
          "Effect": "Allow",
          "Principal": {
            "AWS": ["*"]
          },
          "Resource": ["arn:aws:s3:::%s/*"],
          "Sid": ""
        }
    ]
}`, bucket))
		if err != nil {
			return err
		}
	}
	return nil
}
