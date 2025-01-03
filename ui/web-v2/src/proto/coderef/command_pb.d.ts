// package: bucketeer.coderef
// file: proto/coderef/command.proto

import * as jspb from 'google-protobuf';
import * as proto_coderef_code_reference_pb from '../../proto/coderef/code_reference_pb';

export class CreateCodeReferenceCommand extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFilePath(): string;
  setFilePath(value: string): void;

  getLineNumber(): number;
  setLineNumber(value: number): void;

  getCodeSnippet(): string;
  setCodeSnippet(value: string): void;

  getContentHash(): string;
  setContentHash(value: string): void;

  clearAliasesList(): void;
  getAliasesList(): Array<string>;
  setAliasesList(value: Array<string>): void;
  addAliases(value: string, index?: number): string;

  getRepositoryName(): string;
  setRepositoryName(value: string): void;

  getRepositoryOwner(): string;
  setRepositoryOwner(value: string): void;

  getRepositoryType(): proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap[keyof proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap];
  setRepositoryType(
    value: proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap[keyof proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap]
  ): void;

  getRepositoryBranch(): string;
  setRepositoryBranch(value: string): void;

  getCommitHash(): string;
  setCommitHash(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCodeReferenceCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateCodeReferenceCommand
  ): CreateCodeReferenceCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateCodeReferenceCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateCodeReferenceCommand;
  static deserializeBinaryFromReader(
    message: CreateCodeReferenceCommand,
    reader: jspb.BinaryReader
  ): CreateCodeReferenceCommand;
}

export namespace CreateCodeReferenceCommand {
  export type AsObject = {
    featureId: string;
    filePath: string;
    lineNumber: number;
    codeSnippet: string;
    contentHash: string;
    aliasesList: Array<string>;
    repositoryName: string;
    repositoryOwner: string;
    repositoryType: proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap[keyof proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap];
    repositoryBranch: string;
    commitHash: string;
    environmentId: string;
  };
}

export class UpdateCodeReferenceCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFilePath(): string;
  setFilePath(value: string): void;

  getLineNumber(): number;
  setLineNumber(value: number): void;

  getCodeSnippet(): string;
  setCodeSnippet(value: string): void;

  getContentHash(): string;
  setContentHash(value: string): void;

  clearAliasesList(): void;
  getAliasesList(): Array<string>;
  setAliasesList(value: Array<string>): void;
  addAliases(value: string, index?: number): string;

  getRepositoryBranch(): string;
  setRepositoryBranch(value: string): void;

  getCommitHash(): string;
  setCommitHash(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCodeReferenceCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateCodeReferenceCommand
  ): UpdateCodeReferenceCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateCodeReferenceCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCodeReferenceCommand;
  static deserializeBinaryFromReader(
    message: UpdateCodeReferenceCommand,
    reader: jspb.BinaryReader
  ): UpdateCodeReferenceCommand;
}

export namespace UpdateCodeReferenceCommand {
  export type AsObject = {
    id: string;
    filePath: string;
    lineNumber: number;
    codeSnippet: string;
    contentHash: string;
    aliasesList: Array<string>;
    repositoryBranch: string;
    commitHash: string;
  };
}

export class DeleteCodeReferenceCommand extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCodeReferenceCommand.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteCodeReferenceCommand
  ): DeleteCodeReferenceCommand.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteCodeReferenceCommand,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCodeReferenceCommand;
  static deserializeBinaryFromReader(
    message: DeleteCodeReferenceCommand,
    reader: jspb.BinaryReader
  ): DeleteCodeReferenceCommand;
}

export namespace DeleteCodeReferenceCommand {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}
