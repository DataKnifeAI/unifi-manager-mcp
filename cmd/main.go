package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/surrealwolf/unifi-manager-mcp/internal/mcp"
	"github.com/surrealwolf/unifi-manager-mcp/internal/unifi"
)

func init() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Configure logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Get configuration from environment
	apiKey := os.Getenv("UNIFI_API_KEY")
	if apiKey == "" {
		logrus.Fatal("UNIFI_API_KEY environment variable is required")
	}

	managerClient := unifi.NewClient(apiKey)

	// Initialize MCP server
	server := mcp.NewServer(managerClient)

	// Determine transport mode
	transport := strings.ToLower(os.Getenv("MCP_TRANSPORT"))
	if transport == "" {
		transport = "stdio"
	}

	switch transport {
	case "http":
		httpAddr := os.Getenv("MCP_HTTP_ADDR")
		if httpAddr == "" {
			httpAddr = ":8000"
		}
		logrus.Infof("Starting UniFi Manager MCP Server on HTTP at %s", httpAddr)
		go func() {
			if err := server.ServeHTTP(httpAddr, ctx); err != nil {
				logrus.WithError(err).Fatal("HTTP Server error")
			}
		}()
	default:
		logrus.Info("Starting UniFi Manager MCP Server on stdio transport")
		go func() {
			if err := server.ServeStdio(ctx); err != nil {
				logrus.WithError(err).Fatal("Server error")
			}
		}()
	}

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down gracefully...")
	cancel()
	logrus.Info("UniFi Manager MCP Server stopped")
}
