// package: bucketeer.coderef
// file: proto/coderef/code_reference.proto

import * as jspb from 'google-protobuf';

export class CodeReference extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getRepositoryType(): CodeReference.RepositoryTypeMap[keyof CodeReference.RepositoryTypeMap];
  setRepositoryType(
    value: CodeReference.RepositoryTypeMap[keyof CodeReference.RepositoryTypeMap]
  ): void;

  getRepositoryBranch(): string;
  setRepositoryBranch(value: string): void;

  getCommitHash(): string;
  setCommitHash(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CodeReference.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CodeReference
  ): CodeReference.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CodeReference,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CodeReference;
  static deserializeBinaryFromReader(
    message: CodeReference,
    reader: jspb.BinaryReader
  ): CodeReference;
}

export namespace CodeReference {
  export type AsObject = {
    id: string;
    featureId: string;
    filePath: string;
    lineNumber: number;
    codeSnippet: string;
    contentHash: string;
    aliasesList: Array<string>;
    repositoryName: string;
    repositoryOwner: string;
    repositoryType: CodeReference.RepositoryTypeMap[keyof CodeReference.RepositoryTypeMap];
    repositoryBranch: string;
    commitHash: string;
    environmentId: string;
    createdAt: number;
    updatedAt: number;
  };

  export interface RepositoryTypeMap {
    REPOSITORY_TYPE_UNSPECIFIED: 0;
    GITHUB: 1;
    GITLAB: 2;
    BITBUCKET: 3;
    CUSTOM: 4;
  }

  export const RepositoryType: RepositoryTypeMap;
}
