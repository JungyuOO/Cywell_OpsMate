package appserver

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
)

type StoredObject struct {
	URI       string
	SizeBytes int64
}

type DocumentStorage interface {
	Store(ctx context.Context, documentID string, filename string, reader io.Reader) (StoredObject, error)
}

type LocalDocumentStorage struct {
	BasePath string
}

func (s LocalDocumentStorage) Store(ctx context.Context, documentID string, filename string, reader io.Reader) (StoredObject, error) {
	if s.BasePath == "" {
		return StoredObject{}, fmt.Errorf("base path is required")
	}
	size, err := io.Copy(io.Discard, reader)
	if err != nil {
		return StoredObject{}, err
	}
	return StoredObject{
		URI:       filepath.ToSlash(filepath.Join(s.BasePath, documentID, filepath.Base(filename))),
		SizeBytes: size,
	}, nil
}
