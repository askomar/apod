package minio

import (
	"bytes"
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileDataType struct {
	FileName string
	Data     []byte
}

type OperationError struct {
	ObjectID string
	Error    error
}

// CreateOne создает один объект в бакете Minio.
// Метод принимает структуру fileData, которая содержит имя файла и его данные.
// В случае успешной загрузки данных в бакет, метод возвращает nil, иначе возвращает ошибку.
func (m *minioClient) CreateOne(file FileDataType) (string, error) {
	reader := bytes.NewReader(file.Data)

	_, err := m.mc.PutObject(context.Background(), m.bucket, file.FileName, reader, int64(len(file.Data)), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	url, err := m.mc.PresignedGetObject(context.Background(), m.bucket, file.FileName, time.Second*24*60*60, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// GetOne получает один объект из бакета Minio по его идентификатору.
// Он принимает строку `objectID` в качестве параметра и возвращает срез байт данных объекта и ошибку, если такая возникает.
func (m *minioClient) GetOne(objectID string) (string, error) {
	url, err := m.mc.PresignedGetObject(context.Background(), m.bucket, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
