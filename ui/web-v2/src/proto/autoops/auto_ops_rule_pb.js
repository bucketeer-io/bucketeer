// source: proto/autoops/auto_ops_rule.proto
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

var proto_autoops_clause_pb = require('../../proto/autoops/clause_pb.js');
goog.object.extend(proto, proto_autoops_clause_pb);
goog.exportSymbol('proto.bucketeer.autoops.AutoOpsRule', null, global);
goog.exportSymbol('proto.bucketeer.autoops.AutoOpsRules', null, global);
goog.exportSymbol('proto.bucketeer.autoops.AutoOpsStatus', null, global);
goog.exportSymbol('proto.bucketeer.autoops.OpsType', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.autoops.AutoOpsRule = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.AutoOpsRule.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.AutoOpsRule, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.AutoOpsRule.displayName =
    'proto.bucketeer.autoops.AutoOpsRule';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.bucketeer.autoops.AutoOpsRules = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.autoops.AutoOpsRules.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.autoops.AutoOpsRules, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.AutoOpsRules.displayName =
    'proto.bucketeer.autoops.AutoOpsRules';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.AutoOpsRule.repeatedFields_ = [4];

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.bucketeer.autoops.AutoOpsRule.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.AutoOpsRule.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.AutoOpsRule} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.AutoOpsRule.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        opsType: jspb.Message.getFieldWithDefault(msg, 3, 0),
        clausesList: jspb.Message.toObjectList(
          msg.getClausesList(),
          proto_autoops_clause_pb.Clause.toObject,
          includeInstance
        ),
        createdAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 8, 0),
        deleted: jspb.Message.getBooleanFieldWithDefault(msg, 9, false),
        autoOpsStatus: jspb.Message.getFieldWithDefault(msg, 10, 0),
        featureName: jspb.Message.getFieldWithDefault(msg, 11, '')
      };

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.autoops.AutoOpsRule}
 */
proto.bucketeer.autoops.AutoOpsRule.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.AutoOpsRule();
  return proto.bucketeer.autoops.AutoOpsRule.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.AutoOpsRule} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.AutoOpsRule}
 */
proto.bucketeer.autoops.AutoOpsRule.deserializeBinaryFromReader = function (
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
        msg.setFeatureId(value);
        break;
      case 3:
        var value = /** @type {!proto.bucketeer.autoops.OpsType} */ (
          reader.readEnum()
        );
        msg.setOpsType(value);
        break;
      case 4:
        var value = new proto_autoops_clause_pb.Clause();
        reader.readMessage(
          value,
          proto_autoops_clause_pb.Clause.deserializeBinaryFromReader
        );
        msg.addClauses(value);
        break;
      case 7:
        var value = /** @type {number} */ (reader.readInt64());
        msg.setCreatedAt(value);
        break;
      case 8:
        var value = /** @type {number} */ (reader.readInt64());
        msg.setUpdatedAt(value);
        break;
      case 9:
        var value = /** @type {boolean} */ (reader.readBool());
        msg.setDeleted(value);
        break;
      case 10:
        var value = /** @type {!proto.bucketeer.autoops.AutoOpsStatus} */ (
          reader.readEnum()
        );
        msg.setAutoOpsStatus(value);
        break;
      case 11:
        var value = /** @type {string} */ (reader.readString());
        msg.setFeatureName(value);
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
proto.bucketeer.autoops.AutoOpsRule.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.autoops.AutoOpsRule.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.AutoOpsRule} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.AutoOpsRule.serializeBinaryToWriter = function (
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
  f = message.getOpsType();
  if (f !== 0.0) {
    writer.writeEnum(3, f);
  }
  f = message.getClausesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto_autoops_clause_pb.Clause.serializeBinaryToWriter
    );
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(7, f);
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(8, f);
  }
  f = message.getDeleted();
  if (f) {
    writer.writeBool(9, f);
  }
  f = message.getAutoOpsStatus();
  if (f !== 0.0) {
    writer.writeEnum(10, f);
  }
  f = message.getFeatureName();
  if (f.length > 0) {
    writer.writeString(11, f);
  }
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getFeatureId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setFeatureId = function (value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional OpsType ops_type = 3;
 * @return {!proto.bucketeer.autoops.OpsType}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getOpsType = function () {
  return /** @type {!proto.bucketeer.autoops.OpsType} */ (
    jspb.Message.getFieldWithDefault(this, 3, 0)
  );
};

/**
 * @param {!proto.bucketeer.autoops.OpsType} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setOpsType = function (value) {
  return jspb.Message.setProto3EnumField(this, 3, value);
};

/**
 * repeated Clause clauses = 4;
 * @return {!Array<!proto.bucketeer.autoops.Clause>}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getClausesList = function () {
  return /** @type{!Array<!proto.bucketeer.autoops.Clause>} */ (
    jspb.Message.getRepeatedWrapperField(
      this,
      proto_autoops_clause_pb.Clause,
      4
    )
  );
};

/**
 * @param {!Array<!proto.bucketeer.autoops.Clause>} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setClausesList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};

/**
 * @param {!proto.bucketeer.autoops.Clause=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.Clause}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.addClauses = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    4,
    opt_value,
    proto.bucketeer.autoops.Clause,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.clearClausesList = function () {
  return this.setClausesList([]);
};

/**
 * optional int64 created_at = 7;
 * @return {number}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getCreatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setCreatedAt = function (value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};

/**
 * optional int64 updated_at = 8;
 * @return {number}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getUpdatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setUpdatedAt = function (value) {
  return jspb.Message.setProto3IntField(this, 8, value);
};

/**
 * optional bool deleted = 9;
 * @return {boolean}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getDeleted = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 9, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setDeleted = function (value) {
  return jspb.Message.setProto3BooleanField(this, 9, value);
};

/**
 * optional AutoOpsStatus auto_ops_status = 10;
 * @return {!proto.bucketeer.autoops.AutoOpsStatus}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getAutoOpsStatus = function () {
  return /** @type {!proto.bucketeer.autoops.AutoOpsStatus} */ (
    jspb.Message.getFieldWithDefault(this, 10, 0)
  );
};

/**
 * @param {!proto.bucketeer.autoops.AutoOpsStatus} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setAutoOpsStatus = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 10, value);
};

/**
 * optional string feature_name = 11;
 * @return {string}
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.getFeatureName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 11, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.AutoOpsRule} returns this
 */
proto.bucketeer.autoops.AutoOpsRule.prototype.setFeatureName = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 11, value);
};

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.autoops.AutoOpsRules.repeatedFields_ = [1];

if (jspb.Message.GENERATE_TO_OBJECT) {
  /**
   * Creates an object representation of this proto.
   * Field names that are reserved in JavaScript and will be renamed to pb_name.
   * Optional fields that are not set will be set to undefined.
   * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
   * For the list of reserved names please see:
   *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
   * @param {boolean=} opt_includeInstance Deprecated. whether to include the
   *     JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @return {!Object}
   */
  proto.bucketeer.autoops.AutoOpsRules.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.AutoOpsRules.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.AutoOpsRules} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.AutoOpsRules.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        autoOpsRulesList: jspb.Message.toObjectList(
          msg.getAutoOpsRulesList(),
          proto.bucketeer.autoops.AutoOpsRule.toObject,
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
 * @return {!proto.bucketeer.autoops.AutoOpsRules}
 */
proto.bucketeer.autoops.AutoOpsRules.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.AutoOpsRules();
  return proto.bucketeer.autoops.AutoOpsRules.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.AutoOpsRules} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.AutoOpsRules}
 */
proto.bucketeer.autoops.AutoOpsRules.deserializeBinaryFromReader = function (
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
        var value = new proto.bucketeer.autoops.AutoOpsRule();
        reader.readMessage(
          value,
          proto.bucketeer.autoops.AutoOpsRule.deserializeBinaryFromReader
        );
        msg.addAutoOpsRules(value);
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
proto.bucketeer.autoops.AutoOpsRules.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.autoops.AutoOpsRules.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.AutoOpsRules} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.AutoOpsRules.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getAutoOpsRulesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.bucketeer.autoops.AutoOpsRule.serializeBinaryToWriter
    );
  }
};

/**
 * repeated AutoOpsRule auto_ops_rules = 1;
 * @return {!Array<!proto.bucketeer.autoops.AutoOpsRule>}
 */
proto.bucketeer.autoops.AutoOpsRules.prototype.getAutoOpsRulesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.autoops.AutoOpsRule>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto.bucketeer.autoops.AutoOpsRule,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.autoops.AutoOpsRule>} value
 * @return {!proto.bucketeer.autoops.AutoOpsRules} returns this
 */
proto.bucketeer.autoops.AutoOpsRules.prototype.setAutoOpsRulesList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};

/**
 * @param {!proto.bucketeer.autoops.AutoOpsRule=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.autoops.AutoOpsRule}
 */
proto.bucketeer.autoops.AutoOpsRules.prototype.addAutoOpsRules = function (
  opt_value,
  opt_index
) {
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
 * @return {!proto.bucketeer.autoops.AutoOpsRules} returns this
 */
proto.bucketeer.autoops.AutoOpsRules.prototype.clearAutoOpsRulesList =
  function () {
    return this.setAutoOpsRulesList([]);
  };

/**
 * @enum {number}
 */
proto.bucketeer.autoops.OpsType = {
  TYPE_UNKNOWN: 0,
  SCHEDULE: 2,
  EVENT_RATE: 3
};

/**
 * @enum {number}
 */
proto.bucketeer.autoops.AutoOpsStatus = {
  WAITING: 0,
  RUNNING: 1,
  FINISHED: 2,
  STOPPED: 3
};

goog.object.extend(exports, proto.bucketeer.autoops);
