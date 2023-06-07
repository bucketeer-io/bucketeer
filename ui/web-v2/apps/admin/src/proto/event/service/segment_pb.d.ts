// package: bucketeer.event.service
// file: proto/event/service/segment.proto

import * as jspb from "google-protobuf";
import * as proto_feature_segment_pb from "../../../proto/feature/segment_pb";
import * as proto_event_domain_event_pb from "../../../proto/event/domain/event_pb";

export class BulkSegmentUsersReceivedEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getSegmentId(): string;
  setSegmentId(value: string): void;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  getState(): proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap];
  setState(value: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap]): void;

  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): proto_event_domain_event_pb.Editor | undefined;
  setEditor(value?: proto_event_domain_event_pb.Editor): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSegmentUsersReceivedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSegmentUsersReceivedEvent): BulkSegmentUsersReceivedEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BulkSegmentUsersReceivedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSegmentUsersReceivedEvent;
  static deserializeBinaryFromReader(message: BulkSegmentUsersReceivedEvent, reader: jspb.BinaryReader): BulkSegmentUsersReceivedEvent;
}

export namespace BulkSegmentUsersReceivedEvent {
  export type AsObject = {
    id: string,
    environmentNamespace: string,
    segmentId: string,
    data: Uint8Array | string,
    state: proto_feature_segment_pb.SegmentUser.StateMap[keyof proto_feature_segment_pb.SegmentUser.StateMap],
    editor?: proto_event_domain_event_pb.Editor.AsObject,
  }
}

