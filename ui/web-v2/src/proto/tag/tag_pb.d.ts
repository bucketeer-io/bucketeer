// package: bucketeer.tag
// file: proto/tag/tag.proto

import * as jspb from 'google-protobuf';

export class Tag extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getCreatedAt(): number;
  setCreatedAt(value: number): void;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): void;

  getEntityType(): Tag.EntityTypeMap[keyof Tag.EntityTypeMap];
  setEntityType(value: Tag.EntityTypeMap[keyof Tag.EntityTypeMap]): void;

  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  getEnvironmentName(): string;
  setEnvironmentName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tag.AsObject;
  static toObject(includeInstance: boolean, msg: Tag): Tag.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(message: Tag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tag;
  static deserializeBinaryFromReader(
    message: Tag,
    reader: jspb.BinaryReader
  ): Tag;
}

export namespace Tag {
  export type AsObject = {
    id: string;
    name: string;
    createdAt: number;
    updatedAt: number;
    entityType: Tag.EntityTypeMap[keyof Tag.EntityTypeMap];
    environmentId: string;
    environmentName: string;
  };

  export interface EntityTypeMap {
    UNSPECIFIED: 0;
    FEATURE_FLAG: 1;
    ACCOUNT: 2;
  }

  export const EntityType: EntityTypeMap;
}

export class EnvironmentTag extends jspb.Message {
  getEnvironmentId(): string;
  setEnvironmentId(value: string): void;

  clearTagsList(): void;
  getTagsList(): Array<Tag>;
  setTagsList(value: Array<Tag>): void;
  addTags(value?: Tag, index?: number): Tag;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvironmentTag.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: EnvironmentTag
  ): EnvironmentTag.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: EnvironmentTag,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): EnvironmentTag;
  static deserializeBinaryFromReader(
    message: EnvironmentTag,
    reader: jspb.BinaryReader
  ): EnvironmentTag;
}

export namespace EnvironmentTag {
  export type AsObject = {
    environmentId: string;
    tagsList: Array<Tag.AsObject>;
  };
}
