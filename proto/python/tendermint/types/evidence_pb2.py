# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/types/evidence.proto
# Protobuf Python Version: 5.29.0
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
    0,
    '',
    'tendermint/types/evidence.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from tendermint.types import types_pb2 as tendermint_dot_types_dot_types__pb2
from tendermint.types import validator_pb2 as tendermint_dot_types_dot_validator__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1ftendermint/types/evidence.proto\x12\x10tendermint.types\x1a\x14gogoproto/gogo.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1ctendermint/types/types.proto\x1a tendermint/types/validator.proto\"\xe4\x01\n\x08\x45vidence\x12\x61\n\x17\x64uplicate_vote_evidence\x18\x01 \x01(\x0b\x32\'.tendermint.types.DuplicateVoteEvidenceH\x00R\x15\x64uplicateVoteEvidence\x12n\n\x1clight_client_attack_evidence\x18\x02 \x01(\x0b\x32+.tendermint.types.LightClientAttackEvidenceH\x00R\x19lightClientAttackEvidenceB\x05\n\x03sum\"\x90\x02\n\x15\x44uplicateVoteEvidence\x12-\n\x06vote_a\x18\x01 \x01(\x0b\x32\x16.tendermint.types.VoteR\x05voteA\x12-\n\x06vote_b\x18\x02 \x01(\x0b\x32\x16.tendermint.types.VoteR\x05voteB\x12,\n\x12total_voting_power\x18\x03 \x01(\x03R\x10totalVotingPower\x12\'\n\x0fvalidator_power\x18\x04 \x01(\x03R\x0evalidatorPower\x12\x42\n\ttimestamp\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\"\xcd\x02\n\x19LightClientAttackEvidence\x12I\n\x11\x63onflicting_block\x18\x01 \x01(\x0b\x32\x1c.tendermint.types.LightBlockR\x10\x63onflictingBlock\x12#\n\rcommon_height\x18\x02 \x01(\x03R\x0c\x63ommonHeight\x12N\n\x14\x62yzantine_validators\x18\x03 \x03(\x0b\x32\x1b.tendermint.types.ValidatorR\x13\x62yzantineValidators\x12,\n\x12total_voting_power\x18\x04 \x01(\x03R\x10totalVotingPower\x12\x42\n\ttimestamp\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\"L\n\x0c\x45videnceList\x12<\n\x08\x65vidence\x18\x01 \x03(\x0b\x32\x1a.tendermint.types.EvidenceB\x04\xc8\xde\x1f\x00R\x08\x65videnceB5Z3github.com/cometbft/cometbft/proto/tendermint/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.types.evidence_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z3github.com/cometbft/cometbft/proto/tendermint/types'
  _globals['_DUPLICATEVOTEEVIDENCE'].fields_by_name['timestamp']._loaded_options = None
  _globals['_DUPLICATEVOTEEVIDENCE'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_LIGHTCLIENTATTACKEVIDENCE'].fields_by_name['timestamp']._loaded_options = None
  _globals['_LIGHTCLIENTATTACKEVIDENCE'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_EVIDENCELIST'].fields_by_name['evidence']._loaded_options = None
  _globals['_EVIDENCELIST'].fields_by_name['evidence']._serialized_options = b'\310\336\037\000'
  _globals['_EVIDENCE']._serialized_start=173
  _globals['_EVIDENCE']._serialized_end=401
  _globals['_DUPLICATEVOTEEVIDENCE']._serialized_start=404
  _globals['_DUPLICATEVOTEEVIDENCE']._serialized_end=676
  _globals['_LIGHTCLIENTATTACKEVIDENCE']._serialized_start=679
  _globals['_LIGHTCLIENTATTACKEVIDENCE']._serialized_end=1012
  _globals['_EVIDENCELIST']._serialized_start=1014
  _globals['_EVIDENCELIST']._serialized_end=1090
# @@protoc_insertion_point(module_scope)
