# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/gov/v1/tx.proto
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
    'cosmos/gov/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos.gov.v1 import gov_pb2 as cosmos_dot_gov_dot_v1_dot_gov__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from amino import amino_pb2 as amino_dot_amino__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x16\x63osmos/gov/v1/tx.proto\x12\rcosmos.gov.v1\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x17\x63osmos/gov/v1/gov.proto\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x19google/protobuf/any.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x11\x61mino/amino.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xe4\x03\n\x11MsgSubmitProposal\x12\x30\n\x08messages\x18\x01 \x03(\x0b\x32\x14.google.protobuf.AnyR\x08messages\x12\x8a\x01\n\x0finitial_deposit\x18\x02 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x0einitialDeposit\x12\x34\n\x08proposer\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08proposer\x12\x1a\n\x08metadata\x18\x04 \x01(\tR\x08metadata\x12)\n\x05title\x18\x05 \x01(\tB\x13\xda\xb4-\x0f\x63osmos-sdk 0.47R\x05title\x12-\n\x07summary\x18\x06 \x01(\tB\x13\xda\xb4-\x0f\x63osmos-sdk 0.47R\x07summary\x12\x31\n\texpedited\x18\x07 \x01(\x08\x42\x13\xda\xb4-\x0f\x63osmos-sdk 0.50R\texpedited:1\x82\xe7\xb0*\x08proposer\x8a\xe7\xb0*\x1f\x63osmos-sdk/v1/MsgSubmitProposal\"<\n\x19MsgSubmitProposalResponse\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"\xbb\x01\n\x14MsgExecLegacyContent\x12N\n\x07\x63ontent\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyB\x1e\xca\xb4-\x1a\x63osmos.gov.v1beta1.ContentR\x07\x63ontent\x12\x1c\n\tauthority\x18\x02 \x01(\tR\tauthority:5\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*\"cosmos-sdk/v1/MsgExecLegacyContent\"\x1e\n\x1cMsgExecLegacyContentResponse\"\xe5\x01\n\x07MsgVote\x12\x35\n\x0bproposal_id\x18\x01 \x01(\x04\x42\x14\xea\xde\x1f\x0bproposal_id\xa8\xe7\xb0*\x01R\nproposalId\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter\x12\x31\n\x06option\x18\x03 \x01(\x0e\x32\x19.cosmos.gov.v1.VoteOptionR\x06option\x12\x1a\n\x08metadata\x18\x04 \x01(\tR\x08metadata:$\x82\xe7\xb0*\x05voter\x8a\xe7\xb0*\x15\x63osmos-sdk/v1/MsgVote\"\x11\n\x0fMsgVoteResponse\"\xff\x01\n\x0fMsgVoteWeighted\x12\x35\n\x0bproposal_id\x18\x01 \x01(\x04\x42\x14\xea\xde\x1f\x0bproposal_id\xa8\xe7\xb0*\x01R\nproposalId\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter\x12;\n\x07options\x18\x03 \x03(\x0b\x32!.cosmos.gov.v1.WeightedVoteOptionR\x07options\x12\x1a\n\x08metadata\x18\x04 \x01(\tR\x08metadata:,\x82\xe7\xb0*\x05voter\x8a\xe7\xb0*\x1d\x63osmos-sdk/v1/MsgVoteWeighted\"\x19\n\x17MsgVoteWeightedResponse\"\xe6\x01\n\nMsgDeposit\x12\x35\n\x0bproposal_id\x18\x01 \x01(\x04\x42\x14\xea\xde\x1f\x0bproposal_id\xa8\xe7\xb0*\x01R\nproposalId\x12\x36\n\tdepositor\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor\x12<\n\x06\x61mount\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06\x61mount:+\x82\xe7\xb0*\tdepositor\x8a\xe7\xb0*\x18\x63osmos-sdk/v1/MsgDeposit\"\x14\n\x12MsgDepositResponse\"\xce\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x38\n\x06params\x18\x02 \x01(\x0b\x32\x15.cosmos.gov.v1.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:I\xd2\xb4-\x0f\x63osmos-sdk 0.47\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*#cosmos-sdk/x/gov/v1/MsgUpdateParams\".\n\x17MsgUpdateParamsResponse:\x13\xd2\xb4-\x0f\x63osmos-sdk 0.47\"\x9d\x01\n\x11MsgCancelProposal\x12\x30\n\x0bproposal_id\x18\x01 \x01(\x04\x42\x0f\xea\xde\x1f\x0bproposal_idR\nproposalId\x12\x34\n\x08proposer\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08proposer: \xd2\xb4-\x0f\x63osmos-sdk 0.50\x82\xe7\xb0*\x08proposer\"\xd6\x01\n\x19MsgCancelProposalResponse\x12\x30\n\x0bproposal_id\x18\x01 \x01(\x04\x42\x0f\xea\xde\x1f\x0bproposal_idR\nproposalId\x12I\n\rcanceled_time\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x0c\x63\x61nceledTime\x12\'\n\x0f\x63\x61nceled_height\x18\x03 \x01(\x04R\x0e\x63\x61nceledHeight:\x13\xd2\xb4-\x0f\x63osmos-sdk 0.502\xfb\x05\n\x03Msg\x12q\n\x0eSubmitProposal\x12 .cosmos.gov.v1.MsgSubmitProposal\x1a(.cosmos.gov.v1.MsgSubmitProposalResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.46\x12z\n\x11\x45xecLegacyContent\x12#.cosmos.gov.v1.MsgExecLegacyContent\x1a+.cosmos.gov.v1.MsgExecLegacyContentResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.46\x12S\n\x04Vote\x12\x16.cosmos.gov.v1.MsgVote\x1a\x1e.cosmos.gov.v1.MsgVoteResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.46\x12k\n\x0cVoteWeighted\x12\x1e.cosmos.gov.v1.MsgVoteWeighted\x1a&.cosmos.gov.v1.MsgVoteWeightedResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.46\x12\\\n\x07\x44\x65posit\x12\x19.cosmos.gov.v1.MsgDeposit\x1a!.cosmos.gov.v1.MsgDepositResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.46\x12k\n\x0cUpdateParams\x12\x1e.cosmos.gov.v1.MsgUpdateParams\x1a&.cosmos.gov.v1.MsgUpdateParamsResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.47\x12q\n\x0e\x43\x61ncelProposal\x12 .cosmos.gov.v1.MsgCancelProposal\x1a(.cosmos.gov.v1.MsgCancelProposalResponse\"\x13\xca\xb4-\x0f\x63osmos-sdk 0.50\x1a\x05\x80\xe7\xb0*\x01\x42-Z+github.com/cosmos/cosmos-sdk/x/gov/types/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.gov.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+github.com/cosmos/cosmos-sdk/x/gov/types/v1'
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['initial_deposit']._loaded_options = None
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['initial_deposit']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['proposer']._loaded_options = None
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['proposer']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['title']._loaded_options = None
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['title']._serialized_options = b'\332\264-\017cosmos-sdk 0.47'
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['summary']._loaded_options = None
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['summary']._serialized_options = b'\332\264-\017cosmos-sdk 0.47'
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['expedited']._loaded_options = None
  _globals['_MSGSUBMITPROPOSAL'].fields_by_name['expedited']._serialized_options = b'\332\264-\017cosmos-sdk 0.50'
  _globals['_MSGSUBMITPROPOSAL']._loaded_options = None
  _globals['_MSGSUBMITPROPOSAL']._serialized_options = b'\202\347\260*\010proposer\212\347\260*\037cosmos-sdk/v1/MsgSubmitProposal'
  _globals['_MSGEXECLEGACYCONTENT'].fields_by_name['content']._loaded_options = None
  _globals['_MSGEXECLEGACYCONTENT'].fields_by_name['content']._serialized_options = b'\312\264-\032cosmos.gov.v1beta1.Content'
  _globals['_MSGEXECLEGACYCONTENT']._loaded_options = None
  _globals['_MSGEXECLEGACYCONTENT']._serialized_options = b'\202\347\260*\tauthority\212\347\260*\"cosmos-sdk/v1/MsgExecLegacyContent'
  _globals['_MSGVOTE'].fields_by_name['proposal_id']._loaded_options = None
  _globals['_MSGVOTE'].fields_by_name['proposal_id']._serialized_options = b'\352\336\037\013proposal_id\250\347\260*\001'
  _globals['_MSGVOTE'].fields_by_name['voter']._loaded_options = None
  _globals['_MSGVOTE'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGVOTE']._loaded_options = None
  _globals['_MSGVOTE']._serialized_options = b'\202\347\260*\005voter\212\347\260*\025cosmos-sdk/v1/MsgVote'
  _globals['_MSGVOTEWEIGHTED'].fields_by_name['proposal_id']._loaded_options = None
  _globals['_MSGVOTEWEIGHTED'].fields_by_name['proposal_id']._serialized_options = b'\352\336\037\013proposal_id\250\347\260*\001'
  _globals['_MSGVOTEWEIGHTED'].fields_by_name['voter']._loaded_options = None
  _globals['_MSGVOTEWEIGHTED'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGVOTEWEIGHTED']._loaded_options = None
  _globals['_MSGVOTEWEIGHTED']._serialized_options = b'\202\347\260*\005voter\212\347\260*\035cosmos-sdk/v1/MsgVoteWeighted'
  _globals['_MSGDEPOSIT'].fields_by_name['proposal_id']._loaded_options = None
  _globals['_MSGDEPOSIT'].fields_by_name['proposal_id']._serialized_options = b'\352\336\037\013proposal_id\250\347\260*\001'
  _globals['_MSGDEPOSIT'].fields_by_name['depositor']._loaded_options = None
  _globals['_MSGDEPOSIT'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGDEPOSIT'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGDEPOSIT'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGDEPOSIT']._loaded_options = None
  _globals['_MSGDEPOSIT']._serialized_options = b'\202\347\260*\tdepositor\212\347\260*\030cosmos-sdk/v1/MsgDeposit'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\322\264-\017cosmos-sdk 0.47\202\347\260*\tauthority\212\347\260*#cosmos-sdk/x/gov/v1/MsgUpdateParams'
  _globals['_MSGUPDATEPARAMSRESPONSE']._loaded_options = None
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_options = b'\322\264-\017cosmos-sdk 0.47'
  _globals['_MSGCANCELPROPOSAL'].fields_by_name['proposal_id']._loaded_options = None
  _globals['_MSGCANCELPROPOSAL'].fields_by_name['proposal_id']._serialized_options = b'\352\336\037\013proposal_id'
  _globals['_MSGCANCELPROPOSAL'].fields_by_name['proposer']._loaded_options = None
  _globals['_MSGCANCELPROPOSAL'].fields_by_name['proposer']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCANCELPROPOSAL']._loaded_options = None
  _globals['_MSGCANCELPROPOSAL']._serialized_options = b'\322\264-\017cosmos-sdk 0.50\202\347\260*\010proposer'
  _globals['_MSGCANCELPROPOSALRESPONSE'].fields_by_name['proposal_id']._loaded_options = None
  _globals['_MSGCANCELPROPOSALRESPONSE'].fields_by_name['proposal_id']._serialized_options = b'\352\336\037\013proposal_id'
  _globals['_MSGCANCELPROPOSALRESPONSE'].fields_by_name['canceled_time']._loaded_options = None
  _globals['_MSGCANCELPROPOSALRESPONSE'].fields_by_name['canceled_time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_MSGCANCELPROPOSALRESPONSE']._loaded_options = None
  _globals['_MSGCANCELPROPOSALRESPONSE']._serialized_options = b'\322\264-\017cosmos-sdk 0.50'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSG'].methods_by_name['SubmitProposal']._loaded_options = None
  _globals['_MSG'].methods_by_name['SubmitProposal']._serialized_options = b'\312\264-\017cosmos-sdk 0.46'
  _globals['_MSG'].methods_by_name['ExecLegacyContent']._loaded_options = None
  _globals['_MSG'].methods_by_name['ExecLegacyContent']._serialized_options = b'\312\264-\017cosmos-sdk 0.46'
  _globals['_MSG'].methods_by_name['Vote']._loaded_options = None
  _globals['_MSG'].methods_by_name['Vote']._serialized_options = b'\312\264-\017cosmos-sdk 0.46'
  _globals['_MSG'].methods_by_name['VoteWeighted']._loaded_options = None
  _globals['_MSG'].methods_by_name['VoteWeighted']._serialized_options = b'\312\264-\017cosmos-sdk 0.46'
  _globals['_MSG'].methods_by_name['Deposit']._loaded_options = None
  _globals['_MSG'].methods_by_name['Deposit']._serialized_options = b'\312\264-\017cosmos-sdk 0.46'
  _globals['_MSG'].methods_by_name['UpdateParams']._loaded_options = None
  _globals['_MSG'].methods_by_name['UpdateParams']._serialized_options = b'\312\264-\017cosmos-sdk 0.47'
  _globals['_MSG'].methods_by_name['CancelProposal']._loaded_options = None
  _globals['_MSG'].methods_by_name['CancelProposal']._serialized_options = b'\312\264-\017cosmos-sdk 0.50'
  _globals['_MSGSUBMITPROPOSAL']._serialized_start=252
  _globals['_MSGSUBMITPROPOSAL']._serialized_end=736
  _globals['_MSGSUBMITPROPOSALRESPONSE']._serialized_start=738
  _globals['_MSGSUBMITPROPOSALRESPONSE']._serialized_end=798
  _globals['_MSGEXECLEGACYCONTENT']._serialized_start=801
  _globals['_MSGEXECLEGACYCONTENT']._serialized_end=988
  _globals['_MSGEXECLEGACYCONTENTRESPONSE']._serialized_start=990
  _globals['_MSGEXECLEGACYCONTENTRESPONSE']._serialized_end=1020
  _globals['_MSGVOTE']._serialized_start=1023
  _globals['_MSGVOTE']._serialized_end=1252
  _globals['_MSGVOTERESPONSE']._serialized_start=1254
  _globals['_MSGVOTERESPONSE']._serialized_end=1271
  _globals['_MSGVOTEWEIGHTED']._serialized_start=1274
  _globals['_MSGVOTEWEIGHTED']._serialized_end=1529
  _globals['_MSGVOTEWEIGHTEDRESPONSE']._serialized_start=1531
  _globals['_MSGVOTEWEIGHTEDRESPONSE']._serialized_end=1556
  _globals['_MSGDEPOSIT']._serialized_start=1559
  _globals['_MSGDEPOSIT']._serialized_end=1789
  _globals['_MSGDEPOSITRESPONSE']._serialized_start=1791
  _globals['_MSGDEPOSITRESPONSE']._serialized_end=1811
  _globals['_MSGUPDATEPARAMS']._serialized_start=1814
  _globals['_MSGUPDATEPARAMS']._serialized_end=2020
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=2022
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=2068
  _globals['_MSGCANCELPROPOSAL']._serialized_start=2071
  _globals['_MSGCANCELPROPOSAL']._serialized_end=2228
  _globals['_MSGCANCELPROPOSALRESPONSE']._serialized_start=2231
  _globals['_MSGCANCELPROPOSALRESPONSE']._serialized_end=2445
  _globals['_MSG']._serialized_start=2448
  _globals['_MSG']._serialized_end=3211
# @@protoc_insertion_point(module_scope)
