// source: proto/environment/service.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global =
    (typeof globalThis !== 'undefined' && globalThis) ||
    (typeof window !== 'undefined' && window) ||
    (typeof global !== 'undefined' && global) ||
    (typeof self !== 'undefined' && self) ||
    (function () { return this; }).call(null) ||
    Function('return this')();

var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.object.extend(proto, google_protobuf_wrappers_pb);
var proto_environment_environment_pb = require('../../proto/environment/environment_pb.js');
goog.object.extend(proto, proto_environment_environment_pb);
var proto_environment_project_pb = require('../../proto/environment/project_pb.js');
goog.object.extend(proto, proto_environment_project_pb);
var proto_environment_command_pb = require('../../proto/environment/command_pb.js');
goog.object.extend(proto, proto_environment_command_pb);
goog.exportSymbol('proto.bucketeer.environment.ArchiveEnvironmentV2Request', null, global);
goog.exportSymbol('proto.bucketeer.environment.ArchiveEnvironmentV2Response', null, global);
goog.exportSymbol('proto.bucketeer.environment.ConvertTrialProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.ConvertTrialProjectResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateEnvironmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateEnvironmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateEnvironmentV2Request', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateEnvironmentV2Response', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateProjectResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateTrialProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.CreateTrialProjectResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.DeleteEnvironmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.DeleteEnvironmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.DisableProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.DisableProjectResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.EnableProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.EnableProjectResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetEnvironmentByNamespaceRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetEnvironmentByNamespaceResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetEnvironmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetEnvironmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetEnvironmentV2Request', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetEnvironmentV2Response', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.GetProjectResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsRequest.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsRequest.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsV2Request', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsV2Request.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsV2Request.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListEnvironmentsV2Response', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListProjectsRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListProjectsRequest.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListProjectsRequest.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.environment.ListProjectsResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.UnarchiveEnvironmentV2Request', null, global);
goog.exportSymbol('proto.bucketeer.environment.UnarchiveEnvironmentV2Response', null, global);
goog.exportSymbol('proto.bucketeer.environment.UpdateEnvironmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.UpdateEnvironmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.environment.UpdateEnvironmentV2Request', null, global);
goog.exportSymbol('proto.bucketeer.environment.UpdateEnvironmentV2Response', null, global);
goog.exportSymbol('proto.bucketeer.environment.UpdateProjectRequest', null, global);
goog.exportSymbol('proto.bucketeer.environment.UpdateProjectResponse', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetEnvironmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetEnvironmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetEnvironmentRequest.displayName = 'proto.bucketeer.environment.GetEnvironmentRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetEnvironmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetEnvironmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetEnvironmentResponse.displayName = 'proto.bucketeer.environment.GetEnvironmentResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetEnvironmentByNamespaceRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.displayName = 'proto.bucketeer.environment.GetEnvironmentByNamespaceRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetEnvironmentByNamespaceResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.displayName = 'proto.bucketeer.environment.GetEnvironmentByNamespaceResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ListEnvironmentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ListEnvironmentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ListEnvironmentsRequest.displayName = 'proto.bucketeer.environment.ListEnvironmentsRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ListEnvironmentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.environment.ListEnvironmentsResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.environment.ListEnvironmentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ListEnvironmentsResponse.displayName = 'proto.bucketeer.environment.ListEnvironmentsResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateEnvironmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateEnvironmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateEnvironmentRequest.displayName = 'proto.bucketeer.environment.CreateEnvironmentRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateEnvironmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateEnvironmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateEnvironmentResponse.displayName = 'proto.bucketeer.environment.CreateEnvironmentResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UpdateEnvironmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UpdateEnvironmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UpdateEnvironmentRequest.displayName = 'proto.bucketeer.environment.UpdateEnvironmentRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UpdateEnvironmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UpdateEnvironmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UpdateEnvironmentResponse.displayName = 'proto.bucketeer.environment.UpdateEnvironmentResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.DeleteEnvironmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.DeleteEnvironmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.DeleteEnvironmentRequest.displayName = 'proto.bucketeer.environment.DeleteEnvironmentRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.DeleteEnvironmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.DeleteEnvironmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.DeleteEnvironmentResponse.displayName = 'proto.bucketeer.environment.DeleteEnvironmentResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetEnvironmentV2Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetEnvironmentV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetEnvironmentV2Request.displayName = 'proto.bucketeer.environment.GetEnvironmentV2Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetEnvironmentV2Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetEnvironmentV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetEnvironmentV2Response.displayName = 'proto.bucketeer.environment.GetEnvironmentV2Response';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ListEnvironmentsV2Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ListEnvironmentsV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ListEnvironmentsV2Request.displayName = 'proto.bucketeer.environment.ListEnvironmentsV2Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ListEnvironmentsV2Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.environment.ListEnvironmentsV2Response.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.environment.ListEnvironmentsV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ListEnvironmentsV2Response.displayName = 'proto.bucketeer.environment.ListEnvironmentsV2Response';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateEnvironmentV2Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateEnvironmentV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateEnvironmentV2Request.displayName = 'proto.bucketeer.environment.CreateEnvironmentV2Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateEnvironmentV2Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateEnvironmentV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateEnvironmentV2Response.displayName = 'proto.bucketeer.environment.CreateEnvironmentV2Response';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UpdateEnvironmentV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UpdateEnvironmentV2Request.displayName = 'proto.bucketeer.environment.UpdateEnvironmentV2Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UpdateEnvironmentV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UpdateEnvironmentV2Response.displayName = 'proto.bucketeer.environment.UpdateEnvironmentV2Response';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ArchiveEnvironmentV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ArchiveEnvironmentV2Request.displayName = 'proto.bucketeer.environment.ArchiveEnvironmentV2Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ArchiveEnvironmentV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ArchiveEnvironmentV2Response.displayName = 'proto.bucketeer.environment.ArchiveEnvironmentV2Response';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UnarchiveEnvironmentV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UnarchiveEnvironmentV2Request.displayName = 'proto.bucketeer.environment.UnarchiveEnvironmentV2Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UnarchiveEnvironmentV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UnarchiveEnvironmentV2Response.displayName = 'proto.bucketeer.environment.UnarchiveEnvironmentV2Response';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetProjectRequest.displayName = 'proto.bucketeer.environment.GetProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.GetProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.GetProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.GetProjectResponse.displayName = 'proto.bucketeer.environment.GetProjectResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ListProjectsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ListProjectsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ListProjectsRequest.displayName = 'proto.bucketeer.environment.ListProjectsRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ListProjectsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.environment.ListProjectsResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.environment.ListProjectsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ListProjectsResponse.displayName = 'proto.bucketeer.environment.ListProjectsResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateProjectRequest.displayName = 'proto.bucketeer.environment.CreateProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateProjectResponse.displayName = 'proto.bucketeer.environment.CreateProjectResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateTrialProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateTrialProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateTrialProjectRequest.displayName = 'proto.bucketeer.environment.CreateTrialProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.CreateTrialProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.CreateTrialProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.CreateTrialProjectResponse.displayName = 'proto.bucketeer.environment.CreateTrialProjectResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UpdateProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UpdateProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UpdateProjectRequest.displayName = 'proto.bucketeer.environment.UpdateProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.UpdateProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.UpdateProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.UpdateProjectResponse.displayName = 'proto.bucketeer.environment.UpdateProjectResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.EnableProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.EnableProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.EnableProjectRequest.displayName = 'proto.bucketeer.environment.EnableProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.EnableProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.EnableProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.EnableProjectResponse.displayName = 'proto.bucketeer.environment.EnableProjectResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.DisableProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.DisableProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.DisableProjectRequest.displayName = 'proto.bucketeer.environment.DisableProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.DisableProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.DisableProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.DisableProjectResponse.displayName = 'proto.bucketeer.environment.DisableProjectResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ConvertTrialProjectRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ConvertTrialProjectRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ConvertTrialProjectRequest.displayName = 'proto.bucketeer.environment.ConvertTrialProjectRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.ConvertTrialProjectResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.ConvertTrialProjectResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.ConvertTrialProjectResponse.displayName = 'proto.bucketeer.environment.ConvertTrialProjectResponse';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetEnvironmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetEnvironmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetEnvironmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetEnvironmentRequest}
 */
proto.bucketeer.environment.GetEnvironmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetEnvironmentRequest;
  return proto.bucketeer.environment.GetEnvironmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetEnvironmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetEnvironmentRequest}
 */
proto.bucketeer.environment.GetEnvironmentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetEnvironmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetEnvironmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetEnvironmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.GetEnvironmentRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.GetEnvironmentRequest} returns this
 */
proto.bucketeer.environment.GetEnvironmentRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetEnvironmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetEnvironmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetEnvironmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    environment: (f = msg.getEnvironment()) && proto_environment_environment_pb.Environment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetEnvironmentResponse}
 */
proto.bucketeer.environment.GetEnvironmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetEnvironmentResponse;
  return proto.bucketeer.environment.GetEnvironmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetEnvironmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetEnvironmentResponse}
 */
proto.bucketeer.environment.GetEnvironmentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_environment_pb.Environment;
      reader.readMessage(value,proto_environment_environment_pb.Environment.deserializeBinaryFromReader);
      msg.setEnvironment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetEnvironmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetEnvironmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetEnvironmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_environment_pb.Environment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Environment environment = 1;
 * @return {?proto.bucketeer.environment.Environment}
 */
proto.bucketeer.environment.GetEnvironmentResponse.prototype.getEnvironment = function() {
  return /** @type{?proto.bucketeer.environment.Environment} */ (
    jspb.Message.getWrapperField(this, proto_environment_environment_pb.Environment, 1));
};


/**
 * @param {?proto.bucketeer.environment.Environment|undefined} value
 * @return {!proto.bucketeer.environment.GetEnvironmentResponse} returns this
*/
proto.bucketeer.environment.GetEnvironmentResponse.prototype.setEnvironment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.GetEnvironmentResponse} returns this
 */
proto.bucketeer.environment.GetEnvironmentResponse.prototype.clearEnvironment = function() {
  return this.setEnvironment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.GetEnvironmentResponse.prototype.hasEnvironment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetEnvironmentByNamespaceRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    namespace: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceRequest}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetEnvironmentByNamespaceRequest;
  return proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetEnvironmentByNamespaceRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceRequest}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setNamespace(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetEnvironmentByNamespaceRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string namespace = 1;
 * @return {string}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.prototype.getNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceRequest} returns this
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceRequest.prototype.setNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    environment: (f = msg.getEnvironment()) && proto_environment_environment_pb.Environment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetEnvironmentByNamespaceResponse;
  return proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_environment_pb.Environment;
      reader.readMessage(value,proto_environment_environment_pb.Environment.deserializeBinaryFromReader);
      msg.setEnvironment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_environment_pb.Environment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Environment environment = 1;
 * @return {?proto.bucketeer.environment.Environment}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.prototype.getEnvironment = function() {
  return /** @type{?proto.bucketeer.environment.Environment} */ (
    jspb.Message.getWrapperField(this, proto_environment_environment_pb.Environment, 1));
};


/**
 * @param {?proto.bucketeer.environment.Environment|undefined} value
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse} returns this
*/
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.prototype.setEnvironment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.GetEnvironmentByNamespaceResponse} returns this
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.prototype.clearEnvironment = function() {
  return this.setEnvironment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.GetEnvironmentByNamespaceResponse.prototype.hasEnvironment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ListEnvironmentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ListEnvironmentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    projectId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    orderBy: jspb.Message.getFieldWithDefault(msg, 4, 0),
    orderDirection: jspb.Message.getFieldWithDefault(msg, 5, 0),
    searchKeyword: jspb.Message.getFieldWithDefault(msg, 6, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ListEnvironmentsRequest;
  return proto.bucketeer.environment.ListEnvironmentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ListEnvironmentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageSize(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 5:
      var value = /** @type {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderDirection} */ (reader.readEnum());
      msg.setOrderDirection(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setSearchKeyword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ListEnvironmentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ListEnvironmentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getOrderBy();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getOrderDirection();
  if (f !== 0.0) {
    writer.writeEnum(
      5,
      f
    );
  }
  f = message.getSearchKeyword();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.OrderBy = {
  DEFAULT: 0,
  ID: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3
};

/**
 * @enum {number}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest} returns this
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest} returns this
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string project_id = 3;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest} returns this
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderBy}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderBy} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest} returns this
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderDirection}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.bucketeer.environment.ListEnvironmentsRequest.OrderDirection} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest} returns this
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional string search_keyword = 6;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.getSearchKeyword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsRequest} returns this
 */
proto.bucketeer.environment.ListEnvironmentsRequest.prototype.setSearchKeyword = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.environment.ListEnvironmentsResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ListEnvironmentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ListEnvironmentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentsList: jspb.Message.toObjectList(msg.getEnvironmentsList(),
    proto_environment_environment_pb.Environment.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    totalCount: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ListEnvironmentsResponse}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ListEnvironmentsResponse;
  return proto.bucketeer.environment.ListEnvironmentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ListEnvironmentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ListEnvironmentsResponse}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_environment_pb.Environment;
      reader.readMessage(value,proto_environment_environment_pb.Environment.deserializeBinaryFromReader);
      msg.addEnvironments(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTotalCount(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ListEnvironmentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ListEnvironmentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_environment_environment_pb.Environment.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getTotalCount();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
};


/**
 * repeated Environment environments = 1;
 * @return {!Array<!proto.bucketeer.environment.Environment>}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.getEnvironmentsList = function() {
  return /** @type{!Array<!proto.bucketeer.environment.Environment>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_environment_environment_pb.Environment, 1));
};


/**
 * @param {!Array<!proto.bucketeer.environment.Environment>} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsResponse} returns this
*/
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.setEnvironmentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.environment.Environment=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.environment.Environment}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.addEnvironments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.environment.Environment, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.environment.ListEnvironmentsResponse} returns this
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.clearEnvironmentsList = function() {
  return this.setEnvironmentsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsResponse} returns this
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsResponse} returns this
 */
proto.bucketeer.environment.ListEnvironmentsResponse.prototype.setTotalCount = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateEnvironmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateEnvironmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateEnvironmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && proto_environment_command_pb.CreateEnvironmentCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateEnvironmentRequest}
 */
proto.bucketeer.environment.CreateEnvironmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateEnvironmentRequest;
  return proto.bucketeer.environment.CreateEnvironmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateEnvironmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateEnvironmentRequest}
 */
proto.bucketeer.environment.CreateEnvironmentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_command_pb.CreateEnvironmentCommand;
      reader.readMessage(value,proto_environment_command_pb.CreateEnvironmentCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateEnvironmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateEnvironmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateEnvironmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_command_pb.CreateEnvironmentCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional CreateEnvironmentCommand command = 1;
 * @return {?proto.bucketeer.environment.CreateEnvironmentCommand}
 */
proto.bucketeer.environment.CreateEnvironmentRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.CreateEnvironmentCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.CreateEnvironmentCommand, 1));
};


/**
 * @param {?proto.bucketeer.environment.CreateEnvironmentCommand|undefined} value
 * @return {!proto.bucketeer.environment.CreateEnvironmentRequest} returns this
*/
proto.bucketeer.environment.CreateEnvironmentRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.CreateEnvironmentRequest} returns this
 */
proto.bucketeer.environment.CreateEnvironmentRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.CreateEnvironmentRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateEnvironmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateEnvironmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateEnvironmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateEnvironmentResponse}
 */
proto.bucketeer.environment.CreateEnvironmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateEnvironmentResponse;
  return proto.bucketeer.environment.CreateEnvironmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateEnvironmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateEnvironmentResponse}
 */
proto.bucketeer.environment.CreateEnvironmentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateEnvironmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateEnvironmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateEnvironmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UpdateEnvironmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UpdateEnvironmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    renameCommand: (f = msg.getRenameCommand()) && proto_environment_command_pb.RenameEnvironmentCommand.toObject(includeInstance, f),
    changeDescriptionCommand: (f = msg.getChangeDescriptionCommand()) && proto_environment_command_pb.ChangeDescriptionEnvironmentCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UpdateEnvironmentRequest;
  return proto.bucketeer.environment.UpdateEnvironmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.RenameEnvironmentCommand;
      reader.readMessage(value,proto_environment_command_pb.RenameEnvironmentCommand.deserializeBinaryFromReader);
      msg.setRenameCommand(value);
      break;
    case 3:
      var value = new proto_environment_command_pb.ChangeDescriptionEnvironmentCommand;
      reader.readMessage(value,proto_environment_command_pb.ChangeDescriptionEnvironmentCommand.deserializeBinaryFromReader);
      msg.setChangeDescriptionCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UpdateEnvironmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRenameCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.RenameEnvironmentCommand.serializeBinaryToWriter
    );
  }
  f = message.getChangeDescriptionCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_environment_command_pb.ChangeDescriptionEnvironmentCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest} returns this
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional RenameEnvironmentCommand rename_command = 2;
 * @return {?proto.bucketeer.environment.RenameEnvironmentCommand}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.getRenameCommand = function() {
  return /** @type{?proto.bucketeer.environment.RenameEnvironmentCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.RenameEnvironmentCommand, 2));
};


/**
 * @param {?proto.bucketeer.environment.RenameEnvironmentCommand|undefined} value
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest} returns this
*/
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.setRenameCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest} returns this
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.clearRenameCommand = function() {
  return this.setRenameCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.hasRenameCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional ChangeDescriptionEnvironmentCommand change_description_command = 3;
 * @return {?proto.bucketeer.environment.ChangeDescriptionEnvironmentCommand}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.getChangeDescriptionCommand = function() {
  return /** @type{?proto.bucketeer.environment.ChangeDescriptionEnvironmentCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.ChangeDescriptionEnvironmentCommand, 3));
};


/**
 * @param {?proto.bucketeer.environment.ChangeDescriptionEnvironmentCommand|undefined} value
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest} returns this
*/
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.setChangeDescriptionCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentRequest} returns this
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.clearChangeDescriptionCommand = function() {
  return this.setChangeDescriptionCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UpdateEnvironmentRequest.prototype.hasChangeDescriptionCommand = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UpdateEnvironmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UpdateEnvironmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UpdateEnvironmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentResponse}
 */
proto.bucketeer.environment.UpdateEnvironmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UpdateEnvironmentResponse;
  return proto.bucketeer.environment.UpdateEnvironmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentResponse}
 */
proto.bucketeer.environment.UpdateEnvironmentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UpdateEnvironmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UpdateEnvironmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.DeleteEnvironmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.DeleteEnvironmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_environment_command_pb.DeleteEnvironmentCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.DeleteEnvironmentRequest}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.DeleteEnvironmentRequest;
  return proto.bucketeer.environment.DeleteEnvironmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.DeleteEnvironmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.DeleteEnvironmentRequest}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.DeleteEnvironmentCommand;
      reader.readMessage(value,proto_environment_command_pb.DeleteEnvironmentCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.DeleteEnvironmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.DeleteEnvironmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.DeleteEnvironmentCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.DeleteEnvironmentRequest} returns this
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DeleteEnvironmentCommand command = 2;
 * @return {?proto.bucketeer.environment.DeleteEnvironmentCommand}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.DeleteEnvironmentCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.DeleteEnvironmentCommand, 2));
};


/**
 * @param {?proto.bucketeer.environment.DeleteEnvironmentCommand|undefined} value
 * @return {!proto.bucketeer.environment.DeleteEnvironmentRequest} returns this
*/
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.DeleteEnvironmentRequest} returns this
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.DeleteEnvironmentRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.DeleteEnvironmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.DeleteEnvironmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.DeleteEnvironmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DeleteEnvironmentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.DeleteEnvironmentResponse}
 */
proto.bucketeer.environment.DeleteEnvironmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.DeleteEnvironmentResponse;
  return proto.bucketeer.environment.DeleteEnvironmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.DeleteEnvironmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.DeleteEnvironmentResponse}
 */
proto.bucketeer.environment.DeleteEnvironmentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.DeleteEnvironmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.DeleteEnvironmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.DeleteEnvironmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DeleteEnvironmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetEnvironmentV2Request.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetEnvironmentV2Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetEnvironmentV2Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentV2Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Request}
 */
proto.bucketeer.environment.GetEnvironmentV2Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetEnvironmentV2Request;
  return proto.bucketeer.environment.GetEnvironmentV2Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetEnvironmentV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Request}
 */
proto.bucketeer.environment.GetEnvironmentV2Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetEnvironmentV2Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetEnvironmentV2Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetEnvironmentV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentV2Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.GetEnvironmentV2Request.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.GetEnvironmentV2Request.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetEnvironmentV2Response.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetEnvironmentV2Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetEnvironmentV2Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentV2Response.toObject = function(includeInstance, msg) {
  var f, obj = {
    environment: (f = msg.getEnvironment()) && proto_environment_environment_pb.EnvironmentV2.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Response}
 */
proto.bucketeer.environment.GetEnvironmentV2Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetEnvironmentV2Response;
  return proto.bucketeer.environment.GetEnvironmentV2Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetEnvironmentV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Response}
 */
proto.bucketeer.environment.GetEnvironmentV2Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_environment_pb.EnvironmentV2;
      reader.readMessage(value,proto_environment_environment_pb.EnvironmentV2.deserializeBinaryFromReader);
      msg.setEnvironment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetEnvironmentV2Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetEnvironmentV2Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetEnvironmentV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetEnvironmentV2Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_environment_pb.EnvironmentV2.serializeBinaryToWriter
    );
  }
};


/**
 * optional EnvironmentV2 environment = 1;
 * @return {?proto.bucketeer.environment.EnvironmentV2}
 */
proto.bucketeer.environment.GetEnvironmentV2Response.prototype.getEnvironment = function() {
  return /** @type{?proto.bucketeer.environment.EnvironmentV2} */ (
    jspb.Message.getWrapperField(this, proto_environment_environment_pb.EnvironmentV2, 1));
};


/**
 * @param {?proto.bucketeer.environment.EnvironmentV2|undefined} value
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Response} returns this
*/
proto.bucketeer.environment.GetEnvironmentV2Response.prototype.setEnvironment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.GetEnvironmentV2Response} returns this
 */
proto.bucketeer.environment.GetEnvironmentV2Response.prototype.clearEnvironment = function() {
  return this.setEnvironment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.GetEnvironmentV2Response.prototype.hasEnvironment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ListEnvironmentsV2Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    orderBy: jspb.Message.getFieldWithDefault(msg, 3, 0),
    orderDirection: jspb.Message.getFieldWithDefault(msg, 4, 0),
    projectId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    archived: (f = msg.getArchived()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    searchKeyword: jspb.Message.getFieldWithDefault(msg, 7, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ListEnvironmentsV2Request;
  return proto.bucketeer.environment.ListEnvironmentsV2Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageSize(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderDirection} */ (reader.readEnum());
      msg.setOrderDirection(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
      break;
    case 6:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setArchived(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setSearchKeyword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ListEnvironmentsV2Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getOrderBy();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getOrderDirection();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getArchived();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getSearchKeyword();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.OrderBy = {
  DEFAULT: 0,
  ID: 1,
  NAME: 2,
  URL_CODE: 3,
  CREATED_AT: 4,
  UPDATED_AT: 5
};

/**
 * @enum {number}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional OrderBy order_by = 3;
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderBy}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderBy} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional OrderDirection order_direction = 4;
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderDirection}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Request.OrderDirection} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional string project_id = 5;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional google.protobuf.BoolValue archived = 6;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getArchived = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 6));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
*/
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setArchived = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.clearArchived = function() {
  return this.setArchived(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.hasArchived = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional string search_keyword = 7;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.getSearchKeyword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Request} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Request.prototype.setSearchKeyword = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ListEnvironmentsV2Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentsList: jspb.Message.toObjectList(msg.getEnvironmentsList(),
    proto_environment_environment_pb.EnvironmentV2.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    totalCount: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Response}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ListEnvironmentsV2Response;
  return proto.bucketeer.environment.ListEnvironmentsV2Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Response}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_environment_pb.EnvironmentV2;
      reader.readMessage(value,proto_environment_environment_pb.EnvironmentV2.deserializeBinaryFromReader);
      msg.addEnvironments(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTotalCount(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ListEnvironmentsV2Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ListEnvironmentsV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_environment_environment_pb.EnvironmentV2.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getTotalCount();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
};


/**
 * repeated EnvironmentV2 environments = 1;
 * @return {!Array<!proto.bucketeer.environment.EnvironmentV2>}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.getEnvironmentsList = function() {
  return /** @type{!Array<!proto.bucketeer.environment.EnvironmentV2>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_environment_environment_pb.EnvironmentV2, 1));
};


/**
 * @param {!Array<!proto.bucketeer.environment.EnvironmentV2>} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Response} returns this
*/
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.setEnvironmentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.environment.EnvironmentV2=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.environment.EnvironmentV2}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.addEnvironments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.environment.EnvironmentV2, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Response} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.clearEnvironmentsList = function() {
  return this.setEnvironmentsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Response} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.ListEnvironmentsV2Response} returns this
 */
proto.bucketeer.environment.ListEnvironmentsV2Response.prototype.setTotalCount = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateEnvironmentV2Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateEnvironmentV2Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && proto_environment_command_pb.CreateEnvironmentV2Command.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateEnvironmentV2Request}
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateEnvironmentV2Request;
  return proto.bucketeer.environment.CreateEnvironmentV2Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateEnvironmentV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateEnvironmentV2Request}
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_command_pb.CreateEnvironmentV2Command;
      reader.readMessage(value,proto_environment_command_pb.CreateEnvironmentV2Command.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateEnvironmentV2Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateEnvironmentV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_command_pb.CreateEnvironmentV2Command.serializeBinaryToWriter
    );
  }
};


/**
 * optional CreateEnvironmentV2Command command = 1;
 * @return {?proto.bucketeer.environment.CreateEnvironmentV2Command}
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.CreateEnvironmentV2Command} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.CreateEnvironmentV2Command, 1));
};


/**
 * @param {?proto.bucketeer.environment.CreateEnvironmentV2Command|undefined} value
 * @return {!proto.bucketeer.environment.CreateEnvironmentV2Request} returns this
*/
proto.bucketeer.environment.CreateEnvironmentV2Request.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.CreateEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.CreateEnvironmentV2Request.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateEnvironmentV2Response.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateEnvironmentV2Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateEnvironmentV2Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentV2Response.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateEnvironmentV2Response}
 */
proto.bucketeer.environment.CreateEnvironmentV2Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateEnvironmentV2Response;
  return proto.bucketeer.environment.CreateEnvironmentV2Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateEnvironmentV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateEnvironmentV2Response}
 */
proto.bucketeer.environment.CreateEnvironmentV2Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateEnvironmentV2Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateEnvironmentV2Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateEnvironmentV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateEnvironmentV2Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UpdateEnvironmentV2Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UpdateEnvironmentV2Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    renameCommand: (f = msg.getRenameCommand()) && proto_environment_command_pb.RenameEnvironmentV2Command.toObject(includeInstance, f),
    changeDescriptionCommand: (f = msg.getChangeDescriptionCommand()) && proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UpdateEnvironmentV2Request;
  return proto.bucketeer.environment.UpdateEnvironmentV2Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.RenameEnvironmentV2Command;
      reader.readMessage(value,proto_environment_command_pb.RenameEnvironmentV2Command.deserializeBinaryFromReader);
      msg.setRenameCommand(value);
      break;
    case 3:
      var value = new proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command;
      reader.readMessage(value,proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command.deserializeBinaryFromReader);
      msg.setChangeDescriptionCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UpdateEnvironmentV2Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRenameCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.RenameEnvironmentV2Command.serializeBinaryToWriter
    );
  }
  f = message.getChangeDescriptionCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional RenameEnvironmentV2Command rename_command = 2;
 * @return {?proto.bucketeer.environment.RenameEnvironmentV2Command}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.getRenameCommand = function() {
  return /** @type{?proto.bucketeer.environment.RenameEnvironmentV2Command} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.RenameEnvironmentV2Command, 2));
};


/**
 * @param {?proto.bucketeer.environment.RenameEnvironmentV2Command|undefined} value
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request} returns this
*/
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.setRenameCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.clearRenameCommand = function() {
  return this.setRenameCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.hasRenameCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional ChangeDescriptionEnvironmentV2Command change_description_command = 3;
 * @return {?proto.bucketeer.environment.ChangeDescriptionEnvironmentV2Command}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.getChangeDescriptionCommand = function() {
  return /** @type{?proto.bucketeer.environment.ChangeDescriptionEnvironmentV2Command} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.ChangeDescriptionEnvironmentV2Command, 3));
};


/**
 * @param {?proto.bucketeer.environment.ChangeDescriptionEnvironmentV2Command|undefined} value
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request} returns this
*/
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.setChangeDescriptionCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.clearChangeDescriptionCommand = function() {
  return this.setChangeDescriptionCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Request.prototype.hasChangeDescriptionCommand = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UpdateEnvironmentV2Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UpdateEnvironmentV2Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Response}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UpdateEnvironmentV2Response;
  return proto.bucketeer.environment.UpdateEnvironmentV2Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UpdateEnvironmentV2Response}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UpdateEnvironmentV2Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UpdateEnvironmentV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateEnvironmentV2Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ArchiveEnvironmentV2Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ArchiveEnvironmentV2Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_environment_command_pb.ArchiveEnvironmentV2Command.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Request}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ArchiveEnvironmentV2Request;
  return proto.bucketeer.environment.ArchiveEnvironmentV2Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ArchiveEnvironmentV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Request}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.ArchiveEnvironmentV2Command;
      reader.readMessage(value,proto_environment_command_pb.ArchiveEnvironmentV2Command.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ArchiveEnvironmentV2Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ArchiveEnvironmentV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.ArchiveEnvironmentV2Command.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ArchiveEnvironmentV2Command command = 2;
 * @return {?proto.bucketeer.environment.ArchiveEnvironmentV2Command}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.ArchiveEnvironmentV2Command} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.ArchiveEnvironmentV2Command, 2));
};


/**
 * @param {?proto.bucketeer.environment.ArchiveEnvironmentV2Command|undefined} value
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Request} returns this
*/
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Request.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ArchiveEnvironmentV2Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ArchiveEnvironmentV2Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Response}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ArchiveEnvironmentV2Response;
  return proto.bucketeer.environment.ArchiveEnvironmentV2Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ArchiveEnvironmentV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ArchiveEnvironmentV2Response}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ArchiveEnvironmentV2Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ArchiveEnvironmentV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ArchiveEnvironmentV2Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UnarchiveEnvironmentV2Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_environment_command_pb.UnarchiveEnvironmentV2Command.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UnarchiveEnvironmentV2Request;
  return proto.bucketeer.environment.UnarchiveEnvironmentV2Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.UnarchiveEnvironmentV2Command;
      reader.readMessage(value,proto_environment_command_pb.UnarchiveEnvironmentV2Command.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UnarchiveEnvironmentV2Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.UnarchiveEnvironmentV2Command.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional UnarchiveEnvironmentV2Command command = 2;
 * @return {?proto.bucketeer.environment.UnarchiveEnvironmentV2Command}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.UnarchiveEnvironmentV2Command} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.UnarchiveEnvironmentV2Command, 2));
};


/**
 * @param {?proto.bucketeer.environment.UnarchiveEnvironmentV2Command|undefined} value
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request} returns this
*/
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Request} returns this
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Request.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UnarchiveEnvironmentV2Response.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UnarchiveEnvironmentV2Response} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Response}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UnarchiveEnvironmentV2Response;
  return proto.bucketeer.environment.UnarchiveEnvironmentV2Response.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UnarchiveEnvironmentV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UnarchiveEnvironmentV2Response}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UnarchiveEnvironmentV2Response.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UnarchiveEnvironmentV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UnarchiveEnvironmentV2Response.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetProjectRequest}
 */
proto.bucketeer.environment.GetProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetProjectRequest;
  return proto.bucketeer.environment.GetProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetProjectRequest}
 */
proto.bucketeer.environment.GetProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.GetProjectRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.GetProjectRequest} returns this
 */
proto.bucketeer.environment.GetProjectRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.GetProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.GetProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.GetProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    project: (f = msg.getProject()) && proto_environment_project_pb.Project.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.GetProjectResponse}
 */
proto.bucketeer.environment.GetProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.GetProjectResponse;
  return proto.bucketeer.environment.GetProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.GetProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.GetProjectResponse}
 */
proto.bucketeer.environment.GetProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_project_pb.Project;
      reader.readMessage(value,proto_environment_project_pb.Project.deserializeBinaryFromReader);
      msg.setProject(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.GetProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.GetProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.GetProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.GetProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProject();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_project_pb.Project.serializeBinaryToWriter
    );
  }
};


/**
 * optional Project project = 1;
 * @return {?proto.bucketeer.environment.Project}
 */
proto.bucketeer.environment.GetProjectResponse.prototype.getProject = function() {
  return /** @type{?proto.bucketeer.environment.Project} */ (
    jspb.Message.getWrapperField(this, proto_environment_project_pb.Project, 1));
};


/**
 * @param {?proto.bucketeer.environment.Project|undefined} value
 * @return {!proto.bucketeer.environment.GetProjectResponse} returns this
*/
proto.bucketeer.environment.GetProjectResponse.prototype.setProject = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.GetProjectResponse} returns this
 */
proto.bucketeer.environment.GetProjectResponse.prototype.clearProject = function() {
  return this.setProject(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.GetProjectResponse.prototype.hasProject = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ListProjectsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ListProjectsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListProjectsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    orderBy: jspb.Message.getFieldWithDefault(msg, 3, 0),
    orderDirection: jspb.Message.getFieldWithDefault(msg, 4, 0),
    searchKeyword: jspb.Message.getFieldWithDefault(msg, 5, ""),
    disabled: (f = msg.getDisabled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ListProjectsRequest}
 */
proto.bucketeer.environment.ListProjectsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ListProjectsRequest;
  return proto.bucketeer.environment.ListProjectsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ListProjectsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ListProjectsRequest}
 */
proto.bucketeer.environment.ListProjectsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageSize(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.environment.ListProjectsRequest.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.environment.ListProjectsRequest.OrderDirection} */ (reader.readEnum());
      msg.setOrderDirection(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setSearchKeyword(value);
      break;
    case 6:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setDisabled(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ListProjectsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ListProjectsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListProjectsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getOrderBy();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getOrderDirection();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getSearchKeyword();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getDisabled();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.environment.ListProjectsRequest.OrderBy = {
  DEFAULT: 0,
  ID: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3,
  NAME: 4,
  URL_CODE: 5
};

/**
 * @enum {number}
 */
proto.bucketeer.environment.ListProjectsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional OrderBy order_by = 3;
 * @return {!proto.bucketeer.environment.ListProjectsRequest.OrderBy}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.environment.ListProjectsRequest.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.environment.ListProjectsRequest.OrderBy} value
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional OrderDirection order_direction = 4;
 * @return {!proto.bucketeer.environment.ListProjectsRequest.OrderDirection}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.environment.ListProjectsRequest.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.environment.ListProjectsRequest.OrderDirection} value
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional string search_keyword = 5;
 * @return {string}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.getSearchKeyword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.setSearchKeyword = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional google.protobuf.BoolValue disabled = 6;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.getDisabled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 6));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
*/
proto.bucketeer.environment.ListProjectsRequest.prototype.setDisabled = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.ListProjectsRequest} returns this
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.clearDisabled = function() {
  return this.setDisabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.ListProjectsRequest.prototype.hasDisabled = function() {
  return jspb.Message.getField(this, 6) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.environment.ListProjectsResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ListProjectsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ListProjectsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListProjectsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    projectsList: jspb.Message.toObjectList(msg.getProjectsList(),
    proto_environment_project_pb.Project.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    totalCount: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ListProjectsResponse}
 */
proto.bucketeer.environment.ListProjectsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ListProjectsResponse;
  return proto.bucketeer.environment.ListProjectsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ListProjectsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ListProjectsResponse}
 */
proto.bucketeer.environment.ListProjectsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_project_pb.Project;
      reader.readMessage(value,proto_environment_project_pb.Project.deserializeBinaryFromReader);
      msg.addProjects(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTotalCount(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ListProjectsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ListProjectsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ListProjectsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProjectsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_environment_project_pb.Project.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getTotalCount();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
};


/**
 * repeated Project projects = 1;
 * @return {!Array<!proto.bucketeer.environment.Project>}
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.getProjectsList = function() {
  return /** @type{!Array<!proto.bucketeer.environment.Project>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_environment_project_pb.Project, 1));
};


/**
 * @param {!Array<!proto.bucketeer.environment.Project>} value
 * @return {!proto.bucketeer.environment.ListProjectsResponse} returns this
*/
proto.bucketeer.environment.ListProjectsResponse.prototype.setProjectsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.environment.Project=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.environment.Project}
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.addProjects = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.environment.Project, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.environment.ListProjectsResponse} returns this
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.clearProjectsList = function() {
  return this.setProjectsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ListProjectsResponse} returns this
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.ListProjectsResponse} returns this
 */
proto.bucketeer.environment.ListProjectsResponse.prototype.setTotalCount = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && proto_environment_command_pb.CreateProjectCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateProjectRequest}
 */
proto.bucketeer.environment.CreateProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateProjectRequest;
  return proto.bucketeer.environment.CreateProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateProjectRequest}
 */
proto.bucketeer.environment.CreateProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_command_pb.CreateProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.CreateProjectCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_command_pb.CreateProjectCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional CreateProjectCommand command = 1;
 * @return {?proto.bucketeer.environment.CreateProjectCommand}
 */
proto.bucketeer.environment.CreateProjectRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.CreateProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.CreateProjectCommand, 1));
};


/**
 * @param {?proto.bucketeer.environment.CreateProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.CreateProjectRequest} returns this
*/
proto.bucketeer.environment.CreateProjectRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.CreateProjectRequest} returns this
 */
proto.bucketeer.environment.CreateProjectRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.CreateProjectRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateProjectResponse}
 */
proto.bucketeer.environment.CreateProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateProjectResponse;
  return proto.bucketeer.environment.CreateProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateProjectResponse}
 */
proto.bucketeer.environment.CreateProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateTrialProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateTrialProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateTrialProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateTrialProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && proto_environment_command_pb.CreateTrialProjectCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateTrialProjectRequest}
 */
proto.bucketeer.environment.CreateTrialProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateTrialProjectRequest;
  return proto.bucketeer.environment.CreateTrialProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateTrialProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateTrialProjectRequest}
 */
proto.bucketeer.environment.CreateTrialProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_environment_command_pb.CreateTrialProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.CreateTrialProjectCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateTrialProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateTrialProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateTrialProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateTrialProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_command_pb.CreateTrialProjectCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional CreateTrialProjectCommand command = 1;
 * @return {?proto.bucketeer.environment.CreateTrialProjectCommand}
 */
proto.bucketeer.environment.CreateTrialProjectRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.CreateTrialProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.CreateTrialProjectCommand, 1));
};


/**
 * @param {?proto.bucketeer.environment.CreateTrialProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.CreateTrialProjectRequest} returns this
*/
proto.bucketeer.environment.CreateTrialProjectRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.CreateTrialProjectRequest} returns this
 */
proto.bucketeer.environment.CreateTrialProjectRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.CreateTrialProjectRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.CreateTrialProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.CreateTrialProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.CreateTrialProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateTrialProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.CreateTrialProjectResponse}
 */
proto.bucketeer.environment.CreateTrialProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.CreateTrialProjectResponse;
  return proto.bucketeer.environment.CreateTrialProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.CreateTrialProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.CreateTrialProjectResponse}
 */
proto.bucketeer.environment.CreateTrialProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.CreateTrialProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.CreateTrialProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.CreateTrialProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.CreateTrialProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UpdateProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UpdateProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    changeDescriptionCommand: (f = msg.getChangeDescriptionCommand()) && proto_environment_command_pb.ChangeDescriptionProjectCommand.toObject(includeInstance, f),
    renameCommand: (f = msg.getRenameCommand()) && proto_environment_command_pb.RenameProjectCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UpdateProjectRequest}
 */
proto.bucketeer.environment.UpdateProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UpdateProjectRequest;
  return proto.bucketeer.environment.UpdateProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UpdateProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UpdateProjectRequest}
 */
proto.bucketeer.environment.UpdateProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.ChangeDescriptionProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.ChangeDescriptionProjectCommand.deserializeBinaryFromReader);
      msg.setChangeDescriptionCommand(value);
      break;
    case 3:
      var value = new proto_environment_command_pb.RenameProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.RenameProjectCommand.deserializeBinaryFromReader);
      msg.setRenameCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UpdateProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UpdateProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getChangeDescriptionCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.ChangeDescriptionProjectCommand.serializeBinaryToWriter
    );
  }
  f = message.getRenameCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_environment_command_pb.RenameProjectCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.UpdateProjectRequest} returns this
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ChangeDescriptionProjectCommand change_description_command = 2;
 * @return {?proto.bucketeer.environment.ChangeDescriptionProjectCommand}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.getChangeDescriptionCommand = function() {
  return /** @type{?proto.bucketeer.environment.ChangeDescriptionProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.ChangeDescriptionProjectCommand, 2));
};


/**
 * @param {?proto.bucketeer.environment.ChangeDescriptionProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.UpdateProjectRequest} returns this
*/
proto.bucketeer.environment.UpdateProjectRequest.prototype.setChangeDescriptionCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UpdateProjectRequest} returns this
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.clearChangeDescriptionCommand = function() {
  return this.setChangeDescriptionCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.hasChangeDescriptionCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional RenameProjectCommand rename_command = 3;
 * @return {?proto.bucketeer.environment.RenameProjectCommand}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.getRenameCommand = function() {
  return /** @type{?proto.bucketeer.environment.RenameProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.RenameProjectCommand, 3));
};


/**
 * @param {?proto.bucketeer.environment.RenameProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.UpdateProjectRequest} returns this
*/
proto.bucketeer.environment.UpdateProjectRequest.prototype.setRenameCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.UpdateProjectRequest} returns this
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.clearRenameCommand = function() {
  return this.setRenameCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.UpdateProjectRequest.prototype.hasRenameCommand = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.UpdateProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.UpdateProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.UpdateProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.UpdateProjectResponse}
 */
proto.bucketeer.environment.UpdateProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.UpdateProjectResponse;
  return proto.bucketeer.environment.UpdateProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.UpdateProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.UpdateProjectResponse}
 */
proto.bucketeer.environment.UpdateProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.UpdateProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.UpdateProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.UpdateProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.UpdateProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.EnableProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.EnableProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.EnableProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_environment_command_pb.EnableProjectCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.EnableProjectRequest}
 */
proto.bucketeer.environment.EnableProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.EnableProjectRequest;
  return proto.bucketeer.environment.EnableProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.EnableProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.EnableProjectRequest}
 */
proto.bucketeer.environment.EnableProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.EnableProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.EnableProjectCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.EnableProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.EnableProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.EnableProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.EnableProjectCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.EnableProjectRequest} returns this
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional EnableProjectCommand command = 2;
 * @return {?proto.bucketeer.environment.EnableProjectCommand}
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.EnableProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.EnableProjectCommand, 2));
};


/**
 * @param {?proto.bucketeer.environment.EnableProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.EnableProjectRequest} returns this
*/
proto.bucketeer.environment.EnableProjectRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.EnableProjectRequest} returns this
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.EnableProjectRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.EnableProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.EnableProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.EnableProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.EnableProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.EnableProjectResponse}
 */
proto.bucketeer.environment.EnableProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.EnableProjectResponse;
  return proto.bucketeer.environment.EnableProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.EnableProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.EnableProjectResponse}
 */
proto.bucketeer.environment.EnableProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.EnableProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.EnableProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.EnableProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.EnableProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.DisableProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.DisableProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DisableProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_environment_command_pb.DisableProjectCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.DisableProjectRequest}
 */
proto.bucketeer.environment.DisableProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.DisableProjectRequest;
  return proto.bucketeer.environment.DisableProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.DisableProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.DisableProjectRequest}
 */
proto.bucketeer.environment.DisableProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.DisableProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.DisableProjectCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.DisableProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.DisableProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DisableProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.DisableProjectCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.DisableProjectRequest} returns this
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DisableProjectCommand command = 2;
 * @return {?proto.bucketeer.environment.DisableProjectCommand}
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.DisableProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.DisableProjectCommand, 2));
};


/**
 * @param {?proto.bucketeer.environment.DisableProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.DisableProjectRequest} returns this
*/
proto.bucketeer.environment.DisableProjectRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.DisableProjectRequest} returns this
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.DisableProjectRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.DisableProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.DisableProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.DisableProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DisableProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.DisableProjectResponse}
 */
proto.bucketeer.environment.DisableProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.DisableProjectResponse;
  return proto.bucketeer.environment.DisableProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.DisableProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.DisableProjectResponse}
 */
proto.bucketeer.environment.DisableProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.DisableProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.DisableProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.DisableProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.DisableProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ConvertTrialProjectRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ConvertTrialProjectRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_environment_command_pb.ConvertTrialProjectCommand.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ConvertTrialProjectRequest}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ConvertTrialProjectRequest;
  return proto.bucketeer.environment.ConvertTrialProjectRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ConvertTrialProjectRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ConvertTrialProjectRequest}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto_environment_command_pb.ConvertTrialProjectCommand;
      reader.readMessage(value,proto_environment_command_pb.ConvertTrialProjectCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ConvertTrialProjectRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ConvertTrialProjectRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_environment_command_pb.ConvertTrialProjectCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.ConvertTrialProjectRequest} returns this
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ConvertTrialProjectCommand command = 2;
 * @return {?proto.bucketeer.environment.ConvertTrialProjectCommand}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.environment.ConvertTrialProjectCommand} */ (
    jspb.Message.getWrapperField(this, proto_environment_command_pb.ConvertTrialProjectCommand, 2));
};


/**
 * @param {?proto.bucketeer.environment.ConvertTrialProjectCommand|undefined} value
 * @return {!proto.bucketeer.environment.ConvertTrialProjectRequest} returns this
*/
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.environment.ConvertTrialProjectRequest} returns this
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.environment.ConvertTrialProjectRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.environment.ConvertTrialProjectResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.ConvertTrialProjectResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.ConvertTrialProjectResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ConvertTrialProjectResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.environment.ConvertTrialProjectResponse}
 */
proto.bucketeer.environment.ConvertTrialProjectResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.ConvertTrialProjectResponse;
  return proto.bucketeer.environment.ConvertTrialProjectResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.ConvertTrialProjectResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.ConvertTrialProjectResponse}
 */
proto.bucketeer.environment.ConvertTrialProjectResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.environment.ConvertTrialProjectResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.ConvertTrialProjectResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.ConvertTrialProjectResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.ConvertTrialProjectResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};


goog.object.extend(exports, proto.bucketeer.environment);
