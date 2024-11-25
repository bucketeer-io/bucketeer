// package: bucketeer.account
// file: proto/account/service.proto

var proto_account_service_pb = require('../../proto/account/service_pb');
var grpc = require('@improbable-eng/grpc-web').grpc;

var AccountService = (function () {
  function AccountService() {}
  AccountService.serviceName = 'bucketeer.account.AccountService';
  return AccountService;
})();

AccountService.GetMe = {
  methodName: 'GetMe',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetMeRequest,
  responseType: proto_account_service_pb.GetMeResponse
};

AccountService.GetMyOrganizations = {
  methodName: 'GetMyOrganizations',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetMyOrganizationsRequest,
  responseType: proto_account_service_pb.GetMyOrganizationsResponse
};

AccountService.GetMyOrganizationsByEmail = {
  methodName: 'GetMyOrganizationsByEmail',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetMyOrganizationsByEmailRequest,
  responseType: proto_account_service_pb.GetMyOrganizationsResponse
};

AccountService.CreateAccountV2 = {
  methodName: 'CreateAccountV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.CreateAccountV2Request,
  responseType: proto_account_service_pb.CreateAccountV2Response
};

AccountService.EnableAccountV2 = {
  methodName: 'EnableAccountV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.EnableAccountV2Request,
  responseType: proto_account_service_pb.EnableAccountV2Response
};

AccountService.DisableAccountV2 = {
  methodName: 'DisableAccountV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DisableAccountV2Request,
  responseType: proto_account_service_pb.DisableAccountV2Response
};

AccountService.UpdateAccountV2 = {
  methodName: 'UpdateAccountV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.UpdateAccountV2Request,
  responseType: proto_account_service_pb.UpdateAccountV2Response
};

AccountService.DeleteAccountV2 = {
  methodName: 'DeleteAccountV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DeleteAccountV2Request,
  responseType: proto_account_service_pb.DeleteAccountV2Response
};

AccountService.GetAccountV2 = {
  methodName: 'GetAccountV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAccountV2Request,
  responseType: proto_account_service_pb.GetAccountV2Response
};

AccountService.GetAccountV2ByEnvironmentID = {
  methodName: 'GetAccountV2ByEnvironmentID',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAccountV2ByEnvironmentIDRequest,
  responseType: proto_account_service_pb.GetAccountV2ByEnvironmentIDResponse
};

AccountService.ListAccountsV2 = {
  methodName: 'ListAccountsV2',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ListAccountsV2Request,
  responseType: proto_account_service_pb.ListAccountsV2Response
};

AccountService.CreateAPIKey = {
  methodName: 'CreateAPIKey',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.CreateAPIKeyRequest,
  responseType: proto_account_service_pb.CreateAPIKeyResponse
};

AccountService.ChangeAPIKeyName = {
  methodName: 'ChangeAPIKeyName',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ChangeAPIKeyNameRequest,
  responseType: proto_account_service_pb.ChangeAPIKeyNameResponse
};

AccountService.EnableAPIKey = {
  methodName: 'EnableAPIKey',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.EnableAPIKeyRequest,
  responseType: proto_account_service_pb.EnableAPIKeyResponse
};

AccountService.DisableAPIKey = {
  methodName: 'DisableAPIKey',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DisableAPIKeyRequest,
  responseType: proto_account_service_pb.DisableAPIKeyResponse
};

AccountService.GetAPIKey = {
  methodName: 'GetAPIKey',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.GetAPIKeyRequest,
  responseType: proto_account_service_pb.GetAPIKeyResponse
};

AccountService.ListAPIKeys = {
  methodName: 'ListAPIKeys',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.ListAPIKeysRequest,
  responseType: proto_account_service_pb.ListAPIKeysResponse
};

AccountService.GetAPIKeyBySearchingAllEnvironments = {
  methodName: 'GetAPIKeyBySearchingAllEnvironments',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType:
    proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsRequest,
  responseType:
    proto_account_service_pb.GetAPIKeyBySearchingAllEnvironmentsResponse
};

AccountService.CreateSearchFilter = {
  methodName: 'CreateSearchFilter',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.CreateSearchFilterRequest,
  responseType: proto_account_service_pb.CreateSearchFilterResponse
};

AccountService.UpdateSearchFilter = {
  methodName: 'UpdateSearchFilter',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.UpdateSearchFilterRequest,
  responseType: proto_account_service_pb.UpdateSearchFilterResponse
};

AccountService.DeleteSearchFilter = {
  methodName: 'DeleteSearchFilter',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.DeleteSearchFilterRequest,
  responseType: proto_account_service_pb.DeleteSearchFilterResponse
};

AccountService.UpdateAPIKey = {
  methodName: 'UpdateAPIKey',
  service: AccountService,
  requestStream: false,
  responseStream: false,
  requestType: proto_account_service_pb.UpdateAPIKeyRequest,
  responseType: proto_account_service_pb.UpdateAPIKeyResponse
};

exports.AccountService = AccountService;

function AccountServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

AccountServiceClient.prototype.getMe = function getMe(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.getMyOrganizations = function getMyOrganizations(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetMyOrganizations, {
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

AccountServiceClient.prototype.getMyOrganizationsByEmail =
  function getMyOrganizationsByEmail(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(AccountService.GetMyOrganizationsByEmail, {
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

AccountServiceClient.prototype.createAccountV2 = function createAccountV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.CreateAccountV2, {
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

AccountServiceClient.prototype.enableAccountV2 = function enableAccountV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.EnableAccountV2, {
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

AccountServiceClient.prototype.disableAccountV2 = function disableAccountV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.DisableAccountV2, {
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

AccountServiceClient.prototype.updateAccountV2 = function updateAccountV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.UpdateAccountV2, {
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

AccountServiceClient.prototype.deleteAccountV2 = function deleteAccountV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.DeleteAccountV2, {
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

AccountServiceClient.prototype.getAccountV2 = function getAccountV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.GetAccountV2, {
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

AccountServiceClient.prototype.getAccountV2ByEnvironmentID =
  function getAccountV2ByEnvironmentID(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(AccountService.GetAccountV2ByEnvironmentID, {
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

AccountServiceClient.prototype.listAccountsV2 = function listAccountsV2(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.ListAccountsV2, {
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

AccountServiceClient.prototype.createAPIKey = function createAPIKey(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.changeAPIKeyName = function changeAPIKeyName(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.enableAPIKey = function enableAPIKey(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.disableAPIKey = function disableAPIKey(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.getAPIKey = function getAPIKey(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.listAPIKeys = function listAPIKeys(
  requestMessage,
  metadata,
  callback
) {
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

AccountServiceClient.prototype.getAPIKeyBySearchingAllEnvironments =
  function getAPIKeyBySearchingAllEnvironments(
    requestMessage,
    metadata,
    callback
  ) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(
      AccountService.GetAPIKeyBySearchingAllEnvironments,
      {
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
      }
    );
    return {
      cancel: function () {
        callback = null;
        client.close();
      }
    };
  };

AccountServiceClient.prototype.createSearchFilter = function createSearchFilter(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.CreateSearchFilter, {
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

AccountServiceClient.prototype.updateSearchFilter = function updateSearchFilter(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.UpdateSearchFilter, {
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

AccountServiceClient.prototype.deleteSearchFilter = function deleteSearchFilter(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.DeleteSearchFilter, {
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

AccountServiceClient.prototype.updateAPIKey = function updateAPIKey(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AccountService.UpdateAPIKey, {
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
