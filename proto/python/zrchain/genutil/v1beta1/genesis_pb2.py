# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/genutil/v1beta1/genesis.proto
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
    'zrchain/genutil/v1beta1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n%zrchain/genutil/v1beta1/genesis.proto\x12\x17zrchain.genutil.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x11\x61mino/amino.proto\"_\n\x0cGenesisState\x12O\n\x07gen_txs\x18\x01 \x03(\x0c\x42\x36\xea\xde\x1f\x06gentxs\xfa\xde\x1f\x18\x65ncoding/json.RawMessage\xa2\xe7\xb0*\x06gentxs\xa8\xe7\xb0*\x01R\x06genTxsB:Z8github.com/Zenrock-Foundation/zrchain/v6/x/genutil/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.genutil.v1beta1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z8github.com/Zenrock-Foundation/zrchain/v6/x/genutil/types'
  _globals['_GENESISSTATE'].fields_by_name['gen_txs']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['gen_txs']._serialized_options = b'\352\336\037\006gentxs\372\336\037\030encoding/json.RawMessage\242\347\260*\006gentxs\250\347\260*\001'
  _globals['_GENESISSTATE']._serialized_start=107
  _globals['_GENESISSTATE']._serialized_end=202
# @@protoc_insertion_point(module_scope)
