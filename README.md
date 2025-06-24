# iCARES Image Service API

A Go-based local server that provides APIs to upload images and apply watermarks using Sharp and Firebase Cloud Storage.

---

## 🔧 Setup

### 1. Requirements

* Go 1.21+
* Node.js (for Sharp image processing)
* Firebase Admin SDK (service account)
* Firebase project with Storage enabled

### Firebase Service Account JSON

This project requires a `firebase-service-account.json` file for local testing and deployment.

#### To generate your own `firebase-service-account.json`:
1. Go to the Firebase Console at https://console.firebase.google.com/
2. Select your project.
3. Navigate to **Project settings → Service accounts**.
4. Click **Generate new private key** under "Firebase Admin SDK".
5. Save the JSON file as `firebase-service-account.json` at the project’s root.

### 2. Folder Structure

```
iCARES/
├── main.go
├── watermark.js
├── firebase-service-account.json
├── upload-image.ps1
├── watermark-image.ps1
├── blank_base64.txt
├── icares_base64.txt
├── go.mod
├── go.sum
```

### 3. Install Node Sharp (for watermark)

```bash
npm install sharp
```

### 4. Run the Go Server

```bash
go run main.go
```

Server will start at `http://localhost:8080`

---

## 🧚‍♂️ Testing

You can test the API using PowerShell scripts:

.\upload-image.ps1        # uploads both the main image and watermark

.\watermark-image.ps1     # applies watermark to the uploaded image

---

## 🌐 API Endpoints

| Function        | Method | Endpoint              | Description                           |
| --------------- | ------ | --------------------- | ------------------------------------- |
| Upload Image    | POST   | `/testImageUpload`    | Uploads blank_base64 image to Storage |
| Apply Watermark | POST   | `/testWatermarkImage` | Merges watermark + original           |

---

## 🚀 GitHub Actions Deployment (Optional for Go)

To deploy Firebase Functions (if you use JavaScript), you can set up a GitHub Action workflow.

### `.github/workflows/deploy.yml`

```yaml
name: Deploy Firebase Functions

on:
  push:
    tags:
      - '*'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - run: npm install -g firebase-tools
      - run: firebase deploy --only functions --token ${{ secrets.FIREBASE_TOKEN }}
```

### Add Firebase Token

Go to GitHub > Settings > Secrets and add `FIREBASE_TOKEN`.

### Trigger Deployment

```bash
git tag v1.0.0
git push origin v1.0.0
```

---

## 📦 Dependencies

* Go 1.21+
* Node.js 18+ with Sharp (`npm install sharp`)
* Firebase Admin SDK (service account)
