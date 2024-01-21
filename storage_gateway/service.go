package storagegateway

import (
	"context"
	"fmt"
	"io"

	"github.com/buraksezer/consistent"
	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
	"github.com/spacelift-io/homework-object-storage/common"
	"go.uber.org/zap"
)

const bucketName = "mybucket"

type Service struct {
	keyLocator *consistent.Consistent
}

func NewService(instances ...*StorageInstance) *Service {
	if len(instances) == 0 {
		panic("no storage instance provided")
	}
	cfg := consistent.Config{
		PartitionCount: len(instances),
		Load:           0.1,
		Hasher:         hasher{},
	}

	keyLocator := consistent.New(lo.Map(instances, func(x *StorageInstance, _ int) consistent.Member { return x }), cfg)
	return &Service{keyLocator}
}

func (s *Service) PutObject(ctx context.Context, id string, dataStream io.Reader, dataSize int64) error {
	instance, err := s.findInstance(id)
	if err != nil {
		return fmt.Errorf("find instance: %w", err)
	}

	_, err = instance.minioClient.PutObject(ctx, bucketName, id, dataStream, dataSize, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	zap.L().Info("object put", zap.String("id", id), zap.String("instance", instance.String()))

	return nil
}

func (s *Service) GetObject(ctx context.Context, id string) (io.Reader, int64, error) {
	instance, err := s.findInstance(id)
	if err != nil {
		return nil, 0, fmt.Errorf("find instance: %w", err)
	}

	obj, err := instance.minioClient.GetObject(ctx, bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, fmt.Errorf("get object: %w", err)
	}

	stat, err := obj.Stat()
	if err != nil && err.Error() == "The specified key does not exist." {
		return nil, 0, common.ErrObjectNotFound
	}
	if err != nil {
		return nil, 0, fmt.Errorf("stat object: %w", err)
	}

	zap.L().Info("object get", zap.String("id", id), zap.String("instance", instance.String()))
	return obj, stat.Size, nil
}

func (s *Service) findInstance(id string) (*StorageInstance, error) {
	instance := s.keyLocator.LocateKey([]byte(id))
	if instance == nil {
		return nil, common.ErrInstanceNotFound
	}
	return instance.(*StorageInstance), nil
}
