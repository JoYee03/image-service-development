# iCARES Internship Test – Image Upload & Watermarking API

This project implements an image upload and watermarking API using the Go programming language. Images are uploaded to Firebase Storage and can be watermarked using an overlay function.

---

## Features

- Upload base64-encoded images to Firebase Storage via `/testImageUpload`
- Apply a watermark to an existing image via `/testWatermarkImage`
- Return public Firebase Storage URLs for uploaded and processed images
- Uses Google Cloud Run for deployment
- Firebase Storage is used to store uploaded and watermarked images in the `images/` and `images/watermarked/` folders

---

## Tech Stack

| Component       | Technology                          |
|----------------|--------------------------------------|
| Language        | Go (Golang)                         |
| Deployment      | Google Cloud Run                    |
| Storage         | Firebase Storage                    |
| Automation      | PowerShell scripts (`upload.ps1`, `watermark.ps1`) |

---

## Deployment Instructions

This project is deployed using **Google Cloud Run** with automatic containerization via `--source`.

### Prerequisites

- Google Cloud SDK installed (`gcloud`)
- Firebase project with billing enabled
- Firebase Storage enabled
- Service account JSON key (`firebase-service-account.json`)

### Environment Variables

Set the following environment variables (via `.env`, manually, or Cloud Run):

```
FIREBASE_BUCKET=celestial-geode-464418-p8.appspot.com
PORT=8080
GOOGLE_APPLICATION_CREDENTIALS=firebase-service-account.json
```

### How to Get a Firebase Service Account Key

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project (e.g., `image-service-development`)
3. Click the gear icon → **Project settings**
4. Open the **Service accounts** tab
5. Click **Generate new private key**
6. Save the downloaded file as `firebase-service-account.json` in the project root

Ensure that this key has permission to access Firebase Storage (default editor role works).

---

### Note About Firebase Bucket

The `FIREBASE_BUCKET` environment variable in this project is set to:

```
celestial-geode-464418-p8.appspot.com
```

This is the default Firebase Storage bucket for my project. If you are testing in a different Firebase project, replace this value with **your own bucket name**, which usually follows the format:

```
<your-project-id>.appspot.com
```

### Deployment Command

Use the following command in Command Prompt or PowerShell:

```cmd
gcloud run deploy image-service-development --source . --region asia-southeast1 --allow-unauthenticated --update-env-vars FIREBASE_BUCKET=celestial-geode-464418-p8.appspot.com
```

---

## API Endpoints

### 1. Image Upload

**POST** `/testImageUpload`

**Request Body:**

```json
{
  "content": "<base64-encoded-image>",
  "type": "image/jpeg",
  "filename": "test.jpg"
}
```

**Response:**

```json
{
  "path": "images/test.jpg",
  "public_url": "https://firebasestorage.googleapis.com/...",
  "success": true
}
```

**Example:**

```powershell
.\upload.ps1 -Base64Path "test.txt" -Filename "test.jpg"
```

---

### 2. Watermarking

**POST** `/testWatermarkImage`

**Request Body:**

```json
{
  "image_path": "images/test.jpg",
  "watermark_path": "images/watermark.jpg"
}
```

**Response:**

```json
{
  "watermarked_path": "https://firebasestorage.googleapis.com/...",
  "success": true
}
```

**Example:**

```powershell
.\watermark.ps1 -ImagePath "images/test.jpg" -WatermarkPath "images/watermark.jpg"
```

---

## Project Structure

```
.
├── main.go
├── handlers/
│   ├── upload.go
│   ├── watermark.go
│   └── helpers.go
├── storage/
│   └── storage.go
├── upload.ps1
├── watermark.ps1
├── firebase-service-account.json
└── go.mod
└── go.sum
└── .env
└── README.md
```

---

## Notes

- Go was chosen for this project even though Firebase Functions do not natively support Go.
- Google Cloud Run was used to deploy the Go server in a serverless environment.
- Firebase Storage is integrated via the official Google Cloud Storage SDK.
- All uploaded and processed images are publicly accessible via signed URLs returned in API responses.
