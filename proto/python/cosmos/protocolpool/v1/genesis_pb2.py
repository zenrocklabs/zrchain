# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/protocolpool/v1/genesis.proto
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
    'cosmos/protocolpool/v1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.protocolpool.v1 import types_pb2 as cosmos_dot_protocolpool_dot_v1_dot_types__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n$cosmos/protocolpool/v1/genesis.proto\x12\x16\x63osmos.protocolpool.v1\x1a\"cosmos/protocolpool/v1/types.proto\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\xe9\x01\n\x0cGenesisState\x12O\n\x0f\x63ontinuous_fund\x18\x01 \x03(\x0b\x32&.cosmos.protocolpool.v1.ContinuousFundR\x0e\x63ontinuousFund\x12\x36\n\x06\x62udget\x18\x02 \x03(\x0b\x32\x1e.cosmos.protocolpool.v1.BudgetR\x06\x62udget\x12P\n\rto_distribute\x18\x03 \x01(\tB+\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xd2\xb4-\ncosmos.IntR\x0ctoDistributeB#Z!cosmossdk.io/x/protocolpool/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.protocolpool.v1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z!cosmossdk.io/x/protocolpool/types'
  _globals['_GENESISSTATE'].fields_by_name['to_distribute']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['to_distribute']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\322\264-\ncosmos.Int'
  _globals['_GENESISSTATE']._serialized_start=150
  _globals['_GENESISSTATE']._serialized_end=383
# @@protoc_insertion_point(module_scope)
