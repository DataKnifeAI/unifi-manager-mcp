package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://api.ui.com/v1"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new UniFi Manager API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 0,
		},
	}
}

// Site represents a UniFi site
type Site struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Host represents a UniFi host
type Host struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	SiteID string `json:"site_id"`
	Status string `json:"status"`
}

// Device represents a UniFi device
type Device struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	SiteID string `json:"site_id"`
}

// Deployment represents a UniFi deployment
type Deployment struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ListSites returns all sites
func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	var response struct {
		Data []Site `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", "/sites", nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	return response.Data, nil
}

// GetSiteDetails returns detailed information about a specific site
func (c *Client) GetSiteDetails(ctx context.Context, siteID string) (*Site, error) {
	var response struct {
		Data Site `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/sites/%s", siteID), nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get site details: %w", err)
	}

	return &response.Data, nil
}

// ListHosts returns all hosts
func (c *Client) ListHosts(ctx context.Context) ([]Host, error) {
	var response struct {
		Data []Host `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", "/hosts", nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list hosts: %w", err)
	}

	return response.Data, nil
}

// GetHostDetails returns detailed information about a specific host
func (c *Client) GetHostDetails(ctx context.Context, hostID string) (*Host, error) {
	var response struct {
		Data Host `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/hosts/%s", hostID), nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get host details: %w", err)
	}

	return &response.Data, nil
}

// ListDevices returns all devices
func (c *Client) ListDevices(ctx context.Context) ([]Device, error) {
	var response struct {
		Data []Device `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", "/devices", nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	return response.Data, nil
}

// GetDeviceDetails returns detailed information about a specific device
func (c *Client) GetDeviceDetails(ctx context.Context, deviceID string) (*Device, error) {
	var response struct {
		Data Device `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/devices/%s", deviceID), nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get device details: %w", err)
	}

	return &response.Data, nil
}

// ListDeployments returns all deployments
func (c *Client) ListDeployments(ctx context.Context) ([]Deployment, error) {
	var response struct {
		Data []Deployment `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", "/deployments", nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	return response.Data, nil
}

// GetDeploymentDetails returns detailed information about a specific deployment
func (c *Client) GetDeploymentDetails(ctx context.Context, deploymentID string) (*Deployment, error) {
	var response struct {
		Data Deployment `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/deployments/%s", deploymentID), nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get deployment details: %w", err)
	}

	return &response.Data, nil
}

// doRequest performs an HTTP request to the API
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}
