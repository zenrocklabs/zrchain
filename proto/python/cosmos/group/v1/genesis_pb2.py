# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/group/v1/genesis.proto
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
    'cosmos/group/v1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.group.v1 import types_pb2 as cosmos_dot_group_dot_v1_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1d\x63osmos/group/v1/genesis.proto\x12\x0f\x63osmos.group.v1\x1a\x1b\x63osmos/group/v1/types.proto\"\x9e\x03\n\x0cGenesisState\x12\x1b\n\tgroup_seq\x18\x01 \x01(\x04R\x08groupSeq\x12\x32\n\x06groups\x18\x02 \x03(\x0b\x32\x1a.cosmos.group.v1.GroupInfoR\x06groups\x12\x41\n\rgroup_members\x18\x03 \x03(\x0b\x32\x1c.cosmos.group.v1.GroupMemberR\x0cgroupMembers\x12(\n\x10group_policy_seq\x18\x04 \x01(\x04R\x0egroupPolicySeq\x12G\n\x0egroup_policies\x18\x05 \x03(\x0b\x32 .cosmos.group.v1.GroupPolicyInfoR\rgroupPolicies\x12!\n\x0cproposal_seq\x18\x06 \x01(\x04R\x0bproposalSeq\x12\x37\n\tproposals\x18\x07 \x03(\x0b\x32\x19.cosmos.group.v1.ProposalR\tproposals\x12+\n\x05votes\x18\x08 \x03(\x0b\x32\x15.cosmos.group.v1.VoteR\x05votesB\x16Z\x14\x63osmossdk.io/x/groupb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.group.v1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\024cosmossdk.io/x/group'
  _globals['_GENESISSTATE']._serialized_start=80
  _globals['_GENESISSTATE']._serialized_end=494
# @@protoc_insertion_point(module_scope)
