package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spacelift-io/homework-object-storage/common"
)

type StorageService interface {
	PutObject(ctx context.Context, id string, dataStream io.Reader, dataSize int64) error
	GetObject(ctx context.Context, id string) (io.Reader, int64, error)
}

type StorageGatewayAPI struct {
	service StorageService
}

// (GET /object/{id})
func (sga *StorageGatewayAPI) ObjectGet(w http.ResponseWriter, r *http.Request, id string) {
	dataStream, size, err := sga.service.GetObject(r.Context(), id)
	if err == common.ErrObjectNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, dataStream)
}

// (PUT /object/{id})
func (sga *StorageGatewayAPI) ObjectPut(w http.ResponseWriter, r *http.Request, id string) {
	err := sga.service.PutObject(r.Context(), id, r.Body, r.ContentLength)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
