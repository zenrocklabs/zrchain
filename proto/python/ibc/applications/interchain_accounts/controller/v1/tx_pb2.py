# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/interchain_accounts/controller/v1/tx.proto
# Protobuf Python Version: 5.29.2
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
    2,
    '',
    'ibc/applications/interchain_accounts/controller/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from ibc.applications.interchain_accounts.v1 import packet_pb2 as ibc_dot_applications_dot_interchain__accounts_dot_v1_dot_packet__pb2
from ibc.applications.interchain_accounts.controller.v1 import controller_pb2 as ibc_dot_applications_dot_interchain__accounts_dot_controller_dot_v1_dot_controller__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from ibc.core.channel.v1 import channel_pb2 as ibc_dot_core_dot_channel_dot_v1_dot_channel__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n;ibc/applications/interchain_accounts/controller/v1/tx.proto\x12\x32ibc.applications.interchain_accounts.controller.v1\x1a\x14gogoproto/gogo.proto\x1a\x34ibc/applications/interchain_accounts/v1/packet.proto\x1a\x43ibc/applications/interchain_accounts/controller/v1/controller.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a!ibc/core/channel/v1/channel.proto\"\xbb\x01\n\x1cMsgRegisterInterchainAccount\x12\x14\n\x05owner\x18\x01 \x01(\tR\x05owner\x12#\n\rconnection_id\x18\x02 \x01(\tR\x0c\x63onnectionId\x12\x18\n\x07version\x18\x03 \x01(\tR\x07version\x12\x36\n\x08ordering\x18\x04 \x01(\x0e\x32\x1a.ibc.core.channel.v1.OrderR\x08ordering:\x0e\x88\xa0\x1f\x00\x82\xe7\xb0*\x05owner\"d\n$MsgRegisterInterchainAccountResponse\x12\x1d\n\nchannel_id\x18\x01 \x01(\tR\tchannelId\x12\x17\n\x07port_id\x18\x02 \x01(\tR\x06portId:\x04\x88\xa0\x1f\x00\"\xee\x01\n\tMsgSendTx\x12\x14\n\x05owner\x18\x01 \x01(\tR\x05owner\x12#\n\rconnection_id\x18\x02 \x01(\tR\x0c\x63onnectionId\x12k\n\x0bpacket_data\x18\x03 \x01(\x0b\x32\x44.ibc.applications.interchain_accounts.v1.InterchainAccountPacketDataB\x04\xc8\xde\x1f\x00R\npacketData\x12)\n\x10relative_timeout\x18\x04 \x01(\x04R\x0frelativeTimeout:\x0e\x88\xa0\x1f\x00\x82\xe7\xb0*\x05owner\"5\n\x11MsgSendTxResponse\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequence:\x04\x88\xa0\x1f\x00\"\x94\x01\n\x0fMsgUpdateParams\x12\x16\n\x06signer\x18\x01 \x01(\tR\x06signer\x12X\n\x06params\x18\x02 \x01(\x0b\x32:.ibc.applications.interchain_accounts.controller.v1.ParamsB\x04\xc8\xde\x1f\x00R\x06params:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x19\n\x17MsgUpdateParamsResponse2\x8a\x04\n\x03Msg\x12\xc7\x01\n\x19RegisterInterchainAccount\x12P.ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount\x1aX.ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccountResponse\x12\x8e\x01\n\x06SendTx\x12=.ibc.applications.interchain_accounts.controller.v1.MsgSendTx\x1a\x45.ibc.applications.interchain_accounts.controller.v1.MsgSendTxResponse\x12\xa0\x01\n\x0cUpdateParams\x12\x43.ibc.applications.interchain_accounts.controller.v1.MsgUpdateParams\x1aK.ibc.applications.interchain_accounts.controller.v1.MsgUpdateParamsResponse\x1a\x05\x80\xe7\xb0*\x01\x42RZPgithub.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/controller/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.interchain_accounts.controller.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZPgithub.com/cosmos/ibc-go/v9/modules/apps/27-interchain-accounts/controller/types'
  _globals['_MSGREGISTERINTERCHAINACCOUNT']._loaded_options = None
  _globals['_MSGREGISTERINTERCHAINACCOUNT']._serialized_options = b'\210\240\037\000\202\347\260*\005owner'
  _globals['_MSGREGISTERINTERCHAINACCOUNTRESPONSE']._loaded_options = None
  _globals['_MSGREGISTERINTERCHAINACCOUNTRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSGSENDTX'].fields_by_name['packet_data']._loaded_options = None
  _globals['_MSGSENDTX'].fields_by_name['packet_data']._serialized_options = b'\310\336\037\000'
  _globals['_MSGSENDTX']._loaded_options = None
  _globals['_MSGSENDTX']._serialized_options = b'\210\240\037\000\202\347\260*\005owner'
  _globals['_MSGSENDTXRESPONSE']._loaded_options = None
  _globals['_MSGSENDTXRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGREGISTERINTERCHAINACCOUNT']._serialized_start=321
  _globals['_MSGREGISTERINTERCHAINACCOUNT']._serialized_end=508
  _globals['_MSGREGISTERINTERCHAINACCOUNTRESPONSE']._serialized_start=510
  _globals['_MSGREGISTERINTERCHAINACCOUNTRESPONSE']._serialized_end=610
  _globals['_MSGSENDTX']._serialized_start=613
  _globals['_MSGSENDTX']._serialized_end=851
  _globals['_MSGSENDTXRESPONSE']._serialized_start=853
  _globals['_MSGSENDTXRESPONSE']._serialized_end=906
  _globals['_MSGUPDATEPARAMS']._serialized_start=909
  _globals['_MSGUPDATEPARAMS']._serialized_end=1057
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=1059
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=1084
  _globals['_MSG']._serialized_start=1087
  _globals['_MSG']._serialized_end=1609
# @@protoc_insertion_point(module_scope)
