# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: proto/feature/strategy.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1cproto/feature/strategy.proto\x12\x11\x62ucketeer.feature\"\"\n\rFixedStrategy\x12\x11\n\tvariation\x18\x01 \x01(\t\"\x83\x01\n\x0fRolloutStrategy\x12@\n\nvariations\x18\x01 \x03(\x0b\x32,.bucketeer.feature.RolloutStrategy.Variation\x1a.\n\tVariation\x12\x11\n\tvariation\x18\x01 \x01(\t\x12\x0e\n\x06weight\x18\x02 \x01(\x05\"\xd2\x01\n\x08Strategy\x12.\n\x04type\x18\x01 \x01(\x0e\x32 .bucketeer.feature.Strategy.Type\x12\x38\n\x0e\x66ixed_strategy\x18\x02 \x01(\x0b\x32 .bucketeer.feature.FixedStrategy\x12<\n\x10rollout_strategy\x18\x03 \x01(\x0b\x32\".bucketeer.feature.RolloutStrategy\"\x1e\n\x04Type\x12\t\n\x05\x46IXED\x10\x00\x12\x0b\n\x07ROLLOUT\x10\x01\x42\x31Z/github.com/bucketeer-io/bucketeer/proto/featureb\x06proto3')



_FIXEDSTRATEGY = DESCRIPTOR.message_types_by_name['FixedStrategy']
_ROLLOUTSTRATEGY = DESCRIPTOR.message_types_by_name['RolloutStrategy']
_ROLLOUTSTRATEGY_VARIATION = _ROLLOUTSTRATEGY.nested_types_by_name['Variation']
_STRATEGY = DESCRIPTOR.message_types_by_name['Strategy']
_STRATEGY_TYPE = _STRATEGY.enum_types_by_name['Type']
FixedStrategy = _reflection.GeneratedProtocolMessageType('FixedStrategy', (_message.Message,), {
  'DESCRIPTOR' : _FIXEDSTRATEGY,
  '__module__' : 'proto.feature.strategy_pb2'
  # @@protoc_insertion_point(class_scope:bucketeer.feature.FixedStrategy)
  })
_sym_db.RegisterMessage(FixedStrategy)

RolloutStrategy = _reflection.GeneratedProtocolMessageType('RolloutStrategy', (_message.Message,), {

  'Variation' : _reflection.GeneratedProtocolMessageType('Variation', (_message.Message,), {
    'DESCRIPTOR' : _ROLLOUTSTRATEGY_VARIATION,
    '__module__' : 'proto.feature.strategy_pb2'
    # @@protoc_insertion_point(class_scope:bucketeer.feature.RolloutStrategy.Variation)
    })
  ,
  'DESCRIPTOR' : _ROLLOUTSTRATEGY,
  '__module__' : 'proto.feature.strategy_pb2'
  # @@protoc_insertion_point(class_scope:bucketeer.feature.RolloutStrategy)
  })
_sym_db.RegisterMessage(RolloutStrategy)
_sym_db.RegisterMessage(RolloutStrategy.Variation)

Strategy = _reflection.GeneratedProtocolMessageType('Strategy', (_message.Message,), {
  'DESCRIPTOR' : _STRATEGY,
  '__module__' : 'proto.feature.strategy_pb2'
  # @@protoc_insertion_point(class_scope:bucketeer.feature.Strategy)
  })
_sym_db.RegisterMessage(Strategy)

if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z/github.com/bucketeer-io/bucketeer/proto/feature'
  _FIXEDSTRATEGY._serialized_start=51
  _FIXEDSTRATEGY._serialized_end=85
  _ROLLOUTSTRATEGY._serialized_start=88
  _ROLLOUTSTRATEGY._serialized_end=219
  _ROLLOUTSTRATEGY_VARIATION._serialized_start=173
  _ROLLOUTSTRATEGY_VARIATION._serialized_end=219
  _STRATEGY._serialized_start=222
  _STRATEGY._serialized_end=432
  _STRATEGY_TYPE._serialized_start=402
  _STRATEGY_TYPE._serialized_end=432
# @@protoc_insertion_point(module_scope)
