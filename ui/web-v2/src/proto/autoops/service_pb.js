// source: proto/autoops/service.proto
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
  function () {
    return this;
  }.call(null) ||
  Function('return this')();

var google_api_annotations_pb = require('../../google/api/annotations_pb.js');
goog.object.extend(proto, google_api_annotations_pb);
var google_api_field_behavior_pb = require('../../google/api/field_behavior_pb.js');
goog.object.extend(proto, google_api_field_behavior_pb);
var protoc$gen$openapiv2_options_annotations_pb = require('../../protoc-gen-openapiv2/options/annotations_pb.js');
goog.object.extend(proto, protoc$gen$openapiv2_options_annotations_pb);
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.object.extend(proto, google_protobuf_wrappers_pb);
var proto_autoops_auto_ops_rule_pb = require('../../proto/autoops/auto_ops_rule_pb.js');
goog.object.extend(proto, proto_autoops_auto_ops_rule_pb);
var proto_autoops_clause_pb = require('../../proto/autoops/clause_pb.js');
goog.object.extend(proto, proto_autoops_clause_pb);
var proto_autoops_command_pb = require('../../proto/autoops/command_pb.js');
goog.object.extend(proto, proto_autoops_command_pb);
var proto_autoops_ops_count_pb = require('../../proto/autoops/ops_count_pb.js');
goog.object.extend(proto, proto_autoops_ops_count_pb);
var proto_autoops_progressive_rollout_pb = require('../../proto/autoops/progressive_rollout_pb.js');
goog.object.extend(proto, proto_autoops_progressive_rollout_pb);
goog.exportSymbol('proto.bucketeer.autoops.ChangeType', null, global);
goog.exportSymbol(
  'proto.bucketeer.autoops.CreateAutoOpsRuleRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.CreateAutoOpsRuleResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.CreateProgressiveRolloutRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.CreateProgressiveRolloutResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.autoops.DatetimeClauseChange', null, global);
goog.exportSymbol(
  'proto.bucketeer.autoops.DeleteAutoOpsRuleRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.DeleteAutoOpsRuleResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.DeleteProgressiveRolloutRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.DeleteProgressiveRolloutResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ExecuteAutoOpsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ExecuteAutoOpsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.GetAutoOpsRuleRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.GetAutoOpsRuleResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.GetProgressiveRolloutRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.GetProgressiveRolloutResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListAutoOpsRulesRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListAutoOpsRulesResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.autoops.ListOpsCountsRequest', null, global);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListOpsCountsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListProgressiveRolloutsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderBy',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirection',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ListProgressiveRolloutsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.OpsEventRateClauseChange',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.StopAutoOpsRuleRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.StopAutoOpsRuleResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.StopProgressiveRolloutRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.StopProgressiveRolloutResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.UpdateAutoOpsRuleRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.UpdateAutoOpsRuleResponse',
  null,
  global
);
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
proto.bucketeer.autoops.GetAutoOpsRuleRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.GetAutoOpsRuleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.GetAutoOpsRuleRequest.displayName =
    'proto.bucketeer.autoops.GetAutoOpsRuleRequest';
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
proto.bucketeer.autoops.GetAutoOpsRuleResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.GetAutoOpsRuleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.GetAutoOpsRuleResponse.displayName =
    'proto.bucketeer.autoops.GetAutoOpsRuleResponse';
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
proto.bucketeer.autoops.CreateAutoOpsRuleRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.CreateAutoOpsRuleRequest.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.CreateAutoOpsRuleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.CreateAutoOpsRuleRequest.displayName =
    'proto.bucketeer.autoops.CreateAutoOpsRuleRequest';
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
proto.bucketeer.autoops.CreateAutoOpsRuleResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.CreateAutoOpsRuleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.CreateAutoOpsRuleResponse.displayName =
    'proto.bucketeer.autoops.CreateAutoOpsRuleResponse';
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
proto.bucketeer.autoops.ListAutoOpsRulesRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.ListAutoOpsRulesRequest.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.ListAutoOpsRulesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ListAutoOpsRulesRequest.displayName =
    'proto.bucketeer.autoops.ListAutoOpsRulesRequest';
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
proto.bucketeer.autoops.ListAutoOpsRulesResponse = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.ListAutoOpsRulesResponse.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.ListAutoOpsRulesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ListAutoOpsRulesResponse.displayName =
    'proto.bucketeer.autoops.ListAutoOpsRulesResponse';
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
proto.bucketeer.autoops.StopAutoOpsRuleRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.StopAutoOpsRuleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.StopAutoOpsRuleRequest.displayName =
    'proto.bucketeer.autoops.StopAutoOpsRuleRequest';
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
proto.bucketeer.autoops.StopAutoOpsRuleResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.StopAutoOpsRuleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.StopAutoOpsRuleResponse.displayName =
    'proto.bucketeer.autoops.StopAutoOpsRuleResponse';
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
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.DeleteAutoOpsRuleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.displayName =
    'proto.bucketeer.autoops.DeleteAutoOpsRuleRequest';
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
proto.bucketeer.autoops.DeleteAutoOpsRuleResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.DeleteAutoOpsRuleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.displayName =
    'proto.bucketeer.autoops.DeleteAutoOpsRuleResponse';
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
proto.bucketeer.autoops.OpsEventRateClauseChange = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.OpsEventRateClauseChange, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.OpsEventRateClauseChange.displayName =
    'proto.bucketeer.autoops.OpsEventRateClauseChange';
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
proto.bucketeer.autoops.DatetimeClauseChange = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.DatetimeClauseChange, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.DatetimeClauseChange.displayName =
    'proto.bucketeer.autoops.DatetimeClauseChange';
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
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.UpdateAutoOpsRuleRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.displayName =
    'proto.bucketeer.autoops.UpdateAutoOpsRuleRequest';
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
proto.bucketeer.autoops.UpdateAutoOpsRuleResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.UpdateAutoOpsRuleResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.displayName =
    'proto.bucketeer.autoops.UpdateAutoOpsRuleResponse';
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
proto.bucketeer.autoops.ExecuteAutoOpsRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.ExecuteAutoOpsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ExecuteAutoOpsRequest.displayName =
    'proto.bucketeer.autoops.ExecuteAutoOpsRequest';
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
proto.bucketeer.autoops.ExecuteAutoOpsResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.ExecuteAutoOpsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ExecuteAutoOpsResponse.displayName =
    'proto.bucketeer.autoops.ExecuteAutoOpsResponse';
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
proto.bucketeer.autoops.ListOpsCountsRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.ListOpsCountsRequest.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.ListOpsCountsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ListOpsCountsRequest.displayName =
    'proto.bucketeer.autoops.ListOpsCountsRequest';
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
proto.bucketeer.autoops.ListOpsCountsResponse = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.ListOpsCountsResponse.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.ListOpsCountsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ListOpsCountsResponse.displayName =
    'proto.bucketeer.autoops.ListOpsCountsResponse';
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
proto.bucketeer.autoops.CreateProgressiveRolloutRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.CreateProgressiveRolloutRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.CreateProgressiveRolloutRequest.displayName =
    'proto.bucketeer.autoops.CreateProgressiveRolloutRequest';
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
proto.bucketeer.autoops.CreateProgressiveRolloutResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.CreateProgressiveRolloutResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.CreateProgressiveRolloutResponse.displayName =
    'proto.bucketeer.autoops.CreateProgressiveRolloutResponse';
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
proto.bucketeer.autoops.GetProgressiveRolloutRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.GetProgressiveRolloutRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.GetProgressiveRolloutRequest.displayName =
    'proto.bucketeer.autoops.GetProgressiveRolloutRequest';
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
proto.bucketeer.autoops.GetProgressiveRolloutResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.GetProgressiveRolloutResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.GetProgressiveRolloutResponse.displayName =
    'proto.bucketeer.autoops.GetProgressiveRolloutResponse';
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
proto.bucketeer.autoops.StopProgressiveRolloutRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.StopProgressiveRolloutRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.StopProgressiveRolloutRequest.displayName =
    'proto.bucketeer.autoops.StopProgressiveRolloutRequest';
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
proto.bucketeer.autoops.StopProgressiveRolloutResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.StopProgressiveRolloutResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.StopProgressiveRolloutResponse.displayName =
    'proto.bucketeer.autoops.StopProgressiveRolloutResponse';
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
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.DeleteProgressiveRolloutRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.displayName =
    'proto.bucketeer.autoops.DeleteProgressiveRolloutRequest';
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
proto.bucketeer.autoops.DeleteProgressiveRolloutResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.DeleteProgressiveRolloutResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.displayName =
    'proto.bucketeer.autoops.DeleteProgressiveRolloutResponse';
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
proto.bucketeer.autoops.ListProgressiveRolloutsRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.ListProgressiveRolloutsRequest.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.autoops.ListProgressiveRolloutsRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ListProgressiveRolloutsRequest.displayName =
    'proto.bucketeer.autoops.ListProgressiveRolloutsRequest';
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
proto.bucketeer.autoops.ListProgressiveRolloutsResponse = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.ListProgressiveRolloutsResponse.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.autoops.ListProgressiveRolloutsResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ListProgressiveRolloutsResponse.displayName =
    'proto.bucketeer.autoops.ListProgressiveRolloutsResponse';
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
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.displayName =
    'proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest';
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
proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.displayName =
    'proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse';
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
  proto.bucketeer.autoops.GetAutoOpsRuleRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.GetAutoOpsRuleRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.GetAutoOpsRuleRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.GetAutoOpsRuleRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, '')
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
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.GetAutoOpsRuleRequest();
  return proto.bucketeer.autoops.GetAutoOpsRuleRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.GetAutoOpsRuleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.GetAutoOpsRuleRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.GetAutoOpsRuleRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.GetAutoOpsRuleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.GetAutoOpsRuleRequest.prototype.setEnvironmentId =
  function (value) {
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
  proto.bucketeer.autoops.GetAutoOpsRuleResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.GetAutoOpsRuleResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.GetAutoOpsRuleResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.GetAutoOpsRuleResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        autoOpsRule:
          (f = msg.getAutoOpsRule()) &&
          proto_autoops_auto_ops_rule_pb.AutoOpsRule.toObject(
            includeInstance,
            f
          )
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
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.GetAutoOpsRuleResponse();
  return proto.bucketeer.autoops.GetAutoOpsRuleResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.GetAutoOpsRuleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_autoops_auto_ops_rule_pb.AutoOpsRule();
          reader.readMessage(
            value,
            proto_autoops_auto_ops_rule_pb.AutoOpsRule
              .deserializeBinaryFromReader
          );
          msg.setAutoOpsRule(value);
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
proto.bucketeer.autoops.GetAutoOpsRuleResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.GetAutoOpsRuleResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.GetAutoOpsRuleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAutoOpsRule();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_autoops_auto_ops_rule_pb.AutoOpsRule.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AutoOpsRule auto_ops_rule = 1;
 * @return {?proto.bucketeer.autoops.AutoOpsRule}
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.prototype.getAutoOpsRule =
  function () {
    return /** @type{?proto.bucketeer.autoops.AutoOpsRule} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_auto_ops_rule_pb.AutoOpsRule,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.AutoOpsRule|undefined} value
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleResponse} returns this
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.prototype.setAutoOpsRule =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.GetAutoOpsRuleResponse} returns this
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.prototype.clearAutoOpsRule =
  function () {
    return this.setAutoOpsRule(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.GetAutoOpsRuleResponse.prototype.hasAutoOpsRule =
  function () {
    return jspb.Message.getField(this, 1) != null;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.repeatedFields_ = [6, 7];

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
  proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.CreateAutoOpsRuleRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.CreateAutoOpsRuleRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        command:
          (f = msg.getCommand()) &&
          proto_autoops_command_pb.CreateAutoOpsRuleCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 4, ''),
        opsType: jspb.Message.getFieldWithDefault(msg, 5, 0),
        opsEventRateClausesList: jspb.Message.toObjectList(
          msg.getOpsEventRateClausesList(),
          proto_autoops_clause_pb.OpsEventRateClause.toObject,
          includeInstance
        ),
        datetimeClausesList: jspb.Message.toObjectList(
          msg.getDatetimeClausesList(),
          proto_autoops_clause_pb.DatetimeClause.toObject,
          includeInstance
        )
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
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.CreateAutoOpsRuleRequest();
  return proto.bucketeer.autoops.CreateAutoOpsRuleRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = new proto_autoops_command_pb.CreateAutoOpsRuleCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.CreateAutoOpsRuleCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setFeatureId(value);
          break;
        case 5:
          var value = /** @type {!proto.bucketeer.autoops.OpsType} */ (
            reader.readEnum()
          );
          msg.setOpsType(value);
          break;
        case 6:
          var value = new proto_autoops_clause_pb.OpsEventRateClause();
          reader.readMessage(
            value,
            proto_autoops_clause_pb.OpsEventRateClause
              .deserializeBinaryFromReader
          );
          msg.addOpsEventRateClauses(value);
          break;
        case 7:
          var value = new proto_autoops_clause_pb.DatetimeClause();
          reader.readMessage(
            value,
            proto_autoops_clause_pb.DatetimeClause.deserializeBinaryFromReader
          );
          msg.addDatetimeClauses(value);
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
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.CreateAutoOpsRuleRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_autoops_command_pb.CreateAutoOpsRuleCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getFeatureId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
    f = message.getOpsType();
    if (f !== 0.0) {
      writer.writeEnum(5, f);
    }
    f = message.getOpsEventRateClausesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        6,
        f,
        proto_autoops_clause_pb.OpsEventRateClause.serializeBinaryToWriter
      );
    }
    f = message.getDatetimeClausesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        7,
        f,
        proto_autoops_clause_pb.DatetimeClause.serializeBinaryToWriter
      );
    }
  };

/**
 * optional CreateAutoOpsRuleCommand command = 2;
 * @return {?proto.bucketeer.autoops.CreateAutoOpsRuleCommand}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.CreateAutoOpsRuleCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.CreateAutoOpsRuleCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.CreateAutoOpsRuleCommand|undefined} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional string feature_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.getFeatureId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.setFeatureId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
  };

/**
 * optional OpsType ops_type = 5;
 * @return {!proto.bucketeer.autoops.OpsType}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.getOpsType =
  function () {
    return /** @type {!proto.bucketeer.autoops.OpsType} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.OpsType} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.setOpsType =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
  };

/**
 * repeated OpsEventRateClause ops_event_rate_clauses = 6;
 * @return {!Array<!proto.bucketeer.autoops.OpsEventRateClause>}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.getOpsEventRateClausesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.OpsEventRateClause>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_clause_pb.OpsEventRateClause,
        6
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.OpsEventRateClause>} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.setOpsEventRateClausesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 6, value);
  };

/**
 * @param {!proto.bucketeer.autoops.OpsEventRateClause=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.OpsEventRateClause}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.addOpsEventRateClauses =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      6,
      opt_value,
      proto.bucketeer.autoops.OpsEventRateClause,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.clearOpsEventRateClausesList =
  function () {
    return this.setOpsEventRateClausesList([]);
  };

/**
 * repeated DatetimeClause datetime_clauses = 7;
 * @return {!Array<!proto.bucketeer.autoops.DatetimeClause>}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.getDatetimeClausesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.DatetimeClause>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_clause_pb.DatetimeClause,
        7
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.DatetimeClause>} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.setDatetimeClausesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 7, value);
  };

/**
 * @param {!proto.bucketeer.autoops.DatetimeClause=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.DatetimeClause}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.addDatetimeClauses =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      7,
      opt_value,
      proto.bucketeer.autoops.DatetimeClause,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleRequest.prototype.clearDatetimeClausesList =
  function () {
    return this.setDatetimeClausesList([]);
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
  proto.bucketeer.autoops.CreateAutoOpsRuleResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.CreateAutoOpsRuleResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.CreateAutoOpsRuleResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        autoOpsRule:
          (f = msg.getAutoOpsRule()) &&
          proto_autoops_auto_ops_rule_pb.AutoOpsRule.toObject(
            includeInstance,
            f
          )
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
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.CreateAutoOpsRuleResponse();
  return proto.bucketeer.autoops.CreateAutoOpsRuleResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_autoops_auto_ops_rule_pb.AutoOpsRule();
          reader.readMessage(
            value,
            proto_autoops_auto_ops_rule_pb.AutoOpsRule
              .deserializeBinaryFromReader
          );
          msg.setAutoOpsRule(value);
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
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.CreateAutoOpsRuleResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAutoOpsRule();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_autoops_auto_ops_rule_pb.AutoOpsRule.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AutoOpsRule auto_ops_rule = 1;
 * @return {?proto.bucketeer.autoops.AutoOpsRule}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.prototype.getAutoOpsRule =
  function () {
    return /** @type{?proto.bucketeer.autoops.AutoOpsRule} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_auto_ops_rule_pb.AutoOpsRule,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.AutoOpsRule|undefined} value
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.prototype.setAutoOpsRule =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.CreateAutoOpsRuleResponse} returns this
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.prototype.clearAutoOpsRule =
  function () {
    return this.setAutoOpsRule(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.CreateAutoOpsRuleResponse.prototype.hasAutoOpsRule =
  function () {
    return jspb.Message.getField(this, 1) != null;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.repeatedFields_ = [4];

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
  proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.ListAutoOpsRulesRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ListAutoOpsRulesRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
        cursor: jspb.Message.getFieldWithDefault(msg, 3, ''),
        featureIdsList:
          (f = jspb.Message.getRepeatedField(msg, 4)) == null ? undefined : f,
        environmentId: jspb.Message.getFieldWithDefault(msg, 5, '')
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
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest}
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ListAutoOpsRulesRequest();
  return proto.bucketeer.autoops.ListAutoOpsRulesRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest}
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.addFeatureIds(value);
          break;
        case 5:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ListAutoOpsRulesRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getPageSize();
    if (f !== 0) {
      writer.writeInt64(2, f);
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getFeatureIdsList();
    if (f.length > 0) {
      writer.writeRepeatedString(4, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(5, f);
    }
  };

/**
 * optional int64 page_size = 2;
 * @return {number}
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.setPageSize =
  function (value) {
    return jspb.Message.setProto3IntField(this, 2, value);
  };

/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * repeated string feature_ids = 4;
 * @return {!Array<string>}
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.getFeatureIdsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 4)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.setFeatureIdsList =
  function (value) {
    return jspb.Message.setField(this, 4, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.addFeatureIds =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 4, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.clearFeatureIdsList =
  function () {
    return this.setFeatureIdsList([]);
  };

/**
 * optional string environment_id = 5;
 * @return {string}
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 5, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesRequest} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 5, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.repeatedFields_ = [1];

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
  proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.ListAutoOpsRulesResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ListAutoOpsRulesResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ListAutoOpsRulesResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        autoOpsRulesList: jspb.Message.toObjectList(
          msg.getAutoOpsRulesList(),
          proto_autoops_auto_ops_rule_pb.AutoOpsRule.toObject,
          includeInstance
        ),
        cursor: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesResponse}
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ListAutoOpsRulesResponse();
  return proto.bucketeer.autoops.ListAutoOpsRulesResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ListAutoOpsRulesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesResponse}
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_autoops_auto_ops_rule_pb.AutoOpsRule();
          reader.readMessage(
            value,
            proto_autoops_auto_ops_rule_pb.AutoOpsRule
              .deserializeBinaryFromReader
          );
          msg.addAutoOpsRules(value);
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
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ListAutoOpsRulesResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ListAutoOpsRulesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAutoOpsRulesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_autoops_auto_ops_rule_pb.AutoOpsRule.serializeBinaryToWriter
      );
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
  };

/**
 * repeated AutoOpsRule auto_ops_rules = 1;
 * @return {!Array<!proto.bucketeer.autoops.AutoOpsRule>}
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.getAutoOpsRulesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.AutoOpsRule>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_auto_ops_rule_pb.AutoOpsRule,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.AutoOpsRule>} value
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesResponse} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.setAutoOpsRulesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.autoops.AutoOpsRule=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.AutoOpsRule}
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.addAutoOpsRules =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.autoops.AutoOpsRule,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesResponse} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.clearAutoOpsRulesList =
  function () {
    return this.setAutoOpsRulesList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListAutoOpsRulesResponse} returns this
 */
proto.bucketeer.autoops.ListAutoOpsRulesResponse.prototype.setCursor =
  function (value) {
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
  proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.StopAutoOpsRuleRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.StopAutoOpsRuleRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_autoops_command_pb.StopAutoOpsRuleCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, '')
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
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.StopAutoOpsRuleRequest();
  return proto.bucketeer.autoops.StopAutoOpsRuleRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value = new proto_autoops_command_pb.StopAutoOpsRuleCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.StopAutoOpsRuleCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.StopAutoOpsRuleRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_autoops_command_pb.StopAutoOpsRuleCommand.serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional StopAutoOpsRuleCommand command = 3;
 * @return {?proto.bucketeer.autoops.StopAutoOpsRuleCommand}
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.StopAutoOpsRuleCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.StopAutoOpsRuleCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.StopAutoOpsRuleCommand|undefined} value
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 3, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.StopAutoOpsRuleRequest.prototype.setEnvironmentId =
  function (value) {
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
  proto.bucketeer.autoops.StopAutoOpsRuleResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.StopAutoOpsRuleResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.StopAutoOpsRuleResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.StopAutoOpsRuleResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.StopAutoOpsRuleResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.StopAutoOpsRuleResponse();
  return proto.bucketeer.autoops.StopAutoOpsRuleResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.StopAutoOpsRuleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.StopAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.StopAutoOpsRuleResponse.deserializeBinaryFromReader =
  function (msg, reader) {
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
proto.bucketeer.autoops.StopAutoOpsRuleResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.StopAutoOpsRuleResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.StopAutoOpsRuleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.StopAutoOpsRuleResponse.serializeBinaryToWriter =
  function (message, writer) {
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
  proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_autoops_command_pb.DeleteAutoOpsRuleCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, '')
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
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.DeleteAutoOpsRuleRequest();
  return proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value = new proto_autoops_command_pb.DeleteAutoOpsRuleCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.DeleteAutoOpsRuleCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_autoops_command_pb.DeleteAutoOpsRuleCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional DeleteAutoOpsRuleCommand command = 3;
 * @return {?proto.bucketeer.autoops.DeleteAutoOpsRuleCommand}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.DeleteAutoOpsRuleCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.DeleteAutoOpsRuleCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.DeleteAutoOpsRuleCommand|undefined} value
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleRequest.prototype.setEnvironmentId =
  function (value) {
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
  proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.DeleteAutoOpsRuleResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.DeleteAutoOpsRuleResponse();
  return proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.DeleteAutoOpsRuleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.DeleteAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.deserializeBinaryFromReader =
  function (msg, reader) {
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
proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.DeleteAutoOpsRuleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.DeleteAutoOpsRuleResponse.serializeBinaryToWriter =
  function (message, writer) {
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
  proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.OpsEventRateClauseChange.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.OpsEventRateClauseChange} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.OpsEventRateClauseChange.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        clause:
          (f = msg.getClause()) &&
          proto_autoops_clause_pb.OpsEventRateClause.toObject(
            includeInstance,
            f
          ),
        changeType: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange}
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.OpsEventRateClauseChange();
  return proto.bucketeer.autoops.OpsEventRateClauseChange.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.OpsEventRateClauseChange} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange}
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.deserializeBinaryFromReader =
  function (msg, reader) {
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
          var value = new proto_autoops_clause_pb.OpsEventRateClause();
          reader.readMessage(
            value,
            proto_autoops_clause_pb.OpsEventRateClause
              .deserializeBinaryFromReader
          );
          msg.setClause(value);
          break;
        case 3:
          var value = /** @type {!proto.bucketeer.autoops.ChangeType} */ (
            reader.readEnum()
          );
          msg.setChangeType(value);
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
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.OpsEventRateClauseChange.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.OpsEventRateClauseChange} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getClause();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_autoops_clause_pb.OpsEventRateClause.serializeBinaryToWriter
      );
    }
    f = message.getChangeType();
    if (f !== 0.0) {
      writer.writeEnum(3, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange} returns this
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional OpsEventRateClause clause = 2;
 * @return {?proto.bucketeer.autoops.OpsEventRateClause}
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.getClause =
  function () {
    return /** @type{?proto.bucketeer.autoops.OpsEventRateClause} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_clause_pb.OpsEventRateClause,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.OpsEventRateClause|undefined} value
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange} returns this
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.setClause =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange} returns this
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.clearClause =
  function () {
    return this.setClause(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.hasClause =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional ChangeType change_type = 3;
 * @return {!proto.bucketeer.autoops.ChangeType}
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.getChangeType =
  function () {
    return /** @type {!proto.bucketeer.autoops.ChangeType} */ (
      jspb.Message.getFieldWithDefault(this, 3, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ChangeType} value
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange} returns this
 */
proto.bucketeer.autoops.OpsEventRateClauseChange.prototype.setChangeType =
  function (value) {
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
  proto.bucketeer.autoops.DatetimeClauseChange.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.DatetimeClauseChange.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.DatetimeClauseChange} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.DatetimeClauseChange.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        clause:
          (f = msg.getClause()) &&
          proto_autoops_clause_pb.DatetimeClause.toObject(includeInstance, f),
        changeType: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange}
 */
proto.bucketeer.autoops.DatetimeClauseChange.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.DatetimeClauseChange();
  return proto.bucketeer.autoops.DatetimeClauseChange.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.DatetimeClauseChange} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange}
 */
proto.bucketeer.autoops.DatetimeClauseChange.deserializeBinaryFromReader =
  function (msg, reader) {
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
          var value = new proto_autoops_clause_pb.DatetimeClause();
          reader.readMessage(
            value,
            proto_autoops_clause_pb.DatetimeClause.deserializeBinaryFromReader
          );
          msg.setClause(value);
          break;
        case 3:
          var value = /** @type {!proto.bucketeer.autoops.ChangeType} */ (
            reader.readEnum()
          );
          msg.setChangeType(value);
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
proto.bucketeer.autoops.DatetimeClauseChange.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.DatetimeClauseChange.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.DatetimeClauseChange} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.DatetimeClauseChange.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getClause();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_autoops_clause_pb.DatetimeClause.serializeBinaryToWriter
      );
    }
    f = message.getChangeType();
    if (f !== 0.0) {
      writer.writeEnum(3, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange} returns this
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional DatetimeClause clause = 2;
 * @return {?proto.bucketeer.autoops.DatetimeClause}
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.getClause = function () {
  return /** @type{?proto.bucketeer.autoops.DatetimeClause} */ (
    jspb.Message.getWrapperField(
      this,
      proto_autoops_clause_pb.DatetimeClause,
      2
    )
  );
};

/**
 * @param {?proto.bucketeer.autoops.DatetimeClause|undefined} value
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange} returns this
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.setClause = function (
  value
) {
  return jspb.Message.setWrapperField(this, 2, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange} returns this
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.clearClause =
  function () {
    return this.setClause(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.hasClause = function () {
  return jspb.Message.getField(this, 2) != null;
};

/**
 * optional ChangeType change_type = 3;
 * @return {!proto.bucketeer.autoops.ChangeType}
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.getChangeType =
  function () {
    return /** @type {!proto.bucketeer.autoops.ChangeType} */ (
      jspb.Message.getFieldWithDefault(this, 3, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ChangeType} value
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange} returns this
 */
proto.bucketeer.autoops.DatetimeClauseChange.prototype.setChangeType =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 3, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.repeatedFields_ = [
  4, 5, 6, 7, 8, 10, 11
];

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
  proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        addOpsEventRateClauseCommandsList: jspb.Message.toObjectList(
          msg.getAddOpsEventRateClauseCommandsList(),
          proto_autoops_command_pb.AddOpsEventRateClauseCommand.toObject,
          includeInstance
        ),
        changeOpsEventRateClauseCommandsList: jspb.Message.toObjectList(
          msg.getChangeOpsEventRateClauseCommandsList(),
          proto_autoops_command_pb.ChangeOpsEventRateClauseCommand.toObject,
          includeInstance
        ),
        deleteClauseCommandsList: jspb.Message.toObjectList(
          msg.getDeleteClauseCommandsList(),
          proto_autoops_command_pb.DeleteClauseCommand.toObject,
          includeInstance
        ),
        addDatetimeClauseCommandsList: jspb.Message.toObjectList(
          msg.getAddDatetimeClauseCommandsList(),
          proto_autoops_command_pb.AddDatetimeClauseCommand.toObject,
          includeInstance
        ),
        changeDatetimeClauseCommandsList: jspb.Message.toObjectList(
          msg.getChangeDatetimeClauseCommandsList(),
          proto_autoops_command_pb.ChangeDatetimeClauseCommand.toObject,
          includeInstance
        ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 9, ''),
        opsEventRateClauseChangesList: jspb.Message.toObjectList(
          msg.getOpsEventRateClauseChangesList(),
          proto.bucketeer.autoops.OpsEventRateClauseChange.toObject,
          includeInstance
        ),
        datetimeClauseChangesList: jspb.Message.toObjectList(
          msg.getDatetimeClauseChangesList(),
          proto.bucketeer.autoops.DatetimeClauseChange.toObject,
          includeInstance
        )
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
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.UpdateAutoOpsRuleRequest();
  return proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 4:
          var value =
            new proto_autoops_command_pb.AddOpsEventRateClauseCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.AddOpsEventRateClauseCommand
              .deserializeBinaryFromReader
          );
          msg.addAddOpsEventRateClauseCommands(value);
          break;
        case 5:
          var value =
            new proto_autoops_command_pb.ChangeOpsEventRateClauseCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.ChangeOpsEventRateClauseCommand
              .deserializeBinaryFromReader
          );
          msg.addChangeOpsEventRateClauseCommands(value);
          break;
        case 6:
          var value = new proto_autoops_command_pb.DeleteClauseCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.DeleteClauseCommand
              .deserializeBinaryFromReader
          );
          msg.addDeleteClauseCommands(value);
          break;
        case 7:
          var value = new proto_autoops_command_pb.AddDatetimeClauseCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.AddDatetimeClauseCommand
              .deserializeBinaryFromReader
          );
          msg.addAddDatetimeClauseCommands(value);
          break;
        case 8:
          var value =
            new proto_autoops_command_pb.ChangeDatetimeClauseCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.ChangeDatetimeClauseCommand
              .deserializeBinaryFromReader
          );
          msg.addChangeDatetimeClauseCommands(value);
          break;
        case 9:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 10:
          var value = new proto.bucketeer.autoops.OpsEventRateClauseChange();
          reader.readMessage(
            value,
            proto.bucketeer.autoops.OpsEventRateClauseChange
              .deserializeBinaryFromReader
          );
          msg.addOpsEventRateClauseChanges(value);
          break;
        case 11:
          var value = new proto.bucketeer.autoops.DatetimeClauseChange();
          reader.readMessage(
            value,
            proto.bucketeer.autoops.DatetimeClauseChange
              .deserializeBinaryFromReader
          );
          msg.addDatetimeClauseChanges(value);
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
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getAddOpsEventRateClauseCommandsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        4,
        f,
        proto_autoops_command_pb.AddOpsEventRateClauseCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeOpsEventRateClauseCommandsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        5,
        f,
        proto_autoops_command_pb.ChangeOpsEventRateClauseCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getDeleteClauseCommandsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        6,
        f,
        proto_autoops_command_pb.DeleteClauseCommand.serializeBinaryToWriter
      );
    }
    f = message.getAddDatetimeClauseCommandsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        7,
        f,
        proto_autoops_command_pb.AddDatetimeClauseCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeDatetimeClauseCommandsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        8,
        f,
        proto_autoops_command_pb.ChangeDatetimeClauseCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(9, f);
    }
    f = message.getOpsEventRateClauseChangesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        10,
        f,
        proto.bucketeer.autoops.OpsEventRateClauseChange.serializeBinaryToWriter
      );
    }
    f = message.getDatetimeClauseChangesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        11,
        f,
        proto.bucketeer.autoops.DatetimeClauseChange.serializeBinaryToWriter
      );
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * repeated AddOpsEventRateClauseCommand add_ops_event_rate_clause_commands = 4;
 * @return {!Array<!proto.bucketeer.autoops.AddOpsEventRateClauseCommand>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getAddOpsEventRateClauseCommandsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.AddOpsEventRateClauseCommand>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_command_pb.AddOpsEventRateClauseCommand,
        4
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.AddOpsEventRateClauseCommand>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setAddOpsEventRateClauseCommandsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 4, value);
  };

/**
 * @param {!proto.bucketeer.autoops.AddOpsEventRateClauseCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.AddOpsEventRateClauseCommand}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addAddOpsEventRateClauseCommands =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      4,
      opt_value,
      proto.bucketeer.autoops.AddOpsEventRateClauseCommand,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearAddOpsEventRateClauseCommandsList =
  function () {
    return this.setAddOpsEventRateClauseCommandsList([]);
  };

/**
 * repeated ChangeOpsEventRateClauseCommand change_ops_event_rate_clause_commands = 5;
 * @return {!Array<!proto.bucketeer.autoops.ChangeOpsEventRateClauseCommand>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getChangeOpsEventRateClauseCommandsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.ChangeOpsEventRateClauseCommand>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_command_pb.ChangeOpsEventRateClauseCommand,
        5
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.ChangeOpsEventRateClauseCommand>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setChangeOpsEventRateClauseCommandsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 5, value);
  };

/**
 * @param {!proto.bucketeer.autoops.ChangeOpsEventRateClauseCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ChangeOpsEventRateClauseCommand}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addChangeOpsEventRateClauseCommands =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      5,
      opt_value,
      proto.bucketeer.autoops.ChangeOpsEventRateClauseCommand,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearChangeOpsEventRateClauseCommandsList =
  function () {
    return this.setChangeOpsEventRateClauseCommandsList([]);
  };

/**
 * repeated DeleteClauseCommand delete_clause_commands = 6;
 * @return {!Array<!proto.bucketeer.autoops.DeleteClauseCommand>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getDeleteClauseCommandsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.DeleteClauseCommand>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_command_pb.DeleteClauseCommand,
        6
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.DeleteClauseCommand>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setDeleteClauseCommandsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 6, value);
  };

/**
 * @param {!proto.bucketeer.autoops.DeleteClauseCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.DeleteClauseCommand}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addDeleteClauseCommands =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      6,
      opt_value,
      proto.bucketeer.autoops.DeleteClauseCommand,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearDeleteClauseCommandsList =
  function () {
    return this.setDeleteClauseCommandsList([]);
  };

/**
 * repeated AddDatetimeClauseCommand add_datetime_clause_commands = 7;
 * @return {!Array<!proto.bucketeer.autoops.AddDatetimeClauseCommand>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getAddDatetimeClauseCommandsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.AddDatetimeClauseCommand>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_command_pb.AddDatetimeClauseCommand,
        7
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.AddDatetimeClauseCommand>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setAddDatetimeClauseCommandsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 7, value);
  };

/**
 * @param {!proto.bucketeer.autoops.AddDatetimeClauseCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.AddDatetimeClauseCommand}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addAddDatetimeClauseCommands =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      7,
      opt_value,
      proto.bucketeer.autoops.AddDatetimeClauseCommand,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearAddDatetimeClauseCommandsList =
  function () {
    return this.setAddDatetimeClauseCommandsList([]);
  };

/**
 * repeated ChangeDatetimeClauseCommand change_datetime_clause_commands = 8;
 * @return {!Array<!proto.bucketeer.autoops.ChangeDatetimeClauseCommand>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getChangeDatetimeClauseCommandsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.ChangeDatetimeClauseCommand>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_command_pb.ChangeDatetimeClauseCommand,
        8
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.ChangeDatetimeClauseCommand>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setChangeDatetimeClauseCommandsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 8, value);
  };

/**
 * @param {!proto.bucketeer.autoops.ChangeDatetimeClauseCommand=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ChangeDatetimeClauseCommand}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addChangeDatetimeClauseCommands =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      8,
      opt_value,
      proto.bucketeer.autoops.ChangeDatetimeClauseCommand,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearChangeDatetimeClauseCommandsList =
  function () {
    return this.setChangeDatetimeClauseCommandsList([]);
  };

/**
 * optional string environment_id = 9;
 * @return {string}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 9, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 9, value);
  };

/**
 * repeated OpsEventRateClauseChange ops_event_rate_clause_changes = 10;
 * @return {!Array<!proto.bucketeer.autoops.OpsEventRateClauseChange>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getOpsEventRateClauseChangesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.OpsEventRateClauseChange>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto.bucketeer.autoops.OpsEventRateClauseChange,
        10
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.OpsEventRateClauseChange>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setOpsEventRateClauseChangesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 10, value);
  };

/**
 * @param {!proto.bucketeer.autoops.OpsEventRateClauseChange=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.OpsEventRateClauseChange}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addOpsEventRateClauseChanges =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      10,
      opt_value,
      proto.bucketeer.autoops.OpsEventRateClauseChange,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearOpsEventRateClauseChangesList =
  function () {
    return this.setOpsEventRateClauseChangesList([]);
  };

/**
 * repeated DatetimeClauseChange datetime_clause_changes = 11;
 * @return {!Array<!proto.bucketeer.autoops.DatetimeClauseChange>}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.getDatetimeClauseChangesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.DatetimeClauseChange>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto.bucketeer.autoops.DatetimeClauseChange,
        11
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.DatetimeClauseChange>} value
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.setDatetimeClauseChangesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 11, value);
  };

/**
 * @param {!proto.bucketeer.autoops.DatetimeClauseChange=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.DatetimeClauseChange}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.addDatetimeClauseChanges =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      11,
      opt_value,
      proto.bucketeer.autoops.DatetimeClauseChange,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleRequest} returns this
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleRequest.prototype.clearDatetimeClauseChangesList =
  function () {
    return this.setDatetimeClauseChangesList([]);
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
  proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.UpdateAutoOpsRuleResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.UpdateAutoOpsRuleResponse();
  return proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.UpdateAutoOpsRuleResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.UpdateAutoOpsRuleResponse}
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.deserializeBinaryFromReader =
  function (msg, reader) {
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
proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.UpdateAutoOpsRuleResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.UpdateAutoOpsRuleResponse.serializeBinaryToWriter =
  function (message, writer) {
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
  proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.ExecuteAutoOpsRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ExecuteAutoOpsRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        executeAutoOpsRuleCommand:
          (f = msg.getExecuteAutoOpsRuleCommand()) &&
          proto_autoops_command_pb.ExecuteAutoOpsRuleCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 5, ''),
        clauseId: jspb.Message.getFieldWithDefault(msg, 6, '')
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
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ExecuteAutoOpsRequest();
  return proto.bucketeer.autoops.ExecuteAutoOpsRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 4:
          var value = new proto_autoops_command_pb.ExecuteAutoOpsRuleCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.ExecuteAutoOpsRuleCommand
              .deserializeBinaryFromReader
          );
          msg.setExecuteAutoOpsRuleCommand(value);
          break;
        case 5:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setClauseId(value);
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
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ExecuteAutoOpsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getExecuteAutoOpsRuleCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_autoops_command_pb.ExecuteAutoOpsRuleCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(5, f);
    }
    f = message.getClauseId();
    if (f.length > 0) {
      writer.writeString(6, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} returns this
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional ExecuteAutoOpsRuleCommand execute_auto_ops_rule_command = 4;
 * @return {?proto.bucketeer.autoops.ExecuteAutoOpsRuleCommand}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.getExecuteAutoOpsRuleCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.ExecuteAutoOpsRuleCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.ExecuteAutoOpsRuleCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.ExecuteAutoOpsRuleCommand|undefined} value
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} returns this
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.setExecuteAutoOpsRuleCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} returns this
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.clearExecuteAutoOpsRuleCommand =
  function () {
    return this.setExecuteAutoOpsRuleCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.hasExecuteAutoOpsRuleCommand =
  function () {
    return jspb.Message.getField(this, 4) != null;
  };

/**
 * optional string environment_id = 5;
 * @return {string}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 5, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} returns this
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 5, value);
  };

/**
 * optional string clause_id = 6;
 * @return {string}
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.getClauseId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsRequest} returns this
 */
proto.bucketeer.autoops.ExecuteAutoOpsRequest.prototype.setClauseId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 6, value);
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
  proto.bucketeer.autoops.ExecuteAutoOpsResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.ExecuteAutoOpsResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ExecuteAutoOpsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ExecuteAutoOpsResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        alreadyTriggered: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
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
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsResponse}
 */
proto.bucketeer.autoops.ExecuteAutoOpsResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ExecuteAutoOpsResponse();
  return proto.bucketeer.autoops.ExecuteAutoOpsResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ExecuteAutoOpsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsResponse}
 */
proto.bucketeer.autoops.ExecuteAutoOpsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {boolean} */ (reader.readBool());
          msg.setAlreadyTriggered(value);
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
proto.bucketeer.autoops.ExecuteAutoOpsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ExecuteAutoOpsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ExecuteAutoOpsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ExecuteAutoOpsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAlreadyTriggered();
    if (f) {
      writer.writeBool(1, f);
    }
  };

/**
 * optional bool already_triggered = 1;
 * @return {boolean}
 */
proto.bucketeer.autoops.ExecuteAutoOpsResponse.prototype.getAlreadyTriggered =
  function () {
    return /** @type {boolean} */ (
      jspb.Message.getBooleanFieldWithDefault(this, 1, false)
    );
  };

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.autoops.ExecuteAutoOpsResponse} returns this
 */
proto.bucketeer.autoops.ExecuteAutoOpsResponse.prototype.setAlreadyTriggered =
  function (value) {
    return jspb.Message.setProto3BooleanField(this, 1, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.ListOpsCountsRequest.repeatedFields_ = [4, 5];

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
  proto.bucketeer.autoops.ListOpsCountsRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.ListOpsCountsRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ListOpsCountsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ListOpsCountsRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
        cursor: jspb.Message.getFieldWithDefault(msg, 3, ''),
        autoOpsRuleIdsList:
          (f = jspb.Message.getRepeatedField(msg, 4)) == null ? undefined : f,
        featureIdsList:
          (f = jspb.Message.getRepeatedField(msg, 5)) == null ? undefined : f,
        environmentId: jspb.Message.getFieldWithDefault(msg, 6, '')
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
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ListOpsCountsRequest();
  return proto.bucketeer.autoops.ListOpsCountsRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ListOpsCountsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.addAutoOpsRuleIds(value);
          break;
        case 5:
          var value = /** @type {string} */ (reader.readString());
          msg.addFeatureIds(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ListOpsCountsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ListOpsCountsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ListOpsCountsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getPageSize();
    if (f !== 0) {
      writer.writeInt64(2, f);
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getAutoOpsRuleIdsList();
    if (f.length > 0) {
      writer.writeRepeatedString(4, f);
    }
    f = message.getFeatureIdsList();
    if (f.length > 0) {
      writer.writeRepeatedString(5, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(6, f);
    }
  };

/**
 * optional int64 page_size = 2;
 * @return {number}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.setPageSize = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 2, value);
};

/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.getCursor = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * repeated string auto_ops_rule_ids = 4;
 * @return {!Array<string>}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.getAutoOpsRuleIdsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 4)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.setAutoOpsRuleIdsList =
  function (value) {
    return jspb.Message.setField(this, 4, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.addAutoOpsRuleIds =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 4, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.clearAutoOpsRuleIdsList =
  function () {
    return this.setAutoOpsRuleIdsList([]);
  };

/**
 * repeated string feature_ids = 5;
 * @return {!Array<string>}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.getFeatureIdsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 5)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.setFeatureIdsList =
  function (value) {
    return jspb.Message.setField(this, 5, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.addFeatureIds =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 5, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.clearFeatureIdsList =
  function () {
    return this.setFeatureIdsList([]);
  };

/**
 * optional string environment_id = 6;
 * @return {string}
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsRequest} returns this
 */
proto.bucketeer.autoops.ListOpsCountsRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 6, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.ListOpsCountsResponse.repeatedFields_ = [2];

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
  proto.bucketeer.autoops.ListOpsCountsResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.ListOpsCountsResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ListOpsCountsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ListOpsCountsResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        cursor: jspb.Message.getFieldWithDefault(msg, 1, ''),
        opsCountsList: jspb.Message.toObjectList(
          msg.getOpsCountsList(),
          proto_autoops_ops_count_pb.OpsCount.toObject,
          includeInstance
        )
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
 * @return {!proto.bucketeer.autoops.ListOpsCountsResponse}
 */
proto.bucketeer.autoops.ListOpsCountsResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ListOpsCountsResponse();
  return proto.bucketeer.autoops.ListOpsCountsResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ListOpsCountsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ListOpsCountsResponse}
 */
proto.bucketeer.autoops.ListOpsCountsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 2:
          var value = new proto_autoops_ops_count_pb.OpsCount();
          reader.readMessage(
            value,
            proto_autoops_ops_count_pb.OpsCount.deserializeBinaryFromReader
          );
          msg.addOpsCounts(value);
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
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ListOpsCountsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ListOpsCountsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ListOpsCountsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOpsCountsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        2,
        f,
        proto_autoops_ops_count_pb.OpsCount.serializeBinaryToWriter
      );
    }
  };

/**
 * optional string cursor = 1;
 * @return {string}
 */
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsResponse} returns this
 */
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * repeated OpsCount ops_counts = 2;
 * @return {!Array<!proto.bucketeer.autoops.OpsCount>}
 */
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.getOpsCountsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.OpsCount>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_ops_count_pb.OpsCount,
        2
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.OpsCount>} value
 * @return {!proto.bucketeer.autoops.ListOpsCountsResponse} returns this
 */
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.setOpsCountsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 2, value);
  };

/**
 * @param {!proto.bucketeer.autoops.OpsCount=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.OpsCount}
 */
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.addOpsCounts =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      2,
      opt_value,
      proto.bucketeer.autoops.OpsCount,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListOpsCountsResponse} returns this
 */
proto.bucketeer.autoops.ListOpsCountsResponse.prototype.clearOpsCountsList =
  function () {
    return this.setOpsCountsList([]);
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
  proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.CreateProgressiveRolloutRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.CreateProgressiveRolloutRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        command:
          (f = msg.getCommand()) &&
          proto_autoops_command_pb.CreateProgressiveRolloutCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 4, ''),
        progressiveRolloutManualScheduleClause:
          (f = msg.getProgressiveRolloutManualScheduleClause()) &&
          proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause.toObject(
            includeInstance,
            f
          ),
        progressiveRolloutTemplateScheduleClause:
          (f = msg.getProgressiveRolloutTemplateScheduleClause()) &&
          proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause.toObject(
            includeInstance,
            f
          )
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
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.CreateProgressiveRolloutRequest();
    return proto.bucketeer.autoops.CreateProgressiveRolloutRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value =
            new proto_autoops_command_pb.CreateProgressiveRolloutCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.CreateProgressiveRolloutCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setFeatureId(value);
          break;
        case 5:
          var value =
            new proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause();
          reader.readMessage(
            value,
            proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause
              .deserializeBinaryFromReader
          );
          msg.setProgressiveRolloutManualScheduleClause(value);
          break;
        case 6:
          var value =
            new proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause();
          reader.readMessage(
            value,
            proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause
              .deserializeBinaryFromReader
          );
          msg.setProgressiveRolloutTemplateScheduleClause(value);
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
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.CreateProgressiveRolloutRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_autoops_command_pb.CreateProgressiveRolloutCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getFeatureId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
    f = message.getProgressiveRolloutManualScheduleClause();
    if (f != null) {
      writer.writeMessage(
        5,
        f,
        proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause
          .serializeBinaryToWriter
      );
    }
    f = message.getProgressiveRolloutTemplateScheduleClause();
    if (f != null) {
      writer.writeMessage(
        6,
        f,
        proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional CreateProgressiveRolloutCommand command = 2;
 * @return {?proto.bucketeer.autoops.CreateProgressiveRolloutCommand}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.CreateProgressiveRolloutCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.CreateProgressiveRolloutCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.CreateProgressiveRolloutCommand|undefined} value
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional string feature_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.getFeatureId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.setFeatureId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
  };

/**
 * optional ProgressiveRolloutManualScheduleClause progressive_rollout_manual_schedule_clause = 5;
 * @return {?proto.bucketeer.autoops.ProgressiveRolloutManualScheduleClause}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.getProgressiveRolloutManualScheduleClause =
  function () {
    return /** @type{?proto.bucketeer.autoops.ProgressiveRolloutManualScheduleClause} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_clause_pb.ProgressiveRolloutManualScheduleClause,
        5
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.ProgressiveRolloutManualScheduleClause|undefined} value
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.setProgressiveRolloutManualScheduleClause =
  function (value) {
    return jspb.Message.setWrapperField(this, 5, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.clearProgressiveRolloutManualScheduleClause =
  function () {
    return this.setProgressiveRolloutManualScheduleClause(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.hasProgressiveRolloutManualScheduleClause =
  function () {
    return jspb.Message.getField(this, 5) != null;
  };

/**
 * optional ProgressiveRolloutTemplateScheduleClause progressive_rollout_template_schedule_clause = 6;
 * @return {?proto.bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.getProgressiveRolloutTemplateScheduleClause =
  function () {
    return /** @type{?proto.bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_clause_pb.ProgressiveRolloutTemplateScheduleClause,
        6
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause|undefined} value
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.setProgressiveRolloutTemplateScheduleClause =
  function (value) {
    return jspb.Message.setWrapperField(this, 6, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.clearProgressiveRolloutTemplateScheduleClause =
  function () {
    return this.setProgressiveRolloutTemplateScheduleClause(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutRequest.prototype.hasProgressiveRolloutTemplateScheduleClause =
  function () {
    return jspb.Message.getField(this, 6) != null;
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
  proto.bucketeer.autoops.CreateProgressiveRolloutResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.CreateProgressiveRolloutResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.CreateProgressiveRolloutResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        progressiveRollout:
          (f = msg.getProgressiveRollout()) &&
          proto_autoops_progressive_rollout_pb.ProgressiveRollout.toObject(
            includeInstance,
            f
          )
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
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.CreateProgressiveRolloutResponse();
    return proto.bucketeer.autoops.CreateProgressiveRolloutResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value =
            new proto_autoops_progressive_rollout_pb.ProgressiveRollout();
          reader.readMessage(
            value,
            proto_autoops_progressive_rollout_pb.ProgressiveRollout
              .deserializeBinaryFromReader
          );
          msg.setProgressiveRollout(value);
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
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.CreateProgressiveRolloutResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getProgressiveRollout();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_autoops_progressive_rollout_pb.ProgressiveRollout
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional ProgressiveRollout progressive_rollout = 1;
 * @return {?proto.bucketeer.autoops.ProgressiveRollout}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.prototype.getProgressiveRollout =
  function () {
    return /** @type{?proto.bucketeer.autoops.ProgressiveRollout} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_progressive_rollout_pb.ProgressiveRollout,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.ProgressiveRollout|undefined} value
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.prototype.setProgressiveRollout =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.CreateProgressiveRolloutResponse} returns this
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.prototype.clearProgressiveRollout =
  function () {
    return this.setProgressiveRollout(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.CreateProgressiveRolloutResponse.prototype.hasProgressiveRollout =
  function () {
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
  proto.bucketeer.autoops.GetProgressiveRolloutRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.GetProgressiveRolloutRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.GetProgressiveRolloutRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.GetProgressiveRolloutRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, '')
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
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.GetProgressiveRolloutRequest();
    return proto.bucketeer.autoops.GetProgressiveRolloutRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.GetProgressiveRolloutRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.GetProgressiveRolloutRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.GetProgressiveRolloutRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.GetProgressiveRolloutRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.GetProgressiveRolloutRequest.prototype.setEnvironmentId =
  function (value) {
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
  proto.bucketeer.autoops.GetProgressiveRolloutResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.GetProgressiveRolloutResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.GetProgressiveRolloutResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.GetProgressiveRolloutResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        progressiveRollout:
          (f = msg.getProgressiveRollout()) &&
          proto_autoops_progressive_rollout_pb.ProgressiveRollout.toObject(
            includeInstance,
            f
          )
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
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.GetProgressiveRolloutResponse();
    return proto.bucketeer.autoops.GetProgressiveRolloutResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.GetProgressiveRolloutResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value =
            new proto_autoops_progressive_rollout_pb.ProgressiveRollout();
          reader.readMessage(
            value,
            proto_autoops_progressive_rollout_pb.ProgressiveRollout
              .deserializeBinaryFromReader
          );
          msg.setProgressiveRollout(value);
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
proto.bucketeer.autoops.GetProgressiveRolloutResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.GetProgressiveRolloutResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.GetProgressiveRolloutResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getProgressiveRollout();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_autoops_progressive_rollout_pb.ProgressiveRollout
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional ProgressiveRollout progressive_rollout = 1;
 * @return {?proto.bucketeer.autoops.ProgressiveRollout}
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.prototype.getProgressiveRollout =
  function () {
    return /** @type{?proto.bucketeer.autoops.ProgressiveRollout} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_progressive_rollout_pb.ProgressiveRollout,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.ProgressiveRollout|undefined} value
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutResponse} returns this
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.prototype.setProgressiveRollout =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.GetProgressiveRolloutResponse} returns this
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.prototype.clearProgressiveRollout =
  function () {
    return this.setProgressiveRollout(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.GetProgressiveRolloutResponse.prototype.hasProgressiveRollout =
  function () {
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
  proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.StopProgressiveRolloutRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.StopProgressiveRolloutRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_autoops_command_pb.StopProgressiveRolloutCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, ''),
        stoppedBy: jspb.Message.getFieldWithDefault(msg, 5, 0)
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
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.StopProgressiveRolloutRequest();
    return proto.bucketeer.autoops.StopProgressiveRolloutRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value =
            new proto_autoops_command_pb.StopProgressiveRolloutCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.StopProgressiveRolloutCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 5:
          var value =
            /** @type {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy} */ (
              reader.readEnum()
            );
          msg.setStoppedBy(value);
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
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.StopProgressiveRolloutRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_autoops_command_pb.StopProgressiveRolloutCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
    f = message.getStoppedBy();
    if (f !== 0.0) {
      writer.writeEnum(5, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional StopProgressiveRolloutCommand command = 3;
 * @return {?proto.bucketeer.autoops.StopProgressiveRolloutCommand}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.StopProgressiveRolloutCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.StopProgressiveRolloutCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.StopProgressiveRolloutCommand|undefined} value
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
  };

/**
 * optional ProgressiveRollout.StoppedBy stopped_by = 5;
 * @return {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy}
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.getStoppedBy =
  function () {
    return /** @type {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy} value
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.StopProgressiveRolloutRequest.prototype.setStoppedBy =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
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
  proto.bucketeer.autoops.StopProgressiveRolloutResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.StopProgressiveRolloutResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.StopProgressiveRolloutResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.StopProgressiveRolloutResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.StopProgressiveRolloutResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.StopProgressiveRolloutResponse();
    return proto.bucketeer.autoops.StopProgressiveRolloutResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.StopProgressiveRolloutResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.StopProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.StopProgressiveRolloutResponse.deserializeBinaryFromReader =
  function (msg, reader) {
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
proto.bucketeer.autoops.StopProgressiveRolloutResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.StopProgressiveRolloutResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.StopProgressiveRolloutResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.StopProgressiveRolloutResponse.serializeBinaryToWriter =
  function (message, writer) {
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
  proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_autoops_command_pb.DeleteProgressiveRolloutCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, '')
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
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.DeleteProgressiveRolloutRequest();
    return proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value =
            new proto_autoops_command_pb.DeleteProgressiveRolloutCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb.DeleteProgressiveRolloutCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_autoops_command_pb.DeleteProgressiveRolloutCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional DeleteProgressiveRolloutCommand command = 3;
 * @return {?proto.bucketeer.autoops.DeleteProgressiveRolloutCommand}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.DeleteProgressiveRolloutCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.DeleteProgressiveRolloutCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.DeleteProgressiveRolloutCommand|undefined} value
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutRequest.prototype.setEnvironmentId =
  function (value) {
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
  proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.DeleteProgressiveRolloutResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.DeleteProgressiveRolloutResponse();
    return proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.DeleteProgressiveRolloutResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.DeleteProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.deserializeBinaryFromReader =
  function (msg, reader) {
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
proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.DeleteProgressiveRolloutResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.DeleteProgressiveRolloutResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.repeatedFields_ = [4];

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
  proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.ListProgressiveRolloutsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ListProgressiveRolloutsRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
        cursor: jspb.Message.getFieldWithDefault(msg, 3, ''),
        featureIdsList:
          (f = jspb.Message.getRepeatedField(msg, 4)) == null ? undefined : f,
        orderBy: jspb.Message.getFieldWithDefault(msg, 5, 0),
        orderDirection: jspb.Message.getFieldWithDefault(msg, 6, 0),
        status: jspb.Message.getFieldWithDefault(msg, 7, 0),
        type: jspb.Message.getFieldWithDefault(msg, 8, 0),
        environmentId: jspb.Message.getFieldWithDefault(msg, 9, '')
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
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.ListProgressiveRolloutsRequest();
    return proto.bucketeer.autoops.ListProgressiveRolloutsRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.addFeatureIds(value);
          break;
        case 5:
          var value =
            /** @type {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderBy} */ (
              reader.readEnum()
            );
          msg.setOrderBy(value);
          break;
        case 6:
          var value =
            /** @type {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirection} */ (
              reader.readEnum()
            );
          msg.setOrderDirection(value);
          break;
        case 7:
          var value =
            /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Status} */ (
              reader.readEnum()
            );
          msg.setStatus(value);
          break;
        case 8:
          var value =
            /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Type} */ (
              reader.readEnum()
            );
          msg.setType(value);
          break;
        case 9:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ListProgressiveRolloutsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getPageSize();
    if (f !== 0) {
      writer.writeInt64(2, f);
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getFeatureIdsList();
    if (f.length > 0) {
      writer.writeRepeatedString(4, f);
    }
    f = message.getOrderBy();
    if (f !== 0.0) {
      writer.writeEnum(5, f);
    }
    f = message.getOrderDirection();
    if (f !== 0.0) {
      writer.writeEnum(6, f);
    }
    f = /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Status} */ (
      jspb.Message.getField(message, 7)
    );
    if (f != null) {
      writer.writeEnum(7, f);
    }
    f = /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Type} */ (
      jspb.Message.getField(message, 8)
    );
    if (f != null) {
      writer.writeEnum(8, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(9, f);
    }
  };

/**
 * @enum {number}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderBy = {
  DEFAULT: 0,
  CREATED_AT: 1,
  UPDATED_AT: 2
};

/**
 * @enum {number}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 2;
 * @return {number}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setPageSize =
  function (value) {
    return jspb.Message.setProto3IntField(this, 2, value);
  };

/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * repeated string feature_ids = 4;
 * @return {!Array<string>}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getFeatureIdsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 4)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setFeatureIdsList =
  function (value) {
    return jspb.Message.setField(this, 4, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.addFeatureIds =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 4, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.clearFeatureIdsList =
  function () {
    return this.setFeatureIdsList([]);
  };

/**
 * optional OrderBy order_by = 5;
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderBy}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getOrderBy =
  function () {
    return /** @type {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderBy} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderBy} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setOrderBy =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
  };

/**
 * optional OrderDirection order_direction = 6;
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirection}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getOrderDirection =
  function () {
    return /** @type {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirection} */ (
      jspb.Message.getFieldWithDefault(this, 6, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest.OrderDirection} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setOrderDirection =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 6, value);
  };

/**
 * optional ProgressiveRollout.Status status = 7;
 * @return {!proto.bucketeer.autoops.ProgressiveRollout.Status}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getStatus =
  function () {
    return /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Status} */ (
      jspb.Message.getFieldWithDefault(this, 7, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout.Status} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setStatus =
  function (value) {
    return jspb.Message.setField(this, 7, value);
  };

/**
 * Clears the field making it undefined.
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.clearStatus =
  function () {
    return jspb.Message.setField(this, 7, undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.hasStatus =
  function () {
    return jspb.Message.getField(this, 7) != null;
  };

/**
 * optional ProgressiveRollout.Type type = 8;
 * @return {!proto.bucketeer.autoops.ProgressiveRollout.Type}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getType =
  function () {
    return /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Type} */ (
      jspb.Message.getFieldWithDefault(this, 8, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout.Type} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setType =
  function (value) {
    return jspb.Message.setField(this, 8, value);
  };

/**
 * Clears the field making it undefined.
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.clearType =
  function () {
    return jspb.Message.setField(this, 8, undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.hasType =
  function () {
    return jspb.Message.getField(this, 8) != null;
  };

/**
 * optional string environment_id = 9;
 * @return {string}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 9, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsRequest} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 9, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.repeatedFields_ = [1];

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
  proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.ListProgressiveRolloutsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ListProgressiveRolloutsResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        progressiveRolloutsList: jspb.Message.toObjectList(
          msg.getProgressiveRolloutsList(),
          proto_autoops_progressive_rollout_pb.ProgressiveRollout.toObject,
          includeInstance
        ),
        cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
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
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.ListProgressiveRolloutsResponse();
    return proto.bucketeer.autoops.ListProgressiveRolloutsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value =
            new proto_autoops_progressive_rollout_pb.ProgressiveRollout();
          reader.readMessage(
            value,
            proto_autoops_progressive_rollout_pb.ProgressiveRollout
              .deserializeBinaryFromReader
          );
          msg.addProgressiveRollouts(value);
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
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ListProgressiveRolloutsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getProgressiveRolloutsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_autoops_progressive_rollout_pb.ProgressiveRollout
          .serializeBinaryToWriter
      );
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getTotalCount();
    if (f !== 0) {
      writer.writeInt64(3, f);
    }
  };

/**
 * repeated ProgressiveRollout progressive_rollouts = 1;
 * @return {!Array<!proto.bucketeer.autoops.ProgressiveRollout>}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.getProgressiveRolloutsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.ProgressiveRollout>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_autoops_progressive_rollout_pb.ProgressiveRollout,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.ProgressiveRollout>} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.setProgressiveRolloutsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.ProgressiveRollout}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.addProgressiveRollouts =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.autoops.ProgressiveRollout,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.clearProgressiveRolloutsList =
  function () {
    return this.setProgressiveRolloutsList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.getTotalCount =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ListProgressiveRolloutsResponse} returns this
 */
proto.bucketeer.autoops.ListProgressiveRolloutsResponse.prototype.setTotalCount =
  function (value) {
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
  proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        changeProgressiveRolloutTriggeredAtCommand:
          (f = msg.getChangeProgressiveRolloutTriggeredAtCommand()) &&
          proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, ''),
        scheduleId: jspb.Message.getFieldWithDefault(msg, 5, '')
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
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest();
    return proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value =
            new proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand();
          reader.readMessage(
            value,
            proto_autoops_command_pb
              .ChangeProgressiveRolloutScheduleTriggeredAtCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeProgressiveRolloutTriggeredAtCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 5:
          var value = /** @type {string} */ (reader.readString());
          msg.setScheduleId(value);
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
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getChangeProgressiveRolloutTriggeredAtCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_autoops_command_pb
          .ChangeProgressiveRolloutScheduleTriggeredAtCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
    f = message.getScheduleId();
    if (f.length > 0) {
      writer.writeString(5, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional ChangeProgressiveRolloutScheduleTriggeredAtCommand change_progressive_rollout_triggered_at_command = 3;
 * @return {?proto.bucketeer.autoops.ChangeProgressiveRolloutScheduleTriggeredAtCommand}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.getChangeProgressiveRolloutTriggeredAtCommand =
  function () {
    return /** @type{?proto.bucketeer.autoops.ChangeProgressiveRolloutScheduleTriggeredAtCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_autoops_command_pb.ChangeProgressiveRolloutScheduleTriggeredAtCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.autoops.ChangeProgressiveRolloutScheduleTriggeredAtCommand|undefined} value
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.setChangeProgressiveRolloutTriggeredAtCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.clearChangeProgressiveRolloutTriggeredAtCommand =
  function () {
    return this.setChangeProgressiveRolloutTriggeredAtCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.hasChangeProgressiveRolloutTriggeredAtCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
  };

/**
 * optional string schedule_id = 5;
 * @return {string}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.getScheduleId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 5, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest} returns this
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutRequest.prototype.setScheduleId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 5, value);
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
  proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {};

      if (includeInstance) {
        obj.$jspbMessageInstance = msg;
      }
      return obj;
    };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse();
    return proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse}
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.deserializeBinaryFromReader =
  function (msg, reader) {
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
proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ExecuteProgressiveRolloutResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
  };

/**
 * @enum {number}
 */
proto.bucketeer.autoops.ChangeType = {
  UNSPECIFIED: 0,
  CREATE: 1,
  UPDATE: 2,
  DELETE: 3
};

goog.object.extend(exports, proto.bucketeer.autoops);
