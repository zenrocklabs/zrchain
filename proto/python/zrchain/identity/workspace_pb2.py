# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/identity/workspace.proto
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
    'zrchain/identity/workspace.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n zrchain/identity/workspace.proto\x12\x10zrchain.identity\"\xe6\x01\n\tWorkspace\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12\x18\n\x07\x63reator\x18\x02 \x01(\tR\x07\x63reator\x12\x16\n\x06owners\x18\x03 \x03(\tR\x06owners\x12)\n\x10\x63hild_workspaces\x18\x04 \x03(\tR\x0f\x63hildWorkspaces\x12&\n\x0f\x61\x64min_policy_id\x18\x05 \x01(\x04R\radminPolicyId\x12$\n\x0esign_policy_id\x18\x06 \x01(\x04R\x0csignPolicyId\x12\x14\n\x05\x61lias\x18\x07 \x01(\tR\x05\x61liasB;Z9github.com/Zenrock-Foundation/zrchain/v5/x/identity/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.identity.workspace_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/identity/types'
  _globals['_WORKSPACE']._serialized_start=55
  _globals['_WORKSPACE']._serialized_end=285
# @@protoc_insertion_point(module_scope)
