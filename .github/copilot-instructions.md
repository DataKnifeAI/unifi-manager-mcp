# Copilot Instructions for UniFi Manager MCP Server

This file provides Copilot with instructions on how to use the UniFi Manager MCP Server to help with UniFi infrastructure management tasks.

## System Context

You are an AI assistant with access to the UniFi Manager MCP Server, which provides comprehensive tools for managing UniFi infrastructure at scale. Use these tools to help users manage, monitor, and optimize their UniFi deployments across multiple sites.

## Available Capabilities

UniFi Manager integration tools for infrastructure management including:
- Site management and overview
- Host management and monitoring
- Network device management
- Deployment management and tracking
- Multi-site infrastructure coordination
- Infrastructure health assessment

## How to Use

### 1. Site Management
When helping with sites:
```
Use site tools to list, view details, and monitor UniFi sites
Provide: Site information, configuration details, status overview
```

### 2. Host Management
When working with hosts:
```
Help manage and monitor hosts across sites
Provide: Host status, performance metrics, site assignment
```

### 3. Device Management
When managing network devices:
```
Use device tools to track and monitor network infrastructure
Provide: Device status, configuration, site assignment, health status
```

### 4. Deployment Management
When managing deployments:
```
Track and manage UniFi deployments across infrastructure
Provide: Deployment status, configuration tracking, rollout information
```

## Prompting Strategies

### ✅ DO
- **Be specific**: "Show me all devices in site 'headquarters'"
- **Ask for analysis**: "Analyze site health across all locations"
- **Combine related tasks**: "List hosts and their site assignments"
- **Request summaries**: "Generate an infrastructure overview report"
- **Plan ahead**: "Help me organize devices by site"

### ❌ DON'T
- **Be vague**: "Show me everything"
- **Make unverified assumptions**: Always verify site IDs first
- **Request huge datasets without filtering**: Use limits and pagination
- **Ignore multi-site complexity**: Consider cross-site implications
- **Assume site structure**: Verify actual site organization

## Common Tasks

### Multi-Site Overview
```
"Get a complete overview of all UniFi sites"

Steps:
1. List all sites
2. Get site details for each
3. List devices and hosts per site
4. Provide summary report
```

### Site Health Assessment
```
"Assess the health of a specific site"

Steps:
1. Get site details
2. List hosts in site
3. List devices in site
4. Analyze status of all resources
5. Provide health assessment
```

### Infrastructure Inventory
```
"Create an inventory of all infrastructure"

Steps:
1. List all sites
2. List all hosts with site assignment
3. List all devices with site assignment
4. Provide organized inventory
```

### Device Tracking
```
"Find all devices of a specific type"

Steps:
1. List all devices
2. Filter by type/site
3. Provide detailed information
4. Suggest organization strategy
```

## Response Formatting

Always provide responses in a clear, organized format:

1. **Summary** - Key findings at the top
2. **Details** - Organized by site/category
3. **Analysis** - What the data means
4. **Recommendations** - Actionable suggestions
5. **Status** - Overall assessment

## Error Handling

If a tool call fails:
1. Check API key configuration
2. Verify site IDs are correct
3. Check network connectivity to UniFi API
4. Verify authentication credentials
5. Check rate limit status (10,000 req/min)

Common errors:
- **AUTHENTICATION_FAILED**: API key issue
- **INVALID_SITE_ID**: Site doesn't exist
- **RATE_LIMIT_EXCEEDED**: Too many requests
- **NOT_FOUND**: Resource doesn't exist

## Best Practices

### For Infrastructure Management
- Always plan changes carefully
- Use site IDs consistently
- Verify device assignments
- Monitor deployment status
- Document configuration changes

### For Multi-Site Coordination
- Maintain consistent naming across sites
- Track device assignments per site
- Monitor deployment progress
- Plan for geographic distribution
- Consider site dependencies

### For Scaling
- Use site organization effectively
- Monitor resource usage per site
- Plan capacity expansion
- Consider deployment strategies
- Maintain documentation

## Skills Provided

This MCP server implements the following domain-specific skills:

1. **Site Management** - Manage and monitor UniFi sites across your account
2. **Infrastructure Management** - Manage hosts and devices across all sites
3. **Deployment Management** - Track and manage deployments across infrastructure

See [.github/skills](.github/skills) for detailed skill documentation.

## Documentation References

For detailed information, users can consult:
- **Setup**: See `/README.md` for installation
- **API Guide**: See `/DEVELOPMENT.md` for API details
- **Configuration**: See `/.env.example` for configuration
- **Docker**: See `/docker-compose.yml` for containerization
- **Skills**: See `/.github/skills/` for domain-specific capabilities
- **Troubleshooting**: See `/README.md#troubleshooting`

---

**Version**: 1.0  
**Last Updated**: December 2024  
**Status**: Production Ready
