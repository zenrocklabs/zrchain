# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/crypto/keyring/v1/record.proto
<<<<<<< HEAD
# Protobuf Python Version: 6.30.1
=======
# Protobuf Python Version: 6.30.0
>>>>>>> main
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    30,
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'cosmos/crypto/keyring/v1/record.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from cosmos.crypto.hd.v1 import hd_pb2 as cosmos_dot_crypto_dot_hd_dot_v1_dot_hd__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n%cosmos/crypto/keyring/v1/record.proto\x12\x18\x63osmos.crypto.keyring.v1\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a\x1c\x63osmos/crypto/hd/v1/hd.proto\"\xea\x03\n\x06Record\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12-\n\x07pub_key\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x06pubKey\x12>\n\x05local\x18\x03 \x01(\x0b\x32&.cosmos.crypto.keyring.v1.Record.LocalH\x00R\x05local\x12\x41\n\x06ledger\x18\x04 \x01(\x0b\x32\'.cosmos.crypto.keyring.v1.Record.LedgerH\x00R\x06ledger\x12>\n\x05multi\x18\x05 \x01(\x0b\x32&.cosmos.crypto.keyring.v1.Record.MultiH\x00R\x05multi\x12\x44\n\x07offline\x18\x06 \x01(\x0b\x32(.cosmos.crypto.keyring.v1.Record.OfflineH\x00R\x07offline\x1a\x38\n\x05Local\x12/\n\x08priv_key\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\x07privKey\x1a>\n\x06Ledger\x12\x34\n\x04path\x18\x01 \x01(\x0b\x32 .cosmos.crypto.hd.v1.BIP44ParamsR\x04path\x1a\x07\n\x05Multi\x1a\t\n\x07OfflineB\x06\n\x04itemB5Z+github.com/cosmos/cosmos-sdk/crypto/keyring\xc8\xe1\x1e\x00\x98\xe3\x1e\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.crypto.keyring.v1.record_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+github.com/cosmos/cosmos-sdk/crypto/keyring\310\341\036\000\230\343\036\000'
  _globals['_RECORD']._serialized_start=147
  _globals['_RECORD']._serialized_end=637
  _globals['_RECORD_LOCAL']._serialized_start=489
  _globals['_RECORD_LOCAL']._serialized_end=545
  _globals['_RECORD_LEDGER']._serialized_start=547
  _globals['_RECORD_LEDGER']._serialized_end=609
  _globals['_RECORD_MULTI']._serialized_start=611
  _globals['_RECORD_MULTI']._serialized_end=618
  _globals['_RECORD_OFFLINE']._serialized_start=620
  _globals['_RECORD_OFFLINE']._serialized_end=629
# @@protoc_insertion_point(module_scope)
