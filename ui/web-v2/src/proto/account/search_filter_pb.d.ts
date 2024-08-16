// package: bucketeer.account
// file: proto/account/search_filter.proto

import * as jspb from 'google-protobuf';

export class SearchFilter extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getQuery(): string;
  setQuery(value: string): void;

  getFilterTargetType(): FilterTargetTypeMap[keyof FilterTargetTypeMap];
  setFilterTargetType(
    value: FilterTargetTypeMap[keyof FilterTargetTypeMap]
  ): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getDefaultFilter(): boolean;
  setDefaultFilter(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchFilter.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: SearchFilter
  ): SearchFilter.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: SearchFilter,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): SearchFilter;
  static deserializeBinaryFromReader(
    message: SearchFilter,
    reader: jspb.BinaryReader
  ): SearchFilter;
}

export namespace SearchFilter {
  export type AsObject = {
    id: string;
    name: string;
    query: string;
    filterTargetType: FilterTargetTypeMap[keyof FilterTargetTypeMap];
    environmentId: string;
    defaultFilter: boolean;
  };
}

export interface FilterTargetTypeMap {
  UNKNOWN: 0;
  FEATURE_FLAG: 1;
  GOAL: 2;
}

export const FilterTargetType: FilterTargetTypeMap;
