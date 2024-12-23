# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/vesting/v1beta1/vesting.proto
# Protobuf Python Version: 5.29.2
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
    2,
    '',
    'cosmos/vesting/v1beta1/vesting.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos.auth.v1beta1 import auth_pb2 as cosmos_dot_auth_dot_v1beta1_dot_auth__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n$cosmos/vesting/v1beta1/vesting.proto\x12\x16\x63osmos.vesting.v1beta1\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x1e\x63osmos/auth/v1beta1/auth.proto\"\xcd\x04\n\x12\x42\x61seVestingAccount\x12I\n\x0c\x62\x61se_account\x18\x01 \x01(\x0b\x32 .cosmos.auth.v1beta1.BaseAccountB\x04\xd0\xde\x1f\x01R\x0b\x62\x61seAccount\x12\x8c\x01\n\x10original_vesting\x18\x02 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x0foriginalVesting\x12\x88\x01\n\x0e\x64\x65legated_free\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\rdelegatedFree\x12\x8e\x01\n\x11\x64\x65legated_vesting\x18\x04 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x10\x64\x65legatedVesting\x12\x19\n\x08\x65nd_time\x18\x05 \x01(\x03R\x07\x65ndTime:&\x88\xa0\x1f\x00\x8a\xe7\xb0*\x1d\x63osmos-sdk/BaseVestingAccount\"\xcb\x01\n\x18\x43ontinuousVestingAccount\x12\x62\n\x14\x62\x61se_vesting_account\x18\x01 \x01(\x0b\x32*.cosmos.vesting.v1beta1.BaseVestingAccountB\x04\xd0\xde\x1f\x01R\x12\x62\x61seVestingAccount\x12\x1d\n\nstart_time\x18\x02 \x01(\x03R\tstartTime:,\x88\xa0\x1f\x00\x8a\xe7\xb0*#cosmos-sdk/ContinuousVestingAccount\"\xa6\x01\n\x15\x44\x65layedVestingAccount\x12\x62\n\x14\x62\x61se_vesting_account\x18\x01 \x01(\x0b\x32*.cosmos.vesting.v1beta1.BaseVestingAccountB\x04\xd0\xde\x1f\x01R\x12\x62\x61seVestingAccount:)\x88\xa0\x1f\x00\x8a\xe7\xb0* cosmos-sdk/DelayedVestingAccount\"\x9b\x01\n\x06Period\x12\x16\n\x06length\x18\x01 \x01(\x03R\x06length\x12y\n\x06\x61mount\x18\x02 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x06\x61mount\"\x9b\x02\n\x16PeriodicVestingAccount\x12\x62\n\x14\x62\x61se_vesting_account\x18\x01 \x01(\x0b\x32*.cosmos.vesting.v1beta1.BaseVestingAccountB\x04\xd0\xde\x1f\x01R\x12\x62\x61seVestingAccount\x12\x1d\n\nstart_time\x18\x02 \x01(\x03R\tstartTime\x12R\n\x0fvesting_periods\x18\x03 \x03(\x0b\x32\x1e.cosmos.vesting.v1beta1.PeriodB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0evestingPeriods:*\x88\xa0\x1f\x00\x8a\xe7\xb0*!cosmos-sdk/PeriodicVestingAccount\"\xa8\x01\n\x16PermanentLockedAccount\x12\x62\n\x14\x62\x61se_vesting_account\x18\x01 \x01(\x0b\x32*.cosmos.vesting.v1beta1.BaseVestingAccountB\x04\xd0\xde\x1f\x01R\x12\x62\x61seVestingAccount:*\x88\xa0\x1f\x00\x8a\xe7\xb0*!cosmos-sdk/PermanentLockedAccountB#Z!cosmossdk.io/x/auth/vesting/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.vesting.v1beta1.vesting_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z!cosmossdk.io/x/auth/vesting/types'
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['base_account']._loaded_options = None
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['base_account']._serialized_options = b'\320\336\037\001'
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['original_vesting']._loaded_options = None
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['original_vesting']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['delegated_free']._loaded_options = None
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['delegated_free']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['delegated_vesting']._loaded_options = None
  _globals['_BASEVESTINGACCOUNT'].fields_by_name['delegated_vesting']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_BASEVESTINGACCOUNT']._loaded_options = None
  _globals['_BASEVESTINGACCOUNT']._serialized_options = b'\210\240\037\000\212\347\260*\035cosmos-sdk/BaseVestingAccount'
  _globals['_CONTINUOUSVESTINGACCOUNT'].fields_by_name['base_vesting_account']._loaded_options = None
  _globals['_CONTINUOUSVESTINGACCOUNT'].fields_by_name['base_vesting_account']._serialized_options = b'\320\336\037\001'
  _globals['_CONTINUOUSVESTINGACCOUNT']._loaded_options = None
  _globals['_CONTINUOUSVESTINGACCOUNT']._serialized_options = b'\210\240\037\000\212\347\260*#cosmos-sdk/ContinuousVestingAccount'
  _globals['_DELAYEDVESTINGACCOUNT'].fields_by_name['base_vesting_account']._loaded_options = None
  _globals['_DELAYEDVESTINGACCOUNT'].fields_by_name['base_vesting_account']._serialized_options = b'\320\336\037\001'
  _globals['_DELAYEDVESTINGACCOUNT']._loaded_options = None
  _globals['_DELAYEDVESTINGACCOUNT']._serialized_options = b'\210\240\037\000\212\347\260* cosmos-sdk/DelayedVestingAccount'
  _globals['_PERIOD'].fields_by_name['amount']._loaded_options = None
  _globals['_PERIOD'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_PERIODICVESTINGACCOUNT'].fields_by_name['base_vesting_account']._loaded_options = None
  _globals['_PERIODICVESTINGACCOUNT'].fields_by_name['base_vesting_account']._serialized_options = b'\320\336\037\001'
  _globals['_PERIODICVESTINGACCOUNT'].fields_by_name['vesting_periods']._loaded_options = None
  _globals['_PERIODICVESTINGACCOUNT'].fields_by_name['vesting_periods']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_PERIODICVESTINGACCOUNT']._loaded_options = None
  _globals['_PERIODICVESTINGACCOUNT']._serialized_options = b'\210\240\037\000\212\347\260*!cosmos-sdk/PeriodicVestingAccount'
  _globals['_PERMANENTLOCKEDACCOUNT'].fields_by_name['base_vesting_account']._loaded_options = None
  _globals['_PERMANENTLOCKEDACCOUNT'].fields_by_name['base_vesting_account']._serialized_options = b'\320\336\037\001'
  _globals['_PERMANENTLOCKEDACCOUNT']._loaded_options = None
  _globals['_PERMANENTLOCKEDACCOUNT']._serialized_options = b'\210\240\037\000\212\347\260*!cosmos-sdk/PermanentLockedAccount'
  _globals['_BASEVESTINGACCOUNT']._serialized_start=170
  _globals['_BASEVESTINGACCOUNT']._serialized_end=759
  _globals['_CONTINUOUSVESTINGACCOUNT']._serialized_start=762
  _globals['_CONTINUOUSVESTINGACCOUNT']._serialized_end=965
  _globals['_DELAYEDVESTINGACCOUNT']._serialized_start=968
  _globals['_DELAYEDVESTINGACCOUNT']._serialized_end=1134
  _globals['_PERIOD']._serialized_start=1137
  _globals['_PERIOD']._serialized_end=1292
  _globals['_PERIODICVESTINGACCOUNT']._serialized_start=1295
  _globals['_PERIODICVESTINGACCOUNT']._serialized_end=1578
  _globals['_PERMANENTLOCKEDACCOUNT']._serialized_start=1581
  _globals['_PERMANENTLOCKEDACCOUNT']._serialized_end=1749
# @@protoc_insertion_point(module_scope)
