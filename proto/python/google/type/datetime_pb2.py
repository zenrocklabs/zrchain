# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: google/type/datetime.proto
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
    'google/type/datetime.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1agoogle/type/datetime.proto\x12\x0bgoogle.type\x1a\x1egoogle/protobuf/duration.proto\"\xa7\x02\n\x08\x44\x61teTime\x12\x12\n\x04year\x18\x01 \x01(\x05R\x04year\x12\x14\n\x05month\x18\x02 \x01(\x05R\x05month\x12\x10\n\x03\x64\x61y\x18\x03 \x01(\x05R\x03\x64\x61y\x12\x14\n\x05hours\x18\x04 \x01(\x05R\x05hours\x12\x18\n\x07minutes\x18\x05 \x01(\x05R\x07minutes\x12\x18\n\x07seconds\x18\x06 \x01(\x05R\x07seconds\x12\x14\n\x05nanos\x18\x07 \x01(\x05R\x05nanos\x12:\n\nutc_offset\x18\x08 \x01(\x0b\x32\x19.google.protobuf.DurationH\x00R\tutcOffset\x12\x34\n\ttime_zone\x18\t \x01(\x0b\x32\x15.google.type.TimeZoneH\x00R\x08timeZoneB\r\n\x0btime_offset\"4\n\x08TimeZone\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\x12\x18\n\x07version\x18\x02 \x01(\tR\x07versionBi\n\x0f\x63om.google.typeB\rDateTimeProtoP\x01Z<google.golang.org/genproto/googleapis/type/datetime;datetime\xf8\x01\x01\xa2\x02\x03GTPb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'google.type.datetime_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\017com.google.typeB\rDateTimeProtoP\001Z<google.golang.org/genproto/googleapis/type/datetime;datetime\370\001\001\242\002\003GTP'
  _globals['_DATETIME']._serialized_start=76
  _globals['_DATETIME']._serialized_end=371
  _globals['_TIMEZONE']._serialized_start=373
  _globals['_TIMEZONE']._serialized_end=425
# @@protoc_insertion_point(module_scope)
