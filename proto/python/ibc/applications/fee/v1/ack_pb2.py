# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/fee/v1/ack.proto
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
    'ibc/applications/fee/v1/ack.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n!ibc/applications/fee/v1/ack.proto\x12\x17ibc.applications.fee.v1\"\xbc\x01\n\x1bIncentivizedAcknowledgement\x12/\n\x13\x61pp_acknowledgement\x18\x01 \x01(\x0cR\x12\x61ppAcknowledgement\x12\x36\n\x17\x66orward_relayer_address\x18\x02 \x01(\tR\x15\x66orwardRelayerAddress\x12\x34\n\x16underlying_app_success\x18\x03 \x01(\x08R\x14underlyingAppSuccessB7Z5github.com/cosmos/ibc-go/v9/modules/apps/29-fee/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.fee.v1.ack_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/cosmos/ibc-go/v9/modules/apps/29-fee/types'
  _globals['_INCENTIVIZEDACKNOWLEDGEMENT']._serialized_start=63
  _globals['_INCENTIVIZEDACKNOWLEDGEMENT']._serialized_end=251
# @@protoc_insertion_point(module_scope)