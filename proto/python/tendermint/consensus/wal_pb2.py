# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/consensus/wal.proto
# Protobuf Python Version: 6.30.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    30,
    0,
    '',
    'tendermint/consensus/wal.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from tendermint.consensus import types_pb2 as tendermint_dot_consensus_dot_types__pb2
from tendermint.types import events_pb2 as tendermint_dot_types_dot_events__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1etendermint/consensus/wal.proto\x12\x14tendermint.consensus\x1a\x14gogoproto/gogo.proto\x1a tendermint/consensus/types.proto\x1a\x1dtendermint/types/events.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"e\n\x07MsgInfo\x12\x35\n\x03msg\x18\x01 \x01(\x0b\x32\x1d.tendermint.consensus.MessageB\x04\xc8\xde\x1f\x00R\x03msg\x12#\n\x07peer_id\x18\x02 \x01(\tB\n\xe2\xde\x1f\x06PeerIDR\x06peerId\"\x90\x01\n\x0bTimeoutInfo\x12?\n\x08\x64uration\x18\x01 \x01(\x0b\x32\x19.google.protobuf.DurationB\x08\xc8\xde\x1f\x00\x98\xdf\x1f\x01R\x08\x64uration\x12\x16\n\x06height\x18\x02 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x03 \x01(\x05R\x05round\x12\x12\n\x04step\x18\x04 \x01(\rR\x04step\"#\n\tEndHeight\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\"\xb7\x02\n\nWALMessage\x12\\\n\x16\x65vent_data_round_state\x18\x01 \x01(\x0b\x32%.tendermint.types.EventDataRoundStateH\x00R\x13\x65ventDataRoundState\x12:\n\x08msg_info\x18\x02 \x01(\x0b\x32\x1d.tendermint.consensus.MsgInfoH\x00R\x07msgInfo\x12\x46\n\x0ctimeout_info\x18\x03 \x01(\x0b\x32!.tendermint.consensus.TimeoutInfoH\x00R\x0btimeoutInfo\x12@\n\nend_height\x18\x04 \x01(\x0b\x32\x1f.tendermint.consensus.EndHeightH\x00R\tendHeightB\x05\n\x03sum\"\x7f\n\x0fTimedWALMessage\x12\x38\n\x04time\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x04time\x12\x32\n\x03msg\x18\x02 \x01(\x0b\x32 .tendermint.consensus.WALMessageR\x03msgB9Z7github.com/cometbft/cometbft/proto/tendermint/consensusb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.consensus.wal_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/cometbft/cometbft/proto/tendermint/consensus'
  _globals['_MSGINFO'].fields_by_name['msg']._loaded_options = None
  _globals['_MSGINFO'].fields_by_name['msg']._serialized_options = b'\310\336\037\000'
  _globals['_MSGINFO'].fields_by_name['peer_id']._loaded_options = None
  _globals['_MSGINFO'].fields_by_name['peer_id']._serialized_options = b'\342\336\037\006PeerID'
  _globals['_TIMEOUTINFO'].fields_by_name['duration']._loaded_options = None
  _globals['_TIMEOUTINFO'].fields_by_name['duration']._serialized_options = b'\310\336\037\000\230\337\037\001'
  _globals['_TIMEDWALMESSAGE'].fields_by_name['time']._loaded_options = None
  _globals['_TIMEDWALMESSAGE'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_MSGINFO']._serialized_start=208
  _globals['_MSGINFO']._serialized_end=309
  _globals['_TIMEOUTINFO']._serialized_start=312
  _globals['_TIMEOUTINFO']._serialized_end=456
  _globals['_ENDHEIGHT']._serialized_start=458
  _globals['_ENDHEIGHT']._serialized_end=493
  _globals['_WALMESSAGE']._serialized_start=496
  _globals['_WALMESSAGE']._serialized_end=807
  _globals['_TIMEDWALMESSAGE']._serialized_start=809
  _globals['_TIMEDWALMESSAGE']._serialized_end=936
# @@protoc_insertion_point(module_scope)
