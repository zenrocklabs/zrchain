# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/policy/additional_signature_passkey.proto
# Protobuf Python Version: 6.30.1
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
    1,
    '',
    'zrchain/policy/additional_signature_passkey.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n1zrchain/policy/additional_signature_passkey.proto\x12\x0ezrchain.policy\"\xaa\x01\n\x1a\x41\x64\x64itionalSignaturePasskey\x12\x15\n\x06raw_id\x18\x01 \x01(\x0cR\x05rawId\x12-\n\x12\x61uthenticator_data\x18\x02 \x01(\x0cR\x11\x61uthenticatorData\x12(\n\x10\x63lient_data_json\x18\x03 \x01(\x0cR\x0e\x63lientDataJson\x12\x1c\n\tsignature\x18\x04 \x01(\x0cR\tsignatureB9Z7github.com/Zenrock-Foundation/zrchain/v6/x/policy/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.policy.additional_signature_passkey_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/Zenrock-Foundation/zrchain/v6/x/policy/types'
  _globals['_ADDITIONALSIGNATUREPASSKEY']._serialized_start=70
  _globals['_ADDITIONALSIGNATUREPASSKEY']._serialized_end=240
# @@protoc_insertion_point(module_scope)
