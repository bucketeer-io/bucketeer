// package: bucketeer.notification.sender
// file: proto/notification/sender/notification.proto

import * as jspb from "google-protobuf";
import * as proto_event_domain_event_pb from "../../../proto/event/domain/event_pb";
import * as proto_feature_feature_pb from "../../../proto/feature/feature_pb";
import * as proto_experiment_experiment_pb from "../../../proto/experiment/experiment_pb";

export class Notification extends jspb.Message {
  getType(): Notification.TypeMap[keyof Notification.TypeMap];
  setType(value: Notification.TypeMap[keyof Notification.TypeMap]): void;

  hasDomainEventNotification(): boolean;
  clearDomainEventNotification(): void;
  getDomainEventNotification(): DomainEventNotification | undefined;
  setDomainEventNotification(value?: DomainEventNotification): void;

  hasFeatureStaleNotification(): boolean;
  clearFeatureStaleNotification(): void;
  getFeatureStaleNotification(): FeatureStaleNotification | undefined;
  setFeatureStaleNotification(value?: FeatureStaleNotification): void;

  hasExperimentRunningNotification(): boolean;
  clearExperimentRunningNotification(): void;
  getExperimentRunningNotification(): ExperimentRunningNotification | undefined;
  setExperimentRunningNotification(value?: ExperimentRunningNotification): void;

  hasMauCountNotification(): boolean;
  clearMauCountNotification(): void;
  getMauCountNotification(): MauCountNotification | undefined;
  setMauCountNotification(value?: MauCountNotification): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Notification.AsObject;
  static toObject(includeInstance: boolean, msg: Notification): Notification.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Notification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Notification;
  static deserializeBinaryFromReader(message: Notification, reader: jspb.BinaryReader): Notification;
}

export namespace Notification {
  export type AsObject = {
    type: Notification.TypeMap[keyof Notification.TypeMap],
    domainEventNotification?: DomainEventNotification.AsObject,
    featureStaleNotification?: FeatureStaleNotification.AsObject,
    experimentRunningNotification?: ExperimentRunningNotification.AsObject,
    mauCountNotification?: MauCountNotification.AsObject,
  }

  export interface TypeMap {
    DOMAINEVENT: 0;
    FEATURESTALE: 1;
    EXPERIMENTRUNNING: 2;
    MAUCOUNT: 3;
  }

  export const Type: TypeMap;
}

export class DomainEventNotification extends jspb.Message {
  hasEditor(): boolean;
  clearEditor(): void;
  getEditor(): proto_event_domain_event_pb.Editor | undefined;
  setEditor(value?: proto_event_domain_event_pb.Editor): void;

  getEntityType(): proto_event_domain_event_pb.Event.EntityTypeMap[keyof proto_event_domain_event_pb.Event.EntityTypeMap];
  setEntityType(value: proto_event_domain_event_pb.Event.EntityTypeMap[keyof proto_event_domain_event_pb.Event.EntityTypeMap]): void;

  getEntityId(): string;
  setEntityId(value: string): void;

  getType(): proto_event_domain_event_pb.Event.TypeMap[keyof proto_event_domain_event_pb.Event.TypeMap];
  setType(value: proto_event_domain_event_pb.Event.TypeMap[keyof proto_event_domain_event_pb.Event.TypeMap]): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainEventNotification.AsObject;
  static toObject(includeInstance: boolean, msg: DomainEventNotification): DomainEventNotification.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DomainEventNotification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainEventNotification;
  static deserializeBinaryFromReader(message: DomainEventNotification, reader: jspb.BinaryReader): DomainEventNotification;
}

export namespace DomainEventNotification {
  export type AsObject = {
    editor?: proto_event_domain_event_pb.Editor.AsObject,
    entityType: proto_event_domain_event_pb.Event.EntityTypeMap[keyof proto_event_domain_event_pb.Event.EntityTypeMap],
    entityId: string,
    type: proto_event_domain_event_pb.Event.TypeMap[keyof proto_event_domain_event_pb.Event.TypeMap],
    environmentId: string,
  }
}

export class FeatureStaleNotification extends jspb.Message {
  clearFeaturesList(): void;
  getFeaturesList(): Array<proto_feature_feature_pb.Feature>;
  setFeaturesList(value: Array<proto_feature_feature_pb.Feature>): void;
  addFeatures(value?: proto_feature_feature_pb.Feature, index?: number): proto_feature_feature_pb.Feature;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureStaleNotification.AsObject;
  static toObject(includeInstance: boolean, msg: FeatureStaleNotification): FeatureStaleNotification.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FeatureStaleNotification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FeatureStaleNotification;
  static deserializeBinaryFromReader(message: FeatureStaleNotification, reader: jspb.BinaryReader): FeatureStaleNotification;
}

export namespace FeatureStaleNotification {
  export type AsObject = {
    featuresList: Array<proto_feature_feature_pb.Feature.AsObject>,
    environmentId: string,
  }
}

export class ExperimentRunningNotification extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  clearExperimentsList(): void;
  getExperimentsList(): Array<proto_experiment_experiment_pb.Experiment>;
  setExperimentsList(value: Array<proto_experiment_experiment_pb.Experiment>): void;
  addExperiments(value?: proto_experiment_experiment_pb.Experiment, index?: number): proto_experiment_experiment_pb.Experiment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExperimentRunningNotification.AsObject;
  static toObject(includeInstance: boolean, msg: ExperimentRunningNotification): ExperimentRunningNotification.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExperimentRunningNotification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExperimentRunningNotification;
  static deserializeBinaryFromReader(message: ExperimentRunningNotification, reader: jspb.BinaryReader): ExperimentRunningNotification;
}

export namespace ExperimentRunningNotification {
  export type AsObject = {
    environmentId: string,
    experimentsList: Array<proto_experiment_experiment_pb.Experiment.AsObject>,
  }
}

export class MauCountNotification extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getEventCount(): number;
  setEventCount(value: number): void;

  getUserCount(): number;
  setUserCount(value: number): void;

  getMonth(): number;
  setMonth(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MauCountNotification.AsObject;
  static toObject(includeInstance: boolean, msg: MauCountNotification): MauCountNotification.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MauCountNotification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MauCountNotification;
  static deserializeBinaryFromReader(message: MauCountNotification, reader: jspb.BinaryReader): MauCountNotification;
}

export namespace MauCountNotification {
  export type AsObject = {
    environmentId: string,
    eventCount: number,
    userCount: number,
    month: number,
  }
}

