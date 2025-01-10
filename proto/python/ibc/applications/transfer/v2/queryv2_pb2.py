# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/transfer/v2/queryv2.proto
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
    'ibc/applications/transfer/v2/queryv2.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from ibc.applications.transfer.v2 import token_pb2 as ibc_dot_applications_dot_transfer_dot_v2_dot_token__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n*ibc/applications/transfer/v2/queryv2.proto\x12\x1cibc.applications.transfer.v2\x1a\x14gogoproto/gogo.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a(ibc/applications/transfer/v2/token.proto\x1a\x1cgoogle/api/annotations.proto\"\'\n\x11QueryDenomRequest\x12\x12\n\x04hash\x18\x01 \x01(\tR\x04hash\"O\n\x12QueryDenomResponse\x12\x39\n\x05\x64\x65nom\x18\x01 \x01(\x0b\x32#.ibc.applications.transfer.v2.DenomR\x05\x64\x65nom\"\\\n\x12QueryDenomsRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xab\x01\n\x13QueryDenomsResponse\x12K\n\x06\x64\x65noms\x18\x01 \x03(\x0b\x32#.ibc.applications.transfer.v2.DenomB\x0e\xc8\xde\x1f\x00\xaa\xdf\x1f\x06\x44\x65nomsR\x06\x64\x65noms\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination2\xbc\x02\n\x07QueryV2\x12\x93\x01\n\x06\x44\x65noms\x12\x30.ibc.applications.transfer.v2.QueryDenomsRequest\x1a\x31.ibc.applications.transfer.v2.QueryDenomsResponse\"$\x82\xd3\xe4\x93\x02\x1e\x12\x1c/ibc/apps/transfer/v2/denoms\x12\x9a\x01\n\x05\x44\x65nom\x12/.ibc.applications.transfer.v2.QueryDenomRequest\x1a\x30.ibc.applications.transfer.v2.QueryDenomResponse\".\x82\xd3\xe4\x93\x02(\x12&/ibc/apps/transfer/v2/denoms/{hash=**}B9Z7github.com/cosmos/ibc-go/v9/modules/apps/transfer/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.transfer.v2.queryv2_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/cosmos/ibc-go/v9/modules/apps/transfer/types'
  _globals['_QUERYDENOMSRESPONSE'].fields_by_name['denoms']._loaded_options = None
  _globals['_QUERYDENOMSRESPONSE'].fields_by_name['denoms']._serialized_options = b'\310\336\037\000\252\337\037\006Denoms'
  _globals['_QUERYV2'].methods_by_name['Denoms']._loaded_options = None
  _globals['_QUERYV2'].methods_by_name['Denoms']._serialized_options = b'\202\323\344\223\002\036\022\034/ibc/apps/transfer/v2/denoms'
  _globals['_QUERYV2'].methods_by_name['Denom']._loaded_options = None
  _globals['_QUERYV2'].methods_by_name['Denom']._serialized_options = b'\202\323\344\223\002(\022&/ibc/apps/transfer/v2/denoms/{hash=**}'
  _globals['_QUERYDENOMREQUEST']._serialized_start=214
  _globals['_QUERYDENOMREQUEST']._serialized_end=253
  _globals['_QUERYDENOMRESPONSE']._serialized_start=255
  _globals['_QUERYDENOMRESPONSE']._serialized_end=334
  _globals['_QUERYDENOMSREQUEST']._serialized_start=336
  _globals['_QUERYDENOMSREQUEST']._serialized_end=428
  _globals['_QUERYDENOMSRESPONSE']._serialized_start=431
  _globals['_QUERYDENOMSRESPONSE']._serialized_end=602
  _globals['_QUERYV2']._serialized_start=605
  _globals['_QUERYV2']._serialized_end=921
# @@protoc_insertion_point(module_scope)
