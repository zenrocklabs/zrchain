# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/core/client/v2/genesis.proto
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
    'ibc/core/client/v2/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from ibc.core.client.v2 import counterparty_pb2 as ibc_dot_core_dot_client_dot_v2_dot_counterparty__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n ibc/core/client/v2/genesis.proto\x12\x12ibc.core.client.v2\x1a%ibc/core/client/v2/counterparty.proto\x1a\x14gogoproto/gogo.proto\"\x8f\x01\n\x17GenesisCounterpartyInfo\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\x12W\n\x11\x63ounterparty_info\x18\x02 \x01(\x0b\x32$.ibc.core.client.v2.CounterpartyInfoB\x04\xc8\xde\x1f\x00R\x10\x63ounterpartyInfo\"p\n\x0cGenesisState\x12`\n\x12\x63ounterparty_infos\x18\x01 \x03(\x0b\x32+.ibc.core.client.v2.GenesisCounterpartyInfoB\x04\xc8\xde\x1f\x00R\x11\x63ounterpartyInfosB>Z<github.com/cosmos/ibc-go/v10/modules/core/02-client/v2/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.core.client.v2.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z<github.com/cosmos/ibc-go/v10/modules/core/02-client/v2/types'
  _globals['_GENESISCOUNTERPARTYINFO'].fields_by_name['counterparty_info']._loaded_options = None
  _globals['_GENESISCOUNTERPARTYINFO'].fields_by_name['counterparty_info']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['counterparty_infos']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['counterparty_infos']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISCOUNTERPARTYINFO']._serialized_start=118
  _globals['_GENESISCOUNTERPARTYINFO']._serialized_end=261
  _globals['_GENESISSTATE']._serialized_start=263
  _globals['_GENESISSTATE']._serialized_end=375
# @@protoc_insertion_point(module_scope)
