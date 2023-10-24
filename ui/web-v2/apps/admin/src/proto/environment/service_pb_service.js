// package: bucketeer.environment
// file: proto/environment/service.proto

var proto_environment_service_pb = require("../../proto/environment/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var EnvironmentService = (function () {
  function EnvironmentService() {}
  EnvironmentService.serviceName = "bucketeer.environment.EnvironmentService";
  return EnvironmentService;
}());

EnvironmentService.GetEnvironmentV2 = {
  methodName: "GetEnvironmentV2",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.GetEnvironmentV2Request,
  responseType: proto_environment_service_pb.GetEnvironmentV2Response
};

EnvironmentService.ListEnvironmentsV2 = {
  methodName: "ListEnvironmentsV2",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.ListEnvironmentsV2Request,
  responseType: proto_environment_service_pb.ListEnvironmentsV2Response
};

EnvironmentService.CreateEnvironmentV2 = {
  methodName: "CreateEnvironmentV2",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.CreateEnvironmentV2Request,
  responseType: proto_environment_service_pb.CreateEnvironmentV2Response
};

EnvironmentService.UpdateEnvironmentV2 = {
  methodName: "UpdateEnvironmentV2",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.UpdateEnvironmentV2Request,
  responseType: proto_environment_service_pb.UpdateEnvironmentV2Response
};

EnvironmentService.ArchiveEnvironmentV2 = {
  methodName: "ArchiveEnvironmentV2",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.ArchiveEnvironmentV2Request,
  responseType: proto_environment_service_pb.ArchiveEnvironmentV2Response
};

EnvironmentService.UnarchiveEnvironmentV2 = {
  methodName: "UnarchiveEnvironmentV2",
  service: EnvironmentService,
  requestStream: false,
  responseStream: false,
  requestType: proto_environment_service_pb.UnarchiveEnvironmentV2Request,
  responseType: proto_environment_service_pb.UnarchiveEnvironmentV2Response
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

EnvironmentServiceClient.prototype.getEnvironmentV2 = function getEnvironmentV2(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.GetEnvironmentV2, {
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

EnvironmentServiceClient.prototype.listEnvironmentsV2 = function listEnvironmentsV2(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.ListEnvironmentsV2, {
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

EnvironmentServiceClient.prototype.createEnvironmentV2 = function createEnvironmentV2(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.CreateEnvironmentV2, {
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

EnvironmentServiceClient.prototype.updateEnvironmentV2 = function updateEnvironmentV2(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.UpdateEnvironmentV2, {
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

EnvironmentServiceClient.prototype.archiveEnvironmentV2 = function archiveEnvironmentV2(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.ArchiveEnvironmentV2, {
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

EnvironmentServiceClient.prototype.unarchiveEnvironmentV2 = function unarchiveEnvironmentV2(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(EnvironmentService.UnarchiveEnvironmentV2, {
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

