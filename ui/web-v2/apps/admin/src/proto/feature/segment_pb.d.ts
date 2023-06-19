// package: bucketeer.feature
// file: proto/feature/segment.proto

import * as jspb from "google-protobuf";
import * as proto_feature_rule_pb from "../../proto/feature/rule_pb";
import * as proto_feature_feature_pb from "../../proto/feature/feature_pb";

export class Segment extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearRulesList(): void;
  getRulesList(): Array<proto_feature_rule_pb.Rule>;
  setRulesList(value: Array<proto_feature_rule_pb.Rule>): void;
  addRules(value?: proto_feature_rule_pb.Rule, index?: number): proto_feature_rule_pb.Rule;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getVersion(): number;
  setVersion(value: number): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getIncludedUserCount(): number;
  setIncludedUserCount(value: number): void;

  getExcludedUserCount(): number;
  setExcludedUserCount(value: number): void;

  getStatus(): Segment.StatusMap[keyof Segment.StatusMap];
  setStatus(value: Segment.StatusMap[keyof Segment.StatusMap]): void;

  getIsInUseStatus(): boolean;
  setIsInUseStatus(value: boolean): void;

  clearFeaturesList(): void;
  getFeaturesList(): Array<proto_feature_feature_pb.Feature>;
  setFeaturesList(value: Array<proto_feature_feature_pb.Feature>): void;
  addFeatures(value?: proto_feature_feature_pb.Feature, index?: number): proto_feature_feature_pb.Feature;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Segment.AsObject;
  static toObject(includeInstance: boolean, msg: Segment): Segment.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Segment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Segment;
  static deserializeBinaryFromReader(message: Segment, reader: jspb.BinaryReader): Segment;
}

export namespace Segment {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
    rulesList: Array<proto_feature_rule_pb.Rule.AsObject>,
    createdAt: number,
    updatedAt: number,
    version: number,
    deleted: boolean,
    includedUserCount: number,
    excludedUserCount: number,
    status: Segment.StatusMap[keyof Segment.StatusMap],
    isInUseStatus: boolean,
    featuresList: Array<proto_feature_feature_pb.Feature.AsObject>,
  }

  export interface StatusMap {
    INITIAL: 0;
    UPLOADING: 1;
    SUCEEDED: 2;
    FAILED: 3;
  }

  export const Status: StatusMap;
}

export class SegmentUser extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getSegmentId(): string;
  setSegmentId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getState(): SegmentUser.StateMap[keyof SegmentUser.StateMap];
  setState(value: SegmentUser.StateMap[keyof SegmentUser.StateMap]): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentUser.AsObject;
  static toObject(includeInstance: boolean, msg: SegmentUser): SegmentUser.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SegmentUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SegmentUser;
  static deserializeBinaryFromReader(message: SegmentUser, reader: jspb.BinaryReader): SegmentUser;
}

export namespace SegmentUser {
  export type AsObject = {
    id: string,
    segmentId: string,
    userId: string,
    state: SegmentUser.StateMap[keyof SegmentUser.StateMap],
    deleted: boolean,
  }

  export interface StateMap {
    INCLUDED: 0;
    EXCLUDED: 1;
  }

  export const State: StateMap;
}

export class SegmentUsers extends jspb.Message {
  getSegmentId(): string;
  setSegmentId(value: string): void;

  clearUsersList(): void;
  getUsersList(): Array<SegmentUser>;
  setUsersList(value: Array<SegmentUser>): void;
  addUsers(value?: SegmentUser, index?: number): SegmentUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SegmentUsers.AsObject;
  static toObject(includeInstance: boolean, msg: SegmentUsers): SegmentUsers.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SegmentUsers, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SegmentUsers;
  static deserializeBinaryFromReader(message: SegmentUsers, reader: jspb.BinaryReader): SegmentUsers;
}

export namespace SegmentUsers {
  export type AsObject = {
    segmentId: string,
    usersList: Array<SegmentUser.AsObject>,
  }
}

