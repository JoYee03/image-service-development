# ---------- Build Stage ----------
FROM golang:1.24.4 AS builder
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the Go binary
COPY . .
RUN go build -o image-service main.go

# ---------- Runtime Stage ----------
FROM debian:12-slim
WORKDIR /app

# Install Node.js and Sharp
RUN apt-get update && \
    apt-get install -y nodejs npm python3 g++ make && \
    npm install sharp && \
    apt-get remove -y python3 g++ make && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy Go binary and other files
COPY --from=builder /app/image-service .
COPY watermark.js .
COPY firebase-service-account.json .

# Set port
ENV PORT=8080

# Start the Go binary
CMD ["./image-service"]