// package: bucketeer.auth
// file: proto/auth/service.proto

var proto_auth_service_pb = require('../../proto/auth/service_pb');
var grpc = require('@improbable-eng/grpc-web').grpc;

var AuthService = (function () {
  function AuthService() {}
  AuthService.serviceName = 'bucketeer.auth.AuthService';
  return AuthService;
})();

AuthService.ExchangeToken = {
  methodName: 'ExchangeToken',
  service: AuthService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auth_service_pb.ExchangeTokenRequest,
  responseType: proto_auth_service_pb.ExchangeTokenResponse
};

AuthService.GetAuthenticationURL = {
  methodName: 'GetAuthenticationURL',
  service: AuthService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auth_service_pb.GetAuthenticationURLRequest,
  responseType: proto_auth_service_pb.GetAuthenticationURLResponse
};

AuthService.RefreshToken = {
  methodName: 'RefreshToken',
  service: AuthService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auth_service_pb.RefreshTokenRequest,
  responseType: proto_auth_service_pb.RefreshTokenResponse
};

AuthService.SignIn = {
  methodName: 'SignIn',
  service: AuthService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auth_service_pb.SignInRequest,
  responseType: proto_auth_service_pb.SignInResponse
};

AuthService.SwitchOrganization = {
  methodName: 'SwitchOrganization',
  service: AuthService,
  requestStream: false,
  responseStream: false,
  requestType: proto_auth_service_pb.SwitchOrganizationRequest,
  responseType: proto_auth_service_pb.SwitchOrganizationResponse
};

exports.AuthService = AuthService;

function AuthServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

AuthServiceClient.prototype.exchangeToken = function exchangeToken(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthService.ExchangeToken, {
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

AuthServiceClient.prototype.getAuthenticationURL =
  function getAuthenticationURL(requestMessage, metadata, callback) {
    if (arguments.length === 2) {
      callback = arguments[1];
    }
    var client = grpc.unary(AuthService.GetAuthenticationURL, {
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

AuthServiceClient.prototype.refreshToken = function refreshToken(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthService.RefreshToken, {
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

AuthServiceClient.prototype.signIn = function signIn(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthService.SignIn, {
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

AuthServiceClient.prototype.switchOrganization = function switchOrganization(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthService.SwitchOrganization, {
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

exports.AuthServiceClient = AuthServiceClient;
