package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/surrealwolf/unifi-manager-mcp/internal/unifi"
)

// Server represents the MCP server
type Server struct {
	unifiClient *unifi.Client
	server      *server.MCPServer
	logger      *logrus.Entry
}

// NewServer creates a new MCP server
func NewServer(unifiClient *unifi.Client) *Server {
	s := &Server{
		unifiClient: unifiClient,
		server:      server.NewMCPServer("unifi-manager-mcp", "0.1.0"),
		logger:      logrus.WithField("component", "MCPServer"),
	}

	s.registerTools()
	return s
}

func (s *Server) registerTools() {
	// Helper to create tool definitions
	addTool := func(name, desc string, handler server.ToolHandlerFunc, properties map[string]any) {
		s.server.AddTool(mcp.Tool{
			Name:        name,
			Description: desc,
			InputSchema: mcp.ToolInputSchema{
				Type:       "object",
				Properties: properties,
			},
		}, handler)
	}

	// Sites tools
	addTool("list_sites", "List all UniFi sites in the account", s.listSites, map[string]any{})
	addTool("get_site_details", "Get detailed information about a specific site", s.getSiteDetails, map[string]any{
		"site_id": map[string]any{"type": "string", "description": "Site ID"},
	})
	addTool("get_site_overview", "Get overview statistics for a specific site", s.getSiteOverview, map[string]any{
		"site_id": map[string]any{"type": "string", "description": "Site ID"},
	})

	// Hosts tools
	addTool("list_hosts", "List all hosts across all sites", s.listHosts, map[string]any{})
	addTool("get_host_details", "Get detailed information about a specific host", s.getHostDetails, map[string]any{
		"host_id": map[string]any{"type": "string", "description": "Host ID"},
	})
	addTool("get_hosts_by_site", "Get all hosts in a specific site", s.getHostsBySite, map[string]any{
		"site_id": map[string]any{"type": "string", "description": "Site ID"},
	})

	// Devices tools
	addTool("list_devices", "List all network devices across all sites", s.listDevices, map[string]any{})
	addTool("get_device_details", "Get detailed information about a specific device", s.getDeviceDetails, map[string]any{
		"device_id": map[string]any{"type": "string", "description": "Device ID"},
	})
	addTool("get_devices_by_site", "Get all devices in a specific site", s.getDevicesBySite, map[string]any{
		"site_id": map[string]any{"type": "string", "description": "Site ID"},
	})

	// Deployments tools
	addTool("list_deployments", "List all deployments across all sites", s.listDeployments, map[string]any{})
	addTool("get_deployment_details", "Get detailed information about a specific deployment", s.getDeploymentDetails, map[string]any{
		"deployment_id": map[string]any{"type": "string", "description": "Deployment ID"},
	})

	s.logger.Info("Registered 11 tools")
}

// ServeStdio starts the MCP server with stdio transport
func (s *Server) ServeStdio(ctx context.Context) error {
	s.logger.Info("Starting UniFi Manager MCP Server on stdio")
	return server.ServeStdio(s.server)
}

// ServeHTTP starts the MCP server with HTTP transport
func (s *Server) ServeHTTP(addr string, ctx context.Context) error {
	s.logger.Infof("Starting UniFi Manager MCP Server on HTTP at %s", addr)

	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse the MCP request
		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Log the request
		s.logger.Debugf("HTTP MCP request received: %v", requestData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{
			"status": "MCP HTTP transport is available",
			"info":   "This is an HTTP endpoint. Use stdio transport for full MCP protocol support.",
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	return http.ListenAndServe(addr, nil)
}

// Tool handlers

func (s *Server) listSites(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: list_sites")

	sites, err := s.unifiClient.ListSites(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list sites: %v", err)), nil
	}

	// Convert to interface for JSON marshaling
	var sitesInterface []interface{}
	for _, s := range sites {
		sitesInterface = append(sitesInterface, s)
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"sites": sitesInterface,
		"count": len(sitesInterface),
	})
}

func (s *Server) getSiteDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_site_details")

	siteID := request.GetString("site_id", "")
	if siteID == "" {
		return mcp.NewToolResultError("site_id parameter is required"), nil
	}

	site, err := s.unifiClient.GetSiteDetails(ctx, siteID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get site details: %v", err)), nil
	}

	return mcp.NewToolResultJSON(site)
}

func (s *Server) getSiteOverview(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_site_overview")

	siteID := request.GetString("site_id", "")
	if siteID == "" {
		return mcp.NewToolResultError("site_id parameter is required"), nil
	}

	site, err := s.unifiClient.GetSiteDetails(ctx, siteID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get site overview: %v", err)), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"site_id":     site.ID,
		"name":        site.Name,
		"description": site.Description,
	})
}

func (s *Server) listHosts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: list_hosts")

	hosts, err := s.unifiClient.ListHosts(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list hosts: %v", err)), nil
	}

	// Convert to interface for JSON marshaling
	var hostsInterface []interface{}
	for _, h := range hosts {
		hostsInterface = append(hostsInterface, h)
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"hosts": hostsInterface,
		"count": len(hostsInterface),
	})
}

func (s *Server) getHostDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_host_details")

	hostID := request.GetString("host_id", "")
	if hostID == "" {
		return mcp.NewToolResultError("host_id parameter is required"), nil
	}

	host, err := s.unifiClient.GetHostDetails(ctx, hostID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get host details: %v", err)), nil
	}

	return mcp.NewToolResultJSON(host)
}

func (s *Server) getHostsBySite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_hosts_by_site")

	siteID := request.GetString("site_id", "")
	if siteID == "" {
		return mcp.NewToolResultError("site_id parameter is required"), nil
	}

	hosts, err := s.unifiClient.ListHosts(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hosts for site: %v", err)), nil
	}

	// Filter hosts by site_id (client-side filtering)
	var filtered []interface{}
	for _, host := range hosts {
		if host.SiteID == siteID {
			filtered = append(filtered, host)
		}
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"hosts":   filtered,
		"count":   len(filtered),
		"site_id": siteID,
	})
}

func (s *Server) listDevices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: list_devices")

	devices, err := s.unifiClient.ListDevices(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list devices: %v", err)), nil
	}

	// Convert to interface for JSON marshaling
	var devicesInterface []interface{}
	for _, d := range devices {
		devicesInterface = append(devicesInterface, d)
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"devices": devicesInterface,
		"count":   len(devicesInterface),
	})
}

func (s *Server) getDeviceDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_device_details")

	deviceID := request.GetString("device_id", "")
	if deviceID == "" {
		return mcp.NewToolResultError("device_id parameter is required"), nil
	}

	device, err := s.unifiClient.GetDeviceDetails(ctx, deviceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get device details: %v", err)), nil
	}

	return mcp.NewToolResultJSON(device)
}

func (s *Server) getDevicesBySite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_devices_by_site")

	siteID := request.GetString("site_id", "")
	if siteID == "" {
		return mcp.NewToolResultError("site_id parameter is required"), nil
	}

	devices, err := s.unifiClient.ListDevices(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get devices for site: %v", err)), nil
	}

	// Filter devices by site_id (client-side filtering)
	var filtered []interface{}
	for _, device := range devices {
		if device.SiteID == siteID {
			filtered = append(filtered, device)
		}
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"devices": filtered,
		"count":   len(filtered),
		"site_id": siteID,
	})
}

func (s *Server) listDeployments(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: list_deployments")

	deployments, err := s.unifiClient.ListDeployments(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list deployments: %v", err)), nil
	}

	// Convert to interface for JSON marshaling
	var deploymentsInterface []interface{}
	for _, d := range deployments {
		deploymentsInterface = append(deploymentsInterface, d)
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"deployments": deploymentsInterface,
		"count":       len(deploymentsInterface),
	})
}

func (s *Server) getDeploymentDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_deployment_details")

	deploymentID := request.GetString("deployment_id", "")
	if deploymentID == "" {
		return mcp.NewToolResultError("deployment_id parameter is required"), nil
	}

	deployment, err := s.unifiClient.GetDeploymentDetails(ctx, deploymentID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get deployment details: %v", err)), nil
	}

	return mcp.NewToolResultJSON(deployment)
}
