# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/app/v1alpha1/module.proto
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
    'cosmos/app/v1alpha1/module.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import descriptor_pb2 as google_dot_protobuf_dot_descriptor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n cosmos/app/v1alpha1/module.proto\x12\x13\x63osmos.app.v1alpha1\x1a google/protobuf/descriptor.proto\"\xc7\x01\n\x10ModuleDescriptor\x12\x1b\n\tgo_import\x18\x01 \x01(\tR\x08goImport\x12\x46\n\x0buse_package\x18\x02 \x03(\x0b\x32%.cosmos.app.v1alpha1.PackageReferenceR\nusePackage\x12N\n\x10\x63\x61n_migrate_from\x18\x03 \x03(\x0b\x32$.cosmos.app.v1alpha1.MigrateFromInfoR\x0e\x63\x61nMigrateFrom\"B\n\x10PackageReference\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x1a\n\x08revision\x18\x02 \x01(\rR\x08revision\")\n\x0fMigrateFromInfo\x12\x16\n\x06module\x18\x01 \x01(\tR\x06module:a\n\x06module\x12\x1f.google.protobuf.MessageOptions\x18\x87\xe8\xa2\x1b \x01(\x0b\x32%.cosmos.app.v1alpha1.ModuleDescriptorR\x06moduleb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.app.v1alpha1.module_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_MODULEDESCRIPTOR']._serialized_start=92
  _globals['_MODULEDESCRIPTOR']._serialized_end=291
  _globals['_PACKAGEREFERENCE']._serialized_start=293
  _globals['_PACKAGEREFERENCE']._serialized_end=359
  _globals['_MIGRATEFROMINFO']._serialized_start=361
  _globals['_MIGRATEFROMINFO']._serialized_end=402
# @@protoc_insertion_point(module_scope)
