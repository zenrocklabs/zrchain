# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/accounts/testing/rotation/v1/partial.proto
# Protobuf Python Version: 5.28.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    2,
    '',
    'cosmos/accounts/testing/rotation/v1/partial.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n1cosmos/accounts/testing/rotation/v1/partial.proto\x12#cosmos.accounts.testing.rotation.v1\"-\n\x07MsgInit\x12\"\n\rpub_key_bytes\x18\x01 \x01(\x0cR\x0bpubKeyBytes\"\x11\n\x0fMsgInitResponse\"<\n\x0fMsgRotatePubKey\x12)\n\x11new_pub_key_bytes\x18\x01 \x01(\x0cR\x0enewPubKeyBytes\"\x19\n\x17MsgRotatePubKeyResponseB-Z+cosmossdk.io/x/accounts/testing/rotation/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.accounts.testing.rotation.v1.partial_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+cosmossdk.io/x/accounts/testing/rotation/v1'
  _globals['_MSGINIT']._serialized_start=90
  _globals['_MSGINIT']._serialized_end=135
  _globals['_MSGINITRESPONSE']._serialized_start=137
  _globals['_MSGINITRESPONSE']._serialized_end=154
  _globals['_MSGROTATEPUBKEY']._serialized_start=156
  _globals['_MSGROTATEPUBKEY']._serialized_end=216
  _globals['_MSGROTATEPUBKEYRESPONSE']._serialized_start=218
  _globals['_MSGROTATEPUBKEYRESPONSE']._serialized_end=243
# @@protoc_insertion_point(module_scope)