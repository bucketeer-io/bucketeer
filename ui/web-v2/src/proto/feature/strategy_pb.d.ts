// package: bucketeer.feature
// file: proto/feature/strategy.proto

import * as jspb from "google-protobuf";

export class FixedStrategy extends jspb.Message {
  getVariation(): string;
  setVariation(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FixedStrategy.AsObject;
  static toObject(includeInstance: boolean, msg: FixedStrategy): FixedStrategy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FixedStrategy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FixedStrategy;
  static deserializeBinaryFromReader(message: FixedStrategy, reader: jspb.BinaryReader): FixedStrategy;
}

export namespace FixedStrategy {
  export type AsObject = {
    variation: string,
  }
}

export class RolloutStrategy extends jspb.Message {
  clearVariationsList(): void;
  getVariationsList(): Array<RolloutStrategy.Variation>;
  setVariationsList(value: Array<RolloutStrategy.Variation>): void;
  addVariations(value?: RolloutStrategy.Variation, index?: number): RolloutStrategy.Variation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RolloutStrategy.AsObject;
  static toObject(includeInstance: boolean, msg: RolloutStrategy): RolloutStrategy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RolloutStrategy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RolloutStrategy;
  static deserializeBinaryFromReader(message: RolloutStrategy, reader: jspb.BinaryReader): RolloutStrategy;
}

export namespace RolloutStrategy {
  export type AsObject = {
    variationsList: Array<RolloutStrategy.Variation.AsObject>,
  }

  export class Variation extends jspb.Message {
    getVariation(): string;
    setVariation(value: string): void;

    getWeight(): number;
    setWeight(value: number): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Variation.AsObject;
    static toObject(includeInstance: boolean, msg: Variation): Variation.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Variation, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Variation;
    static deserializeBinaryFromReader(message: Variation, reader: jspb.BinaryReader): Variation;
  }

  export namespace Variation {
    export type AsObject = {
      variation: string,
      weight: number,
    }
  }
}

export class Strategy extends jspb.Message {
  getType(): Strategy.TypeMap[keyof Strategy.TypeMap];
  setType(value: Strategy.TypeMap[keyof Strategy.TypeMap]): void;

  hasFixedStrategy(): boolean;
  clearFixedStrategy(): void;
  getFixedStrategy(): FixedStrategy | undefined;
  setFixedStrategy(value?: FixedStrategy): void;

  hasRolloutStrategy(): boolean;
  clearRolloutStrategy(): void;
  getRolloutStrategy(): RolloutStrategy | undefined;
  setRolloutStrategy(value?: RolloutStrategy): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Strategy.AsObject;
  static toObject(includeInstance: boolean, msg: Strategy): Strategy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Strategy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Strategy;
  static deserializeBinaryFromReader(message: Strategy, reader: jspb.BinaryReader): Strategy;
}

export namespace Strategy {
  export type AsObject = {
    type: Strategy.TypeMap[keyof Strategy.TypeMap],
    fixedStrategy?: FixedStrategy.AsObject,
    rolloutStrategy?: RolloutStrategy.AsObject,
  }

  export interface TypeMap {
    FIXED: 0;
    ROLLOUT: 1;
  }

  export const Type: TypeMap;
}

