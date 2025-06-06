# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/mint/v1beta1/query.proto
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
    'cosmos/mint/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.mint.v1beta1 import mint_pb2 as cosmos_dot_mint_dot_v1beta1_dot_mint__pb2
from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/mint/v1beta1/query.proto\x12\x13\x63osmos.mint.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1e\x63osmos/mint/v1beta1/mint.proto\x1a\x11\x61mino/amino.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\x14\n\x12QueryParamsRequest\"U\n\x13QueryParamsResponse\x12>\n\x06params\x18\x01 \x01(\x0b\x32\x1b.cosmos.mint.v1beta1.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\"\x17\n\x15QueryInflationRequest\"n\n\x16QueryInflationResponse\x12T\n\tinflation\x18\x01 \x01(\x0c\x42\x36\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\tinflation\"\x1e\n\x1cQueryAnnualProvisionsRequest\"\x84\x01\n\x1dQueryAnnualProvisionsResponse\x12\x63\n\x11\x61nnual_provisions\x18\x01 \x01(\x0c\x42\x36\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x10\x61nnualProvisions2\xc5\x03\n\x05Query\x12\x80\x01\n\x06Params\x12\'.cosmos.mint.v1beta1.QueryParamsRequest\x1a(.cosmos.mint.v1beta1.QueryParamsResponse\"#\x82\xd3\xe4\x93\x02\x1d\x12\x1b/cosmos/mint/v1beta1/params\x12\x8c\x01\n\tInflation\x12*.cosmos.mint.v1beta1.QueryInflationRequest\x1a+.cosmos.mint.v1beta1.QueryInflationResponse\"&\x82\xd3\xe4\x93\x02 \x12\x1e/cosmos/mint/v1beta1/inflation\x12\xa9\x01\n\x10\x41nnualProvisions\x12\x31.cosmos.mint.v1beta1.QueryAnnualProvisionsRequest\x1a\x32.cosmos.mint.v1beta1.QueryAnnualProvisionsResponse\".\x82\xd3\xe4\x93\x02(\x12&/cosmos/mint/v1beta1/annual_provisionsB+Z)github.com/cosmos/cosmos-sdk/x/mint/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.mint.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)github.com/cosmos/cosmos-sdk/x/mint/types'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYINFLATIONRESPONSE'].fields_by_name['inflation']._loaded_options = None
  _globals['_QUERYINFLATIONRESPONSE'].fields_by_name['inflation']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_QUERYANNUALPROVISIONSRESPONSE'].fields_by_name['annual_provisions']._loaded_options = None
  _globals['_QUERYANNUALPROVISIONSRESPONSE'].fields_by_name['annual_provisions']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\035\022\033/cosmos/mint/v1beta1/params'
  _globals['_QUERY'].methods_by_name['Inflation']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Inflation']._serialized_options = b'\202\323\344\223\002 \022\036/cosmos/mint/v1beta1/inflation'
  _globals['_QUERY'].methods_by_name['AnnualProvisions']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AnnualProvisions']._serialized_options = b'\202\323\344\223\002(\022&/cosmos/mint/v1beta1/annual_provisions'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=186
  _globals['_QUERYPARAMSREQUEST']._serialized_end=206
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=208
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=293
  _globals['_QUERYINFLATIONREQUEST']._serialized_start=295
  _globals['_QUERYINFLATIONREQUEST']._serialized_end=318
  _globals['_QUERYINFLATIONRESPONSE']._serialized_start=320
  _globals['_QUERYINFLATIONRESPONSE']._serialized_end=430
  _globals['_QUERYANNUALPROVISIONSREQUEST']._serialized_start=432
  _globals['_QUERYANNUALPROVISIONSREQUEST']._serialized_end=462
  _globals['_QUERYANNUALPROVISIONSRESPONSE']._serialized_start=465
  _globals['_QUERYANNUALPROVISIONSRESPONSE']._serialized_end=597
  _globals['_QUERY']._serialized_start=600
  _globals['_QUERY']._serialized_end=1053
# @@protoc_insertion_point(module_scope)
