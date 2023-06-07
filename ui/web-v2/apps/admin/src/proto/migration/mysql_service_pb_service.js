// package: bucketeer.migration
// file: proto/migration/mysql_service.proto

var proto_migration_mysql_service_pb = require("../../proto/migration/mysql_service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var MigrationMySQLService = (function () {
  function MigrationMySQLService() {}
  MigrationMySQLService.serviceName = "bucketeer.migration.MigrationMySQLService";
  return MigrationMySQLService;
}());

MigrationMySQLService.MigrateAllMasterSchema = {
  methodName: "MigrateAllMasterSchema",
  service: MigrationMySQLService,
  requestStream: false,
  responseStream: false,
  requestType: proto_migration_mysql_service_pb.MigrateAllMasterSchemaRequest,
  responseType: proto_migration_mysql_service_pb.MigrateAllMasterSchemaResponse
};

MigrationMySQLService.RollbackMasterSchema = {
  methodName: "RollbackMasterSchema",
  service: MigrationMySQLService,
  requestStream: false,
  responseStream: false,
  requestType: proto_migration_mysql_service_pb.RollbackMasterSchemaRequest,
  responseType: proto_migration_mysql_service_pb.RollbackMasterSchemaResponse
};

exports.MigrationMySQLService = MigrationMySQLService;

function MigrationMySQLServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

MigrationMySQLServiceClient.prototype.migrateAllMasterSchema = function migrateAllMasterSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MigrationMySQLService.MigrateAllMasterSchema, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

MigrationMySQLServiceClient.prototype.rollbackMasterSchema = function rollbackMasterSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MigrationMySQLService.RollbackMasterSchema, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.MigrationMySQLServiceClient = MigrationMySQLServiceClient;

