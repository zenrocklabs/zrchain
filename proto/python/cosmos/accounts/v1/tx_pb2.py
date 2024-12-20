# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/accounts/v1/tx.proto
# Protobuf Python Version: 5.29.1
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
    1,
    '',
    'cosmos/accounts/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos.tx.v1beta1 import tx_pb2 as cosmos_dot_tx_dot_v1beta1_dot_tx__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1b\x63osmos/accounts/v1/tx.proto\x12\x12\x63osmos.accounts.v1\x1a\x19google/protobuf/any.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x1a\x63osmos/tx/v1beta1/tx.proto\x1a\x14gogoproto/gogo.proto\"\xe4\x01\n\x07MsgInit\x12\x16\n\x06sender\x18\x01 \x01(\tR\x06sender\x12!\n\x0c\x61\x63\x63ount_type\x18\x02 \x01(\tR\x0b\x61\x63\x63ountType\x12.\n\x07message\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyR\x07message\x12\x61\n\x05\x66unds\x18\x04 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x05\x66unds:\x0b\x82\xe7\xb0*\x06sender\"l\n\x0fMsgInitResponse\x12\'\n\x0f\x61\x63\x63ount_address\x18\x01 \x01(\tR\x0e\x61\x63\x63ountAddress\x12\x30\n\x08response\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x08response\"\xdc\x01\n\nMsgExecute\x12\x16\n\x06sender\x18\x01 \x01(\tR\x06sender\x12\x16\n\x06target\x18\x02 \x01(\tR\x06target\x12.\n\x07message\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyR\x07message\x12\x61\n\x05\x66unds\x18\x04 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\x05\x66unds:\x0b\x82\xe7\xb0*\x06sender\"F\n\x12MsgExecuteResponse\x12\x30\n\x08response\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\x08response\"f\n\x10MsgExecuteBundle\x12\x18\n\x07\x62undler\x18\x01 \x01(\tR\x07\x62undler\x12*\n\x03txs\x18\x02 \x03(\x0b\x32\x18.cosmos.tx.v1beta1.TxRawR\x03txs:\x0c\x82\xe7\xb0*\x07\x62undler\"f\n\x11\x42undledTxResponse\x12;\n\x0e\x65xec_responses\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\rexecResponses\x12\x14\n\x05\x65rror\x18\x02 \x01(\tR\x05\x65rror\"_\n\x18MsgExecuteBundleResponse\x12\x43\n\tresponses\x18\x01 \x03(\x0b\x32%.cosmos.accounts.v1.BundledTxResponseR\tresponses2\x8e\x02\n\x03Msg\x12H\n\x04Init\x12\x1b.cosmos.accounts.v1.MsgInit\x1a#.cosmos.accounts.v1.MsgInitResponse\x12Q\n\x07\x45xecute\x12\x1e.cosmos.accounts.v1.MsgExecute\x1a&.cosmos.accounts.v1.MsgExecuteResponse\x12\x63\n\rExecuteBundle\x12$.cosmos.accounts.v1.MsgExecuteBundle\x1a,.cosmos.accounts.v1.MsgExecuteBundleResponse\x1a\x05\x80\xe7\xb0*\x01\x42\x1cZ\x1a\x63osmossdk.io/x/accounts/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.accounts.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\032cosmossdk.io/x/accounts/v1'
  _globals['_MSGINIT'].fields_by_name['funds']._loaded_options = None
  _globals['_MSGINIT'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGINIT']._loaded_options = None
  _globals['_MSGINIT']._serialized_options = b'\202\347\260*\006sender'
  _globals['_MSGEXECUTE'].fields_by_name['funds']._loaded_options = None
  _globals['_MSGEXECUTE'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_MSGEXECUTE']._loaded_options = None
  _globals['_MSGEXECUTE']._serialized_options = b'\202\347\260*\006sender'
  _globals['_MSGEXECUTEBUNDLE']._loaded_options = None
  _globals['_MSGEXECUTEBUNDLE']._serialized_options = b'\202\347\260*\007bundler'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGINIT']._serialized_start=186
  _globals['_MSGINIT']._serialized_end=414
  _globals['_MSGINITRESPONSE']._serialized_start=416
  _globals['_MSGINITRESPONSE']._serialized_end=524
  _globals['_MSGEXECUTE']._serialized_start=527
  _globals['_MSGEXECUTE']._serialized_end=747
  _globals['_MSGEXECUTERESPONSE']._serialized_start=749
  _globals['_MSGEXECUTERESPONSE']._serialized_end=819
  _globals['_MSGEXECUTEBUNDLE']._serialized_start=821
  _globals['_MSGEXECUTEBUNDLE']._serialized_end=923
  _globals['_BUNDLEDTXRESPONSE']._serialized_start=925
  _globals['_BUNDLEDTXRESPONSE']._serialized_end=1027
  _globals['_MSGEXECUTEBUNDLERESPONSE']._serialized_start=1029
  _globals['_MSGEXECUTEBUNDLERESPONSE']._serialized_end=1124
  _globals['_MSG']._serialized_start=1127
  _globals['_MSG']._serialized_end=1397
# @@protoc_insertion_point(module_scope)
