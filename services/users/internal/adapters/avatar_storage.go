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

func (s *LocalAvatarStorage) Save(ctx context.Context, ownerID uuid.UUID, data []byte) (string, error) {
	log.Println("LocalAvatarStorage.Save")
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	key := ownerID.String() + ".png"

	filename := filepath.Join(s.baseDir, key)

	log.Printf("Trying to store an avatar: %s\n", filename)

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("save avatar: %w", err)
	}

	return key, nil
}

func (s *LocalAvatarStorage) GetByKey(ctx context.Context, key string) ([]byte, error) {
	filename := filepath.Join(s.baseDir, key)
	log.Printf("Trying to get an avatar: %s\n", filename)
	return os.ReadFile(filename)
}

func (s *LocalAvatarStorage) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]byte, error) {
	filename := filepath.Join(s.baseDir, ownerID.String()+".png")
	log.Printf("Trying to get an avatar: %s\n", filename)
	return os.ReadFile(filename)
}

func (s *LocalAvatarStorage) Delete(ctx context.Context, key string) error {
	return nil
}
