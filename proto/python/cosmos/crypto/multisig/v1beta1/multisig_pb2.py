# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/crypto/multisig/v1beta1/multisig.proto
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
    'cosmos/crypto/multisig/v1beta1/multisig.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n-cosmos/crypto/multisig/v1beta1/multisig.proto\x12\x1e\x63osmos.crypto.multisig.v1beta1\x1a\x14gogoproto/gogo.proto\"6\n\x0eMultiSignature\x12\x1e\n\nsignatures\x18\x01 \x03(\x0cR\nsignatures:\x04\xd0\xa1\x1f\x01\"Y\n\x0f\x43ompactBitArray\x12*\n\x11\x65xtra_bits_stored\x18\x01 \x01(\rR\x0f\x65xtraBitsStored\x12\x14\n\x05\x65lems\x18\x02 \x01(\x0cR\x05\x65lems:\x04\x98\xa0\x1f\x00\x42+Z)github.com/cosmos/cosmos-sdk/crypto/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.crypto.multisig.v1beta1.multisig_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)github.com/cosmos/cosmos-sdk/crypto/types'
  _globals['_MULTISIGNATURE']._loaded_options = None
  _globals['_MULTISIGNATURE']._serialized_options = b'\320\241\037\001'
  _globals['_COMPACTBITARRAY']._loaded_options = None
  _globals['_COMPACTBITARRAY']._serialized_options = b'\230\240\037\000'
  _globals['_MULTISIGNATURE']._serialized_start=103
  _globals['_MULTISIGNATURE']._serialized_end=157
  _globals['_COMPACTBITARRAY']._serialized_start=159
  _globals['_COMPACTBITARRAY']._serialized_end=248
# @@protoc_insertion_point(module_scope)
