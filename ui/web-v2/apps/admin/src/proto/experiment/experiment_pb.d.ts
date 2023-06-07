// package: bucketeer.experiment
// file: proto/experiment/experiment.proto

import * as jspb from "google-protobuf";
import * as proto_feature_variation_pb from "../../proto/feature/variation_pb";

export class Experiment extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  clearVariationsList(): void;
  getVariationsList(): Array<proto_feature_variation_pb.Variation>;
  setVariationsList(value: Array<proto_feature_variation_pb.Variation>): void;
  addVariations(value?: proto_feature_variation_pb.Variation, index?: number): proto_feature_variation_pb.Variation;

  getStartAt(): number;
  setStartAt(value: number): void;

  getStopAt(): number;
  setStopAt(value: number): void;

  getStopped(): boolean;
  setStopped(value: boolean): void;

  getStoppedAt(): string;
  setStoppedAt(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  clearGoalIdsList(): void;
  getGoalIdsList(): Array<string>;
  setGoalIdsList(value: Array<string>): void;
  addGoalIds(value: string, index?: number): string;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getBaseVariationId(): string;
  setBaseVariationId(value: string): void;

  getStatus(): Experiment.StatusMap[keyof Experiment.StatusMap];
  setStatus(value: Experiment.StatusMap[keyof Experiment.StatusMap]): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Experiment.AsObject;
  static toObject(includeInstance: boolean, msg: Experiment): Experiment.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Experiment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Experiment;
  static deserializeBinaryFromReader(message: Experiment, reader: jspb.BinaryReader): Experiment;
}

export namespace Experiment {
  export type AsObject = {
    id: string,
    goalId: string,
    featureId: string,
    featureVersion: number,
    variationsList: Array<proto_feature_variation_pb.Variation.AsObject>,
    startAt: number,
    stopAt: number,
    stopped: boolean,
    stoppedAt: string,
    createdAt: number,
    updatedAt: number,
    deleted: boolean,
    goalIdsList: Array<string>,
    name: string,
    description: string,
    baseVariationId: string,
    status: Experiment.StatusMap[keyof Experiment.StatusMap],
    maintainer: string,
    archived: boolean,
  }

  export interface StatusMap {
    WAITING: 0;
    RUNNING: 1;
    STOPPED: 2;
    FORCE_STOPPED: 3;
  }

  export const Status: StatusMap;
}

export class Experiments extends jspb.Message {
  clearExperimentsList(): void;
  getExperimentsList(): Array<Experiment>;
  setExperimentsList(value: Array<Experiment>): void;
  addExperiments(value?: Experiment, index?: number): Experiment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Experiments.AsObject;
  static toObject(includeInstance: boolean, msg: Experiments): Experiments.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Experiments, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Experiments;
  static deserializeBinaryFromReader(message: Experiments, reader: jspb.BinaryReader): Experiments;
}

export namespace Experiments {
  export type AsObject = {
    experimentsList: Array<Experiment.AsObject>,
  }
}

