package storagegateway

import (
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

	return &StorageInstance{
		id:          "minio",
		minioClient: minioClient,
	}
}
