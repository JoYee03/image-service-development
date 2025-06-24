package function

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "time"
)

type WatermarkRequest struct {
    ImagePath     string `json:"image_path"`
    WatermarkPath string `json:"watermark_path"`
    Filename      string `json:"filename"`
}
type WatermarkResponse struct {
    Path    string `json:"path"`
    Success bool   `json:"success"`
}

// Cloud Function entry point
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

    os.WriteFile("/tmp/temp-image.png", imageData, 0644)
    defer os.Remove("/tmp/temp-image.png")
    os.WriteFile("/tmp/temp-watermark.png", watermarkData, 0644)
    defer os.Remove("/tmp/temp-watermark.png")

    cmd := exec.Command("node", "watermark.js")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        http.Error(w, "Watermarking failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer os.Remove("/tmp/output.png")

    outputData, err := os.ReadFile("/tmp/output.png")
    if err != nil {
        http.Error(w, "Error reading output: "+err.Error(), http.StatusInternalServerError)
        return
    }

    outputPath := fmt.Sprintf("image/%d/watermarked-%d.png", time.Now().Year(), time.Now().Unix())
    writer := bucket.Object(outputPath).NewWriter(ctx)
    writer.ContentType = "image/png"
    if _, err := writer.Write(outputData); err != nil {
        http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    writer.Close()

    json.NewEncoder(w).Encode(WatermarkResponse{Path: outputPath, Success: true})
}
