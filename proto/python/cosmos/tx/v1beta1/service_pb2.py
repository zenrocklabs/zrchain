# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/tx/v1beta1/service.proto
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
    'cosmos/tx/v1beta1/service.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.base.abci.v1beta1 import abci_pb2 as cosmos_dot_base_dot_abci_dot_v1beta1_dot_abci__pb2
from cosmos.tx.v1beta1 import tx_pb2 as cosmos_dot_tx_dot_v1beta1_dot_tx__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from tendermint.types import block_pb2 as tendermint_dot_types_dot_block__pb2
from tendermint.types import types_pb2 as tendermint_dot_types_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/tx/v1beta1/service.proto\x12\x11\x63osmos.tx.v1beta1\x1a\x1cgoogle/api/annotations.proto\x1a#cosmos/base/abci/v1beta1/abci.proto\x1a\x1a\x63osmos/tx/v1beta1/tx.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x1ctendermint/types/block.proto\x1a\x1ctendermint/types/types.proto\"\xf3\x01\n\x12GetTxsEventRequest\x12\x1a\n\x06\x65vents\x18\x01 \x03(\tB\x02\x18\x01R\x06\x65vents\x12J\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestB\x02\x18\x01R\npagination\x12\x35\n\x08order_by\x18\x03 \x01(\x0e\x32\x1a.cosmos.tx.v1beta1.OrderByR\x07orderBy\x12\x12\n\x04page\x18\x04 \x01(\x04R\x04page\x12\x14\n\x05limit\x18\x05 \x01(\x04R\x05limit\x12\x14\n\x05query\x18\x06 \x01(\tR\x05query\"\xea\x01\n\x13GetTxsEventResponse\x12\'\n\x03txs\x18\x01 \x03(\x0b\x32\x15.cosmos.tx.v1beta1.TxR\x03txs\x12G\n\x0ctx_responses\x18\x02 \x03(\x0b\x32$.cosmos.base.abci.v1beta1.TxResponseR\x0btxResponses\x12K\n\npagination\x18\x03 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseB\x02\x18\x01R\npagination\x12\x14\n\x05total\x18\x04 \x01(\x04R\x05total\"e\n\x12\x42roadcastTxRequest\x12\x19\n\x08tx_bytes\x18\x01 \x01(\x0cR\x07txBytes\x12\x34\n\x04mode\x18\x02 \x01(\x0e\x32 .cosmos.tx.v1beta1.BroadcastModeR\x04mode\"\\\n\x13\x42roadcastTxResponse\x12\x45\n\x0btx_response\x18\x01 \x01(\x0b\x32$.cosmos.base.abci.v1beta1.TxResponseR\ntxResponse\"W\n\x0fSimulateRequest\x12)\n\x02tx\x18\x01 \x01(\x0b\x32\x15.cosmos.tx.v1beta1.TxB\x02\x18\x01R\x02tx\x12\x19\n\x08tx_bytes\x18\x02 \x01(\x0cR\x07txBytes\"\x8a\x01\n\x10SimulateResponse\x12<\n\x08gas_info\x18\x01 \x01(\x0b\x32!.cosmos.base.abci.v1beta1.GasInfoR\x07gasInfo\x12\x38\n\x06result\x18\x02 \x01(\x0b\x32 .cosmos.base.abci.v1beta1.ResultR\x06result\"\"\n\x0cGetTxRequest\x12\x12\n\x04hash\x18\x01 \x01(\tR\x04hash\"}\n\rGetTxResponse\x12%\n\x02tx\x18\x01 \x01(\x0b\x32\x15.cosmos.tx.v1beta1.TxR\x02tx\x12\x45\n\x0btx_response\x18\x02 \x01(\x0b\x32$.cosmos.base.abci.v1beta1.TxResponseR\ntxResponse\"x\n\x16GetBlockWithTxsRequest\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xf0\x01\n\x17GetBlockWithTxsResponse\x12\'\n\x03txs\x18\x01 \x03(\x0b\x32\x15.cosmos.tx.v1beta1.TxR\x03txs\x12\x34\n\x08\x62lock_id\x18\x02 \x01(\x0b\x32\x19.tendermint.types.BlockIDR\x07\x62lockId\x12-\n\x05\x62lock\x18\x03 \x01(\x0b\x32\x17.tendermint.types.BlockR\x05\x62lock\x12G\n\npagination\x18\x04 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\",\n\x0fTxDecodeRequest\x12\x19\n\x08tx_bytes\x18\x01 \x01(\x0cR\x07txBytes\"9\n\x10TxDecodeResponse\x12%\n\x02tx\x18\x01 \x01(\x0b\x32\x15.cosmos.tx.v1beta1.TxR\x02tx\"8\n\x0fTxEncodeRequest\x12%\n\x02tx\x18\x01 \x01(\x0b\x32\x15.cosmos.tx.v1beta1.TxR\x02tx\"-\n\x10TxEncodeResponse\x12\x19\n\x08tx_bytes\x18\x01 \x01(\x0cR\x07txBytes\"5\n\x14TxEncodeAminoRequest\x12\x1d\n\namino_json\x18\x01 \x01(\tR\taminoJson\":\n\x15TxEncodeAminoResponse\x12!\n\x0c\x61mino_binary\x18\x01 \x01(\x0cR\x0b\x61minoBinary\"9\n\x14TxDecodeAminoRequest\x12!\n\x0c\x61mino_binary\x18\x01 \x01(\x0cR\x0b\x61minoBinary\"6\n\x15TxDecodeAminoResponse\x12\x1d\n\namino_json\x18\x01 \x01(\tR\taminoJson*H\n\x07OrderBy\x12\x18\n\x14ORDER_BY_UNSPECIFIED\x10\x00\x12\x10\n\x0cORDER_BY_ASC\x10\x01\x12\x11\n\rORDER_BY_DESC\x10\x02*\x80\x01\n\rBroadcastMode\x12\x1e\n\x1a\x42ROADCAST_MODE_UNSPECIFIED\x10\x00\x12\x1c\n\x14\x42ROADCAST_MODE_BLOCK\x10\x01\x1a\x02\x08\x01\x12\x17\n\x13\x42ROADCAST_MODE_SYNC\x10\x02\x12\x18\n\x14\x42ROADCAST_MODE_ASYNC\x10\x03\x32\xaa\t\n\x07Service\x12{\n\x08Simulate\x12\".cosmos.tx.v1beta1.SimulateRequest\x1a#.cosmos.tx.v1beta1.SimulateResponse\"&\x82\xd3\xe4\x93\x02 \"\x1b/cosmos/tx/v1beta1/simulate:\x01*\x12q\n\x05GetTx\x12\x1f.cosmos.tx.v1beta1.GetTxRequest\x1a .cosmos.tx.v1beta1.GetTxResponse\"%\x82\xd3\xe4\x93\x02\x1f\x12\x1d/cosmos/tx/v1beta1/txs/{hash}\x12\x7f\n\x0b\x42roadcastTx\x12%.cosmos.tx.v1beta1.BroadcastTxRequest\x1a&.cosmos.tx.v1beta1.BroadcastTxResponse\"!\x82\xd3\xe4\x93\x02\x1b\"\x16/cosmos/tx/v1beta1/txs:\x01*\x12|\n\x0bGetTxsEvent\x12%.cosmos.tx.v1beta1.GetTxsEventRequest\x1a&.cosmos.tx.v1beta1.GetTxsEventResponse\"\x1e\x82\xd3\xe4\x93\x02\x18\x12\x16/cosmos/tx/v1beta1/txs\x12\x97\x01\n\x0fGetBlockWithTxs\x12).cosmos.tx.v1beta1.GetBlockWithTxsRequest\x1a*.cosmos.tx.v1beta1.GetBlockWithTxsResponse\"-\x82\xd3\xe4\x93\x02\'\x12%/cosmos/tx/v1beta1/txs/block/{height}\x12y\n\x08TxDecode\x12\".cosmos.tx.v1beta1.TxDecodeRequest\x1a#.cosmos.tx.v1beta1.TxDecodeResponse\"$\x82\xd3\xe4\x93\x02\x1e\"\x19/cosmos/tx/v1beta1/decode:\x01*\x12y\n\x08TxEncode\x12\".cosmos.tx.v1beta1.TxEncodeRequest\x1a#.cosmos.tx.v1beta1.TxEncodeResponse\"$\x82\xd3\xe4\x93\x02\x1e\"\x19/cosmos/tx/v1beta1/encode:\x01*\x12\x8e\x01\n\rTxEncodeAmino\x12\'.cosmos.tx.v1beta1.TxEncodeAminoRequest\x1a(.cosmos.tx.v1beta1.TxEncodeAminoResponse\"*\x82\xd3\xe4\x93\x02$\"\x1f/cosmos/tx/v1beta1/encode/amino:\x01*\x12\x8e\x01\n\rTxDecodeAmino\x12\'.cosmos.tx.v1beta1.TxDecodeAminoRequest\x1a(.cosmos.tx.v1beta1.TxDecodeAminoResponse\"*\x82\xd3\xe4\x93\x02$\"\x1f/cosmos/tx/v1beta1/decode/amino:\x01*B\'Z%github.com/cosmos/cosmos-sdk/types/txb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.tx.v1beta1.service_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z%github.com/cosmos/cosmos-sdk/types/tx'
  _globals['_BROADCASTMODE'].values_by_name["BROADCAST_MODE_BLOCK"]._loaded_options = None
  _globals['_BROADCASTMODE'].values_by_name["BROADCAST_MODE_BLOCK"]._serialized_options = b'\010\001'
  _globals['_GETTXSEVENTREQUEST'].fields_by_name['events']._loaded_options = None
  _globals['_GETTXSEVENTREQUEST'].fields_by_name['events']._serialized_options = b'\030\001'
  _globals['_GETTXSEVENTREQUEST'].fields_by_name['pagination']._loaded_options = None
  _globals['_GETTXSEVENTREQUEST'].fields_by_name['pagination']._serialized_options = b'\030\001'
  _globals['_GETTXSEVENTRESPONSE'].fields_by_name['pagination']._loaded_options = None
  _globals['_GETTXSEVENTRESPONSE'].fields_by_name['pagination']._serialized_options = b'\030\001'
  _globals['_SIMULATEREQUEST'].fields_by_name['tx']._loaded_options = None
  _globals['_SIMULATEREQUEST'].fields_by_name['tx']._serialized_options = b'\030\001'
  _globals['_SERVICE'].methods_by_name['Simulate']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['Simulate']._serialized_options = b'\202\323\344\223\002 \"\033/cosmos/tx/v1beta1/simulate:\001*'
  _globals['_SERVICE'].methods_by_name['GetTx']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['GetTx']._serialized_options = b'\202\323\344\223\002\037\022\035/cosmos/tx/v1beta1/txs/{hash}'
  _globals['_SERVICE'].methods_by_name['BroadcastTx']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['BroadcastTx']._serialized_options = b'\202\323\344\223\002\033\"\026/cosmos/tx/v1beta1/txs:\001*'
  _globals['_SERVICE'].methods_by_name['GetTxsEvent']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['GetTxsEvent']._serialized_options = b'\202\323\344\223\002\030\022\026/cosmos/tx/v1beta1/txs'
  _globals['_SERVICE'].methods_by_name['GetBlockWithTxs']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['GetBlockWithTxs']._serialized_options = b'\202\323\344\223\002\'\022%/cosmos/tx/v1beta1/txs/block/{height}'
  _globals['_SERVICE'].methods_by_name['TxDecode']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['TxDecode']._serialized_options = b'\202\323\344\223\002\036\"\031/cosmos/tx/v1beta1/decode:\001*'
  _globals['_SERVICE'].methods_by_name['TxEncode']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['TxEncode']._serialized_options = b'\202\323\344\223\002\036\"\031/cosmos/tx/v1beta1/encode:\001*'
  _globals['_SERVICE'].methods_by_name['TxEncodeAmino']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['TxEncodeAmino']._serialized_options = b'\202\323\344\223\002$\"\037/cosmos/tx/v1beta1/encode/amino:\001*'
  _globals['_SERVICE'].methods_by_name['TxDecodeAmino']._loaded_options = None
  _globals['_SERVICE'].methods_by_name['TxDecodeAmino']._serialized_options = b'\202\323\344\223\002$\"\037/cosmos/tx/v1beta1/decode/amino:\001*'
  _globals['_ORDERBY']._serialized_start=2131
  _globals['_ORDERBY']._serialized_end=2203
  _globals['_BROADCASTMODE']._serialized_start=2206
  _globals['_BROADCASTMODE']._serialized_end=2334
  _globals['_GETTXSEVENTREQUEST']._serialized_start=254
  _globals['_GETTXSEVENTREQUEST']._serialized_end=497
  _globals['_GETTXSEVENTRESPONSE']._serialized_start=500
  _globals['_GETTXSEVENTRESPONSE']._serialized_end=734
  _globals['_BROADCASTTXREQUEST']._serialized_start=736
  _globals['_BROADCASTTXREQUEST']._serialized_end=837
  _globals['_BROADCASTTXRESPONSE']._serialized_start=839
  _globals['_BROADCASTTXRESPONSE']._serialized_end=931
  _globals['_SIMULATEREQUEST']._serialized_start=933
  _globals['_SIMULATEREQUEST']._serialized_end=1020
  _globals['_SIMULATERESPONSE']._serialized_start=1023
  _globals['_SIMULATERESPONSE']._serialized_end=1161
  _globals['_GETTXREQUEST']._serialized_start=1163
  _globals['_GETTXREQUEST']._serialized_end=1197
  _globals['_GETTXRESPONSE']._serialized_start=1199
  _globals['_GETTXRESPONSE']._serialized_end=1324
  _globals['_GETBLOCKWITHTXSREQUEST']._serialized_start=1326
  _globals['_GETBLOCKWITHTXSREQUEST']._serialized_end=1446
  _globals['_GETBLOCKWITHTXSRESPONSE']._serialized_start=1449
  _globals['_GETBLOCKWITHTXSRESPONSE']._serialized_end=1689
  _globals['_TXDECODEREQUEST']._serialized_start=1691
  _globals['_TXDECODEREQUEST']._serialized_end=1735
  _globals['_TXDECODERESPONSE']._serialized_start=1737
  _globals['_TXDECODERESPONSE']._serialized_end=1794
  _globals['_TXENCODEREQUEST']._serialized_start=1796
  _globals['_TXENCODEREQUEST']._serialized_end=1852
  _globals['_TXENCODERESPONSE']._serialized_start=1854
  _globals['_TXENCODERESPONSE']._serialized_end=1899
  _globals['_TXENCODEAMINOREQUEST']._serialized_start=1901
  _globals['_TXENCODEAMINOREQUEST']._serialized_end=1954
  _globals['_TXENCODEAMINORESPONSE']._serialized_start=1956
  _globals['_TXENCODEAMINORESPONSE']._serialized_end=2014
  _globals['_TXDECODEAMINOREQUEST']._serialized_start=2016
  _globals['_TXDECODEAMINOREQUEST']._serialized_end=2073
  _globals['_TXDECODEAMINORESPONSE']._serialized_start=2075
  _globals['_TXDECODEAMINORESPONSE']._serialized_end=2129
  _globals['_SERVICE']._serialized_start=2337
  _globals['_SERVICE']._serialized_end=3531
# @@protoc_insertion_point(module_scope)
