# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/core/client/v1/tx.proto
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
    'ibc/core/client/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos.upgrade.v1beta1 import upgrade_pb2 as cosmos_dot_upgrade_dot_v1beta1_dot_upgrade__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from ibc.core.client.v1 import client_pb2 as ibc_dot_core_dot_client_dot_v1_dot_client__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1bibc/core/client/v1/tx.proto\x12\x12ibc.core.client.v1\x1a\x17\x63osmos/msg/v1/msg.proto\x1a$cosmos/upgrade/v1beta1/upgrade.proto\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a\x1fibc/core/client/v1/client.proto\"\xb2\x01\n\x0fMsgCreateClient\x12\x37\n\x0c\x63lient_state\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0b\x63lientState\x12=\n\x0f\x63onsensus_state\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0e\x63onsensusState\x12\x16\n\x06signer\x18\x03 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"<\n\x17MsgCreateClientResponse\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId:\x04\x88\xa0\x1f\x00\"\x94\x01\n\x0fMsgUpdateClient\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\x12;\n\x0e\x63lient_message\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\rclientMessage\x12\x16\n\x06signer\x18\x03 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x19\n\x17MsgUpdateClientResponse\"\xc5\x02\n\x10MsgUpgradeClient\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\x12\x37\n\x0c\x63lient_state\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0b\x63lientState\x12=\n\x0f\x63onsensus_state\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0e\x63onsensusState\x12\x30\n\x14proof_upgrade_client\x18\x04 \x01(\x0cR\x12proofUpgradeClient\x12\x41\n\x1dproof_upgrade_consensus_state\x18\x05 \x01(\x0cR\x1aproofUpgradeConsensusState\x12\x16\n\x06signer\x18\x06 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x1a\n\x18MsgUpgradeClientResponse\"\x99\x01\n\x15MsgSubmitMisbehaviour\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\x12\x38\n\x0cmisbehaviour\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0cmisbehaviour\x12\x16\n\x06signer\x18\x03 \x01(\tR\x06signer:\x11\x18\x01\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x1f\n\x1dMsgSubmitMisbehaviourResponse\"\x99\x01\n\x10MsgRecoverClient\x12*\n\x11subject_client_id\x18\x01 \x01(\tR\x0fsubjectClientId\x12\x30\n\x14substitute_client_id\x18\x02 \x01(\tR\x12substituteClientId\x12\x16\n\x06signer\x18\x03 \x01(\tR\x06signer:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x1a\n\x18MsgRecoverClientResponse\"\xbe\x01\n\x15MsgIBCSoftwareUpgrade\x12\x36\n\x04plan\x18\x01 \x01(\x0b\x32\x1c.cosmos.upgrade.v1beta1.PlanB\x04\xc8\xde\x1f\x00R\x04plan\x12H\n\x15upgraded_client_state\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x13upgradedClientState\x12\x16\n\x06signer\x18\x03 \x01(\tR\x06signer:\x0b\x82\xe7\xb0*\x06signer\"\x1f\n\x1dMsgIBCSoftwareUpgradeResponse\"t\n\x0fMsgUpdateParams\x12\x16\n\x06signer\x18\x01 \x01(\tR\x06signer\x12\x38\n\x06params\x18\x02 \x01(\x0b\x32\x1a.ibc.core.client.v1.ParamsB\x04\xc8\xde\x1f\x00R\x06params:\x0f\x88\xa0\x1f\x00\x82\xe7\xb0*\x06signer\"\x19\n\x17MsgUpdateParamsResponse2\xea\x05\n\x03Msg\x12`\n\x0c\x43reateClient\x12#.ibc.core.client.v1.MsgCreateClient\x1a+.ibc.core.client.v1.MsgCreateClientResponse\x12`\n\x0cUpdateClient\x12#.ibc.core.client.v1.MsgUpdateClient\x1a+.ibc.core.client.v1.MsgUpdateClientResponse\x12\x63\n\rUpgradeClient\x12$.ibc.core.client.v1.MsgUpgradeClient\x1a,.ibc.core.client.v1.MsgUpgradeClientResponse\x12r\n\x12SubmitMisbehaviour\x12).ibc.core.client.v1.MsgSubmitMisbehaviour\x1a\x31.ibc.core.client.v1.MsgSubmitMisbehaviourResponse\x12\x63\n\rRecoverClient\x12$.ibc.core.client.v1.MsgRecoverClient\x1a,.ibc.core.client.v1.MsgRecoverClientResponse\x12r\n\x12IBCSoftwareUpgrade\x12).ibc.core.client.v1.MsgIBCSoftwareUpgrade\x1a\x31.ibc.core.client.v1.MsgIBCSoftwareUpgradeResponse\x12\x66\n\x12UpdateClientParams\x12#.ibc.core.client.v1.MsgUpdateParams\x1a+.ibc.core.client.v1.MsgUpdateParamsResponse\x1a\x05\x80\xe7\xb0*\x01\x42:Z8github.com/cosmos/ibc-go/v9/modules/core/02-client/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.core.client.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z8github.com/cosmos/ibc-go/v9/modules/core/02-client/types'
  _globals['_MSGCREATECLIENT']._loaded_options = None
  _globals['_MSGCREATECLIENT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGCREATECLIENTRESPONSE']._loaded_options = None
  _globals['_MSGCREATECLIENTRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_MSGUPDATECLIENT']._loaded_options = None
  _globals['_MSGUPDATECLIENT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGUPGRADECLIENT']._loaded_options = None
  _globals['_MSGUPGRADECLIENT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGSUBMITMISBEHAVIOUR']._loaded_options = None
  _globals['_MSGSUBMITMISBEHAVIOUR']._serialized_options = b'\030\001\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGRECOVERCLIENT']._loaded_options = None
  _globals['_MSGRECOVERCLIENT']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSGIBCSOFTWAREUPGRADE'].fields_by_name['plan']._loaded_options = None
  _globals['_MSGIBCSOFTWAREUPGRADE'].fields_by_name['plan']._serialized_options = b'\310\336\037\000'
  _globals['_MSGIBCSOFTWAREUPGRADE']._loaded_options = None
  _globals['_MSGIBCSOFTWAREUPGRADE']._serialized_options = b'\202\347\260*\006signer'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\210\240\037\000\202\347\260*\006signer'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGCREATECLIENT']._serialized_start=197
  _globals['_MSGCREATECLIENT']._serialized_end=375
  _globals['_MSGCREATECLIENTRESPONSE']._serialized_start=377
  _globals['_MSGCREATECLIENTRESPONSE']._serialized_end=437
  _globals['_MSGUPDATECLIENT']._serialized_start=440
  _globals['_MSGUPDATECLIENT']._serialized_end=588
  _globals['_MSGUPDATECLIENTRESPONSE']._serialized_start=590
  _globals['_MSGUPDATECLIENTRESPONSE']._serialized_end=615
  _globals['_MSGUPGRADECLIENT']._serialized_start=618
  _globals['_MSGUPGRADECLIENT']._serialized_end=943
  _globals['_MSGUPGRADECLIENTRESPONSE']._serialized_start=945
  _globals['_MSGUPGRADECLIENTRESPONSE']._serialized_end=971
  _globals['_MSGSUBMITMISBEHAVIOUR']._serialized_start=974
  _globals['_MSGSUBMITMISBEHAVIOUR']._serialized_end=1127
  _globals['_MSGSUBMITMISBEHAVIOURRESPONSE']._serialized_start=1129
  _globals['_MSGSUBMITMISBEHAVIOURRESPONSE']._serialized_end=1160
  _globals['_MSGRECOVERCLIENT']._serialized_start=1163
  _globals['_MSGRECOVERCLIENT']._serialized_end=1316
  _globals['_MSGRECOVERCLIENTRESPONSE']._serialized_start=1318
  _globals['_MSGRECOVERCLIENTRESPONSE']._serialized_end=1344
  _globals['_MSGIBCSOFTWAREUPGRADE']._serialized_start=1347
  _globals['_MSGIBCSOFTWAREUPGRADE']._serialized_end=1537
  _globals['_MSGIBCSOFTWAREUPGRADERESPONSE']._serialized_start=1539
  _globals['_MSGIBCSOFTWAREUPGRADERESPONSE']._serialized_end=1570
  _globals['_MSGUPDATEPARAMS']._serialized_start=1572
  _globals['_MSGUPDATEPARAMS']._serialized_end=1688
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=1690
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=1715
  _globals['_MSG']._serialized_start=1718
  _globals['_MSG']._serialized_end=2464
# @@protoc_insertion_point(module_scope)
