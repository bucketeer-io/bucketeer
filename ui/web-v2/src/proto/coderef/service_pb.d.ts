// package: bucketeer.coderef
// file: proto/coderef/service.proto

import * as jspb from 'google-protobuf';
import * as proto_coderef_code_reference_pb from '../../proto/coderef/code_reference_pb';
import * as google_api_annotations_pb from '../../google/api/annotations_pb';
import * as protoc_gen_openapiv2_options_annotations_pb from '../../protoc-gen-openapiv2/options/annotations_pb';

export class GetCodeReferenceRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCodeReferenceRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetCodeReferenceRequest
  ): GetCodeReferenceRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetCodeReferenceRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetCodeReferenceRequest;
  static deserializeBinaryFromReader(
    message: GetCodeReferenceRequest,
    reader: jspb.BinaryReader
  ): GetCodeReferenceRequest;
}

export namespace GetCodeReferenceRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}

export class GetCodeReferenceResponse extends jspb.Message {
  hasCodeReference(): boolean;
  clearCodeReference(): void;
  getCodeReference(): proto_coderef_code_reference_pb.CodeReference | undefined;
  setCodeReference(value?: proto_coderef_code_reference_pb.CodeReference): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCodeReferenceResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: GetCodeReferenceResponse
  ): GetCodeReferenceResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: GetCodeReferenceResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): GetCodeReferenceResponse;
  static deserializeBinaryFromReader(
    message: GetCodeReferenceResponse,
    reader: jspb.BinaryReader
  ): GetCodeReferenceResponse;
}

export namespace GetCodeReferenceResponse {
  export type AsObject = {
    codeReference?: proto_coderef_code_reference_pb.CodeReference.AsObject;
  };
}

export class ListCodeReferencesRequest extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

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

  getCursor(): string;
  setCursor(value: string): void;

  getPageSize(): number;
  setPageSize(value: number): void;

  getOrderBy(): ListCodeReferencesRequest.OrderByMap[keyof ListCodeReferencesRequest.OrderByMap];
  setOrderBy(
    value: ListCodeReferencesRequest.OrderByMap[keyof ListCodeReferencesRequest.OrderByMap]
  ): void;

  getOrderDirection(): ListCodeReferencesRequest.OrderDirectionMap[keyof ListCodeReferencesRequest.OrderDirectionMap];
  setOrderDirection(
    value: ListCodeReferencesRequest.OrderDirectionMap[keyof ListCodeReferencesRequest.OrderDirectionMap]
  ): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCodeReferencesRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListCodeReferencesRequest
  ): ListCodeReferencesRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListCodeReferencesRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListCodeReferencesRequest;
  static deserializeBinaryFromReader(
    message: ListCodeReferencesRequest,
    reader: jspb.BinaryReader
  ): ListCodeReferencesRequest;
}

export namespace ListCodeReferencesRequest {
  export type AsObject = {
    environmentId: string;
    featureId: string;
    repositoryName: string;
    repositoryOwner: string;
    repositoryType: proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap[keyof proto_coderef_code_reference_pb.CodeReference.RepositoryTypeMap];
    repositoryBranch: string;
    cursor: string;
    pageSize: number;
    orderBy: ListCodeReferencesRequest.OrderByMap[keyof ListCodeReferencesRequest.OrderByMap];
    orderDirection: ListCodeReferencesRequest.OrderDirectionMap[keyof ListCodeReferencesRequest.OrderDirectionMap];
  };

  export interface OrderByMap {
    DEFAULT: 0;
    CREATED_AT: 1;
    UPDATED_AT: 2;
  }

  export const OrderBy: OrderByMap;

  export interface OrderDirectionMap {
    ASC: 0;
    DESC: 1;
  }

  export const OrderDirection: OrderDirectionMap;
}

export class ListCodeReferencesResponse extends jspb.Message {
  clearCodeReferencesList(): void;
  getCodeReferencesList(): Array<proto_coderef_code_reference_pb.CodeReference>;
  setCodeReferencesList(
    value: Array<proto_coderef_code_reference_pb.CodeReference>
  ): void;
  addCodeReferences(
    value?: proto_coderef_code_reference_pb.CodeReference,
    index?: number
  ): proto_coderef_code_reference_pb.CodeReference;

  getCursor(): string;
  setCursor(value: string): void;

  getTotalCount(): number;
  setTotalCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCodeReferencesResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ListCodeReferencesResponse
  ): ListCodeReferencesResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ListCodeReferencesResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ListCodeReferencesResponse;
  static deserializeBinaryFromReader(
    message: ListCodeReferencesResponse,
    reader: jspb.BinaryReader
  ): ListCodeReferencesResponse;
}

export namespace ListCodeReferencesResponse {
  export type AsObject = {
    codeReferencesList: Array<proto_coderef_code_reference_pb.CodeReference.AsObject>;
    cursor: string;
    totalCount: number;
  };
}

export class CreateCodeReferenceRequest extends jspb.Message {
  getFeatureId(): string;
  setFeatureId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

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

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCodeReferenceRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateCodeReferenceRequest
  ): CreateCodeReferenceRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateCodeReferenceRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateCodeReferenceRequest;
  static deserializeBinaryFromReader(
    message: CreateCodeReferenceRequest,
    reader: jspb.BinaryReader
  ): CreateCodeReferenceRequest;
}

export namespace CreateCodeReferenceRequest {
  export type AsObject = {
    featureId: string;
    environmentId: string;
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
  };
}

export class CreateCodeReferenceResponse extends jspb.Message {
  hasCodeReference(): boolean;
  clearCodeReference(): void;
  getCodeReference(): proto_coderef_code_reference_pb.CodeReference | undefined;
  setCodeReference(value?: proto_coderef_code_reference_pb.CodeReference): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCodeReferenceResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: CreateCodeReferenceResponse
  ): CreateCodeReferenceResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: CreateCodeReferenceResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): CreateCodeReferenceResponse;
  static deserializeBinaryFromReader(
    message: CreateCodeReferenceResponse,
    reader: jspb.BinaryReader
  ): CreateCodeReferenceResponse;
}

export namespace CreateCodeReferenceResponse {
  export type AsObject = {
    codeReference?: proto_coderef_code_reference_pb.CodeReference.AsObject;
  };
}

export class UpdateCodeReferenceRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

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

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCodeReferenceRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateCodeReferenceRequest
  ): UpdateCodeReferenceRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateCodeReferenceRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCodeReferenceRequest;
  static deserializeBinaryFromReader(
    message: UpdateCodeReferenceRequest,
    reader: jspb.BinaryReader
  ): UpdateCodeReferenceRequest;
}

export namespace UpdateCodeReferenceRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
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
  };
}

export class UpdateCodeReferenceResponse extends jspb.Message {
  hasCodeReference(): boolean;
  clearCodeReference(): void;
  getCodeReference(): proto_coderef_code_reference_pb.CodeReference | undefined;
  setCodeReference(value?: proto_coderef_code_reference_pb.CodeReference): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCodeReferenceResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: UpdateCodeReferenceResponse
  ): UpdateCodeReferenceResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: UpdateCodeReferenceResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCodeReferenceResponse;
  static deserializeBinaryFromReader(
    message: UpdateCodeReferenceResponse,
    reader: jspb.BinaryReader
  ): UpdateCodeReferenceResponse;
}

export namespace UpdateCodeReferenceResponse {
  export type AsObject = {
    codeReference?: proto_coderef_code_reference_pb.CodeReference.AsObject;
  };
}

export class DeleteCodeReferenceRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCodeReferenceRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteCodeReferenceRequest
  ): DeleteCodeReferenceRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteCodeReferenceRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCodeReferenceRequest;
  static deserializeBinaryFromReader(
    message: DeleteCodeReferenceRequest,
    reader: jspb.BinaryReader
  ): DeleteCodeReferenceRequest;
}

export namespace DeleteCodeReferenceRequest {
  export type AsObject = {
    id: string;
    environmentId: string;
  };
}

export class DeleteCodeReferenceResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteCodeReferenceResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: DeleteCodeReferenceResponse
  ): DeleteCodeReferenceResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: DeleteCodeReferenceResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): DeleteCodeReferenceResponse;
  static deserializeBinaryFromReader(
    message: DeleteCodeReferenceResponse,
    reader: jspb.BinaryReader
  ): DeleteCodeReferenceResponse;
}

export namespace DeleteCodeReferenceResponse {
  export type AsObject = {};
}
