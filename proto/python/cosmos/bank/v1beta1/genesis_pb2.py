# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/bank/v1beta1/genesis.proto
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
    'cosmos/bank/v1beta1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos.bank.v1beta1 import bank_pb2 as cosmos_dot_bank_dot_v1beta1_dot_bank__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n!cosmos/bank/v1beta1/genesis.proto\x12\x13\x63osmos.bank.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x1e\x63osmos/bank/v1beta1/bank.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x11\x61mino/amino.proto\"\xaf\x03\n\x0cGenesisState\x12>\n\x06params\x18\x01 \x01(\x0b\x32\x1b.cosmos.bank.v1beta1.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\x12\x43\n\x08\x62\x61lances\x18\x02 \x03(\x0b\x32\x1c.cosmos.bank.v1beta1.BalanceB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x08\x62\x61lances\x12y\n\x06supply\x18\x03 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x06supply\x12O\n\x0e\x64\x65nom_metadata\x18\x04 \x03(\x0b\x32\x1d.cosmos.bank.v1beta1.MetadataB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\rdenomMetadata\x12N\n\x0csend_enabled\x18\x05 \x03(\x0b\x32 .cosmos.bank.v1beta1.SendEnabledB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0bsendEnabled\"\xc0\x01\n\x07\x42\x61lance\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\x12w\n\x05\x63oins\x18\x02 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x05\x63oins:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\x42\x1bZ\x19\x63osmossdk.io/x/bank/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.bank.v1beta1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\031cosmossdk.io/x/bank/types'
  _globals['_GENESISSTATE'].fields_by_name['params']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_GENESISSTATE'].fields_by_name['balances']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['balances']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_GENESISSTATE'].fields_by_name['supply']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['supply']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_GENESISSTATE'].fields_by_name['denom_metadata']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['denom_metadata']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_GENESISSTATE'].fields_by_name['send_enabled']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['send_enabled']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_BALANCE'].fields_by_name['address']._loaded_options = None
  _globals['_BALANCE'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_BALANCE'].fields_by_name['coins']._loaded_options = None
  _globals['_BALANCE'].fields_by_name['coins']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_BALANCE']._loaded_options = None
  _globals['_BALANCE']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_GENESISSTATE']._serialized_start=191
  _globals['_GENESISSTATE']._serialized_end=622
  _globals['_BALANCE']._serialized_start=625
  _globals['_BALANCE']._serialized_end=817
# @@protoc_insertion_point(module_scope)
