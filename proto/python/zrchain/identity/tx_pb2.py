# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/identity/tx.proto
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
    'zrchain/identity/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from zrchain.identity import params_pb2 as zrchain_dot_identity_dot_params__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x19zrchain/identity/tx.proto\x12\x10zrchain.identity\x1a\x11\x61mino/amino.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x1dzrchain/identity/params.proto\"\xde\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12;\n\x06params\x18\x02 \x01(\x0b\x32\x18.zrchain.identity.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:V\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*Cgithub.com/Zenrock-Foundation/zrchain/v6/x/identity/MsgUpdateParams\"\x19\n\x17MsgUpdateParamsResponse\"\xb4\x01\n\x0fMsgNewWorkspace\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12&\n\x0f\x61\x64min_policy_id\x18\x02 \x01(\x04R\radminPolicyId\x12$\n\x0esign_policy_id\x18\x03 \x01(\x04R\x0csignPolicyId\x12+\n\x11\x61\x64\x64itional_owners\x18\x04 \x03(\tR\x10\x61\x64\x64itionalOwners:\x0c\x82\xe7\xb0*\x07\x63reator\"-\n\x17MsgNewWorkspaceResponse\x12\x12\n\x04\x61\x64\x64r\x18\x01 \x01(\tR\x04\x61\x64\x64r\"\x94\x01\n\x14MsgAddWorkspaceOwner\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12%\n\x0eworkspace_addr\x18\x02 \x01(\tR\rworkspaceAddr\x12\x1b\n\tnew_owner\x18\x03 \x01(\tR\x08newOwner\x12\x10\n\x03\x62tl\x18\x04 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1e\n\x1cMsgAddWorkspaceOwnerResponse\"\xb9\x01\n\x17MsgAppendChildWorkspace\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x32\n\x15parent_workspace_addr\x18\x02 \x01(\tR\x13parentWorkspaceAddr\x12\x30\n\x14\x63hild_workspace_addr\x18\x03 \x01(\tR\x12\x63hildWorkspaceAddr\x12\x10\n\x03\x62tl\x18\x04 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"!\n\x1fMsgAppendChildWorkspaceResponse\"\x84\x01\n\x14MsgNewChildWorkspace\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x32\n\x15parent_workspace_addr\x18\x02 \x01(\tR\x13parentWorkspaceAddr\x12\x10\n\x03\x62tl\x18\x03 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"8\n\x1cMsgNewChildWorkspaceResponse\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\"\x90\x01\n\x17MsgRemoveWorkspaceOwner\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12%\n\x0eworkspace_addr\x18\x02 \x01(\tR\rworkspaceAddr\x12\x14\n\x05owner\x18\x03 \x01(\tR\x05owner\x12\x10\n\x03\x62tl\x18\x04 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"!\n\x1fMsgRemoveWorkspaceOwnerResponse\"\xe7\x01\n\rMsgNewKeyring\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\'\n\x0fparty_threshold\x18\x03 \x01(\rR\x0epartyThreshold\x12\x1e\n\x0bkey_req_fee\x18\x04 \x01(\x04R\tkeyReqFee\x12\x1e\n\x0bsig_req_fee\x18\x05 \x01(\x04R\tsigReqFee\x12#\n\rdelegate_fees\x18\x06 \x01(\x08R\x0c\x64\x65legateFees:\x0c\x82\xe7\xb0*\x07\x63reator\"+\n\x15MsgNewKeyringResponse\x12\x12\n\x04\x61\x64\x64r\x18\x01 \x01(\tR\x04\x61\x64\x64r\"\xa4\x01\n\x12MsgAddKeyringParty\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr\x12\x14\n\x05party\x18\x03 \x01(\tR\x05party\x12-\n\x12increase_threshold\x18\x04 \x01(\x08R\x11increaseThreshold:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1c\n\x1aMsgAddKeyringPartyResponse\"\x85\x02\n\x10MsgUpdateKeyring\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr\x12\'\n\x0fparty_threshold\x18\x03 \x01(\rR\x0epartyThreshold\x12\x1e\n\x0bkey_req_fee\x18\x04 \x01(\x04R\tkeyReqFee\x12\x1e\n\x0bsig_req_fee\x18\x05 \x01(\x04R\tsigReqFee\x12 \n\x0b\x64\x65scription\x18\x06 \x01(\tR\x0b\x64\x65scription\x12\x1b\n\tis_active\x18\x07 \x01(\x08R\x08isActive:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1a\n\x18MsgUpdateKeyringResponse\"\xa7\x01\n\x15MsgRemoveKeyringParty\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr\x12\x14\n\x05party\x18\x03 \x01(\tR\x05party\x12-\n\x12\x64\x65\x63rease_threshold\x18\x04 \x01(\x08R\x11\x64\x65\x63reaseThreshold:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1f\n\x1dMsgRemoveKeyringPartyResponse\"u\n\x12MsgAddKeyringAdmin\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr\x12\x14\n\x05\x61\x64min\x18\x03 \x01(\tR\x05\x61\x64min:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1c\n\x1aMsgAddKeyringAdminResponse\"x\n\x15MsgRemoveKeyringAdmin\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr\x12\x14\n\x05\x61\x64min\x18\x03 \x01(\tR\x05\x61\x64min:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1f\n\x1dMsgRemoveKeyringAdminResponse\"\xc3\x01\n\x12MsgUpdateWorkspace\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12%\n\x0eworkspace_addr\x18\x02 \x01(\tR\rworkspaceAddr\x12&\n\x0f\x61\x64min_policy_id\x18\x03 \x01(\x04R\radminPolicyId\x12$\n\x0esign_policy_id\x18\x04 \x01(\x04R\x0csignPolicyId\x12\x10\n\x03\x62tl\x18\x05 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1c\n\x1aMsgUpdateWorkspaceResponse\"a\n\x14MsgDeactivateKeyring\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1e\n\x1cMsgDeactivateKeyringResponse2\xc9\x0b\n\x03Msg\x12\\\n\x0cUpdateParams\x12!.zrchain.identity.MsgUpdateParams\x1a).zrchain.identity.MsgUpdateParamsResponse\x12\\\n\x0cNewWorkspace\x12!.zrchain.identity.MsgNewWorkspace\x1a).zrchain.identity.MsgNewWorkspaceResponse\x12k\n\x11\x41\x64\x64WorkspaceOwner\x12&.zrchain.identity.MsgAddWorkspaceOwner\x1a..zrchain.identity.MsgAddWorkspaceOwnerResponse\x12t\n\x14\x41ppendChildWorkspace\x12).zrchain.identity.MsgAppendChildWorkspace\x1a\x31.zrchain.identity.MsgAppendChildWorkspaceResponse\x12k\n\x11NewChildWorkspace\x12&.zrchain.identity.MsgNewChildWorkspace\x1a..zrchain.identity.MsgNewChildWorkspaceResponse\x12t\n\x14RemoveWorkspaceOwner\x12).zrchain.identity.MsgRemoveWorkspaceOwner\x1a\x31.zrchain.identity.MsgRemoveWorkspaceOwnerResponse\x12V\n\nNewKeyring\x12\x1f.zrchain.identity.MsgNewKeyring\x1a\'.zrchain.identity.MsgNewKeyringResponse\x12\x65\n\x0f\x41\x64\x64KeyringParty\x12$.zrchain.identity.MsgAddKeyringParty\x1a,.zrchain.identity.MsgAddKeyringPartyResponse\x12_\n\rUpdateKeyring\x12\".zrchain.identity.MsgUpdateKeyring\x1a*.zrchain.identity.MsgUpdateKeyringResponse\x12n\n\x12RemoveKeyringParty\x12\'.zrchain.identity.MsgRemoveKeyringParty\x1a/.zrchain.identity.MsgRemoveKeyringPartyResponse\x12\x65\n\x0f\x41\x64\x64KeyringAdmin\x12$.zrchain.identity.MsgAddKeyringAdmin\x1a,.zrchain.identity.MsgAddKeyringAdminResponse\x12n\n\x12RemoveKeyringAdmin\x12\'.zrchain.identity.MsgRemoveKeyringAdmin\x1a/.zrchain.identity.MsgRemoveKeyringAdminResponse\x12\x65\n\x0fUpdateWorkspace\x12$.zrchain.identity.MsgUpdateWorkspace\x1a,.zrchain.identity.MsgUpdateWorkspaceResponse\x12k\n\x11\x44\x65\x61\x63tivateKeyring\x12&.zrchain.identity.MsgDeactivateKeyring\x1a..zrchain.identity.MsgDeactivateKeyringResponse\x1a\x05\x80\xe7\xb0*\x01\x42;Z9github.com/Zenrock-Foundation/zrchain/v6/x/identity/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.identity.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v6/x/identity/types'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\202\347\260*\tauthority\212\347\260*Cgithub.com/Zenrock-Foundation/zrchain/v6/x/identity/MsgUpdateParams'
  _globals['_MSGNEWWORKSPACE']._loaded_options = None
  _globals['_MSGNEWWORKSPACE']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGADDWORKSPACEOWNER']._loaded_options = None
  _globals['_MSGADDWORKSPACEOWNER']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGAPPENDCHILDWORKSPACE']._loaded_options = None
  _globals['_MSGAPPENDCHILDWORKSPACE']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGNEWCHILDWORKSPACE']._loaded_options = None
  _globals['_MSGNEWCHILDWORKSPACE']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGREMOVEWORKSPACEOWNER']._loaded_options = None
  _globals['_MSGREMOVEWORKSPACEOWNER']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGNEWKEYRING']._loaded_options = None
  _globals['_MSGNEWKEYRING']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGADDKEYRINGPARTY']._loaded_options = None
  _globals['_MSGADDKEYRINGPARTY']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGUPDATEKEYRING']._loaded_options = None
  _globals['_MSGUPDATEKEYRING']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGREMOVEKEYRINGPARTY']._loaded_options = None
  _globals['_MSGREMOVEKEYRINGPARTY']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGADDKEYRINGADMIN']._loaded_options = None
  _globals['_MSGADDKEYRINGADMIN']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGREMOVEKEYRINGADMIN']._loaded_options = None
  _globals['_MSGREMOVEKEYRINGADMIN']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGUPDATEWORKSPACE']._loaded_options = None
  _globals['_MSGUPDATEWORKSPACE']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGDEACTIVATEKEYRING']._loaded_options = None
  _globals['_MSGDEACTIVATEKEYRING']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._serialized_start=172
  _globals['_MSGUPDATEPARAMS']._serialized_end=394
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=396
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=421
  _globals['_MSGNEWWORKSPACE']._serialized_start=424
  _globals['_MSGNEWWORKSPACE']._serialized_end=604
  _globals['_MSGNEWWORKSPACERESPONSE']._serialized_start=606
  _globals['_MSGNEWWORKSPACERESPONSE']._serialized_end=651
  _globals['_MSGADDWORKSPACEOWNER']._serialized_start=654
  _globals['_MSGADDWORKSPACEOWNER']._serialized_end=802
  _globals['_MSGADDWORKSPACEOWNERRESPONSE']._serialized_start=804
  _globals['_MSGADDWORKSPACEOWNERRESPONSE']._serialized_end=834
  _globals['_MSGAPPENDCHILDWORKSPACE']._serialized_start=837
  _globals['_MSGAPPENDCHILDWORKSPACE']._serialized_end=1022
  _globals['_MSGAPPENDCHILDWORKSPACERESPONSE']._serialized_start=1024
  _globals['_MSGAPPENDCHILDWORKSPACERESPONSE']._serialized_end=1057
  _globals['_MSGNEWCHILDWORKSPACE']._serialized_start=1060
  _globals['_MSGNEWCHILDWORKSPACE']._serialized_end=1192
  _globals['_MSGNEWCHILDWORKSPACERESPONSE']._serialized_start=1194
  _globals['_MSGNEWCHILDWORKSPACERESPONSE']._serialized_end=1250
  _globals['_MSGREMOVEWORKSPACEOWNER']._serialized_start=1253
  _globals['_MSGREMOVEWORKSPACEOWNER']._serialized_end=1397
  _globals['_MSGREMOVEWORKSPACEOWNERRESPONSE']._serialized_start=1399
  _globals['_MSGREMOVEWORKSPACEOWNERRESPONSE']._serialized_end=1432
  _globals['_MSGNEWKEYRING']._serialized_start=1435
  _globals['_MSGNEWKEYRING']._serialized_end=1666
  _globals['_MSGNEWKEYRINGRESPONSE']._serialized_start=1668
  _globals['_MSGNEWKEYRINGRESPONSE']._serialized_end=1711
  _globals['_MSGADDKEYRINGPARTY']._serialized_start=1714
  _globals['_MSGADDKEYRINGPARTY']._serialized_end=1878
  _globals['_MSGADDKEYRINGPARTYRESPONSE']._serialized_start=1880
  _globals['_MSGADDKEYRINGPARTYRESPONSE']._serialized_end=1908
  _globals['_MSGUPDATEKEYRING']._serialized_start=1911
  _globals['_MSGUPDATEKEYRING']._serialized_end=2172
  _globals['_MSGUPDATEKEYRINGRESPONSE']._serialized_start=2174
  _globals['_MSGUPDATEKEYRINGRESPONSE']._serialized_end=2200
  _globals['_MSGREMOVEKEYRINGPARTY']._serialized_start=2203
  _globals['_MSGREMOVEKEYRINGPARTY']._serialized_end=2370
  _globals['_MSGREMOVEKEYRINGPARTYRESPONSE']._serialized_start=2372
  _globals['_MSGREMOVEKEYRINGPARTYRESPONSE']._serialized_end=2403
  _globals['_MSGADDKEYRINGADMIN']._serialized_start=2405
  _globals['_MSGADDKEYRINGADMIN']._serialized_end=2522
  _globals['_MSGADDKEYRINGADMINRESPONSE']._serialized_start=2524
  _globals['_MSGADDKEYRINGADMINRESPONSE']._serialized_end=2552
  _globals['_MSGREMOVEKEYRINGADMIN']._serialized_start=2554
  _globals['_MSGREMOVEKEYRINGADMIN']._serialized_end=2674
  _globals['_MSGREMOVEKEYRINGADMINRESPONSE']._serialized_start=2676
  _globals['_MSGREMOVEKEYRINGADMINRESPONSE']._serialized_end=2707
  _globals['_MSGUPDATEWORKSPACE']._serialized_start=2710
  _globals['_MSGUPDATEWORKSPACE']._serialized_end=2905
  _globals['_MSGUPDATEWORKSPACERESPONSE']._serialized_start=2907
  _globals['_MSGUPDATEWORKSPACERESPONSE']._serialized_end=2935
  _globals['_MSGDEACTIVATEKEYRING']._serialized_start=2937
  _globals['_MSGDEACTIVATEKEYRING']._serialized_end=3034
  _globals['_MSGDEACTIVATEKEYRINGRESPONSE']._serialized_start=3036
  _globals['_MSGDEACTIVATEKEYRINGRESPONSE']._serialized_end=3066
  _globals['_MSG']._serialized_start=3069
  _globals['_MSG']._serialized_end=4550
# @@protoc_insertion_point(module_scope)
