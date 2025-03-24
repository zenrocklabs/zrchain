# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/accounts/v1/query.proto
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
    'cosmos/accounts/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1e\x63osmos/accounts/v1/query.proto\x12\x12\x63osmos.accounts.v1\x1a\x19google/protobuf/any.proto\"]\n\x13\x41\x63\x63ountQueryRequest\x12\x16\n\x06target\x18\x01 \x01(\tR\x06target\x12.\n\x07request\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x07request\"H\n\x14\x41\x63\x63ountQueryResponse\x12\x30\n\x08response\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\x08response\"2\n\rSchemaRequest\x12!\n\x0c\x61\x63\x63ount_type\x18\x01 \x01(\tR\x0b\x61\x63\x63ountType\"\xc8\x02\n\x0eSchemaResponse\x12K\n\x0binit_schema\x18\x01 \x01(\x0b\x32*.cosmos.accounts.v1.SchemaResponse.HandlerR\ninitSchema\x12U\n\x10\x65xecute_handlers\x18\x02 \x03(\x0b\x32*.cosmos.accounts.v1.SchemaResponse.HandlerR\x0f\x65xecuteHandlers\x12Q\n\x0equery_handlers\x18\x03 \x03(\x0b\x32*.cosmos.accounts.v1.SchemaResponse.HandlerR\rqueryHandlers\x1a?\n\x07Handler\x12\x18\n\x07request\x18\x01 \x01(\tR\x07request\x12\x1a\n\x08response\x18\x02 \x01(\tR\x08response\".\n\x12\x41\x63\x63ountTypeRequest\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\"8\n\x13\x41\x63\x63ountTypeResponse\x12!\n\x0c\x61\x63\x63ount_type\x18\x01 \x01(\tR\x0b\x61\x63\x63ountType\"0\n\x14\x41\x63\x63ountNumberRequest\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\"/\n\x15\x41\x63\x63ountNumberResponse\x12\x16\n\x06number\x18\x01 \x01(\x04R\x06number2\x89\x03\n\x05Query\x12\x63\n\x0c\x41\x63\x63ountQuery\x12\'.cosmos.accounts.v1.AccountQueryRequest\x1a(.cosmos.accounts.v1.AccountQueryResponse\"\x00\x12Q\n\x06Schema\x12!.cosmos.accounts.v1.SchemaRequest\x1a\".cosmos.accounts.v1.SchemaResponse\"\x00\x12`\n\x0b\x41\x63\x63ountType\x12&.cosmos.accounts.v1.AccountTypeRequest\x1a\'.cosmos.accounts.v1.AccountTypeResponse\"\x00\x12\x66\n\rAccountNumber\x12(.cosmos.accounts.v1.AccountNumberRequest\x1a).cosmos.accounts.v1.AccountNumberResponse\"\x00\x42\x1cZ\x1a\x63osmossdk.io/x/accounts/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.accounts.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\032cosmossdk.io/x/accounts/v1'
  _globals['_ACCOUNTQUERYREQUEST']._serialized_start=81
  _globals['_ACCOUNTQUERYREQUEST']._serialized_end=174
  _globals['_ACCOUNTQUERYRESPONSE']._serialized_start=176
  _globals['_ACCOUNTQUERYRESPONSE']._serialized_end=248
  _globals['_SCHEMAREQUEST']._serialized_start=250
  _globals['_SCHEMAREQUEST']._serialized_end=300
  _globals['_SCHEMARESPONSE']._serialized_start=303
  _globals['_SCHEMARESPONSE']._serialized_end=631
  _globals['_SCHEMARESPONSE_HANDLER']._serialized_start=568
  _globals['_SCHEMARESPONSE_HANDLER']._serialized_end=631
  _globals['_ACCOUNTTYPEREQUEST']._serialized_start=633
  _globals['_ACCOUNTTYPEREQUEST']._serialized_end=679
  _globals['_ACCOUNTTYPERESPONSE']._serialized_start=681
  _globals['_ACCOUNTTYPERESPONSE']._serialized_end=737
  _globals['_ACCOUNTNUMBERREQUEST']._serialized_start=739
  _globals['_ACCOUNTNUMBERREQUEST']._serialized_end=787
  _globals['_ACCOUNTNUMBERRESPONSE']._serialized_start=789
  _globals['_ACCOUNTNUMBERRESPONSE']._serialized_end=836
  _globals['_QUERY']._serialized_start=839
  _globals['_QUERY']._serialized_end=1232
# @@protoc_insertion_point(module_scope)
