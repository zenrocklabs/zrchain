# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: google/api/field_info.proto
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
    'google/api/field_info.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import descriptor_pb2 as google_dot_protobuf_dot_descriptor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1bgoogle/api/field_info.proto\x12\ngoogle.api\x1a google/protobuf/descriptor.proto\"\xda\x01\n\tFieldInfo\x12\x34\n\x06\x66ormat\x18\x01 \x01(\x0e\x32\x1c.google.api.FieldInfo.FormatR\x06\x66ormat\x12\x44\n\x10referenced_types\x18\x02 \x03(\x0b\x32\x19.google.api.TypeReferenceR\x0freferencedTypes\"Q\n\x06\x46ormat\x12\x16\n\x12\x46ORMAT_UNSPECIFIED\x10\x00\x12\t\n\x05UUID4\x10\x01\x12\x08\n\x04IPV4\x10\x02\x12\x08\n\x04IPV6\x10\x03\x12\x10\n\x0cIPV4_OR_IPV6\x10\x04\",\n\rTypeReference\x12\x1b\n\ttype_name\x18\x01 \x01(\tR\x08typeName:W\n\nfield_info\x12\x1d.google.protobuf.FieldOptions\x18\xcc\xf1\xf9\x8a\x01 \x01(\x0b\x32\x15.google.api.FieldInfoR\tfieldInfoBl\n\x0e\x63om.google.apiB\x0e\x46ieldInfoProtoP\x01ZAgoogle.golang.org/genproto/googleapis/api/annotations;annotations\xa2\x02\x04GAPIb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'google.api.field_info_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\016com.google.apiB\016FieldInfoProtoP\001ZAgoogle.golang.org/genproto/googleapis/api/annotations;annotations\242\002\004GAPI'
  _globals['_FIELDINFO']._serialized_start=78
  _globals['_FIELDINFO']._serialized_end=296
  _globals['_FIELDINFO_FORMAT']._serialized_start=215
  _globals['_FIELDINFO_FORMAT']._serialized_end=296
  _globals['_TYPEREFERENCE']._serialized_start=298
  _globals['_TYPEREFERENCE']._serialized_end=342
# @@protoc_insertion_point(module_scope)
