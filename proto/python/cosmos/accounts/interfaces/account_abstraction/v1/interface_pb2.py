# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/accounts/interfaces/account_abstraction/v1/interface.proto
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
    'cosmos/accounts/interfaces/account_abstraction/v1/interface.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.tx.v1beta1 import tx_pb2 as cosmos_dot_tx_dot_v1beta1_dot_tx__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\nAcosmos/accounts/interfaces/account_abstraction/v1/interface.proto\x12\x31\x63osmos.accounts.interfaces.account_abstraction.v1\x1a\x1a\x63osmos/tx/v1beta1/tx.proto\"\xa6\x01\n\x0fMsgAuthenticate\x12\x18\n\x07\x62undler\x18\x01 \x01(\tR\x07\x62undler\x12/\n\x06raw_tx\x18\x02 \x01(\x0b\x32\x18.cosmos.tx.v1beta1.TxRawR\x05rawTx\x12%\n\x02tx\x18\x03 \x01(\x0b\x32\x15.cosmos.tx.v1beta1.TxR\x02tx\x12!\n\x0csigner_index\x18\x04 \x01(\rR\x0bsignerIndex\"\x19\n\x17MsgAuthenticateResponse\"\x1c\n\x1aQueryAuthenticationMethods\"[\n\"QueryAuthenticationMethodsResponse\x12\x35\n\x16\x61uthentication_methods\x18\x01 \x03(\tR\x15\x61uthenticationMethodsB;Z9cosmossdk.io/x/accounts/interfaces/account_abstraction/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.accounts.interfaces.account_abstraction.v1.interface_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9cosmossdk.io/x/accounts/interfaces/account_abstraction/v1'
  _globals['_MSGAUTHENTICATE']._serialized_start=149
  _globals['_MSGAUTHENTICATE']._serialized_end=315
  _globals['_MSGAUTHENTICATERESPONSE']._serialized_start=317
  _globals['_MSGAUTHENTICATERESPONSE']._serialized_end=342
  _globals['_QUERYAUTHENTICATIONMETHODS']._serialized_start=344
  _globals['_QUERYAUTHENTICATIONMETHODS']._serialized_end=372
  _globals['_QUERYAUTHENTICATIONMETHODSRESPONSE']._serialized_start=374
  _globals['_QUERYAUTHENTICATIONMETHODSRESPONSE']._serialized_end=465
# @@protoc_insertion_point(module_scope)
