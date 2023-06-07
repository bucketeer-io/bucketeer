// Copyright 2022 The Bucketeer Authors.
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
// 	protoc        v3.20.3
// source: proto/notification/sender/notification_event.proto

package sender

import (
	notification "github.com/bucketeer-io/bucketeer/proto/notification"
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

type NotificationEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                   string                               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	EnvironmentNamespace string                               `protobuf:"bytes,2,opt,name=environment_namespace,json=environmentNamespace,proto3" json:"environment_namespace,omitempty"`
	SourceType           notification.Subscription_SourceType `protobuf:"varint,3,opt,name=source_type,json=sourceType,proto3,enum=bucketeer.notification.Subscription_SourceType" json:"source_type,omitempty"`
	Notification         *Notification                        `protobuf:"bytes,4,opt,name=notification,proto3" json:"notification,omitempty"`
	IsAdminEvent         bool                                 `protobuf:"varint,5,opt,name=is_admin_event,json=isAdminEvent,proto3" json:"is_admin_event,omitempty"`
}

func (x *NotificationEvent) Reset() {
	*x = NotificationEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_notification_sender_notification_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationEvent) ProtoMessage() {}

func (x *NotificationEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_notification_sender_notification_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationEvent.ProtoReflect.Descriptor instead.
func (*NotificationEvent) Descriptor() ([]byte, []int) {
	return file_proto_notification_sender_notification_event_proto_rawDescGZIP(), []int{0}
}

func (x *NotificationEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *NotificationEvent) GetEnvironmentNamespace() string {
	if x != nil {
		return x.EnvironmentNamespace
	}
	return ""
}

func (x *NotificationEvent) GetSourceType() notification.Subscription_SourceType {
	if x != nil {
		return x.SourceType
	}
	return notification.Subscription_DOMAIN_EVENT_FEATURE
}

func (x *NotificationEvent) GetNotification() *Notification {
	if x != nil {
		return x.Notification
	}
	return nil
}

func (x *NotificationEvent) GetIsAdminEvent() bool {
	if x != nil {
		return x.IsAdminEvent
	}
	return false
}

var File_proto_notification_sender_notification_event_proto protoreflect.FileDescriptor

var file_proto_notification_sender_notification_event_proto_rawDesc = []byte{
	0x0a, 0x32, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x1a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x25, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x02, 0x0a, 0x11, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x33,
	0x0a, 0x15, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0e, 0x69, 0x73, 0x5f, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x69, 0x73, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x3d, 0x5a, 0x3b,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_notification_sender_notification_event_proto_rawDescOnce sync.Once
	file_proto_notification_sender_notification_event_proto_rawDescData = file_proto_notification_sender_notification_event_proto_rawDesc
)

func file_proto_notification_sender_notification_event_proto_rawDescGZIP() []byte {
	file_proto_notification_sender_notification_event_proto_rawDescOnce.Do(func() {
		file_proto_notification_sender_notification_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_notification_sender_notification_event_proto_rawDescData)
	})
	return file_proto_notification_sender_notification_event_proto_rawDescData
}

var file_proto_notification_sender_notification_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_notification_sender_notification_event_proto_goTypes = []interface{}{
	(*NotificationEvent)(nil),                 // 0: bucketeer.notification.sender.NotificationEvent
	(notification.Subscription_SourceType)(0), // 1: bucketeer.notification.Subscription.SourceType
	(*Notification)(nil),                      // 2: bucketeer.notification.sender.Notification
}
var file_proto_notification_sender_notification_event_proto_depIdxs = []int32{
	1, // 0: bucketeer.notification.sender.NotificationEvent.source_type:type_name -> bucketeer.notification.Subscription.SourceType
	2, // 1: bucketeer.notification.sender.NotificationEvent.notification:type_name -> bucketeer.notification.sender.Notification
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_notification_sender_notification_event_proto_init() }
func file_proto_notification_sender_notification_event_proto_init() {
	if File_proto_notification_sender_notification_event_proto != nil {
		return
	}
	file_proto_notification_sender_notification_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_notification_sender_notification_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationEvent); i {
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
			RawDescriptor: file_proto_notification_sender_notification_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_notification_sender_notification_event_proto_goTypes,
		DependencyIndexes: file_proto_notification_sender_notification_event_proto_depIdxs,
		MessageInfos:      file_proto_notification_sender_notification_event_proto_msgTypes,
	}.Build()
	File_proto_notification_sender_notification_event_proto = out.File
	file_proto_notification_sender_notification_event_proto_rawDesc = nil
	file_proto_notification_sender_notification_event_proto_goTypes = nil
	file_proto_notification_sender_notification_event_proto_depIdxs = nil
}
