# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/consensus/v1/query.proto
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
    'cosmos/consensus/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from tendermint.types import params_pb2 as tendermint_dot_types_dot_params__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/consensus/v1/query.proto\x12\x13\x63osmos.consensus.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x1dtendermint/types/params.proto\"\x14\n\x12QueryParamsRequest\"P\n\x13QueryParamsResponse\x12\x39\n\x06params\x18\x01 \x01(\x0b\x32!.tendermint.types.ConsensusParamsR\x06params2\x8a\x01\n\x05Query\x12\x80\x01\n\x06Params\x12\'.cosmos.consensus.v1.QueryParamsRequest\x1a(.cosmos.consensus.v1.QueryParamsResponse\"#\x82\xd3\xe4\x93\x02\x1d\x12\x1b/cosmos/consensus/v1/paramsB0Z.github.com/cosmos/cosmos-sdk/x/consensus/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.consensus.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z.github.com/cosmos/cosmos-sdk/x/consensus/types'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\035\022\033/cosmos/consensus/v1/params'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=117
  _globals['_QUERYPARAMSREQUEST']._serialized_end=137
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=139
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=219
  _globals['_QUERY']._serialized_start=222
  _globals['_QUERY']._serialized_end=360
# @@protoc_insertion_point(module_scope)