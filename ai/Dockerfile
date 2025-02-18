# -------------------------------
# Stage 1: Build the binary
# -------------------------------
FROM golang:1.23.1 AS builder

    # Set the working directory inside the container
WORKDIR /app
    
    # Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download
    
    # Copy the rest of the source code
COPY . .
    
    # Build the Go binary.
    # Adjust the build command if your main package is located elsewhere.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd
    
    # -------------------------------
    # Stage 2: Create the runtime image
    # -------------------------------
FROM alpine:latest
    
    # Install CA certificates (if needed for HTTPS requests)
# RUN apk --no-cache add ca-certificates
    
    # Set working directory for the runtime image
WORKDIR /app
    
    # Copy the built binary from the builder stage
COPY --from=builder /app/app .
    
    # Copy configuration files if required
COPY --from=builder /app/config ./config
    
    # Expose the gRPC port (adjust if your app uses a different port)
EXPOSE 50051
    
    # Command to run the application
CMD ["./app"]
    