# Dev Container Dependency Caching

This document explains the caching optimizations implemented in the Bucketeer dev container setup to significantly reduce startup times when reopening the dev container.

## Caching Strategies

### 1. Persistent Volume Mounts
The dev container uses Docker named volumes to persist dependencies across container rebuilds:

- **Go modules cache** (`/go/pkg/mod`): Prevents re-downloading Go dependencies
- **Go tools** (`/home/codespace/go-tools`): Keeps installed Go development tools
- **Node.js modules** (various `node_modules` directories): Preserves installed yarn packages
  - `/workspaces/bucketeer/ui/dashboard/node_modules`
  - `/workspaces/bucketeer/evaluation/typescript/node_modules`
- **Package manager caches**: Yarn local cache at `/home/codespace/.yarn/cache`

### 2. Intelligent Setup Script
The `setup.sh` script implements smart dependency checking:

- **Go tools**: Checks if tools are already installed before running `go install`
- **Go vendor**: Compares `go.mod`/`go.sum` timestamps with `vendor/` directory
- **Node modules**: Compares `yarn.lock` timestamps with `node_modules/` directories
- **Docker cleanup**: Only runs when disk space is low (< 5GB)
- **Cache status reporting**: Shows what will be installed vs what's cached

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
  "source=bucketeer-eval-ts-node-modules,target=/workspaces/bucketeer/evaluation/typescript/node_modules,type=volume",
  "source=bucketeer-yarn-cache,target=/home/codespace/.yarn/cache,type=volume"
]
```

**Important:** Use absolute paths instead of variables like `${containerHome}` as they are not resolved in mount configurations.

## Cache Management

### Volume Organization
The cache manager organizes volumes by type for easier management:

```bash
# Go-related volumes
GO_VOLUMES=(
    "bucketeer-go-mod-cache"
    "bucketeer-go-tools"
)

# JavaScript/Node.js-related volumes
JS_VOLUMES=(
    "bucketeer-dashboard-node-modules"
    "bucketeer-eval-ts-node-modules"
    "bucketeer-yarn-cache"
)

# Combined array (automatically generated)
CACHE_VOLUMES=("${GO_VOLUMES[@]}" "${JS_VOLUMES[@]}")
```

This organization makes it easy to:
- Add/remove volumes by editing only one array
- Clear specific technology stacks (Go vs JS)
- Maintain consistency across all cache operations

### Check cache status
The setup script provides a detailed status showing what dependencies are cached vs what needs to be installed:

```bash
bash .devcontainer/setup.sh
```

### Using the cache manager
The cache manager provides convenient commands for managing cache volumes:

```bash
# Check status of all cache volumes
.devcontainer/cache-manager.sh status

# Show disk usage of cache volumes
.devcontainer/cache-manager.sh size

# Clear all caches (will trigger full reinstall next time)
.devcontainer/cache-manager.sh clear

# Clear only Go-related caches
.devcontainer/cache-manager.sh clear-go

# Clear only Node.js-related caches
.devcontainer/cache-manager.sh clear-js
```

### Manual cache management
```bash
# Remove specific Docker volumes (will trigger full reinstall next time)
docker volume rm bucketeer-go-mod-cache
docker volume rm bucketeer-go-tools
docker volume rm bucketeer-dashboard-node-modules
docker volume rm bucketeer-eval-ts-node-modules
docker volume rm bucketeer-yarn-cache
```

### Force dependency refresh
```bash
# From inside the dev container
rm -rf ui/dashboard/node_modules evaluation/typescript/node_modules
rm -rf vendor/
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

## Additional Optimizations

For even better performance, consider:
1. Using a faster SSD for Docker volumes
2. Increasing Docker Desktop resource allocation
3. Using Docker BuildKit for faster image builds
4. Pre-warming the dev container image with common dependencies

### Adding New Cache Volumes
Thanks to the organized volume structure, adding new cache volumes is simple:

```bash
# Add a new Node.js project to caching
JS_VOLUMES+=(
    "bucketeer-new-project-node-modules"
)

# Add a new Go cache (e.g., build cache)
GO_VOLUMES+=(
    "bucketeer-go-build-cache"
)
```

The `CACHE_VOLUMES` array will automatically include new volumes, and all cache manager commands (`status`, `clear-go`, `clear-js`) will work with the new volumes without any additional changes.

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
[SUCCESS] Evaluation TypeScript: Node.js dependencies are up to date
[SUCCESS] 🎉 All dependencies are cached and up to date! Setup completed in seconds.
[INFO] Cache volumes will persist dependencies for next container restart
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
[INFO] Installing Go development tools...
[SUCCESS] Go tools installed successfully
[INFO] Installing Dashboard dependencies...
[SUCCESS] Dashboard dependencies installed
[SUCCESS] 🎉 Post-attach setup completed!
[INFO] Cache volumes will persist dependencies for next container restart
``` 