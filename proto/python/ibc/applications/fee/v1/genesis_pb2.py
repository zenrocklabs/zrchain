# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/applications/fee/v1/genesis.proto
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
    'ibc/applications/fee/v1/genesis.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from ibc.applications.fee.v1 import fee_pb2 as ibc_dot_applications_dot_fee_dot_v1_dot_fee__pb2
from ibc.core.channel.v1 import channel_pb2 as ibc_dot_core_dot_channel_dot_v1_dot_channel__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n%ibc/applications/fee/v1/genesis.proto\x12\x17ibc.applications.fee.v1\x1a\x14gogoproto/gogo.proto\x1a!ibc/applications/fee/v1/fee.proto\x1a!ibc/core/channel/v1/channel.proto\"\x91\x04\n\x0cGenesisState\x12\\\n\x0fidentified_fees\x18\x01 \x03(\x0b\x32-.ibc.applications.fee.v1.IdentifiedPacketFeesB\x04\xc8\xde\x1f\x00R\x0eidentifiedFees\x12\x62\n\x14\x66\x65\x65_enabled_channels\x18\x02 \x03(\x0b\x32*.ibc.applications.fee.v1.FeeEnabledChannelB\x04\xc8\xde\x1f\x00R\x12\x66\x65\x65\x45nabledChannels\x12[\n\x11registered_payees\x18\x03 \x03(\x0b\x32(.ibc.applications.fee.v1.RegisteredPayeeB\x04\xc8\xde\x1f\x00R\x10registeredPayees\x12\x80\x01\n\x1eregistered_counterparty_payees\x18\x04 \x03(\x0b\x32\x34.ibc.applications.fee.v1.RegisteredCounterpartyPayeeB\x04\xc8\xde\x1f\x00R\x1cregisteredCounterpartyPayees\x12_\n\x10\x66orward_relayers\x18\x05 \x03(\x0b\x32..ibc.applications.fee.v1.ForwardRelayerAddressB\x04\xc8\xde\x1f\x00R\x0f\x66orwardRelayers\"K\n\x11\x46\x65\x65\x45nabledChannel\x12\x17\n\x07port_id\x18\x01 \x01(\tR\x06portId\x12\x1d\n\nchannel_id\x18\x02 \x01(\tR\tchannelId\"`\n\x0fRegisteredPayee\x12\x1d\n\nchannel_id\x18\x01 \x01(\tR\tchannelId\x12\x18\n\x07relayer\x18\x02 \x01(\tR\x07relayer\x12\x14\n\x05payee\x18\x03 \x01(\tR\x05payee\"\x85\x01\n\x1bRegisteredCounterpartyPayee\x12\x1d\n\nchannel_id\x18\x01 \x01(\tR\tchannelId\x12\x18\n\x07relayer\x18\x02 \x01(\tR\x07relayer\x12-\n\x12\x63ounterparty_payee\x18\x03 \x01(\tR\x11\x63ounterpartyPayee\"s\n\x15\x46orwardRelayerAddress\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12@\n\tpacket_id\x18\x02 \x01(\x0b\x32\x1d.ibc.core.channel.v1.PacketIdB\x04\xc8\xde\x1f\x00R\x08packetIdB7Z5github.com/cosmos/ibc-go/v9/modules/apps/29-fee/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.applications.fee.v1.genesis_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z5github.com/cosmos/ibc-go/v9/modules/apps/29-fee/types'
  _globals['_GENESISSTATE'].fields_by_name['identified_fees']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['identified_fees']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['fee_enabled_channels']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['fee_enabled_channels']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['registered_payees']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['registered_payees']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['registered_counterparty_payees']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['registered_counterparty_payees']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE'].fields_by_name['forward_relayers']._loaded_options = None
  _globals['_GENESISSTATE'].fields_by_name['forward_relayers']._serialized_options = b'\310\336\037\000'
  _globals['_FORWARDRELAYERADDRESS'].fields_by_name['packet_id']._loaded_options = None
  _globals['_FORWARDRELAYERADDRESS'].fields_by_name['packet_id']._serialized_options = b'\310\336\037\000'
  _globals['_GENESISSTATE']._serialized_start=159
  _globals['_GENESISSTATE']._serialized_end=688
  _globals['_FEEENABLEDCHANNEL']._serialized_start=690
  _globals['_FEEENABLEDCHANNEL']._serialized_end=765
  _globals['_REGISTEREDPAYEE']._serialized_start=767
  _globals['_REGISTEREDPAYEE']._serialized_end=863
  _globals['_REGISTEREDCOUNTERPARTYPAYEE']._serialized_start=866
  _globals['_REGISTEREDCOUNTERPARTYPAYEE']._serialized_end=999
  _globals['_FORWARDRELAYERADDRESS']._serialized_start=1001
  _globals['_FORWARDRELAYERADDRESS']._serialized_end=1116
# @@protoc_insertion_point(module_scope)
