// source: proto/feature/service.proto
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

var google_api_annotations_pb = require('../../google/api/annotations_pb.js');
goog.object.extend(proto, google_api_annotations_pb);
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.object.extend(proto, google_protobuf_wrappers_pb);
var proto_feature_command_pb = require('../../proto/feature/command_pb.js');
goog.object.extend(proto, proto_feature_command_pb);
var proto_feature_feature_pb = require('../../proto/feature/feature_pb.js');
goog.object.extend(proto, proto_feature_feature_pb);
var proto_feature_evaluation_pb = require('../../proto/feature/evaluation_pb.js');
goog.object.extend(proto, proto_feature_evaluation_pb);
var proto_user_user_pb = require('../../proto/user/user_pb.js');
goog.object.extend(proto, proto_user_user_pb);
var proto_feature_segment_pb = require('../../proto/feature/segment_pb.js');
goog.object.extend(proto, proto_feature_segment_pb);
var proto_feature_flag_trigger_pb = require('../../proto/feature/flag_trigger_pb.js');
goog.object.extend(proto, proto_feature_flag_trigger_pb);
goog.exportSymbol('proto.bucketeer.feature.AddSegmentUserRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddSegmentUserResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ArchiveFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ArchiveFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.BulkDownloadSegmentUsersRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.BulkDownloadSegmentUsersResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.BulkUploadSegmentUsersRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.BulkUploadSegmentUsersResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.CloneFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.CloneFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateSegmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateSegmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteSegmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteSegmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteSegmentUserRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteSegmentUserResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.DisableFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.DisableFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.DisableFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.DisableFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.EnableFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.EnableFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.EnableFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.EnableFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.EvaluateFeaturesRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.EvaluateFeaturesResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.FlagTriggerWebhookRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.FlagTriggerWebhookResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetFeaturesRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetFeaturesResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetSegmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetSegmentResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetSegmentUserRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.GetSegmentUserResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListEnabledFeaturesRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListEnabledFeaturesResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFeaturesRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFeaturesRequest.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFeaturesRequest.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFeaturesResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFlagTriggersRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFlagTriggersRequest.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFlagTriggersRequest.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFlagTriggersResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListSegmentUsersRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListSegmentUsersResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListSegmentsRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListSegmentsRequest.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListSegmentsRequest.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListSegmentsResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListTagsRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListTagsRequest.OrderBy', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListTagsRequest.OrderDirection', null, global);
goog.exportSymbol('proto.bucketeer.feature.ListTagsResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.ResetFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.ResetFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.UnarchiveFeatureRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.UnarchiveFeatureResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFeatureDetailsRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFeatureDetailsResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFeatureTargetingRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFeatureTargetingResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFeatureVariationsRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFeatureVariationsResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFlagTriggerRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFlagTriggerResponse', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateSegmentRequest', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateSegmentResponse', null, global);
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
proto.bucketeer.feature.GetFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetFeatureRequest.displayName = 'proto.bucketeer.feature.GetFeatureRequest';
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
proto.bucketeer.feature.GetFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetFeatureResponse.displayName = 'proto.bucketeer.feature.GetFeatureResponse';
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
proto.bucketeer.feature.GetFeaturesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.GetFeaturesRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.GetFeaturesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetFeaturesRequest.displayName = 'proto.bucketeer.feature.GetFeaturesRequest';
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
proto.bucketeer.feature.GetFeaturesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.GetFeaturesResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.GetFeaturesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetFeaturesResponse.displayName = 'proto.bucketeer.feature.GetFeaturesResponse';
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
proto.bucketeer.feature.ListFeaturesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListFeaturesRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListFeaturesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListFeaturesRequest.displayName = 'proto.bucketeer.feature.ListFeaturesRequest';
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
proto.bucketeer.feature.ListFeaturesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListFeaturesResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListFeaturesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListFeaturesResponse.displayName = 'proto.bucketeer.feature.ListFeaturesResponse';
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
proto.bucketeer.feature.ListEnabledFeaturesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListEnabledFeaturesRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListEnabledFeaturesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListEnabledFeaturesRequest.displayName = 'proto.bucketeer.feature.ListEnabledFeaturesRequest';
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
proto.bucketeer.feature.ListEnabledFeaturesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListEnabledFeaturesResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListEnabledFeaturesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListEnabledFeaturesResponse.displayName = 'proto.bucketeer.feature.ListEnabledFeaturesResponse';
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
proto.bucketeer.feature.CreateFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateFeatureRequest.displayName = 'proto.bucketeer.feature.CreateFeatureRequest';
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
proto.bucketeer.feature.CreateFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateFeatureResponse.displayName = 'proto.bucketeer.feature.CreateFeatureResponse';
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
proto.bucketeer.feature.EnableFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EnableFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EnableFeatureRequest.displayName = 'proto.bucketeer.feature.EnableFeatureRequest';
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
proto.bucketeer.feature.EnableFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EnableFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EnableFeatureResponse.displayName = 'proto.bucketeer.feature.EnableFeatureResponse';
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
proto.bucketeer.feature.DisableFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DisableFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DisableFeatureRequest.displayName = 'proto.bucketeer.feature.DisableFeatureRequest';
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
proto.bucketeer.feature.DisableFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DisableFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DisableFeatureResponse.displayName = 'proto.bucketeer.feature.DisableFeatureResponse';
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
proto.bucketeer.feature.ArchiveFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ArchiveFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ArchiveFeatureRequest.displayName = 'proto.bucketeer.feature.ArchiveFeatureRequest';
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
proto.bucketeer.feature.ArchiveFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ArchiveFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ArchiveFeatureResponse.displayName = 'proto.bucketeer.feature.ArchiveFeatureResponse';
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
proto.bucketeer.feature.UnarchiveFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UnarchiveFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UnarchiveFeatureRequest.displayName = 'proto.bucketeer.feature.UnarchiveFeatureRequest';
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
proto.bucketeer.feature.UnarchiveFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UnarchiveFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UnarchiveFeatureResponse.displayName = 'proto.bucketeer.feature.UnarchiveFeatureResponse';
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
proto.bucketeer.feature.DeleteFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteFeatureRequest.displayName = 'proto.bucketeer.feature.DeleteFeatureRequest';
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
proto.bucketeer.feature.DeleteFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteFeatureResponse.displayName = 'proto.bucketeer.feature.DeleteFeatureResponse';
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
proto.bucketeer.feature.UpdateFeatureDetailsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.UpdateFeatureDetailsRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFeatureDetailsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFeatureDetailsRequest.displayName = 'proto.bucketeer.feature.UpdateFeatureDetailsRequest';
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
proto.bucketeer.feature.UpdateFeatureDetailsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFeatureDetailsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFeatureDetailsResponse.displayName = 'proto.bucketeer.feature.UpdateFeatureDetailsResponse';
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
proto.bucketeer.feature.UpdateFeatureVariationsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.UpdateFeatureVariationsRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFeatureVariationsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFeatureVariationsRequest.displayName = 'proto.bucketeer.feature.UpdateFeatureVariationsRequest';
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
proto.bucketeer.feature.UpdateFeatureVariationsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFeatureVariationsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFeatureVariationsResponse.displayName = 'proto.bucketeer.feature.UpdateFeatureVariationsResponse';
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
proto.bucketeer.feature.UpdateFeatureTargetingRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.UpdateFeatureTargetingRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFeatureTargetingRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFeatureTargetingRequest.displayName = 'proto.bucketeer.feature.UpdateFeatureTargetingRequest';
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
proto.bucketeer.feature.UpdateFeatureTargetingResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFeatureTargetingResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFeatureTargetingResponse.displayName = 'proto.bucketeer.feature.UpdateFeatureTargetingResponse';
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
proto.bucketeer.feature.CloneFeatureRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CloneFeatureRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CloneFeatureRequest.displayName = 'proto.bucketeer.feature.CloneFeatureRequest';
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
proto.bucketeer.feature.CloneFeatureResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CloneFeatureResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CloneFeatureResponse.displayName = 'proto.bucketeer.feature.CloneFeatureResponse';
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
proto.bucketeer.feature.CreateSegmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateSegmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateSegmentRequest.displayName = 'proto.bucketeer.feature.CreateSegmentRequest';
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
proto.bucketeer.feature.CreateSegmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateSegmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateSegmentResponse.displayName = 'proto.bucketeer.feature.CreateSegmentResponse';
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
proto.bucketeer.feature.GetSegmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetSegmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetSegmentRequest.displayName = 'proto.bucketeer.feature.GetSegmentRequest';
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
proto.bucketeer.feature.GetSegmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetSegmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetSegmentResponse.displayName = 'proto.bucketeer.feature.GetSegmentResponse';
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
proto.bucketeer.feature.ListSegmentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ListSegmentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListSegmentsRequest.displayName = 'proto.bucketeer.feature.ListSegmentsRequest';
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
proto.bucketeer.feature.ListSegmentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListSegmentsResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListSegmentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListSegmentsResponse.displayName = 'proto.bucketeer.feature.ListSegmentsResponse';
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
proto.bucketeer.feature.DeleteSegmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteSegmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteSegmentRequest.displayName = 'proto.bucketeer.feature.DeleteSegmentRequest';
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
proto.bucketeer.feature.DeleteSegmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteSegmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteSegmentResponse.displayName = 'proto.bucketeer.feature.DeleteSegmentResponse';
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
proto.bucketeer.feature.UpdateSegmentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.UpdateSegmentRequest.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.UpdateSegmentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateSegmentRequest.displayName = 'proto.bucketeer.feature.UpdateSegmentRequest';
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
proto.bucketeer.feature.UpdateSegmentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateSegmentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateSegmentResponse.displayName = 'proto.bucketeer.feature.UpdateSegmentResponse';
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
proto.bucketeer.feature.AddSegmentUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddSegmentUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddSegmentUserRequest.displayName = 'proto.bucketeer.feature.AddSegmentUserRequest';
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
proto.bucketeer.feature.AddSegmentUserResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddSegmentUserResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddSegmentUserResponse.displayName = 'proto.bucketeer.feature.AddSegmentUserResponse';
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
proto.bucketeer.feature.DeleteSegmentUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteSegmentUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteSegmentUserRequest.displayName = 'proto.bucketeer.feature.DeleteSegmentUserRequest';
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
proto.bucketeer.feature.DeleteSegmentUserResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteSegmentUserResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteSegmentUserResponse.displayName = 'proto.bucketeer.feature.DeleteSegmentUserResponse';
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
proto.bucketeer.feature.GetSegmentUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetSegmentUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetSegmentUserRequest.displayName = 'proto.bucketeer.feature.GetSegmentUserRequest';
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
proto.bucketeer.feature.GetSegmentUserResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetSegmentUserResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetSegmentUserResponse.displayName = 'proto.bucketeer.feature.GetSegmentUserResponse';
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
proto.bucketeer.feature.ListSegmentUsersRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ListSegmentUsersRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListSegmentUsersRequest.displayName = 'proto.bucketeer.feature.ListSegmentUsersRequest';
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
proto.bucketeer.feature.ListSegmentUsersResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListSegmentUsersResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListSegmentUsersResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListSegmentUsersResponse.displayName = 'proto.bucketeer.feature.ListSegmentUsersResponse';
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
proto.bucketeer.feature.BulkUploadSegmentUsersRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.BulkUploadSegmentUsersRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.BulkUploadSegmentUsersRequest.displayName = 'proto.bucketeer.feature.BulkUploadSegmentUsersRequest';
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
proto.bucketeer.feature.BulkUploadSegmentUsersResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.BulkUploadSegmentUsersResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.BulkUploadSegmentUsersResponse.displayName = 'proto.bucketeer.feature.BulkUploadSegmentUsersResponse';
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
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.BulkDownloadSegmentUsersRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.displayName = 'proto.bucketeer.feature.BulkDownloadSegmentUsersRequest';
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
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.BulkDownloadSegmentUsersResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.displayName = 'proto.bucketeer.feature.BulkDownloadSegmentUsersResponse';
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
proto.bucketeer.feature.EvaluateFeaturesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EvaluateFeaturesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EvaluateFeaturesRequest.displayName = 'proto.bucketeer.feature.EvaluateFeaturesRequest';
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
proto.bucketeer.feature.EvaluateFeaturesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EvaluateFeaturesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EvaluateFeaturesResponse.displayName = 'proto.bucketeer.feature.EvaluateFeaturesResponse';
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
proto.bucketeer.feature.ListTagsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ListTagsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListTagsRequest.displayName = 'proto.bucketeer.feature.ListTagsRequest';
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
proto.bucketeer.feature.ListTagsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListTagsResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListTagsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListTagsResponse.displayName = 'proto.bucketeer.feature.ListTagsResponse';
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
proto.bucketeer.feature.CreateFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateFlagTriggerRequest.displayName = 'proto.bucketeer.feature.CreateFlagTriggerRequest';
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
proto.bucketeer.feature.CreateFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateFlagTriggerResponse.displayName = 'proto.bucketeer.feature.CreateFlagTriggerResponse';
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
proto.bucketeer.feature.DeleteFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteFlagTriggerRequest.displayName = 'proto.bucketeer.feature.DeleteFlagTriggerRequest';
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
proto.bucketeer.feature.DeleteFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteFlagTriggerResponse.displayName = 'proto.bucketeer.feature.DeleteFlagTriggerResponse';
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
proto.bucketeer.feature.UpdateFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFlagTriggerRequest.displayName = 'proto.bucketeer.feature.UpdateFlagTriggerRequest';
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
proto.bucketeer.feature.UpdateFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFlagTriggerResponse.displayName = 'proto.bucketeer.feature.UpdateFlagTriggerResponse';
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
proto.bucketeer.feature.EnableFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EnableFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EnableFlagTriggerRequest.displayName = 'proto.bucketeer.feature.EnableFlagTriggerRequest';
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
proto.bucketeer.feature.EnableFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EnableFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EnableFlagTriggerResponse.displayName = 'proto.bucketeer.feature.EnableFlagTriggerResponse';
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
proto.bucketeer.feature.DisableFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DisableFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DisableFlagTriggerRequest.displayName = 'proto.bucketeer.feature.DisableFlagTriggerRequest';
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
proto.bucketeer.feature.DisableFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DisableFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DisableFlagTriggerResponse.displayName = 'proto.bucketeer.feature.DisableFlagTriggerResponse';
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
proto.bucketeer.feature.ResetFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ResetFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ResetFlagTriggerRequest.displayName = 'proto.bucketeer.feature.ResetFlagTriggerRequest';
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
proto.bucketeer.feature.ResetFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ResetFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ResetFlagTriggerResponse.displayName = 'proto.bucketeer.feature.ResetFlagTriggerResponse';
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
proto.bucketeer.feature.GetFlagTriggerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetFlagTriggerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetFlagTriggerRequest.displayName = 'proto.bucketeer.feature.GetFlagTriggerRequest';
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
proto.bucketeer.feature.GetFlagTriggerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.GetFlagTriggerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.GetFlagTriggerResponse.displayName = 'proto.bucketeer.feature.GetFlagTriggerResponse';
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
proto.bucketeer.feature.ListFlagTriggersRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ListFlagTriggersRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListFlagTriggersRequest.displayName = 'proto.bucketeer.feature.ListFlagTriggersRequest';
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
proto.bucketeer.feature.ListFlagTriggersResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.ListFlagTriggersResponse.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.ListFlagTriggersResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListFlagTriggersResponse.displayName = 'proto.bucketeer.feature.ListFlagTriggersResponse';
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
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.displayName = 'proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl';
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
proto.bucketeer.feature.FlagTriggerWebhookRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.FlagTriggerWebhookRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.FlagTriggerWebhookRequest.displayName = 'proto.bucketeer.feature.FlagTriggerWebhookRequest';
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
proto.bucketeer.feature.FlagTriggerWebhookResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.FlagTriggerWebhookResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.FlagTriggerWebhookResponse.displayName = 'proto.bucketeer.feature.FlagTriggerWebhookResponse';
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
proto.bucketeer.feature.GetFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.GetFeatureRequest}
 */
proto.bucketeer.feature.GetFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetFeatureRequest;
  return proto.bucketeer.feature.GetFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetFeatureRequest}
 */
proto.bucketeer.feature.GetFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.GetFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeatureRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.GetFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetFeatureRequest} returns this
 */
proto.bucketeer.feature.GetFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.GetFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetFeatureRequest} returns this
 */
proto.bucketeer.feature.GetFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.GetFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeatureResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    feature: (f = msg.getFeature()) && proto_feature_feature_pb.Feature.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.GetFeatureResponse}
 */
proto.bucketeer.feature.GetFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetFeatureResponse;
  return proto.bucketeer.feature.GetFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetFeatureResponse}
 */
proto.bucketeer.feature.GetFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_feature_pb.Feature;
      reader.readMessage(value,proto_feature_feature_pb.Feature.deserializeBinaryFromReader);
      msg.setFeature(value);
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
proto.bucketeer.feature.GetFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeatureResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeature();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_feature_pb.Feature.serializeBinaryToWriter
    );
  }
};


/**
 * optional Feature feature = 1;
 * @return {?proto.bucketeer.feature.Feature}
 */
proto.bucketeer.feature.GetFeatureResponse.prototype.getFeature = function() {
  return /** @type{?proto.bucketeer.feature.Feature} */ (
    jspb.Message.getWrapperField(this, proto_feature_feature_pb.Feature, 1));
};


/**
 * @param {?proto.bucketeer.feature.Feature|undefined} value
 * @return {!proto.bucketeer.feature.GetFeatureResponse} returns this
*/
proto.bucketeer.feature.GetFeatureResponse.prototype.setFeature = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.GetFeatureResponse} returns this
 */
proto.bucketeer.feature.GetFeatureResponse.prototype.clearFeature = function() {
  return this.setFeature(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.GetFeatureResponse.prototype.hasFeature = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.GetFeaturesRequest.repeatedFields_ = [2];



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
proto.bucketeer.feature.GetFeaturesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetFeaturesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetFeaturesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeaturesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 1, ""),
    idsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f
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
 * @return {!proto.bucketeer.feature.GetFeaturesRequest}
 */
proto.bucketeer.feature.GetFeaturesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetFeaturesRequest;
  return proto.bucketeer.feature.GetFeaturesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetFeaturesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetFeaturesRequest}
 */
proto.bucketeer.feature.GetFeaturesRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addIds(value);
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
proto.bucketeer.feature.GetFeaturesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetFeaturesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetFeaturesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeaturesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
};


/**
 * optional string environment_namespace = 1;
 * @return {string}
 */
proto.bucketeer.feature.GetFeaturesRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetFeaturesRequest} returns this
 */
proto.bucketeer.feature.GetFeaturesRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string ids = 2;
 * @return {!Array<string>}
 */
proto.bucketeer.feature.GetFeaturesRequest.prototype.getIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.feature.GetFeaturesRequest} returns this
 */
proto.bucketeer.feature.GetFeaturesRequest.prototype.setIdsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.GetFeaturesRequest} returns this
 */
proto.bucketeer.feature.GetFeaturesRequest.prototype.addIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.GetFeaturesRequest} returns this
 */
proto.bucketeer.feature.GetFeaturesRequest.prototype.clearIdsList = function() {
  return this.setIdsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.GetFeaturesResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.GetFeaturesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetFeaturesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetFeaturesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeaturesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    featuresList: jspb.Message.toObjectList(msg.getFeaturesList(),
    proto_feature_feature_pb.Feature.toObject, includeInstance)
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
 * @return {!proto.bucketeer.feature.GetFeaturesResponse}
 */
proto.bucketeer.feature.GetFeaturesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetFeaturesResponse;
  return proto.bucketeer.feature.GetFeaturesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetFeaturesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetFeaturesResponse}
 */
proto.bucketeer.feature.GetFeaturesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_feature_pb.Feature;
      reader.readMessage(value,proto_feature_feature_pb.Feature.deserializeBinaryFromReader);
      msg.addFeatures(value);
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
proto.bucketeer.feature.GetFeaturesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetFeaturesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetFeaturesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFeaturesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeaturesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_feature_feature_pb.Feature.serializeBinaryToWriter
    );
  }
};


/**
 * repeated Feature features = 1;
 * @return {!Array<!proto.bucketeer.feature.Feature>}
 */
proto.bucketeer.feature.GetFeaturesResponse.prototype.getFeaturesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Feature>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_feature_pb.Feature, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Feature>} value
 * @return {!proto.bucketeer.feature.GetFeaturesResponse} returns this
*/
proto.bucketeer.feature.GetFeaturesResponse.prototype.setFeaturesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.Feature=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Feature}
 */
proto.bucketeer.feature.GetFeaturesResponse.prototype.addFeatures = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.Feature, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.GetFeaturesResponse} returns this
 */
proto.bucketeer.feature.GetFeaturesResponse.prototype.clearFeaturesList = function() {
  return this.setFeaturesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListFeaturesRequest.repeatedFields_ = [3];



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
proto.bucketeer.feature.ListFeaturesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListFeaturesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListFeaturesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFeaturesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    tagsList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f,
    orderBy: jspb.Message.getFieldWithDefault(msg, 4, 0),
    orderDirection: jspb.Message.getFieldWithDefault(msg, 5, 0),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 6, ""),
    maintainer: jspb.Message.getFieldWithDefault(msg, 7, ""),
    enabled: (f = msg.getEnabled()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    hasExperiment: (f = msg.getHasExperiment()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    searchKeyword: jspb.Message.getFieldWithDefault(msg, 10, ""),
    archived: (f = msg.getArchived()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
    hasPrerequisites: (f = msg.getHasPrerequisites()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ListFeaturesRequest}
 */
proto.bucketeer.feature.ListFeaturesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListFeaturesRequest;
  return proto.bucketeer.feature.ListFeaturesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListFeaturesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListFeaturesRequest}
 */
proto.bucketeer.feature.ListFeaturesRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.addTags(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.feature.ListFeaturesRequest.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 5:
      var value = /** @type {!proto.bucketeer.feature.ListFeaturesRequest.OrderDirection} */ (reader.readEnum());
      msg.setOrderDirection(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setMaintainer(value);
      break;
    case 8:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setEnabled(value);
      break;
    case 9:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setHasExperiment(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setSearchKeyword(value);
      break;
    case 11:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setArchived(value);
      break;
    case 12:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setHasPrerequisites(value);
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
proto.bucketeer.feature.ListFeaturesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListFeaturesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListFeaturesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFeaturesRequest.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
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
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getMaintainer();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getEnabled();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getHasExperiment();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getSearchKeyword();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getArchived();
  if (f != null) {
    writer.writeMessage(
      11,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getHasPrerequisites();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.feature.ListFeaturesRequest.OrderBy = {
  DEFAULT: 0,
  NAME: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3,
  TAGS: 4,
  ENABLED: 5
};

/**
 * @enum {number}
 */
proto.bucketeer.feature.ListFeaturesRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated string tags = 3;
 * @return {!Array<string>}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.feature.ListFeaturesRequest.OrderBy}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.feature.ListFeaturesRequest.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListFeaturesRequest.OrderBy} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.feature.ListFeaturesRequest.OrderDirection}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.feature.ListFeaturesRequest.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListFeaturesRequest.OrderDirection} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional string environment_namespace = 6;
 * @return {string}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string maintainer = 7;
 * @return {string}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getMaintainer = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setMaintainer = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional google.protobuf.BoolValue enabled = 8;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getEnabled = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 8));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
*/
proto.bucketeer.feature.ListFeaturesRequest.prototype.setEnabled = function(value) {
  return jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.clearEnabled = function() {
  return this.setEnabled(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.hasEnabled = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional google.protobuf.BoolValue has_experiment = 9;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getHasExperiment = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 9));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
*/
proto.bucketeer.feature.ListFeaturesRequest.prototype.setHasExperiment = function(value) {
  return jspb.Message.setWrapperField(this, 9, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.clearHasExperiment = function() {
  return this.setHasExperiment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.hasHasExperiment = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional string search_keyword = 10;
 * @return {string}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getSearchKeyword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.setSearchKeyword = function(value) {
  return jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional google.protobuf.BoolValue archived = 11;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getArchived = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 11));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
*/
proto.bucketeer.feature.ListFeaturesRequest.prototype.setArchived = function(value) {
  return jspb.Message.setWrapperField(this, 11, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.clearArchived = function() {
  return this.setArchived(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.hasArchived = function() {
  return jspb.Message.getField(this, 11) != null;
};


/**
 * optional google.protobuf.BoolValue has_prerequisites = 12;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.getHasPrerequisites = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 12));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
*/
proto.bucketeer.feature.ListFeaturesRequest.prototype.setHasPrerequisites = function(value) {
  return jspb.Message.setWrapperField(this, 12, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.clearHasPrerequisites = function() {
  return this.setHasPrerequisites(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListFeaturesRequest.prototype.hasHasPrerequisites = function() {
  return jspb.Message.getField(this, 12) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListFeaturesResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.ListFeaturesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListFeaturesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListFeaturesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFeaturesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    featuresList: jspb.Message.toObjectList(msg.getFeaturesList(),
    proto_feature_feature_pb.Feature.toObject, includeInstance),
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
 * @return {!proto.bucketeer.feature.ListFeaturesResponse}
 */
proto.bucketeer.feature.ListFeaturesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListFeaturesResponse;
  return proto.bucketeer.feature.ListFeaturesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListFeaturesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListFeaturesResponse}
 */
proto.bucketeer.feature.ListFeaturesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_feature_pb.Feature;
      reader.readMessage(value,proto_feature_feature_pb.Feature.deserializeBinaryFromReader);
      msg.addFeatures(value);
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
proto.bucketeer.feature.ListFeaturesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListFeaturesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListFeaturesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFeaturesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeaturesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_feature_feature_pb.Feature.serializeBinaryToWriter
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
 * repeated Feature features = 1;
 * @return {!Array<!proto.bucketeer.feature.Feature>}
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.getFeaturesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Feature>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_feature_pb.Feature, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Feature>} value
 * @return {!proto.bucketeer.feature.ListFeaturesResponse} returns this
*/
proto.bucketeer.feature.ListFeaturesResponse.prototype.setFeaturesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.Feature=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Feature}
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.addFeatures = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.Feature, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListFeaturesResponse} returns this
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.clearFeaturesList = function() {
  return this.setFeaturesList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFeaturesResponse} returns this
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListFeaturesResponse} returns this
 */
proto.bucketeer.feature.ListFeaturesResponse.prototype.setTotalCount = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.repeatedFields_ = [3];



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
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListEnabledFeaturesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListEnabledFeaturesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    tagsList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f,
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest}
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListEnabledFeaturesRequest;
  return proto.bucketeer.feature.ListEnabledFeaturesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListEnabledFeaturesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest}
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.addTags(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListEnabledFeaturesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListEnabledFeaturesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      3,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated string tags = 3;
 * @return {!Array<string>}
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * optional string environment_namespace = 4;
 * @return {string}
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesRequest} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListEnabledFeaturesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListEnabledFeaturesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    featuresList: jspb.Message.toObjectList(msg.getFeaturesList(),
    proto_feature_feature_pb.Feature.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesResponse}
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListEnabledFeaturesResponse;
  return proto.bucketeer.feature.ListEnabledFeaturesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListEnabledFeaturesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesResponse}
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_feature_pb.Feature;
      reader.readMessage(value,proto_feature_feature_pb.Feature.deserializeBinaryFromReader);
      msg.addFeatures(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
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
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListEnabledFeaturesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListEnabledFeaturesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeaturesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_feature_feature_pb.Feature.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated Feature features = 1;
 * @return {!Array<!proto.bucketeer.feature.Feature>}
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.getFeaturesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Feature>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_feature_pb.Feature, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Feature>} value
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesResponse} returns this
*/
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.setFeaturesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.Feature=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Feature}
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.addFeatures = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.Feature, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesResponse} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.clearFeaturesList = function() {
  return this.setFeaturesList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListEnabledFeaturesResponse} returns this
 */
proto.bucketeer.feature.ListEnabledFeaturesResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.CreateFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && proto_feature_command_pb.CreateFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.CreateFeatureRequest}
 */
proto.bucketeer.feature.CreateFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateFeatureRequest;
  return proto.bucketeer.feature.CreateFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateFeatureRequest}
 */
proto.bucketeer.feature.CreateFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_command_pb.CreateFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.CreateFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.CreateFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFeatureRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_command_pb.CreateFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional CreateFeatureCommand command = 1;
 * @return {?proto.bucketeer.feature.CreateFeatureCommand}
 */
proto.bucketeer.feature.CreateFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.CreateFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.CreateFeatureCommand, 1));
};


/**
 * @param {?proto.bucketeer.feature.CreateFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.CreateFeatureRequest} returns this
*/
proto.bucketeer.feature.CreateFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateFeatureRequest} returns this
 */
proto.bucketeer.feature.CreateFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.CreateFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFeatureRequest} returns this
 */
proto.bucketeer.feature.CreateFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.CreateFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.CreateFeatureResponse}
 */
proto.bucketeer.feature.CreateFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateFeatureResponse;
  return proto.bucketeer.feature.CreateFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateFeatureResponse}
 */
proto.bucketeer.feature.CreateFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.CreateFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFeatureResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.EnableFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EnableFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EnableFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.EnableFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.EnableFeatureRequest}
 */
proto.bucketeer.feature.EnableFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EnableFeatureRequest;
  return proto.bucketeer.feature.EnableFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EnableFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EnableFeatureRequest}
 */
proto.bucketeer.feature.EnableFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.EnableFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.EnableFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.EnableFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EnableFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EnableFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFeatureRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.EnableFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EnableFeatureRequest} returns this
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional EnableFeatureCommand command = 2;
 * @return {?proto.bucketeer.feature.EnableFeatureCommand}
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.EnableFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.EnableFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.EnableFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.EnableFeatureRequest} returns this
*/
proto.bucketeer.feature.EnableFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.EnableFeatureRequest} returns this
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EnableFeatureRequest} returns this
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EnableFeatureRequest} returns this
 */
proto.bucketeer.feature.EnableFeatureRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.EnableFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EnableFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EnableFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.EnableFeatureResponse}
 */
proto.bucketeer.feature.EnableFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EnableFeatureResponse;
  return proto.bucketeer.feature.EnableFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EnableFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EnableFeatureResponse}
 */
proto.bucketeer.feature.EnableFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.EnableFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EnableFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EnableFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFeatureResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DisableFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DisableFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DisableFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.DisableFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.DisableFeatureRequest}
 */
proto.bucketeer.feature.DisableFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DisableFeatureRequest;
  return proto.bucketeer.feature.DisableFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DisableFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DisableFeatureRequest}
 */
proto.bucketeer.feature.DisableFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.DisableFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.DisableFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.DisableFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DisableFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DisableFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFeatureRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.DisableFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DisableFeatureRequest} returns this
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DisableFeatureCommand command = 2;
 * @return {?proto.bucketeer.feature.DisableFeatureCommand}
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.DisableFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.DisableFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.DisableFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.DisableFeatureRequest} returns this
*/
proto.bucketeer.feature.DisableFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.DisableFeatureRequest} returns this
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DisableFeatureRequest} returns this
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DisableFeatureRequest} returns this
 */
proto.bucketeer.feature.DisableFeatureRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.DisableFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DisableFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DisableFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DisableFeatureResponse}
 */
proto.bucketeer.feature.DisableFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DisableFeatureResponse;
  return proto.bucketeer.feature.DisableFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DisableFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DisableFeatureResponse}
 */
proto.bucketeer.feature.DisableFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DisableFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DisableFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DisableFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFeatureResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ArchiveFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ArchiveFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ArchiveFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.ArchiveFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ArchiveFeatureRequest;
  return proto.bucketeer.feature.ArchiveFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ArchiveFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.ArchiveFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.ArchiveFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ArchiveFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ArchiveFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ArchiveFeatureRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.ArchiveFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional ArchiveFeatureCommand command = 2;
 * @return {?proto.bucketeer.feature.ArchiveFeatureCommand}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.ArchiveFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.ArchiveFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.ArchiveFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest} returns this
*/
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ArchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.ArchiveFeatureRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.ArchiveFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ArchiveFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ArchiveFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ArchiveFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.ArchiveFeatureResponse}
 */
proto.bucketeer.feature.ArchiveFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ArchiveFeatureResponse;
  return proto.bucketeer.feature.ArchiveFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ArchiveFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ArchiveFeatureResponse}
 */
proto.bucketeer.feature.ArchiveFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.ArchiveFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ArchiveFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ArchiveFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ArchiveFeatureResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UnarchiveFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UnarchiveFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.UnarchiveFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UnarchiveFeatureRequest;
  return proto.bucketeer.feature.UnarchiveFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UnarchiveFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.UnarchiveFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.UnarchiveFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UnarchiveFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UnarchiveFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.UnarchiveFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional UnarchiveFeatureCommand command = 2;
 * @return {?proto.bucketeer.feature.UnarchiveFeatureCommand}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.UnarchiveFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.UnarchiveFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.UnarchiveFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest} returns this
*/
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UnarchiveFeatureRequest} returns this
 */
proto.bucketeer.feature.UnarchiveFeatureRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.UnarchiveFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UnarchiveFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UnarchiveFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UnarchiveFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UnarchiveFeatureResponse}
 */
proto.bucketeer.feature.UnarchiveFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UnarchiveFeatureResponse;
  return proto.bucketeer.feature.UnarchiveFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UnarchiveFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UnarchiveFeatureResponse}
 */
proto.bucketeer.feature.UnarchiveFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UnarchiveFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UnarchiveFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UnarchiveFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UnarchiveFeatureResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DeleteFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.DeleteFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest}
 */
proto.bucketeer.feature.DeleteFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteFeatureRequest;
  return proto.bucketeer.feature.DeleteFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest}
 */
proto.bucketeer.feature.DeleteFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.DeleteFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.DeleteFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.DeleteFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFeatureRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.DeleteFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest} returns this
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DeleteFeatureCommand command = 2;
 * @return {?proto.bucketeer.feature.DeleteFeatureCommand}
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.DeleteFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.DeleteFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.DeleteFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest} returns this
*/
proto.bucketeer.feature.DeleteFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest} returns this
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest} returns this
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteFeatureRequest} returns this
 */
proto.bucketeer.feature.DeleteFeatureRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.DeleteFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteFeatureResponse}
 */
proto.bucketeer.feature.DeleteFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteFeatureResponse;
  return proto.bucketeer.feature.DeleteFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteFeatureResponse}
 */
proto.bucketeer.feature.DeleteFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFeatureResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.repeatedFields_ = [4,5];



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
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFeatureDetailsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    renameFeatureCommand: (f = msg.getRenameFeatureCommand()) && proto_feature_command_pb.RenameFeatureCommand.toObject(includeInstance, f),
    changeDescriptionCommand: (f = msg.getChangeDescriptionCommand()) && proto_feature_command_pb.ChangeDescriptionCommand.toObject(includeInstance, f),
    addTagCommandsList: jspb.Message.toObjectList(msg.getAddTagCommandsList(),
    proto_feature_command_pb.AddTagCommand.toObject, includeInstance),
    removeTagCommandsList: jspb.Message.toObjectList(msg.getRemoveTagCommandsList(),
    proto_feature_command_pb.RemoveTagCommand.toObject, includeInstance),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 6, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 7, "")
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
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFeatureDetailsRequest;
  return proto.bucketeer.feature.UpdateFeatureDetailsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.RenameFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.RenameFeatureCommand.deserializeBinaryFromReader);
      msg.setRenameFeatureCommand(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.ChangeDescriptionCommand;
      reader.readMessage(value,proto_feature_command_pb.ChangeDescriptionCommand.deserializeBinaryFromReader);
      msg.setChangeDescriptionCommand(value);
      break;
    case 4:
      var value = new proto_feature_command_pb.AddTagCommand;
      reader.readMessage(value,proto_feature_command_pb.AddTagCommand.deserializeBinaryFromReader);
      msg.addAddTagCommands(value);
      break;
    case 5:
      var value = new proto_feature_command_pb.RemoveTagCommand;
      reader.readMessage(value,proto_feature_command_pb.RemoveTagCommand.deserializeBinaryFromReader);
      msg.addRemoveTagCommands(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFeatureDetailsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRenameFeatureCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_command_pb.RenameFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getChangeDescriptionCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.ChangeDescriptionCommand.serializeBinaryToWriter
    );
  }
  f = message.getAddTagCommandsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto_feature_command_pb.AddTagCommand.serializeBinaryToWriter
    );
  }
  f = message.getRemoveTagCommandsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto_feature_command_pb.RemoveTagCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional RenameFeatureCommand rename_feature_command = 2;
 * @return {?proto.bucketeer.feature.RenameFeatureCommand}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getRenameFeatureCommand = function() {
  return /** @type{?proto.bucketeer.feature.RenameFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.RenameFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.RenameFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
*/
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setRenameFeatureCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.clearRenameFeatureCommand = function() {
  return this.setRenameFeatureCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.hasRenameFeatureCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional ChangeDescriptionCommand change_description_command = 3;
 * @return {?proto.bucketeer.feature.ChangeDescriptionCommand}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getChangeDescriptionCommand = function() {
  return /** @type{?proto.bucketeer.feature.ChangeDescriptionCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.ChangeDescriptionCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.ChangeDescriptionCommand|undefined} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
*/
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setChangeDescriptionCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.clearChangeDescriptionCommand = function() {
  return this.setChangeDescriptionCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.hasChangeDescriptionCommand = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * repeated AddTagCommand add_tag_commands = 4;
 * @return {!Array<!proto.bucketeer.feature.AddTagCommand>}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getAddTagCommandsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.AddTagCommand>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_command_pb.AddTagCommand, 4));
};


/**
 * @param {!Array<!proto.bucketeer.feature.AddTagCommand>} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
*/
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setAddTagCommandsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.bucketeer.feature.AddTagCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.AddTagCommand}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.addAddTagCommands = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.bucketeer.feature.AddTagCommand, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.clearAddTagCommandsList = function() {
  return this.setAddTagCommandsList([]);
};


/**
 * repeated RemoveTagCommand remove_tag_commands = 5;
 * @return {!Array<!proto.bucketeer.feature.RemoveTagCommand>}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getRemoveTagCommandsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.RemoveTagCommand>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_command_pb.RemoveTagCommand, 5));
};


/**
 * @param {!Array<!proto.bucketeer.feature.RemoveTagCommand>} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
*/
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setRemoveTagCommandsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.bucketeer.feature.RemoveTagCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.RemoveTagCommand}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.addRemoveTagCommands = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.bucketeer.feature.RemoveTagCommand, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.clearRemoveTagCommandsList = function() {
  return this.setRemoveTagCommandsList([]);
};


/**
 * optional string environment_namespace = 6;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string comment = 7;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureDetailsRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
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
proto.bucketeer.feature.UpdateFeatureDetailsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFeatureDetailsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFeatureDetailsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureDetailsResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsResponse}
 */
proto.bucketeer.feature.UpdateFeatureDetailsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFeatureDetailsResponse;
  return proto.bucketeer.feature.UpdateFeatureDetailsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFeatureDetailsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFeatureDetailsResponse}
 */
proto.bucketeer.feature.UpdateFeatureDetailsResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UpdateFeatureDetailsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFeatureDetailsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFeatureDetailsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureDetailsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.repeatedFields_ = [2];



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
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFeatureVariationsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    commandsList: jspb.Message.toObjectList(msg.getCommandsList(),
    proto_feature_command_pb.Command.toObject, includeInstance),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFeatureVariationsRequest;
  return proto.bucketeer.feature.UpdateFeatureVariationsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.Command;
      reader.readMessage(value,proto_feature_command_pb.Command.deserializeBinaryFromReader);
      msg.addCommands(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFeatureVariationsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommandsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto_feature_command_pb.Command.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated Command commands = 2;
 * @return {!Array<!proto.bucketeer.feature.Command>}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.getCommandsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Command>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_command_pb.Command, 2));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Command>} value
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} returns this
*/
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.setCommandsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.bucketeer.feature.Command=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Command}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.addCommands = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.bucketeer.feature.Command, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.clearCommandsList = function() {
  return this.setCommandsList([]);
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureVariationsRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.UpdateFeatureVariationsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFeatureVariationsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFeatureVariationsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureVariationsResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsResponse}
 */
proto.bucketeer.feature.UpdateFeatureVariationsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFeatureVariationsResponse;
  return proto.bucketeer.feature.UpdateFeatureVariationsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFeatureVariationsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFeatureVariationsResponse}
 */
proto.bucketeer.feature.UpdateFeatureVariationsResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UpdateFeatureVariationsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFeatureVariationsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFeatureVariationsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureVariationsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.repeatedFields_ = [2];



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
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFeatureTargetingRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    commandsList: jspb.Message.toObjectList(msg.getCommandsList(),
    proto_feature_command_pb.Command.toObject, includeInstance),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    comment: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFeatureTargetingRequest;
  return proto.bucketeer.feature.UpdateFeatureTargetingRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.Command;
      reader.readMessage(value,proto_feature_command_pb.Command.deserializeBinaryFromReader);
      msg.addCommands(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
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
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFeatureTargetingRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommandsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto_feature_command_pb.Command.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated Command commands = 2;
 * @return {!Array<!proto.bucketeer.feature.Command>}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.getCommandsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Command>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_command_pb.Command, 2));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Command>} value
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} returns this
*/
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.setCommandsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.bucketeer.feature.Command=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Command}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.addCommands = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.bucketeer.feature.Command, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.clearCommandsList = function() {
  return this.setCommandsList([]);
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string comment = 4;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingRequest} returns this
 */
proto.bucketeer.feature.UpdateFeatureTargetingRequest.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.UpdateFeatureTargetingResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFeatureTargetingResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFeatureTargetingResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureTargetingResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingResponse}
 */
proto.bucketeer.feature.UpdateFeatureTargetingResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFeatureTargetingResponse;
  return proto.bucketeer.feature.UpdateFeatureTargetingResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFeatureTargetingResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFeatureTargetingResponse}
 */
proto.bucketeer.feature.UpdateFeatureTargetingResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UpdateFeatureTargetingResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFeatureTargetingResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFeatureTargetingResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFeatureTargetingResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.CloneFeatureRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CloneFeatureRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CloneFeatureRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CloneFeatureRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.CloneFeatureCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.CloneFeatureRequest}
 */
proto.bucketeer.feature.CloneFeatureRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CloneFeatureRequest;
  return proto.bucketeer.feature.CloneFeatureRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CloneFeatureRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CloneFeatureRequest}
 */
proto.bucketeer.feature.CloneFeatureRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.CloneFeatureCommand;
      reader.readMessage(value,proto_feature_command_pb.CloneFeatureCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.CloneFeatureRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CloneFeatureRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CloneFeatureRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CloneFeatureRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.CloneFeatureCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CloneFeatureRequest} returns this
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional CloneFeatureCommand command = 2;
 * @return {?proto.bucketeer.feature.CloneFeatureCommand}
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.CloneFeatureCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.CloneFeatureCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.CloneFeatureCommand|undefined} value
 * @return {!proto.bucketeer.feature.CloneFeatureRequest} returns this
*/
proto.bucketeer.feature.CloneFeatureRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CloneFeatureRequest} returns this
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CloneFeatureRequest} returns this
 */
proto.bucketeer.feature.CloneFeatureRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
proto.bucketeer.feature.CloneFeatureResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CloneFeatureResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CloneFeatureResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CloneFeatureResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.CloneFeatureResponse}
 */
proto.bucketeer.feature.CloneFeatureResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CloneFeatureResponse;
  return proto.bucketeer.feature.CloneFeatureResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CloneFeatureResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CloneFeatureResponse}
 */
proto.bucketeer.feature.CloneFeatureResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.CloneFeatureResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CloneFeatureResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CloneFeatureResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CloneFeatureResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.CreateSegmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateSegmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateSegmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateSegmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && proto_feature_command_pb.CreateSegmentCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.CreateSegmentRequest}
 */
proto.bucketeer.feature.CreateSegmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateSegmentRequest;
  return proto.bucketeer.feature.CreateSegmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateSegmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateSegmentRequest}
 */
proto.bucketeer.feature.CreateSegmentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_command_pb.CreateSegmentCommand;
      reader.readMessage(value,proto_feature_command_pb.CreateSegmentCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.CreateSegmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateSegmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateSegmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateSegmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_command_pb.CreateSegmentCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional CreateSegmentCommand command = 1;
 * @return {?proto.bucketeer.feature.CreateSegmentCommand}
 */
proto.bucketeer.feature.CreateSegmentRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.CreateSegmentCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.CreateSegmentCommand, 1));
};


/**
 * @param {?proto.bucketeer.feature.CreateSegmentCommand|undefined} value
 * @return {!proto.bucketeer.feature.CreateSegmentRequest} returns this
*/
proto.bucketeer.feature.CreateSegmentRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateSegmentRequest} returns this
 */
proto.bucketeer.feature.CreateSegmentRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateSegmentRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.CreateSegmentRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateSegmentRequest} returns this
 */
proto.bucketeer.feature.CreateSegmentRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.CreateSegmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateSegmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateSegmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateSegmentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    segment: (f = msg.getSegment()) && proto_feature_segment_pb.Segment.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.CreateSegmentResponse}
 */
proto.bucketeer.feature.CreateSegmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateSegmentResponse;
  return proto.bucketeer.feature.CreateSegmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateSegmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateSegmentResponse}
 */
proto.bucketeer.feature.CreateSegmentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_segment_pb.Segment;
      reader.readMessage(value,proto_feature_segment_pb.Segment.deserializeBinaryFromReader);
      msg.setSegment(value);
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
proto.bucketeer.feature.CreateSegmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateSegmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateSegmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateSegmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_segment_pb.Segment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Segment segment = 1;
 * @return {?proto.bucketeer.feature.Segment}
 */
proto.bucketeer.feature.CreateSegmentResponse.prototype.getSegment = function() {
  return /** @type{?proto.bucketeer.feature.Segment} */ (
    jspb.Message.getWrapperField(this, proto_feature_segment_pb.Segment, 1));
};


/**
 * @param {?proto.bucketeer.feature.Segment|undefined} value
 * @return {!proto.bucketeer.feature.CreateSegmentResponse} returns this
*/
proto.bucketeer.feature.CreateSegmentResponse.prototype.setSegment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateSegmentResponse} returns this
 */
proto.bucketeer.feature.CreateSegmentResponse.prototype.clearSegment = function() {
  return this.setSegment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateSegmentResponse.prototype.hasSegment = function() {
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
proto.bucketeer.feature.GetSegmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetSegmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetSegmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.GetSegmentRequest}
 */
proto.bucketeer.feature.GetSegmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetSegmentRequest;
  return proto.bucketeer.feature.GetSegmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetSegmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetSegmentRequest}
 */
proto.bucketeer.feature.GetSegmentRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.GetSegmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetSegmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetSegmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.GetSegmentRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetSegmentRequest} returns this
 */
proto.bucketeer.feature.GetSegmentRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.GetSegmentRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetSegmentRequest} returns this
 */
proto.bucketeer.feature.GetSegmentRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.GetSegmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetSegmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetSegmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    segment: (f = msg.getSegment()) && proto_feature_segment_pb.Segment.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.GetSegmentResponse}
 */
proto.bucketeer.feature.GetSegmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetSegmentResponse;
  return proto.bucketeer.feature.GetSegmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetSegmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetSegmentResponse}
 */
proto.bucketeer.feature.GetSegmentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_segment_pb.Segment;
      reader.readMessage(value,proto_feature_segment_pb.Segment.deserializeBinaryFromReader);
      msg.setSegment(value);
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
proto.bucketeer.feature.GetSegmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetSegmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetSegmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_segment_pb.Segment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Segment segment = 1;
 * @return {?proto.bucketeer.feature.Segment}
 */
proto.bucketeer.feature.GetSegmentResponse.prototype.getSegment = function() {
  return /** @type{?proto.bucketeer.feature.Segment} */ (
    jspb.Message.getWrapperField(this, proto_feature_segment_pb.Segment, 1));
};


/**
 * @param {?proto.bucketeer.feature.Segment|undefined} value
 * @return {!proto.bucketeer.feature.GetSegmentResponse} returns this
*/
proto.bucketeer.feature.GetSegmentResponse.prototype.setSegment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.GetSegmentResponse} returns this
 */
proto.bucketeer.feature.GetSegmentResponse.prototype.clearSegment = function() {
  return this.setSegment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.GetSegmentResponse.prototype.hasSegment = function() {
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
proto.bucketeer.feature.ListSegmentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListSegmentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListSegmentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, ""),
    orderBy: jspb.Message.getFieldWithDefault(msg, 4, 0),
    orderDirection: jspb.Message.getFieldWithDefault(msg, 5, 0),
    searchKeyword: jspb.Message.getFieldWithDefault(msg, 6, ""),
    status: (f = msg.getStatus()) && google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
    isInUseStatus: (f = msg.getIsInUseStatus()) && google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ListSegmentsRequest}
 */
proto.bucketeer.feature.ListSegmentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListSegmentsRequest;
  return proto.bucketeer.feature.ListSegmentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListSegmentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListSegmentsRequest}
 */
proto.bucketeer.feature.ListSegmentsRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setEnvironmentNamespace(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.feature.ListSegmentsRequest.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 5:
      var value = /** @type {!proto.bucketeer.feature.ListSegmentsRequest.OrderDirection} */ (reader.readEnum());
      msg.setOrderDirection(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setSearchKeyword(value);
      break;
    case 7:
      var value = new google_protobuf_wrappers_pb.Int32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.Int32Value.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 8:
      var value = new google_protobuf_wrappers_pb.BoolValue;
      reader.readMessage(value,google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader);
      msg.setIsInUseStatus(value);
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
proto.bucketeer.feature.ListSegmentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListSegmentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListSegmentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentsRequest.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getEnvironmentNamespace();
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
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_wrappers_pb.Int32Value.serializeBinaryToWriter
    );
  }
  f = message.getIsInUseStatus();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.feature.ListSegmentsRequest.OrderBy = {
  DEFAULT: 0,
  NAME: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3
};

/**
 * @enum {number}
 */
proto.bucketeer.feature.ListSegmentsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.feature.ListSegmentsRequest.OrderBy}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.feature.ListSegmentsRequest.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListSegmentsRequest.OrderBy} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.feature.ListSegmentsRequest.OrderDirection}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.feature.ListSegmentsRequest.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListSegmentsRequest.OrderDirection} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional string search_keyword = 6;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getSearchKeyword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.setSearchKeyword = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional google.protobuf.Int32Value status = 7;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getStatus = function() {
  return /** @type{?proto.google.protobuf.Int32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.Int32Value, 7));
};


/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
*/
proto.bucketeer.feature.ListSegmentsRequest.prototype.setStatus = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.clearStatus = function() {
  return this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional google.protobuf.BoolValue is_in_use_status = 8;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.getIsInUseStatus = function() {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 8));
};


/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
*/
proto.bucketeer.feature.ListSegmentsRequest.prototype.setIsInUseStatus = function(value) {
  return jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListSegmentsRequest} returns this
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.clearIsInUseStatus = function() {
  return this.setIsInUseStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListSegmentsRequest.prototype.hasIsInUseStatus = function() {
  return jspb.Message.getField(this, 8) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListSegmentsResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.ListSegmentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListSegmentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListSegmentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentsList: jspb.Message.toObjectList(msg.getSegmentsList(),
    proto_feature_segment_pb.Segment.toObject, includeInstance),
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
 * @return {!proto.bucketeer.feature.ListSegmentsResponse}
 */
proto.bucketeer.feature.ListSegmentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListSegmentsResponse;
  return proto.bucketeer.feature.ListSegmentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListSegmentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListSegmentsResponse}
 */
proto.bucketeer.feature.ListSegmentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_segment_pb.Segment;
      reader.readMessage(value,proto_feature_segment_pb.Segment.deserializeBinaryFromReader);
      msg.addSegments(value);
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
proto.bucketeer.feature.ListSegmentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListSegmentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListSegmentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_feature_segment_pb.Segment.serializeBinaryToWriter
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
 * repeated Segment segments = 1;
 * @return {!Array<!proto.bucketeer.feature.Segment>}
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.getSegmentsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Segment>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_segment_pb.Segment, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Segment>} value
 * @return {!proto.bucketeer.feature.ListSegmentsResponse} returns this
*/
proto.bucketeer.feature.ListSegmentsResponse.prototype.setSegmentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.Segment=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Segment}
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.addSegments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.Segment, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListSegmentsResponse} returns this
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.clearSegmentsList = function() {
  return this.setSegmentsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentsResponse} returns this
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListSegmentsResponse} returns this
 */
proto.bucketeer.feature.ListSegmentsResponse.prototype.setTotalCount = function(value) {
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
proto.bucketeer.feature.DeleteSegmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteSegmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteSegmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.DeleteSegmentCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.DeleteSegmentRequest}
 */
proto.bucketeer.feature.DeleteSegmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteSegmentRequest;
  return proto.bucketeer.feature.DeleteSegmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteSegmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteSegmentRequest}
 */
proto.bucketeer.feature.DeleteSegmentRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.DeleteSegmentCommand;
      reader.readMessage(value,proto_feature_command_pb.DeleteSegmentCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.DeleteSegmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteSegmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteSegmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.DeleteSegmentCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteSegmentRequest} returns this
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DeleteSegmentCommand command = 2;
 * @return {?proto.bucketeer.feature.DeleteSegmentCommand}
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.DeleteSegmentCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.DeleteSegmentCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.DeleteSegmentCommand|undefined} value
 * @return {!proto.bucketeer.feature.DeleteSegmentRequest} returns this
*/
proto.bucketeer.feature.DeleteSegmentRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.DeleteSegmentRequest} returns this
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteSegmentRequest} returns this
 */
proto.bucketeer.feature.DeleteSegmentRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
proto.bucketeer.feature.DeleteSegmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteSegmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteSegmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteSegmentResponse}
 */
proto.bucketeer.feature.DeleteSegmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteSegmentResponse;
  return proto.bucketeer.feature.DeleteSegmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteSegmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteSegmentResponse}
 */
proto.bucketeer.feature.DeleteSegmentResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteSegmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteSegmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteSegmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.UpdateSegmentRequest.repeatedFields_ = [2];



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
proto.bucketeer.feature.UpdateSegmentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateSegmentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateSegmentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateSegmentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    commandsList: jspb.Message.toObjectList(msg.getCommandsList(),
    proto_feature_command_pb.Command.toObject, includeInstance),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.UpdateSegmentRequest}
 */
proto.bucketeer.feature.UpdateSegmentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateSegmentRequest;
  return proto.bucketeer.feature.UpdateSegmentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateSegmentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateSegmentRequest}
 */
proto.bucketeer.feature.UpdateSegmentRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.Command;
      reader.readMessage(value,proto_feature_command_pb.Command.deserializeBinaryFromReader);
      msg.addCommands(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.UpdateSegmentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateSegmentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateSegmentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateSegmentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCommandsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto_feature_command_pb.Command.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateSegmentRequest} returns this
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated Command commands = 2;
 * @return {!Array<!proto.bucketeer.feature.Command>}
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.getCommandsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Command>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_command_pb.Command, 2));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Command>} value
 * @return {!proto.bucketeer.feature.UpdateSegmentRequest} returns this
*/
proto.bucketeer.feature.UpdateSegmentRequest.prototype.setCommandsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.bucketeer.feature.Command=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Command}
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.addCommands = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.bucketeer.feature.Command, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.UpdateSegmentRequest} returns this
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.clearCommandsList = function() {
  return this.setCommandsList([]);
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateSegmentRequest} returns this
 */
proto.bucketeer.feature.UpdateSegmentRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
proto.bucketeer.feature.UpdateSegmentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateSegmentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateSegmentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateSegmentResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UpdateSegmentResponse}
 */
proto.bucketeer.feature.UpdateSegmentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateSegmentResponse;
  return proto.bucketeer.feature.UpdateSegmentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateSegmentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateSegmentResponse}
 */
proto.bucketeer.feature.UpdateSegmentResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UpdateSegmentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateSegmentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateSegmentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateSegmentResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.AddSegmentUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddSegmentUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddSegmentUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddSegmentUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.AddSegmentUserCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.AddSegmentUserRequest}
 */
proto.bucketeer.feature.AddSegmentUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddSegmentUserRequest;
  return proto.bucketeer.feature.AddSegmentUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddSegmentUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddSegmentUserRequest}
 */
proto.bucketeer.feature.AddSegmentUserRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.AddSegmentUserCommand;
      reader.readMessage(value,proto_feature_command_pb.AddSegmentUserCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.AddSegmentUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddSegmentUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddSegmentUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddSegmentUserRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.AddSegmentUserCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddSegmentUserRequest} returns this
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional AddSegmentUserCommand command = 2;
 * @return {?proto.bucketeer.feature.AddSegmentUserCommand}
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.AddSegmentUserCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.AddSegmentUserCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.AddSegmentUserCommand|undefined} value
 * @return {!proto.bucketeer.feature.AddSegmentUserRequest} returns this
*/
proto.bucketeer.feature.AddSegmentUserRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.AddSegmentUserRequest} returns this
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddSegmentUserRequest} returns this
 */
proto.bucketeer.feature.AddSegmentUserRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
proto.bucketeer.feature.AddSegmentUserResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddSegmentUserResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddSegmentUserResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddSegmentUserResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.AddSegmentUserResponse}
 */
proto.bucketeer.feature.AddSegmentUserResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddSegmentUserResponse;
  return proto.bucketeer.feature.AddSegmentUserResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddSegmentUserResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddSegmentUserResponse}
 */
proto.bucketeer.feature.AddSegmentUserResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.AddSegmentUserResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddSegmentUserResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddSegmentUserResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddSegmentUserResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteSegmentUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteSegmentUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.DeleteSegmentUserCommand.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.DeleteSegmentUserRequest}
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteSegmentUserRequest;
  return proto.bucketeer.feature.DeleteSegmentUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteSegmentUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteSegmentUserRequest}
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_command_pb.DeleteSegmentUserCommand;
      reader.readMessage(value,proto_feature_command_pb.DeleteSegmentUserCommand.deserializeBinaryFromReader);
      msg.setCommand(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteSegmentUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteSegmentUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.serializeBinaryToWriter = function(message, writer) {
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
      proto_feature_command_pb.DeleteSegmentUserCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteSegmentUserRequest} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional DeleteSegmentUserCommand command = 2;
 * @return {?proto.bucketeer.feature.DeleteSegmentUserCommand}
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.DeleteSegmentUserCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.DeleteSegmentUserCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.DeleteSegmentUserCommand|undefined} value
 * @return {!proto.bucketeer.feature.DeleteSegmentUserRequest} returns this
*/
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.DeleteSegmentUserRequest} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string environment_namespace = 3;
 * @return {string}
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteSegmentUserRequest} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
proto.bucketeer.feature.DeleteSegmentUserResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteSegmentUserResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteSegmentUserResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentUserResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteSegmentUserResponse}
 */
proto.bucketeer.feature.DeleteSegmentUserResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteSegmentUserResponse;
  return proto.bucketeer.feature.DeleteSegmentUserResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteSegmentUserResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteSegmentUserResponse}
 */
proto.bucketeer.feature.DeleteSegmentUserResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteSegmentUserResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteSegmentUserResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteSegmentUserResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentUserResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.GetSegmentUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetSegmentUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetSegmentUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    state: jspb.Message.getFieldWithDefault(msg, 3, 0),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.GetSegmentUserRequest}
 */
proto.bucketeer.feature.GetSegmentUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetSegmentUserRequest;
  return proto.bucketeer.feature.GetSegmentUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetSegmentUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetSegmentUserRequest}
 */
proto.bucketeer.feature.GetSegmentUserRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSegmentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.GetSegmentUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetSegmentUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetSegmentUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentUserRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetSegmentUserRequest} returns this
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string user_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetSegmentUserRequest} returns this
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.setUserId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional SegmentUser.State state = 3;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.feature.GetSegmentUserRequest} returns this
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional string environment_namespace = 4;
 * @return {string}
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetSegmentUserRequest} returns this
 */
proto.bucketeer.feature.GetSegmentUserRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.GetSegmentUserResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetSegmentUserResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetSegmentUserResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentUserResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    user: (f = msg.getUser()) && proto_feature_segment_pb.SegmentUser.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.GetSegmentUserResponse}
 */
proto.bucketeer.feature.GetSegmentUserResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetSegmentUserResponse;
  return proto.bucketeer.feature.GetSegmentUserResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetSegmentUserResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetSegmentUserResponse}
 */
proto.bucketeer.feature.GetSegmentUserResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_segment_pb.SegmentUser;
      reader.readMessage(value,proto_feature_segment_pb.SegmentUser.deserializeBinaryFromReader);
      msg.setUser(value);
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
proto.bucketeer.feature.GetSegmentUserResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetSegmentUserResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetSegmentUserResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetSegmentUserResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUser();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_segment_pb.SegmentUser.serializeBinaryToWriter
    );
  }
};


/**
 * optional SegmentUser user = 1;
 * @return {?proto.bucketeer.feature.SegmentUser}
 */
proto.bucketeer.feature.GetSegmentUserResponse.prototype.getUser = function() {
  return /** @type{?proto.bucketeer.feature.SegmentUser} */ (
    jspb.Message.getWrapperField(this, proto_feature_segment_pb.SegmentUser, 1));
};


/**
 * @param {?proto.bucketeer.feature.SegmentUser|undefined} value
 * @return {!proto.bucketeer.feature.GetSegmentUserResponse} returns this
*/
proto.bucketeer.feature.GetSegmentUserResponse.prototype.setUser = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.GetSegmentUserResponse} returns this
 */
proto.bucketeer.feature.GetSegmentUserResponse.prototype.clearUser = function() {
  return this.setUser(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.GetSegmentUserResponse.prototype.hasUser = function() {
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
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListSegmentUsersRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListSegmentUsersRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentUsersRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, ""),
    segmentId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    state: (f = msg.getState()) && google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
    userId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 6, "")
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
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListSegmentUsersRequest;
  return proto.bucketeer.feature.ListSegmentUsersRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListSegmentUsersRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setSegmentId(value);
      break;
    case 4:
      var value = new google_protobuf_wrappers_pb.Int32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.Int32Value.deserializeBinaryFromReader);
      msg.setState(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListSegmentUsersRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListSegmentUsersRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentUsersRequest.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getState();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_wrappers_pb.Int32Value.serializeBinaryToWriter
    );
  }
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
};


/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string segment_id = 3;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional google.protobuf.Int32Value state = 4;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.getState = function() {
  return /** @type{?proto.google.protobuf.Int32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.Int32Value, 4));
};


/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
*/
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.setState = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.clearState = function() {
  return this.setState(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.hasState = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string user_id = 5;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.setUserId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string environment_namespace = 6;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.ListSegmentUsersRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListSegmentUsersResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListSegmentUsersResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListSegmentUsersResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentUsersResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    usersList: jspb.Message.toObjectList(msg.getUsersList(),
    proto_feature_segment_pb.SegmentUser.toObject, includeInstance),
    cursor: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ListSegmentUsersResponse}
 */
proto.bucketeer.feature.ListSegmentUsersResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListSegmentUsersResponse;
  return proto.bucketeer.feature.ListSegmentUsersResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListSegmentUsersResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListSegmentUsersResponse}
 */
proto.bucketeer.feature.ListSegmentUsersResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_segment_pb.SegmentUser;
      reader.readMessage(value,proto_feature_segment_pb.SegmentUser.deserializeBinaryFromReader);
      msg.addUsers(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
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
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListSegmentUsersResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListSegmentUsersResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListSegmentUsersResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUsersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_feature_segment_pb.SegmentUser.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * repeated SegmentUser users = 1;
 * @return {!Array<!proto.bucketeer.feature.SegmentUser>}
 */
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.getUsersList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.SegmentUser>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_segment_pb.SegmentUser, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.SegmentUser>} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersResponse} returns this
*/
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.setUsersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.SegmentUser}
 */
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.addUsers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.SegmentUser, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListSegmentUsersResponse} returns this
 */
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.clearUsersList = function() {
  return this.setUsersList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListSegmentUsersResponse} returns this
 */
proto.bucketeer.feature.ListSegmentUsersResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.BulkUploadSegmentUsersRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 1, ""),
    segmentId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    command: (f = msg.getCommand()) && proto_feature_command_pb.BulkUploadSegmentUsersCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.BulkUploadSegmentUsersRequest;
  return proto.bucketeer.feature.BulkUploadSegmentUsersRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setSegmentId(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.BulkUploadSegmentUsersCommand;
      reader.readMessage(value,proto_feature_command_pb.BulkUploadSegmentUsersCommand.deserializeBinaryFromReader);
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
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.BulkUploadSegmentUsersRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.BulkUploadSegmentUsersCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string environment_namespace = 1;
 * @return {string}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string segment_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional BulkUploadSegmentUsersCommand command = 3;
 * @return {?proto.bucketeer.feature.BulkUploadSegmentUsersCommand}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.getCommand = function() {
  return /** @type{?proto.bucketeer.feature.BulkUploadSegmentUsersCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.BulkUploadSegmentUsersCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.BulkUploadSegmentUsersCommand|undefined} value
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} returns this
*/
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersRequest.prototype.hasCommand = function() {
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
proto.bucketeer.feature.BulkUploadSegmentUsersResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.BulkUploadSegmentUsersResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkUploadSegmentUsersResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersResponse}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.BulkUploadSegmentUsersResponse;
  return proto.bucketeer.feature.BulkUploadSegmentUsersResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersResponse}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.BulkUploadSegmentUsersResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.BulkUploadSegmentUsersResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkUploadSegmentUsersResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 1, ""),
    segmentId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    state: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.BulkDownloadSegmentUsersRequest;
  return proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setSegmentId(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (reader.readEnum());
      msg.setState(value);
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
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
};


/**
 * optional string environment_namespace = 1;
 * @return {string}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string segment_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional SegmentUser.State state = 3;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersRequest} returns this
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersRequest.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
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
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.BulkDownloadSegmentUsersResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    data: msg.getData_asB64()
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
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersResponse}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.BulkDownloadSegmentUsersResponse;
  return proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.BulkDownloadSegmentUsersResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersResponse}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
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
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.BulkDownloadSegmentUsersResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      1,
      f
    );
  }
};


/**
 * optional bytes data = 1;
 * @return {!(string|Uint8Array)}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.prototype.getData = function() {
  return /** @type {!(string|Uint8Array)} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * optional bytes data = 1;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes data = 1;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.bucketeer.feature.BulkDownloadSegmentUsersResponse} returns this
 */
proto.bucketeer.feature.BulkDownloadSegmentUsersResponse.prototype.setData = function(value) {
  return jspb.Message.setProto3BytesField(this, 1, value);
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
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EvaluateFeaturesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EvaluateFeaturesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    user: (f = msg.getUser()) && proto_user_user_pb.User.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    tag: jspb.Message.getFieldWithDefault(msg, 3, ""),
    featureId: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EvaluateFeaturesRequest;
  return proto.bucketeer.feature.EvaluateFeaturesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EvaluateFeaturesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_user_user_pb.User;
      reader.readMessage(value,proto_user_user_pb.User.deserializeBinaryFromReader);
      msg.setUser(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setTag(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setFeatureId(value);
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
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EvaluateFeaturesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EvaluateFeaturesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUser();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_user_user_pb.User.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getTag();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional bucketeer.user.User user = 1;
 * @return {?proto.bucketeer.user.User}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.getUser = function() {
  return /** @type{?proto.bucketeer.user.User} */ (
    jspb.Message.getWrapperField(this, proto_user_user_pb.User, 1));
};


/**
 * @param {?proto.bucketeer.user.User|undefined} value
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest} returns this
*/
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.setUser = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest} returns this
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.clearUser = function() {
  return this.setUser(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.hasUser = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest} returns this
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string tag = 3;
 * @return {string}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.getTag = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest} returns this
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.setTag = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string feature_id = 4;
 * @return {string}
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EvaluateFeaturesRequest} returns this
 */
proto.bucketeer.feature.EvaluateFeaturesRequest.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
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
proto.bucketeer.feature.EvaluateFeaturesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EvaluateFeaturesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EvaluateFeaturesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    userEvaluations: (f = msg.getUserEvaluations()) && proto_feature_evaluation_pb.UserEvaluations.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.EvaluateFeaturesResponse}
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EvaluateFeaturesResponse;
  return proto.bucketeer.feature.EvaluateFeaturesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EvaluateFeaturesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EvaluateFeaturesResponse}
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_evaluation_pb.UserEvaluations;
      reader.readMessage(value,proto_feature_evaluation_pb.UserEvaluations.deserializeBinaryFromReader);
      msg.setUserEvaluations(value);
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
proto.bucketeer.feature.EvaluateFeaturesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EvaluateFeaturesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EvaluateFeaturesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserEvaluations();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_evaluation_pb.UserEvaluations.serializeBinaryToWriter
    );
  }
};


/**
 * optional UserEvaluations user_evaluations = 1;
 * @return {?proto.bucketeer.feature.UserEvaluations}
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.prototype.getUserEvaluations = function() {
  return /** @type{?proto.bucketeer.feature.UserEvaluations} */ (
    jspb.Message.getWrapperField(this, proto_feature_evaluation_pb.UserEvaluations, 1));
};


/**
 * @param {?proto.bucketeer.feature.UserEvaluations|undefined} value
 * @return {!proto.bucketeer.feature.EvaluateFeaturesResponse} returns this
*/
proto.bucketeer.feature.EvaluateFeaturesResponse.prototype.setUserEvaluations = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.EvaluateFeaturesResponse} returns this
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.prototype.clearUserEvaluations = function() {
  return this.setUserEvaluations(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.EvaluateFeaturesResponse.prototype.hasUserEvaluations = function() {
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
proto.bucketeer.feature.ListTagsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListTagsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListTagsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListTagsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 1, ""),
    pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
    cursor: jspb.Message.getFieldWithDefault(msg, 3, ""),
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
 * @return {!proto.bucketeer.feature.ListTagsRequest}
 */
proto.bucketeer.feature.ListTagsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListTagsRequest;
  return proto.bucketeer.feature.ListTagsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListTagsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListTagsRequest}
 */
proto.bucketeer.feature.ListTagsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPageSize(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.feature.ListTagsRequest.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 5:
      var value = /** @type {!proto.bucketeer.feature.ListTagsRequest.OrderDirection} */ (reader.readEnum());
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
proto.bucketeer.feature.ListTagsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListTagsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListTagsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListTagsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getCursor();
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
proto.bucketeer.feature.ListTagsRequest.OrderBy = {
  DEFAULT: 0,
  ID: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3
};

/**
 * @enum {number}
 */
proto.bucketeer.feature.ListTagsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional string environment_namespace = 1;
 * @return {string}
 */
proto.bucketeer.feature.ListTagsRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListTagsRequest} returns this
 */
proto.bucketeer.feature.ListTagsRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 page_size = 2;
 * @return {number}
 */
proto.bucketeer.feature.ListTagsRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListTagsRequest} returns this
 */
proto.bucketeer.feature.ListTagsRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.feature.ListTagsRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListTagsRequest} returns this
 */
proto.bucketeer.feature.ListTagsRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.feature.ListTagsRequest.OrderBy}
 */
proto.bucketeer.feature.ListTagsRequest.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.feature.ListTagsRequest.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListTagsRequest.OrderBy} value
 * @return {!proto.bucketeer.feature.ListTagsRequest} returns this
 */
proto.bucketeer.feature.ListTagsRequest.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.feature.ListTagsRequest.OrderDirection}
 */
proto.bucketeer.feature.ListTagsRequest.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.feature.ListTagsRequest.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListTagsRequest.OrderDirection} value
 * @return {!proto.bucketeer.feature.ListTagsRequest} returns this
 */
proto.bucketeer.feature.ListTagsRequest.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional string search_keyword = 6;
 * @return {string}
 */
proto.bucketeer.feature.ListTagsRequest.prototype.getSearchKeyword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListTagsRequest} returns this
 */
proto.bucketeer.feature.ListTagsRequest.prototype.setSearchKeyword = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListTagsResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.ListTagsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListTagsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListTagsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListTagsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    tagsList: jspb.Message.toObjectList(msg.getTagsList(),
    proto_feature_feature_pb.Tag.toObject, includeInstance),
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
 * @return {!proto.bucketeer.feature.ListTagsResponse}
 */
proto.bucketeer.feature.ListTagsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListTagsResponse;
  return proto.bucketeer.feature.ListTagsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListTagsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListTagsResponse}
 */
proto.bucketeer.feature.ListTagsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_feature_pb.Tag;
      reader.readMessage(value,proto_feature_feature_pb.Tag.deserializeBinaryFromReader);
      msg.addTags(value);
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
proto.bucketeer.feature.ListTagsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListTagsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListTagsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListTagsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_feature_feature_pb.Tag.serializeBinaryToWriter
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
 * repeated Tag tags = 1;
 * @return {!Array<!proto.bucketeer.feature.Tag>}
 */
proto.bucketeer.feature.ListTagsResponse.prototype.getTagsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Tag>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_feature_pb.Tag, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Tag>} value
 * @return {!proto.bucketeer.feature.ListTagsResponse} returns this
*/
proto.bucketeer.feature.ListTagsResponse.prototype.setTagsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.Tag=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Tag}
 */
proto.bucketeer.feature.ListTagsResponse.prototype.addTags = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.Tag, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListTagsResponse} returns this
 */
proto.bucketeer.feature.ListTagsResponse.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListTagsResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListTagsResponse} returns this
 */
proto.bucketeer.feature.ListTagsResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.feature.ListTagsResponse.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListTagsResponse} returns this
 */
proto.bucketeer.feature.ListTagsResponse.prototype.setTotalCount = function(value) {
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
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 1, ""),
    createFlagTriggerCommand: (f = msg.getCreateFlagTriggerCommand()) && proto_feature_command_pb.CreateFlagTriggerCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.CreateFlagTriggerRequest}
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateFlagTriggerRequest;
  return proto.bucketeer.feature.CreateFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateFlagTriggerRequest}
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 2:
      var value = new proto_feature_command_pb.CreateFlagTriggerCommand;
      reader.readMessage(value,proto_feature_command_pb.CreateFlagTriggerCommand.deserializeBinaryFromReader);
      msg.setCreateFlagTriggerCommand(value);
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
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCreateFlagTriggerCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_command_pb.CreateFlagTriggerCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string environment_namespace = 1;
 * @return {string}
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional CreateFlagTriggerCommand create_flag_trigger_command = 2;
 * @return {?proto.bucketeer.feature.CreateFlagTriggerCommand}
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.getCreateFlagTriggerCommand = function() {
  return /** @type{?proto.bucketeer.feature.CreateFlagTriggerCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.CreateFlagTriggerCommand, 2));
};


/**
 * @param {?proto.bucketeer.feature.CreateFlagTriggerCommand|undefined} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerRequest} returns this
*/
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.setCreateFlagTriggerCommand = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.clearCreateFlagTriggerCommand = function() {
  return this.setCreateFlagTriggerCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateFlagTriggerRequest.prototype.hasCreateFlagTriggerCommand = function() {
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
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    flagTrigger: (f = msg.getFlagTrigger()) && proto_feature_flag_trigger_pb.FlagTrigger.toObject(includeInstance, f),
    url: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.CreateFlagTriggerResponse}
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateFlagTriggerResponse;
  return proto.bucketeer.feature.CreateFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateFlagTriggerResponse}
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_flag_trigger_pb.FlagTrigger;
      reader.readMessage(value,proto_feature_flag_trigger_pb.FlagTrigger.deserializeBinaryFromReader);
      msg.setFlagTrigger(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrl(value);
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
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFlagTrigger();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_flag_trigger_pb.FlagTrigger.serializeBinaryToWriter
    );
  }
  f = message.getUrl();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional FlagTrigger flag_trigger = 1;
 * @return {?proto.bucketeer.feature.FlagTrigger}
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.getFlagTrigger = function() {
  return /** @type{?proto.bucketeer.feature.FlagTrigger} */ (
    jspb.Message.getWrapperField(this, proto_feature_flag_trigger_pb.FlagTrigger, 1));
};


/**
 * @param {?proto.bucketeer.feature.FlagTrigger|undefined} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerResponse} returns this
*/
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.setFlagTrigger = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateFlagTriggerResponse} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.clearFlagTrigger = function() {
  return this.setFlagTrigger(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.hasFlagTrigger = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string url = 2;
 * @return {string}
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.getUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerResponse} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerResponse.prototype.setUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    deleteFlagTriggerCommand: (f = msg.getDeleteFlagTriggerCommand()) && proto_feature_command_pb.DeleteFlagTriggerCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerRequest}
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteFlagTriggerRequest;
  return proto.bucketeer.feature.DeleteFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerRequest}
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.DeleteFlagTriggerCommand;
      reader.readMessage(value,proto_feature_command_pb.DeleteFlagTriggerCommand.deserializeBinaryFromReader);
      msg.setDeleteFlagTriggerCommand(value);
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
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getDeleteFlagTriggerCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.DeleteFlagTriggerCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional DeleteFlagTriggerCommand delete_flag_trigger_command = 3;
 * @return {?proto.bucketeer.feature.DeleteFlagTriggerCommand}
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.getDeleteFlagTriggerCommand = function() {
  return /** @type{?proto.bucketeer.feature.DeleteFlagTriggerCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.DeleteFlagTriggerCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.DeleteFlagTriggerCommand|undefined} value
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerRequest} returns this
*/
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.setDeleteFlagTriggerCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.clearDeleteFlagTriggerCommand = function() {
  return this.setDeleteFlagTriggerCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.DeleteFlagTriggerRequest.prototype.hasDeleteFlagTriggerCommand = function() {
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
proto.bucketeer.feature.DeleteFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFlagTriggerResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerResponse}
 */
proto.bucketeer.feature.DeleteFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteFlagTriggerResponse;
  return proto.bucketeer.feature.DeleteFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerResponse}
 */
proto.bucketeer.feature.DeleteFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    changeFlagTriggerDescriptionCommand: (f = msg.getChangeFlagTriggerDescriptionCommand()) && proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerRequest}
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFlagTriggerRequest;
  return proto.bucketeer.feature.UpdateFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerRequest}
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand;
      reader.readMessage(value,proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand.deserializeBinaryFromReader);
      msg.setChangeFlagTriggerDescriptionCommand(value);
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
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getChangeFlagTriggerDescriptionCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional ChangeFlagTriggerDescriptionCommand change_flag_trigger_description_command = 3;
 * @return {?proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand}
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.getChangeFlagTriggerDescriptionCommand = function() {
  return /** @type{?proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.ChangeFlagTriggerDescriptionCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand|undefined} value
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerRequest} returns this
*/
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.setChangeFlagTriggerDescriptionCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.clearChangeFlagTriggerDescriptionCommand = function() {
  return this.setChangeFlagTriggerDescriptionCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.UpdateFlagTriggerRequest.prototype.hasChangeFlagTriggerDescriptionCommand = function() {
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
proto.bucketeer.feature.UpdateFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFlagTriggerResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerResponse}
 */
proto.bucketeer.feature.UpdateFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFlagTriggerResponse;
  return proto.bucketeer.feature.UpdateFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerResponse}
 */
proto.bucketeer.feature.UpdateFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UpdateFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EnableFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EnableFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    enableFlagTriggerCommand: (f = msg.getEnableFlagTriggerCommand()) && proto_feature_command_pb.EnableFlagTriggerCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.EnableFlagTriggerRequest}
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EnableFlagTriggerRequest;
  return proto.bucketeer.feature.EnableFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EnableFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EnableFlagTriggerRequest}
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.EnableFlagTriggerCommand;
      reader.readMessage(value,proto_feature_command_pb.EnableFlagTriggerCommand.deserializeBinaryFromReader);
      msg.setEnableFlagTriggerCommand(value);
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
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EnableFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EnableFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getEnableFlagTriggerCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.EnableFlagTriggerCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EnableFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.EnableFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional EnableFlagTriggerCommand enable_flag_trigger_command = 3;
 * @return {?proto.bucketeer.feature.EnableFlagTriggerCommand}
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.getEnableFlagTriggerCommand = function() {
  return /** @type{?proto.bucketeer.feature.EnableFlagTriggerCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.EnableFlagTriggerCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.EnableFlagTriggerCommand|undefined} value
 * @return {!proto.bucketeer.feature.EnableFlagTriggerRequest} returns this
*/
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.setEnableFlagTriggerCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.EnableFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.clearEnableFlagTriggerCommand = function() {
  return this.setEnableFlagTriggerCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.EnableFlagTriggerRequest.prototype.hasEnableFlagTriggerCommand = function() {
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
proto.bucketeer.feature.EnableFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EnableFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EnableFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFlagTriggerResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.EnableFlagTriggerResponse}
 */
proto.bucketeer.feature.EnableFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EnableFlagTriggerResponse;
  return proto.bucketeer.feature.EnableFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EnableFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EnableFlagTriggerResponse}
 */
proto.bucketeer.feature.EnableFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.EnableFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EnableFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EnableFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DisableFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DisableFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    disableFlagTriggerCommand: (f = msg.getDisableFlagTriggerCommand()) && proto_feature_command_pb.DisableFlagTriggerCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.DisableFlagTriggerRequest}
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DisableFlagTriggerRequest;
  return proto.bucketeer.feature.DisableFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DisableFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DisableFlagTriggerRequest}
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.DisableFlagTriggerCommand;
      reader.readMessage(value,proto_feature_command_pb.DisableFlagTriggerCommand.deserializeBinaryFromReader);
      msg.setDisableFlagTriggerCommand(value);
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
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DisableFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DisableFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getDisableFlagTriggerCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.DisableFlagTriggerCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DisableFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DisableFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional DisableFlagTriggerCommand disable_flag_trigger_command = 3;
 * @return {?proto.bucketeer.feature.DisableFlagTriggerCommand}
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.getDisableFlagTriggerCommand = function() {
  return /** @type{?proto.bucketeer.feature.DisableFlagTriggerCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.DisableFlagTriggerCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.DisableFlagTriggerCommand|undefined} value
 * @return {!proto.bucketeer.feature.DisableFlagTriggerRequest} returns this
*/
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.setDisableFlagTriggerCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.DisableFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.clearDisableFlagTriggerCommand = function() {
  return this.setDisableFlagTriggerCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.DisableFlagTriggerRequest.prototype.hasDisableFlagTriggerCommand = function() {
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
proto.bucketeer.feature.DisableFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DisableFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DisableFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFlagTriggerResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DisableFlagTriggerResponse}
 */
proto.bucketeer.feature.DisableFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DisableFlagTriggerResponse;
  return proto.bucketeer.feature.DisableFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DisableFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DisableFlagTriggerResponse}
 */
proto.bucketeer.feature.DisableFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DisableFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DisableFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DisableFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ResetFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ResetFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    resetFlagTriggerCommand: (f = msg.getResetFlagTriggerCommand()) && proto_feature_command_pb.ResetFlagTriggerCommand.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ResetFlagTriggerRequest}
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ResetFlagTriggerRequest;
  return proto.bucketeer.feature.ResetFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ResetFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ResetFlagTriggerRequest}
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = new proto_feature_command_pb.ResetFlagTriggerCommand;
      reader.readMessage(value,proto_feature_command_pb.ResetFlagTriggerCommand.deserializeBinaryFromReader);
      msg.setResetFlagTriggerCommand(value);
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
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ResetFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ResetFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getResetFlagTriggerCommand();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_command_pb.ResetFlagTriggerCommand.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ResetFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ResetFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional ResetFlagTriggerCommand reset_flag_trigger_command = 3;
 * @return {?proto.bucketeer.feature.ResetFlagTriggerCommand}
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.getResetFlagTriggerCommand = function() {
  return /** @type{?proto.bucketeer.feature.ResetFlagTriggerCommand} */ (
    jspb.Message.getWrapperField(this, proto_feature_command_pb.ResetFlagTriggerCommand, 3));
};


/**
 * @param {?proto.bucketeer.feature.ResetFlagTriggerCommand|undefined} value
 * @return {!proto.bucketeer.feature.ResetFlagTriggerRequest} returns this
*/
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.setResetFlagTriggerCommand = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ResetFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.clearResetFlagTriggerCommand = function() {
  return this.setResetFlagTriggerCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ResetFlagTriggerRequest.prototype.hasResetFlagTriggerCommand = function() {
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
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ResetFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ResetFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    flagTrigger: (f = msg.getFlagTrigger()) && proto_feature_flag_trigger_pb.FlagTrigger.toObject(includeInstance, f),
    url: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ResetFlagTriggerResponse}
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ResetFlagTriggerResponse;
  return proto.bucketeer.feature.ResetFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ResetFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ResetFlagTriggerResponse}
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_flag_trigger_pb.FlagTrigger;
      reader.readMessage(value,proto_feature_flag_trigger_pb.FlagTrigger.deserializeBinaryFromReader);
      msg.setFlagTrigger(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrl(value);
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
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ResetFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ResetFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFlagTrigger();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_flag_trigger_pb.FlagTrigger.serializeBinaryToWriter
    );
  }
  f = message.getUrl();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional FlagTrigger flag_trigger = 1;
 * @return {?proto.bucketeer.feature.FlagTrigger}
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.getFlagTrigger = function() {
  return /** @type{?proto.bucketeer.feature.FlagTrigger} */ (
    jspb.Message.getWrapperField(this, proto_feature_flag_trigger_pb.FlagTrigger, 1));
};


/**
 * @param {?proto.bucketeer.feature.FlagTrigger|undefined} value
 * @return {!proto.bucketeer.feature.ResetFlagTriggerResponse} returns this
*/
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.setFlagTrigger = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ResetFlagTriggerResponse} returns this
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.clearFlagTrigger = function() {
  return this.setFlagTrigger(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.hasFlagTrigger = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string url = 2;
 * @return {string}
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.getUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ResetFlagTriggerResponse} returns this
 */
proto.bucketeer.feature.ResetFlagTriggerResponse.prototype.setUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.GetFlagTriggerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetFlagTriggerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetFlagTriggerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFlagTriggerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.GetFlagTriggerRequest}
 */
proto.bucketeer.feature.GetFlagTriggerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetFlagTriggerRequest;
  return proto.bucketeer.feature.GetFlagTriggerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetFlagTriggerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetFlagTriggerRequest}
 */
proto.bucketeer.feature.GetFlagTriggerRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
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
proto.bucketeer.feature.GetFlagTriggerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetFlagTriggerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetFlagTriggerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFlagTriggerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.GetFlagTriggerRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.GetFlagTriggerRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.GetFlagTriggerRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetFlagTriggerRequest} returns this
 */
proto.bucketeer.feature.GetFlagTriggerRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.GetFlagTriggerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.GetFlagTriggerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFlagTriggerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    flagTrigger: (f = msg.getFlagTrigger()) && proto_feature_flag_trigger_pb.FlagTrigger.toObject(includeInstance, f),
    url: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.GetFlagTriggerResponse}
 */
proto.bucketeer.feature.GetFlagTriggerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.GetFlagTriggerResponse;
  return proto.bucketeer.feature.GetFlagTriggerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.GetFlagTriggerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.GetFlagTriggerResponse}
 */
proto.bucketeer.feature.GetFlagTriggerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_flag_trigger_pb.FlagTrigger;
      reader.readMessage(value,proto_feature_flag_trigger_pb.FlagTrigger.deserializeBinaryFromReader);
      msg.setFlagTrigger(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrl(value);
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
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.GetFlagTriggerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.GetFlagTriggerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.GetFlagTriggerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFlagTrigger();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_flag_trigger_pb.FlagTrigger.serializeBinaryToWriter
    );
  }
  f = message.getUrl();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional FlagTrigger flag_trigger = 1;
 * @return {?proto.bucketeer.feature.FlagTrigger}
 */
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.getFlagTrigger = function() {
  return /** @type{?proto.bucketeer.feature.FlagTrigger} */ (
    jspb.Message.getWrapperField(this, proto_feature_flag_trigger_pb.FlagTrigger, 1));
};


/**
 * @param {?proto.bucketeer.feature.FlagTrigger|undefined} value
 * @return {!proto.bucketeer.feature.GetFlagTriggerResponse} returns this
*/
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.setFlagTrigger = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.GetFlagTriggerResponse} returns this
 */
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.clearFlagTrigger = function() {
  return this.setFlagTrigger(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.hasFlagTrigger = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string url = 2;
 * @return {string}
 */
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.getUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.GetFlagTriggerResponse} returns this
 */
proto.bucketeer.feature.GetFlagTriggerResponse.prototype.setUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
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
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListFlagTriggersRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListFlagTriggersRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFlagTriggersRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    cursor: jspb.Message.getFieldWithDefault(msg, 3, ""),
    pageSize: jspb.Message.getFieldWithDefault(msg, 4, 0),
    orderBy: jspb.Message.getFieldWithDefault(msg, 5, 0),
    orderDirection: jspb.Message.getFieldWithDefault(msg, 6, 0)
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
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListFlagTriggersRequest;
  return proto.bucketeer.feature.ListFlagTriggersRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListFlagTriggersRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setFeatureId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCursor(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPageSize(value);
      break;
    case 5:
      var value = /** @type {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderBy} */ (reader.readEnum());
      msg.setOrderBy(value);
      break;
    case 6:
      var value = /** @type {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderDirection} */ (reader.readEnum());
      msg.setOrderDirection(value);
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
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListFlagTriggersRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListFlagTriggersRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFlagTriggersRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt32(
      4,
      f
    );
  }
  f = message.getOrderBy();
  if (f !== 0.0) {
    writer.writeEnum(
      5,
      f
    );
  }
  f = message.getOrderDirection();
  if (f !== 0.0) {
    writer.writeEnum(
      6,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.OrderBy = {
  DEFAULT: 0,
  CREATED_AT: 1,
  UPDATED_AT: 2
};

/**
 * @enum {number}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest} returns this
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string environment_namespace = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest} returns this
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest} returns this
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional int32 page_size = 4;
 * @return {number}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.getPageSize = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest} returns this
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.setPageSize = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional OrderBy order_by = 5;
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderBy}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.getOrderBy = function() {
  return /** @type {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderBy} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderBy} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest} returns this
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.setOrderBy = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional OrderDirection order_direction = 6;
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderDirection}
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.getOrderDirection = function() {
  return /** @type {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderDirection} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {!proto.bucketeer.feature.ListFlagTriggersRequest.OrderDirection} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersRequest} returns this
 */
proto.bucketeer.feature.ListFlagTriggersRequest.prototype.setOrderDirection = function(value) {
  return jspb.Message.setProto3EnumField(this, 6, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ListFlagTriggersResponse.repeatedFields_ = [1];



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
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListFlagTriggersResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFlagTriggersResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    flagTriggersList: jspb.Message.toObjectList(msg.getFlagTriggersList(),
    proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.toObject, includeInstance),
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
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListFlagTriggersResponse;
  return proto.bucketeer.feature.ListFlagTriggersResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl;
      reader.readMessage(value,proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.deserializeBinaryFromReader);
      msg.addFlagTriggers(value);
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
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListFlagTriggersResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFlagTriggersResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFlagTriggersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.serializeBinaryToWriter
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
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.toObject = function(includeInstance, msg) {
  var f, obj = {
    flagTrigger: (f = msg.getFlagTrigger()) && proto_feature_flag_trigger_pb.FlagTrigger.toObject(includeInstance, f),
    url: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl;
  return proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_flag_trigger_pb.FlagTrigger;
      reader.readMessage(value,proto_feature_flag_trigger_pb.FlagTrigger.deserializeBinaryFromReader);
      msg.setFlagTrigger(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrl(value);
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
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFlagTrigger();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_flag_trigger_pb.FlagTrigger.serializeBinaryToWriter
    );
  }
  f = message.getUrl();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional FlagTrigger flag_trigger = 1;
 * @return {?proto.bucketeer.feature.FlagTrigger}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.getFlagTrigger = function() {
  return /** @type{?proto.bucketeer.feature.FlagTrigger} */ (
    jspb.Message.getWrapperField(this, proto_feature_flag_trigger_pb.FlagTrigger, 1));
};


/**
 * @param {?proto.bucketeer.feature.FlagTrigger|undefined} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl} returns this
*/
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.setFlagTrigger = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl} returns this
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.clearFlagTrigger = function() {
  return this.setFlagTrigger(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.hasFlagTrigger = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string url = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.getUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl} returns this
 */
proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl.prototype.setUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated FlagTriggerWithUrl flag_triggers = 1;
 * @return {!Array<!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl>}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.getFlagTriggersList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl, 1));
};


/**
 * @param {!Array<!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl>} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse} returns this
*/
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.setFlagTriggersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.addFlagTriggers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.bucketeer.feature.ListFlagTriggersResponse.FlagTriggerWithUrl, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse} returns this
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.clearFlagTriggersList = function() {
  return this.setFlagTriggersList([]);
};


/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.getCursor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse} returns this
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.setCursor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.getTotalCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ListFlagTriggersResponse} returns this
 */
proto.bucketeer.feature.ListFlagTriggersResponse.prototype.setTotalCount = function(value) {
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
proto.bucketeer.feature.FlagTriggerWebhookRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.FlagTriggerWebhookRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.FlagTriggerWebhookRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.FlagTriggerWebhookRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    token: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.FlagTriggerWebhookRequest}
 */
proto.bucketeer.feature.FlagTriggerWebhookRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.FlagTriggerWebhookRequest;
  return proto.bucketeer.feature.FlagTriggerWebhookRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.FlagTriggerWebhookRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.FlagTriggerWebhookRequest}
 */
proto.bucketeer.feature.FlagTriggerWebhookRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setToken(value);
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
proto.bucketeer.feature.FlagTriggerWebhookRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.FlagTriggerWebhookRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.FlagTriggerWebhookRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.FlagTriggerWebhookRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getToken();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string token = 1;
 * @return {string}
 */
proto.bucketeer.feature.FlagTriggerWebhookRequest.prototype.getToken = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.FlagTriggerWebhookRequest} returns this
 */
proto.bucketeer.feature.FlagTriggerWebhookRequest.prototype.setToken = function(value) {
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
proto.bucketeer.feature.FlagTriggerWebhookResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.FlagTriggerWebhookResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.FlagTriggerWebhookResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.FlagTriggerWebhookResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.FlagTriggerWebhookResponse}
 */
proto.bucketeer.feature.FlagTriggerWebhookResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.FlagTriggerWebhookResponse;
  return proto.bucketeer.feature.FlagTriggerWebhookResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.FlagTriggerWebhookResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.FlagTriggerWebhookResponse}
 */
proto.bucketeer.feature.FlagTriggerWebhookResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.FlagTriggerWebhookResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.FlagTriggerWebhookResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.FlagTriggerWebhookResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.FlagTriggerWebhookResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};


goog.object.extend(exports, proto.bucketeer.feature);
