# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/tx/signing/v1beta1/signing.proto
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
    'cosmos/tx/signing/v1beta1/signing.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.crypto.multisig.v1beta1 import multisig_pb2 as cosmos_dot_crypto_dot_multisig_dot_v1beta1_dot_multisig__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\'cosmos/tx/signing/v1beta1/signing.proto\x12\x19\x63osmos.tx.signing.v1beta1\x1a-cosmos/crypto/multisig/v1beta1/multisig.proto\x1a\x19google/protobuf/any.proto\"f\n\x14SignatureDescriptors\x12N\n\nsignatures\x18\x01 \x03(\x0b\x32..cosmos.tx.signing.v1beta1.SignatureDescriptorR\nsignatures\"\xf5\x04\n\x13SignatureDescriptor\x12\x33\n\npublic_key\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\tpublicKey\x12G\n\x04\x64\x61ta\x18\x02 \x01(\x0b\x32\x33.cosmos.tx.signing.v1beta1.SignatureDescriptor.DataR\x04\x64\x61ta\x12\x1a\n\x08sequence\x18\x03 \x01(\x04R\x08sequence\x1a\xc3\x03\n\x04\x44\x61ta\x12T\n\x06single\x18\x01 \x01(\x0b\x32:.cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.SingleH\x00R\x06single\x12Q\n\x05multi\x18\x02 \x01(\x0b\x32\x39.cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.MultiH\x00R\x05multi\x1a_\n\x06Single\x12\x37\n\x04mode\x18\x01 \x01(\x0e\x32#.cosmos.tx.signing.v1beta1.SignModeR\x04mode\x12\x1c\n\tsignature\x18\x02 \x01(\x0cR\tsignature\x1a\xa9\x01\n\x05Multi\x12K\n\x08\x62itarray\x18\x01 \x01(\x0b\x32/.cosmos.crypto.multisig.v1beta1.CompactBitArrayR\x08\x62itarray\x12S\n\nsignatures\x18\x02 \x03(\x0b\x32\x33.cosmos.tx.signing.v1beta1.SignatureDescriptor.DataR\nsignaturesB\x05\n\x03sum*\xa9\x01\n\x08SignMode\x12\x19\n\x15SIGN_MODE_UNSPECIFIED\x10\x00\x12\x14\n\x10SIGN_MODE_DIRECT\x10\x01\x12\x15\n\x11SIGN_MODE_TEXTUAL\x10\x02\x12\x18\n\x14SIGN_MODE_DIRECT_AUX\x10\x03\x12\x1f\n\x1bSIGN_MODE_LEGACY_AMINO_JSON\x10\x7f\x12\x1a\n\x11SIGN_MODE_EIP_191\x10\xbf\x01\x1a\x02\x08\x01\x42/Z-github.com/cosmos/cosmos-sdk/types/tx/signingb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.tx.signing.v1beta1.signing_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z-github.com/cosmos/cosmos-sdk/types/tx/signing'
  _globals['_SIGNMODE'].values_by_name["SIGN_MODE_EIP_191"]._loaded_options = None
  _globals['_SIGNMODE'].values_by_name["SIGN_MODE_EIP_191"]._serialized_options = b'\010\001'
  _globals['_SIGNMODE']._serialized_start=881
  _globals['_SIGNMODE']._serialized_end=1050
  _globals['_SIGNATUREDESCRIPTORS']._serialized_start=144
  _globals['_SIGNATUREDESCRIPTORS']._serialized_end=246
  _globals['_SIGNATUREDESCRIPTOR']._serialized_start=249
  _globals['_SIGNATUREDESCRIPTOR']._serialized_end=878
  _globals['_SIGNATUREDESCRIPTOR_DATA']._serialized_start=427
  _globals['_SIGNATUREDESCRIPTOR_DATA']._serialized_end=878
  _globals['_SIGNATUREDESCRIPTOR_DATA_SINGLE']._serialized_start=604
  _globals['_SIGNATUREDESCRIPTOR_DATA_SINGLE']._serialized_end=699
  _globals['_SIGNATUREDESCRIPTOR_DATA_MULTI']._serialized_start=702
  _globals['_SIGNATUREDESCRIPTOR_DATA_MULTI']._serialized_end=871
# @@protoc_insertion_point(module_scope)
