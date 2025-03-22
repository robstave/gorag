# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Cache dependencies

ENV  CGO_ENABLED=1 GOOS=linux GOARCH=amd64 
 
# Install necessary packages for building (gcc, musl-dev, etc.)
RUN apk update && apk add --no-cache git gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the binary
COPY . .

RUN go build -o service ./cmd/main/main.go


# Debug output to verify the binary was created
RUN ls -la /app/

# Final stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the compiled binary from the builder stage
COPY --from=builder /app/service .

# Debug output to verify the binary was copied
RUN ls -la /app/

# Ensure the binary is executable
RUN chmod +x /app/service

# Expose the port defined in the app (default 8711)
EXPOSE 8711

# Run the binary
CMD ["/app/service"]