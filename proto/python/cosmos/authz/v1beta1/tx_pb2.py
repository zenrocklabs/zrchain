# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/authz/v1beta1/tx.proto
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
    'cosmos/authz/v1beta1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from cosmos.authz.v1beta1 import authz_pb2 as cosmos_dot_authz_dot_v1beta1_dot_authz__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1d\x63osmos/authz/v1beta1/tx.proto\x12\x14\x63osmos.authz.v1beta1\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a cosmos/authz/v1beta1/authz.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x11\x61mino/amino.proto\"\xd6\x01\n\x08MsgGrant\x12\x32\n\x07granter\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07granter\x12\x32\n\x07grantee\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07grantee\x12<\n\x05grant\x18\x03 \x01(\x0b\x32\x1b.cosmos.authz.v1beta1.GrantB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x05grant:$\x82\xe7\xb0*\x07granter\x8a\xe7\xb0*\x13\x63osmos-sdk/MsgGrant\"\x12\n\x10MsgGrantResponse\"\xa9\x01\n\x07MsgExec\x12\x32\n\x07grantee\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07grantee\x12\x45\n\x04msgs\x18\x02 \x03(\x0b\x32\x14.google.protobuf.AnyB\x1b\xca\xb4-\x17\x63osmos.base.v1beta1.MsgR\x04msgs:#\x82\xe7\xb0*\x07grantee\x8a\xe7\xb0*\x12\x63osmos-sdk/MsgExec\"+\n\x0fMsgExecResponse\x12\x18\n\x07results\x18\x01 \x03(\x0cR\x07results\"\xbc\x01\n\tMsgRevoke\x12\x32\n\x07granter\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07granter\x12\x32\n\x07grantee\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07grantee\x12 \n\x0cmsg_type_url\x18\x03 \x01(\tR\nmsgTypeUrl:%\x82\xe7\xb0*\x07granter\x8a\xe7\xb0*\x14\x63osmos-sdk/MsgRevoke\"\x13\n\x11MsgRevokeResponse2\xff\x01\n\x03Msg\x12O\n\x05Grant\x12\x1e.cosmos.authz.v1beta1.MsgGrant\x1a&.cosmos.authz.v1beta1.MsgGrantResponse\x12L\n\x04\x45xec\x12\x1d.cosmos.authz.v1beta1.MsgExec\x1a%.cosmos.authz.v1beta1.MsgExecResponse\x12R\n\x06Revoke\x12\x1f.cosmos.authz.v1beta1.MsgRevoke\x1a\'.cosmos.authz.v1beta1.MsgRevokeResponse\x1a\x05\x80\xe7\xb0*\x01\x42*Z$github.com/cosmos/cosmos-sdk/x/authz\xc8\xe1\x1e\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.authz.v1beta1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z$github.com/cosmos/cosmos-sdk/x/authz\310\341\036\000'
  _globals['_MSGGRANT'].fields_by_name['granter']._loaded_options = None
  _globals['_MSGGRANT'].fields_by_name['granter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGGRANT'].fields_by_name['grantee']._loaded_options = None
  _globals['_MSGGRANT'].fields_by_name['grantee']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGGRANT'].fields_by_name['grant']._loaded_options = None
  _globals['_MSGGRANT'].fields_by_name['grant']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGGRANT']._loaded_options = None
  _globals['_MSGGRANT']._serialized_options = b'\202\347\260*\007granter\212\347\260*\023cosmos-sdk/MsgGrant'
  _globals['_MSGEXEC'].fields_by_name['grantee']._loaded_options = None
  _globals['_MSGEXEC'].fields_by_name['grantee']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGEXEC'].fields_by_name['msgs']._loaded_options = None
  _globals['_MSGEXEC'].fields_by_name['msgs']._serialized_options = b'\312\264-\027cosmos.base.v1beta1.Msg'
  _globals['_MSGEXEC']._loaded_options = None
  _globals['_MSGEXEC']._serialized_options = b'\202\347\260*\007grantee\212\347\260*\022cosmos-sdk/MsgExec'
  _globals['_MSGREVOKE'].fields_by_name['granter']._loaded_options = None
  _globals['_MSGREVOKE'].fields_by_name['granter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGREVOKE'].fields_by_name['grantee']._loaded_options = None
  _globals['_MSGREVOKE'].fields_by_name['grantee']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGREVOKE']._loaded_options = None
  _globals['_MSGREVOKE']._serialized_options = b'\202\347\260*\007granter\212\347\260*\024cosmos-sdk/MsgRevoke'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGGRANT']._serialized_start=210
  _globals['_MSGGRANT']._serialized_end=424
  _globals['_MSGGRANTRESPONSE']._serialized_start=426
  _globals['_MSGGRANTRESPONSE']._serialized_end=444
  _globals['_MSGEXEC']._serialized_start=447
  _globals['_MSGEXEC']._serialized_end=616
  _globals['_MSGEXECRESPONSE']._serialized_start=618
  _globals['_MSGEXECRESPONSE']._serialized_end=661
  _globals['_MSGREVOKE']._serialized_start=664
  _globals['_MSGREVOKE']._serialized_end=852
  _globals['_MSGREVOKERESPONSE']._serialized_start=854
  _globals['_MSGREVOKERESPONSE']._serialized_end=873
  _globals['_MSG']._serialized_start=876
  _globals['_MSG']._serialized_end=1131
# @@protoc_insertion_point(module_scope)
