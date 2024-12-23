# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: google/rpc/context/attribute_context.proto
# Protobuf Python Version: 5.29.2
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
    2,
    '',
    'google/rpc/context/attribute_context.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n*google/rpc/context/attribute_context.proto\x12\x12google.rpc.context\x1a\x19google/protobuf/any.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x1cgoogle/protobuf/struct.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\x81\x14\n\x10\x41ttributeContext\x12\x41\n\x06origin\x18\x07 \x01(\x0b\x32).google.rpc.context.AttributeContext.PeerR\x06origin\x12\x41\n\x06source\x18\x01 \x01(\x0b\x32).google.rpc.context.AttributeContext.PeerR\x06source\x12K\n\x0b\x64\x65stination\x18\x02 \x01(\x0b\x32).google.rpc.context.AttributeContext.PeerR\x0b\x64\x65stination\x12\x46\n\x07request\x18\x03 \x01(\x0b\x32,.google.rpc.context.AttributeContext.RequestR\x07request\x12I\n\x08response\x18\x04 \x01(\x0b\x32-.google.rpc.context.AttributeContext.ResponseR\x08response\x12I\n\x08resource\x18\x05 \x01(\x0b\x32-.google.rpc.context.AttributeContext.ResourceR\x08resource\x12:\n\x03\x61pi\x18\x06 \x01(\x0b\x32(.google.rpc.context.AttributeContext.ApiR\x03\x61pi\x12\x34\n\nextensions\x18\x08 \x03(\x0b\x32\x14.google.protobuf.AnyR\nextensions\x1a\xf3\x01\n\x04Peer\x12\x0e\n\x02ip\x18\x01 \x01(\tR\x02ip\x12\x12\n\x04port\x18\x02 \x01(\x03R\x04port\x12M\n\x06labels\x18\x06 \x03(\x0b\x32\x35.google.rpc.context.AttributeContext.Peer.LabelsEntryR\x06labels\x12\x1c\n\tprincipal\x18\x07 \x01(\tR\tprincipal\x12\x1f\n\x0bregion_code\x18\x08 \x01(\tR\nregionCode\x1a\x39\n\x0bLabelsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\x1as\n\x03\x41pi\x12\x18\n\x07service\x18\x01 \x01(\tR\x07service\x12\x1c\n\toperation\x18\x02 \x01(\tR\toperation\x12\x1a\n\x08protocol\x18\x03 \x01(\tR\x08protocol\x12\x18\n\x07version\x18\x04 \x01(\tR\x07version\x1a\xb6\x01\n\x04\x41uth\x12\x1c\n\tprincipal\x18\x01 \x01(\tR\tprincipal\x12\x1c\n\taudiences\x18\x02 \x03(\tR\taudiences\x12\x1c\n\tpresenter\x18\x03 \x01(\tR\tpresenter\x12/\n\x06\x63laims\x18\x04 \x01(\x0b\x32\x17.google.protobuf.StructR\x06\x63laims\x12#\n\raccess_levels\x18\x05 \x03(\tR\x0c\x61\x63\x63\x65ssLevels\x1a\xcf\x03\n\x07Request\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\x12\x16\n\x06method\x18\x02 \x01(\tR\x06method\x12S\n\x07headers\x18\x03 \x03(\x0b\x32\x39.google.rpc.context.AttributeContext.Request.HeadersEntryR\x07headers\x12\x12\n\x04path\x18\x04 \x01(\tR\x04path\x12\x12\n\x04host\x18\x05 \x01(\tR\x04host\x12\x16\n\x06scheme\x18\x06 \x01(\tR\x06scheme\x12\x14\n\x05query\x18\x07 \x01(\tR\x05query\x12.\n\x04time\x18\t \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x04time\x12\x12\n\x04size\x18\n \x01(\x03R\x04size\x12\x1a\n\x08protocol\x18\x0b \x01(\tR\x08protocol\x12\x16\n\x06reason\x18\x0c \x01(\tR\x06reason\x12=\n\x04\x61uth\x18\r \x01(\x0b\x32).google.rpc.context.AttributeContext.AuthR\x04\x61uth\x1a:\n\x0cHeadersEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\x1a\xb8\x02\n\x08Response\x12\x12\n\x04\x63ode\x18\x01 \x01(\x03R\x04\x63ode\x12\x12\n\x04size\x18\x02 \x01(\x03R\x04size\x12T\n\x07headers\x18\x03 \x03(\x0b\x32:.google.rpc.context.AttributeContext.Response.HeadersEntryR\x07headers\x12.\n\x04time\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\x04time\x12\x42\n\x0f\x62\x61\x63kend_latency\x18\x05 \x01(\x0b\x32\x19.google.protobuf.DurationR\x0e\x62\x61\x63kendLatency\x1a:\n\x0cHeadersEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\x1a\x98\x05\n\x08Resource\x12\x18\n\x07service\x18\x01 \x01(\tR\x07service\x12\x12\n\x04name\x18\x02 \x01(\tR\x04name\x12\x12\n\x04type\x18\x03 \x01(\tR\x04type\x12Q\n\x06labels\x18\x04 \x03(\x0b\x32\x39.google.rpc.context.AttributeContext.Resource.LabelsEntryR\x06labels\x12\x10\n\x03uid\x18\x05 \x01(\tR\x03uid\x12`\n\x0b\x61nnotations\x18\x06 \x03(\x0b\x32>.google.rpc.context.AttributeContext.Resource.AnnotationsEntryR\x0b\x61nnotations\x12!\n\x0c\x64isplay_name\x18\x07 \x01(\tR\x0b\x64isplayName\x12;\n\x0b\x63reate_time\x18\x08 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\ncreateTime\x12;\n\x0bupdate_time\x18\t \x01(\x0b\x32\x1a.google.protobuf.TimestampR\nupdateTime\x12;\n\x0b\x64\x65lete_time\x18\n \x01(\x0b\x32\x1a.google.protobuf.TimestampR\ndeleteTime\x12\x12\n\x04\x65tag\x18\x0b \x01(\tR\x04\x65tag\x12\x1a\n\x08location\x18\x0c \x01(\tR\x08location\x1a\x39\n\x0bLabelsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\x1a>\n\x10\x41nnotationsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\x42\x8b\x01\n\x16\x63om.google.rpc.contextB\x15\x41ttributeContextProtoP\x01ZUgoogle.golang.org/genproto/googleapis/rpc/context/attribute_context;attribute_context\xf8\x01\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'google.rpc.context.attribute_context_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\026com.google.rpc.contextB\025AttributeContextProtoP\001ZUgoogle.golang.org/genproto/googleapis/rpc/context/attribute_context;attribute_context\370\001\001'
  _globals['_ATTRIBUTECONTEXT_PEER_LABELSENTRY']._loaded_options = None
  _globals['_ATTRIBUTECONTEXT_PEER_LABELSENTRY']._serialized_options = b'8\001'
  _globals['_ATTRIBUTECONTEXT_REQUEST_HEADERSENTRY']._loaded_options = None
  _globals['_ATTRIBUTECONTEXT_REQUEST_HEADERSENTRY']._serialized_options = b'8\001'
  _globals['_ATTRIBUTECONTEXT_RESPONSE_HEADERSENTRY']._loaded_options = None
  _globals['_ATTRIBUTECONTEXT_RESPONSE_HEADERSENTRY']._serialized_options = b'8\001'
  _globals['_ATTRIBUTECONTEXT_RESOURCE_LABELSENTRY']._loaded_options = None
  _globals['_ATTRIBUTECONTEXT_RESOURCE_LABELSENTRY']._serialized_options = b'8\001'
  _globals['_ATTRIBUTECONTEXT_RESOURCE_ANNOTATIONSENTRY']._loaded_options = None
  _globals['_ATTRIBUTECONTEXT_RESOURCE_ANNOTATIONSENTRY']._serialized_options = b'8\001'
  _globals['_ATTRIBUTECONTEXT']._serialized_start=189
  _globals['_ATTRIBUTECONTEXT']._serialized_end=2750
  _globals['_ATTRIBUTECONTEXT_PEER']._serialized_start=757
  _globals['_ATTRIBUTECONTEXT_PEER']._serialized_end=1000
  _globals['_ATTRIBUTECONTEXT_PEER_LABELSENTRY']._serialized_start=943
  _globals['_ATTRIBUTECONTEXT_PEER_LABELSENTRY']._serialized_end=1000
  _globals['_ATTRIBUTECONTEXT_API']._serialized_start=1002
  _globals['_ATTRIBUTECONTEXT_API']._serialized_end=1117
  _globals['_ATTRIBUTECONTEXT_AUTH']._serialized_start=1120
  _globals['_ATTRIBUTECONTEXT_AUTH']._serialized_end=1302
  _globals['_ATTRIBUTECONTEXT_REQUEST']._serialized_start=1305
  _globals['_ATTRIBUTECONTEXT_REQUEST']._serialized_end=1768
  _globals['_ATTRIBUTECONTEXT_REQUEST_HEADERSENTRY']._serialized_start=1710
  _globals['_ATTRIBUTECONTEXT_REQUEST_HEADERSENTRY']._serialized_end=1768
  _globals['_ATTRIBUTECONTEXT_RESPONSE']._serialized_start=1771
  _globals['_ATTRIBUTECONTEXT_RESPONSE']._serialized_end=2083
  _globals['_ATTRIBUTECONTEXT_RESPONSE_HEADERSENTRY']._serialized_start=1710
  _globals['_ATTRIBUTECONTEXT_RESPONSE_HEADERSENTRY']._serialized_end=1768
  _globals['_ATTRIBUTECONTEXT_RESOURCE']._serialized_start=2086
  _globals['_ATTRIBUTECONTEXT_RESOURCE']._serialized_end=2750
  _globals['_ATTRIBUTECONTEXT_RESOURCE_LABELSENTRY']._serialized_start=943
  _globals['_ATTRIBUTECONTEXT_RESOURCE_LABELSENTRY']._serialized_end=1000
  _globals['_ATTRIBUTECONTEXT_RESOURCE_ANNOTATIONSENTRY']._serialized_start=2688
  _globals['_ATTRIBUTECONTEXT_RESOURCE_ANNOTATIONSENTRY']._serialized_end=2750
# @@protoc_insertion_point(module_scope)
