# ========== Build Stage ==========
FROM golang:1.24 AS builder
WORKDIR /app

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go binary
RUN go build -o image-service main.go

# ========== Deploy Stage ==========
FROM debian:12-slim

# Install node.js and npm so we can run watermark.js with sharp
RUN apt-get update && apt-get install -y nodejs npm \
    && rm -rf /var/lib/apt/lists/*

# Set work dir
WORKDIR /app

# Copy Go binary and other files
COPY --from=builder /app/image-service /app/image-service
COPY firebase-service-account.json /app/firebase-service-account.json
COPY watermark.js /app/watermark.js

# Install sharp
RUN npm init -y && npm install sharp

# Expose port (Cloud Run expects this)
ENV PORT=8080

# Run the Go binary
CMD ["/app/image-service"]
