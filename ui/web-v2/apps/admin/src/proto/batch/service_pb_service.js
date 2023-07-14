/*
 * Copyright 2023 The Bucketeer Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// package: bucketeer.batch
// file: proto/batch/service.proto

var proto_batch_service_pb = require("../../proto/batch/service_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var BatchService = (function () {
  function BatchService() {}
  BatchService.serviceName = "bucketeer.batch.BatchService";
  return BatchService;
}());

BatchService.ExecuteBatchJob = {
  methodName: "ExecuteBatchJob",
  service: BatchService,
  requestStream: false,
  responseStream: false,
  requestType: proto_batch_service_pb.BatchJobRequest,
  responseType: proto_batch_service_pb.BatchJobResponse
};

exports.BatchService = BatchService;

function BatchServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

BatchServiceClient.prototype.executeBatchJob = function executeBatchJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(BatchService.ExecuteBatchJob, {
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

exports.BatchServiceClient = BatchServiceClient;

