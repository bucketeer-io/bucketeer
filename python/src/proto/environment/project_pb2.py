# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/environment/project.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1fproto/environment/project.proto\x12\x15\x62ucketeer.environment\"\x8a\x01\n\x07Project\x12\n\n\x02id\x18\x01 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x02 \x01(\t\x12\x10\n\x08\x64isabled\x18\x03 \x01(\x08\x12\r\n\x05trial\x18\x04 \x01(\x08\x12\x15\n\rcreator_email\x18\x05 \x01(\t\x12\x12\n\ncreated_at\x18\x06 \x01(\x03\x12\x12\n\nupdated_at\x18\x07 \x01(\x03\x42\x35Z3github.com/bucketeer-io/bucketeer/proto/environmentb\x06proto3')



_PROJECT = DESCRIPTOR.message_types_by_name['Project']
Project = _reflection.GeneratedProtocolMessageType('Project', (_message.Message,), {
  'DESCRIPTOR' : _PROJECT,
  '__module__' : 'proto.environment.project_pb2'
  # @@protoc_insertion_point(class_scope:bucketeer.environment.Project)
  })
_sym_db.RegisterMessage(Project)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z3github.com/bucketeer-io/bucketeer/proto/environment'
  _PROJECT._serialized_start=59
  _PROJECT._serialized_end=197
# @@protoc_insertion_point(module_scope)
