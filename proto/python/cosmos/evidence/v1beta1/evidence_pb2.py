# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/evidence/v1beta1/evidence.proto
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
    'cosmos/evidence/v1beta1/evidence.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n&cosmos/evidence/v1beta1/evidence.proto\x12\x17\x63osmos.evidence.v1beta1\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\xe8\x01\n\x0c\x45quivocation\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12=\n\x04time\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\r\xc8\xde\x1f\x00\x90\xdf\x1f\x01\xa8\xe7\xb0*\x01R\x04time\x12\x14\n\x05power\x18\x03 \x01(\x03R\x05power\x12\x45\n\x11\x63onsensus_address\x18\x04 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10\x63onsensusAddress:$\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\x8a\xe7\xb0*\x17\x63osmos-sdk/EquivocationB#Z\x1d\x63osmossdk.io/x/evidence/types\xa8\xe2\x1e\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.evidence.v1beta1.evidence_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\035cosmossdk.io/x/evidence/types\250\342\036\001'
  _globals['_EQUIVOCATION'].fields_by_name['time']._loaded_options = None
  _globals['_EQUIVOCATION'].fields_by_name['time']._serialized_options = b'\310\336\037\000\220\337\037\001\250\347\260*\001'
  _globals['_EQUIVOCATION'].fields_by_name['consensus_address']._loaded_options = None
  _globals['_EQUIVOCATION'].fields_by_name['consensus_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_EQUIVOCATION']._loaded_options = None
  _globals['_EQUIVOCATION']._serialized_options = b'\210\240\037\000\350\240\037\000\212\347\260*\027cosmos-sdk/Equivocation'
  _globals['_EQUIVOCATION']._serialized_start=169
  _globals['_EQUIVOCATION']._serialized_end=401
# @@protoc_insertion_point(module_scope)
