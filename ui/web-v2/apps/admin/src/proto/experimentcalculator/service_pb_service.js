// package: bucketeer.experimentcalculator
// file: proto/experimentcalculator/service.proto

var proto_experimentcalculator_service_pb = require("../../proto/experimentcalculator/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ExperimentCalculatorService = (function () {
  function ExperimentCalculatorService() {}
  ExperimentCalculatorService.serviceName = "bucketeer.experimentcalculator.ExperimentCalculatorService";
  return ExperimentCalculatorService;
}());

ExperimentCalculatorService.CalcExperiment = {
  methodName: "CalcExperiment",
  service: ExperimentCalculatorService,
  requestStream: false,
  responseStream: false,
  requestType: proto_experimentcalculator_service_pb.BatchCalcRequest,
  responseType: proto_experimentcalculator_service_pb.BatchCalcResponse
};

exports.ExperimentCalculatorService = ExperimentCalculatorService;

function ExperimentCalculatorServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ExperimentCalculatorServiceClient.prototype.calcExperiment = function calcExperiment(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ExperimentCalculatorService.CalcExperiment, {
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

exports.ExperimentCalculatorServiceClient = ExperimentCalculatorServiceClient;

