# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos_proto/cosmos.proto
# Protobuf Python Version: 6.30.1
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
    1,
    '',
    'cosmos_proto/cosmos.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import descriptor_pb2 as google_dot_protobuf_dot_descriptor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x19\x63osmos_proto/cosmos.proto\x12\x0c\x63osmos_proto\x1a google/protobuf/descriptor.proto\"K\n\x13InterfaceDescriptor\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\"\x81\x01\n\x10ScalarDescriptor\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x37\n\nfield_type\x18\x03 \x03(\x0e\x32\x18.cosmos_proto.ScalarTypeR\tfieldType*X\n\nScalarType\x12\x1b\n\x17SCALAR_TYPE_UNSPECIFIED\x10\x00\x12\x16\n\x12SCALAR_TYPE_STRING\x10\x01\x12\x15\n\x11SCALAR_TYPE_BYTES\x10\x02:H\n\x0fmethod_added_in\x12\x1e.google.protobuf.MethodOptions\x18\xc9\xd6\x05 \x01(\tR\rmethodAddedIn:T\n\x14implements_interface\x12\x1f.google.protobuf.MessageOptions\x18\xc9\xd6\x05 \x03(\tR\x13implementsInterface:K\n\x10message_added_in\x12\x1f.google.protobuf.MessageOptions\x18\xca\xd6\x05 \x01(\tR\x0emessageAddedIn:L\n\x11\x61\x63\x63\x65pts_interface\x12\x1d.google.protobuf.FieldOptions\x18\xc9\xd6\x05 \x01(\tR\x10\x61\x63\x63\x65ptsInterface:7\n\x06scalar\x12\x1d.google.protobuf.FieldOptions\x18\xca\xd6\x05 \x01(\tR\x06scalar:E\n\x0e\x66ield_added_in\x12\x1d.google.protobuf.FieldOptions\x18\xcb\xd6\x05 \x01(\tR\x0c\x66ieldAddedIn:n\n\x11\x64\x65\x63lare_interface\x12\x1c.google.protobuf.FileOptions\x18\xbd\xb3\x30 \x03(\x0b\x32!.cosmos_proto.InterfaceDescriptorR\x10\x64\x65\x63lareInterface:e\n\x0e\x64\x65\x63lare_scalar\x12\x1c.google.protobuf.FileOptions\x18\xbe\xb3\x30 \x03(\x0b\x32\x1e.cosmos_proto.ScalarDescriptorR\rdeclareScalar:B\n\rfile_added_in\x12\x1c.google.protobuf.FileOptions\x18\xbf\xb3\x30 \x01(\tR\x0b\x66ileAddedInB-Z+github.com/cosmos/cosmos-proto;cosmos_protob\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos_proto.cosmos_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+github.com/cosmos/cosmos-proto;cosmos_proto'
  _globals['_SCALARTYPE']._serialized_start=286
  _globals['_SCALARTYPE']._serialized_end=374
  _globals['_INTERFACEDESCRIPTOR']._serialized_start=77
  _globals['_INTERFACEDESCRIPTOR']._serialized_end=152
  _globals['_SCALARDESCRIPTOR']._serialized_start=155
  _globals['_SCALARDESCRIPTOR']._serialized_end=284
# @@protoc_insertion_point(module_scope)
