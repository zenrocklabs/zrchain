# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/rpc/grpc/types.proto
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
    'tendermint/rpc/grpc/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from tendermint.abci import types_pb2 as tendermint_dot_abci_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1ftendermint/rpc/grpc/types.proto\x12\x13tendermint.rpc.grpc\x1a\x1btendermint/abci/types.proto\"\r\n\x0bRequestPing\"$\n\x12RequestBroadcastTx\x12\x0e\n\x02tx\x18\x01 \x01(\x0cR\x02tx\"\x0e\n\x0cResponsePing\"\x8e\x01\n\x13ResponseBroadcastTx\x12;\n\x08\x63heck_tx\x18\x01 \x01(\x0b\x32 .tendermint.abci.ResponseCheckTxR\x07\x63heckTx\x12:\n\ttx_result\x18\x02 \x01(\x0b\x32\x1d.tendermint.abci.ExecTxResultR\x08txResult2\xbd\x01\n\x0c\x42roadcastAPI\x12K\n\x04Ping\x12 .tendermint.rpc.grpc.RequestPing\x1a!.tendermint.rpc.grpc.ResponsePing\x12`\n\x0b\x42roadcastTx\x12\'.tendermint.rpc.grpc.RequestBroadcastTx\x1a(.tendermint.rpc.grpc.ResponseBroadcastTxB0Z.github.com/cometbft/cometbft/rpc/grpc;coregrpcb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.rpc.grpc.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z.github.com/cometbft/cometbft/rpc/grpc;coregrpc'
  _globals['_REQUESTPING']._serialized_start=85
  _globals['_REQUESTPING']._serialized_end=98
  _globals['_REQUESTBROADCASTTX']._serialized_start=100
  _globals['_REQUESTBROADCASTTX']._serialized_end=136
  _globals['_RESPONSEPING']._serialized_start=138
  _globals['_RESPONSEPING']._serialized_end=152
  _globals['_RESPONSEBROADCASTTX']._serialized_start=155
  _globals['_RESPONSEBROADCASTTX']._serialized_end=297
  _globals['_BROADCASTAPI']._serialized_start=300
  _globals['_BROADCASTAPI']._serialized_end=489
# @@protoc_insertion_point(module_scope)
