# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/mint/v1beta1/mint.proto
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
    'zrchain/mint/v1beta1/mint.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1fzrchain/mint/v1beta1/mint.proto\x12\x14zrchain.mint.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x11\x61mino/amino.proto\"\xb9\x01\n\x06Minter\x12O\n\tinflation\x18\x01 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\tinflation\x12^\n\x11\x61nnual_provisions\x18\x02 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\x10\x61nnualProvisions\"\xa2\t\n\x06Params\x12\x1d\n\nmint_denom\x18\x01 \x01(\tR\tmintDenom\x12j\n\x15inflation_rate_change\x18\x02 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x13inflationRateChange\x12[\n\rinflation_max\x18\x03 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x0cinflationMax\x12[\n\rinflation_min\x18\x04 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x0cinflationMin\x12W\n\x0bgoal_bonded\x18\x05 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\ngoalBonded\x12&\n\x0f\x62locks_per_year\x18\x06 \x01(\x04R\rblocksPerYear\x12[\n\rstaking_yield\x18\x07 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x0cstakingYield\x12S\n\tburn_rate\x18\x08 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x08\x62urnRate\x12h\n\x14protocol_wallet_rate\x18\t \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x12protocolWalletRate\x12\x36\n\x17protocol_wallet_address\x18\n \x01(\tR\x15protocolWalletAddress\x12t\n\x1a\x61\x64\x64itional_staking_rewards\x18\x0b \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x18\x61\x64\x64itionalStakingRewards\x12l\n\x16\x61\x64\x64itional_mpc_rewards\x18\x0c \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x14\x61\x64\x64itionalMpcRewards\x12h\n\x14\x61\x64\x64itional_burn_rate\x18\r \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x12\x61\x64\x64itionalBurnRate:0\x8a\xe7\xb0*+Zenrock-Foundation/zrchain/v5/x/mint/ParamsB7Z5github.com/Zenrock-Foundation/zrchain/v5/x/mint/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.mint.v1beta1.mint_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/Zenrock-Foundation/zrchain/v5/x/mint/types'
  _globals['_MINTER'].fields_by_name['inflation']._loaded_options = None
  _globals['_MINTER'].fields_by_name['inflation']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_MINTER'].fields_by_name['annual_provisions']._loaded_options = None
  _globals['_MINTER'].fields_by_name['annual_provisions']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['inflation_rate_change']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['inflation_rate_change']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['inflation_max']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['inflation_max']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['inflation_min']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['inflation_min']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['goal_bonded']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['goal_bonded']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['staking_yield']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['staking_yield']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['burn_rate']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['burn_rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['protocol_wallet_rate']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['protocol_wallet_rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['additional_staking_rewards']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['additional_staking_rewards']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['additional_mpc_rewards']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['additional_mpc_rewards']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['additional_burn_rate']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['additional_burn_rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\212\347\260*+Zenrock-Foundation/zrchain/v5/x/mint/Params'
  _globals['_MINTER']._serialized_start=126
  _globals['_MINTER']._serialized_end=311
  _globals['_PARAMS']._serialized_start=314
  _globals['_PARAMS']._serialized_end=1500
# @@protoc_insertion_point(module_scope)
