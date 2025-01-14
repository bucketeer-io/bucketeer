// package: bucketeer.tag
// file: proto/tag/service.proto

var proto_tag_service_pb = require('../../proto/tag/service_pb');
var grpc = require('@improbable-eng/grpc-web').grpc;

var TagService = (function () {
  function TagService() {}
  TagService.serviceName = 'bucketeer.tag.TagService';
  return TagService;
})();

TagService.ListTags = {
  methodName: 'ListTags',
  service: TagService,
  requestStream: false,
  responseStream: false,
  requestType: proto_tag_service_pb.ListTagsRequest,
  responseType: proto_tag_service_pb.ListTagsResponse
};

TagService.CreateTag = {
  methodName: 'CreateTag',
  service: TagService,
  requestStream: false,
  responseStream: false,
  requestType: proto_tag_service_pb.CreateTagRequest,
  responseType: proto_tag_service_pb.CreateTagResponse
};

TagService.DeleteTag = {
  methodName: 'DeleteTag',
  service: TagService,
  requestStream: false,
  responseStream: false,
  requestType: proto_tag_service_pb.DeleteTagRequest,
  responseType: proto_tag_service_pb.DeleteTagResponse
};

exports.TagService = TagService;

function TagServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

TagServiceClient.prototype.listTags = function listTags(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TagService.ListTags, {
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

TagServiceClient.prototype.createTag = function createTag(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TagService.CreateTag, {
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

TagServiceClient.prototype.deleteTag = function deleteTag(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TagService.DeleteTag, {
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

exports.TagServiceClient = TagServiceClient;
