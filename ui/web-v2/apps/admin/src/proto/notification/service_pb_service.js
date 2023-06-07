// package: bucketeer.notification
// file: proto/notification/service.proto

var proto_notification_service_pb = require("../../proto/notification/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var NotificationService = (function () {
  function NotificationService() {}
  NotificationService.serviceName = "bucketeer.notification.NotificationService";
  return NotificationService;
}());

NotificationService.GetAdminSubscription = {
  methodName: "GetAdminSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.GetAdminSubscriptionRequest,
  responseType: proto_notification_service_pb.GetAdminSubscriptionResponse
};

NotificationService.ListAdminSubscriptions = {
  methodName: "ListAdminSubscriptions",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.ListAdminSubscriptionsRequest,
  responseType: proto_notification_service_pb.ListAdminSubscriptionsResponse
};

NotificationService.ListEnabledAdminSubscriptions = {
  methodName: "ListEnabledAdminSubscriptions",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.ListEnabledAdminSubscriptionsRequest,
  responseType: proto_notification_service_pb.ListEnabledAdminSubscriptionsResponse
};

NotificationService.CreateAdminSubscription = {
  methodName: "CreateAdminSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.CreateAdminSubscriptionRequest,
  responseType: proto_notification_service_pb.CreateAdminSubscriptionResponse
};

NotificationService.DeleteAdminSubscription = {
  methodName: "DeleteAdminSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.DeleteAdminSubscriptionRequest,
  responseType: proto_notification_service_pb.DeleteAdminSubscriptionResponse
};

NotificationService.EnableAdminSubscription = {
  methodName: "EnableAdminSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.EnableAdminSubscriptionRequest,
  responseType: proto_notification_service_pb.EnableAdminSubscriptionResponse
};

NotificationService.DisableAdminSubscription = {
  methodName: "DisableAdminSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.DisableAdminSubscriptionRequest,
  responseType: proto_notification_service_pb.DisableAdminSubscriptionResponse
};

NotificationService.UpdateAdminSubscription = {
  methodName: "UpdateAdminSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.UpdateAdminSubscriptionRequest,
  responseType: proto_notification_service_pb.UpdateAdminSubscriptionResponse
};

NotificationService.GetSubscription = {
  methodName: "GetSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.GetSubscriptionRequest,
  responseType: proto_notification_service_pb.GetSubscriptionResponse
};

NotificationService.ListSubscriptions = {
  methodName: "ListSubscriptions",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.ListSubscriptionsRequest,
  responseType: proto_notification_service_pb.ListSubscriptionsResponse
};

NotificationService.ListEnabledSubscriptions = {
  methodName: "ListEnabledSubscriptions",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.ListEnabledSubscriptionsRequest,
  responseType: proto_notification_service_pb.ListEnabledSubscriptionsResponse
};

NotificationService.CreateSubscription = {
  methodName: "CreateSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.CreateSubscriptionRequest,
  responseType: proto_notification_service_pb.CreateSubscriptionResponse
};

NotificationService.DeleteSubscription = {
  methodName: "DeleteSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.DeleteSubscriptionRequest,
  responseType: proto_notification_service_pb.DeleteSubscriptionResponse
};

NotificationService.EnableSubscription = {
  methodName: "EnableSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.EnableSubscriptionRequest,
  responseType: proto_notification_service_pb.EnableSubscriptionResponse
};

NotificationService.DisableSubscription = {
  methodName: "DisableSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.DisableSubscriptionRequest,
  responseType: proto_notification_service_pb.DisableSubscriptionResponse
};

NotificationService.UpdateSubscription = {
  methodName: "UpdateSubscription",
  service: NotificationService,
  requestStream: false,
  responseStream: false,
  requestType: proto_notification_service_pb.UpdateSubscriptionRequest,
  responseType: proto_notification_service_pb.UpdateSubscriptionResponse
};

exports.NotificationService = NotificationService;

function NotificationServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

NotificationServiceClient.prototype.getAdminSubscription = function getAdminSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.GetAdminSubscription, {
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

NotificationServiceClient.prototype.listAdminSubscriptions = function listAdminSubscriptions(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.ListAdminSubscriptions, {
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

NotificationServiceClient.prototype.listEnabledAdminSubscriptions = function listEnabledAdminSubscriptions(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.ListEnabledAdminSubscriptions, {
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

NotificationServiceClient.prototype.createAdminSubscription = function createAdminSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.CreateAdminSubscription, {
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

NotificationServiceClient.prototype.deleteAdminSubscription = function deleteAdminSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.DeleteAdminSubscription, {
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

NotificationServiceClient.prototype.enableAdminSubscription = function enableAdminSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.EnableAdminSubscription, {
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

NotificationServiceClient.prototype.disableAdminSubscription = function disableAdminSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.DisableAdminSubscription, {
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

NotificationServiceClient.prototype.updateAdminSubscription = function updateAdminSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.UpdateAdminSubscription, {
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

NotificationServiceClient.prototype.getSubscription = function getSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.GetSubscription, {
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

NotificationServiceClient.prototype.listSubscriptions = function listSubscriptions(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.ListSubscriptions, {
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

NotificationServiceClient.prototype.listEnabledSubscriptions = function listEnabledSubscriptions(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.ListEnabledSubscriptions, {
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

NotificationServiceClient.prototype.createSubscription = function createSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.CreateSubscription, {
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

NotificationServiceClient.prototype.deleteSubscription = function deleteSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.DeleteSubscription, {
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

NotificationServiceClient.prototype.enableSubscription = function enableSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.EnableSubscription, {
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

NotificationServiceClient.prototype.disableSubscription = function disableSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.DisableSubscription, {
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

NotificationServiceClient.prototype.updateSubscription = function updateSubscription(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(NotificationService.UpdateSubscription, {
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

exports.NotificationServiceClient = NotificationServiceClient;

