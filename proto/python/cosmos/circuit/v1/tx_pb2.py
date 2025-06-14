# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/circuit/v1/tx.proto
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
    'cosmos/circuit/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos.circuit.v1 import types_pb2 as cosmos_dot_circuit_dot_v1_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1a\x63osmos/circuit/v1/tx.proto\x12\x11\x63osmos.circuit.v1\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x1d\x63osmos/circuit/v1/types.proto\"\xa0\x01\n\x1aMsgAuthorizeCircuitBreaker\x12\x18\n\x07granter\x18\x01 \x01(\tR\x07granter\x12\x18\n\x07grantee\x18\x02 \x01(\tR\x07grantee\x12@\n\x0bpermissions\x18\x03 \x01(\x0b\x32\x1e.cosmos.circuit.v1.PermissionsR\x0bpermissions:\x0c\x82\xe7\xb0*\x07granter\">\n\"MsgAuthorizeCircuitBreakerResponse\x12\x18\n\x07success\x18\x01 \x01(\x08R\x07success\"i\n\x15MsgTripCircuitBreaker\x12\x1c\n\tauthority\x18\x01 \x01(\tR\tauthority\x12\"\n\rmsg_type_urls\x18\x02 \x03(\tR\x0bmsgTypeUrls:\x0e\x82\xe7\xb0*\tauthority\"9\n\x1dMsgTripCircuitBreakerResponse\x12\x18\n\x07success\x18\x01 \x01(\x08R\x07success\"j\n\x16MsgResetCircuitBreaker\x12\x1c\n\tauthority\x18\x01 \x01(\tR\tauthority\x12\"\n\rmsg_type_urls\x18\x03 \x03(\tR\x0bmsgTypeUrls:\x0e\x82\xe7\xb0*\tauthority\":\n\x1eMsgResetCircuitBreakerResponse\x12\x18\n\x07success\x18\x01 \x01(\x08R\x07success2\xf4\x02\n\x03Msg\x12\x7f\n\x17\x41uthorizeCircuitBreaker\x12-.cosmos.circuit.v1.MsgAuthorizeCircuitBreaker\x1a\x35.cosmos.circuit.v1.MsgAuthorizeCircuitBreakerResponse\x12p\n\x12TripCircuitBreaker\x12(.cosmos.circuit.v1.MsgTripCircuitBreaker\x1a\x30.cosmos.circuit.v1.MsgTripCircuitBreakerResponse\x12s\n\x13ResetCircuitBreaker\x12).cosmos.circuit.v1.MsgResetCircuitBreaker\x1a\x31.cosmos.circuit.v1.MsgResetCircuitBreakerResponse\x1a\x05\x80\xe7\xb0*\x01\x42.Z,github.com/cosmos/cosmos-sdk/x/circuit/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.circuit.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z,github.com/cosmos/cosmos-sdk/x/circuit/types'
  _globals['_MSGAUTHORIZECIRCUITBREAKER']._loaded_options = None
  _globals['_MSGAUTHORIZECIRCUITBREAKER']._serialized_options = b'\202\347\260*\007granter'
  _globals['_MSGTRIPCIRCUITBREAKER']._loaded_options = None
  _globals['_MSGTRIPCIRCUITBREAKER']._serialized_options = b'\202\347\260*\tauthority'
  _globals['_MSGRESETCIRCUITBREAKER']._loaded_options = None
  _globals['_MSGRESETCIRCUITBREAKER']._serialized_options = b'\202\347\260*\tauthority'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGAUTHORIZECIRCUITBREAKER']._serialized_start=106
  _globals['_MSGAUTHORIZECIRCUITBREAKER']._serialized_end=266
  _globals['_MSGAUTHORIZECIRCUITBREAKERRESPONSE']._serialized_start=268
  _globals['_MSGAUTHORIZECIRCUITBREAKERRESPONSE']._serialized_end=330
  _globals['_MSGTRIPCIRCUITBREAKER']._serialized_start=332
  _globals['_MSGTRIPCIRCUITBREAKER']._serialized_end=437
  _globals['_MSGTRIPCIRCUITBREAKERRESPONSE']._serialized_start=439
  _globals['_MSGTRIPCIRCUITBREAKERRESPONSE']._serialized_end=496
  _globals['_MSGRESETCIRCUITBREAKER']._serialized_start=498
  _globals['_MSGRESETCIRCUITBREAKER']._serialized_end=604
  _globals['_MSGRESETCIRCUITBREAKERRESPONSE']._serialized_start=606
  _globals['_MSGRESETCIRCUITBREAKERRESPONSE']._serialized_end=664
  _globals['_MSG']._serialized_start=667
  _globals['_MSG']._serialized_end=1039
# @@protoc_insertion_point(module_scope)
