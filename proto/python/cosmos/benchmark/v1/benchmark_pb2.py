# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/benchmark/v1/benchmark.proto
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
    'cosmos/benchmark/v1/benchmark.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n#cosmos/benchmark/v1/benchmark.proto\x12\x13\x63osmos.benchmark.v1\"\xc0\x01\n\x02Op\x12\x12\n\x04seed\x18\x01 \x01(\x04R\x04seed\x12\x14\n\x05\x61\x63tor\x18\x02 \x01(\tR\x05\x61\x63tor\x12\x1d\n\nkey_length\x18\x03 \x01(\x04R\tkeyLength\x12!\n\x0cvalue_length\x18\x04 \x01(\x04R\x0bvalueLength\x12\x1e\n\niterations\x18\x05 \x01(\rR\niterations\x12\x16\n\x06\x64\x65lete\x18\x06 \x01(\x08R\x06\x64\x65lete\x12\x16\n\x06\x65xists\x18\x07 \x01(\x08R\x06\x65xistsB\x1eZ\x1c\x63osmossdk.io/tools/benchmarkb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.benchmark.v1.benchmark_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\034cosmossdk.io/tools/benchmark'
  _globals['_OP']._serialized_start=61
  _globals['_OP']._serialized_end=253
# @@protoc_insertion_point(module_scope)
