// package: bucketeer.environment
// file: proto/environment/command.proto

import * as jspb from "google-protobuf";

export class CreateEnvironmentV2Command extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnvironmentV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnvironmentV2Command): CreateEnvironmentV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateEnvironmentV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnvironmentV2Command;
  static deserializeBinaryFromReader(message: CreateEnvironmentV2Command, reader: jspb.BinaryReader): CreateEnvironmentV2Command;
}

export namespace CreateEnvironmentV2Command {
  export type AsObject = {
    name: string,
    urlCode: string,
    description: string,
    projectId: string,
  }
}

export class RenameEnvironmentV2Command extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenameEnvironmentV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: RenameEnvironmentV2Command): RenameEnvironmentV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenameEnvironmentV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenameEnvironmentV2Command;
  static deserializeBinaryFromReader(message: RenameEnvironmentV2Command, reader: jspb.BinaryReader): RenameEnvironmentV2Command;
}

export namespace RenameEnvironmentV2Command {
  export type AsObject = {
    name: string,
  }
}

export class ChangeDescriptionEnvironmentV2Command extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeDescriptionEnvironmentV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeDescriptionEnvironmentV2Command): ChangeDescriptionEnvironmentV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeDescriptionEnvironmentV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDescriptionEnvironmentV2Command;
  static deserializeBinaryFromReader(message: ChangeDescriptionEnvironmentV2Command, reader: jspb.BinaryReader): ChangeDescriptionEnvironmentV2Command;
}

export namespace ChangeDescriptionEnvironmentV2Command {
  export type AsObject = {
    description: string,
  }
}

export class ArchiveEnvironmentV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArchiveEnvironmentV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: ArchiveEnvironmentV2Command): ArchiveEnvironmentV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ArchiveEnvironmentV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArchiveEnvironmentV2Command;
  static deserializeBinaryFromReader(message: ArchiveEnvironmentV2Command, reader: jspb.BinaryReader): ArchiveEnvironmentV2Command;
}

export namespace ArchiveEnvironmentV2Command {
  export type AsObject = {
  }
}

export class UnarchiveEnvironmentV2Command extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnarchiveEnvironmentV2Command.AsObject;
  static toObject(includeInstance: boolean, msg: UnarchiveEnvironmentV2Command): UnarchiveEnvironmentV2Command.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnarchiveEnvironmentV2Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnarchiveEnvironmentV2Command;
  static deserializeBinaryFromReader(message: UnarchiveEnvironmentV2Command, reader: jspb.BinaryReader): UnarchiveEnvironmentV2Command;
}

export namespace UnarchiveEnvironmentV2Command {
  export type AsObject = {
  }
}

export class CreateProjectCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateProjectCommand): CreateProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateProjectCommand;
  static deserializeBinaryFromReader(message: CreateProjectCommand, reader: jspb.BinaryReader): CreateProjectCommand;
}

export namespace CreateProjectCommand {
  export type AsObject = {
    id: string,
    description: string,
    name: string,
    urlCode: string,
  }
}

export class CreateTrialProjectCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getName(): string;
  setName(value: string): void;

  getUrlCode(): string;
  setUrlCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTrialProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTrialProjectCommand): CreateTrialProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateTrialProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTrialProjectCommand;
  static deserializeBinaryFromReader(message: CreateTrialProjectCommand, reader: jspb.BinaryReader): CreateTrialProjectCommand;
}

export namespace CreateTrialProjectCommand {
  export type AsObject = {
    id: string,
    email: string,
    name: string,
    urlCode: string,
  }
}

export class RenameProjectCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenameProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RenameProjectCommand): RenameProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenameProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenameProjectCommand;
  static deserializeBinaryFromReader(message: RenameProjectCommand, reader: jspb.BinaryReader): RenameProjectCommand;
}

export namespace RenameProjectCommand {
  export type AsObject = {
    name: string,
  }
}

export class ChangeDescriptionProjectCommand extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeDescriptionProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeDescriptionProjectCommand): ChangeDescriptionProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeDescriptionProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDescriptionProjectCommand;
  static deserializeBinaryFromReader(message: ChangeDescriptionProjectCommand, reader: jspb.BinaryReader): ChangeDescriptionProjectCommand;
}

export namespace ChangeDescriptionProjectCommand {
  export type AsObject = {
    description: string,
  }
}

export class EnableProjectCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: EnableProjectCommand): EnableProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableProjectCommand;
  static deserializeBinaryFromReader(message: EnableProjectCommand, reader: jspb.BinaryReader): EnableProjectCommand;
}

export namespace EnableProjectCommand {
  export type AsObject = {
  }
}

export class DisableProjectCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DisableProjectCommand): DisableProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableProjectCommand;
  static deserializeBinaryFromReader(message: DisableProjectCommand, reader: jspb.BinaryReader): DisableProjectCommand;
}

export namespace DisableProjectCommand {
  export type AsObject = {
  }
}

export class ConvertTrialProjectCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertTrialProjectCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertTrialProjectCommand): ConvertTrialProjectCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConvertTrialProjectCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertTrialProjectCommand;
  static deserializeBinaryFromReader(message: ConvertTrialProjectCommand, reader: jspb.BinaryReader): ConvertTrialProjectCommand;
}

export namespace ConvertTrialProjectCommand {
  export type AsObject = {
  }
}

