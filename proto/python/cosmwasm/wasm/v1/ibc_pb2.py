# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmwasm/wasm/v1/ibc.proto
# Protobuf Python Version: 6.30.2
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
    2,
    '',
    'cosmwasm/wasm/v1/ibc.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1a\x63osmwasm/wasm/v1/ibc.proto\x12\x10\x63osmwasm.wasm.v1\x1a\x14gogoproto/gogo.proto\"\xe2\x01\n\nMsgIBCSend\x12\x33\n\x07\x63hannel\x18\x02 \x01(\tB\x19\xf2\xde\x1f\x15yaml:\"source_channel\"R\x07\x63hannel\x12@\n\x0etimeout_height\x18\x04 \x01(\x04\x42\x19\xf2\xde\x1f\x15yaml:\"timeout_height\"R\rtimeoutHeight\x12I\n\x11timeout_timestamp\x18\x05 \x01(\x04\x42\x1c\xf2\xde\x1f\x18yaml:\"timeout_timestamp\"R\x10timeoutTimestamp\x12\x12\n\x04\x64\x61ta\x18\x06 \x01(\x0cR\x04\x64\x61ta\"0\n\x12MsgIBCSendResponse\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequence\"$\n\"MsgIBCWriteAcknowledgementResponse\"I\n\x12MsgIBCCloseChannel\x12\x33\n\x07\x63hannel\x18\x02 \x01(\tB\x19\xf2\xde\x1f\x15yaml:\"source_channel\"R\x07\x63hannelB,Z&github.com/CosmWasm/wasmd/x/wasm/types\xc8\xe1\x1e\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmwasm.wasm.v1.ibc_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z&github.com/CosmWasm/wasmd/x/wasm/types\310\341\036\000'
  _globals['_MSGIBCSEND'].fields_by_name['channel']._loaded_options = None
  _globals['_MSGIBCSEND'].fields_by_name['channel']._serialized_options = b'\362\336\037\025yaml:\"source_channel\"'
  _globals['_MSGIBCSEND'].fields_by_name['timeout_height']._loaded_options = None
  _globals['_MSGIBCSEND'].fields_by_name['timeout_height']._serialized_options = b'\362\336\037\025yaml:\"timeout_height\"'
  _globals['_MSGIBCSEND'].fields_by_name['timeout_timestamp']._loaded_options = None
  _globals['_MSGIBCSEND'].fields_by_name['timeout_timestamp']._serialized_options = b'\362\336\037\030yaml:\"timeout_timestamp\"'
  _globals['_MSGIBCCLOSECHANNEL'].fields_by_name['channel']._loaded_options = None
  _globals['_MSGIBCCLOSECHANNEL'].fields_by_name['channel']._serialized_options = b'\362\336\037\025yaml:\"source_channel\"'
  _globals['_MSGIBCSEND']._serialized_start=71
  _globals['_MSGIBCSEND']._serialized_end=297
  _globals['_MSGIBCSENDRESPONSE']._serialized_start=299
  _globals['_MSGIBCSENDRESPONSE']._serialized_end=347
  _globals['_MSGIBCWRITEACKNOWLEDGEMENTRESPONSE']._serialized_start=349
  _globals['_MSGIBCWRITEACKNOWLEDGEMENTRESPONSE']._serialized_end=385
  _globals['_MSGIBCCLOSECHANNEL']._serialized_start=387
  _globals['_MSGIBCCLOSECHANNEL']._serialized_end=460
# @@protoc_insertion_point(module_scope)
