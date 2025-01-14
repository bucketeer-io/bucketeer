// package: bucketeer.tag
// file: proto/tag/service.proto

import * as proto_tag_service_pb from '../../proto/tag/service_pb';
import { grpc } from '@improbable-eng/grpc-web';

type TagServiceListTags = {
  readonly methodName: string;
  readonly service: typeof TagService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_tag_service_pb.ListTagsRequest;
  readonly responseType: typeof proto_tag_service_pb.ListTagsResponse;
};

type TagServiceCreateTag = {
  readonly methodName: string;
  readonly service: typeof TagService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_tag_service_pb.CreateTagRequest;
  readonly responseType: typeof proto_tag_service_pb.CreateTagResponse;
};

type TagServiceDeleteTag = {
  readonly methodName: string;
  readonly service: typeof TagService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_tag_service_pb.DeleteTagRequest;
  readonly responseType: typeof proto_tag_service_pb.DeleteTagResponse;
};

export class TagService {
  static readonly serviceName: string;
  static readonly ListTags: TagServiceListTags;
  static readonly CreateTag: TagServiceCreateTag;
  static readonly DeleteTag: TagServiceDeleteTag;
}

export type ServiceError = {
  message: string;
  code: number;
  metadata: grpc.Metadata;
};
export type Status = { details: string; code: number; metadata: grpc.Metadata };

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(
    type: 'data',
    handler: (message: ResT) => void
  ): BidirectionalStream<ReqT, ResT>;
  on(
    type: 'end',
    handler: (status?: Status) => void
  ): BidirectionalStream<ReqT, ResT>;
  on(
    type: 'status',
    handler: (status: Status) => void
  ): BidirectionalStream<ReqT, ResT>;
}

export class TagServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listTags(
    requestMessage: proto_tag_service_pb.ListTagsRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_tag_service_pb.ListTagsResponse | null
    ) => void
  ): UnaryResponse;
  listTags(
    requestMessage: proto_tag_service_pb.ListTagsRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_tag_service_pb.ListTagsResponse | null
    ) => void
  ): UnaryResponse;
  createTag(
    requestMessage: proto_tag_service_pb.CreateTagRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_tag_service_pb.CreateTagResponse | null
    ) => void
  ): UnaryResponse;
  createTag(
    requestMessage: proto_tag_service_pb.CreateTagRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_tag_service_pb.CreateTagResponse | null
    ) => void
  ): UnaryResponse;
  deleteTag(
    requestMessage: proto_tag_service_pb.DeleteTagRequest,
    metadata: grpc.Metadata,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_tag_service_pb.DeleteTagResponse | null
    ) => void
  ): UnaryResponse;
  deleteTag(
    requestMessage: proto_tag_service_pb.DeleteTagRequest,
    callback: (
      error: ServiceError | null,
      responseMessage: proto_tag_service_pb.DeleteTagResponse | null
    ) => void
  ): UnaryResponse;
}
