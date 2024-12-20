# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmwasm/wasm/v1/types.proto
# Protobuf Python Version: 5.29.1
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
    1,
    '',
    'cosmwasm/wasm/v1/types.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1c\x63osmwasm/wasm/v1/types.proto\x12\x10\x63osmwasm.wasm.v1\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a\x11\x61mino/amino.proto\"]\n\x0f\x41\x63\x63\x65ssTypeParam\x12\x44\n\x05value\x18\x01 \x01(\x0e\x32\x1c.cosmwasm.wasm.v1.AccessTypeB\x10\xf2\xde\x1f\x0cyaml:\"value\"R\x05value:\x04\x98\xa0\x1f\x01\"\xa7\x01\n\x0c\x41\x63\x63\x65ssConfig\x12S\n\npermission\x18\x01 \x01(\x0e\x32\x1c.cosmwasm.wasm.v1.AccessTypeB\x15\xf2\xde\x1f\x11yaml:\"permission\"R\npermission\x12\x36\n\taddresses\x18\x03 \x03(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\taddresses:\x04\x98\xa0\x1f\x01J\x04\x08\x02\x10\x03\"\x94\x02\n\x06Params\x12t\n\x12\x63ode_upload_access\x18\x01 \x01(\x0b\x32\x1e.cosmwasm.wasm.v1.AccessConfigB&\xc8\xde\x1f\x00\xf2\xde\x1f\x19yaml:\"code_upload_access\"\xa8\xe7\xb0*\x01R\x10\x63odeUploadAccess\x12\x8d\x01\n\x1einstantiate_default_permission\x18\x02 \x01(\x0e\x32\x1c.cosmwasm.wasm.v1.AccessTypeB)\xf2\xde\x1f%yaml:\"instantiate_default_permission\"R\x1cinstantiateDefaultPermission:\x04\x98\xa0\x1f\x00\"\xc1\x01\n\x08\x43odeInfo\x12\x1b\n\tcode_hash\x18\x01 \x01(\x0cR\x08\x63odeHash\x12\x32\n\x07\x63reator\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x63reator\x12X\n\x12instantiate_config\x18\x05 \x01(\x0b\x32\x1e.cosmwasm.wasm.v1.AccessConfigB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x11instantiateConfigJ\x04\x08\x03\x10\x04J\x04\x08\x04\x10\x05\"\x82\x03\n\x0c\x43ontractInfo\x12#\n\x07\x63ode_id\x18\x01 \x01(\x04\x42\n\xe2\xde\x1f\x06\x43odeIDR\x06\x63odeId\x12\x32\n\x07\x63reator\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x63reator\x12.\n\x05\x61\x64min\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05\x61\x64min\x12\x14\n\x05label\x18\x04 \x01(\tR\x05label\x12>\n\x07\x63reated\x18\x05 \x01(\x0b\x32$.cosmwasm.wasm.v1.AbsoluteTxPositionR\x07\x63reated\x12-\n\x0bibc_port_id\x18\x06 \x01(\tB\r\xe2\xde\x1f\tIBCPortIDR\tibcPortId\x12^\n\textension\x18\x07 \x01(\x0b\x32\x14.google.protobuf.AnyB*\xca\xb4-&cosmwasm.wasm.v1.ContractInfoExtensionR\textension:\x04\xe8\xa0\x1f\x01\"\x8b\x02\n\x18\x43ontractCodeHistoryEntry\x12P\n\toperation\x18\x01 \x01(\x0e\x32\x32.cosmwasm.wasm.v1.ContractCodeHistoryOperationTypeR\toperation\x12#\n\x07\x63ode_id\x18\x02 \x01(\x04\x42\n\xe2\xde\x1f\x06\x43odeIDR\x06\x63odeId\x12>\n\x07updated\x18\x03 \x01(\x0b\x32$.cosmwasm.wasm.v1.AbsoluteTxPositionR\x07updated\x12\x38\n\x03msg\x18\x04 \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg\"R\n\x12\x41\x62soluteTxPosition\x12!\n\x0c\x62lock_height\x18\x01 \x01(\x04R\x0b\x62lockHeight\x12\x19\n\x08tx_index\x18\x02 \x01(\x04R\x07txIndex\"e\n\x05Model\x12\x46\n\x03key\x18\x01 \x01(\x0c\x42\x34\xfa\xde\x1f\x30github.com/cometbft/cometbft/libs/bytes.HexBytesR\x03key\x12\x14\n\x05value\x18\x02 \x01(\x0cR\x05value*\xf6\x01\n\nAccessType\x12\x36\n\x17\x41\x43\x43\x45SS_TYPE_UNSPECIFIED\x10\x00\x1a\x19\x8a\x9d \x15\x41\x63\x63\x65ssTypeUnspecified\x12,\n\x12\x41\x43\x43\x45SS_TYPE_NOBODY\x10\x01\x1a\x14\x8a\x9d \x10\x41\x63\x63\x65ssTypeNobody\x12\x32\n\x15\x41\x43\x43\x45SS_TYPE_EVERYBODY\x10\x03\x1a\x17\x8a\x9d \x13\x41\x63\x63\x65ssTypeEverybody\x12>\n\x1c\x41\x43\x43\x45SS_TYPE_ANY_OF_ADDRESSES\x10\x04\x1a\x1c\x8a\x9d \x18\x41\x63\x63\x65ssTypeAnyOfAddresses\x1a\x08\x88\xa3\x1e\x00\xa8\xa4\x1e\x00\"\x04\x08\x02\x10\x02*\xa6\x03\n ContractCodeHistoryOperationType\x12\x65\n0CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED\x10\x00\x1a/\x8a\x9d +ContractCodeHistoryOperationTypeUnspecified\x12W\n)CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT\x10\x01\x1a(\x8a\x9d $ContractCodeHistoryOperationTypeInit\x12]\n,CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE\x10\x02\x1a+\x8a\x9d \'ContractCodeHistoryOperationTypeMigrate\x12]\n,CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS\x10\x03\x1a+\x8a\x9d \'ContractCodeHistoryOperationTypeGenesis\x1a\x04\x88\xa3\x1e\x00\x42\x30Z&github.com/CosmWasm/wasmd/x/wasm/types\xc8\xe1\x1e\x00\xa8\xe2\x1e\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmwasm.wasm.v1.types_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z&github.com/CosmWasm/wasmd/x/wasm/types\310\341\036\000\250\342\036\001'
  _globals['_ACCESSTYPE']._loaded_options = None
  _globals['_ACCESSTYPE']._serialized_options = b'\210\243\036\000\250\244\036\000'
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_UNSPECIFIED"]._loaded_options = None
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_UNSPECIFIED"]._serialized_options = b'\212\235 \025AccessTypeUnspecified'
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_NOBODY"]._loaded_options = None
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_NOBODY"]._serialized_options = b'\212\235 \020AccessTypeNobody'
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_EVERYBODY"]._loaded_options = None
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_EVERYBODY"]._serialized_options = b'\212\235 \023AccessTypeEverybody'
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_ANY_OF_ADDRESSES"]._loaded_options = None
  _globals['_ACCESSTYPE'].values_by_name["ACCESS_TYPE_ANY_OF_ADDRESSES"]._serialized_options = b'\212\235 \030AccessTypeAnyOfAddresses'
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE']._loaded_options = None
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE']._serialized_options = b'\210\243\036\000'
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED"]._loaded_options = None
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED"]._serialized_options = b'\212\235 +ContractCodeHistoryOperationTypeUnspecified'
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT"]._loaded_options = None
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT"]._serialized_options = b'\212\235 $ContractCodeHistoryOperationTypeInit'
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE"]._loaded_options = None
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE"]._serialized_options = b'\212\235 \'ContractCodeHistoryOperationTypeMigrate'
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS"]._loaded_options = None
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE'].values_by_name["CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS"]._serialized_options = b'\212\235 \'ContractCodeHistoryOperationTypeGenesis'
  _globals['_ACCESSTYPEPARAM'].fields_by_name['value']._loaded_options = None
  _globals['_ACCESSTYPEPARAM'].fields_by_name['value']._serialized_options = b'\362\336\037\014yaml:\"value\"'
  _globals['_ACCESSTYPEPARAM']._loaded_options = None
  _globals['_ACCESSTYPEPARAM']._serialized_options = b'\230\240\037\001'
  _globals['_ACCESSCONFIG'].fields_by_name['permission']._loaded_options = None
  _globals['_ACCESSCONFIG'].fields_by_name['permission']._serialized_options = b'\362\336\037\021yaml:\"permission\"'
  _globals['_ACCESSCONFIG'].fields_by_name['addresses']._loaded_options = None
  _globals['_ACCESSCONFIG'].fields_by_name['addresses']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_ACCESSCONFIG']._loaded_options = None
  _globals['_ACCESSCONFIG']._serialized_options = b'\230\240\037\001'
  _globals['_PARAMS'].fields_by_name['code_upload_access']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['code_upload_access']._serialized_options = b'\310\336\037\000\362\336\037\031yaml:\"code_upload_access\"\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['instantiate_default_permission']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['instantiate_default_permission']._serialized_options = b'\362\336\037%yaml:\"instantiate_default_permission\"'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\230\240\037\000'
  _globals['_CODEINFO'].fields_by_name['creator']._loaded_options = None
  _globals['_CODEINFO'].fields_by_name['creator']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_CODEINFO'].fields_by_name['instantiate_config']._loaded_options = None
  _globals['_CODEINFO'].fields_by_name['instantiate_config']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_CONTRACTINFO'].fields_by_name['code_id']._loaded_options = None
  _globals['_CONTRACTINFO'].fields_by_name['code_id']._serialized_options = b'\342\336\037\006CodeID'
  _globals['_CONTRACTINFO'].fields_by_name['creator']._loaded_options = None
  _globals['_CONTRACTINFO'].fields_by_name['creator']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_CONTRACTINFO'].fields_by_name['admin']._loaded_options = None
  _globals['_CONTRACTINFO'].fields_by_name['admin']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_CONTRACTINFO'].fields_by_name['ibc_port_id']._loaded_options = None
  _globals['_CONTRACTINFO'].fields_by_name['ibc_port_id']._serialized_options = b'\342\336\037\tIBCPortID'
  _globals['_CONTRACTINFO'].fields_by_name['extension']._loaded_options = None
  _globals['_CONTRACTINFO'].fields_by_name['extension']._serialized_options = b'\312\264-&cosmwasm.wasm.v1.ContractInfoExtension'
  _globals['_CONTRACTINFO']._loaded_options = None
  _globals['_CONTRACTINFO']._serialized_options = b'\350\240\037\001'
  _globals['_CONTRACTCODEHISTORYENTRY'].fields_by_name['code_id']._loaded_options = None
  _globals['_CONTRACTCODEHISTORYENTRY'].fields_by_name['code_id']._serialized_options = b'\342\336\037\006CodeID'
  _globals['_CONTRACTCODEHISTORYENTRY'].fields_by_name['msg']._loaded_options = None
  _globals['_CONTRACTCODEHISTORYENTRY'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_MODEL'].fields_by_name['key']._loaded_options = None
  _globals['_MODEL'].fields_by_name['key']._serialized_options = b'\372\336\0370github.com/cometbft/cometbft/libs/bytes.HexBytes'
  _globals['_ACCESSTYPE']._serialized_start=1732
  _globals['_ACCESSTYPE']._serialized_end=1978
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE']._serialized_start=1981
  _globals['_CONTRACTCODEHISTORYOPERATIONTYPE']._serialized_end=2403
  _globals['_ACCESSTYPEPARAM']._serialized_start=145
  _globals['_ACCESSTYPEPARAM']._serialized_end=238
  _globals['_ACCESSCONFIG']._serialized_start=241
  _globals['_ACCESSCONFIG']._serialized_end=408
  _globals['_PARAMS']._serialized_start=411
  _globals['_PARAMS']._serialized_end=687
  _globals['_CODEINFO']._serialized_start=690
  _globals['_CODEINFO']._serialized_end=883
  _globals['_CONTRACTINFO']._serialized_start=886
  _globals['_CONTRACTINFO']._serialized_end=1272
  _globals['_CONTRACTCODEHISTORYENTRY']._serialized_start=1275
  _globals['_CONTRACTCODEHISTORYENTRY']._serialized_end=1542
  _globals['_ABSOLUTETXPOSITION']._serialized_start=1544
  _globals['_ABSOLUTETXPOSITION']._serialized_end=1626
  _globals['_MODEL']._serialized_start=1628
  _globals['_MODEL']._serialized_end=1729
# @@protoc_insertion_point(module_scope)
