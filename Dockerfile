# 1️⃣ Use the official Golang image as a builder
FROM golang:1.24-alpine AS builder

# 2️⃣ Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 3️⃣ Set the working directory
WORKDIR /app

# 4️⃣ Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# 5️⃣ Copy the source code
COPY . .

# 6️⃣ Build the Go application
RUN go build -o agent ./main.go

# 7️⃣ Use a minimal runtime base image
FROM alpine:latest

# 8️⃣ Set working directory
WORKDIR /root/

# 9️⃣ Copy the compiled binary from the builder stage
COPY --from=builder /app/agent .

# 🔟 Copy the environment variables and configuration files
COPY .env .
COPY configs/ ./configs/

# 1️⃣1️⃣ Set the entrypoint for running bots dynamically with config
ENTRYPOINT ["./agent"]