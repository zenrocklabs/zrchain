# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/feegrant/v1beta1/query.proto
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
    'cosmos/feegrant/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.feegrant.v1beta1 import feegrant_pb2 as cosmos_dot_feegrant_dot_v1beta1_dot_feegrant__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n#cosmos/feegrant/v1beta1/query.proto\x12\x17\x63osmos.feegrant.v1beta1\x1a&cosmos/feegrant/v1beta1/feegrant.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\x7f\n\x15QueryAllowanceRequest\x12\x32\n\x07granter\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07granter\x12\x32\n\x07grantee\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07grantee\"V\n\x16QueryAllowanceResponse\x12<\n\tallowance\x18\x01 \x01(\x0b\x32\x1e.cosmos.feegrant.v1beta1.GrantR\tallowance\"\x94\x01\n\x16QueryAllowancesRequest\x12\x32\n\x07grantee\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07grantee\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xa2\x01\n\x17QueryAllowancesResponse\x12>\n\nallowances\x18\x01 \x03(\x0b\x32\x1e.cosmos.feegrant.v1beta1.GrantR\nallowances\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"\xb2\x01\n\x1fQueryAllowancesByGranterRequest\x12\x32\n\x07granter\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07granter\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination:\x13\xd2\xb4-\x0f\x63osmos-sdk 0.46\"\xc0\x01\n QueryAllowancesByGranterResponse\x12>\n\nallowances\x18\x01 \x03(\x0b\x32\x1e.cosmos.feegrant.v1beta1.GrantR\nallowances\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination:\x13\xd2\xb4-\x0f\x63osmos-sdk 0.462\xb2\x04\n\x05Query\x12\xac\x01\n\tAllowance\x12..cosmos.feegrant.v1beta1.QueryAllowanceRequest\x1a/.cosmos.feegrant.v1beta1.QueryAllowanceResponse\">\x82\xd3\xe4\x93\x02\x38\x12\x36/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}\x12\xa6\x01\n\nAllowances\x12/.cosmos.feegrant.v1beta1.QueryAllowancesRequest\x1a\x30.cosmos.feegrant.v1beta1.QueryAllowancesResponse\"5\x82\xd3\xe4\x93\x02/\x12-/cosmos/feegrant/v1beta1/allowances/{grantee}\x12\xd0\x01\n\x13\x41llowancesByGranter\x12\x38.cosmos.feegrant.v1beta1.QueryAllowancesByGranterRequest\x1a\x39.cosmos.feegrant.v1beta1.QueryAllowancesByGranterResponse\"D\xca\xb4-\x0f\x63osmos-sdk 0.46\x82\xd3\xe4\x93\x02+\x12)/cosmos/feegrant/v1beta1/issued/{granter}B)Z\'github.com/cosmos/cosmos-sdk/x/feegrantb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.feegrant.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\'github.com/cosmos/cosmos-sdk/x/feegrant'
  _globals['_QUERYALLOWANCEREQUEST'].fields_by_name['granter']._loaded_options = None
  _globals['_QUERYALLOWANCEREQUEST'].fields_by_name['granter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYALLOWANCEREQUEST'].fields_by_name['grantee']._loaded_options = None
  _globals['_QUERYALLOWANCEREQUEST'].fields_by_name['grantee']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYALLOWANCESREQUEST'].fields_by_name['grantee']._loaded_options = None
  _globals['_QUERYALLOWANCESREQUEST'].fields_by_name['grantee']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYALLOWANCESBYGRANTERREQUEST'].fields_by_name['granter']._loaded_options = None
  _globals['_QUERYALLOWANCESBYGRANTERREQUEST'].fields_by_name['granter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYALLOWANCESBYGRANTERREQUEST']._loaded_options = None
  _globals['_QUERYALLOWANCESBYGRANTERREQUEST']._serialized_options = b'\322\264-\017cosmos-sdk 0.46'
  _globals['_QUERYALLOWANCESBYGRANTERRESPONSE']._loaded_options = None
  _globals['_QUERYALLOWANCESBYGRANTERRESPONSE']._serialized_options = b'\322\264-\017cosmos-sdk 0.46'
  _globals['_QUERY'].methods_by_name['Allowance']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Allowance']._serialized_options = b'\202\323\344\223\0028\0226/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}'
  _globals['_QUERY'].methods_by_name['Allowances']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Allowances']._serialized_options = b'\202\323\344\223\002/\022-/cosmos/feegrant/v1beta1/allowances/{grantee}'
  _globals['_QUERY'].methods_by_name['AllowancesByGranter']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AllowancesByGranter']._serialized_options = b'\312\264-\017cosmos-sdk 0.46\202\323\344\223\002+\022)/cosmos/feegrant/v1beta1/issued/{granter}'
  _globals['_QUERYALLOWANCEREQUEST']._serialized_start=205
  _globals['_QUERYALLOWANCEREQUEST']._serialized_end=332
  _globals['_QUERYALLOWANCERESPONSE']._serialized_start=334
  _globals['_QUERYALLOWANCERESPONSE']._serialized_end=420
  _globals['_QUERYALLOWANCESREQUEST']._serialized_start=423
  _globals['_QUERYALLOWANCESREQUEST']._serialized_end=571
  _globals['_QUERYALLOWANCESRESPONSE']._serialized_start=574
  _globals['_QUERYALLOWANCESRESPONSE']._serialized_end=736
  _globals['_QUERYALLOWANCESBYGRANTERREQUEST']._serialized_start=739
  _globals['_QUERYALLOWANCESBYGRANTERREQUEST']._serialized_end=917
  _globals['_QUERYALLOWANCESBYGRANTERRESPONSE']._serialized_start=920
  _globals['_QUERYALLOWANCESBYGRANTERRESPONSE']._serialized_end=1112
  _globals['_QUERY']._serialized_start=1115
  _globals['_QUERY']._serialized_end=1677
# @@protoc_insertion_point(module_scope)
