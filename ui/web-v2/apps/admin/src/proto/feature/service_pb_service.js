// package: bucketeer.feature
// file: proto/feature/service.proto

var proto_feature_service_pb = require("../../proto/feature/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var FeatureService = (function () {
  function FeatureService() {}
  FeatureService.serviceName = "bucketeer.feature.FeatureService";
  return FeatureService;
}());

FeatureService.GetFeature = {
  methodName: "GetFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.GetFeatureRequest,
  responseType: proto_feature_service_pb.GetFeatureResponse
};

FeatureService.GetFeatures = {
  methodName: "GetFeatures",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.GetFeaturesRequest,
  responseType: proto_feature_service_pb.GetFeaturesResponse
};

FeatureService.ListFeatures = {
  methodName: "ListFeatures",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.ListFeaturesRequest,
  responseType: proto_feature_service_pb.ListFeaturesResponse
};

FeatureService.ListEnabledFeatures = {
  methodName: "ListEnabledFeatures",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.ListEnabledFeaturesRequest,
  responseType: proto_feature_service_pb.ListEnabledFeaturesResponse
};

FeatureService.CreateFeature = {
  methodName: "CreateFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.CreateFeatureRequest,
  responseType: proto_feature_service_pb.CreateFeatureResponse
};

FeatureService.EnableFeature = {
  methodName: "EnableFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.EnableFeatureRequest,
  responseType: proto_feature_service_pb.EnableFeatureResponse
};

FeatureService.DisableFeature = {
  methodName: "DisableFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.DisableFeatureRequest,
  responseType: proto_feature_service_pb.DisableFeatureResponse
};

FeatureService.ArchiveFeature = {
  methodName: "ArchiveFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.ArchiveFeatureRequest,
  responseType: proto_feature_service_pb.ArchiveFeatureResponse
};

FeatureService.UnarchiveFeature = {
  methodName: "UnarchiveFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.UnarchiveFeatureRequest,
  responseType: proto_feature_service_pb.UnarchiveFeatureResponse
};

FeatureService.DeleteFeature = {
  methodName: "DeleteFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.DeleteFeatureRequest,
  responseType: proto_feature_service_pb.DeleteFeatureResponse
};

FeatureService.UpdateFeatureDetails = {
  methodName: "UpdateFeatureDetails",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.UpdateFeatureDetailsRequest,
  responseType: proto_feature_service_pb.UpdateFeatureDetailsResponse
};

FeatureService.UpdateFeatureVariations = {
  methodName: "UpdateFeatureVariations",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.UpdateFeatureVariationsRequest,
  responseType: proto_feature_service_pb.UpdateFeatureVariationsResponse
};

FeatureService.UpdateFeatureTargeting = {
  methodName: "UpdateFeatureTargeting",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.UpdateFeatureTargetingRequest,
  responseType: proto_feature_service_pb.UpdateFeatureTargetingResponse
};

FeatureService.CloneFeature = {
  methodName: "CloneFeature",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.CloneFeatureRequest,
  responseType: proto_feature_service_pb.CloneFeatureResponse
};

FeatureService.CreateSegment = {
  methodName: "CreateSegment",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.CreateSegmentRequest,
  responseType: proto_feature_service_pb.CreateSegmentResponse
};

FeatureService.GetSegment = {
  methodName: "GetSegment",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.GetSegmentRequest,
  responseType: proto_feature_service_pb.GetSegmentResponse
};

FeatureService.ListSegments = {
  methodName: "ListSegments",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.ListSegmentsRequest,
  responseType: proto_feature_service_pb.ListSegmentsResponse
};

FeatureService.DeleteSegment = {
  methodName: "DeleteSegment",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.DeleteSegmentRequest,
  responseType: proto_feature_service_pb.DeleteSegmentResponse
};

FeatureService.UpdateSegment = {
  methodName: "UpdateSegment",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.UpdateSegmentRequest,
  responseType: proto_feature_service_pb.UpdateSegmentResponse
};

FeatureService.AddSegmentUser = {
  methodName: "AddSegmentUser",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.AddSegmentUserRequest,
  responseType: proto_feature_service_pb.AddSegmentUserResponse
};

FeatureService.DeleteSegmentUser = {
  methodName: "DeleteSegmentUser",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.DeleteSegmentUserRequest,
  responseType: proto_feature_service_pb.DeleteSegmentUserResponse
};

FeatureService.GetSegmentUser = {
  methodName: "GetSegmentUser",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.GetSegmentUserRequest,
  responseType: proto_feature_service_pb.GetSegmentUserResponse
};

FeatureService.ListSegmentUsers = {
  methodName: "ListSegmentUsers",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.ListSegmentUsersRequest,
  responseType: proto_feature_service_pb.ListSegmentUsersResponse
};

FeatureService.BulkUploadSegmentUsers = {
  methodName: "BulkUploadSegmentUsers",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.BulkUploadSegmentUsersRequest,
  responseType: proto_feature_service_pb.BulkUploadSegmentUsersResponse
};

FeatureService.BulkDownloadSegmentUsers = {
  methodName: "BulkDownloadSegmentUsers",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.BulkDownloadSegmentUsersRequest,
  responseType: proto_feature_service_pb.BulkDownloadSegmentUsersResponse
};

FeatureService.EvaluateFeatures = {
  methodName: "EvaluateFeatures",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.EvaluateFeaturesRequest,
  responseType: proto_feature_service_pb.EvaluateFeaturesResponse
};

FeatureService.ListTags = {
  methodName: "ListTags",
  service: FeatureService,
  requestStream: false,
  responseStream: false,
  requestType: proto_feature_service_pb.ListTagsRequest,
  responseType: proto_feature_service_pb.ListTagsResponse
};

exports.FeatureService = FeatureService;

function FeatureServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

FeatureServiceClient.prototype.getFeature = function getFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.GetFeature, {
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

FeatureServiceClient.prototype.getFeatures = function getFeatures(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.GetFeatures, {
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

FeatureServiceClient.prototype.listFeatures = function listFeatures(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.ListFeatures, {
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

FeatureServiceClient.prototype.listEnabledFeatures = function listEnabledFeatures(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.ListEnabledFeatures, {
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

FeatureServiceClient.prototype.createFeature = function createFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.CreateFeature, {
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

FeatureServiceClient.prototype.enableFeature = function enableFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.EnableFeature, {
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

FeatureServiceClient.prototype.disableFeature = function disableFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.DisableFeature, {
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

FeatureServiceClient.prototype.archiveFeature = function archiveFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.ArchiveFeature, {
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

FeatureServiceClient.prototype.unarchiveFeature = function unarchiveFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.UnarchiveFeature, {
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

FeatureServiceClient.prototype.deleteFeature = function deleteFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.DeleteFeature, {
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

FeatureServiceClient.prototype.updateFeatureDetails = function updateFeatureDetails(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.UpdateFeatureDetails, {
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

FeatureServiceClient.prototype.updateFeatureVariations = function updateFeatureVariations(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.UpdateFeatureVariations, {
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

FeatureServiceClient.prototype.updateFeatureTargeting = function updateFeatureTargeting(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.UpdateFeatureTargeting, {
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

FeatureServiceClient.prototype.cloneFeature = function cloneFeature(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.CloneFeature, {
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

FeatureServiceClient.prototype.createSegment = function createSegment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.CreateSegment, {
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

FeatureServiceClient.prototype.getSegment = function getSegment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.GetSegment, {
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

FeatureServiceClient.prototype.listSegments = function listSegments(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.ListSegments, {
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

FeatureServiceClient.prototype.deleteSegment = function deleteSegment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.DeleteSegment, {
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

FeatureServiceClient.prototype.updateSegment = function updateSegment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.UpdateSegment, {
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

FeatureServiceClient.prototype.addSegmentUser = function addSegmentUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.AddSegmentUser, {
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

FeatureServiceClient.prototype.deleteSegmentUser = function deleteSegmentUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.DeleteSegmentUser, {
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

FeatureServiceClient.prototype.getSegmentUser = function getSegmentUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.GetSegmentUser, {
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

FeatureServiceClient.prototype.listSegmentUsers = function listSegmentUsers(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.ListSegmentUsers, {
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

FeatureServiceClient.prototype.bulkUploadSegmentUsers = function bulkUploadSegmentUsers(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.BulkUploadSegmentUsers, {
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

FeatureServiceClient.prototype.bulkDownloadSegmentUsers = function bulkDownloadSegmentUsers(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.BulkDownloadSegmentUsers, {
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

FeatureServiceClient.prototype.evaluateFeatures = function evaluateFeatures(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.EvaluateFeatures, {
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

FeatureServiceClient.prototype.listTags = function listTags(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(FeatureService.ListTags, {
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

exports.FeatureServiceClient = FeatureServiceClient;

