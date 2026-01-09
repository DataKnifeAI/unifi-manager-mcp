# Harbor Registry Setup

This project uses Harbor (harbor.dataknife.net) as the container registry for Docker images.

## GitHub Secrets Configuration

To enable automated builds and pushes to Harbor, you need to configure the following GitHub secrets:

### Required Secrets

1. **HARBOR_USERNAME**
   - Description: Harbor robot account username for CI/CD
   - Format: `robot$library+ci-builder` (replace with your actual robot account username)

2. **HARBOR_PASSWORD**
   - Description: Harbor robot account password/token
   - Format: Your actual Harbor robot account password/token

### How to Add GitHub Secrets

**Using GitHub CLI (recommended):**
```bash
gh secret set HARBOR_USERNAME --body "robot\$library+ci-builder"
gh secret set HARBOR_PASSWORD --body "your-password-here"
```

**Using GitHub Web UI:**
1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret with the appropriate name and value

## Local Development

### Manual Docker Login

```bash
docker login harbor.dataknife.net \
  -u '<your-harbor-username>' \
  -p '<your-harbor-password>'
```

### Build and Push with Make

**Important**: When using `make` with variables on the command line, escape `$` characters with `$$` because Make interprets `$` as a variable reference.

**Option 1: Using environment variables**
```bash
# Set environment variables (replace with your actual credentials)
export HARBOR_USERNAME='robot$library+ci-builder'
export HARBOR_PASSWORD='<your-harbor-password>'

# Build and push
make docker-push
```

**Option 2: Passing variables directly (escape $ with $$)**
```bash
make docker-push HARBOR_USERNAME='robot$$library+ci-builder' HARBOR_PASSWORD='<your-harbor-password>'
```

**Option 3: Using docker login directly**
```bash
# Login first (replace with your actual credentials)
docker login harbor.dataknife.net -u '<your-harbor-username>' -p '<your-harbor-password>'

# Then build and push (will use cached credentials)
make docker-build
docker push harbor.dataknife.net/library/unifi-manager-mcp:latest
```

### Pull from Harbor

```bash
# Pull the latest image
docker pull harbor.dataknife.net/library/unifi-manager-mcp:latest

# Or use docker-compose (will pull automatically)
docker-compose pull
docker-compose up -d
```

## Image Naming Convention

Images are pushed to: `harbor.dataknife.net/library/unifi-manager-mcp:<tag>`

- `latest` - Latest build from main/master branch
- `<branch-name>` - Branch-specific builds
- `v<version>` - Semantic version tags (e.g., v1.0.0)
- `<branch>-<sha>` - Commit-specific builds

## Base Images

The Dockerfile uses Harbor-cached base images from DockerHub:
- `harbor.dataknife.net/dockerhub/library/golang:1.23.2-alpine`
- `harbor.dataknife.net/dockerhub/library/alpine:latest`

These images are automatically cached in Harbor when first pulled from DockerHub.
