package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/JoYee03/image-service-development/storage"
)

type WatermarkRequest struct {
	ImagePath     string `json:"image_path"`
	WatermarkPath string `json:"watermark_path"`
}

type WatermarkResponse struct {
	WatermarkedPath string `json:"watermarked_path"`
	Success         bool   `json:"success"`
}

func WatermarkHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	var req WatermarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ImagePath == "" || req.WatermarkPath == "" {
		sendJSONError(w, "Both paths required", http.StatusBadRequest)
		return
	}

	storageService, err := storage.NewStorageService(ctx, storageBucket, credentialsFile)
	if err != nil {
		sendJSONError(w, "Storage error", http.StatusInternalServerError)
		return
	}

	var imgData, watermarkData []byte
	var downloadErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		data, err := storageService.DownloadImage(ctx, req.ImagePath)
		if err != nil {
			downloadErr = err
			return
		}
		imgData = data
	}()

	go func() {
		defer wg.Done()
		data, err := storageService.DownloadImage(ctx, req.WatermarkPath)
		if err != nil {
			downloadErr = err
			return
		}
		watermarkData = data
	}()

	wg.Wait()
	if downloadErr != nil {
		sendJSONError(w, "Download failed", http.StatusInternalServerError)
		return
	}

	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		sendJSONError(w, "Invalid image", http.StatusBadRequest)
		return
	}

	watermark, _, err := image.Decode(bytes.NewReader(watermarkData))
	if err != nil {
		sendJSONError(w, "Invalid watermark", http.StatusBadRequest)
		return
	}

	resizedWatermark := imaging.Resize(watermark, img.Bounds().Dx()/4, 0, imaging.Lanczos)
	watermarked := imaging.Clone(img)
	spacing := resizedWatermark.Bounds().Dy() / 2

	for y := 0; y < img.Bounds().Dy(); y += resizedWatermark.Bounds().Dy() + spacing {
		for x := 0; x < img.Bounds().Dx(); x += resizedWatermark.Bounds().Dx() + spacing {
			watermarked = imaging.Overlay(watermarked, resizedWatermark, image.Pt(x, y), 0.5)
		}
	}

	var buf bytes.Buffer
	switch format {
	case "png":
		err = png.Encode(&buf, watermarked)
	default:
		err = jpeg.Encode(&buf, watermarked, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		sendJSONError(w, "Encoding failed", http.StatusInternalServerError)
		return
	}

	outputPath := filepath.Join(filepath.Dir(req.ImagePath), "watermarked", filepath.Base(req.ImagePath))
	publicURL, err := storageService.UploadImage(ctx, buf.Bytes(), outputPath, "image/"+format)
	if err != nil {
		sendJSONError(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, WatermarkResponse{
		WatermarkedPath: publicURL,
		Success:         true,
	})
}