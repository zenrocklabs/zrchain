# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/treasury/genesis.proto
# Protobuf Python Version: 5.29.0
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
    0,
    '',
    'zrchain/treasury/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from zrchain.treasury import params_pb2 as zrchain_dot_treasury_dot_params__pb2
from zrchain.treasury import key_pb2 as zrchain_dot_treasury_dot_key__pb2
from zrchain.treasury import mpcsign_pb2 as zrchain_dot_treasury_dot_mpcsign__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1ezrchain/treasury/genesis.proto\x12\x10zrchain.treasury\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x1dzrchain/treasury/params.proto\x1a\x1azrchain/treasury/key.proto\x1a\x1ezrchain/treasury/mpcsign.proto\"\xd7\x03\n\x0cGenesisState\x12;\n\x06params\x18\x01 \x01(\x0b\x32\x18.zrchain.treasury.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\x12\x17\n\x07port_id\x18\x02 \x01(\tR\x06portId\x12/\n\x04keys\x18\x03 \x03(\x0b\x32\x15.zrchain.treasury.KeyB\x04\xc8\xde\x1f\x00R\x04keys\x12\x45\n\x0ckey_requests\x18\x04 \x03(\x0b\x32\x1c.zrchain.treasury.KeyRequestB\x04\xc8\xde\x1f\x00R\x0bkeyRequests\x12H\n\rsign_requests\x18\x05 \x03(\x0b\x32\x1d.zrchain.treasury.SignRequestB\x04\xc8\xde\x1f\x00R\x0csignRequests\x12X\n\x10sign_tx_requests\x18\x06 \x03(\x0b\x32(.zrchain.treasury.SignTransactionRequestB\x04\xc8\xde\x1f\x00R\x0esignTxRequests\x12U\n\x0fica_tx_requests\x18\x07 \x03(\x0b\x32\'.zrchain.treasury.ICATransactionRequestB\x04\xc8\xde\x1f\x00R\ricaTxRequestsB;Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.treasury.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types'
  _globals['_GENESISSTATE'].fields_by_name['params']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_GENESISSTATE'].fields_by_name['keys']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['keys']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['key_requests']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['key_requests']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['sign_requests']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['sign_requests']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['sign_tx_requests']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['sign_tx_requests']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['ica_tx_requests']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['ica_tx_requests']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE']._serialized_start=185
  _globals['_GENESISSTATE']._serialized_end=656
# @@protoc_insertion_point(module_scope)
