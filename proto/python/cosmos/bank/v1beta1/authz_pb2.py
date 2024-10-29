# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/bank/v1beta1/authz.proto
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
    'cosmos/bank/v1beta1/authz.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/bank/v1beta1/authz.proto\x12\x13\x63osmos.bank.v1beta1\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\"\x9a\x02\n\x11SendAuthorization\x12\x82\x01\n\x0bspend_limit\x18\x01 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\nspendLimit\x12\x37\n\nallow_list\x18\x02 \x03(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tallowList:G\xca\xb4-\"cosmos.authz.v1beta1.Authorization\x8a\xe7\xb0*\x1c\x63osmos-sdk/SendAuthorizationB\x1bZ\x19\x63osmossdk.io/x/bank/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.bank.v1beta1.authz_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\031cosmossdk.io/x/bank/types'
  _globals['_SENDAUTHORIZATION'].fields_by_name['spend_limit']._loaded_options = None
  _globals['_SENDAUTHORIZATION'].fields_by_name['spend_limit']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_SENDAUTHORIZATION'].fields_by_name['allow_list']._loaded_options = None
  _globals['_SENDAUTHORIZATION'].fields_by_name['allow_list']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_SENDAUTHORIZATION']._loaded_options = None
  _globals['_SENDAUTHORIZATION']._serialized_options = b'\312\264-\"cosmos.authz.v1beta1.Authorization\212\347\260*\034cosmos-sdk/SendAuthorization'
  _globals['_SENDAUTHORIZATION']._serialized_start=157
  _globals['_SENDAUTHORIZATION']._serialized_end=439
# @@protoc_insertion_point(module_scope)