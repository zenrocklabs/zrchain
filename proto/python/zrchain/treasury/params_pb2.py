# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/treasury/params.proto
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
    'zrchain/treasury/params.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1dzrchain/treasury/params.proto\x12\x10zrchain.treasury\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\"\xab\x02\n\x06Params\x12\x1f\n\x0bmpc_keyring\x18\x01 \x01(\tR\nmpcKeyring\x12&\n\x0fzr_sign_address\x18\x02 \x01(\tR\rzrSignAddress\x12-\n\x12keyring_commission\x18\x03 \x01(\x04R\x11keyringCommission\x12\x44\n\x1ekeyring_commission_destination\x18\x04 \x01(\tR\x1ckeyringCommissionDestination\x12\x1e\n\x0bmin_gas_fee\x18\x05 \x01(\tR\tminGasFee:C\xe8\xa0\x1f\x01\x8a\xe7\xb0*:github.com/Zenrock-Foundation/zrchain/v5/x/treasury/ParamsB;Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.treasury.params_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\350\240\037\001\212\347\260*:github.com/Zenrock-Foundation/zrchain/v5/x/treasury/Params'
  _globals['_PARAMS']._serialized_start=93
  _globals['_PARAMS']._serialized_end=392
# @@protoc_insertion_point(module_scope)
