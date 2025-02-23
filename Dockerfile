# 1️ Use the official Golang image as a base
FROM golang:1.24-alpine AS builder

# 2️ Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 3️ Set the working directory
WORKDIR /app

# 4️ Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# 5️ Copy the source code
COPY . .

# 6️ Build the Go application
RUN go build -o agent ./main.go

# 7️ Use a minimal base image for running the binary
FROM alpine:latest

# 8️ Set working directory
WORKDIR /root/

# 9️ Copy the compiled binary from the builder stage
COPY --from=builder /app/agent .

# 10️ Copy the .env file for environment variables (if you have one)
COPY .env .

# 11 Run the agent by default (can be overridden by command)
ENTRYPOINT ["./agent"]