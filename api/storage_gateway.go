package api

import "net/http"

type StorageGatewayAPI struct {
}

// (GET /object/{id})
func (sga *StorageGatewayAPI) ObjectGet(w http.ResponseWriter, r *http.Request, id string) {
	panic("object get")
}

// (PUT /object/{id})
func (sga *StorageGatewayAPI) ObjectPut(w http.ResponseWriter, r *http.Request, id string) {
	panic("object put")
}
