# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/mint/v1beta1/tx.proto
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
    'zrchain/mint/v1beta1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from amino import amino_pb2 as amino_dot_amino__pb2
from zrchain.mint.v1beta1 import mint_pb2 as zrchain_dot_mint_dot_v1beta1_dot_mint__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1dzrchain/mint/v1beta1/tx.proto\x12\x14zrchain.mint.v1beta1\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x11\x61mino/amino.proto\x1a\x1fzrchain/mint/v1beta1/mint.proto\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\xbd\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12?\n\x06params\x18\x02 \x01(\x0b\x32\x1c.zrchain.mint.v1beta1.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:1\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*\x1ezrchain/x/mint/MsgUpdateParams\"\x19\n\x17MsgUpdateParamsResponse2r\n\x03Msg\x12\x64\n\x0cUpdateParams\x12%.zrchain.mint.v1beta1.MsgUpdateParams\x1a-.zrchain.mint.v1beta1.MsgUpdateParamsResponse\x1a\x05\x80\xe7\xb0*\x01\x42\x37Z5github.com/Zenrock-Foundation/zrchain/v6/x/mint/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.mint.v1beta1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/Zenrock-Foundation/zrchain/v6/x/mint/types'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\202\347\260*\tauthority\212\347\260*\036zrchain/x/mint/MsgUpdateParams'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._serialized_start=182
  _globals['_MSGUPDATEPARAMS']._serialized_end=371
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=373
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=398
  _globals['_MSG']._serialized_start=400
  _globals['_MSG']._serialized_end=514
# @@protoc_insertion_point(module_scope)
