# Database Schema Migration

For the database schema migration, we use the [Atlas tool](https://github.com/ariga/atlas).

When you install the Bucketeer Helm application, the migration runs automatically before the Bucketeer application is
installed.
If the migration fails, it won't install the application.

## Prerequisite

Ensure that you
have [started the Minikube](https://github.com/bucketeer-io/bucketeer/blob/main/DEVELOPMENT.md#set-up-minikube-and-services-that-bucketeer-depends-on)
in the dev container and that the `localenv-mysql-0` pod is running without errors.

## 1- Connecting to MySQL

To connect to MySQL, you must enter the `localenv-mysql-0` pod and then connect to MySQL.
The password can be found in
the [values.dev.yaml](https://github.com/bucketeer-io/bucketeer/blob/main/manifests/bucketeer/values.dev.yaml) in the
field `mysqlPass`.<br />
In case you changed the default value, then use the password you set.

```shell
kubectl exec -it localenv-mysql-0 -- /bin/sh
mysql -h localhost -u bucketeer -p bucketeer
```

Once you are logged in, you can make your changes in the Database.

## 2- Creating Migration File

### 2.1 - Generate Migration File

To create the migration file, you must port-forward the `localenv-mysql` service before creating it.

```shell
kubectl port-forward svc/localenv-mysql 3306:3306
```

The following command will create the migration file.

```shell
make create-migration NAME=<MIGRATION_FILE_NAME> HOST=localhost USER=bucketeer PASS=bucketeer PORT=3306 DB=bucketeer
```

For the migration file name, please use one of the following prefixes.

- **create:** Used for new tables. E.g. `create_xxx_table`
- **update:** Used when you alter an existing table. E.g. `update_xxxx_table`
- **drop:** Used when you drop a table. E.g. `drop_xxx_table`

After creating it, ensure you see the new file in the `migration/mysql` directory.

### 2.2 - Create Migration File manually

In case no structure changes are executed, you can create the migration file manually.

Create a new migration file

```shell
atlas migrate new <file_name> --dir "file://migration/mysql"
```

Edit the migration file then update atlas migration files hash

```shell
atlas migrate hash --dir "file://migration/mysql"
```

If the migration sql statement changed again, rerun the above command to update the hash.

## 3- Pushing Migration File

Create a Pull Request to push the file so that the Bucketeer helm chart can migrate it when installed.

## 4- Rolling Back Migrations

Atlas supports dynamic migration rollback, which automatically computes the necessary changes to revert migrations based
on the current database state.

### 4.1 - Prerequisites

For local development or Kubernetes environments, ensure you can connect to the target MySQL database:

```shell
# For Kubernetes/Minikube environments
kubectl port-forward svc/localenv-mysql 3306:3306
```

### 4.2 - Preview Rollback (Dry Run)

**Always preview rollback changes before executing them.**

#### Preview last migration rollback:

```shell
make check-rollback-migration USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

#### Preview last N migrations rollback:

```shell
# Rollback last 3 migrations
make check-rollback-migration COUNT=3 USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

#### Preview rollback to specific version:

```shell
# Rollback to version 20240815043128 (reverts all migrations after this version)
make check-rollback-migration VERSION=20240815043128 USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

### 4.3 - Execute Rollback

⚠️ **WARNING**: Rolling back migrations can result in data loss. Always:

1. Take a database backup before rollback
2. Preview changes with dry-run first
3. Test on non-production environment
4. Verify application compatibility with rolled-back schema

#### Rollback last migration:

```shell
make rollback-migration USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

#### Rollback last N migrations:

```shell
# Rollback last 3 migrations
make rollback-migration COUNT=3 USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

#### Rollback to specific version:

```shell
# Rollback to version 20240815043128
make rollback-migration VERSION=20240815043128 USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

### 4.4 - Check Migration Status

To view the current migration status and see which migrations are applied:

```shell
make migration-status USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

### 4.5 - How Rollback Works

Atlas uses **dynamic rollback computation**:

1. **No manual down files needed**: Atlas analyzes your current database state and migration files
2. **Automatic computation**: Calculates necessary SQL to revert to target state
3. **Safety checks**: Runs pre-migration checks to detect potential issues
4. **Drift detection**: Validates database state matches expected state
5. **Transactional**: Rollback is atomic when supported by database

### 4.6 - Finding Version Numbers

To find available migration versions:

```shell
# List all migration files with timestamps
ls -1 migration/mysql/*.sql

# Check current applied migrations
make migration-status USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
```

Version format: `YYYYMMDDHHMMSS` (e.g., `20240815043128`)

### 4.7 - Rollback Best Practices

1. **Always preview first**: Use `check-rollback-migration` to see what will be rolled back
   ```shell
   make check-rollback-migration USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
   ```

2. **Backup before rollback**: Create a database backup before executing any rollback
   ```shell
   mysqldump -h localhost -P 3306 -u bucketeer -p bucketeer > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

3. **Verify status**: Check migration state before and after rollback
   ```shell
   make migration-status USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer
   ```

4. **Test in staging**: Never rollback production without testing in staging first

5. **Coordinate with team**: Ensure no active migrations or deployments

6. **Document reason**: Record why rollback was necessary for future reference

### 4.8 - Troubleshooting Rollback

#### Error: "cannot rollback, database has drift"

**Solution**: Database state doesn't match expected state. Review manual changes or fix drift issues before rollback.

#### Error: "version not found"

**Solution**: Verify version number exists in `migration/mysql/` directory.

#### Error: "destructive changes detected"

**Solution**: Atlas detected potential data loss. Review dry-run output carefully and ensure you have a backup.

#### Rollback failed partially

**Solution**:

1. Check migration status: `make migration-status USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer`
2. Review database state manually
3. Restore from backup if needed
4. Contact team if in production environment

### 4.9 - Kubernetes/Production Rollback

⚠️ **Manual Operation Only**: Migration rollback in Kubernetes is **intentionally manual** to prevent accidental data
loss.

**Process**:

1. **Port-forward to MySQL**:
   ```shell
   kubectl port-forward -n <namespace> svc/localenv-mysql 3306:3306
   ```

2. **Preview rollback**:
   ```shell
   make check-rollback-migration USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer COUNT=<N or VERSION=version>
   ```

3. **Create database backup**:
   ```shell
   mysqldump -h localhost -P 3306 -u bucketeer -p bucketeer > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

4. **Execute rollback**:
   ```shell
   make rollback-migration USER=bucketeer PASS=bucketeer HOST=localhost PORT=3306 DB=bucketeer COUNT=<N or VERSION=version>
   ```

5. **Verify application compatibility**: Ensure services work with rolled-back schema

### 4.10 - Using Rollback Containers

For automated rollback workflows, Bucketeer provides pre-configured rollback containers for Kubernetes environments.

#### Kubernetes Rollback Job

The Kubernetes rollback uses a Job resource that must be manually enabled and applied.

**Prerequisites**:

- MySQL service must be accessible
- Helm values must be configured

**Usage:**

```bash
# Dry-run preview
helm template migration-rollback ./manifests/bucketeer-migration \
  --set rollback.enabled=true \
  --set rollback.dryRun=true \
  --set rollback.count=1 \
  --set dbUrl="mysql://bucketeer:bucketeer@localenv-mysql.default.svc:3306/bucketeer" \
  | kubectl apply -f -

# Execute rollback
helm template migration-rollback ./manifests/bucketeer-migration \
  --set rollback.enabled=true \
  --set rollback.dryRun=false \
  --set rollback.count=3 \
  --set dbUrl="mysql://bucketeer:bucketeer@localenv-mysql.default.svc:3306/bucketeer" \
  | kubectl apply -f -

# Rollback to specific version
helm template migration-rollback ./manifests/bucketeer-migration \
  --set rollback.enabled=true \
  --set rollback.dryRun=false \
  --set rollback.version=20240815043128 \
  --set dbUrl="mysql://bucketeer:bucketeer@localenv-mysql.default.svc:3306/bucketeer" \
  | kubectl apply -f -
```

**Note:** The migration chart uses the main `dbUrl` value for both migration and rollback operations.

#### Safety Features

The Kubernetes rollback container includes the following safety mechanisms:

1. **Dry-run Mode**: Default mode only previews changes (`rollback.dryRun=true`)
2. **Pre-flight Status Check**: Shows current migration state before rollback
3. **Dry-run Preview**: Displays SQL to be executed before actual rollback
4. **Post-rollback Status**: Shows final migration state after rollback (when executed)
5. **No Automatic Hooks**: Must be manually triggered
6. **Clear Error Messages**: Provides usage examples if safety checks fail

#### Rollback Container Workflow

```
1. Pre-flight   → Show current migration status
                  ↓
2. Dry Run      → Preview rollback changes
                  ↓
3. Execute      → Apply rollback (if dryRun=false)
                  ↓
4. Verify       → Show final status
```