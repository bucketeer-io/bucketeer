# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/eventcounter/timeseries.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n#proto/eventcounter/timeseries.proto\x12\x16\x62ucketeer.eventcounter\"c\n\x13VariationTimeseries\x12\x14\n\x0cvariation_id\x18\x01 \x01(\t\x12\x36\n\ntimeseries\x18\x02 \x01(\x0b\x32\".bucketeer.eventcounter.Timeseries\"\x98\x01\n\nTimeseries\x12\x12\n\ntimestamps\x18\x01 \x03(\x03\x12\x0e\n\x06values\x18\x02 \x03(\x01\x12\x35\n\x04unit\x18\x03 \x01(\x0e\x32\'.bucketeer.eventcounter.Timeseries.Unit\x12\x14\n\x0ctotal_counts\x18\x04 \x01(\x03\"\x19\n\x04Unit\x12\x08\n\x04HOUR\x10\x00\x12\x07\n\x03\x44\x41Y\x10\x01\x42\x36Z4github.com/bucketeer-io/bucketeer/proto/eventcounterb\x06proto3')



_VARIATIONTIMESERIES = DESCRIPTOR.message_types_by_name['VariationTimeseries']
_TIMESERIES = DESCRIPTOR.message_types_by_name['Timeseries']
_TIMESERIES_UNIT = _TIMESERIES.enum_types_by_name['Unit']
VariationTimeseries = _reflection.GeneratedProtocolMessageType('VariationTimeseries', (_message.Message,), {
  'DESCRIPTOR' : _VARIATIONTIMESERIES,
  '__module__' : 'proto.eventcounter.timeseries_pb2'
  # @@protoc_insertion_point(class_scope:bucketeer.eventcounter.VariationTimeseries)
  })
_sym_db.RegisterMessage(VariationTimeseries)

Timeseries = _reflection.GeneratedProtocolMessageType('Timeseries', (_message.Message,), {
  'DESCRIPTOR' : _TIMESERIES,
  '__module__' : 'proto.eventcounter.timeseries_pb2'
  # @@protoc_insertion_point(class_scope:bucketeer.eventcounter.Timeseries)
  })
_sym_db.RegisterMessage(Timeseries)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z4github.com/bucketeer-io/bucketeer/proto/eventcounter'
  _VARIATIONTIMESERIES._serialized_start=63
  _VARIATIONTIMESERIES._serialized_end=162
  _TIMESERIES._serialized_start=165
  _TIMESERIES._serialized_end=317
  _TIMESERIES_UNIT._serialized_start=292
  _TIMESERIES_UNIT._serialized_end=317
# @@protoc_insertion_point(module_scope)
