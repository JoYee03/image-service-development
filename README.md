# iCARES Image Service API

A Go + Node.js Cloud Run API that supports image uploading and watermarking via Sharp and Firebase Cloud Storage.

---

## Requirements

* Go 1.21+
* Node.js (for `sharp`)
* Firebase project with Cloud Storage
* Firebase Admin SDK JSON key

---

## Folder Structure

```
iCARES/
├── main.go
├── watermark.js
├── firebase-service-account.json
├── Dockerfile
├── upload-image.ps1
├── upload-image2.ps1
├── watermark-image.ps1
├── blank_base64.txt
├── icares_base64.txt
├── go.mod
├── go.sum
```

---

## Local Setup

### 1. Install Go & Node

Install [Go](https://go.dev/dl) and [Node.js](https://nodejs.org/) (v18+).
Install `sharp` using:

```bash
npm install sharp
```

### 2. Add Firebase Admin SDK JSON

* Go to **Firebase Console** → **Project Settings** → **Service accounts**
* Click **Generate new private key**
* Save it as: `firebase-service-account.json`

---

##  Run Locally

```bash
go run main.go
```

Server runs on: `http://localhost:8080`

---

## Test with PowerShell

1. **Upload image + watermark**

```powershell
.\upload-image.ps1
.\upload-image2.ps1
```

2. **Apply watermark**

```powershell
.\watermark-image.ps1
```

---

## API Endpoints (Cloud Run)

| Endpoint              | Method | Description                           |
| --------------------- | ------ | ------------------------------------- |
| `/testImageUpload`    | POST   | Uploads an image to Firebase Storage  |
| `/testWatermarkImage` | POST   | Applies watermark over uploaded image |

Cloud Run URL:
`https://image-service-development-735683043266.asia-southeast1.run.app`

---

## Docker Deployment (Cloud Run)

### 1. Build and Submit

```bash
gcloud builds submit --tag gcr.io/image-service-development/image-service
```

### 2. Deploy to Cloud Run

```bash
gcloud run deploy image-service-development \
  --image gcr.io/image-service-development/image-service \
  --platform managed \
  --region asia-southeast1 \
  --allow-unauthenticated
```

---

## GitHub Actions Deployment

Automatically deploy when a Git tag is pushed.

### `.github/workflows/deploy.yml`

```yaml
name: Deploy to Cloud Run on Tag

on:
  push:
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Google Cloud CLI
        uses: google-github-actions/setup-gcloud@v2
        with:
          credentials: ${{ secrets.GCP_SA_KEY }}
          project_id: ${{ secrets.GCP_PROJECT_ID }}
          export_default_credentials: true

      - name: Build and push Docker image
        run: |
          gcloud builds submit --tag gcr.io/${{ secrets.GCP_PROJECT_ID }}/image-service

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy image-service-development \
            --image gcr.io/${{ secrets.GCP_PROJECT_ID }}/image-service \
            --platform managed \
            --region asia-southeast1 \
            --allow-unauthenticated
```

### Secrets to Add

| Secret Name      | Description                                                                             |
| ---------------- | --------------------------------------------------------------------------------------- |
| `GCP_SA_KEY`     | Contents of your service account JSON (base64-encoded or plain string in GitHub Secret) |
| `GCP_PROJECT_ID` | Your Google Cloud Project ID                                                            |

### Trigger Deployment

```bash
git tag v1.0.0
git push origin v1.0.0
```