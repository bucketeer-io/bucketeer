// source: proto/environment/project.proto
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

goog.exportSymbol('proto.bucketeer.environment.Project', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.environment.Project = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.environment.Project, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.environment.Project.displayName = 'proto.bucketeer.environment.Project';
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
proto.bucketeer.environment.Project.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.environment.Project.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.environment.Project} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.Project.toObject = function(includeInstance, msg) {
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
 * @return {!proto.bucketeer.environment.Project}
 */
proto.bucketeer.environment.Project.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.environment.Project;
  return proto.bucketeer.environment.Project.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.environment.Project} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.environment.Project}
 */
proto.bucketeer.environment.Project.deserializeBinaryFromReader = function(msg, reader) {
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
proto.bucketeer.environment.Project.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.environment.Project.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.environment.Project} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.environment.Project.serializeBinaryToWriter = function(message, writer) {
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
proto.bucketeer.environment.Project.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string description = 2;
 * @return {string}
 */
proto.bucketeer.environment.Project.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setDescription = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool disabled = 3;
 * @return {boolean}
 */
proto.bucketeer.environment.Project.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};


/**
 * optional bool trial = 4;
 * @return {boolean}
 */
proto.bucketeer.environment.Project.prototype.getTrial = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setTrial = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional string creator_email = 5;
 * @return {string}
 */
proto.bucketeer.environment.Project.prototype.getCreatorEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setCreatorEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.environment.Project.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.environment.Project.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional string name = 8;
 * @return {string}
 */
proto.bucketeer.environment.Project.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string url_code = 9;
 * @return {string}
 */
proto.bucketeer.environment.Project.prototype.getUrlCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.environment.Project} returns this
 */
proto.bucketeer.environment.Project.prototype.setUrlCode = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};


goog.object.extend(exports, proto.bucketeer.environment);
