# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: google/api/visibility.proto
# Protobuf Python Version: 5.29.3
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
    3,
    '',
    'google/api/visibility.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import descriptor_pb2 as google_dot_protobuf_dot_descriptor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1bgoogle/api/visibility.proto\x12\ngoogle.api\x1a google/protobuf/descriptor.proto\">\n\nVisibility\x12\x30\n\x05rules\x18\x01 \x03(\x0b\x32\x1a.google.api.VisibilityRuleR\x05rules\"N\n\x0eVisibilityRule\x12\x1a\n\x08selector\x18\x01 \x01(\tR\x08selector\x12 \n\x0brestriction\x18\x02 \x01(\tR\x0brestriction:d\n\x0f\x65num_visibility\x12\x1c.google.protobuf.EnumOptions\x18\xaf\xca\xbc\" \x01(\x0b\x32\x1a.google.api.VisibilityRuleR\x0e\x65numVisibility:k\n\x10value_visibility\x12!.google.protobuf.EnumValueOptions\x18\xaf\xca\xbc\" \x01(\x0b\x32\x1a.google.api.VisibilityRuleR\x0fvalueVisibility:g\n\x10\x66ield_visibility\x12\x1d.google.protobuf.FieldOptions\x18\xaf\xca\xbc\" \x01(\x0b\x32\x1a.google.api.VisibilityRuleR\x0f\x66ieldVisibility:m\n\x12message_visibility\x12\x1f.google.protobuf.MessageOptions\x18\xaf\xca\xbc\" \x01(\x0b\x32\x1a.google.api.VisibilityRuleR\x11messageVisibility:j\n\x11method_visibility\x12\x1e.google.protobuf.MethodOptions\x18\xaf\xca\xbc\" \x01(\x0b\x32\x1a.google.api.VisibilityRuleR\x10methodVisibility:e\n\x0e\x61pi_visibility\x12\x1f.google.protobuf.ServiceOptions\x18\xaf\xca\xbc\" \x01(\x0b\x32\x1a.google.api.VisibilityRuleR\rapiVisibilityBn\n\x0e\x63om.google.apiB\x0fVisibilityProtoP\x01Z?google.golang.org/genproto/googleapis/api/visibility;visibility\xf8\x01\x01\xa2\x02\x04GAPIb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'google.api.visibility_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\016com.google.apiB\017VisibilityProtoP\001Z?google.golang.org/genproto/googleapis/api/visibility;visibility\370\001\001\242\002\004GAPI'
  _globals['_VISIBILITY']._serialized_start=77
  _globals['_VISIBILITY']._serialized_end=139
  _globals['_VISIBILITYRULE']._serialized_start=141
  _globals['_VISIBILITYRULE']._serialized_end=219
# @@protoc_insertion_point(module_scope)
