# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/slashing/v1beta1/tx.proto
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
    'cosmos/slashing/v1beta1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.slashing.v1beta1 import slashing_pb2 as cosmos_dot_slashing_dot_v1beta1_dot_slashing__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n cosmos/slashing/v1beta1/tx.proto\x12\x17\x63osmos.slashing.v1beta1\x1a\x14gogoproto/gogo.proto\x1a&cosmos/slashing/v1beta1/slashing.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x11\x61mino/amino.proto\"\xa3\x01\n\tMsgUnjail\x12\x64\n\x0evalidator_addr\x18\x01 \x01(\tB=\xea\xde\x1f\x07\x61\x64\x64ress\xd2\xb4-\x1d\x63osmos.ValidatorAddressString\xa2\xe7\xb0*\x07\x61\x64\x64ress\xa8\xe7\xb0*\x01R\rvalidatorAddr:0\x88\xa0\x1f\x00\x82\xe7\xb0*\x0evalidator_addr\x8a\xe7\xb0*\x14\x63osmos-sdk/MsgUnjail\"\x13\n\x11MsgUnjailResponse\"\xc7\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12\x42\n\x06params\x18\x02 \x01(\x0b\x32\x1f.cosmos.slashing.v1beta1.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:8\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*%cosmos-sdk/x/slashing/MsgUpdateParams\"\x19\n\x17MsgUpdateParamsResponse2\xd2\x01\n\x03Msg\x12X\n\x06Unjail\x12\".cosmos.slashing.v1beta1.MsgUnjail\x1a*.cosmos.slashing.v1beta1.MsgUnjailResponse\x12j\n\x0cUpdateParams\x12(.cosmos.slashing.v1beta1.MsgUpdateParams\x1a\x30.cosmos.slashing.v1beta1.MsgUpdateParamsResponse\x1a\x05\x80\xe7\xb0*\x01\x42#Z\x1d\x63osmossdk.io/x/slashing/types\xa8\xe2\x1e\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.slashing.v1beta1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\035cosmossdk.io/x/slashing/types\250\342\036\001'
  _globals['_MSGUNJAIL'].fields_by_name['validator_addr']._loaded_options = None
  _globals['_MSGUNJAIL'].fields_by_name['validator_addr']._serialized_options = b'\352\336\037\007address\322\264-\035cosmos.ValidatorAddressString\242\347\260*\007address\250\347\260*\001'
  _globals['_MSGUNJAIL']._loaded_options = None
  _globals['_MSGUNJAIL']._serialized_options = b'\210\240\037\000\202\347\260*\016validator_addr\212\347\260*\024cosmos-sdk/MsgUnjail'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\202\347\260*\tauthority\212\347\260*%cosmos-sdk/x/slashing/MsgUpdateParams'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGUNJAIL']._serialized_start=195
  _globals['_MSGUNJAIL']._serialized_end=358
  _globals['_MSGUNJAILRESPONSE']._serialized_start=360
  _globals['_MSGUNJAILRESPONSE']._serialized_end=379
  _globals['_MSGUPDATEPARAMS']._serialized_start=382
  _globals['_MSGUPDATEPARAMS']._serialized_end=581
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=583
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=608
  _globals['_MSG']._serialized_start=611
  _globals['_MSG']._serialized_end=821
# @@protoc_insertion_point(module_scope)
