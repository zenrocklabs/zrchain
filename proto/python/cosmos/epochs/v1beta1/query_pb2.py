# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/epochs/v1beta1/query.proto
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
    'cosmos/epochs/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.epochs.v1beta1 import genesis_pb2 as cosmos_dot_epochs_dot_v1beta1_dot_genesis__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n!cosmos/epochs/v1beta1/query.proto\x12\x15\x63osmos.epochs.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a#cosmos/epochs/v1beta1/genesis.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\x18\n\x16QueryEpochInfosRequest\"Y\n\x17QueryEpochInfosResponse\x12>\n\x06\x65pochs\x18\x01 \x03(\x0b\x32 .cosmos.epochs.v1beta1.EpochInfoB\x04\xc8\xde\x1f\x00R\x06\x65pochs\":\n\x18QueryCurrentEpochRequest\x12\x1e\n\nidentifier\x18\x01 \x01(\tR\nidentifier\"@\n\x19QueryCurrentEpochResponse\x12#\n\rcurrent_epoch\x18\x01 \x01(\x03R\x0c\x63urrentEpoch2\xe4\x02\n\x05Query\x12\xa5\x01\n\nEpochInfos\x12-.cosmos.epochs.v1beta1.QueryEpochInfosRequest\x1a..cosmos.epochs.v1beta1.QueryEpochInfosResponse\"8\xca\xb4-\x0f\x63osmos-sdk 0.53\x82\xd3\xe4\x93\x02\x1f\x12\x1d/cosmos/epochs/v1beta1/epochs\x12\xb2\x01\n\x0c\x43urrentEpoch\x12/.cosmos.epochs.v1beta1.QueryCurrentEpochRequest\x1a\x30.cosmos.epochs.v1beta1.QueryCurrentEpochResponse\"?\xca\xb4-\x0f\x63osmos-sdk 0.53\x82\xd3\xe4\x93\x02&\x12$/cosmos/epochs/v1beta1/current_epochB-Z+github.com/cosmos/cosmos-sdk/x/epochs/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.epochs.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+github.com/cosmos/cosmos-sdk/x/epochs/types'
  _globals['_QUERYEPOCHINFOSRESPONSE'].fields_by_name['epochs']._loaded_options = None
  _globals['_QUERYEPOCHINFOSRESPONSE'].fields_by_name['epochs']._serialized_options = b'\310\336\037\000'
  _globals['_QUERY'].methods_by_name['EpochInfos']._loaded_options = None
  _globals['_QUERY'].methods_by_name['EpochInfos']._serialized_options = b'\312\264-\017cosmos-sdk 0.53\202\323\344\223\002\037\022\035/cosmos/epochs/v1beta1/epochs'
  _globals['_QUERY'].methods_by_name['CurrentEpoch']._loaded_options = None
  _globals['_QUERY'].methods_by_name['CurrentEpoch']._serialized_options = b'\312\264-\017cosmos-sdk 0.53\202\323\344\223\002&\022$/cosmos/epochs/v1beta1/current_epoch'
  _globals['_QUERYEPOCHINFOSREQUEST']._serialized_start=176
  _globals['_QUERYEPOCHINFOSREQUEST']._serialized_end=200
  _globals['_QUERYEPOCHINFOSRESPONSE']._serialized_start=202
  _globals['_QUERYEPOCHINFOSRESPONSE']._serialized_end=291
  _globals['_QUERYCURRENTEPOCHREQUEST']._serialized_start=293
  _globals['_QUERYCURRENTEPOCHREQUEST']._serialized_end=351
  _globals['_QUERYCURRENTEPOCHRESPONSE']._serialized_start=353
  _globals['_QUERYCURRENTEPOCHRESPONSE']._serialized_end=417
  _globals['_QUERY']._serialized_start=420
  _globals['_QUERY']._serialized_end=776
# @@protoc_insertion_point(module_scope)
