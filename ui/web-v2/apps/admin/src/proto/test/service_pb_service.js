// package: bucketeer.test
// file: proto/test/service.proto

var proto_test_service_pb = require("../../proto/test/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var TestService = (function () {
  function TestService() {}
  TestService.serviceName = "bucketeer.test.TestService";
  return TestService;
}());

TestService.Test = {
  methodName: "Test",
  service: TestService,
  requestStream: false,
  responseStream: false,
  requestType: proto_test_service_pb.TestRequest,
  responseType: proto_test_service_pb.TestResponse
};

exports.TestService = TestService;

function TestServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

TestServiceClient.prototype.test = function test(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TestService.Test, {
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

exports.TestServiceClient = TestServiceClient;

