# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/core/channel/v2/tx.proto
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
    'ibc/core/channel/v2/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from ibc.core.channel.v2 import packet_pb2 as ibc_dot_core_dot_channel_dot_v2_dot_packet__pb2
from ibc.core.client.v1 import client_pb2 as ibc_dot_core_dot_client_dot_v1_dot_client__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1cibc/core/channel/v2/tx.proto\x12\x13ibc.core.channel.v2\x1a\x14gogoproto/gogo.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a ibc/core/channel/v2/packet.proto\x1a\x1fibc/core/client/v1/client.proto\"\xca\x01\n\rMsgSendPacket\x12#\n\rsource_client\x18\x01 \x01(\tR\x0csourceClient\x12+\n\x11timeout_timestamp\x18\x02 \x01(\x04R\x10timeoutTimestamp\x12>\n\x08payloads\x18\x03 \x03(\x0b\x32\x1c.ibc.core.channel.v2.PayloadB\x04\xc8\xde\x1f\x00R\x08payloads\x12\x16\n\x06signer\x18\x04 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"9\n\x15MsgSendPacketResponse\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequence:\x04\x88\xa0\x1f\x00\"\xe3\x01\n\rMsgRecvPacket\x12\x39\n\x06packet\x18\x01 \x01(\x0b\x32\x1b.ibc.core.channel.v2.PacketB\x04\xc8\xde\x1f\x00R\x06packet\x12)\n\x10proof_commitment\x18\x02 \x01(\x0cR\x0fproofCommitment\x12\x43\n\x0cproof_height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\x12\x16\n\x06signer\x18\x04 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"^\n\x15MsgRecvPacketResponse\x12?\n\x06result\x18\x01 \x01(\x0e\x32\'.ibc.core.channel.v2.ResponseResultTypeR\x06result:\x04\x88\xa0\x1f\x00\"\xe0\x01\n\nMsgTimeout\x12\x39\n\x06packet\x18\x01 \x01(\x0b\x32\x1b.ibc.core.channel.v2.PacketB\x04\xc8\xde\x1f\x00R\x06packet\x12)\n\x10proof_unreceived\x18\x02 \x01(\x0cR\x0fproofUnreceived\x12\x43\n\x0cproof_height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\x12\x16\n\x06signer\x18\x05 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"[\n\x12MsgTimeoutResponse\x12?\n\x06result\x18\x01 \x01(\x0e\x32\'.ibc.core.channel.v2.ResponseResultTypeR\x06result:\x04\x88\xa0\x1f\x00\"\xb4\x02\n\x12MsgAcknowledgement\x12\x39\n\x06packet\x18\x01 \x01(\x0b\x32\x1b.ibc.core.channel.v2.PacketB\x04\xc8\xde\x1f\x00R\x06packet\x12T\n\x0f\x61\x63knowledgement\x18\x02 \x01(\x0b\x32$.ibc.core.channel.v2.AcknowledgementB\x04\xc8\xde\x1f\x00R\x0f\x61\x63knowledgement\x12\x1f\n\x0bproof_acked\x18\x03 \x01(\x0cR\nproofAcked\x12\x43\n\x0cproof_height\x18\x04 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\x12\x16\n\x06signer\x18\x05 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"c\n\x1aMsgAcknowledgementResponse\x12?\n\x06result\x18\x01 \x01(\x0e\x32\'.ibc.core.channel.v2.ResponseResultTypeR\x06result:\x04\x88\xa0\x1f\x00*\xd8\x01\n\x12ResponseResultType\x12\x35\n RESPONSE_RESULT_TYPE_UNSPECIFIED\x10\x00\x1a\x0f\x8a\x9d \x0bUNSPECIFIED\x12\'\n\x19RESPONSE_RESULT_TYPE_NOOP\x10\x01\x1a\x08\x8a\x9d \x04NOOP\x12-\n\x1cRESPONSE_RESULT_TYPE_SUCCESS\x10\x02\x1a\x0b\x8a\x9d \x07SUCCESS\x12-\n\x1cRESPONSE_RESULT_TYPE_FAILURE\x10\x03\x1a\x0b\x8a\x9d \x07\x46\x41ILURE\x1a\x04\x88\xa3\x1e\x00\x32\x8a\x03\n\x03Msg\x12\\\n\nSendPacket\x12\".ibc.core.channel.v2.MsgSendPacket\x1a*.ibc.core.channel.v2.MsgSendPacketResponse\x12\\\n\nRecvPacket\x12\".ibc.core.channel.v2.MsgRecvPacket\x1a*.ibc.core.channel.v2.MsgRecvPacketResponse\x12S\n\x07Timeout\x12\x1f.ibc.core.channel.v2.MsgTimeout\x1a\'.ibc.core.channel.v2.MsgTimeoutResponse\x12k\n\x0f\x41\x63knowledgement\x12\'.ibc.core.channel.v2.MsgAcknowledgement\x1a/.ibc.core.channel.v2.MsgAcknowledgementResponse\x1a\x05\x80\xe7\xb0*\x01\x42?Z=github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.core.channel.v2.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z=github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types'
  _globals['_RESPONSERESULTTYPE']._loaded_options = None
  _globals['_RESPONSERESULTTYPE']._serialized_options = b'\210\243\036\000'
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_UNSPECIFIED"]._loaded_options = None
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_UNSPECIFIED"]._serialized_options = b'\212\235 \013UNSPECIFIED'
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_NOOP"]._loaded_options = None
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_NOOP"]._serialized_options = b'\212\235 \004NOOP'
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_SUCCESS"]._loaded_options = None
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_SUCCESS"]._serialized_options = b'\212\235 \007SUCCESS'
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_FAILURE"]._loaded_options = None
  _globals['_RESPONSERESULTTYPE'].values_by_name["RESPONSE_RESULT_TYPE_FAILURE"]._serialized_options = b'\212\235 \007FAILURE'
  _globals['_MSGSENDPACKET'].fields_by_name['payloads']._loaded_options = None
  _globals['_MSGSENDPACKET'].fields_by_name['payloads']._serialized_options = b'\310\336\037\000'
  _globals['_MSGSENDPACKET']._loaded_options = None
  _globals['_MSGSENDPACKET']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGSENDPACKETRESPONSE']._loaded_options = None
  _globals['_MSGSENDPACKETRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSGRECVPACKET'].fields_by_name['packet']._loaded_options = None
  _globals['_MSGRECVPACKET'].fields_by_name['packet']._serialized_options = b'\310\336\037\000'
  _globals['_MSGRECVPACKET'].fields_by_name['proof_height']._loaded_options = None
  _globals['_MSGRECVPACKET'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_MSGRECVPACKET']._loaded_options = None
  _globals['_MSGRECVPACKET']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGRECVPACKETRESPONSE']._loaded_options = None
  _globals['_MSGRECVPACKETRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSGTIMEOUT'].fields_by_name['packet']._loaded_options = None
  _globals['_MSGTIMEOUT'].fields_by_name['packet']._serialized_options = b'\310\336\037\000'
  _globals['_MSGTIMEOUT'].fields_by_name['proof_height']._loaded_options = None
  _globals['_MSGTIMEOUT'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_MSGTIMEOUT']._loaded_options = None
  _globals['_MSGTIMEOUT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGTIMEOUTRESPONSE']._loaded_options = None
  _globals['_MSGTIMEOUTRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSGACKNOWLEDGEMENT'].fields_by_name['packet']._loaded_options = None
  _globals['_MSGACKNOWLEDGEMENT'].fields_by_name['packet']._serialized_options = b'\310\336\037\000'
  _globals['_MSGACKNOWLEDGEMENT'].fields_by_name['acknowledgement']._loaded_options = None
  _globals['_MSGACKNOWLEDGEMENT'].fields_by_name['acknowledgement']._serialized_options = b'\310\336\037\000'
  _globals['_MSGACKNOWLEDGEMENT'].fields_by_name['proof_height']._loaded_options = None
  _globals['_MSGACKNOWLEDGEMENT'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_MSGACKNOWLEDGEMENT']._loaded_options = None
  _globals['_MSGACKNOWLEDGEMENT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGACKNOWLEDGEMENTRESPONSE']._loaded_options = None
  _globals['_MSGACKNOWLEDGEMENTRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_RESPONSERESULTTYPE']._serialized_start=1490
  _globals['_RESPONSERESULTTYPE']._serialized_end=1706
  _globals['_MSGSENDPACKET']._serialized_start=168
  _globals['_MSGSENDPACKET']._serialized_end=370
  _globals['_MSGSENDPACKETRESPONSE']._serialized_start=372
  _globals['_MSGSENDPACKETRESPONSE']._serialized_end=429
  _globals['_MSGRECVPACKET']._serialized_start=432
  _globals['_MSGRECVPACKET']._serialized_end=659
  _globals['_MSGRECVPACKETRESPONSE']._serialized_start=661
  _globals['_MSGRECVPACKETRESPONSE']._serialized_end=755
  _globals['_MSGTIMEOUT']._serialized_start=758
  _globals['_MSGTIMEOUT']._serialized_end=982
  _globals['_MSGTIMEOUTRESPONSE']._serialized_start=984
  _globals['_MSGTIMEOUTRESPONSE']._serialized_end=1075
  _globals['_MSGACKNOWLEDGEMENT']._serialized_start=1078
  _globals['_MSGACKNOWLEDGEMENT']._serialized_end=1386
  _globals['_MSGACKNOWLEDGEMENTRESPONSE']._serialized_start=1388
  _globals['_MSGACKNOWLEDGEMENTRESPONSE']._serialized_end=1487
  _globals['_MSG']._serialized_start=1709
  _globals['_MSG']._serialized_end=2103
# @@protoc_insertion_point(module_scope)
