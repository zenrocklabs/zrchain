# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/vesting/v1beta1/tx.proto
# Protobuf Python Version: 5.28.3
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
    3,
    '',
    'cosmos/vesting/v1beta1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.vesting.v1beta1 import vesting_pb2 as cosmos_dot_vesting_dot_v1beta1_dot_vesting__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/vesting/v1beta1/tx.proto\x12\x16\x63osmos.vesting.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a$cosmos/vesting/v1beta1/vesting.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x11\x61mino/amino.proto\"\x9c\x03\n\x17MsgCreateVestingAccount\x12;\n\x0c\x66rom_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x0b\x66romAddress\x12\x37\n\nto_address\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\ttoAddress\x12y\n\x06\x61mount\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x06\x61mount\x12\x19\n\x08\x65nd_time\x18\x04 \x01(\x03R\x07\x65ndTime\x12\x18\n\x07\x64\x65layed\x18\x05 \x01(\x08R\x07\x64\x65layed\x12\x1d\n\nstart_time\x18\x06 \x01(\x03R\tstartTime:<\xe8\xa0\x1f\x01\x82\xe7\xb0*\x0c\x66rom_address\x8a\xe7\xb0*\"cosmos-sdk/MsgCreateVestingAccount\"!\n\x1fMsgCreateVestingAccountResponse\"\xcf\x02\n\x1fMsgCreatePermanentLockedAccount\x12:\n\x0c\x66rom_address\x18\x01 \x01(\tB\x17\xf2\xde\x1f\x13yaml:\"from_address\"R\x0b\x66romAddress\x12\x34\n\nto_address\x18\x02 \x01(\tB\x15\xf2\xde\x1f\x11yaml:\"to_address\"R\ttoAddress\x12y\n\x06\x61mount\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x06\x61mount:?\xe8\xa0\x1f\x01\x82\xe7\xb0*\x0c\x66rom_address\x8a\xe7\xb0*%cosmos-sdk/MsgCreatePermLockedAccount\")\n\'MsgCreatePermanentLockedAccountResponse\"\x97\x02\n\x1fMsgCreatePeriodicVestingAccount\x12!\n\x0c\x66rom_address\x18\x01 \x01(\tR\x0b\x66romAddress\x12\x1d\n\nto_address\x18\x02 \x01(\tR\ttoAddress\x12\x1d\n\nstart_time\x18\x03 \x01(\x03R\tstartTime\x12R\n\x0fvesting_periods\x18\x04 \x03(\x0b\x32\x1e.cosmos.vesting.v1beta1.PeriodB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0evestingPeriods:?\xe8\xa0\x1f\x00\x82\xe7\xb0*\x0c\x66rom_address\x8a\xe7\xb0*%cosmos-sdk/MsgCreatePeriodVestAccount\")\n\'MsgCreatePeriodicVestingAccountResponse2\xc5\x03\n\x03Msg\x12\x80\x01\n\x14\x43reateVestingAccount\x12/.cosmos.vesting.v1beta1.MsgCreateVestingAccount\x1a\x37.cosmos.vesting.v1beta1.MsgCreateVestingAccountResponse\x12\x98\x01\n\x1c\x43reatePermanentLockedAccount\x12\x37.cosmos.vesting.v1beta1.MsgCreatePermanentLockedAccount\x1a?.cosmos.vesting.v1beta1.MsgCreatePermanentLockedAccountResponse\x12\x98\x01\n\x1c\x43reatePeriodicVestingAccount\x12\x37.cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccount\x1a?.cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccountResponse\x1a\x05\x80\xe7\xb0*\x01\x42#Z!cosmossdk.io/x/auth/vesting/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.vesting.v1beta1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z!cosmossdk.io/x/auth/vesting/types'
  _globals['_MSGCREATEVESTINGACCOUNT'].fields_by_name['from_address']._loaded_options = None
  _globals['_MSGCREATEVESTINGACCOUNT'].fields_by_name['from_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCREATEVESTINGACCOUNT'].fields_by_name['to_address']._loaded_options = None
  _globals['_MSGCREATEVESTINGACCOUNT'].fields_by_name['to_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCREATEVESTINGACCOUNT'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGCREATEVESTINGACCOUNT'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_MSGCREATEVESTINGACCOUNT']._loaded_options = None
  _globals['_MSGCREATEVESTINGACCOUNT']._serialized_options = b'\350\240\037\001\202\347\260*\014from_address\212\347\260*\"cosmos-sdk/MsgCreateVestingAccount'
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT'].fields_by_name['from_address']._loaded_options = None
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT'].fields_by_name['from_address']._serialized_options = b'\362\336\037\023yaml:\"from_address\"'
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT'].fields_by_name['to_address']._loaded_options = None
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT'].fields_by_name['to_address']._serialized_options = b'\362\336\037\021yaml:\"to_address\"'
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT']._loaded_options = None
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT']._serialized_options = b'\350\240\037\001\202\347\260*\014from_address\212\347\260*%cosmos-sdk/MsgCreatePermLockedAccount'
  _globals['_MSGCREATEPERIODICVESTINGACCOUNT'].fields_by_name['vesting_periods']._loaded_options = None
  _globals['_MSGCREATEPERIODICVESTINGACCOUNT'].fields_by_name['vesting_periods']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGCREATEPERIODICVESTINGACCOUNT']._loaded_options = None
  _globals['_MSGCREATEPERIODICVESTINGACCOUNT']._serialized_options = b'\350\240\037\000\202\347\260*\014from_address\212\347\260*%cosmos-sdk/MsgCreatePeriodVestAccount'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGCREATEVESTINGACCOUNT']._serialized_start=223
  _globals['_MSGCREATEVESTINGACCOUNT']._serialized_end=635
  _globals['_MSGCREATEVESTINGACCOUNTRESPONSE']._serialized_start=637
  _globals['_MSGCREATEVESTINGACCOUNTRESPONSE']._serialized_end=670
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT']._serialized_start=673
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNT']._serialized_end=1008
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNTRESPONSE']._serialized_start=1010
  _globals['_MSGCREATEPERMANENTLOCKEDACCOUNTRESPONSE']._serialized_end=1051
  _globals['_MSGCREATEPERIODICVESTINGACCOUNT']._serialized_start=1054
  _globals['_MSGCREATEPERIODICVESTINGACCOUNT']._serialized_end=1333
  _globals['_MSGCREATEPERIODICVESTINGACCOUNTRESPONSE']._serialized_start=1335
  _globals['_MSGCREATEPERIODICVESTINGACCOUNTRESPONSE']._serialized_end=1376
  _globals['_MSG']._serialized_start=1379
  _globals['_MSG']._serialized_end=1832
# @@protoc_insertion_point(module_scope)
