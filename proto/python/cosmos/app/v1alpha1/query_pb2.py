# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/app/v1alpha1/query.proto
<<<<<<< HEAD
# Protobuf Python Version: 5.29.1
=======
# Protobuf Python Version: 5.29.0
>>>>>>> main
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
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'cosmos/app/v1alpha1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.app.v1alpha1 import config_pb2 as cosmos_dot_app_dot_v1alpha1_dot_config__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/app/v1alpha1/query.proto\x12\x13\x63osmos.app.v1alpha1\x1a cosmos/app/v1alpha1/config.proto\"\x14\n\x12QueryConfigRequest\"J\n\x13QueryConfigResponse\x12\x33\n\x06\x63onfig\x18\x01 \x01(\x0b\x32\x1b.cosmos.app.v1alpha1.ConfigR\x06\x63onfig2i\n\x05Query\x12`\n\x06\x43onfig\x12\'.cosmos.app.v1alpha1.QueryConfigRequest\x1a(.cosmos.app.v1alpha1.QueryConfigResponse\"\x03\x88\x02\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.app.v1alpha1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_QUERY'].methods_by_name['Config']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Config']._serialized_options = b'\210\002\001'
  _globals['_QUERYCONFIGREQUEST']._serialized_start=90
  _globals['_QUERYCONFIGREQUEST']._serialized_end=110
  _globals['_QUERYCONFIGRESPONSE']._serialized_start=112
  _globals['_QUERYCONFIGRESPONSE']._serialized_end=186
  _globals['_QUERY']._serialized_start=188
  _globals['_QUERY']._serialized_end=293
# @@protoc_insertion_point(module_scope)
