package common

import "errors"

var (
	ErrObjectNotFound   = errors.New("object not found")
	ErrInstanceNotFound = errors.New("instance not found")
)
