// Copyright 2023 The Bucketeer Authors.
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
// source: proto/auth/service.proto

package auth

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

type GetAuthCodeURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State       string `protobuf:"bytes,1,opt,name=state,proto3" json:"state"`
	RedirectUrl string `protobuf:"bytes,2,opt,name=redirect_url,json=redirectUrl,proto3" json:"redirect_url"`
}

func (x *GetAuthCodeURLRequest) Reset() {
	*x = GetAuthCodeURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthCodeURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthCodeURLRequest) ProtoMessage() {}

func (x *GetAuthCodeURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthCodeURLRequest.ProtoReflect.Descriptor instead.
func (*GetAuthCodeURLRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAuthCodeURLRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *GetAuthCodeURLRequest) GetRedirectUrl() string {
	if x != nil {
		return x.RedirectUrl
	}
	return ""
}

type GetAuthCodeURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url"`
}

func (x *GetAuthCodeURLResponse) Reset() {
	*x = GetAuthCodeURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthCodeURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthCodeURLResponse) ProtoMessage() {}

func (x *GetAuthCodeURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthCodeURLResponse.ProtoReflect.Descriptor instead.
func (*GetAuthCodeURLResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetAuthCodeURLResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ExchangeTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code        string `protobuf:"bytes,1,opt,name=code,proto3" json:"code"`
	RedirectUrl string `protobuf:"bytes,2,opt,name=redirect_url,json=redirectUrl,proto3" json:"redirect_url"`
}

func (x *ExchangeTokenRequest) Reset() {
	*x = ExchangeTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExchangeTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExchangeTokenRequest) ProtoMessage() {}

func (x *ExchangeTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExchangeTokenRequest.ProtoReflect.Descriptor instead.
func (*ExchangeTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_service_proto_rawDescGZIP(), []int{2}
}

func (x *ExchangeTokenRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *ExchangeTokenRequest) GetRedirectUrl() string {
	if x != nil {
		return x.RedirectUrl
	}
	return ""
}

type ExchangeTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
}

func (x *ExchangeTokenResponse) Reset() {
	*x = ExchangeTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExchangeTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExchangeTokenResponse) ProtoMessage() {}

func (x *ExchangeTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExchangeTokenResponse.ProtoReflect.Descriptor instead.
func (*ExchangeTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_service_proto_rawDescGZIP(), []int{3}
}

func (x *ExchangeTokenResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type RefreshTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefreshToken string `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token"`
	RedirectUrl  string `protobuf:"bytes,2,opt,name=redirect_url,json=redirectUrl,proto3" json:"redirect_url"`
}

func (x *RefreshTokenRequest) Reset() {
	*x = RefreshTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenRequest) ProtoMessage() {}

func (x *RefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*RefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_service_proto_rawDescGZIP(), []int{4}
}

func (x *RefreshTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *RefreshTokenRequest) GetRedirectUrl() string {
	if x != nil {
		return x.RedirectUrl
	}
	return ""
}

type RefreshTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
}

func (x *RefreshTokenResponse) Reset() {
	*x = RefreshTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_auth_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenResponse) ProtoMessage() {}

func (x *RefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*RefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_service_proto_rawDescGZIP(), []int{5}
}

func (x *RefreshTokenResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

var File_proto_auth_service_proto protoreflect.FileDescriptor

var file_proto_auth_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x1a, 0x16, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x50, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64,
	0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x55, 0x72, 0x6c, 0x22, 0x2a, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x43,
	0x6f, 0x64, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x22, 0x4d, 0x0a, 0x14, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x6c, 0x22,
	0x44, 0x0a, 0x15, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x5d, 0x0a, 0x13, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x55, 0x72, 0x6c, 0x22, 0x43, 0x0a, 0x14, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xa7, 0x02, 0x0a, 0x0b, 0x41, 0x75,
	0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x25, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5c, 0x0a, 0x0d, 0x45, 0x78,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x24, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x25, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x23, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_auth_service_proto_rawDescOnce sync.Once
	file_proto_auth_service_proto_rawDescData = file_proto_auth_service_proto_rawDesc
)

func file_proto_auth_service_proto_rawDescGZIP() []byte {
	file_proto_auth_service_proto_rawDescOnce.Do(func() {
		file_proto_auth_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_auth_service_proto_rawDescData)
	})
	return file_proto_auth_service_proto_rawDescData
}

var file_proto_auth_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_auth_service_proto_goTypes = []interface{}{
	(*GetAuthCodeURLRequest)(nil),  // 0: bucketeer.auth.GetAuthCodeURLRequest
	(*GetAuthCodeURLResponse)(nil), // 1: bucketeer.auth.GetAuthCodeURLResponse
	(*ExchangeTokenRequest)(nil),   // 2: bucketeer.auth.ExchangeTokenRequest
	(*ExchangeTokenResponse)(nil),  // 3: bucketeer.auth.ExchangeTokenResponse
	(*RefreshTokenRequest)(nil),    // 4: bucketeer.auth.RefreshTokenRequest
	(*RefreshTokenResponse)(nil),   // 5: bucketeer.auth.RefreshTokenResponse
	(*Token)(nil),                  // 6: bucketeer.auth.Token
}
var file_proto_auth_service_proto_depIdxs = []int32{
	6, // 0: bucketeer.auth.ExchangeTokenResponse.token:type_name -> bucketeer.auth.Token
	6, // 1: bucketeer.auth.RefreshTokenResponse.token:type_name -> bucketeer.auth.Token
	0, // 2: bucketeer.auth.AuthService.GetAuthCodeURL:input_type -> bucketeer.auth.GetAuthCodeURLRequest
	2, // 3: bucketeer.auth.AuthService.ExchangeToken:input_type -> bucketeer.auth.ExchangeTokenRequest
	4, // 4: bucketeer.auth.AuthService.RefreshToken:input_type -> bucketeer.auth.RefreshTokenRequest
	1, // 5: bucketeer.auth.AuthService.GetAuthCodeURL:output_type -> bucketeer.auth.GetAuthCodeURLResponse
	3, // 6: bucketeer.auth.AuthService.ExchangeToken:output_type -> bucketeer.auth.ExchangeTokenResponse
	5, // 7: bucketeer.auth.AuthService.RefreshToken:output_type -> bucketeer.auth.RefreshTokenResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_auth_service_proto_init() }
func file_proto_auth_service_proto_init() {
	if File_proto_auth_service_proto != nil {
		return
	}
	file_proto_auth_token_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_auth_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthCodeURLRequest); i {
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
		file_proto_auth_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthCodeURLResponse); i {
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
		file_proto_auth_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExchangeTokenRequest); i {
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
		file_proto_auth_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExchangeTokenResponse); i {
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
		file_proto_auth_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshTokenRequest); i {
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
		file_proto_auth_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshTokenResponse); i {
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
			RawDescriptor: file_proto_auth_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_auth_service_proto_goTypes,
		DependencyIndexes: file_proto_auth_service_proto_depIdxs,
		MessageInfos:      file_proto_auth_service_proto_msgTypes,
	}.Build()
	File_proto_auth_service_proto = out.File
	file_proto_auth_service_proto_rawDesc = nil
	file_proto_auth_service_proto_goTypes = nil
	file_proto_auth_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthServiceClient interface {
	GetAuthCodeURL(ctx context.Context, in *GetAuthCodeURLRequest, opts ...grpc.CallOption) (*GetAuthCodeURLResponse, error)
	ExchangeToken(ctx context.Context, in *ExchangeTokenRequest, opts ...grpc.CallOption) (*ExchangeTokenResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) GetAuthCodeURL(ctx context.Context, in *GetAuthCodeURLRequest, opts ...grpc.CallOption) (*GetAuthCodeURLResponse, error) {
	out := new(GetAuthCodeURLResponse)
	err := c.cc.Invoke(ctx, "/bucketeer.auth.AuthService/GetAuthCodeURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ExchangeToken(ctx context.Context, in *ExchangeTokenRequest, opts ...grpc.CallOption) (*ExchangeTokenResponse, error) {
	out := new(ExchangeTokenResponse)
	err := c.cc.Invoke(ctx, "/bucketeer.auth.AuthService/ExchangeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error) {
	out := new(RefreshTokenResponse)
	err := c.cc.Invoke(ctx, "/bucketeer.auth.AuthService/RefreshToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
type AuthServiceServer interface {
	GetAuthCodeURL(context.Context, *GetAuthCodeURLRequest) (*GetAuthCodeURLResponse, error)
	ExchangeToken(context.Context, *ExchangeTokenRequest) (*ExchangeTokenResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error)
}

// UnimplementedAuthServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (*UnimplementedAuthServiceServer) GetAuthCodeURL(context.Context, *GetAuthCodeURLRequest) (*GetAuthCodeURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthCodeURL not implemented")
}
func (*UnimplementedAuthServiceServer) ExchangeToken(context.Context, *ExchangeTokenRequest) (*ExchangeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExchangeToken not implemented")
}
func (*UnimplementedAuthServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
}

func _AuthService_GetAuthCodeURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthCodeURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GetAuthCodeURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bucketeer.auth.AuthService/GetAuthCodeURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetAuthCodeURL(ctx, req.(*GetAuthCodeURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ExchangeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExchangeTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ExchangeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bucketeer.auth.AuthService/ExchangeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ExchangeToken(ctx, req.(*ExchangeTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bucketeer.auth.AuthService/RefreshToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bucketeer.auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthCodeURL",
			Handler:    _AuthService_GetAuthCodeURL_Handler,
		},
		{
			MethodName: "ExchangeToken",
			Handler:    _AuthService_ExchangeToken_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _AuthService_RefreshToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/auth/service.proto",
}
