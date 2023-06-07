// package: bucketeer.event.service
// file: proto/event/service/feature.proto

import * as jspb from "google-protobuf";
import * as proto_user_user_pb from "../../../proto/user/user_pb";

export class EvaluationRequestEvent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  hasUser(): boolean;
  clearUser(): void;
  getUser(): proto_user_user_pb.User | undefined;
  setUser(value?: proto_user_user_pb.User): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluationRequestEvent.AsObject;
  static toObject(includeInstance: boolean, msg: EvaluationRequestEvent): EvaluationRequestEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EvaluationRequestEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EvaluationRequestEvent;
  static deserializeBinaryFromReader(message: EvaluationRequestEvent, reader: jspb.BinaryReader): EvaluationRequestEvent;
}

export namespace EvaluationRequestEvent {
  export type AsObject = {
    id: string,
    timestamp: number,
    user?: proto_user_user_pb.User.AsObject,
    environmentNamespace: string,
    tag: string,
  }
}

