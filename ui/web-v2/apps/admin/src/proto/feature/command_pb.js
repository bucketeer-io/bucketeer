// source: proto/feature/command.proto
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

var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');
goog.object.extend(proto, google_protobuf_any_pb);
var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.object.extend(proto, google_protobuf_wrappers_pb);
var proto_feature_clause_pb = require('../../proto/feature/clause_pb.js');
goog.object.extend(proto, proto_feature_clause_pb);
var proto_feature_feature_pb = require('../../proto/feature/feature_pb.js');
goog.object.extend(proto, proto_feature_feature_pb);
var proto_feature_rule_pb = require('../../proto/feature/rule_pb.js');
goog.object.extend(proto, proto_feature_rule_pb);
var proto_feature_variation_pb = require('../../proto/feature/variation_pb.js');
goog.object.extend(proto, proto_feature_variation_pb);
var proto_feature_strategy_pb = require('../../proto/feature/strategy_pb.js');
goog.object.extend(proto, proto_feature_strategy_pb);
var proto_feature_segment_pb = require('../../proto/feature/segment_pb.js');
goog.object.extend(proto, proto_feature_segment_pb);
var proto_feature_prerequisite_pb = require('../../proto/feature/prerequisite_pb.js');
goog.object.extend(proto, proto_feature_prerequisite_pb);
var proto_feature_flag_trigger_pb = require('../../proto/feature/flag_trigger_pb.js');
goog.object.extend(proto, proto_feature_flag_trigger_pb);
goog.exportSymbol('proto.bucketeer.feature.AddClauseCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddClauseValueCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddPrerequisiteCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddRuleCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddSegmentUserCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddTagCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddUserToVariationCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.AddVariationCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ArchiveFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.BulkUploadSegmentUsersCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeClauseAttributeCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeClauseOperatorCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeDefaultStrategyCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeDescriptionCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeFixedStrategyCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeOffVariationCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangePrerequisiteVariationCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeRolloutStrategyCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeRuleStrategyCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeSegmentDescriptionCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeSegmentNameCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeVariationDescriptionCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeVariationNameCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ChangeVariationValueCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.CloneFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.Command', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateFlagTriggerCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.CreateSegmentCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteClauseCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteFlagTriggerCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteRuleCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteSegmentCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DeleteSegmentUserCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DisableFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.DisableFlagTriggerCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.EnableFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.EnableFlagTriggerCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.IncrementFeatureVersionCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.RemoveClauseValueCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.RemovePrerequisiteCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.RemoveTagCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.RemoveUserFromVariationCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.RemoveVariationCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.RenameFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ResetFlagTriggerCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.ResetSamplingSeedCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.UnarchiveFeatureCommand', null, global);
goog.exportSymbol('proto.bucketeer.feature.UpdateFlagTriggerUsageCommand', null, global);
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
proto.bucketeer.feature.Command = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.Command, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.Command.displayName = 'proto.bucketeer.feature.Command';
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
proto.bucketeer.feature.CreateFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.CreateFeatureCommand.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.CreateFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateFeatureCommand.displayName = 'proto.bucketeer.feature.CreateFeatureCommand';
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
proto.bucketeer.feature.ArchiveFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ArchiveFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ArchiveFeatureCommand.displayName = 'proto.bucketeer.feature.ArchiveFeatureCommand';
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
proto.bucketeer.feature.UnarchiveFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UnarchiveFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UnarchiveFeatureCommand.displayName = 'proto.bucketeer.feature.UnarchiveFeatureCommand';
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
proto.bucketeer.feature.DeleteFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteFeatureCommand.displayName = 'proto.bucketeer.feature.DeleteFeatureCommand';
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
proto.bucketeer.feature.RenameFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.RenameFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.RenameFeatureCommand.displayName = 'proto.bucketeer.feature.RenameFeatureCommand';
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
proto.bucketeer.feature.ChangeDescriptionCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeDescriptionCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeDescriptionCommand.displayName = 'proto.bucketeer.feature.ChangeDescriptionCommand';
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
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.displayName = 'proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand';
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
proto.bucketeer.feature.AddTagCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddTagCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddTagCommand.displayName = 'proto.bucketeer.feature.AddTagCommand';
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
proto.bucketeer.feature.RemoveTagCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.RemoveTagCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.RemoveTagCommand.displayName = 'proto.bucketeer.feature.RemoveTagCommand';
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
proto.bucketeer.feature.EnableFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EnableFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EnableFeatureCommand.displayName = 'proto.bucketeer.feature.EnableFeatureCommand';
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
proto.bucketeer.feature.DisableFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DisableFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DisableFeatureCommand.displayName = 'proto.bucketeer.feature.DisableFeatureCommand';
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
proto.bucketeer.feature.AddVariationCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddVariationCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddVariationCommand.displayName = 'proto.bucketeer.feature.AddVariationCommand';
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
proto.bucketeer.feature.RemoveVariationCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.RemoveVariationCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.RemoveVariationCommand.displayName = 'proto.bucketeer.feature.RemoveVariationCommand';
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
proto.bucketeer.feature.ChangeVariationValueCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeVariationValueCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeVariationValueCommand.displayName = 'proto.bucketeer.feature.ChangeVariationValueCommand';
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
proto.bucketeer.feature.ChangeVariationNameCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeVariationNameCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeVariationNameCommand.displayName = 'proto.bucketeer.feature.ChangeVariationNameCommand';
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
proto.bucketeer.feature.ChangeVariationDescriptionCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeVariationDescriptionCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeVariationDescriptionCommand.displayName = 'proto.bucketeer.feature.ChangeVariationDescriptionCommand';
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
proto.bucketeer.feature.ChangeOffVariationCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeOffVariationCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeOffVariationCommand.displayName = 'proto.bucketeer.feature.ChangeOffVariationCommand';
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
proto.bucketeer.feature.AddUserToVariationCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddUserToVariationCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddUserToVariationCommand.displayName = 'proto.bucketeer.feature.AddUserToVariationCommand';
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
proto.bucketeer.feature.RemoveUserFromVariationCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.RemoveUserFromVariationCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.RemoveUserFromVariationCommand.displayName = 'proto.bucketeer.feature.RemoveUserFromVariationCommand';
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
proto.bucketeer.feature.ChangeDefaultStrategyCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeDefaultStrategyCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeDefaultStrategyCommand.displayName = 'proto.bucketeer.feature.ChangeDefaultStrategyCommand';
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
proto.bucketeer.feature.AddRuleCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddRuleCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddRuleCommand.displayName = 'proto.bucketeer.feature.AddRuleCommand';
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
proto.bucketeer.feature.ChangeRuleStrategyCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeRuleStrategyCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeRuleStrategyCommand.displayName = 'proto.bucketeer.feature.ChangeRuleStrategyCommand';
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
proto.bucketeer.feature.DeleteRuleCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteRuleCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteRuleCommand.displayName = 'proto.bucketeer.feature.DeleteRuleCommand';
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
proto.bucketeer.feature.AddClauseCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddClauseCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddClauseCommand.displayName = 'proto.bucketeer.feature.AddClauseCommand';
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
proto.bucketeer.feature.DeleteClauseCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteClauseCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteClauseCommand.displayName = 'proto.bucketeer.feature.DeleteClauseCommand';
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
proto.bucketeer.feature.ChangeClauseAttributeCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeClauseAttributeCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeClauseAttributeCommand.displayName = 'proto.bucketeer.feature.ChangeClauseAttributeCommand';
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
proto.bucketeer.feature.ChangeClauseOperatorCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeClauseOperatorCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeClauseOperatorCommand.displayName = 'proto.bucketeer.feature.ChangeClauseOperatorCommand';
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
proto.bucketeer.feature.AddClauseValueCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddClauseValueCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddClauseValueCommand.displayName = 'proto.bucketeer.feature.AddClauseValueCommand';
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
proto.bucketeer.feature.RemoveClauseValueCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.RemoveClauseValueCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.RemoveClauseValueCommand.displayName = 'proto.bucketeer.feature.RemoveClauseValueCommand';
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
proto.bucketeer.feature.ChangeFixedStrategyCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeFixedStrategyCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeFixedStrategyCommand.displayName = 'proto.bucketeer.feature.ChangeFixedStrategyCommand';
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
proto.bucketeer.feature.ChangeRolloutStrategyCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeRolloutStrategyCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeRolloutStrategyCommand.displayName = 'proto.bucketeer.feature.ChangeRolloutStrategyCommand';
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
proto.bucketeer.feature.CreateSegmentCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateSegmentCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateSegmentCommand.displayName = 'proto.bucketeer.feature.CreateSegmentCommand';
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
proto.bucketeer.feature.DeleteSegmentCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteSegmentCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteSegmentCommand.displayName = 'proto.bucketeer.feature.DeleteSegmentCommand';
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
proto.bucketeer.feature.ChangeSegmentNameCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeSegmentNameCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeSegmentNameCommand.displayName = 'proto.bucketeer.feature.ChangeSegmentNameCommand';
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
proto.bucketeer.feature.ChangeSegmentDescriptionCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeSegmentDescriptionCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeSegmentDescriptionCommand.displayName = 'proto.bucketeer.feature.ChangeSegmentDescriptionCommand';
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
proto.bucketeer.feature.AddSegmentUserCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.AddSegmentUserCommand.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.AddSegmentUserCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddSegmentUserCommand.displayName = 'proto.bucketeer.feature.AddSegmentUserCommand';
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
proto.bucketeer.feature.DeleteSegmentUserCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.feature.DeleteSegmentUserCommand.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.feature.DeleteSegmentUserCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteSegmentUserCommand.displayName = 'proto.bucketeer.feature.DeleteSegmentUserCommand';
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
proto.bucketeer.feature.BulkUploadSegmentUsersCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.BulkUploadSegmentUsersCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.BulkUploadSegmentUsersCommand.displayName = 'proto.bucketeer.feature.BulkUploadSegmentUsersCommand';
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
proto.bucketeer.feature.IncrementFeatureVersionCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.IncrementFeatureVersionCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.IncrementFeatureVersionCommand.displayName = 'proto.bucketeer.feature.IncrementFeatureVersionCommand';
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
proto.bucketeer.feature.CloneFeatureCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CloneFeatureCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CloneFeatureCommand.displayName = 'proto.bucketeer.feature.CloneFeatureCommand';
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
proto.bucketeer.feature.ResetSamplingSeedCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ResetSamplingSeedCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ResetSamplingSeedCommand.displayName = 'proto.bucketeer.feature.ResetSamplingSeedCommand';
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
proto.bucketeer.feature.AddPrerequisiteCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.AddPrerequisiteCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.AddPrerequisiteCommand.displayName = 'proto.bucketeer.feature.AddPrerequisiteCommand';
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
proto.bucketeer.feature.RemovePrerequisiteCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.RemovePrerequisiteCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.RemovePrerequisiteCommand.displayName = 'proto.bucketeer.feature.RemovePrerequisiteCommand';
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
proto.bucketeer.feature.ChangePrerequisiteVariationCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangePrerequisiteVariationCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangePrerequisiteVariationCommand.displayName = 'proto.bucketeer.feature.ChangePrerequisiteVariationCommand';
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
proto.bucketeer.feature.CreateFlagTriggerCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.CreateFlagTriggerCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.CreateFlagTriggerCommand.displayName = 'proto.bucketeer.feature.CreateFlagTriggerCommand';
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
proto.bucketeer.feature.ResetFlagTriggerCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ResetFlagTriggerCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ResetFlagTriggerCommand.displayName = 'proto.bucketeer.feature.ResetFlagTriggerCommand';
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
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.displayName = 'proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand';
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
proto.bucketeer.feature.EnableFlagTriggerCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.EnableFlagTriggerCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.EnableFlagTriggerCommand.displayName = 'proto.bucketeer.feature.EnableFlagTriggerCommand';
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
proto.bucketeer.feature.DisableFlagTriggerCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DisableFlagTriggerCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DisableFlagTriggerCommand.displayName = 'proto.bucketeer.feature.DisableFlagTriggerCommand';
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
proto.bucketeer.feature.DeleteFlagTriggerCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.DeleteFlagTriggerCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.DeleteFlagTriggerCommand.displayName = 'proto.bucketeer.feature.DeleteFlagTriggerCommand';
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
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.UpdateFlagTriggerUsageCommand, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.displayName = 'proto.bucketeer.feature.UpdateFlagTriggerUsageCommand';
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
proto.bucketeer.feature.Command.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.Command.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.Command} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.Command.toObject = function(includeInstance, msg) {
  var f, obj = {
    command: (f = msg.getCommand()) && google_protobuf_any_pb.Any.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.Command}
 */
proto.bucketeer.feature.Command.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.Command;
  return proto.bucketeer.feature.Command.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.Command} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.Command}
 */
proto.bucketeer.feature.Command.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
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
proto.bucketeer.feature.Command.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.Command.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.Command} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.Command.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.Any command = 1;
 * @return {?proto.google.protobuf.Any}
 */
proto.bucketeer.feature.Command.prototype.getCommand = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 1));
};


/**
 * @param {?proto.google.protobuf.Any|undefined} value
 * @return {!proto.bucketeer.feature.Command} returns this
*/
proto.bucketeer.feature.Command.prototype.setCommand = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.Command} returns this
 */
proto.bucketeer.feature.Command.prototype.clearCommand = function() {
  return this.setCommand(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.Command.prototype.hasCommand = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.CreateFeatureCommand.repeatedFields_ = [4,5];



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
proto.bucketeer.feature.CreateFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFeatureCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, ""),
    variationsList: jspb.Message.toObjectList(msg.getVariationsList(),
    proto_feature_variation_pb.Variation.toObject, includeInstance),
    tagsList: (f = jspb.Message.getRepeatedField(msg, 5)) == null ? undefined : f,
    defaultOnVariationIndex: (f = msg.getDefaultOnVariationIndex()) && google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
    defaultOffVariationIndex: (f = msg.getDefaultOffVariationIndex()) && google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
    variationType: jspb.Message.getFieldWithDefault(msg, 8, 0)
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
 * @return {!proto.bucketeer.feature.CreateFeatureCommand}
 */
proto.bucketeer.feature.CreateFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateFeatureCommand;
  return proto.bucketeer.feature.CreateFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateFeatureCommand}
 */
proto.bucketeer.feature.CreateFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 4:
      var value = new proto_feature_variation_pb.Variation;
      reader.readMessage(value,proto_feature_variation_pb.Variation.deserializeBinaryFromReader);
      msg.addVariations(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.addTags(value);
      break;
    case 6:
      var value = new google_protobuf_wrappers_pb.Int32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.Int32Value.deserializeBinaryFromReader);
      msg.setDefaultOnVariationIndex(value);
      break;
    case 7:
      var value = new google_protobuf_wrappers_pb.Int32Value;
      reader.readMessage(value,google_protobuf_wrappers_pb.Int32Value.deserializeBinaryFromReader);
      msg.setDefaultOffVariationIndex(value);
      break;
    case 8:
      var value = /** @type {!proto.bucketeer.feature.Feature.VariationType} */ (reader.readEnum());
      msg.setVariationType(value);
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
proto.bucketeer.feature.CreateFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFeatureCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getVariationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto_feature_variation_pb.Variation.serializeBinaryToWriter
    );
  }
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      5,
      f
    );
  }
  f = message.getDefaultOnVariationIndex();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_wrappers_pb.Int32Value.serializeBinaryToWriter
    );
  }
  f = message.getDefaultOffVariationIndex();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_wrappers_pb.Int32Value.serializeBinaryToWriter
    );
  }
  f = message.getVariationType();
  if (f !== 0.0) {
    writer.writeEnum(
      8,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated Variation variations = 4;
 * @return {!Array<!proto.bucketeer.feature.Variation>}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getVariationsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Variation>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_variation_pb.Variation, 4));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Variation>} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
*/
proto.bucketeer.feature.CreateFeatureCommand.prototype.setVariationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.bucketeer.feature.Variation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Variation}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.addVariations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.bucketeer.feature.Variation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.clearVariationsList = function() {
  return this.setVariationsList([]);
};


/**
 * repeated string tags = 5;
 * @return {!Array<string>}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 5));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 5, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 5, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * optional google.protobuf.Int32Value default_on_variation_index = 6;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getDefaultOnVariationIndex = function() {
  return /** @type{?proto.google.protobuf.Int32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.Int32Value, 6));
};


/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
*/
proto.bucketeer.feature.CreateFeatureCommand.prototype.setDefaultOnVariationIndex = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.clearDefaultOnVariationIndex = function() {
  return this.setDefaultOnVariationIndex(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.hasDefaultOnVariationIndex = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional google.protobuf.Int32Value default_off_variation_index = 7;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getDefaultOffVariationIndex = function() {
  return /** @type{?proto.google.protobuf.Int32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.Int32Value, 7));
};


/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
*/
proto.bucketeer.feature.CreateFeatureCommand.prototype.setDefaultOffVariationIndex = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.clearDefaultOffVariationIndex = function() {
  return this.setDefaultOffVariationIndex(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.hasDefaultOffVariationIndex = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional Feature.VariationType variation_type = 8;
 * @return {!proto.bucketeer.feature.Feature.VariationType}
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.getVariationType = function() {
  return /** @type {!proto.bucketeer.feature.Feature.VariationType} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {!proto.bucketeer.feature.Feature.VariationType} value
 * @return {!proto.bucketeer.feature.CreateFeatureCommand} returns this
 */
proto.bucketeer.feature.CreateFeatureCommand.prototype.setVariationType = function(value) {
  return jspb.Message.setProto3EnumField(this, 8, value);
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
proto.bucketeer.feature.ArchiveFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ArchiveFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ArchiveFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ArchiveFeatureCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.ArchiveFeatureCommand}
 */
proto.bucketeer.feature.ArchiveFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ArchiveFeatureCommand;
  return proto.bucketeer.feature.ArchiveFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ArchiveFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ArchiveFeatureCommand}
 */
proto.bucketeer.feature.ArchiveFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.ArchiveFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ArchiveFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ArchiveFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ArchiveFeatureCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.UnarchiveFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UnarchiveFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UnarchiveFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UnarchiveFeatureCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UnarchiveFeatureCommand}
 */
proto.bucketeer.feature.UnarchiveFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UnarchiveFeatureCommand;
  return proto.bucketeer.feature.UnarchiveFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UnarchiveFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UnarchiveFeatureCommand}
 */
proto.bucketeer.feature.UnarchiveFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UnarchiveFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UnarchiveFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UnarchiveFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UnarchiveFeatureCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DeleteFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFeatureCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteFeatureCommand}
 */
proto.bucketeer.feature.DeleteFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteFeatureCommand;
  return proto.bucketeer.feature.DeleteFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteFeatureCommand}
 */
proto.bucketeer.feature.DeleteFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFeatureCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.RenameFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.RenameFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.RenameFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RenameFeatureCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.RenameFeatureCommand}
 */
proto.bucketeer.feature.RenameFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.RenameFeatureCommand;
  return proto.bucketeer.feature.RenameFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.RenameFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.RenameFeatureCommand}
 */
proto.bucketeer.feature.RenameFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
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
proto.bucketeer.feature.RenameFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.RenameFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.RenameFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RenameFeatureCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.bucketeer.feature.RenameFeatureCommand.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RenameFeatureCommand} returns this
 */
proto.bucketeer.feature.RenameFeatureCommand.prototype.setName = function(value) {
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
proto.bucketeer.feature.ChangeDescriptionCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeDescriptionCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeDescriptionCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeDescriptionCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    description: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.ChangeDescriptionCommand}
 */
proto.bucketeer.feature.ChangeDescriptionCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeDescriptionCommand;
  return proto.bucketeer.feature.ChangeDescriptionCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeDescriptionCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeDescriptionCommand}
 */
proto.bucketeer.feature.ChangeDescriptionCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
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
proto.bucketeer.feature.ChangeDescriptionCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeDescriptionCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeDescriptionCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeDescriptionCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string description = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeDescriptionCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeDescriptionCommand} returns this
 */
proto.bucketeer.feature.ChangeDescriptionCommand.prototype.setDescription = function(value) {
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
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    status: jspb.Message.getFieldWithDefault(msg, 1, 0),
    state: jspb.Message.getFieldWithDefault(msg, 2, 0),
    count: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand}
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand;
  return proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand}
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.bucketeer.feature.Segment.Status} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 2:
      var value = /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCount(value);
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
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getCount();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
};


/**
 * optional Segment.Status status = 1;
 * @return {!proto.bucketeer.feature.Segment.Status}
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.getStatus = function() {
  return /** @type {!proto.bucketeer.feature.Segment.Status} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.bucketeer.feature.Segment.Status} value
 * @return {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand} returns this
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional SegmentUser.State state = 2;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand} returns this
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional int64 count = 3;
 * @return {number}
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.getCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand} returns this
 */
proto.bucketeer.feature.ChangeBulkUploadSegmentUsersStatusCommand.prototype.setCount = function(value) {
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
proto.bucketeer.feature.AddTagCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddTagCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddTagCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddTagCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    tag: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.AddTagCommand}
 */
proto.bucketeer.feature.AddTagCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddTagCommand;
  return proto.bucketeer.feature.AddTagCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddTagCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddTagCommand}
 */
proto.bucketeer.feature.AddTagCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setTag(value);
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
proto.bucketeer.feature.AddTagCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddTagCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddTagCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddTagCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTag();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string tag = 1;
 * @return {string}
 */
proto.bucketeer.feature.AddTagCommand.prototype.getTag = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddTagCommand} returns this
 */
proto.bucketeer.feature.AddTagCommand.prototype.setTag = function(value) {
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
proto.bucketeer.feature.RemoveTagCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.RemoveTagCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.RemoveTagCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveTagCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    tag: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.RemoveTagCommand}
 */
proto.bucketeer.feature.RemoveTagCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.RemoveTagCommand;
  return proto.bucketeer.feature.RemoveTagCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.RemoveTagCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.RemoveTagCommand}
 */
proto.bucketeer.feature.RemoveTagCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setTag(value);
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
proto.bucketeer.feature.RemoveTagCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.RemoveTagCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.RemoveTagCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveTagCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTag();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string tag = 1;
 * @return {string}
 */
proto.bucketeer.feature.RemoveTagCommand.prototype.getTag = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveTagCommand} returns this
 */
proto.bucketeer.feature.RemoveTagCommand.prototype.setTag = function(value) {
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
proto.bucketeer.feature.EnableFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EnableFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EnableFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFeatureCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.EnableFeatureCommand}
 */
proto.bucketeer.feature.EnableFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EnableFeatureCommand;
  return proto.bucketeer.feature.EnableFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EnableFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EnableFeatureCommand}
 */
proto.bucketeer.feature.EnableFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.EnableFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EnableFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EnableFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFeatureCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DisableFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DisableFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DisableFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFeatureCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DisableFeatureCommand}
 */
proto.bucketeer.feature.DisableFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DisableFeatureCommand;
  return proto.bucketeer.feature.DisableFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DisableFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DisableFeatureCommand}
 */
proto.bucketeer.feature.DisableFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DisableFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DisableFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DisableFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFeatureCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.AddVariationCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddVariationCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddVariationCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddVariationCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    value: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.AddVariationCommand}
 */
proto.bucketeer.feature.AddVariationCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddVariationCommand;
  return proto.bucketeer.feature.AddVariationCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddVariationCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddVariationCommand}
 */
proto.bucketeer.feature.AddVariationCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
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
proto.bucketeer.feature.AddVariationCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddVariationCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddVariationCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddVariationCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string value = 1;
 * @return {string}
 */
proto.bucketeer.feature.AddVariationCommand.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddVariationCommand} returns this
 */
proto.bucketeer.feature.AddVariationCommand.prototype.setValue = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.feature.AddVariationCommand.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddVariationCommand} returns this
 */
proto.bucketeer.feature.AddVariationCommand.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.feature.AddVariationCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddVariationCommand} returns this
 */
proto.bucketeer.feature.AddVariationCommand.prototype.setDescription = function(value) {
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
proto.bucketeer.feature.RemoveVariationCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.RemoveVariationCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.RemoveVariationCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveVariationCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.RemoveVariationCommand}
 */
proto.bucketeer.feature.RemoveVariationCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.RemoveVariationCommand;
  return proto.bucketeer.feature.RemoveVariationCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.RemoveVariationCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.RemoveVariationCommand}
 */
proto.bucketeer.feature.RemoveVariationCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.RemoveVariationCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.RemoveVariationCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.RemoveVariationCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveVariationCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.RemoveVariationCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveVariationCommand} returns this
 */
proto.bucketeer.feature.RemoveVariationCommand.prototype.setId = function(value) {
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
proto.bucketeer.feature.ChangeVariationValueCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeVariationValueCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeVariationValueCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeVariationValueCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    value: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ChangeVariationValueCommand}
 */
proto.bucketeer.feature.ChangeVariationValueCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeVariationValueCommand;
  return proto.bucketeer.feature.ChangeVariationValueCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeVariationValueCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeVariationValueCommand}
 */
proto.bucketeer.feature.ChangeVariationValueCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setValue(value);
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
proto.bucketeer.feature.ChangeVariationValueCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeVariationValueCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeVariationValueCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeVariationValueCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.feature.ChangeVariationValueCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeVariationValueCommand} returns this
 */
proto.bucketeer.feature.ChangeVariationValueCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string value = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeVariationValueCommand.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeVariationValueCommand} returns this
 */
proto.bucketeer.feature.ChangeVariationValueCommand.prototype.setValue = function(value) {
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
proto.bucketeer.feature.ChangeVariationNameCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeVariationNameCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeVariationNameCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeVariationNameCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ChangeVariationNameCommand}
 */
proto.bucketeer.feature.ChangeVariationNameCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeVariationNameCommand;
  return proto.bucketeer.feature.ChangeVariationNameCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeVariationNameCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeVariationNameCommand}
 */
proto.bucketeer.feature.ChangeVariationNameCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setName(value);
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
proto.bucketeer.feature.ChangeVariationNameCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeVariationNameCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeVariationNameCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeVariationNameCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
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
proto.bucketeer.feature.ChangeVariationNameCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeVariationNameCommand} returns this
 */
proto.bucketeer.feature.ChangeVariationNameCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeVariationNameCommand.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeVariationNameCommand} returns this
 */
proto.bucketeer.feature.ChangeVariationNameCommand.prototype.setName = function(value) {
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
proto.bucketeer.feature.ChangeVariationDescriptionCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeVariationDescriptionCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeVariationDescriptionCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    description: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.ChangeVariationDescriptionCommand}
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeVariationDescriptionCommand;
  return proto.bucketeer.feature.ChangeVariationDescriptionCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeVariationDescriptionCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeVariationDescriptionCommand}
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setDescription(value);
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
proto.bucketeer.feature.ChangeVariationDescriptionCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeVariationDescriptionCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeVariationDescriptionCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDescription();
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
proto.bucketeer.feature.ChangeVariationDescriptionCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeVariationDescriptionCommand} returns this
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeVariationDescriptionCommand} returns this
 */
proto.bucketeer.feature.ChangeVariationDescriptionCommand.prototype.setDescription = function(value) {
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
proto.bucketeer.feature.ChangeOffVariationCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeOffVariationCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeOffVariationCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeOffVariationCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.ChangeOffVariationCommand}
 */
proto.bucketeer.feature.ChangeOffVariationCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeOffVariationCommand;
  return proto.bucketeer.feature.ChangeOffVariationCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeOffVariationCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeOffVariationCommand}
 */
proto.bucketeer.feature.ChangeOffVariationCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.ChangeOffVariationCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeOffVariationCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeOffVariationCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeOffVariationCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.ChangeOffVariationCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeOffVariationCommand} returns this
 */
proto.bucketeer.feature.ChangeOffVariationCommand.prototype.setId = function(value) {
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
proto.bucketeer.feature.AddUserToVariationCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddUserToVariationCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddUserToVariationCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddUserToVariationCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    user: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.AddUserToVariationCommand}
 */
proto.bucketeer.feature.AddUserToVariationCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddUserToVariationCommand;
  return proto.bucketeer.feature.AddUserToVariationCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddUserToVariationCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddUserToVariationCommand}
 */
proto.bucketeer.feature.AddUserToVariationCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.AddUserToVariationCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddUserToVariationCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddUserToVariationCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddUserToVariationCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUser();
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
proto.bucketeer.feature.AddUserToVariationCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddUserToVariationCommand} returns this
 */
proto.bucketeer.feature.AddUserToVariationCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string user = 2;
 * @return {string}
 */
proto.bucketeer.feature.AddUserToVariationCommand.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddUserToVariationCommand} returns this
 */
proto.bucketeer.feature.AddUserToVariationCommand.prototype.setUser = function(value) {
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
proto.bucketeer.feature.RemoveUserFromVariationCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.RemoveUserFromVariationCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.RemoveUserFromVariationCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    user: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.RemoveUserFromVariationCommand}
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.RemoveUserFromVariationCommand;
  return proto.bucketeer.feature.RemoveUserFromVariationCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.RemoveUserFromVariationCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.RemoveUserFromVariationCommand}
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.RemoveUserFromVariationCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.RemoveUserFromVariationCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.RemoveUserFromVariationCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUser();
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
proto.bucketeer.feature.RemoveUserFromVariationCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveUserFromVariationCommand} returns this
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string user = 2;
 * @return {string}
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveUserFromVariationCommand} returns this
 */
proto.bucketeer.feature.RemoveUserFromVariationCommand.prototype.setUser = function(value) {
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
proto.bucketeer.feature.ChangeDefaultStrategyCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeDefaultStrategyCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeDefaultStrategyCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    strategy: (f = msg.getStrategy()) && proto_feature_strategy_pb.Strategy.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ChangeDefaultStrategyCommand}
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeDefaultStrategyCommand;
  return proto.bucketeer.feature.ChangeDefaultStrategyCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeDefaultStrategyCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeDefaultStrategyCommand}
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_strategy_pb.Strategy;
      reader.readMessage(value,proto_feature_strategy_pb.Strategy.deserializeBinaryFromReader);
      msg.setStrategy(value);
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
proto.bucketeer.feature.ChangeDefaultStrategyCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeDefaultStrategyCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeDefaultStrategyCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStrategy();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_strategy_pb.Strategy.serializeBinaryToWriter
    );
  }
};


/**
 * optional Strategy strategy = 1;
 * @return {?proto.bucketeer.feature.Strategy}
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.Strategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.Strategy, 1));
};


/**
 * @param {?proto.bucketeer.feature.Strategy|undefined} value
 * @return {!proto.bucketeer.feature.ChangeDefaultStrategyCommand} returns this
*/
proto.bucketeer.feature.ChangeDefaultStrategyCommand.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ChangeDefaultStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ChangeDefaultStrategyCommand.prototype.hasStrategy = function() {
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
proto.bucketeer.feature.AddRuleCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddRuleCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddRuleCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddRuleCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    rule: (f = msg.getRule()) && proto_feature_rule_pb.Rule.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.AddRuleCommand}
 */
proto.bucketeer.feature.AddRuleCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddRuleCommand;
  return proto.bucketeer.feature.AddRuleCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddRuleCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddRuleCommand}
 */
proto.bucketeer.feature.AddRuleCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_rule_pb.Rule;
      reader.readMessage(value,proto_feature_rule_pb.Rule.deserializeBinaryFromReader);
      msg.setRule(value);
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
proto.bucketeer.feature.AddRuleCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddRuleCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddRuleCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddRuleCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRule();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_rule_pb.Rule.serializeBinaryToWriter
    );
  }
};


/**
 * optional Rule rule = 1;
 * @return {?proto.bucketeer.feature.Rule}
 */
proto.bucketeer.feature.AddRuleCommand.prototype.getRule = function() {
  return /** @type{?proto.bucketeer.feature.Rule} */ (
    jspb.Message.getWrapperField(this, proto_feature_rule_pb.Rule, 1));
};


/**
 * @param {?proto.bucketeer.feature.Rule|undefined} value
 * @return {!proto.bucketeer.feature.AddRuleCommand} returns this
*/
proto.bucketeer.feature.AddRuleCommand.prototype.setRule = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.AddRuleCommand} returns this
 */
proto.bucketeer.feature.AddRuleCommand.prototype.clearRule = function() {
  return this.setRule(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.AddRuleCommand.prototype.hasRule = function() {
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
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeRuleStrategyCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeRuleStrategyCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    strategy: (f = msg.getStrategy()) && proto_feature_strategy_pb.Strategy.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ChangeRuleStrategyCommand}
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeRuleStrategyCommand;
  return proto.bucketeer.feature.ChangeRuleStrategyCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeRuleStrategyCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeRuleStrategyCommand}
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = new proto_feature_strategy_pb.Strategy;
      reader.readMessage(value,proto_feature_strategy_pb.Strategy.deserializeBinaryFromReader);
      msg.setStrategy(value);
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
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeRuleStrategyCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeRuleStrategyCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getStrategy();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_strategy_pb.Strategy.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeRuleStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeRuleStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional Strategy strategy = 3;
 * @return {?proto.bucketeer.feature.Strategy}
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.Strategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.Strategy, 3));
};


/**
 * @param {?proto.bucketeer.feature.Strategy|undefined} value
 * @return {!proto.bucketeer.feature.ChangeRuleStrategyCommand} returns this
*/
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ChangeRuleStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ChangeRuleStrategyCommand.prototype.hasStrategy = function() {
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
proto.bucketeer.feature.DeleteRuleCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteRuleCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteRuleCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteRuleCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteRuleCommand}
 */
proto.bucketeer.feature.DeleteRuleCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteRuleCommand;
  return proto.bucketeer.feature.DeleteRuleCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteRuleCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteRuleCommand}
 */
proto.bucketeer.feature.DeleteRuleCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteRuleCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteRuleCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteRuleCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteRuleCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DeleteRuleCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteRuleCommand} returns this
 */
proto.bucketeer.feature.DeleteRuleCommand.prototype.setId = function(value) {
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
proto.bucketeer.feature.AddClauseCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddClauseCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddClauseCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddClauseCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    ruleId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clause: (f = msg.getClause()) && proto_feature_clause_pb.Clause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.AddClauseCommand}
 */
proto.bucketeer.feature.AddClauseCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddClauseCommand;
  return proto.bucketeer.feature.AddClauseCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddClauseCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddClauseCommand}
 */
proto.bucketeer.feature.AddClauseCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRuleId(value);
      break;
    case 2:
      var value = new proto_feature_clause_pb.Clause;
      reader.readMessage(value,proto_feature_clause_pb.Clause.deserializeBinaryFromReader);
      msg.setClause(value);
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
proto.bucketeer.feature.AddClauseCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddClauseCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddClauseCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddClauseCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_clause_pb.Clause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string rule_id = 1;
 * @return {string}
 */
proto.bucketeer.feature.AddClauseCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddClauseCommand} returns this
 */
proto.bucketeer.feature.AddClauseCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Clause clause = 2;
 * @return {?proto.bucketeer.feature.Clause}
 */
proto.bucketeer.feature.AddClauseCommand.prototype.getClause = function() {
  return /** @type{?proto.bucketeer.feature.Clause} */ (
    jspb.Message.getWrapperField(this, proto_feature_clause_pb.Clause, 2));
};


/**
 * @param {?proto.bucketeer.feature.Clause|undefined} value
 * @return {!proto.bucketeer.feature.AddClauseCommand} returns this
*/
proto.bucketeer.feature.AddClauseCommand.prototype.setClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.AddClauseCommand} returns this
 */
proto.bucketeer.feature.AddClauseCommand.prototype.clearClause = function() {
  return this.setClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.AddClauseCommand.prototype.hasClause = function() {
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
proto.bucketeer.feature.DeleteClauseCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteClauseCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteClauseCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteClauseCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.DeleteClauseCommand}
 */
proto.bucketeer.feature.DeleteClauseCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteClauseCommand;
  return proto.bucketeer.feature.DeleteClauseCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteClauseCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteClauseCommand}
 */
proto.bucketeer.feature.DeleteClauseCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
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
proto.bucketeer.feature.DeleteClauseCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteClauseCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteClauseCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteClauseCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
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
proto.bucketeer.feature.DeleteClauseCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteClauseCommand} returns this
 */
proto.bucketeer.feature.DeleteClauseCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.DeleteClauseCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.DeleteClauseCommand} returns this
 */
proto.bucketeer.feature.DeleteClauseCommand.prototype.setRuleId = function(value) {
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
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeClauseAttributeCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeClauseAttributeCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    attribute: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.ChangeClauseAttributeCommand}
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeClauseAttributeCommand;
  return proto.bucketeer.feature.ChangeClauseAttributeCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeClauseAttributeCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeClauseAttributeCommand}
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setAttribute(value);
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
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeClauseAttributeCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeClauseAttributeCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getAttribute();
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
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeClauseAttributeCommand} returns this
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeClauseAttributeCommand} returns this
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string attribute = 3;
 * @return {string}
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.getAttribute = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeClauseAttributeCommand} returns this
 */
proto.bucketeer.feature.ChangeClauseAttributeCommand.prototype.setAttribute = function(value) {
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
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeClauseOperatorCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeClauseOperatorCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    operator: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.feature.ChangeClauseOperatorCommand}
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeClauseOperatorCommand;
  return proto.bucketeer.feature.ChangeClauseOperatorCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeClauseOperatorCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeClauseOperatorCommand}
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.feature.Clause.Operator} */ (reader.readEnum());
      msg.setOperator(value);
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
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeClauseOperatorCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeClauseOperatorCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getOperator();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeClauseOperatorCommand} returns this
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeClauseOperatorCommand} returns this
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional Clause.Operator operator = 3;
 * @return {!proto.bucketeer.feature.Clause.Operator}
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.getOperator = function() {
  return /** @type {!proto.bucketeer.feature.Clause.Operator} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.Clause.Operator} value
 * @return {!proto.bucketeer.feature.ChangeClauseOperatorCommand} returns this
 */
proto.bucketeer.feature.ChangeClauseOperatorCommand.prototype.setOperator = function(value) {
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
proto.bucketeer.feature.AddClauseValueCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddClauseValueCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddClauseValueCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddClauseValueCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    value: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.AddClauseValueCommand}
 */
proto.bucketeer.feature.AddClauseValueCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddClauseValueCommand;
  return proto.bucketeer.feature.AddClauseValueCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddClauseValueCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddClauseValueCommand}
 */
proto.bucketeer.feature.AddClauseValueCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
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
proto.bucketeer.feature.AddClauseValueCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddClauseValueCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddClauseValueCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddClauseValueCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.feature.AddClauseValueCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddClauseValueCommand} returns this
 */
proto.bucketeer.feature.AddClauseValueCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.AddClauseValueCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddClauseValueCommand} returns this
 */
proto.bucketeer.feature.AddClauseValueCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string value = 3;
 * @return {string}
 */
proto.bucketeer.feature.AddClauseValueCommand.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.AddClauseValueCommand} returns this
 */
proto.bucketeer.feature.AddClauseValueCommand.prototype.setValue = function(value) {
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
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.RemoveClauseValueCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.RemoveClauseValueCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveClauseValueCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    value: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.feature.RemoveClauseValueCommand}
 */
proto.bucketeer.feature.RemoveClauseValueCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.RemoveClauseValueCommand;
  return proto.bucketeer.feature.RemoveClauseValueCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.RemoveClauseValueCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.RemoveClauseValueCommand}
 */
proto.bucketeer.feature.RemoveClauseValueCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
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
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.RemoveClauseValueCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.RemoveClauseValueCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemoveClauseValueCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveClauseValueCommand} returns this
 */
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveClauseValueCommand} returns this
 */
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string value = 3;
 * @return {string}
 */
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemoveClauseValueCommand} returns this
 */
proto.bucketeer.feature.RemoveClauseValueCommand.prototype.setValue = function(value) {
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
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeFixedStrategyCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeFixedStrategyCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    strategy: (f = msg.getStrategy()) && proto_feature_strategy_pb.FixedStrategy.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ChangeFixedStrategyCommand}
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeFixedStrategyCommand;
  return proto.bucketeer.feature.ChangeFixedStrategyCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeFixedStrategyCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeFixedStrategyCommand}
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = new proto_feature_strategy_pb.FixedStrategy;
      reader.readMessage(value,proto_feature_strategy_pb.FixedStrategy.deserializeBinaryFromReader);
      msg.setStrategy(value);
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
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeFixedStrategyCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeFixedStrategyCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getStrategy();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_strategy_pb.FixedStrategy.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeFixedStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeFixedStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional FixedStrategy strategy = 3;
 * @return {?proto.bucketeer.feature.FixedStrategy}
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.FixedStrategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.FixedStrategy, 3));
};


/**
 * @param {?proto.bucketeer.feature.FixedStrategy|undefined} value
 * @return {!proto.bucketeer.feature.ChangeFixedStrategyCommand} returns this
*/
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ChangeFixedStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ChangeFixedStrategyCommand.prototype.hasStrategy = function() {
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
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeRolloutStrategyCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    strategy: (f = msg.getStrategy()) && proto_feature_strategy_pb.RolloutStrategy.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ChangeRolloutStrategyCommand}
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeRolloutStrategyCommand;
  return proto.bucketeer.feature.ChangeRolloutStrategyCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeRolloutStrategyCommand}
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = new proto_feature_strategy_pb.RolloutStrategy;
      reader.readMessage(value,proto_feature_strategy_pb.RolloutStrategy.deserializeBinaryFromReader);
      msg.setStrategy(value);
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
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeRolloutStrategyCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRuleId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getStrategy();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_strategy_pb.RolloutStrategy.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional RolloutStrategy strategy = 3;
 * @return {?proto.bucketeer.feature.RolloutStrategy}
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.RolloutStrategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.RolloutStrategy, 3));
};


/**
 * @param {?proto.bucketeer.feature.RolloutStrategy|undefined} value
 * @return {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} returns this
*/
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ChangeRolloutStrategyCommand} returns this
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ChangeRolloutStrategyCommand.prototype.hasStrategy = function() {
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
proto.bucketeer.feature.CreateSegmentCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateSegmentCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateSegmentCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateSegmentCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    description: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.feature.CreateSegmentCommand}
 */
proto.bucketeer.feature.CreateSegmentCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateSegmentCommand;
  return proto.bucketeer.feature.CreateSegmentCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateSegmentCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateSegmentCommand}
 */
proto.bucketeer.feature.CreateSegmentCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
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
proto.bucketeer.feature.CreateSegmentCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateSegmentCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateSegmentCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateSegmentCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.bucketeer.feature.CreateSegmentCommand.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateSegmentCommand} returns this
 */
proto.bucketeer.feature.CreateSegmentCommand.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.feature.CreateSegmentCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateSegmentCommand} returns this
 */
proto.bucketeer.feature.CreateSegmentCommand.prototype.setDescription = function(value) {
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
proto.bucketeer.feature.DeleteSegmentCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteSegmentCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteSegmentCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteSegmentCommand}
 */
proto.bucketeer.feature.DeleteSegmentCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteSegmentCommand;
  return proto.bucketeer.feature.DeleteSegmentCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteSegmentCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteSegmentCommand}
 */
proto.bucketeer.feature.DeleteSegmentCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteSegmentCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteSegmentCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteSegmentCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.ChangeSegmentNameCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeSegmentNameCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeSegmentNameCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeSegmentNameCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.ChangeSegmentNameCommand}
 */
proto.bucketeer.feature.ChangeSegmentNameCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeSegmentNameCommand;
  return proto.bucketeer.feature.ChangeSegmentNameCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeSegmentNameCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeSegmentNameCommand}
 */
proto.bucketeer.feature.ChangeSegmentNameCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
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
proto.bucketeer.feature.ChangeSegmentNameCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeSegmentNameCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeSegmentNameCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeSegmentNameCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeSegmentNameCommand.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeSegmentNameCommand} returns this
 */
proto.bucketeer.feature.ChangeSegmentNameCommand.prototype.setName = function(value) {
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
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeSegmentDescriptionCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeSegmentDescriptionCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    description: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.ChangeSegmentDescriptionCommand}
 */
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeSegmentDescriptionCommand;
  return proto.bucketeer.feature.ChangeSegmentDescriptionCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeSegmentDescriptionCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeSegmentDescriptionCommand}
 */
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
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
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeSegmentDescriptionCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeSegmentDescriptionCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string description = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeSegmentDescriptionCommand} returns this
 */
proto.bucketeer.feature.ChangeSegmentDescriptionCommand.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.AddSegmentUserCommand.repeatedFields_ = [1];



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
proto.bucketeer.feature.AddSegmentUserCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddSegmentUserCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddSegmentUserCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddSegmentUserCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    userIdsList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f,
    state: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.feature.AddSegmentUserCommand}
 */
proto.bucketeer.feature.AddSegmentUserCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddSegmentUserCommand;
  return proto.bucketeer.feature.AddSegmentUserCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddSegmentUserCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddSegmentUserCommand}
 */
proto.bucketeer.feature.AddSegmentUserCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addUserIds(value);
      break;
    case 2:
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
proto.bucketeer.feature.AddSegmentUserCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddSegmentUserCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddSegmentUserCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddSegmentUserCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * repeated string user_ids = 1;
 * @return {!Array<string>}
 */
proto.bucketeer.feature.AddSegmentUserCommand.prototype.getUserIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.feature.AddSegmentUserCommand} returns this
 */
proto.bucketeer.feature.AddSegmentUserCommand.prototype.setUserIdsList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.AddSegmentUserCommand} returns this
 */
proto.bucketeer.feature.AddSegmentUserCommand.prototype.addUserIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.AddSegmentUserCommand} returns this
 */
proto.bucketeer.feature.AddSegmentUserCommand.prototype.clearUserIdsList = function() {
  return this.setUserIdsList([]);
};


/**
 * optional SegmentUser.State state = 2;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.feature.AddSegmentUserCommand.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.feature.AddSegmentUserCommand} returns this
 */
proto.bucketeer.feature.AddSegmentUserCommand.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.repeatedFields_ = [1];



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
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteSegmentUserCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteSegmentUserCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    userIdsList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f,
    state: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.feature.DeleteSegmentUserCommand}
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteSegmentUserCommand;
  return proto.bucketeer.feature.DeleteSegmentUserCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteSegmentUserCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteSegmentUserCommand}
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addUserIds(value);
      break;
    case 2:
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
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteSegmentUserCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteSegmentUserCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * repeated string user_ids = 1;
 * @return {!Array<string>}
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.getUserIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.feature.DeleteSegmentUserCommand} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.setUserIdsList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.DeleteSegmentUserCommand} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.addUserIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.DeleteSegmentUserCommand} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.clearUserIdsList = function() {
  return this.setUserIdsList([]);
};


/**
 * optional SegmentUser.State state = 2;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.feature.DeleteSegmentUserCommand} returns this
 */
proto.bucketeer.feature.DeleteSegmentUserCommand.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
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
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.BulkUploadSegmentUsersCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    data: msg.getData_asB64(),
    state: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.BulkUploadSegmentUsersCommand;
  return proto.bucketeer.feature.BulkUploadSegmentUsersCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.deserializeBinaryFromReader = function(msg, reader) {
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
    case 2:
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
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.BulkUploadSegmentUsersCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional bytes data = 1;
 * @return {!(string|Uint8Array)}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.getData = function() {
  return /** @type {!(string|Uint8Array)} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * optional bytes data = 1;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.getData_asB64 = function() {
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
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand} returns this
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.setData = function(value) {
  return jspb.Message.setProto3BytesField(this, 1, value);
};


/**
 * optional SegmentUser.State state = 2;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.feature.BulkUploadSegmentUsersCommand} returns this
 */
proto.bucketeer.feature.BulkUploadSegmentUsersCommand.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
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
proto.bucketeer.feature.IncrementFeatureVersionCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.IncrementFeatureVersionCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.IncrementFeatureVersionCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.IncrementFeatureVersionCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.IncrementFeatureVersionCommand}
 */
proto.bucketeer.feature.IncrementFeatureVersionCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.IncrementFeatureVersionCommand;
  return proto.bucketeer.feature.IncrementFeatureVersionCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.IncrementFeatureVersionCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.IncrementFeatureVersionCommand}
 */
proto.bucketeer.feature.IncrementFeatureVersionCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.IncrementFeatureVersionCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.IncrementFeatureVersionCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.IncrementFeatureVersionCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.IncrementFeatureVersionCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.CloneFeatureCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CloneFeatureCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CloneFeatureCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CloneFeatureCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.CloneFeatureCommand}
 */
proto.bucketeer.feature.CloneFeatureCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CloneFeatureCommand;
  return proto.bucketeer.feature.CloneFeatureCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CloneFeatureCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CloneFeatureCommand}
 */
proto.bucketeer.feature.CloneFeatureCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.CloneFeatureCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CloneFeatureCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CloneFeatureCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CloneFeatureCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string environment_namespace = 1;
 * @return {string}
 */
proto.bucketeer.feature.CloneFeatureCommand.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CloneFeatureCommand} returns this
 */
proto.bucketeer.feature.CloneFeatureCommand.prototype.setEnvironmentNamespace = function(value) {
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
proto.bucketeer.feature.ResetSamplingSeedCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ResetSamplingSeedCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ResetSamplingSeedCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetSamplingSeedCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.ResetSamplingSeedCommand}
 */
proto.bucketeer.feature.ResetSamplingSeedCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ResetSamplingSeedCommand;
  return proto.bucketeer.feature.ResetSamplingSeedCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ResetSamplingSeedCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ResetSamplingSeedCommand}
 */
proto.bucketeer.feature.ResetSamplingSeedCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.ResetSamplingSeedCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ResetSamplingSeedCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ResetSamplingSeedCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetSamplingSeedCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.AddPrerequisiteCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.AddPrerequisiteCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.AddPrerequisiteCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddPrerequisiteCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    prerequisite: (f = msg.getPrerequisite()) && proto_feature_prerequisite_pb.Prerequisite.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.AddPrerequisiteCommand}
 */
proto.bucketeer.feature.AddPrerequisiteCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.AddPrerequisiteCommand;
  return proto.bucketeer.feature.AddPrerequisiteCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.AddPrerequisiteCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.AddPrerequisiteCommand}
 */
proto.bucketeer.feature.AddPrerequisiteCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_prerequisite_pb.Prerequisite;
      reader.readMessage(value,proto_feature_prerequisite_pb.Prerequisite.deserializeBinaryFromReader);
      msg.setPrerequisite(value);
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
proto.bucketeer.feature.AddPrerequisiteCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.AddPrerequisiteCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.AddPrerequisiteCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.AddPrerequisiteCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPrerequisite();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_prerequisite_pb.Prerequisite.serializeBinaryToWriter
    );
  }
};


/**
 * optional Prerequisite prerequisite = 1;
 * @return {?proto.bucketeer.feature.Prerequisite}
 */
proto.bucketeer.feature.AddPrerequisiteCommand.prototype.getPrerequisite = function() {
  return /** @type{?proto.bucketeer.feature.Prerequisite} */ (
    jspb.Message.getWrapperField(this, proto_feature_prerequisite_pb.Prerequisite, 1));
};


/**
 * @param {?proto.bucketeer.feature.Prerequisite|undefined} value
 * @return {!proto.bucketeer.feature.AddPrerequisiteCommand} returns this
*/
proto.bucketeer.feature.AddPrerequisiteCommand.prototype.setPrerequisite = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.AddPrerequisiteCommand} returns this
 */
proto.bucketeer.feature.AddPrerequisiteCommand.prototype.clearPrerequisite = function() {
  return this.setPrerequisite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.AddPrerequisiteCommand.prototype.hasPrerequisite = function() {
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
proto.bucketeer.feature.RemovePrerequisiteCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.RemovePrerequisiteCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.RemovePrerequisiteCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemovePrerequisiteCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.RemovePrerequisiteCommand}
 */
proto.bucketeer.feature.RemovePrerequisiteCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.RemovePrerequisiteCommand;
  return proto.bucketeer.feature.RemovePrerequisiteCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.RemovePrerequisiteCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.RemovePrerequisiteCommand}
 */
proto.bucketeer.feature.RemovePrerequisiteCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.RemovePrerequisiteCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.RemovePrerequisiteCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.RemovePrerequisiteCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.RemovePrerequisiteCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.feature.RemovePrerequisiteCommand.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.RemovePrerequisiteCommand} returns this
 */
proto.bucketeer.feature.RemovePrerequisiteCommand.prototype.setFeatureId = function(value) {
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
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangePrerequisiteVariationCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    prerequisite: (f = msg.getPrerequisite()) && proto_feature_prerequisite_pb.Prerequisite.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand}
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangePrerequisiteVariationCommand;
  return proto.bucketeer.feature.ChangePrerequisiteVariationCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand}
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto_feature_prerequisite_pb.Prerequisite;
      reader.readMessage(value,proto_feature_prerequisite_pb.Prerequisite.deserializeBinaryFromReader);
      msg.setPrerequisite(value);
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
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangePrerequisiteVariationCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPrerequisite();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_feature_prerequisite_pb.Prerequisite.serializeBinaryToWriter
    );
  }
};


/**
 * optional Prerequisite prerequisite = 1;
 * @return {?proto.bucketeer.feature.Prerequisite}
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.prototype.getPrerequisite = function() {
  return /** @type{?proto.bucketeer.feature.Prerequisite} */ (
    jspb.Message.getWrapperField(this, proto_feature_prerequisite_pb.Prerequisite, 1));
};


/**
 * @param {?proto.bucketeer.feature.Prerequisite|undefined} value
 * @return {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand} returns this
*/
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.prototype.setPrerequisite = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.feature.ChangePrerequisiteVariationCommand} returns this
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.prototype.clearPrerequisite = function() {
  return this.setPrerequisite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.feature.ChangePrerequisiteVariationCommand.prototype.hasPrerequisite = function() {
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
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.CreateFlagTriggerCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.CreateFlagTriggerCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0),
    action: jspb.Message.getFieldWithDefault(msg, 3, 0),
    description: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.feature.CreateFlagTriggerCommand}
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.CreateFlagTriggerCommand;
  return proto.bucketeer.feature.CreateFlagTriggerCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.CreateFlagTriggerCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.CreateFlagTriggerCommand}
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.bucketeer.feature.FlagTrigger.Type} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.feature.FlagTrigger.Action} */ (reader.readEnum());
      msg.setAction(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
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
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.CreateFlagTriggerCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.CreateFlagTriggerCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getAction();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerCommand} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional FlagTrigger.Type type = 2;
 * @return {!proto.bucketeer.feature.FlagTrigger.Type}
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.getType = function() {
  return /** @type {!proto.bucketeer.feature.FlagTrigger.Type} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.FlagTrigger.Type} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerCommand} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.setType = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional FlagTrigger.Action action = 3;
 * @return {!proto.bucketeer.feature.FlagTrigger.Action}
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.getAction = function() {
  return /** @type {!proto.bucketeer.feature.FlagTrigger.Action} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.FlagTrigger.Action} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerCommand} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.setAction = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional string description = 4;
 * @return {string}
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.CreateFlagTriggerCommand} returns this
 */
proto.bucketeer.feature.CreateFlagTriggerCommand.prototype.setDescription = function(value) {
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
proto.bucketeer.feature.ResetFlagTriggerCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ResetFlagTriggerCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ResetFlagTriggerCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetFlagTriggerCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.ResetFlagTriggerCommand}
 */
proto.bucketeer.feature.ResetFlagTriggerCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ResetFlagTriggerCommand;
  return proto.bucketeer.feature.ResetFlagTriggerCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ResetFlagTriggerCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ResetFlagTriggerCommand}
 */
proto.bucketeer.feature.ResetFlagTriggerCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.ResetFlagTriggerCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ResetFlagTriggerCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ResetFlagTriggerCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ResetFlagTriggerCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.toObject = function(includeInstance, msg) {
  var f, obj = {
    description: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand}
 */
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand;
  return proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand}
 */
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
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
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string description = 1;
 * @return {string}
 */
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand} returns this
 */
proto.bucketeer.feature.ChangeFlagTriggerDescriptionCommand.prototype.setDescription = function(value) {
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
proto.bucketeer.feature.EnableFlagTriggerCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.EnableFlagTriggerCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.EnableFlagTriggerCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFlagTriggerCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.EnableFlagTriggerCommand}
 */
proto.bucketeer.feature.EnableFlagTriggerCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.EnableFlagTriggerCommand;
  return proto.bucketeer.feature.EnableFlagTriggerCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.EnableFlagTriggerCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.EnableFlagTriggerCommand}
 */
proto.bucketeer.feature.EnableFlagTriggerCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.EnableFlagTriggerCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.EnableFlagTriggerCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.EnableFlagTriggerCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.EnableFlagTriggerCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DisableFlagTriggerCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DisableFlagTriggerCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DisableFlagTriggerCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFlagTriggerCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DisableFlagTriggerCommand}
 */
proto.bucketeer.feature.DisableFlagTriggerCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DisableFlagTriggerCommand;
  return proto.bucketeer.feature.DisableFlagTriggerCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DisableFlagTriggerCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DisableFlagTriggerCommand}
 */
proto.bucketeer.feature.DisableFlagTriggerCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DisableFlagTriggerCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DisableFlagTriggerCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DisableFlagTriggerCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DisableFlagTriggerCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.DeleteFlagTriggerCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.DeleteFlagTriggerCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFlagTriggerCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerCommand}
 */
proto.bucketeer.feature.DeleteFlagTriggerCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.DeleteFlagTriggerCommand;
  return proto.bucketeer.feature.DeleteFlagTriggerCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.DeleteFlagTriggerCommand}
 */
proto.bucketeer.feature.DeleteFlagTriggerCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.DeleteFlagTriggerCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.DeleteFlagTriggerCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.DeleteFlagTriggerCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.DeleteFlagTriggerCommand.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerUsageCommand} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerUsageCommand}
 */
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.UpdateFlagTriggerUsageCommand;
  return proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerUsageCommand} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.UpdateFlagTriggerUsageCommand}
 */
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.UpdateFlagTriggerUsageCommand} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.UpdateFlagTriggerUsageCommand.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};


goog.object.extend(exports, proto.bucketeer.feature);
