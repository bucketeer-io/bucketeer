// source: proto/account/api_key.proto
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

var proto_environment_environment_pb = require('../../proto/environment/environment_pb.js');
goog.object.extend(proto, proto_environment_environment_pb);
goog.exportSymbol('proto.bucketeer.account.APIKey', null, global);
goog.exportSymbol('proto.bucketeer.account.APIKey.Role', null, global);
goog.exportSymbol('proto.bucketeer.account.EnvironmentAPIKey', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.account.APIKey = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.APIKey, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.APIKey.displayName = 'proto.bucketeer.account.APIKey';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.account.EnvironmentAPIKey = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.EnvironmentAPIKey, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.EnvironmentAPIKey.displayName =
    'proto.bucketeer.account.EnvironmentAPIKey';
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
  proto.bucketeer.account.APIKey.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.APIKey.toObject(opt_includeInstance, this);
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.APIKey} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.APIKey.toObject = function (includeInstance, msg) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        name: jspb.Message.getFieldWithDefault(msg, 2, ''),
        role: jspb.Message.getFieldWithDefault(msg, 3, 0),
        disabled: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
        createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
        maintainer: jspb.Message.getFieldWithDefault(msg, 7, ''),
        apiKey: jspb.Message.getFieldWithDefault(msg, 8, ''),
        description: jspb.Message.getFieldWithDefault(msg, 9, '')
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.APIKey}
 */
proto.bucketeer.account.APIKey.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.APIKey();
  return proto.bucketeer.account.APIKey.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.APIKey} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.APIKey}
 */
proto.bucketeer.account.APIKey.deserializeBinaryFromReader = function (
  msg,
  reader
) {
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
        var value = /** @type {!proto.bucketeer.account.APIKey.Role} */ (
          reader.readEnum()
        );
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
      case 7:
        var value = /** @type {string} */ (reader.readString());
        msg.setMaintainer(value);
        break;
      case 8:
        var value = /** @type {string} */ (reader.readString());
        msg.setApiKey(value);
        break;
      case 9:
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
proto.bucketeer.account.APIKey.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.account.APIKey.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.APIKey} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.APIKey.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(3, f);
  }
  f = message.getDisabled();
  if (f) {
    writer.writeBool(4, f);
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(5, f);
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(6, f);
  }
  f = message.getMaintainer();
  if (f.length > 0) {
    writer.writeString(7, f);
  }
  f = message.getApiKey();
  if (f.length > 0) {
    writer.writeString(8, f);
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(9, f);
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.account.APIKey.Role = {
  UNKNOWN: 0,
  SDK_CLIENT: 1,
  SDK_SERVER: 2,
  PUBLIC_API_READ_ONLY: 3,
  PUBLIC_API_WRITE: 4,
  PUBLIC_API_ADMIN: 5
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.APIKey.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.account.APIKey.prototype.getName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setName = function (value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional Role role = 3;
 * @return {!proto.bucketeer.account.APIKey.Role}
 */
proto.bucketeer.account.APIKey.prototype.getRole = function () {
  return /** @type {!proto.bucketeer.account.APIKey.Role} */ (
    jspb.Message.getFieldWithDefault(this, 3, 0)
  );
};

/**
 * @param {!proto.bucketeer.account.APIKey.Role} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setRole = function (value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};

/**
 * optional bool disabled = 4;
 * @return {boolean}
 */
proto.bucketeer.account.APIKey.prototype.getDisabled = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 4, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setDisabled = function (value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};

/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.account.APIKey.prototype.getCreatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setCreatedAt = function (value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};

/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.account.APIKey.prototype.getUpdatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setUpdatedAt = function (value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};

/**
 * optional string maintainer = 7;
 * @return {string}
 */
proto.bucketeer.account.APIKey.prototype.getMaintainer = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setMaintainer = function (value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};

/**
 * optional string api_key = 8;
 * @return {string}
 */
proto.bucketeer.account.APIKey.prototype.getApiKey = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setApiKey = function (value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};

/**
 * optional string description = 9;
 * @return {string}
 */
proto.bucketeer.account.APIKey.prototype.getDescription = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.APIKey} returns this
 */
proto.bucketeer.account.APIKey.prototype.setDescription = function (value) {
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
  proto.bucketeer.account.EnvironmentAPIKey.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.EnvironmentAPIKey.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.EnvironmentAPIKey} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.EnvironmentAPIKey.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        apiKey:
          (f = msg.getApiKey()) &&
          proto.bucketeer.account.APIKey.toObject(includeInstance, f),
        environmentDisabled: jspb.Message.getBooleanFieldWithDefault(
          msg,
          3,
          false
        ),
        projectId: jspb.Message.getFieldWithDefault(msg, 4, ''),
        environment:
          (f = msg.getEnvironment()) &&
          proto_environment_environment_pb.EnvironmentV2.toObject(
            includeInstance,
            f
          ),
        projectUrlCode: jspb.Message.getFieldWithDefault(msg, 6, '')
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.EnvironmentAPIKey}
 */
proto.bucketeer.account.EnvironmentAPIKey.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.EnvironmentAPIKey();
  return proto.bucketeer.account.EnvironmentAPIKey.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.EnvironmentAPIKey} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.EnvironmentAPIKey}
 */
proto.bucketeer.account.EnvironmentAPIKey.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = new proto.bucketeer.account.APIKey();
          reader.readMessage(
            value,
            proto.bucketeer.account.APIKey.deserializeBinaryFromReader
          );
          msg.setApiKey(value);
          break;
        case 3:
          var value = /** @type {boolean} */ (reader.readBool());
          msg.setEnvironmentDisabled(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setProjectId(value);
          break;
        case 5:
          var value = new proto_environment_environment_pb.EnvironmentV2();
          reader.readMessage(
            value,
            proto_environment_environment_pb.EnvironmentV2
              .deserializeBinaryFromReader
          );
          msg.setEnvironment(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setProjectUrlCode(value);
          break;
        default:
          reader.skipField();
          break;
      }
    }
    return msg;
  };

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.EnvironmentAPIKey.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.EnvironmentAPIKey} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnvironmentAPIKey.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getApiKey();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.bucketeer.account.APIKey.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentDisabled();
  if (f) {
    writer.writeBool(3, f);
  }
  f = message.getProjectId();
  if (f.length > 0) {
    writer.writeString(4, f);
  }
  f = message.getEnvironment();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto_environment_environment_pb.EnvironmentV2.serializeBinaryToWriter
    );
  }
  f = message.getProjectUrlCode();
  if (f.length > 0) {
    writer.writeString(6, f);
  }
};

/**
 * optional APIKey api_key = 2;
 * @return {?proto.bucketeer.account.APIKey}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.getApiKey = function () {
  return /** @type{?proto.bucketeer.account.APIKey} */ (
    jspb.Message.getWrapperField(this, proto.bucketeer.account.APIKey, 2)
  );
};

/**
 * @param {?proto.bucketeer.account.APIKey|undefined} value
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.setApiKey = function (
  value
) {
  return jspb.Message.setWrapperField(this, 2, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.clearApiKey = function () {
  return this.setApiKey(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.hasApiKey = function () {
  return jspb.Message.getField(this, 2) != null;
};

/**
 * optional bool environment_disabled = 3;
 * @return {boolean}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.getEnvironmentDisabled =
  function () {
    return /** @type {boolean} */ (
      jspb.Message.getBooleanFieldWithDefault(this, 3, false)
    );
  };

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.setEnvironmentDisabled =
  function (value) {
    return jspb.Message.setProto3BooleanField(this, 3, value);
  };

/**
 * optional string project_id = 4;
 * @return {string}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.getProjectId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.setProjectId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 4, value);
};

/**
 * optional bucketeer.environment.EnvironmentV2 environment = 5;
 * @return {?proto.bucketeer.environment.EnvironmentV2}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.getEnvironment =
  function () {
    return /** @type{?proto.bucketeer.environment.EnvironmentV2} */ (
      jspb.Message.getWrapperField(
        this,
        proto_environment_environment_pb.EnvironmentV2,
        5
      )
    );
  };

/**
 * @param {?proto.bucketeer.environment.EnvironmentV2|undefined} value
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.setEnvironment = function (
  value
) {
  return jspb.Message.setWrapperField(this, 5, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.clearEnvironment =
  function () {
    return this.setEnvironment(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.hasEnvironment =
  function () {
    return jspb.Message.getField(this, 5) != null;
  };

/**
 * optional string project_url_code = 6;
 * @return {string}
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.getProjectUrlCode =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.EnvironmentAPIKey} returns this
 */
proto.bucketeer.account.EnvironmentAPIKey.prototype.setProjectUrlCode =
  function (value) {
    return jspb.Message.setProto3StringField(this, 6, value);
  };

goog.object.extend(exports, proto.bucketeer.account);
