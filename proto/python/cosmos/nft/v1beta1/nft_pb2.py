# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/nft/v1beta1/nft.proto
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
    'cosmos/nft/v1beta1/nft.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1c\x63osmos/nft/v1beta1/nft.proto\x12\x12\x63osmos.nft.v1beta1\x1a\x19google/protobuf/any.proto\"\xbc\x01\n\x05\x43lass\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n\x06symbol\x18\x03 \x01(\tR\x06symbol\x12 \n\x0b\x64\x65scription\x18\x04 \x01(\tR\x0b\x64\x65scription\x12\x10\n\x03uri\x18\x05 \x01(\tR\x03uri\x12\x19\n\x08uri_hash\x18\x06 \x01(\tR\x07uriHash\x12(\n\x04\x64\x61ta\x18\x07 \x01(\x0b\x32\x14.google.protobuf.AnyR\x04\x64\x61ta\"\x87\x01\n\x03NFT\x12\x19\n\x08\x63lass_id\x18\x01 \x01(\tR\x07\x63lassId\x12\x0e\n\x02id\x18\x02 \x01(\tR\x02id\x12\x10\n\x03uri\x18\x03 \x01(\tR\x03uri\x12\x19\n\x08uri_hash\x18\x04 \x01(\tR\x07uriHash\x12(\n\x04\x64\x61ta\x18\n \x01(\x0b\x32\x14.google.protobuf.AnyR\x04\x64\x61taB\x14Z\x12\x63osmossdk.io/x/nftb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.nft.v1beta1.nft_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\022cosmossdk.io/x/nft'
  _globals['_CLASS']._serialized_start=80
  _globals['_CLASS']._serialized_end=268
  _globals['_NFT']._serialized_start=271
  _globals['_NFT']._serialized_end=406
# @@protoc_insertion_point(module_scope)
