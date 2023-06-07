// package: bucketeer.push
// file: proto/push/service.proto

var proto_push_service_pb = require("../../proto/push/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PushService = (function () {
  function PushService() {}
  PushService.serviceName = "bucketeer.push.PushService";
  return PushService;
}());

PushService.ListPushes = {
  methodName: "ListPushes",
  service: PushService,
  requestStream: false,
  responseStream: false,
  requestType: proto_push_service_pb.ListPushesRequest,
  responseType: proto_push_service_pb.ListPushesResponse
};

PushService.CreatePush = {
  methodName: "CreatePush",
  service: PushService,
  requestStream: false,
  responseStream: false,
  requestType: proto_push_service_pb.CreatePushRequest,
  responseType: proto_push_service_pb.CreatePushResponse
};

PushService.DeletePush = {
  methodName: "DeletePush",
  service: PushService,
  requestStream: false,
  responseStream: false,
  requestType: proto_push_service_pb.DeletePushRequest,
  responseType: proto_push_service_pb.DeletePushResponse
};

PushService.UpdatePush = {
  methodName: "UpdatePush",
  service: PushService,
  requestStream: false,
  responseStream: false,
  requestType: proto_push_service_pb.UpdatePushRequest,
  responseType: proto_push_service_pb.UpdatePushResponse
};

exports.PushService = PushService;

function PushServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PushServiceClient.prototype.listPushes = function listPushes(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PushService.ListPushes, {
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

PushServiceClient.prototype.createPush = function createPush(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PushService.CreatePush, {
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

PushServiceClient.prototype.deletePush = function deletePush(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PushService.DeletePush, {
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

PushServiceClient.prototype.updatePush = function updatePush(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PushService.UpdatePush, {
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

exports.PushServiceClient = PushServiceClient;

