// Copyright 2025 The Bucketeer Authors.
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
// source: proto/coderef/code_reference.proto

package coderef

import (
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

type CodeReference_RepositoryType int32

const (
	CodeReference_REPOSITORY_TYPE_UNSPECIFIED CodeReference_RepositoryType = 0
	CodeReference_GITHUB                      CodeReference_RepositoryType = 1
	CodeReference_GITLAB                      CodeReference_RepositoryType = 2
	CodeReference_BITBUCKET                   CodeReference_RepositoryType = 3
	CodeReference_CUSTOM                      CodeReference_RepositoryType = 4
)

// Enum value maps for CodeReference_RepositoryType.
var (
	CodeReference_RepositoryType_name = map[int32]string{
		0: "REPOSITORY_TYPE_UNSPECIFIED",
		1: "GITHUB",
		2: "GITLAB",
		3: "BITBUCKET",
		4: "CUSTOM",
	}
	CodeReference_RepositoryType_value = map[string]int32{
		"REPOSITORY_TYPE_UNSPECIFIED": 0,
		"GITHUB":                      1,
		"GITLAB":                      2,
		"BITBUCKET":                   3,
		"CUSTOM":                      4,
	}
)

func (x CodeReference_RepositoryType) Enum() *CodeReference_RepositoryType {
	p := new(CodeReference_RepositoryType)
	*p = x
	return p
}

func (x CodeReference_RepositoryType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CodeReference_RepositoryType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_coderef_code_reference_proto_enumTypes[0].Descriptor()
}

func (CodeReference_RepositoryType) Type() protoreflect.EnumType {
	return &file_proto_coderef_code_reference_proto_enumTypes[0]
}

func (x CodeReference_RepositoryType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CodeReference_RepositoryType.Descriptor instead.
func (CodeReference_RepositoryType) EnumDescriptor() ([]byte, []int) {
	return file_proto_coderef_code_reference_proto_rawDescGZIP(), []int{0, 0}
}

type CodeReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string                       `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FeatureId        string                       `protobuf:"bytes,2,opt,name=feature_id,json=featureId,proto3" json:"feature_id"`
	FilePath         string                       `protobuf:"bytes,3,opt,name=file_path,json=filePath,proto3" json:"file_path"`
	LineNumber       int32                        `protobuf:"varint,4,opt,name=line_number,json=lineNumber,proto3" json:"line_number"`
	CodeSnippet      string                       `protobuf:"bytes,5,opt,name=code_snippet,json=codeSnippet,proto3" json:"code_snippet"`
	ContentHash      string                       `protobuf:"bytes,6,opt,name=content_hash,json=contentHash,proto3" json:"content_hash"`
	Aliases          []string                     `protobuf:"bytes,7,rep,name=aliases,proto3" json:"aliases"`
	RepositoryName   string                       `protobuf:"bytes,8,opt,name=repository_name,json=repositoryName,proto3" json:"repository_name"`
	RepositoryOwner  string                       `protobuf:"bytes,9,opt,name=repository_owner,json=repositoryOwner,proto3" json:"repository_owner"`
	RepositoryType   CodeReference_RepositoryType `protobuf:"varint,10,opt,name=repository_type,json=repositoryType,proto3,enum=bucketeer.coderef.CodeReference_RepositoryType" json:"repository_type"`
	RepositoryBranch string                       `protobuf:"bytes,11,opt,name=repository_branch,json=repositoryBranch,proto3" json:"repository_branch"`
	CommitHash       string                       `protobuf:"bytes,12,opt,name=commit_hash,json=commitHash,proto3" json:"commit_hash"`
	EnvironmentId    string                       `protobuf:"bytes,13,opt,name=environment_id,json=environmentId,proto3" json:"environment_id"`
	CreatedAt        int64                        `protobuf:"varint,14,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt        int64                        `protobuf:"varint,15,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	// URL to view the code in the repository
	SourceUrl     string `protobuf:"bytes,16,opt,name=source_url,json=sourceUrl,proto3" json:"source_url"`
	BranchUrl     string `protobuf:"bytes,17,opt,name=branch_url,json=branchUrl,proto3" json:"branch_url"`
	FileExtension string `protobuf:"bytes,18,opt,name=file_extension,json=fileExtension,proto3" json:"file_extension"` // File extension (e.g., go, ts, cpp)
}

func (x *CodeReference) Reset() {
	*x = CodeReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_coderef_code_reference_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CodeReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CodeReference) ProtoMessage() {}

func (x *CodeReference) ProtoReflect() protoreflect.Message {
	mi := &file_proto_coderef_code_reference_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CodeReference.ProtoReflect.Descriptor instead.
func (*CodeReference) Descriptor() ([]byte, []int) {
	return file_proto_coderef_code_reference_proto_rawDescGZIP(), []int{0}
}

func (x *CodeReference) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CodeReference) GetFeatureId() string {
	if x != nil {
		return x.FeatureId
	}
	return ""
}

func (x *CodeReference) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *CodeReference) GetLineNumber() int32 {
	if x != nil {
		return x.LineNumber
	}
	return 0
}

func (x *CodeReference) GetCodeSnippet() string {
	if x != nil {
		return x.CodeSnippet
	}
	return ""
}

func (x *CodeReference) GetContentHash() string {
	if x != nil {
		return x.ContentHash
	}
	return ""
}

func (x *CodeReference) GetAliases() []string {
	if x != nil {
		return x.Aliases
	}
	return nil
}

func (x *CodeReference) GetRepositoryName() string {
	if x != nil {
		return x.RepositoryName
	}
	return ""
}

func (x *CodeReference) GetRepositoryOwner() string {
	if x != nil {
		return x.RepositoryOwner
	}
	return ""
}

func (x *CodeReference) GetRepositoryType() CodeReference_RepositoryType {
	if x != nil {
		return x.RepositoryType
	}
	return CodeReference_REPOSITORY_TYPE_UNSPECIFIED
}

func (x *CodeReference) GetRepositoryBranch() string {
	if x != nil {
		return x.RepositoryBranch
	}
	return ""
}

func (x *CodeReference) GetCommitHash() string {
	if x != nil {
		return x.CommitHash
	}
	return ""
}

func (x *CodeReference) GetEnvironmentId() string {
	if x != nil {
		return x.EnvironmentId
	}
	return ""
}

func (x *CodeReference) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *CodeReference) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *CodeReference) GetSourceUrl() string {
	if x != nil {
		return x.SourceUrl
	}
	return ""
}

func (x *CodeReference) GetBranchUrl() string {
	if x != nil {
		return x.BranchUrl
	}
	return ""
}

func (x *CodeReference) GetFileExtension() string {
	if x != nil {
		return x.FileExtension
	}
	return ""
}

var File_proto_coderef_code_reference_proto protoreflect.FileDescriptor

var file_proto_coderef_code_reference_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x65, 0x66, 0x2f,
	0x63, 0x6f, 0x64, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x63, 0x6f, 0x64, 0x65, 0x72, 0x65, 0x66, 0x22, 0x88, 0x06, 0x0a, 0x0d, 0x43, 0x6f, 0x64, 0x65,
	0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x6c, 0x69, 0x6e, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x73,
	0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f,
	0x64, 0x65, 0x53, 0x6e, 0x69, 0x70, 0x70, 0x65, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x6c, 0x69, 0x61, 0x73, 0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x6c, 0x69, 0x61, 0x73, 0x65, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x29, 0x0a, 0x10, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x58, 0x0a, 0x0f, 0x72, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x63, 0x6f, 0x64, 0x65, 0x72, 0x65, 0x66, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x66, 0x65,
	0x72, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x10, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x42, 0x72, 0x61, 0x6e, 0x63,
	0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x6e, 0x76, 0x69,
	0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x55, 0x72, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66,
	0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x64, 0x0a, 0x0e,
	0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f,
	0x0a, 0x1b, 0x52, 0x45, 0x50, 0x4f, 0x53, 0x49, 0x54, 0x4f, 0x52, 0x59, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x0a, 0x0a, 0x06, 0x47, 0x49, 0x54, 0x48, 0x55, 0x42, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x47,
	0x49, 0x54, 0x4c, 0x41, 0x42, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x42, 0x49, 0x54, 0x42, 0x55,
	0x43, 0x4b, 0x45, 0x54, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x55, 0x53, 0x54, 0x4f, 0x4d,
	0x10, 0x04, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f,
	0x64, 0x65, 0x72, 0x65, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_coderef_code_reference_proto_rawDescOnce sync.Once
	file_proto_coderef_code_reference_proto_rawDescData = file_proto_coderef_code_reference_proto_rawDesc
)

func file_proto_coderef_code_reference_proto_rawDescGZIP() []byte {
	file_proto_coderef_code_reference_proto_rawDescOnce.Do(func() {
		file_proto_coderef_code_reference_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_coderef_code_reference_proto_rawDescData)
	})
	return file_proto_coderef_code_reference_proto_rawDescData
}

var file_proto_coderef_code_reference_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_coderef_code_reference_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_coderef_code_reference_proto_goTypes = []interface{}{
	(CodeReference_RepositoryType)(0), // 0: bucketeer.coderef.CodeReference.RepositoryType
	(*CodeReference)(nil),             // 1: bucketeer.coderef.CodeReference
}
var file_proto_coderef_code_reference_proto_depIdxs = []int32{
	0, // 0: bucketeer.coderef.CodeReference.repository_type:type_name -> bucketeer.coderef.CodeReference.RepositoryType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_coderef_code_reference_proto_init() }
func file_proto_coderef_code_reference_proto_init() {
	if File_proto_coderef_code_reference_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_coderef_code_reference_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CodeReference); i {
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
			RawDescriptor: file_proto_coderef_code_reference_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_coderef_code_reference_proto_goTypes,
		DependencyIndexes: file_proto_coderef_code_reference_proto_depIdxs,
		EnumInfos:         file_proto_coderef_code_reference_proto_enumTypes,
		MessageInfos:      file_proto_coderef_code_reference_proto_msgTypes,
	}.Build()
	File_proto_coderef_code_reference_proto = out.File
	file_proto_coderef_code_reference_proto_rawDesc = nil
	file_proto_coderef_code_reference_proto_goTypes = nil
	file_proto_coderef_code_reference_proto_depIdxs = nil
}
