# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/autocli/v1/options.proto
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
    'cosmos/autocli/v1/options.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/autocli/v1/options.proto\x12\x11\x63osmos.autocli.v1\"\x8f\x01\n\rModuleOptions\x12;\n\x02tx\x18\x01 \x01(\x0b\x32+.cosmos.autocli.v1.ServiceCommandDescriptorR\x02tx\x12\x41\n\x05query\x18\x02 \x01(\x0b\x32+.cosmos.autocli.v1.ServiceCommandDescriptorR\x05query\"\x8e\x03\n\x18ServiceCommandDescriptor\x12\x18\n\x07service\x18\x01 \x01(\tR\x07service\x12T\n\x13rpc_command_options\x18\x02 \x03(\x0b\x32$.cosmos.autocli.v1.RpcCommandOptionsR\x11rpcCommandOptions\x12_\n\x0csub_commands\x18\x03 \x03(\x0b\x32<.cosmos.autocli.v1.ServiceCommandDescriptor.SubCommandsEntryR\x0bsubCommands\x12\x34\n\x16\x65nhance_custom_command\x18\x04 \x01(\x08R\x14\x65nhanceCustomCommand\x1ak\n\x10SubCommandsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x41\n\x05value\x18\x02 \x01(\x0b\x32+.cosmos.autocli.v1.ServiceCommandDescriptorR\x05value:\x02\x38\x01\"\xbf\x04\n\x11RpcCommandOptions\x12\x1d\n\nrpc_method\x18\x01 \x01(\tR\trpcMethod\x12\x10\n\x03use\x18\x02 \x01(\tR\x03use\x12\x12\n\x04long\x18\x03 \x01(\tR\x04long\x12\x14\n\x05short\x18\x04 \x01(\tR\x05short\x12\x18\n\x07\x65xample\x18\x05 \x01(\tR\x07\x65xample\x12\x14\n\x05\x61lias\x18\x06 \x03(\tR\x05\x61lias\x12\x1f\n\x0bsuggest_for\x18\x07 \x03(\tR\nsuggestFor\x12\x1e\n\ndeprecated\x18\x08 \x01(\tR\ndeprecated\x12\x18\n\x07version\x18\t \x01(\tR\x07version\x12X\n\x0c\x66lag_options\x18\n \x03(\x0b\x32\x35.cosmos.autocli.v1.RpcCommandOptions.FlagOptionsEntryR\x0b\x66lagOptions\x12S\n\x0fpositional_args\x18\x0b \x03(\x0b\x32*.cosmos.autocli.v1.PositionalArgDescriptorR\x0epositionalArgs\x12\x12\n\x04skip\x18\x0c \x01(\x08R\x04skip\x12!\n\x0cgov_proposal\x18\r \x01(\x08R\x0bgovProposal\x1a^\n\x10\x46lagOptionsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x34\n\x05value\x18\x02 \x01(\x0b\x32\x1e.cosmos.autocli.v1.FlagOptionsR\x05value:\x02\x38\x01\"\xe5\x01\n\x0b\x46lagOptions\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x1c\n\tshorthand\x18\x02 \x01(\tR\tshorthand\x12\x14\n\x05usage\x18\x03 \x01(\tR\x05usage\x12#\n\rdefault_value\x18\x04 \x01(\tR\x0c\x64\x65\x66\x61ultValue\x12\x1e\n\ndeprecated\x18\x06 \x01(\tR\ndeprecated\x12\x31\n\x14shorthand_deprecated\x18\x07 \x01(\tR\x13shorthandDeprecated\x12\x16\n\x06hidden\x18\x08 \x01(\x08R\x06hidden\"p\n\x17PositionalArgDescriptor\x12\x1f\n\x0bproto_field\x18\x01 \x01(\tR\nprotoField\x12\x18\n\x07varargs\x18\x02 \x01(\x08R\x07varargs\x12\x1a\n\x08optional\x18\x03 \x01(\x08R\x08optionalB+Z)cosmossdk.io/api/cosmos/base/cli/v1;cliv1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.autocli.v1.options_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)cosmossdk.io/api/cosmos/base/cli/v1;cliv1'
  _globals['_SERVICECOMMANDDESCRIPTOR_SUBCOMMANDSENTRY']._loaded_options = None
  _globals['_SERVICECOMMANDDESCRIPTOR_SUBCOMMANDSENTRY']._serialized_options = b'8\001'
  _globals['_RPCCOMMANDOPTIONS_FLAGOPTIONSENTRY']._loaded_options = None
  _globals['_RPCCOMMANDOPTIONS_FLAGOPTIONSENTRY']._serialized_options = b'8\001'
  _globals['_MODULEOPTIONS']._serialized_start=55
  _globals['_MODULEOPTIONS']._serialized_end=198
  _globals['_SERVICECOMMANDDESCRIPTOR']._serialized_start=201
  _globals['_SERVICECOMMANDDESCRIPTOR']._serialized_end=599
  _globals['_SERVICECOMMANDDESCRIPTOR_SUBCOMMANDSENTRY']._serialized_start=492
  _globals['_SERVICECOMMANDDESCRIPTOR_SUBCOMMANDSENTRY']._serialized_end=599
  _globals['_RPCCOMMANDOPTIONS']._serialized_start=602
  _globals['_RPCCOMMANDOPTIONS']._serialized_end=1177
  _globals['_RPCCOMMANDOPTIONS_FLAGOPTIONSENTRY']._serialized_start=1083
  _globals['_RPCCOMMANDOPTIONS_FLAGOPTIONSENTRY']._serialized_end=1177
  _globals['_FLAGOPTIONS']._serialized_start=1180
  _globals['_FLAGOPTIONS']._serialized_end=1409
  _globals['_POSITIONALARGDESCRIPTOR']._serialized_start=1411
  _globals['_POSITIONALARGDESCRIPTOR']._serialized_end=1523
# @@protoc_insertion_point(module_scope)
