.PHONY: build test run docker-build docker-run clean help

help:
	@echo "UniFi Manager MCP - Available targets:"
	@echo ""
	@echo "  build          - Build the project"
	@echo "  test           - Run tests"
	@echo "  run            - Run the server locally"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  clean          - Clean build artifacts"
	@echo ""

build:
	go build -o bin/unifi-manager-mcp ./cmd

test:
	go test -v ./...

run: build
	UNIFI_API_KEY=test ./bin/unifi-manager-mcp

docker-build:
	docker build -t unifi-manager-mcp:latest .

docker-run: docker-build
	docker run -it \
		-e UNIFI_API_KEY=your-key \
		-e LOG_LEVEL=info \
		-p 8000:8000 \
		unifi-manager-mcp:latest

clean:
	rm -rf bin/ dist/ build/

.DEFAULT_GOAL := help
