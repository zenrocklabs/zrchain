# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/orm/v1alpha1/schema.proto
# Protobuf Python Version: 5.28.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    2,
    '',
    'cosmos/orm/v1alpha1/schema.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import descriptor_pb2 as google_dot_protobuf_dot_descriptor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n cosmos/orm/v1alpha1/schema.proto\x12\x13\x63osmos.orm.v1alpha1\x1a google/protobuf/descriptor.proto\"\x93\x02\n\x16ModuleSchemaDescriptor\x12V\n\x0bschema_file\x18\x01 \x03(\x0b\x32\x35.cosmos.orm.v1alpha1.ModuleSchemaDescriptor.FileEntryR\nschemaFile\x12\x16\n\x06prefix\x18\x02 \x01(\x0cR\x06prefix\x1a\x88\x01\n\tFileEntry\x12\x0e\n\x02id\x18\x01 \x01(\rR\x02id\x12&\n\x0fproto_file_name\x18\x02 \x01(\tR\rprotoFileName\x12\x43\n\x0cstorage_type\x18\x03 \x01(\x0e\x32 .cosmos.orm.v1alpha1.StorageTypeR\x0bstorageType*h\n\x0bStorageType\x12$\n STORAGE_TYPE_DEFAULT_UNSPECIFIED\x10\x00\x12\x17\n\x13STORAGE_TYPE_MEMORY\x10\x01\x12\x1a\n\x16STORAGE_TYPE_TRANSIENT\x10\x02:t\n\rmodule_schema\x12\x1f.google.protobuf.MessageOptions\x18\xf0\xb3\xea\x31 \x01(\x0b\x32+.cosmos.orm.v1alpha1.ModuleSchemaDescriptorR\x0cmoduleSchemab\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.orm.v1alpha1.schema_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_STORAGETYPE']._serialized_start=369
  _globals['_STORAGETYPE']._serialized_end=473
  _globals['_MODULESCHEMADESCRIPTOR']._serialized_start=92
  _globals['_MODULESCHEMADESCRIPTOR']._serialized_end=367
  _globals['_MODULESCHEMADESCRIPTOR_FILEENTRY']._serialized_start=231
  _globals['_MODULESCHEMADESCRIPTOR_FILEENTRY']._serialized_end=367
# @@protoc_insertion_point(module_scope)