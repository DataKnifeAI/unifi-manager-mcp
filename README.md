# UniFi Manager MCP

Model Context Protocol (MCP) server for Ubiquiti UniFi Site Manager API integration. Monitor and manage your UniFi infrastructure at scale through an AI-powered interface.

**Focused on:** Site Manager API integration for monitoring hosts, devices, deployments, and network infrastructure at scale.

⚠️ **Early Development Warning**: This project is in early development and may contain bugs that could cause unexpected behavior. Use with caution.

⚠️ **Prompt Injection Risk**: You are responsible for guarding against prompt injection when using these tools. Exercise extreme caution or use MCP tools only on systems and data you trust.

## Features

- **Infrastructure Monitoring**: Query and monitor hosts, devices, and deployments
- **Site Management**: Access detailed information about UniFi sites
- **Deployment Overview**: Track deployment status and configuration
- **Device Discovery**: List and query all managed devices
- **API Key Authentication**: Secure X-API-KEY header-based authentication
- **Stdio Transport**: MCP protocol over standard input/output for seamless integration
- **HTTP Transport**: Optional HTTP API for remote connections and integration
- **Rate Limit Awareness**: Built-in support for API rate limiting (10,000 req/min)

## Quick Start

### Installation

```bash
# Clone and build
git clone https://github.com/surrealwolf/unifi-manager-mcp.git
cd unifi-manager-mcp
go build -o bin/unifi-manager-mcp ./cmd
```

### Configuration

Create a `.env` file:

```bash
UNIFI_API_KEY=your-api-key-here
LOG_LEVEL=info
MCP_TRANSPORT=stdio
```

### Obtaining API Key

1. Sign in to the UniFi Site Manager at https://unifi.ui.com
2. Navigate to the API section from the left navigation bar
3. Select "Create API Key"
4. Copy the generated key and store it securely (displayed only once)
5. Use the key in the `UNIFI_API_KEY` environment variable

### Running the Server

**Stdio Transport (Default):**
```bash
UNIFI_API_KEY=your-key ./bin/unifi-manager-mcp
```

**HTTP Transport:**
```bash
MCP_TRANSPORT=http MCP_HTTP_ADDR=:8000 UNIFI_API_KEY=your-key ./bin/unifi-manager-mcp
```

Then access the endpoints:
```bash
# Health check
curl http://localhost:8000/health

# MCP endpoint
curl -X POST http://localhost:8000/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/list"}'
```

**Environment Variables:**
- `UNIFI_API_KEY` - API key from Site Manager (Required)
- `MCP_TRANSPORT` - Transport mode: `stdio` (default) or `http`
- `MCP_HTTP_ADDR` - HTTP server address (default: `:8000`)
- `LOG_LEVEL` - Logging level: debug, info, warn, error (default: info)

## Available Tools

### Sites Management (3 tools)
- `list_sites` - List all UniFi sites
- `get_site_details` - Get detailed site information
- `get_site_overview` - Get site overview and statistics

### Hosts Management (3 tools)
- `list_hosts` - List all hosts in the infrastructure
- `get_host_details` - Get detailed host information
- `get_host_by_site` - Get hosts filtered by site

### Devices Management (3 tools)
- `list_devices` - List all managed devices
- `get_device_details` - Get detailed device information
- `get_devices_by_site` - Get devices filtered by site

### Deployments Management (2 tools)
- `list_deployments` - List all deployments
- `get_deployment_details` - Get detailed deployment information

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `UNIFI_API_KEY` | API key from Site Manager | Required |
| `MCP_TRANSPORT` | Transport protocol (stdio/http) | stdio |
| `MCP_HTTP_ADDR` | HTTP server address | :8000 |
| `LOG_LEVEL` | Logging level | info |

## API Reference

For detailed UniFi Site Manager API documentation, see:
- **Official Documentation**: https://developer.ui.com/site-manager-api/gettingstarted
- **API Base URL**: https://api.ui.com/v1
- **Rate Limit**: 10,000 requests per minute (v1 stable)

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run with Docker

```bash
make docker-build
make docker-run
```

## License

MIT License - See LICENSE file for details

## Support

For issues and questions:
- Check the [UniFi Site Manager API Documentation](https://developer.ui.com/site-manager-api/gettingstarted)
- Review tool examples in `docs/examples/`

---

**Built with Claude Haiku 4.5** - Crafted by AI to extend your infrastructure possibilities. 🤖✨
