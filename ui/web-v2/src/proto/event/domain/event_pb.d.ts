// package: bucketeer.event.domain
// file: proto/event/domain/event.proto

import * as jspb from 'google-protobuf';
import * as google_protobuf_any_pb from 'google-protobuf/google/protobuf/any_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as proto_feature_clause_pb from '../../../proto/feature/clause_pb';
import * as proto_feature_feature_pb from '../../../proto/feature/feature_pb';
import * as proto_feature_rule_pb from '../../../proto/feature/rule_pb';
import * as proto_feature_variation_pb from '../../../proto/feature/variation_pb';
import * as proto_feature_strategy_pb from '../../../proto/feature/strategy_pb';
import * as proto_feature_segment_pb from '../../../proto/feature/segment_pb';
import * as proto_feature_target_pb from '../../../proto/feature/target_pb';
import * as proto_account_account_pb from '../../../proto/account/account_pb';
import * as proto_account_api_key_pb from '../../../proto/account/api_key_pb';
import * as proto_autoops_auto_ops_rule_pb from '../../../proto/autoops/auto_ops_rule_pb';
import * as proto_autoops_clause_pb from '../../../proto/autoops/clause_pb';
import * as proto_notification_subscription_pb from '../../../proto/notification/subscription_pb';
import * as proto_notification_recipient_pb from '../../../proto/notification/recipient_pb';
import * as proto_feature_prerequisite_pb from '../../../proto/feature/prerequisite_pb';
import * as proto_autoops_progressive_rollout_pb from '../../../proto/autoops/progressive_rollout_pb';
import * as proto_feature_flag_trigger_pb from '../../../proto/feature/flag_trigger_pb';
import * as proto_account_search_filter_pb from '../../../proto/account/search_filter_pb';
import * as proto_tag_tag_pb from '../../../proto/tag/tag_pb';

export class Event extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  getEntityType(): Event.EntityTypeMap[keyof Event.EntityTypeMap];
  setEntityType(value: Event.EntityTypeMap[keyof Event.EntityTypeMap]): void;

  getEntityId(): string;
  setEntityId(value: string): void;

  getType(): Event.TypeMap[keyof Event.TypeMap];
  setType(value: Event.TypeMap[keyof Event.TypeMap]): void;

  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): Editor | undefined;
  setEditor(value?: Editor): void;

  hasData(): boolean;
  clearData(): void;
  getData(): google_protobuf_any_pb.Any | undefined;
  setData(value?: google_protobuf_any_pb.Any): void;

  getIsAdminEvent(): boolean;
  setIsAdminEvent(value: boolean): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): Options | undefined;
  setOptions(value?: Options): void;

  getEntityData(): string;
  setEntityData(value: string): void;

  getPreviousEntityData(): string;
  setPreviousEntityData(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Event,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(
    message: Event,
    reader: jspb.BinaryReader
  ): Event;
}

export namespace Event {
  export type AsObject = {
    id: string;
    timestamp: number;
    entityType: Event.EntityTypeMap[keyof Event.EntityTypeMap];
    entityId: string;
    type: Event.TypeMap[keyof Event.TypeMap];
    editor?: Editor.AsObject;
    data?: google_protobuf_any_pb.Any.AsObject;
    isAdminEvent: boolean;
    options?: Options.AsObject;
    entityData: string;
    previousEntityData: string;
    environmentId: string;
  };

  export interface EntityTypeMap {
    FEATURE: 0;
    GOAL: 1;
    EXPERIMENT: 2;
    ACCOUNT: 3;
    APIKEY: 4;
    SEGMENT: 5;
    ENVIRONMENT: 6;
    ADMIN_ACCOUNT: 7;
    AUTOOPS_RULE: 8;
    PUSH: 9;
    SUBSCRIPTION: 10;
    ADMIN_SUBSCRIPTION: 11;
    PROJECT: 12;
    PROGRESSIVE_ROLLOUT: 14;
    ORGANIZATION: 15;
    FLAG_TRIGGER: 16;
    TAG: 17;
  }

  export const EntityType: EntityTypeMap;

  export interface TypeMap {
    UNKNOWN: 0;
    FEATURE_CREATED: 1;
    FEATURE_RENAMED: 2;
    FEATURE_ENABLED: 3;
    FEATURE_DISABLED: 4;
    FEATURE_DELETED: 5;
    FEATURE_DESCRIPTION_CHANGED: 8;
    FEATURE_VARIATION_ADDED: 9;
    FEATURE_VARIATION_REMOVED: 10;
    FEATURE_OFF_VARIATION_CHANGED: 11;
    VARIATION_VALUE_CHANGED: 12;
    VARIATION_NAME_CHANGED: 13;
    VARIATION_DESCRIPTION_CHANGED: 14;
    VARIATION_USER_ADDED: 15;
    VARIATION_USER_REMOVED: 16;
    FEATURE_RULE_ADDED: 17;
    FEATURE_RULE_STRATEGY_CHANGED: 18;
    FEATURE_RULE_DELETED: 19;
    RULE_CLAUSE_ADDED: 20;
    RULE_CLAUSE_DELETED: 21;
    RULE_FIXED_STRATEGY_CHANGED: 22;
    RULE_ROLLOUT_STRATEGY_CHANGED: 23;
    CLAUSE_ATTRIBUTE_CHANGED: 24;
    CLAUSE_OPERATOR_CHANGED: 25;
    CLAUSE_VALUE_ADDED: 26;
    CLAUSE_VALUE_REMOVED: 27;
    FEATURE_DEFAULT_STRATEGY_CHANGED: 28;
    FEATURE_TAG_ADDED: 29;
    FEATURE_TAG_REMOVED: 30;
    FEATURE_VERSION_INCREMENTED: 31;
    FEATURE_ARCHIVED: 32;
    FEATURE_CLONED: 33;
    FEATURE_UNARCHIVED: 35;
    SAMPLING_SEED_RESET: 34;
    PREREQUISITE_ADDED: 36;
    PREREQUISITE_REMOVED: 37;
    PREREQUISITE_VARIATION_CHANGED: 38;
    FEATURE_RULES_ORDER_CHANGED: 39;
    FEATURE_UPDATED: 40;
    GOAL_CREATED: 100;
    GOAL_RENAMED: 101;
    GOAL_DESCRIPTION_CHANGED: 102;
    GOAL_DELETED: 103;
    GOAL_ARCHIVED: 104;
    EXPERIMENT_CREATED: 200;
    EXPERIMENT_STOPPED: 201;
    EXPERIMENT_START_AT_CHANGED: 202;
    EXPERIMENT_STOP_AT_CHANGED: 203;
    EXPERIMENT_DELETED: 204;
    EXPERIMENT_PERIOD_CHANGED: 205;
    EXPERIMENT_NAME_CHANGED: 206;
    EXPERIMENT_DESCRIPTION_CHANGED: 207;
    EXPERIMENT_STARTED: 208;
    EXPERIMENT_FINISHED: 209;
    EXPERIMENT_ARCHIVED: 210;
    ACCOUNT_CREATED: 300;
    ACCOUNT_ROLE_CHANGED: 301;
    ACCOUNT_ENABLED: 302;
    ACCOUNT_DISABLED: 303;
    ACCOUNT_DELETED: 304;
    ACCOUNT_V2_CREATED: 305;
    ACCOUNT_V2_NAME_CHANGED: 306;
    ACCOUNT_V2_AVATAR_IMAGE_URL_CHANGED: 307;
    ACCOUNT_V2_ORGANIZATION_ROLE_CHANGED: 308;
    ACCOUNT_V2_ENVIRONMENT_ROLES_CHANGED: 309;
    ACCOUNT_V2_ENABLED: 310;
    ACCOUNT_V2_DISABLED: 311;
    ACCOUNT_V2_DELETED: 312;
    ACCOUNT_V2_CREATED_SEARCH_FILTER: 313;
    ACCOUNT_V2_SEARCH_FILTER_NANE_CHANGED: 314;
    ACCOUNT_V2_SEARCH_FILTER_QUERY_CHANGED: 315;
    ACCOUNT_V2_SEARCH_FILTER_DEFAULT_CHANGED: 316;
    ACCOUNT_V2_SEARCH_FILTER_DELETED: 317;
    ACCOUNT_V2_FIRST_NAME_CHANGED: 318;
    ACCOUNT_V2_LAST_NAME_CHANGED: 319;
    ACCOUNT_V2_LANGUAGE_CHANGED: 320;
    ACCOUNT_V2_UPDATED: 321;
    ACCOUNT_V2_TAGS_CHANGED: 322;
    APIKEY_CREATED: 400;
    APIKEY_NAME_CHANGED: 401;
    APIKEY_ENABLED: 402;
    APIKEY_DISABLED: 403;
    APIKEY_CHANGED: 404;
    SEGMENT_CREATED: 500;
    SEGMENT_DELETED: 501;
    SEGMENT_NAME_CHANGED: 502;
    SEGMENT_DESCRIPTION_CHANGED: 503;
    SEGMENT_RULE_ADDED: 504;
    SEGMENT_RULE_DELETED: 505;
    SEGMENT_RULE_CLAUSE_ADDED: 506;
    SEGMENT_RULE_CLAUSE_DELETED: 507;
    SEGMENT_CLAUSE_ATTRIBUTE_CHANGED: 508;
    SEGMENT_CLAUSE_OPERATOR_CHANGED: 509;
    SEGMENT_CLAUSE_VALUE_ADDED: 510;
    SEGMENT_CLAUSE_VALUE_REMOVED: 511;
    SEGMENT_USER_ADDED: 512;
    SEGMENT_USER_DELETED: 513;
    SEGMENT_BULK_UPLOAD_USERS: 514;
    SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED: 515;
    SEGMENT_UPDATED: 516;
    ENVIRONMENT_CREATED: 600;
    ENVIRONMENT_RENAMED: 601;
    ENVIRONMENT_DESCRIPTION_CHANGED: 602;
    ENVIRONMENT_DELETED: 603;
    ENVIRONMENT_V2_CREATED: 604;
    ENVIRONMENT_V2_RENAMED: 605;
    ENVIRONMENT_V2_DESCRIPTION_CHANGED: 606;
    ENVIRONMENT_V2_ARCHIVED: 607;
    ENVIRONMENT_V2_UNARCHIVED: 608;
    ENVIRONMENT_V2_REQUIRE_COMMENT_CHANGED: 609;
    ADMIN_ACCOUNT_CREATED: 700;
    ADMIN_ACCOUNT_ENABLED: 702;
    ADMIN_ACCOUNT_DISABLED: 703;
    AUTOOPS_RULE_CREATED: 800;
    AUTOOPS_RULE_DELETED: 801;
    AUTOOPS_RULE_OPS_TYPE_CHANGED: 802;
    AUTOOPS_RULE_CLAUSE_DELETED: 803;
    AUTOOPS_RULE_TRIGGERED_AT_CHANGED: 804;
    OPS_EVENT_RATE_CLAUSE_ADDED: 805;
    OPS_EVENT_RATE_CLAUSE_CHANGED: 806;
    DATETIME_CLAUSE_ADDED: 807;
    DATETIME_CLAUSE_CHANGED: 808;
    AUTOOPS_RULE_STOPPED: 809;
    AUTOOPS_RULE_OPS_STATUS_CHANGED: 810;
    PUSH_CREATED: 900;
    PUSH_DELETED: 901;
    PUSH_TAGS_ADDED: 902;
    PUSH_TAGS_DELETED: 903;
    PUSH_RENAMED: 904;
    PUSH_UPDATED: 905;
    SUBSCRIPTION_CREATED: 1000;
    SUBSCRIPTION_DELETED: 1001;
    SUBSCRIPTION_ENABLED: 1002;
    SUBSCRIPTION_DISABLED: 1003;
    SUBSCRIPTION_SOURCE_TYPE_ADDED: 1004;
    SUBSCRIPTION_SOURCE_TYPE_DELETED: 1005;
    SUBSCRIPTION_RENAMED: 1006;
    SUBSCRIPTION_UPDATED: 1007;
    ADMIN_SUBSCRIPTION_CREATED: 1100;
    ADMIN_SUBSCRIPTION_DELETED: 1101;
    ADMIN_SUBSCRIPTION_ENABLED: 1102;
    ADMIN_SUBSCRIPTION_DISABLED: 1103;
    ADMIN_SUBSCRIPTION_SOURCE_TYPE_ADDED: 1104;
    ADMIN_SUBSCRIPTION_SOURCE_TYPE_DELETED: 1105;
    ADMIN_SUBSCRIPTION_RENAMED: 1106;
    PROJECT_CREATED: 1200;
    PROJECT_DESCRIPTION_CHANGED: 1201;
    PROJECT_ENABLED: 1202;
    PROJECT_DISABLED: 1203;
    PROJECT_TRIAL_CREATED: 1204;
    PROJECT_TRIAL_CONVERTED: 1205;
    PROJECT_RENAMED: 1206;
    PROGRESSIVE_ROLLOUT_CREATED: 1400;
    PROGRESSIVE_ROLLOUT_DELETED: 1401;
    PROGRESSIVE_ROLLOUT_SCHEDULE_TRIGGERED_AT_CHANGED: 1402;
    PROGRESSIVE_ROLLOUT_STOPPED: 1403;
    ORGANIZATION_CREATED: 1500;
    ORGANIZATION_NAME_CHANGED: 1501;
    ORGANIZATION_DESCRIPTION_CHANGED: 1502;
    ORGANIZATION_ENABLED: 1503;
    ORGANIZATION_DISABLED: 1504;
    ORGANIZATION_ARCHIVED: 1505;
    ORGANIZATION_UNARCHIVED: 1506;
    ORGANIZATION_TRIAL_CONVERTED: 1507;
    ORGANIZATION_OWNER_EMAIL_CHANGED: 1508;
    FLAG_TRIGGER_CREATED: 1601;
    FLAG_TRIGGER_RESET: 1602;
    FLAG_TRIGGER_DESCRIPTION_CHANGED: 1603;
    FLAG_TRIGGER_DISABLED: 1604;
    FLAG_TRIGGER_ENABLED: 1605;
    FLAG_TRIGGER_DELETED: 1606;
    FLAG_TRIGGER_USAGE_UPDATED: 1607;
    TAG_CREATED: 1701;
    TAG_DELETED: 1702;
  }

  export const Type: TypeMap;
}

export class Editor extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getIsAdmin(): boolean;
  setIsAdmin(value: boolean): void;

  getName(): string;
  setName(value: string): void;

  hasPublicApiEditor(): boolean;
  clearPublicApiEditor(): void;
  getPublicApiEditor(): Editor.PublicAPIEditor | undefined;
  setPublicApiEditor(value?: Editor.PublicAPIEditor): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Editor.AsObject;
  static toObject(includeInstance: boolean, msg: Editor): Editor.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Editor,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Editor;
  static deserializeBinaryFromReader(
    message: Editor,
    reader: jspb.BinaryReader
  ): Editor;
}

export namespace Editor {
  export type AsObject = {
    email: string;
    isAdmin: boolean;
    name: string;
    publicApiEditor?: Editor.PublicAPIEditor.AsObject;
  };

  export class PublicAPIEditor extends jspb.Message {
    getToken(): string;
    setToken(value: string): void;

    getMaintainer(): string;
    setMaintainer(value: string): void;

    getName(): string;
    setName(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PublicAPIEditor.AsObject;
    static toObject(
      includeInstance: boolean,
      msg: PublicAPIEditor
    ): PublicAPIEditor.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: PublicAPIEditor,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): PublicAPIEditor;
    static deserializeBinaryFromReader(
      message: PublicAPIEditor,
      reader: jspb.BinaryReader
    ): PublicAPIEditor;
  }

  export namespace PublicAPIEditor {
    export type AsObject = {
      token: string;
      maintainer: string;
      name: string;
    };
  }
}

export class Options extends jspb.Message {
  getComment(): string;
  setComment(value: string): void;

  getNewVersion(): number;
  setNewVersion(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Options.AsObject;
  static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Options,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Options;
  static deserializeBinaryFromReader(
    message: Options,
    reader: jspb.BinaryReader
  ): Options;
}

export namespace Options {
  export type AsObject = {
    comment: string;
    newVersion: number;
  };
}

export class FeatureCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getUser(): string;
  setUser(value: string): void;

  clearVariationsList(): void;
  getVariationsList(): Array<proto_feature_variation_pb.Variation>;
  setVariationsList(value: Array<proto_feature_variation_pb.Variation>): void;
  addVariations(
    value?: proto_feature_variation_pb.Variation,
    index?: number
  ): proto_feature_variation_pb.Variation;

  hasDefaultOnVariationIndex(): boolean;
  clearDefaultOnVariationIndex(): void;
  getDefaultOnVariationIndex():
    | google_protobuf_wrappers_pb.Int32Value
    | undefined;
  setDefaultOnVariationIndex(
    value?: google_protobuf_wrappers_pb.Int32Value
  ): void;

  hasDefaultOffVariationIndex(): boolean;
  clearDefaultOffVariationIndex(): void;
  getDefaultOffVariationIndex():
    | google_protobuf_wrappers_pb.Int32Value
    | undefined;
  setDefaultOffVariationIndex(
    value?: google_protobuf_wrappers_pb.Int32Value
  ): void;

  getVariationType(): proto_feature_feature_pb.Feature.VariationTypeMap[keyof proto_feature_feature_pb.Feature.VariationTypeMap];
  setVariationType(
    value: proto_feature_feature_pb.Feature.VariationTypeMap[keyof proto_feature_feature_pb.Feature.VariationTypeMap]
  ): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  clearPrerequisitesList(): void;
  getPrerequisitesList(): Array<proto_feature_prerequisite_pb.Prerequisite>;
  setPrerequisitesList(
    value: Array<proto_feature_prerequisite_pb.Prerequisite>
  ): void;
  addPrerequisites(
    value?: proto_feature_prerequisite_pb.Prerequisite,
    index?: number
  ): proto_feature_prerequisite_pb.Prerequisite;

  clearRulesList(): void;
  getRulesList(): Array<proto_feature_rule_pb.Rule>;
  setRulesList(value: Array<proto_feature_rule_pb.Rule>): void;
  addRules(
    value?: proto_feature_rule_pb.Rule,
    index?: number
  ): proto_feature_rule_pb.Rule;

  clearTargetsList(): void;
  getTargetsList(): Array<proto_feature_target_pb.Target>;
  setTargetsList(value: Array<proto_feature_target_pb.Target>): void;
  addTargets(
    value?: proto_feature_target_pb.Target,
    index?: number
  ): proto_feature_target_pb.Target;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureCreatedEvent
  ): FeatureCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureCreatedEvent;
  static deserializeBinaryFromReader(
    message: FeatureCreatedEvent,
    reader: jspb.BinaryReader
  ): FeatureCreatedEvent;
}

export namespace FeatureCreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    user: string;
    variationsList: Array<proto_feature_variation_pb.Variation.AsObject>;
    defaultOnVariationIndex?: google_protobuf_wrappers_pb.Int32Value.AsObject;
    defaultOffVariationIndex?: google_protobuf_wrappers_pb.Int32Value.AsObject;
    variationType: proto_feature_feature_pb.Feature.VariationTypeMap[keyof proto_feature_feature_pb.Feature.VariationTypeMap];
    tagsList: Array<string>;
    prerequisitesList: Array<proto_feature_prerequisite_pb.Prerequisite.AsObject>;
    rulesList: Array<proto_feature_rule_pb.Rule.AsObject>;
    targetsList: Array<proto_feature_target_pb.Target.AsObject>;
  };
}

export class FeatureUpdatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getData(): string;
  setData(value: string): void;

  getPreviousData(): string;
  setPreviousData(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureUpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureUpdatedEvent
  ): FeatureUpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureUpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureUpdatedEvent;
  static deserializeBinaryFromReader(
    message: FeatureUpdatedEvent,
    reader: jspb.BinaryReader
  ): FeatureUpdatedEvent;
}

export namespace FeatureUpdatedEvent {
  export type AsObject = {
    id: string;
    data: string;
    previousData: string;
  };
}

export class FeatureEnabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureEnabledEvent
  ): FeatureEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureEnabledEvent;
  static deserializeBinaryFromReader(
    message: FeatureEnabledEvent,
    reader: jspb.BinaryReader
  ): FeatureEnabledEvent;
}

export namespace FeatureEnabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class FeatureDisabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureDisabledEvent
  ): FeatureDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureDisabledEvent;
  static deserializeBinaryFromReader(
    message: FeatureDisabledEvent,
    reader: jspb.BinaryReader
  ): FeatureDisabledEvent;
}

export namespace FeatureDisabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class FeatureArchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureArchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureArchivedEvent
  ): FeatureArchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureArchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureArchivedEvent;
  static deserializeBinaryFromReader(
    message: FeatureArchivedEvent,
    reader: jspb.BinaryReader
  ): FeatureArchivedEvent;
}

export namespace FeatureArchivedEvent {
  export type AsObject = {
    id: string;
  };
}

export class FeatureUnarchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureUnarchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureUnarchivedEvent
  ): FeatureUnarchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureUnarchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureUnarchivedEvent;
  static deserializeBinaryFromReader(
    message: FeatureUnarchivedEvent,
    reader: jspb.BinaryReader
  ): FeatureUnarchivedEvent;
}

export namespace FeatureUnarchivedEvent {
  export type AsObject = {
    id: string;
  };
}

export class FeatureDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureDeletedEvent
  ): FeatureDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureDeletedEvent;
  static deserializeBinaryFromReader(
    message: FeatureDeletedEvent,
    reader: jspb.BinaryReader
  ): FeatureDeletedEvent;
}

export namespace FeatureDeletedEvent {
  export type AsObject = {
    id: string;
  };
}

export class EvaluationDelayableSetEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluationDelayableSetEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EvaluationDelayableSetEvent
  ): EvaluationDelayableSetEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EvaluationDelayableSetEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EvaluationDelayableSetEvent;
  static deserializeBinaryFromReader(
    message: EvaluationDelayableSetEvent,
    reader: jspb.BinaryReader
  ): EvaluationDelayableSetEvent;
}

export namespace EvaluationDelayableSetEvent {
  export type AsObject = {
    id: string;
  };
}

export class EvaluationUndelayableSetEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluationUndelayableSetEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EvaluationUndelayableSetEvent
  ): EvaluationUndelayableSetEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EvaluationUndelayableSetEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EvaluationUndelayableSetEvent;
  static deserializeBinaryFromReader(
    message: EvaluationUndelayableSetEvent,
    reader: jspb.BinaryReader
  ): EvaluationUndelayableSetEvent;
}

export namespace EvaluationUndelayableSetEvent {
  export type AsObject = {
    id: string;
  };
}

export class FeatureRenamedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureRenamedEvent
  ): FeatureRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureRenamedEvent;
  static deserializeBinaryFromReader(
    message: FeatureRenamedEvent,
    reader: jspb.BinaryReader
  ): FeatureRenamedEvent;
}

export namespace FeatureRenamedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class FeatureDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureDescriptionChangedEvent
  ): FeatureDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: FeatureDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): FeatureDescriptionChangedEvent;
}

export namespace FeatureDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class FeatureOffVariationChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOffVariation(): string;
  setOffVariation(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureOffVariationChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureOffVariationChangedEvent
  ): FeatureOffVariationChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureOffVariationChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureOffVariationChangedEvent;
  static deserializeBinaryFromReader(
    message: FeatureOffVariationChangedEvent,
    reader: jspb.BinaryReader
  ): FeatureOffVariationChangedEvent;
}

export namespace FeatureOffVariationChangedEvent {
  export type AsObject = {
    id: string;
    offVariation: string;
  };
}

export class FeatureVariationAddedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasVariation(): boolean;
  clearVariation(): void;
  getVariation(): proto_feature_variation_pb.Variation | undefined;
  setVariation(value?: proto_feature_variation_pb.Variation): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureVariationAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureVariationAddedEvent
  ): FeatureVariationAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureVariationAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureVariationAddedEvent;
  static deserializeBinaryFromReader(
    message: FeatureVariationAddedEvent,
    reader: jspb.BinaryReader
  ): FeatureVariationAddedEvent;
}

export namespace FeatureVariationAddedEvent {
  export type AsObject = {
    id: string;
    variation?: proto_feature_variation_pb.Variation.AsObject;
  };
}

export class FeatureVariationRemovedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureVariationRemovedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureVariationRemovedEvent
  ): FeatureVariationRemovedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureVariationRemovedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureVariationRemovedEvent;
  static deserializeBinaryFromReader(
    message: FeatureVariationRemovedEvent,
    reader: jspb.BinaryReader
  ): FeatureVariationRemovedEvent;
}

export namespace FeatureVariationRemovedEvent {
  export type AsObject = {
    id: string;
    variationId: string;
  };
}

export class VariationValueChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationValueChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: VariationValueChangedEvent
  ): VariationValueChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: VariationValueChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): VariationValueChangedEvent;
  static deserializeBinaryFromReader(
    message: VariationValueChangedEvent,
    reader: jspb.BinaryReader
  ): VariationValueChangedEvent;
}

export namespace VariationValueChangedEvent {
  export type AsObject = {
    featureId: string;
    id: string;
    value: string;
  };
}

export class VariationNameChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: VariationNameChangedEvent
  ): VariationNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: VariationNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): VariationNameChangedEvent;
  static deserializeBinaryFromReader(
    message: VariationNameChangedEvent,
    reader: jspb.BinaryReader
  ): VariationNameChangedEvent;
}

export namespace VariationNameChangedEvent {
  export type AsObject = {
    featureId: string;
    id: string;
    name: string;
  };
}

export class VariationDescriptionChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): VariationDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: VariationDescriptionChangedEvent
  ): VariationDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: VariationDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): VariationDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: VariationDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): VariationDescriptionChangedEvent;
}

export namespace VariationDescriptionChangedEvent {
  export type AsObject = {
    featureId: string;
    id: string;
    description: string;
  };
}

export class VariationUserAddedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getUser(): string;
  setUser(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationUserAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: VariationUserAddedEvent
  ): VariationUserAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: VariationUserAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): VariationUserAddedEvent;
  static deserializeBinaryFromReader(
    message: VariationUserAddedEvent,
    reader: jspb.BinaryReader
  ): VariationUserAddedEvent;
}

export namespace VariationUserAddedEvent {
  export type AsObject = {
    featureId: string;
    id: string;
    user: string;
  };
}

export class VariationUserRemovedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getUser(): string;
  setUser(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VariationUserRemovedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: VariationUserRemovedEvent
  ): VariationUserRemovedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: VariationUserRemovedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): VariationUserRemovedEvent;
  static deserializeBinaryFromReader(
    message: VariationUserRemovedEvent,
    reader: jspb.BinaryReader
  ): VariationUserRemovedEvent;
}

export namespace VariationUserRemovedEvent {
  export type AsObject = {
    featureId: string;
    id: string;
    user: string;
  };
}

export class FeatureRuleAddedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRule(): boolean;
  clearRule(): void;
  getRule(): proto_feature_rule_pb.Rule | undefined;
  setRule(value?: proto_feature_rule_pb.Rule): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureRuleAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureRuleAddedEvent
  ): FeatureRuleAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureRuleAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureRuleAddedEvent;
  static deserializeBinaryFromReader(
    message: FeatureRuleAddedEvent,
    reader: jspb.BinaryReader
  ): FeatureRuleAddedEvent;
}

export namespace FeatureRuleAddedEvent {
  export type AsObject = {
    id: string;
    rule?: proto_feature_rule_pb.Rule.AsObject;
  };
}

export class FeatureChangeRuleStrategyEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  hasStrategy(): boolean;
  clearStrategy(): void;
  getStrategy(): proto_feature_strategy_pb.Strategy | undefined;
  setStrategy(value?: proto_feature_strategy_pb.Strategy): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureChangeRuleStrategyEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureChangeRuleStrategyEvent
  ): FeatureChangeRuleStrategyEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureChangeRuleStrategyEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureChangeRuleStrategyEvent;
  static deserializeBinaryFromReader(
    message: FeatureChangeRuleStrategyEvent,
    reader: jspb.BinaryReader
  ): FeatureChangeRuleStrategyEvent;
}

export namespace FeatureChangeRuleStrategyEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    strategy?: proto_feature_strategy_pb.Strategy.AsObject;
  };
}

export class FeatureRulesOrderChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  clearRuleIdsList(): void;
  getRuleIdsList(): Array<string>;
  setRuleIdsList(value: Array<string>): void;
  addRuleIds(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureRulesOrderChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureRulesOrderChangedEvent
  ): FeatureRulesOrderChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureRulesOrderChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureRulesOrderChangedEvent;
  static deserializeBinaryFromReader(
    message: FeatureRulesOrderChangedEvent,
    reader: jspb.BinaryReader
  ): FeatureRulesOrderChangedEvent;
}

export namespace FeatureRulesOrderChangedEvent {
  export type AsObject = {
    featureId: string;
    ruleIdsList: Array<string>;
  };
}

export class FeatureRuleDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureRuleDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureRuleDeletedEvent
  ): FeatureRuleDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureRuleDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureRuleDeletedEvent;
  static deserializeBinaryFromReader(
    message: FeatureRuleDeletedEvent,
    reader: jspb.BinaryReader
  ): FeatureRuleDeletedEvent;
}

export namespace FeatureRuleDeletedEvent {
  export type AsObject = {
    id: string;
    ruleId: string;
  };
}

export class FeatureFixedStrategyChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  hasStrategy(): boolean;
  clearStrategy(): void;
  getStrategy(): proto_feature_strategy_pb.FixedStrategy | undefined;
  setStrategy(value?: proto_feature_strategy_pb.FixedStrategy): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): FeatureFixedStrategyChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureFixedStrategyChangedEvent
  ): FeatureFixedStrategyChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureFixedStrategyChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureFixedStrategyChangedEvent;
  static deserializeBinaryFromReader(
    message: FeatureFixedStrategyChangedEvent,
    reader: jspb.BinaryReader
  ): FeatureFixedStrategyChangedEvent;
}

export namespace FeatureFixedStrategyChangedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    strategy?: proto_feature_strategy_pb.FixedStrategy.AsObject;
  };
}

export class FeatureRolloutStrategyChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  hasStrategy(): boolean;
  clearStrategy(): void;
  getStrategy(): proto_feature_strategy_pb.RolloutStrategy | undefined;
  setStrategy(value?: proto_feature_strategy_pb.RolloutStrategy): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): FeatureRolloutStrategyChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureRolloutStrategyChangedEvent
  ): FeatureRolloutStrategyChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureRolloutStrategyChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): FeatureRolloutStrategyChangedEvent;
  static deserializeBinaryFromReader(
    message: FeatureRolloutStrategyChangedEvent,
    reader: jspb.BinaryReader
  ): FeatureRolloutStrategyChangedEvent;
}

export namespace FeatureRolloutStrategyChangedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    strategy?: proto_feature_strategy_pb.RolloutStrategy.AsObject;
  };
}

export class RuleClauseAddedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  hasClause(): boolean;
  clearClause(): void;
  getClause(): proto_feature_clause_pb.Clause | undefined;
  setClause(value?: proto_feature_clause_pb.Clause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuleClauseAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RuleClauseAddedEvent
  ): RuleClauseAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RuleClauseAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RuleClauseAddedEvent;
  static deserializeBinaryFromReader(
    message: RuleClauseAddedEvent,
    reader: jspb.BinaryReader
  ): RuleClauseAddedEvent;
}

export namespace RuleClauseAddedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    clause?: proto_feature_clause_pb.Clause.AsObject;
  };
}

export class RuleClauseDeletedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuleClauseDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RuleClauseDeletedEvent
  ): RuleClauseDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RuleClauseDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RuleClauseDeletedEvent;
  static deserializeBinaryFromReader(
    message: RuleClauseDeletedEvent,
    reader: jspb.BinaryReader
  ): RuleClauseDeletedEvent;
}

export namespace RuleClauseDeletedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    id: string;
  };
}

export class ClauseAttributeChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getAttribute(): string;
  setAttribute(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClauseAttributeChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ClauseAttributeChangedEvent
  ): ClauseAttributeChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ClauseAttributeChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ClauseAttributeChangedEvent;
  static deserializeBinaryFromReader(
    message: ClauseAttributeChangedEvent,
    reader: jspb.BinaryReader
  ): ClauseAttributeChangedEvent;
}

export namespace ClauseAttributeChangedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    id: string;
    attribute: string;
  };
}

export class ClauseOperatorChangedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getOperator(): proto_feature_clause_pb.Clause.OperatorMap[keyof proto_feature_clause_pb.Clause.OperatorMap];
  setOperator(
    value: proto_feature_clause_pb.Clause.OperatorMap[keyof proto_feature_clause_pb.Clause.OperatorMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClauseOperatorChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ClauseOperatorChangedEvent
  ): ClauseOperatorChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ClauseOperatorChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ClauseOperatorChangedEvent;
  static deserializeBinaryFromReader(
    message: ClauseOperatorChangedEvent,
    reader: jspb.BinaryReader
  ): ClauseOperatorChangedEvent;
}

export namespace ClauseOperatorChangedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    id: string;
    operator: proto_feature_clause_pb.Clause.OperatorMap[keyof proto_feature_clause_pb.Clause.OperatorMap];
  };
}

export class ClauseValueAddedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClauseValueAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ClauseValueAddedEvent
  ): ClauseValueAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ClauseValueAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ClauseValueAddedEvent;
  static deserializeBinaryFromReader(
    message: ClauseValueAddedEvent,
    reader: jspb.BinaryReader
  ): ClauseValueAddedEvent;
}

export namespace ClauseValueAddedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    id: string;
    value: string;
  };
}

export class ClauseValueRemovedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClauseValueRemovedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ClauseValueRemovedEvent
  ): ClauseValueRemovedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ClauseValueRemovedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ClauseValueRemovedEvent;
  static deserializeBinaryFromReader(
    message: ClauseValueRemovedEvent,
    reader: jspb.BinaryReader
  ): ClauseValueRemovedEvent;
}

export namespace ClauseValueRemovedEvent {
  export type AsObject = {
    featureId: string;
    ruleId: string;
    id: string;
    value: string;
  };
}

export class FeatureDefaultStrategyChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasStrategy(): boolean;
  clearStrategy(): void;
  getStrategy(): proto_feature_strategy_pb.Strategy | undefined;
  setStrategy(value?: proto_feature_strategy_pb.Strategy): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): FeatureDefaultStrategyChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureDefaultStrategyChangedEvent
  ): FeatureDefaultStrategyChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureDefaultStrategyChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): FeatureDefaultStrategyChangedEvent;
  static deserializeBinaryFromReader(
    message: FeatureDefaultStrategyChangedEvent,
    reader: jspb.BinaryReader
  ): FeatureDefaultStrategyChangedEvent;
}

export namespace FeatureDefaultStrategyChangedEvent {
  export type AsObject = {
    id: string;
    strategy?: proto_feature_strategy_pb.Strategy.AsObject;
  };
}

export class FeatureTagAddedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureTagAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureTagAddedEvent
  ): FeatureTagAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureTagAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureTagAddedEvent;
  static deserializeBinaryFromReader(
    message: FeatureTagAddedEvent,
    reader: jspb.BinaryReader
  ): FeatureTagAddedEvent;
}

export namespace FeatureTagAddedEvent {
  export type AsObject = {
    id: string;
    tag: string;
  };
}

export class FeatureTagRemovedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureTagRemovedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureTagRemovedEvent
  ): FeatureTagRemovedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureTagRemovedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureTagRemovedEvent;
  static deserializeBinaryFromReader(
    message: FeatureTagRemovedEvent,
    reader: jspb.BinaryReader
  ): FeatureTagRemovedEvent;
}

export namespace FeatureTagRemovedEvent {
  export type AsObject = {
    id: string;
    tag: string;
  };
}

export class FeatureVersionIncrementedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getVersion(): number;
  setVersion(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureVersionIncrementedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureVersionIncrementedEvent
  ): FeatureVersionIncrementedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureVersionIncrementedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureVersionIncrementedEvent;
  static deserializeBinaryFromReader(
    message: FeatureVersionIncrementedEvent,
    reader: jspb.BinaryReader
  ): FeatureVersionIncrementedEvent;
}

export namespace FeatureVersionIncrementedEvent {
  export type AsObject = {
    id: string;
    version: number;
  };
}

export class FeatureClonedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearVariationsList(): void;
  getVariationsList(): Array<proto_feature_variation_pb.Variation>;
  setVariationsList(value: Array<proto_feature_variation_pb.Variation>): void;
  addVariations(
    value?: proto_feature_variation_pb.Variation,
    index?: number
  ): proto_feature_variation_pb.Variation;

  clearTargetsList(): void;
  getTargetsList(): Array<proto_feature_target_pb.Target>;
  setTargetsList(value: Array<proto_feature_target_pb.Target>): void;
  addTargets(
    value?: proto_feature_target_pb.Target,
    index?: number
  ): proto_feature_target_pb.Target;

  clearRulesList(): void;
  getRulesList(): Array<proto_feature_rule_pb.Rule>;
  setRulesList(value: Array<proto_feature_rule_pb.Rule>): void;
  addRules(
    value?: proto_feature_rule_pb.Rule,
    index?: number
  ): proto_feature_rule_pb.Rule;

  hasDefaultStrategy(): boolean;
  clearDefaultStrategy(): void;
  getDefaultStrategy(): proto_feature_strategy_pb.Strategy | undefined;
  setDefaultStrategy(value?: proto_feature_strategy_pb.Strategy): void;

  getOffVariation(): string;
  setOffVariation(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getVariationType(): proto_feature_feature_pb.Feature.VariationTypeMap[keyof proto_feature_feature_pb.Feature.VariationTypeMap];
  setVariationType(
    value: proto_feature_feature_pb.Feature.VariationTypeMap[keyof proto_feature_feature_pb.Feature.VariationTypeMap]
  ): void;

  clearPrerequisitesList(): void;
  getPrerequisitesList(): Array<proto_feature_prerequisite_pb.Prerequisite>;
  setPrerequisitesList(
    value: Array<proto_feature_prerequisite_pb.Prerequisite>
  ): void;
  addPrerequisites(
    value?: proto_feature_prerequisite_pb.Prerequisite,
    index?: number
  ): proto_feature_prerequisite_pb.Prerequisite;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureClonedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureClonedEvent
  ): FeatureClonedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureClonedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureClonedEvent;
  static deserializeBinaryFromReader(
    message: FeatureClonedEvent,
    reader: jspb.BinaryReader
  ): FeatureClonedEvent;
}

export namespace FeatureClonedEvent {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    variationsList: Array<proto_feature_variation_pb.Variation.AsObject>;
    targetsList: Array<proto_feature_target_pb.Target.AsObject>;
    rulesList: Array<proto_feature_rule_pb.Rule.AsObject>;
    defaultStrategy?: proto_feature_strategy_pb.Strategy.AsObject;
    offVariation: string;
    tagsList: Array<string>;
    maintainer: string;
    variationType: proto_feature_feature_pb.Feature.VariationTypeMap[keyof proto_feature_feature_pb.Feature.VariationTypeMap];
    prerequisitesList: Array<proto_feature_prerequisite_pb.Prerequisite.AsObject>;
  };
}

export class FeatureSamplingSeedResetEvent extends jspb.Message {
  getSamplingSeed(): string;
  setSamplingSeed(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureSamplingSeedResetEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FeatureSamplingSeedResetEvent
  ): FeatureSamplingSeedResetEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FeatureSamplingSeedResetEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FeatureSamplingSeedResetEvent;
  static deserializeBinaryFromReader(
    message: FeatureSamplingSeedResetEvent,
    reader: jspb.BinaryReader
  ): FeatureSamplingSeedResetEvent;
}

export namespace FeatureSamplingSeedResetEvent {
  export type AsObject = {
    samplingSeed: string;
  };
}

export class GoalCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GoalCreatedEvent
  ): GoalCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GoalCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GoalCreatedEvent;
  static deserializeBinaryFromReader(
    message: GoalCreatedEvent,
    reader: jspb.BinaryReader
  ): GoalCreatedEvent;
}

export namespace GoalCreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    deleted: boolean;
    createdAt: number;
    updatedAt: number;
  };
}

export class GoalRenamedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GoalRenamedEvent
  ): GoalRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GoalRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GoalRenamedEvent;
  static deserializeBinaryFromReader(
    message: GoalRenamedEvent,
    reader: jspb.BinaryReader
  ): GoalRenamedEvent;
}

export namespace GoalRenamedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class GoalDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GoalDescriptionChangedEvent
  ): GoalDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GoalDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GoalDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: GoalDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): GoalDescriptionChangedEvent;
}

export namespace GoalDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class GoalArchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalArchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GoalArchivedEvent
  ): GoalArchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GoalArchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GoalArchivedEvent;
  static deserializeBinaryFromReader(
    message: GoalArchivedEvent,
    reader: jspb.BinaryReader
  ): GoalArchivedEvent;
}

export namespace GoalArchivedEvent {
  export type AsObject = {
    id: string;
  };
}

export class GoalDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GoalDeletedEvent
  ): GoalDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GoalDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GoalDeletedEvent;
  static deserializeBinaryFromReader(
    message: GoalDeletedEvent,
    reader: jspb.BinaryReader
  ): GoalDeletedEvent;
}

export namespace GoalDeletedEvent {
  export type AsObject = {
    id: string;
  };
}

export class ExperimentCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  clearVariationsList(): void;
  getVariationsList(): Array<proto_feature_variation_pb.Variation>;
  setVariationsList(value: Array<proto_feature_variation_pb.Variation>): void;
  addVariations(
    value?: proto_feature_variation_pb.Variation,
    index?: number
  ): proto_feature_variation_pb.Variation;

  getGoalId(): string;
  setGoalId(value: string): void;

  getStartAt(): number;
  setStartAt(value: number): void;

  getStopAt(): number;
  setStopAt(value: number): void;

  getStopped(): boolean;
  setStopped(value: boolean): void;

  getStoppedAt(): number;
  setStoppedAt(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  clearGoalIdsList(): void;
  getGoalIdsList(): Array<string>;
  setGoalIdsList(value: Array<string>): void;
  addGoalIds(value: string, index?: number): string;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getBaseVariationId(): string;
  setBaseVariationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentCreatedEvent
  ): ExperimentCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentCreatedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentCreatedEvent,
    reader: jspb.BinaryReader
  ): ExperimentCreatedEvent;
}

export namespace ExperimentCreatedEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    featureVersion: number;
    variationsList: Array<proto_feature_variation_pb.Variation.AsObject>;
    goalId: string;
    startAt: number;
    stopAt: number;
    stopped: boolean;
    stoppedAt: number;
    createdAt: number;
    updatedAt: number;
    goalIdsList: Array<string>;
    name: string;
    description: string;
    baseVariationId: string;
  };
}

export class ExperimentStoppedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getStoppedAt(): number;
  setStoppedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentStoppedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentStoppedEvent
  ): ExperimentStoppedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentStoppedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentStoppedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentStoppedEvent,
    reader: jspb.BinaryReader
  ): ExperimentStoppedEvent;
}

export namespace ExperimentStoppedEvent {
  export type AsObject = {
    id: string;
    stoppedAt: number;
  };
}

export class ExperimentArchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentArchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentArchivedEvent
  ): ExperimentArchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentArchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentArchivedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentArchivedEvent,
    reader: jspb.BinaryReader
  ): ExperimentArchivedEvent;
}

export namespace ExperimentArchivedEvent {
  export type AsObject = {
    id: string;
  };
}

export class ExperimentDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentDeletedEvent
  ): ExperimentDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentDeletedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentDeletedEvent,
    reader: jspb.BinaryReader
  ): ExperimentDeletedEvent;
}

export namespace ExperimentDeletedEvent {
  export type AsObject = {
    id: string;
  };
}

export class ExperimentStartAtChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getStartAt(): number;
  setStartAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentStartAtChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentStartAtChangedEvent
  ): ExperimentStartAtChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentStartAtChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentStartAtChangedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentStartAtChangedEvent,
    reader: jspb.BinaryReader
  ): ExperimentStartAtChangedEvent;
}

export namespace ExperimentStartAtChangedEvent {
  export type AsObject = {
    id: string;
    startAt: number;
  };
}

export class ExperimentStopAtChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getStopAt(): number;
  setStopAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentStopAtChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentStopAtChangedEvent
  ): ExperimentStopAtChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentStopAtChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentStopAtChangedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentStopAtChangedEvent,
    reader: jspb.BinaryReader
  ): ExperimentStopAtChangedEvent;
}

export namespace ExperimentStopAtChangedEvent {
  export type AsObject = {
    id: string;
    stopAt: number;
  };
}

export class ExperimentPeriodChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getStartAt(): number;
  setStartAt(value: number): void;

  getStopAt(): number;
  setStopAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentPeriodChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentPeriodChangedEvent
  ): ExperimentPeriodChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentPeriodChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentPeriodChangedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentPeriodChangedEvent,
    reader: jspb.BinaryReader
  ): ExperimentPeriodChangedEvent;
}

export namespace ExperimentPeriodChangedEvent {
  export type AsObject = {
    id: string;
    startAt: number;
    stopAt: number;
  };
}

export class ExperimentNameChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentNameChangedEvent
  ): ExperimentNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentNameChangedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentNameChangedEvent,
    reader: jspb.BinaryReader
  ): ExperimentNameChangedEvent;
}

export namespace ExperimentNameChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class ExperimentDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ExperimentDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentDescriptionChangedEvent
  ): ExperimentDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ExperimentDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): ExperimentDescriptionChangedEvent;
}

export namespace ExperimentDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class ExperimentStartedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentStartedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentStartedEvent
  ): ExperimentStartedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentStartedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentStartedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentStartedEvent,
    reader: jspb.BinaryReader
  ): ExperimentStartedEvent;
}

export namespace ExperimentStartedEvent {
  export type AsObject = {};
}

export class ExperimentFinishedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentFinishedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ExperimentFinishedEvent
  ): ExperimentFinishedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ExperimentFinishedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentFinishedEvent;
  static deserializeBinaryFromReader(
    message: ExperimentFinishedEvent,
    reader: jspb.BinaryReader
  ): ExperimentFinishedEvent;
}

export namespace ExperimentFinishedEvent {
  export type AsObject = {};
}

export class AccountV2CreatedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  getOrganizationRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  setOrganizationRole(
    value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]
  ): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<proto_account_account_pb.AccountV2.EnvironmentRole>;
  setEnvironmentRolesList(
    value: Array<proto_account_account_pb.AccountV2.EnvironmentRole>
  ): void;
  addEnvironmentRoles(
    value?: proto_account_account_pb.AccountV2.EnvironmentRole,
    index?: number
  ): proto_account_account_pb.AccountV2.EnvironmentRole;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getLanguage(): string;
  setLanguage(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2CreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2CreatedEvent
  ): AccountV2CreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2CreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2CreatedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2CreatedEvent,
    reader: jspb.BinaryReader
  ): AccountV2CreatedEvent;
}

export namespace AccountV2CreatedEvent {
  export type AsObject = {
    email: string;
    name: string;
    avatarImageUrl: string;
    organizationId: string;
    organizationRole: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
    environmentRolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>;
    disabled: boolean;
    createdAt: number;
    updatedAt: number;
    firstName: string;
    lastName: string;
    language: string;
    tagsList: Array<string>;
  };
}

export class AccountV2UpdatedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationId(): string;
  setOrganizationId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2UpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2UpdatedEvent
  ): AccountV2UpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2UpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2UpdatedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2UpdatedEvent,
    reader: jspb.BinaryReader
  ): AccountV2UpdatedEvent;
}

export namespace AccountV2UpdatedEvent {
  export type AsObject = {
    email: string;
    organizationId: string;
  };
}

export class AccountV2NameChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2NameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2NameChangedEvent
  ): AccountV2NameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2NameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2NameChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2NameChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2NameChangedEvent;
}

export namespace AccountV2NameChangedEvent {
  export type AsObject = {
    email: string;
    name: string;
  };
}

export class AccountV2AvatarImageURLChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getAvatarImageUrl(): string;
  setAvatarImageUrl(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AccountV2AvatarImageURLChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2AvatarImageURLChangedEvent
  ): AccountV2AvatarImageURLChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2AvatarImageURLChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AccountV2AvatarImageURLChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2AvatarImageURLChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2AvatarImageURLChangedEvent;
}

export namespace AccountV2AvatarImageURLChangedEvent {
  export type AsObject = {
    email: string;
    avatarImageUrl: string;
  };
}

export class AccountV2TagsChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2TagsChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2TagsChangedEvent
  ): AccountV2TagsChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2TagsChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2TagsChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2TagsChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2TagsChangedEvent;
}

export namespace AccountV2TagsChangedEvent {
  export type AsObject = {
    email: string;
    tagsList: Array<string>;
  };
}

export class AccountV2OrganizationRoleChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getOrganizationRole(): proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  setOrganizationRole(
    value: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AccountV2OrganizationRoleChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2OrganizationRoleChangedEvent
  ): AccountV2OrganizationRoleChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2OrganizationRoleChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AccountV2OrganizationRoleChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2OrganizationRoleChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2OrganizationRoleChangedEvent;
}

export namespace AccountV2OrganizationRoleChangedEvent {
  export type AsObject = {
    email: string;
    organizationRole: proto_account_account_pb.AccountV2.Role.OrganizationMap[keyof proto_account_account_pb.AccountV2.Role.OrganizationMap];
  };
}

export class AccountV2EnvironmentRolesChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  clearEnvironmentRolesList(): void;
  getEnvironmentRolesList(): Array<proto_account_account_pb.AccountV2.EnvironmentRole>;
  setEnvironmentRolesList(
    value: Array<proto_account_account_pb.AccountV2.EnvironmentRole>
  ): void;
  addEnvironmentRoles(
    value?: proto_account_account_pb.AccountV2.EnvironmentRole,
    index?: number
  ): proto_account_account_pb.AccountV2.EnvironmentRole;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AccountV2EnvironmentRolesChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2EnvironmentRolesChangedEvent
  ): AccountV2EnvironmentRolesChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2EnvironmentRolesChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AccountV2EnvironmentRolesChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2EnvironmentRolesChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2EnvironmentRolesChangedEvent;
}

export namespace AccountV2EnvironmentRolesChangedEvent {
  export type AsObject = {
    email: string;
    environmentRolesList: Array<proto_account_account_pb.AccountV2.EnvironmentRole.AsObject>;
  };
}

export class AccountV2EnabledEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2EnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2EnabledEvent
  ): AccountV2EnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2EnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2EnabledEvent;
  static deserializeBinaryFromReader(
    message: AccountV2EnabledEvent,
    reader: jspb.BinaryReader
  ): AccountV2EnabledEvent;
}

export namespace AccountV2EnabledEvent {
  export type AsObject = {
    email: string;
  };
}

export class AccountV2DisabledEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2DisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2DisabledEvent
  ): AccountV2DisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2DisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2DisabledEvent;
  static deserializeBinaryFromReader(
    message: AccountV2DisabledEvent,
    reader: jspb.BinaryReader
  ): AccountV2DisabledEvent;
}

export namespace AccountV2DisabledEvent {
  export type AsObject = {
    email: string;
  };
}

export class AccountV2DeletedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2DeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2DeletedEvent
  ): AccountV2DeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2DeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2DeletedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2DeletedEvent,
    reader: jspb.BinaryReader
  ): AccountV2DeletedEvent;
}

export namespace AccountV2DeletedEvent {
  export type AsObject = {
    email: string;
  };
}

export class SearchFilterCreatedEvent extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getTargettype(): proto_account_search_filter_pb.FilterTargetTypeMap[keyof proto_account_search_filter_pb.FilterTargetTypeMap];
  setTargettype(
    value: proto_account_search_filter_pb.FilterTargetTypeMap[keyof proto_account_search_filter_pb.FilterTargetTypeMap]
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getDefaultFilter(): boolean;
  setDefaultFilter(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchFilterCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SearchFilterCreatedEvent
  ): SearchFilterCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SearchFilterCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SearchFilterCreatedEvent;
  static deserializeBinaryFromReader(
    message: SearchFilterCreatedEvent,
    reader: jspb.BinaryReader
  ): SearchFilterCreatedEvent;
}

export namespace SearchFilterCreatedEvent {
  export type AsObject = {
    name: string;
    query: string;
    targettype: proto_account_search_filter_pb.FilterTargetTypeMap[keyof proto_account_search_filter_pb.FilterTargetTypeMap];
    environmentId: string;
    defaultFilter: boolean;
  };
}

export class SearchFilterNameChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchFilterNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SearchFilterNameChangedEvent
  ): SearchFilterNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SearchFilterNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SearchFilterNameChangedEvent;
  static deserializeBinaryFromReader(
    message: SearchFilterNameChangedEvent,
    reader: jspb.BinaryReader
  ): SearchFilterNameChangedEvent;
}

export namespace SearchFilterNameChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class SearchFilterQueryChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchFilterQueryChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SearchFilterQueryChangedEvent
  ): SearchFilterQueryChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SearchFilterQueryChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SearchFilterQueryChangedEvent;
  static deserializeBinaryFromReader(
    message: SearchFilterQueryChangedEvent,
    reader: jspb.BinaryReader
  ): SearchFilterQueryChangedEvent;
}

export namespace SearchFilterQueryChangedEvent {
  export type AsObject = {
    id: string;
    query: string;
  };
}

export class SearchFilterDefaultChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDefaultFilter(): boolean;
  setDefaultFilter(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchFilterDefaultChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SearchFilterDefaultChangedEvent
  ): SearchFilterDefaultChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SearchFilterDefaultChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SearchFilterDefaultChangedEvent;
  static deserializeBinaryFromReader(
    message: SearchFilterDefaultChangedEvent,
    reader: jspb.BinaryReader
  ): SearchFilterDefaultChangedEvent;
}

export namespace SearchFilterDefaultChangedEvent {
  export type AsObject = {
    id: string;
    defaultFilter: boolean;
  };
}

export class SearchFilterDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchFilterDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SearchFilterDeletedEvent
  ): SearchFilterDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SearchFilterDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SearchFilterDeletedEvent;
  static deserializeBinaryFromReader(
    message: SearchFilterDeletedEvent,
    reader: jspb.BinaryReader
  ): SearchFilterDeletedEvent;
}

export namespace SearchFilterDeletedEvent {
  export type AsObject = {
    id: string;
  };
}

export class APIKeyCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getRole(): proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
  setRole(
    value: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap]
  ): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getApiKey(): string;
  setApiKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKeyCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: APIKeyCreatedEvent
  ): APIKeyCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: APIKeyCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): APIKeyCreatedEvent;
  static deserializeBinaryFromReader(
    message: APIKeyCreatedEvent,
    reader: jspb.BinaryReader
  ): APIKeyCreatedEvent;
}

export namespace APIKeyCreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    role: proto_account_api_key_pb.APIKey.RoleMap[keyof proto_account_api_key_pb.APIKey.RoleMap];
    disabled: boolean;
    createdAt: number;
    updatedAt: number;
    maintainer: string;
    apiKey: string;
  };
}

export class APIKeyChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKeyChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: APIKeyChangedEvent
  ): APIKeyChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: APIKeyChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): APIKeyChangedEvent;
  static deserializeBinaryFromReader(
    message: APIKeyChangedEvent,
    reader: jspb.BinaryReader
  ): APIKeyChangedEvent;
}

export namespace APIKeyChangedEvent {
  export type AsObject = {
    id: string;
  };
}

export class APIKeyNameChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKeyNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: APIKeyNameChangedEvent
  ): APIKeyNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: APIKeyNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): APIKeyNameChangedEvent;
  static deserializeBinaryFromReader(
    message: APIKeyNameChangedEvent,
    reader: jspb.BinaryReader
  ): APIKeyNameChangedEvent;
}

export namespace APIKeyNameChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class APIKeyEnabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKeyEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: APIKeyEnabledEvent
  ): APIKeyEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: APIKeyEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): APIKeyEnabledEvent;
  static deserializeBinaryFromReader(
    message: APIKeyEnabledEvent,
    reader: jspb.BinaryReader
  ): APIKeyEnabledEvent;
}

export namespace APIKeyEnabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class APIKeyDisabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKeyDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: APIKeyDisabledEvent
  ): APIKeyDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: APIKeyDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): APIKeyDisabledEvent;
  static deserializeBinaryFromReader(
    message: APIKeyDisabledEvent,
    reader: jspb.BinaryReader
  ): APIKeyDisabledEvent;
}

export namespace APIKeyDisabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class SegmentCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentCreatedEvent
  ): SegmentCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentCreatedEvent;
  static deserializeBinaryFromReader(
    message: SegmentCreatedEvent,
    reader: jspb.BinaryReader
  ): SegmentCreatedEvent;
}

export namespace SegmentCreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
  };
}

export class SegmentDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentDeletedEvent
  ): SegmentDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentDeletedEvent;
  static deserializeBinaryFromReader(
    message: SegmentDeletedEvent,
    reader: jspb.BinaryReader
  ): SegmentDeletedEvent;
}

export namespace SegmentDeletedEvent {
  export type AsObject = {
    id: string;
  };
}

export class SegmentUpdatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasDescription(): boolean;
  clearDescription(): void;
  getDescription(): google_protobuf_wrappers_pb.StringValue | undefined;
  setDescription(value?: google_protobuf_wrappers_pb.StringValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentUpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentUpdatedEvent
  ): SegmentUpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentUpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentUpdatedEvent;
  static deserializeBinaryFromReader(
    message: SegmentUpdatedEvent,
    reader: jspb.BinaryReader
  ): SegmentUpdatedEvent;
}

export namespace SegmentUpdatedEvent {
  export type AsObject = {
    id: string;
    name?: google_protobuf_wrappers_pb.StringValue.AsObject;
    description?: google_protobuf_wrappers_pb.StringValue.AsObject;
  };
}

export class SegmentNameChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentNameChangedEvent
  ): SegmentNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentNameChangedEvent;
  static deserializeBinaryFromReader(
    message: SegmentNameChangedEvent,
    reader: jspb.BinaryReader
  ): SegmentNameChangedEvent;
}

export namespace SegmentNameChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class SegmentDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentDescriptionChangedEvent
  ): SegmentDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: SegmentDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): SegmentDescriptionChangedEvent;
}

export namespace SegmentDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class SegmentRuleAddedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasRule(): boolean;
  clearRule(): void;
  getRule(): proto_feature_rule_pb.Rule | undefined;
  setRule(value?: proto_feature_rule_pb.Rule): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentRuleAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentRuleAddedEvent
  ): SegmentRuleAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentRuleAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentRuleAddedEvent;
  static deserializeBinaryFromReader(
    message: SegmentRuleAddedEvent,
    reader: jspb.BinaryReader
  ): SegmentRuleAddedEvent;
}

export namespace SegmentRuleAddedEvent {
  export type AsObject = {
    id: string;
    rule?: proto_feature_rule_pb.Rule.AsObject;
  };
}

export class SegmentRuleDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentRuleDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentRuleDeletedEvent
  ): SegmentRuleDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentRuleDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentRuleDeletedEvent;
  static deserializeBinaryFromReader(
    message: SegmentRuleDeletedEvent,
    reader: jspb.BinaryReader
  ): SegmentRuleDeletedEvent;
}

export namespace SegmentRuleDeletedEvent {
  export type AsObject = {
    id: string;
    ruleId: string;
  };
}

export class SegmentRuleClauseAddedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  hasClause(): boolean;
  clearClause(): void;
  getClause(): proto_feature_clause_pb.Clause | undefined;
  setClause(value?: proto_feature_clause_pb.Clause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentRuleClauseAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentRuleClauseAddedEvent
  ): SegmentRuleClauseAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentRuleClauseAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentRuleClauseAddedEvent;
  static deserializeBinaryFromReader(
    message: SegmentRuleClauseAddedEvent,
    reader: jspb.BinaryReader
  ): SegmentRuleClauseAddedEvent;
}

export namespace SegmentRuleClauseAddedEvent {
  export type AsObject = {
    segmentId: string;
    ruleId: string;
    clause?: proto_feature_clause_pb.Clause.AsObject;
  };
}

export class SegmentRuleClauseDeletedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentRuleClauseDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentRuleClauseDeletedEvent
  ): SegmentRuleClauseDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentRuleClauseDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentRuleClauseDeletedEvent;
  static deserializeBinaryFromReader(
    message: SegmentRuleClauseDeletedEvent,
    reader: jspb.BinaryReader
  ): SegmentRuleClauseDeletedEvent;
}

export namespace SegmentRuleClauseDeletedEvent {
  export type AsObject = {
    segmentId: string;
    ruleId: string;
    clauseId: string;
  };
}

export class SegmentClauseAttributeChangedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getAttribute(): string;
  setAttribute(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): SegmentClauseAttributeChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentClauseAttributeChangedEvent
  ): SegmentClauseAttributeChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentClauseAttributeChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): SegmentClauseAttributeChangedEvent;
  static deserializeBinaryFromReader(
    message: SegmentClauseAttributeChangedEvent,
    reader: jspb.BinaryReader
  ): SegmentClauseAttributeChangedEvent;
}

export namespace SegmentClauseAttributeChangedEvent {
  export type AsObject = {
    segmentId: string;
    ruleId: string;
    clauseId: string;
    attribute: string;
  };
}

export class SegmentClauseOperatorChangedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getOperator(): proto_feature_clause_pb.Clause.OperatorMap[keyof proto_feature_clause_pb.Clause.OperatorMap];
  setOperator(
    value: proto_feature_clause_pb.Clause.OperatorMap[keyof proto_feature_clause_pb.Clause.OperatorMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): SegmentClauseOperatorChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentClauseOperatorChangedEvent
  ): SegmentClauseOperatorChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentClauseOperatorChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): SegmentClauseOperatorChangedEvent;
  static deserializeBinaryFromReader(
    message: SegmentClauseOperatorChangedEvent,
    reader: jspb.BinaryReader
  ): SegmentClauseOperatorChangedEvent;
}

export namespace SegmentClauseOperatorChangedEvent {
  export type AsObject = {
    segmentId: string;
    ruleId: string;
    clauseId: string;
    operator: proto_feature_clause_pb.Clause.OperatorMap[keyof proto_feature_clause_pb.Clause.OperatorMap];
  };
}

export class SegmentClauseValueAddedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentClauseValueAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentClauseValueAddedEvent
  ): SegmentClauseValueAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentClauseValueAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentClauseValueAddedEvent;
  static deserializeBinaryFromReader(
    message: SegmentClauseValueAddedEvent,
    reader: jspb.BinaryReader
  ): SegmentClauseValueAddedEvent;
}

export namespace SegmentClauseValueAddedEvent {
  export type AsObject = {
    segmentId: string;
    ruleId: string;
    clauseId: string;
    value: string;
  };
}

export class SegmentClauseValueRemovedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getRuleId(): string;
  setRuleId(value: string): void;

  getClauseId(): string;
  setClauseId(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentClauseValueRemovedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentClauseValueRemovedEvent
  ): SegmentClauseValueRemovedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentClauseValueRemovedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentClauseValueRemovedEvent;
  static deserializeBinaryFromReader(
    message: SegmentClauseValueRemovedEvent,
    reader: jspb.BinaryReader
  ): SegmentClauseValueRemovedEvent;
}

export namespace SegmentClauseValueRemovedEvent {
  export type AsObject = {
    segmentId: string;
    ruleId: string;
    clauseId: string;
    value: string;
  };
}

export class SegmentUserAddedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(
    value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentUserAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentUserAddedEvent
  ): SegmentUserAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentUserAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentUserAddedEvent;
  static deserializeBinaryFromReader(
    message: SegmentUserAddedEvent,
    reader: jspb.BinaryReader
  ): SegmentUserAddedEvent;
}

export namespace SegmentUserAddedEvent {
  export type AsObject = {
    segmentId: string;
    userIdsList: Array<string>;
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  };
}

export class SegmentUserDeletedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  clearUserIdsList(): void;
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): void;
  addUserIds(value: string, index?: number): string;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(
    value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentUserDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentUserDeletedEvent
  ): SegmentUserDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentUserDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentUserDeletedEvent;
  static deserializeBinaryFromReader(
    message: SegmentUserDeletedEvent,
    reader: jspb.BinaryReader
  ): SegmentUserDeletedEvent;
}

export namespace SegmentUserDeletedEvent {
  export type AsObject = {
    segmentId: string;
    userIdsList: Array<string>;
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  };
}

export class SegmentBulkUploadUsersEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getStatus(): proto_feature_segment_pb.Segment.StatusMap[keyof proto_feature_segment_pb.Segment.StatusMap];
  setStatus(
    value: proto_feature_segment_pb.Segment.StatusMap[keyof proto_feature_segment_pb.Segment.StatusMap]
  ): void;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(
    value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentBulkUploadUsersEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentBulkUploadUsersEvent
  ): SegmentBulkUploadUsersEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentBulkUploadUsersEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SegmentBulkUploadUsersEvent;
  static deserializeBinaryFromReader(
    message: SegmentBulkUploadUsersEvent,
    reader: jspb.BinaryReader
  ): SegmentBulkUploadUsersEvent;
}

export namespace SegmentBulkUploadUsersEvent {
  export type AsObject = {
    segmentId: string;
    status: proto_feature_segment_pb.Segment.StatusMap[keyof proto_feature_segment_pb.Segment.StatusMap];
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  };
}

export class SegmentBulkUploadUsersStatusChangedEvent extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  getStatus(): proto_feature_segment_pb.Segment.StatusMap[keyof proto_feature_segment_pb.Segment.StatusMap];
  setStatus(
    value: proto_feature_segment_pb.Segment.StatusMap[keyof proto_feature_segment_pb.Segment.StatusMap]
  ): void;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(
    value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]
  ): void;

  getCount(): number;
  setCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): SegmentBulkUploadUsersStatusChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SegmentBulkUploadUsersStatusChangedEvent
  ): SegmentBulkUploadUsersStatusChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SegmentBulkUploadUsersStatusChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): SegmentBulkUploadUsersStatusChangedEvent;
  static deserializeBinaryFromReader(
    message: SegmentBulkUploadUsersStatusChangedEvent,
    reader: jspb.BinaryReader
  ): SegmentBulkUploadUsersStatusChangedEvent;
}

export namespace SegmentBulkUploadUsersStatusChangedEvent {
  export type AsObject = {
    segmentId: string;
    status: proto_feature_segment_pb.Segment.StatusMap[keyof proto_feature_segment_pb.Segment.StatusMap];
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
    count: number;
  };
}

export class EnvironmentCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentCreatedEvent
  ): EnvironmentCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentCreatedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentCreatedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentCreatedEvent;
}

export namespace EnvironmentCreatedEvent {
  export type AsObject = {
    id: string;
    namespace: string;
    name: string;
    description: string;
    deleted: boolean;
    createdAt: number;
    updatedAt: number;
    projectId: string;
  };
}

export class EnvironmentRenamedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentRenamedEvent
  ): EnvironmentRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentRenamedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentRenamedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentRenamedEvent;
}

export namespace EnvironmentRenamedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class EnvironmentDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): EnvironmentDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentDescriptionChangedEvent
  ): EnvironmentDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): EnvironmentDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentDescriptionChangedEvent;
}

export namespace EnvironmentDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class EnvironmentDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentDeletedEvent
  ): EnvironmentDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentDeletedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentDeletedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentDeletedEvent;
}

export namespace EnvironmentDeletedEvent {
  export type AsObject = {
    id: string;
    namespace: string;
  };
}

export class EnvironmentV2CreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getRequireComment(): boolean;
  setRequireComment(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentV2CreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentV2CreatedEvent
  ): EnvironmentV2CreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentV2CreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentV2CreatedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentV2CreatedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentV2CreatedEvent;
}

export namespace EnvironmentV2CreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    urlCode: string;
    description: string;
    projectId: string;
    archived: boolean;
    createdAt: number;
    updatedAt: number;
    requireComment: boolean;
  };
}

export class EnvironmentV2RenamedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getOldName(): string;
  setOldName(value: string): void;

  getNewName(): string;
  setNewName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentV2RenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentV2RenamedEvent
  ): EnvironmentV2RenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentV2RenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentV2RenamedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentV2RenamedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentV2RenamedEvent;
}

export namespace EnvironmentV2RenamedEvent {
  export type AsObject = {
    id: string;
    projectId: string;
    oldName: string;
    newName: string;
  };
}

export class EnvironmentV2DescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getOldDescription(): string;
  setOldDescription(value: string): void;

  getNewDescription(): string;
  setNewDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): EnvironmentV2DescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentV2DescriptionChangedEvent
  ): EnvironmentV2DescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentV2DescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): EnvironmentV2DescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentV2DescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentV2DescriptionChangedEvent;
}

export namespace EnvironmentV2DescriptionChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
    projectId: string;
    oldDescription: string;
    newDescription: string;
  };
}

export class EnvironmentV2RequireCommentChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getRequireComment(): boolean;
  setRequireComment(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): EnvironmentV2RequireCommentChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentV2RequireCommentChangedEvent
  ): EnvironmentV2RequireCommentChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentV2RequireCommentChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): EnvironmentV2RequireCommentChangedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentV2RequireCommentChangedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentV2RequireCommentChangedEvent;
}

export namespace EnvironmentV2RequireCommentChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
    projectId: string;
    requireComment: boolean;
  };
}

export class EnvironmentV2ArchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentV2ArchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentV2ArchivedEvent
  ): EnvironmentV2ArchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentV2ArchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentV2ArchivedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentV2ArchivedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentV2ArchivedEvent;
}

export namespace EnvironmentV2ArchivedEvent {
  export type AsObject = {
    id: string;
    name: string;
    projectId: string;
  };
}

export class EnvironmentV2UnarchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentV2UnarchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentV2UnarchivedEvent
  ): EnvironmentV2UnarchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentV2UnarchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentV2UnarchivedEvent;
  static deserializeBinaryFromReader(
    message: EnvironmentV2UnarchivedEvent,
    reader: jspb.BinaryReader
  ): EnvironmentV2UnarchivedEvent;
}

export namespace EnvironmentV2UnarchivedEvent {
  export type AsObject = {
    id: string;
    name: string;
    projectId: string;
  };
}

export class AutoOpsRuleCreatedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getOpsType(): proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  setOpsType(
    value: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap]
  ): void;

  clearClausesList(): void;
  getClausesList(): Array<proto_autoops_clause_pb.Clause>;
  setClausesList(value: Array<proto_autoops_clause_pb.Clause>): void;
  addClauses(
    value?: proto_autoops_clause_pb.Clause,
    index?: number
  ): proto_autoops_clause_pb.Clause;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getOpsStatus(): proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
  setOpsStatus(
    value: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoOpsRuleCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleCreatedEvent
  ): AutoOpsRuleCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRuleCreatedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleCreatedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleCreatedEvent;
}

export namespace AutoOpsRuleCreatedEvent {
  export type AsObject = {
    featureId: string;
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
    clausesList: Array<proto_autoops_clause_pb.Clause.AsObject>;
    createdAt: number;
    updatedAt: number;
    opsStatus: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
  };
}

export class AutoOpsRuleStoppedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoOpsRuleStoppedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleStoppedEvent
  ): AutoOpsRuleStoppedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleStoppedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRuleStoppedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleStoppedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleStoppedEvent;
}

export namespace AutoOpsRuleStoppedEvent {
  export type AsObject = {};
}

export class AutoOpsRuleDeletedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoOpsRuleDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleDeletedEvent
  ): AutoOpsRuleDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRuleDeletedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleDeletedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleDeletedEvent;
}

export namespace AutoOpsRuleDeletedEvent {
  export type AsObject = {};
}

export class AutoOpsRuleOpsTypeChangedEvent extends jspb.Message {
  getOpsType(): proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  setOpsType(
    value: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoOpsRuleOpsTypeChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleOpsTypeChangedEvent
  ): AutoOpsRuleOpsTypeChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleOpsTypeChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRuleOpsTypeChangedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleOpsTypeChangedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleOpsTypeChangedEvent;
}

export namespace AutoOpsRuleOpsTypeChangedEvent {
  export type AsObject = {
    opsType: proto_autoops_auto_ops_rule_pb.OpsTypeMap[keyof proto_autoops_auto_ops_rule_pb.OpsTypeMap];
  };
}

export class AutoOpsRuleTriggeredAtChangedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AutoOpsRuleTriggeredAtChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleTriggeredAtChangedEvent
  ): AutoOpsRuleTriggeredAtChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleTriggeredAtChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AutoOpsRuleTriggeredAtChangedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleTriggeredAtChangedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleTriggeredAtChangedEvent;
}

export namespace AutoOpsRuleTriggeredAtChangedEvent {
  export type AsObject = {};
}

export class AutoOpsRuleOpsStatusChangedEvent extends jspb.Message {
  getOpsStatus(): proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
  setOpsStatus(
    value: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AutoOpsRuleOpsStatusChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleOpsStatusChangedEvent
  ): AutoOpsRuleOpsStatusChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleOpsStatusChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRuleOpsStatusChangedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleOpsStatusChangedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleOpsStatusChangedEvent;
}

export namespace AutoOpsRuleOpsStatusChangedEvent {
  export type AsObject = {
    opsStatus: proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap[keyof proto_autoops_auto_ops_rule_pb.AutoOpsStatusMap];
  };
}

export class OpsEventRateClauseAddedEvent extends jspb.Message {
  getClauseId(): string;
  setClauseId(value: string): void;

  hasOpsEventRateClause(): boolean;
  clearOpsEventRateClause(): void;
  getOpsEventRateClause():
    | proto_autoops_clause_pb.OpsEventRateClause
    | undefined;
  setOpsEventRateClause(
    value?: proto_autoops_clause_pb.OpsEventRateClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpsEventRateClauseAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OpsEventRateClauseAddedEvent
  ): OpsEventRateClauseAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OpsEventRateClauseAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OpsEventRateClauseAddedEvent;
  static deserializeBinaryFromReader(
    message: OpsEventRateClauseAddedEvent,
    reader: jspb.BinaryReader
  ): OpsEventRateClauseAddedEvent;
}

export namespace OpsEventRateClauseAddedEvent {
  export type AsObject = {
    clauseId: string;
    opsEventRateClause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject;
  };
}

export class OpsEventRateClauseChangedEvent extends jspb.Message {
  getClauseId(): string;
  setClauseId(value: string): void;

  hasOpsEventRateClause(): boolean;
  clearOpsEventRateClause(): void;
  getOpsEventRateClause():
    | proto_autoops_clause_pb.OpsEventRateClause
    | undefined;
  setOpsEventRateClause(
    value?: proto_autoops_clause_pb.OpsEventRateClause
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpsEventRateClauseChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OpsEventRateClauseChangedEvent
  ): OpsEventRateClauseChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OpsEventRateClauseChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OpsEventRateClauseChangedEvent;
  static deserializeBinaryFromReader(
    message: OpsEventRateClauseChangedEvent,
    reader: jspb.BinaryReader
  ): OpsEventRateClauseChangedEvent;
}

export namespace OpsEventRateClauseChangedEvent {
  export type AsObject = {
    clauseId: string;
    opsEventRateClause?: proto_autoops_clause_pb.OpsEventRateClause.AsObject;
  };
}

export class AutoOpsRuleClauseDeletedEvent extends jspb.Message {
  getClauseId(): string;
  setClauseId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AutoOpsRuleClauseDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AutoOpsRuleClauseDeletedEvent
  ): AutoOpsRuleClauseDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AutoOpsRuleClauseDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AutoOpsRuleClauseDeletedEvent;
  static deserializeBinaryFromReader(
    message: AutoOpsRuleClauseDeletedEvent,
    reader: jspb.BinaryReader
  ): AutoOpsRuleClauseDeletedEvent;
}

export namespace AutoOpsRuleClauseDeletedEvent {
  export type AsObject = {
    clauseId: string;
  };
}

export class DatetimeClauseAddedEvent extends jspb.Message {
  getClauseId(): string;
  setClauseId(value: string): void;

  hasDatetimeClause(): boolean;
  clearDatetimeClause(): void;
  getDatetimeClause(): proto_autoops_clause_pb.DatetimeClause | undefined;
  setDatetimeClause(value?: proto_autoops_clause_pb.DatetimeClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatetimeClauseAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DatetimeClauseAddedEvent
  ): DatetimeClauseAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DatetimeClauseAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DatetimeClauseAddedEvent;
  static deserializeBinaryFromReader(
    message: DatetimeClauseAddedEvent,
    reader: jspb.BinaryReader
  ): DatetimeClauseAddedEvent;
}

export namespace DatetimeClauseAddedEvent {
  export type AsObject = {
    clauseId: string;
    datetimeClause?: proto_autoops_clause_pb.DatetimeClause.AsObject;
  };
}

export class DatetimeClauseChangedEvent extends jspb.Message {
  getClauseId(): string;
  setClauseId(value: string): void;

  hasDatetimeClause(): boolean;
  clearDatetimeClause(): void;
  getDatetimeClause(): proto_autoops_clause_pb.DatetimeClause | undefined;
  setDatetimeClause(value?: proto_autoops_clause_pb.DatetimeClause): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatetimeClauseChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DatetimeClauseChangedEvent
  ): DatetimeClauseChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DatetimeClauseChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DatetimeClauseChangedEvent;
  static deserializeBinaryFromReader(
    message: DatetimeClauseChangedEvent,
    reader: jspb.BinaryReader
  ): DatetimeClauseChangedEvent;
}

export namespace DatetimeClauseChangedEvent {
  export type AsObject = {
    clauseId: string;
    datetimeClause?: proto_autoops_clause_pb.DatetimeClause.AsObject;
  };
}

export class PushCreatedEvent extends jspb.Message {
  getFcmServiceAccount(): string;
  setFcmServiceAccount(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushCreatedEvent
  ): PushCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushCreatedEvent;
  static deserializeBinaryFromReader(
    message: PushCreatedEvent,
    reader: jspb.BinaryReader
  ): PushCreatedEvent;
}

export namespace PushCreatedEvent {
  export type AsObject = {
    fcmServiceAccount: string;
    tagsList: Array<string>;
    name: string;
  };
}

export class PushDeletedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushDeletedEvent
  ): PushDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushDeletedEvent;
  static deserializeBinaryFromReader(
    message: PushDeletedEvent,
    reader: jspb.BinaryReader
  ): PushDeletedEvent;
}

export namespace PushDeletedEvent {
  export type AsObject = {};
}

export class PushTagsAddedEvent extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushTagsAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushTagsAddedEvent
  ): PushTagsAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushTagsAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushTagsAddedEvent;
  static deserializeBinaryFromReader(
    message: PushTagsAddedEvent,
    reader: jspb.BinaryReader
  ): PushTagsAddedEvent;
}

export namespace PushTagsAddedEvent {
  export type AsObject = {
    tagsList: Array<string>;
  };
}

export class PushTagsDeletedEvent extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushTagsDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushTagsDeletedEvent
  ): PushTagsDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushTagsDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushTagsDeletedEvent;
  static deserializeBinaryFromReader(
    message: PushTagsDeletedEvent,
    reader: jspb.BinaryReader
  ): PushTagsDeletedEvent;
}

export namespace PushTagsDeletedEvent {
  export type AsObject = {
    tagsList: Array<string>;
  };
}

export class PushRenamedEvent extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushRenamedEvent
  ): PushRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushRenamedEvent;
  static deserializeBinaryFromReader(
    message: PushRenamedEvent,
    reader: jspb.BinaryReader
  ): PushRenamedEvent;
}

export namespace PushRenamedEvent {
  export type AsObject = {
    name: string;
  };
}

export class PushTagsUpdatedEvent extends jspb.Message {
  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushTagsUpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushTagsUpdatedEvent
  ): PushTagsUpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushTagsUpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushTagsUpdatedEvent;
  static deserializeBinaryFromReader(
    message: PushTagsUpdatedEvent,
    reader: jspb.BinaryReader
  ): PushTagsUpdatedEvent;
}

export namespace PushTagsUpdatedEvent {
  export type AsObject = {
    tagsList: Array<string>;
  };
}

export class PushUpdatedEvent extends jspb.Message {
  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushUpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PushUpdatedEvent
  ): PushUpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PushUpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PushUpdatedEvent;
  static deserializeBinaryFromReader(
    message: PushUpdatedEvent,
    reader: jspb.BinaryReader
  ): PushUpdatedEvent;
}

export namespace PushUpdatedEvent {
  export type AsObject = {
    name?: google_protobuf_wrappers_pb.StringValue.AsObject;
    tagsList: Array<string>;
  };
}

export class SubscriptionCreatedEvent extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  hasRecipient(): boolean;
  clearRecipient(): void;
  getRecipient(): proto_notification_recipient_pb.Recipient | undefined;
  setRecipient(value?: proto_notification_recipient_pb.Recipient): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscriptionCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionCreatedEvent
  ): SubscriptionCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SubscriptionCreatedEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionCreatedEvent,
    reader: jspb.BinaryReader
  ): SubscriptionCreatedEvent;
}

export namespace SubscriptionCreatedEvent {
  export type AsObject = {
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
    recipient?: proto_notification_recipient_pb.Recipient.AsObject;
    name: string;
  };
}

export class SubscriptionDeletedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscriptionDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionDeletedEvent
  ): SubscriptionDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SubscriptionDeletedEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionDeletedEvent,
    reader: jspb.BinaryReader
  ): SubscriptionDeletedEvent;
}

export namespace SubscriptionDeletedEvent {
  export type AsObject = {};
}

export class SubscriptionEnabledEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscriptionEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionEnabledEvent
  ): SubscriptionEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SubscriptionEnabledEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionEnabledEvent,
    reader: jspb.BinaryReader
  ): SubscriptionEnabledEvent;
}

export namespace SubscriptionEnabledEvent {
  export type AsObject = {};
}

export class SubscriptionDisabledEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscriptionDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionDisabledEvent
  ): SubscriptionDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SubscriptionDisabledEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionDisabledEvent,
    reader: jspb.BinaryReader
  ): SubscriptionDisabledEvent;
}

export namespace SubscriptionDisabledEvent {
  export type AsObject = {};
}

export class SubscriptionSourceTypesAddedEvent extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): SubscriptionSourceTypesAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionSourceTypesAddedEvent
  ): SubscriptionSourceTypesAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionSourceTypesAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): SubscriptionSourceTypesAddedEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionSourceTypesAddedEvent,
    reader: jspb.BinaryReader
  ): SubscriptionSourceTypesAddedEvent;
}

export namespace SubscriptionSourceTypesAddedEvent {
  export type AsObject = {
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
  };
}

export class SubscriptionSourceTypesDeletedEvent extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): SubscriptionSourceTypesDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionSourceTypesDeletedEvent
  ): SubscriptionSourceTypesDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionSourceTypesDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): SubscriptionSourceTypesDeletedEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionSourceTypesDeletedEvent,
    reader: jspb.BinaryReader
  ): SubscriptionSourceTypesDeletedEvent;
}

export namespace SubscriptionSourceTypesDeletedEvent {
  export type AsObject = {
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
  };
}

export class SubscriptionRenamedEvent extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscriptionRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionRenamedEvent
  ): SubscriptionRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SubscriptionRenamedEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionRenamedEvent,
    reader: jspb.BinaryReader
  ): SubscriptionRenamedEvent;
}

export namespace SubscriptionRenamedEvent {
  export type AsObject = {
    name: string;
  };
}

export class SubscriptionUpdatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasName(): boolean;
  clearName(): void;
  getName(): google_protobuf_wrappers_pb.StringValue | undefined;
  setName(value?: google_protobuf_wrappers_pb.StringValue): void;

  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisabled(value?: google_protobuf_wrappers_pb.BoolValue): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscriptionUpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SubscriptionUpdatedEvent
  ): SubscriptionUpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SubscriptionUpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SubscriptionUpdatedEvent;
  static deserializeBinaryFromReader(
    message: SubscriptionUpdatedEvent,
    reader: jspb.BinaryReader
  ): SubscriptionUpdatedEvent;
}

export namespace SubscriptionUpdatedEvent {
  export type AsObject = {
    id: string;
    name?: google_protobuf_wrappers_pb.StringValue.AsObject;
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
    disabled?: google_protobuf_wrappers_pb.BoolValue.AsObject;
  };
}

export class AdminSubscriptionCreatedEvent extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  hasRecipient(): boolean;
  clearRecipient(): void;
  getRecipient(): proto_notification_recipient_pb.Recipient | undefined;
  setRecipient(value?: proto_notification_recipient_pb.Recipient): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AdminSubscriptionCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionCreatedEvent
  ): AdminSubscriptionCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AdminSubscriptionCreatedEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionCreatedEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionCreatedEvent;
}

export namespace AdminSubscriptionCreatedEvent {
  export type AsObject = {
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
    recipient?: proto_notification_recipient_pb.Recipient.AsObject;
    name: string;
  };
}

export class AdminSubscriptionDeletedEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AdminSubscriptionDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionDeletedEvent
  ): AdminSubscriptionDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AdminSubscriptionDeletedEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionDeletedEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionDeletedEvent;
}

export namespace AdminSubscriptionDeletedEvent {
  export type AsObject = {};
}

export class AdminSubscriptionEnabledEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AdminSubscriptionEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionEnabledEvent
  ): AdminSubscriptionEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AdminSubscriptionEnabledEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionEnabledEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionEnabledEvent;
}

export namespace AdminSubscriptionEnabledEvent {
  export type AsObject = {};
}

export class AdminSubscriptionDisabledEvent extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AdminSubscriptionDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionDisabledEvent
  ): AdminSubscriptionDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AdminSubscriptionDisabledEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionDisabledEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionDisabledEvent;
}

export namespace AdminSubscriptionDisabledEvent {
  export type AsObject = {};
}

export class AdminSubscriptionSourceTypesAddedEvent extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AdminSubscriptionSourceTypesAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionSourceTypesAddedEvent
  ): AdminSubscriptionSourceTypesAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionSourceTypesAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AdminSubscriptionSourceTypesAddedEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionSourceTypesAddedEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionSourceTypesAddedEvent;
}

export namespace AdminSubscriptionSourceTypesAddedEvent {
  export type AsObject = {
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
  };
}

export class AdminSubscriptionSourceTypesDeletedEvent extends jspb.Message {
  clearSourceTypesList(): void;
  getSourceTypesList(): Array<
    proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
  >;
  setSourceTypesList(
    value: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >
  ): void;
  addSourceTypes(
    value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    index?: number
  ): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): AdminSubscriptionSourceTypesDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionSourceTypesDeletedEvent
  ): AdminSubscriptionSourceTypesDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionSourceTypesDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): AdminSubscriptionSourceTypesDeletedEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionSourceTypesDeletedEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionSourceTypesDeletedEvent;
}

export namespace AdminSubscriptionSourceTypesDeletedEvent {
  export type AsObject = {
    sourceTypesList: Array<
      proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]
    >;
  };
}

export class AdminSubscriptionRenamedEvent extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AdminSubscriptionRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AdminSubscriptionRenamedEvent
  ): AdminSubscriptionRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AdminSubscriptionRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AdminSubscriptionRenamedEvent;
  static deserializeBinaryFromReader(
    message: AdminSubscriptionRenamedEvent,
    reader: jspb.BinaryReader
  ): AdminSubscriptionRenamedEvent;
}

export namespace AdminSubscriptionRenamedEvent {
  export type AsObject = {
    name: string;
  };
}

export class ProjectCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getTrial(): boolean;
  setTrial(value: boolean): void;

  getCreatorEmail(): string;
  setCreatorEmail(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectCreatedEvent
  ): ProjectCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectCreatedEvent;
  static deserializeBinaryFromReader(
    message: ProjectCreatedEvent,
    reader: jspb.BinaryReader
  ): ProjectCreatedEvent;
}

export namespace ProjectCreatedEvent {
  export type AsObject = {
    id: string;
    description: string;
    disabled: boolean;
    trial: boolean;
    creatorEmail: string;
    createdAt: number;
    updatedAt: number;
    name: string;
    urlCode: string;
  };
}

export class ProjectDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectDescriptionChangedEvent
  ): ProjectDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: ProjectDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): ProjectDescriptionChangedEvent;
}

export namespace ProjectDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class ProjectRenamedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRenamedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectRenamedEvent
  ): ProjectRenamedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectRenamedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRenamedEvent;
  static deserializeBinaryFromReader(
    message: ProjectRenamedEvent,
    reader: jspb.BinaryReader
  ): ProjectRenamedEvent;
}

export namespace ProjectRenamedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class ProjectEnabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectEnabledEvent
  ): ProjectEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectEnabledEvent;
  static deserializeBinaryFromReader(
    message: ProjectEnabledEvent,
    reader: jspb.BinaryReader
  ): ProjectEnabledEvent;
}

export namespace ProjectEnabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class ProjectDisabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectDisabledEvent
  ): ProjectDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectDisabledEvent;
  static deserializeBinaryFromReader(
    message: ProjectDisabledEvent,
    reader: jspb.BinaryReader
  ): ProjectDisabledEvent;
}

export namespace ProjectDisabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class ProjectTrialCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getTrial(): boolean;
  setTrial(value: boolean): void;

  getCreatorEmail(): string;
  setCreatorEmail(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectTrialCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectTrialCreatedEvent
  ): ProjectTrialCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectTrialCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectTrialCreatedEvent;
  static deserializeBinaryFromReader(
    message: ProjectTrialCreatedEvent,
    reader: jspb.BinaryReader
  ): ProjectTrialCreatedEvent;
}

export namespace ProjectTrialCreatedEvent {
  export type AsObject = {
    id: string;
    description: string;
    disabled: boolean;
    trial: boolean;
    creatorEmail: string;
    createdAt: number;
    updatedAt: number;
    name: string;
    urlCode: string;
  };
}

export class ProjectTrialConvertedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectTrialConvertedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProjectTrialConvertedEvent
  ): ProjectTrialConvertedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProjectTrialConvertedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProjectTrialConvertedEvent;
  static deserializeBinaryFromReader(
    message: ProjectTrialConvertedEvent,
    reader: jspb.BinaryReader
  ): ProjectTrialConvertedEvent;
}

export namespace ProjectTrialConvertedEvent {
  export type AsObject = {
    id: string;
  };
}

export class PrerequisiteAddedEvent extends jspb.Message {
  hasPrerequisite(): boolean;
  clearPrerequisite(): void;
  getPrerequisite(): proto_feature_prerequisite_pb.Prerequisite | undefined;
  setPrerequisite(value?: proto_feature_prerequisite_pb.Prerequisite): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PrerequisiteAddedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PrerequisiteAddedEvent
  ): PrerequisiteAddedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PrerequisiteAddedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PrerequisiteAddedEvent;
  static deserializeBinaryFromReader(
    message: PrerequisiteAddedEvent,
    reader: jspb.BinaryReader
  ): PrerequisiteAddedEvent;
}

export namespace PrerequisiteAddedEvent {
  export type AsObject = {
    prerequisite?: proto_feature_prerequisite_pb.Prerequisite.AsObject;
  };
}

export class PrerequisiteVariationChangedEvent extends jspb.Message {
  hasPrerequisite(): boolean;
  clearPrerequisite(): void;
  getPrerequisite(): proto_feature_prerequisite_pb.Prerequisite | undefined;
  setPrerequisite(value?: proto_feature_prerequisite_pb.Prerequisite): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): PrerequisiteVariationChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PrerequisiteVariationChangedEvent
  ): PrerequisiteVariationChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PrerequisiteVariationChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): PrerequisiteVariationChangedEvent;
  static deserializeBinaryFromReader(
    message: PrerequisiteVariationChangedEvent,
    reader: jspb.BinaryReader
  ): PrerequisiteVariationChangedEvent;
}

export namespace PrerequisiteVariationChangedEvent {
  export type AsObject = {
    prerequisite?: proto_feature_prerequisite_pb.Prerequisite.AsObject;
  };
}

export class PrerequisiteRemovedEvent extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PrerequisiteRemovedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: PrerequisiteRemovedEvent
  ): PrerequisiteRemovedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: PrerequisiteRemovedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): PrerequisiteRemovedEvent;
  static deserializeBinaryFromReader(
    message: PrerequisiteRemovedEvent,
    reader: jspb.BinaryReader
  ): PrerequisiteRemovedEvent;
}

export namespace PrerequisiteRemovedEvent {
  export type AsObject = {
    featureId: string;
  };
}

export class ProgressiveRolloutCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  hasClause(): boolean;
  clearClause(): void;
  getClause(): google_protobuf_any_pb.Any | undefined;
  setClause(value?: google_protobuf_any_pb.Any): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getType(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap];
  setType(
    value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRolloutCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutCreatedEvent
  ): ProgressiveRolloutCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutCreatedEvent;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutCreatedEvent,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutCreatedEvent;
}

export namespace ProgressiveRolloutCreatedEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    clause?: google_protobuf_any_pb.Any.AsObject;
    createdAt: number;
    updatedAt: number;
    type: proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.TypeMap];
  };
}

export class ProgressiveRolloutStoppedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getStatus(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap];
  setStatus(
    value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap]
  ): void;

  getStoppedBy(): proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap];
  setStoppedBy(
    value: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap]
  ): void;

  getStoppedAt(): number;
  setStoppedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRolloutStoppedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutStoppedEvent
  ): ProgressiveRolloutStoppedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutStoppedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutStoppedEvent;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutStoppedEvent,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutStoppedEvent;
}

export namespace ProgressiveRolloutStoppedEvent {
  export type AsObject = {
    id: string;
    status: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StatusMap];
    stoppedBy: proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap[keyof proto_autoops_progressive_rollout_pb.ProgressiveRollout.StoppedByMap];
    stoppedAt: number;
  };
}

export class ProgressiveRolloutDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRolloutDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutDeletedEvent
  ): ProgressiveRolloutDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRolloutDeletedEvent;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutDeletedEvent,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutDeletedEvent;
}

export namespace ProgressiveRolloutDeletedEvent {
  export type AsObject = {
    id: string;
  };
}

export class ProgressiveRolloutScheduleTriggeredAtChangedEvent extends jspb.Message {
  getScheduleId(): string;
  setScheduleId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): ProgressiveRolloutScheduleTriggeredAtChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ProgressiveRolloutScheduleTriggeredAtChangedEvent
  ): ProgressiveRolloutScheduleTriggeredAtChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ProgressiveRolloutScheduleTriggeredAtChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): ProgressiveRolloutScheduleTriggeredAtChangedEvent;
  static deserializeBinaryFromReader(
    message: ProgressiveRolloutScheduleTriggeredAtChangedEvent,
    reader: jspb.BinaryReader
  ): ProgressiveRolloutScheduleTriggeredAtChangedEvent;
}

export namespace ProgressiveRolloutScheduleTriggeredAtChangedEvent {
  export type AsObject = {
    scheduleId: string;
  };
}

export class OrganizationCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  getTrial(): boolean;
  setTrial(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getOwnerEmail(): string;
  setOwnerEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationCreatedEvent
  ): OrganizationCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationCreatedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationCreatedEvent,
    reader: jspb.BinaryReader
  ): OrganizationCreatedEvent;
}

export namespace OrganizationCreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    urlCode: string;
    description: string;
    disabled: boolean;
    archived: boolean;
    trial: boolean;
    createdAt: number;
    updatedAt: number;
    ownerEmail: string;
  };
}

export class OrganizationDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): OrganizationDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationDescriptionChangedEvent
  ): OrganizationDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): OrganizationDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): OrganizationDescriptionChangedEvent;
}

export namespace OrganizationDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    description: string;
  };
}

export class OrganizationOwnerEmailChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOwnerEmail(): string;
  setOwnerEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): OrganizationOwnerEmailChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationOwnerEmailChangedEvent
  ): OrganizationOwnerEmailChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationOwnerEmailChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): OrganizationOwnerEmailChangedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationOwnerEmailChangedEvent,
    reader: jspb.BinaryReader
  ): OrganizationOwnerEmailChangedEvent;
}

export namespace OrganizationOwnerEmailChangedEvent {
  export type AsObject = {
    id: string;
    ownerEmail: string;
  };
}

export class OrganizationNameChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationNameChangedEvent
  ): OrganizationNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationNameChangedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationNameChangedEvent,
    reader: jspb.BinaryReader
  ): OrganizationNameChangedEvent;
}

export namespace OrganizationNameChangedEvent {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class OrganizationEnabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationEnabledEvent
  ): OrganizationEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationEnabledEvent;
  static deserializeBinaryFromReader(
    message: OrganizationEnabledEvent,
    reader: jspb.BinaryReader
  ): OrganizationEnabledEvent;
}

export namespace OrganizationEnabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class OrganizationDisabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationDisabledEvent
  ): OrganizationDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationDisabledEvent;
  static deserializeBinaryFromReader(
    message: OrganizationDisabledEvent,
    reader: jspb.BinaryReader
  ): OrganizationDisabledEvent;
}

export namespace OrganizationDisabledEvent {
  export type AsObject = {
    id: string;
  };
}

export class OrganizationArchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationArchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationArchivedEvent
  ): OrganizationArchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationArchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationArchivedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationArchivedEvent,
    reader: jspb.BinaryReader
  ): OrganizationArchivedEvent;
}

export namespace OrganizationArchivedEvent {
  export type AsObject = {
    id: string;
  };
}

export class OrganizationUnarchivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationUnarchivedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationUnarchivedEvent
  ): OrganizationUnarchivedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationUnarchivedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationUnarchivedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationUnarchivedEvent,
    reader: jspb.BinaryReader
  ): OrganizationUnarchivedEvent;
}

export namespace OrganizationUnarchivedEvent {
  export type AsObject = {
    id: string;
  };
}

export class OrganizationTrialConvertedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationTrialConvertedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: OrganizationTrialConvertedEvent
  ): OrganizationTrialConvertedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: OrganizationTrialConvertedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationTrialConvertedEvent;
  static deserializeBinaryFromReader(
    message: OrganizationTrialConvertedEvent,
    reader: jspb.BinaryReader
  ): OrganizationTrialConvertedEvent;
}

export namespace OrganizationTrialConvertedEvent {
  export type AsObject = {
    id: string;
  };
}

export class FlagTriggerCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getType(): proto_feature_flag_trigger_pb.FlagTrigger.TypeMap[keyof proto_feature_flag_trigger_pb.FlagTrigger.TypeMap];
  setType(
    value: proto_feature_flag_trigger_pb.FlagTrigger.TypeMap[keyof proto_feature_flag_trigger_pb.FlagTrigger.TypeMap]
  ): void;

  getAction(): proto_feature_flag_trigger_pb.FlagTrigger.ActionMap[keyof proto_feature_flag_trigger_pb.FlagTrigger.ActionMap];
  setAction(
    value: proto_feature_flag_trigger_pb.FlagTrigger.ActionMap[keyof proto_feature_flag_trigger_pb.FlagTrigger.ActionMap]
  ): void;

  getDescription(): string;
  setDescription(value: string): void;

  getToken(): string;
  setToken(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerCreatedEvent
  ): FlagTriggerCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerCreatedEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerCreatedEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerCreatedEvent;
}

export namespace FlagTriggerCreatedEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    type: proto_feature_flag_trigger_pb.FlagTrigger.TypeMap[keyof proto_feature_flag_trigger_pb.FlagTrigger.TypeMap];
    action: proto_feature_flag_trigger_pb.FlagTrigger.ActionMap[keyof proto_feature_flag_trigger_pb.FlagTrigger.ActionMap];
    description: string;
    token: string;
    createdAt: number;
    updatedAt: number;
    environmentId: string;
  };
}

export class FlagTriggerResetEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getToken(): string;
  setToken(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerResetEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerResetEvent
  ): FlagTriggerResetEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerResetEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerResetEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerResetEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerResetEvent;
}

export namespace FlagTriggerResetEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    token: string;
    environmentId: string;
  };
}

export class FlagTriggerDescriptionChangedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(
    includeInstance?: boolean
  ): FlagTriggerDescriptionChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerDescriptionChangedEvent
  ): FlagTriggerDescriptionChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerDescriptionChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(
    bytes: Uint8Array
  ): FlagTriggerDescriptionChangedEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerDescriptionChangedEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerDescriptionChangedEvent;
}

export namespace FlagTriggerDescriptionChangedEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    description: string;
    environmentId: string;
  };
}

export class FlagTriggerDisabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerDisabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerDisabledEvent
  ): FlagTriggerDisabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerDisabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerDisabledEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerDisabledEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerDisabledEvent;
}

export namespace FlagTriggerDisabledEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    environmentId: string;
  };
}

export class FlagTriggerEnabledEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerEnabledEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerEnabledEvent
  ): FlagTriggerEnabledEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerEnabledEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerEnabledEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerEnabledEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerEnabledEvent;
}

export namespace FlagTriggerEnabledEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    environmentId: string;
  };
}

export class FlagTriggerDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerDeletedEvent
  ): FlagTriggerDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerDeletedEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerDeletedEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerDeletedEvent;
}

export namespace FlagTriggerDeletedEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    environmentId: string;
  };
}

export class FlagTriggerUsageUpdatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getLastTriggeredAt(): number;
  setLastTriggeredAt(value: number): void;

  getTriggerTimes(): number;
  setTriggerTimes(value: number): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlagTriggerUsageUpdatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FlagTriggerUsageUpdatedEvent
  ): FlagTriggerUsageUpdatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FlagTriggerUsageUpdatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FlagTriggerUsageUpdatedEvent;
  static deserializeBinaryFromReader(
    message: FlagTriggerUsageUpdatedEvent,
    reader: jspb.BinaryReader
  ): FlagTriggerUsageUpdatedEvent;
}

export namespace FlagTriggerUsageUpdatedEvent {
  export type AsObject = {
    id: string;
    featureId: string;
    lastTriggeredAt: number;
    triggerTimes: number;
    environmentId: string;
  };
}

export class AccountV2FirstNameChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2FirstNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2FirstNameChangedEvent
  ): AccountV2FirstNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2FirstNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2FirstNameChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2FirstNameChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2FirstNameChangedEvent;
}

export namespace AccountV2FirstNameChangedEvent {
  export type AsObject = {
    email: string;
    firstName: string;
  };
}

export class AccountV2LastNameChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2LastNameChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2LastNameChangedEvent
  ): AccountV2LastNameChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2LastNameChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2LastNameChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2LastNameChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2LastNameChangedEvent;
}

export namespace AccountV2LastNameChangedEvent {
  export type AsObject = {
    email: string;
    lastName: string;
  };
}

export class AccountV2LanguageChangedEvent extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getLanguage(): string;
  setLanguage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccountV2LanguageChangedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: AccountV2LanguageChangedEvent
  ): AccountV2LanguageChangedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: AccountV2LanguageChangedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): AccountV2LanguageChangedEvent;
  static deserializeBinaryFromReader(
    message: AccountV2LanguageChangedEvent,
    reader: jspb.BinaryReader
  ): AccountV2LanguageChangedEvent;
}

export namespace AccountV2LanguageChangedEvent {
  export type AsObject = {
    email: string;
    language: string;
  };
}

export class TagCreatedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getEntityType(): proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap];
  setEntityType(
    value: proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap]
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TagCreatedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: TagCreatedEvent
  ): TagCreatedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: TagCreatedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): TagCreatedEvent;
  static deserializeBinaryFromReader(
    message: TagCreatedEvent,
    reader: jspb.BinaryReader
  ): TagCreatedEvent;
}

export namespace TagCreatedEvent {
  export type AsObject = {
    id: string;
    name: string;
    createdAt: number;
    updatedAt: number;
    entityType: proto_tag_tag_pb.Tag.EntityTypeMap[keyof proto_tag_tag_pb.Tag.EntityTypeMap];
    environmentId: string;
  };
}

export class TagDeletedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TagDeletedEvent.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: TagDeletedEvent
  ): TagDeletedEvent.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: TagDeletedEvent,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): TagDeletedEvent;
  static deserializeBinaryFromReader(
    message: TagDeletedEvent,
    reader: jspb.BinaryReader
  ): TagDeletedEvent;
}

export namespace TagDeletedEvent {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}
