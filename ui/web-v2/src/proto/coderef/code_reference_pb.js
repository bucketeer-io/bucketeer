// source: proto/coderef/code_reference.proto
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

goog.exportSymbol('proto.bucketeer.coderef.CodeReference', null, global);
goog.exportSymbol(
  'proto.bucketeer.coderef.CodeReference.RepositoryType',
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
proto.bucketeer.coderef.CodeReference = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.coderef.CodeReference.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.coderef.CodeReference, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.coderef.CodeReference.displayName =
    'proto.bucketeer.coderef.CodeReference';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.coderef.CodeReference.repeatedFields_ = [7];

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
  proto.bucketeer.coderef.CodeReference.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.coderef.CodeReference.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.coderef.CodeReference} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.coderef.CodeReference.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        featureId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        filePath: jspb.Message.getFieldWithDefault(msg, 3, ''),
        lineNumber: jspb.Message.getFieldWithDefault(msg, 4, 0),
        codeSnippet: jspb.Message.getFieldWithDefault(msg, 5, ''),
        contentHash: jspb.Message.getFieldWithDefault(msg, 6, ''),
        aliasesList:
          (f = jspb.Message.getRepeatedField(msg, 7)) == null ? undefined : f,
        repositoryName: jspb.Message.getFieldWithDefault(msg, 8, ''),
        repositoryOwner: jspb.Message.getFieldWithDefault(msg, 9, ''),
        repositoryType: jspb.Message.getFieldWithDefault(msg, 10, 0),
        repositoryBranch: jspb.Message.getFieldWithDefault(msg, 11, ''),
        commitHash: jspb.Message.getFieldWithDefault(msg, 12, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 13, ''),
        createdAt: jspb.Message.getFieldWithDefault(msg, 14, 0),
        updatedAt: jspb.Message.getFieldWithDefault(msg, 15, 0),
        sourceUrl: jspb.Message.getFieldWithDefault(msg, 16, ''),
        branchUrl: jspb.Message.getFieldWithDefault(msg, 17, ''),
        fileExtension: jspb.Message.getFieldWithDefault(msg, 18, '')
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
 * @return {!proto.bucketeer.coderef.CodeReference}
 */
proto.bucketeer.coderef.CodeReference.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.coderef.CodeReference();
  return proto.bucketeer.coderef.CodeReference.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.coderef.CodeReference} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.coderef.CodeReference}
 */
proto.bucketeer.coderef.CodeReference.deserializeBinaryFromReader = function (
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
        var value = /** @type {string} */ (reader.readString());
        msg.setFilePath(value);
        break;
      case 4:
        var value = /** @type {number} */ (reader.readInt32());
        msg.setLineNumber(value);
        break;
      case 5:
        var value = /** @type {string} */ (reader.readString());
        msg.setCodeSnippet(value);
        break;
      case 6:
        var value = /** @type {string} */ (reader.readString());
        msg.setContentHash(value);
        break;
      case 7:
        var value = /** @type {string} */ (reader.readString());
        msg.addAliases(value);
        break;
      case 8:
        var value = /** @type {string} */ (reader.readString());
        msg.setRepositoryName(value);
        break;
      case 9:
        var value = /** @type {string} */ (reader.readString());
        msg.setRepositoryOwner(value);
        break;
      case 10:
        var value =
          /** @type {!proto.bucketeer.coderef.CodeReference.RepositoryType} */ (
            reader.readEnum()
          );
        msg.setRepositoryType(value);
        break;
      case 11:
        var value = /** @type {string} */ (reader.readString());
        msg.setRepositoryBranch(value);
        break;
      case 12:
        var value = /** @type {string} */ (reader.readString());
        msg.setCommitHash(value);
        break;
      case 13:
        var value = /** @type {string} */ (reader.readString());
        msg.setEnvironmentId(value);
        break;
      case 14:
        var value = /** @type {number} */ (reader.readInt64());
        msg.setCreatedAt(value);
        break;
      case 15:
        var value = /** @type {number} */ (reader.readInt64());
        msg.setUpdatedAt(value);
        break;
      case 16:
        var value = /** @type {string} */ (reader.readString());
        msg.setSourceUrl(value);
        break;
      case 17:
        var value = /** @type {string} */ (reader.readString());
        msg.setBranchUrl(value);
        break;
      case 18:
        var value = /** @type {string} */ (reader.readString());
        msg.setFileExtension(value);
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
proto.bucketeer.coderef.CodeReference.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.coderef.CodeReference.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.coderef.CodeReference} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.coderef.CodeReference.serializeBinaryToWriter = function (
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
  f = message.getFilePath();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
  f = message.getLineNumber();
  if (f !== 0) {
    writer.writeInt32(4, f);
  }
  f = message.getCodeSnippet();
  if (f.length > 0) {
    writer.writeString(5, f);
  }
  f = message.getContentHash();
  if (f.length > 0) {
    writer.writeString(6, f);
  }
  f = message.getAliasesList();
  if (f.length > 0) {
    writer.writeRepeatedString(7, f);
  }
  f = message.getRepositoryName();
  if (f.length > 0) {
    writer.writeString(8, f);
  }
  f = message.getRepositoryOwner();
  if (f.length > 0) {
    writer.writeString(9, f);
  }
  f = message.getRepositoryType();
  if (f !== 0.0) {
    writer.writeEnum(10, f);
  }
  f = message.getRepositoryBranch();
  if (f.length > 0) {
    writer.writeString(11, f);
  }
  f = message.getCommitHash();
  if (f.length > 0) {
    writer.writeString(12, f);
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(13, f);
  }
  f = message.getCreatedAt();
  if (f !== 0) {
    writer.writeInt64(14, f);
  }
  f = message.getUpdatedAt();
  if (f !== 0) {
    writer.writeInt64(15, f);
  }
  f = message.getSourceUrl();
  if (f.length > 0) {
    writer.writeString(16, f);
  }
  f = message.getBranchUrl();
  if (f.length > 0) {
    writer.writeString(17, f);
  }
  f = message.getFileExtension();
  if (f.length > 0) {
    writer.writeString(18, f);
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.coderef.CodeReference.RepositoryType = {
  REPOSITORY_TYPE_UNSPECIFIED: 0,
  GITHUB: 1,
  GITLAB: 2,
  BITBUCKET: 3,
  CUSTOM: 4
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string feature_id = 2;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getFeatureId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setFeatureId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string file_path = 3;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getFilePath = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setFilePath = function (value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * optional int32 line_number = 4;
 * @return {number}
 */
proto.bucketeer.coderef.CodeReference.prototype.getLineNumber = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setLineNumber = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 4, value);
};

/**
 * optional string code_snippet = 5;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getCodeSnippet = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setCodeSnippet = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 5, value);
};

/**
 * optional string content_hash = 6;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getContentHash = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setContentHash = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 6, value);
};

/**
 * repeated string aliases = 7;
 * @return {!Array<string>}
 */
proto.bucketeer.coderef.CodeReference.prototype.getAliasesList = function () {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 7));
};

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setAliasesList = function (
  value
) {
  return jspb.Message.setField(this, 7, value || []);
};

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.addAliases = function (
  value,
  opt_index
) {
  return jspb.Message.addToRepeatedField(this, 7, value, opt_index);
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.clearAliasesList = function () {
  return this.setAliasesList([]);
};

/**
 * optional string repository_name = 8;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getRepositoryName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 8, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setRepositoryName = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 8, value);
};

/**
 * optional string repository_owner = 9;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getRepositoryOwner =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 9, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setRepositoryOwner = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 9, value);
};

/**
 * optional RepositoryType repository_type = 10;
 * @return {!proto.bucketeer.coderef.CodeReference.RepositoryType}
 */
proto.bucketeer.coderef.CodeReference.prototype.getRepositoryType =
  function () {
    return /** @type {!proto.bucketeer.coderef.CodeReference.RepositoryType} */ (
      jspb.Message.getFieldWithDefault(this, 10, 0)
    );
  };

/**
 * @param {!proto.bucketeer.coderef.CodeReference.RepositoryType} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setRepositoryType = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 10, value);
};

/**
 * optional string repository_branch = 11;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getRepositoryBranch =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 11, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setRepositoryBranch = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 11, value);
};

/**
 * optional string commit_hash = 12;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getCommitHash = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 12, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setCommitHash = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 12, value);
};

/**
 * optional string environment_id = 13;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getEnvironmentId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 13, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setEnvironmentId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 13, value);
};

/**
 * optional int64 created_at = 14;
 * @return {number}
 */
proto.bucketeer.coderef.CodeReference.prototype.getCreatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 14, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setCreatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 14, value);
};

/**
 * optional int64 updated_at = 15;
 * @return {number}
 */
proto.bucketeer.coderef.CodeReference.prototype.getUpdatedAt = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 15, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setUpdatedAt = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 15, value);
};

/**
 * optional string source_url = 16;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getSourceUrl = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 16, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setSourceUrl = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 16, value);
};

/**
 * optional string branch_url = 17;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getBranchUrl = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 17, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setBranchUrl = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 17, value);
};

/**
 * optional string file_extension = 18;
 * @return {string}
 */
proto.bucketeer.coderef.CodeReference.prototype.getFileExtension = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 18, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.coderef.CodeReference} returns this
 */
proto.bucketeer.coderef.CodeReference.prototype.setFileExtension = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 18, value);
};

goog.object.extend(exports, proto.bucketeer.coderef);
