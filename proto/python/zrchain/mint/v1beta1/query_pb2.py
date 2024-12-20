# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/mint/v1beta1/query.proto
# Protobuf Python Version: 5.29.0
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
    0,
    '',
    'zrchain/mint/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from zrchain.mint.v1beta1 import mint_pb2 as zrchain_dot_mint_dot_v1beta1_dot_mint__pb2
from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n zrchain/mint/v1beta1/query.proto\x12\x14zrchain.mint.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1fzrchain/mint/v1beta1/mint.proto\x1a\x11\x61mino/amino.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\x14\n\x12QueryParamsRequest\"V\n\x13QueryParamsResponse\x12?\n\x06params\x18\x01 \x01(\x0b\x32\x1c.zrchain.mint.v1beta1.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\"\x17\n\x15QueryInflationRequest\"n\n\x16QueryInflationResponse\x12T\n\tinflation\x18\x01 \x01(\x0c\x42\x36\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\tinflation\"\x1e\n\x1cQueryAnnualProvisionsRequest\"\x84\x01\n\x1dQueryAnnualProvisionsResponse\x12\x63\n\x11\x61nnual_provisions\x18\x01 \x01(\x0c\x42\x36\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x10\x61nnualProvisions2\xce\x03\n\x05Query\x12\x83\x01\n\x06Params\x12(.zrchain.mint.v1beta1.QueryParamsRequest\x1a).zrchain.mint.v1beta1.QueryParamsResponse\"$\x82\xd3\xe4\x93\x02\x1e\x12\x1c/zrchain/mint/v1beta1/params\x12\x8f\x01\n\tInflation\x12+.zrchain.mint.v1beta1.QueryInflationRequest\x1a,.zrchain.mint.v1beta1.QueryInflationResponse\"\'\x82\xd3\xe4\x93\x02!\x12\x1f/zrchain/mint/v1beta1/inflation\x12\xac\x01\n\x10\x41nnualProvisions\x12\x32.zrchain.mint.v1beta1.QueryAnnualProvisionsRequest\x1a\x33.zrchain.mint.v1beta1.QueryAnnualProvisionsResponse\"/\x82\xd3\xe4\x93\x02)\x12\'/zrchain/mint/v1beta1/annual_provisionsB7Z5github.com/Zenrock-Foundation/zrchain/v5/x/mint/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.mint.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/Zenrock-Foundation/zrchain/v5/x/mint/types'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYINFLATIONRESPONSE'].fields_by_name['inflation']._loaded_options = None
  _globals['_QUERYINFLATIONRESPONSE'].fields_by_name['inflation']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_QUERYANNUALPROVISIONSRESPONSE'].fields_by_name['annual_provisions']._loaded_options = None
  _globals['_QUERYANNUALPROVISIONSRESPONSE'].fields_by_name['annual_provisions']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\036\022\034/zrchain/mint/v1beta1/params'
  _globals['_QUERY'].methods_by_name['Inflation']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Inflation']._serialized_options = b'\202\323\344\223\002!\022\037/zrchain/mint/v1beta1/inflation'
  _globals['_QUERY'].methods_by_name['AnnualProvisions']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AnnualProvisions']._serialized_options = b'\202\323\344\223\002)\022\'/zrchain/mint/v1beta1/annual_provisions'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=189
  _globals['_QUERYPARAMSREQUEST']._serialized_end=209
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=211
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=297
  _globals['_QUERYINFLATIONREQUEST']._serialized_start=299
  _globals['_QUERYINFLATIONREQUEST']._serialized_end=322
  _globals['_QUERYINFLATIONRESPONSE']._serialized_start=324
  _globals['_QUERYINFLATIONRESPONSE']._serialized_end=434
  _globals['_QUERYANNUALPROVISIONSREQUEST']._serialized_start=436
  _globals['_QUERYANNUALPROVISIONSREQUEST']._serialized_end=466
  _globals['_QUERYANNUALPROVISIONSRESPONSE']._serialized_start=469
  _globals['_QUERYANNUALPROVISIONSRESPONSE']._serialized_end=601
  _globals['_QUERY']._serialized_start=604
  _globals['_QUERY']._serialized_end=1066
# @@protoc_insertion_point(module_scope)
