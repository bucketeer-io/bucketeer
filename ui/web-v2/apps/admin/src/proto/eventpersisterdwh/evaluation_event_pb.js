// source: proto/eventpersisterdwh/evaluation_event.proto
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
var global = Function('return this')();

goog.exportSymbol('proto.bucketeer.eventcounter.EvaluationEvent', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.eventcounter.EvaluationEvent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.eventcounter.EvaluationEvent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.eventcounter.EvaluationEvent.displayName = 'proto.bucketeer.eventcounter.EvaluationEvent';
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
proto.bucketeer.eventcounter.EvaluationEvent.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.eventcounter.EvaluationEvent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.eventcounter.EvaluationEvent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.eventcounter.EvaluationEvent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    featureId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    featureVersion: jspb.Message.getFieldWithDefault(msg, 3, 0),
    userData: jspb.Message.getFieldWithDefault(msg, 4, ""),
    userId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    variationId: jspb.Message.getFieldWithDefault(msg, 6, ""),
    reason: jspb.Message.getFieldWithDefault(msg, 7, ""),
    tag: jspb.Message.getFieldWithDefault(msg, 8, ""),
    sourceId: jspb.Message.getFieldWithDefault(msg, 9, ""),
    environmentNamespace: jspb.Message.getFieldWithDefault(msg, 10, ""),
    timestamp: jspb.Message.getFieldWithDefault(msg, 11, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent}
 */
proto.bucketeer.eventcounter.EvaluationEvent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.eventcounter.EvaluationEvent;
  return proto.bucketeer.eventcounter.EvaluationEvent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.eventcounter.EvaluationEvent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent}
 */
proto.bucketeer.eventcounter.EvaluationEvent.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setUserData(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setVariationId(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setReason(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setTag(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setSourceId(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setEnvironmentNamespace(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTimestamp(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.eventcounter.EvaluationEvent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.eventcounter.EvaluationEvent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.eventcounter.EvaluationEvent.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getUserData();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getVariationId();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getReason();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getTag();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getSourceId();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getEnvironmentNamespace();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getTimestamp();
  if (f !== 0) {
    writer.writeInt64(
      11,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getFeatureId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setFeatureId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional int32 feature_version = 3;
 * @return {number}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getFeatureVersion = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setFeatureVersion = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional string user_data = 4;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getUserData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setUserData = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string user_id = 5;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setUserId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string variation_id = 6;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getVariationId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setVariationId = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string reason = 7;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getReason = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setReason = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string tag = 8;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getTag = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setTag = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string source_id = 9;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getSourceId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setSourceId = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * optional string environment_namespace = 10;
 * @return {string}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getEnvironmentNamespace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setEnvironmentNamespace = function(value) {
  return jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional int64 timestamp = 11;
 * @return {number}
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.getTimestamp = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.eventcounter.EvaluationEvent} returns this
 */
proto.bucketeer.eventcounter.EvaluationEvent.prototype.setTimestamp = function(value) {
  return jspb.Message.setProto3IntField(this, 11, value);
};


goog.object.extend(exports, proto.bucketeer.eventcounter);
