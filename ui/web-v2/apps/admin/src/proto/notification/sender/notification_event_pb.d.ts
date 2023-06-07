// package: bucketeer.notification.sender
// file: proto/notification/sender/notification_event.proto

import * as jspb from "google-protobuf";
import * as proto_notification_sender_notification_pb from "../../../proto/notification/sender/notification_pb";
import * as proto_notification_subscription_pb from "../../../proto/notification/subscription_pb";

export class NotificationEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getSourceType(): proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap];
  setSourceType(value: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap]): void;

  hasNotification(): boolean;
  clearNotification(): void;
  getNotification(): proto_notification_sender_notification_pb.Notification | undefined;
  setNotification(value?: proto_notification_sender_notification_pb.Notification): void;

  getIsAdminEvent(): boolean;
  setIsAdminEvent(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEvent.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEvent): NotificationEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotificationEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEvent;
  static deserializeBinaryFromReader(message: NotificationEvent, reader: jspb.BinaryReader): NotificationEvent;
}

export namespace NotificationEvent {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    sourceType: proto_notification_subscription_pb.Subscription.SourceTypeMap[keyof proto_notification_subscription_pb.Subscription.SourceTypeMap],
    notification?: proto_notification_sender_notification_pb.Notification.AsObject,
    isAdminEvent: boolean,
  }
}

