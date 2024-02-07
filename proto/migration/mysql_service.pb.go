// Copyright 2024 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.23.4
// source: proto/migration/mysql_service.proto

package migration

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MigrateAllMasterSchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MigrateAllMasterSchemaRequest) Reset() {
	*x = MigrateAllMasterSchemaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_migration_mysql_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MigrateAllMasterSchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MigrateAllMasterSchemaRequest) ProtoMessage() {}

func (x *MigrateAllMasterSchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_migration_mysql_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MigrateAllMasterSchemaRequest.ProtoReflect.Descriptor instead.
func (*MigrateAllMasterSchemaRequest) Descriptor() ([]byte, []int) {
	return file_proto_migration_mysql_service_proto_rawDescGZIP(), []int{0}
}

type MigrateAllMasterSchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MigrateAllMasterSchemaResponse) Reset() {
	*x = MigrateAllMasterSchemaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_migration_mysql_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MigrateAllMasterSchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MigrateAllMasterSchemaResponse) ProtoMessage() {}

func (x *MigrateAllMasterSchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_migration_mysql_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MigrateAllMasterSchemaResponse.ProtoReflect.Descriptor instead.
func (*MigrateAllMasterSchemaResponse) Descriptor() ([]byte, []int) {
	return file_proto_migration_mysql_service_proto_rawDescGZIP(), []int{1}
}

type RollbackMasterSchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Step int64 `protobuf:"varint,1,opt,name=step,proto3" json:"step"`
}

func (x *RollbackMasterSchemaRequest) Reset() {
	*x = RollbackMasterSchemaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_migration_mysql_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RollbackMasterSchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RollbackMasterSchemaRequest) ProtoMessage() {}

func (x *RollbackMasterSchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_migration_mysql_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RollbackMasterSchemaRequest.ProtoReflect.Descriptor instead.
func (*RollbackMasterSchemaRequest) Descriptor() ([]byte, []int) {
	return file_proto_migration_mysql_service_proto_rawDescGZIP(), []int{2}
}

func (x *RollbackMasterSchemaRequest) GetStep() int64 {
	if x != nil {
		return x.Step
	}
	return 0
}

type RollbackMasterSchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RollbackMasterSchemaResponse) Reset() {
	*x = RollbackMasterSchemaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_migration_mysql_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RollbackMasterSchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RollbackMasterSchemaResponse) ProtoMessage() {}

func (x *RollbackMasterSchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_migration_mysql_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RollbackMasterSchemaResponse.ProtoReflect.Descriptor instead.
func (*RollbackMasterSchemaResponse) Descriptor() ([]byte, []int) {
	return file_proto_migration_mysql_service_proto_rawDescGZIP(), []int{3}
}

var File_proto_migration_mysql_service_proto protoreflect.FileDescriptor

var file_proto_migration_mysql_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x1f, 0x0a, 0x1d, 0x4d, 0x69,
	0x67, 0x72, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x20, 0x0a, 0x1e, 0x4d,
	0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x0a,
	0x1b, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x74, 0x65, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x74, 0x65, 0x70,
	0x22, 0x1e, 0x0a, 0x1c, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x9c, 0x02, 0x0a, 0x15, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x79,
	0x53, 0x51, 0x4c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x16, 0x4d,
	0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x32, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x69, 0x67, 0x72,
	0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x7d, 0x0a, 0x14, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x30, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2e, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52,
	0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x69, 0x67, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_migration_mysql_service_proto_rawDescOnce sync.Once
	file_proto_migration_mysql_service_proto_rawDescData = file_proto_migration_mysql_service_proto_rawDesc
)

func file_proto_migration_mysql_service_proto_rawDescGZIP() []byte {
	file_proto_migration_mysql_service_proto_rawDescOnce.Do(func() {
		file_proto_migration_mysql_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_migration_mysql_service_proto_rawDescData)
	})
	return file_proto_migration_mysql_service_proto_rawDescData
}

var file_proto_migration_mysql_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_migration_mysql_service_proto_goTypes = []interface{}{
	(*MigrateAllMasterSchemaRequest)(nil),  // 0: bucketeer.migration.MigrateAllMasterSchemaRequest
	(*MigrateAllMasterSchemaResponse)(nil), // 1: bucketeer.migration.MigrateAllMasterSchemaResponse
	(*RollbackMasterSchemaRequest)(nil),    // 2: bucketeer.migration.RollbackMasterSchemaRequest
	(*RollbackMasterSchemaResponse)(nil),   // 3: bucketeer.migration.RollbackMasterSchemaResponse
}
var file_proto_migration_mysql_service_proto_depIdxs = []int32{
	0, // 0: bucketeer.migration.MigrationMySQLService.MigrateAllMasterSchema:input_type -> bucketeer.migration.MigrateAllMasterSchemaRequest
	2, // 1: bucketeer.migration.MigrationMySQLService.RollbackMasterSchema:input_type -> bucketeer.migration.RollbackMasterSchemaRequest
	1, // 2: bucketeer.migration.MigrationMySQLService.MigrateAllMasterSchema:output_type -> bucketeer.migration.MigrateAllMasterSchemaResponse
	3, // 3: bucketeer.migration.MigrationMySQLService.RollbackMasterSchema:output_type -> bucketeer.migration.RollbackMasterSchemaResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_migration_mysql_service_proto_init() }
func file_proto_migration_mysql_service_proto_init() {
	if File_proto_migration_mysql_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_migration_mysql_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MigrateAllMasterSchemaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_migration_mysql_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MigrateAllMasterSchemaResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_migration_mysql_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RollbackMasterSchemaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_migration_mysql_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RollbackMasterSchemaResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_migration_mysql_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_migration_mysql_service_proto_goTypes,
		DependencyIndexes: file_proto_migration_mysql_service_proto_depIdxs,
		MessageInfos:      file_proto_migration_mysql_service_proto_msgTypes,
	}.Build()
	File_proto_migration_mysql_service_proto = out.File
	file_proto_migration_mysql_service_proto_rawDesc = nil
	file_proto_migration_mysql_service_proto_goTypes = nil
	file_proto_migration_mysql_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MigrationMySQLServiceClient is the client API for MigrationMySQLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MigrationMySQLServiceClient interface {
	MigrateAllMasterSchema(ctx context.Context, in *MigrateAllMasterSchemaRequest, opts ...grpc.CallOption) (*MigrateAllMasterSchemaResponse, error)
	RollbackMasterSchema(ctx context.Context, in *RollbackMasterSchemaRequest, opts ...grpc.CallOption) (*RollbackMasterSchemaResponse, error)
}

type migrationMySQLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMigrationMySQLServiceClient(cc grpc.ClientConnInterface) MigrationMySQLServiceClient {
	return &migrationMySQLServiceClient{cc}
}

func (c *migrationMySQLServiceClient) MigrateAllMasterSchema(ctx context.Context, in *MigrateAllMasterSchemaRequest, opts ...grpc.CallOption) (*MigrateAllMasterSchemaResponse, error) {
	out := new(MigrateAllMasterSchemaResponse)
	err := c.cc.Invoke(ctx, "/bucketeer.migration.MigrationMySQLService/MigrateAllMasterSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *migrationMySQLServiceClient) RollbackMasterSchema(ctx context.Context, in *RollbackMasterSchemaRequest, opts ...grpc.CallOption) (*RollbackMasterSchemaResponse, error) {
	out := new(RollbackMasterSchemaResponse)
	err := c.cc.Invoke(ctx, "/bucketeer.migration.MigrationMySQLService/RollbackMasterSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MigrationMySQLServiceServer is the server API for MigrationMySQLService service.
type MigrationMySQLServiceServer interface {
	MigrateAllMasterSchema(context.Context, *MigrateAllMasterSchemaRequest) (*MigrateAllMasterSchemaResponse, error)
	RollbackMasterSchema(context.Context, *RollbackMasterSchemaRequest) (*RollbackMasterSchemaResponse, error)
}

// UnimplementedMigrationMySQLServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMigrationMySQLServiceServer struct {
}

func (*UnimplementedMigrationMySQLServiceServer) MigrateAllMasterSchema(context.Context, *MigrateAllMasterSchemaRequest) (*MigrateAllMasterSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MigrateAllMasterSchema not implemented")
}
func (*UnimplementedMigrationMySQLServiceServer) RollbackMasterSchema(context.Context, *RollbackMasterSchemaRequest) (*RollbackMasterSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RollbackMasterSchema not implemented")
}

func RegisterMigrationMySQLServiceServer(s *grpc.Server, srv MigrationMySQLServiceServer) {
	s.RegisterService(&_MigrationMySQLService_serviceDesc, srv)
}

func _MigrationMySQLService_MigrateAllMasterSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MigrateAllMasterSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MigrationMySQLServiceServer).MigrateAllMasterSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bucketeer.migration.MigrationMySQLService/MigrateAllMasterSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MigrationMySQLServiceServer).MigrateAllMasterSchema(ctx, req.(*MigrateAllMasterSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MigrationMySQLService_RollbackMasterSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RollbackMasterSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MigrationMySQLServiceServer).RollbackMasterSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bucketeer.migration.MigrationMySQLService/RollbackMasterSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MigrationMySQLServiceServer).RollbackMasterSchema(ctx, req.(*RollbackMasterSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MigrationMySQLService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bucketeer.migration.MigrationMySQLService",
	HandlerType: (*MigrationMySQLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MigrateAllMasterSchema",
			Handler:    _MigrationMySQLService_MigrateAllMasterSchema_Handler,
		},
		{
			MethodName: "RollbackMasterSchema",
			Handler:    _MigrationMySQLService_RollbackMasterSchema_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/migration/mysql_service.proto",
}
