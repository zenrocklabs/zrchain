# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/validation/asset_data.proto
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
    'zrchain/validation/asset_data.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n#zrchain/validation/asset_data.proto\x12\x12zrchain.validation\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\"\xa9\x01\n\tAssetData\x12/\n\x05\x61sset\x18\x01 \x01(\x0e\x32\x19.zrchain.validation.AssetR\x05\x61sset\x12M\n\x08priceUSD\x18\x02 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\x08priceUSD\x12\x1c\n\tprecision\x18\x03 \x01(\rR\tprecision*4\n\x05\x41sset\x12\x0f\n\x0bUNSPECIFIED\x10\x00\x12\x08\n\x04ROCK\x10\x01\x12\x07\n\x03\x42TC\x10\x02\x12\x07\n\x03\x45TH\x10\x03\x42=Z;github.com/Zenrock-Foundation/zrchain/v5/x/validation/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.validation.asset_data_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z;github.com/Zenrock-Foundation/zrchain/v5/x/validation/types'
  _globals['_ASSETDATA'].fields_by_name['priceUSD']._loaded_options = None
  _globals['_ASSETDATA'].fields_by_name['priceUSD']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_ASSET']._serialized_start=280
  _globals['_ASSET']._serialized_end=332
  _globals['_ASSETDATA']._serialized_start=109
  _globals['_ASSETDATA']._serialized_end=278
# @@protoc_insertion_point(module_scope)
