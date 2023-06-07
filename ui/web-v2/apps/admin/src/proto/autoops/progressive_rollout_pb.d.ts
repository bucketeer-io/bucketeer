// package: bucketeer.autoops
// file: proto/autoops/progressive_rollout.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";

export class ProgressiveRollout extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  hasClause(): boolean;
  clearClause(): void;
  getClause(): google_protobuf_any_pb.Any | undefined;
  setClause(value?: google_protobuf_any_pb.Any): void;

  getStatus(): ProgressiveRollout.StatusMap[keyof ProgressiveRollout.StatusMap];
  setStatus(value: ProgressiveRollout.StatusMap[keyof ProgressiveRollout.StatusMap]): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getType(): ProgressiveRollout.TypeMap[keyof ProgressiveRollout.TypeMap];
  setType(value: ProgressiveRollout.TypeMap[keyof ProgressiveRollout.TypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgressiveRollout.AsObject;
  static toObject(includeInstance: boolean, msg: ProgressiveRollout): ProgressiveRollout.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProgressiveRollout, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProgressiveRollout;
  static deserializeBinaryFromReader(message: ProgressiveRollout, reader: jspb.BinaryReader): ProgressiveRollout;
}

export namespace ProgressiveRollout {
  export type AsObject = {
    id: string,
    featureId: string,
    clause?: google_protobuf_any_pb.Any.AsObject,
    status: ProgressiveRollout.StatusMap[keyof ProgressiveRollout.StatusMap],
    createdAt: number,
    updatedAt: number,
    type: ProgressiveRollout.TypeMap[keyof ProgressiveRollout.TypeMap],
  }

  export interface TypeMap {
    MANUAL_SCHEDULE: 0;
    TEMPLATE_SCHEDULE: 1;
  }

  export const Type: TypeMap;

  export interface StatusMap {
    WAITING: 0;
    RUNNING: 1;
    FINISHED: 2;
  }

  export const Status: StatusMap;
}

