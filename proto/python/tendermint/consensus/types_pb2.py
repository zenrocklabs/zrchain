# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/consensus/types.proto
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
    'tendermint/consensus/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from tendermint.types import types_pb2 as tendermint_dot_types_dot_types__pb2
from tendermint.libs.bits import types_pb2 as tendermint_dot_libs_dot_bits_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n tendermint/consensus/types.proto\x12\x14tendermint.consensus\x1a\x14gogoproto/gogo.proto\x1a\x1ctendermint/types/types.proto\x1a tendermint/libs/bits/types.proto\"\xb5\x01\n\x0cNewRoundStep\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x12\n\x04step\x18\x03 \x01(\rR\x04step\x12\x37\n\x18seconds_since_start_time\x18\x04 \x01(\x03R\x15secondsSinceStartTime\x12*\n\x11last_commit_round\x18\x05 \x01(\x05R\x0flastCommitRound\"\xf5\x01\n\rNewValidBlock\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12X\n\x15\x62lock_part_set_header\x18\x03 \x01(\x0b\x32\x1f.tendermint.types.PartSetHeaderB\x04\xc8\xde\x1f\x00R\x12\x62lockPartSetHeader\x12?\n\x0b\x62lock_parts\x18\x04 \x01(\x0b\x32\x1e.tendermint.libs.bits.BitArrayR\nblockParts\x12\x1b\n\tis_commit\x18\x05 \x01(\x08R\x08isCommit\"H\n\x08Proposal\x12<\n\x08proposal\x18\x01 \x01(\x0b\x32\x1a.tendermint.types.ProposalB\x04\xc8\xde\x1f\x00R\x08proposal\"\x9c\x01\n\x0bProposalPOL\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12,\n\x12proposal_pol_round\x18\x02 \x01(\x05R\x10proposalPolRound\x12G\n\x0cproposal_pol\x18\x03 \x01(\x0b\x32\x1e.tendermint.libs.bits.BitArrayB\x04\xc8\xde\x1f\x00R\x0bproposalPol\"k\n\tBlockPart\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x30\n\x04part\x18\x03 \x01(\x0b\x32\x16.tendermint.types.PartB\x04\xc8\xde\x1f\x00R\x04part\"2\n\x04Vote\x12*\n\x04vote\x18\x01 \x01(\x0b\x32\x16.tendermint.types.VoteR\x04vote\"\x82\x01\n\x07HasVote\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x33\n\x04type\x18\x03 \x01(\x0e\x32\x1f.tendermint.types.SignedMsgTypeR\x04type\x12\x14\n\x05index\x18\x04 \x01(\x05R\x05index\"\xb8\x01\n\x0cVoteSetMaj23\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x33\n\x04type\x18\x03 \x01(\x0e\x32\x1f.tendermint.types.SignedMsgTypeR\x04type\x12\x45\n\x08\x62lock_id\x18\x04 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\"\xf3\x01\n\x0bVoteSetBits\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12\x33\n\x04type\x18\x03 \x01(\x0e\x32\x1f.tendermint.types.SignedMsgTypeR\x04type\x12\x45\n\x08\x62lock_id\x18\x04 \x01(\x0b\x32\x19.tendermint.types.BlockIDB\x0f\xc8\xde\x1f\x00\xe2\xde\x1f\x07\x42lockIDR\x07\x62lockId\x12:\n\x05votes\x18\x05 \x01(\x0b\x32\x1e.tendermint.libs.bits.BitArrayB\x04\xc8\xde\x1f\x00R\x05votes\"\xf6\x04\n\x07Message\x12J\n\x0enew_round_step\x18\x01 \x01(\x0b\x32\".tendermint.consensus.NewRoundStepH\x00R\x0cnewRoundStep\x12M\n\x0fnew_valid_block\x18\x02 \x01(\x0b\x32#.tendermint.consensus.NewValidBlockH\x00R\rnewValidBlock\x12<\n\x08proposal\x18\x03 \x01(\x0b\x32\x1e.tendermint.consensus.ProposalH\x00R\x08proposal\x12\x46\n\x0cproposal_pol\x18\x04 \x01(\x0b\x32!.tendermint.consensus.ProposalPOLH\x00R\x0bproposalPol\x12@\n\nblock_part\x18\x05 \x01(\x0b\x32\x1f.tendermint.consensus.BlockPartH\x00R\tblockPart\x12\x30\n\x04vote\x18\x06 \x01(\x0b\x32\x1a.tendermint.consensus.VoteH\x00R\x04vote\x12:\n\x08has_vote\x18\x07 \x01(\x0b\x32\x1d.tendermint.consensus.HasVoteH\x00R\x07hasVote\x12J\n\x0evote_set_maj23\x18\x08 \x01(\x0b\x32\".tendermint.consensus.VoteSetMaj23H\x00R\x0cvoteSetMaj23\x12G\n\rvote_set_bits\x18\t \x01(\x0b\x32!.tendermint.consensus.VoteSetBitsH\x00R\x0bvoteSetBitsB\x05\n\x03sumB9Z7github.com/cometbft/cometbft/proto/tendermint/consensusb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.consensus.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/cometbft/cometbft/proto/tendermint/consensus'
  _globals['_NEWVALIDBLOCK'].fields_by_name['block_part_set_header']._loaded_options = None
  _globals['_NEWVALIDBLOCK'].fields_by_name['block_part_set_header']._serialized_options = b'\310\336\037\000'
  _globals['_PROPOSAL'].fields_by_name['proposal']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['proposal']._serialized_options = b'\310\336\037\000'
  _globals['_PROPOSALPOL'].fields_by_name['proposal_pol']._loaded_options = None
  _globals['_PROPOSALPOL'].fields_by_name['proposal_pol']._serialized_options = b'\310\336\037\000'
  _globals['_BLOCKPART'].fields_by_name['part']._loaded_options = None
  _globals['_BLOCKPART'].fields_by_name['part']._serialized_options = b'\310\336\037\000'
  _globals['_VOTESETMAJ23'].fields_by_name['block_id']._loaded_options = None
  _globals['_VOTESETMAJ23'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_VOTESETBITS'].fields_by_name['block_id']._loaded_options = None
  _globals['_VOTESETBITS'].fields_by_name['block_id']._serialized_options = b'\310\336\037\000\342\336\037\007BlockID'
  _globals['_VOTESETBITS'].fields_by_name['votes']._loaded_options = None
  _globals['_VOTESETBITS'].fields_by_name['votes']._serialized_options = b'\310\336\037\000'
  _globals['_NEWROUNDSTEP']._serialized_start=145
  _globals['_NEWROUNDSTEP']._serialized_end=326
  _globals['_NEWVALIDBLOCK']._serialized_start=329
  _globals['_NEWVALIDBLOCK']._serialized_end=574
  _globals['_PROPOSAL']._serialized_start=576
  _globals['_PROPOSAL']._serialized_end=648
  _globals['_PROPOSALPOL']._serialized_start=651
  _globals['_PROPOSALPOL']._serialized_end=807
  _globals['_BLOCKPART']._serialized_start=809
  _globals['_BLOCKPART']._serialized_end=916
  _globals['_VOTE']._serialized_start=918
  _globals['_VOTE']._serialized_end=968
  _globals['_HASVOTE']._serialized_start=971
  _globals['_HASVOTE']._serialized_end=1101
  _globals['_VOTESETMAJ23']._serialized_start=1104
  _globals['_VOTESETMAJ23']._serialized_end=1288
  _globals['_VOTESETBITS']._serialized_start=1291
  _globals['_VOTESETBITS']._serialized_end=1534
  _globals['_MESSAGE']._serialized_start=1537
  _globals['_MESSAGE']._serialized_end=2167
# @@protoc_insertion_point(module_scope)
