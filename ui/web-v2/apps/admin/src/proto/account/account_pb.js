// source: proto/account/account.proto
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

var proto_environment_environment_pb = require('../../proto/environment/environment_pb.js');
goog.object.extend(proto, proto_environment_environment_pb);
goog.exportSymbol('proto.bucketeer.account.Account', null, global);
goog.exportSymbol('proto.bucketeer.account.Account.Role', null, global);
goog.exportSymbol('proto.bucketeer.account.EnvironmentRole', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.account.Account = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.Account, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.Account.displayName = 'proto.bucketeer.account.Account';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.account.EnvironmentRole = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.EnvironmentRole, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.EnvironmentRole.displayName = 'proto.bucketeer.account.EnvironmentRole';
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
proto.bucketeer.account.Account.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.account.Account.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.account.Account} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.Account.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    email: jspb.Message.getFieldWithDefault(msg, 2, ""),
    name: jspb.Message.getFieldWithDefault(msg, 3, ""),
    role: jspb.Message.getFieldWithDefault(msg, 4, 0),
    disabled: jspb.Message.getBooleanFieldWithDefault(msg, 5, false),
    createdAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
    updatedAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    deleted: jspb.Message.getBooleanFieldWithDefault(msg, 8, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.Account}
 */
proto.bucketeer.account.Account.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.Account;
  return proto.bucketeer.account.Account.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.Account} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.Account}
 */
proto.bucketeer.account.Account.deserializeBinaryFromReader = function(msg, reader) {
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
    case 8:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDeleted(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.account.Account.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.account.Account.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.Account} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.Account.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getDeleted();
  if (f) {
    writer.writeBool(
      8,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.bucketeer.account.Account.Role = {
  VIEWER: 0,
  EDITOR: 1,
  OWNER: 2,
  UNASSIGNED: 99
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.Account.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string email = 2;
 * @return {string}
 */
proto.bucketeer.account.Account.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setEmail = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string name = 3;
 * @return {string}
 */
proto.bucketeer.account.Account.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional Role role = 4;
 * @return {!proto.bucketeer.account.Account.Role}
 */
proto.bucketeer.account.Account.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.Account.Role} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.bucketeer.account.Account.Role} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional bool disabled = 5;
 * @return {boolean}
 */
proto.bucketeer.account.Account.prototype.getDisabled = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setDisabled = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional int64 created_at = 6;
 * @return {number}
 */
proto.bucketeer.account.Account.prototype.getCreatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setCreatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * optional int64 updated_at = 7;
 * @return {number}
 */
proto.bucketeer.account.Account.prototype.getUpdatedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};


/**
 * optional bool deleted = 8;
 * @return {boolean}
 */
proto.bucketeer.account.Account.prototype.getDeleted = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 8, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.account.Account} returns this
 */
proto.bucketeer.account.Account.prototype.setDeleted = function(value) {
  return jspb.Message.setProto3BooleanField(this, 8, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.bucketeer.account.EnvironmentRole.prototype.toObject = function(opt_includeInstance) {
  return proto.bucketeer.account.EnvironmentRole.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.bucketeer.account.EnvironmentRole} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnvironmentRole.toObject = function(includeInstance, msg) {
  var f, obj = {
    environment: (f = msg.getEnvironment()) && proto_environment_environment_pb.Environment.toObject(includeInstance, f),
    role: jspb.Message.getFieldWithDefault(msg, 2, 0),
    trialProject: jspb.Message.getBooleanFieldWithDefault(msg, 3, false),
    trialStartedAt: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.EnvironmentRole}
 */
proto.bucketeer.account.EnvironmentRole.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.EnvironmentRole;
  return proto.bucketeer.account.EnvironmentRole.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.EnvironmentRole} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.EnvironmentRole}
 */
proto.bucketeer.account.EnvironmentRole.deserializeBinaryFromReader = function(msg, reader) {
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
    case 2:
      var value = /** @type {!proto.bucketeer.account.Account.Role} */ (reader.readEnum());
      msg.setRole(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setTrialProject(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setTrialStartedAt(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.account.EnvironmentRole.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.account.EnvironmentRole.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.EnvironmentRole} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnvironmentRole.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnvironment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_environment_environment_pb.Environment.serializeBinaryToWriter
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getTrialProject();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
  f = message.getTrialStartedAt();
  if (f !== 0) {
    writer.writeInt64(
      4,
      f
    );
  }
};


/**
 * optional bucketeer.environment.Environment environment = 1;
 * @return {?proto.bucketeer.environment.Environment}
 */
proto.bucketeer.account.EnvironmentRole.prototype.getEnvironment = function() {
  return /** @type{?proto.bucketeer.environment.Environment} */ (
    jspb.Message.getWrapperField(this, proto_environment_environment_pb.Environment, 1));
};


/**
 * @param {?proto.bucketeer.environment.Environment|undefined} value
 * @return {!proto.bucketeer.account.EnvironmentRole} returns this
*/
proto.bucketeer.account.EnvironmentRole.prototype.setEnvironment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.EnvironmentRole} returns this
 */
proto.bucketeer.account.EnvironmentRole.prototype.clearEnvironment = function() {
  return this.setEnvironment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.EnvironmentRole.prototype.hasEnvironment = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Account.Role role = 2;
 * @return {!proto.bucketeer.account.Account.Role}
 */
proto.bucketeer.account.EnvironmentRole.prototype.getRole = function() {
  return /** @type {!proto.bucketeer.account.Account.Role} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.bucketeer.account.Account.Role} value
 * @return {!proto.bucketeer.account.EnvironmentRole} returns this
 */
proto.bucketeer.account.EnvironmentRole.prototype.setRole = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional bool trial_project = 3;
 * @return {boolean}
 */
proto.bucketeer.account.EnvironmentRole.prototype.getTrialProject = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.bucketeer.account.EnvironmentRole} returns this
 */
proto.bucketeer.account.EnvironmentRole.prototype.setTrialProject = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};


/**
 * optional int64 trial_started_at = 4;
 * @return {number}
 */
proto.bucketeer.account.EnvironmentRole.prototype.getTrialStartedAt = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.bucketeer.account.EnvironmentRole} returns this
 */
proto.bucketeer.account.EnvironmentRole.prototype.setTrialStartedAt = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


goog.object.extend(exports, proto.bucketeer.account);
