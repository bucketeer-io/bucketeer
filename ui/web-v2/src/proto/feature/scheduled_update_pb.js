// source: proto/feature/scheduled_update.proto
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

goog.exportSymbol('proto.bucketeer.feature.ScheduledChange', null, global);
goog.exportSymbol(
  'proto.bucketeer.feature.ScheduledChange.ChangeType',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.feature.ScheduledChange.FieldType',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.feature.ScheduledFlagUpdate', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.feature.ScheduledChange = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.feature.ScheduledChange, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ScheduledChange.displayName =
    'proto.bucketeer.feature.ScheduledChange';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.feature.ScheduledFlagUpdate = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.feature.ScheduledFlagUpdate.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.feature.ScheduledFlagUpdate, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.feature.ScheduledFlagUpdate.displayName =
    'proto.bucketeer.feature.ScheduledFlagUpdate';
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
  proto.bucketeer.feature.ScheduledChange.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.feature.ScheduledChange.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.feature.ScheduledChange} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.feature.ScheduledChange.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        changeType: jspb.Message.getFieldWithDefault(msg, 2, 0),
        fieldType: jspb.Message.getFieldWithDefault(msg, 3, 0),
        fieldValue: jspb.Message.getFieldWithDefault(msg, 4, '')
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.feature.ScheduledChange}
 */
proto.bucketeer.feature.ScheduledChange.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ScheduledChange();
  return proto.bucketeer.feature.ScheduledChange.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ScheduledChange} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ScheduledChange}
 */
proto.bucketeer.feature.ScheduledChange.deserializeBinaryFromReader = function (
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
        var value =
          /** @type {!proto.bucketeer.feature.ScheduledChange.ChangeType} */ (
            reader.readEnum()
          );
        msg.setChangeType(value);
        break;
      case 3:
        var value =
          /** @type {!proto.bucketeer.feature.ScheduledChange.FieldType} */ (
            reader.readEnum()
          );
        msg.setFieldType(value);
        break;
      case 4:
        var value = /** @type {string} */ (reader.readString());
        msg.setFieldValue(value);
        break;
      default:
        reader.skipField();
        break;
    }
  }
  return msg;
};

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.feature.ScheduledChange.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.feature.ScheduledChange.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ScheduledChange} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ScheduledChange.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getChangeType();
  if (f !== 0.0) {
    writer.writeEnum(2, f);
  }
  f = message.getFieldType();
  if (f !== 0.0) {
    writer.writeEnum(3, f);
  }
  f = message.getFieldValue();
  if (f.length > 0) {
    writer.writeString(4, f);
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.feature.ScheduledChange.FieldType = {
  UNSPECIFIED: 0,
  PREREQUISITES: 1,
  TARGETS: 2,
  RULES: 3,
  DEFAULT_STRATEGY: 4,
  OFF_VARIATION: 5,
  VARIATIONS: 6
};

/**
 * @enum {number}
 */
proto.bucketeer.feature.ScheduledChange.ChangeType = {
  CHANGE_UNSPECIFIED: 0,
  CHANGE_CREATE: 1,
  CHANGE_UPDATE: 2,
  CHANGE_DELETE: 3
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ScheduledChange.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ScheduledChange} returns this
 */
proto.bucketeer.feature.ScheduledChange.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional ChangeType change_type = 2;
 * @return {!proto.bucketeer.feature.ScheduledChange.ChangeType}
 */
proto.bucketeer.feature.ScheduledChange.prototype.getChangeType = function () {
  return /** @type {!proto.bucketeer.feature.ScheduledChange.ChangeType} */ (
    jspb.Message.getFieldWithDefault(this, 2, 0)
  );
};

/**
 * @param {!proto.bucketeer.feature.ScheduledChange.ChangeType} value
 * @return {!proto.bucketeer.feature.ScheduledChange} returns this
 */
proto.bucketeer.feature.ScheduledChange.prototype.setChangeType = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};

/**
 * optional FieldType field_type = 3;
 * @return {!proto.bucketeer.feature.ScheduledChange.FieldType}
 */
proto.bucketeer.feature.ScheduledChange.prototype.getFieldType = function () {
  return /** @type {!proto.bucketeer.feature.ScheduledChange.FieldType} */ (
    jspb.Message.getFieldWithDefault(this, 3, 0)
  );
};

/**
 * @param {!proto.bucketeer.feature.ScheduledChange.FieldType} value
 * @return {!proto.bucketeer.feature.ScheduledChange} returns this
 */
proto.bucketeer.feature.ScheduledChange.prototype.setFieldType = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};

/**
 * optional string field_value = 4;
 * @return {string}
 */
proto.bucketeer.feature.ScheduledChange.prototype.getFieldValue = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ScheduledChange} returns this
 */
proto.bucketeer.feature.ScheduledChange.prototype.setFieldValue = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 4, value);
};

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.feature.ScheduledFlagUpdate.repeatedFields_ = [7];

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.bucketeer.feature.ScheduledFlagUpdate.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.feature.ScheduledFlagUpdate.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.feature.ScheduledFlagUpdate} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.feature.ScheduledFlagUpdate.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        scheduledAt: jspb.Message.getFieldWithDefault(msg, 4, 0),
        createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
        changesList: jspb.Message.toObjectList(
          msg.getChangesList(),
          proto.bucketeer.feature.ScheduledChange.toObject,
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
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.feature.ScheduledFlagUpdate();
  return proto.bucketeer.feature.ScheduledFlagUpdate.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.feature.ScheduledFlagUpdate} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.deserializeBinaryFromReader =
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
          var value = /** @type {string} */ (reader.readString());
          msg.setFeatureId(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setScheduledAt(value);
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
          var value = new proto.bucketeer.feature.ScheduledChange();
          reader.readMessage(
            value,
            proto.bucketeer.feature.ScheduledChange.deserializeBinaryFromReader
          );
          msg.addChanges(value);
          break;
        default:
          reader.skipField();
          break;
      }
    }
    return msg;
  };

/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.feature.ScheduledFlagUpdate.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.feature.ScheduledFlagUpdate} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.feature.ScheduledFlagUpdate.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
  f = message.getScheduledAt();
  if (f !== 0) {
    writer.writeInt64(4, f);
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(5, f);
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(6, f);
  }
  f = message.getChangesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      7,
      f,
      proto.bucketeer.feature.ScheduledChange.serializeBinaryToWriter
    );
  }
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getFeatureId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setFeatureId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional int64 scheduled_at = 4;
 * @return {number}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getScheduledAt =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setScheduledAt =
  function (value) {
    return jspb.Message.setProto3IntField(this, 4, value);
  };

/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getCreatedAt =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setCreatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 5, value);
};

/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getUpdatedAt =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setUpdatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 6, value);
};

/**
 * repeated ScheduledChange changes = 7;
 * @return {!Array<!proto.bucketeer.feature.ScheduledChange>}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.getChangesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.feature.ScheduledChange>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto.bucketeer.feature.ScheduledChange,
        7
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.feature.ScheduledChange>} value
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.setChangesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 7, value);
  };

/**
 * @param {!proto.bucketeer.feature.ScheduledChange=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.ScheduledChange}
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.addChanges = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    7,
    opt_value,
    proto.bucketeer.feature.ScheduledChange,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.feature.ScheduledFlagUpdate} returns this
 */
proto.bucketeer.feature.ScheduledFlagUpdate.prototype.clearChangesList =
  function () {
    return this.setChangesList([]);
  };

goog.object.extend(exports, proto.bucketeer.feature);
