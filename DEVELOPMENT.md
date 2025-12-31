# UniFi Manager MCP - Development Guide

## Project Structure

```
unifi-manager-mcp/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── mcp/
│   │   └── server.go        # MCP server implementation
│   └── unifi/
│       └── client.go        # UniFi Manager API client
├── go.mod                   # Go module definition
├── README.md                # Project documentation
├── LICENSE                  # MIT License
├── CONTRIBUTING.md          # Contributing guidelines
├── CODE_OF_CONDUCT.md       # Community guidelines
└── CHANGELOG.md             # Version history
```

## Development

### Prerequisites

- Go 1.23.2 or higher
- UniFi Site Manager API Key

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/surrealwolf/unifi-manager-mcp.git
   cd unifi-manager-mcp
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build -o bin/unifi-manager-mcp ./cmd
   ```

4. Create `.env` file:
   ```bash
   cat > .env << EOF
   UNIFI_API_KEY=your-api-key
   LOG_LEVEL=debug
   MCP_TRANSPORT=stdio
   EOF
   ```

5. Run the server:
   ```bash
   ./bin/unifi-manager-mcp
   ```

### API Endpoints

The UniFi Manager API provides the following endpoints:

- `GET /v1/sites` - List all sites
- `GET /v1/hosts` - List all hosts
- `GET /v1/devices` - List all devices
- `GET /v1/deployments` - List all deployments

### MCP Tools Implementation

Each API endpoint is exposed as an MCP tool. Tools are registered in `internal/mcp/server.go` with proper input validation and error handling.

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./...
```

### Code Quality

Follow Go best practices:
- Use `gofmt` for formatting
- Run `go vet` for static analysis
- Add tests for new features
- Document exported functions

## API Integration Details

### Authentication

The UniFi Manager API uses X-API-KEY header authentication:

```go
client.Request(ctx, "GET", "/v1/sites", nil, 
    client.WithHeader("X-API-KEY", apiKey))
```

### Rate Limiting

The API enforces rate limits:
- v1 stable: 10,000 requests per minute
- Response includes `Retry-After` header if rate limited

Handle with:
```go
if resp.StatusCode == 429 {
    retryAfter := resp.Header.Get("Retry-After")
    // Wait and retry
}
```

### Response Format

All API responses follow a consistent JSON format with data and metadata.

## Deploying

### Docker

```bash
# Build image
docker build -t unifi-manager-mcp .

# Run container
docker run -e UNIFI_API_KEY=your-key \
  -p 8000:8000 \
  unifi-manager-mcp
```

### Environment Variables

Set these before running:
- `UNIFI_API_KEY` - Required API key
- `MCP_TRANSPORT` - Transport type (stdio/http)
- `MCP_HTTP_ADDR` - HTTP server address if using HTTP transport
- `LOG_LEVEL` - Logging verbosity

## Troubleshooting

### API Key Issues
- Ensure API key is valid and has access
- Verify X-API-KEY header is being sent correctly
- Check API key hasn't expired

### Rate Limiting
- Monitor request frequency
- Implement exponential backoff
- Check `Retry-After` headers

### Connection Issues
- Verify network connectivity to api.ui.com
- Check for firewall restrictions
- Review API documentation for endpoint status

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Support

For questions or issues:
1. Check the [UniFi Site Manager API documentation](https://developer.ui.com/site-manager-api/gettingstarted)
2. Review existing GitHub issues
3. Open a new issue with detailed information

---

Happy coding! 🚀
