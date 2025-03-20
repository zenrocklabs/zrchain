# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/base/abci/v1beta1/abci.proto
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
    'cosmos/base/abci/v1beta1/abci.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from tendermint.abci import types_pb2 as tendermint_dot_abci_dot_types__pb2
from tendermint.types import block_pb2 as tendermint_dot_types_dot_block__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n#cosmos/base/abci/v1beta1/abci.proto\x12\x18\x63osmos.base.abci.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1btendermint/abci/types.proto\x1a\x1ctendermint/types/block.proto\x1a\x19google/protobuf/any.proto\"\xcc\x03\n\nTxResponse\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\x12\"\n\x06txhash\x18\x02 \x01(\tB\n\xe2\xde\x1f\x06TxHashR\x06txhash\x12\x1c\n\tcodespace\x18\x03 \x01(\tR\tcodespace\x12\x12\n\x04\x63ode\x18\x04 \x01(\rR\x04\x63ode\x12\x12\n\x04\x64\x61ta\x18\x05 \x01(\tR\x04\x64\x61ta\x12\x17\n\x07raw_log\x18\x06 \x01(\tR\x06rawLog\x12U\n\x04logs\x18\x07 \x03(\x0b\x32(.cosmos.base.abci.v1beta1.ABCIMessageLogB\x17\xc8\xde\x1f\x00\xaa\xdf\x1f\x0f\x41\x42\x43IMessageLogsR\x04logs\x12\x12\n\x04info\x18\x08 \x01(\tR\x04info\x12\x1d\n\ngas_wanted\x18\t \x01(\x03R\tgasWanted\x12\x19\n\x08gas_used\x18\n \x01(\x03R\x07gasUsed\x12$\n\x02tx\x18\x0b \x01(\x0b\x32\x14.google.protobuf.AnyR\x02tx\x12\x1c\n\ttimestamp\x18\x0c \x01(\tR\ttimestamp\x12\x34\n\x06\x65vents\x18\r \x03(\x0b\x32\x16.tendermint.abci.EventB\x04\xc8\xde\x1f\x00R\x06\x65vents:\x04\x88\xa0\x1f\x00\"\xa9\x01\n\x0e\x41\x42\x43IMessageLog\x12*\n\tmsg_index\x18\x01 \x01(\rB\r\xea\xde\x1f\tmsg_indexR\x08msgIndex\x12\x10\n\x03log\x18\x02 \x01(\tR\x03log\x12S\n\x06\x65vents\x18\x03 \x03(\x0b\x32%.cosmos.base.abci.v1beta1.StringEventB\x14\xc8\xde\x1f\x00\xaa\xdf\x1f\x0cStringEventsR\x06\x65vents:\x04\x80\xdc \x01\"r\n\x0bStringEvent\x12\x12\n\x04type\x18\x01 \x01(\tR\x04type\x12I\n\nattributes\x18\x02 \x03(\x0b\x32#.cosmos.base.abci.v1beta1.AttributeB\x04\xc8\xde\x1f\x00R\nattributes:\x04\x80\xdc \x01\"3\n\tAttribute\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value\"C\n\x07GasInfo\x12\x1d\n\ngas_wanted\x18\x01 \x01(\x04R\tgasWanted\x12\x19\n\x08gas_used\x18\x02 \x01(\x04R\x07gasUsed\"\xa9\x01\n\x06Result\x12\x16\n\x04\x64\x61ta\x18\x01 \x01(\x0c\x42\x02\x18\x01R\x04\x64\x61ta\x12\x10\n\x03log\x18\x02 \x01(\tR\x03log\x12\x34\n\x06\x65vents\x18\x03 \x03(\x0b\x32\x16.tendermint.abci.EventB\x04\xc8\xde\x1f\x00R\x06\x65vents\x12\x39\n\rmsg_responses\x18\x04 \x03(\x0b\x32\x14.google.protobuf.AnyR\x0cmsgResponses:\x04\x88\xa0\x1f\x00\"\x96\x01\n\x12SimulationResponse\x12\x46\n\x08gas_info\x18\x01 \x01(\x0b\x32!.cosmos.base.abci.v1beta1.GasInfoB\x08\xc8\xde\x1f\x00\xd0\xde\x1f\x01R\x07gasInfo\x12\x38\n\x06result\x18\x02 \x01(\x0b\x32 .cosmos.base.abci.v1beta1.ResultR\x06result\"@\n\x07MsgData\x12\x19\n\x08msg_type\x18\x01 \x01(\tR\x07msgType\x12\x12\n\x04\x64\x61ta\x18\x02 \x01(\x0cR\x04\x64\x61ta:\x06\x18\x01\x80\xdc \x01\"\x87\x01\n\tTxMsgData\x12\x39\n\x04\x64\x61ta\x18\x01 \x03(\x0b\x32!.cosmos.base.abci.v1beta1.MsgDataB\x02\x18\x01R\x04\x64\x61ta\x12\x39\n\rmsg_responses\x18\x02 \x03(\x0b\x32\x14.google.protobuf.AnyR\x0cmsgResponses:\x04\x80\xdc \x01\"\xdc\x01\n\x0fSearchTxsResult\x12\x1f\n\x0btotal_count\x18\x01 \x01(\x04R\ntotalCount\x12\x14\n\x05\x63ount\x18\x02 \x01(\x04R\x05\x63ount\x12\x1f\n\x0bpage_number\x18\x03 \x01(\x04R\npageNumber\x12\x1d\n\npage_total\x18\x04 \x01(\x04R\tpageTotal\x12\x14\n\x05limit\x18\x05 \x01(\x04R\x05limit\x12\x36\n\x03txs\x18\x06 \x03(\x0b\x32$.cosmos.base.abci.v1beta1.TxResponseR\x03txs:\x04\x80\xdc \x01\"\xd8\x01\n\x12SearchBlocksResult\x12\x1f\n\x0btotal_count\x18\x01 \x01(\x03R\ntotalCount\x12\x14\n\x05\x63ount\x18\x02 \x01(\x03R\x05\x63ount\x12\x1f\n\x0bpage_number\x18\x03 \x01(\x03R\npageNumber\x12\x1d\n\npage_total\x18\x04 \x01(\x03R\tpageTotal\x12\x14\n\x05limit\x18\x05 \x01(\x03R\x05limit\x12/\n\x06\x62locks\x18\x06 \x03(\x0b\x32\x17.tendermint.types.BlockR\x06\x62locks:\x04\x80\xdc \x01\x42(Z\"github.com/cosmos/cosmos-sdk/types\xd8\xe1\x1e\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.base.abci.v1beta1.abci_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\"github.com/cosmos/cosmos-sdk/types\330\341\036\000'
  _globals['_TXRESPONSE'].fields_by_name['txhash']._loaded_options = None
  _globals['_TXRESPONSE'].fields_by_name['txhash']._serialized_options = b'\342\336\037\006TxHash'
  _globals['_TXRESPONSE'].fields_by_name['logs']._loaded_options = None
  _globals['_TXRESPONSE'].fields_by_name['logs']._serialized_options = b'\310\336\037\000\252\337\037\017ABCIMessageLogs'
  _globals['_TXRESPONSE'].fields_by_name['events']._loaded_options = None
  _globals['_TXRESPONSE'].fields_by_name['events']._serialized_options = b'\310\336\037\000'
  _globals['_TXRESPONSE']._loaded_options = None
  _globals['_TXRESPONSE']._serialized_options = b'\210\240\037\000'
  _globals['_ABCIMESSAGELOG'].fields_by_name['msg_index']._loaded_options = None
  _globals['_ABCIMESSAGELOG'].fields_by_name['msg_index']._serialized_options = b'\352\336\037\tmsg_index'
  _globals['_ABCIMESSAGELOG'].fields_by_name['events']._loaded_options = None
  _globals['_ABCIMESSAGELOG'].fields_by_name['events']._serialized_options = b'\310\336\037\000\252\337\037\014StringEvents'
  _globals['_ABCIMESSAGELOG']._loaded_options = None
  _globals['_ABCIMESSAGELOG']._serialized_options = b'\200\334 \001'
  _globals['_STRINGEVENT'].fields_by_name['attributes']._loaded_options = None
  _globals['_STRINGEVENT'].fields_by_name['attributes']._serialized_options = b'\310\336\037\000'
  _globals['_STRINGEVENT']._loaded_options = None
  _globals['_STRINGEVENT']._serialized_options = b'\200\334 \001'
  _globals['_RESULT'].fields_by_name['data']._loaded_options = None
  _globals['_RESULT'].fields_by_name['data']._serialized_options = b'\030\001'
  _globals['_RESULT'].fields_by_name['events']._loaded_options = None
  _globals['_RESULT'].fields_by_name['events']._serialized_options = b'\310\336\037\000'
  _globals['_RESULT']._loaded_options = None
  _globals['_RESULT']._serialized_options = b'\210\240\037\000'
  _globals['_SIMULATIONRESPONSE'].fields_by_name['gas_info']._loaded_options = None
  _globals['_SIMULATIONRESPONSE'].fields_by_name['gas_info']._serialized_options = b'\310\336\037\000\320\336\037\001'
  _globals['_MSGDATA']._loaded_options = None
  _globals['_MSGDATA']._serialized_options = b'\030\001\200\334 \001'
  _globals['_TXMSGDATA'].fields_by_name['data']._loaded_options = None
  _globals['_TXMSGDATA'].fields_by_name['data']._serialized_options = b'\030\001'
  _globals['_TXMSGDATA']._loaded_options = None
  _globals['_TXMSGDATA']._serialized_options = b'\200\334 \001'
  _globals['_SEARCHTXSRESULT']._loaded_options = None
  _globals['_SEARCHTXSRESULT']._serialized_options = b'\200\334 \001'
  _globals['_SEARCHBLOCKSRESULT']._loaded_options = None
  _globals['_SEARCHBLOCKSRESULT']._serialized_options = b'\200\334 \001'
  _globals['_TXRESPONSE']._serialized_start=174
  _globals['_TXRESPONSE']._serialized_end=634
  _globals['_ABCIMESSAGELOG']._serialized_start=637
  _globals['_ABCIMESSAGELOG']._serialized_end=806
  _globals['_STRINGEVENT']._serialized_start=808
  _globals['_STRINGEVENT']._serialized_end=922
  _globals['_ATTRIBUTE']._serialized_start=924
  _globals['_ATTRIBUTE']._serialized_end=975
  _globals['_GASINFO']._serialized_start=977
  _globals['_GASINFO']._serialized_end=1044
  _globals['_RESULT']._serialized_start=1047
  _globals['_RESULT']._serialized_end=1216
  _globals['_SIMULATIONRESPONSE']._serialized_start=1219
  _globals['_SIMULATIONRESPONSE']._serialized_end=1369
  _globals['_MSGDATA']._serialized_start=1371
  _globals['_MSGDATA']._serialized_end=1435
  _globals['_TXMSGDATA']._serialized_start=1438
  _globals['_TXMSGDATA']._serialized_end=1573
  _globals['_SEARCHTXSRESULT']._serialized_start=1576
  _globals['_SEARCHTXSRESULT']._serialized_end=1796
  _globals['_SEARCHBLOCKSRESULT']._serialized_start=1799
  _globals['_SEARCHBLOCKSRESULT']._serialized_end=2015
# @@protoc_insertion_point(module_scope)
