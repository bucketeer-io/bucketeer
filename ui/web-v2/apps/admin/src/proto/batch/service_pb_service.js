// package: bucketeer.batch
// file: proto/batch/service.proto

var proto_batch_service_pb = require("../../proto/batch/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var BatchService = (function () {
  function BatchService() {}
  BatchService.serviceName = "bucketeer.batch.BatchService";
  return BatchService;
}());

BatchService.ExecuteBatchJob = {
  methodName: "ExecuteBatchJob",
  service: BatchService,
  requestStream: false,
  responseStream: false,
  requestType: proto_batch_service_pb.BatchJobRequest,
  responseType: proto_batch_service_pb.BatchJobResponse
};

BatchService.CurrentMigrationVersion = {
  methodName: "CurrentMigrationVersion",
  service: BatchService,
  requestStream: false,
  responseStream: false,
  requestType: proto_batch_service_pb.MigrationVersionRequest,
  responseType: proto_batch_service_pb.MigrationVersionResponse
};

BatchService.Migrate = {
  methodName: "Migrate",
  service: BatchService,
  requestStream: false,
  responseStream: false,
  requestType: proto_batch_service_pb.MigrationRequest,
  responseType: proto_batch_service_pb.MigrationResponse
};

exports.BatchService = BatchService;

function BatchServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

BatchServiceClient.prototype.executeBatchJob = function executeBatchJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BatchService.ExecuteBatchJob, {
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

BatchServiceClient.prototype.currentMigrationVersion = function currentMigrationVersion(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BatchService.CurrentMigrationVersion, {
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

BatchServiceClient.prototype.migrate = function migrate(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BatchService.Migrate, {
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

exports.BatchServiceClient = BatchServiceClient;

