// source: proto/experiment/goal.proto
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

var proto_experiment_experiment_pb = require('../../proto/experiment/experiment_pb.js');
goog.object.extend(proto, proto_experiment_experiment_pb);
var proto_autoops_auto_ops_rule_pb = require('../../proto/autoops/auto_ops_rule_pb.js');
goog.object.extend(proto, proto_autoops_auto_ops_rule_pb);
goog.exportSymbol('proto.bucketeer.experiment.Goal', null, global);
goog.exportSymbol(
  'proto.bucketeer.experiment.Goal.AutoOpsRuleReference',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.experiment.Goal.ConnectionType',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.experiment.Goal.ExperimentReference',
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
proto.bucketeer.experiment.Goal = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.experiment.Goal.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.experiment.Goal, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.experiment.Goal.displayName =
    'proto.bucketeer.experiment.Goal';
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
proto.bucketeer.experiment.Goal.ExperimentReference = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.experiment.Goal.ExperimentReference,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.experiment.Goal.ExperimentReference.displayName =
    'proto.bucketeer.experiment.Goal.ExperimentReference';
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
proto.bucketeer.experiment.Goal.AutoOpsRuleReference = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.experiment.Goal.AutoOpsRuleReference,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.experiment.Goal.AutoOpsRuleReference.displayName =
    'proto.bucketeer.experiment.Goal.AutoOpsRuleReference';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.experiment.Goal.repeatedFields_ = [10, 11];

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
  proto.bucketeer.experiment.Goal.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.experiment.Goal.toObject(opt_includeInstance, this);
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.experiment.Goal} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.experiment.Goal.toObject = function (includeInstance, msg) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        name: jspb.Message.getFieldWithDefault(msg, 2, ''),
        description: jspb.Message.getFieldWithDefault(msg, 3, ''),
        deleted: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
        createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
        isInUseStatus: jspb.Message.getBooleanFieldWithDefault(msg, 7, false),
        archived: jspb.Message.getBooleanFieldWithDefault(msg, 8, false),
        connectionType: jspb.Message.getFieldWithDefault(msg, 9, 0),
        experimentsList: jspb.Message.toObjectList(
          msg.getExperimentsList(),
          proto.bucketeer.experiment.Goal.ExperimentReference.toObject,
          includeInstance
        ),
        autoOpsRulesList: jspb.Message.toObjectList(
          msg.getAutoOpsRulesList(),
          proto.bucketeer.experiment.Goal.AutoOpsRuleReference.toObject,
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
 * @return {!proto.bucketeer.experiment.Goal}
 */
proto.bucketeer.experiment.Goal.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.experiment.Goal();
  return proto.bucketeer.experiment.Goal.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.experiment.Goal} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.experiment.Goal}
 */
proto.bucketeer.experiment.Goal.deserializeBinaryFromReader = function (
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
      case 7:
        var value = /** @type {boolean} */ (reader.readBool());
        msg.setIsInUseStatus(value);
        break;
      case 8:
        var value = /** @type {boolean} */ (reader.readBool());
        msg.setArchived(value);
        break;
      case 9:
        var value =
          /** @type {!proto.bucketeer.experiment.Goal.ConnectionType} */ (
            reader.readEnum()
          );
        msg.setConnectionType(value);
        break;
      case 10:
        var value = new proto.bucketeer.experiment.Goal.ExperimentReference();
        reader.readMessage(
          value,
          proto.bucketeer.experiment.Goal.ExperimentReference
            .deserializeBinaryFromReader
        );
        msg.addExperiments(value);
        break;
      case 11:
        var value = new proto.bucketeer.experiment.Goal.AutoOpsRuleReference();
        reader.readMessage(
          value,
          proto.bucketeer.experiment.Goal.AutoOpsRuleReference
            .deserializeBinaryFromReader
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
proto.bucketeer.experiment.Goal.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.experiment.Goal.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.experiment.Goal} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.experiment.Goal.serializeBinaryToWriter = function (
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
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
  f = message.getDeleted();
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
  f = message.getIsInUseStatus();
  if (f) {
    writer.writeBool(7, f);
  }
  f = message.getArchived();
  if (f) {
    writer.writeBool(8, f);
  }
  f = message.getConnectionType();
  if (f !== 0.0) {
    writer.writeEnum(9, f);
  }
  f = message.getExperimentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      10,
      f,
      proto.bucketeer.experiment.Goal.ExperimentReference
        .serializeBinaryToWriter
    );
  }
  f = message.getAutoOpsRulesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      11,
      f,
      proto.bucketeer.experiment.Goal.AutoOpsRuleReference
        .serializeBinaryToWriter
    );
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.experiment.Goal.ConnectionType = {
  UNKNOWN: 0,
  EXPERIMENT: 1,
  OPERATION: 2
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
  proto.bucketeer.experiment.Goal.ExperimentReference.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.experiment.Goal.ExperimentReference.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.experiment.Goal.ExperimentReference} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.experiment.Goal.ExperimentReference.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        name: jspb.Message.getFieldWithDefault(msg, 2, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        status: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference}
 */
proto.bucketeer.experiment.Goal.ExperimentReference.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.experiment.Goal.ExperimentReference();
    return proto.bucketeer.experiment.Goal.ExperimentReference.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.experiment.Goal.ExperimentReference} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference}
 */
proto.bucketeer.experiment.Goal.ExperimentReference.deserializeBinaryFromReader =
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
          msg.setName(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setFeatureId(value);
          break;
        case 4:
          var value =
            /** @type {!proto.bucketeer.experiment.Experiment.Status} */ (
              reader.readEnum()
            );
          msg.setStatus(value);
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
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.experiment.Goal.ExperimentReference.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.experiment.Goal.ExperimentReference} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.experiment.Goal.ExperimentReference.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getName();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getFeatureId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getStatus();
    if (f !== 0.0) {
      writer.writeEnum(4, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference} returns this
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.getName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference} returns this
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.setName =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string feature_id = 3;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.getFeatureId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference} returns this
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.setFeatureId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional Experiment.Status status = 4;
 * @return {!proto.bucketeer.experiment.Experiment.Status}
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.getStatus =
  function () {
    return /** @type {!proto.bucketeer.experiment.Experiment.Status} */ (
      jspb.Message.getFieldWithDefault(this, 4, 0)
    );
  };

/**
 * @param {!proto.bucketeer.experiment.Experiment.Status} value
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference} returns this
 */
proto.bucketeer.experiment.Goal.ExperimentReference.prototype.setStatus =
  function (value) {
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
  proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.experiment.Goal.AutoOpsRuleReference.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.experiment.Goal.AutoOpsRuleReference.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        featureName: jspb.Message.getFieldWithDefault(msg, 3, ''),
        autoOpsStatus: jspb.Message.getFieldWithDefault(msg, 4, 0)
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
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference}
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.experiment.Goal.AutoOpsRuleReference();
    return proto.bucketeer.experiment.Goal.AutoOpsRuleReference.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference}
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.deserializeBinaryFromReader =
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
          msg.setFeatureName(value);
          break;
        case 4:
          var value = /** @type {!proto.bucketeer.autoops.AutoOpsStatus} */ (
            reader.readEnum()
          );
          msg.setAutoOpsStatus(value);
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
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.experiment.Goal.AutoOpsRuleReference.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getFeatureId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getFeatureName();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getAutoOpsStatus();
    if (f !== 0.0) {
      writer.writeEnum(4, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} returns this
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.getFeatureId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} returns this
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.setFeatureId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string feature_name = 3;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.getFeatureName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} returns this
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.setFeatureName =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional bucketeer.autoops.AutoOpsStatus auto_ops_status = 4;
 * @return {!proto.bucketeer.autoops.AutoOpsStatus}
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.getAutoOpsStatus =
  function () {
    return /** @type {!proto.bucketeer.autoops.AutoOpsStatus} */ (
      jspb.Message.getFieldWithDefault(this, 4, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.AutoOpsStatus} value
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference} returns this
 */
proto.bucketeer.experiment.Goal.AutoOpsRuleReference.prototype.setAutoOpsStatus =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 4, value);
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.prototype.getName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setName = function (value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string description = 3;
 * @return {string}
 */
proto.bucketeer.experiment.Goal.prototype.getDescription = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setDescription = function (value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * optional bool deleted = 4;
 * @return {boolean}
 */
proto.bucketeer.experiment.Goal.prototype.getDeleted = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 4, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setDeleted = function (value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};

/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.experiment.Goal.prototype.getCreatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setCreatedAt = function (value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};

/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.experiment.Goal.prototype.getUpdatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setUpdatedAt = function (value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};

/**
 * optional bool is_in_use_status = 7;
 * @return {boolean}
 */
proto.bucketeer.experiment.Goal.prototype.getIsInUseStatus = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 7, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setIsInUseStatus = function (value) {
  return jspb.Message.setProto3BooleanField(this, 7, value);
};

/**
 * optional bool archived = 8;
 * @return {boolean}
 */
proto.bucketeer.experiment.Goal.prototype.getArchived = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 8, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setArchived = function (value) {
  return jspb.Message.setProto3BooleanField(this, 8, value);
};

/**
 * optional ConnectionType connection_type = 9;
 * @return {!proto.bucketeer.experiment.Goal.ConnectionType}
 */
proto.bucketeer.experiment.Goal.prototype.getConnectionType = function () {
  return /** @type {!proto.bucketeer.experiment.Goal.ConnectionType} */ (
    jspb.Message.getFieldWithDefault(this, 9, 0)
  );
};

/**
 * @param {!proto.bucketeer.experiment.Goal.ConnectionType} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setConnectionType = function (value) {
  return jspb.Message.setProto3EnumField(this, 9, value);
};

/**
 * repeated ExperimentReference experiments = 10;
 * @return {!Array<!proto.bucketeer.experiment.Goal.ExperimentReference>}
 */
proto.bucketeer.experiment.Goal.prototype.getExperimentsList = function () {
  return /** @type{!Array<!proto.bucketeer.experiment.Goal.ExperimentReference>} */ (
    jspb.Message.getRepeatedWrapperField(
      this,
      proto.bucketeer.experiment.Goal.ExperimentReference,
      10
    )
  );
};

/**
 * @param {!Array<!proto.bucketeer.experiment.Goal.ExperimentReference>} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setExperimentsList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 10, value);
};

/**
 * @param {!proto.bucketeer.experiment.Goal.ExperimentReference=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.experiment.Goal.ExperimentReference}
 */
proto.bucketeer.experiment.Goal.prototype.addExperiments = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    10,
    opt_value,
    proto.bucketeer.experiment.Goal.ExperimentReference,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.clearExperimentsList = function () {
  return this.setExperimentsList([]);
};

/**
 * repeated AutoOpsRuleReference auto_ops_rules = 11;
 * @return {!Array<!proto.bucketeer.experiment.Goal.AutoOpsRuleReference>}
 */
proto.bucketeer.experiment.Goal.prototype.getAutoOpsRulesList = function () {
  return /** @type{!Array<!proto.bucketeer.experiment.Goal.AutoOpsRuleReference>} */ (
    jspb.Message.getRepeatedWrapperField(
      this,
      proto.bucketeer.experiment.Goal.AutoOpsRuleReference,
      11
    )
  );
};

/**
 * @param {!Array<!proto.bucketeer.experiment.Goal.AutoOpsRuleReference>} value
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.setAutoOpsRulesList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 11, value);
};

/**
 * @param {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.experiment.Goal.AutoOpsRuleReference}
 */
proto.bucketeer.experiment.Goal.prototype.addAutoOpsRules = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    11,
    opt_value,
    proto.bucketeer.experiment.Goal.AutoOpsRuleReference,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.experiment.Goal} returns this
 */
proto.bucketeer.experiment.Goal.prototype.clearAutoOpsRulesList = function () {
  return this.setAutoOpsRulesList([]);
};

goog.object.extend(exports, proto.bucketeer.experiment);
