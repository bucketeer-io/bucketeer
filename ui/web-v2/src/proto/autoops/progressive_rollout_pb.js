// source: proto/autoops/progressive_rollout.proto
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

var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');
goog.object.extend(proto, google_protobuf_any_pb);
goog.exportSymbol('proto.bucketeer.autoops.ProgressiveRollout', null, global);
goog.exportSymbol(
  'proto.bucketeer.autoops.ProgressiveRollout.Status',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ProgressiveRollout.StoppedBy',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.autoops.ProgressiveRollout.Type',
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
proto.bucketeer.autoops.ProgressiveRollout = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.autoops.ProgressiveRollout, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.autoops.ProgressiveRollout.displayName =
    'proto.bucketeer.autoops.ProgressiveRollout';
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
  proto.bucketeer.autoops.ProgressiveRollout.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.autoops.ProgressiveRollout.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.autoops.ProgressiveRollout} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.autoops.ProgressiveRollout.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        clause:
          (f = msg.getClause()) &&
          google_protobuf_any_pb.Any.toObject(includeInstance, f),
        status: jspb.Message.getFieldWithDefault(msg, 4, 0),
        createdAt: jspb.Message.getFieldWithDefault(msg, 5, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 6, 0),
        type: jspb.Message.getFieldWithDefault(msg, 7, 0),
        stoppedBy: jspb.Message.getFieldWithDefault(msg, 8, 0),
        stoppedAt: jspb.Message.getFieldWithDefault(msg, 9, 0)
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
 * @return {!proto.bucketeer.autoops.ProgressiveRollout}
 */
proto.bucketeer.autoops.ProgressiveRollout.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.autoops.ProgressiveRollout();
  return proto.bucketeer.autoops.ProgressiveRollout.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.autoops.ProgressiveRollout} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.autoops.ProgressiveRollout}
 */
proto.bucketeer.autoops.ProgressiveRollout.deserializeBinaryFromReader =
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
          var value = new google_protobuf_any_pb.Any();
          reader.readMessage(
            value,
            google_protobuf_any_pb.Any.deserializeBinaryFromReader
          );
          msg.setClause(value);
          break;
        case 4:
          var value =
            /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Status} */ (
              reader.readEnum()
            );
          msg.setStatus(value);
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
          var value =
            /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Type} */ (
              reader.readEnum()
            );
          msg.setType(value);
          break;
        case 8:
          var value =
            /** @type {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy} */ (
              reader.readEnum()
            );
          msg.setStoppedBy(value);
          break;
        case 9:
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
proto.bucketeer.autoops.ProgressiveRollout.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.autoops.ProgressiveRollout.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.autoops.ProgressiveRollout} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.autoops.ProgressiveRollout.serializeBinaryToWriter = function (
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
  f = message.getClause();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
  f = message.getStatus();
  if (f !== 0.0) {
    writer.writeEnum(4, f);
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(5, f);
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(6, f);
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(7, f);
  }
  f = message.getStoppedBy();
  if (f !== 0.0) {
    writer.writeEnum(8, f);
  }
  f = message.getStoppedAt();
  if (f !== 0) {
    writer.writeInt64(9, f);
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.autoops.ProgressiveRollout.Type = {
  MANUAL_SCHEDULE: 0,
  TEMPLATE_SCHEDULE: 1
};

/**
 * @enum {number}
 */
proto.bucketeer.autoops.ProgressiveRollout.Status = {
  WAITING: 0,
  RUNNING: 1,
  FINISHED: 2,
  STOPPED: 3
};

/**
 * @enum {number}
 */
proto.bucketeer.autoops.ProgressiveRollout.StoppedBy = {
  UNKNOWN: 0,
  USER: 1,
  OPS_SCHEDULE: 2,
  OPS_KILL_SWITCH: 3
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getFeatureId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setFeatureId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional google.protobuf.Any clause = 3;
 * @return {?proto.google.protobuf.Any}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getClause = function () {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 3)
  );
};

/**
 * @param {?proto.google.protobuf.Any|undefined} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setClause = function (
  value
) {
  return jspb.Message.setWrapperField(this, 3, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.clearClause = function () {
  return this.setClause(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.hasClause = function () {
  return jspb.Message.getField(this, 3) != null;
};

/**
 * optional Status status = 4;
 * @return {!proto.bucketeer.autoops.ProgressiveRollout.Status}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getStatus = function () {
  return /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Status} */ (
    jspb.Message.getFieldWithDefault(this, 4, 0)
  );
};

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout.Status} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setStatus = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};

/**
 * optional int64 created_at = 5;
 * @return {number}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getCreatedAt =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setCreatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 5, value);
};

/**
 * optional int64 updated_at = 6;
 * @return {number}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getUpdatedAt =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setUpdatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 6, value);
};

/**
 * optional Type type = 7;
 * @return {!proto.bucketeer.autoops.ProgressiveRollout.Type}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getType = function () {
  return /** @type {!proto.bucketeer.autoops.ProgressiveRollout.Type} */ (
    jspb.Message.getFieldWithDefault(this, 7, 0)
  );
};

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout.Type} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setType = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 7, value);
};

/**
 * optional StoppedBy stopped_by = 8;
 * @return {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getStoppedBy =
  function () {
    return /** @type {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy} */ (
      jspb.Message.getFieldWithDefault(this, 8, 0)
    );
  };

/**
 * @param {!proto.bucketeer.autoops.ProgressiveRollout.StoppedBy} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setStoppedBy = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 8, value);
};

/**
 * optional int64 stopped_at = 9;
 * @return {number}
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.getStoppedAt =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.autoops.ProgressiveRollout} returns this
 */
proto.bucketeer.autoops.ProgressiveRollout.prototype.setStoppedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 9, value);
};

goog.object.extend(exports, proto.bucketeer.autoops);
