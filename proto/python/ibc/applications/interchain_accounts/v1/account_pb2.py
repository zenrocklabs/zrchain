# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/interchain_accounts/v1/account.proto
# Protobuf Python Version: 5.29.1
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
    1,
    '',
    'ibc/applications/interchain_accounts/v1/account.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.auth.v1beta1 import auth_pb2 as cosmos_dot_auth_dot_v1beta1_dot_auth__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n5ibc/applications/interchain_accounts/v1/account.proto\x12\'ibc.applications.interchain_accounts.v1\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x1e\x63osmos/auth/v1beta1/auth.proto\"\xcb\x01\n\x11InterchainAccount\x12I\n\x0c\x62\x61se_account\x18\x01 \x01(\x0b\x32 .cosmos.auth.v1beta1.BaseAccountB\x04\xd0\xde\x1f\x01R\x0b\x62\x61seAccount\x12#\n\raccount_owner\x18\x02 \x01(\tR\x0c\x61\x63\x63ountOwner:F\x88\xa0\x1f\x00\x98\xa0\x1f\x00\xca\xb4-:ibc.applications.interchain_accounts.v1.InterchainAccountIBGZEgithub.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.interchain_accounts.v1.account_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZEgithub.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/types'
  _globals['_INTERCHAINACCOUNT'].fields_by_name['base_account']._loaded_options = None
  _globals['_INTERCHAINACCOUNT'].fields_by_name['base_account']._serialized_options = b'\320\336\037\001'
  _globals['_INTERCHAINACCOUNT']._loaded_options = None
  _globals['_INTERCHAINACCOUNT']._serialized_options = b'\210\240\037\000\230\240\037\000\312\264-:ibc.applications.interchain_accounts.v1.InterchainAccountI'
  _globals['_INTERCHAINACCOUNT']._serialized_start=180
  _globals['_INTERCHAINACCOUNT']._serialized_end=383
# @@protoc_insertion_point(module_scope)
