// package: bucketeer.auditlog
// file: proto/auditlog/auditlog.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as proto_event_domain_event_pb from "../../proto/event/domain/event_pb";
import * as proto_event_domain_localized_message_pb from "../../proto/event/domain/localized_message_pb";

export class AuditLog extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  getEntityType(): proto_event_domain_event_pb.Event.EntityTypeMap[keyof proto_event_domain_event_pb.Event.EntityTypeMap];
  setEntityType(value: proto_event_domain_event_pb.Event.EntityTypeMap[keyof proto_event_domain_event_pb.Event.EntityTypeMap]): void;

  getEntityId(): string;
  setEntityId(value: string): void;

  getType(): proto_event_domain_event_pb.Event.TypeMap[keyof proto_event_domain_event_pb.Event.TypeMap];
  setType(value: proto_event_domain_event_pb.Event.TypeMap[keyof proto_event_domain_event_pb.Event.TypeMap]): void;

  hasEvent(): boolean;
  clearEvent(): void;
  getEvent(): google_protobuf_any_pb.Any | undefined;
  setEvent(value?: google_protobuf_any_pb.Any): void;

  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): proto_event_domain_event_pb.Editor | undefined;
  setEditor(value?: proto_event_domain_event_pb.Editor): void;

  hasOptions(): boolean;
  clearOptions(): void;
  getOptions(): proto_event_domain_event_pb.Options | undefined;
  setOptions(value?: proto_event_domain_event_pb.Options): void;

  hasLocalizedMessage(): boolean;
  clearLocalizedMessage(): void;
  getLocalizedMessage(): proto_event_domain_localized_message_pb.LocalizedMessage | undefined;
  setLocalizedMessage(value?: proto_event_domain_localized_message_pb.LocalizedMessage): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuditLog.AsObject;
  static toObject(includeInstance: boolean, msg: AuditLog): AuditLog.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuditLog, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuditLog;
  static deserializeBinaryFromReader(message: AuditLog, reader: jspb.BinaryReader): AuditLog;
}

export namespace AuditLog {
  export type AsObject = {
    id: string,
    timestamp: number,
    entityType: proto_event_domain_event_pb.Event.EntityTypeMap[keyof proto_event_domain_event_pb.Event.EntityTypeMap],
    entityId: string,
    type: proto_event_domain_event_pb.Event.TypeMap[keyof proto_event_domain_event_pb.Event.TypeMap],
    event?: google_protobuf_any_pb.Any.AsObject,
    editor?: proto_event_domain_event_pb.Editor.AsObject,
    options?: proto_event_domain_event_pb.Options.AsObject,
    localizedMessage?: proto_event_domain_localized_message_pb.LocalizedMessage.AsObject,
  }
}

