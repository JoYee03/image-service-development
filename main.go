package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io"
	"os/exec"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type UploadRequest struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	Filename string `json:"filename"`
}

type UploadResponse struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
}

type WatermarkRequest struct {
	ImagePath     string `json:"image_path"`
	WatermarkPath string `json:"watermark_path"`
	Filename string `json:"filename"`
}

type WatermarkResponse struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
}

// firebase initialization
func initFirebase(ctx context.Context) (*firebase.App, error) {
	opt := option.WithCredentialsFile("firebase-service-account.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("Firebase init error: %v", err)
	}
	return app, nil
}

// image upload handler
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	app, err := initFirebase(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		http.Error(w, "Storage client error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bucket, err := storageClient.Bucket("image-service-development.firebasestorage.app")
	if err != nil {
		http.Error(w, "Bucket error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		http.Error(w, "Base64 decode error: "+err.Error(), http.StatusBadRequest)
		return
	}

		objectPath := req.Filename
	if objectPath == "" {
		objectPath = fmt.Sprintf("image/%d/uploaded-%d.png", time.Now().Year(), time.Now().Unix())
	}
	writer := bucket.Object(objectPath).NewWriter(ctx)
	defer writer.Close()
	writer.ContentType = req.Type

	if _, err := writer.Write(data); err != nil {
		http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(UploadResponse{
		Path:    objectPath,
		Success: true,
	})
}

// watermark handler
func watermarkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	app, err := initFirebase(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		http.Error(w, "Storage client error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bucket, err := storageClient.Bucket("image-service-development.firebasestorage.app")
	if err != nil {
		http.Error(w, "Bucket error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req WatermarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// download image and watermark
	imageData, err := downloadFile(ctx, bucket, req.ImagePath)
	if err != nil {
		http.Error(w, "Error downloading image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	watermarkData, err := downloadFile(ctx, bucket, req.WatermarkPath)
	if err != nil {
		http.Error(w, "Error downloading watermark: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// save to temporary files
	if err := os.WriteFile("temp-image.png", imageData, 0644); err != nil {
		http.Error(w, "Error saving temp image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove("temp-image.png")

	if err := os.WriteFile("temp-watermark.png", watermarkData, 0644); err != nil {
		http.Error(w, "Error saving temp watermark: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove("temp-watermark.png")

	cmd := exec.Command("node", "watermark.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		http.Error(w, "Watermarking failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove("output.png")

	// upload watermarked image
	outputData, err := os.ReadFile("output.png")
	if err != nil {
		http.Error(w, "Error reading output: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outputPath := fmt.Sprintf("image/%d/watermarked-%d.png", time.Now().Year(), time.Now().Unix())
	writer := bucket.Object(outputPath).NewWriter(ctx)
	defer writer.Close()
	writer.ContentType = "image/png"

	if _, err := writer.Write(outputData); err != nil {
		http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(WatermarkResponse{
		Path:    outputPath,
		Success: true,
	})

}

// helper: download file from firebase storage
func downloadFile(ctx context.Context, bucket *storage.BucketHandle, path string) ([]byte, error) {
	reader, err := bucket.Object(path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// main function
func main() {
    // handler for the home page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to iCARES Image Service \n\nAvailable endpoints:\n- POST /testImageUpload\n- POST /testWatermarkImage\n")
    })

    http.HandleFunc("/testImageUpload", uploadHandler)
    http.HandleFunc("/testWatermarkImage", watermarkHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server started on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}