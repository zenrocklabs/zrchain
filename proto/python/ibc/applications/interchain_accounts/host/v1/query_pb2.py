# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/interchain_accounts/host/v1/query.proto
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
    'ibc/applications/interchain_accounts/host/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from ibc.applications.interchain_accounts.host.v1 import host_pb2 as ibc_dot_applications_dot_interchain__accounts_dot_host_dot_v1_dot_host__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n8ibc/applications/interchain_accounts/host/v1/query.proto\x12,ibc.applications.interchain_accounts.host.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x37ibc/applications/interchain_accounts/host/v1/host.proto\"\x14\n\x12QueryParamsRequest\"c\n\x13QueryParamsResponse\x12L\n\x06params\x18\x01 \x01(\x0b\x32\x34.ibc.applications.interchain_accounts.host.v1.ParamsR\x06params2\xcd\x01\n\x05Query\x12\xc3\x01\n\x06Params\x12@.ibc.applications.interchain_accounts.host.v1.QueryParamsRequest\x1a\x41.ibc.applications.interchain_accounts.host.v1.QueryParamsResponse\"4\x82\xd3\xe4\x93\x02.\x12,/ibc/apps/interchain_accounts/host/v1/paramsBMZKgithub.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.interchain_accounts.host.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZKgithub.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002.\022,/ibc/apps/interchain_accounts/host/v1/params'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=193
  _globals['_QUERYPARAMSREQUEST']._serialized_end=213
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=215
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=314
  _globals['_QUERY']._serialized_start=317
  _globals['_QUERY']._serialized_end=522
# @@protoc_insertion_point(module_scope)
