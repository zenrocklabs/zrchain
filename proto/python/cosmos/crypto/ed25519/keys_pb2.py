# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/crypto/ed25519/keys.proto
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
    'cosmos/crypto/ed25519/keys.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n cosmos/crypto/ed25519/keys.proto\x12\x15\x63osmos.crypto.ed25519\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\"i\n\x06PubKey\x12.\n\x03key\x18\x01 \x01(\x0c\x42\x1c\xfa\xde\x1f\x18\x63rypto/ed25519.PublicKeyR\x03key:/\x98\xa0\x1f\x00\x8a\xe7\xb0*\x18tendermint/PubKeyEd25519\x92\xe7\xb0*\tkey_field\"h\n\x07PrivKey\x12/\n\x03key\x18\x01 \x01(\x0c\x42\x1d\xfa\xde\x1f\x19\x63rypto/ed25519.PrivateKeyR\x03key:,\x8a\xe7\xb0*\x19tendermint/PrivKeyEd25519\x92\xe7\xb0*\tkey_fieldB2Z0github.com/cosmos/cosmos-sdk/crypto/keys/ed25519b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.crypto.ed25519.keys_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z0github.com/cosmos/cosmos-sdk/crypto/keys/ed25519'
  _globals['_PUBKEY'].fields_by_name['key']._loaded_options = None
  _globals['_PUBKEY'].fields_by_name['key']._serialized_options = b'\372\336\037\030crypto/ed25519.PublicKey'
  _globals['_PUBKEY']._loaded_options = None
  _globals['_PUBKEY']._serialized_options = b'\230\240\037\000\212\347\260*\030tendermint/PubKeyEd25519\222\347\260*\tkey_field'
  _globals['_PRIVKEY'].fields_by_name['key']._loaded_options = None
  _globals['_PRIVKEY'].fields_by_name['key']._serialized_options = b'\372\336\037\031crypto/ed25519.PrivateKey'
  _globals['_PRIVKEY']._loaded_options = None
  _globals['_PRIVKEY']._serialized_options = b'\212\347\260*\031tendermint/PrivKeyEd25519\222\347\260*\tkey_field'
  _globals['_PUBKEY']._serialized_start=100
  _globals['_PUBKEY']._serialized_end=205
  _globals['_PRIVKEY']._serialized_start=207
  _globals['_PRIVKEY']._serialized_end=311
# @@protoc_insertion_point(module_scope)
