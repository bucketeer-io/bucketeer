// source: proto/experiment/experiment.proto
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

var proto_feature_variation_pb = require('../../proto/feature/variation_pb.js');
goog.object.extend(proto, proto_feature_variation_pb);
goog.exportSymbol('proto.bucketeer.experiment.Experiment', null, global);
goog.exportSymbol(
  'proto.bucketeer.experiment.Experiment.GoalReference',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.experiment.Experiment.Status', null, global);
goog.exportSymbol('proto.bucketeer.experiment.Experiments', null, global);
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
proto.bucketeer.experiment.Experiment = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.experiment.Experiment.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.experiment.Experiment, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.experiment.Experiment.displayName =
    'proto.bucketeer.experiment.Experiment';
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
proto.bucketeer.experiment.Experiment.GoalReference = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.experiment.Experiment.GoalReference,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.experiment.Experiment.GoalReference.displayName =
    'proto.bucketeer.experiment.Experiment.GoalReference';
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
proto.bucketeer.experiment.Experiments = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.experiment.Experiments.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.experiment.Experiments, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.experiment.Experiments.displayName =
    'proto.bucketeer.experiment.Experiments';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.experiment.Experiment.repeatedFields_ = [5, 13, 21];

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
  proto.bucketeer.experiment.Experiment.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.experiment.Experiment.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.experiment.Experiment} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.experiment.Experiment.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        goalId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        featureVersion: jspb.Message.getFieldWithDefault(msg, 4, 0),
        variationsList: jspb.Message.toObjectList(
          msg.getVariationsList(),
          proto_feature_variation_pb.Variation.toObject,
          includeInstance
        ),
        startAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
        stopAt: jspb.Message.getFieldWithDefault(msg, 7, 0),
        stopped: jspb.Message.getBooleanFieldWithDefault(msg, 8, false),
        stoppedAt: jspb.Message.getFieldWithDefault(msg, 9, '0'),
        createdAt: jspb.Message.getFieldWithDefault(msg, 10, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 11, 0),
        deleted: jspb.Message.getBooleanFieldWithDefault(msg, 12, false),
        goalIdsList:
          (f = jspb.Message.getRepeatedField(msg, 13)) == null ? undefined : f,
        name: jspb.Message.getFieldWithDefault(msg, 14, ''),
        description: jspb.Message.getFieldWithDefault(msg, 15, ''),
        baseVariationId: jspb.Message.getFieldWithDefault(msg, 16, ''),
        status: jspb.Message.getFieldWithDefault(msg, 18, 0),
        maintainer: jspb.Message.getFieldWithDefault(msg, 19, ''),
        archived: jspb.Message.getBooleanFieldWithDefault(msg, 20, false),
        goalsList: jspb.Message.toObjectList(
          msg.getGoalsList(),
          proto.bucketeer.experiment.Experiment.GoalReference.toObject,
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
 * @return {!proto.bucketeer.experiment.Experiment}
 */
proto.bucketeer.experiment.Experiment.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.experiment.Experiment();
  return proto.bucketeer.experiment.Experiment.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.experiment.Experiment} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.experiment.Experiment}
 */
proto.bucketeer.experiment.Experiment.deserializeBinaryFromReader = function (
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
        msg.setGoalId(value);
        break;
      case 3:
        var value = /** @type {string} */ (reader.readString());
        msg.setFeatureId(value);
        break;
      case 4:
        var value = /** @type {number} */ (reader.readInt32());
        msg.setFeatureVersion(value);
        break;
      case 5:
        var value = new proto_feature_variation_pb.Variation();
        reader.readMessage(
          value,
          proto_feature_variation_pb.Variation.deserializeBinaryFromReader
        );
        msg.addVariations(value);
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
        var value = /** @type {string} */ (reader.readInt64String());
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
        var value = /** @type {boolean} */ (reader.readBool());
        msg.setDeleted(value);
        break;
      case 13:
        var value = /** @type {string} */ (reader.readString());
        msg.addGoalIds(value);
        break;
      case 14:
        var value = /** @type {string} */ (reader.readString());
        msg.setName(value);
        break;
      case 15:
        var value = /** @type {string} */ (reader.readString());
        msg.setDescription(value);
        break;
      case 16:
        var value = /** @type {string} */ (reader.readString());
        msg.setBaseVariationId(value);
        break;
      case 18:
        var value =
          /** @type {!proto.bucketeer.experiment.Experiment.Status} */ (
            reader.readEnum()
          );
        msg.setStatus(value);
        break;
      case 19:
        var value = /** @type {string} */ (reader.readString());
        msg.setMaintainer(value);
        break;
      case 20:
        var value = /** @type {boolean} */ (reader.readBool());
        msg.setArchived(value);
        break;
      case 21:
        var value = new proto.bucketeer.experiment.Experiment.GoalReference();
        reader.readMessage(
          value,
          proto.bucketeer.experiment.Experiment.GoalReference
            .deserializeBinaryFromReader
        );
        msg.addGoals(value);
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
proto.bucketeer.experiment.Experiment.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.experiment.Experiment.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.experiment.Experiment} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.experiment.Experiment.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getGoalId();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getFeatureId();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
  f = message.getFeatureVersion();
  if (f !== 0) {
    writer.writeInt32(4, f);
  }
  f = message.getVariationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto_feature_variation_pb.Variation.serializeBinaryToWriter
    );
  }
  f = message.getStartAt();
  if (f !== 0) {
    writer.writeInt64(6, f);
  }
  f = message.getStopAt();
  if (f !== 0) {
    writer.writeInt64(7, f);
  }
  f = message.getStopped();
  if (f) {
    writer.writeBool(8, f);
  }
  f = message.getStoppedAt();
  if (parseInt(f, 10) !== 0) {
    writer.writeInt64String(9, f);
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(10, f);
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(11, f);
  }
  f = message.getDeleted();
  if (f) {
    writer.writeBool(12, f);
  }
  f = message.getGoalIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(13, f);
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(14, f);
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(15, f);
  }
  f = message.getBaseVariationId();
  if (f.length > 0) {
    writer.writeString(16, f);
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(18, f);
  }
  f = message.getMaintainer();
  if (f.length > 0) {
    writer.writeString(19, f);
  }
  f = message.getArchived();
  if (f) {
    writer.writeBool(20, f);
  }
  f = message.getGoalsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      21,
      f,
      proto.bucketeer.experiment.Experiment.GoalReference
        .serializeBinaryToWriter
    );
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.experiment.Experiment.Status = {
  WAITING: 0,
  RUNNING: 1,
  STOPPED: 2,
  FORCE_STOPPED: 3
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
  proto.bucketeer.experiment.Experiment.GoalReference.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.experiment.Experiment.GoalReference.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.experiment.Experiment.GoalReference} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.experiment.Experiment.GoalReference.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        name: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.experiment.Experiment.GoalReference}
 */
proto.bucketeer.experiment.Experiment.GoalReference.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.experiment.Experiment.GoalReference();
    return proto.bucketeer.experiment.Experiment.GoalReference.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.experiment.Experiment.GoalReference} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.experiment.Experiment.GoalReference}
 */
proto.bucketeer.experiment.Experiment.GoalReference.deserializeBinaryFromReader =
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
proto.bucketeer.experiment.Experiment.GoalReference.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.experiment.Experiment.GoalReference.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.experiment.Experiment.GoalReference} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.experiment.Experiment.GoalReference.serializeBinaryToWriter =
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
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.GoalReference.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment.GoalReference} returns this
 */
proto.bucketeer.experiment.Experiment.GoalReference.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string name = 2;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.GoalReference.prototype.getName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment.GoalReference} returns this
 */
proto.bucketeer.experiment.Experiment.GoalReference.prototype.setName =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string goal_id = 2;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getGoalId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setGoalId = function (value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string feature_id = 3;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getFeatureId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setFeatureId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * optional int32 feature_version = 4;
 * @return {number}
 */
proto.bucketeer.experiment.Experiment.prototype.getFeatureVersion =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setFeatureVersion = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 4, value);
};

/**
 * repeated bucketeer.feature.Variation variations = 5;
 * @return {!Array<!proto.bucketeer.feature.Variation>}
 */
proto.bucketeer.experiment.Experiment.prototype.getVariationsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.feature.Variation>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_feature_variation_pb.Variation,
        5
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.feature.Variation>} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setVariationsList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 5, value);
};

/**
 * @param {!proto.bucketeer.feature.Variation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.feature.Variation}
 */
proto.bucketeer.experiment.Experiment.prototype.addVariations = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    5,
    opt_value,
    proto.bucketeer.feature.Variation,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.clearVariationsList =
  function () {
    return this.setVariationsList([]);
  };

/**
 * optional int64 start_at = 6;
 * @return {number}
 */
proto.bucketeer.experiment.Experiment.prototype.getStartAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setStartAt = function (value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};

/**
 * optional int64 stop_at = 7;
 * @return {number}
 */
proto.bucketeer.experiment.Experiment.prototype.getStopAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setStopAt = function (value) {
  return jspb.Message.setProto3IntField(this, 7, value);
};

/**
 * optional bool stopped = 8;
 * @return {boolean}
 */
proto.bucketeer.experiment.Experiment.prototype.getStopped = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 8, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setStopped = function (value) {
  return jspb.Message.setProto3BooleanField(this, 8, value);
};

/**
 * optional int64 stopped_at = 9;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getStoppedAt = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, '0'));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setStoppedAt = function (
  value
) {
  return jspb.Message.setProto3StringIntField(this, 9, value);
};

/**
 * optional int64 created_at = 10;
 * @return {number}
 */
proto.bucketeer.experiment.Experiment.prototype.getCreatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 10, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setCreatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 10, value);
};

/**
 * optional int64 updated_at = 11;
 * @return {number}
 */
proto.bucketeer.experiment.Experiment.prototype.getUpdatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setUpdatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 11, value);
};

/**
 * optional bool deleted = 12;
 * @return {boolean}
 */
proto.bucketeer.experiment.Experiment.prototype.getDeleted = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 12, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setDeleted = function (value) {
  return jspb.Message.setProto3BooleanField(this, 12, value);
};

/**
 * repeated string goal_ids = 13;
 * @return {!Array<string>}
 */
proto.bucketeer.experiment.Experiment.prototype.getGoalIdsList = function () {
  return /** @type {!Array<string>} */ (
    jspb.Message.getRepeatedField(this, 13)
  );
};

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setGoalIdsList = function (
  value
) {
  return jspb.Message.setField(this, 13, value || []);
};

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.addGoalIds = function (
  value,
  opt_index
) {
  return jspb.Message.addToRepeatedField(this, 13, value, opt_index);
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.clearGoalIdsList = function () {
  return this.setGoalIdsList([]);
};

/**
 * optional string name = 14;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 14, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setName = function (value) {
  return jspb.Message.setProto3StringField(this, 14, value);
};

/**
 * optional string description = 15;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getDescription = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 15, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setDescription = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 15, value);
};

/**
 * optional string base_variation_id = 16;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getBaseVariationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 16, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setBaseVariationId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 16, value);
};

/**
 * optional Status status = 18;
 * @return {!proto.bucketeer.experiment.Experiment.Status}
 */
proto.bucketeer.experiment.Experiment.prototype.getStatus = function () {
  return /** @type {!proto.bucketeer.experiment.Experiment.Status} */ (
    jspb.Message.getFieldWithDefault(this, 18, 0)
  );
};

/**
 * @param {!proto.bucketeer.experiment.Experiment.Status} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setStatus = function (value) {
  return jspb.Message.setProto3EnumField(this, 18, value);
};

/**
 * optional string maintainer = 19;
 * @return {string}
 */
proto.bucketeer.experiment.Experiment.prototype.getMaintainer = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 19, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setMaintainer = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 19, value);
};

/**
 * optional bool archived = 20;
 * @return {boolean}
 */
proto.bucketeer.experiment.Experiment.prototype.getArchived = function () {
  return /** @type {boolean} */ (
    jspb.Message.getBooleanFieldWithDefault(this, 20, false)
  );
};

/**
 * @param {boolean} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setArchived = function (value) {
  return jspb.Message.setProto3BooleanField(this, 20, value);
};

/**
 * repeated GoalReference goals = 21;
 * @return {!Array<!proto.bucketeer.experiment.Experiment.GoalReference>}
 */
proto.bucketeer.experiment.Experiment.prototype.getGoalsList = function () {
  return /** @type{!Array<!proto.bucketeer.experiment.Experiment.GoalReference>} */ (
    jspb.Message.getRepeatedWrapperField(
      this,
      proto.bucketeer.experiment.Experiment.GoalReference,
      21
    )
  );
};

/**
 * @param {!Array<!proto.bucketeer.experiment.Experiment.GoalReference>} value
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.setGoalsList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 21, value);
};

/**
 * @param {!proto.bucketeer.experiment.Experiment.GoalReference=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.experiment.Experiment.GoalReference}
 */
proto.bucketeer.experiment.Experiment.prototype.addGoals = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    21,
    opt_value,
    proto.bucketeer.experiment.Experiment.GoalReference,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.experiment.Experiment} returns this
 */
proto.bucketeer.experiment.Experiment.prototype.clearGoalsList = function () {
  return this.setGoalsList([]);
};

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.experiment.Experiments.repeatedFields_ = [1];

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
  proto.bucketeer.experiment.Experiments.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.experiment.Experiments.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.experiment.Experiments} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.experiment.Experiments.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        experimentsList: jspb.Message.toObjectList(
          msg.getExperimentsList(),
          proto.bucketeer.experiment.Experiment.toObject,
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
 * @return {!proto.bucketeer.experiment.Experiments}
 */
proto.bucketeer.experiment.Experiments.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.experiment.Experiments();
  return proto.bucketeer.experiment.Experiments.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.experiment.Experiments} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.experiment.Experiments}
 */
proto.bucketeer.experiment.Experiments.deserializeBinaryFromReader = function (
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
        var value = new proto.bucketeer.experiment.Experiment();
        reader.readMessage(
          value,
          proto.bucketeer.experiment.Experiment.deserializeBinaryFromReader
        );
        msg.addExperiments(value);
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
proto.bucketeer.experiment.Experiments.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.experiment.Experiments.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.experiment.Experiments} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.experiment.Experiments.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getExperimentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.bucketeer.experiment.Experiment.serializeBinaryToWriter
    );
  }
};

/**
 * repeated Experiment experiments = 1;
 * @return {!Array<!proto.bucketeer.experiment.Experiment>}
 */
proto.bucketeer.experiment.Experiments.prototype.getExperimentsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.experiment.Experiment>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto.bucketeer.experiment.Experiment,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.experiment.Experiment>} value
 * @return {!proto.bucketeer.experiment.Experiments} returns this
 */
proto.bucketeer.experiment.Experiments.prototype.setExperimentsList = function (
  value
) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};

/**
 * @param {!proto.bucketeer.experiment.Experiment=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.experiment.Experiment}
 */
proto.bucketeer.experiment.Experiments.prototype.addExperiments = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    1,
    opt_value,
    proto.bucketeer.experiment.Experiment,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.experiment.Experiments} returns this
 */
proto.bucketeer.experiment.Experiments.prototype.clearExperimentsList =
  function () {
    return this.setExperimentsList([]);
  };

goog.object.extend(exports, proto.bucketeer.experiment);
