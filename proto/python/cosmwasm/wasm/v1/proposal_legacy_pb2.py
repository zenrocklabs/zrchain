# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmwasm/wasm/v1/proposal_legacy.proto
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
    'cosmwasm/wasm/v1/proposal_legacy.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmwasm.wasm.v1 import types_pb2 as cosmwasm_dot_wasm_dot_v1_dot_types__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n&cosmwasm/wasm/v1/proposal_legacy.proto\x12\x10\x63osmwasm.wasm.v1\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x1c\x63osmwasm/wasm/v1/types.proto\x1a\x11\x61mino/amino.proto\"\xc2\x03\n\x11StoreCodeProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12/\n\x06run_as\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05runAs\x12\x36\n\x0ewasm_byte_code\x18\x04 \x01(\x0c\x42\x10\xe2\xde\x1f\x0cWASMByteCodeR\x0cwasmByteCode\x12U\n\x16instantiate_permission\x18\x07 \x01(\x0b\x32\x1e.cosmwasm.wasm.v1.AccessConfigR\x15instantiatePermission\x12\x1d\n\nunpin_code\x18\x08 \x01(\x08R\tunpinCode\x12\x16\n\x06source\x18\t \x01(\tR\x06source\x12\x18\n\x07\x62uilder\x18\n \x01(\tR\x07\x62uilder\x12\x1b\n\tcode_hash\x18\x0b \x01(\x0cR\x08\x63odeHash:;\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x16wasm/StoreCodeProposalJ\x04\x08\x05\x10\x06J\x04\x08\x06\x10\x07\"\xeb\x03\n\x1bInstantiateContractProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12/\n\x06run_as\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05runAs\x12.\n\x05\x61\x64min\x18\x04 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05\x61\x64min\x12#\n\x07\x63ode_id\x18\x05 \x01(\x04\x42\n\xe2\xde\x1f\x06\x43odeIDR\x06\x63odeId\x12\x14\n\x05label\x18\x06 \x01(\tR\x05label\x12\x38\n\x03msg\x18\x07 \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg\x12w\n\x05\x66unds\x18\x08 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x05\x66unds:E\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0* wasm/InstantiateContractProposal\"\x9a\x04\n\x1cInstantiateContract2Proposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12/\n\x06run_as\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05runAs\x12.\n\x05\x61\x64min\x18\x04 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05\x61\x64min\x12#\n\x07\x63ode_id\x18\x05 \x01(\x04\x42\n\xe2\xde\x1f\x06\x43odeIDR\x06\x63odeId\x12\x14\n\x05label\x18\x06 \x01(\tR\x05label\x12\x38\n\x03msg\x18\x07 \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg\x12w\n\x05\x66unds\x18\x08 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x05\x66unds\x12\x12\n\x04salt\x18\t \x01(\x0cR\x04salt\x12\x17\n\x07\x66ix_msg\x18\n \x01(\x08R\x06\x66ixMsg:F\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*!wasm/InstantiateContract2Proposal\"\xa9\x02\n\x17MigrateContractProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x34\n\x08\x63ontract\x18\x04 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08\x63ontract\x12#\n\x07\x63ode_id\x18\x05 \x01(\x04\x42\n\xe2\xde\x1f\x06\x43odeIDR\x06\x63odeId\x12\x38\n\x03msg\x18\x06 \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg:A\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x1cwasm/MigrateContractProposal\"\xfe\x01\n\x14SudoContractProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x34\n\x08\x63ontract\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08\x63ontract\x12\x38\n\x03msg\x18\x04 \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg:>\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x19wasm/SudoContractProposal\"\xae\x03\n\x17\x45xecuteContractProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12/\n\x06run_as\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05runAs\x12\x34\n\x08\x63ontract\x18\x04 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08\x63ontract\x12\x38\n\x03msg\x18\x05 \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg\x12w\n\x05\x66unds\x18\x06 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x05\x66unds:A\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x1cwasm/ExecuteContractProposal\"\x8d\x02\n\x13UpdateAdminProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12I\n\tnew_admin\x18\x03 \x01(\tB,\xf2\xde\x1f\x10yaml:\"new_admin\"\xd2\xb4-\x14\x63osmos.AddressStringR\x08newAdmin\x12\x34\n\x08\x63ontract\x18\x04 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08\x63ontract:=\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x18wasm/UpdateAdminProposal\"\xc0\x01\n\x12\x43learAdminProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x34\n\x08\x63ontract\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08\x63ontract:<\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x17wasm/ClearAdminProposal\"\xc1\x01\n\x10PinCodesProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x39\n\x08\x63ode_ids\x18\x03 \x03(\x04\x42\x1e\xe2\xde\x1f\x07\x43odeIDs\xf2\xde\x1f\x0fyaml:\"code_ids\"R\x07\x63odeIds::\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x15wasm/PinCodesProposal\"\xc5\x01\n\x12UnpinCodesProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x39\n\x08\x63ode_ids\x18\x03 \x03(\x04\x42\x1e\xe2\xde\x1f\x07\x43odeIDs\xf2\xde\x1f\x0fyaml:\"code_ids\"R\x07\x63odeIds:<\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*\x17wasm/UnpinCodesProposal\"\x9b\x01\n\x12\x41\x63\x63\x65ssConfigUpdate\x12#\n\x07\x63ode_id\x18\x01 \x01(\x04\x42\n\xe2\xde\x1f\x06\x43odeIDR\x06\x63odeId\x12`\n\x16instantiate_permission\x18\x02 \x01(\x0b\x32\x1e.cosmwasm.wasm.v1.AccessConfigB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x15instantiatePermission\"\xb3\x02\n\x1fUpdateInstantiateConfigProposal\x12&\n\x05title\x18\x01 \x01(\tB\x10\xf2\xde\x1f\x0cyaml:\"title\"R\x05title\x12\x38\n\x0b\x64\x65scription\x18\x02 \x01(\tB\x16\xf2\xde\x1f\x12yaml:\"description\"R\x0b\x64\x65scription\x12\x63\n\x15\x61\x63\x63\x65ss_config_updates\x18\x03 \x03(\x0b\x32$.cosmwasm.wasm.v1.AccessConfigUpdateB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x13\x61\x63\x63\x65ssConfigUpdates:I\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*$wasm/UpdateInstantiateConfigProposal\"\xb9\x05\n#StoreAndInstantiateContractProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12/\n\x06run_as\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05runAs\x12\x36\n\x0ewasm_byte_code\x18\x04 \x01(\x0c\x42\x10\xe2\xde\x1f\x0cWASMByteCodeR\x0cwasmByteCode\x12U\n\x16instantiate_permission\x18\x05 \x01(\x0b\x32\x1e.cosmwasm.wasm.v1.AccessConfigR\x15instantiatePermission\x12\x1d\n\nunpin_code\x18\x06 \x01(\x08R\tunpinCode\x12\x14\n\x05\x61\x64min\x18\x07 \x01(\tR\x05\x61\x64min\x12\x14\n\x05label\x18\x08 \x01(\tR\x05label\x12\x38\n\x03msg\x18\t \x01(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x03msg\x12w\n\x05\x66unds\x18\n \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x05\x66unds\x12\x16\n\x06source\x18\x0b \x01(\tR\x06source\x12\x18\n\x07\x62uilder\x18\x0c \x01(\tR\x07\x62uilder\x12\x1b\n\tcode_hash\x18\r \x01(\x0cR\x08\x63odeHash:M\x18\x01\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\x8a\xe7\xb0*(wasm/StoreAndInstantiateContractProposalB4Z&github.com/CosmWasm/wasmd/x/wasm/types\xc8\xe1\x1e\x00\xd8\xe1\x1e\x00\xa8\xe2\x1e\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmwasm.wasm.v1.proposal_legacy_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z&github.com/CosmWasm/wasmd/x/wasm/types\310\341\036\000\330\341\036\000\250\342\036\001'
  _globals['_STORECODEPROPOSAL'].fields_by_name['run_as']._loaded_options = None
  _globals['_STORECODEPROPOSAL'].fields_by_name['run_as']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_STORECODEPROPOSAL'].fields_by_name['wasm_byte_code']._loaded_options = None
  _globals['_STORECODEPROPOSAL'].fields_by_name['wasm_byte_code']._serialized_options = b'\342\336\037\014WASMByteCode'
  _globals['_STORECODEPROPOSAL']._loaded_options = None
  _globals['_STORECODEPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\026wasm/StoreCodeProposal'
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['run_as']._loaded_options = None
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['run_as']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['admin']._loaded_options = None
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['admin']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['code_id']._loaded_options = None
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['code_id']._serialized_options = b'\342\336\037\006CodeID'
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['msg']._loaded_options = None
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['funds']._loaded_options = None
  _globals['_INSTANTIATECONTRACTPROPOSAL'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_INSTANTIATECONTRACTPROPOSAL']._loaded_options = None
  _globals['_INSTANTIATECONTRACTPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260* wasm/InstantiateContractProposal'
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['run_as']._loaded_options = None
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['run_as']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['admin']._loaded_options = None
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['admin']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['code_id']._loaded_options = None
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['code_id']._serialized_options = b'\342\336\037\006CodeID'
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['msg']._loaded_options = None
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['funds']._loaded_options = None
  _globals['_INSTANTIATECONTRACT2PROPOSAL'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_INSTANTIATECONTRACT2PROPOSAL']._loaded_options = None
  _globals['_INSTANTIATECONTRACT2PROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*!wasm/InstantiateContract2Proposal'
  _globals['_MIGRATECONTRACTPROPOSAL'].fields_by_name['contract']._loaded_options = None
  _globals['_MIGRATECONTRACTPROPOSAL'].fields_by_name['contract']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MIGRATECONTRACTPROPOSAL'].fields_by_name['code_id']._loaded_options = None
  _globals['_MIGRATECONTRACTPROPOSAL'].fields_by_name['code_id']._serialized_options = b'\342\336\037\006CodeID'
  _globals['_MIGRATECONTRACTPROPOSAL'].fields_by_name['msg']._loaded_options = None
  _globals['_MIGRATECONTRACTPROPOSAL'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_MIGRATECONTRACTPROPOSAL']._loaded_options = None
  _globals['_MIGRATECONTRACTPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\034wasm/MigrateContractProposal'
  _globals['_SUDOCONTRACTPROPOSAL'].fields_by_name['contract']._loaded_options = None
  _globals['_SUDOCONTRACTPROPOSAL'].fields_by_name['contract']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_SUDOCONTRACTPROPOSAL'].fields_by_name['msg']._loaded_options = None
  _globals['_SUDOCONTRACTPROPOSAL'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_SUDOCONTRACTPROPOSAL']._loaded_options = None
  _globals['_SUDOCONTRACTPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\031wasm/SudoContractProposal'
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['run_as']._loaded_options = None
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['run_as']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['contract']._loaded_options = None
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['contract']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['msg']._loaded_options = None
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['funds']._loaded_options = None
  _globals['_EXECUTECONTRACTPROPOSAL'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_EXECUTECONTRACTPROPOSAL']._loaded_options = None
  _globals['_EXECUTECONTRACTPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\034wasm/ExecuteContractProposal'
  _globals['_UPDATEADMINPROPOSAL'].fields_by_name['new_admin']._loaded_options = None
  _globals['_UPDATEADMINPROPOSAL'].fields_by_name['new_admin']._serialized_options = b'\362\336\037\020yaml:\"new_admin\"\322\264-\024cosmos.AddressString'
  _globals['_UPDATEADMINPROPOSAL'].fields_by_name['contract']._loaded_options = None
  _globals['_UPDATEADMINPROPOSAL'].fields_by_name['contract']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_UPDATEADMINPROPOSAL']._loaded_options = None
  _globals['_UPDATEADMINPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\030wasm/UpdateAdminProposal'
  _globals['_CLEARADMINPROPOSAL'].fields_by_name['contract']._loaded_options = None
  _globals['_CLEARADMINPROPOSAL'].fields_by_name['contract']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_CLEARADMINPROPOSAL']._loaded_options = None
  _globals['_CLEARADMINPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\027wasm/ClearAdminProposal'
  _globals['_PINCODESPROPOSAL'].fields_by_name['code_ids']._loaded_options = None
  _globals['_PINCODESPROPOSAL'].fields_by_name['code_ids']._serialized_options = b'\342\336\037\007CodeIDs\362\336\037\017yaml:\"code_ids\"'
  _globals['_PINCODESPROPOSAL']._loaded_options = None
  _globals['_PINCODESPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\025wasm/PinCodesProposal'
  _globals['_UNPINCODESPROPOSAL'].fields_by_name['code_ids']._loaded_options = None
  _globals['_UNPINCODESPROPOSAL'].fields_by_name['code_ids']._serialized_options = b'\342\336\037\007CodeIDs\362\336\037\017yaml:\"code_ids\"'
  _globals['_UNPINCODESPROPOSAL']._loaded_options = None
  _globals['_UNPINCODESPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*\027wasm/UnpinCodesProposal'
  _globals['_ACCESSCONFIGUPDATE'].fields_by_name['code_id']._loaded_options = None
  _globals['_ACCESSCONFIGUPDATE'].fields_by_name['code_id']._serialized_options = b'\342\336\037\006CodeID'
  _globals['_ACCESSCONFIGUPDATE'].fields_by_name['instantiate_permission']._loaded_options = None
  _globals['_ACCESSCONFIGUPDATE'].fields_by_name['instantiate_permission']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL'].fields_by_name['title']._loaded_options = None
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL'].fields_by_name['title']._serialized_options = b'\362\336\037\014yaml:\"title\"'
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL'].fields_by_name['description']._loaded_options = None
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL'].fields_by_name['description']._serialized_options = b'\362\336\037\022yaml:\"description\"'
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL'].fields_by_name['access_config_updates']._loaded_options = None
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL'].fields_by_name['access_config_updates']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL']._loaded_options = None
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*$wasm/UpdateInstantiateConfigProposal'
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['run_as']._loaded_options = None
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['run_as']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['wasm_byte_code']._loaded_options = None
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['wasm_byte_code']._serialized_options = b'\342\336\037\014WASMByteCode'
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['msg']._loaded_options = None
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['msg']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['funds']._loaded_options = None
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL'].fields_by_name['funds']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL']._loaded_options = None
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL']._serialized_options = b'\030\001\312\264-\032cosmos.gov.v1beta1.Content\212\347\260*(wasm/StoreAndInstantiateContractProposal'
  _globals['_STORECODEPROPOSAL']._serialized_start=191
  _globals['_STORECODEPROPOSAL']._serialized_end=641
  _globals['_INSTANTIATECONTRACTPROPOSAL']._serialized_start=644
  _globals['_INSTANTIATECONTRACTPROPOSAL']._serialized_end=1135
  _globals['_INSTANTIATECONTRACT2PROPOSAL']._serialized_start=1138
  _globals['_INSTANTIATECONTRACT2PROPOSAL']._serialized_end=1676
  _globals['_MIGRATECONTRACTPROPOSAL']._serialized_start=1679
  _globals['_MIGRATECONTRACTPROPOSAL']._serialized_end=1976
  _globals['_SUDOCONTRACTPROPOSAL']._serialized_start=1979
  _globals['_SUDOCONTRACTPROPOSAL']._serialized_end=2233
  _globals['_EXECUTECONTRACTPROPOSAL']._serialized_start=2236
  _globals['_EXECUTECONTRACTPROPOSAL']._serialized_end=2666
  _globals['_UPDATEADMINPROPOSAL']._serialized_start=2669
  _globals['_UPDATEADMINPROPOSAL']._serialized_end=2938
  _globals['_CLEARADMINPROPOSAL']._serialized_start=2941
  _globals['_CLEARADMINPROPOSAL']._serialized_end=3133
  _globals['_PINCODESPROPOSAL']._serialized_start=3136
  _globals['_PINCODESPROPOSAL']._serialized_end=3329
  _globals['_UNPINCODESPROPOSAL']._serialized_start=3332
  _globals['_UNPINCODESPROPOSAL']._serialized_end=3529
  _globals['_ACCESSCONFIGUPDATE']._serialized_start=3532
  _globals['_ACCESSCONFIGUPDATE']._serialized_end=3687
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL']._serialized_start=3690
  _globals['_UPDATEINSTANTIATECONFIGPROPOSAL']._serialized_end=3997
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL']._serialized_start=4000
  _globals['_STOREANDINSTANTIATECONTRACTPROPOSAL']._serialized_end=4697
# @@protoc_insertion_point(module_scope)
