# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/validation/authz.proto
# Protobuf Python Version: 6.30.1
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
    1,
    '',
    'zrchain/validation/authz.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1ezrchain/validation/authz.proto\x12\x12zrchain.validation\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x11\x61mino/amino.proto\"\xe3\x04\n\x12StakeAuthorization\x12\x65\n\nmax_tokens\x18\x01 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinB+\xaa\xdf\x1f\'github.com/cosmos/cosmos-sdk/types.CoinR\tmaxTokens\x12}\n\nallow_list\x18\x02 \x01(\x0b\x32\x31.zrchain.validation.StakeAuthorization.ValidatorsB)\xb2\xe7\xb0*$zrchain/StakeAuthorization/AllowListH\x00R\tallowList\x12z\n\tdeny_list\x18\x03 \x01(\x0b\x32\x31.zrchain.validation.StakeAuthorization.ValidatorsB(\xb2\xe7\xb0*#zrchain/StakeAuthorization/DenyListH\x00R\x08\x64\x65nyList\x12T\n\x12\x61uthorization_type\x18\x04 \x01(\x0e\x32%.zrchain.validation.AuthorizationTypeR\x11\x61uthorizationType\x1a@\n\nValidators\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x03(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress:E\xca\xb4-\"cosmos.authz.v1beta1.Authorization\x8a\xe7\xb0*\x1azrchain/StakeAuthorizationB\x0c\n\nvalidators*\xd2\x01\n\x11\x41uthorizationType\x12\"\n\x1e\x41UTHORIZATION_TYPE_UNSPECIFIED\x10\x00\x12\x1f\n\x1b\x41UTHORIZATION_TYPE_DELEGATE\x10\x01\x12!\n\x1d\x41UTHORIZATION_TYPE_UNDELEGATE\x10\x02\x12!\n\x1d\x41UTHORIZATION_TYPE_REDELEGATE\x10\x03\x12\x32\n.AUTHORIZATION_TYPE_CANCEL_UNBONDING_DELEGATION\x10\x04\x42=Z;github.com/Zenrock-Foundation/zrchain/v5/x/validation/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.validation.authz_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z;github.com/Zenrock-Foundation/zrchain/v5/x/validation/types'
  _globals['_STAKEAUTHORIZATION_VALIDATORS'].fields_by_name['address']._loaded_options = None
  _globals['_STAKEAUTHORIZATION_VALIDATORS'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_STAKEAUTHORIZATION'].fields_by_name['max_tokens']._loaded_options = None
  _globals['_STAKEAUTHORIZATION'].fields_by_name['max_tokens']._serialized_options = b'\252\337\037\'github.com/cosmos/cosmos-sdk/types.Coin'
  _globals['_STAKEAUTHORIZATION'].fields_by_name['allow_list']._loaded_options = None
  _globals['_STAKEAUTHORIZATION'].fields_by_name['allow_list']._serialized_options = b'\262\347\260*$zrchain/StakeAuthorization/AllowList'
  _globals['_STAKEAUTHORIZATION'].fields_by_name['deny_list']._loaded_options = None
  _globals['_STAKEAUTHORIZATION'].fields_by_name['deny_list']._serialized_options = b'\262\347\260*#zrchain/StakeAuthorization/DenyList'
  _globals['_STAKEAUTHORIZATION']._loaded_options = None
  _globals['_STAKEAUTHORIZATION']._serialized_options = b'\312\264-\"cosmos.authz.v1beta1.Authorization\212\347\260*\032zrchain/StakeAuthorization'
  _globals['_AUTHORIZATIONTYPE']._serialized_start=769
  _globals['_AUTHORIZATIONTYPE']._serialized_end=979
  _globals['_STAKEAUTHORIZATION']._serialized_start=155
  _globals['_STAKEAUTHORIZATION']._serialized_end=766
  _globals['_STAKEAUTHORIZATION_VALIDATORS']._serialized_start=617
  _globals['_STAKEAUTHORIZATION_VALIDATORS']._serialized_end=681
# @@protoc_insertion_point(module_scope)
