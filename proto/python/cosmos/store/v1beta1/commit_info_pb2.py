# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/store/v1beta1/commit_info.proto
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
    'cosmos/store/v1beta1/commit_info.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n&cosmos/store/v1beta1/commit_info.proto\x12\x14\x63osmos.store.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xb2\x01\n\nCommitInfo\x12\x18\n\x07version\x18\x01 \x01(\x03R\x07version\x12\x46\n\x0bstore_infos\x18\x02 \x03(\x0b\x32\x1f.cosmos.store.v1beta1.StoreInfoB\x04\xc8\xde\x1f\x00R\nstoreInfos\x12\x42\n\ttimestamp\x18\x03 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\x08\xc8\xde\x1f\x00\x90\xdf\x1f\x01R\ttimestamp\"b\n\tStoreInfo\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x41\n\tcommit_id\x18\x02 \x01(\x0b\x32\x1e.cosmos.store.v1beta1.CommitIDB\x04\xc8\xde\x1f\x00R\x08\x63ommitId\">\n\x08\x43ommitID\x12\x18\n\x07version\x18\x01 \x01(\x03R\x07version\x12\x12\n\x04hash\x18\x02 \x01(\x0cR\x04hash:\x04\x98\xa0\x1f\x00\x42\x1aZ\x18\x63osmossdk.io/store/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.store.v1beta1.commit_info_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\030cosmossdk.io/store/types'
  _globals['_COMMITINFO'].fields_by_name['store_infos']._loaded_options = None
  _globals['_COMMITINFO'].fields_by_name['store_infos']._serialized_options = b'\310\336\037\000'
  _globals['_COMMITINFO'].fields_by_name['timestamp']._loaded_options = None
  _globals['_COMMITINFO'].fields_by_name['timestamp']._serialized_options = b'\310\336\037\000\220\337\037\001'
  _globals['_STOREINFO'].fields_by_name['commit_id']._loaded_options = None
  _globals['_STOREINFO'].fields_by_name['commit_id']._serialized_options = b'\310\336\037\000'
  _globals['_COMMITID']._loaded_options = None
  _globals['_COMMITID']._serialized_options = b'\230\240\037\000'
  _globals['_COMMITINFO']._serialized_start=120
  _globals['_COMMITINFO']._serialized_end=298
  _globals['_STOREINFO']._serialized_start=300
  _globals['_STOREINFO']._serialized_end=398
  _globals['_COMMITID']._serialized_start=400
  _globals['_COMMITID']._serialized_end=462
# @@protoc_insertion_point(module_scope)
