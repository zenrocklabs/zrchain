# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: google/api/expr/v1beta1/eval.proto
<<<<<<< HEAD
# Protobuf Python Version: 6.30.1
=======
# Protobuf Python Version: 6.30.0
>>>>>>> main
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
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'google/api/expr/v1beta1/eval.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api.expr.v1beta1 import value_pb2 as google_dot_api_dot_expr_dot_v1beta1_dot_value__pb2
from google.rpc import status_pb2 as google_dot_rpc_dot_status__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\"google/api/expr/v1beta1/eval.proto\x12\x17google.api.expr.v1beta1\x1a#google/api/expr/v1beta1/value.proto\x1a\x17google/rpc/status.proto\"\xe0\x01\n\tEvalState\x12:\n\x06values\x18\x01 \x03(\x0b\x32\".google.api.expr.v1beta1.ExprValueR\x06values\x12\x43\n\x07results\x18\x03 \x03(\x0b\x32).google.api.expr.v1beta1.EvalState.ResultR\x07results\x1aR\n\x06Result\x12\x32\n\x04\x65xpr\x18\x01 \x01(\x0b\x32\x1e.google.api.expr.v1beta1.IdRefR\x04\x65xpr\x12\x14\n\x05value\x18\x02 \x01(\x05R\x05value\"\xc7\x01\n\tExprValue\x12\x36\n\x05value\x18\x01 \x01(\x0b\x32\x1e.google.api.expr.v1beta1.ValueH\x00R\x05value\x12\x39\n\x05\x65rror\x18\x02 \x01(\x0b\x32!.google.api.expr.v1beta1.ErrorSetH\x00R\x05\x65rror\x12?\n\x07unknown\x18\x03 \x01(\x0b\x32#.google.api.expr.v1beta1.UnknownSetH\x00R\x07unknownB\x06\n\x04kind\"6\n\x08\x45rrorSet\x12*\n\x06\x65rrors\x18\x01 \x03(\x0b\x32\x12.google.rpc.StatusR\x06\x65rrors\"B\n\nUnknownSet\x12\x34\n\x05\x65xprs\x18\x01 \x03(\x0b\x32\x1e.google.api.expr.v1beta1.IdRefR\x05\x65xprs\"\x17\n\x05IdRef\x12\x0e\n\x02id\x18\x01 \x01(\x05R\x02idBj\n\x1b\x63om.google.api.expr.v1beta1B\tEvalProtoP\x01Z;google.golang.org/genproto/googleapis/api/expr/v1beta1;expr\xf8\x01\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'google.api.expr.v1beta1.eval_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\033com.google.api.expr.v1beta1B\tEvalProtoP\001Z;google.golang.org/genproto/googleapis/api/expr/v1beta1;expr\370\001\001'
  _globals['_EVALSTATE']._serialized_start=126
  _globals['_EVALSTATE']._serialized_end=350
  _globals['_EVALSTATE_RESULT']._serialized_start=268
  _globals['_EVALSTATE_RESULT']._serialized_end=350
  _globals['_EXPRVALUE']._serialized_start=353
  _globals['_EXPRVALUE']._serialized_end=552
  _globals['_ERRORSET']._serialized_start=554
  _globals['_ERRORSET']._serialized_end=608
  _globals['_UNKNOWNSET']._serialized_start=610
  _globals['_UNKNOWNSET']._serialized_end=676
  _globals['_IDREF']._serialized_start=678
  _globals['_IDREF']._serialized_end=701
# @@protoc_insertion_point(module_scope)
