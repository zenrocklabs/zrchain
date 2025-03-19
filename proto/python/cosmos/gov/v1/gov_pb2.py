# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/gov/v1/gov.proto
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
    'cosmos/gov/v1/gov.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x17\x63osmos/gov/v1/gov.proto\x12\rcosmos.gov.v1\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x14gogoproto/gogo.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x19google/protobuf/any.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x11\x61mino/amino.proto\"o\n\x12WeightedVoteOption\x12\x31\n\x06option\x18\x01 \x01(\x0e\x32\x19.cosmos.gov.v1.VoteOptionR\x06option\x12&\n\x06weight\x18\x02 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x06weight\"\xa0\x01\n\x07\x44\x65posit\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x36\n\tdepositor\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor\x12<\n\x06\x61mount\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06\x61mount\"\xca\x06\n\x08Proposal\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x30\n\x08messages\x18\x02 \x03(\x0b\x32\x14.google.protobuf.AnyR\x08messages\x12\x35\n\x06status\x18\x03 \x01(\x0e\x32\x1d.cosmos.gov.v1.ProposalStatusR\x06status\x12H\n\x12\x66inal_tally_result\x18\x04 \x01(\x0b\x32\x1a.cosmos.gov.v1.TallyResultR\x10\x66inalTallyResult\x12\x41\n\x0bsubmit_time\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\nsubmitTime\x12J\n\x10\x64\x65posit_end_time\x18\x06 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\x0e\x64\x65positEndTime\x12I\n\rtotal_deposit\x18\x07 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0ctotalDeposit\x12L\n\x11voting_start_time\x18\x08 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\x0fvotingStartTime\x12H\n\x0fvoting_end_time\x18\t \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\rvotingEndTime\x12\x1a\n\x08metadata\x18\n \x01(\tR\x08metadata\x12\x14\n\x05title\x18\x0b \x01(\tR\x05title\x12\x18\n\x07summary\x18\x0c \x01(\tR\x07summary\x12\x34\n\x08proposer\x18\r \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08proposer\x12 \n\texpedited\x18\x0e \x01(\x08\x42\x02\x18\x01R\texpedited\x12#\n\rfailed_reason\x18\x0f \x01(\tR\x0c\x66\x61iledReason\x12@\n\rproposal_type\x18\x10 \x01(\x0e\x32\x1b.cosmos.gov.v1.ProposalTypeR\x0cproposalType\"\xb8\x01\n\x13ProposalVoteOptions\x12\x1d\n\noption_one\x18\x01 \x01(\tR\toptionOne\x12\x1d\n\noption_two\x18\x02 \x01(\tR\toptionTwo\x12!\n\x0coption_three\x18\x03 \x01(\tR\x0boptionThree\x12\x1f\n\x0boption_four\x18\x04 \x01(\tR\noptionFour\x12\x1f\n\x0boption_spam\x18\x05 \x01(\tR\noptionSpam\"\xfc\x03\n\x0bTallyResult\x12-\n\tyes_count\x18\x01 \x01(\tB\x10\x18\x01\xd2\xb4-\ncosmos.IntR\x08yesCount\x12\x35\n\rabstain_count\x18\x02 \x01(\tB\x10\x18\x01\xd2\xb4-\ncosmos.IntR\x0c\x61\x62stainCount\x12+\n\x08no_count\x18\x03 \x01(\tB\x10\x18\x01\xd2\xb4-\ncosmos.IntR\x07noCount\x12=\n\x12no_with_veto_count\x18\x04 \x01(\tB\x10\x18\x01\xd2\xb4-\ncosmos.IntR\x0fnoWithVetoCount\x12\x38\n\x10option_one_count\x18\x05 \x01(\tB\x0e\xd2\xb4-\ncosmos.IntR\x0eoptionOneCount\x12\x38\n\x10option_two_count\x18\x06 \x01(\tB\x0e\xd2\xb4-\ncosmos.IntR\x0eoptionTwoCount\x12<\n\x12option_three_count\x18\x07 \x01(\tB\x0e\xd2\xb4-\ncosmos.IntR\x10optionThreeCount\x12:\n\x11option_four_count\x18\x08 \x01(\tB\x0e\xd2\xb4-\ncosmos.IntR\x0foptionFourCount\x12-\n\nspam_count\x18\t \x01(\tB\x0e\xd2\xb4-\ncosmos.IntR\tspamCount\"\xb6\x01\n\x04Vote\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter\x12;\n\x07options\x18\x04 \x03(\x0b\x32!.cosmos.gov.v1.WeightedVoteOptionR\x07options\x12\x1a\n\x08metadata\x18\x05 \x01(\tR\x08metadataJ\x04\x08\x03\x10\x04\"\xdd\x01\n\rDepositParams\x12Y\n\x0bmin_deposit\x18\x01 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\x1d\xc8\xde\x1f\x00\xea\xde\x1f\x15min_deposit,omitemptyR\nminDeposit\x12m\n\x12max_deposit_period\x18\x02 \x01(\x0b\x32\x19.google.protobuf.DurationB$\xea\xde\x1f\x1cmax_deposit_period,omitempty\x98\xdf\x1f\x01R\x10maxDepositPeriod:\x02\x18\x01\"X\n\x0cVotingParams\x12\x44\n\rvoting_period\x18\x01 \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x0cvotingPeriod:\x02\x18\x01\"\x9e\x01\n\x0bTallyParams\x12&\n\x06quorum\x18\x01 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x06quorum\x12,\n\tthreshold\x18\x02 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\tthreshold\x12\x35\n\x0eveto_threshold\x18\x03 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\rvetoThreshold:\x02\x18\x01\"\xc1\n\n\x06Params\x12\x45\n\x0bmin_deposit\x18\x01 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\nminDeposit\x12M\n\x12max_deposit_period\x18\x02 \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x10maxDepositPeriod\x12\x44\n\rvoting_period\x18\x03 \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x0cvotingPeriod\x12&\n\x06quorum\x18\x04 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x06quorum\x12,\n\tthreshold\x18\x05 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\tthreshold\x12\x35\n\x0eveto_threshold\x18\x06 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\rvetoThreshold\x12I\n\x19min_initial_deposit_ratio\x18\x07 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x16minInitialDepositRatio\x12\x42\n\x15proposal_cancel_ratio\x18\x08 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x13proposalCancelRatio\x12J\n\x14proposal_cancel_dest\x18\t \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x12proposalCancelDest\x12W\n\x17\x65xpedited_voting_period\x18\n \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x15\x65xpeditedVotingPeriod\x12?\n\x13\x65xpedited_threshold\x18\x0b \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x12\x65xpeditedThreshold\x12X\n\x15\x65xpedited_min_deposit\x18\x0c \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x13\x65xpeditedMinDeposit\x12(\n\x10\x62urn_vote_quorum\x18\r \x01(\x08R\x0e\x62urnVoteQuorum\x12\x41\n\x1d\x62urn_proposal_deposit_prevote\x18\x0e \x01(\x08R\x1a\x62urnProposalDepositPrevote\x12$\n\x0e\x62urn_vote_veto\x18\x0f \x01(\x08R\x0c\x62urnVoteVeto\x12:\n\x11min_deposit_ratio\x18\x10 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x0fminDepositRatio\x12K\n\x1aproposal_cancel_max_period\x18\x11 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x17proposalCancelMaxPeriod\x12`\n\x1foptimistic_authorized_addresses\x18\x12 \x03(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x1doptimisticAuthorizedAddresses\x12R\n\x1doptimistic_rejected_threshold\x18\x13 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x1boptimisticRejectedThreshold\x12-\n\nyes_quorum\x18\x14 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\tyesQuorum\"\x96\x02\n\x12MessageBasedParams\x12\x44\n\rvoting_period\x18\x01 \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x0cvotingPeriod\x12&\n\x06quorum\x18\x02 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\x06quorum\x12-\n\nyes_quorum\x18\x14 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\tyesQuorum\x12,\n\tthreshold\x18\x03 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\tthreshold\x12\x35\n\x0eveto_threshold\x18\x04 \x01(\tB\x0e\xd2\xb4-\ncosmos.DecR\rvetoThreshold*\xa7\x01\n\x0cProposalType\x12\x1d\n\x19PROPOSAL_TYPE_UNSPECIFIED\x10\x00\x12\x1a\n\x16PROPOSAL_TYPE_STANDARD\x10\x01\x12!\n\x1dPROPOSAL_TYPE_MULTIPLE_CHOICE\x10\x02\x12\x1c\n\x18PROPOSAL_TYPE_OPTIMISTIC\x10\x03\x12\x1b\n\x17PROPOSAL_TYPE_EXPEDITED\x10\x04*\xfa\x01\n\nVoteOption\x12\x1b\n\x17VOTE_OPTION_UNSPECIFIED\x10\x00\x12\x13\n\x0fVOTE_OPTION_ONE\x10\x01\x12\x13\n\x0fVOTE_OPTION_YES\x10\x01\x12\x13\n\x0fVOTE_OPTION_TWO\x10\x02\x12\x17\n\x13VOTE_OPTION_ABSTAIN\x10\x02\x12\x15\n\x11VOTE_OPTION_THREE\x10\x03\x12\x12\n\x0eVOTE_OPTION_NO\x10\x03\x12\x14\n\x10VOTE_OPTION_FOUR\x10\x04\x12\x1c\n\x18VOTE_OPTION_NO_WITH_VETO\x10\x04\x12\x14\n\x10VOTE_OPTION_SPAM\x10\x05\x1a\x02\x10\x01*\xce\x01\n\x0eProposalStatus\x12\x1f\n\x1bPROPOSAL_STATUS_UNSPECIFIED\x10\x00\x12\"\n\x1ePROPOSAL_STATUS_DEPOSIT_PERIOD\x10\x01\x12!\n\x1dPROPOSAL_STATUS_VOTING_PERIOD\x10\x02\x12\x1a\n\x16PROPOSAL_STATUS_PASSED\x10\x03\x12\x1c\n\x18PROPOSAL_STATUS_REJECTED\x10\x04\x12\x1a\n\x16PROPOSAL_STATUS_FAILED\x10\x05\x42\x1dZ\x1b\x63osmossdk.io/x/gov/types/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.gov.v1.gov_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\033cosmossdk.io/x/gov/types/v1'
  _globals['_VOTEOPTION']._loaded_options = None
  _globals['_VOTEOPTION']._serialized_options = b'\020\001'
  _globals['_WEIGHTEDVOTEOPTION'].fields_by_name['weight']._loaded_options = None
  _globals['_WEIGHTEDVOTEOPTION'].fields_by_name['weight']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_DEPOSIT'].fields_by_name['depositor']._loaded_options = None
  _globals['_DEPOSIT'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_DEPOSIT'].fields_by_name['amount']._loaded_options = None
  _globals['_DEPOSIT'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_PROPOSAL'].fields_by_name['submit_time']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['submit_time']._serialized_options = b'\220\337\037\001'
  _globals['_PROPOSAL'].fields_by_name['deposit_end_time']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['deposit_end_time']._serialized_options = b'\220\337\037\001'
  _globals['_PROPOSAL'].fields_by_name['total_deposit']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['total_deposit']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_PROPOSAL'].fields_by_name['voting_start_time']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['voting_start_time']._serialized_options = b'\220\337\037\001'
  _globals['_PROPOSAL'].fields_by_name['voting_end_time']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['voting_end_time']._serialized_options = b'\220\337\037\001'
  _globals['_PROPOSAL'].fields_by_name['proposer']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['proposer']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_PROPOSAL'].fields_by_name['expedited']._loaded_options = None
  _globals['_PROPOSAL'].fields_by_name['expedited']._serialized_options = b'\030\001'
  _globals['_TALLYRESULT'].fields_by_name['yes_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['yes_count']._serialized_options = b'\030\001\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['abstain_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['abstain_count']._serialized_options = b'\030\001\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['no_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['no_count']._serialized_options = b'\030\001\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['no_with_veto_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['no_with_veto_count']._serialized_options = b'\030\001\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['option_one_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['option_one_count']._serialized_options = b'\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['option_two_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['option_two_count']._serialized_options = b'\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['option_three_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['option_three_count']._serialized_options = b'\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['option_four_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['option_four_count']._serialized_options = b'\322\264-\ncosmos.Int'
  _globals['_TALLYRESULT'].fields_by_name['spam_count']._loaded_options = None
  _globals['_TALLYRESULT'].fields_by_name['spam_count']._serialized_options = b'\322\264-\ncosmos.Int'
  _globals['_VOTE'].fields_by_name['voter']._loaded_options = None
  _globals['_VOTE'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_DEPOSITPARAMS'].fields_by_name['min_deposit']._loaded_options = None
  _globals['_DEPOSITPARAMS'].fields_by_name['min_deposit']._serialized_options = b'\310\336\037\000\352\336\037\025min_deposit,omitempty'
  _globals['_DEPOSITPARAMS'].fields_by_name['max_deposit_period']._loaded_options = None
  _globals['_DEPOSITPARAMS'].fields_by_name['max_deposit_period']._serialized_options = b'\352\336\037\034max_deposit_period,omitempty\230\337\037\001'
  _globals['_DEPOSITPARAMS']._loaded_options = None
  _globals['_DEPOSITPARAMS']._serialized_options = b'\030\001'
  _globals['_VOTINGPARAMS'].fields_by_name['voting_period']._loaded_options = None
  _globals['_VOTINGPARAMS'].fields_by_name['voting_period']._serialized_options = b'\230\337\037\001'
  _globals['_VOTINGPARAMS']._loaded_options = None
  _globals['_VOTINGPARAMS']._serialized_options = b'\030\001'
  _globals['_TALLYPARAMS'].fields_by_name['quorum']._loaded_options = None
  _globals['_TALLYPARAMS'].fields_by_name['quorum']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_TALLYPARAMS'].fields_by_name['threshold']._loaded_options = None
  _globals['_TALLYPARAMS'].fields_by_name['threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_TALLYPARAMS'].fields_by_name['veto_threshold']._loaded_options = None
  _globals['_TALLYPARAMS'].fields_by_name['veto_threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_TALLYPARAMS']._loaded_options = None
  _globals['_TALLYPARAMS']._serialized_options = b'\030\001'
  _globals['_PARAMS'].fields_by_name['min_deposit']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['min_deposit']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['max_deposit_period']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['max_deposit_period']._serialized_options = b'\230\337\037\001'
  _globals['_PARAMS'].fields_by_name['voting_period']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['voting_period']._serialized_options = b'\230\337\037\001'
  _globals['_PARAMS'].fields_by_name['quorum']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['quorum']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['threshold']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['veto_threshold']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['veto_threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['min_initial_deposit_ratio']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['min_initial_deposit_ratio']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['proposal_cancel_ratio']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['proposal_cancel_ratio']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['proposal_cancel_dest']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['proposal_cancel_dest']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_PARAMS'].fields_by_name['expedited_voting_period']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['expedited_voting_period']._serialized_options = b'\230\337\037\001'
  _globals['_PARAMS'].fields_by_name['expedited_threshold']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['expedited_threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['expedited_min_deposit']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['expedited_min_deposit']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['min_deposit_ratio']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['min_deposit_ratio']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['proposal_cancel_max_period']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['proposal_cancel_max_period']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['optimistic_authorized_addresses']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['optimistic_authorized_addresses']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_PARAMS'].fields_by_name['optimistic_rejected_threshold']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['optimistic_rejected_threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PARAMS'].fields_by_name['yes_quorum']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['yes_quorum']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['voting_period']._loaded_options = None
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['voting_period']._serialized_options = b'\230\337\037\001'
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['quorum']._loaded_options = None
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['quorum']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['yes_quorum']._loaded_options = None
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['yes_quorum']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['threshold']._loaded_options = None
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['veto_threshold']._loaded_options = None
  _globals['_MESSAGEBASEDPARAMS'].fields_by_name['veto_threshold']._serialized_options = b'\322\264-\ncosmos.Dec'
  _globals['_PROPOSALTYPE']._serialized_start=4343
  _globals['_PROPOSALTYPE']._serialized_end=4510
  _globals['_VOTEOPTION']._serialized_start=4513
  _globals['_VOTEOPTION']._serialized_end=4763
  _globals['_PROPOSALSTATUS']._serialized_start=4766
  _globals['_PROPOSALSTATUS']._serialized_end=4972
  _globals['_WEIGHTEDVOTEOPTION']._serialized_start=234
  _globals['_WEIGHTEDVOTEOPTION']._serialized_end=345
  _globals['_DEPOSIT']._serialized_start=348
  _globals['_DEPOSIT']._serialized_end=508
  _globals['_PROPOSAL']._serialized_start=511
  _globals['_PROPOSAL']._serialized_end=1353
  _globals['_PROPOSALVOTEOPTIONS']._serialized_start=1356
  _globals['_PROPOSALVOTEOPTIONS']._serialized_end=1540
  _globals['_TALLYRESULT']._serialized_start=1543
  _globals['_TALLYRESULT']._serialized_end=2051
  _globals['_VOTE']._serialized_start=2054
  _globals['_VOTE']._serialized_end=2236
  _globals['_DEPOSITPARAMS']._serialized_start=2239
  _globals['_DEPOSITPARAMS']._serialized_end=2460
  _globals['_VOTINGPARAMS']._serialized_start=2462
  _globals['_VOTINGPARAMS']._serialized_end=2550
  _globals['_TALLYPARAMS']._serialized_start=2553
  _globals['_TALLYPARAMS']._serialized_end=2711
  _globals['_PARAMS']._serialized_start=2714
  _globals['_PARAMS']._serialized_end=4059
  _globals['_MESSAGEBASEDPARAMS']._serialized_start=4062
  _globals['_MESSAGEBASEDPARAMS']._serialized_end=4340
# @@protoc_insertion_point(module_scope)
