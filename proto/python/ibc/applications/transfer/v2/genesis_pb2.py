# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/transfer/v2/genesis.proto
<<<<<<< HEAD
# Protobuf Python Version: 5.29.1
=======
# Protobuf Python Version: 5.29.0
>>>>>>> main
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
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'ibc/applications/transfer/v2/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from ibc.applications.transfer.v1 import transfer_pb2 as ibc_dot_applications_dot_transfer_dot_v1_dot_transfer__pb2
from ibc.applications.transfer.v2 import token_pb2 as ibc_dot_applications_dot_transfer_dot_v2_dot_token__pb2
from ibc.core.channel.v1 import channel_pb2 as ibc_dot_core_dot_channel_dot_v1_dot_channel__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n*ibc/applications/transfer/v2/genesis.proto\x12\x1cibc.applications.transfer.v2\x1a+ibc/applications/transfer/v1/transfer.proto\x1a(ibc/applications/transfer/v2/token.proto\x1a!ibc/core/channel/v1/channel.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x14gogoproto/gogo.proto\"\x8e\x03\n\x0cGenesisState\x12\x17\n\x07port_id\x18\x01 \x01(\tR\x06portId\x12K\n\x06\x64\x65noms\x18\x02 \x03(\x0b\x32#.ibc.applications.transfer.v2.DenomB\x0e\xc8\xde\x1f\x00\xaa\xdf\x1f\x06\x44\x65nomsR\x06\x64\x65noms\x12\x42\n\x06params\x18\x03 \x01(\x0b\x32$.ibc.applications.transfer.v1.ParamsB\x04\xc8\xde\x1f\x00R\x06params\x12r\n\x0etotal_escrowed\x18\x04 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinB0\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.CoinsR\rtotalEscrowed\x12`\n\x11\x66orwarded_packets\x18\x05 \x03(\x0b\x32-.ibc.applications.transfer.v2.ForwardedPacketB\x04\xc8\xde\x1f\x00R\x10\x66orwardedPackets\"\x92\x01\n\x0f\x46orwardedPacket\x12\x44\n\x0b\x66orward_key\x18\x01 \x01(\x0b\x32\x1d.ibc.core.channel.v1.PacketIdB\x04\xc8\xde\x1f\x00R\nforwardKey\x12\x39\n\x06packet\x18\x02 \x01(\x0b\x32\x1b.ibc.core.channel.v1.PacketB\x04\xc8\xde\x1f\x00R\x06packetB9Z7github.com/cosmos/ibc-go/v9/modules/apps/transfer/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.transfer.v2.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/cosmos/ibc-go/v9/modules/apps/transfer/types'
  _globals['_GENESISSTATE'].fields_by_name['denoms']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['denoms']._serialized_options = b'\310\336\037\000\252\337\037\006Denoms'
  _globals['_GENESISSTATE'].fields_by_name['params']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['total_escrowed']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['total_escrowed']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins'
  _globals['_GENESISSTATE'].fields_by_name['forwarded_packets']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['forwarded_packets']._serialized_options = b'\310\336\037\000'
  _globals['_FORWARDEDPACKET'].fields_by_name['forward_key']._loaded_options = None
  _globals['_FORWARDEDPACKET'].fields_by_name['forward_key']._serialized_options = b'\310\336\037\000'
  _globals['_FORWARDEDPACKET'].fields_by_name['packet']._loaded_options = None
  _globals['_FORWARDEDPACKET'].fields_by_name['packet']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE']._serialized_start=253
  _globals['_GENESISSTATE']._serialized_end=651
  _globals['_FORWARDEDPACKET']._serialized_start=654
  _globals['_FORWARDEDPACKET']._serialized_end=800
# @@protoc_insertion_point(module_scope)
