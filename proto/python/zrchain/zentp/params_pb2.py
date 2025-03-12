# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/zentp/params.proto
# Protobuf Python Version: 6.30.0
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
    0,
    '',
    'zrchain/zentp/params.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1azrchain/zentp/params.proto\x12\rzrchain.zentp\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\"\x8d\x01\n\x06Params\x12\x33\n\x16zrchain_relayer_key_id\x18\x01 \x01(\x04R\x13zrchainRelayerKeyId\x12-\n\x06solana\x18\x02 \x01(\x0b\x32\x15.zrchain.zentp.SolanaR\x06solana:\x1f\xe8\xa0\x1f\x01\x8a\xe7\xb0*\x16zrchain/x/zentp/Params\"\x9a\x02\n\x06Solana\x12\"\n\rsigner_key_id\x18\x01 \x01(\x04R\x0bsignerKeyId\x12\x1d\n\nprogram_id\x18\x02 \x01(\tR\tprogramId\x12*\n\x11nonce_account_key\x18\x03 \x01(\x04R\x0fnonceAccountKey\x12.\n\x13nonce_authority_key\x18\x04 \x01(\x04R\x11nonceAuthorityKey\x12\x17\n\x07rpc_url\x18\x05 \x01(\tR\x06rpcUrl\x12!\n\x0cmint_address\x18\x06 \x01(\tR\x0bmintAddress\x12\x1d\n\nfee_wallet\x18\x07 \x01(\tR\tfeeWallet\x12\x10\n\x03\x66\x65\x65\x18\x08 \x01(\x04R\x03\x66\x65\x65:\x04\xe8\xa0\x1f\x01\x42\x38Z6github.com/Zenrock-Foundation/zrchain/v5/x/zentp/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.zentp.params_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z6github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\350\240\037\001\212\347\260*\026zrchain/x/zentp/Params'
  _globals['_SOLANA']._loaded_options = None
  _globals['_SOLANA']._serialized_options = b'\350\240\037\001'
  _globals['_PARAMS']._serialized_start=87
  _globals['_PARAMS']._serialized_end=228
  _globals['_SOLANA']._serialized_start=231
  _globals['_SOLANA']._serialized_end=513
# @@protoc_insertion_point(module_scope)
