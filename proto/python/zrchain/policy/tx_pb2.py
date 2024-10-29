# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/policy/tx.proto
# Protobuf Python Version: 5.28.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    2,
    '',
    'zrchain/policy/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from zrchain.policy import params_pb2 as zrchain_dot_policy_dot_params__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x17zrchain/policy/tx.proto\x12\x0ezrchain.policy\x1a\x11\x61mino/amino.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x1bzrchain/policy/params.proto\x1a\x19google/protobuf/any.proto\"\xb9\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x39\n\x06params\x18\x02 \x01(\x0b\x32\x16.zrchain.policy.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:3\x82\xe7\xb0*\tauthority\x8a\xe7\xb0* zrchain/x/policy/MsgUpdateParams\"\x19\n\x17MsgUpdateParamsResponse\"\x8a\x01\n\x0cMsgNewPolicy\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x12\n\x04name\x18\x02 \x01(\tR\x04name\x12,\n\x06policy\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyR\x06policy\x12\x10\n\x03\x62tl\x18\x04 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"&\n\x14MsgNewPolicyResponse\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\"V\n\x0fMsgRevokeAction\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x1b\n\taction_id\x18\x02 \x01(\x04R\x08\x61\x63tionId:\x0c\x82\xe7\xb0*\x07\x63reator\"\x19\n\x17MsgRevokeActionResponse\"\xc3\x01\n\x10MsgApproveAction\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x1f\n\x0b\x61\x63tion_type\x18\x02 \x01(\tR\nactionType\x12\x1b\n\taction_id\x18\x03 \x01(\x04R\x08\x61\x63tionId\x12I\n\x15\x61\x64\x64itional_signatures\x18\x04 \x03(\x0b\x32\x14.google.protobuf.AnyR\x14\x61\x64\x64itionalSignatures:\x0c\x82\xe7\xb0*\x07\x63reator\"2\n\x18MsgApproveActionResponse\x12\x16\n\x06status\x18\x01 \x01(\tR\x06status\"h\n\x10MsgAddSignMethod\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12,\n\x06\x63onfig\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x06\x63onfig:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1a\n\x18MsgAddSignMethodResponse\"M\n\x13MsgRemoveSignMethod\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x0e\n\x02id\x18\x02 \x01(\tR\x02id:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1d\n\x1bMsgRemoveSignMethodResponse\"h\n\x10MsgAddMultiGrant\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x18\n\x07grantee\x18\x02 \x01(\tR\x07grantee\x12\x12\n\x04msgs\x18\x03 \x03(\tR\x04msgs:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1a\n\x18MsgAddMultiGrantResponse\"k\n\x13MsgRemoveMultiGrant\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x18\n\x07grantee\x18\x02 \x01(\tR\x07grantee\x12\x12\n\x04msgs\x18\x03 \x03(\tR\x04msgs:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1d\n\x1bMsgRemoveMultiGrantResponse2\xf4\x05\n\x03Msg\x12X\n\x0cUpdateParams\x12\x1f.zrchain.policy.MsgUpdateParams\x1a\'.zrchain.policy.MsgUpdateParamsResponse\x12O\n\tNewPolicy\x12\x1c.zrchain.policy.MsgNewPolicy\x1a$.zrchain.policy.MsgNewPolicyResponse\x12X\n\x0cRevokeAction\x12\x1f.zrchain.policy.MsgRevokeAction\x1a\'.zrchain.policy.MsgRevokeActionResponse\x12[\n\rApproveAction\x12 .zrchain.policy.MsgApproveAction\x1a(.zrchain.policy.MsgApproveActionResponse\x12[\n\rAddSignMethod\x12 .zrchain.policy.MsgAddSignMethod\x1a(.zrchain.policy.MsgAddSignMethodResponse\x12\x64\n\x10RemoveSignMethod\x12#.zrchain.policy.MsgRemoveSignMethod\x1a+.zrchain.policy.MsgRemoveSignMethodResponse\x12[\n\rAddMultiGrant\x12 .zrchain.policy.MsgAddMultiGrant\x1a(.zrchain.policy.MsgAddMultiGrantResponse\x12\x64\n\x10RemoveMultiGrant\x12#.zrchain.policy.MsgRemoveMultiGrant\x1a+.zrchain.policy.MsgRemoveMultiGrantResponse\x1a\x05\x80\xe7\xb0*\x01\x42\x39Z7github.com/Zenrock-Foundation/zrchain/v4/x/policy/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.policy.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/Zenrock-Foundation/zrchain/v4/x/policy/types'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\202\347\260*\tauthority\212\347\260* zrchain/x/policy/MsgUpdateParams'
  _globals['_MSGNEWPOLICY']._loaded_options = None
  _globals['_MSGNEWPOLICY']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGREVOKEACTION']._loaded_options = None
  _globals['_MSGREVOKEACTION']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGAPPROVEACTION']._loaded_options = None
  _globals['_MSGAPPROVEACTION']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGADDSIGNMETHOD']._loaded_options = None
  _globals['_MSGADDSIGNMETHOD']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGREMOVESIGNMETHOD']._loaded_options = None
  _globals['_MSGREMOVESIGNMETHOD']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGADDMULTIGRANT']._loaded_options = None
  _globals['_MSGADDMULTIGRANT']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGREMOVEMULTIGRANT']._loaded_options = None
  _globals['_MSGREMOVEMULTIGRANT']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._serialized_start=193
  _globals['_MSGUPDATEPARAMS']._serialized_end=378
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=380
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=405
  _globals['_MSGNEWPOLICY']._serialized_start=408
  _globals['_MSGNEWPOLICY']._serialized_end=546
  _globals['_MSGNEWPOLICYRESPONSE']._serialized_start=548
  _globals['_MSGNEWPOLICYRESPONSE']._serialized_end=586
  _globals['_MSGREVOKEACTION']._serialized_start=588
  _globals['_MSGREVOKEACTION']._serialized_end=674
  _globals['_MSGREVOKEACTIONRESPONSE']._serialized_start=676
  _globals['_MSGREVOKEACTIONRESPONSE']._serialized_end=701
  _globals['_MSGAPPROVEACTION']._serialized_start=704
  _globals['_MSGAPPROVEACTION']._serialized_end=899
  _globals['_MSGAPPROVEACTIONRESPONSE']._serialized_start=901
  _globals['_MSGAPPROVEACTIONRESPONSE']._serialized_end=951
  _globals['_MSGADDSIGNMETHOD']._serialized_start=953
  _globals['_MSGADDSIGNMETHOD']._serialized_end=1057
  _globals['_MSGADDSIGNMETHODRESPONSE']._serialized_start=1059
  _globals['_MSGADDSIGNMETHODRESPONSE']._serialized_end=1085
  _globals['_MSGREMOVESIGNMETHOD']._serialized_start=1087
  _globals['_MSGREMOVESIGNMETHOD']._serialized_end=1164
  _globals['_MSGREMOVESIGNMETHODRESPONSE']._serialized_start=1166
  _globals['_MSGREMOVESIGNMETHODRESPONSE']._serialized_end=1195
  _globals['_MSGADDMULTIGRANT']._serialized_start=1197
  _globals['_MSGADDMULTIGRANT']._serialized_end=1301
  _globals['_MSGADDMULTIGRANTRESPONSE']._serialized_start=1303
  _globals['_MSGADDMULTIGRANTRESPONSE']._serialized_end=1329
  _globals['_MSGREMOVEMULTIGRANT']._serialized_start=1331
  _globals['_MSGREMOVEMULTIGRANT']._serialized_end=1438
  _globals['_MSGREMOVEMULTIGRANTRESPONSE']._serialized_start=1440
  _globals['_MSGREMOVEMULTIGRANTRESPONSE']._serialized_end=1469
  _globals['_MSG']._serialized_start=1472
  _globals['_MSG']._serialized_end=2228
# @@protoc_insertion_point(module_scope)