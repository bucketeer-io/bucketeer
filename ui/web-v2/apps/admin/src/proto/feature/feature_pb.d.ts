// package: bucketeer.feature
// file: proto/feature/feature.proto

import * as jspb from "google-protobuf";
import * as proto_feature_rule_pb from "../../proto/feature/rule_pb";
import * as proto_feature_target_pb from "../../proto/feature/target_pb";
import * as proto_feature_variation_pb from "../../proto/feature/variation_pb";
import * as proto_feature_strategy_pb from "../../proto/feature/strategy_pb";
import * as proto_feature_feature_last_used_info_pb from "../../proto/feature/feature_last_used_info_pb";
import * as proto_feature_prerequisite_pb from "../../proto/feature/prerequisite_pb";

export class Feature extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getEnabled(): boolean;
  setEnabled(value: boolean): void;

  getDeleted(): boolean;
  setDeleted(value: boolean): void;

  getEvaluationUndelayable(): boolean;
  setEvaluationUndelayable(value: boolean): void;

  getTtl(): number;
  setTtl(value: number): void;

  getVersion(): number;
  setVersion(value: number): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  clearVariationsList(): void;
  getVariationsList(): Array<proto_feature_variation_pb.Variation>;
  setVariationsList(value: Array<proto_feature_variation_pb.Variation>): void;
  addVariations(value?: proto_feature_variation_pb.Variation, index?: number): proto_feature_variation_pb.Variation;

  clearTargetsList(): void;
  getTargetsList(): Array<proto_feature_target_pb.Target>;
  setTargetsList(value: Array<proto_feature_target_pb.Target>): void;
  addTargets(value?: proto_feature_target_pb.Target, index?: number): proto_feature_target_pb.Target;

  clearRulesList(): void;
  getRulesList(): Array<proto_feature_rule_pb.Rule>;
  setRulesList(value: Array<proto_feature_rule_pb.Rule>): void;
  addRules(value?: proto_feature_rule_pb.Rule, index?: number): proto_feature_rule_pb.Rule;

  hasDefaultStrategy(): boolean;
  clearDefaultStrategy(): void;
  getDefaultStrategy(): proto_feature_strategy_pb.Strategy | undefined;
  setDefaultStrategy(value?: proto_feature_strategy_pb.Strategy): void;

  getOffVariation(): string;
  setOffVariation(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): void;
  addTags(value: string, index?: number): string;

  hasLastUsedInfo(): boolean;
  clearLastUsedInfo(): void;
  getLastUsedInfo(): proto_feature_feature_last_used_info_pb.FeatureLastUsedInfo | undefined;
  setLastUsedInfo(value?: proto_feature_feature_last_used_info_pb.FeatureLastUsedInfo): void;

  getMaintainer(): string;
  setMaintainer(value: string): void;

  getVariationType(): Feature.VariationTypeMap[keyof Feature.VariationTypeMap];
  setVariationType(value: Feature.VariationTypeMap[keyof Feature.VariationTypeMap]): void;

  getArchived(): boolean;
  setArchived(value: boolean): void;

  clearPrerequisitesList(): void;
  getPrerequisitesList(): Array<proto_feature_prerequisite_pb.Prerequisite>;
  setPrerequisitesList(value: Array<proto_feature_prerequisite_pb.Prerequisite>): void;
  addPrerequisites(value?: proto_feature_prerequisite_pb.Prerequisite, index?: number): proto_feature_prerequisite_pb.Prerequisite;

  getSamplingSeed(): string;
  setSamplingSeed(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Feature.AsObject;
  static toObject(includeInstance: boolean, msg: Feature): Feature.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Feature, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Feature;
  static deserializeBinaryFromReader(message: Feature, reader: jspb.BinaryReader): Feature;
}

export namespace Feature {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
    enabled: boolean,
    deleted: boolean,
    evaluationUndelayable: boolean,
    ttl: number,
    version: number,
    createdAt: number,
    updatedAt: number,
    variationsList: Array<proto_feature_variation_pb.Variation.AsObject>,
    targetsList: Array<proto_feature_target_pb.Target.AsObject>,
    rulesList: Array<proto_feature_rule_pb.Rule.AsObject>,
    defaultStrategy?: proto_feature_strategy_pb.Strategy.AsObject,
    offVariation: string,
    tagsList: Array<string>,
    lastUsedInfo?: proto_feature_feature_last_used_info_pb.FeatureLastUsedInfo.AsObject,
    maintainer: string,
    variationType: Feature.VariationTypeMap[keyof Feature.VariationTypeMap],
    archived: boolean,
    prerequisitesList: Array<proto_feature_prerequisite_pb.Prerequisite.AsObject>,
    samplingSeed: string,
  }

  export interface VariationTypeMap {
    STRING: 0;
    BOOLEAN: 1;
    NUMBER: 2;
    JSON: 3;
  }

  export const VariationType: VariationTypeMap;
}

export class Features extends jspb.Message {
  clearFeaturesList(): void;
  getFeaturesList(): Array<Feature>;
  setFeaturesList(value: Array<Feature>): void;
  addFeatures(value?: Feature, index?: number): Feature;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Features.AsObject;
  static toObject(includeInstance: boolean, msg: Features): Features.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Features, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Features;
  static deserializeBinaryFromReader(message: Features, reader: jspb.BinaryReader): Features;
}

export namespace Features {
  export type AsObject = {
    featuresList: Array<Feature.AsObject>,
  }
}

export class Tag extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tag.AsObject;
  static toObject(includeInstance: boolean, msg: Tag): Tag.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Tag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tag;
  static deserializeBinaryFromReader(message: Tag, reader: jspb.BinaryReader): Tag;
}

export namespace Tag {
  export type AsObject = {
    id: string,
    createdAt: number,
    updatedAt: number,
  }
}

