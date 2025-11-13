# Migration Guide: Bitnami to Official Images

This guide explains what developers need to do when pulling the changes that migrate from Bitnami Helm charts to custom templates using official Docker images.

## What Changed

- **Removed**: Bitnami Helm chart dependencies for Redis and MySQL
- **Added**: Custom Helm templates using official `redis:7.2-alpine` and `mysql:8.0` images
- **Updated**: `docker-compose/compose.yml` to remove deprecated MySQL authentication plugin

## Migration Steps

### For Developers Using Minikube

**You do NOT need to:**
- ❌ Rebuild the dev container
- ❌ Delete minikube (your data will be preserved in PVCs)
- ❌ Reinstall everything from scratch

**Option 1: Automatic (Recommended)**

Just run the updated `make start-minikube` command - it now handles upgrades automatically:

```bash
git pull
make start-minikube
```

The Makefile will:
1. Update Helm dependencies
2. Upgrade the existing release (or install if new)
3. Handle the migration automatically

**Option 2: Manual Steps**

If you prefer to do it manually or if the automatic upgrade has issues:

1. **Pull the latest changes:**
   ```bash
   git pull
   ```

2. **Delete old Bitnami StatefulSets** (they can't be upgraded directly):
   ```bash
   kubectl delete statefulset localenv-mysql localenv-redis-master -n default
   ```

3. **Update Helm dependencies and upgrade:**
   ```bash
   cd manifests/localenv
   helm dependency update
   helm upgrade localenv . -n default
   ```

4. **Verify everything is running:**
   ```bash
   kubectl get pods -n default | grep -E "mysql|redis"
   # Should show:
   # localenv-localenv-mysql-0    1/1   Running
   # localenv-localenv-redis-0    1/1   Running
   ```

### For Developers Using Docker Compose

**No action needed!** The docker-compose setup already uses official images and will continue working as before.

### Troubleshooting

**If pods fail to start:**

1. **Check if old PVCs are causing issues:**
   ```bash
   kubectl get pvc -n default | grep -E "mysql|redis"
   ```

2. **If needed, delete and recreate PVCs** (⚠️ This will delete your data):
   ```bash
   kubectl delete pvc mysql-data-localenv-mysql-0 redis-data-localenv-redis-0 -n default
   kubectl delete pod localenv-mysql-0 localenv-redis-0 -n default
   # Pods will be recreated with fresh PVCs
   ```

3. **Check pod logs:**
   ```bash
   kubectl logs localenv-mysql-0 -n default
   kubectl logs localenv-redis-0 -n default
   ```

### Data Migration (Optional)

If you had important data in the old Bitnami deployments and want to migrate it:

1. **Backup data from old PVCs** (if they still exist)
2. **Restore to new PVCs** using standard MySQL/Redis backup/restore procedures

For most development environments, starting fresh is usually fine.

## Benefits of This Change

- ✅ **No dependency on Bitnami Legacy repository** (which stops receiving updates after Aug 2025)
- ✅ **Uses official, maintained images** (same as docker-compose)
- ✅ **Better long-term support** and security updates
- ✅ **Consistent** between docker-compose and Kubernetes setups

