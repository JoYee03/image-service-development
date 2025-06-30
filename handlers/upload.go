package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/JoYee03/image-service-development/storage"
)

var (
	storageBucket   string
	credentialsFile string
)

type UploadRequest struct {
	Content  string `json:"content"`  // Base64
	Type     string `json:"type"`     // MIME type
	Filename string `json:"filename"` // Optional
}

type UploadResponse struct {
	Path      string `json:"path"`
	PublicURL string `json:"public_url,omitempty"`
	Success   bool   `json:"success"`
}

func InitUploadHandler(bucket, credsFile string) {
	storageBucket = bucket
	credentialsFile = credsFile
}

func ImageUploadHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if isRateLimited(r) {
		sendJSONError(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	var req UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	if req.Content == "" {
		sendJSONError(w, "Base64 content required", http.StatusBadRequest)
		return
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[req.Type] {
		sendJSONError(w, "Unsupported image type", http.StatusBadRequest)
		return
	}

	data, err := base64.StdEncoding.Strict().DecodeString(req.Content)
	if err != nil {
		sendJSONError(w, "Invalid base64", http.StatusBadRequest)
		return
	}

	if len(data) > 10*1024*1024 { // 10MB limit
		sendJSONError(w, "Image too large", http.StatusBadRequest)
		return
	}

	if !isValidImage(data, req.Type) {
		sendJSONError(w, "Invalid image content", http.StatusBadRequest)
		return
	}

	ext := ".jpg"
	switch req.Type {
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	}

	filePath := generateFilePath(req.Filename, ext)

	storageService, err := storage.NewStorageService(ctx, storageBucket, credentialsFile)
	if err != nil {
		log.Printf("Storage init failed: %v", err)
		sendJSONError(w, "Storage error", http.StatusInternalServerError)
		return
	}

	publicURL, err := storageService.UploadImage(ctx, data, filePath, req.Type)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		sendJSONError(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, UploadResponse{
		Path:      filePath,
		PublicURL: publicURL,
		Success:   true,
	})
}

func generateFilePath(filename, ext string) string {
	if filename != "" {
		return filepath.Join("images", filename)
	}
	return filepath.Join("images", storage.GenerateUUID()+ext)
}

func isValidImage(data []byte, mimeType string) bool {
	if len(data) < 12 {
		return false
	}
	switch mimeType {
	case "image/jpeg":
		return bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF})
	case "image/png":
		return bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47})
	case "image/webp":
		return bytes.HasPrefix(data, []byte{0x52, 0x49, 0x46, 0x46}) &&
			bytes.HasPrefix(data[8:], []byte{0x57, 0x45, 0x42, 0x50})
	}
	return true
}

func isRateLimited(r *http.Request) bool {
	return false // Implement your rate limiting logic
}