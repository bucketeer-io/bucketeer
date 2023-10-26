// package: bucketeer.eventcounter
// file: proto/eventcounter/mau_summary.proto

import * as jspb from "google-protobuf";
import * as proto_event_client_event_pb from "../../proto/event/client/event_pb";

export class MAUSummary extends jspb.Message {
  getYearmonth(): string;
  setYearmonth(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getSourceId(): proto_event_client_event_pb.SourceIdMap[keyof proto_event_client_event_pb.SourceIdMap];
  setSourceId(value: proto_event_client_event_pb.SourceIdMap[keyof proto_event_client_event_pb.SourceIdMap]): void;

  getUserCount(): number;
  setUserCount(value: number): void;

  getRequestCount(): number;
  setRequestCount(value: number): void;

  getEvaluationCount(): number;
  setEvaluationCount(value: number): void;

  getGoalCount(): number;
  setGoalCount(value: number): void;

  getIsAll(): boolean;
  setIsAll(value: boolean): void;

  getIsFinished(): boolean;
  setIsFinished(value: boolean): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MAUSummary.AsObject;
  static toObject(includeInstance: boolean, msg: MAUSummary): MAUSummary.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MAUSummary, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MAUSummary;
  static deserializeBinaryFromReader(message: MAUSummary, reader: jspb.BinaryReader): MAUSummary;
}

export namespace MAUSummary {
  export type AsObject = {
    yearmonth: string,
    environmentId: string,
    sourceId: proto_event_client_event_pb.SourceIdMap[keyof proto_event_client_event_pb.SourceIdMap],
    userCount: number,
    requestCount: number,
    evaluationCount: number,
    goalCount: number,
    isAll: boolean,
    isFinished: boolean,
    createdAt: number,
    updatedAt: number,
  }
}

