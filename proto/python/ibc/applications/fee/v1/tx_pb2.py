# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/fee/v1/tx.proto
<<<<<<< HEAD
# Protobuf Python Version: 5.29.1
=======
# Protobuf Python Version: 5.29.0
>>>>>>> main
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
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'ibc/applications/fee/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from ibc.applications.fee.v1 import fee_pb2 as ibc_dot_applications_dot_fee_dot_v1_dot_fee__pb2
from ibc.core.channel.v1 import channel_pb2 as ibc_dot_core_dot_channel_dot_v1_dot_channel__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n ibc/applications/fee/v1/tx.proto\x12\x17ibc.applications.fee.v1\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a!ibc/applications/fee/v1/fee.proto\x1a!ibc/core/channel/v1/channel.proto\x1a\x17\x63osmos/msg/v1/msg.proto\"\xac\x01\n\x10MsgRegisterPayee\x12\x17\n\x07port_id\x18\x01 \x01(\tR\x06portId\x12\x1d\n\nchannel_id\x18\x02 \x01(\tR\tchannelId\x12\x18\n\x07relayer\x18\x03 \x01(\tR\x07relayer\x12\x14\n\x05payee\x18\x04 \x01(\tR\x05payee:0\x88\xa0\x1f\x00\x82\xe7\xb0*\x07relayer\x8a\xe7\xb0*\x1b\x63osmos-sdk/MsgRegisterPayee\"\x1a\n\x18MsgRegisterPayeeResponse\"\xdd\x01\n\x1cMsgRegisterCounterpartyPayee\x12\x17\n\x07port_id\x18\x01 \x01(\tR\x06portId\x12\x1d\n\nchannel_id\x18\x02 \x01(\tR\tchannelId\x12\x18\n\x07relayer\x18\x03 \x01(\tR\x07relayer\x12-\n\x12\x63ounterparty_payee\x18\x04 \x01(\tR\x11\x63ounterpartyPayee:<\x88\xa0\x1f\x00\x82\xe7\xb0*\x07relayer\x8a\xe7\xb0*\'cosmos-sdk/MsgRegisterCounterpartyPayee\"&\n$MsgRegisterCounterpartyPayeeResponse\"\x82\x02\n\x0fMsgPayPacketFee\x12\x39\n\x03\x66\x65\x65\x18\x01 \x01(\x0b\x32\x1c.ibc.applications.fee.v1.FeeB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x03\x66\x65\x65\x12$\n\x0esource_port_id\x18\x02 \x01(\tR\x0csourcePortId\x12*\n\x11source_channel_id\x18\x03 \x01(\tR\x0fsourceChannelId\x12\x16\n\x06signer\x18\x04 \x01(\tR\x06signer\x12\x1a\n\x08relayers\x18\x05 \x03(\tR\x08relayers:.\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\x8a\xe7\xb0*\x1a\x63osmos-sdk/MsgPayPacketFee\"\x19\n\x17MsgPayPacketFeeResponse\"\xe4\x01\n\x14MsgPayPacketFeeAsync\x12\x45\n\tpacket_id\x18\x01 \x01(\x0b\x32\x1d.ibc.core.channel.v1.PacketIdB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x08packetId\x12L\n\npacket_fee\x18\x02 \x01(\x0b\x32\".ibc.applications.fee.v1.PacketFeeB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\tpacketFee:7\x88\xa0\x1f\x00\x82\xe7\xb0*\npacket_fee\x8a\xe7\xb0*\x1f\x63osmos-sdk/MsgPayPacketFeeAsync\"\x1e\n\x1cMsgPayPacketFeeAsyncResponse2\xf6\x03\n\x03Msg\x12m\n\rRegisterPayee\x12).ibc.applications.fee.v1.MsgRegisterPayee\x1a\x31.ibc.applications.fee.v1.MsgRegisterPayeeResponse\x12\x91\x01\n\x19RegisterCounterpartyPayee\x12\x35.ibc.applications.fee.v1.MsgRegisterCounterpartyPayee\x1a=.ibc.applications.fee.v1.MsgRegisterCounterpartyPayeeResponse\x12j\n\x0cPayPacketFee\x12(.ibc.applications.fee.v1.MsgPayPacketFee\x1a\x30.ibc.applications.fee.v1.MsgPayPacketFeeResponse\x12y\n\x11PayPacketFeeAsync\x12-.ibc.applications.fee.v1.MsgPayPacketFeeAsync\x1a\x35.ibc.applications.fee.v1.MsgPayPacketFeeAsyncResponse\x1a\x05\x80\xe7\xb0*\x01\x42\x37Z5github.com/cosmos/ibc-go/v9/modules/apps/29-fee/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.fee.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/cosmos/ibc-go/v9/modules/apps/29-fee/types'
  _globals['_MSGREGISTERPAYEE']._loaded_options = None
  _globals['_MSGREGISTERPAYEE']._serialized_options = b'\210\240\037\000\202\347\260*\007relayer\212\347\260*\033cosmos-sdk/MsgRegisterPayee'
  _globals['_MSGREGISTERCOUNTERPARTYPAYEE']._loaded_options = None
  _globals['_MSGREGISTERCOUNTERPARTYPAYEE']._serialized_options = b'\210\240\037\000\202\347\260*\007relayer\212\347\260*\'cosmos-sdk/MsgRegisterCounterpartyPayee'
  _globals['_MSGPAYPACKETFEE'].fields_by_name['fee']._loaded_options = None
  _globals['_MSGPAYPACKETFEE'].fields_by_name['fee']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGPAYPACKETFEE']._loaded_options = None
  _globals['_MSGPAYPACKETFEE']._serialized_options = b'\210\240\037\000\202\347\260*\006signer\212\347\260*\032cosmos-sdk/MsgPayPacketFee'
  _globals['_MSGPAYPACKETFEEASYNC'].fields_by_name['packet_id']._loaded_options = None
  _globals['_MSGPAYPACKETFEEASYNC'].fields_by_name['packet_id']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGPAYPACKETFEEASYNC'].fields_by_name['packet_fee']._loaded_options = None
  _globals['_MSGPAYPACKETFEEASYNC'].fields_by_name['packet_fee']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGPAYPACKETFEEASYNC']._loaded_options = None
  _globals['_MSGPAYPACKETFEEASYNC']._serialized_options = b'\210\240\037\000\202\347\260*\npacket_fee\212\347\260*\037cosmos-sdk/MsgPayPacketFeeAsync'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGREGISTERPAYEE']._serialized_start=198
  _globals['_MSGREGISTERPAYEE']._serialized_end=370
  _globals['_MSGREGISTERPAYEERESPONSE']._serialized_start=372
  _globals['_MSGREGISTERPAYEERESPONSE']._serialized_end=398
  _globals['_MSGREGISTERCOUNTERPARTYPAYEE']._serialized_start=401
  _globals['_MSGREGISTERCOUNTERPARTYPAYEE']._serialized_end=622
  _globals['_MSGREGISTERCOUNTERPARTYPAYEERESPONSE']._serialized_start=624
  _globals['_MSGREGISTERCOUNTERPARTYPAYEERESPONSE']._serialized_end=662
  _globals['_MSGPAYPACKETFEE']._serialized_start=665
  _globals['_MSGPAYPACKETFEE']._serialized_end=923
  _globals['_MSGPAYPACKETFEERESPONSE']._serialized_start=925
  _globals['_MSGPAYPACKETFEERESPONSE']._serialized_end=950
  _globals['_MSGPAYPACKETFEEASYNC']._serialized_start=953
  _globals['_MSGPAYPACKETFEEASYNC']._serialized_end=1181
  _globals['_MSGPAYPACKETFEEASYNCRESPONSE']._serialized_start=1183
  _globals['_MSGPAYPACKETFEEASYNCRESPONSE']._serialized_end=1213
  _globals['_MSG']._serialized_start=1216
  _globals['_MSG']._serialized_end=1718
# @@protoc_insertion_point(module_scope)
