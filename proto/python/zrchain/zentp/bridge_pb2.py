# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/zentp/bridge.proto
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
    'zrchain/zentp/bridge.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from zrchain.zentp import params_pb2 as zrchain_dot_zentp_dot_params__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1azrchain/zentp/bridge.proto\x12\rzrchain.zentp\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x1azrchain/zentp/params.proto\"\xba\x03\n\x06\x42ridge\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x14\n\x05\x64\x65nom\x18\x02 \x01(\tR\x05\x64\x65nom\x12\x18\n\x07\x63reator\x18\x03 \x01(\tR\x07\x63reator\x12%\n\x0esource_address\x18\x04 \x01(\tR\rsourceAddress\x12!\n\x0csource_chain\x18\x05 \x01(\tR\x0bsourceChain\x12+\n\x11\x64\x65stination_chain\x18\x06 \x01(\tR\x10\x64\x65stinationChain\x12\x16\n\x06\x61mount\x18\x07 \x01(\x04R\x06\x61mount\x12+\n\x11recipient_address\x18\x08 \x01(\tR\x10recipientAddress\x12\x13\n\x05tx_id\x18\t \x01(\x04R\x04txId\x12\x17\n\x07tx_hash\x18\n \x01(\tR\x06txHash\x12\x31\n\x05state\x18\x0b \x01(\x0e\x32\x1b.zrchain.zentp.BridgeStatusR\x05state\x12!\n\x0c\x62lock_height\x18\x0c \x01(\x03R\x0b\x62lockHeight\x12\x30\n\x14\x61waiting_event_since\x18\r \x01(\x03R\x12\x61waitingEventSince*\x7f\n\x0c\x42ridgeStatus\x12\x1d\n\x19\x42RIDGE_STATUS_UNSPECIFIED\x10\x00\x12\x19\n\x15\x42RIDGE_STATUS_PENDING\x10\x01\x12\x1b\n\x17\x42RIDGE_STATUS_COMPLETED\x10\x02\x12\x18\n\x14\x42RIDGE_STATUS_FAILED\x10\x04\x42\x38Z6github.com/Zenrock-Foundation/zrchain/v6/x/zentp/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.zentp.bridge_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z6github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types'
  _globals['_BRIDGESTATUS']._serialized_start=559
  _globals['_BRIDGESTATUS']._serialized_end=686
  _globals['_BRIDGE']._serialized_start=115
  _globals['_BRIDGE']._serialized_end=557
# @@protoc_insertion_point(module_scope)
