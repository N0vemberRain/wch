package adapters

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalAvatarStorage struct {
	baseDir string
}

func NewLocalAvatarStorage(baseDir string) (*LocalAvatarStorage, error) {
	if baseDir == "" {
		return nil, fmt.Errorf("NewLocalAvatarStorage: base dir path is empty")
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("NewLocalAvatarStorage: ", err)
	}

	return &LocalAvatarStorage{baseDir: baseDir}, nil
}

func (s *LocalAvatarStorage) Save(ctx context.Context, userID uuid.UUID, data []byte) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	key := userID.String() + ".png"

	filename := filepath.Join(s.baseDir, key)

	log.Printf("Trying to store an avatar: %s\n", filename)

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("save avatar: %w", err)
	}

	return key, nil
}

func (s *LocalAvatarStorage) Get(ctx context.Context, key string) ([]byte, error) {
	filename := filepath.Join(s.baseDir, key)
	return os.ReadFile(filename)
}

func (s *LocalAvatarStorage) Delete(ctx context.Context, key string) error {
	return nil
}
