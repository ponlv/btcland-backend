package minio

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOUploader struct {
	client *minio.Client
}

var client *MinIOUploader

func NewClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	client = &MinIOUploader{
		client: minioClient,
	}

	return minioClient, nil
}

func PutObject(bucket, key string, file io.Reader) (minio.UploadInfo, error) {
	object, err := client.client.PutObject(context.Background(), bucket, key, file, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return object, fmt.Errorf("failed to upload object to %v %w", client.client.EndpointURL(), err)
	}

	return object, nil
}

func GetObject(bucket, key string) (string, error) {
	object, err := client.client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get object %s: %w", key, err)
	}

	// We expect an error when reading back.
	buf, err := io.ReadAll(object)
	if err != nil {
		return "", fmt.Errorf("failed to read object %s: %w", key, err)
	}

	uEnc := base64.StdEncoding.EncodeToString(buf)

	return uEnc, nil
}

func DeleteObject(bucket, key string) error {
	err := client.client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}

	return nil
}

func EnsureBucketExists(bucketName string) error {
	exists, err := client.client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket %s exists: %w", bucketName, err)
	}

	if !exists {
		err = client.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	return nil
}

func GetClient() *minio.Client {
	return client.client
}
