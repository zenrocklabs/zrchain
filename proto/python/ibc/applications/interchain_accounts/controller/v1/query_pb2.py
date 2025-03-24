# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/interchain_accounts/controller/v1/query.proto
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
    'ibc/applications/interchain_accounts/controller/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from ibc.applications.interchain_accounts.controller.v1 import controller_pb2 as ibc_dot_applications_dot_interchain__accounts_dot_controller_dot_v1_dot_controller__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n>ibc/applications/interchain_accounts/controller/v1/query.proto\x12\x32ibc.applications.interchain_accounts.controller.v1\x1a\x43ibc/applications/interchain_accounts/controller/v1/controller.proto\x1a\x1cgoogle/api/annotations.proto\"Z\n\x1dQueryInterchainAccountRequest\x12\x14\n\x05owner\x18\x01 \x01(\tR\x05owner\x12#\n\rconnection_id\x18\x02 \x01(\tR\x0c\x63onnectionId\":\n\x1eQueryInterchainAccountResponse\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\"\x14\n\x12QueryParamsRequest\"i\n\x13QueryParamsResponse\x12R\n\x06params\x18\x01 \x01(\x0b\x32:.ibc.applications.interchain_accounts.controller.v1.ParamsR\x06params2\xfc\x03\n\x05Query\x12\x9a\x02\n\x11InterchainAccount\x12Q.ibc.applications.interchain_accounts.controller.v1.QueryInterchainAccountRequest\x1aR.ibc.applications.interchain_accounts.controller.v1.QueryInterchainAccountResponse\"^\x82\xd3\xe4\x93\x02X\x12V/ibc/apps/interchain_accounts/controller/v1/owners/{owner}/connections/{connection_id}\x12\xd5\x01\n\x06Params\x12\x46.ibc.applications.interchain_accounts.controller.v1.QueryParamsRequest\x1aG.ibc.applications.interchain_accounts.controller.v1.QueryParamsResponse\":\x82\xd3\xe4\x93\x02\x34\x12\x32/ibc/apps/interchain_accounts/controller/v1/paramsBSZQgithub.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.interchain_accounts.controller.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZQgithub.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/types'
  _globals['_QUERY'].methods_by_name['InterchainAccount']._loaded_options = None
  _globals['_QUERY'].methods_by_name['InterchainAccount']._serialized_options = b'\202\323\344\223\002X\022V/ibc/apps/interchain_accounts/controller/v1/owners/{owner}/connections/{connection_id}'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\0024\0222/ibc/apps/interchain_accounts/controller/v1/params'
  _globals['_QUERYINTERCHAINACCOUNTREQUEST']._serialized_start=217
  _globals['_QUERYINTERCHAINACCOUNTREQUEST']._serialized_end=307
  _globals['_QUERYINTERCHAINACCOUNTRESPONSE']._serialized_start=309
  _globals['_QUERYINTERCHAINACCOUNTRESPONSE']._serialized_end=367
  _globals['_QUERYPARAMSREQUEST']._serialized_start=369
  _globals['_QUERYPARAMSREQUEST']._serialized_end=389
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=391
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=496
  _globals['_QUERY']._serialized_start=499
  _globals['_QUERY']._serialized_end=1007
# @@protoc_insertion_point(module_scope)
