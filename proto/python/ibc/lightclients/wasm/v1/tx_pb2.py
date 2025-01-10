# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/lightclients/wasm/v1/tx.proto
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
    'ibc/lightclients/wasm/v1/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n!ibc/lightclients/wasm/v1/tx.proto\x12\x18ibc.lightclients.wasm.v1\x1a\x17\x63osmos/msg/v1/msg.proto\"Y\n\x0cMsgStoreCode\x12\x16\n\x06signer\x18\x01 \x01(\tR\x06signer\x12$\n\x0ewasm_byte_code\x18\x02 \x01(\x0cR\x0cwasmByteCode:\x0b\x82\xe7\xb0*\x06signer\"2\n\x14MsgStoreCodeResponse\x12\x1a\n\x08\x63hecksum\x18\x01 \x01(\x0cR\x08\x63hecksum\"T\n\x11MsgRemoveChecksum\x12\x16\n\x06signer\x18\x01 \x01(\tR\x06signer\x12\x1a\n\x08\x63hecksum\x18\x02 \x01(\x0cR\x08\x63hecksum:\x0b\x82\xe7\xb0*\x06signer\"\x1b\n\x19MsgRemoveChecksumResponse\"\x84\x01\n\x12MsgMigrateContract\x12\x16\n\x06signer\x18\x01 \x01(\tR\x06signer\x12\x1b\n\tclient_id\x18\x02 \x01(\tR\x08\x63lientId\x12\x1a\n\x08\x63hecksum\x18\x03 \x01(\x0cR\x08\x63hecksum\x12\x10\n\x03msg\x18\x04 \x01(\x0cR\x03msg:\x0b\x82\xe7\xb0*\x06signer\"\x1c\n\x1aMsgMigrateContractResponse2\xdc\x02\n\x03Msg\x12\x63\n\tStoreCode\x12&.ibc.lightclients.wasm.v1.MsgStoreCode\x1a..ibc.lightclients.wasm.v1.MsgStoreCodeResponse\x12r\n\x0eRemoveChecksum\x12+.ibc.lightclients.wasm.v1.MsgRemoveChecksum\x1a\x33.ibc.lightclients.wasm.v1.MsgRemoveChecksumResponse\x12u\n\x0fMigrateContract\x12,.ibc.lightclients.wasm.v1.MsgMigrateContract\x1a\x34.ibc.lightclients.wasm.v1.MsgMigrateContractResponse\x1a\x05\x80\xe7\xb0*\x01\x42>Z<github.com/cosmos/ibc-go/modules/light-clients/08-wasm/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.lightclients.wasm.v1.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z<github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types'
  _globals['_MSGSTORECODE']._loaded_options = None
  _globals['_MSGSTORECODE']._serialized_options = b'\202\347\260*\006signer'
  _globals['_MSGREMOVECHECKSUM']._loaded_options = None
  _globals['_MSGREMOVECHECKSUM']._serialized_options = b'\202\347\260*\006signer'
  _globals['_MSGMIGRATECONTRACT']._loaded_options = None
  _globals['_MSGMIGRATECONTRACT']._serialized_options = b'\202\347\260*\006signer'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_MSGSTORECODE']._serialized_start=88
  _globals['_MSGSTORECODE']._serialized_end=177
  _globals['_MSGSTORECODERESPONSE']._serialized_start=179
  _globals['_MSGSTORECODERESPONSE']._serialized_end=229
  _globals['_MSGREMOVECHECKSUM']._serialized_start=231
  _globals['_MSGREMOVECHECKSUM']._serialized_end=315
  _globals['_MSGREMOVECHECKSUMRESPONSE']._serialized_start=317
  _globals['_MSGREMOVECHECKSUMRESPONSE']._serialized_end=344
  _globals['_MSGMIGRATECONTRACT']._serialized_start=347
  _globals['_MSGMIGRATECONTRACT']._serialized_end=479
  _globals['_MSGMIGRATECONTRACTRESPONSE']._serialized_start=481
  _globals['_MSGMIGRATECONTRACTRESPONSE']._serialized_end=509
  _globals['_MSG']._serialized_start=512
  _globals['_MSG']._serialized_end=860
# @@protoc_insertion_point(module_scope)
