package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	firebase "firebase.google.com/go/v4"
)

type UploadRequest struct {
	Content  string `json:"content"`
	Type     string `json:"type"`
	Filename string `json:"filename"`
}

type UploadResponse struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
}

type WatermarkRequest struct {
	ImagePath     string `json:"image_path"`
	WatermarkPath string `json:"watermark_path"`
	Filename      string `json:"filename"`
}

type WatermarkResponse struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
}

func initFirebase(ctx context.Context) (*firebase.App, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase: %w", err)
	}
	return app, nil
}

func TestImageUpload(w http.ResponseWriter, r *http.Request) {
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
		objectPath = fmt.Sprintf("image/%d/uploaded-%d.jpg", time.Now().Year(), time.Now().Unix())
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

func TestWatermarkImage(w http.ResponseWriter, r *http.Request) {
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

	if err := os.WriteFile("temp-image.jpg", imageData, 0644); err != nil {
		http.Error(w, "Error saving temp image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove("temp-image.jpg")

	if err := os.WriteFile("temp-watermark.jpg", watermarkData, 0644); err != nil {
		http.Error(w, "Error saving temp watermark: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove("temp-watermark.jpg")

	cmd := exec.Command("node", "watermark.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		http.Error(w, "Watermarking failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove("output.jpg")

	outputData, err := os.ReadFile("output.jpg")
	if err != nil {
		http.Error(w, "Error reading output: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outputPath := req.Filename
	if outputPath == "" {
		outputPath = fmt.Sprintf("image/%d/watermarked-%d.jpg", time.Now().Year(), time.Now().Unix())
	}

	writer := bucket.Object(outputPath).NewWriter(ctx)
	defer writer.Close()
	writer.ContentType = "image/jpg"

	if _, err := writer.Write(outputData); err != nil {
		http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(WatermarkResponse{
		Path:    outputPath,
		Success: true,
	})
}

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

func main() {
	http.HandleFunc("/testImageUpload", TestImageUpload)
	http.HandleFunc("/testWatermarkImage", TestWatermarkImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
