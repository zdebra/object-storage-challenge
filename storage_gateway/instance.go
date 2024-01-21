package storagegateway

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageInstance struct {
	id          string
	minioClient *minio.Client
}

func (i *StorageInstance) String() string {
	return i.id
}

func NewStorageInstance(endpoint, accessKeyID, secretAccessKey string) *StorageInstance {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		panic(err.Error())
	}

	err = ensureBucketIsCreated(context.Background(), minioClient, bucketName)
	if err != nil {
		panic(err.Error())
	}

	return &StorageInstance{
		id:          "minio",
		minioClient: minioClient,
	}
}

func ensureBucketIsCreated(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("check bucket existence: %w", err)
	}
	if exists {
		return nil
	}
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}
	return nil
}
