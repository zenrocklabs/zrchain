# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/base/node/v1beta1/query.proto
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
    'cosmos/base/node/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n$cosmos/base/node/v1beta1/query.proto\x12\x18\x63osmos.base.node.v1beta1\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x14gogoproto/gogo.proto\"\x0f\n\rConfigRequest\"\xb8\x01\n\x0e\x43onfigResponse\x12*\n\x11minimum_gas_price\x18\x01 \x01(\tR\x0fminimumGasPrice\x12.\n\x13pruning_keep_recent\x18\x02 \x01(\tR\x11pruningKeepRecent\x12)\n\x10pruning_interval\x18\x03 \x01(\tR\x0fpruningInterval\x12\x1f\n\x0bhalt_height\x18\x04 \x01(\x04R\nhaltHeight\"\x0f\n\rStatusRequest\"\xde\x01\n\x0eStatusResponse\x12\x32\n\x15\x65\x61rliest_store_height\x18\x01 \x01(\x04R\x13\x65\x61rliestStoreHeight\x12\x16\n\x06height\x18\x02 \x01(\x04R\x06height\x12>\n\ttimestamp\x18\x03 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\ttimestamp\x12\x19\n\x08\x61pp_hash\x18\x04 \x01(\x0cR\x07\x61ppHash\x12%\n\x0evalidator_hash\x18\x05 \x01(\x0cR\rvalidatorHash2\x99\x02\n\x07Service\x12\x85\x01\n\x06\x43onfig\x12\'.cosmos.base.node.v1beta1.ConfigRequest\x1a(.cosmos.base.node.v1beta1.ConfigResponse\"(\x82\xd3\xe4\x93\x02\"\x12 /cosmos/base/node/v1beta1/config\x12\x85\x01\n\x06Status\x12\'.cosmos.base.node.v1beta1.StatusRequest\x1a(.cosmos.base.node.v1beta1.StatusResponse\"(\x82\xd3\xe4\x93\x02\"\x12 /cosmos/base/node/v1beta1/statusB/Z-github.com/cosmos/cosmos-sdk/client/grpc/nodeb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.base.node.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z-github.com/cosmos/cosmos-sdk/client/grpc/node'
  _globals['_STATUSRESPONSE'].fields_by_name['timestamp']._loaded_options = None
  _globals['_STATUSRESPONSE'].fields_by_name['timestamp']._serialized_options = b'\220\337\037\001'
  _globals['_SERVICE'].methods_by_name['Config']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['Config']._serialized_options = b'\202\323\344\223\002\"\022 /cosmos/base/node/v1beta1/config'
  _globals['_SERVICE'].methods_by_name['Status']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['Status']._serialized_options = b'\202\323\344\223\002\"\022 /cosmos/base/node/v1beta1/status'
  _globals['_CONFIGREQUEST']._serialized_start=151
  _globals['_CONFIGREQUEST']._serialized_end=166
  _globals['_CONFIGRESPONSE']._serialized_start=169
  _globals['_CONFIGRESPONSE']._serialized_end=353
  _globals['_STATUSREQUEST']._serialized_start=355
  _globals['_STATUSREQUEST']._serialized_end=370
  _globals['_STATUSRESPONSE']._serialized_start=373
  _globals['_STATUSRESPONSE']._serialized_end=595
  _globals['_SERVICE']._serialized_start=598
  _globals['_SERVICE']._serialized_end=879
# @@protoc_insertion_point(module_scope)
