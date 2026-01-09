# Build stage - Use Harbor-cached base image
FROM harbor.dataknife.net/dockerhub/library/golang:1.23.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/unifi-manager-mcp ./cmd

# Runtime stage - Use Harbor-cached base image
FROM harbor.dataknife.net/dockerhub/library/alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bin/unifi-manager-mcp .

# Copy .env.example for reference
COPY .env.example .

# Expose default HTTP port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8000/health || exit 1

# Default to HTTP transport
ENV MCP_TRANSPORT=http
ENV MCP_HTTP_ADDR=:8000

# Run the server
CMD ["./unifi-manager-mcp"]
