# Database Schema Migration

For the database schema migration, we use the [Atlas tool](https://github.com/ariga/atlas).

When you install the Bucketeer Helm application, the migration runs automatically before the Bucketeer application is installed.
If the migration fails, it won't install the application.

## Supported Databases

- **MySQL** - Primary storage (production ready)
- **PostgreSQL** - Primary storage (alternative option)

## Prerequisite

Ensure that you have [started the Minikube](https://github.com/bucketeer-io/bucketeer/blob/main/DEVELOPMENT.md#set-up-minikube-and-services-that-bucketeer-depends-on) in the dev container and that the `localenv-mysql-0` or `localenv-postgres-0` pod is running without errors.

# MySQL Migration

## 1- Connecting to MySQL

To connect to MySQL, you must enter the `localenv-mysql-0` pod and then connect to MySQL.
The password can be found in the [values.dev.yaml](https://github.com/bucketeer-io/bucketeer/blob/main/manifests/bucketeer/values.dev.yaml) in the field `global.operationalDatabase.mysql`.<br />
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

---

# PostgreSQL Migration

## Overview

PostgreSQL migrations are provided as a consolidated initialization file that creates the complete schema in one step. This is the recommended approach for new PostgreSQL deployments.

## 1- Connecting to PostgreSQL

```shell
kubectl exec -it localenv-postgres-0 -- /bin/sh
psql -h localhost -U bucketeer -d bucketeer
```

## 2- Running the Migration

Migrations ship as **two images** (same tag scheme as the MySQL image):

- **`ghcr.io/bucketeer-io/bucketeer-migration`** — built from `migration/Dockerfile` (`migration/mysql` → `/migrations`).
- **`ghcr.io/bucketeer-io/bucketeer-migration-postgres`** — built from `migration/Dockerfile.postgres` (`migration/postgres` → `/migrations`).

For Helm, set `migration.image.repository` to the postgres image when applying PostgreSQL migrations (and set `dbUrl` / `dbBaseline` for Postgres as in `values.yaml` comments).

When installing Bucketeer, the pre-install job runs Atlas against the image you select and creates or updates the schema.
