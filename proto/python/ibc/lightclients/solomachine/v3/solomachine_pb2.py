# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/lightclients/solomachine/v3/solomachine.proto
# Protobuf Python Version: 6.31.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    31,
    1,
    '',
    'ibc/lightclients/solomachine/v3/solomachine.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n1ibc/lightclients/solomachine/v3/solomachine.proto\x12\x1fibc.lightclients.solomachine.v3\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\"\xa6\x01\n\x0b\x43lientState\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequence\x12\x1b\n\tis_frozen\x18\x02 \x01(\x08R\x08isFrozen\x12X\n\x0f\x63onsensus_state\x18\x03 \x01(\x0b\x32/.ibc.lightclients.solomachine.v3.ConsensusStateR\x0e\x63onsensusState:\x04\x88\xa0\x1f\x00\"\x8b\x01\n\x0e\x43onsensusState\x12\x33\n\npublic_key\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\tpublicKey\x12 \n\x0b\x64iversifier\x18\x02 \x01(\tR\x0b\x64iversifier\x12\x1c\n\ttimestamp\x18\x03 \x01(\x04R\ttimestamp:\x04\x88\xa0\x1f\x00\"\xaf\x01\n\x06Header\x12\x1c\n\ttimestamp\x18\x01 \x01(\x04R\ttimestamp\x12\x1c\n\tsignature\x18\x02 \x01(\x0cR\tsignature\x12:\n\x0enew_public_key\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0cnewPublicKey\x12\'\n\x0fnew_diversifier\x18\x04 \x01(\tR\x0enewDiversifier:\x04\x88\xa0\x1f\x00\"\xe0\x01\n\x0cMisbehaviour\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequence\x12V\n\rsignature_one\x18\x02 \x01(\x0b\x32\x31.ibc.lightclients.solomachine.v3.SignatureAndDataR\x0csignatureOne\x12V\n\rsignature_two\x18\x03 \x01(\x0b\x32\x31.ibc.lightclients.solomachine.v3.SignatureAndDataR\x0csignatureTwo:\x04\x88\xa0\x1f\x00\"|\n\x10SignatureAndData\x12\x1c\n\tsignature\x18\x01 \x01(\x0cR\tsignature\x12\x12\n\x04path\x18\x02 \x01(\x0cR\x04path\x12\x12\n\x04\x64\x61ta\x18\x03 \x01(\x0cR\x04\x64\x61ta\x12\x1c\n\ttimestamp\x18\x04 \x01(\x04R\ttimestamp:\x04\x88\xa0\x1f\x00\"e\n\x18TimestampedSignatureData\x12%\n\x0esignature_data\x18\x01 \x01(\x0cR\rsignatureData\x12\x1c\n\ttimestamp\x18\x02 \x01(\x04R\ttimestamp:\x04\x88\xa0\x1f\x00\"\x95\x01\n\tSignBytes\x12\x1a\n\x08sequence\x18\x01 \x01(\x04R\x08sequence\x12\x1c\n\ttimestamp\x18\x02 \x01(\x04R\ttimestamp\x12 \n\x0b\x64iversifier\x18\x03 \x01(\tR\x0b\x64iversifier\x12\x12\n\x04path\x18\x04 \x01(\x0cR\x04path\x12\x12\n\x04\x64\x61ta\x18\x05 \x01(\x0cR\x04\x64\x61ta:\x04\x88\xa0\x1f\x00\"q\n\nHeaderData\x12\x34\n\x0bnew_pub_key\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\tnewPubKey\x12\'\n\x0fnew_diversifier\x18\x02 \x01(\tR\x0enewDiversifier:\x04\x88\xa0\x1f\x00\x42OZMgithub.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine;solomachineb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.lightclients.solomachine.v3.solomachine_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZMgithub.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine;solomachine'
  _globals['_CLIENTSTATE']._loaded_options = None
  _globals['_CLIENTSTATE']._serialized_options = b'\210\240\037\000'
  _globals['_CONSENSUSSTATE']._loaded_options = None
  _globals['_CONSENSUSSTATE']._serialized_options = b'\210\240\037\000'
  _globals['_HEADER']._loaded_options = None
  _globals['_HEADER']._serialized_options = b'\210\240\037\000'
  _globals['_MISBEHAVIOUR']._loaded_options = None
  _globals['_MISBEHAVIOUR']._serialized_options = b'\210\240\037\000'
  _globals['_SIGNATUREANDDATA']._loaded_options = None
  _globals['_SIGNATUREANDDATA']._serialized_options = b'\210\240\037\000'
  _globals['_TIMESTAMPEDSIGNATUREDATA']._loaded_options = None
  _globals['_TIMESTAMPEDSIGNATUREDATA']._serialized_options = b'\210\240\037\000'
  _globals['_SIGNBYTES']._loaded_options = None
  _globals['_SIGNBYTES']._serialized_options = b'\210\240\037\000'
  _globals['_HEADERDATA']._loaded_options = None
  _globals['_HEADERDATA']._serialized_options = b'\210\240\037\000'
  _globals['_CLIENTSTATE']._serialized_start=136
  _globals['_CLIENTSTATE']._serialized_end=302
  _globals['_CONSENSUSSTATE']._serialized_start=305
  _globals['_CONSENSUSSTATE']._serialized_end=444
  _globals['_HEADER']._serialized_start=447
  _globals['_HEADER']._serialized_end=622
  _globals['_MISBEHAVIOUR']._serialized_start=625
  _globals['_MISBEHAVIOUR']._serialized_end=849
  _globals['_SIGNATUREANDDATA']._serialized_start=851
  _globals['_SIGNATUREANDDATA']._serialized_end=975
  _globals['_TIMESTAMPEDSIGNATUREDATA']._serialized_start=977
  _globals['_TIMESTAMPEDSIGNATUREDATA']._serialized_end=1078
  _globals['_SIGNBYTES']._serialized_start=1081
  _globals['_SIGNBYTES']._serialized_end=1230
  _globals['_HEADERDATA']._serialized_start=1232
  _globals['_HEADERDATA']._serialized_end=1345
# @@protoc_insertion_point(module_scope)
