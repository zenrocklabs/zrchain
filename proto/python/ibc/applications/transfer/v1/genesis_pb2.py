# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/transfer/v1/genesis.proto
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
    'ibc/applications/transfer/v1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from ibc.applications.transfer.v1 import transfer_pb2 as ibc_dot_applications_dot_transfer_dot_v1_dot_transfer__pb2
from ibc.applications.transfer.v1 import token_pb2 as ibc_dot_applications_dot_transfer_dot_v1_dot_token__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n*ibc/applications/transfer/v1/genesis.proto\x12\x1cibc.applications.transfer.v1\x1a+ibc/applications/transfer/v1/transfer.proto\x1a(ibc/applications/transfer/v1/token.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x14gogoproto/gogo.proto\"\xac\x02\n\x0cGenesisState\x12\x17\n\x07port_id\x18\x01 \x01(\tR\x06portId\x12K\n\x06\x64\x65noms\x18\x02 \x03(\x0b\x32#.ibc.applications.transfer.v1.DenomB\x0e\xc8\xde\x1f\x00\xaa\xdf\x1f\x06\x44\x65nomsR\x06\x64\x65noms\x12\x42\n\x06params\x18\x03 \x01(\x0b\x32$.ibc.applications.transfer.v1.ParamsB\x04\xc8\xde\x1f\x00R\x06params\x12r\n\x0etotal_escrowed\x18\x04 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\rtotalEscrowedB:Z8github.com/cosmos/ibc-go/v10/modules/apps/transfer/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.transfer.v1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z8github.com/cosmos/ibc-go/v10/modules/apps/transfer/types'
  _globals['_GENESISSTATE'].fields_by_name['denoms']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['denoms']._serialized_options = b'\310\336\037\000\252\337\037\006Denoms'
  _globals['_GENESISSTATE'].fields_by_name['params']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['total_escrowed']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['total_escrowed']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_GENESISSTATE']._serialized_start=218
  _globals['_GENESISSTATE']._serialized_end=518
# @@protoc_insertion_point(module_scope)
