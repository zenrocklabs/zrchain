# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/interchain_accounts/v1/metadata.proto
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
    'ibc/applications/interchain_accounts/v1/metadata.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n6ibc/applications/interchain_accounts/v1/metadata.proto\x12\'ibc.applications.interchain_accounts.v1\"\xdb\x01\n\x08Metadata\x12\x18\n\x07version\x18\x01 \x01(\tR\x07version\x12\x38\n\x18\x63ontroller_connection_id\x18\x02 \x01(\tR\x16\x63ontrollerConnectionId\x12,\n\x12host_connection_id\x18\x03 \x01(\tR\x10hostConnectionId\x12\x18\n\x07\x61\x64\x64ress\x18\x04 \x01(\tR\x07\x61\x64\x64ress\x12\x1a\n\x08\x65ncoding\x18\x05 \x01(\tR\x08\x65ncoding\x12\x17\n\x07tx_type\x18\x06 \x01(\tR\x06txTypeBHZFgithub.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.interchain_accounts.v1.metadata_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZFgithub.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types'
  _globals['_METADATA']._serialized_start=100
  _globals['_METADATA']._serialized_end=319
# @@protoc_insertion_point(module_scope)
