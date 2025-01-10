# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/identity/query.proto
# Protobuf Python Version: 5.29.3
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
    3,
    '',
    'zrchain/identity/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from zrchain.identity import params_pb2 as zrchain_dot_identity_dot_params__pb2
from zrchain.identity import workspace_pb2 as zrchain_dot_identity_dot_workspace__pb2
from zrchain.identity import keyring_pb2 as zrchain_dot_identity_dot_keyring__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1czrchain/identity/query.proto\x12\x10zrchain.identity\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x1dzrchain/identity/params.proto\x1a zrchain/identity/workspace.proto\x1a\x1ezrchain/identity/keyring.proto\"\x14\n\x12QueryParamsRequest\"R\n\x13QueryParamsResponse\x12;\n\x06params\x18\x01 \x01(\x0b\x32\x18.zrchain.identity.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\"\x90\x01\n\x16QueryWorkspacesRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\x12\x14\n\x05owner\x18\x02 \x01(\tR\x05owner\x12\x18\n\x07\x63reator\x18\x03 \x01(\tR\x07\x63reator\"\x9f\x01\n\x17QueryWorkspacesResponse\x12;\n\nworkspaces\x18\x01 \x03(\x0b\x32\x1b.zrchain.identity.WorkspaceR\nworkspaces\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"G\n\x1eQueryWorkspaceByAddressRequest\x12%\n\x0eworkspace_addr\x18\x01 \x01(\tR\rworkspaceAddr\"\\\n\x1fQueryWorkspaceByAddressResponse\x12\x39\n\tworkspace\x18\x01 \x01(\x0b\x32\x1b.zrchain.identity.WorkspaceR\tworkspace\"^\n\x14QueryKeyringsRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x97\x01\n\x15QueryKeyringsResponse\x12\x35\n\x08keyrings\x18\x01 \x03(\x0b\x32\x19.zrchain.identity.KeyringR\x08keyrings\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"A\n\x1cQueryKeyringByAddressRequest\x12!\n\x0ckeyring_addr\x18\x01 \x01(\tR\x0bkeyringAddr\"T\n\x1dQueryKeyringByAddressResponse\x12\x33\n\x07keyring\x18\x01 \x01(\x0b\x32\x19.zrchain.identity.KeyringR\x07keyring2\xfb\x05\n\x05Query\x12w\n\x06Params\x12$.zrchain.identity.QueryParamsRequest\x1a%.zrchain.identity.QueryParamsResponse\" \x82\xd3\xe4\x93\x02\x1a\x12\x18/zrchain/identity/params\x12\x87\x01\n\nWorkspaces\x12(.zrchain.identity.QueryWorkspacesRequest\x1a).zrchain.identity.QueryWorkspacesResponse\"$\x82\xd3\xe4\x93\x02\x1e\x12\x1c/zrchain/identity/workspaces\x12\xba\x01\n\x12WorkspaceByAddress\x12\x30.zrchain.identity.QueryWorkspaceByAddressRequest\x1a\x31.zrchain.identity.QueryWorkspaceByAddressResponse\"?\x82\xd3\xe4\x93\x02\x39\x12\x37/zrchain/identity/workspace_by_address/{workspace_addr}\x12\x7f\n\x08Keyrings\x12&.zrchain.identity.QueryKeyringsRequest\x1a\'.zrchain.identity.QueryKeyringsResponse\"\"\x82\xd3\xe4\x93\x02\x1c\x12\x1a/zrchain/identity/keyrings\x12\xb0\x01\n\x10KeyringByAddress\x12..zrchain.identity.QueryKeyringByAddressRequest\x1a/.zrchain.identity.QueryKeyringByAddressResponse\";\x82\xd3\xe4\x93\x02\x35\x12\x33/zrchain/identity/keyring_by_address/{keyring_addr}B;Z9github.com/Zenrock-Foundation/zrchain/v5/x/identity/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.identity.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/identity/types'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\032\022\030/zrchain/identity/params'
  _globals['_QUERY'].methods_by_name['Workspaces']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Workspaces']._serialized_options = b'\202\323\344\223\002\036\022\034/zrchain/identity/workspaces'
  _globals['_QUERY'].methods_by_name['WorkspaceByAddress']._loaded_options = None
  _globals['_QUERY'].methods_by_name['WorkspaceByAddress']._serialized_options = b'\202\323\344\223\0029\0227/zrchain/identity/workspace_by_address/{workspace_addr}'
  _globals['_QUERY'].methods_by_name['Keyrings']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Keyrings']._serialized_options = b'\202\323\344\223\002\034\022\032/zrchain/identity/keyrings'
  _globals['_QUERY'].methods_by_name['KeyringByAddress']._loaded_options = None
  _globals['_QUERY'].methods_by_name['KeyringByAddress']._serialized_options = b'\202\323\344\223\0025\0223/zrchain/identity/keyring_by_address/{keyring_addr}'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=262
  _globals['_QUERYPARAMSREQUEST']._serialized_end=282
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=284
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=366
  _globals['_QUERYWORKSPACESREQUEST']._serialized_start=369
  _globals['_QUERYWORKSPACESREQUEST']._serialized_end=513
  _globals['_QUERYWORKSPACESRESPONSE']._serialized_start=516
  _globals['_QUERYWORKSPACESRESPONSE']._serialized_end=675
  _globals['_QUERYWORKSPACEBYADDRESSREQUEST']._serialized_start=677
  _globals['_QUERYWORKSPACEBYADDRESSREQUEST']._serialized_end=748
  _globals['_QUERYWORKSPACEBYADDRESSRESPONSE']._serialized_start=750
  _globals['_QUERYWORKSPACEBYADDRESSRESPONSE']._serialized_end=842
  _globals['_QUERYKEYRINGSREQUEST']._serialized_start=844
  _globals['_QUERYKEYRINGSREQUEST']._serialized_end=938
  _globals['_QUERYKEYRINGSRESPONSE']._serialized_start=941
  _globals['_QUERYKEYRINGSRESPONSE']._serialized_end=1092
  _globals['_QUERYKEYRINGBYADDRESSREQUEST']._serialized_start=1094
  _globals['_QUERYKEYRINGBYADDRESSREQUEST']._serialized_end=1159
  _globals['_QUERYKEYRINGBYADDRESSRESPONSE']._serialized_start=1161
  _globals['_QUERYKEYRINGBYADDRESSRESPONSE']._serialized_end=1245
  _globals['_QUERY']._serialized_start=1248
  _globals['_QUERY']._serialized_end=2011
# @@protoc_insertion_point(module_scope)
