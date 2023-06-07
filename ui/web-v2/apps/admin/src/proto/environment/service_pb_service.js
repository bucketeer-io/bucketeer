// package: bucketeer.environment
// file: proto/environment/service.proto

var proto_environment_service_pb = require("../../proto/environment/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var EnvironmentService = (function () {
  function EnvironmentService() {}
  EnvironmentService.serviceName = "bucketeer.environment.EnvironmentService";
  return EnvironmentService;
}());

EnvironmentService.GetEnvironment = {
  methodName: "GetEnvironment",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.GetEnvironmentRequest,
  responseType: proto_environment_service_pb.GetEnvironmentResponse
};

EnvironmentService.GetEnvironmentByNamespace = {
  methodName: "GetEnvironmentByNamespace",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.GetEnvironmentByNamespaceRequest,
  responseType: proto_environment_service_pb.GetEnvironmentByNamespaceResponse
};

EnvironmentService.ListEnvironments = {
  methodName: "ListEnvironments",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.ListEnvironmentsRequest,
  responseType: proto_environment_service_pb.ListEnvironmentsResponse
};

EnvironmentService.CreateEnvironment = {
  methodName: "CreateEnvironment",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.CreateEnvironmentRequest,
  responseType: proto_environment_service_pb.CreateEnvironmentResponse
};

EnvironmentService.UpdateEnvironment = {
  methodName: "UpdateEnvironment",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.UpdateEnvironmentRequest,
  responseType: proto_environment_service_pb.UpdateEnvironmentResponse
};

EnvironmentService.DeleteEnvironment = {
  methodName: "DeleteEnvironment",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.DeleteEnvironmentRequest,
  responseType: proto_environment_service_pb.DeleteEnvironmentResponse
};

EnvironmentService.GetProject = {
  methodName: "GetProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.GetProjectRequest,
  responseType: proto_environment_service_pb.GetProjectResponse
};

EnvironmentService.ListProjects = {
  methodName: "ListProjects",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.ListProjectsRequest,
  responseType: proto_environment_service_pb.ListProjectsResponse
};

EnvironmentService.CreateProject = {
  methodName: "CreateProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.CreateProjectRequest,
  responseType: proto_environment_service_pb.CreateProjectResponse
};

EnvironmentService.CreateTrialProject = {
  methodName: "CreateTrialProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.CreateTrialProjectRequest,
  responseType: proto_environment_service_pb.CreateTrialProjectResponse
};

EnvironmentService.UpdateProject = {
  methodName: "UpdateProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.UpdateProjectRequest,
  responseType: proto_environment_service_pb.UpdateProjectResponse
};

EnvironmentService.EnableProject = {
  methodName: "EnableProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.EnableProjectRequest,
  responseType: proto_environment_service_pb.EnableProjectResponse
};

EnvironmentService.DisableProject = {
  methodName: "DisableProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.DisableProjectRequest,
  responseType: proto_environment_service_pb.DisableProjectResponse
};

EnvironmentService.ConvertTrialProject = {
  methodName: "ConvertTrialProject",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.ConvertTrialProjectRequest,
  responseType: proto_environment_service_pb.ConvertTrialProjectResponse
};

exports.EnvironmentService = EnvironmentService;

function EnvironmentServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

EnvironmentServiceClient.prototype.getEnvironment = function getEnvironment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.GetEnvironment, {
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

EnvironmentServiceClient.prototype.getEnvironmentByNamespace = function getEnvironmentByNamespace(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.GetEnvironmentByNamespace, {
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

EnvironmentServiceClient.prototype.listEnvironments = function listEnvironments(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.ListEnvironments, {
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

EnvironmentServiceClient.prototype.createEnvironment = function createEnvironment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.CreateEnvironment, {
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

EnvironmentServiceClient.prototype.updateEnvironment = function updateEnvironment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.UpdateEnvironment, {
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

EnvironmentServiceClient.prototype.deleteEnvironment = function deleteEnvironment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.DeleteEnvironment, {
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

EnvironmentServiceClient.prototype.getProject = function getProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.GetProject, {
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

EnvironmentServiceClient.prototype.listProjects = function listProjects(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.ListProjects, {
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

EnvironmentServiceClient.prototype.createProject = function createProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.CreateProject, {
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

EnvironmentServiceClient.prototype.createTrialProject = function createTrialProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.CreateTrialProject, {
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

EnvironmentServiceClient.prototype.updateProject = function updateProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.UpdateProject, {
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

EnvironmentServiceClient.prototype.enableProject = function enableProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.EnableProject, {
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

EnvironmentServiceClient.prototype.disableProject = function disableProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.DisableProject, {
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

EnvironmentServiceClient.prototype.convertTrialProject = function convertTrialProject(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.ConvertTrialProject, {
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

exports.EnvironmentServiceClient = EnvironmentServiceClient;

