// source: proto/account/service.proto
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

var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js');
goog.object.extend(proto, google_protobuf_wrappers_pb);
var google_api_annotations_pb = require('../../google/api/annotations_pb.js');
goog.object.extend(proto, google_api_annotations_pb);
var google_api_field_behavior_pb = require('../../google/api/field_behavior_pb.js');
goog.object.extend(proto, google_api_field_behavior_pb);
var protoc$gen$openapiv2_options_annotations_pb = require('../../protoc-gen-openapiv2/options/annotations_pb.js');
goog.object.extend(proto, protoc$gen$openapiv2_options_annotations_pb);
var proto_account_account_pb = require('../../proto/account/account_pb.js');
goog.object.extend(proto, proto_account_account_pb);
var proto_account_api_key_pb = require('../../proto/account/api_key_pb.js');
goog.object.extend(proto, proto_account_api_key_pb);
var proto_account_command_pb = require('../../proto/account/command_pb.js');
goog.object.extend(proto, proto_account_command_pb);
var proto_environment_organization_pb = require('../../proto/environment/organization_pb.js');
goog.object.extend(proto, proto_environment_organization_pb);
goog.exportSymbol(
  'proto.bucketeer.account.ChangeAPIKeyNameRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.ChangeAPIKeyNameResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.CreateAPIKeyRequest', null, global);
goog.exportSymbol('proto.bucketeer.account.CreateAPIKeyResponse', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.CreateAccountV2Request',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.CreateAccountV2Response',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.CreateSearchFilterRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.CreateSearchFilterResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.DeleteAccountV2Request',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.DeleteAccountV2Response',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.DeleteSearchFilterRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.DeleteSearchFilterResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.DisableAPIKeyRequest', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.DisableAPIKeyResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.DisableAccountV2Request',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.DisableAccountV2Response',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.EnableAPIKeyRequest', null, global);
goog.exportSymbol('proto.bucketeer.account.EnableAPIKeyResponse', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.EnableAccountV2Request',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.EnableAccountV2Response',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.GetAPIKeyRequest', null, global);
goog.exportSymbol('proto.bucketeer.account.GetAPIKeyResponse', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.GetAccountV2Request', null, global);
goog.exportSymbol('proto.bucketeer.account.GetAccountV2Response', null, global);
goog.exportSymbol('proto.bucketeer.account.GetMeRequest', null, global);
goog.exportSymbol('proto.bucketeer.account.GetMeResponse', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.GetMyOrganizationsByEmailRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.GetMyOrganizationsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.GetMyOrganizationsResponse',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.ListAPIKeysRequest', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.ListAPIKeysRequest.OrderBy',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.ListAPIKeysRequest.OrderDirection',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.ListAPIKeysResponse', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.ListAccountsV2Request',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.ListAccountsV2Request.OrderBy',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.ListAccountsV2Request.OrderDirection',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.ListAccountsV2Response',
  null,
  global
);
goog.exportSymbol('proto.bucketeer.account.UpdateAPIKeyRequest', null, global);
goog.exportSymbol('proto.bucketeer.account.UpdateAPIKeyResponse', null, global);
goog.exportSymbol(
  'proto.bucketeer.account.UpdateAccountV2Request',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.UpdateAccountV2Response',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.UpdateSearchFilterRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.account.UpdateSearchFilterResponse',
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
proto.bucketeer.account.GetMeRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetMeRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetMeRequest.displayName =
    'proto.bucketeer.account.GetMeRequest';
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
proto.bucketeer.account.GetMeResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetMeResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetMeResponse.displayName =
    'proto.bucketeer.account.GetMeResponse';
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
proto.bucketeer.account.GetMyOrganizationsRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetMyOrganizationsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetMyOrganizationsRequest.displayName =
    'proto.bucketeer.account.GetMyOrganizationsRequest';
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
proto.bucketeer.account.GetMyOrganizationsByEmailRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.account.GetMyOrganizationsByEmailRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetMyOrganizationsByEmailRequest.displayName =
    'proto.bucketeer.account.GetMyOrganizationsByEmailRequest';
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
proto.bucketeer.account.GetMyOrganizationsResponse = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.account.GetMyOrganizationsResponse.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.account.GetMyOrganizationsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetMyOrganizationsResponse.displayName =
    'proto.bucketeer.account.GetMyOrganizationsResponse';
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
proto.bucketeer.account.CreateAccountV2Request = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.account.CreateAccountV2Request.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.account.CreateAccountV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.CreateAccountV2Request.displayName =
    'proto.bucketeer.account.CreateAccountV2Request';
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
proto.bucketeer.account.CreateAccountV2Response = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.CreateAccountV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.CreateAccountV2Response.displayName =
    'proto.bucketeer.account.CreateAccountV2Response';
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
proto.bucketeer.account.EnableAccountV2Request = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.EnableAccountV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.EnableAccountV2Request.displayName =
    'proto.bucketeer.account.EnableAccountV2Request';
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
proto.bucketeer.account.EnableAccountV2Response = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.EnableAccountV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.EnableAccountV2Response.displayName =
    'proto.bucketeer.account.EnableAccountV2Response';
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
proto.bucketeer.account.DisableAccountV2Request = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DisableAccountV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DisableAccountV2Request.displayName =
    'proto.bucketeer.account.DisableAccountV2Request';
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
proto.bucketeer.account.DisableAccountV2Response = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DisableAccountV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DisableAccountV2Response.displayName =
    'proto.bucketeer.account.DisableAccountV2Response';
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
proto.bucketeer.account.DeleteAccountV2Request = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DeleteAccountV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DeleteAccountV2Request.displayName =
    'proto.bucketeer.account.DeleteAccountV2Request';
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
proto.bucketeer.account.DeleteAccountV2Response = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DeleteAccountV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DeleteAccountV2Response.displayName =
    'proto.bucketeer.account.DeleteAccountV2Response';
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
proto.bucketeer.account.UpdateAccountV2Request = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.account.UpdateAccountV2Request.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.account.UpdateAccountV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateAccountV2Request.displayName =
    'proto.bucketeer.account.UpdateAccountV2Request';
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
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.displayName =
    'proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar';
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
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue =
  function (opt_data) {
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  };
goog.inherits(
  proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.displayName =
    'proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue';
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
proto.bucketeer.account.UpdateAccountV2Response = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.UpdateAccountV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateAccountV2Response.displayName =
    'proto.bucketeer.account.UpdateAccountV2Response';
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
proto.bucketeer.account.GetAccountV2Request = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetAccountV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAccountV2Request.displayName =
    'proto.bucketeer.account.GetAccountV2Request';
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
proto.bucketeer.account.GetAccountV2Response = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetAccountV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAccountV2Response.displayName =
    'proto.bucketeer.account.GetAccountV2Response';
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
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.displayName =
    'proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest';
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
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.displayName =
    'proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse';
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
proto.bucketeer.account.ListAccountsV2Request = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.ListAccountsV2Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.ListAccountsV2Request.displayName =
    'proto.bucketeer.account.ListAccountsV2Request';
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
proto.bucketeer.account.ListAccountsV2Response = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.account.ListAccountsV2Response.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.account.ListAccountsV2Response, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.ListAccountsV2Response.displayName =
    'proto.bucketeer.account.ListAccountsV2Response';
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
proto.bucketeer.account.CreateAPIKeyRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.CreateAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.CreateAPIKeyRequest.displayName =
    'proto.bucketeer.account.CreateAPIKeyRequest';
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
proto.bucketeer.account.CreateAPIKeyResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.CreateAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.CreateAPIKeyResponse.displayName =
    'proto.bucketeer.account.CreateAPIKeyResponse';
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
proto.bucketeer.account.ChangeAPIKeyNameRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.ChangeAPIKeyNameRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.ChangeAPIKeyNameRequest.displayName =
    'proto.bucketeer.account.ChangeAPIKeyNameRequest';
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
proto.bucketeer.account.ChangeAPIKeyNameResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.ChangeAPIKeyNameResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.ChangeAPIKeyNameResponse.displayName =
    'proto.bucketeer.account.ChangeAPIKeyNameResponse';
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
proto.bucketeer.account.EnableAPIKeyRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.EnableAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.EnableAPIKeyRequest.displayName =
    'proto.bucketeer.account.EnableAPIKeyRequest';
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
proto.bucketeer.account.EnableAPIKeyResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.EnableAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.EnableAPIKeyResponse.displayName =
    'proto.bucketeer.account.EnableAPIKeyResponse';
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
proto.bucketeer.account.DisableAPIKeyRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DisableAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DisableAPIKeyRequest.displayName =
    'proto.bucketeer.account.DisableAPIKeyRequest';
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
proto.bucketeer.account.DisableAPIKeyResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DisableAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DisableAPIKeyResponse.displayName =
    'proto.bucketeer.account.DisableAPIKeyResponse';
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
proto.bucketeer.account.GetAPIKeyRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAPIKeyRequest.displayName =
    'proto.bucketeer.account.GetAPIKeyRequest';
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
proto.bucketeer.account.GetAPIKeyResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.GetAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAPIKeyResponse.displayName =
    'proto.bucketeer.account.GetAPIKeyResponse';
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
proto.bucketeer.account.ListAPIKeysRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.account.ListAPIKeysRequest.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.account.ListAPIKeysRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.ListAPIKeysRequest.displayName =
    'proto.bucketeer.account.ListAPIKeysRequest';
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
proto.bucketeer.account.ListAPIKeysResponse = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.account.ListAPIKeysResponse.repeatedFields_,
    null
  );
};
goog.inherits(proto.bucketeer.account.ListAPIKeysResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.ListAPIKeysResponse.displayName =
    'proto.bucketeer.account.ListAPIKeysResponse';
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
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.displayName =
    'proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest';
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
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.displayName =
    'proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse';
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
proto.bucketeer.account.CreateSearchFilterRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.CreateSearchFilterRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.CreateSearchFilterRequest.displayName =
    'proto.bucketeer.account.CreateSearchFilterRequest';
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
proto.bucketeer.account.CreateSearchFilterResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.CreateSearchFilterResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.CreateSearchFilterResponse.displayName =
    'proto.bucketeer.account.CreateSearchFilterResponse';
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
proto.bucketeer.account.UpdateSearchFilterRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.UpdateSearchFilterRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateSearchFilterRequest.displayName =
    'proto.bucketeer.account.UpdateSearchFilterRequest';
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
proto.bucketeer.account.UpdateSearchFilterResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.UpdateSearchFilterResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateSearchFilterResponse.displayName =
    'proto.bucketeer.account.UpdateSearchFilterResponse';
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
proto.bucketeer.account.DeleteSearchFilterRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DeleteSearchFilterRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DeleteSearchFilterRequest.displayName =
    'proto.bucketeer.account.DeleteSearchFilterRequest';
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
proto.bucketeer.account.DeleteSearchFilterResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.DeleteSearchFilterResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.DeleteSearchFilterResponse.displayName =
    'proto.bucketeer.account.DeleteSearchFilterResponse';
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
proto.bucketeer.account.UpdateAPIKeyRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.UpdateAPIKeyRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateAPIKeyRequest.displayName =
    'proto.bucketeer.account.UpdateAPIKeyRequest';
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
proto.bucketeer.account.UpdateAPIKeyResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.bucketeer.account.UpdateAPIKeyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.account.UpdateAPIKeyResponse.displayName =
    'proto.bucketeer.account.UpdateAPIKeyResponse';
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
  proto.bucketeer.account.GetMeRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.GetMeRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetMeRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetMeRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        organizationId: jspb.Message.getFieldWithDefault(msg, 1, '')
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
 * @return {!proto.bucketeer.account.GetMeRequest}
 */
proto.bucketeer.account.GetMeRequest.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetMeRequest();
  return proto.bucketeer.account.GetMeRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetMeRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetMeRequest}
 */
proto.bucketeer.account.GetMeRequest.deserializeBinaryFromReader = function (
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
        msg.setOrganizationId(value);
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
proto.bucketeer.account.GetMeRequest.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.account.GetMeRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetMeRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetMeRequest.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getOrganizationId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
};

/**
 * optional string organization_id = 1;
 * @return {string}
 */
proto.bucketeer.account.GetMeRequest.prototype.getOrganizationId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetMeRequest} returns this
 */
proto.bucketeer.account.GetMeRequest.prototype.setOrganizationId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
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
  proto.bucketeer.account.GetMeResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.GetMeResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetMeResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetMeResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        account:
          (f = msg.getAccount()) &&
          proto_account_account_pb.ConsoleAccount.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.GetMeResponse}
 */
proto.bucketeer.account.GetMeResponse.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetMeResponse();
  return proto.bucketeer.account.GetMeResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetMeResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetMeResponse}
 */
proto.bucketeer.account.GetMeResponse.deserializeBinaryFromReader = function (
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
        var value = new proto_account_account_pb.ConsoleAccount();
        reader.readMessage(
          value,
          proto_account_account_pb.ConsoleAccount.deserializeBinaryFromReader
        );
        msg.setAccount(value);
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
proto.bucketeer.account.GetMeResponse.prototype.serializeBinary = function () {
  var writer = new jspb.BinaryWriter();
  proto.bucketeer.account.GetMeResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetMeResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetMeResponse.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getAccount();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_account_account_pb.ConsoleAccount.serializeBinaryToWriter
    );
  }
};

/**
 * optional ConsoleAccount account = 1;
 * @return {?proto.bucketeer.account.ConsoleAccount}
 */
proto.bucketeer.account.GetMeResponse.prototype.getAccount = function () {
  return /** @type{?proto.bucketeer.account.ConsoleAccount} */ (
    jspb.Message.getWrapperField(
      this,
      proto_account_account_pb.ConsoleAccount,
      1
    )
  );
};

/**
 * @param {?proto.bucketeer.account.ConsoleAccount|undefined} value
 * @return {!proto.bucketeer.account.GetMeResponse} returns this
 */
proto.bucketeer.account.GetMeResponse.prototype.setAccount = function (value) {
  return jspb.Message.setWrapperField(this, 1, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.GetMeResponse} returns this
 */
proto.bucketeer.account.GetMeResponse.prototype.clearAccount = function () {
  return this.setAccount(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.GetMeResponse.prototype.hasAccount = function () {
  return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.GetMyOrganizationsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetMyOrganizationsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetMyOrganizationsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetMyOrganizationsRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.GetMyOrganizationsRequest}
 */
proto.bucketeer.account.GetMyOrganizationsRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetMyOrganizationsRequest();
  return proto.bucketeer.account.GetMyOrganizationsRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetMyOrganizationsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetMyOrganizationsRequest}
 */
proto.bucketeer.account.GetMyOrganizationsRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.GetMyOrganizationsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetMyOrganizationsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetMyOrganizationsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetMyOrganizationsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.GetMyOrganizationsByEmailRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetMyOrganizationsByEmailRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetMyOrganizationsByEmailRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetMyOrganizationsByEmailRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, '')
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
 * @return {!proto.bucketeer.account.GetMyOrganizationsByEmailRequest}
 */
proto.bucketeer.account.GetMyOrganizationsByEmailRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.GetMyOrganizationsByEmailRequest();
    return proto.bucketeer.account.GetMyOrganizationsByEmailRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetMyOrganizationsByEmailRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetMyOrganizationsByEmailRequest}
 */
proto.bucketeer.account.GetMyOrganizationsByEmailRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
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
proto.bucketeer.account.GetMyOrganizationsByEmailRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetMyOrganizationsByEmailRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetMyOrganizationsByEmailRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetMyOrganizationsByEmailRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.GetMyOrganizationsByEmailRequest.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetMyOrganizationsByEmailRequest} returns this
 */
proto.bucketeer.account.GetMyOrganizationsByEmailRequest.prototype.setEmail =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.account.GetMyOrganizationsResponse.repeatedFields_ = [1];

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
  proto.bucketeer.account.GetMyOrganizationsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetMyOrganizationsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetMyOrganizationsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetMyOrganizationsResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        organizationsList: jspb.Message.toObjectList(
          msg.getOrganizationsList(),
          proto_environment_organization_pb.Organization.toObject,
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
 * @return {!proto.bucketeer.account.GetMyOrganizationsResponse}
 */
proto.bucketeer.account.GetMyOrganizationsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.GetMyOrganizationsResponse();
    return proto.bucketeer.account.GetMyOrganizationsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetMyOrganizationsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetMyOrganizationsResponse}
 */
proto.bucketeer.account.GetMyOrganizationsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_environment_organization_pb.Organization();
          reader.readMessage(
            value,
            proto_environment_organization_pb.Organization
              .deserializeBinaryFromReader
          );
          msg.addOrganizations(value);
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
proto.bucketeer.account.GetMyOrganizationsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetMyOrganizationsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetMyOrganizationsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetMyOrganizationsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getOrganizationsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_environment_organization_pb.Organization.serializeBinaryToWriter
      );
    }
  };

/**
 * repeated bucketeer.environment.Organization organizations = 1;
 * @return {!Array<!proto.bucketeer.environment.Organization>}
 */
proto.bucketeer.account.GetMyOrganizationsResponse.prototype.getOrganizationsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.environment.Organization>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_environment_organization_pb.Organization,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.environment.Organization>} value
 * @return {!proto.bucketeer.account.GetMyOrganizationsResponse} returns this
 */
proto.bucketeer.account.GetMyOrganizationsResponse.prototype.setOrganizationsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.environment.Organization=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.environment.Organization}
 */
proto.bucketeer.account.GetMyOrganizationsResponse.prototype.addOrganizations =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.environment.Organization,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.account.GetMyOrganizationsResponse} returns this
 */
proto.bucketeer.account.GetMyOrganizationsResponse.prototype.clearOrganizationsList =
  function () {
    return this.setOrganizationsList([]);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.account.CreateAccountV2Request.repeatedFields_ = [7];

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
  proto.bucketeer.account.CreateAccountV2Request.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.CreateAccountV2Request.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.CreateAccountV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.CreateAccountV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        organizationId: jspb.Message.getFieldWithDefault(msg, 1, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.CreateAccountV2Command.toObject(
            includeInstance,
            f
          ),
        email: jspb.Message.getFieldWithDefault(msg, 3, ''),
        name: jspb.Message.getFieldWithDefault(msg, 4, ''),
        avatarImageUrl: jspb.Message.getFieldWithDefault(msg, 5, ''),
        organizationRole: jspb.Message.getFieldWithDefault(msg, 6, 0),
        environmentRolesList: jspb.Message.toObjectList(
          msg.getEnvironmentRolesList(),
          proto_account_account_pb.AccountV2.EnvironmentRole.toObject,
          includeInstance
        ),
        firstName: jspb.Message.getFieldWithDefault(msg, 8, ''),
        lastName: jspb.Message.getFieldWithDefault(msg, 9, ''),
        language: jspb.Message.getFieldWithDefault(msg, 10, '')
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
 * @return {!proto.bucketeer.account.CreateAccountV2Request}
 */
proto.bucketeer.account.CreateAccountV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.CreateAccountV2Request();
  return proto.bucketeer.account.CreateAccountV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.CreateAccountV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.CreateAccountV2Request}
 */
proto.bucketeer.account.CreateAccountV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 2:
          var value = new proto_account_command_pb.CreateAccountV2Command();
          reader.readMessage(
            value,
            proto_account_command_pb.CreateAccountV2Command
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setName(value);
          break;
        case 5:
          var value = /** @type {string} */ (reader.readString());
          msg.setAvatarImageUrl(value);
          break;
        case 6:
          var value =
            /** @type {!proto.bucketeer.account.AccountV2.Role.Organization} */ (
              reader.readEnum()
            );
          msg.setOrganizationRole(value);
          break;
        case 7:
          var value = new proto_account_account_pb.AccountV2.EnvironmentRole();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.EnvironmentRole
              .deserializeBinaryFromReader
          );
          msg.addEnvironmentRoles(value);
          break;
        case 8:
          var value = /** @type {string} */ (reader.readString());
          msg.setFirstName(value);
          break;
        case 9:
          var value = /** @type {string} */ (reader.readString());
          msg.setLastName(value);
          break;
        case 10:
          var value = /** @type {string} */ (reader.readString());
          msg.setLanguage(value);
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
proto.bucketeer.account.CreateAccountV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.CreateAccountV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.CreateAccountV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.CreateAccountV2Request.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_account_command_pb.CreateAccountV2Command.serializeBinaryToWriter
      );
    }
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getName();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
    f = message.getAvatarImageUrl();
    if (f.length > 0) {
      writer.writeString(5, f);
    }
    f = message.getOrganizationRole();
    if (f !== 0.0) {
      writer.writeEnum(6, f);
    }
    f = message.getEnvironmentRolesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        7,
        f,
        proto_account_account_pb.AccountV2.EnvironmentRole
          .serializeBinaryToWriter
      );
    }
    f = message.getFirstName();
    if (f.length > 0) {
      writer.writeString(8, f);
    }
    f = message.getLastName();
    if (f.length > 0) {
      writer.writeString(9, f);
    }
    f = message.getLanguage();
    if (f.length > 0) {
      writer.writeString(10, f);
    }
  };

/**
 * optional string organization_id = 1;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional CreateAccountV2Command command = 2;
 * @return {?proto.bucketeer.account.CreateAccountV2Command}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.CreateAccountV2Command} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.CreateAccountV2Command,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.CreateAccountV2Command|undefined} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 2, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional string email = 3;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setEmail = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 3, value);
};

/**
 * optional string name = 4;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setName = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 4, value);
};

/**
 * optional string avatar_image_url = 5;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getAvatarImageUrl =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 5, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setAvatarImageUrl =
  function (value) {
    return jspb.Message.setProto3StringField(this, 5, value);
  };

/**
 * optional AccountV2.Role.Organization organization_role = 6;
 * @return {!proto.bucketeer.account.AccountV2.Role.Organization}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getOrganizationRole =
  function () {
    return /** @type {!proto.bucketeer.account.AccountV2.Role.Organization} */ (
      jspb.Message.getFieldWithDefault(this, 6, 0)
    );
  };

/**
 * @param {!proto.bucketeer.account.AccountV2.Role.Organization} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setOrganizationRole =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 6, value);
  };

/**
 * repeated AccountV2.EnvironmentRole environment_roles = 7;
 * @return {!Array<!proto.bucketeer.account.AccountV2.EnvironmentRole>}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getEnvironmentRolesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.account.AccountV2.EnvironmentRole>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_account_account_pb.AccountV2.EnvironmentRole,
        7
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.account.AccountV2.EnvironmentRole>} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setEnvironmentRolesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 7, value);
  };

/**
 * @param {!proto.bucketeer.account.AccountV2.EnvironmentRole=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.account.AccountV2.EnvironmentRole}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.addEnvironmentRoles =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      7,
      opt_value,
      proto.bucketeer.account.AccountV2.EnvironmentRole,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.clearEnvironmentRolesList =
  function () {
    return this.setEnvironmentRolesList([]);
  };

/**
 * optional string first_name = 8;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getFirstName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 8, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setFirstName =
  function (value) {
    return jspb.Message.setProto3StringField(this, 8, value);
  };

/**
 * optional string last_name = 9;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getLastName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 9, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setLastName =
  function (value) {
    return jspb.Message.setProto3StringField(this, 9, value);
  };

/**
 * optional string language = 10;
 * @return {string}
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.getLanguage =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 10, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAccountV2Request} returns this
 */
proto.bucketeer.account.CreateAccountV2Request.prototype.setLanguage =
  function (value) {
    return jspb.Message.setProto3StringField(this, 10, value);
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
  proto.bucketeer.account.CreateAccountV2Response.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.CreateAccountV2Response.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.CreateAccountV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.CreateAccountV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        account:
          (f = msg.getAccount()) &&
          proto_account_account_pb.AccountV2.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.CreateAccountV2Response}
 */
proto.bucketeer.account.CreateAccountV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.CreateAccountV2Response();
  return proto.bucketeer.account.CreateAccountV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.CreateAccountV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.CreateAccountV2Response}
 */
proto.bucketeer.account.CreateAccountV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.setAccount(value);
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
proto.bucketeer.account.CreateAccountV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.CreateAccountV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.CreateAccountV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.CreateAccountV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccount();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AccountV2 account = 1;
 * @return {?proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.CreateAccountV2Response.prototype.getAccount =
  function () {
    return /** @type{?proto.bucketeer.account.AccountV2} */ (
      jspb.Message.getWrapperField(this, proto_account_account_pb.AccountV2, 1)
    );
  };

/**
 * @param {?proto.bucketeer.account.AccountV2|undefined} value
 * @return {!proto.bucketeer.account.CreateAccountV2Response} returns this
 */
proto.bucketeer.account.CreateAccountV2Response.prototype.setAccount =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.CreateAccountV2Response} returns this
 */
proto.bucketeer.account.CreateAccountV2Response.prototype.clearAccount =
  function () {
    return this.setAccount(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.CreateAccountV2Response.prototype.hasAccount =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.EnableAccountV2Request.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.EnableAccountV2Request.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.EnableAccountV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.EnableAccountV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.EnableAccountV2Command.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.EnableAccountV2Request}
 */
proto.bucketeer.account.EnableAccountV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.EnableAccountV2Request();
  return proto.bucketeer.account.EnableAccountV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.EnableAccountV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.EnableAccountV2Request}
 */
proto.bucketeer.account.EnableAccountV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = new proto_account_command_pb.EnableAccountV2Command();
          reader.readMessage(
            value,
            proto_account_command_pb.EnableAccountV2Command
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
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
proto.bucketeer.account.EnableAccountV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.EnableAccountV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.EnableAccountV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnableAccountV2Request.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_account_command_pb.EnableAccountV2Command.serializeBinaryToWriter
      );
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.EnableAccountV2Request} returns this
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.setEmail = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.EnableAccountV2Request} returns this
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional EnableAccountV2Command command = 3;
 * @return {?proto.bucketeer.account.EnableAccountV2Command}
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.EnableAccountV2Command} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.EnableAccountV2Command,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.EnableAccountV2Command|undefined} value
 * @return {!proto.bucketeer.account.EnableAccountV2Request} returns this
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 3, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.EnableAccountV2Request} returns this
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.EnableAccountV2Request.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
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
  proto.bucketeer.account.EnableAccountV2Response.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.EnableAccountV2Response.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.EnableAccountV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.EnableAccountV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        account:
          (f = msg.getAccount()) &&
          proto_account_account_pb.AccountV2.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.EnableAccountV2Response}
 */
proto.bucketeer.account.EnableAccountV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.EnableAccountV2Response();
  return proto.bucketeer.account.EnableAccountV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.EnableAccountV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.EnableAccountV2Response}
 */
proto.bucketeer.account.EnableAccountV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.setAccount(value);
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
proto.bucketeer.account.EnableAccountV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.EnableAccountV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.EnableAccountV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnableAccountV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccount();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AccountV2 account = 1;
 * @return {?proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.EnableAccountV2Response.prototype.getAccount =
  function () {
    return /** @type{?proto.bucketeer.account.AccountV2} */ (
      jspb.Message.getWrapperField(this, proto_account_account_pb.AccountV2, 1)
    );
  };

/**
 * @param {?proto.bucketeer.account.AccountV2|undefined} value
 * @return {!proto.bucketeer.account.EnableAccountV2Response} returns this
 */
proto.bucketeer.account.EnableAccountV2Response.prototype.setAccount =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.EnableAccountV2Response} returns this
 */
proto.bucketeer.account.EnableAccountV2Response.prototype.clearAccount =
  function () {
    return this.setAccount(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.EnableAccountV2Response.prototype.hasAccount =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.DisableAccountV2Request.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.DisableAccountV2Request.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DisableAccountV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DisableAccountV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.DisableAccountV2Command.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.DisableAccountV2Request}
 */
proto.bucketeer.account.DisableAccountV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DisableAccountV2Request();
  return proto.bucketeer.account.DisableAccountV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DisableAccountV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DisableAccountV2Request}
 */
proto.bucketeer.account.DisableAccountV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = new proto_account_command_pb.DisableAccountV2Command();
          reader.readMessage(
            value,
            proto_account_command_pb.DisableAccountV2Command
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
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
proto.bucketeer.account.DisableAccountV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DisableAccountV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DisableAccountV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DisableAccountV2Request.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_account_command_pb.DisableAccountV2Command.serializeBinaryToWriter
      );
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DisableAccountV2Request} returns this
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.setEmail = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DisableAccountV2Request} returns this
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional DisableAccountV2Command command = 3;
 * @return {?proto.bucketeer.account.DisableAccountV2Command}
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.DisableAccountV2Command} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.DisableAccountV2Command,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.DisableAccountV2Command|undefined} value
 * @return {!proto.bucketeer.account.DisableAccountV2Request} returns this
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.DisableAccountV2Request} returns this
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.DisableAccountV2Request.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
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
  proto.bucketeer.account.DisableAccountV2Response.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.DisableAccountV2Response.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DisableAccountV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DisableAccountV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        account:
          (f = msg.getAccount()) &&
          proto_account_account_pb.AccountV2.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.DisableAccountV2Response}
 */
proto.bucketeer.account.DisableAccountV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DisableAccountV2Response();
  return proto.bucketeer.account.DisableAccountV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DisableAccountV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DisableAccountV2Response}
 */
proto.bucketeer.account.DisableAccountV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.setAccount(value);
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
proto.bucketeer.account.DisableAccountV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DisableAccountV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DisableAccountV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DisableAccountV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccount();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AccountV2 account = 1;
 * @return {?proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.DisableAccountV2Response.prototype.getAccount =
  function () {
    return /** @type{?proto.bucketeer.account.AccountV2} */ (
      jspb.Message.getWrapperField(this, proto_account_account_pb.AccountV2, 1)
    );
  };

/**
 * @param {?proto.bucketeer.account.AccountV2|undefined} value
 * @return {!proto.bucketeer.account.DisableAccountV2Response} returns this
 */
proto.bucketeer.account.DisableAccountV2Response.prototype.setAccount =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.DisableAccountV2Response} returns this
 */
proto.bucketeer.account.DisableAccountV2Response.prototype.clearAccount =
  function () {
    return this.setAccount(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.DisableAccountV2Response.prototype.hasAccount =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.DeleteAccountV2Request.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.DeleteAccountV2Request.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DeleteAccountV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DeleteAccountV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.DeleteAccountV2Command.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.DeleteAccountV2Request}
 */
proto.bucketeer.account.DeleteAccountV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DeleteAccountV2Request();
  return proto.bucketeer.account.DeleteAccountV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DeleteAccountV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DeleteAccountV2Request}
 */
proto.bucketeer.account.DeleteAccountV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = new proto_account_command_pb.DeleteAccountV2Command();
          reader.readMessage(
            value,
            proto_account_command_pb.DeleteAccountV2Command
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
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
proto.bucketeer.account.DeleteAccountV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DeleteAccountV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DeleteAccountV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DeleteAccountV2Request.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_account_command_pb.DeleteAccountV2Command.serializeBinaryToWriter
      );
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DeleteAccountV2Request} returns this
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.setEmail = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DeleteAccountV2Request} returns this
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional DeleteAccountV2Command command = 3;
 * @return {?proto.bucketeer.account.DeleteAccountV2Command}
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.DeleteAccountV2Command} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.DeleteAccountV2Command,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.DeleteAccountV2Command|undefined} value
 * @return {!proto.bucketeer.account.DeleteAccountV2Request} returns this
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 3, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.DeleteAccountV2Request} returns this
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.DeleteAccountV2Request.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
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
  proto.bucketeer.account.DeleteAccountV2Response.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.DeleteAccountV2Response.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DeleteAccountV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DeleteAccountV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.DeleteAccountV2Response}
 */
proto.bucketeer.account.DeleteAccountV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DeleteAccountV2Response();
  return proto.bucketeer.account.DeleteAccountV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DeleteAccountV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DeleteAccountV2Response}
 */
proto.bucketeer.account.DeleteAccountV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.DeleteAccountV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DeleteAccountV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DeleteAccountV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DeleteAccountV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.account.UpdateAccountV2Request.repeatedFields_ = [15];

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
  proto.bucketeer.account.UpdateAccountV2Request.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.UpdateAccountV2Request.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateAccountV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateAccountV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        changeNameCommand:
          (f = msg.getChangeNameCommand()) &&
          proto_account_command_pb.ChangeAccountV2NameCommand.toObject(
            includeInstance,
            f
          ),
        changeAvatarUrlCommand:
          (f = msg.getChangeAvatarUrlCommand()) &&
          proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand.toObject(
            includeInstance,
            f
          ),
        changeOrganizationRoleCommand:
          (f = msg.getChangeOrganizationRoleCommand()) &&
          proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand.toObject(
            includeInstance,
            f
          ),
        changeEnvironmentRolesCommand:
          (f = msg.getChangeEnvironmentRolesCommand()) &&
          proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand.toObject(
            includeInstance,
            f
          ),
        changeFirstNameCommand:
          (f = msg.getChangeFirstNameCommand()) &&
          proto_account_command_pb.ChangeAccountV2FirstNameCommand.toObject(
            includeInstance,
            f
          ),
        changeLastNameCommand:
          (f = msg.getChangeLastNameCommand()) &&
          proto_account_command_pb.ChangeAccountV2LastNameCommand.toObject(
            includeInstance,
            f
          ),
        changeLanguageCommand:
          (f = msg.getChangeLanguageCommand()) &&
          proto_account_command_pb.ChangeAccountV2LanguageCommand.toObject(
            includeInstance,
            f
          ),
        changeLastSeenCommand:
          (f = msg.getChangeLastSeenCommand()) &&
          proto_account_command_pb.ChangeAccountV2LastSeenCommand.toObject(
            includeInstance,
            f
          ),
        changeAvatarCommand:
          (f = msg.getChangeAvatarCommand()) &&
          proto_account_command_pb.ChangeAccountV2AvatarCommand.toObject(
            includeInstance,
            f
          ),
        name:
          (f = msg.getName()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        avatarImageUrl:
          (f = msg.getAvatarImageUrl()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        organizationRole:
          (f = msg.getOrganizationRole()) &&
          proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.toObject(
            includeInstance,
            f
          ),
        environmentRolesList: jspb.Message.toObjectList(
          msg.getEnvironmentRolesList(),
          proto_account_account_pb.AccountV2.EnvironmentRole.toObject,
          includeInstance
        ),
        firstName:
          (f = msg.getFirstName()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        lastName:
          (f = msg.getLastName()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        language:
          (f = msg.getLanguage()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        lastSeen:
          (f = msg.getLastSeen()) &&
          google_protobuf_wrappers_pb.Int64Value.toObject(includeInstance, f),
        avatar:
          (f = msg.getAvatar()) &&
          proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.UpdateAccountV2Request}
 */
proto.bucketeer.account.UpdateAccountV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.UpdateAccountV2Request();
  return proto.bucketeer.account.UpdateAccountV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateAccountV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request}
 */
proto.bucketeer.account.UpdateAccountV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = new proto_account_command_pb.ChangeAccountV2NameCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2NameCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeNameCommand(value);
          break;
        case 4:
          var value =
            new proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeAvatarUrlCommand(value);
          break;
        case 5:
          var value =
            new proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeOrganizationRoleCommand(value);
          break;
        case 6:
          var value =
            new proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeEnvironmentRolesCommand(value);
          break;
        case 7:
          var value =
            new proto_account_command_pb.ChangeAccountV2FirstNameCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2FirstNameCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeFirstNameCommand(value);
          break;
        case 8:
          var value =
            new proto_account_command_pb.ChangeAccountV2LastNameCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2LastNameCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeLastNameCommand(value);
          break;
        case 9:
          var value =
            new proto_account_command_pb.ChangeAccountV2LanguageCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2LanguageCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeLanguageCommand(value);
          break;
        case 10:
          var value =
            new proto_account_command_pb.ChangeAccountV2LastSeenCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2LastSeenCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeLastSeenCommand(value);
          break;
        case 11:
          var value =
            new proto_account_command_pb.ChangeAccountV2AvatarCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAccountV2AvatarCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeAvatarCommand(value);
          break;
        case 12:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setName(value);
          break;
        case 13:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setAvatarImageUrl(value);
          break;
        case 14:
          var value =
            new proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue();
          reader.readMessage(
            value,
            proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue
              .deserializeBinaryFromReader
          );
          msg.setOrganizationRole(value);
          break;
        case 15:
          var value = new proto_account_account_pb.AccountV2.EnvironmentRole();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.EnvironmentRole
              .deserializeBinaryFromReader
          );
          msg.addEnvironmentRoles(value);
          break;
        case 16:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setFirstName(value);
          break;
        case 17:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setLastName(value);
          break;
        case 18:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setLanguage(value);
          break;
        case 19:
          var value = new google_protobuf_wrappers_pb.Int64Value();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.Int64Value.deserializeBinaryFromReader
          );
          msg.setLastSeen(value);
          break;
        case 20:
          var value =
            new proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar();
          reader.readMessage(
            value,
            proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar
              .deserializeBinaryFromReader
          );
          msg.setAvatar(value);
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
proto.bucketeer.account.UpdateAccountV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateAccountV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateAccountV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateAccountV2Request.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getChangeNameCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_account_command_pb.ChangeAccountV2NameCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeAvatarUrlCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeOrganizationRoleCommand();
    if (f != null) {
      writer.writeMessage(
        5,
        f,
        proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeEnvironmentRolesCommand();
    if (f != null) {
      writer.writeMessage(
        6,
        f,
        proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeFirstNameCommand();
    if (f != null) {
      writer.writeMessage(
        7,
        f,
        proto_account_command_pb.ChangeAccountV2FirstNameCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeLastNameCommand();
    if (f != null) {
      writer.writeMessage(
        8,
        f,
        proto_account_command_pb.ChangeAccountV2LastNameCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeLanguageCommand();
    if (f != null) {
      writer.writeMessage(
        9,
        f,
        proto_account_command_pb.ChangeAccountV2LanguageCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeLastSeenCommand();
    if (f != null) {
      writer.writeMessage(
        10,
        f,
        proto_account_command_pb.ChangeAccountV2LastSeenCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeAvatarCommand();
    if (f != null) {
      writer.writeMessage(
        11,
        f,
        proto_account_command_pb.ChangeAccountV2AvatarCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getName();
    if (f != null) {
      writer.writeMessage(
        12,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getAvatarImageUrl();
    if (f != null) {
      writer.writeMessage(
        13,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getOrganizationRole();
    if (f != null) {
      writer.writeMessage(
        14,
        f,
        proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentRolesList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        15,
        f,
        proto_account_account_pb.AccountV2.EnvironmentRole
          .serializeBinaryToWriter
      );
    }
    f = message.getFirstName();
    if (f != null) {
      writer.writeMessage(
        16,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getLastName();
    if (f != null) {
      writer.writeMessage(
        17,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getLanguage();
    if (f != null) {
      writer.writeMessage(
        18,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getLastSeen();
    if (f != null) {
      writer.writeMessage(
        19,
        f,
        google_protobuf_wrappers_pb.Int64Value.serializeBinaryToWriter
      );
    }
    f = message.getAvatar();
    if (f != null) {
      writer.writeMessage(
        20,
        f,
        proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar
          .serializeBinaryToWriter
      );
    }
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
  proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          avatarImage: msg.getAvatarImage_asB64(),
          avatarFileType: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar}
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar();
    return proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar}
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {!Uint8Array} */ (reader.readBytes());
          msg.setAvatarImage(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setAvatarFileType(value);
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
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAvatarImage_asU8();
    if (f.length > 0) {
      writer.writeBytes(1, f);
    }
    f = message.getAvatarFileType();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
  };

/**
 * optional bytes avatar_image = 1;
 * @return {!(string|Uint8Array)}
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.getAvatarImage =
  function () {
    return /** @type {!(string|Uint8Array)} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * optional bytes avatar_image = 1;
 * This is a type-conversion wrapper around `getAvatarImage()`
 * @return {string}
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.getAvatarImage_asB64 =
  function () {
    return /** @type {string} */ (
      jspb.Message.bytesAsB64(this.getAvatarImage())
    );
  };

/**
 * optional bytes avatar_image = 1;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getAvatarImage()`
 * @return {!Uint8Array}
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.getAvatarImage_asU8 =
  function () {
    return /** @type {!Uint8Array} */ (
      jspb.Message.bytesAsU8(this.getAvatarImage())
    );
  };

/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.setAvatarImage =
  function (value) {
    return jspb.Message.setProto3BytesField(this, 1, value);
  };

/**
 * optional string avatar_file_type = 2;
 * @return {string}
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.getAvatarFileType =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar.prototype.setAvatarFileType =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
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
  proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          role: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue();
    return proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value =
            /** @type {!proto.bucketeer.account.AccountV2.Role.Organization} */ (
              reader.readEnum()
            );
          msg.setRole(value);
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
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getRole();
    if (f !== 0.0) {
      writer.writeEnum(1, f);
    }
  };

/**
 * optional AccountV2.Role.Organization role = 1;
 * @return {!proto.bucketeer.account.AccountV2.Role.Organization}
 */
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.prototype.getRole =
  function () {
    return /** @type {!proto.bucketeer.account.AccountV2.Role.Organization} */ (
      jspb.Message.getFieldWithDefault(this, 1, 0)
    );
  };

/**
 * @param {!proto.bucketeer.account.AccountV2.Role.Organization} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue.prototype.setRole =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 1, value);
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setEmail = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional ChangeAccountV2NameCommand change_name_command = 3;
 * @return {?proto.bucketeer.account.ChangeAccountV2NameCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeNameCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2NameCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2NameCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2NameCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeNameCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeNameCommand =
  function () {
    return this.setChangeNameCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeNameCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional ChangeAccountV2AvatarImageUrlCommand change_avatar_url_command = 4;
 * @return {?proto.bucketeer.account.ChangeAccountV2AvatarImageUrlCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeAvatarUrlCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2AvatarImageUrlCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2AvatarImageUrlCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2AvatarImageUrlCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeAvatarUrlCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeAvatarUrlCommand =
  function () {
    return this.setChangeAvatarUrlCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeAvatarUrlCommand =
  function () {
    return jspb.Message.getField(this, 4) != null;
  };

/**
 * optional ChangeAccountV2OrganizationRoleCommand change_organization_role_command = 5;
 * @return {?proto.bucketeer.account.ChangeAccountV2OrganizationRoleCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeOrganizationRoleCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2OrganizationRoleCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2OrganizationRoleCommand,
        5
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2OrganizationRoleCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeOrganizationRoleCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 5, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeOrganizationRoleCommand =
  function () {
    return this.setChangeOrganizationRoleCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeOrganizationRoleCommand =
  function () {
    return jspb.Message.getField(this, 5) != null;
  };

/**
 * optional ChangeAccountV2EnvironmentRolesCommand change_environment_roles_command = 6;
 * @return {?proto.bucketeer.account.ChangeAccountV2EnvironmentRolesCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeEnvironmentRolesCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2EnvironmentRolesCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2EnvironmentRolesCommand,
        6
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2EnvironmentRolesCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeEnvironmentRolesCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 6, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeEnvironmentRolesCommand =
  function () {
    return this.setChangeEnvironmentRolesCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeEnvironmentRolesCommand =
  function () {
    return jspb.Message.getField(this, 6) != null;
  };

/**
 * optional ChangeAccountV2FirstNameCommand change_first_name_command = 7;
 * @return {?proto.bucketeer.account.ChangeAccountV2FirstNameCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeFirstNameCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2FirstNameCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2FirstNameCommand,
        7
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2FirstNameCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeFirstNameCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 7, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeFirstNameCommand =
  function () {
    return this.setChangeFirstNameCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeFirstNameCommand =
  function () {
    return jspb.Message.getField(this, 7) != null;
  };

/**
 * optional ChangeAccountV2LastNameCommand change_last_name_command = 8;
 * @return {?proto.bucketeer.account.ChangeAccountV2LastNameCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeLastNameCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2LastNameCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2LastNameCommand,
        8
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2LastNameCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeLastNameCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 8, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeLastNameCommand =
  function () {
    return this.setChangeLastNameCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeLastNameCommand =
  function () {
    return jspb.Message.getField(this, 8) != null;
  };

/**
 * optional ChangeAccountV2LanguageCommand change_language_command = 9;
 * @return {?proto.bucketeer.account.ChangeAccountV2LanguageCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeLanguageCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2LanguageCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2LanguageCommand,
        9
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2LanguageCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeLanguageCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 9, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeLanguageCommand =
  function () {
    return this.setChangeLanguageCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeLanguageCommand =
  function () {
    return jspb.Message.getField(this, 9) != null;
  };

/**
 * optional ChangeAccountV2LastSeenCommand change_last_seen_command = 10;
 * @return {?proto.bucketeer.account.ChangeAccountV2LastSeenCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeLastSeenCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2LastSeenCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2LastSeenCommand,
        10
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2LastSeenCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeLastSeenCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 10, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeLastSeenCommand =
  function () {
    return this.setChangeLastSeenCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeLastSeenCommand =
  function () {
    return jspb.Message.getField(this, 10) != null;
  };

/**
 * optional ChangeAccountV2AvatarCommand change_avatar_command = 11;
 * @return {?proto.bucketeer.account.ChangeAccountV2AvatarCommand}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getChangeAvatarCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAccountV2AvatarCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAccountV2AvatarCommand,
        11
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAccountV2AvatarCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setChangeAvatarCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 11, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearChangeAvatarCommand =
  function () {
    return this.setChangeAvatarCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasChangeAvatarCommand =
  function () {
    return jspb.Message.getField(this, 11) != null;
  };

/**
 * optional google.protobuf.StringValue name = 12;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getName = function () {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(
      this,
      google_protobuf_wrappers_pb.StringValue,
      12
    )
  );
};

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setName = function (
  value
) {
  return jspb.Message.setWrapperField(this, 12, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearName =
  function () {
    return this.setName(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasName = function () {
  return jspb.Message.getField(this, 12) != null;
};

/**
 * optional google.protobuf.StringValue avatar_image_url = 13;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getAvatarImageUrl =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        13
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setAvatarImageUrl =
  function (value) {
    return jspb.Message.setWrapperField(this, 13, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearAvatarImageUrl =
  function () {
    return this.setAvatarImageUrl(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasAvatarImageUrl =
  function () {
    return jspb.Message.getField(this, 13) != null;
  };

/**
 * optional OrganizationRoleValue organization_role = 14;
 * @return {?proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getOrganizationRole =
  function () {
    return /** @type{?proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue} */ (
      jspb.Message.getWrapperField(
        this,
        proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue,
        14
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.UpdateAccountV2Request.OrganizationRoleValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setOrganizationRole =
  function (value) {
    return jspb.Message.setWrapperField(this, 14, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearOrganizationRole =
  function () {
    return this.setOrganizationRole(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasOrganizationRole =
  function () {
    return jspb.Message.getField(this, 14) != null;
  };

/**
 * repeated AccountV2.EnvironmentRole environment_roles = 15;
 * @return {!Array<!proto.bucketeer.account.AccountV2.EnvironmentRole>}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getEnvironmentRolesList =
  function () {
    return /** @type{!Array<!proto.bucketeer.account.AccountV2.EnvironmentRole>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_account_account_pb.AccountV2.EnvironmentRole,
        15
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.account.AccountV2.EnvironmentRole>} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setEnvironmentRolesList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 15, value);
  };

/**
 * @param {!proto.bucketeer.account.AccountV2.EnvironmentRole=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.account.AccountV2.EnvironmentRole}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.addEnvironmentRoles =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      15,
      opt_value,
      proto.bucketeer.account.AccountV2.EnvironmentRole,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearEnvironmentRolesList =
  function () {
    return this.setEnvironmentRolesList([]);
  };

/**
 * optional google.protobuf.StringValue first_name = 16;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getFirstName =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        16
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setFirstName =
  function (value) {
    return jspb.Message.setWrapperField(this, 16, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearFirstName =
  function () {
    return this.setFirstName(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasFirstName =
  function () {
    return jspb.Message.getField(this, 16) != null;
  };

/**
 * optional google.protobuf.StringValue last_name = 17;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getLastName =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        17
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setLastName =
  function (value) {
    return jspb.Message.setWrapperField(this, 17, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearLastName =
  function () {
    return this.setLastName(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasLastName =
  function () {
    return jspb.Message.getField(this, 17) != null;
  };

/**
 * optional google.protobuf.StringValue language = 18;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getLanguage =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        18
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setLanguage =
  function (value) {
    return jspb.Message.setWrapperField(this, 18, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearLanguage =
  function () {
    return this.setLanguage(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasLanguage =
  function () {
    return jspb.Message.getField(this, 18) != null;
  };

/**
 * optional google.protobuf.Int64Value last_seen = 19;
 * @return {?proto.google.protobuf.Int64Value}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getLastSeen =
  function () {
    return /** @type{?proto.google.protobuf.Int64Value} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.Int64Value,
        19
      )
    );
  };

/**
 * @param {?proto.google.protobuf.Int64Value|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setLastSeen =
  function (value) {
    return jspb.Message.setWrapperField(this, 19, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearLastSeen =
  function () {
    return this.setLastSeen(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasLastSeen =
  function () {
    return jspb.Message.getField(this, 19) != null;
  };

/**
 * optional AccountV2Avatar avatar = 20;
 * @return {?proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.getAvatar =
  function () {
    return /** @type{?proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar} */ (
      jspb.Message.getWrapperField(
        this,
        proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar,
        20
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.UpdateAccountV2Request.AccountV2Avatar|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.setAvatar = function (
  value
) {
  return jspb.Message.setWrapperField(this, 20, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Request} returns this
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.clearAvatar =
  function () {
    return this.setAvatar(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Request.prototype.hasAvatar =
  function () {
    return jspb.Message.getField(this, 20) != null;
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
  proto.bucketeer.account.UpdateAccountV2Response.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.UpdateAccountV2Response.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateAccountV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateAccountV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        account:
          (f = msg.getAccount()) &&
          proto_account_account_pb.AccountV2.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.UpdateAccountV2Response}
 */
proto.bucketeer.account.UpdateAccountV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.UpdateAccountV2Response();
  return proto.bucketeer.account.UpdateAccountV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateAccountV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateAccountV2Response}
 */
proto.bucketeer.account.UpdateAccountV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.setAccount(value);
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
proto.bucketeer.account.UpdateAccountV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateAccountV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateAccountV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateAccountV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccount();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AccountV2 account = 1;
 * @return {?proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.UpdateAccountV2Response.prototype.getAccount =
  function () {
    return /** @type{?proto.bucketeer.account.AccountV2} */ (
      jspb.Message.getWrapperField(this, proto_account_account_pb.AccountV2, 1)
    );
  };

/**
 * @param {?proto.bucketeer.account.AccountV2|undefined} value
 * @return {!proto.bucketeer.account.UpdateAccountV2Response} returns this
 */
proto.bucketeer.account.UpdateAccountV2Response.prototype.setAccount =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAccountV2Response} returns this
 */
proto.bucketeer.account.UpdateAccountV2Response.prototype.clearAccount =
  function () {
    return this.setAccount(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAccountV2Response.prototype.hasAccount =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.GetAccountV2Request.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.GetAccountV2Request.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAccountV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAccountV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.account.GetAccountV2Request}
 */
proto.bucketeer.account.GetAccountV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetAccountV2Request();
  return proto.bucketeer.account.GetAccountV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAccountV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAccountV2Request}
 */
proto.bucketeer.account.GetAccountV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
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
proto.bucketeer.account.GetAccountV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAccountV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAccountV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAccountV2Request.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getOrganizationId();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
};

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.GetAccountV2Request.prototype.getEmail = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAccountV2Request} returns this
 */
proto.bucketeer.account.GetAccountV2Request.prototype.setEmail = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.GetAccountV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAccountV2Request} returns this
 */
proto.bucketeer.account.GetAccountV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
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
  proto.bucketeer.account.GetAccountV2Response.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.GetAccountV2Response.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAccountV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAccountV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        account:
          (f = msg.getAccount()) &&
          proto_account_account_pb.AccountV2.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.GetAccountV2Response}
 */
proto.bucketeer.account.GetAccountV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetAccountV2Response();
  return proto.bucketeer.account.GetAccountV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAccountV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAccountV2Response}
 */
proto.bucketeer.account.GetAccountV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.setAccount(value);
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
proto.bucketeer.account.GetAccountV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAccountV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAccountV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAccountV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccount();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AccountV2 account = 1;
 * @return {?proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.GetAccountV2Response.prototype.getAccount =
  function () {
    return /** @type{?proto.bucketeer.account.AccountV2} */ (
      jspb.Message.getWrapperField(this, proto_account_account_pb.AccountV2, 1)
    );
  };

/**
 * @param {?proto.bucketeer.account.AccountV2|undefined} value
 * @return {!proto.bucketeer.account.GetAccountV2Response} returns this
 */
proto.bucketeer.account.GetAccountV2Response.prototype.setAccount = function (
  value
) {
  return jspb.Message.setWrapperField(this, 1, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.GetAccountV2Response} returns this
 */
proto.bucketeer.account.GetAccountV2Response.prototype.clearAccount =
  function () {
    return this.setAccount(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.GetAccountV2Response.prototype.hasAccount =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          email: jspb.Message.getFieldWithDefault(msg, 1, ''),
          environmentId: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest();
    return proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest} returns this
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.prototype.setEmail =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional string environment_id = 2;
 * @return {string}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest} returns this
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
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
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          account:
            (f = msg.getAccount()) &&
            proto_account_account_pb.AccountV2.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse();
    return proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.setAccount(value);
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
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccount();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
  };

/**
 * optional AccountV2 account = 1;
 * @return {?proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.prototype.getAccount =
  function () {
    return /** @type{?proto.bucketeer.account.AccountV2} */ (
      jspb.Message.getWrapperField(this, proto_account_account_pb.AccountV2, 1)
    );
  };

/**
 * @param {?proto.bucketeer.account.AccountV2|undefined} value
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse} returns this
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.prototype.setAccount =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse} returns this
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.prototype.clearAccount =
  function () {
    return this.setAccount(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.GetAccountV2ByEnvironmentIDResponse.prototype.hasAccount =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.ListAccountsV2Request.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.ListAccountsV2Request.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.ListAccountsV2Request} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.ListAccountsV2Request.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
        cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        orderBy: jspb.Message.getFieldWithDefault(msg, 4, 0),
        orderDirection: jspb.Message.getFieldWithDefault(msg, 5, 0),
        searchKeyword: jspb.Message.getFieldWithDefault(msg, 6, ''),
        disabled:
          (f = msg.getDisabled()) &&
          google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
        organizationRole:
          (f = msg.getOrganizationRole()) &&
          google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f),
        environmentId:
          (f = msg.getEnvironmentId()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        environmentRole:
          (f = msg.getEnvironmentRole()) &&
          google_protobuf_wrappers_pb.Int32Value.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.ListAccountsV2Request}
 */
proto.bucketeer.account.ListAccountsV2Request.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.ListAccountsV2Request();
  return proto.bucketeer.account.ListAccountsV2Request.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.ListAccountsV2Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.ListAccountsV2Request}
 */
proto.bucketeer.account.ListAccountsV2Request.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 4:
          var value =
            /** @type {!proto.bucketeer.account.ListAccountsV2Request.OrderBy} */ (
              reader.readEnum()
            );
          msg.setOrderBy(value);
          break;
        case 5:
          var value =
            /** @type {!proto.bucketeer.account.ListAccountsV2Request.OrderDirection} */ (
              reader.readEnum()
            );
          msg.setOrderDirection(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setSearchKeyword(value);
          break;
        case 7:
          var value = new google_protobuf_wrappers_pb.BoolValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader
          );
          msg.setDisabled(value);
          break;
        case 8:
          var value = new google_protobuf_wrappers_pb.Int32Value();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.Int32Value.deserializeBinaryFromReader
          );
          msg.setOrganizationRole(value);
          break;
        case 9:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setEnvironmentId(value);
          break;
        case 10:
          var value = new google_protobuf_wrappers_pb.Int32Value();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.Int32Value.deserializeBinaryFromReader
          );
          msg.setEnvironmentRole(value);
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
proto.bucketeer.account.ListAccountsV2Request.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.ListAccountsV2Request.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.ListAccountsV2Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.ListAccountsV2Request.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getPageSize();
    if (f !== 0) {
      writer.writeInt64(1, f);
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getOrderBy();
    if (f !== 0.0) {
      writer.writeEnum(4, f);
    }
    f = message.getOrderDirection();
    if (f !== 0.0) {
      writer.writeEnum(5, f);
    }
    f = message.getSearchKeyword();
    if (f.length > 0) {
      writer.writeString(6, f);
    }
    f = message.getDisabled();
    if (f != null) {
      writer.writeMessage(
        7,
        f,
        google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
      );
    }
    f = message.getOrganizationRole();
    if (f != null) {
      writer.writeMessage(
        8,
        f,
        google_protobuf_wrappers_pb.Int32Value.serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f != null) {
      writer.writeMessage(
        9,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentRole();
    if (f != null) {
      writer.writeMessage(
        10,
        f,
        google_protobuf_wrappers_pb.Int32Value.serializeBinaryToWriter
      );
    }
  };

/**
 * @enum {number}
 */
proto.bucketeer.account.ListAccountsV2Request.OrderBy = {
  DEFAULT: 0,
  EMAIL: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3,
  ORGANIZATION_ROLE: 4,
  ENVIRONMENT_COUNT: 5,
  LAST_SEEN: 6,
  STATE: 7
};

/**
 * @enum {number}
 */
proto.bucketeer.account.ListAccountsV2Request.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setPageSize = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 1, value);
};

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string organization_id = 3;
 * @return {string}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.account.ListAccountsV2Request.OrderBy}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getOrderBy =
  function () {
    return /** @type {!proto.bucketeer.account.ListAccountsV2Request.OrderBy} */ (
      jspb.Message.getFieldWithDefault(this, 4, 0)
    );
  };

/**
 * @param {!proto.bucketeer.account.ListAccountsV2Request.OrderBy} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setOrderBy = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};

/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.account.ListAccountsV2Request.OrderDirection}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getOrderDirection =
  function () {
    return /** @type {!proto.bucketeer.account.ListAccountsV2Request.OrderDirection} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.account.ListAccountsV2Request.OrderDirection} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setOrderDirection =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
  };

/**
 * optional string search_keyword = 6;
 * @return {string}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getSearchKeyword =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setSearchKeyword =
  function (value) {
    return jspb.Message.setProto3StringField(this, 6, value);
  };

/**
 * optional google.protobuf.BoolValue disabled = 7;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getDisabled =
  function () {
    return /** @type{?proto.google.protobuf.BoolValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.BoolValue,
        7
      )
    );
  };

/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setDisabled = function (
  value
) {
  return jspb.Message.setWrapperField(this, 7, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.clearDisabled =
  function () {
    return this.setDisabled(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.hasDisabled =
  function () {
    return jspb.Message.getField(this, 7) != null;
  };

/**
 * optional google.protobuf.Int32Value organization_role = 8;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getOrganizationRole =
  function () {
    return /** @type{?proto.google.protobuf.Int32Value} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.Int32Value,
        8
      )
    );
  };

/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setOrganizationRole =
  function (value) {
    return jspb.Message.setWrapperField(this, 8, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.clearOrganizationRole =
  function () {
    return this.setOrganizationRole(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.hasOrganizationRole =
  function () {
    return jspb.Message.getField(this, 8) != null;
  };

/**
 * optional google.protobuf.StringValue environment_id = 9;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getEnvironmentId =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        9
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setWrapperField(this, 9, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.clearEnvironmentId =
  function () {
    return this.setEnvironmentId(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.hasEnvironmentId =
  function () {
    return jspb.Message.getField(this, 9) != null;
  };

/**
 * optional google.protobuf.Int32Value environment_role = 10;
 * @return {?proto.google.protobuf.Int32Value}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.getEnvironmentRole =
  function () {
    return /** @type{?proto.google.protobuf.Int32Value} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.Int32Value,
        10
      )
    );
  };

/**
 * @param {?proto.google.protobuf.Int32Value|undefined} value
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.setEnvironmentRole =
  function (value) {
    return jspb.Message.setWrapperField(this, 10, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.ListAccountsV2Request} returns this
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.clearEnvironmentRole =
  function () {
    return this.setEnvironmentRole(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.ListAccountsV2Request.prototype.hasEnvironmentRole =
  function () {
    return jspb.Message.getField(this, 10) != null;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.account.ListAccountsV2Response.repeatedFields_ = [1];

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
  proto.bucketeer.account.ListAccountsV2Response.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.ListAccountsV2Response.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.ListAccountsV2Response} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.ListAccountsV2Response.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        accountsList: jspb.Message.toObjectList(
          msg.getAccountsList(),
          proto_account_account_pb.AccountV2.toObject,
          includeInstance
        ),
        cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
        totalCount: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.account.ListAccountsV2Response}
 */
proto.bucketeer.account.ListAccountsV2Response.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.ListAccountsV2Response();
  return proto.bucketeer.account.ListAccountsV2Response.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.ListAccountsV2Response} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.ListAccountsV2Response}
 */
proto.bucketeer.account.ListAccountsV2Response.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_account_pb.AccountV2();
          reader.readMessage(
            value,
            proto_account_account_pb.AccountV2.deserializeBinaryFromReader
          );
          msg.addAccounts(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 3:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setTotalCount(value);
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
proto.bucketeer.account.ListAccountsV2Response.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.ListAccountsV2Response.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.ListAccountsV2Response} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.ListAccountsV2Response.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getAccountsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_account_account_pb.AccountV2.serializeBinaryToWriter
      );
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getTotalCount();
    if (f !== 0) {
      writer.writeInt64(3, f);
    }
  };

/**
 * repeated AccountV2 accounts = 1;
 * @return {!Array<!proto.bucketeer.account.AccountV2>}
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.getAccountsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.account.AccountV2>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_account_account_pb.AccountV2,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.account.AccountV2>} value
 * @return {!proto.bucketeer.account.ListAccountsV2Response} returns this
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.setAccountsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.account.AccountV2=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.account.AccountV2}
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.addAccounts =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.account.AccountV2,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.account.ListAccountsV2Response} returns this
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.clearAccountsList =
  function () {
    return this.setAccountsList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAccountsV2Response} returns this
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.getTotalCount =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.account.ListAccountsV2Response} returns this
 */
proto.bucketeer.account.ListAccountsV2Response.prototype.setTotalCount =
  function (value) {
    return jspb.Message.setProto3IntField(this, 3, value);
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
  proto.bucketeer.account.CreateAPIKeyRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.CreateAPIKeyRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.CreateAPIKeyRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.CreateAPIKeyRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.CreateAPIKeyCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        name: jspb.Message.getFieldWithDefault(msg, 4, ''),
        role: jspb.Message.getFieldWithDefault(msg, 5, 0),
        maintainer: jspb.Message.getFieldWithDefault(msg, 6, ''),
        description: jspb.Message.getFieldWithDefault(msg, 7, '')
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
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest}
 */
proto.bucketeer.account.CreateAPIKeyRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.CreateAPIKeyRequest();
  return proto.bucketeer.account.CreateAPIKeyRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.CreateAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest}
 */
proto.bucketeer.account.CreateAPIKeyRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_command_pb.CreateAPIKeyCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.CreateAPIKeyCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setName(value);
          break;
        case 5:
          var value = /** @type {!proto.bucketeer.account.APIKey.Role} */ (
            reader.readEnum()
          );
          msg.setRole(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setMaintainer(value);
          break;
        case 7:
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
proto.bucketeer.account.CreateAPIKeyRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.CreateAPIKeyRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.CreateAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.CreateAPIKeyRequest.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_account_command_pb.CreateAPIKeyCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(4, f);
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(5, f);
  }
  f = message.getMaintainer();
  if (f.length > 0) {
    writer.writeString(6, f);
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(7, f);
  }
};

/**
 * optional CreateAPIKeyCommand command = 1;
 * @return {?proto.bucketeer.account.CreateAPIKeyCommand}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.getCommand = function () {
  return /** @type{?proto.bucketeer.account.CreateAPIKeyCommand} */ (
    jspb.Message.getWrapperField(
      this,
      proto_account_command_pb.CreateAPIKeyCommand,
      1
    )
  );
};

/**
 * @param {?proto.bucketeer.account.CreateAPIKeyCommand|undefined} value
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 1, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.hasCommand = function () {
  return jspb.Message.getField(this, 1) != null;
};

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional string name = 4;
 * @return {string}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.getName = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.setName = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 4, value);
};

/**
 * optional APIKey.Role role = 5;
 * @return {!proto.bucketeer.account.APIKey.Role}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.getRole = function () {
  return /** @type {!proto.bucketeer.account.APIKey.Role} */ (
    jspb.Message.getFieldWithDefault(this, 5, 0)
  );
};

/**
 * @param {!proto.bucketeer.account.APIKey.Role} value
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.setRole = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};

/**
 * optional string maintainer = 6;
 * @return {string}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.getMaintainer =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.setMaintainer = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 6, value);
};

/**
 * optional string description = 7;
 * @return {string}
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.getDescription =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 7, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateAPIKeyRequest} returns this
 */
proto.bucketeer.account.CreateAPIKeyRequest.prototype.setDescription =
  function (value) {
    return jspb.Message.setProto3StringField(this, 7, value);
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
  proto.bucketeer.account.CreateAPIKeyResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.CreateAPIKeyResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.CreateAPIKeyResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.CreateAPIKeyResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        apiKey:
          (f = msg.getApiKey()) &&
          proto_account_api_key_pb.APIKey.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.CreateAPIKeyResponse}
 */
proto.bucketeer.account.CreateAPIKeyResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.CreateAPIKeyResponse();
  return proto.bucketeer.account.CreateAPIKeyResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.CreateAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.CreateAPIKeyResponse}
 */
proto.bucketeer.account.CreateAPIKeyResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_api_key_pb.APIKey();
          reader.readMessage(
            value,
            proto_account_api_key_pb.APIKey.deserializeBinaryFromReader
          );
          msg.setApiKey(value);
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
proto.bucketeer.account.CreateAPIKeyResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.CreateAPIKeyResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.CreateAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.CreateAPIKeyResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getApiKey();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_api_key_pb.APIKey.serializeBinaryToWriter
      );
    }
  };

/**
 * optional APIKey api_key = 1;
 * @return {?proto.bucketeer.account.APIKey}
 */
proto.bucketeer.account.CreateAPIKeyResponse.prototype.getApiKey = function () {
  return /** @type{?proto.bucketeer.account.APIKey} */ (
    jspb.Message.getWrapperField(this, proto_account_api_key_pb.APIKey, 1)
  );
};

/**
 * @param {?proto.bucketeer.account.APIKey|undefined} value
 * @return {!proto.bucketeer.account.CreateAPIKeyResponse} returns this
 */
proto.bucketeer.account.CreateAPIKeyResponse.prototype.setApiKey = function (
  value
) {
  return jspb.Message.setWrapperField(this, 1, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.CreateAPIKeyResponse} returns this
 */
proto.bucketeer.account.CreateAPIKeyResponse.prototype.clearApiKey =
  function () {
    return this.setApiKey(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.CreateAPIKeyResponse.prototype.hasApiKey = function () {
  return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.ChangeAPIKeyNameRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.ChangeAPIKeyNameRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.ChangeAPIKeyNameRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.ChangeAPIKeyNameCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, '')
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
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameRequest}
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.ChangeAPIKeyNameRequest();
  return proto.bucketeer.account.ChangeAPIKeyNameRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.ChangeAPIKeyNameRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameRequest}
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.deserializeBinaryFromReader =
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
          var value = new proto_account_command_pb.ChangeAPIKeyNameCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeAPIKeyNameCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.ChangeAPIKeyNameRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.ChangeAPIKeyNameRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_account_command_pb.ChangeAPIKeyNameCommand.serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameRequest} returns this
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional ChangeAPIKeyNameCommand command = 2;
 * @return {?proto.bucketeer.account.ChangeAPIKeyNameCommand}
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeAPIKeyNameCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeAPIKeyNameCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeAPIKeyNameCommand|undefined} value
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameRequest} returns this
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameRequest} returns this
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameRequest} returns this
 */
proto.bucketeer.account.ChangeAPIKeyNameRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
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
  proto.bucketeer.account.ChangeAPIKeyNameResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.ChangeAPIKeyNameResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.ChangeAPIKeyNameResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.ChangeAPIKeyNameResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameResponse}
 */
proto.bucketeer.account.ChangeAPIKeyNameResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.ChangeAPIKeyNameResponse();
  return proto.bucketeer.account.ChangeAPIKeyNameResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.ChangeAPIKeyNameResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.ChangeAPIKeyNameResponse}
 */
proto.bucketeer.account.ChangeAPIKeyNameResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.ChangeAPIKeyNameResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.ChangeAPIKeyNameResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.ChangeAPIKeyNameResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.ChangeAPIKeyNameResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.EnableAPIKeyRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.EnableAPIKeyRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.EnableAPIKeyRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.EnableAPIKeyRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.EnableAPIKeyCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, '')
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
 * @return {!proto.bucketeer.account.EnableAPIKeyRequest}
 */
proto.bucketeer.account.EnableAPIKeyRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.EnableAPIKeyRequest();
  return proto.bucketeer.account.EnableAPIKeyRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.EnableAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.EnableAPIKeyRequest}
 */
proto.bucketeer.account.EnableAPIKeyRequest.deserializeBinaryFromReader =
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
          var value = new proto_account_command_pb.EnableAPIKeyCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.EnableAPIKeyCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.account.EnableAPIKeyRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.EnableAPIKeyRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.EnableAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnableAPIKeyRequest.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getCommand();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto_account_command_pb.EnableAPIKeyCommand.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(4, f);
  }
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.EnableAPIKeyRequest} returns this
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional EnableAPIKeyCommand command = 2;
 * @return {?proto.bucketeer.account.EnableAPIKeyCommand}
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.getCommand = function () {
  return /** @type{?proto.bucketeer.account.EnableAPIKeyCommand} */ (
    jspb.Message.getWrapperField(
      this,
      proto_account_command_pb.EnableAPIKeyCommand,
      2
    )
  );
};

/**
 * @param {?proto.bucketeer.account.EnableAPIKeyCommand|undefined} value
 * @return {!proto.bucketeer.account.EnableAPIKeyRequest} returns this
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 2, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.EnableAPIKeyRequest} returns this
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.hasCommand = function () {
  return jspb.Message.getField(this, 2) != null;
};

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.EnableAPIKeyRequest} returns this
 */
proto.bucketeer.account.EnableAPIKeyRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
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
  proto.bucketeer.account.EnableAPIKeyResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.EnableAPIKeyResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.EnableAPIKeyResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.EnableAPIKeyResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.EnableAPIKeyResponse}
 */
proto.bucketeer.account.EnableAPIKeyResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.EnableAPIKeyResponse();
  return proto.bucketeer.account.EnableAPIKeyResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.EnableAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.EnableAPIKeyResponse}
 */
proto.bucketeer.account.EnableAPIKeyResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.EnableAPIKeyResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.EnableAPIKeyResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.EnableAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.EnableAPIKeyResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.DisableAPIKeyRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.DisableAPIKeyRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DisableAPIKeyRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DisableAPIKeyRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.DisableAPIKeyCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 4, '')
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
 * @return {!proto.bucketeer.account.DisableAPIKeyRequest}
 */
proto.bucketeer.account.DisableAPIKeyRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DisableAPIKeyRequest();
  return proto.bucketeer.account.DisableAPIKeyRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DisableAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DisableAPIKeyRequest}
 */
proto.bucketeer.account.DisableAPIKeyRequest.deserializeBinaryFromReader =
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
          var value = new proto_account_command_pb.DisableAPIKeyCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.DisableAPIKeyCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
          break;
        case 4:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.account.DisableAPIKeyRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DisableAPIKeyRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DisableAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DisableAPIKeyRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_account_command_pb.DisableAPIKeyCommand.serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DisableAPIKeyRequest} returns this
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional DisableAPIKeyCommand command = 2;
 * @return {?proto.bucketeer.account.DisableAPIKeyCommand}
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.DisableAPIKeyCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.DisableAPIKeyCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.DisableAPIKeyCommand|undefined} value
 * @return {!proto.bucketeer.account.DisableAPIKeyRequest} returns this
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.setCommand = function (
  value
) {
  return jspb.Message.setWrapperField(this, 2, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.DisableAPIKeyRequest} returns this
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DisableAPIKeyRequest} returns this
 */
proto.bucketeer.account.DisableAPIKeyRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
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
  proto.bucketeer.account.DisableAPIKeyResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.DisableAPIKeyResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DisableAPIKeyResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DisableAPIKeyResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.DisableAPIKeyResponse}
 */
proto.bucketeer.account.DisableAPIKeyResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DisableAPIKeyResponse();
  return proto.bucketeer.account.DisableAPIKeyResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DisableAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DisableAPIKeyResponse}
 */
proto.bucketeer.account.DisableAPIKeyResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.DisableAPIKeyResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DisableAPIKeyResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DisableAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DisableAPIKeyResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.GetAPIKeyRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.GetAPIKeyRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAPIKeyRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAPIKeyRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, '')
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
 * @return {!proto.bucketeer.account.GetAPIKeyRequest}
 */
proto.bucketeer.account.GetAPIKeyRequest.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetAPIKeyRequest();
  return proto.bucketeer.account.GetAPIKeyRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAPIKeyRequest}
 */
proto.bucketeer.account.GetAPIKeyRequest.deserializeBinaryFromReader =
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
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
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
proto.bucketeer.account.GetAPIKeyRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAPIKeyRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAPIKeyRequest.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(3, f);
  }
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.GetAPIKeyRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAPIKeyRequest} returns this
 */
proto.bucketeer.account.GetAPIKeyRequest.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.account.GetAPIKeyRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAPIKeyRequest} returns this
 */
proto.bucketeer.account.GetAPIKeyRequest.prototype.setEnvironmentId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 3, value);
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
  proto.bucketeer.account.GetAPIKeyResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.GetAPIKeyResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAPIKeyResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAPIKeyResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        apiKey:
          (f = msg.getApiKey()) &&
          proto_account_api_key_pb.APIKey.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.GetAPIKeyResponse}
 */
proto.bucketeer.account.GetAPIKeyResponse.deserializeBinary = function (bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.GetAPIKeyResponse();
  return proto.bucketeer.account.GetAPIKeyResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAPIKeyResponse}
 */
proto.bucketeer.account.GetAPIKeyResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_api_key_pb.APIKey();
          reader.readMessage(
            value,
            proto_account_api_key_pb.APIKey.deserializeBinaryFromReader
          );
          msg.setApiKey(value);
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
proto.bucketeer.account.GetAPIKeyResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAPIKeyResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAPIKeyResponse.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getApiKey();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto_account_api_key_pb.APIKey.serializeBinaryToWriter
    );
  }
};

/**
 * optional APIKey api_key = 1;
 * @return {?proto.bucketeer.account.APIKey}
 */
proto.bucketeer.account.GetAPIKeyResponse.prototype.getApiKey = function () {
  return /** @type{?proto.bucketeer.account.APIKey} */ (
    jspb.Message.getWrapperField(this, proto_account_api_key_pb.APIKey, 1)
  );
};

/**
 * @param {?proto.bucketeer.account.APIKey|undefined} value
 * @return {!proto.bucketeer.account.GetAPIKeyResponse} returns this
 */
proto.bucketeer.account.GetAPIKeyResponse.prototype.setApiKey = function (
  value
) {
  return jspb.Message.setWrapperField(this, 1, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.GetAPIKeyResponse} returns this
 */
proto.bucketeer.account.GetAPIKeyResponse.prototype.clearApiKey = function () {
  return this.setApiKey(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.GetAPIKeyResponse.prototype.hasApiKey = function () {
  return jspb.Message.getField(this, 1) != null;
};

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.account.ListAPIKeysRequest.repeatedFields_ = [9];

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
  proto.bucketeer.account.ListAPIKeysRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.ListAPIKeysRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.ListAPIKeysRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.ListAPIKeysRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
        cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
        orderBy: jspb.Message.getFieldWithDefault(msg, 4, 0),
        orderDirection: jspb.Message.getFieldWithDefault(msg, 5, 0),
        searchKeyword: jspb.Message.getFieldWithDefault(msg, 6, ''),
        disabled:
          (f = msg.getDisabled()) &&
          google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
        environmentId: jspb.Message.getFieldWithDefault(msg, 8, ''),
        environmentIdsList:
          (f = jspb.Message.getRepeatedField(msg, 9)) == null ? undefined : f,
        organizationId: jspb.Message.getFieldWithDefault(msg, 10, '')
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
 * @return {!proto.bucketeer.account.ListAPIKeysRequest}
 */
proto.bucketeer.account.ListAPIKeysRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.ListAPIKeysRequest();
  return proto.bucketeer.account.ListAPIKeysRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.ListAPIKeysRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.ListAPIKeysRequest}
 */
proto.bucketeer.account.ListAPIKeysRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 4:
          var value =
            /** @type {!proto.bucketeer.account.ListAPIKeysRequest.OrderBy} */ (
              reader.readEnum()
            );
          msg.setOrderBy(value);
          break;
        case 5:
          var value =
            /** @type {!proto.bucketeer.account.ListAPIKeysRequest.OrderDirection} */ (
              reader.readEnum()
            );
          msg.setOrderDirection(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setSearchKeyword(value);
          break;
        case 7:
          var value = new google_protobuf_wrappers_pb.BoolValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader
          );
          msg.setDisabled(value);
          break;
        case 8:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 9:
          var value = /** @type {string} */ (reader.readString());
          msg.addEnvironmentIds(value);
          break;
        case 10:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
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
proto.bucketeer.account.ListAPIKeysRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.ListAPIKeysRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.ListAPIKeysRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.ListAPIKeysRequest.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getPageSize();
  if (f !== 0) {
    writer.writeInt64(1, f);
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getOrderBy();
  if (f !== 0.0) {
    writer.writeEnum(4, f);
  }
  f = message.getOrderDirection();
  if (f !== 0.0) {
    writer.writeEnum(5, f);
  }
  f = message.getSearchKeyword();
  if (f.length > 0) {
    writer.writeString(6, f);
  }
  f = message.getDisabled();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(8, f);
  }
  f = message.getEnvironmentIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(9, f);
  }
  f = message.getOrganizationId();
  if (f.length > 0) {
    writer.writeString(10, f);
  }
};

/**
 * @enum {number}
 */
proto.bucketeer.account.ListAPIKeysRequest.OrderBy = {
  DEFAULT: 0,
  NAME: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3,
  ROLE: 4,
  ENVIRONMENT: 5
};

/**
 * @enum {number}
 */
proto.bucketeer.account.ListAPIKeysRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getPageSize = function () {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};

/**
 * @param {number} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setPageSize = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 1, value);
};

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getCursor = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.account.ListAPIKeysRequest.OrderBy}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getOrderBy = function () {
  return /** @type {!proto.bucketeer.account.ListAPIKeysRequest.OrderBy} */ (
    jspb.Message.getFieldWithDefault(this, 4, 0)
  );
};

/**
 * @param {!proto.bucketeer.account.ListAPIKeysRequest.OrderBy} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setOrderBy = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};

/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.account.ListAPIKeysRequest.OrderDirection}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getOrderDirection =
  function () {
    return /** @type {!proto.bucketeer.account.ListAPIKeysRequest.OrderDirection} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.account.ListAPIKeysRequest.OrderDirection} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setOrderDirection =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
  };

/**
 * optional string search_keyword = 6;
 * @return {string}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getSearchKeyword =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setSearchKeyword =
  function (value) {
    return jspb.Message.setProto3StringField(this, 6, value);
  };

/**
 * optional google.protobuf.BoolValue disabled = 7;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getDisabled = function () {
  return /** @type{?proto.google.protobuf.BoolValue} */ (
    jspb.Message.getWrapperField(this, google_protobuf_wrappers_pb.BoolValue, 7)
  );
};

/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setDisabled = function (
  value
) {
  return jspb.Message.setWrapperField(this, 7, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.clearDisabled =
  function () {
    return this.setDisabled(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.hasDisabled = function () {
  return jspb.Message.getField(this, 7) != null;
};

/**
 * optional string environment_id = 8;
 * @return {string}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 8, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 8, value);
  };

/**
 * repeated string environment_ids = 9;
 * @return {!Array<string>}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getEnvironmentIdsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 9)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setEnvironmentIdsList =
  function (value) {
    return jspb.Message.setField(this, 9, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.addEnvironmentIds =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 9, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.clearEnvironmentIdsList =
  function () {
    return this.setEnvironmentIdsList([]);
  };

/**
 * optional string organization_id = 10;
 * @return {string}
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 10, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAPIKeysRequest} returns this
 */
proto.bucketeer.account.ListAPIKeysRequest.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 10, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.account.ListAPIKeysResponse.repeatedFields_ = [1];

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
  proto.bucketeer.account.ListAPIKeysResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.ListAPIKeysResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.ListAPIKeysResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.ListAPIKeysResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        apiKeysList: jspb.Message.toObjectList(
          msg.getApiKeysList(),
          proto_account_api_key_pb.APIKey.toObject,
          includeInstance
        ),
        cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
        totalCount: jspb.Message.getFieldWithDefault(msg, 3, 0)
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
 * @return {!proto.bucketeer.account.ListAPIKeysResponse}
 */
proto.bucketeer.account.ListAPIKeysResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.ListAPIKeysResponse();
  return proto.bucketeer.account.ListAPIKeysResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.ListAPIKeysResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.ListAPIKeysResponse}
 */
proto.bucketeer.account.ListAPIKeysResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_api_key_pb.APIKey();
          reader.readMessage(
            value,
            proto_account_api_key_pb.APIKey.deserializeBinaryFromReader
          );
          msg.addApiKeys(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 3:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setTotalCount(value);
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
proto.bucketeer.account.ListAPIKeysResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.ListAPIKeysResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.ListAPIKeysResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.ListAPIKeysResponse.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getApiKeysList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto_account_api_key_pb.APIKey.serializeBinaryToWriter
    );
  }
  f = message.getCursor();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getTotalCount();
  if (f !== 0) {
    writer.writeInt64(3, f);
  }
};

/**
 * repeated APIKey api_keys = 1;
 * @return {!Array<!proto.bucketeer.account.APIKey>}
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.getApiKeysList =
  function () {
    return /** @type{!Array<!proto.bucketeer.account.APIKey>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_account_api_key_pb.APIKey,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.account.APIKey>} value
 * @return {!proto.bucketeer.account.ListAPIKeysResponse} returns this
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.setApiKeysList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.account.APIKey=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.account.APIKey}
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.addApiKeys = function (
  opt_value,
  opt_index
) {
  return jspb.Message.addToRepeatedWrapperField(
    this,
    1,
    opt_value,
    proto.bucketeer.account.APIKey,
    opt_index
  );
};

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.account.ListAPIKeysResponse} returns this
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.clearApiKeysList =
  function () {
    return this.setApiKeysList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.getCursor = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.ListAPIKeysResponse} returns this
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.setCursor = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.getTotalCount =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.account.ListAPIKeysResponse} returns this
 */
proto.bucketeer.account.ListAPIKeysResponse.prototype.setTotalCount = function (
  value
) {
  return jspb.Message.setProto3IntField(this, 3, value);
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
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          id: jspb.Message.getFieldWithDefault(msg, 1, ''),
          apiKey: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest();
    return proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.deserializeBinaryFromReader =
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
          msg.setApiKey(value);
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
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getApiKey();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest} returns this
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional string api_key = 2;
 * @return {string}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.prototype.getApiKey =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest} returns this
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsRequest.prototype.setApiKey =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
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
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          environmentApiKey:
            (f = msg.getEnvironmentApiKey()) &&
            proto_account_api_key_pb.EnvironmentAPIKey.toObject(
              includeInstance,
              f
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
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse();
    return proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_account_api_key_pb.EnvironmentAPIKey();
          reader.readMessage(
            value,
            proto_account_api_key_pb.EnvironmentAPIKey
              .deserializeBinaryFromReader
          );
          msg.setEnvironmentApiKey(value);
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
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEnvironmentApiKey();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_account_api_key_pb.EnvironmentAPIKey.serializeBinaryToWriter
      );
    }
  };

/**
 * optional EnvironmentAPIKey environment_api_key = 1;
 * @return {?proto.bucketeer.account.EnvironmentAPIKey}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.prototype.getEnvironmentApiKey =
  function () {
    return /** @type{?proto.bucketeer.account.EnvironmentAPIKey} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_api_key_pb.EnvironmentAPIKey,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.EnvironmentAPIKey|undefined} value
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse} returns this
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.prototype.setEnvironmentApiKey =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse} returns this
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.prototype.clearEnvironmentApiKey =
  function () {
    return this.setEnvironmentApiKey(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.GetAPIKeyBySearchingAllEnvironmentsResponse.prototype.hasEnvironmentApiKey =
  function () {
    return jspb.Message.getField(this, 1) != null;
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
  proto.bucketeer.account.CreateSearchFilterRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.CreateSearchFilterRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.CreateSearchFilterRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.CreateSearchFilterRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.CreateSearchFilterCommand.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest}
 */
proto.bucketeer.account.CreateSearchFilterRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.CreateSearchFilterRequest();
  return proto.bucketeer.account.CreateSearchFilterRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.CreateSearchFilterRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest}
 */
proto.bucketeer.account.CreateSearchFilterRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value = new proto_account_command_pb.CreateSearchFilterCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.CreateSearchFilterCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
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
proto.bucketeer.account.CreateSearchFilterRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.CreateSearchFilterRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.CreateSearchFilterRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.CreateSearchFilterRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_account_command_pb.CreateSearchFilterCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest} returns this
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.setEmail =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest} returns this
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest} returns this
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional CreateSearchFilterCommand command = 4;
 * @return {?proto.bucketeer.account.CreateSearchFilterCommand}
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.CreateSearchFilterCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.CreateSearchFilterCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.CreateSearchFilterCommand|undefined} value
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest} returns this
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.CreateSearchFilterRequest} returns this
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.CreateSearchFilterRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 4) != null;
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
  proto.bucketeer.account.CreateSearchFilterResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.CreateSearchFilterResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.CreateSearchFilterResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.CreateSearchFilterResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.CreateSearchFilterResponse}
 */
proto.bucketeer.account.CreateSearchFilterResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.CreateSearchFilterResponse();
    return proto.bucketeer.account.CreateSearchFilterResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.CreateSearchFilterResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.CreateSearchFilterResponse}
 */
proto.bucketeer.account.CreateSearchFilterResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.CreateSearchFilterResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.CreateSearchFilterResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.CreateSearchFilterResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.CreateSearchFilterResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.UpdateSearchFilterRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.UpdateSearchFilterRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateSearchFilterRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateSearchFilterRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        changeNameCommand:
          (f = msg.getChangeNameCommand()) &&
          proto_account_command_pb.ChangeSearchFilterNameCommand.toObject(
            includeInstance,
            f
          ),
        changeQueryCommand:
          (f = msg.getChangeQueryCommand()) &&
          proto_account_command_pb.ChangeSearchFilterQueryCommand.toObject(
            includeInstance,
            f
          ),
        changeDefaultFilterCommand:
          (f = msg.getChangeDefaultFilterCommand()) &&
          proto_account_command_pb.ChangeDefaultSearchFilterCommand.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.UpdateSearchFilterRequest();
  return proto.bucketeer.account.UpdateSearchFilterRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateSearchFilterRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value =
            new proto_account_command_pb.ChangeSearchFilterNameCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeSearchFilterNameCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeNameCommand(value);
          break;
        case 5:
          var value =
            new proto_account_command_pb.ChangeSearchFilterQueryCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeSearchFilterQueryCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeQueryCommand(value);
          break;
        case 6:
          var value =
            new proto_account_command_pb.ChangeDefaultSearchFilterCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.ChangeDefaultSearchFilterCommand
              .deserializeBinaryFromReader
          );
          msg.setChangeDefaultFilterCommand(value);
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
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateSearchFilterRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateSearchFilterRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateSearchFilterRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getChangeNameCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_account_command_pb.ChangeSearchFilterNameCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeQueryCommand();
    if (f != null) {
      writer.writeMessage(
        5,
        f,
        proto_account_command_pb.ChangeSearchFilterQueryCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getChangeDefaultFilterCommand();
    if (f != null) {
      writer.writeMessage(
        6,
        f,
        proto_account_command_pb.ChangeDefaultSearchFilterCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.setEmail =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional ChangeSearchFilterNameCommand change_name_command = 4;
 * @return {?proto.bucketeer.account.ChangeSearchFilterNameCommand}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.getChangeNameCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeSearchFilterNameCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeSearchFilterNameCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeSearchFilterNameCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.setChangeNameCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.clearChangeNameCommand =
  function () {
    return this.setChangeNameCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.hasChangeNameCommand =
  function () {
    return jspb.Message.getField(this, 4) != null;
  };

/**
 * optional ChangeSearchFilterQueryCommand change_query_command = 5;
 * @return {?proto.bucketeer.account.ChangeSearchFilterQueryCommand}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.getChangeQueryCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeSearchFilterQueryCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeSearchFilterQueryCommand,
        5
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeSearchFilterQueryCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.setChangeQueryCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 5, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.clearChangeQueryCommand =
  function () {
    return this.setChangeQueryCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.hasChangeQueryCommand =
  function () {
    return jspb.Message.getField(this, 5) != null;
  };

/**
 * optional ChangeDefaultSearchFilterCommand change_default_filter_command = 6;
 * @return {?proto.bucketeer.account.ChangeDefaultSearchFilterCommand}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.getChangeDefaultFilterCommand =
  function () {
    return /** @type{?proto.bucketeer.account.ChangeDefaultSearchFilterCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.ChangeDefaultSearchFilterCommand,
        6
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.ChangeDefaultSearchFilterCommand|undefined} value
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.setChangeDefaultFilterCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 6, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateSearchFilterRequest} returns this
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.clearChangeDefaultFilterCommand =
  function () {
    return this.setChangeDefaultFilterCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateSearchFilterRequest.prototype.hasChangeDefaultFilterCommand =
  function () {
    return jspb.Message.getField(this, 6) != null;
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
  proto.bucketeer.account.UpdateSearchFilterResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.UpdateSearchFilterResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateSearchFilterResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateSearchFilterResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.UpdateSearchFilterResponse}
 */
proto.bucketeer.account.UpdateSearchFilterResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.UpdateSearchFilterResponse();
    return proto.bucketeer.account.UpdateSearchFilterResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateSearchFilterResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateSearchFilterResponse}
 */
proto.bucketeer.account.UpdateSearchFilterResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.UpdateSearchFilterResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateSearchFilterResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateSearchFilterResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateSearchFilterResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.DeleteSearchFilterRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.DeleteSearchFilterRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DeleteSearchFilterRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DeleteSearchFilterRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        email: jspb.Message.getFieldWithDefault(msg, 1, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        command:
          (f = msg.getCommand()) &&
          proto_account_command_pb.DeleteSearchFilterCommand.toObject(
            includeInstance,
            f
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
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.DeleteSearchFilterRequest();
  return proto.bucketeer.account.DeleteSearchFilterRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DeleteSearchFilterRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = /** @type {string} */ (reader.readString());
          msg.setEmail(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 4:
          var value = new proto_account_command_pb.DeleteSearchFilterCommand();
          reader.readMessage(
            value,
            proto_account_command_pb.DeleteSearchFilterCommand
              .deserializeBinaryFromReader
          );
          msg.setCommand(value);
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
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DeleteSearchFilterRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DeleteSearchFilterRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DeleteSearchFilterRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getEmail();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_account_command_pb.DeleteSearchFilterCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string email = 1;
 * @return {string}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.getEmail =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest} returns this
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.setEmail =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional string organization_id = 2;
 * @return {string}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest} returns this
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest} returns this
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional DeleteSearchFilterCommand command = 4;
 * @return {?proto.bucketeer.account.DeleteSearchFilterCommand}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.account.DeleteSearchFilterCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_account_command_pb.DeleteSearchFilterCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.account.DeleteSearchFilterCommand|undefined} value
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest} returns this
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.DeleteSearchFilterRequest} returns this
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.DeleteSearchFilterRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 4) != null;
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
  proto.bucketeer.account.DeleteSearchFilterResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.account.DeleteSearchFilterResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.DeleteSearchFilterResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.DeleteSearchFilterResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.DeleteSearchFilterResponse}
 */
proto.bucketeer.account.DeleteSearchFilterResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.account.DeleteSearchFilterResponse();
    return proto.bucketeer.account.DeleteSearchFilterResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.DeleteSearchFilterResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.DeleteSearchFilterResponse}
 */
proto.bucketeer.account.DeleteSearchFilterResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.DeleteSearchFilterResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.DeleteSearchFilterResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.DeleteSearchFilterResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.DeleteSearchFilterResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
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
  proto.bucketeer.account.UpdateAPIKeyRequest.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.UpdateAPIKeyRequest.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateAPIKeyRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateAPIKeyRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, ''),
        environmentId: jspb.Message.getFieldWithDefault(msg, 2, ''),
        name:
          (f = msg.getName()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        description:
          (f = msg.getDescription()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        role: jspb.Message.getFieldWithDefault(msg, 5, 0),
        disabled:
          (f = msg.getDisabled()) &&
          google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
        maintainer:
          (f = msg.getMaintainer()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.UpdateAPIKeyRequest();
  return proto.bucketeer.account.UpdateAPIKeyRequest.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateAPIKeyRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.deserializeBinaryFromReader =
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
          msg.setEnvironmentId(value);
          break;
        case 3:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setName(value);
          break;
        case 4:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setDescription(value);
          break;
        case 5:
          var value = /** @type {!proto.bucketeer.account.APIKey.Role} */ (
            reader.readEnum()
          );
          msg.setRole(value);
          break;
        case 6:
          var value = new google_protobuf_wrappers_pb.BoolValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader
          );
          msg.setDisabled(value);
          break;
        case 7:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setMaintainer(value);
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
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateAPIKeyRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateAPIKeyRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateAPIKeyRequest.serializeBinaryToWriter = function (
  message,
  writer
) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(1, f);
  }
  f = message.getEnvironmentId();
  if (f.length > 0) {
    writer.writeString(2, f);
  }
  f = message.getName();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
  f = message.getDescription();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
  f = message.getRole();
  if (f !== 0.0) {
    writer.writeEnum(5, f);
  }
  f = message.getDisabled();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
    );
  }
  f = message.getMaintainer();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
    );
  }
};

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getId = function () {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ''));
};

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setId = function (value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};

/**
 * optional string environment_id = 2;
 * @return {string}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional google.protobuf.StringValue name = 3;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getName = function () {
  return /** @type{?proto.google.protobuf.StringValue} */ (
    jspb.Message.getWrapperField(
      this,
      google_protobuf_wrappers_pb.StringValue,
      3
    )
  );
};

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setName = function (
  value
) {
  return jspb.Message.setWrapperField(this, 3, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.clearName = function () {
  return this.setName(undefined);
};

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.hasName = function () {
  return jspb.Message.getField(this, 3) != null;
};

/**
 * optional google.protobuf.StringValue description = 4;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getDescription =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        4
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setDescription =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.clearDescription =
  function () {
    return this.setDescription(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.hasDescription =
  function () {
    return jspb.Message.getField(this, 4) != null;
  };

/**
 * optional APIKey.Role role = 5;
 * @return {!proto.bucketeer.account.APIKey.Role}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getRole = function () {
  return /** @type {!proto.bucketeer.account.APIKey.Role} */ (
    jspb.Message.getFieldWithDefault(this, 5, 0)
  );
};

/**
 * @param {!proto.bucketeer.account.APIKey.Role} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setRole = function (
  value
) {
  return jspb.Message.setProto3EnumField(this, 5, value);
};

/**
 * optional google.protobuf.BoolValue disabled = 6;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getDisabled =
  function () {
    return /** @type{?proto.google.protobuf.BoolValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.BoolValue,
        6
      )
    );
  };

/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setDisabled = function (
  value
) {
  return jspb.Message.setWrapperField(this, 6, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.clearDisabled =
  function () {
    return this.setDisabled(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.hasDisabled =
  function () {
    return jspb.Message.getField(this, 6) != null;
  };

/**
 * optional google.protobuf.StringValue maintainer = 7;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.getMaintainer =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        7
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.setMaintainer = function (
  value
) {
  return jspb.Message.setWrapperField(this, 7, value);
};

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.account.UpdateAPIKeyRequest} returns this
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.clearMaintainer =
  function () {
    return this.setMaintainer(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.account.UpdateAPIKeyRequest.prototype.hasMaintainer =
  function () {
    return jspb.Message.getField(this, 7) != null;
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
  proto.bucketeer.account.UpdateAPIKeyResponse.prototype.toObject = function (
    opt_includeInstance
  ) {
    return proto.bucketeer.account.UpdateAPIKeyResponse.toObject(
      opt_includeInstance,
      this
    );
  };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.account.UpdateAPIKeyResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.account.UpdateAPIKeyResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {};

    if (includeInstance) {
      obj.$jspbMessageInstance = msg;
    }
    return obj;
  };
}

/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.bucketeer.account.UpdateAPIKeyResponse}
 */
proto.bucketeer.account.UpdateAPIKeyResponse.deserializeBinary = function (
  bytes
) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.bucketeer.account.UpdateAPIKeyResponse();
  return proto.bucketeer.account.UpdateAPIKeyResponse.deserializeBinaryFromReader(
    msg,
    reader
  );
};

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.account.UpdateAPIKeyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.account.UpdateAPIKeyResponse}
 */
proto.bucketeer.account.UpdateAPIKeyResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
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
proto.bucketeer.account.UpdateAPIKeyResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.account.UpdateAPIKeyResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.account.UpdateAPIKeyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.account.UpdateAPIKeyResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
  };

goog.object.extend(exports, proto.bucketeer.account);
