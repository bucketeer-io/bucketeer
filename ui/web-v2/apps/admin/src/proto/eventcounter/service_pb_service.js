// package: bucketeer.eventcounter
// file: proto/eventcounter/service.proto

var proto_eventcounter_service_pb = require("../../proto/eventcounter/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var EventCounterService = (function () {
  function EventCounterService() {}
  EventCounterService.serviceName = "bucketeer.eventcounter.EventCounterService";
  return EventCounterService;
}());

EventCounterService.GetExperimentEvaluationCount = {
  methodName: "GetExperimentEvaluationCount",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetExperimentEvaluationCountRequest,
  responseType: proto_eventcounter_service_pb.GetExperimentEvaluationCountResponse
};

EventCounterService.GetEvaluationTimeseriesCount = {
  methodName: "GetEvaluationTimeseriesCount",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetEvaluationTimeseriesCountRequest,
  responseType: proto_eventcounter_service_pb.GetEvaluationTimeseriesCountResponse
};

EventCounterService.GetExperimentResult = {
  methodName: "GetExperimentResult",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetExperimentResultRequest,
  responseType: proto_eventcounter_service_pb.GetExperimentResultResponse
};

EventCounterService.ListExperimentResults = {
  methodName: "ListExperimentResults",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.ListExperimentResultsRequest,
  responseType: proto_eventcounter_service_pb.ListExperimentResultsResponse
};

EventCounterService.GetExperimentGoalCount = {
  methodName: "GetExperimentGoalCount",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetExperimentGoalCountRequest,
  responseType: proto_eventcounter_service_pb.GetExperimentGoalCountResponse
};

EventCounterService.GetMAUCount = {
  methodName: "GetMAUCount",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetMAUCountRequest,
  responseType: proto_eventcounter_service_pb.GetMAUCountResponse
};

EventCounterService.SummarizeMAUCounts = {
  methodName: "SummarizeMAUCounts",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.SummarizeMAUCountsRequest,
  responseType: proto_eventcounter_service_pb.SummarizeMAUCountsResponse
};

EventCounterService.GetOpsEvaluationUserCount = {
  methodName: "GetOpsEvaluationUserCount",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetOpsEvaluationUserCountRequest,
  responseType: proto_eventcounter_service_pb.GetOpsEvaluationUserCountResponse
};

EventCounterService.GetOpsGoalUserCount = {
  methodName: "GetOpsGoalUserCount",
  service: EventCounterService,
  requestStream: false,
  responseStream: false,
  requestType: proto_eventcounter_service_pb.GetOpsGoalUserCountRequest,
  responseType: proto_eventcounter_service_pb.GetOpsGoalUserCountResponse
};

exports.EventCounterService = EventCounterService;

function EventCounterServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

EventCounterServiceClient.prototype.getExperimentEvaluationCount = function getExperimentEvaluationCount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetExperimentEvaluationCount, {
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

EventCounterServiceClient.prototype.getEvaluationTimeseriesCount = function getEvaluationTimeseriesCount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetEvaluationTimeseriesCount, {
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

EventCounterServiceClient.prototype.getExperimentResult = function getExperimentResult(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetExperimentResult, {
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

EventCounterServiceClient.prototype.listExperimentResults = function listExperimentResults(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.ListExperimentResults, {
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

EventCounterServiceClient.prototype.getExperimentGoalCount = function getExperimentGoalCount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetExperimentGoalCount, {
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

EventCounterServiceClient.prototype.getMAUCount = function getMAUCount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetMAUCount, {
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

EventCounterServiceClient.prototype.summarizeMAUCounts = function summarizeMAUCounts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.SummarizeMAUCounts, {
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

EventCounterServiceClient.prototype.getOpsEvaluationUserCount = function getOpsEvaluationUserCount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetOpsEvaluationUserCount, {
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

EventCounterServiceClient.prototype.getOpsGoalUserCount = function getOpsGoalUserCount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EventCounterService.GetOpsGoalUserCount, {
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

exports.EventCounterServiceClient = EventCounterServiceClient;

