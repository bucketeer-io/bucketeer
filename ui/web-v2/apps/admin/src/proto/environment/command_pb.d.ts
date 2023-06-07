// package: bucketeer.environment
// file: proto/environment/command.proto

import * as jspb from "google-protobuf";

export class CreateEnvironmentCommand extends jspb.Message {
  getNamespace(): string;
  setNamespace(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getId(): string;
  setId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateEnvironmentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: CreateEnvironmentCommand): CreateEnvironmentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateEnvironmentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateEnvironmentCommand;
  static deserializeBinaryFromReader(message: CreateEnvironmentCommand, reader: jspb.BinaryReader): CreateEnvironmentCommand;
}

export namespace CreateEnvironmentCommand {
  export type AsObject = {
    namespace: string,
    name: string,
    description: string,
    id: string,
    projectId: string,
  }
}

export class RenameEnvironmentCommand extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenameEnvironmentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: RenameEnvironmentCommand): RenameEnvironmentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RenameEnvironmentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenameEnvironmentCommand;
  static deserializeBinaryFromReader(message: RenameEnvironmentCommand, reader: jspb.BinaryReader): RenameEnvironmentCommand;
}

export namespace RenameEnvironmentCommand {
  export type AsObject = {
    name: string,
  }
}

export class ChangeDescriptionEnvironmentCommand extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeDescriptionEnvironmentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeDescriptionEnvironmentCommand): ChangeDescriptionEnvironmentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangeDescriptionEnvironmentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeDescriptionEnvironmentCommand;
  static deserializeBinaryFromReader(message: ChangeDescriptionEnvironmentCommand, reader: jspb.BinaryReader): ChangeDescriptionEnvironmentCommand;
}

export namespace ChangeDescriptionEnvironmentCommand {
  export type AsObject = {
    description: string,
  }
}

export class DeleteEnvironmentCommand extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteEnvironmentCommand.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteEnvironmentCommand): DeleteEnvironmentCommand.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteEnvironmentCommand, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteEnvironmentCommand;
  static deserializeBinaryFromReader(message: DeleteEnvironmentCommand, reader: jspb.BinaryReader): DeleteEnvironmentCommand;
}

export namespace DeleteEnvironmentCommand {
  export type AsObject = {
  }
}

export class CreateProjectCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

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
  }
}

export class CreateTrialProjectCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

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

