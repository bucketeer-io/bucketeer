# Database Migration

For migration, we use the [Atlas tool](https://github.com/ariga/atlas).

When you install the Bucketeer Helm application, the migration runs automatically before the Bucketeer application is installed.
If the migration fails, it won't install the application.

## Prerequisite

Ensure that you have started the Minikube in the dev container and that the `localenv-mysql-0` pod is running without errors.

## Connecting to MySQL

To connect to MySQL, you must enter the `localenv-mysql-0` pod and then connect to MySQL.

```shell
kubectl exec -it localenv-mysql-0 -- /bin/sh
mysql -h localhost -u bucketeer -p bucketeer
```

Once you are logged in, you can make your changes in the Database.

## Creating Migration

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

## Pushing the migration file

Create a Pull Request to push the file so that the Bucketeer helm chart can migrate it when installed.
