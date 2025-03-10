# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/core/connection/v1/tx.proto
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
    'ibc/core/connection/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from ibc.core.client.v1 import client_pb2 as ibc_dot_core_dot_client_dot_v1_dot_client__pb2
from ibc.core.connection.v1 import connection_pb2 as ibc_dot_core_dot_connection_dot_v1_dot_connection__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1fibc/core/connection/v1/tx.proto\x12\x16ibc.core.connection.v1\x1a\x14gogoproto/gogo.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x19google/protobuf/any.proto\x1a\x1fibc/core/client/v1/client.proto\x1a\'ibc/core/connection/v1/connection.proto\"\x8b\x02\n\x15MsgConnectionOpenInit\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\x12N\n\x0c\x63ounterparty\x18\x02 \x01(\x0b\x32$.ibc.core.connection.v1.CounterpartyB\x04\xc8\xde\x1f\x00R\x0c\x63ounterparty\x12\x39\n\x07version\x18\x03 \x01(\x0b\x32\x1f.ibc.core.connection.v1.VersionR\x07version\x12!\n\x0c\x64\x65lay_period\x18\x04 \x01(\x04R\x0b\x64\x65layPeriod\x12\x16\n\x06signer\x18\x05 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x1f\n\x1dMsgConnectionOpenInitResponse\"\xe4\x05\n\x14MsgConnectionOpenTry\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\x12\x38\n\x16previous_connection_id\x18\x02 \x01(\tB\x02\x18\x01R\x14previousConnectionId\x12;\n\x0c\x63lient_state\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyB\x02\x18\x01R\x0b\x63lientState\x12N\n\x0c\x63ounterparty\x18\x04 \x01(\x0b\x32$.ibc.core.connection.v1.CounterpartyB\x04\xc8\xde\x1f\x00R\x0c\x63ounterparty\x12!\n\x0c\x64\x65lay_period\x18\x05 \x01(\x04R\x0b\x64\x65layPeriod\x12T\n\x15\x63ounterparty_versions\x18\x06 \x03(\x0b\x32\x1f.ibc.core.connection.v1.VersionR\x14\x63ounterpartyVersions\x12\x43\n\x0cproof_height\x18\x07 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\x12\x1d\n\nproof_init\x18\x08 \x01(\x0cR\tproofInit\x12%\n\x0cproof_client\x18\t \x01(\x0c\x42\x02\x18\x01R\x0bproofClient\x12+\n\x0fproof_consensus\x18\n \x01(\x0c\x42\x02\x18\x01R\x0eproofConsensus\x12M\n\x10\x63onsensus_height\x18\x0b \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x06\x18\x01\xc8\xde\x1f\x00R\x0f\x63onsensusHeight\x12\x16\n\x06signer\x18\x0c \x01(\tR\x06signer\x12?\n\x1ahost_consensus_state_proof\x18\r \x01(\x0c\x42\x02\x18\x01R\x17hostConsensusStateProof:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x1e\n\x1cMsgConnectionOpenTryResponse\"\xe0\x04\n\x14MsgConnectionOpenAck\x12#\n\rconnection_id\x18\x01 \x01(\tR\x0c\x63onnectionId\x12<\n\x1a\x63ounterparty_connection_id\x18\x02 \x01(\tR\x18\x63ounterpartyConnectionId\x12\x39\n\x07version\x18\x03 \x01(\x0b\x32\x1f.ibc.core.connection.v1.VersionR\x07version\x12;\n\x0c\x63lient_state\x18\x04 \x01(\x0b\x32\x14.google.protobuf.AnyB\x02\x18\x01R\x0b\x63lientState\x12\x43\n\x0cproof_height\x18\x05 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\x12\x1b\n\tproof_try\x18\x06 \x01(\x0cR\x08proofTry\x12%\n\x0cproof_client\x18\x07 \x01(\x0c\x42\x02\x18\x01R\x0bproofClient\x12+\n\x0fproof_consensus\x18\x08 \x01(\x0c\x42\x02\x18\x01R\x0eproofConsensus\x12M\n\x10\x63onsensus_height\x18\t \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x06\x18\x01\xc8\xde\x1f\x00R\x0f\x63onsensusHeight\x12\x16\n\x06signer\x18\n \x01(\tR\x06signer\x12?\n\x1ahost_consensus_state_proof\x18\x0b \x01(\x0c\x42\x02\x18\x01R\x17hostConsensusStateProof:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x1e\n\x1cMsgConnectionOpenAckResponse\"\xca\x01\n\x18MsgConnectionOpenConfirm\x12#\n\rconnection_id\x18\x01 \x01(\tR\x0c\x63onnectionId\x12\x1b\n\tproof_ack\x18\x02 \x01(\x0cR\x08proofAck\x12\x43\n\x0cproof_height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\x12\x16\n\x06signer\x18\x04 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\"\n MsgConnectionOpenConfirmResponse\"x\n\x0fMsgUpdateParams\x12\x16\n\x06signer\x18\x01 \x01(\tR\x06signer\x12<\n\x06params\x18\x02 \x01(\x0b\x32\x1e.ibc.core.connection.v1.ParamsB\x04\xc8\xde\x1f\x00R\x06params:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x19\n\x17MsgUpdateParamsResponse2\xf4\x04\n\x03Msg\x12z\n\x12\x43onnectionOpenInit\x12-.ibc.core.connection.v1.MsgConnectionOpenInit\x1a\x35.ibc.core.connection.v1.MsgConnectionOpenInitResponse\x12w\n\x11\x43onnectionOpenTry\x12,.ibc.core.connection.v1.MsgConnectionOpenTry\x1a\x34.ibc.core.connection.v1.MsgConnectionOpenTryResponse\x12w\n\x11\x43onnectionOpenAck\x12,.ibc.core.connection.v1.MsgConnectionOpenAck\x1a\x34.ibc.core.connection.v1.MsgConnectionOpenAckResponse\x12\x83\x01\n\x15\x43onnectionOpenConfirm\x12\x30.ibc.core.connection.v1.MsgConnectionOpenConfirm\x1a\x38.ibc.core.connection.v1.MsgConnectionOpenConfirmResponse\x12r\n\x16UpdateConnectionParams\x12\'.ibc.core.connection.v1.MsgUpdateParams\x1a/.ibc.core.connection.v1.MsgUpdateParamsResponse\x1a\x05\x80\xe7\xb0*\x01\x42?Z=github.com/cosmos/ibc-go/v10/modules/core/03-connection/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.core.connection.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z=github.com/cosmos/ibc-go/v10/modules/core/03-connection/types'
  _globals['_MSGCONNECTIONOPENINIT'].fields_by_name['counterparty']._loaded_options = None
  _globals['_MSGCONNECTIONOPENINIT'].fields_by_name['counterparty']._serialized_options = b'\310\336\037\000'
  _globals['_MSGCONNECTIONOPENINIT']._loaded_options = None
  _globals['_MSGCONNECTIONOPENINIT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['previous_connection_id']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['previous_connection_id']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['client_state']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['client_state']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['counterparty']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['counterparty']._serialized_options = b'\310\336\037\000'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['proof_height']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['proof_client']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['proof_client']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['proof_consensus']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['proof_consensus']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['consensus_height']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['consensus_height']._serialized_options = b'\030\001\310\336\037\000'
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['host_consensus_state_proof']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY'].fields_by_name['host_consensus_state_proof']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENTRY']._loaded_options = None
  _globals['_MSGCONNECTIONOPENTRY']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['client_state']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['client_state']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['proof_height']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['proof_client']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['proof_client']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['proof_consensus']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['proof_consensus']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['consensus_height']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['consensus_height']._serialized_options = b'\030\001\310\336\037\000'
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['host_consensus_state_proof']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK'].fields_by_name['host_consensus_state_proof']._serialized_options = b'\030\001'
  _globals['_MSGCONNECTIONOPENACK']._loaded_options = None
  _globals['_MSGCONNECTIONOPENACK']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGCONNECTIONOPENCONFIRM'].fields_by_name['proof_height']._loaded_options = None
  _globals['_MSGCONNECTIONOPENCONFIRM'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_MSGCONNECTIONOPENCONFIRM']._loaded_options = None
  _globals['_MSGCONNECTIONOPENCONFIRM']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGCONNECTIONOPENINIT']._serialized_start=208
  _globals['_MSGCONNECTIONOPENINIT']._serialized_end=475
  _globals['_MSGCONNECTIONOPENINITRESPONSE']._serialized_start=477
  _globals['_MSGCONNECTIONOPENINITRESPONSE']._serialized_end=508
  _globals['_MSGCONNECTIONOPENTRY']._serialized_start=511
  _globals['_MSGCONNECTIONOPENTRY']._serialized_end=1251
  _globals['_MSGCONNECTIONOPENTRYRESPONSE']._serialized_start=1253
  _globals['_MSGCONNECTIONOPENTRYRESPONSE']._serialized_end=1283
  _globals['_MSGCONNECTIONOPENACK']._serialized_start=1286
  _globals['_MSGCONNECTIONOPENACK']._serialized_end=1894
  _globals['_MSGCONNECTIONOPENACKRESPONSE']._serialized_start=1896
  _globals['_MSGCONNECTIONOPENACKRESPONSE']._serialized_end=1926
  _globals['_MSGCONNECTIONOPENCONFIRM']._serialized_start=1929
  _globals['_MSGCONNECTIONOPENCONFIRM']._serialized_end=2131
  _globals['_MSGCONNECTIONOPENCONFIRMRESPONSE']._serialized_start=2133
  _globals['_MSGCONNECTIONOPENCONFIRMRESPONSE']._serialized_end=2167
  _globals['_MSGUPDATEPARAMS']._serialized_start=2169
  _globals['_MSGUPDATEPARAMS']._serialized_end=2289
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=2291
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=2316
  _globals['_MSG']._serialized_start=2319
  _globals['_MSG']._serialized_end=2947
# @@protoc_insertion_point(module_scope)
