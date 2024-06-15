// package: bucketeer.experiment
// file: proto/experiment/command.proto

import * as jspb from "google-protobuf";

export class CreateGoalCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateGoalCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateGoalCommand): CreateGoalCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateGoalCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateGoalCommand;
  static deserializeBinaryFromReader(message: CreateGoalCommand, reader: jspb.BinaryReader): CreateGoalCommand;
}

export namespace CreateGoalCommand {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
  }
}

export class RenameGoalCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenameGoalCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RenameGoalCommand): RenameGoalCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenameGoalCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenameGoalCommand;
  static deserializeBinaryFromReader(message: RenameGoalCommand, reader: jspb.BinaryReader): RenameGoalCommand;
}

export namespace RenameGoalCommand {
  export type AsObject = {
    name: string,
  }
}

export class ChangeDescriptionGoalCommand extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeDescriptionGoalCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeDescriptionGoalCommand): ChangeDescriptionGoalCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeDescriptionGoalCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDescriptionGoalCommand;
  static deserializeBinaryFromReader(message: ChangeDescriptionGoalCommand, reader: jspb.BinaryReader): ChangeDescriptionGoalCommand;
}

export namespace ChangeDescriptionGoalCommand {
  export type AsObject = {
    description: string,
  }
}

export class ArchiveGoalCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveGoalCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveGoalCommand): ArchiveGoalCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveGoalCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveGoalCommand;
  static deserializeBinaryFromReader(message: ArchiveGoalCommand, reader: jspb.BinaryReader): ArchiveGoalCommand;
}

export namespace ArchiveGoalCommand {
  export type AsObject = {
  }
}

export class DeleteGoalCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteGoalCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteGoalCommand): DeleteGoalCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteGoalCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteGoalCommand;
  static deserializeBinaryFromReader(message: DeleteGoalCommand, reader: jspb.BinaryReader): DeleteGoalCommand;
}

export namespace DeleteGoalCommand {
  export type AsObject = {
  }
}

export class CreateExperimentCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getStartAt(): number;
  setStartAt(value: number): void;

  getStopAt(): number;
  setStopAt(value: number): void;

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

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateExperimentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateExperimentCommand): CreateExperimentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateExperimentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateExperimentCommand;
  static deserializeBinaryFromReader(message: CreateExperimentCommand, reader: jspb.BinaryReader): CreateExperimentCommand;
}

export namespace CreateExperimentCommand {
  export type AsObject = {
    featureId: string,
    startAt: number,
    stopAt: number,
    goalIdsList: Array<string>,
    name: string,
    description: string,
    baseVariationId: string,
  }
}

export class ChangeExperimentPeriodCommand extends jspb.Message {
  getStartAt(): number;
  setStartAt(value: number): void;

  getStopAt(): number;
  setStopAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeExperimentPeriodCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeExperimentPeriodCommand): ChangeExperimentPeriodCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeExperimentPeriodCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeExperimentPeriodCommand;
  static deserializeBinaryFromReader(message: ChangeExperimentPeriodCommand, reader: jspb.BinaryReader): ChangeExperimentPeriodCommand;
}

export namespace ChangeExperimentPeriodCommand {
  export type AsObject = {
    startAt: number,
    stopAt: number,
  }
}

export class ChangeExperimentNameCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeExperimentNameCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeExperimentNameCommand): ChangeExperimentNameCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeExperimentNameCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeExperimentNameCommand;
  static deserializeBinaryFromReader(message: ChangeExperimentNameCommand, reader: jspb.BinaryReader): ChangeExperimentNameCommand;
}

export namespace ChangeExperimentNameCommand {
  export type AsObject = {
    name: string,
  }
}

export class ChangeExperimentDescriptionCommand extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeExperimentDescriptionCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeExperimentDescriptionCommand): ChangeExperimentDescriptionCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeExperimentDescriptionCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeExperimentDescriptionCommand;
  static deserializeBinaryFromReader(message: ChangeExperimentDescriptionCommand, reader: jspb.BinaryReader): ChangeExperimentDescriptionCommand;
}

export namespace ChangeExperimentDescriptionCommand {
  export type AsObject = {
    description: string,
  }
}

export class StopExperimentCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StopExperimentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: StopExperimentCommand): StopExperimentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StopExperimentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StopExperimentCommand;
  static deserializeBinaryFromReader(message: StopExperimentCommand, reader: jspb.BinaryReader): StopExperimentCommand;
}

export namespace StopExperimentCommand {
  export type AsObject = {
  }
}

export class ArchiveExperimentCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveExperimentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveExperimentCommand): ArchiveExperimentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveExperimentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveExperimentCommand;
  static deserializeBinaryFromReader(message: ArchiveExperimentCommand, reader: jspb.BinaryReader): ArchiveExperimentCommand;
}

export namespace ArchiveExperimentCommand {
  export type AsObject = {
  }
}

export class DeleteExperimentCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteExperimentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteExperimentCommand): DeleteExperimentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteExperimentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteExperimentCommand;
  static deserializeBinaryFromReader(message: DeleteExperimentCommand, reader: jspb.BinaryReader): DeleteExperimentCommand;
}

export namespace DeleteExperimentCommand {
  export type AsObject = {
  }
}

export class StartExperimentCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartExperimentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: StartExperimentCommand): StartExperimentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartExperimentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartExperimentCommand;
  static deserializeBinaryFromReader(message: StartExperimentCommand, reader: jspb.BinaryReader): StartExperimentCommand;
}

export namespace StartExperimentCommand {
  export type AsObject = {
  }
}

export class FinishExperimentCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FinishExperimentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: FinishExperimentCommand): FinishExperimentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FinishExperimentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FinishExperimentCommand;
  static deserializeBinaryFromReader(message: FinishExperimentCommand, reader: jspb.BinaryReader): FinishExperimentCommand;
}

export namespace FinishExperimentCommand {
  export type AsObject = {
  }
}

