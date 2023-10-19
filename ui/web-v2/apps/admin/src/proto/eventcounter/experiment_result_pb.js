// source: proto/eventcounter/experiment_result.proto
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

var proto_eventcounter_goal_result_pb = require('../../proto/eventcounter/goal_result_pb.js');
goog.object.extend(proto, proto_eventcounter_goal_result_pb);
goog.exportSymbol('proto.bucketeer.eventcounter.ExperimentResult', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.eventcounter.ExperimentResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.bucketeer.eventcounter.ExperimentResult.repeatedFields_, null);
};
goog.inherits(proto.bucketeer.eventcounter.ExperimentResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.eventcounter.ExperimentResult.displayName = 'proto.bucketeer.eventcounter.ExperimentResult';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.eventcounter.ExperimentResult.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.eventcounter.ExperimentResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.eventcounter.ExperimentResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.eventcounter.ExperimentResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    experimentId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 3, 0),
    goalResultsList: jspb.Message.toObjectList(msg.getGoalResultsList(),
    proto_eventcounter_goal_result_pb.GoalResult.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.eventcounter.ExperimentResult}
 */
proto.bucketeer.eventcounter.ExperimentResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.eventcounter.ExperimentResult;
  return proto.bucketeer.eventcounter.ExperimentResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.eventcounter.ExperimentResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.eventcounter.ExperimentResult}
 */
proto.bucketeer.eventcounter.ExperimentResult.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setExperimentId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUpdatedAt(value);
      break;
    case 4:
      var value = new proto_eventcounter_goal_result_pb.GoalResult;
      reader.readMessage(value,proto_eventcounter_goal_result_pb.GoalResult.deserializeBinaryFromReader);
      msg.addGoalResults(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.eventcounter.ExperimentResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.eventcounter.ExperimentResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.eventcounter.ExperimentResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getExperimentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
  f = message.getGoalResultsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto_eventcounter_goal_result_pb.GoalResult.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.ExperimentResult} returns this
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string experiment_id = 2;
 * @return {string}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.getExperimentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.ExperimentResult} returns this
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.setExperimentId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int64 updated_at = 3;
 * @return {number}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.eventcounter.ExperimentResult} returns this
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * repeated GoalResult goal_results = 4;
 * @return {!Array<!proto.bucketeer.eventcounter.GoalResult>}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.getGoalResultsList = function() {
  return /** @type{!Array<!proto.bucketeer.eventcounter.GoalResult>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto_eventcounter_goal_result_pb.GoalResult, 4));
};


/**
 * @param {!Array<!proto.bucketeer.eventcounter.GoalResult>} value
 * @return {!proto.bucketeer.eventcounter.ExperimentResult} returns this
*/
proto.bucketeer.eventcounter.ExperimentResult.prototype.setGoalResultsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.bucketeer.eventcounter.GoalResult=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.eventcounter.GoalResult}
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.addGoalResults = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.bucketeer.eventcounter.GoalResult, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.eventcounter.ExperimentResult} returns this
 */
proto.bucketeer.eventcounter.ExperimentResult.prototype.clearGoalResultsList = function() {
  return this.setGoalResultsList([]);
};


goog.object.extend(exports, proto.bucketeer.eventcounter);
