.PHONY: build test run docker-build docker-run docker-push docker-login clean help

# Harbor registry configuration
HARBOR_REGISTRY ?= harbor.dataknife.net
HARBOR_PROJECT ?= library
IMAGE_NAME ?= unifi-manager-mcp
IMAGE_TAG ?= latest
FULL_IMAGE_NAME = $(HARBOR_REGISTRY)/$(HARBOR_PROJECT)/$(IMAGE_NAME):$(IMAGE_TAG)

help:
	@echo "UniFi Manager MCP - Available targets:"
	@echo ""
	@echo "  build          - Build the project"
	@echo "  test           - Run tests"
	@echo "  run            - Run the server locally"
	@echo "  docker-build   - Build Docker image (local)"
	@echo "  docker-run     - Run Docker container (local)"
	@echo "  docker-login   - Login to Harbor registry"
	@echo "  docker-push    - Build and push Docker image to Harbor"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "Environment variables:"
	@echo "  HARBOR_REGISTRY - Harbor registry URL (default: harbor.dataknife.net)"
	@echo "  HARBOR_PROJECT  - Harbor project name (default: library)"
	@echo "  IMAGE_NAME      - Image name (default: unifi-manager-mcp)"
	@echo "  IMAGE_TAG       - Image tag (default: latest)"
	@echo ""

build:
	go build -o bin/unifi-manager-mcp ./cmd

test:
	go test -v ./...

run: build
	UNIFI_API_KEY=test ./bin/unifi-manager-mcp

docker-build:
	docker build -t $(FULL_IMAGE_NAME) -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run: docker-build
	docker run -it \
		-e UNIFI_API_KEY=your-key \
		-e LOG_LEVEL=info \
		-p 8000:8000 \
		$(IMAGE_NAME):$(IMAGE_TAG)

docker-login:
	@if [ -z "$(HARBOR_USERNAME)" ] || [ -z "$(HARBOR_PASSWORD)" ]; then \
		echo "Error: HARBOR_USERNAME and HARBOR_PASSWORD must be set"; \
		echo "Usage: make docker-login HARBOR_USERNAME='user' HARBOR_PASSWORD='pass'"; \
		exit 1; \
	fi
	@echo "Logging into Harbor registry..."
	@printf '%s\n' '$(HARBOR_PASSWORD)' | docker login $(HARBOR_REGISTRY) \
		-u '$(HARBOR_USERNAME)' \
		--password-stdin

docker-push: docker-build
	@if [ -z "$(HARBOR_USERNAME)" ] || [ -z "$(HARBOR_PASSWORD)" ]; then \
		echo "Error: HARBOR_USERNAME and HARBOR_PASSWORD must be set"; \
		echo "Usage: make docker-push HARBOR_USERNAME='user' HARBOR_PASSWORD='pass'"; \
		exit 1; \
	fi
	@echo "Logging into Harbor registry..."
	@printf '%s\n' '$(HARBOR_PASSWORD)' | docker login $(HARBOR_REGISTRY) \
		-u '$(HARBOR_USERNAME)' \
		--password-stdin
	@echo "Pushing $(FULL_IMAGE_NAME) to Harbor..."
	docker push $(FULL_IMAGE_NAME)
	@echo "Successfully pushed $(FULL_IMAGE_NAME)"

clean:
	rm -rf bin/ dist/ build/

.DEFAULT_GOAL := help
