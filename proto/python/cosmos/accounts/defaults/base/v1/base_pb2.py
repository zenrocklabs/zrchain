# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/accounts/defaults/base/v1/base.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'cosmos/accounts/defaults/base/v1/base.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n+cosmos/accounts/defaults/base/v1/base.proto\x12 cosmos.accounts.defaults.base.v1\"\"\n\x07MsgInit\x12\x17\n\x07pub_key\x18\x01 \x01(\x0cR\x06pubKey\"\x11\n\x0fMsgInitResponse\"/\n\rMsgSwapPubKey\x12\x1e\n\x0bnew_pub_key\x18\x01 \x01(\x0cR\tnewPubKey\"\x17\n\x15MsgSwapPubKeyResponse\"\x0f\n\rQuerySequence\"3\n\x15QuerySequenceResponse\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequenceB*Z(cosmossdk.io/x/accounts/defaults/base/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.accounts.defaults.base.v1.base_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z(cosmossdk.io/x/accounts/defaults/base/v1'
  _globals['_MSGINIT']._serialized_start=81
  _globals['_MSGINIT']._serialized_end=115
  _globals['_MSGINITRESPONSE']._serialized_start=117
  _globals['_MSGINITRESPONSE']._serialized_end=134
  _globals['_MSGSWAPPUBKEY']._serialized_start=136
  _globals['_MSGSWAPPUBKEY']._serialized_end=183
  _globals['_MSGSWAPPUBKEYRESPONSE']._serialized_start=185
  _globals['_MSGSWAPPUBKEYRESPONSE']._serialized_end=208
  _globals['_QUERYSEQUENCE']._serialized_start=210
  _globals['_QUERYSEQUENCE']._serialized_end=225
  _globals['_QUERYSEQUENCERESPONSE']._serialized_start=227
  _globals['_QUERYSEQUENCERESPONSE']._serialized_end=278
# @@protoc_insertion_point(module_scope)
