# Minikube Image Caching

This document explains the image caching optimizations implemented to significantly reduce minikube startup times and avoid re-downloading images after `minikube delete`.

## 🚀 Quick Start

### **Fresh Devcontainer (Recommended)**
```bash
# One command does everything: pre-cache + setup minikube
make -C tools/dev setup-minikube-cached
```

### **Alternative: Manual Steps**
```bash
# 1. Pre-cache images first (optional, setup-minikube-cached does this)
make -C tools/dev pre-cache-docker-images

# 2. Then setup minikube
make -C tools/dev setup-minikube-cached
```

## 📋 Caching Strategy

### 1. Persistent Docker Volume Cache
The devcontainer now mounts a persistent Docker volume that survives container rebuilds:

```json
"source=bucketeer-docker-cache,target=/var/lib/docker,type=volume"
```

This means all Docker images remain cached even when the devcontainer is rebuilt.

### 2. Persistent Minikube Configuration
Minikube configuration and cache persist across devcontainer rebuilds:

```json
"source=bucketeer-minikube-cache,target=/home/codespace/.minikube,type=volume"
```

### 3. Pre-Cache Critical Images
New Makefile target `pre-cache-docker-images` downloads all necessary images to the Docker cache:

- `gcr.io/k8s-minikube/kicbase:v0.0.47` - Minikube base image (largest download)
- `registry.k8s.io/ingress-nginx/controller:v1.12.2` - Ingress controller
- `registry.k8s.io/ingress-nginx/kube-webhook-certgen:v1.5.3` - Webhook cert generator
- `gcr.io/k8s-minikube/storage-provisioner:v5` - Storage provisioner
- `mysql:8.0` - Database
- `ghcr.io/bucketeer-io/bigquery-emulator:0.6.6` - BigQuery emulator (forked version with ARM64 support)
- `gcr.io/google.com/cloudsdktool/google-cloud-cli:545.0.0` - Cloud SDK (full version, includes Java for PubSub emulator)

## 🛠️ Available Commands

### New Commands

```bash
# Pre-cache all Docker images (recommended for first-time setup)
make -C tools/dev pre-cache-docker-images

# Setup minikube with pre-cached images
make -C tools/dev setup-minikube-cached

# Show cache status and disk usage
make -C tools/dev cache-status

# Light cleanup: remove unused images only (keeps volumes)
make -C tools/dev clean-cache-light

# Full cleanup: remove all cache and volumes (requires confirmation)
make -C tools/dev clean-cache
```

### Enhanced Existing Commands

```bash
# Now uses --cache-images flag for faster startups
make start-minikube

# Original minikube image caching (requires running minikube)
make -C tools/dev cache-minikube-images
```

## ⚡ Performance Impact

**Before optimization:**
- Fresh minikube setup: ~3-5 minutes (downloading 400+ MB)
- After `minikube delete`: ~3-5 minutes (re-downloading everything)

**After optimization:**
- Fresh minikube setup: ~3-5 minutes (same, initial download required)
- After `minikube delete`: ~30-60 seconds (using cached images!)
- After devcontainer rebuild: ~30-60 seconds (persistent cache!)

## 🔄 Workflow Recommendations

### First Time Setup
1. Start devcontainer
2. Run `make -C tools/dev setup-minikube-cached` (does everything in one command)

### Daily Development
- Use `make start-minikube` for regular startup
- If you need to reset minikube: `minikube delete && make start-minikube`
- Images remain cached, so restarts are much faster

### Troubleshooting Network Issues
If you encounter network issues and need to delete minikube:

1. `minikube delete` - Only deletes the cluster, not Docker images
2. `make start-minikube` - Uses cached images for fast restart

If Docker images are corrupted or you need to free space:
1. `make -C tools/dev clean-cache` - Interactive cleanup (removes everything)
2. `make -C tools/dev pre-cache-docker-images` - Re-cache images
3. `make start-minikube` - Start with fresh cache

For lighter cleanup (keeps volumes):
1. `make -C tools/dev clean-cache-light` - Remove unused images only

## 🔍 Cache Status

Check what's cached:

```bash
# Comprehensive cache status (recommended)
make -C tools/dev cache-status

# Manual checks
docker images              # List cached Docker images
minikube cache list        # Check minikube cache (if running)
docker system df           # Check disk usage
```

## 📁 Volume Management

The persistent volumes can be managed via Docker:

```bash
# List volumes
docker volume ls | grep bucketeer

# Inspect cache volume
docker volume inspect bucketeer-docker-cache
docker volume inspect bucketeer-minikube-cache

# Remove cache volumes (forces fresh download next time)
docker volume rm bucketeer-docker-cache bucketeer-minikube-cache
```

## 🎯 Benefits

1. **Faster Development Cycles**: No more waiting for image downloads
2. **Offline Development**: Works without internet after initial cache
3. **Reliable Builds**: Cached images reduce external dependency failures
4. **Resource Efficient**: Shared cache across all devcontainer instances
