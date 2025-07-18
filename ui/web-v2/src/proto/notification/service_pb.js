// source: proto/notification/service.proto
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
var proto_notification_subscription_pb = require('../../proto/notification/subscription_pb.js');
goog.object.extend(proto, proto_notification_subscription_pb);
var proto_notification_recipient_pb = require('../../proto/notification/recipient_pb.js');
goog.object.extend(proto, proto_notification_recipient_pb);
var proto_notification_command_pb = require('../../proto/notification/command_pb.js');
goog.object.extend(proto, proto_notification_command_pb);
goog.exportSymbol(
  'proto.bucketeer.notification.CreateAdminSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.CreateAdminSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.CreateSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.CreateSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DeleteAdminSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DeleteAdminSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DeleteSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DeleteSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DisableAdminSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DisableAdminSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DisableSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.DisableSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.EnableAdminSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.EnableAdminSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.EnableSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.EnableSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.GetAdminSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.GetAdminSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.GetSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.GetSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListAdminSubscriptionsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderBy',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderDirection',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListAdminSubscriptionsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListEnabledSubscriptionsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListEnabledSubscriptionsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListSubscriptionsRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListSubscriptionsRequest.OrderBy',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListSubscriptionsRequest.OrderDirection',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.ListSubscriptionsResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.UpdateAdminSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.UpdateAdminSubscriptionResponse',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.UpdateSubscriptionRequest',
  null,
  global
);
goog.exportSymbol(
  'proto.bucketeer.notification.UpdateSubscriptionResponse',
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
proto.bucketeer.notification.GetAdminSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.GetAdminSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.GetAdminSubscriptionRequest.displayName =
    'proto.bucketeer.notification.GetAdminSubscriptionRequest';
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
proto.bucketeer.notification.GetAdminSubscriptionResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.GetAdminSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.GetAdminSubscriptionResponse.displayName =
    'proto.bucketeer.notification.GetAdminSubscriptionResponse';
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
proto.bucketeer.notification.ListAdminSubscriptionsRequest = function (
  opt_data
) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListAdminSubscriptionsRequest.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListAdminSubscriptionsRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListAdminSubscriptionsRequest.displayName =
    'proto.bucketeer.notification.ListAdminSubscriptionsRequest';
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
proto.bucketeer.notification.ListAdminSubscriptionsResponse = function (
  opt_data
) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListAdminSubscriptionsResponse.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListAdminSubscriptionsResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListAdminSubscriptionsResponse.displayName =
    'proto.bucketeer.notification.ListAdminSubscriptionsResponse';
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
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest = function (
  opt_data
) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest
      .repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.displayName =
    'proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest';
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
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse = function (
  opt_data
) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse
      .repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.displayName =
    'proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse';
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
proto.bucketeer.notification.CreateAdminSubscriptionRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.CreateAdminSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.CreateAdminSubscriptionRequest.displayName =
    'proto.bucketeer.notification.CreateAdminSubscriptionRequest';
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
proto.bucketeer.notification.CreateAdminSubscriptionResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.CreateAdminSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.CreateAdminSubscriptionResponse.displayName =
    'proto.bucketeer.notification.CreateAdminSubscriptionResponse';
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
proto.bucketeer.notification.DeleteAdminSubscriptionRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DeleteAdminSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DeleteAdminSubscriptionRequest.displayName =
    'proto.bucketeer.notification.DeleteAdminSubscriptionRequest';
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
proto.bucketeer.notification.DeleteAdminSubscriptionResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DeleteAdminSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DeleteAdminSubscriptionResponse.displayName =
    'proto.bucketeer.notification.DeleteAdminSubscriptionResponse';
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
proto.bucketeer.notification.EnableAdminSubscriptionRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.EnableAdminSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.EnableAdminSubscriptionRequest.displayName =
    'proto.bucketeer.notification.EnableAdminSubscriptionRequest';
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
proto.bucketeer.notification.EnableAdminSubscriptionResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.EnableAdminSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.EnableAdminSubscriptionResponse.displayName =
    'proto.bucketeer.notification.EnableAdminSubscriptionResponse';
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
proto.bucketeer.notification.DisableAdminSubscriptionRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DisableAdminSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DisableAdminSubscriptionRequest.displayName =
    'proto.bucketeer.notification.DisableAdminSubscriptionRequest';
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
proto.bucketeer.notification.DisableAdminSubscriptionResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DisableAdminSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DisableAdminSubscriptionResponse.displayName =
    'proto.bucketeer.notification.DisableAdminSubscriptionResponse';
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
proto.bucketeer.notification.UpdateAdminSubscriptionRequest = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.UpdateAdminSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.UpdateAdminSubscriptionRequest.displayName =
    'proto.bucketeer.notification.UpdateAdminSubscriptionRequest';
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
proto.bucketeer.notification.UpdateAdminSubscriptionResponse = function (
  opt_data
) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.UpdateAdminSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.UpdateAdminSubscriptionResponse.displayName =
    'proto.bucketeer.notification.UpdateAdminSubscriptionResponse';
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
proto.bucketeer.notification.GetSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.GetSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.GetSubscriptionRequest.displayName =
    'proto.bucketeer.notification.GetSubscriptionRequest';
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
proto.bucketeer.notification.GetSubscriptionResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.GetSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.GetSubscriptionResponse.displayName =
    'proto.bucketeer.notification.GetSubscriptionResponse';
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
proto.bucketeer.notification.ListSubscriptionsRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListSubscriptionsRequest.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListSubscriptionsRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListSubscriptionsRequest.displayName =
    'proto.bucketeer.notification.ListSubscriptionsRequest';
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
proto.bucketeer.notification.ListSubscriptionsResponse = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListSubscriptionsResponse.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListSubscriptionsResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListSubscriptionsResponse.displayName =
    'proto.bucketeer.notification.ListSubscriptionsResponse';
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
proto.bucketeer.notification.ListEnabledSubscriptionsRequest = function (
  opt_data
) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListEnabledSubscriptionsRequest
      .repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListEnabledSubscriptionsRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListEnabledSubscriptionsRequest.displayName =
    'proto.bucketeer.notification.ListEnabledSubscriptionsRequest';
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
proto.bucketeer.notification.ListEnabledSubscriptionsResponse = function (
  opt_data
) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.ListEnabledSubscriptionsResponse
      .repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.ListEnabledSubscriptionsResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.ListEnabledSubscriptionsResponse.displayName =
    'proto.bucketeer.notification.ListEnabledSubscriptionsResponse';
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
proto.bucketeer.notification.CreateSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.CreateSubscriptionRequest.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.CreateSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.CreateSubscriptionRequest.displayName =
    'proto.bucketeer.notification.CreateSubscriptionRequest';
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
proto.bucketeer.notification.CreateSubscriptionResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.CreateSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.CreateSubscriptionResponse.displayName =
    'proto.bucketeer.notification.CreateSubscriptionResponse';
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
proto.bucketeer.notification.DeleteSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DeleteSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DeleteSubscriptionRequest.displayName =
    'proto.bucketeer.notification.DeleteSubscriptionRequest';
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
proto.bucketeer.notification.DeleteSubscriptionResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DeleteSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DeleteSubscriptionResponse.displayName =
    'proto.bucketeer.notification.DeleteSubscriptionResponse';
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
proto.bucketeer.notification.EnableSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.EnableSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.EnableSubscriptionRequest.displayName =
    'proto.bucketeer.notification.EnableSubscriptionRequest';
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
proto.bucketeer.notification.EnableSubscriptionResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.EnableSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.EnableSubscriptionResponse.displayName =
    'proto.bucketeer.notification.EnableSubscriptionResponse';
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
proto.bucketeer.notification.DisableSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DisableSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DisableSubscriptionRequest.displayName =
    'proto.bucketeer.notification.DisableSubscriptionRequest';
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
proto.bucketeer.notification.DisableSubscriptionResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.DisableSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.DisableSubscriptionResponse.displayName =
    'proto.bucketeer.notification.DisableSubscriptionResponse';
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
proto.bucketeer.notification.UpdateSubscriptionRequest = function (opt_data) {
  jspb.Message.initialize(
    this,
    opt_data,
    0,
    -1,
    proto.bucketeer.notification.UpdateSubscriptionRequest.repeatedFields_,
    null
  );
};
goog.inherits(
  proto.bucketeer.notification.UpdateSubscriptionRequest,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.UpdateSubscriptionRequest.displayName =
    'proto.bucketeer.notification.UpdateSubscriptionRequest';
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
proto.bucketeer.notification.UpdateSubscriptionResponse = function (opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(
  proto.bucketeer.notification.UpdateSubscriptionResponse,
  jspb.Message
);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.bucketeer.notification.UpdateSubscriptionResponse.displayName =
    'proto.bucketeer.notification.UpdateSubscriptionResponse';
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
  proto.bucketeer.notification.GetAdminSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.GetAdminSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.GetAdminSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.GetAdminSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 1, '')
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
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionRequest}
 */
proto.bucketeer.notification.GetAdminSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.GetAdminSubscriptionRequest();
    return proto.bucketeer.notification.GetAdminSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.GetAdminSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionRequest}
 */
proto.bucketeer.notification.GetAdminSubscriptionRequest.deserializeBinaryFromReader =
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
proto.bucketeer.notification.GetAdminSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.GetAdminSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.GetAdminSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.GetAdminSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.notification.GetAdminSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.GetAdminSubscriptionRequest.prototype.setId =
  function (value) {
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
  proto.bucketeer.notification.GetAdminSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.GetAdminSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.GetAdminSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.GetAdminSubscriptionResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          subscription:
            (f = msg.getSubscription()) &&
            proto_notification_subscription_pb.Subscription.toObject(
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
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionResponse}
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.GetAdminSubscriptionResponse();
    return proto.bucketeer.notification.GetAdminSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.GetAdminSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionResponse}
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.setSubscription(value);
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
proto.bucketeer.notification.GetAdminSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.GetAdminSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.GetAdminSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscription();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
      );
    }
  };

/**
 * optional Subscription subscription = 1;
 * @return {?proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.prototype.getSubscription =
  function () {
    return /** @type{?proto.bucketeer.notification.Subscription} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.Subscription|undefined} value
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionResponse} returns this
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.prototype.setSubscription =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.GetAdminSubscriptionResponse} returns this
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.prototype.clearSubscription =
  function () {
    return this.setSubscription(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.GetAdminSubscriptionResponse.prototype.hasSubscription =
  function () {
    return jspb.Message.getField(this, 1) != null;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.repeatedFields_ = [
  3
];

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
  proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListAdminSubscriptionsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListAdminSubscriptionsRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
          cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
          sourceTypesList:
            (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f,
          orderBy: jspb.Message.getFieldWithDefault(msg, 4, 0),
          orderDirection: jspb.Message.getFieldWithDefault(msg, 5, 0),
          searchKeyword: jspb.Message.getFieldWithDefault(msg, 6, ''),
          disabled:
            (f = msg.getDisabled()) &&
            google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f)
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
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.ListAdminSubscriptionsRequest();
    return proto.bucketeer.notification.ListAdminSubscriptionsRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.deserializeBinaryFromReader =
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
          var values =
            /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
              reader.isDelimited()
                ? reader.readPackedEnum()
                : [reader.readEnum()]
            );
          for (var i = 0; i < values.length; i++) {
            msg.addSourceTypes(values[i]);
          }
          break;
        case 4:
          var value =
            /** @type {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderBy} */ (
              reader.readEnum()
            );
          msg.setOrderBy(value);
          break;
        case 5:
          var value =
            /** @type {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderDirection} */ (
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
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListAdminSubscriptionsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.serializeBinaryToWriter =
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
    f = message.getSourceTypesList();
    if (f.length > 0) {
      writer.writePackedEnum(3, f);
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
  };

/**
 * @enum {number}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderBy = {
  DEFAULT: 0,
  NAME: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3
};

/**
 * @enum {number}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setPageSize =
  function (value) {
    return jspb.Message.setProto3IntField(this, 1, value);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * repeated Subscription.SourceType source_types = 3;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getSourceTypesList =
  function () {
    return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
      jspb.Message.getRepeatedField(this, 3)
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setSourceTypesList =
  function (value) {
    return jspb.Message.setField(this, 3, value || []);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.addSourceTypes =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.clearSourceTypesList =
  function () {
    return this.setSourceTypesList([]);
  };

/**
 * optional OrderBy order_by = 4;
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderBy}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getOrderBy =
  function () {
    return /** @type {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderBy} */ (
      jspb.Message.getFieldWithDefault(this, 4, 0)
    );
  };

/**
 * @param {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderBy} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setOrderBy =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 4, value);
  };

/**
 * optional OrderDirection order_direction = 5;
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderDirection}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getOrderDirection =
  function () {
    return /** @type {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderDirection} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.notification.ListAdminSubscriptionsRequest.OrderDirection} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setOrderDirection =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
  };

/**
 * optional string search_keyword = 6;
 * @return {string}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getSearchKeyword =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setSearchKeyword =
  function (value) {
    return jspb.Message.setProto3StringField(this, 6, value);
  };

/**
 * optional google.protobuf.BoolValue disabled = 7;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.getDisabled =
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
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.setDisabled =
  function (value) {
    return jspb.Message.setWrapperField(this, 7, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.clearDisabled =
  function () {
    return this.setDisabled(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.ListAdminSubscriptionsRequest.prototype.hasDisabled =
  function () {
    return jspb.Message.getField(this, 7) != null;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.repeatedFields_ = [
  1
];

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
  proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListAdminSubscriptionsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListAdminSubscriptionsResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          subscriptionsList: jspb.Message.toObjectList(
            msg.getSubscriptionsList(),
            proto_notification_subscription_pb.Subscription.toObject,
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
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsResponse}
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.ListAdminSubscriptionsResponse();
    return proto.bucketeer.notification.ListAdminSubscriptionsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsResponse}
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.addSubscriptions(value);
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
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListAdminSubscriptionsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscriptionsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
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
 * repeated Subscription subscriptions = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription>}
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.getSubscriptionsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.notification.Subscription>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription>} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.setSubscriptionsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.addSubscriptions =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.notification.Subscription,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.clearSubscriptionsList =
  function () {
    return this.setSubscriptionsList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.getTotalCount =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.notification.ListAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListAdminSubscriptionsResponse.prototype.setTotalCount =
  function (value) {
    return jspb.Message.setProto3IntField(this, 3, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.repeatedFields_ =
  [3];

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
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          pageSize: jspb.Message.getFieldWithDefault(msg, 1, 0),
          cursor: jspb.Message.getFieldWithDefault(msg, 2, ''),
          sourceTypesList:
            (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f
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
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest();
    return proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.deserializeBinaryFromReader =
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
          var values =
            /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
              reader.isDelimited()
                ? reader.readPackedEnum()
                : [reader.readEnum()]
            );
          for (var i = 0; i < values.length; i++) {
            msg.addSourceTypes(values[i]);
          }
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
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.serializeBinaryToWriter =
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
    f = message.getSourceTypesList();
    if (f.length > 0) {
      writer.writePackedEnum(3, f);
    }
  };

/**
 * optional int64 page_size = 1;
 * @return {number}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.setPageSize =
  function (value) {
    return jspb.Message.setProto3IntField(this, 1, value);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * repeated Subscription.SourceType source_types = 3;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.getSourceTypesList =
  function () {
    return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
      jspb.Message.getRepeatedField(this, 3)
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.setSourceTypesList =
  function (value) {
    return jspb.Message.setField(this, 3, value || []);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.addSourceTypes =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsRequest.prototype.clearSourceTypesList =
  function () {
    return this.setSourceTypesList([]);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.repeatedFields_ =
  [1];

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
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          subscriptionsList: jspb.Message.toObjectList(
            msg.getSubscriptionsList(),
            proto_notification_subscription_pb.Subscription.toObject,
            includeInstance
          ),
          cursor: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse();
    return proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.addSubscriptions(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
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
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscriptionsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
      );
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
  };

/**
 * repeated Subscription subscriptions = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription>}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.getSubscriptionsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.notification.Subscription>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription>} value
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.setSubscriptionsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.addSubscriptions =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.notification.Subscription,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.clearSubscriptionsList =
  function () {
    return this.setSubscriptionsList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListEnabledAdminSubscriptionsResponse.prototype.setCursor =
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
  proto.bucketeer.notification.CreateAdminSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.CreateAdminSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.CreateAdminSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.CreateAdminSubscriptionRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          command:
            (f = msg.getCommand()) &&
            proto_notification_command_pb.CreateAdminSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.CreateAdminSubscriptionRequest}
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.CreateAdminSubscriptionRequest();
    return proto.bucketeer.notification.CreateAdminSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.CreateAdminSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.CreateAdminSubscriptionRequest}
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value =
            new proto_notification_command_pb.CreateAdminSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.CreateAdminSubscriptionCommand
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
proto.bucketeer.notification.CreateAdminSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.CreateAdminSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.CreateAdminSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_notification_command_pb.CreateAdminSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional CreateAdminSubscriptionCommand command = 1;
 * @return {?proto.bucketeer.notification.CreateAdminSubscriptionCommand}
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.CreateAdminSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.CreateAdminSubscriptionCommand,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.CreateAdminSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.CreateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.CreateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.CreateAdminSubscriptionRequest.prototype.hasCommand =
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
  proto.bucketeer.notification.CreateAdminSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.CreateAdminSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.CreateAdminSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.CreateAdminSubscriptionResponse.toObject =
    function (includeInstance, msg) {
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
 * @return {!proto.bucketeer.notification.CreateAdminSubscriptionResponse}
 */
proto.bucketeer.notification.CreateAdminSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.CreateAdminSubscriptionResponse();
    return proto.bucketeer.notification.CreateAdminSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.CreateAdminSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.CreateAdminSubscriptionResponse}
 */
proto.bucketeer.notification.CreateAdminSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.CreateAdminSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.CreateAdminSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.CreateAdminSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.CreateAdminSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DeleteAdminSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DeleteAdminSubscriptionRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          id: jspb.Message.getFieldWithDefault(msg, 1, ''),
          command:
            (f = msg.getCommand()) &&
            proto_notification_command_pb.DeleteAdminSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.DeleteAdminSubscriptionRequest();
    return proto.bucketeer.notification.DeleteAdminSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.deserializeBinaryFromReader =
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
          var value =
            new proto_notification_command_pb.DeleteAdminSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.DeleteAdminSubscriptionCommand
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
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DeleteAdminSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.serializeBinaryToWriter =
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
        proto_notification_command_pb.DeleteAdminSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional DeleteAdminSubscriptionCommand command = 2;
 * @return {?proto.bucketeer.notification.DeleteAdminSubscriptionCommand}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.DeleteAdminSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.DeleteAdminSubscriptionCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.DeleteAdminSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
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
  proto.bucketeer.notification.DeleteAdminSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DeleteAdminSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DeleteAdminSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DeleteAdminSubscriptionResponse.toObject =
    function (includeInstance, msg) {
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
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionResponse}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.DeleteAdminSubscriptionResponse();
    return proto.bucketeer.notification.DeleteAdminSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DeleteAdminSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DeleteAdminSubscriptionResponse}
 */
proto.bucketeer.notification.DeleteAdminSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.DeleteAdminSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DeleteAdminSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DeleteAdminSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DeleteAdminSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.EnableAdminSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.EnableAdminSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.EnableAdminSubscriptionRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          id: jspb.Message.getFieldWithDefault(msg, 1, ''),
          command:
            (f = msg.getCommand()) &&
            proto_notification_command_pb.EnableAdminSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionRequest}
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.EnableAdminSubscriptionRequest();
    return proto.bucketeer.notification.EnableAdminSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.EnableAdminSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionRequest}
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.deserializeBinaryFromReader =
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
          var value =
            new proto_notification_command_pb.EnableAdminSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.EnableAdminSubscriptionCommand
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
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.EnableAdminSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.EnableAdminSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.serializeBinaryToWriter =
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
        proto_notification_command_pb.EnableAdminSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional EnableAdminSubscriptionCommand command = 2;
 * @return {?proto.bucketeer.notification.EnableAdminSubscriptionCommand}
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.EnableAdminSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.EnableAdminSubscriptionCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.EnableAdminSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.EnableAdminSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
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
  proto.bucketeer.notification.EnableAdminSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.EnableAdminSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.EnableAdminSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.EnableAdminSubscriptionResponse.toObject =
    function (includeInstance, msg) {
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
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionResponse}
 */
proto.bucketeer.notification.EnableAdminSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.EnableAdminSubscriptionResponse();
    return proto.bucketeer.notification.EnableAdminSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.EnableAdminSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.EnableAdminSubscriptionResponse}
 */
proto.bucketeer.notification.EnableAdminSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.EnableAdminSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.EnableAdminSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.EnableAdminSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.EnableAdminSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DisableAdminSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DisableAdminSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DisableAdminSubscriptionRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          id: jspb.Message.getFieldWithDefault(msg, 1, ''),
          command:
            (f = msg.getCommand()) &&
            proto_notification_command_pb.DisableAdminSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionRequest}
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.DisableAdminSubscriptionRequest();
    return proto.bucketeer.notification.DisableAdminSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DisableAdminSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionRequest}
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.deserializeBinaryFromReader =
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
          var value =
            new proto_notification_command_pb.DisableAdminSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.DisableAdminSubscriptionCommand
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
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DisableAdminSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DisableAdminSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.serializeBinaryToWriter =
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
        proto_notification_command_pb.DisableAdminSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional DisableAdminSubscriptionCommand command = 2;
 * @return {?proto.bucketeer.notification.DisableAdminSubscriptionCommand}
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.DisableAdminSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.DisableAdminSubscriptionCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.DisableAdminSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.DisableAdminSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
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
  proto.bucketeer.notification.DisableAdminSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DisableAdminSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DisableAdminSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DisableAdminSubscriptionResponse.toObject =
    function (includeInstance, msg) {
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
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionResponse}
 */
proto.bucketeer.notification.DisableAdminSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.DisableAdminSubscriptionResponse();
    return proto.bucketeer.notification.DisableAdminSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DisableAdminSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DisableAdminSubscriptionResponse}
 */
proto.bucketeer.notification.DisableAdminSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.DisableAdminSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DisableAdminSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DisableAdminSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DisableAdminSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.UpdateAdminSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.UpdateAdminSubscriptionRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          id: jspb.Message.getFieldWithDefault(msg, 1, ''),
          addSourceTypesCommand:
            (f = msg.getAddSourceTypesCommand()) &&
            proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand.toObject(
              includeInstance,
              f
            ),
          deleteSourceTypesCommand:
            (f = msg.getDeleteSourceTypesCommand()) &&
            proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand.toObject(
              includeInstance,
              f
            ),
          renameSubscriptionCommand:
            (f = msg.getRenameSubscriptionCommand()) &&
            proto_notification_command_pb.RenameAdminSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.UpdateAdminSubscriptionRequest();
    return proto.bucketeer.notification.UpdateAdminSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.deserializeBinaryFromReader =
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
          var value =
            new proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand
              .deserializeBinaryFromReader
          );
          msg.setAddSourceTypesCommand(value);
          break;
        case 3:
          var value =
            new proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb
              .DeleteAdminSubscriptionSourceTypesCommand
              .deserializeBinaryFromReader
          );
          msg.setDeleteSourceTypesCommand(value);
          break;
        case 4:
          var value =
            new proto_notification_command_pb.RenameAdminSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.RenameAdminSubscriptionCommand
              .deserializeBinaryFromReader
          );
          msg.setRenameSubscriptionCommand(value);
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
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.UpdateAdminSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(1, f);
    }
    f = message.getAddSourceTypesCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getDeleteSourceTypesCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getRenameSubscriptionCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_notification_command_pb.RenameAdminSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string id = 1;
 * @return {string}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 1, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
  };

/**
 * optional AddAdminSubscriptionSourceTypesCommand add_source_types_command = 2;
 * @return {?proto.bucketeer.notification.AddAdminSubscriptionSourceTypesCommand}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.getAddSourceTypesCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.AddAdminSubscriptionSourceTypesCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.AddAdminSubscriptionSourceTypesCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.AddAdminSubscriptionSourceTypesCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.setAddSourceTypesCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.clearAddSourceTypesCommand =
  function () {
    return this.setAddSourceTypesCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.hasAddSourceTypesCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional DeleteAdminSubscriptionSourceTypesCommand delete_source_types_command = 3;
 * @return {?proto.bucketeer.notification.DeleteAdminSubscriptionSourceTypesCommand}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.getDeleteSourceTypesCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.DeleteAdminSubscriptionSourceTypesCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.DeleteAdminSubscriptionSourceTypesCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.DeleteAdminSubscriptionSourceTypesCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.setDeleteSourceTypesCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.clearDeleteSourceTypesCommand =
  function () {
    return this.setDeleteSourceTypesCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.hasDeleteSourceTypesCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional RenameAdminSubscriptionCommand rename_subscription_command = 4;
 * @return {?proto.bucketeer.notification.RenameAdminSubscriptionCommand}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.getRenameSubscriptionCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.RenameAdminSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.RenameAdminSubscriptionCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.RenameAdminSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.setRenameSubscriptionCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.clearRenameSubscriptionCommand =
  function () {
    return this.setRenameSubscriptionCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionRequest.prototype.hasRenameSubscriptionCommand =
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
  proto.bucketeer.notification.UpdateAdminSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.UpdateAdminSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.UpdateAdminSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.UpdateAdminSubscriptionResponse.toObject =
    function (includeInstance, msg) {
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
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionResponse}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.UpdateAdminSubscriptionResponse();
    return proto.bucketeer.notification.UpdateAdminSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.UpdateAdminSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.UpdateAdminSubscriptionResponse}
 */
proto.bucketeer.notification.UpdateAdminSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.UpdateAdminSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.UpdateAdminSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.UpdateAdminSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.UpdateAdminSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.GetSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.GetSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.GetSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.GetSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
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
 * @return {!proto.bucketeer.notification.GetSubscriptionRequest}
 */
proto.bucketeer.notification.GetSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.GetSubscriptionRequest();
    return proto.bucketeer.notification.GetSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.GetSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.GetSubscriptionRequest}
 */
proto.bucketeer.notification.GetSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
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
proto.bucketeer.notification.GetSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.GetSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.GetSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.GetSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.notification.GetSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.GetSubscriptionRequest} returns this
 */
proto.bucketeer.notification.GetSubscriptionRequest.prototype.setId = function (
  value
) {
  return jspb.Message.setProto3StringField(this, 2, value);
};

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.notification.GetSubscriptionRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.GetSubscriptionRequest} returns this
 */
proto.bucketeer.notification.GetSubscriptionRequest.prototype.setEnvironmentId =
  function (value) {
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
  proto.bucketeer.notification.GetSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.GetSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.GetSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.GetSubscriptionResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        subscription:
          (f = msg.getSubscription()) &&
          proto_notification_subscription_pb.Subscription.toObject(
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
 * @return {!proto.bucketeer.notification.GetSubscriptionResponse}
 */
proto.bucketeer.notification.GetSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.GetSubscriptionResponse();
    return proto.bucketeer.notification.GetSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.GetSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.GetSubscriptionResponse}
 */
proto.bucketeer.notification.GetSubscriptionResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.setSubscription(value);
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
proto.bucketeer.notification.GetSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.GetSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.GetSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.GetSubscriptionResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscription();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
      );
    }
  };

/**
 * optional Subscription subscription = 1;
 * @return {?proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.GetSubscriptionResponse.prototype.getSubscription =
  function () {
    return /** @type{?proto.bucketeer.notification.Subscription} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.Subscription|undefined} value
 * @return {!proto.bucketeer.notification.GetSubscriptionResponse} returns this
 */
proto.bucketeer.notification.GetSubscriptionResponse.prototype.setSubscription =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.GetSubscriptionResponse} returns this
 */
proto.bucketeer.notification.GetSubscriptionResponse.prototype.clearSubscription =
  function () {
    return this.setSubscription(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.GetSubscriptionResponse.prototype.hasSubscription =
  function () {
    return jspb.Message.getField(this, 1) != null;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListSubscriptionsRequest.repeatedFields_ = [4, 11];

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
  proto.bucketeer.notification.ListSubscriptionsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListSubscriptionsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListSubscriptionsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListSubscriptionsRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
        cursor: jspb.Message.getFieldWithDefault(msg, 3, ''),
        sourceTypesList:
          (f = jspb.Message.getRepeatedField(msg, 4)) == null ? undefined : f,
        orderBy: jspb.Message.getFieldWithDefault(msg, 5, 0),
        orderDirection: jspb.Message.getFieldWithDefault(msg, 6, 0),
        searchKeyword: jspb.Message.getFieldWithDefault(msg, 7, ''),
        disabled:
          (f = msg.getDisabled()) &&
          google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
        environmentId: jspb.Message.getFieldWithDefault(msg, 9, ''),
        organizationId: jspb.Message.getFieldWithDefault(msg, 10, ''),
        environmentIdsList:
          (f = jspb.Message.getRepeatedField(msg, 11)) == null ? undefined : f
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
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.ListSubscriptionsRequest();
    return proto.bucketeer.notification.ListSubscriptionsRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListSubscriptionsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 4:
          var values =
            /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
              reader.isDelimited()
                ? reader.readPackedEnum()
                : [reader.readEnum()]
            );
          for (var i = 0; i < values.length; i++) {
            msg.addSourceTypes(values[i]);
          }
          break;
        case 5:
          var value =
            /** @type {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderBy} */ (
              reader.readEnum()
            );
          msg.setOrderBy(value);
          break;
        case 6:
          var value =
            /** @type {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderDirection} */ (
              reader.readEnum()
            );
          msg.setOrderDirection(value);
          break;
        case 7:
          var value = /** @type {string} */ (reader.readString());
          msg.setSearchKeyword(value);
          break;
        case 8:
          var value = new google_protobuf_wrappers_pb.BoolValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader
          );
          msg.setDisabled(value);
          break;
        case 9:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 10:
          var value = /** @type {string} */ (reader.readString());
          msg.setOrganizationId(value);
          break;
        case 11:
          var value = /** @type {string} */ (reader.readString());
          msg.addEnvironmentIds(value);
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
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListSubscriptionsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListSubscriptionsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListSubscriptionsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getPageSize();
    if (f !== 0) {
      writer.writeInt64(2, f);
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getSourceTypesList();
    if (f.length > 0) {
      writer.writePackedEnum(4, f);
    }
    f = message.getOrderBy();
    if (f !== 0.0) {
      writer.writeEnum(5, f);
    }
    f = message.getOrderDirection();
    if (f !== 0.0) {
      writer.writeEnum(6, f);
    }
    f = message.getSearchKeyword();
    if (f.length > 0) {
      writer.writeString(7, f);
    }
    f = message.getDisabled();
    if (f != null) {
      writer.writeMessage(
        8,
        f,
        google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(9, f);
    }
    f = message.getOrganizationId();
    if (f.length > 0) {
      writer.writeString(10, f);
    }
    f = message.getEnvironmentIdsList();
    if (f.length > 0) {
      writer.writeRepeatedString(11, f);
    }
  };

/**
 * @enum {number}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.OrderBy = {
  DEFAULT: 0,
  NAME: 1,
  CREATED_AT: 2,
  UPDATED_AT: 3,
  ENVIRONMENT: 4,
  STATE: 5
};

/**
 * @enum {number}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.OrderDirection = {
  ASC: 0,
  DESC: 1
};

/**
 * optional int64 page_size = 2;
 * @return {number}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setPageSize =
  function (value) {
    return jspb.Message.setProto3IntField(this, 2, value);
  };

/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * repeated Subscription.SourceType source_types = 4;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getSourceTypesList =
  function () {
    return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
      jspb.Message.getRepeatedField(this, 4)
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setSourceTypesList =
  function (value) {
    return jspb.Message.setField(this, 4, value || []);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.addSourceTypes =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 4, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.clearSourceTypesList =
  function () {
    return this.setSourceTypesList([]);
  };

/**
 * optional OrderBy order_by = 5;
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderBy}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getOrderBy =
  function () {
    return /** @type {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderBy} */ (
      jspb.Message.getFieldWithDefault(this, 5, 0)
    );
  };

/**
 * @param {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderBy} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setOrderBy =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 5, value);
  };

/**
 * optional OrderDirection order_direction = 6;
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderDirection}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getOrderDirection =
  function () {
    return /** @type {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderDirection} */ (
      jspb.Message.getFieldWithDefault(this, 6, 0)
    );
  };

/**
 * @param {!proto.bucketeer.notification.ListSubscriptionsRequest.OrderDirection} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setOrderDirection =
  function (value) {
    return jspb.Message.setProto3EnumField(this, 6, value);
  };

/**
 * optional string search_keyword = 7;
 * @return {string}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getSearchKeyword =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 7, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setSearchKeyword =
  function (value) {
    return jspb.Message.setProto3StringField(this, 7, value);
  };

/**
 * optional google.protobuf.BoolValue disabled = 8;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getDisabled =
  function () {
    return /** @type{?proto.google.protobuf.BoolValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.BoolValue,
        8
      )
    );
  };

/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setDisabled =
  function (value) {
    return jspb.Message.setWrapperField(this, 8, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.clearDisabled =
  function () {
    return this.setDisabled(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.hasDisabled =
  function () {
    return jspb.Message.getField(this, 8) != null;
  };

/**
 * optional string environment_id = 9;
 * @return {string}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 9, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 9, value);
  };

/**
 * optional string organization_id = 10;
 * @return {string}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getOrganizationId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 10, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setOrganizationId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 10, value);
  };

/**
 * repeated string environment_ids = 11;
 * @return {!Array<string>}
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.getEnvironmentIdsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 11)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.setEnvironmentIdsList =
  function (value) {
    return jspb.Message.setField(this, 11, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.addEnvironmentIds =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 11, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListSubscriptionsRequest.prototype.clearEnvironmentIdsList =
  function () {
    return this.setEnvironmentIdsList([]);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListSubscriptionsResponse.repeatedFields_ = [1];

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
  proto.bucketeer.notification.ListSubscriptionsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListSubscriptionsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListSubscriptionsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListSubscriptionsResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        subscriptionsList: jspb.Message.toObjectList(
          msg.getSubscriptionsList(),
          proto_notification_subscription_pb.Subscription.toObject,
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
 * @return {!proto.bucketeer.notification.ListSubscriptionsResponse}
 */
proto.bucketeer.notification.ListSubscriptionsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.ListSubscriptionsResponse();
    return proto.bucketeer.notification.ListSubscriptionsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListSubscriptionsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListSubscriptionsResponse}
 */
proto.bucketeer.notification.ListSubscriptionsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.addSubscriptions(value);
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
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListSubscriptionsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListSubscriptionsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListSubscriptionsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscriptionsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
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
 * repeated Subscription subscriptions = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription>}
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.getSubscriptionsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.notification.Subscription>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription>} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.setSubscriptionsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.addSubscriptions =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.notification.Subscription,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.clearSubscriptionsList =
  function () {
    return this.setSubscriptionsList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional int64 total_count = 3;
 * @return {number}
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.getTotalCount =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.notification.ListSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListSubscriptionsResponse.prototype.setTotalCount =
  function (value) {
    return jspb.Message.setProto3IntField(this, 3, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.repeatedFields_ = [
  4
];

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
  proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListEnabledSubscriptionsRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListEnabledSubscriptionsRequest.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          pageSize: jspb.Message.getFieldWithDefault(msg, 2, 0),
          cursor: jspb.Message.getFieldWithDefault(msg, 3, ''),
          sourceTypesList:
            (f = jspb.Message.getRepeatedField(msg, 4)) == null ? undefined : f,
          environmentId: jspb.Message.getFieldWithDefault(msg, 5, '')
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
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.ListEnabledSubscriptionsRequest();
    return proto.bucketeer.notification.ListEnabledSubscriptionsRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {number} */ (reader.readInt64());
          msg.setPageSize(value);
          break;
        case 3:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
          break;
        case 4:
          var values =
            /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
              reader.isDelimited()
                ? reader.readPackedEnum()
                : [reader.readEnum()]
            );
          for (var i = 0; i < values.length; i++) {
            msg.addSourceTypes(values[i]);
          }
          break;
        case 5:
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
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListEnabledSubscriptionsRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getPageSize();
    if (f !== 0) {
      writer.writeInt64(2, f);
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(3, f);
    }
    f = message.getSourceTypesList();
    if (f.length > 0) {
      writer.writePackedEnum(4, f);
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(5, f);
    }
  };

/**
 * optional int64 page_size = 2;
 * @return {number}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.getPageSize =
  function () {
    return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
  };

/**
 * @param {number} value
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.setPageSize =
  function (value) {
    return jspb.Message.setProto3IntField(this, 2, value);
  };

/**
 * optional string cursor = 3;
 * @return {string}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * repeated Subscription.SourceType source_types = 4;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.getSourceTypesList =
  function () {
    return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
      jspb.Message.getRepeatedField(this, 4)
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.setSourceTypesList =
  function (value) {
    return jspb.Message.setField(this, 4, value || []);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.addSourceTypes =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 4, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.clearSourceTypesList =
  function () {
    return this.setSourceTypesList([]);
  };

/**
 * optional string environment_id = 5;
 * @return {string}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 5, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsRequest} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 5, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.repeatedFields_ =
  [1];

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
  proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.ListEnabledSubscriptionsResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.ListEnabledSubscriptionsResponse.toObject =
    function (includeInstance, msg) {
      var f,
        obj = {
          subscriptionsList: jspb.Message.toObjectList(
            msg.getSubscriptionsList(),
            proto_notification_subscription_pb.Subscription.toObject,
            includeInstance
          ),
          cursor: jspb.Message.getFieldWithDefault(msg, 2, '')
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
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg =
      new proto.bucketeer.notification.ListEnabledSubscriptionsResponse();
    return proto.bucketeer.notification.ListEnabledSubscriptionsResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.addSubscriptions(value);
          break;
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setCursor(value);
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
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.ListEnabledSubscriptionsResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscriptionsList();
    if (f.length > 0) {
      writer.writeRepeatedMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
      );
    }
    f = message.getCursor();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
  };

/**
 * repeated Subscription subscriptions = 1;
 * @return {!Array<!proto.bucketeer.notification.Subscription>}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.getSubscriptionsList =
  function () {
    return /** @type{!Array<!proto.bucketeer.notification.Subscription>} */ (
      jspb.Message.getRepeatedWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription>} value
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.setSubscriptionsList =
  function (value) {
    return jspb.Message.setRepeatedWrapperField(this, 1, value);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription=} opt_value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.addSubscriptions =
  function (opt_value, opt_index) {
    return jspb.Message.addToRepeatedWrapperField(
      this,
      1,
      opt_value,
      proto.bucketeer.notification.Subscription,
      opt_index
    );
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.clearSubscriptionsList =
  function () {
    return this.setSubscriptionsList([]);
  };

/**
 * optional string cursor = 2;
 * @return {string}
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.getCursor =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.ListEnabledSubscriptionsResponse} returns this
 */
proto.bucketeer.notification.ListEnabledSubscriptionsResponse.prototype.setCursor =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.CreateSubscriptionRequest.repeatedFields_ = [5, 7];

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
  proto.bucketeer.notification.CreateSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.CreateSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.CreateSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.CreateSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        command:
          (f = msg.getCommand()) &&
          proto_notification_command_pb.CreateSubscriptionCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 3, ''),
        name: jspb.Message.getFieldWithDefault(msg, 4, ''),
        sourceTypesList:
          (f = jspb.Message.getRepeatedField(msg, 5)) == null ? undefined : f,
        recipient:
          (f = msg.getRecipient()) &&
          proto_notification_recipient_pb.Recipient.toObject(
            includeInstance,
            f
          ),
        featureFlagTagsList:
          (f = jspb.Message.getRepeatedField(msg, 7)) == null ? undefined : f
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
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.CreateSubscriptionRequest();
    return proto.bucketeer.notification.CreateSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.CreateSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value =
            new proto_notification_command_pb.CreateSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.CreateSubscriptionCommand
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
          var values =
            /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
              reader.isDelimited()
                ? reader.readPackedEnum()
                : [reader.readEnum()]
            );
          for (var i = 0; i < values.length; i++) {
            msg.addSourceTypes(values[i]);
          }
          break;
        case 6:
          var value = new proto_notification_recipient_pb.Recipient();
          reader.readMessage(
            value,
            proto_notification_recipient_pb.Recipient
              .deserializeBinaryFromReader
          );
          msg.setRecipient(value);
          break;
        case 7:
          var value = /** @type {string} */ (reader.readString());
          msg.addFeatureFlagTags(value);
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
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.CreateSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.CreateSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.CreateSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        2,
        f,
        proto_notification_command_pb.CreateSubscriptionCommand
          .serializeBinaryToWriter
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
    f = message.getSourceTypesList();
    if (f.length > 0) {
      writer.writePackedEnum(5, f);
    }
    f = message.getRecipient();
    if (f != null) {
      writer.writeMessage(
        6,
        f,
        proto_notification_recipient_pb.Recipient.serializeBinaryToWriter
      );
    }
    f = message.getFeatureFlagTagsList();
    if (f.length > 0) {
      writer.writeRepeatedString(7, f);
    }
  };

/**
 * optional CreateSubscriptionCommand command = 2;
 * @return {?proto.bucketeer.notification.CreateSubscriptionCommand}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.CreateSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.CreateSubscriptionCommand,
        2
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.CreateSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 2, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 2) != null;
  };

/**
 * optional string environment_id = 3;
 * @return {string}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 3, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
  };

/**
 * optional string name = 4;
 * @return {string}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.getName =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.setName =
  function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
  };

/**
 * repeated Subscription.SourceType source_types = 5;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.getSourceTypesList =
  function () {
    return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
      jspb.Message.getRepeatedField(this, 5)
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.setSourceTypesList =
  function (value) {
    return jspb.Message.setField(this, 5, value || []);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.addSourceTypes =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 5, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.clearSourceTypesList =
  function () {
    return this.setSourceTypesList([]);
  };

/**
 * optional Recipient recipient = 6;
 * @return {?proto.bucketeer.notification.Recipient}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.getRecipient =
  function () {
    return /** @type{?proto.bucketeer.notification.Recipient} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_recipient_pb.Recipient,
        6
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.Recipient|undefined} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.setRecipient =
  function (value) {
    return jspb.Message.setWrapperField(this, 6, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.clearRecipient =
  function () {
    return this.setRecipient(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.hasRecipient =
  function () {
    return jspb.Message.getField(this, 6) != null;
  };

/**
 * repeated string feature_flag_tags = 7;
 * @return {!Array<string>}
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.getFeatureFlagTagsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 7)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.setFeatureFlagTagsList =
  function (value) {
    return jspb.Message.setField(this, 7, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.addFeatureFlagTags =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 7, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.CreateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.CreateSubscriptionRequest.prototype.clearFeatureFlagTagsList =
  function () {
    return this.setFeatureFlagTagsList([]);
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
  proto.bucketeer.notification.CreateSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.CreateSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.CreateSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.CreateSubscriptionResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        subscription:
          (f = msg.getSubscription()) &&
          proto_notification_subscription_pb.Subscription.toObject(
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
 * @return {!proto.bucketeer.notification.CreateSubscriptionResponse}
 */
proto.bucketeer.notification.CreateSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.CreateSubscriptionResponse();
    return proto.bucketeer.notification.CreateSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.CreateSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.CreateSubscriptionResponse}
 */
proto.bucketeer.notification.CreateSubscriptionResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.setSubscription(value);
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
proto.bucketeer.notification.CreateSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.CreateSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.CreateSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.CreateSubscriptionResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscription();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
      );
    }
  };

/**
 * optional Subscription subscription = 1;
 * @return {?proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.CreateSubscriptionResponse.prototype.getSubscription =
  function () {
    return /** @type{?proto.bucketeer.notification.Subscription} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.Subscription|undefined} value
 * @return {!proto.bucketeer.notification.CreateSubscriptionResponse} returns this
 */
proto.bucketeer.notification.CreateSubscriptionResponse.prototype.setSubscription =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.CreateSubscriptionResponse} returns this
 */
proto.bucketeer.notification.CreateSubscriptionResponse.prototype.clearSubscription =
  function () {
    return this.setSubscription(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.CreateSubscriptionResponse.prototype.hasSubscription =
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
  proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DeleteSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DeleteSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DeleteSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_notification_command_pb.DeleteSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.DeleteSubscriptionRequest}
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.DeleteSubscriptionRequest();
    return proto.bucketeer.notification.DeleteSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DeleteSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DeleteSubscriptionRequest}
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value =
            new proto_notification_command_pb.DeleteSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.DeleteSubscriptionCommand
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
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DeleteSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DeleteSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_notification_command_pb.DeleteSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.DeleteSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional DeleteSubscriptionCommand command = 3;
 * @return {?proto.bucketeer.notification.DeleteSubscriptionCommand}
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.DeleteSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.DeleteSubscriptionCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.DeleteSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.DeleteSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.DeleteSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.DeleteSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DeleteSubscriptionRequest.prototype.setEnvironmentId =
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
  proto.bucketeer.notification.DeleteSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DeleteSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DeleteSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DeleteSubscriptionResponse.toObject = function (
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
 * @return {!proto.bucketeer.notification.DeleteSubscriptionResponse}
 */
proto.bucketeer.notification.DeleteSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.DeleteSubscriptionResponse();
    return proto.bucketeer.notification.DeleteSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DeleteSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DeleteSubscriptionResponse}
 */
proto.bucketeer.notification.DeleteSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.DeleteSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DeleteSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DeleteSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DeleteSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.EnableSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.EnableSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.EnableSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.EnableSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_notification_command_pb.EnableSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.EnableSubscriptionRequest}
 */
proto.bucketeer.notification.EnableSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.EnableSubscriptionRequest();
    return proto.bucketeer.notification.EnableSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.EnableSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.EnableSubscriptionRequest}
 */
proto.bucketeer.notification.EnableSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value =
            new proto_notification_command_pb.EnableSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.EnableSubscriptionCommand
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
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.EnableSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.EnableSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.EnableSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_notification_command_pb.EnableSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.EnableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional EnableSubscriptionCommand command = 3;
 * @return {?proto.bucketeer.notification.EnableSubscriptionCommand}
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.EnableSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.EnableSubscriptionCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.EnableSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.EnableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.EnableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.EnableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.EnableSubscriptionRequest.prototype.setEnvironmentId =
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
  proto.bucketeer.notification.EnableSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.EnableSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.EnableSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.EnableSubscriptionResponse.toObject = function (
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
 * @return {!proto.bucketeer.notification.EnableSubscriptionResponse}
 */
proto.bucketeer.notification.EnableSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.EnableSubscriptionResponse();
    return proto.bucketeer.notification.EnableSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.EnableSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.EnableSubscriptionResponse}
 */
proto.bucketeer.notification.EnableSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.EnableSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.EnableSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.EnableSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.EnableSubscriptionResponse.serializeBinaryToWriter =
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
  proto.bucketeer.notification.DisableSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DisableSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DisableSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DisableSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        command:
          (f = msg.getCommand()) &&
          proto_notification_command_pb.DisableSubscriptionCommand.toObject(
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
 * @return {!proto.bucketeer.notification.DisableSubscriptionRequest}
 */
proto.bucketeer.notification.DisableSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.DisableSubscriptionRequest();
    return proto.bucketeer.notification.DisableSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DisableSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DisableSubscriptionRequest}
 */
proto.bucketeer.notification.DisableSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value =
            new proto_notification_command_pb.DisableSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.DisableSubscriptionCommand
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
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DisableSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DisableSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DisableSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_notification_command_pb.DisableSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(4, f);
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.DisableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional DisableSubscriptionCommand command = 3;
 * @return {?proto.bucketeer.notification.DisableSubscriptionCommand}
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.getCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.DisableSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.DisableSubscriptionCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.DisableSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.DisableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.setCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.DisableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.clearCommand =
  function () {
    return this.setCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.hasCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional string environment_id = 4;
 * @return {string}
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 4, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.DisableSubscriptionRequest} returns this
 */
proto.bucketeer.notification.DisableSubscriptionRequest.prototype.setEnvironmentId =
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
  proto.bucketeer.notification.DisableSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.DisableSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.DisableSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.DisableSubscriptionResponse.toObject = function (
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
 * @return {!proto.bucketeer.notification.DisableSubscriptionResponse}
 */
proto.bucketeer.notification.DisableSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.DisableSubscriptionResponse();
    return proto.bucketeer.notification.DisableSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.DisableSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.DisableSubscriptionResponse}
 */
proto.bucketeer.notification.DisableSubscriptionResponse.deserializeBinaryFromReader =
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
proto.bucketeer.notification.DisableSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.DisableSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.DisableSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.DisableSubscriptionResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
  };

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.repeatedFields_ = [
  7, 10
];

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
  proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.UpdateSubscriptionRequest.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.UpdateSubscriptionRequest} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.UpdateSubscriptionRequest.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        id: jspb.Message.getFieldWithDefault(msg, 2, ''),
        addSourceTypesCommand:
          (f = msg.getAddSourceTypesCommand()) &&
          proto_notification_command_pb.AddSourceTypesCommand.toObject(
            includeInstance,
            f
          ),
        deleteSourceTypesCommand:
          (f = msg.getDeleteSourceTypesCommand()) &&
          proto_notification_command_pb.DeleteSourceTypesCommand.toObject(
            includeInstance,
            f
          ),
        renameSubscriptionCommand:
          (f = msg.getRenameSubscriptionCommand()) &&
          proto_notification_command_pb.RenameSubscriptionCommand.toObject(
            includeInstance,
            f
          ),
        environmentId: jspb.Message.getFieldWithDefault(msg, 6, ''),
        sourceTypesList:
          (f = jspb.Message.getRepeatedField(msg, 7)) == null ? undefined : f,
        name:
          (f = msg.getName()) &&
          google_protobuf_wrappers_pb.StringValue.toObject(includeInstance, f),
        disabled:
          (f = msg.getDisabled()) &&
          google_protobuf_wrappers_pb.BoolValue.toObject(includeInstance, f),
        featureFlagTagsList:
          (f = jspb.Message.getRepeatedField(msg, 10)) == null ? undefined : f,
        updateSubscriptionFeatureTagsCommand:
          (f = msg.getUpdateSubscriptionFeatureTagsCommand()) &&
          proto_notification_command_pb.UpdateSubscriptionFeatureFlagTagsCommand.toObject(
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
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.UpdateSubscriptionRequest();
    return proto.bucketeer.notification.UpdateSubscriptionRequest.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.UpdateSubscriptionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 2:
          var value = /** @type {string} */ (reader.readString());
          msg.setId(value);
          break;
        case 3:
          var value = new proto_notification_command_pb.AddSourceTypesCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.AddSourceTypesCommand
              .deserializeBinaryFromReader
          );
          msg.setAddSourceTypesCommand(value);
          break;
        case 4:
          var value =
            new proto_notification_command_pb.DeleteSourceTypesCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.DeleteSourceTypesCommand
              .deserializeBinaryFromReader
          );
          msg.setDeleteSourceTypesCommand(value);
          break;
        case 5:
          var value =
            new proto_notification_command_pb.RenameSubscriptionCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb.RenameSubscriptionCommand
              .deserializeBinaryFromReader
          );
          msg.setRenameSubscriptionCommand(value);
          break;
        case 6:
          var value = /** @type {string} */ (reader.readString());
          msg.setEnvironmentId(value);
          break;
        case 7:
          var values =
            /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
              reader.isDelimited()
                ? reader.readPackedEnum()
                : [reader.readEnum()]
            );
          for (var i = 0; i < values.length; i++) {
            msg.addSourceTypes(values[i]);
          }
          break;
        case 8:
          var value = new google_protobuf_wrappers_pb.StringValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.StringValue.deserializeBinaryFromReader
          );
          msg.setName(value);
          break;
        case 9:
          var value = new google_protobuf_wrappers_pb.BoolValue();
          reader.readMessage(
            value,
            google_protobuf_wrappers_pb.BoolValue.deserializeBinaryFromReader
          );
          msg.setDisabled(value);
          break;
        case 10:
          var value = /** @type {string} */ (reader.readString());
          msg.addFeatureFlagTags(value);
          break;
        case 11:
          var value =
            new proto_notification_command_pb.UpdateSubscriptionFeatureFlagTagsCommand();
          reader.readMessage(
            value,
            proto_notification_command_pb
              .UpdateSubscriptionFeatureFlagTagsCommand
              .deserializeBinaryFromReader
          );
          msg.setUpdateSubscriptionFeatureTagsCommand(value);
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
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.UpdateSubscriptionRequest.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.UpdateSubscriptionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getId();
    if (f.length > 0) {
      writer.writeString(2, f);
    }
    f = message.getAddSourceTypesCommand();
    if (f != null) {
      writer.writeMessage(
        3,
        f,
        proto_notification_command_pb.AddSourceTypesCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getDeleteSourceTypesCommand();
    if (f != null) {
      writer.writeMessage(
        4,
        f,
        proto_notification_command_pb.DeleteSourceTypesCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getRenameSubscriptionCommand();
    if (f != null) {
      writer.writeMessage(
        5,
        f,
        proto_notification_command_pb.RenameSubscriptionCommand
          .serializeBinaryToWriter
      );
    }
    f = message.getEnvironmentId();
    if (f.length > 0) {
      writer.writeString(6, f);
    }
    f = message.getSourceTypesList();
    if (f.length > 0) {
      writer.writePackedEnum(7, f);
    }
    f = message.getName();
    if (f != null) {
      writer.writeMessage(
        8,
        f,
        google_protobuf_wrappers_pb.StringValue.serializeBinaryToWriter
      );
    }
    f = message.getDisabled();
    if (f != null) {
      writer.writeMessage(
        9,
        f,
        google_protobuf_wrappers_pb.BoolValue.serializeBinaryToWriter
      );
    }
    f = message.getFeatureFlagTagsList();
    if (f.length > 0) {
      writer.writeRepeatedString(10, f);
    }
    f = message.getUpdateSubscriptionFeatureTagsCommand();
    if (f != null) {
      writer.writeMessage(
        11,
        f,
        proto_notification_command_pb.UpdateSubscriptionFeatureFlagTagsCommand
          .serializeBinaryToWriter
      );
    }
  };

/**
 * optional string id = 2;
 * @return {string}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 2, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
  };

/**
 * optional AddSourceTypesCommand add_source_types_command = 3;
 * @return {?proto.bucketeer.notification.AddSourceTypesCommand}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getAddSourceTypesCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.AddSourceTypesCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.AddSourceTypesCommand,
        3
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.AddSourceTypesCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setAddSourceTypesCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 3, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearAddSourceTypesCommand =
  function () {
    return this.setAddSourceTypesCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.hasAddSourceTypesCommand =
  function () {
    return jspb.Message.getField(this, 3) != null;
  };

/**
 * optional DeleteSourceTypesCommand delete_source_types_command = 4;
 * @return {?proto.bucketeer.notification.DeleteSourceTypesCommand}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getDeleteSourceTypesCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.DeleteSourceTypesCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.DeleteSourceTypesCommand,
        4
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.DeleteSourceTypesCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setDeleteSourceTypesCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearDeleteSourceTypesCommand =
  function () {
    return this.setDeleteSourceTypesCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.hasDeleteSourceTypesCommand =
  function () {
    return jspb.Message.getField(this, 4) != null;
  };

/**
 * optional RenameSubscriptionCommand rename_subscription_command = 5;
 * @return {?proto.bucketeer.notification.RenameSubscriptionCommand}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getRenameSubscriptionCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.RenameSubscriptionCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.RenameSubscriptionCommand,
        5
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.RenameSubscriptionCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setRenameSubscriptionCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 5, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearRenameSubscriptionCommand =
  function () {
    return this.setRenameSubscriptionCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.hasRenameSubscriptionCommand =
  function () {
    return jspb.Message.getField(this, 5) != null;
  };

/**
 * optional string environment_id = 6;
 * @return {string}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getEnvironmentId =
  function () {
    return /** @type {string} */ (
      jspb.Message.getFieldWithDefault(this, 6, '')
    );
  };

/**
 * @param {string} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setEnvironmentId =
  function (value) {
    return jspb.Message.setProto3StringField(this, 6, value);
  };

/**
 * repeated Subscription.SourceType source_types = 7;
 * @return {!Array<!proto.bucketeer.notification.Subscription.SourceType>}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getSourceTypesList =
  function () {
    return /** @type {!Array<!proto.bucketeer.notification.Subscription.SourceType>} */ (
      jspb.Message.getRepeatedField(this, 7)
    );
  };

/**
 * @param {!Array<!proto.bucketeer.notification.Subscription.SourceType>} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setSourceTypesList =
  function (value) {
    return jspb.Message.setField(this, 7, value || []);
  };

/**
 * @param {!proto.bucketeer.notification.Subscription.SourceType} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.addSourceTypes =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 7, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearSourceTypesList =
  function () {
    return this.setSourceTypesList([]);
  };

/**
 * optional google.protobuf.StringValue name = 8;
 * @return {?proto.google.protobuf.StringValue}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getName =
  function () {
    return /** @type{?proto.google.protobuf.StringValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.StringValue,
        8
      )
    );
  };

/**
 * @param {?proto.google.protobuf.StringValue|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setName =
  function (value) {
    return jspb.Message.setWrapperField(this, 8, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearName =
  function () {
    return this.setName(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.hasName =
  function () {
    return jspb.Message.getField(this, 8) != null;
  };

/**
 * optional google.protobuf.BoolValue disabled = 9;
 * @return {?proto.google.protobuf.BoolValue}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getDisabled =
  function () {
    return /** @type{?proto.google.protobuf.BoolValue} */ (
      jspb.Message.getWrapperField(
        this,
        google_protobuf_wrappers_pb.BoolValue,
        9
      )
    );
  };

/**
 * @param {?proto.google.protobuf.BoolValue|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setDisabled =
  function (value) {
    return jspb.Message.setWrapperField(this, 9, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearDisabled =
  function () {
    return this.setDisabled(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.hasDisabled =
  function () {
    return jspb.Message.getField(this, 9) != null;
  };

/**
 * repeated string feature_flag_tags = 10;
 * @return {!Array<string>}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getFeatureFlagTagsList =
  function () {
    return /** @type {!Array<string>} */ (
      jspb.Message.getRepeatedField(this, 10)
    );
  };

/**
 * @param {!Array<string>} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setFeatureFlagTagsList =
  function (value) {
    return jspb.Message.setField(this, 10, value || []);
  };

/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.addFeatureFlagTags =
  function (value, opt_index) {
    return jspb.Message.addToRepeatedField(this, 10, value, opt_index);
  };

/**
 * Clears the list making it empty but non-null.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearFeatureFlagTagsList =
  function () {
    return this.setFeatureFlagTagsList([]);
  };

/**
 * optional UpdateSubscriptionFeatureFlagTagsCommand update_subscription_feature_tags_command = 11;
 * @return {?proto.bucketeer.notification.UpdateSubscriptionFeatureFlagTagsCommand}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.getUpdateSubscriptionFeatureTagsCommand =
  function () {
    return /** @type{?proto.bucketeer.notification.UpdateSubscriptionFeatureFlagTagsCommand} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_command_pb.UpdateSubscriptionFeatureFlagTagsCommand,
        11
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.UpdateSubscriptionFeatureFlagTagsCommand|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.setUpdateSubscriptionFeatureTagsCommand =
  function (value) {
    return jspb.Message.setWrapperField(this, 11, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionRequest} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.clearUpdateSubscriptionFeatureTagsCommand =
  function () {
    return this.setUpdateSubscriptionFeatureTagsCommand(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionRequest.prototype.hasUpdateSubscriptionFeatureTagsCommand =
  function () {
    return jspb.Message.getField(this, 11) != null;
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
  proto.bucketeer.notification.UpdateSubscriptionResponse.prototype.toObject =
    function (opt_includeInstance) {
      return proto.bucketeer.notification.UpdateSubscriptionResponse.toObject(
        opt_includeInstance,
        this
      );
    };

  /**
   * Static version of the {@see toObject} method.
   * @param {boolean|undefined} includeInstance Deprecated. Whether to include
   *     the JSPB instance for transitional soy proto support:
   *     http://goto/soy-param-migration
   * @param {!proto.bucketeer.notification.UpdateSubscriptionResponse} msg The msg instance to transform.
   * @return {!Object}
   * @suppress {unusedLocalVariables} f is only used for nested messages
   */
  proto.bucketeer.notification.UpdateSubscriptionResponse.toObject = function (
    includeInstance,
    msg
  ) {
    var f,
      obj = {
        subscription:
          (f = msg.getSubscription()) &&
          proto_notification_subscription_pb.Subscription.toObject(
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
 * @return {!proto.bucketeer.notification.UpdateSubscriptionResponse}
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.deserializeBinary =
  function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.bucketeer.notification.UpdateSubscriptionResponse();
    return proto.bucketeer.notification.UpdateSubscriptionResponse.deserializeBinaryFromReader(
      msg,
      reader
    );
  };

/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.bucketeer.notification.UpdateSubscriptionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionResponse}
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.deserializeBinaryFromReader =
  function (msg, reader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      var field = reader.getFieldNumber();
      switch (field) {
        case 1:
          var value = new proto_notification_subscription_pb.Subscription();
          reader.readMessage(
            value,
            proto_notification_subscription_pb.Subscription
              .deserializeBinaryFromReader
          );
          msg.setSubscription(value);
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
proto.bucketeer.notification.UpdateSubscriptionResponse.prototype.serializeBinary =
  function () {
    var writer = new jspb.BinaryWriter();
    proto.bucketeer.notification.UpdateSubscriptionResponse.serializeBinaryToWriter(
      this,
      writer
    );
    return writer.getResultBuffer();
  };

/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.bucketeer.notification.UpdateSubscriptionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.serializeBinaryToWriter =
  function (message, writer) {
    var f = undefined;
    f = message.getSubscription();
    if (f != null) {
      writer.writeMessage(
        1,
        f,
        proto_notification_subscription_pb.Subscription.serializeBinaryToWriter
      );
    }
  };

/**
 * optional Subscription subscription = 1;
 * @return {?proto.bucketeer.notification.Subscription}
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.prototype.getSubscription =
  function () {
    return /** @type{?proto.bucketeer.notification.Subscription} */ (
      jspb.Message.getWrapperField(
        this,
        proto_notification_subscription_pb.Subscription,
        1
      )
    );
  };

/**
 * @param {?proto.bucketeer.notification.Subscription|undefined} value
 * @return {!proto.bucketeer.notification.UpdateSubscriptionResponse} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.prototype.setSubscription =
  function (value) {
    return jspb.Message.setWrapperField(this, 1, value);
  };

/**
 * Clears the message field making it undefined.
 * @return {!proto.bucketeer.notification.UpdateSubscriptionResponse} returns this
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.prototype.clearSubscription =
  function () {
    return this.setSubscription(undefined);
  };

/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.bucketeer.notification.UpdateSubscriptionResponse.prototype.hasSubscription =
  function () {
    return jspb.Message.getField(this, 1) != null;
  };

goog.object.extend(exports, proto.bucketeer.notification);
