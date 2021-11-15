//go:generate mockgen -source storage.go -destination mock/storage_mock.go -package mock
package storage

import (
	"context"
	"io"
)

type Storage interface {
	PutObject(ctx context.Context, file io.Reader, fileName, fileType string) error
	GetObject(ctx context.Context, fileName string) (io.ReadCloser, error)
	GetDownloadURL(fileName string) (string, error)
}
