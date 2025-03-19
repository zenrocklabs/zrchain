# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/protocolpool/v1/tx.proto
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
    'cosmos/protocolpool/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/protocolpool/v1/tx.proto\x12\x16\x63osmos.protocolpool.v1\x1a\x14gogoproto/gogo.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1egoogle/protobuf/duration.proto\"\xcb\x01\n\x14MsgFundCommunityPool\x12\x63\n\x06\x61mount\x18\x01 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x06\x61mount\x12\x36\n\tdepositor\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor:\x16\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\x82\xe7\xb0*\tdepositor\"\x1e\n\x1cMsgFundCommunityPoolResponse\"\xe2\x01\n\x15MsgCommunityPoolSpend\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x1c\n\trecipient\x18\x02 \x01(\tR\trecipient\x12\x63\n\x06\x61mount\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x06\x61mount:\x0e\x82\xe7\xb0*\tauthority\"\x1f\n\x1dMsgCommunityPoolSpendResponse\"\xfc\x02\n\x17MsgSubmitBudgetProposal\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x45\n\x11recipient_address\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10recipientAddress\x12<\n\x0ctotal_budget\x18\x03 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinR\x0btotalBudget\x12?\n\nstart_time\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\tstartTime\x12\x1a\n\x08tranches\x18\x05 \x01(\x04R\x08tranches\x12\x37\n\x06period\x18\x06 \x01(\x0b\x32\x19.google.protobuf.DurationB\x04\x98\xdf\x1f\x01R\x06period:\x0e\x82\xe7\xb0*\tauthority\"!\n\x1fMsgSubmitBudgetProposalResponse\"o\n\x0eMsgClaimBudget\x12\x45\n\x11recipient_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10recipientAddress:\x16\x82\xe7\xb0*\x11recipient_address\"}\n\x16MsgClaimBudgetResponse\x12\x63\n\x06\x61mount\x18\x01 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x06\x61mount\"\xa6\x02\n\x17MsgCreateContinuousFund\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x36\n\trecipient\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\trecipient\x12Q\n\npercentage\x18\x03 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\npercentage\x12\x38\n\x06\x65xpiry\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x04\x90\xdf\x1f\x01R\x06\x65xpiry:\x0e\x82\xe7\xb0*\tauthority\"!\n\x1fMsgCreateContinuousFundResponse\"\xa8\x01\n\x17MsgCancelContinuousFund\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x45\n\x11recipient_address\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10recipientAddress:\x0e\x82\xe7\xb0*\tauthority\"\xe4\x02\n\x1fMsgCancelContinuousFundResponse\x12I\n\rcanceled_time\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\x0c\x63\x61nceledTime\x12\'\n\x0f\x63\x61nceled_height\x18\x02 \x01(\x04R\x0e\x63\x61nceledHeight\x12\x45\n\x11recipient_address\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10recipientAddress\x12\x85\x01\n\x18withdrawn_allocated_fund\x18\x04 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x16withdrawnAllocatedFund\"z\n\x19MsgWithdrawContinuousFund\x12\x45\n\x11recipient_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10recipientAddress:\x16\x82\xe7\xb0*\x11recipient_address\"\x88\x01\n!MsgWithdrawContinuousFundResponse\x12\x63\n\x06\x61mount\x18\x01 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x06\x61mount2\xfa\x06\n\x03Msg\x12w\n\x11\x46undCommunityPool\x12,.cosmos.protocolpool.v1.MsgFundCommunityPool\x1a\x34.cosmos.protocolpool.v1.MsgFundCommunityPoolResponse\x12z\n\x12\x43ommunityPoolSpend\x12-.cosmos.protocolpool.v1.MsgCommunityPoolSpend\x1a\x35.cosmos.protocolpool.v1.MsgCommunityPoolSpendResponse\x12\x80\x01\n\x14SubmitBudgetProposal\x12/.cosmos.protocolpool.v1.MsgSubmitBudgetProposal\x1a\x37.cosmos.protocolpool.v1.MsgSubmitBudgetProposalResponse\x12\x65\n\x0b\x43laimBudget\x12&.cosmos.protocolpool.v1.MsgClaimBudget\x1a..cosmos.protocolpool.v1.MsgClaimBudgetResponse\x12\x80\x01\n\x14\x43reateContinuousFund\x12/.cosmos.protocolpool.v1.MsgCreateContinuousFund\x1a\x37.cosmos.protocolpool.v1.MsgCreateContinuousFundResponse\x12\x86\x01\n\x16WithdrawContinuousFund\x12\x31.cosmos.protocolpool.v1.MsgWithdrawContinuousFund\x1a\x39.cosmos.protocolpool.v1.MsgWithdrawContinuousFundResponse\x12\x80\x01\n\x14\x43\x61ncelContinuousFund\x12/.cosmos.protocolpool.v1.MsgCancelContinuousFund\x1a\x37.cosmos.protocolpool.v1.MsgCancelContinuousFundResponse\x1a\x05\x80\xe7\xb0*\x01\x42#Z!cosmossdk.io/x/protocolpool/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.protocolpool.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z!cosmossdk.io/x/protocolpool/types'
  _globals['_MSGFUNDCOMMUNITYPOOL'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGFUNDCOMMUNITYPOOL'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGFUNDCOMMUNITYPOOL'].fields_by_name['depositor']._loaded_options = None
  _globals['_MSGFUNDCOMMUNITYPOOL'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGFUNDCOMMUNITYPOOL']._loaded_options = None
  _globals['_MSGFUNDCOMMUNITYPOOL']._serialized_options = b'\210\240\037\000\350\240\037\000\202\347\260*\tdepositor'
  _globals['_MSGCOMMUNITYPOOLSPEND'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGCOMMUNITYPOOLSPEND'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCOMMUNITYPOOLSPEND'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGCOMMUNITYPOOLSPEND'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGCOMMUNITYPOOLSPEND']._loaded_options = None
  _globals['_MSGCOMMUNITYPOOLSPEND']._serialized_options = b'\202\347\260*\tauthority'
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['recipient_address']._loaded_options = None
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['recipient_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['start_time']._loaded_options = None
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['start_time']._serialized_options = b'\220\337\037\001'
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['period']._loaded_options = None
  _globals['_MSGSUBMITBUDGETPROPOSAL'].fields_by_name['period']._serialized_options = b'\230\337\037\001'
  _globals['_MSGSUBMITBUDGETPROPOSAL']._loaded_options = None
  _globals['_MSGSUBMITBUDGETPROPOSAL']._serialized_options = b'\202\347\260*\tauthority'
  _globals['_MSGCLAIMBUDGET'].fields_by_name['recipient_address']._loaded_options = None
  _globals['_MSGCLAIMBUDGET'].fields_by_name['recipient_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCLAIMBUDGET']._loaded_options = None
  _globals['_MSGCLAIMBUDGET']._serialized_options = b'\202\347\260*\021recipient_address'
  _globals['_MSGCLAIMBUDGETRESPONSE'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGCLAIMBUDGETRESPONSE'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['recipient']._loaded_options = None
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['recipient']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['percentage']._loaded_options = None
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['percentage']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['expiry']._loaded_options = None
  _globals['_MSGCREATECONTINUOUSFUND'].fields_by_name['expiry']._serialized_options = b'\220\337\037\001'
  _globals['_MSGCREATECONTINUOUSFUND']._loaded_options = None
  _globals['_MSGCREATECONTINUOUSFUND']._serialized_options = b'\202\347\260*\tauthority'
  _globals['_MSGCANCELCONTINUOUSFUND'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGCANCELCONTINUOUSFUND'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCANCELCONTINUOUSFUND'].fields_by_name['recipient_address']._loaded_options = None
  _globals['_MSGCANCELCONTINUOUSFUND'].fields_by_name['recipient_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCANCELCONTINUOUSFUND']._loaded_options = None
  _globals['_MSGCANCELCONTINUOUSFUND']._serialized_options = b'\202\347\260*\tauthority'
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE'].fields_by_name['canceled_time']._loaded_options = None
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE'].fields_by_name['canceled_time']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE'].fields_by_name['recipient_address']._loaded_options = None
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE'].fields_by_name['recipient_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE'].fields_by_name['withdrawn_allocated_fund']._loaded_options = None
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE'].fields_by_name['withdrawn_allocated_fund']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGWITHDRAWCONTINUOUSFUND'].fields_by_name['recipient_address']._loaded_options = None
  _globals['_MSGWITHDRAWCONTINUOUSFUND'].fields_by_name['recipient_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGWITHDRAWCONTINUOUSFUND']._loaded_options = None
  _globals['_MSGWITHDRAWCONTINUOUSFUND']._serialized_options = b'\202\347\260*\021recipient_address'
  _globals['_MSGWITHDRAWCONTINUOUSFUNDRESPONSE'].fields_by_name['amount']._loaded_options = None
  _globals['_MSGWITHDRAWCONTINUOUSFUNDRESPONSE'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGFUNDCOMMUNITYPOOL']._serialized_start=231
  _globals['_MSGFUNDCOMMUNITYPOOL']._serialized_end=434
  _globals['_MSGFUNDCOMMUNITYPOOLRESPONSE']._serialized_start=436
  _globals['_MSGFUNDCOMMUNITYPOOLRESPONSE']._serialized_end=466
  _globals['_MSGCOMMUNITYPOOLSPEND']._serialized_start=469
  _globals['_MSGCOMMUNITYPOOLSPEND']._serialized_end=695
  _globals['_MSGCOMMUNITYPOOLSPENDRESPONSE']._serialized_start=697
  _globals['_MSGCOMMUNITYPOOLSPENDRESPONSE']._serialized_end=728
  _globals['_MSGSUBMITBUDGETPROPOSAL']._serialized_start=731
  _globals['_MSGSUBMITBUDGETPROPOSAL']._serialized_end=1111
  _globals['_MSGSUBMITBUDGETPROPOSALRESPONSE']._serialized_start=1113
  _globals['_MSGSUBMITBUDGETPROPOSALRESPONSE']._serialized_end=1146
  _globals['_MSGCLAIMBUDGET']._serialized_start=1148
  _globals['_MSGCLAIMBUDGET']._serialized_end=1259
  _globals['_MSGCLAIMBUDGETRESPONSE']._serialized_start=1261
  _globals['_MSGCLAIMBUDGETRESPONSE']._serialized_end=1386
  _globals['_MSGCREATECONTINUOUSFUND']._serialized_start=1389
  _globals['_MSGCREATECONTINUOUSFUND']._serialized_end=1683
  _globals['_MSGCREATECONTINUOUSFUNDRESPONSE']._serialized_start=1685
  _globals['_MSGCREATECONTINUOUSFUNDRESPONSE']._serialized_end=1718
  _globals['_MSGCANCELCONTINUOUSFUND']._serialized_start=1721
  _globals['_MSGCANCELCONTINUOUSFUND']._serialized_end=1889
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE']._serialized_start=1892
  _globals['_MSGCANCELCONTINUOUSFUNDRESPONSE']._serialized_end=2248
  _globals['_MSGWITHDRAWCONTINUOUSFUND']._serialized_start=2250
  _globals['_MSGWITHDRAWCONTINUOUSFUND']._serialized_end=2372
  _globals['_MSGWITHDRAWCONTINUOUSFUNDRESPONSE']._serialized_start=2375
  _globals['_MSGWITHDRAWCONTINUOUSFUNDRESPONSE']._serialized_end=2511
  _globals['_MSG']._serialized_start=2514
  _globals['_MSG']._serialized_end=3404
# @@protoc_insertion_point(module_scope)
