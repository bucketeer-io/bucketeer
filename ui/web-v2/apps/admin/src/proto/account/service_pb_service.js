// package: bucketeer.account
// file: proto/account/service.proto

var proto_account_service_pb = require("../../proto/account/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var AccountService = (function () {
  function AccountService() {}
  AccountService.serviceName = "bucketeer.account.AccountService";
  return AccountService;
}());

AccountService.GetMe = {
  methodName: "GetMe",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetMeRequest,
  responseType: proto_account_service_pb.GetMeResponse
};

AccountService.GetMeByEmail = {
  methodName: "GetMeByEmail",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetMeByEmailRequest,
  responseType: proto_account_service_pb.GetMeResponse
};

AccountService.CreateAdminAccount = {
  methodName: "CreateAdminAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.CreateAdminAccountRequest,
  responseType: proto_account_service_pb.CreateAdminAccountResponse
};

AccountService.EnableAdminAccount = {
  methodName: "EnableAdminAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.EnableAdminAccountRequest,
  responseType: proto_account_service_pb.EnableAdminAccountResponse
};

AccountService.DisableAdminAccount = {
  methodName: "DisableAdminAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DisableAdminAccountRequest,
  responseType: proto_account_service_pb.DisableAdminAccountResponse
};

AccountService.GetAdminAccount = {
  methodName: "GetAdminAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAdminAccountRequest,
  responseType: proto_account_service_pb.GetAdminAccountResponse
};

AccountService.ListAdminAccounts = {
  methodName: "ListAdminAccounts",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ListAdminAccountsRequest,
  responseType: proto_account_service_pb.ListAdminAccountsResponse
};

AccountService.ConvertAccount = {
  methodName: "ConvertAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ConvertAccountRequest,
  responseType: proto_account_service_pb.ConvertAccountResponse
};

AccountService.CreateAccount = {
  methodName: "CreateAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.CreateAccountRequest,
  responseType: proto_account_service_pb.CreateAccountResponse
};

AccountService.EnableAccount = {
  methodName: "EnableAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.EnableAccountRequest,
  responseType: proto_account_service_pb.EnableAccountResponse
};

AccountService.DisableAccount = {
  methodName: "DisableAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DisableAccountRequest,
  responseType: proto_account_service_pb.DisableAccountResponse
};

AccountService.ChangeAccountRole = {
  methodName: "ChangeAccountRole",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ChangeAccountRoleRequest,
  responseType: proto_account_service_pb.ChangeAccountRoleResponse
};

AccountService.GetAccount = {
  methodName: "GetAccount",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAccountRequest,
  responseType: proto_account_service_pb.GetAccountResponse
};

AccountService.ListAccounts = {
  methodName: "ListAccounts",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ListAccountsRequest,
  responseType: proto_account_service_pb.ListAccountsResponse
};

AccountService.CreateAPIKey = {
  methodName: "CreateAPIKey",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.CreateAPIKeyRequest,
  responseType: proto_account_service_pb.CreateAPIKeyResponse
};

AccountService.ChangeAPIKeyName = {
  methodName: "ChangeAPIKeyName",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ChangeAPIKeyNameRequest,
  responseType: proto_account_service_pb.ChangeAPIKeyNameResponse
};

AccountService.EnableAPIKey = {
  methodName: "EnableAPIKey",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.EnableAPIKeyRequest,
  responseType: proto_account_service_pb.EnableAPIKeyResponse
};

AccountService.DisableAPIKey = {
  methodName: "DisableAPIKey",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DisableAPIKeyRequest,
  responseType: proto_account_service_pb.DisableAPIKeyResponse
};

AccountService.GetAPIKey = {
  methodName: "GetAPIKey",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAPIKeyRequest,
  responseType: proto_account_service_pb.GetAPIKeyResponse
};

AccountService.ListAPIKeys = {
  methodName: "ListAPIKeys",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ListAPIKeysRequest,
  responseType: proto_account_service_pb.ListAPIKeysResponse
};

AccountService.GetAPIKeyBySearchingAllEnvironments = {
  methodName: "GetAPIKeyBySearchingAllEnvironments",
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsRequest,
  responseType: proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsResponse
};

exports.AccountService = AccountService;

function AccountServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

AccountServiceClient.prototype.getMe = function getMe(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetMe, {
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

AccountServiceClient.prototype.getMeByEmail = function getMeByEmail(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetMeByEmail, {
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

AccountServiceClient.prototype.createAdminAccount = function createAdminAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.CreateAdminAccount, {
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

AccountServiceClient.prototype.enableAdminAccount = function enableAdminAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.EnableAdminAccount, {
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

AccountServiceClient.prototype.disableAdminAccount = function disableAdminAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.DisableAdminAccount, {
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

AccountServiceClient.prototype.getAdminAccount = function getAdminAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetAdminAccount, {
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

AccountServiceClient.prototype.listAdminAccounts = function listAdminAccounts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ListAdminAccounts, {
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

AccountServiceClient.prototype.convertAccount = function convertAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ConvertAccount, {
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

AccountServiceClient.prototype.createAccount = function createAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.CreateAccount, {
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

AccountServiceClient.prototype.enableAccount = function enableAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.EnableAccount, {
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

AccountServiceClient.prototype.disableAccount = function disableAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.DisableAccount, {
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

AccountServiceClient.prototype.changeAccountRole = function changeAccountRole(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ChangeAccountRole, {
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

AccountServiceClient.prototype.getAccount = function getAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetAccount, {
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

AccountServiceClient.prototype.listAccounts = function listAccounts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ListAccounts, {
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

AccountServiceClient.prototype.createAPIKey = function createAPIKey(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.CreateAPIKey, {
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

AccountServiceClient.prototype.changeAPIKeyName = function changeAPIKeyName(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ChangeAPIKeyName, {
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

AccountServiceClient.prototype.enableAPIKey = function enableAPIKey(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.EnableAPIKey, {
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

AccountServiceClient.prototype.disableAPIKey = function disableAPIKey(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.DisableAPIKey, {
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

AccountServiceClient.prototype.getAPIKey = function getAPIKey(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetAPIKey, {
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

AccountServiceClient.prototype.listAPIKeys = function listAPIKeys(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ListAPIKeys, {
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

AccountServiceClient.prototype.getAPIKeyBySearchingAllEnvironments = function getAPIKeyBySearchingAllEnvironments(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetAPIKeyBySearchingAllEnvironments, {
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

exports.AccountServiceClient = AccountServiceClient;

