// package: bucketeer.auditlog
// file: proto/auditlog/service.proto

var proto_auditlog_service_pb = require("../../proto/auditlog/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var AuditLogService = (function () {
  function AuditLogService() {}
  AuditLogService.serviceName = "bucketeer.auditlog.AuditLogService";
  return AuditLogService;
}());

AuditLogService.ListAuditLogs = {
  methodName: "ListAuditLogs",
  service: AuditLogService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auditlog_service_pb.ListAuditLogsRequest,
  responseType: proto_auditlog_service_pb.ListAuditLogsResponse
};

AuditLogService.ListAdminAuditLogs = {
  methodName: "ListAdminAuditLogs",
  service: AuditLogService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auditlog_service_pb.ListAdminAuditLogsRequest,
  responseType: proto_auditlog_service_pb.ListAdminAuditLogsResponse
};

AuditLogService.ListFeatureHistory = {
  methodName: "ListFeatureHistory",
  service: AuditLogService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auditlog_service_pb.ListFeatureHistoryRequest,
  responseType: proto_auditlog_service_pb.ListFeatureHistoryResponse
};

exports.AuditLogService = AuditLogService;

function AuditLogServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

AuditLogServiceClient.prototype.listAuditLogs = function listAuditLogs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuditLogService.ListAuditLogs, {
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

AuditLogServiceClient.prototype.listAdminAuditLogs = function listAdminAuditLogs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuditLogService.ListAdminAuditLogs, {
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

AuditLogServiceClient.prototype.listFeatureHistory = function listFeatureHistory(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuditLogService.ListFeatureHistory, {
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

exports.AuditLogServiceClient = AuditLogServiceClient;

