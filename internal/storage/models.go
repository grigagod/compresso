package storage

import "io"

type UploadInput struct {
	File        io.ReadSeeker
	Name        string
	Size        int64
	ContentType string
}
