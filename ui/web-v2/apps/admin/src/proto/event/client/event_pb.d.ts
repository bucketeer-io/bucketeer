// package: bucketeer.event.client
// file: proto/event/client/event.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as proto_feature_evaluation_pb from "../../../proto/feature/evaluation_pb";
import * as proto_feature_reason_pb from "../../../proto/feature/reason_pb";
import * as proto_user_user_pb from "../../../proto/user/user_pb";

export class Event extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  hasEvent(): boolean;
  clearEvent(): void;
  getEvent(): google_protobuf_any_pb.Any | undefined;
  setEvent(value?: google_protobuf_any_pb.Any): void;

  getEnvironmentNamespace(): string;
  setEnvironmentNamespace(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    id: string,
    event?: google_protobuf_any_pb.Any.AsObject,
    environmentNamespace: string,
  }
}

export class EvaluationEvent extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getUserId(): string;
  setUserId(value: string): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  hasUser(): boolean;
  clearUser(): void;
  getUser(): proto_user_user_pb.User | undefined;
  setUser(value?: proto_user_user_pb.User): void;

  hasReason(): boolean;
  clearReason(): void;
  getReason(): proto_feature_reason_pb.Reason | undefined;
  setReason(value?: proto_feature_reason_pb.Reason): void;

  getTag(): string;
  setTag(value: string): void;

  getSourceId(): SourceIdMap[keyof SourceIdMap];
  setSourceId(value: SourceIdMap[keyof SourceIdMap]): void;

  getSdkVersion(): string;
  setSdkVersion(value: string): void;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EvaluationEvent.AsObject;
  static toObject(includeInstance: boolean, msg: EvaluationEvent): EvaluationEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EvaluationEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EvaluationEvent;
  static deserializeBinaryFromReader(message: EvaluationEvent, reader: jspb.BinaryReader): EvaluationEvent;
}

export namespace EvaluationEvent {
  export type AsObject = {
    timestamp: number,
    featureId: string,
    featureVersion: number,
    userId: string,
    variationId: string,
    user?: proto_user_user_pb.User.AsObject,
    reason?: proto_feature_reason_pb.Reason.AsObject,
    tag: string,
    sourceId: SourceIdMap[keyof SourceIdMap],
    sdkVersion: string,
    metadataMap: Array<[string, string]>,
  }
}

export class GoalEvent extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getValue(): number;
  setValue(value: number): void;

  hasUser(): boolean;
  clearUser(): void;
  getUser(): proto_user_user_pb.User | undefined;
  setUser(value?: proto_user_user_pb.User): void;

  clearEvaluationsList(): void;
  getEvaluationsList(): Array<proto_feature_evaluation_pb.Evaluation>;
  setEvaluationsList(value: Array<proto_feature_evaluation_pb.Evaluation>): void;
  addEvaluations(value?: proto_feature_evaluation_pb.Evaluation, index?: number): proto_feature_evaluation_pb.Evaluation;

  getTag(): string;
  setTag(value: string): void;

  getSourceId(): SourceIdMap[keyof SourceIdMap];
  setSourceId(value: SourceIdMap[keyof SourceIdMap]): void;

  getSdkVersion(): string;
  setSdkVersion(value: string): void;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoalEvent.AsObject;
  static toObject(includeInstance: boolean, msg: GoalEvent): GoalEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GoalEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GoalEvent;
  static deserializeBinaryFromReader(message: GoalEvent, reader: jspb.BinaryReader): GoalEvent;
}

export namespace GoalEvent {
  export type AsObject = {
    timestamp: number,
    goalId: string,
    userId: string,
    value: number,
    user?: proto_user_user_pb.User.AsObject,
    evaluationsList: Array<proto_feature_evaluation_pb.Evaluation.AsObject>,
    tag: string,
    sourceId: SourceIdMap[keyof SourceIdMap],
    sdkVersion: string,
    metadataMap: Array<[string, string]>,
  }
}

export class MetricsEvent extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  hasEvent(): boolean;
  clearEvent(): void;
  getEvent(): google_protobuf_any_pb.Any | undefined;
  setEvent(value?: google_protobuf_any_pb.Any): void;

  getSourceId(): SourceIdMap[keyof SourceIdMap];
  setSourceId(value: SourceIdMap[keyof SourceIdMap]): void;

  getSdkVersion(): string;
  setSdkVersion(value: string): void;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: MetricsEvent): MetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetricsEvent;
  static deserializeBinaryFromReader(message: MetricsEvent, reader: jspb.BinaryReader): MetricsEvent;
}

export namespace MetricsEvent {
  export type AsObject = {
    timestamp: number,
    event?: google_protobuf_any_pb.Any.AsObject,
    sourceId: SourceIdMap[keyof SourceIdMap],
    sdkVersion: string,
    metadataMap: Array<[string, string]>,
  }
}

export class GetEvaluationLatencyMetricsEvent extends jspb.Message {
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  hasDuration(): boolean;
  clearDuration(): void;
  getDuration(): google_protobuf_duration_pb.Duration | undefined;
  setDuration(value?: google_protobuf_duration_pb.Duration): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEvaluationLatencyMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: GetEvaluationLatencyMetricsEvent): GetEvaluationLatencyMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEvaluationLatencyMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEvaluationLatencyMetricsEvent;
  static deserializeBinaryFromReader(message: GetEvaluationLatencyMetricsEvent, reader: jspb.BinaryReader): GetEvaluationLatencyMetricsEvent;
}

export namespace GetEvaluationLatencyMetricsEvent {
  export type AsObject = {
    labelsMap: Array<[string, string]>,
    duration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class GetEvaluationSizeMetricsEvent extends jspb.Message {
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getSizeByte(): number;
  setSizeByte(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetEvaluationSizeMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: GetEvaluationSizeMetricsEvent): GetEvaluationSizeMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetEvaluationSizeMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetEvaluationSizeMetricsEvent;
  static deserializeBinaryFromReader(message: GetEvaluationSizeMetricsEvent, reader: jspb.BinaryReader): GetEvaluationSizeMetricsEvent;
}

export namespace GetEvaluationSizeMetricsEvent {
  export type AsObject = {
    labelsMap: Array<[string, string]>,
    sizeByte: number,
  }
}

export class LatencyMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  hasDuration(): boolean;
  clearDuration(): void;
  getDuration(): google_protobuf_duration_pb.Duration | undefined;
  setDuration(value?: google_protobuf_duration_pb.Duration): void;

  getLatencySecond(): number;
  setLatencySecond(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LatencyMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: LatencyMetricsEvent): LatencyMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LatencyMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LatencyMetricsEvent;
  static deserializeBinaryFromReader(message: LatencyMetricsEvent, reader: jspb.BinaryReader): LatencyMetricsEvent;
}

export namespace LatencyMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
    duration?: google_protobuf_duration_pb.Duration.AsObject,
    latencySecond: number,
  }
}

export class SizeMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  getSizeByte(): number;
  setSizeByte(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SizeMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: SizeMetricsEvent): SizeMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SizeMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SizeMetricsEvent;
  static deserializeBinaryFromReader(message: SizeMetricsEvent, reader: jspb.BinaryReader): SizeMetricsEvent;
}

export namespace SizeMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
    sizeByte: number,
  }
}

export class TimeoutErrorCountMetricsEvent extends jspb.Message {
  getTag(): string;
  setTag(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TimeoutErrorCountMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: TimeoutErrorCountMetricsEvent): TimeoutErrorCountMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TimeoutErrorCountMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TimeoutErrorCountMetricsEvent;
  static deserializeBinaryFromReader(message: TimeoutErrorCountMetricsEvent, reader: jspb.BinaryReader): TimeoutErrorCountMetricsEvent;
}

export namespace TimeoutErrorCountMetricsEvent {
  export type AsObject = {
    tag: string,
  }
}

export class InternalErrorCountMetricsEvent extends jspb.Message {
  getTag(): string;
  setTag(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InternalErrorCountMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: InternalErrorCountMetricsEvent): InternalErrorCountMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InternalErrorCountMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InternalErrorCountMetricsEvent;
  static deserializeBinaryFromReader(message: InternalErrorCountMetricsEvent, reader: jspb.BinaryReader): InternalErrorCountMetricsEvent;
}

export namespace InternalErrorCountMetricsEvent {
  export type AsObject = {
    tag: string,
  }
}

export class RedirectionRequestExceptionEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RedirectionRequestExceptionEvent.AsObject;
  static toObject(includeInstance: boolean, msg: RedirectionRequestExceptionEvent): RedirectionRequestExceptionEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RedirectionRequestExceptionEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RedirectionRequestExceptionEvent;
  static deserializeBinaryFromReader(message: RedirectionRequestExceptionEvent, reader: jspb.BinaryReader): RedirectionRequestExceptionEvent;
}

export namespace RedirectionRequestExceptionEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class BadRequestErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BadRequestErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: BadRequestErrorMetricsEvent): BadRequestErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: BadRequestErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BadRequestErrorMetricsEvent;
  static deserializeBinaryFromReader(message: BadRequestErrorMetricsEvent, reader: jspb.BinaryReader): BadRequestErrorMetricsEvent;
}

export namespace BadRequestErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class UnauthorizedErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnauthorizedErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: UnauthorizedErrorMetricsEvent): UnauthorizedErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnauthorizedErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnauthorizedErrorMetricsEvent;
  static deserializeBinaryFromReader(message: UnauthorizedErrorMetricsEvent, reader: jspb.BinaryReader): UnauthorizedErrorMetricsEvent;
}

export namespace UnauthorizedErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class ForbiddenErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ForbiddenErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: ForbiddenErrorMetricsEvent): ForbiddenErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ForbiddenErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ForbiddenErrorMetricsEvent;
  static deserializeBinaryFromReader(message: ForbiddenErrorMetricsEvent, reader: jspb.BinaryReader): ForbiddenErrorMetricsEvent;
}

export namespace ForbiddenErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class NotFoundErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotFoundErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: NotFoundErrorMetricsEvent): NotFoundErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NotFoundErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotFoundErrorMetricsEvent;
  static deserializeBinaryFromReader(message: NotFoundErrorMetricsEvent, reader: jspb.BinaryReader): NotFoundErrorMetricsEvent;
}

export namespace NotFoundErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class PayloadTooLargeExceptionEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PayloadTooLargeExceptionEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PayloadTooLargeExceptionEvent): PayloadTooLargeExceptionEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PayloadTooLargeExceptionEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PayloadTooLargeExceptionEvent;
  static deserializeBinaryFromReader(message: PayloadTooLargeExceptionEvent, reader: jspb.BinaryReader): PayloadTooLargeExceptionEvent;
}

export namespace PayloadTooLargeExceptionEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class ClientClosedRequestErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClientClosedRequestErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: ClientClosedRequestErrorMetricsEvent): ClientClosedRequestErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClientClosedRequestErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClientClosedRequestErrorMetricsEvent;
  static deserializeBinaryFromReader(message: ClientClosedRequestErrorMetricsEvent, reader: jspb.BinaryReader): ClientClosedRequestErrorMetricsEvent;
}

export namespace ClientClosedRequestErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class InternalServerErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InternalServerErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: InternalServerErrorMetricsEvent): InternalServerErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InternalServerErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InternalServerErrorMetricsEvent;
  static deserializeBinaryFromReader(message: InternalServerErrorMetricsEvent, reader: jspb.BinaryReader): InternalServerErrorMetricsEvent;
}

export namespace InternalServerErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class ServiceUnavailableErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceUnavailableErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceUnavailableErrorMetricsEvent): ServiceUnavailableErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceUnavailableErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceUnavailableErrorMetricsEvent;
  static deserializeBinaryFromReader(message: ServiceUnavailableErrorMetricsEvent, reader: jspb.BinaryReader): ServiceUnavailableErrorMetricsEvent;
}

export namespace ServiceUnavailableErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class TimeoutErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TimeoutErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: TimeoutErrorMetricsEvent): TimeoutErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TimeoutErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TimeoutErrorMetricsEvent;
  static deserializeBinaryFromReader(message: TimeoutErrorMetricsEvent, reader: jspb.BinaryReader): TimeoutErrorMetricsEvent;
}

export namespace TimeoutErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class InternalErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InternalErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: InternalErrorMetricsEvent): InternalErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InternalErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InternalErrorMetricsEvent;
  static deserializeBinaryFromReader(message: InternalErrorMetricsEvent, reader: jspb.BinaryReader): InternalErrorMetricsEvent;
}

export namespace InternalErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class NetworkErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NetworkErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: NetworkErrorMetricsEvent): NetworkErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NetworkErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NetworkErrorMetricsEvent;
  static deserializeBinaryFromReader(message: NetworkErrorMetricsEvent, reader: jspb.BinaryReader): NetworkErrorMetricsEvent;
}

export namespace NetworkErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class InternalSdkErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InternalSdkErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: InternalSdkErrorMetricsEvent): InternalSdkErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InternalSdkErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InternalSdkErrorMetricsEvent;
  static deserializeBinaryFromReader(message: InternalSdkErrorMetricsEvent, reader: jspb.BinaryReader): InternalSdkErrorMetricsEvent;
}

export namespace InternalSdkErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class UnknownErrorMetricsEvent extends jspb.Message {
  getApiId(): ApiIdMap[keyof ApiIdMap];
  setApiId(value: ApiIdMap[keyof ApiIdMap]): void;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnknownErrorMetricsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: UnknownErrorMetricsEvent): UnknownErrorMetricsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UnknownErrorMetricsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnknownErrorMetricsEvent;
  static deserializeBinaryFromReader(message: UnknownErrorMetricsEvent, reader: jspb.BinaryReader): UnknownErrorMetricsEvent;
}

export namespace UnknownErrorMetricsEvent {
  export type AsObject = {
    apiId: ApiIdMap[keyof ApiIdMap],
    labelsMap: Array<[string, string]>,
  }
}

export class OpsEvent extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  getFeatureId(): string;
  setFeatureId(value: string): void;

  getFeatureVersion(): number;
  setFeatureVersion(value: number): void;

  getVariationId(): string;
  setVariationId(value: string): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OpsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: OpsEvent): OpsEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OpsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OpsEvent;
  static deserializeBinaryFromReader(message: OpsEvent, reader: jspb.BinaryReader): OpsEvent;
}

export namespace OpsEvent {
  export type AsObject = {
    timestamp: number,
    featureId: string,
    featureVersion: number,
    variationId: string,
    goalId: string,
    userId: string,
  }
}

export class UserGoalEvent extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  getGoalId(): string;
  setGoalId(value: string): void;

  getValue(): number;
  setValue(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGoalEvent.AsObject;
  static toObject(includeInstance: boolean, msg: UserGoalEvent): UserGoalEvent.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserGoalEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGoalEvent;
  static deserializeBinaryFromReader(message: UserGoalEvent, reader: jspb.BinaryReader): UserGoalEvent;
}

export namespace UserGoalEvent {
  export type AsObject = {
    timestamp: number,
    goalId: string,
    value: number,
  }
}

export interface SourceIdMap {
  UNKNOWN: 0;
  ANDROID: 1;
  IOS: 2;
  WEB: 3;
  GO_SERVER: 5;
  NODE_SERVER: 6;
  JAVASCRIPT: 7;
}

export const SourceId: SourceIdMap;

export interface ApiIdMap {
  UNKNOWN_API: 0;
  GET_EVALUATION: 1;
  GET_EVALUATIONS: 2;
  REGISTER_EVENTS: 3;
}

export const ApiId: ApiIdMap;

