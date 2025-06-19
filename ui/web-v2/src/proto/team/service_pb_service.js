// package: bucketeer.team
// file: proto/team/service.proto

var proto_team_service_pb = require('../../proto/team/service_pb');
var grpc = require('@improbable-eng/grpc-web').grpc;

var TeamService = (function () {
  function TeamService() {}
  TeamService.serviceName = 'bucketeer.team.TeamService';
  return TeamService;
})();

TeamService.CreateTeam = {
  methodName: 'CreateTeam',
  service: TeamService,
  requestStream: false,
  responseStream: false,
  requestType: proto_team_service_pb.CreateTeamRequest,
  responseType: proto_team_service_pb.CreateTeamResponse
};

TeamService.DeleteTeam = {
  methodName: 'DeleteTeam',
  service: TeamService,
  requestStream: false,
  responseStream: false,
  requestType: proto_team_service_pb.DeleteTeamRequest,
  responseType: proto_team_service_pb.DeleteTeamResponse
};

TeamService.ListTeams = {
  methodName: 'ListTeams',
  service: TeamService,
  requestStream: false,
  responseStream: false,
  requestType: proto_team_service_pb.ListTeamsRequest,
  responseType: proto_team_service_pb.ListTeamsResponse
};

exports.TeamService = TeamService;

function TeamServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

TeamServiceClient.prototype.createTeam = function createTeam(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TeamService.CreateTeam, {
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

TeamServiceClient.prototype.deleteTeam = function deleteTeam(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TeamService.DeleteTeam, {
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

TeamServiceClient.prototype.listTeams = function listTeams(
  requestMessage,
  metadata,
  callback
) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(TeamService.ListTeams, {
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

exports.TeamServiceClient = TeamServiceClient;
