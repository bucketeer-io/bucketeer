// package: bucketeer.notification
// file: proto/notification/subscription.proto

import * as jspb from "google-protobuf";
import * as proto_notification_recipient_pb from "../../proto/notification/recipient_pb";

export class Subscription extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getDisabled(): boolean;
  setDisabled(value: boolean): void;

  clearSourceTypesList(): void;
  getSourceTypesList(): Array<Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]>;
  setSourceTypesList(value: Array<Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]>): void;
  addSourceTypes(value: Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap], index?: number): Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap];

  hasRecipient(): boolean;
  clearRecipient(): void;
  getRecipient(): proto_notification_recipient_pb.Recipient | undefined;
  setRecipient(value?: proto_notification_recipient_pb.Recipient): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Subscription.AsObject;
  static toObject(includeInstance: boolean, msg: Subscription): Subscription.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Subscription, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Subscription;
  static deserializeBinaryFromReader(message: Subscription, reader: jspb.BinaryReader): Subscription;
}

export namespace Subscription {
  export type AsObject = {
    id: string,
    createdAt: number,
    updatedAt: number,
    disabled: boolean,
    sourceTypesList: Array<Subscription.SourceTypeMap[keyof Subscription.SourceTypeMap]>,
    recipient?: proto_notification_recipient_pb.Recipient.AsObject,
    name: string,
  }

  export interface SourceTypeMap {
    DOMAIN_EVENT_FEATURE: 0;
    DOMAIN_EVENT_GOAL: 1;
    DOMAIN_EVENT_EXPERIMENT: 2;
    DOMAIN_EVENT_ACCOUNT: 3;
    DOMAIN_EVENT_APIKEY: 4;
    DOMAIN_EVENT_SEGMENT: 5;
    DOMAIN_EVENT_ENVIRONMENT: 6;
    DOMAIN_EVENT_ADMIN_ACCOUNT: 7;
    DOMAIN_EVENT_AUTOOPS_RULE: 8;
    DOMAIN_EVENT_PUSH: 9;
    DOMAIN_EVENT_SUBSCRIPTION: 10;
    DOMAIN_EVENT_ADMIN_SUBSCRIPTION: 11;
    DOMAIN_EVENT_PROJECT: 12;
    DOMAIN_EVENT_WEBHOOK: 13;
    DOMAIN_EVENT_PROGRESSIVE_ROLLOUT: 14;
    FEATURE_STALE: 100;
    EXPERIMENT_RUNNING: 200;
    MAU_COUNT: 300;
  }

  export const SourceType: SourceTypeMap;
}

