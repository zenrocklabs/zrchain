# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/protocolpool/v1/query.proto
<<<<<<< HEAD
# Protobuf Python Version: 5.29.1
=======
# Protobuf Python Version: 5.29.0
>>>>>>> main
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
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'cosmos/protocolpool/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\"cosmos/protocolpool/v1/query.proto\x12\x16\x63osmos.protocolpool.v1\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1egoogle/protobuf/duration.proto\"\x1b\n\x19QueryCommunityPoolRequest\"\x83\x01\n\x1aQueryCommunityPoolResponse\x12\x65\n\x04pool\x18\x01 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB3\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoinsR\x04pool\"Q\n\x1bQueryUnclaimedBudgetRequest\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\"\x8c\x03\n\x1cQueryUnclaimedBudgetResponse\x12<\n\x0ctotal_budget\x18\x01 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinR\x0btotalBudget\x12@\n\x0e\x63laimed_amount\x18\x02 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinR\rclaimedAmount\x12\x44\n\x10unclaimed_amount\x18\x03 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinR\x0funclaimedAmount\x12H\n\x0fnext_claim_from\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\rnextClaimFrom\x12\x37\n\x06period\x18\x05 \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x06period\x12#\n\rtranches_left\x18\x06 \x01(\x04R\x0ctranchesLeft2\xeb\x02\n\x05Query\x12\xa6\x01\n\rCommunityPool\x12\x31.cosmos.protocolpool.v1.QueryCommunityPoolRequest\x1a\x32.cosmos.protocolpool.v1.QueryCommunityPoolResponse\".\x82\xd3\xe4\x93\x02(\x12&/cosmos/protocolpool/v1/community_pool\x12\xb8\x01\n\x0fUnclaimedBudget\x12\x33.cosmos.protocolpool.v1.QueryUnclaimedBudgetRequest\x1a\x34.cosmos.protocolpool.v1.QueryUnclaimedBudgetResponse\":\x82\xd3\xe4\x93\x02\x34\x12\x32/cosmos/protocolpool/v1/unclaimed_budget/{address}B#Z!cosmossdk.io/x/protocolpool/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.protocolpool.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z!cosmossdk.io/x/protocolpool/types'
  _globals['_QUERYCOMMUNITYPOOLRESPONSE'].fields_by_name['pool']._loaded_options = None
  _globals['_QUERYCOMMUNITYPOOLRESPONSE'].fields_by_name['pool']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins'
  _globals['_QUERYUNCLAIMEDBUDGETREQUEST'].fields_by_name['address']._loaded_options = None
  _globals['_QUERYUNCLAIMEDBUDGETREQUEST'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYUNCLAIMEDBUDGETRESPONSE'].fields_by_name['next_claim_from']._loaded_options = None
  _globals['_QUERYUNCLAIMEDBUDGETRESPONSE'].fields_by_name['next_claim_from']._serialized_options = b'\220\337\037\001'
  _globals['_QUERYUNCLAIMEDBUDGETRESPONSE'].fields_by_name['period']._loaded_options = None
  _globals['_QUERYUNCLAIMEDBUDGETRESPONSE'].fields_by_name['period']._serialized_options = b'\230\337\037\001'
  _globals['_QUERY'].methods_by_name['CommunityPool']._loaded_options = None
  _globals['_QUERY'].methods_by_name['CommunityPool']._serialized_options = b'\202\323\344\223\002(\022&/cosmos/protocolpool/v1/community_pool'
  _globals['_QUERY'].methods_by_name['UnclaimedBudget']._loaded_options = None
  _globals['_QUERY'].methods_by_name['UnclaimedBudget']._serialized_options = b'\202\323\344\223\0024\0222/cosmos/protocolpool/v1/unclaimed_budget/{address}'
  _globals['_QUERYCOMMUNITYPOOLREQUEST']._serialized_start=238
  _globals['_QUERYCOMMUNITYPOOLREQUEST']._serialized_end=265
  _globals['_QUERYCOMMUNITYPOOLRESPONSE']._serialized_start=268
  _globals['_QUERYCOMMUNITYPOOLRESPONSE']._serialized_end=399
  _globals['_QUERYUNCLAIMEDBUDGETREQUEST']._serialized_start=401
  _globals['_QUERYUNCLAIMEDBUDGETREQUEST']._serialized_end=482
  _globals['_QUERYUNCLAIMEDBUDGETRESPONSE']._serialized_start=485
  _globals['_QUERYUNCLAIMEDBUDGETRESPONSE']._serialized_end=881
  _globals['_QUERY']._serialized_start=884
  _globals['_QUERY']._serialized_end=1247
# @@protoc_insertion_point(module_scope)
