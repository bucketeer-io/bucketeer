// package: bucketeer.event.service
// file: proto/event/service/user.proto

import * as jspb from "google-protobuf";
import * as proto_event_client_event_pb from "../../../proto/event/client/event_pb";

export class UserEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getLastSeen(): number;
  setLastSeen(value: number): void;

  getDataMap(): jspb.Map<string, string>;
  clearDataMap(): void;
  getSourceId(): proto_event_client_event_pb.SourceIdMap[keyof proto_event_client_event_pb.SourceIdMap];
  setSourceId(value: proto_event_client_event_pb.SourceIdMap[keyof proto_event_client_event_pb.SourceIdMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserEvent.AsObject;
  static toObject(includeInstance: boolean, msg: UserEvent): UserEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserEvent;
  static deserializeBinaryFromReader(message: UserEvent, reader: jspb.BinaryReader): UserEvent;
}

export namespace UserEvent {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    tag: string,
    userId: string,
    lastSeen: number,
    dataMap: Array<[string, string]>,
    sourceId: proto_event_client_event_pb.SourceIdMap[keyof proto_event_client_event_pb.SourceIdMap],
  }
}

