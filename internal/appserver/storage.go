package appserver

import (
	"context"
	"fmt"
	"io"
	"os"
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

	targetDir := filepath.Join(s.BasePath, documentID)
	if err := os.MkdirAll(targetDir, 0o700); err != nil {
		return StoredObject{}, err
	}

	targetPath := filepath.Join(targetDir, filepath.Base(filename))
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return StoredObject{}, err
	}
	defer file.Close()

	size, err := io.Copy(file, reader)
	if err != nil {
		return StoredObject{}, err
	}
	return StoredObject{
		URI:       filepath.ToSlash(targetPath),
		SizeBytes: size,
	}, nil
}
