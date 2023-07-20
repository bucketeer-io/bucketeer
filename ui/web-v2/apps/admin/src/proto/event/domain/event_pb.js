// source: proto/event/domain/event.proto
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
var proto_feature_clause_pb = require('../../../proto/feature/clause_pb.js');
goog.object.extend(proto, proto_feature_clause_pb);
var proto_feature_feature_pb = require('../../../proto/feature/feature_pb.js');
goog.object.extend(proto, proto_feature_feature_pb);
var proto_feature_rule_pb = require('../../../proto/feature/rule_pb.js');
goog.object.extend(proto, proto_feature_rule_pb);
var proto_feature_variation_pb = require('../../../proto/feature/variation_pb.js');
goog.object.extend(proto, proto_feature_variation_pb);
var proto_feature_strategy_pb = require('../../../proto/feature/strategy_pb.js');
goog.object.extend(proto, proto_feature_strategy_pb);
var proto_feature_segment_pb = require('../../../proto/feature/segment_pb.js');
goog.object.extend(proto, proto_feature_segment_pb);
var proto_feature_target_pb = require('../../../proto/feature/target_pb.js');
goog.object.extend(proto, proto_feature_target_pb);
var proto_account_account_pb = require('../../../proto/account/account_pb.js');
goog.object.extend(proto, proto_account_account_pb);
var proto_account_api_key_pb = require('../../../proto/account/api_key_pb.js');
goog.object.extend(proto, proto_account_api_key_pb);
var proto_autoops_auto_ops_rule_pb = require('../../../proto/autoops/auto_ops_rule_pb.js');
goog.object.extend(proto, proto_autoops_auto_ops_rule_pb);
var proto_autoops_clause_pb = require('../../../proto/autoops/clause_pb.js');
goog.object.extend(proto, proto_autoops_clause_pb);
var proto_notification_subscription_pb = require('../../../proto/notification/subscription_pb.js');
goog.object.extend(proto, proto_notification_subscription_pb);
var proto_notification_recipient_pb = require('../../../proto/notification/recipient_pb.js');
goog.object.extend(proto, proto_notification_recipient_pb);
var proto_feature_prerequisite_pb = require('../../../proto/feature/prerequisite_pb.js');
goog.object.extend(proto, proto_feature_prerequisite_pb);
goog.exportSymbol('proto.bucketeer.event.domain.APIKeyCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.APIKeyDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.APIKeyEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.APIKeyNameChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AccountCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AccountDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AccountDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AccountEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AccountRoleChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminAccountCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminAccountDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminAccountDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminAccountEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ClauseAttributeChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ClauseOperatorChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ClauseValueAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ClauseValueRemovedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.DatetimeClauseAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.DatetimeClauseChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.Editor', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.EnvironmentCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.EnvironmentDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.EnvironmentRenamedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.EvaluationDelayableSetEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.EvaluationUndelayableSetEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.Event', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.Event.EntityType', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.Event.Type', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentArchivedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentFinishedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentNameChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentPeriodChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentStartAtChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentStartedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentStopAtChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ExperimentStoppedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureArchivedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureClonedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureOffVariationChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureRenamedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureRuleAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureRuleDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureTagAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureTagRemovedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureUnarchivedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureVariationAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureVariationRemovedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.FeatureVersionIncrementedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.GoalArchivedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.GoalCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.GoalDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.GoalDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.GoalRenamedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.Options', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PrerequisiteAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PrerequisiteRemovedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ProjectCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ProjectDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ProjectDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ProjectEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ProjectTrialConvertedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.ProjectTrialCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PushCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PushDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PushRenamedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PushTagsAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.PushTagsDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.RuleClauseAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.RuleClauseDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentClauseValueAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentNameChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentRuleAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentRuleDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentUserAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SegmentUserDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionDisabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionEnabledEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionRenamedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.VariationDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.VariationNameChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.VariationUserAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.VariationUserRemovedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.VariationValueChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.WebhookClauseAddedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.WebhookClauseChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.WebhookCreatedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.WebhookDeletedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.WebhookDescriptionChangedEvent', null, global);
goog.exportSymbol('proto.bucketeer.event.domain.WebhookNameChangedEvent', null, global);
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
proto.bucketeer.event.domain.Event = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.Event, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.Event.displayName = 'proto.bucketeer.event.domain.Event';
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
proto.bucketeer.event.domain.Editor = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.Editor, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.Editor.displayName = 'proto.bucketeer.event.domain.Editor';
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
proto.bucketeer.event.domain.Options = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.Options, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.Options.displayName = 'proto.bucketeer.event.domain.Options';
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
proto.bucketeer.event.domain.FeatureCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.FeatureCreatedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureCreatedEvent.displayName = 'proto.bucketeer.event.domain.FeatureCreatedEvent';
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
proto.bucketeer.event.domain.FeatureEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureEnabledEvent.displayName = 'proto.bucketeer.event.domain.FeatureEnabledEvent';
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
proto.bucketeer.event.domain.FeatureDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureDisabledEvent.displayName = 'proto.bucketeer.event.domain.FeatureDisabledEvent';
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
proto.bucketeer.event.domain.FeatureArchivedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureArchivedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureArchivedEvent.displayName = 'proto.bucketeer.event.domain.FeatureArchivedEvent';
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
proto.bucketeer.event.domain.FeatureUnarchivedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureUnarchivedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureUnarchivedEvent.displayName = 'proto.bucketeer.event.domain.FeatureUnarchivedEvent';
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
proto.bucketeer.event.domain.FeatureDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureDeletedEvent.displayName = 'proto.bucketeer.event.domain.FeatureDeletedEvent';
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
proto.bucketeer.event.domain.EvaluationDelayableSetEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.EvaluationDelayableSetEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.EvaluationDelayableSetEvent.displayName = 'proto.bucketeer.event.domain.EvaluationDelayableSetEvent';
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
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.EvaluationUndelayableSetEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.displayName = 'proto.bucketeer.event.domain.EvaluationUndelayableSetEvent';
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
proto.bucketeer.event.domain.FeatureRenamedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureRenamedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureRenamedEvent.displayName = 'proto.bucketeer.event.domain.FeatureRenamedEvent';
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
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.FeatureDescriptionChangedEvent';
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
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureOffVariationChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.displayName = 'proto.bucketeer.event.domain.FeatureOffVariationChangedEvent';
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
proto.bucketeer.event.domain.FeatureVariationAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureVariationAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureVariationAddedEvent.displayName = 'proto.bucketeer.event.domain.FeatureVariationAddedEvent';
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
proto.bucketeer.event.domain.FeatureVariationRemovedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureVariationRemovedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureVariationRemovedEvent.displayName = 'proto.bucketeer.event.domain.FeatureVariationRemovedEvent';
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
proto.bucketeer.event.domain.VariationValueChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.VariationValueChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.VariationValueChangedEvent.displayName = 'proto.bucketeer.event.domain.VariationValueChangedEvent';
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
proto.bucketeer.event.domain.VariationNameChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.VariationNameChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.VariationNameChangedEvent.displayName = 'proto.bucketeer.event.domain.VariationNameChangedEvent';
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
proto.bucketeer.event.domain.VariationDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.VariationDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.VariationDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.VariationDescriptionChangedEvent';
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
proto.bucketeer.event.domain.VariationUserAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.VariationUserAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.VariationUserAddedEvent.displayName = 'proto.bucketeer.event.domain.VariationUserAddedEvent';
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
proto.bucketeer.event.domain.VariationUserRemovedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.VariationUserRemovedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.VariationUserRemovedEvent.displayName = 'proto.bucketeer.event.domain.VariationUserRemovedEvent';
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
proto.bucketeer.event.domain.FeatureRuleAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureRuleAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureRuleAddedEvent.displayName = 'proto.bucketeer.event.domain.FeatureRuleAddedEvent';
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
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.displayName = 'proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent';
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
proto.bucketeer.event.domain.FeatureRuleDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureRuleDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureRuleDeletedEvent.displayName = 'proto.bucketeer.event.domain.FeatureRuleDeletedEvent';
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
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.displayName = 'proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent';
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
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.displayName = 'proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent';
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
proto.bucketeer.event.domain.RuleClauseAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.RuleClauseAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.RuleClauseAddedEvent.displayName = 'proto.bucketeer.event.domain.RuleClauseAddedEvent';
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
proto.bucketeer.event.domain.RuleClauseDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.RuleClauseDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.RuleClauseDeletedEvent.displayName = 'proto.bucketeer.event.domain.RuleClauseDeletedEvent';
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
proto.bucketeer.event.domain.ClauseAttributeChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ClauseAttributeChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ClauseAttributeChangedEvent.displayName = 'proto.bucketeer.event.domain.ClauseAttributeChangedEvent';
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
proto.bucketeer.event.domain.ClauseOperatorChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ClauseOperatorChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ClauseOperatorChangedEvent.displayName = 'proto.bucketeer.event.domain.ClauseOperatorChangedEvent';
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
proto.bucketeer.event.domain.ClauseValueAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ClauseValueAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ClauseValueAddedEvent.displayName = 'proto.bucketeer.event.domain.ClauseValueAddedEvent';
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
proto.bucketeer.event.domain.ClauseValueRemovedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ClauseValueRemovedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ClauseValueRemovedEvent.displayName = 'proto.bucketeer.event.domain.ClauseValueRemovedEvent';
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
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.displayName = 'proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent';
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
proto.bucketeer.event.domain.FeatureTagAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureTagAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureTagAddedEvent.displayName = 'proto.bucketeer.event.domain.FeatureTagAddedEvent';
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
proto.bucketeer.event.domain.FeatureTagRemovedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureTagRemovedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureTagRemovedEvent.displayName = 'proto.bucketeer.event.domain.FeatureTagRemovedEvent';
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
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureVersionIncrementedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.displayName = 'proto.bucketeer.event.domain.FeatureVersionIncrementedEvent';
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
proto.bucketeer.event.domain.FeatureClonedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.FeatureClonedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureClonedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureClonedEvent.displayName = 'proto.bucketeer.event.domain.FeatureClonedEvent';
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
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.displayName = 'proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent';
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
proto.bucketeer.event.domain.GoalCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.GoalCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.GoalCreatedEvent.displayName = 'proto.bucketeer.event.domain.GoalCreatedEvent';
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
proto.bucketeer.event.domain.GoalRenamedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.GoalRenamedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.GoalRenamedEvent.displayName = 'proto.bucketeer.event.domain.GoalRenamedEvent';
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
proto.bucketeer.event.domain.GoalDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.GoalDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.GoalDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.GoalDescriptionChangedEvent';
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
proto.bucketeer.event.domain.GoalArchivedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.GoalArchivedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.GoalArchivedEvent.displayName = 'proto.bucketeer.event.domain.GoalArchivedEvent';
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
proto.bucketeer.event.domain.GoalDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.GoalDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.GoalDeletedEvent.displayName = 'proto.bucketeer.event.domain.GoalDeletedEvent';
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
proto.bucketeer.event.domain.ExperimentCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.ExperimentCreatedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentCreatedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentCreatedEvent';
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
proto.bucketeer.event.domain.ExperimentStoppedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentStoppedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentStoppedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentStoppedEvent';
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
proto.bucketeer.event.domain.ExperimentArchivedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentArchivedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentArchivedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentArchivedEvent';
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
proto.bucketeer.event.domain.ExperimentDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentDeletedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentDeletedEvent';
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
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentStartAtChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentStartAtChangedEvent';
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
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentStopAtChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentStopAtChangedEvent';
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
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentPeriodChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentPeriodChangedEvent';
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
proto.bucketeer.event.domain.ExperimentNameChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentNameChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentNameChangedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentNameChangedEvent';
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
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent';
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
proto.bucketeer.event.domain.ExperimentStartedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentStartedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentStartedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentStartedEvent';
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
proto.bucketeer.event.domain.ExperimentFinishedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ExperimentFinishedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ExperimentFinishedEvent.displayName = 'proto.bucketeer.event.domain.ExperimentFinishedEvent';
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
proto.bucketeer.event.domain.AccountCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AccountCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AccountCreatedEvent.displayName = 'proto.bucketeer.event.domain.AccountCreatedEvent';
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
proto.bucketeer.event.domain.AccountRoleChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AccountRoleChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AccountRoleChangedEvent.displayName = 'proto.bucketeer.event.domain.AccountRoleChangedEvent';
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
proto.bucketeer.event.domain.AccountEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AccountEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AccountEnabledEvent.displayName = 'proto.bucketeer.event.domain.AccountEnabledEvent';
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
proto.bucketeer.event.domain.AccountDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AccountDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AccountDisabledEvent.displayName = 'proto.bucketeer.event.domain.AccountDisabledEvent';
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
proto.bucketeer.event.domain.AccountDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AccountDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AccountDeletedEvent.displayName = 'proto.bucketeer.event.domain.AccountDeletedEvent';
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
proto.bucketeer.event.domain.APIKeyCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.APIKeyCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.APIKeyCreatedEvent.displayName = 'proto.bucketeer.event.domain.APIKeyCreatedEvent';
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
proto.bucketeer.event.domain.APIKeyNameChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.APIKeyNameChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.APIKeyNameChangedEvent.displayName = 'proto.bucketeer.event.domain.APIKeyNameChangedEvent';
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
proto.bucketeer.event.domain.APIKeyEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.APIKeyEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.APIKeyEnabledEvent.displayName = 'proto.bucketeer.event.domain.APIKeyEnabledEvent';
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
proto.bucketeer.event.domain.APIKeyDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.APIKeyDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.APIKeyDisabledEvent.displayName = 'proto.bucketeer.event.domain.APIKeyDisabledEvent';
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
proto.bucketeer.event.domain.SegmentCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentCreatedEvent.displayName = 'proto.bucketeer.event.domain.SegmentCreatedEvent';
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
proto.bucketeer.event.domain.SegmentDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentDeletedEvent.displayName = 'proto.bucketeer.event.domain.SegmentDeletedEvent';
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
proto.bucketeer.event.domain.SegmentNameChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentNameChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentNameChangedEvent.displayName = 'proto.bucketeer.event.domain.SegmentNameChangedEvent';
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
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.SegmentDescriptionChangedEvent';
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
proto.bucketeer.event.domain.SegmentRuleAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentRuleAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentRuleAddedEvent.displayName = 'proto.bucketeer.event.domain.SegmentRuleAddedEvent';
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
proto.bucketeer.event.domain.SegmentRuleDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentRuleDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentRuleDeletedEvent.displayName = 'proto.bucketeer.event.domain.SegmentRuleDeletedEvent';
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
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.displayName = 'proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent';
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
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.displayName = 'proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent';
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
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.displayName = 'proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent';
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
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.displayName = 'proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent';
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
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentClauseValueAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.displayName = 'proto.bucketeer.event.domain.SegmentClauseValueAddedEvent';
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
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.displayName = 'proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent';
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
proto.bucketeer.event.domain.SegmentUserAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.SegmentUserAddedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentUserAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentUserAddedEvent.displayName = 'proto.bucketeer.event.domain.SegmentUserAddedEvent';
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
proto.bucketeer.event.domain.SegmentUserDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.SegmentUserDeletedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentUserDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentUserDeletedEvent.displayName = 'proto.bucketeer.event.domain.SegmentUserDeletedEvent';
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
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.displayName = 'proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent';
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
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.displayName = 'proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent';
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
proto.bucketeer.event.domain.EnvironmentCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.EnvironmentCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.EnvironmentCreatedEvent.displayName = 'proto.bucketeer.event.domain.EnvironmentCreatedEvent';
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
proto.bucketeer.event.domain.EnvironmentRenamedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.EnvironmentRenamedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.EnvironmentRenamedEvent.displayName = 'proto.bucketeer.event.domain.EnvironmentRenamedEvent';
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
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent';
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
proto.bucketeer.event.domain.EnvironmentDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.EnvironmentDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.EnvironmentDeletedEvent.displayName = 'proto.bucketeer.event.domain.EnvironmentDeletedEvent';
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
proto.bucketeer.event.domain.AdminAccountCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminAccountCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminAccountCreatedEvent.displayName = 'proto.bucketeer.event.domain.AdminAccountCreatedEvent';
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
proto.bucketeer.event.domain.AdminAccountEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminAccountEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminAccountEnabledEvent.displayName = 'proto.bucketeer.event.domain.AdminAccountEnabledEvent';
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
proto.bucketeer.event.domain.AdminAccountDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminAccountDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminAccountDisabledEvent.displayName = 'proto.bucketeer.event.domain.AdminAccountDisabledEvent';
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
proto.bucketeer.event.domain.AdminAccountDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminAccountDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminAccountDeletedEvent.displayName = 'proto.bucketeer.event.domain.AdminAccountDeletedEvent';
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
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.displayName = 'proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent';
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
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.displayName = 'proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent';
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
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.displayName = 'proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent';
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
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.displayName = 'proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent';
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
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.displayName = 'proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent';
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
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.displayName = 'proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent';
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
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.displayName = 'proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent';
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
proto.bucketeer.event.domain.DatetimeClauseAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.DatetimeClauseAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.DatetimeClauseAddedEvent.displayName = 'proto.bucketeer.event.domain.DatetimeClauseAddedEvent';
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
proto.bucketeer.event.domain.DatetimeClauseChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.DatetimeClauseChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.DatetimeClauseChangedEvent.displayName = 'proto.bucketeer.event.domain.DatetimeClauseChangedEvent';
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
proto.bucketeer.event.domain.PushCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.PushCreatedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.PushCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PushCreatedEvent.displayName = 'proto.bucketeer.event.domain.PushCreatedEvent';
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
proto.bucketeer.event.domain.PushDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.PushDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PushDeletedEvent.displayName = 'proto.bucketeer.event.domain.PushDeletedEvent';
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
proto.bucketeer.event.domain.PushTagsAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.PushTagsAddedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.PushTagsAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PushTagsAddedEvent.displayName = 'proto.bucketeer.event.domain.PushTagsAddedEvent';
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
proto.bucketeer.event.domain.PushTagsDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.PushTagsDeletedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.PushTagsDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PushTagsDeletedEvent.displayName = 'proto.bucketeer.event.domain.PushTagsDeletedEvent';
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
proto.bucketeer.event.domain.PushRenamedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.PushRenamedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PushRenamedEvent.displayName = 'proto.bucketeer.event.domain.PushRenamedEvent';
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
proto.bucketeer.event.domain.SubscriptionCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.SubscriptionCreatedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionCreatedEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionCreatedEvent';
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
proto.bucketeer.event.domain.SubscriptionDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionDeletedEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionDeletedEvent';
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
proto.bucketeer.event.domain.SubscriptionEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionEnabledEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionEnabledEvent';
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
proto.bucketeer.event.domain.SubscriptionDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionDisabledEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionDisabledEvent';
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
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent';
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
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent';
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
proto.bucketeer.event.domain.SubscriptionRenamedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.SubscriptionRenamedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.SubscriptionRenamedEvent.displayName = 'proto.bucketeer.event.domain.SubscriptionRenamedEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent';
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
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.displayName = 'proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent';
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
proto.bucketeer.event.domain.ProjectCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ProjectCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ProjectCreatedEvent.displayName = 'proto.bucketeer.event.domain.ProjectCreatedEvent';
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
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ProjectDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.ProjectDescriptionChangedEvent';
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
proto.bucketeer.event.domain.ProjectEnabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ProjectEnabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ProjectEnabledEvent.displayName = 'proto.bucketeer.event.domain.ProjectEnabledEvent';
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
proto.bucketeer.event.domain.ProjectDisabledEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ProjectDisabledEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ProjectDisabledEvent.displayName = 'proto.bucketeer.event.domain.ProjectDisabledEvent';
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
proto.bucketeer.event.domain.ProjectTrialCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ProjectTrialCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ProjectTrialCreatedEvent.displayName = 'proto.bucketeer.event.domain.ProjectTrialCreatedEvent';
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
proto.bucketeer.event.domain.ProjectTrialConvertedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.ProjectTrialConvertedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.ProjectTrialConvertedEvent.displayName = 'proto.bucketeer.event.domain.ProjectTrialConvertedEvent';
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
proto.bucketeer.event.domain.PrerequisiteAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.PrerequisiteAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PrerequisiteAddedEvent.displayName = 'proto.bucketeer.event.domain.PrerequisiteAddedEvent';
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
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.displayName = 'proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent';
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
proto.bucketeer.event.domain.PrerequisiteRemovedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.PrerequisiteRemovedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.PrerequisiteRemovedEvent.displayName = 'proto.bucketeer.event.domain.PrerequisiteRemovedEvent';
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
proto.bucketeer.event.domain.WebhookCreatedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.WebhookCreatedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.WebhookCreatedEvent.displayName = 'proto.bucketeer.event.domain.WebhookCreatedEvent';
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
proto.bucketeer.event.domain.WebhookDeletedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.WebhookDeletedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.WebhookDeletedEvent.displayName = 'proto.bucketeer.event.domain.WebhookDeletedEvent';
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
proto.bucketeer.event.domain.WebhookNameChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.WebhookNameChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.WebhookNameChangedEvent.displayName = 'proto.bucketeer.event.domain.WebhookNameChangedEvent';
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
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.WebhookDescriptionChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.displayName = 'proto.bucketeer.event.domain.WebhookDescriptionChangedEvent';
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
proto.bucketeer.event.domain.WebhookClauseAddedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.WebhookClauseAddedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.WebhookClauseAddedEvent.displayName = 'proto.bucketeer.event.domain.WebhookClauseAddedEvent';
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
proto.bucketeer.event.domain.WebhookClauseChangedEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.event.domain.WebhookClauseChangedEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.event.domain.WebhookClauseChangedEvent.displayName = 'proto.bucketeer.event.domain.WebhookClauseChangedEvent';
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
proto.bucketeer.event.domain.Event.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.Event.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.Event} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.Event.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    timestamp: jspb.Message.getFieldWithDefault(msg, 2, 0),
    entityType: jspb.Message.getFieldWithDefault(msg, 3, 0),
    entityId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    type: jspb.Message.getFieldWithDefault(msg, 5, 0),
    editor: (f = msg.getEditor()) && proto.bucketeer.event.domain.Editor.toObject(includeInstance, f),
    data: (f = msg.getData()) && google_protobuf_any_pb.Any.toObject(includeInstance, f),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 8, ""),
    isAdminEvent: jspb.Message.getBooleanFieldWithDefault(msg, 9, false),
    options: (f = msg.getOptions()) && proto.bucketeer.event.domain.Options.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.Event}
 */
proto.bucketeer.event.domain.Event.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.Event;
  return proto.bucketeer.event.domain.Event.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.Event} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.Event}
 */
proto.bucketeer.event.domain.Event.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTimestamp(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.event.domain.Event.EntityType} */ (reader.readEnum());
      msg.setEntityType(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setEntityId(value);
      break;
    case 5:
      var value = /** @type {!proto.bucketeer.event.domain.Event.Type} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 6:
      var value = new proto.bucketeer.event.domain.Editor;
      reader.readMessage(value,proto.bucketeer.event.domain.Editor.deserializeBinaryFromReader);
      msg.setEditor(value);
      break;
    case 7:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setData(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 9:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsAdminEvent(value);
      break;
    case 10:
      var value = new proto.bucketeer.event.domain.Options;
      reader.readMessage(value,proto.bucketeer.event.domain.Options.deserializeBinaryFromReader);
      msg.setOptions(value);
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
proto.bucketeer.event.domain.Event.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.Event.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.Event} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.Event.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTimestamp();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getEntityType();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getEntityId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      5,
      f
    );
  }
  f = message.getEditor();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.bucketeer.event.domain.Editor.serializeBinaryToWriter
    );
  }
  f = message.getData();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getIsAdminEvent();
  if (f) {
    writer.writeBool(
      9,
      f
    );
  }
  f = message.getOptions();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      proto.bucketeer.event.domain.Options.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.event.domain.Event.EntityType = {
  FEATURE: 0,
  GOAL: 1,
  EXPERIMENT: 2,
  ACCOUNT: 3,
  APIKEY: 4,
  SEGMENT: 5,
  ENVIRONMENT: 6,
  ADMIN_ACCOUNT: 7,
  AUTOOPS_RULE: 8,
  PUSH: 9,
  SUBSCRIPTION: 10,
  ADMIN_SUBSCRIPTION: 11,
  PROJECT: 12,
  WEBHOOK: 13
};

/**
 * @enum {number}
 */
proto.bucketeer.event.domain.Event.Type = {
  UNKNOWN: 0,
  FEATURE_CREATED: 1,
  FEATURE_RENAMED: 2,
  FEATURE_ENABLED: 3,
  FEATURE_DISABLED: 4,
  FEATURE_DELETED: 5,
  FEATURE_DESCRIPTION_CHANGED: 8,
  FEATURE_VARIATION_ADDED: 9,
  FEATURE_VARIATION_REMOVED: 10,
  FEATURE_OFF_VARIATION_CHANGED: 11,
  VARIATION_VALUE_CHANGED: 12,
  VARIATION_NAME_CHANGED: 13,
  VARIATION_DESCRIPTION_CHANGED: 14,
  VARIATION_USER_ADDED: 15,
  VARIATION_USER_REMOVED: 16,
  FEATURE_RULE_ADDED: 17,
  FEATURE_RULE_STRATEGY_CHANGED: 18,
  FEATURE_RULE_DELETED: 19,
  RULE_CLAUSE_ADDED: 20,
  RULE_CLAUSE_DELETED: 21,
  RULE_FIXED_STRATEGY_CHANGED: 22,
  RULE_ROLLOUT_STRATEGY_CHANGED: 23,
  CLAUSE_ATTRIBUTE_CHANGED: 24,
  CLAUSE_OPERATOR_CHANGED: 25,
  CLAUSE_VALUE_ADDED: 26,
  CLAUSE_VALUE_REMOVED: 27,
  FEATURE_DEFAULT_STRATEGY_CHANGED: 28,
  FEATURE_TAG_ADDED: 29,
  FEATURE_TAG_REMOVED: 30,
  FEATURE_VERSION_INCREMENTED: 31,
  FEATURE_ARCHIVED: 32,
  FEATURE_CLONED: 33,
  FEATURE_UNARCHIVED: 35,
  SAMPLING_SEED_RESET: 34,
  PREREQUISITE_ADDED: 36,
  PREREQUISITE_REMOVED: 37,
  PREREQUISITE_VARIATION_CHANGED: 38,
  GOAL_CREATED: 100,
  GOAL_RENAMED: 101,
  GOAL_DESCRIPTION_CHANGED: 102,
  GOAL_DELETED: 103,
  GOAL_ARCHIVED: 104,
  EXPERIMENT_CREATED: 200,
  EXPERIMENT_STOPPED: 201,
  EXPERIMENT_START_AT_CHANGED: 202,
  EXPERIMENT_STOP_AT_CHANGED: 203,
  EXPERIMENT_DELETED: 204,
  EXPERIMENT_PERIOD_CHANGED: 205,
  EXPERIMENT_NAME_CHANGED: 206,
  EXPERIMENT_DESCRIPTION_CHANGED: 207,
  EXPERIMENT_STARTED: 208,
  EXPERIMENT_FINISHED: 209,
  EXPERIMENT_ARCHIVED: 210,
  ACCOUNT_CREATED: 300,
  ACCOUNT_ROLE_CHANGED: 301,
  ACCOUNT_ENABLED: 302,
  ACCOUNT_DISABLED: 303,
  ACCOUNT_DELETED: 304,
  APIKEY_CREATED: 400,
  APIKEY_NAME_CHANGED: 401,
  APIKEY_ENABLED: 402,
  APIKEY_DISABLED: 403,
  SEGMENT_CREATED: 500,
  SEGMENT_DELETED: 501,
  SEGMENT_NAME_CHANGED: 502,
  SEGMENT_DESCRIPTION_CHANGED: 503,
  SEGMENT_RULE_ADDED: 504,
  SEGMENT_RULE_DELETED: 505,
  SEGMENT_RULE_CLAUSE_ADDED: 506,
  SEGMENT_RULE_CLAUSE_DELETED: 507,
  SEGMENT_CLAUSE_ATTRIBUTE_CHANGED: 508,
  SEGMENT_CLAUSE_OPERATOR_CHANGED: 509,
  SEGMENT_CLAUSE_VALUE_ADDED: 510,
  SEGMENT_CLAUSE_VALUE_REMOVED: 511,
  SEGMENT_USER_ADDED: 512,
  SEGMENT_USER_DELETED: 513,
  SEGMENT_BULK_UPLOAD_USERS: 514,
  SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED: 515,
  ENVIRONMENT_CREATED: 600,
  ENVIRONMENT_RENAMED: 601,
  ENVIRONMENT_DESCRIPTION_CHANGED: 602,
  ENVIRONMENT_DELETED: 603,
  ADMIN_ACCOUNT_CREATED: 700,
  ADMIN_ACCOUNT_ENABLED: 702,
  ADMIN_ACCOUNT_DISABLED: 703,
  AUTOOPS_RULE_CREATED: 800,
  AUTOOPS_RULE_DELETED: 801,
  AUTOOPS_RULE_OPS_TYPE_CHANGED: 802,
  AUTOOPS_RULE_CLAUSE_DELETED: 803,
  AUTOOPS_RULE_TRIGGERED_AT_CHANGED: 804,
  OPS_EVENT_RATE_CLAUSE_ADDED: 805,
  OPS_EVENT_RATE_CLAUSE_CHANGED: 806,
  DATETIME_CLAUSE_ADDED: 807,
  DATETIME_CLAUSE_CHANGED: 808,
  PUSH_CREATED: 900,
  PUSH_DELETED: 901,
  PUSH_TAGS_ADDED: 902,
  PUSH_TAGS_DELETED: 903,
  PUSH_RENAMED: 904,
  SUBSCRIPTION_CREATED: 1000,
  SUBSCRIPTION_DELETED: 1001,
  SUBSCRIPTION_ENABLED: 1002,
  SUBSCRIPTION_DISABLED: 1003,
  SUBSCRIPTION_SOURCE_TYPE_ADDED: 1004,
  SUBSCRIPTION_SOURCE_TYPE_DELETED: 1005,
  SUBSCRIPTION_RENAMED: 1006,
  ADMIN_SUBSCRIPTION_CREATED: 1100,
  ADMIN_SUBSCRIPTION_DELETED: 1101,
  ADMIN_SUBSCRIPTION_ENABLED: 1102,
  ADMIN_SUBSCRIPTION_DISABLED: 1103,
  ADMIN_SUBSCRIPTION_SOURCE_TYPE_ADDED: 1104,
  ADMIN_SUBSCRIPTION_SOURCE_TYPE_DELETED: 1105,
  ADMIN_SUBSCRIPTION_RENAMED: 1106,
  PROJECT_CREATED: 1200,
  PROJECT_DESCRIPTION_CHANGED: 1201,
  PROJECT_ENABLED: 1202,
  PROJECT_DISABLED: 1203,
  PROJECT_TRIAL_CREATED: 1204,
  PROJECT_TRIAL_CONVERTED: 1205,
  WEBHOOK_CREATED: 1300,
  WEBHOOK_DELETED: 1301,
  WEBHOOK_NAME_CHANGED: 1302,
  WEBHOOK_DESCRIPTION_CHANGED: 1303,
  WEBHOOK_CLAUSE_ADDED: 1304,
  WEBHOOK_CLAUSE_CHANGED: 1305
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.Event.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 timestamp = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.Event.prototype.getTimestamp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setTimestamp = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional EntityType entity_type = 3;
 * @return {!proto.bucketeer.event.domain.Event.EntityType}
 */
proto.bucketeer.event.domain.Event.prototype.getEntityType = function() {
  return /** @type {!proto.bucketeer.event.domain.Event.EntityType} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.event.domain.Event.EntityType} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setEntityType = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional string entity_id = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.Event.prototype.getEntityId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setEntityId = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional Type type = 5;
 * @return {!proto.bucketeer.event.domain.Event.Type}
 */
proto.bucketeer.event.domain.Event.prototype.getType = function() {
  return /** @type {!proto.bucketeer.event.domain.Event.Type} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.bucketeer.event.domain.Event.Type} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setType = function(value) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};


/**
 * optional Editor editor = 6;
 * @return {?proto.bucketeer.event.domain.Editor}
 */
proto.bucketeer.event.domain.Event.prototype.getEditor = function() {
  return /** @type{?proto.bucketeer.event.domain.Editor} */ (
    jspb.Message.getWrapperField(this, proto.bucketeer.event.domain.Editor, 6));
};


/**
 * @param {?proto.bucketeer.event.domain.Editor|undefined} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
*/
proto.bucketeer.event.domain.Event.prototype.setEditor = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.clearEditor = function() {
  return this.setEditor(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.Event.prototype.hasEditor = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional google.protobuf.Any data = 7;
 * @return {?proto.google.protobuf.Any}
 */
proto.bucketeer.event.domain.Event.prototype.getData = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 7));
};


/**
 * @param {?proto.google.protobuf.Any|undefined} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
*/
proto.bucketeer.event.domain.Event.prototype.setData = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.clearData = function() {
  return this.setData(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.Event.prototype.hasData = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string environment_namespace = 8;
 * @return {string}
 */
proto.bucketeer.event.domain.Event.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional bool is_admin_event = 9;
 * @return {boolean}
 */
proto.bucketeer.event.domain.Event.prototype.getIsAdminEvent = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 9, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.setIsAdminEvent = function(value) {
  return jspb.Message.setProto3BooleanField(this, 9, value);
};


/**
 * optional Options options = 10;
 * @return {?proto.bucketeer.event.domain.Options}
 */
proto.bucketeer.event.domain.Event.prototype.getOptions = function() {
  return /** @type{?proto.bucketeer.event.domain.Options} */ (
    jspb.Message.getWrapperField(this, proto.bucketeer.event.domain.Options, 10));
};


/**
 * @param {?proto.bucketeer.event.domain.Options|undefined} value
 * @return {!proto.bucketeer.event.domain.Event} returns this
*/
proto.bucketeer.event.domain.Event.prototype.setOptions = function(value) {
  return jspb.Message.setWrapperField(this, 10, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.Event} returns this
 */
proto.bucketeer.event.domain.Event.prototype.clearOptions = function() {
  return this.setOptions(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.Event.prototype.hasOptions = function() {
  return jspb.Message.getField(this, 10) != null;
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
proto.bucketeer.event.domain.Editor.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.Editor.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.Editor} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.Editor.toObject = function(includeInstance, msg) {
  var f, obj = {
    email: jspb.Message.getFieldWithDefault(msg, 1, ""),
    role: jspb.Message.getFieldWithDefault(msg, 2, 0),
    isAdmin: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
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
 * @return {!proto.bucketeer.event.domain.Editor}
 */
proto.bucketeer.event.domain.Editor.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.Editor;
  return proto.bucketeer.event.domain.Editor.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.Editor} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.Editor}
 */
proto.bucketeer.event.domain.Editor.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    case 2:
      var value = /** @type {!proto.bucketeer.account.Account.Role} */ (reader.readEnum());
      msg.setRole(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsAdmin(value);
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
proto.bucketeer.event.domain.Editor.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.Editor.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.Editor} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.Editor.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getIsAdmin();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.Editor.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.Editor} returns this
 */
proto.bucketeer.event.domain.Editor.prototype.setEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.account.Account.Role role = 2;
 * @return {!proto.bucketeer.account.Account.Role}
 */
proto.bucketeer.event.domain.Editor.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.Account.Role} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.account.Account.Role} value
 * @return {!proto.bucketeer.event.domain.Editor} returns this
 */
proto.bucketeer.event.domain.Editor.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional bool is_admin = 3;
 * @return {boolean}
 */
proto.bucketeer.event.domain.Editor.prototype.getIsAdmin = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.Editor} returns this
 */
proto.bucketeer.event.domain.Editor.prototype.setIsAdmin = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
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
proto.bucketeer.event.domain.Options.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.Options.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.Options} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.Options.toObject = function(includeInstance, msg) {
  var f, obj = {
    comment: jspb.Message.getFieldWithDefault(msg, 1, ""),
    newVersion: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.event.domain.Options}
 */
proto.bucketeer.event.domain.Options.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.Options;
  return proto.bucketeer.event.domain.Options.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.Options} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.Options}
 */
proto.bucketeer.event.domain.Options.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setComment(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setNewVersion(value);
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
proto.bucketeer.event.domain.Options.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.Options.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.Options} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.Options.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getComment();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getNewVersion();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
};


/**
 * optional string comment = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.Options.prototype.getComment = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.Options} returns this
 */
proto.bucketeer.event.domain.Options.prototype.setComment = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int32 new_version = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.Options.prototype.getNewVersion = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.Options} returns this
 */
proto.bucketeer.event.domain.Options.prototype.setNewVersion = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.repeatedFields_ = [5,9,10,11,12];



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
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, ""),
    user: jspb.Message.getFieldWithDefault(msg, 4, ""),
    variationsList: jspb.Message.toObjectList(msg.getVariationsList(),
    proto_feature_variation_pb.Variation.toObject, includeInstance),
    defaultOnVariationIndex: (f = msg.getDefaultOnVariationIndex()) && google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
    defaultOffVariationIndex: (f = msg.getDefaultOffVariationIndex()) && google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
    variationType: jspb.Message.getFieldWithDefault(msg, 8, 0),
    tagsList: (f = jspb.Message.getRepeatedField(msg, 9)) == null ? undefined : f,
    prerequisitesList: jspb.Message.toObjectList(msg.getPrerequisitesList(),
    proto_feature_prerequisite_pb.Prerequisite.toObject, includeInstance),
    rulesList: jspb.Message.toObjectList(msg.getRulesList(),
    proto_feature_rule_pb.Rule.toObject, includeInstance),
    targetsList: jspb.Message.toObjectList(msg.getTargetsList(),
    proto_feature_target_pb.Target.toObject, includeInstance)
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
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureCreatedEvent;
  return proto.bucketeer.event.domain.FeatureCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setUser(value);
      break;
    case 5:
      var value = new proto_feature_variation_pb.Variation;
      reader.readMessage(value,proto_feature_variation_pb.Variation.deserializeBinaryFromReader);
      msg.addVariations(value);
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
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addTags(value);
      break;
    case 10:
      var value = new proto_feature_prerequisite_pb.Prerequisite;
      reader.readMessage(value,proto_feature_prerequisite_pb.Prerequisite.deserializeBinaryFromReader);
      msg.addPrerequisites(value);
      break;
    case 11:
      var value = new proto_feature_rule_pb.Rule;
      reader.readMessage(value,proto_feature_rule_pb.Rule.deserializeBinaryFromReader);
      msg.addRules(value);
      break;
    case 12:
      var value = new proto_feature_target_pb.Target;
      reader.readMessage(value,proto_feature_target_pb.Target.deserializeBinaryFromReader);
      msg.addTargets(value);
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
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getUser();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getVariationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto_feature_variation_pb.Variation.serializeBinaryToWriter
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
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
  f = message.getPrerequisitesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      10,
      f,
      proto_feature_prerequisite_pb.Prerequisite.serializeBinaryToWriter
    );
  }
  f = message.getRulesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      11,
      f,
      proto_feature_rule_pb.Rule.serializeBinaryToWriter
    );
  }
  f = message.getTargetsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      12,
      f,
      proto_feature_target_pb.Target.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string user = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setUser = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * repeated bucketeer.feature.Variation variations = 5;
 * @return {!Array<!proto.bucketeer.feature.Variation>}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getVariationsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Variation>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_variation_pb.Variation, 5));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Variation>} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setVariationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.bucketeer.feature.Variation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Variation}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.addVariations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.bucketeer.feature.Variation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearVariationsList = function() {
  return this.setVariationsList([]);
};


/**
 * optional google.protobuf.Int32Value default_on_variation_index = 6;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getDefaultOnVariationIndex = function() {
  return /** @type{?proto.google.protobuf.Int32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.Int32Value, 6));
};


/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setDefaultOnVariationIndex = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearDefaultOnVariationIndex = function() {
  return this.setDefaultOnVariationIndex(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.hasDefaultOnVariationIndex = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional google.protobuf.Int32Value default_off_variation_index = 7;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getDefaultOffVariationIndex = function() {
  return /** @type{?proto.google.protobuf.Int32Value} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.Int32Value, 7));
};


/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setDefaultOffVariationIndex = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearDefaultOffVariationIndex = function() {
  return this.setDefaultOffVariationIndex(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.hasDefaultOffVariationIndex = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional bucketeer.feature.Feature.VariationType variation_type = 8;
 * @return {!proto.bucketeer.feature.Feature.VariationType}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getVariationType = function() {
  return /** @type {!proto.bucketeer.feature.Feature.VariationType} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/**
 * @param {!proto.bucketeer.feature.Feature.VariationType} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setVariationType = function(value) {
  return jspb.Message.setProto3EnumField(this, 8, value);
};


/**
 * repeated string tags = 9;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * repeated bucketeer.feature.Prerequisite prerequisites = 10;
 * @return {!Array<!proto.bucketeer.feature.Prerequisite>}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getPrerequisitesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Prerequisite>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_prerequisite_pb.Prerequisite, 10));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Prerequisite>} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setPrerequisitesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 10, value);
};


/**
 * @param {!proto.bucketeer.feature.Prerequisite=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Prerequisite}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.addPrerequisites = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 10, opt_value, proto.bucketeer.feature.Prerequisite, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearPrerequisitesList = function() {
  return this.setPrerequisitesList([]);
};


/**
 * repeated bucketeer.feature.Rule rules = 11;
 * @return {!Array<!proto.bucketeer.feature.Rule>}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getRulesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Rule>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_rule_pb.Rule, 11));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Rule>} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setRulesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 11, value);
};


/**
 * @param {!proto.bucketeer.feature.Rule=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Rule}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.addRules = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 11, opt_value, proto.bucketeer.feature.Rule, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearRulesList = function() {
  return this.setRulesList([]);
};


/**
 * repeated bucketeer.feature.Target targets = 12;
 * @return {!Array<!proto.bucketeer.feature.Target>}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.getTargetsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Target>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_target_pb.Target, 12));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Target>} value
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.setTargetsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 12, value);
};


/**
 * @param {!proto.bucketeer.feature.Target=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Target}
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.addTargets = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 12, opt_value, proto.bucketeer.feature.Target, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureCreatedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureCreatedEvent.prototype.clearTargetsList = function() {
  return this.setTargetsList([]);
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
proto.bucketeer.event.domain.FeatureEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureEnabledEvent}
 */
proto.bucketeer.event.domain.FeatureEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureEnabledEvent;
  return proto.bucketeer.event.domain.FeatureEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureEnabledEvent}
 */
proto.bucketeer.event.domain.FeatureEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureEnabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureEnabledEvent} returns this
 */
proto.bucketeer.event.domain.FeatureEnabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.FeatureDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureDisabledEvent}
 */
proto.bucketeer.event.domain.FeatureDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureDisabledEvent;
  return proto.bucketeer.event.domain.FeatureDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureDisabledEvent}
 */
proto.bucketeer.event.domain.FeatureDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDisabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureDisabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureDisabledEvent} returns this
 */
proto.bucketeer.event.domain.FeatureDisabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.FeatureArchivedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureArchivedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureArchivedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureArchivedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureArchivedEvent}
 */
proto.bucketeer.event.domain.FeatureArchivedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureArchivedEvent;
  return proto.bucketeer.event.domain.FeatureArchivedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureArchivedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureArchivedEvent}
 */
proto.bucketeer.event.domain.FeatureArchivedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureArchivedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureArchivedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureArchivedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureArchivedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureArchivedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureArchivedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureArchivedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.FeatureUnarchivedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureUnarchivedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureUnarchivedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureUnarchivedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureUnarchivedEvent}
 */
proto.bucketeer.event.domain.FeatureUnarchivedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureUnarchivedEvent;
  return proto.bucketeer.event.domain.FeatureUnarchivedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureUnarchivedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureUnarchivedEvent}
 */
proto.bucketeer.event.domain.FeatureUnarchivedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureUnarchivedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureUnarchivedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureUnarchivedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureUnarchivedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureUnarchivedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureUnarchivedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureUnarchivedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.FeatureDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureDeletedEvent}
 */
proto.bucketeer.event.domain.FeatureDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureDeletedEvent;
  return proto.bucketeer.event.domain.FeatureDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureDeletedEvent}
 */
proto.bucketeer.event.domain.FeatureDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureDeletedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureDeletedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.EvaluationDelayableSetEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.EvaluationDelayableSetEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.EvaluationDelayableSetEvent}
 */
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.EvaluationDelayableSetEvent;
  return proto.bucketeer.event.domain.EvaluationDelayableSetEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.EvaluationDelayableSetEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.EvaluationDelayableSetEvent}
 */
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.EvaluationDelayableSetEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.EvaluationDelayableSetEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EvaluationDelayableSetEvent} returns this
 */
proto.bucketeer.event.domain.EvaluationDelayableSetEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.EvaluationUndelayableSetEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.EvaluationUndelayableSetEvent}
 */
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.EvaluationUndelayableSetEvent;
  return proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.EvaluationUndelayableSetEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.EvaluationUndelayableSetEvent}
 */
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.EvaluationUndelayableSetEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EvaluationUndelayableSetEvent} returns this
 */
proto.bucketeer.event.domain.EvaluationUndelayableSetEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.FeatureRenamedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureRenamedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureRenamedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureRenamedEvent}
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureRenamedEvent;
  return proto.bucketeer.event.domain.FeatureRenamedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureRenamedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureRenamedEvent}
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureRenamedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureRenamedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureRenamedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureRenamedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRenamedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRenamedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRenamedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureDescriptionChangedEvent;
  return proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    offVariation: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent}
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureOffVariationChangedEvent;
  return proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent}
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setOffVariation(value);
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
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getOffVariation();
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
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string off_variation = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.prototype.getOffVariation = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureOffVariationChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureOffVariationChangedEvent.prototype.setOffVariation = function(value) {
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
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureVariationAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureVariationAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    variation: (f = msg.getVariation()) && proto_feature_variation_pb.Variation.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.FeatureVariationAddedEvent}
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureVariationAddedEvent;
  return proto.bucketeer.event.domain.FeatureVariationAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureVariationAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureVariationAddedEvent}
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_variation_pb.Variation;
      reader.readMessage(value,proto_feature_variation_pb.Variation.deserializeBinaryFromReader);
      msg.setVariation(value);
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
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureVariationAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureVariationAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getVariation();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_variation_pb.Variation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureVariationAddedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.feature.Variation variation = 2;
 * @return {?proto.bucketeer.feature.Variation}
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.getVariation = function() {
  return /** @type{?proto.bucketeer.feature.Variation} */ (
    jspb.Message.getWrapperField(this, proto_feature_variation_pb.Variation, 2));
};


/**
 * @param {?proto.bucketeer.feature.Variation|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureVariationAddedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.setVariation = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureVariationAddedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.clearVariation = function() {
  return this.setVariation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureVariationAddedEvent.prototype.hasVariation = function() {
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
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureVariationRemovedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    variationId: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent}
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureVariationRemovedEvent;
  return proto.bucketeer.event.domain.FeatureVariationRemovedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent}
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setVariationId(value);
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
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureVariationRemovedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getVariationId();
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
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string variation_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.prototype.getVariationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureVariationRemovedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureVariationRemovedEvent.prototype.setVariationId = function(value) {
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
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.VariationValueChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.VariationValueChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
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
 * @return {!proto.bucketeer.event.domain.VariationValueChangedEvent}
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.VariationValueChangedEvent;
  return proto.bucketeer.event.domain.VariationValueChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.VariationValueChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.VariationValueChangedEvent}
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setId(value);
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
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.VariationValueChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.VariationValueChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getId();
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
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationValueChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationValueChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string value = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationValueChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationValueChangedEvent.prototype.setValue = function(value) {
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
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.VariationNameChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.VariationNameChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
    name: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.VariationNameChangedEvent}
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.VariationNameChangedEvent;
  return proto.bucketeer.event.domain.VariationNameChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.VariationNameChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.VariationNameChangedEvent}
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.VariationNameChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.VariationNameChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationNameChangedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.VariationDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
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
 * @return {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.VariationDescriptionChangedEvent;
  return proto.bucketeer.event.domain.VariationDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setId(value);
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
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.VariationDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getId();
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
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.VariationDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.VariationUserAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.VariationUserAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
    user: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.VariationUserAddedEvent}
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.VariationUserAddedEvent;
  return proto.bucketeer.event.domain.VariationUserAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.VariationUserAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.VariationUserAddedEvent}
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.VariationUserAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.VariationUserAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUser();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string user = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.VariationUserAddedEvent.prototype.setUser = function(value) {
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
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.VariationUserRemovedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.VariationUserRemovedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
    user: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.VariationUserRemovedEvent}
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.VariationUserRemovedEvent;
  return proto.bucketeer.event.domain.VariationUserRemovedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.VariationUserRemovedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.VariationUserRemovedEvent}
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.VariationUserRemovedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.VariationUserRemovedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUser();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationUserRemovedEvent} returns this
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationUserRemovedEvent} returns this
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string user = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.VariationUserRemovedEvent} returns this
 */
proto.bucketeer.event.domain.VariationUserRemovedEvent.prototype.setUser = function(value) {
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
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureRuleAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureRuleAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.FeatureRuleAddedEvent}
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureRuleAddedEvent;
  return proto.bucketeer.event.domain.FeatureRuleAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureRuleAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureRuleAddedEvent}
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureRuleAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureRuleAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRule();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_rule_pb.Rule.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRuleAddedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.feature.Rule rule = 2;
 * @return {?proto.bucketeer.feature.Rule}
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.getRule = function() {
  return /** @type{?proto.bucketeer.feature.Rule} */ (
    jspb.Message.getWrapperField(this, proto_feature_rule_pb.Rule, 2));
};


/**
 * @param {?proto.bucketeer.feature.Rule|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureRuleAddedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.setRule = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureRuleAddedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.clearRule = function() {
  return this.setRule(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureRuleAddedEvent.prototype.hasRule = function() {
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
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent}
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent;
  return proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent}
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} returns this
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} returns this
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bucketeer.feature.Strategy strategy = 3;
 * @return {?proto.bucketeer.feature.Strategy}
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.Strategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.Strategy, 3));
};


/**
 * @param {?proto.bucketeer.feature.Strategy|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} returns this
*/
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent} returns this
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureChangeRuleStrategyEvent.prototype.hasStrategy = function() {
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
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureRuleDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent}
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureRuleDeletedEvent;
  return proto.bucketeer.event.domain.FeatureRuleDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent}
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureRuleDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRuleDeletedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRuleDeletedEvent.prototype.setRuleId = function(value) {
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
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent}
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent;
  return proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent}
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bucketeer.feature.FixedStrategy strategy = 3;
 * @return {?proto.bucketeer.feature.FixedStrategy}
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.FixedStrategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.FixedStrategy, 3));
};


/**
 * @param {?proto.bucketeer.feature.FixedStrategy|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureFixedStrategyChangedEvent.prototype.hasStrategy = function() {
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
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent}
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent;
  return proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent}
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bucketeer.feature.RolloutStrategy strategy = 3;
 * @return {?proto.bucketeer.feature.RolloutStrategy}
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.RolloutStrategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.RolloutStrategy, 3));
};


/**
 * @param {?proto.bucketeer.feature.RolloutStrategy|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureRolloutStrategyChangedEvent.prototype.hasStrategy = function() {
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
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.RuleClauseAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.RuleClauseAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
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
 * @return {!proto.bucketeer.event.domain.RuleClauseAddedEvent}
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.RuleClauseAddedEvent;
  return proto.bucketeer.event.domain.RuleClauseAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.RuleClauseAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.RuleClauseAddedEvent}
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.RuleClauseAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.RuleClauseAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
  f = message.getClause();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_clause_pb.Clause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.RuleClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.RuleClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bucketeer.feature.Clause clause = 3;
 * @return {?proto.bucketeer.feature.Clause}
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.getClause = function() {
  return /** @type{?proto.bucketeer.feature.Clause} */ (
    jspb.Message.getWrapperField(this, proto_feature_clause_pb.Clause, 3));
};


/**
 * @param {?proto.bucketeer.feature.Clause|undefined} value
 * @return {!proto.bucketeer.event.domain.RuleClauseAddedEvent} returns this
*/
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.setClause = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.RuleClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.clearClause = function() {
  return this.setClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.RuleClauseAddedEvent.prototype.hasClause = function() {
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
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.RuleClauseDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.RuleClauseDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    id: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.RuleClauseDeletedEvent}
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.RuleClauseDeletedEvent;
  return proto.bucketeer.event.domain.RuleClauseDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.RuleClauseDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.RuleClauseDeletedEvent}
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.RuleClauseDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.RuleClauseDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.RuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.RuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.RuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.RuleClauseDeletedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ClauseAttributeChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    id: jspb.Message.getFieldWithDefault(msg, 3, ""),
    attribute: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent}
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ClauseAttributeChangedEvent;
  return proto.bucketeer.event.domain.ClauseAttributeChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent}
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ClauseAttributeChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAttribute();
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
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string attribute = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.getAttribute = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseAttributeChangedEvent.prototype.setAttribute = function(value) {
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
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ClauseOperatorChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    id: jspb.Message.getFieldWithDefault(msg, 3, ""),
    operator: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent}
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ClauseOperatorChangedEvent;
  return proto.bucketeer.event.domain.ClauseOperatorChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent}
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ClauseOperatorChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getOperator();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional bucketeer.feature.Clause.Operator operator = 4;
 * @return {!proto.bucketeer.feature.Clause.Operator}
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.getOperator = function() {
  return /** @type {!proto.bucketeer.feature.Clause.Operator} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.feature.Clause.Operator} value
 * @return {!proto.bucketeer.event.domain.ClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseOperatorChangedEvent.prototype.setOperator = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
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
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ClauseValueAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ClauseValueAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    id: jspb.Message.getFieldWithDefault(msg, 3, ""),
    value: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.ClauseValueAddedEvent}
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ClauseValueAddedEvent;
  return proto.bucketeer.event.domain.ClauseValueAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ClauseValueAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ClauseValueAddedEvent}
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ClauseValueAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ClauseValueAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string value = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueAddedEvent.prototype.setValue = function(value) {
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
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ClauseValueRemovedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    id: jspb.Message.getFieldWithDefault(msg, 3, ""),
    value: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.ClauseValueRemovedEvent}
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ClauseValueRemovedEvent;
  return proto.bucketeer.event.domain.ClauseValueRemovedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ClauseValueRemovedEvent}
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ClauseValueRemovedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
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
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string value = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.ClauseValueRemovedEvent.prototype.setValue = function(value) {
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
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent}
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent;
  return proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent}
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStrategy();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_strategy_pb.Strategy.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.feature.Strategy strategy = 2;
 * @return {?proto.bucketeer.feature.Strategy}
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.getStrategy = function() {
  return /** @type{?proto.bucketeer.feature.Strategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.Strategy, 2));
};


/**
 * @param {?proto.bucketeer.feature.Strategy|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.setStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.clearStrategy = function() {
  return this.setStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureDefaultStrategyChangedEvent.prototype.hasStrategy = function() {
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
proto.bucketeer.event.domain.FeatureTagAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureTagAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureTagAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    tag: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.event.domain.FeatureTagAddedEvent}
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureTagAddedEvent;
  return proto.bucketeer.event.domain.FeatureTagAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureTagAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureTagAddedEvent}
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureTagAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureTagAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureTagAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTag();
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
proto.bucketeer.event.domain.FeatureTagAddedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureTagAddedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string tag = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.prototype.getTag = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureTagAddedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureTagAddedEvent.prototype.setTag = function(value) {
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
proto.bucketeer.event.domain.FeatureTagRemovedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureTagRemovedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureTagRemovedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    tag: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.event.domain.FeatureTagRemovedEvent}
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureTagRemovedEvent;
  return proto.bucketeer.event.domain.FeatureTagRemovedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureTagRemovedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureTagRemovedEvent}
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.FeatureTagRemovedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureTagRemovedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureTagRemovedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTag();
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
proto.bucketeer.event.domain.FeatureTagRemovedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureTagRemovedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string tag = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.prototype.getTag = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureTagRemovedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureTagRemovedEvent.prototype.setTag = function(value) {
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
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    version: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent}
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureVersionIncrementedEvent;
  return proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent}
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt32());
      msg.setVersion(value);
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
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getVersion();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int32 version = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.prototype.getVersion = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.FeatureVersionIncrementedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureVersionIncrementedEvent.prototype.setVersion = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.FeatureClonedEvent.repeatedFields_ = [4,5,6,9,12];



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
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureClonedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureClonedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureClonedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, ""),
    variationsList: jspb.Message.toObjectList(msg.getVariationsList(),
    proto_feature_variation_pb.Variation.toObject, includeInstance),
    targetsList: jspb.Message.toObjectList(msg.getTargetsList(),
    proto_feature_target_pb.Target.toObject, includeInstance),
    rulesList: jspb.Message.toObjectList(msg.getRulesList(),
    proto_feature_rule_pb.Rule.toObject, includeInstance),
    defaultStrategy: (f = msg.getDefaultStrategy()) && proto_feature_strategy_pb.Strategy.toObject(includeInstance, f),
    offVariation: jspb.Message.getFieldWithDefault(msg, 8, ""),
    tagsList: (f = jspb.Message.getRepeatedField(msg, 9)) == null ? undefined : f,
    maintainer: jspb.Message.getFieldWithDefault(msg, 10, ""),
    variationType: jspb.Message.getFieldWithDefault(msg, 11, 0),
    prerequisitesList: jspb.Message.toObjectList(msg.getPrerequisitesList(),
    proto_feature_prerequisite_pb.Prerequisite.toObject, includeInstance)
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
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureClonedEvent;
  return proto.bucketeer.event.domain.FeatureClonedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureClonedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = new proto_feature_target_pb.Target;
      reader.readMessage(value,proto_feature_target_pb.Target.deserializeBinaryFromReader);
      msg.addTargets(value);
      break;
    case 6:
      var value = new proto_feature_rule_pb.Rule;
      reader.readMessage(value,proto_feature_rule_pb.Rule.deserializeBinaryFromReader);
      msg.addRules(value);
      break;
    case 7:
      var value = new proto_feature_strategy_pb.Strategy;
      reader.readMessage(value,proto_feature_strategy_pb.Strategy.deserializeBinaryFromReader);
      msg.setDefaultStrategy(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setOffVariation(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addTags(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setMaintainer(value);
      break;
    case 11:
      var value = /** @type {!proto.bucketeer.feature.Feature.VariationType} */ (reader.readEnum());
      msg.setVariationType(value);
      break;
    case 12:
      var value = new proto_feature_prerequisite_pb.Prerequisite;
      reader.readMessage(value,proto_feature_prerequisite_pb.Prerequisite.deserializeBinaryFromReader);
      msg.addPrerequisites(value);
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
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureClonedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureClonedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureClonedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getTargetsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto_feature_target_pb.Target.serializeBinaryToWriter
    );
  }
  f = message.getRulesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      6,
      f,
      proto_feature_rule_pb.Rule.serializeBinaryToWriter
    );
  }
  f = message.getDefaultStrategy();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto_feature_strategy_pb.Strategy.serializeBinaryToWriter
    );
  }
  f = message.getOffVariation();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
  f = message.getMaintainer();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getVariationType();
  if (f !== 0.0) {
    writer.writeEnum(
      11,
      f
    );
  }
  f = message.getPrerequisitesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      12,
      f,
      proto_feature_prerequisite_pb.Prerequisite.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated bucketeer.feature.Variation variations = 4;
 * @return {!Array<!proto.bucketeer.feature.Variation>}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getVariationsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Variation>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_variation_pb.Variation, 4));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Variation>} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setVariationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.bucketeer.feature.Variation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Variation}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.addVariations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.bucketeer.feature.Variation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.clearVariationsList = function() {
  return this.setVariationsList([]);
};


/**
 * repeated bucketeer.feature.Target targets = 5;
 * @return {!Array<!proto.bucketeer.feature.Target>}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getTargetsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Target>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_target_pb.Target, 5));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Target>} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setTargetsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.bucketeer.feature.Target=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Target}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.addTargets = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.bucketeer.feature.Target, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.clearTargetsList = function() {
  return this.setTargetsList([]);
};


/**
 * repeated bucketeer.feature.Rule rules = 6;
 * @return {!Array<!proto.bucketeer.feature.Rule>}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getRulesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Rule>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_rule_pb.Rule, 6));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Rule>} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setRulesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 6, value);
};


/**
 * @param {!proto.bucketeer.feature.Rule=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Rule}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.addRules = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 6, opt_value, proto.bucketeer.feature.Rule, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.clearRulesList = function() {
  return this.setRulesList([]);
};


/**
 * optional bucketeer.feature.Strategy default_strategy = 7;
 * @return {?proto.bucketeer.feature.Strategy}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getDefaultStrategy = function() {
  return /** @type{?proto.bucketeer.feature.Strategy} */ (
    jspb.Message.getWrapperField(this, proto_feature_strategy_pb.Strategy, 7));
};


/**
 * @param {?proto.bucketeer.feature.Strategy|undefined} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setDefaultStrategy = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.clearDefaultStrategy = function() {
  return this.setDefaultStrategy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.hasDefaultStrategy = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string off_variation = 8;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getOffVariation = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setOffVariation = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * repeated string tags = 9;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * optional string maintainer = 10;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getMaintainer = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setMaintainer = function(value) {
  return jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional bucketeer.feature.Feature.VariationType variation_type = 11;
 * @return {!proto.bucketeer.feature.Feature.VariationType}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getVariationType = function() {
  return /** @type {!proto.bucketeer.feature.Feature.VariationType} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {!proto.bucketeer.feature.Feature.VariationType} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setVariationType = function(value) {
  return jspb.Message.setProto3EnumField(this, 11, value);
};


/**
 * repeated bucketeer.feature.Prerequisite prerequisites = 12;
 * @return {!Array<!proto.bucketeer.feature.Prerequisite>}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.getPrerequisitesList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Prerequisite>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_prerequisite_pb.Prerequisite, 12));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Prerequisite>} value
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
*/
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.setPrerequisitesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 12, value);
};


/**
 * @param {!proto.bucketeer.feature.Prerequisite=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Prerequisite}
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.addPrerequisites = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 12, opt_value, proto.bucketeer.feature.Prerequisite, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.FeatureClonedEvent} returns this
 */
proto.bucketeer.event.domain.FeatureClonedEvent.prototype.clearPrerequisitesList = function() {
  return this.setPrerequisitesList([]);
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
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    samplingSeed: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent}
 */
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent;
  return proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent}
 */
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSamplingSeed(value);
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
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSamplingSeed();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string sampling_seed = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.prototype.getSamplingSeed = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent} returns this
 */
proto.bucketeer.event.domain.FeatureSamplingSeedResetEvent.prototype.setSamplingSeed = function(value) {
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
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.GoalCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.GoalCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, ""),
    deleted: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
    createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0)
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
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.GoalCreatedEvent;
  return proto.bucketeer.event.domain.GoalCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.GoalCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDeleted(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
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
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.GoalCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.GoalCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getDeleted();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent} returns this
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent} returns this
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent} returns this
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional bool deleted = 4;
 * @return {boolean}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.getDeleted = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent} returns this
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.setDeleted = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent} returns this
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.GoalCreatedEvent} returns this
 */
proto.bucketeer.event.domain.GoalCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
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
proto.bucketeer.event.domain.GoalRenamedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.GoalRenamedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.GoalRenamedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalRenamedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.GoalRenamedEvent}
 */
proto.bucketeer.event.domain.GoalRenamedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.GoalRenamedEvent;
  return proto.bucketeer.event.domain.GoalRenamedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.GoalRenamedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.GoalRenamedEvent}
 */
proto.bucketeer.event.domain.GoalRenamedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.GoalRenamedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.GoalRenamedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.GoalRenamedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalRenamedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.GoalRenamedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalRenamedEvent} returns this
 */
proto.bucketeer.event.domain.GoalRenamedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.GoalRenamedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalRenamedEvent} returns this
 */
proto.bucketeer.event.domain.GoalRenamedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.GoalDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.GoalDescriptionChangedEvent;
  return proto.bucketeer.event.domain.GoalDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.GoalDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.GoalDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.GoalArchivedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.GoalArchivedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.GoalArchivedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalArchivedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.GoalArchivedEvent}
 */
proto.bucketeer.event.domain.GoalArchivedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.GoalArchivedEvent;
  return proto.bucketeer.event.domain.GoalArchivedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.GoalArchivedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.GoalArchivedEvent}
 */
proto.bucketeer.event.domain.GoalArchivedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.GoalArchivedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.GoalArchivedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.GoalArchivedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalArchivedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.GoalArchivedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalArchivedEvent} returns this
 */
proto.bucketeer.event.domain.GoalArchivedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.GoalDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.GoalDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.GoalDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.GoalDeletedEvent}
 */
proto.bucketeer.event.domain.GoalDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.GoalDeletedEvent;
  return proto.bucketeer.event.domain.GoalDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.GoalDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.GoalDeletedEvent}
 */
proto.bucketeer.event.domain.GoalDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.GoalDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.GoalDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.GoalDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.GoalDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.GoalDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.GoalDeletedEvent} returns this
 */
proto.bucketeer.event.domain.GoalDeletedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.repeatedFields_ = [4,12];



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
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    featureId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    featureVersion: jspb.Message.getFieldWithDefault(msg, 3, 0),
    variationsList: jspb.Message.toObjectList(msg.getVariationsList(),
    proto_feature_variation_pb.Variation.toObject, includeInstance),
    goalId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    startAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    stopAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    stopped: jspb.Message.getBooleanFieldWithDefault(msg, 8, false),
    stoppedAt: jspb.Message.getFieldWithDefault(msg, 9, 0),
    createdAt: jspb.Message.getFieldWithDefault(msg, 10, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 11, 0),
    goalIdsList: (f = jspb.Message.getRepeatedField(msg, 12)) == null ? undefined : f,
    name: jspb.Message.getFieldWithDefault(msg, 13, ""),
    description: jspb.Message.getFieldWithDefault(msg, 14, ""),
    baseVariationId: jspb.Message.getFieldWithDefault(msg, 15, "")
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
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentCreatedEvent;
  return proto.bucketeer.event.domain.ExperimentCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setFeatureId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setFeatureVersion(value);
      break;
    case 4:
      var value = new proto_feature_variation_pb.Variation;
      reader.readMessage(value,proto_feature_variation_pb.Variation.deserializeBinaryFromReader);
      msg.addVariations(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setGoalId(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStartAt(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStopAt(value);
      break;
    case 8:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setStopped(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStoppedAt(value);
      break;
    case 10:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 12:
      var value = /** @type {string} */ (reader.readString());
      msg.addGoalIds(value);
      break;
    case 13:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 14:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 15:
      var value = /** @type {string} */ (reader.readString());
      msg.setBaseVariationId(value);
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
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getFeatureVersion();
  if (f !== 0) {
    writer.writeInt32(
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
  f = message.getGoalId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getStartAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getStopAt();
  if (f !== 0) {
    writer.writeInt64(
      7,
      f
    );
  }
  f = message.getStopped();
  if (f) {
    writer.writeBool(
      8,
      f
    );
  }
  f = message.getStoppedAt();
  if (f !== 0) {
    writer.writeInt64(
      9,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      10,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      11,
      f
    );
  }
  f = message.getGoalIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      12,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      13,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      14,
      f
    );
  }
  f = message.getBaseVariationId();
  if (f.length > 0) {
    writer.writeString(
      15,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int32 feature_version = 3;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getFeatureVersion = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setFeatureVersion = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * repeated bucketeer.feature.Variation variations = 4;
 * @return {!Array<!proto.bucketeer.feature.Variation>}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getVariationsList = function() {
  return /** @type{!Array<!proto.bucketeer.feature.Variation>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_feature_variation_pb.Variation, 4));
};


/**
 * @param {!Array<!proto.bucketeer.feature.Variation>} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
*/
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setVariationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.bucketeer.feature.Variation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Variation}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.addVariations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.bucketeer.feature.Variation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.clearVariationsList = function() {
  return this.setVariationsList([]);
};


/**
 * optional string goal_id = 5;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getGoalId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setGoalId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional int64 start_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getStartAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setStartAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 stop_at = 7;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getStopAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setStopAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional bool stopped = 8;
 * @return {boolean}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getStopped = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 8, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setStopped = function(value) {
  return jspb.Message.setProto3BooleanField(this, 8, value);
};


/**
 * optional int64 stopped_at = 9;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getStoppedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setStoppedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * optional int64 created_at = 10;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 10, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 10, value);
};


/**
 * optional int64 updated_at = 11;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 11, value);
};


/**
 * repeated string goal_ids = 12;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getGoalIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 12));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setGoalIdsList = function(value) {
  return jspb.Message.setField(this, 12, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.addGoalIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 12, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.clearGoalIdsList = function() {
  return this.setGoalIdsList([]);
};


/**
 * optional string name = 13;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 13, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 13, value);
};


/**
 * optional string description = 14;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 14, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 14, value);
};


/**
 * optional string base_variation_id = 15;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.getBaseVariationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 15, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentCreatedEvent.prototype.setBaseVariationId = function(value) {
  return jspb.Message.setProto3StringField(this, 15, value);
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
proto.bucketeer.event.domain.ExperimentStoppedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentStoppedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentStoppedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stoppedAt: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.event.domain.ExperimentStoppedEvent}
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentStoppedEvent;
  return proto.bucketeer.event.domain.ExperimentStoppedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentStoppedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentStoppedEvent}
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStoppedAt(value);
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
proto.bucketeer.event.domain.ExperimentStoppedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentStoppedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentStoppedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStoppedAt();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentStoppedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 stopped_at = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.prototype.getStoppedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentStoppedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentStoppedEvent.prototype.setStoppedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
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
proto.bucketeer.event.domain.ExperimentArchivedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentArchivedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentArchivedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentArchivedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ExperimentArchivedEvent}
 */
proto.bucketeer.event.domain.ExperimentArchivedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentArchivedEvent;
  return proto.bucketeer.event.domain.ExperimentArchivedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentArchivedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentArchivedEvent}
 */
proto.bucketeer.event.domain.ExperimentArchivedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ExperimentArchivedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentArchivedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentArchivedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentArchivedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ExperimentArchivedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentArchivedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentArchivedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.ExperimentDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ExperimentDeletedEvent}
 */
proto.bucketeer.event.domain.ExperimentDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentDeletedEvent;
  return proto.bucketeer.event.domain.ExperimentDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentDeletedEvent}
 */
proto.bucketeer.event.domain.ExperimentDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ExperimentDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ExperimentDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentDeletedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentDeletedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    startAt: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentStartAtChangedEvent;
  return proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStartAt(value);
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
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStartAt();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 start_at = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.prototype.getStartAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentStartAtChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentStartAtChangedEvent.prototype.setStartAt = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
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
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stopAt: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentStopAtChangedEvent;
  return proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStopAt(value);
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
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStopAt();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 stop_at = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.prototype.getStopAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentStopAtChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentStopAtChangedEvent.prototype.setStopAt = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
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
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    startAt: jspb.Message.getFieldWithDefault(msg, 2, 0),
    stopAt: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentPeriodChangedEvent;
  return proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStartAt(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setStopAt(value);
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
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStartAt();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getStopAt();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 start_at = 2;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.getStartAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.setStartAt = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int64 stop_at = 3;
 * @return {number}
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.getStopAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ExperimentPeriodChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentPeriodChangedEvent.prototype.setStopAt = function(value) {
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
proto.bucketeer.event.domain.ExperimentNameChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentNameChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentNameChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ExperimentNameChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentNameChangedEvent;
  return proto.bucketeer.event.domain.ExperimentNameChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentNameChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentNameChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ExperimentNameChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentNameChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentNameChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ExperimentNameChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentNameChangedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent;
  return proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.ExperimentDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.ExperimentStartedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentStartedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentStartedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStartedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ExperimentStartedEvent}
 */
proto.bucketeer.event.domain.ExperimentStartedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentStartedEvent;
  return proto.bucketeer.event.domain.ExperimentStartedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentStartedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentStartedEvent}
 */
proto.bucketeer.event.domain.ExperimentStartedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ExperimentStartedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentStartedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentStartedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentStartedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ExperimentFinishedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ExperimentFinishedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ExperimentFinishedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentFinishedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ExperimentFinishedEvent}
 */
proto.bucketeer.event.domain.ExperimentFinishedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ExperimentFinishedEvent;
  return proto.bucketeer.event.domain.ExperimentFinishedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ExperimentFinishedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ExperimentFinishedEvent}
 */
proto.bucketeer.event.domain.ExperimentFinishedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ExperimentFinishedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ExperimentFinishedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ExperimentFinishedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ExperimentFinishedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AccountCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AccountCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    email: jspb.Message.getFieldWithDefault(msg, 2, ""),
    name: jspb.Message.getFieldWithDefault(msg, 3, ""),
    role: jspb.Message.getFieldWithDefault(msg, 4, 0),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 5, false),
    createdAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 7, 0)
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
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AccountCreatedEvent;
  return proto.bucketeer.event.domain.AccountCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AccountCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setEmail(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.account.Account.Role} */ (reader.readEnum());
      msg.setRole(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
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
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AccountCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AccountCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      7,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string email = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional bucketeer.account.Account.Role role = 4;
 * @return {!proto.bucketeer.account.Account.Role}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.Account.Role} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.account.Account.Role} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional bool disabled = 5;
 * @return {boolean}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AccountCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
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
proto.bucketeer.event.domain.AccountRoleChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AccountRoleChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AccountRoleChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    role: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.bucketeer.event.domain.AccountRoleChangedEvent}
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AccountRoleChangedEvent;
  return proto.bucketeer.event.domain.AccountRoleChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AccountRoleChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AccountRoleChangedEvent}
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.bucketeer.account.Account.Role} */ (reader.readEnum());
      msg.setRole(value);
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
proto.bucketeer.event.domain.AccountRoleChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AccountRoleChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AccountRoleChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountRoleChangedEvent} returns this
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.account.Account.Role role = 2;
 * @return {!proto.bucketeer.account.Account.Role}
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.Account.Role} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.account.Account.Role} value
 * @return {!proto.bucketeer.event.domain.AccountRoleChangedEvent} returns this
 */
proto.bucketeer.event.domain.AccountRoleChangedEvent.prototype.setRole = function(value) {
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
proto.bucketeer.event.domain.AccountEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AccountEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AccountEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AccountEnabledEvent}
 */
proto.bucketeer.event.domain.AccountEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AccountEnabledEvent;
  return proto.bucketeer.event.domain.AccountEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AccountEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AccountEnabledEvent}
 */
proto.bucketeer.event.domain.AccountEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AccountEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AccountEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AccountEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AccountEnabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountEnabledEvent} returns this
 */
proto.bucketeer.event.domain.AccountEnabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.AccountDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AccountDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AccountDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AccountDisabledEvent}
 */
proto.bucketeer.event.domain.AccountDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AccountDisabledEvent;
  return proto.bucketeer.event.domain.AccountDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AccountDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AccountDisabledEvent}
 */
proto.bucketeer.event.domain.AccountDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AccountDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AccountDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AccountDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountDisabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AccountDisabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountDisabledEvent} returns this
 */
proto.bucketeer.event.domain.AccountDisabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.AccountDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AccountDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AccountDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AccountDeletedEvent}
 */
proto.bucketeer.event.domain.AccountDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AccountDeletedEvent;
  return proto.bucketeer.event.domain.AccountDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AccountDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AccountDeletedEvent}
 */
proto.bucketeer.event.domain.AccountDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AccountDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AccountDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AccountDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AccountDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AccountDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AccountDeletedEvent} returns this
 */
proto.bucketeer.event.domain.AccountDeletedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.APIKeyCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.APIKeyCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    role: jspb.Message.getFieldWithDefault(msg, 3, 0),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
    createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0)
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
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.APIKeyCreatedEvent;
  return proto.bucketeer.event.domain.APIKeyCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.APIKeyCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.bucketeer.account.APIKey.Role} */ (reader.readEnum());
      msg.setRole(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
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
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.APIKeyCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.APIKeyCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bucketeer.account.APIKey.Role role = 3;
 * @return {!proto.bucketeer.account.APIKey.Role}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.APIKey.Role} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.account.APIKey.Role} value
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional bool disabled = 4;
 * @return {boolean}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.APIKeyCreatedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
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
proto.bucketeer.event.domain.APIKeyNameChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.APIKeyNameChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.APIKeyNameChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.APIKeyNameChangedEvent}
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.APIKeyNameChangedEvent;
  return proto.bucketeer.event.domain.APIKeyNameChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.APIKeyNameChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.APIKeyNameChangedEvent}
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.APIKeyNameChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.APIKeyNameChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.APIKeyNameChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.APIKeyNameChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.APIKeyNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.APIKeyNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyNameChangedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.APIKeyEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.APIKeyEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.APIKeyEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.APIKeyEnabledEvent}
 */
proto.bucketeer.event.domain.APIKeyEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.APIKeyEnabledEvent;
  return proto.bucketeer.event.domain.APIKeyEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.APIKeyEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.APIKeyEnabledEvent}
 */
proto.bucketeer.event.domain.APIKeyEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.APIKeyEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.APIKeyEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.APIKeyEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.APIKeyEnabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.APIKeyEnabledEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyEnabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.APIKeyDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.APIKeyDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.APIKeyDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.APIKeyDisabledEvent}
 */
proto.bucketeer.event.domain.APIKeyDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.APIKeyDisabledEvent;
  return proto.bucketeer.event.domain.APIKeyDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.APIKeyDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.APIKeyDisabledEvent}
 */
proto.bucketeer.event.domain.APIKeyDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.APIKeyDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.APIKeyDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.APIKeyDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.APIKeyDisabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.APIKeyDisabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.APIKeyDisabledEvent} returns this
 */
proto.bucketeer.event.domain.APIKeyDisabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.SegmentCreatedEvent}
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentCreatedEvent;
  return proto.bucketeer.event.domain.SegmentCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentCreatedEvent}
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentCreatedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.SegmentDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SegmentDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentDeletedEvent;
  return proto.bucketeer.event.domain.SegmentDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SegmentDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SegmentDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentDeletedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.SegmentNameChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentNameChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentNameChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SegmentNameChangedEvent}
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentNameChangedEvent;
  return proto.bucketeer.event.domain.SegmentNameChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentNameChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentNameChangedEvent}
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SegmentNameChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentNameChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentNameChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SegmentNameChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentNameChangedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentDescriptionChangedEvent;
  return proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentRuleAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentRuleAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
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
 * @return {!proto.bucketeer.event.domain.SegmentRuleAddedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentRuleAddedEvent;
  return proto.bucketeer.event.domain.SegmentRuleAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentRuleAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentRuleAddedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentRuleAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentRuleAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRule();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_feature_rule_pb.Rule.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.feature.Rule rule = 2;
 * @return {?proto.bucketeer.feature.Rule}
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.getRule = function() {
  return /** @type{?proto.bucketeer.feature.Rule} */ (
    jspb.Message.getWrapperField(this, proto_feature_rule_pb.Rule, 2));
};


/**
 * @param {?proto.bucketeer.feature.Rule|undefined} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleAddedEvent} returns this
*/
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.setRule = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.SegmentRuleAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.clearRule = function() {
  return this.setRule(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.SegmentRuleAddedEvent.prototype.hasRule = function() {
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
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentRuleDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentRuleDeletedEvent;
  return proto.bucketeer.event.domain.SegmentRuleDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentRuleDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleDeletedEvent.prototype.setRuleId = function(value) {
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
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
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
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent;
  return proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
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
  f = message.getClause();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto_feature_clause_pb.Clause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bucketeer.feature.Clause clause = 3;
 * @return {?proto.bucketeer.feature.Clause}
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.getClause = function() {
  return /** @type{?proto.bucketeer.feature.Clause} */ (
    jspb.Message.getWrapperField(this, proto_feature_clause_pb.Clause, 3));
};


/**
 * @param {?proto.bucketeer.feature.Clause|undefined} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} returns this
*/
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.setClause = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.clearClause = function() {
  return this.setClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.SegmentRuleClauseAddedEvent.prototype.hasClause = function() {
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
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    clauseId: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent;
  return proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
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
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string clause_id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentRuleClauseDeletedEvent.prototype.setClauseId = function(value) {
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
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    clauseId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    attribute: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent;
  return proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
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
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getAttribute();
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
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string clause_id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string attribute = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.getAttribute = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseAttributeChangedEvent.prototype.setAttribute = function(value) {
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
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    clauseId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    operator: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent;
  return proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
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
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getOperator();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
};


/**
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string clause_id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional bucketeer.feature.Clause.Operator operator = 4;
 * @return {!proto.bucketeer.feature.Clause.Operator}
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.getOperator = function() {
  return /** @type {!proto.bucketeer.feature.Clause.Operator} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.feature.Clause.Operator} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseOperatorChangedEvent.prototype.setOperator = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
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
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    clauseId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    value: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentClauseValueAddedEvent;
  return proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
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
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string clause_id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string value = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueAddedEvent.prototype.setValue = function(value) {
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
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    ruleId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    clauseId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    value: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent;
  return proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent}
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setRuleId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
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
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getValue();
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
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string rule_id = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.getRuleId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.setRuleId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string clause_id = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string value = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentClauseValueRemovedEvent.prototype.setValue = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.repeatedFields_ = [2];



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
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentUserAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentUserAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userIdsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
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
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent}
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentUserAddedEvent;
  return proto.bucketeer.event.domain.SegmentUserAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentUserAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent}
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.addUserIds(value);
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
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentUserAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentUserAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
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
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string user_ids = 2;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.getUserIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.setUserIdsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.addUserIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.clearUserIdsList = function() {
  return this.setUserIdsList([]);
};


/**
 * optional bucketeer.feature.SegmentUser.State state = 3;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.event.domain.SegmentUserAddedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserAddedEvent.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.repeatedFields_ = [2];



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
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentUserDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userIdsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
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
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentUserDeletedEvent;
  return proto.bucketeer.event.domain.SegmentUserDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent}
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.addUserIds(value);
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
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentUserDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
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
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string user_ids = 2;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.getUserIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.setUserIdsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.addUserIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.clearUserIdsList = function() {
  return this.setUserIdsList([]);
};


/**
 * optional bucketeer.feature.SegmentUser.State state = 3;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.event.domain.SegmentUserDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentUserDeletedEvent.prototype.setState = function(value) {
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
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    status: jspb.Message.getFieldWithDefault(msg, 2, 0),
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
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent;
  return proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.bucketeer.feature.Segment.Status} */ (reader.readEnum());
      msg.setStatus(value);
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
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
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
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.feature.Segment.Status status = 2;
 * @return {!proto.bucketeer.feature.Segment.Status}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.getStatus = function() {
  return /** @type {!proto.bucketeer.feature.Segment.Status} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.Segment.Status} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional bucketeer.feature.SegmentUser.State state = 3;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersEvent.prototype.setState = function(value) {
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
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    segmentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    status: jspb.Message.getFieldWithDefault(msg, 2, 0),
    state: jspb.Message.getFieldWithDefault(msg, 3, 0),
    count: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent;
  return proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.bucketeer.feature.Segment.Status} */ (reader.readEnum());
      msg.setStatus(value);
      break;
    case 3:
      var value = /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSegmentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(
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
  f = message.getCount();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
};


/**
 * optional string segment_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.getSegmentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.setSegmentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.feature.Segment.Status status = 2;
 * @return {!proto.bucketeer.feature.Segment.Status}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.getStatus = function() {
  return /** @type {!proto.bucketeer.feature.Segment.Status} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.feature.Segment.Status} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.setStatus = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional bucketeer.feature.SegmentUser.State state = 3;
 * @return {!proto.bucketeer.feature.SegmentUser.State}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.getState = function() {
  return /** @type {!proto.bucketeer.feature.SegmentUser.State} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {!proto.bucketeer.feature.SegmentUser.State} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional int64 count = 4;
 * @return {number}
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.getCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent} returns this
 */
proto.bucketeer.event.domain.SegmentBulkUploadUsersStatusChangedEvent.prototype.setCount = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
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
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.EnvironmentCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    namespace: jspb.Message.getFieldWithDefault(msg, 2, ""),
    name: jspb.Message.getFieldWithDefault(msg, 3, ""),
    description: jspb.Message.getFieldWithDefault(msg, 4, ""),
    deleted: jspb.Message.getBooleanFieldWithDefault(msg, 5, false),
    createdAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    projectId: jspb.Message.getFieldWithDefault(msg, 8, "")
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
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.EnvironmentCreatedEvent;
  return proto.bucketeer.event.domain.EnvironmentCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setNamespace(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDeleted(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectId(value);
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
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.EnvironmentCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getNamespace();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
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
  f = message.getDeleted();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      7,
      f
    );
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string namespace = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string description = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional bool deleted = 5;
 * @return {boolean}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getDeleted = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setDeleted = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional string project_id = 8;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.getProjectId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentCreatedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentCreatedEvent.prototype.setProjectId = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
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
proto.bucketeer.event.domain.EnvironmentRenamedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.EnvironmentRenamedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.EnvironmentRenamedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.EnvironmentRenamedEvent}
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.EnvironmentRenamedEvent;
  return proto.bucketeer.event.domain.EnvironmentRenamedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.EnvironmentRenamedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.EnvironmentRenamedEvent}
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.EnvironmentRenamedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.EnvironmentRenamedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.EnvironmentRenamedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.EnvironmentRenamedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentRenamedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentRenamedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentRenamedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent;
  return proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.EnvironmentDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.EnvironmentDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.EnvironmentDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    namespace: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.bucketeer.event.domain.EnvironmentDeletedEvent}
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.EnvironmentDeletedEvent;
  return proto.bucketeer.event.domain.EnvironmentDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.EnvironmentDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.EnvironmentDeletedEvent}
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.EnvironmentDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.EnvironmentDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.EnvironmentDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getNamespace();
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
proto.bucketeer.event.domain.EnvironmentDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentDeletedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string namespace = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.prototype.getNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.EnvironmentDeletedEvent} returns this
 */
proto.bucketeer.event.domain.EnvironmentDeletedEvent.prototype.setNamespace = function(value) {
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
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminAccountCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    email: jspb.Message.getFieldWithDefault(msg, 2, ""),
    name: jspb.Message.getFieldWithDefault(msg, 3, ""),
    role: jspb.Message.getFieldWithDefault(msg, 4, 0),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 5, false),
    createdAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 7, 0)
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
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminAccountCreatedEvent;
  return proto.bucketeer.event.domain.AdminAccountCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setEmail(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 4:
      var value = /** @type {!proto.bucketeer.account.Account.Role} */ (reader.readEnum());
      msg.setRole(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
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
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminAccountCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      7,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string email = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional bucketeer.account.Account.Role role = 4;
 * @return {!proto.bucketeer.account.Account.Role}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.Account.Role} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.account.Account.Role} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional bool disabled = 5;
 * @return {boolean}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AdminAccountCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
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
proto.bucketeer.event.domain.AdminAccountEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminAccountEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminAccountEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminAccountEnabledEvent}
 */
proto.bucketeer.event.domain.AdminAccountEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminAccountEnabledEvent;
  return proto.bucketeer.event.domain.AdminAccountEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminAccountEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminAccountEnabledEvent}
 */
proto.bucketeer.event.domain.AdminAccountEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminAccountEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminAccountEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminAccountEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AdminAccountEnabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminAccountEnabledEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountEnabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.AdminAccountDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminAccountDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminAccountDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminAccountDisabledEvent}
 */
proto.bucketeer.event.domain.AdminAccountDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminAccountDisabledEvent;
  return proto.bucketeer.event.domain.AdminAccountDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminAccountDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminAccountDisabledEvent}
 */
proto.bucketeer.event.domain.AdminAccountDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminAccountDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminAccountDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminAccountDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountDisabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AdminAccountDisabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminAccountDisabledEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountDisabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.AdminAccountDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminAccountDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminAccountDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminAccountDeletedEvent}
 */
proto.bucketeer.event.domain.AdminAccountDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminAccountDeletedEvent;
  return proto.bucketeer.event.domain.AdminAccountDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminAccountDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminAccountDeletedEvent}
 */
proto.bucketeer.event.domain.AdminAccountDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminAccountDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminAccountDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminAccountDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminAccountDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AdminAccountDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminAccountDeletedEvent} returns this
 */
proto.bucketeer.event.domain.AdminAccountDeletedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.repeatedFields_ = [3];



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
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    featureId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    opsType: jspb.Message.getFieldWithDefault(msg, 2, 0),
    clausesList: jspb.Message.toObjectList(msg.getClausesList(),
    proto_autoops_clause_pb.Clause.toObject, includeInstance),
    triggeredAt: jspb.Message.getFieldWithDefault(msg, 4, 0),
    createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0)
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
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent;
  return proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {!proto.bucketeer.autoops.OpsType} */ (reader.readEnum());
      msg.setOpsType(value);
      break;
    case 3:
      var value = new proto_autoops_clause_pb.Clause;
      reader.readMessage(value,proto_autoops_clause_pb.Clause.deserializeBinaryFromReader);
      msg.addClauses(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTriggeredAt(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
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
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getOpsType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getClausesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto_autoops_clause_pb.Clause.serializeBinaryToWriter
    );
  }
  f = message.getTriggeredAt();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
};


/**
 * optional string feature_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.OpsType ops_type = 2;
 * @return {!proto.bucketeer.autoops.OpsType}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.getOpsType = function() {
  return /** @type {!proto.bucketeer.autoops.OpsType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.autoops.OpsType} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.setOpsType = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * repeated bucketeer.autoops.Clause clauses = 3;
 * @return {!Array<!proto.bucketeer.autoops.Clause>}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.getClausesList = function() {
  return /** @type{!Array<!proto.bucketeer.autoops.Clause>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_autoops_clause_pb.Clause, 3));
};


/**
 * @param {!Array<!proto.bucketeer.autoops.Clause>} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
*/
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.setClausesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.bucketeer.autoops.Clause=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.Clause}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.addClauses = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.bucketeer.autoops.Clause, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.clearClausesList = function() {
  return this.setClausesList([]);
};


/**
 * optional int64 triggered_at = 4;
 * @return {number}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.getTriggeredAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.setTriggeredAt = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
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
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent;
  return proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    opsType: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent;
  return proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.bucketeer.autoops.OpsType} */ (reader.readEnum());
      msg.setOpsType(value);
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
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOpsType();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
};


/**
 * optional bucketeer.autoops.OpsType ops_type = 1;
 * @return {!proto.bucketeer.autoops.OpsType}
 */
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.prototype.getOpsType = function() {
  return /** @type {!proto.bucketeer.autoops.OpsType} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.bucketeer.autoops.OpsType} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleOpsTypeChangedEvent.prototype.setOpsType = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
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
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent;
  return proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleTriggeredAtChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    opsEventRateClause: (f = msg.getOpsEventRateClause()) && proto_autoops_clause_pb.OpsEventRateClause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent}
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent;
  return proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent}
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 2:
      var value = new proto_autoops_clause_pb.OpsEventRateClause;
      reader.readMessage(value,proto_autoops_clause_pb.OpsEventRateClause.deserializeBinaryFromReader);
      msg.setOpsEventRateClause(value);
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
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getOpsEventRateClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_autoops_clause_pb.OpsEventRateClause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.OpsEventRateClause ops_event_rate_clause = 2;
 * @return {?proto.bucketeer.autoops.OpsEventRateClause}
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.getOpsEventRateClause = function() {
  return /** @type{?proto.bucketeer.autoops.OpsEventRateClause} */ (
    jspb.Message.getWrapperField(this, proto_autoops_clause_pb.OpsEventRateClause, 2));
};


/**
 * @param {?proto.bucketeer.autoops.OpsEventRateClause|undefined} value
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent} returns this
*/
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.setOpsEventRateClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.clearOpsEventRateClause = function() {
  return this.setOpsEventRateClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.OpsEventRateClauseAddedEvent.prototype.hasOpsEventRateClause = function() {
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
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    opsEventRateClause: (f = msg.getOpsEventRateClause()) && proto_autoops_clause_pb.OpsEventRateClause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent}
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent;
  return proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent}
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 2:
      var value = new proto_autoops_clause_pb.OpsEventRateClause;
      reader.readMessage(value,proto_autoops_clause_pb.OpsEventRateClause.deserializeBinaryFromReader);
      msg.setOpsEventRateClause(value);
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
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getOpsEventRateClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_autoops_clause_pb.OpsEventRateClause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent} returns this
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.OpsEventRateClause ops_event_rate_clause = 2;
 * @return {?proto.bucketeer.autoops.OpsEventRateClause}
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.getOpsEventRateClause = function() {
  return /** @type{?proto.bucketeer.autoops.OpsEventRateClause} */ (
    jspb.Message.getWrapperField(this, proto_autoops_clause_pb.OpsEventRateClause, 2));
};


/**
 * @param {?proto.bucketeer.autoops.OpsEventRateClause|undefined} value
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent} returns this
*/
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.setOpsEventRateClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent} returns this
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.clearOpsEventRateClause = function() {
  return this.setOpsEventRateClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.OpsEventRateClauseChangedEvent.prototype.hasOpsEventRateClause = function() {
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
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, "")
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
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent;
  return proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent}
 */
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
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
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent} returns this
 */
proto.bucketeer.event.domain.AutoOpsRuleClauseDeletedEvent.prototype.setClauseId = function(value) {
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
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.DatetimeClauseAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    datetimeClause: (f = msg.getDatetimeClause()) && proto_autoops_clause_pb.DatetimeClause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent}
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.DatetimeClauseAddedEvent;
  return proto.bucketeer.event.domain.DatetimeClauseAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent}
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 2:
      var value = new proto_autoops_clause_pb.DatetimeClause;
      reader.readMessage(value,proto_autoops_clause_pb.DatetimeClause.deserializeBinaryFromReader);
      msg.setDatetimeClause(value);
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
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.DatetimeClauseAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDatetimeClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_autoops_clause_pb.DatetimeClause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.DatetimeClause datetime_clause = 2;
 * @return {?proto.bucketeer.autoops.DatetimeClause}
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.getDatetimeClause = function() {
  return /** @type{?proto.bucketeer.autoops.DatetimeClause} */ (
    jspb.Message.getWrapperField(this, proto_autoops_clause_pb.DatetimeClause, 2));
};


/**
 * @param {?proto.bucketeer.autoops.DatetimeClause|undefined} value
 * @return {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent} returns this
*/
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.setDatetimeClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.DatetimeClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.clearDatetimeClause = function() {
  return this.setDatetimeClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.DatetimeClauseAddedEvent.prototype.hasDatetimeClause = function() {
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
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.DatetimeClauseChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    datetimeClause: (f = msg.getDatetimeClause()) && proto_autoops_clause_pb.DatetimeClause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent}
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.DatetimeClauseChangedEvent;
  return proto.bucketeer.event.domain.DatetimeClauseChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent}
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 2:
      var value = new proto_autoops_clause_pb.DatetimeClause;
      reader.readMessage(value,proto_autoops_clause_pb.DatetimeClause.deserializeBinaryFromReader);
      msg.setDatetimeClause(value);
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
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.DatetimeClauseChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDatetimeClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_autoops_clause_pb.DatetimeClause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent} returns this
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.DatetimeClause datetime_clause = 2;
 * @return {?proto.bucketeer.autoops.DatetimeClause}
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.getDatetimeClause = function() {
  return /** @type{?proto.bucketeer.autoops.DatetimeClause} */ (
    jspb.Message.getWrapperField(this, proto_autoops_clause_pb.DatetimeClause, 2));
};


/**
 * @param {?proto.bucketeer.autoops.DatetimeClause|undefined} value
 * @return {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent} returns this
*/
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.setDatetimeClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.DatetimeClauseChangedEvent} returns this
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.clearDatetimeClause = function() {
  return this.setDatetimeClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.DatetimeClauseChangedEvent.prototype.hasDatetimeClause = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.PushCreatedEvent.repeatedFields_ = [3];



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
proto.bucketeer.event.domain.PushCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PushCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PushCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    fcmApiKey: jspb.Message.getFieldWithDefault(msg, 2, ""),
    tagsList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f,
    name: jspb.Message.getFieldWithDefault(msg, 4, "")
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
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent}
 */
proto.bucketeer.event.domain.PushCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PushCreatedEvent;
  return proto.bucketeer.event.domain.PushCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PushCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent}
 */
proto.bucketeer.event.domain.PushCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setFcmApiKey(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.addTags(value);
      break;
    case 4:
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
proto.bucketeer.event.domain.PushCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PushCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PushCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFcmApiKey();
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
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
};


/**
 * optional string fcm_api_key = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.getFcmApiKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent} returns this
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.setFcmApiKey = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated string tags = 3;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent} returns this
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent} returns this
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent} returns this
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};


/**
 * optional string name = 4;
 * @return {string}
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.PushCreatedEvent} returns this
 */
proto.bucketeer.event.domain.PushCreatedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.PushDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PushDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PushDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.PushDeletedEvent}
 */
proto.bucketeer.event.domain.PushDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PushDeletedEvent;
  return proto.bucketeer.event.domain.PushDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PushDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PushDeletedEvent}
 */
proto.bucketeer.event.domain.PushDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.PushDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PushDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PushDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.repeatedFields_ = [2];



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
proto.bucketeer.event.domain.PushTagsAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PushTagsAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PushTagsAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    tagsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f
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
 * @return {!proto.bucketeer.event.domain.PushTagsAddedEvent}
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PushTagsAddedEvent;
  return proto.bucketeer.event.domain.PushTagsAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PushTagsAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PushTagsAddedEvent}
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addTags(value);
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
proto.bucketeer.event.domain.PushTagsAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PushTagsAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PushTagsAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
};


/**
 * repeated string tags = 2;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.PushTagsAddedEvent} returns this
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.PushTagsAddedEvent} returns this
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.PushTagsAddedEvent} returns this
 */
proto.bucketeer.event.domain.PushTagsAddedEvent.prototype.clearTagsList = function() {
  return this.setTagsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.repeatedFields_ = [2];



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
proto.bucketeer.event.domain.PushTagsDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PushTagsDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PushTagsDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    tagsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f
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
 * @return {!proto.bucketeer.event.domain.PushTagsDeletedEvent}
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PushTagsDeletedEvent;
  return proto.bucketeer.event.domain.PushTagsDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PushTagsDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PushTagsDeletedEvent}
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addTags(value);
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
proto.bucketeer.event.domain.PushTagsDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PushTagsDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PushTagsDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTagsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
};


/**
 * repeated string tags = 2;
 * @return {!Array<string>}
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.prototype.getTagsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.event.domain.PushTagsDeletedEvent} returns this
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.prototype.setTagsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.PushTagsDeletedEvent} returns this
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.prototype.addTags = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.PushTagsDeletedEvent} returns this
 */
proto.bucketeer.event.domain.PushTagsDeletedEvent.prototype.clearTagsList = function() {
  return this.setTagsList([]);
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
proto.bucketeer.event.domain.PushRenamedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PushRenamedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PushRenamedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushRenamedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
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
 * @return {!proto.bucketeer.event.domain.PushRenamedEvent}
 */
proto.bucketeer.event.domain.PushRenamedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PushRenamedEvent;
  return proto.bucketeer.event.domain.PushRenamedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PushRenamedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PushRenamedEvent}
 */
proto.bucketeer.event.domain.PushRenamedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
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
proto.bucketeer.event.domain.PushRenamedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PushRenamedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PushRenamedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PushRenamedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.PushRenamedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.PushRenamedEvent} returns this
 */
proto.bucketeer.event.domain.PushRenamedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.repeatedFields_ = [1];



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
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceTypesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f,
    recipient: (f = msg.getRecipient()) && proto_notification_recipient_pb.Recipient.toObject(includeInstance, f),
    name: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent}
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionCreatedEvent;
  return proto.bucketeer.event.domain.SubscriptionCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent}
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSourceTypes(values[i]);
      }
      break;
    case 2:
      var value = new proto_notification_recipient_pb.Recipient;
      reader.readMessage(value,proto_notification_recipient_pb.Recipient.deserializeBinaryFromReader);
      msg.setRecipient(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
  f = message.getRecipient();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_notification_recipient_pb.Recipient.serializeBinaryToWriter
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * repeated bucketeer.notification.Subscription.SourceType source_types = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.getSourceTypesList = function() {
  return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.setSourceTypesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.addSourceTypes = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.clearSourceTypesList = function() {
  return this.setSourceTypesList([]);
};


/**
 * optional bucketeer.notification.Recipient recipient = 2;
 * @return {?proto.bucketeer.notification.Recipient}
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.getRecipient = function() {
  return /** @type{?proto.bucketeer.notification.Recipient} */ (
    jspb.Message.getWrapperField(this, proto_notification_recipient_pb.Recipient, 2));
};


/**
 * @param {?proto.bucketeer.notification.Recipient|undefined} value
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} returns this
*/
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.setRecipient = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.clearRecipient = function() {
  return this.setRecipient(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.hasRecipient = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionCreatedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.SubscriptionDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SubscriptionDeletedEvent}
 */
proto.bucketeer.event.domain.SubscriptionDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionDeletedEvent;
  return proto.bucketeer.event.domain.SubscriptionDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionDeletedEvent}
 */
proto.bucketeer.event.domain.SubscriptionDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SubscriptionDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SubscriptionEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SubscriptionEnabledEvent}
 */
proto.bucketeer.event.domain.SubscriptionEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionEnabledEvent;
  return proto.bucketeer.event.domain.SubscriptionEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionEnabledEvent}
 */
proto.bucketeer.event.domain.SubscriptionEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SubscriptionEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SubscriptionDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SubscriptionDisabledEvent}
 */
proto.bucketeer.event.domain.SubscriptionDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionDisabledEvent;
  return proto.bucketeer.event.domain.SubscriptionDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionDisabledEvent}
 */
proto.bucketeer.event.domain.SubscriptionDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SubscriptionDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionDisabledEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.repeatedFields_ = [1];



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
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceTypesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
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
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent}
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent;
  return proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent}
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSourceTypes(values[i]);
      }
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
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
};


/**
 * repeated bucketeer.notification.Subscription.SourceType source_types = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.prototype.getSourceTypesList = function() {
  return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.prototype.setSourceTypesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.prototype.addSourceTypes = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesAddedEvent.prototype.clearSourceTypesList = function() {
  return this.setSourceTypesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.repeatedFields_ = [1];



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
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceTypesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
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
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent}
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent;
  return proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent}
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSourceTypes(values[i]);
      }
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
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
};


/**
 * repeated bucketeer.notification.Subscription.SourceType source_types = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.prototype.getSourceTypesList = function() {
  return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.prototype.setSourceTypesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.prototype.addSourceTypes = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionSourceTypesDeletedEvent.prototype.clearSourceTypesList = function() {
  return this.setSourceTypesList([]);
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
proto.bucketeer.event.domain.SubscriptionRenamedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.SubscriptionRenamedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.SubscriptionRenamedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionRenamedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.SubscriptionRenamedEvent}
 */
proto.bucketeer.event.domain.SubscriptionRenamedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.SubscriptionRenamedEvent;
  return proto.bucketeer.event.domain.SubscriptionRenamedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.SubscriptionRenamedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.SubscriptionRenamedEvent}
 */
proto.bucketeer.event.domain.SubscriptionRenamedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.SubscriptionRenamedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.SubscriptionRenamedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.SubscriptionRenamedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.SubscriptionRenamedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.SubscriptionRenamedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.SubscriptionRenamedEvent} returns this
 */
proto.bucketeer.event.domain.SubscriptionRenamedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.repeatedFields_ = [1];



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
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceTypesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f,
    recipient: (f = msg.getRecipient()) && proto_notification_recipient_pb.Recipient.toObject(includeInstance, f),
    name: jspb.Message.getFieldWithDefault(msg, 3, "")
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSourceTypes(values[i]);
      }
      break;
    case 2:
      var value = new proto_notification_recipient_pb.Recipient;
      reader.readMessage(value,proto_notification_recipient_pb.Recipient.deserializeBinaryFromReader);
      msg.setRecipient(value);
      break;
    case 3:
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
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
  f = message.getRecipient();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_notification_recipient_pb.Recipient.serializeBinaryToWriter
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * repeated bucketeer.notification.Subscription.SourceType source_types = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.getSourceTypesList = function() {
  return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.setSourceTypesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.addSourceTypes = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.clearSourceTypesList = function() {
  return this.setSourceTypesList([]);
};


/**
 * optional bucketeer.notification.Recipient recipient = 2;
 * @return {?proto.bucketeer.notification.Recipient}
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.getRecipient = function() {
  return /** @type{?proto.bucketeer.notification.Recipient} */ (
    jspb.Message.getWrapperField(this, proto_notification_recipient_pb.Recipient, 2));
};


/**
 * @param {?proto.bucketeer.notification.Recipient|undefined} value
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} returns this
*/
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.setRecipient = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.clearRecipient = function() {
  return this.setRecipient(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.hasRecipient = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionCreatedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionDisabledEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.repeatedFields_ = [1];



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
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceTypesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSourceTypes(values[i]);
      }
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
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
};


/**
 * repeated bucketeer.notification.Subscription.SourceType source_types = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.prototype.getSourceTypesList = function() {
  return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.prototype.setSourceTypesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.prototype.addSourceTypes = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesAddedEvent.prototype.clearSourceTypesList = function() {
  return this.setSourceTypesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.repeatedFields_ = [1];



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
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    sourceTypesList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSourceTypes(values[i]);
      }
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
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSourceTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      1,
      f
    );
  }
};


/**
 * repeated bucketeer.notification.Subscription.SourceType source_types = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.prototype.getSourceTypesList = function() {
  return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.prototype.setSourceTypesList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.prototype.addSourceTypes = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionSourceTypesDeletedEvent.prototype.clearSourceTypesList = function() {
  return this.setSourceTypesList([]);
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
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent;
  return proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent}
 */
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent} returns this
 */
proto.bucketeer.event.domain.AdminSubscriptionRenamedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ProjectCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ProjectCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    description: jspb.Message.getFieldWithDefault(msg, 2, ""),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 3, false),
    trial: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
    creatorEmail: jspb.Message.getFieldWithDefault(msg, 5, ""),
    createdAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    name: jspb.Message.getFieldWithDefault(msg, 8, ""),
    urlCode: jspb.Message.getFieldWithDefault(msg, 9, "")
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
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ProjectCreatedEvent;
  return proto.bucketeer.event.domain.ProjectCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ProjectCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setTrial(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setCreatorEmail(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrlCode(value);
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
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ProjectCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ProjectCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
  f = message.getTrial();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getCreatorEmail();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      7,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getUrlCode();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool disabled = 3;
 * @return {boolean}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};


/**
 * optional bool trial = 4;
 * @return {boolean}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getTrial = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setTrial = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional string creator_email = 5;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getCreatorEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setCreatorEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional string name = 8;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string url_code = 9;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.getUrlCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectCreatedEvent.prototype.setUrlCode = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
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
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ProjectDescriptionChangedEvent;
  return proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.ProjectEnabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ProjectEnabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ProjectEnabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectEnabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ProjectEnabledEvent}
 */
proto.bucketeer.event.domain.ProjectEnabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ProjectEnabledEvent;
  return proto.bucketeer.event.domain.ProjectEnabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ProjectEnabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ProjectEnabledEvent}
 */
proto.bucketeer.event.domain.ProjectEnabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ProjectEnabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ProjectEnabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ProjectEnabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectEnabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ProjectEnabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectEnabledEvent} returns this
 */
proto.bucketeer.event.domain.ProjectEnabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.ProjectDisabledEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ProjectDisabledEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ProjectDisabledEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectDisabledEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ProjectDisabledEvent}
 */
proto.bucketeer.event.domain.ProjectDisabledEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ProjectDisabledEvent;
  return proto.bucketeer.event.domain.ProjectDisabledEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ProjectDisabledEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ProjectDisabledEvent}
 */
proto.bucketeer.event.domain.ProjectDisabledEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ProjectDisabledEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ProjectDisabledEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ProjectDisabledEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectDisabledEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ProjectDisabledEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectDisabledEvent} returns this
 */
proto.bucketeer.event.domain.ProjectDisabledEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ProjectTrialCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    description: jspb.Message.getFieldWithDefault(msg, 2, ""),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 3, false),
    trial: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
    creatorEmail: jspb.Message.getFieldWithDefault(msg, 5, ""),
    createdAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    name: jspb.Message.getFieldWithDefault(msg, 8, ""),
    urlCode: jspb.Message.getFieldWithDefault(msg, 9, "")
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
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ProjectTrialCreatedEvent;
  return proto.bucketeer.event.domain.ProjectTrialCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisabled(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setTrial(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setCreatorEmail(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 7:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrlCode(value);
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
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ProjectTrialCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getDisabled();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
  f = message.getTrial();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getCreatorEmail();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      7,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getUrlCode();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool disabled = 3;
 * @return {boolean}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};


/**
 * optional bool trial = 4;
 * @return {boolean}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getTrial = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setTrial = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional string creator_email = 5;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getCreatorEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setCreatorEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional string name = 8;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string url_code = 9;
 * @return {string}
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.getUrlCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialCreatedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialCreatedEvent.prototype.setUrlCode = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
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
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.ProjectTrialConvertedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.ProjectTrialConvertedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.ProjectTrialConvertedEvent}
 */
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.ProjectTrialConvertedEvent;
  return proto.bucketeer.event.domain.ProjectTrialConvertedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.ProjectTrialConvertedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.ProjectTrialConvertedEvent}
 */
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.ProjectTrialConvertedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.ProjectTrialConvertedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.ProjectTrialConvertedEvent} returns this
 */
proto.bucketeer.event.domain.ProjectTrialConvertedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.PrerequisiteAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PrerequisiteAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PrerequisiteAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.PrerequisiteAddedEvent}
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PrerequisiteAddedEvent;
  return proto.bucketeer.event.domain.PrerequisiteAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PrerequisiteAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PrerequisiteAddedEvent}
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.PrerequisiteAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PrerequisiteAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PrerequisiteAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.serializeBinaryToWriter = function(message, writer) {
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
 * optional bucketeer.feature.Prerequisite prerequisite = 1;
 * @return {?proto.bucketeer.feature.Prerequisite}
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.prototype.getPrerequisite = function() {
  return /** @type{?proto.bucketeer.feature.Prerequisite} */ (
    jspb.Message.getWrapperField(this, proto_feature_prerequisite_pb.Prerequisite, 1));
};


/**
 * @param {?proto.bucketeer.feature.Prerequisite|undefined} value
 * @return {!proto.bucketeer.event.domain.PrerequisiteAddedEvent} returns this
*/
proto.bucketeer.event.domain.PrerequisiteAddedEvent.prototype.setPrerequisite = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.PrerequisiteAddedEvent} returns this
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.prototype.clearPrerequisite = function() {
  return this.setPrerequisite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.PrerequisiteAddedEvent.prototype.hasPrerequisite = function() {
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
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent}
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent;
  return proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent}
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
 * optional bucketeer.feature.Prerequisite prerequisite = 1;
 * @return {?proto.bucketeer.feature.Prerequisite}
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.prototype.getPrerequisite = function() {
  return /** @type{?proto.bucketeer.feature.Prerequisite} */ (
    jspb.Message.getWrapperField(this, proto_feature_prerequisite_pb.Prerequisite, 1));
};


/**
 * @param {?proto.bucketeer.feature.Prerequisite|undefined} value
 * @return {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent} returns this
*/
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.prototype.setPrerequisite = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent} returns this
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.prototype.clearPrerequisite = function() {
  return this.setPrerequisite(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.PrerequisiteVariationChangedEvent.prototype.hasPrerequisite = function() {
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
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.PrerequisiteRemovedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.PrerequisiteRemovedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.PrerequisiteRemovedEvent}
 */
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.PrerequisiteRemovedEvent;
  return proto.bucketeer.event.domain.PrerequisiteRemovedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.PrerequisiteRemovedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.PrerequisiteRemovedEvent}
 */
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.PrerequisiteRemovedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.PrerequisiteRemovedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.PrerequisiteRemovedEvent} returns this
 */
proto.bucketeer.event.domain.PrerequisiteRemovedEvent.prototype.setFeatureId = function(value) {
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
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.WebhookCreatedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.WebhookCreatedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    description: jspb.Message.getFieldWithDefault(msg, 3, ""),
    createdAt: jspb.Message.getFieldWithDefault(msg, 4, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 5, 0)
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
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.WebhookCreatedEvent;
  return proto.bucketeer.event.domain.WebhookCreatedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.WebhookCreatedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreatedAt(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
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
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.WebhookCreatedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.WebhookCreatedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional int64 created_at = 4;
 * @return {number}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional int64 updated_at = 5;
 * @return {number}
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.event.domain.WebhookCreatedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookCreatedEvent.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
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
proto.bucketeer.event.domain.WebhookDeletedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.WebhookDeletedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.WebhookDeletedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookDeletedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.WebhookDeletedEvent}
 */
proto.bucketeer.event.domain.WebhookDeletedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.WebhookDeletedEvent;
  return proto.bucketeer.event.domain.WebhookDeletedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.WebhookDeletedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.WebhookDeletedEvent}
 */
proto.bucketeer.event.domain.WebhookDeletedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.WebhookDeletedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.WebhookDeletedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.WebhookDeletedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookDeletedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.WebhookDeletedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookDeletedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookDeletedEvent.prototype.setId = function(value) {
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
proto.bucketeer.event.domain.WebhookNameChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.WebhookNameChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.WebhookNameChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.WebhookNameChangedEvent}
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.WebhookNameChangedEvent;
  return proto.bucketeer.event.domain.WebhookNameChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.WebhookNameChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.WebhookNameChangedEvent}
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.WebhookNameChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.WebhookNameChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.WebhookNameChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.WebhookNameChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookNameChangedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookNameChangedEvent.prototype.setName = function(value) {
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
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.WebhookDescriptionChangedEvent;
  return proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent}
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookDescriptionChangedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookDescriptionChangedEvent.prototype.setDescription = function(value) {
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
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.WebhookClauseAddedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.WebhookClauseAddedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    webhookClause: (f = msg.getWebhookClause()) && proto_autoops_clause_pb.WebhookClause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.WebhookClauseAddedEvent}
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.WebhookClauseAddedEvent;
  return proto.bucketeer.event.domain.WebhookClauseAddedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.WebhookClauseAddedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.WebhookClauseAddedEvent}
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 2:
      var value = new proto_autoops_clause_pb.WebhookClause;
      reader.readMessage(value,proto_autoops_clause_pb.WebhookClause.deserializeBinaryFromReader);
      msg.setWebhookClause(value);
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
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.WebhookClauseAddedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.WebhookClauseAddedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getWebhookClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_autoops_clause_pb.WebhookClause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.WebhookClause webhook_clause = 2;
 * @return {?proto.bucketeer.autoops.WebhookClause}
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.getWebhookClause = function() {
  return /** @type{?proto.bucketeer.autoops.WebhookClause} */ (
    jspb.Message.getWrapperField(this, proto_autoops_clause_pb.WebhookClause, 2));
};


/**
 * @param {?proto.bucketeer.autoops.WebhookClause|undefined} value
 * @return {!proto.bucketeer.event.domain.WebhookClauseAddedEvent} returns this
*/
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.setWebhookClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.WebhookClauseAddedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.clearWebhookClause = function() {
  return this.setWebhookClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.WebhookClauseAddedEvent.prototype.hasWebhookClause = function() {
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
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.event.domain.WebhookClauseChangedEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.event.domain.WebhookClauseChangedEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    clauseId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    webhookClause: (f = msg.getWebhookClause()) && proto_autoops_clause_pb.WebhookClause.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.event.domain.WebhookClauseChangedEvent}
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.event.domain.WebhookClauseChangedEvent;
  return proto.bucketeer.event.domain.WebhookClauseChangedEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.event.domain.WebhookClauseChangedEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.event.domain.WebhookClauseChangedEvent}
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClauseId(value);
      break;
    case 2:
      var value = new proto_autoops_clause_pb.WebhookClause;
      reader.readMessage(value,proto_autoops_clause_pb.WebhookClause.deserializeBinaryFromReader);
      msg.setWebhookClause(value);
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
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.event.domain.WebhookClauseChangedEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.event.domain.WebhookClauseChangedEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClauseId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getWebhookClause();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_autoops_clause_pb.WebhookClause.serializeBinaryToWriter
    );
  }
};


/**
 * optional string clause_id = 1;
 * @return {string}
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.getClauseId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.event.domain.WebhookClauseChangedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.setClauseId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bucketeer.autoops.WebhookClause webhook_clause = 2;
 * @return {?proto.bucketeer.autoops.WebhookClause}
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.getWebhookClause = function() {
  return /** @type{?proto.bucketeer.autoops.WebhookClause} */ (
    jspb.Message.getWrapperField(this, proto_autoops_clause_pb.WebhookClause, 2));
};


/**
 * @param {?proto.bucketeer.autoops.WebhookClause|undefined} value
 * @return {!proto.bucketeer.event.domain.WebhookClauseChangedEvent} returns this
*/
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.setWebhookClause = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.event.domain.WebhookClauseChangedEvent} returns this
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.clearWebhookClause = function() {
  return this.setWebhookClause(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.event.domain.WebhookClauseChangedEvent.prototype.hasWebhookClause = function() {
  return jspb.Message.getField(this, 2) != null;
};


goog.object.extend(exports, proto.bucketeer.event.domain);
