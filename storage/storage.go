package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"
	"crypto/rand"
	"encoding/hex"
	"net/url"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type StorageService struct {
	client *storage.Client
	bucket string
}

func NewStorageService(ctx context.Context, bucketName string, credentialsFile string) (*StorageService, error) {
	var client *storage.Client
	var err error

	if credentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		// Auto-discovers credentials in Cloud Functions/Cloud Run
		client, err = storage.NewClient(ctx) 
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &StorageService{
		client: client,
		bucket: bucketName,
	}, nil
}


func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *StorageService) UploadImage(ctx context.Context, data []byte, path string, contentType string) (string, error) {
	token := generateToken()

	wc := s.client.Bucket(s.bucket).Object(path).NewWriter(ctx)
	wc.ContentType = contentType
	wc.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": token,
	}

	if _, err := wc.Write(data); err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize upload: %w", err)
	}

	// ✅ Correct public URL with token
	publicURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		s.bucket,
		url.PathEscape(path),
		token,
	)

	return publicURL, nil
}

func (s *StorageService) DownloadImage(ctx context.Context, path string) ([]byte, error) {
	rc, err := s.client.Bucket(s.bucket).Object(path).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

func (s *StorageService) GenerateUploadPath(originalName string) string {
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".jpg"
	}
	return filepath.Join(
		"images",
		time.Now().Format("2006/01/02"),
		fmt.Sprintf("%d%s", time.Now().UnixNano(), ext),
	)
}

func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano()) // fallback
	}
	return hex.EncodeToString(b)
}