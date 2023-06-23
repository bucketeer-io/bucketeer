// package: bucketeer.autoops
// file: proto/autoops/service.proto

var proto_autoops_service_pb = require("../../proto/autoops/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var AutoOpsService = (function () {
  function AutoOpsService() {}
  AutoOpsService.serviceName = "bucketeer.autoops.AutoOpsService";
  return AutoOpsService;
}());

AutoOpsService.GetAutoOpsRule = {
  methodName: "GetAutoOpsRule",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.GetAutoOpsRuleRequest,
  responseType: proto_autoops_service_pb.GetAutoOpsRuleResponse
};

AutoOpsService.ListAutoOpsRules = {
  methodName: "ListAutoOpsRules",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.ListAutoOpsRulesRequest,
  responseType: proto_autoops_service_pb.ListAutoOpsRulesResponse
};

AutoOpsService.CreateAutoOpsRule = {
  methodName: "CreateAutoOpsRule",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.CreateAutoOpsRuleRequest,
  responseType: proto_autoops_service_pb.CreateAutoOpsRuleResponse
};

AutoOpsService.DeleteAutoOpsRule = {
  methodName: "DeleteAutoOpsRule",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.DeleteAutoOpsRuleRequest,
  responseType: proto_autoops_service_pb.DeleteAutoOpsRuleResponse
};

AutoOpsService.UpdateAutoOpsRule = {
  methodName: "UpdateAutoOpsRule",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.UpdateAutoOpsRuleRequest,
  responseType: proto_autoops_service_pb.UpdateAutoOpsRuleResponse
};

AutoOpsService.ExecuteAutoOps = {
  methodName: "ExecuteAutoOps",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.ExecuteAutoOpsRequest,
  responseType: proto_autoops_service_pb.ExecuteAutoOpsResponse
};

AutoOpsService.ListOpsCounts = {
  methodName: "ListOpsCounts",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.ListOpsCountsRequest,
  responseType: proto_autoops_service_pb.ListOpsCountsResponse
};

AutoOpsService.CreateWebhook = {
  methodName: "CreateWebhook",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.CreateWebhookRequest,
  responseType: proto_autoops_service_pb.CreateWebhookResponse
};

AutoOpsService.GetWebhook = {
  methodName: "GetWebhook",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.GetWebhookRequest,
  responseType: proto_autoops_service_pb.GetWebhookResponse
};

AutoOpsService.UpdateWebhook = {
  methodName: "UpdateWebhook",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.UpdateWebhookRequest,
  responseType: proto_autoops_service_pb.UpdateWebhookResponse
};

AutoOpsService.DeleteWebhook = {
  methodName: "DeleteWebhook",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.DeleteWebhookRequest,
  responseType: proto_autoops_service_pb.DeleteWebhookResponse
};

AutoOpsService.ListWebhooks = {
  methodName: "ListWebhooks",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.ListWebhooksRequest,
  responseType: proto_autoops_service_pb.ListWebhooksResponse
};

AutoOpsService.CreateProgressiveRollout = {
  methodName: "CreateProgressiveRollout",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.CreateProgressiveRolloutRequest,
  responseType: proto_autoops_service_pb.CreateProgressiveRolloutResponse
};

AutoOpsService.GetProgressiveRollout = {
  methodName: "GetProgressiveRollout",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.GetProgressiveRolloutRequest,
  responseType: proto_autoops_service_pb.GetProgressiveRolloutResponse
};

AutoOpsService.DeleteProgressiveRollout = {
  methodName: "DeleteProgressiveRollout",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.DeleteProgressiveRolloutRequest,
  responseType: proto_autoops_service_pb.DeleteProgressiveRolloutResponse
};

AutoOpsService.ListProgressiveRollouts = {
  methodName: "ListProgressiveRollouts",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.ListProgressiveRolloutsRequest,
  responseType: proto_autoops_service_pb.ListProgressiveRolloutsResponse
};

AutoOpsService.ExecuteProgressiveRollout = {
  methodName: "ExecuteProgressiveRollout",
  service: AutoOpsService,
  requestStream: false,
  responseStream: false,
  requestType: proto_autoops_service_pb.ExecuteProgressiveRolloutRequest,
  responseType: proto_autoops_service_pb.ExecuteProgressiveRolloutResponse
};

exports.AutoOpsService = AutoOpsService;

function AutoOpsServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

AutoOpsServiceClient.prototype.getAutoOpsRule = function getAutoOpsRule(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.GetAutoOpsRule, {
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

AutoOpsServiceClient.prototype.listAutoOpsRules = function listAutoOpsRules(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.ListAutoOpsRules, {
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

AutoOpsServiceClient.prototype.createAutoOpsRule = function createAutoOpsRule(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.CreateAutoOpsRule, {
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

AutoOpsServiceClient.prototype.deleteAutoOpsRule = function deleteAutoOpsRule(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.DeleteAutoOpsRule, {
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

AutoOpsServiceClient.prototype.updateAutoOpsRule = function updateAutoOpsRule(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.UpdateAutoOpsRule, {
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

AutoOpsServiceClient.prototype.executeAutoOps = function executeAutoOps(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.ExecuteAutoOps, {
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

AutoOpsServiceClient.prototype.listOpsCounts = function listOpsCounts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.ListOpsCounts, {
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

AutoOpsServiceClient.prototype.createWebhook = function createWebhook(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.CreateWebhook, {
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

AutoOpsServiceClient.prototype.getWebhook = function getWebhook(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.GetWebhook, {
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

AutoOpsServiceClient.prototype.updateWebhook = function updateWebhook(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.UpdateWebhook, {
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

AutoOpsServiceClient.prototype.deleteWebhook = function deleteWebhook(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.DeleteWebhook, {
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

AutoOpsServiceClient.prototype.listWebhooks = function listWebhooks(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.ListWebhooks, {
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

AutoOpsServiceClient.prototype.createProgressiveRollout = function createProgressiveRollout(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.CreateProgressiveRollout, {
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

AutoOpsServiceClient.prototype.getProgressiveRollout = function getProgressiveRollout(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.GetProgressiveRollout, {
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

AutoOpsServiceClient.prototype.deleteProgressiveRollout = function deleteProgressiveRollout(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.DeleteProgressiveRollout, {
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

AutoOpsServiceClient.prototype.listProgressiveRollouts = function listProgressiveRollouts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.ListProgressiveRollouts, {
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

AutoOpsServiceClient.prototype.executeProgressiveRollout = function executeProgressiveRollout(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AutoOpsService.ExecuteProgressiveRollout, {
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

exports.AutoOpsServiceClient = AutoOpsServiceClient;

