# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/policy/policy.proto
# Protobuf Python Version: 5.28.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    2,
    '',
    'zrchain/policy/policy.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1bzrchain/policy/policy.proto\x12\x0ezrchain.policy\x1a\x19google/protobuf/any.proto\"\x86\x01\n\x06Policy\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x0e\n\x02id\x18\x02 \x01(\x04R\x02id\x12\x12\n\x04name\x18\x03 \x01(\tR\x04name\x12,\n\x06policy\x18\x04 \x01(\x0b\x32\x14.google.protobuf.AnyR\x06policy\x12\x10\n\x03\x62tl\x18\x05 \x01(\x04R\x03\x62tl\"y\n\x10\x42oolparserPolicy\x12\x1e\n\ndefinition\x18\x01 \x01(\tR\ndefinition\x12\x45\n\x0cparticipants\x18\x02 \x03(\x0b\x32!.zrchain.policy.PolicyParticipantR\x0cparticipants\"U\n\x11PolicyParticipant\x12&\n\x0c\x61\x62\x62reviation\x18\x01 \x01(\tB\x02\x18\x01R\x0c\x61\x62\x62reviation\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ressB9Z7github.com/Zenrock-Foundation/zrchain/v4/x/policy/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.policy.policy_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/Zenrock-Foundation/zrchain/v4/x/policy/types'
  _globals['_POLICYPARTICIPANT'].fields_by_name['abbreviation']._loaded_options = None
  _globals['_POLICYPARTICIPANT'].fields_by_name['abbreviation']._serialized_options = b'\030\001'
  _globals['_POLICY']._serialized_start=75
  _globals['_POLICY']._serialized_end=209
  _globals['_BOOLPARSERPOLICY']._serialized_start=211
  _globals['_BOOLPARSERPOLICY']._serialized_end=332
  _globals['_POLICYPARTICIPANT']._serialized_start=334
  _globals['_POLICYPARTICIPANT']._serialized_end=419
# @@protoc_insertion_point(module_scope)
