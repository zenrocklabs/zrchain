# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/core/connection/v1/genesis.proto
# Protobuf Python Version: 6.31.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    31,
    1,
    '',
    'ibc/core/connection/v1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from ibc.core.connection.v1 import connection_pb2 as ibc_dot_core_dot_connection_dot_v1_dot_connection__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n$ibc/core/connection/v1/genesis.proto\x12\x16ibc.core.connection.v1\x1a\x14gogoproto/gogo.proto\x1a\'ibc/core/connection/v1/connection.proto\"\xc3\x02\n\x0cGenesisState\x12T\n\x0b\x63onnections\x18\x01 \x03(\x0b\x32,.ibc.core.connection.v1.IdentifiedConnectionB\x04\xc8\xde\x1f\x00R\x0b\x63onnections\x12\x65\n\x17\x63lient_connection_paths\x18\x02 \x03(\x0b\x32\'.ibc.core.connection.v1.ConnectionPathsB\x04\xc8\xde\x1f\x00R\x15\x63lientConnectionPaths\x12\x38\n\x18next_connection_sequence\x18\x03 \x01(\x04R\x16nextConnectionSequence\x12<\n\x06params\x18\x04 \x01(\x0b\x32\x1e.ibc.core.connection.v1.ParamsB\x04\xc8\xde\x1f\x00R\x06paramsB?Z=github.com/cosmos/ibc-go/v10/modules/core/03-connection/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.core.connection.v1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z=github.com/cosmos/ibc-go/v10/modules/core/03-connection/types'
  _globals['_GENESISSTATE'].fields_by_name['connections']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['connections']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['client_connection_paths']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['client_connection_paths']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['params']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE']._serialized_start=128
  _globals['_GENESISSTATE']._serialized_end=451
# @@protoc_insertion_point(module_scope)
