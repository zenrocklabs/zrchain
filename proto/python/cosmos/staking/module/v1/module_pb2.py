# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/staking/module/v1/module.proto
# Protobuf Python Version: 5.29.3
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
    3,
    '',
    'cosmos/staking/module/v1/module.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.app.v1alpha1 import module_pb2 as cosmos_dot_app_dot_v1alpha1_dot_module__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n%cosmos/staking/module/v1/module.proto\x12\x18\x63osmos.staking.module.v1\x1a cosmos/app/v1alpha1/module.proto\"\xd7\x01\n\x06Module\x12\x1f\n\x0bhooks_order\x18\x01 \x03(\tR\nhooksOrder\x12\x1c\n\tauthority\x18\x02 \x01(\tR\tauthority\x12\x36\n\x17\x62\x65\x63h32_prefix_validator\x18\x03 \x01(\tR\x15\x62\x65\x63h32PrefixValidator\x12\x36\n\x17\x62\x65\x63h32_prefix_consensus\x18\x04 \x01(\tR\x15\x62\x65\x63h32PrefixConsensus:\x1e\xba\xc0\x96\xda\x01\x18\n\x16\x63osmossdk.io/x/stakingb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.staking.module.v1.module_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_MODULE']._loaded_options = None
  _globals['_MODULE']._serialized_options = b'\272\300\226\332\001\030\n\026cosmossdk.io/x/staking'
  _globals['_MODULE']._serialized_start=102
  _globals['_MODULE']._serialized_end=317
# @@protoc_insertion_point(module_scope)
