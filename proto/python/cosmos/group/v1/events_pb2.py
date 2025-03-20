# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/group/v1/events.proto
# Protobuf Python Version: 6.30.0
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
    0,
    '',
    'cosmos/group/v1/events.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.group.v1 import types_pb2 as cosmos_dot_group_dot_v1_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1c\x63osmos/group/v1/events.proto\x12\x0f\x63osmos.group.v1\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1b\x63osmos/group/v1/types.proto\"-\n\x10\x45ventCreateGroup\x12\x19\n\x08group_id\x18\x01 \x01(\x04R\x07groupId\"-\n\x10\x45ventUpdateGroup\x12\x19\n\x08group_id\x18\x01 \x01(\x04R\x07groupId\"L\n\x16\x45ventCreateGroupPolicy\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\"L\n\x16\x45ventUpdateGroupPolicy\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\"6\n\x13\x45ventSubmitProposal\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"8\n\x15\x45ventWithdrawProposal\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\",\n\tEventVote\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"\x81\x01\n\tEventExec\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12?\n\x06result\x18\x02 \x01(\x0e\x32\'.cosmos.group.v1.ProposalExecutorResultR\x06result\x12\x12\n\x04logs\x18\x03 \x01(\tR\x04logs\"`\n\x0f\x45ventLeaveGroup\x12\x19\n\x08group_id\x18\x01 \x01(\x04R\x07groupId\x12\x32\n\x07\x61\x64\x64ress\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\"\xb0\x01\n\x13\x45ventProposalPruned\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x37\n\x06status\x18\x02 \x01(\x0e\x32\x1f.cosmos.group.v1.ProposalStatusR\x06status\x12?\n\x0ctally_result\x18\x03 \x01(\x0b\x32\x1c.cosmos.group.v1.TallyResultR\x0btallyResultB\x16Z\x14\x63osmossdk.io/x/groupb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.group.v1.events_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\024cosmossdk.io/x/group'
  _globals['_EVENTCREATEGROUPPOLICY'].fields_by_name['address']._loaded_options = None
  _globals['_EVENTCREATEGROUPPOLICY'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_EVENTUPDATEGROUPPOLICY'].fields_by_name['address']._loaded_options = None
  _globals['_EVENTUPDATEGROUPPOLICY'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_EVENTLEAVEGROUP'].fields_by_name['address']._loaded_options = None
  _globals['_EVENTLEAVEGROUP'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_EVENTCREATEGROUP']._serialized_start=105
  _globals['_EVENTCREATEGROUP']._serialized_end=150
  _globals['_EVENTUPDATEGROUP']._serialized_start=152
  _globals['_EVENTUPDATEGROUP']._serialized_end=197
  _globals['_EVENTCREATEGROUPPOLICY']._serialized_start=199
  _globals['_EVENTCREATEGROUPPOLICY']._serialized_end=275
  _globals['_EVENTUPDATEGROUPPOLICY']._serialized_start=277
  _globals['_EVENTUPDATEGROUPPOLICY']._serialized_end=353
  _globals['_EVENTSUBMITPROPOSAL']._serialized_start=355
  _globals['_EVENTSUBMITPROPOSAL']._serialized_end=409
  _globals['_EVENTWITHDRAWPROPOSAL']._serialized_start=411
  _globals['_EVENTWITHDRAWPROPOSAL']._serialized_end=467
  _globals['_EVENTVOTE']._serialized_start=469
  _globals['_EVENTVOTE']._serialized_end=513
  _globals['_EVENTEXEC']._serialized_start=516
  _globals['_EVENTEXEC']._serialized_end=645
  _globals['_EVENTLEAVEGROUP']._serialized_start=647
  _globals['_EVENTLEAVEGROUP']._serialized_end=743
  _globals['_EVENTPROPOSALPRUNED']._serialized_start=746
  _globals['_EVENTPROPOSALPRUNED']._serialized_end=922
# @@protoc_insertion_point(module_scope)
