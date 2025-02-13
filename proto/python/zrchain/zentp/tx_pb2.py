# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/zentp/tx.proto
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
    'zrchain/zentp/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from zrchain.zentp import params_pb2 as zrchain_dot_zentp_dot_params__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x16zrchain/zentp/tx.proto\x12\rzrchain.zentp\x1a\x11\x61mino/amino.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x1azrchain/zentp/params.proto\"\xb7\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x38\n\x06params\x18\x02 \x01(\x0b\x32\x15.zrchain.zentp.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:2\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*\x1fzrchain/x/zentp/MsgUpdateParams\"\x19\n\x17MsgUpdateParamsResponse\"\xc8\x01\n\x0bMsgMintRock\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x16\n\x06\x61mount\x18\x02 \x01(\x04R\x06\x61mount\x12\"\n\rsource_key_id\x18\x03 \x01(\x04R\x0bsourceKeyId\x12+\n\x11\x64\x65stination_chain\x18\x04 \x01(\tR\x10\x64\x65stinationChain\x12(\n\x10recipient_key_id\x18\x05 \x01(\x04R\x0erecipientKeyId:\x0c\x82\xe7\xb0*\x07\x63reator\".\n\x13MsgMintRockResponse\x12\x17\n\x07tx_hash\x18\x01 \x01(\tR\x06txHash\"\xc6\x01\n\x07MsgBurn\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12%\n\x0emodule_account\x18\x02 \x01(\tR\rmoduleAccount\x12\x14\n\x05\x64\x65nom\x18\x03 \x01(\tR\x05\x64\x65nom\x12\x16\n\x06\x61mount\x18\x04 \x01(\x03R\x06\x61mount:.\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*\x1bzrchain/x/zentp/MsgMintRock\"\x11\n\x0fMsgBurnResponse2\xf0\x01\n\x03Msg\x12V\n\x0cUpdateParams\x12\x1e.zrchain.zentp.MsgUpdateParams\x1a&.zrchain.zentp.MsgUpdateParamsResponse\x12J\n\x08MintRock\x12\x1a.zrchain.zentp.MsgMintRock\x1a\".zrchain.zentp.MsgMintRockResponse\x12>\n\x04\x42urn\x12\x16.zrchain.zentp.MsgBurn\x1a\x1e.zrchain.zentp.MsgBurnResponse\x1a\x05\x80\xe7\xb0*\x01\x42\x38Z6github.com/Zenrock-Foundation/zrchain/v5/x/zentp/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.zentp.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z6github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\202\347\260*\tauthority\212\347\260*\037zrchain/x/zentp/MsgUpdateParams'
  _globals['_MSGMINTROCK']._loaded_options = None
  _globals['_MSGMINTROCK']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGBURN'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGBURN'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGBURN']._loaded_options = None
  _globals['_MSGBURN']._serialized_options = b'\202\347\260*\tauthority\212\347\260*\033zrchain/x/zentp/MsgMintRock'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._serialized_start=163
  _globals['_MSGUPDATEPARAMS']._serialized_end=346
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=348
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=373
  _globals['_MSGMINTROCK']._serialized_start=376
  _globals['_MSGMINTROCK']._serialized_end=576
  _globals['_MSGMINTROCKRESPONSE']._serialized_start=578
  _globals['_MSGMINTROCKRESPONSE']._serialized_end=624
  _globals['_MSGBURN']._serialized_start=627
  _globals['_MSGBURN']._serialized_end=825
  _globals['_MSGBURNRESPONSE']._serialized_start=827
  _globals['_MSGBURNRESPONSE']._serialized_end=844
  _globals['_MSG']._serialized_start=847
  _globals['_MSG']._serialized_end=1087
# @@protoc_insertion_point(module_scope)
