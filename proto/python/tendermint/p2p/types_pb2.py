# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: tendermint/p2p/types.proto
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
    'tendermint/p2p/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1atendermint/p2p/types.proto\x12\x0etendermint.p2p\x1a\x14gogoproto/gogo.proto\"P\n\nNetAddress\x12\x16\n\x02id\x18\x01 \x01(\tB\x06\xe2\xde\x1f\x02IDR\x02id\x12\x16\n\x02ip\x18\x02 \x01(\tB\x06\xe2\xde\x1f\x02IPR\x02ip\x12\x12\n\x04port\x18\x03 \x01(\rR\x04port\"T\n\x0fProtocolVersion\x12\x19\n\x03p2p\x18\x01 \x01(\x04\x42\x07\xe2\xde\x1f\x03P2PR\x03p2p\x12\x14\n\x05\x62lock\x18\x02 \x01(\x04R\x05\x62lock\x12\x10\n\x03\x61pp\x18\x03 \x01(\x04R\x03\x61pp\"\xeb\x02\n\x0f\x44\x65\x66\x61ultNodeInfo\x12P\n\x10protocol_version\x18\x01 \x01(\x0b\x32\x1f.tendermint.p2p.ProtocolVersionB\x04\xc8\xde\x1f\x00R\x0fprotocolVersion\x12\x39\n\x0f\x64\x65\x66\x61ult_node_id\x18\x02 \x01(\tB\x11\xe2\xde\x1f\rDefaultNodeIDR\rdefaultNodeId\x12\x1f\n\x0blisten_addr\x18\x03 \x01(\tR\nlistenAddr\x12\x18\n\x07network\x18\x04 \x01(\tR\x07network\x12\x18\n\x07version\x18\x05 \x01(\tR\x07version\x12\x1a\n\x08\x63hannels\x18\x06 \x01(\x0cR\x08\x63hannels\x12\x18\n\x07moniker\x18\x07 \x01(\tR\x07moniker\x12@\n\x05other\x18\x08 \x01(\x0b\x32$.tendermint.p2p.DefaultNodeInfoOtherB\x04\xc8\xde\x1f\x00R\x05other\"b\n\x14\x44\x65\x66\x61ultNodeInfoOther\x12\x19\n\x08tx_index\x18\x01 \x01(\tR\x07txIndex\x12/\n\x0brpc_address\x18\x02 \x01(\tB\x0e\xe2\xde\x1f\nRPCAddressR\nrpcAddressB3Z1github.com/cometbft/cometbft/proto/tendermint/p2pb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'tendermint.p2p.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z1github.com/cometbft/cometbft/proto/tendermint/p2p'
  _globals['_NETADDRESS'].fields_by_name['id']._loaded_options = None
  _globals['_NETADDRESS'].fields_by_name['id']._serialized_options = b'\342\336\037\002ID'
  _globals['_NETADDRESS'].fields_by_name['ip']._loaded_options = None
  _globals['_NETADDRESS'].fields_by_name['ip']._serialized_options = b'\342\336\037\002IP'
  _globals['_PROTOCOLVERSION'].fields_by_name['p2p']._loaded_options = None
  _globals['_PROTOCOLVERSION'].fields_by_name['p2p']._serialized_options = b'\342\336\037\003P2P'
  _globals['_DEFAULTNODEINFO'].fields_by_name['protocol_version']._loaded_options = None
  _globals['_DEFAULTNODEINFO'].fields_by_name['protocol_version']._serialized_options = b'\310\336\037\000'
  _globals['_DEFAULTNODEINFO'].fields_by_name['default_node_id']._loaded_options = None
  _globals['_DEFAULTNODEINFO'].fields_by_name['default_node_id']._serialized_options = b'\342\336\037\rDefaultNodeID'
  _globals['_DEFAULTNODEINFO'].fields_by_name['other']._loaded_options = None
  _globals['_DEFAULTNODEINFO'].fields_by_name['other']._serialized_options = b'\310\336\037\000'
  _globals['_DEFAULTNODEINFOOTHER'].fields_by_name['rpc_address']._loaded_options = None
  _globals['_DEFAULTNODEINFOOTHER'].fields_by_name['rpc_address']._serialized_options = b'\342\336\037\nRPCAddress'
  _globals['_NETADDRESS']._serialized_start=68
  _globals['_NETADDRESS']._serialized_end=148
  _globals['_PROTOCOLVERSION']._serialized_start=150
  _globals['_PROTOCOLVERSION']._serialized_end=234
  _globals['_DEFAULTNODEINFO']._serialized_start=237
  _globals['_DEFAULTNODEINFO']._serialized_end=600
  _globals['_DEFAULTNODEINFOOTHER']._serialized_start=602
  _globals['_DEFAULTNODEINFOOTHER']._serialized_end=700
# @@protoc_insertion_point(module_scope)
