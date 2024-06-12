// package: bucketeer.backend
// file: proto/publicapi/service.proto

var proto_publicapi_service_pb = require("../../proto/publicapi/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PublicAPIService = (function () {
  function PublicAPIService() {}
  PublicAPIService.serviceName = "bucketeer.backend.PublicAPIService";
  return PublicAPIService;
}());

PublicAPIService.GetFeature = {
  methodName: "GetFeature",
  service: PublicAPIService,
  requestStream: false,
  responseStream: false,
  requestType: proto_publicapi_service_pb.GetFeatureRequest,
  responseType: proto_publicapi_service_pb.GetFeatureResponse
};

PublicAPIService.UpdateFeature = {
  methodName: "UpdateFeature",
  service: PublicAPIService,
  requestStream: false,
  responseStream: false,
  requestType: proto_publicapi_service_pb.UpdateFeatureRequest,
  responseType: proto_publicapi_service_pb.UpdateFeatureResponse
};

exports.PublicAPIService = PublicAPIService;

function PublicAPIServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PublicAPIServiceClient.prototype.getFeature = function getFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PublicAPIService.GetFeature, {
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

PublicAPIServiceClient.prototype.updateFeature = function updateFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PublicAPIService.UpdateFeature, {
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

exports.PublicAPIServiceClient = PublicAPIServiceClient;

