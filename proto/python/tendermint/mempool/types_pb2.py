# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/mempool/types.proto
# Protobuf Python Version: 5.29.2
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
    2,
    '',
    'tendermint/mempool/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1etendermint/mempool/types.proto\x12\x12tendermint.mempool\"\x17\n\x03Txs\x12\x10\n\x03txs\x18\x01 \x03(\x0cR\x03txs\"=\n\x07Message\x12+\n\x03txs\x18\x01 \x01(\x0b\x32\x17.tendermint.mempool.TxsH\x00R\x03txsB\x05\n\x03sumB7Z5github.com/cometbft/cometbft/proto/tendermint/mempoolb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.mempool.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/cometbft/cometbft/proto/tendermint/mempool'
  _globals['_TXS']._serialized_start=54
  _globals['_TXS']._serialized_end=77
  _globals['_MESSAGE']._serialized_start=79
  _globals['_MESSAGE']._serialized_end=140
# @@protoc_insertion_point(module_scope)
