# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/store/snapshots/v1/snapshot.proto
# Protobuf Python Version: 6.30.2
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
    2,
    '',
    'cosmos/store/snapshots/v1/snapshot.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n(cosmos/store/snapshots/v1/snapshot.proto\x12\x19\x63osmos.store.snapshots.v1\x1a\x14gogoproto/gogo.proto\"\xad\x01\n\x08Snapshot\x12\x16\n\x06height\x18\x01 \x01(\x04R\x06height\x12\x16\n\x06\x66ormat\x18\x02 \x01(\rR\x06\x66ormat\x12\x16\n\x06\x63hunks\x18\x03 \x01(\rR\x06\x63hunks\x12\x12\n\x04hash\x18\x04 \x01(\x0cR\x04hash\x12\x45\n\x08metadata\x18\x05 \x01(\x0b\x32#.cosmos.store.snapshots.v1.MetadataB\x04\xc8\xde\x1f\x00R\x08metadata\"-\n\x08Metadata\x12!\n\x0c\x63hunk_hashes\x18\x01 \x03(\x0cR\x0b\x63hunkHashes\"\xdf\x02\n\x0cSnapshotItem\x12\x44\n\x05store\x18\x01 \x01(\x0b\x32,.cosmos.store.snapshots.v1.SnapshotStoreItemH\x00R\x05store\x12K\n\x04iavl\x18\x02 \x01(\x0b\x32+.cosmos.store.snapshots.v1.SnapshotIAVLItemB\x08\xe2\xde\x1f\x04IAVLH\x00R\x04iavl\x12P\n\textension\x18\x03 \x01(\x0b\x32\x30.cosmos.store.snapshots.v1.SnapshotExtensionMetaH\x00R\textension\x12\x62\n\x11\x65xtension_payload\x18\x04 \x01(\x0b\x32\x33.cosmos.store.snapshots.v1.SnapshotExtensionPayloadH\x00R\x10\x65xtensionPayloadB\x06\n\x04item\"\'\n\x11SnapshotStoreItem\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\"l\n\x10SnapshotIAVLItem\x12\x10\n\x03key\x18\x01 \x01(\x0cR\x03key\x12\x14\n\x05value\x18\x02 \x01(\x0cR\x05value\x12\x18\n\x07version\x18\x03 \x01(\x03R\x07version\x12\x16\n\x06height\x18\x04 \x01(\x05R\x06height\"C\n\x15SnapshotExtensionMeta\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x16\n\x06\x66ormat\x18\x02 \x01(\rR\x06\x66ormat\"4\n\x18SnapshotExtensionPayload\x12\x18\n\x07payload\x18\x01 \x01(\x0cR\x07payloadB$Z\"cosmossdk.io/store/snapshots/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.store.snapshots.v1.snapshot_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\"cosmossdk.io/store/snapshots/types'
  _globals['_SNAPSHOT'].fields_by_name['metadata']._loaded_options = None
  _globals['_SNAPSHOT'].fields_by_name['metadata']._serialized_options = b'\310\336\037\000'
  _globals['_SNAPSHOTITEM'].fields_by_name['iavl']._loaded_options = None
  _globals['_SNAPSHOTITEM'].fields_by_name['iavl']._serialized_options = b'\342\336\037\004IAVL'
  _globals['_SNAPSHOT']._serialized_start=94
  _globals['_SNAPSHOT']._serialized_end=267
  _globals['_METADATA']._serialized_start=269
  _globals['_METADATA']._serialized_end=314
  _globals['_SNAPSHOTITEM']._serialized_start=317
  _globals['_SNAPSHOTITEM']._serialized_end=668
  _globals['_SNAPSHOTSTOREITEM']._serialized_start=670
  _globals['_SNAPSHOTSTOREITEM']._serialized_end=709
  _globals['_SNAPSHOTIAVLITEM']._serialized_start=711
  _globals['_SNAPSHOTIAVLITEM']._serialized_end=819
  _globals['_SNAPSHOTEXTENSIONMETA']._serialized_start=821
  _globals['_SNAPSHOTEXTENSIONMETA']._serialized_end=888
  _globals['_SNAPSHOTEXTENSIONPAYLOAD']._serialized_start=890
  _globals['_SNAPSHOTEXTENSIONPAYLOAD']._serialized_end=942
# @@protoc_insertion_point(module_scope)
