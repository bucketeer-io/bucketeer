# Dev Container Dependency Caching

This document explains the caching optimizations implemented in the Bucketeer dev container setup to significantly reduce startup times when reopening the dev container.

## 🚀 **Quick Start - Where to Run Commands**

**📍 LOCAL MACHINE (Terminal outside VS Code):**
- `bash .devcontainer/cache-manager.sh pull-images` - Pre-cache 4.5GB images

**📍 DEVCONTAINER (Terminal inside VS Code):**  
- `bash .devcontainer/cache-manager.sh load-to-minikube` - Transfer cache to minikube
- `minikube start` and `make start-minikube` - Development workflow

**Why this matters:** Docker volumes are shared between host and devcontainer, but minikube has its own Docker daemon!

## Caching Strategies

### 1. Persistent Volume Mounts
The dev container uses Docker named volumes to persist dependencies across container rebuilds:

- **Go modules cache** (`/go/pkg/mod`): Prevents re-downloading Go dependencies
- **Go tools** (`/home/codespace/go-tools`): Keeps installed Go development tools
- **Node.js modules** (various `node_modules` directories): Preserves installed npm/yarn packages
  - `/workspaces/bucketeer/ui/dashboard/node_modules`
  - `/workspaces/bucketeer/ui/web-v2/node_modules`
  - `/workspaces/bucketeer/evaluation/typescript/node_modules`
- **Package manager caches**: Yarn and npm local caches at `/home/codespace/.yarn/cache` and `/home/codespace/.npm`
- **Docker image cache** (`/home/codespace/.docker-images`): Stores emulator and development images as tar files
  - BigQuery Emulator: `ghcr.io/goccy/bigquery-emulator:latest` (~500MB)
  - Pub/Sub Emulator: `gcr.io/google.com/cloudsdktool/google-cloud-cli:latest` (~1GB)
  - Database images: MySQL, Redis, Vault
  - Development tools: Atlas migrations, distroless base images

### 2. Intelligent Setup Script
The `setup.sh` script implements smart dependency checking:

- **Go tools**: Checks if tools are already installed before running `go install`
- **Go vendor**: Compares `go.mod`/`go.sum` timestamps with `vendor/` directory
- **Node modules**: Compares `yarn.lock` timestamps with `node_modules/` directories
- **Docker cleanup**: Only runs when disk space is low (< 5GB)
- **Cache status reporting**: Shows what will be installed vs what's cached
- **Emulator image pre-caching**: Downloads and caches large emulator images when internet is available
- **Offline image restoration**: Loads cached images from storage when no internet connection

### 3. Optimized Installation Flags
- `yarn install --frozen-lockfile --silent`: Faster, deterministic installs
- Go proxy optimization with `GOPROXY` environment variable
- Yarn cache folder optimization via `YARN_CACHE_FOLDER` environment variable

## Performance Impact

**Before optimization:**
- Fresh setup: ~5-10 minutes
- Reopen setup: ~3-5 minutes

**After optimization:**
- Fresh setup: ~5-10 minutes (same)
- Reopen setup: ~30-60 seconds (85% faster!)

## Cache Invalidation

Caches are automatically invalidated when:
- `go.mod` or `go.sum` files change (Go dependencies)
- `yarn.lock` files change (Node.js dependencies)
- Go tools are missing from PATH
- Volume mounts are deleted manually

## Volume Configuration

The volumes are configured in `devcontainer.json` with absolute paths:

```json
"mounts": [
  "source=bucketeer-go-mod-cache,target=/go/pkg/mod,type=volume",
  "source=bucketeer-go-tools,target=/home/codespace/go-tools,type=volume",
  "source=bucketeer-dashboard-node-modules,target=/workspaces/bucketeer/ui/dashboard/node_modules,type=volume",
  "source=bucketeer-web-v2-node-modules,target=/workspaces/bucketeer/ui/web-v2/node_modules,type=volume",
  "source=bucketeer-eval-ts-node-modules,target=/workspaces/bucketeer/evaluation/typescript/node_modules,type=volume",
  "source=bucketeer-yarn-cache,target=/home/codespace/.yarn/cache,type=volume",
  "source=bucketeer-npm-cache,target=/home/codespace/.npm,type=volume",
  "source=bucketeer-docker-images,target=/home/codespace/.docker-images,type=volume"
]
```

**Important:** Use absolute paths instead of variables like `${containerHome}` as they are not resolved in mount configurations.

## Correct Minikube Workflow 

**❌ What you tried (doesn't work):**
```bash
minikube delete
make start-minikube  # This re-downloads images from internet
```

**✅ Correct workflow (uses cache):**

### **Step 1: Pre-cache images (One-time setup)**
**Location: Run from LOCAL MACHINE (terminal outside devcontainer)**
```bash
# On your local machine terminal (outside devcontainer)
cd /path/to/your/bucketeer-repo
bash .devcontainer/cache-manager.sh pull-images  # Downloads 4.5GB once
```
*Why local machine? More reliable internet connection and Docker access*

### **Step 2-4: Daily development**
**Location: Run from INSIDE DEVCONTAINER**
```bash
# Inside VS Code devcontainer terminal

# 2. Start minikube (downloads small K8s images ~70MB)
minikube delete
minikube start --memory max --cpus max
minikube addons enable ingress

# 3. CRITICAL: Load cached images into minikube (FULLY AUTOMATED!)
bash .devcontainer/cache-manager.sh load-to-minikube
# ✨ This now automatically handles Docker context issues and restores from cache

# 4. Deploy services (now uses cached images - NO DOWNLOADS!)
make start-minikube
```
*Why devcontainer? Minikube runs inside devcontainer environment*

**🔄 For repeated development cycles:**
**Location: INSIDE DEVCONTAINER**
```bash
# Inside VS Code devcontainer terminal - when you need to restart minikube
minikube delete
minikube start --memory max --cpus max           # Downloads ~70MB K8s images
minikube addons enable ingress                   # Small K8s addon (~20MB)
bash .devcontainer/cache-manager.sh load-to-minikube  # Fully automated: restores + transfers 4.5GB
make start-minikube                              # Uses cached images - fast!
```

## 🔑 **Key Insight: Where to Run Commands**

| Command | Location | Why |
|---------|----------|-----|
| `pull-images` | **Local Machine** | Better internet, more reliable Docker access |
| `load-to-minikube` | **Devcontainer** | Minikube runs inside devcontainer |
| `minikube start` | **Devcontainer** | Minikube environment |
| `make start-minikube` | **Devcontainer** | Project development environment |

**Important:** Docker volumes are **shared** between your local machine and devcontainer (same Docker Desktop daemon), so images cached locally are accessible from devcontainer.

The key insight: `make start-minikube` only uses images that are **already in minikube's Docker daemon**. You must run `load-to-minikube` to transfer cached images from local Docker to minikube.

## Manual Cache Management

### Clear all caches
```bash
# Remove Docker volumes (will trigger full reinstall next time)
docker volume rm bucketeer-go-mod-cache
docker volume rm bucketeer-go-tools
docker volume rm bucketeer-dashboard-node-modules
docker volume rm bucketeer-web-v2-node-modules
docker volume rm bucketeer-eval-ts-node-modules
docker volume rm bucketeer-yarn-cache
docker volume rm bucketeer-npm-cache
docker volume rm bucketeer-docker-images
```

### Manage Docker images
```bash
# Check cached emulator images status (either location)
bash .devcontainer/cache-manager.sh images

# Pre-pull and cache emulator images (LOCAL MACHINE - good internet)
bash .devcontainer/cache-manager.sh pull-images

# Load cached images into Docker daemon (either location)
bash .devcontainer/cache-manager.sh load-images

# Load cached images into minikube (DEVCONTAINER - THE KEY STEP!)
bash .devcontainer/cache-manager.sh load-to-minikube

# Clear cached Docker images (either location)
bash .devcontainer/cache-manager.sh clear-images
```

### Force dependency refresh
```bash
# From inside the dev container
rm -rf ui/dashboard/node_modules ui/web-v2/node_modules evaluation/typescript/node_modules
rm -rf vendor/
bash .devcontainer/setup.sh
```

### Check cache status
The setup script provides a detailed status showing what dependencies are cached vs what needs to be installed:

```bash
bash .devcontainer/setup.sh
```

### Debug mode
If you're having issues with tool detection, enable debug mode:

```bash
DEBUG_SETUP=true bash .devcontainer/setup.sh
```

This will show detailed information about where tools are being looked for and what's found.

## Troubleshooting

**Dependencies seem outdated:**
- Update your `yarn.lock` or `go.sum` files
- The cache will automatically refresh

**Setup script hangs:**
- Check Docker volume space: `docker system df`
- Consider clearing old volumes if space is full

**Want to force fresh install:**
- Delete the relevant volume and restart the container
- Or temporarily rename lock files to trigger reinstall

**Mount path errors:**
- Ensure absolute paths are used in devcontainer.json
- Avoid variables like `${containerHome}` in mount configurations

**Permission errors:**
- The setup script automatically fixes volume mount permissions for all cache directories
- If you see "permission denied" errors, restart the dev container
- Cache volumes are mounted with root ownership by default, but the script corrects this
- Go modules cache (`/go/pkg/mod`) and Go tools directory are automatically fixed

**Emulator image issues:**
- Run `.devcontainer/cache-manager.sh images` (either location) to check image cache status  
- If minikube fails to start emulators, run `.devcontainer/cache-manager.sh load-to-minikube` (devcontainer)
- Pre-cache images when you have good internet: `.devcontainer/cache-manager.sh pull-images` (local machine)
- Large images (1GB+) are automatically cached to avoid re-downloading on 4G connections

**Location-specific issues:**
- If `pull-images` fails from devcontainer, try from local machine terminal
- If `load-to-minikube` fails, ensure you're inside devcontainer and minikube is running
- Docker volumes are shared between local machine and devcontainer - cache once, use everywhere

**✨ NEW: Automatic Docker Context Handling**
- `load-to-minikube` now automatically detects missing images and restores them from cache
- No more manual commands needed when switching between local machine and devcontainer
- Handles Docker daemon synchronization issues transparently

## Emulator Image Caching Benefits

### 🚀 **Solves the Offline Problem:**
- **Before**: `minikube delete` removes all images → re-download 1GB+ on 4G = impossible
- **After**: Images cached persistently → `minikube start` works offline

### 📱 **Perfect for Mobile/Limited Internet:**
- BigQuery Emulator (~500MB) + Pub/Sub Emulator (~1GB) cached once
- No more waiting 10+ minutes on 4G connections
- Cached images survive `minikube delete` and container restarts

### 🔄 **Intelligent Workflow:**
1. **With good WiFi**: Pre-cache images: `.devcontainer/cache-manager.sh pull-images`
2. **After `minikube delete`**: Load cached images: `.devcontainer/cache-manager.sh load-images`
3. **Then start minikube**: `make start-minikube` (uses cached images, no downloads!)

### 🛠 **Manual Control:**
```bash
# When you have good internet - cache everything (LOCAL MACHINE)
.devcontainer/cache-manager.sh pull-images

# Check what's cached (either location)
.devcontainer/cache-manager.sh images

# Load cached images into minikube (DEVCONTAINER - FULLY AUTOMATED!)
.devcontainer/cache-manager.sh load-to-minikube

# Force load cached images into local Docker (either location)
.devcontainer/cache-manager.sh load-images
```

## Additional Optimizations

For even better performance, consider:
1. Using a faster SSD for Docker volumes
2. Increasing Docker Desktop resource allocation
3. Using Docker BuildKit for faster image builds
4. Pre-warming the dev container image with common dependencies 

## Script Output Example

When everything is cached:
```
🚀 Starting post-attach setup with intelligent caching...
[INFO] Ensuring cache directories have correct permissions...
[SUCCESS] Cache permissions fixed
[INFO] Checking cache status...
[SUCCESS] All Go tools are already installed
[SUCCESS] Go vendor directory is up to date
[SUCCESS] Dashboard: Node.js dependencies are up to date
[SUCCESS] Web-v2: Node.js dependencies are up to date
[SUCCESS] Evaluation TypeScript: Node.js dependencies are up to date
[INFO] Restoring cached emulator images to Docker daemon...
[SUCCESS] All cached images already available in Docker daemon
[SUCCESS] All emulator images already cached
[SUCCESS] 🎉 All dependencies are cached and up to date! Setup completed in seconds.
[INFO] Cache volumes will persist dependencies for next container restart
[INFO] Emulator images cached for offline minikube usage
```

When some dependencies need installation:
```
🚀 Starting post-attach setup with intelligent caching...
[INFO] Ensuring cache directories have correct permissions...
[SUCCESS] Cache permissions fixed
[INFO] Checking cache status...
[WARNING] Missing Go tools: goimports golangci-lint
[SUCCESS] Go vendor directory is up to date
[WARNING] Dashboard: yarn.lock newer than node_modules/
[INFO] Cache analysis complete. Will install: Go tools Dashboard dependencies
[INFO] Restoring cached emulator images to Docker daemon...
[SUCCESS] Restored 3 images from cache
[INFO] Installing Go development tools...
[SUCCESS] Go tools installed successfully
[INFO] Installing Dashboard dependencies...
[SUCCESS] Dashboard dependencies installed
[INFO] Pre-pulling and caching emulator images...
[INFO] Pulling gcr.io/google.com/cloudsdktool/google-cloud-cli:latest...
[SUCCESS] Successfully cached 2 new emulator images
[SUCCESS] 🎉 Post-attach setup completed!
[INFO] Cache volumes will persist dependencies for next container restart
[INFO] Emulator images cached for offline minikube usage
``` 