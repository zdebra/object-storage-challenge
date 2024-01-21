package storagegateway

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type StorageInstance struct {
	id          string
	minioClient *minio.Client
}

func (i *StorageInstance) String() string {
	return i.id
}

func newStorageInstance(endpoint, accessKeyID, secretAccessKey string) *StorageInstance {
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
		id:          endpoint,
		minioClient: minioClient,
	}
}

func ensureBucketIsCreated(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("check bucket existence: %w", err)
	}
	if exists {
		zap.L().Info("bucket already exists", zap.String("name", bucketName))
		return nil
	}
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}
	zap.L().Info("bucket created", zap.String("name", bucketName))
	return nil
}

func InitInstances(cfgs []InstanceCfg) []*StorageInstance {
	instances := []*StorageInstance{}
	for _, cfg := range cfgs {
		instances = append(instances, newStorageInstance(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey))
	}
	zap.L().Info("storage instances initialized", zap.Int("count", len(instances)))
	return instances
}
