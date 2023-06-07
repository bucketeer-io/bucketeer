// package: bucketeer.experiment
// file: proto/experiment/service.proto

var proto_experiment_service_pb = require("../../proto/experiment/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ExperimentService = (function () {
  function ExperimentService() {}
  ExperimentService.serviceName = "bucketeer.experiment.ExperimentService";
  return ExperimentService;
}());

ExperimentService.GetGoal = {
  methodName: "GetGoal",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.GetGoalRequest,
  responseType: proto_experiment_service_pb.GetGoalResponse
};

ExperimentService.ListGoals = {
  methodName: "ListGoals",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.ListGoalsRequest,
  responseType: proto_experiment_service_pb.ListGoalsResponse
};

ExperimentService.CreateGoal = {
  methodName: "CreateGoal",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.CreateGoalRequest,
  responseType: proto_experiment_service_pb.CreateGoalResponse
};

ExperimentService.UpdateGoal = {
  methodName: "UpdateGoal",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.UpdateGoalRequest,
  responseType: proto_experiment_service_pb.UpdateGoalResponse
};

ExperimentService.ArchiveGoal = {
  methodName: "ArchiveGoal",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.ArchiveGoalRequest,
  responseType: proto_experiment_service_pb.ArchiveGoalResponse
};

ExperimentService.DeleteGoal = {
  methodName: "DeleteGoal",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.DeleteGoalRequest,
  responseType: proto_experiment_service_pb.DeleteGoalResponse
};

ExperimentService.GetExperiment = {
  methodName: "GetExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.GetExperimentRequest,
  responseType: proto_experiment_service_pb.GetExperimentResponse
};

ExperimentService.ListExperiments = {
  methodName: "ListExperiments",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.ListExperimentsRequest,
  responseType: proto_experiment_service_pb.ListExperimentsResponse
};

ExperimentService.CreateExperiment = {
  methodName: "CreateExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.CreateExperimentRequest,
  responseType: proto_experiment_service_pb.CreateExperimentResponse
};

ExperimentService.UpdateExperiment = {
  methodName: "UpdateExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.UpdateExperimentRequest,
  responseType: proto_experiment_service_pb.UpdateExperimentResponse
};

ExperimentService.StartExperiment = {
  methodName: "StartExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.StartExperimentRequest,
  responseType: proto_experiment_service_pb.StartExperimentResponse
};

ExperimentService.FinishExperiment = {
  methodName: "FinishExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.FinishExperimentRequest,
  responseType: proto_experiment_service_pb.FinishExperimentResponse
};

ExperimentService.StopExperiment = {
  methodName: "StopExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.StopExperimentRequest,
  responseType: proto_experiment_service_pb.StopExperimentResponse
};

ExperimentService.ArchiveExperiment = {
  methodName: "ArchiveExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.ArchiveExperimentRequest,
  responseType: proto_experiment_service_pb.ArchiveExperimentResponse
};

ExperimentService.DeleteExperiment = {
  methodName: "DeleteExperiment",
  service: ExperimentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experiment_service_pb.DeleteExperimentRequest,
  responseType: proto_experiment_service_pb.DeleteExperimentResponse
};

exports.ExperimentService = ExperimentService;

function ExperimentServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ExperimentServiceClient.prototype.getGoal = function getGoal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.GetGoal, {
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

ExperimentServiceClient.prototype.listGoals = function listGoals(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.ListGoals, {
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

ExperimentServiceClient.prototype.createGoal = function createGoal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.CreateGoal, {
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

ExperimentServiceClient.prototype.updateGoal = function updateGoal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.UpdateGoal, {
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

ExperimentServiceClient.prototype.archiveGoal = function archiveGoal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.ArchiveGoal, {
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

ExperimentServiceClient.prototype.deleteGoal = function deleteGoal(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.DeleteGoal, {
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

ExperimentServiceClient.prototype.getExperiment = function getExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.GetExperiment, {
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

ExperimentServiceClient.prototype.listExperiments = function listExperiments(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.ListExperiments, {
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

ExperimentServiceClient.prototype.createExperiment = function createExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.CreateExperiment, {
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

ExperimentServiceClient.prototype.updateExperiment = function updateExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.UpdateExperiment, {
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

ExperimentServiceClient.prototype.startExperiment = function startExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.StartExperiment, {
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

ExperimentServiceClient.prototype.finishExperiment = function finishExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.FinishExperiment, {
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

ExperimentServiceClient.prototype.stopExperiment = function stopExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.StopExperiment, {
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

ExperimentServiceClient.prototype.archiveExperiment = function archiveExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.ArchiveExperiment, {
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

ExperimentServiceClient.prototype.deleteExperiment = function deleteExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentService.DeleteExperiment, {
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

exports.ExperimentServiceClient = ExperimentServiceClient;

