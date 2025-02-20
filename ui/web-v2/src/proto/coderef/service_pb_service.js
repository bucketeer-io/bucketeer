// package: bucketeer.coderef
// file: proto/coderef/service.proto

var proto_coderef_service_pb = require('../../proto/coderef/service_pb');
var grpc = require('@improbable-eng/grpc-web').grpc;

var CodeReferenceService = (function () {
  function CodeReferenceService() {}
  CodeReferenceService.serviceName = 'bucketeer.coderef.CodeReferenceService';
  return CodeReferenceService;
})();

CodeReferenceService.GetCodeReference = {
  methodName: 'GetCodeReference',
  service: CodeReferenceService,
  requestStream: false,
  responseStream: false,
  requestType: proto_coderef_service_pb.GetCodeReferenceRequest,
  responseType: proto_coderef_service_pb.GetCodeReferenceResponse
};

CodeReferenceService.ListCodeReferences = {
  methodName: 'ListCodeReferences',
  service: CodeReferenceService,
  requestStream: false,
  responseStream: false,
  requestType: proto_coderef_service_pb.ListCodeReferencesRequest,
  responseType: proto_coderef_service_pb.ListCodeReferencesResponse
};

CodeReferenceService.CreateCodeReference = {
  methodName: 'CreateCodeReference',
  service: CodeReferenceService,
  requestStream: false,
  responseStream: false,
  requestType: proto_coderef_service_pb.CreateCodeReferenceRequest,
  responseType: proto_coderef_service_pb.CreateCodeReferenceResponse
};

CodeReferenceService.UpdateCodeReference = {
  methodName: 'UpdateCodeReference',
  service: CodeReferenceService,
  requestStream: false,
  responseStream: false,
  requestType: proto_coderef_service_pb.UpdateCodeReferenceRequest,
  responseType: proto_coderef_service_pb.UpdateCodeReferenceResponse
};

CodeReferenceService.DeleteCodeReference = {
  methodName: 'DeleteCodeReference',
  service: CodeReferenceService,
  requestStream: false,
  responseStream: false,
  requestType: proto_coderef_service_pb.DeleteCodeReferenceRequest,
  responseType: proto_coderef_service_pb.DeleteCodeReferenceResponse
};

exports.CodeReferenceService = CodeReferenceService;

function CodeReferenceServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

CodeReferenceServiceClient.prototype.getCodeReference =
  function getCodeReference(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(CodeReferenceService.GetCodeReference, {
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

CodeReferenceServiceClient.prototype.listCodeReferences =
  function listCodeReferences(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(CodeReferenceService.ListCodeReferences, {
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

CodeReferenceServiceClient.prototype.createCodeReference =
  function createCodeReference(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(CodeReferenceService.CreateCodeReference, {
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

CodeReferenceServiceClient.prototype.updateCodeReference =
  function updateCodeReference(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(CodeReferenceService.UpdateCodeReference, {
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

CodeReferenceServiceClient.prototype.deleteCodeReference =
  function deleteCodeReference(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(CodeReferenceService.DeleteCodeReference, {
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

exports.CodeReferenceServiceClient = CodeReferenceServiceClient;
