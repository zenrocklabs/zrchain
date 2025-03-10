# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/accounts/testing/counter/v1/counter.proto
# Protobuf Python Version: 6.30.0
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
    0,
    '',
    'cosmos/accounts/testing/counter/v1/counter.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n0cosmos/accounts/testing/counter/v1/counter.proto\x12\"cosmos.accounts.testing.counter.v1\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x14gogoproto/gogo.proto\".\n\x07MsgInit\x12#\n\rinitial_value\x18\x01 \x01(\x04R\x0cinitialValue\"\x11\n\x0fMsgInitResponse\",\n\x12MsgIncreaseCounter\x12\x16\n\x06\x61mount\x18\x01 \x01(\x04R\x06\x61mount\";\n\x1aMsgIncreaseCounterResponse\x12\x1d\n\nnew_amount\x18\x01 \x01(\x04R\tnewAmount\"\x15\n\x13MsgTestDependencies\"\xf1\x01\n\x1bMsgTestDependenciesResponse\x12\x19\n\x08\x63hain_id\x18\x01 \x01(\tR\x07\x63hainId\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ress\x12\x1d\n\nbefore_gas\x18\x03 \x01(\x04R\tbeforeGas\x12\x1b\n\tafter_gas\x18\x04 \x01(\x04R\x08\x61\x66terGas\x12\x61\n\x05\x66unds\x18\x05 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x05\x66unds\"\x15\n\x13QueryCounterRequest\",\n\x14QueryCounterResponse\x12\x14\n\x05value\x18\x01 \x01(\x04R\x05valueB,Z*cosmossdk.io/x/accounts/testing/counter/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.accounts.testing.counter.v1.counter_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z*cosmossdk.io/x/accounts/testing/counter/v1'
  _globals['_MSGTESTDEPENDENCIESRESPONSE'].fields_by_name['funds']._loaded_options = None
  _globals['_MSGTESTDEPENDENCIESRESPONSE'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGINIT']._serialized_start=142
  _globals['_MSGINIT']._serialized_end=188
  _globals['_MSGINITRESPONSE']._serialized_start=190
  _globals['_MSGINITRESPONSE']._serialized_end=207
  _globals['_MSGINCREASECOUNTER']._serialized_start=209
  _globals['_MSGINCREASECOUNTER']._serialized_end=253
  _globals['_MSGINCREASECOUNTERRESPONSE']._serialized_start=255
  _globals['_MSGINCREASECOUNTERRESPONSE']._serialized_end=314
  _globals['_MSGTESTDEPENDENCIES']._serialized_start=316
  _globals['_MSGTESTDEPENDENCIES']._serialized_end=337
  _globals['_MSGTESTDEPENDENCIESRESPONSE']._serialized_start=340
  _globals['_MSGTESTDEPENDENCIESRESPONSE']._serialized_end=581
  _globals['_QUERYCOUNTERREQUEST']._serialized_start=583
  _globals['_QUERYCOUNTERREQUEST']._serialized_end=604
  _globals['_QUERYCOUNTERRESPONSE']._serialized_start=606
  _globals['_QUERYCOUNTERRESPONSE']._serialized_end=650
# @@protoc_insertion_point(module_scope)
