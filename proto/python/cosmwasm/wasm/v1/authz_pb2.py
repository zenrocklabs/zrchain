# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmwasm/wasm/v1/authz.proto
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
    'cosmwasm/wasm/v1/authz.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmwasm.wasm.v1 import types_pb2 as cosmwasm_dot_wasm_dot_v1_dot_types__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1c\x63osmwasm/wasm/v1/authz.proto\x12\x10\x63osmwasm.wasm.v1\x1a\x14gogoproto/gogo.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x1c\x63osmwasm/wasm/v1/types.proto\x1a\x19google/protobuf/any.proto\x1a\x11\x61mino/amino.proto\"\xa0\x01\n\x16StoreCodeAuthorization\x12>\n\x06grants\x18\x01 \x03(\x0b\x32\x1b.cosmwasm.wasm.v1.CodeGrantB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06grants:F\xca\xb4-\"cosmos.authz.v1beta1.Authorization\x8a\xe7\xb0*\x1bwasm/StoreCodeAuthorization\"\xb4\x01\n\x1e\x43ontractExecutionAuthorization\x12\x42\n\x06grants\x18\x01 \x03(\x0b\x32\x1f.cosmwasm.wasm.v1.ContractGrantB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06grants:N\xca\xb4-\"cosmos.authz.v1beta1.Authorization\x8a\xe7\xb0*#wasm/ContractExecutionAuthorization\"\xb4\x01\n\x1e\x43ontractMigrationAuthorization\x12\x42\n\x06grants\x18\x01 \x03(\x0b\x32\x1f.cosmwasm.wasm.v1.ContractGrantB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06grants:N\xca\xb4-\"cosmos.authz.v1beta1.Authorization\x8a\xe7\xb0*#wasm/ContractMigrationAuthorization\"\x7f\n\tCodeGrant\x12\x1b\n\tcode_hash\x18\x01 \x01(\x0cR\x08\x63odeHash\x12U\n\x16instantiate_permission\x18\x02 \x01(\x0b\x32\x1e.cosmwasm.wasm.v1.AccessConfigR\x15instantiatePermission\"\xf4\x01\n\rContractGrant\x12\x34\n\x08\x63ontract\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x08\x63ontract\x12T\n\x05limit\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyB(\xca\xb4-$cosmwasm.wasm.v1.ContractAuthzLimitXR\x05limit\x12W\n\x06\x66ilter\x18\x03 \x01(\x0b\x32\x14.google.protobuf.AnyB)\xca\xb4-%cosmwasm.wasm.v1.ContractAuthzFilterXR\x06\x66ilter\"n\n\rMaxCallsLimit\x12\x1c\n\tremaining\x18\x01 \x01(\x04R\tremaining:?\xca\xb4-$cosmwasm.wasm.v1.ContractAuthzLimitX\x8a\xe7\xb0*\x12wasm/MaxCallsLimit\"\xcd\x01\n\rMaxFundsLimit\x12{\n\x07\x61mounts\x18\x01 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x07\x61mounts:?\xca\xb4-$cosmwasm.wasm.v1.ContractAuthzLimitX\x8a\xe7\xb0*\x12wasm/MaxFundsLimit\"\xf6\x01\n\rCombinedLimit\x12\'\n\x0f\x63\x61lls_remaining\x18\x01 \x01(\x04R\x0e\x63\x61llsRemaining\x12{\n\x07\x61mounts\x18\x02 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x07\x61mounts:?\xca\xb4-$cosmwasm.wasm.v1.ContractAuthzLimitX\x8a\xe7\xb0*\x12wasm/CombinedLimit\"c\n\x16\x41llowAllMessagesFilter:I\xca\xb4-%cosmwasm.wasm.v1.ContractAuthzFilterX\x8a\xe7\xb0*\x1bwasm/AllowAllMessagesFilter\"}\n\x19\x41\x63\x63\x65ptedMessageKeysFilter\x12\x12\n\x04keys\x18\x01 \x03(\tR\x04keys:L\xca\xb4-%cosmwasm.wasm.v1.ContractAuthzFilterX\x8a\xe7\xb0*\x1ewasm/AcceptedMessageKeysFilter\"\xa7\x01\n\x16\x41\x63\x63\x65ptedMessagesFilter\x12\x42\n\x08messages\x18\x01 \x03(\x0c\x42&\xfa\xde\x1f\x12RawContractMessage\x9a\xe7\xb0*\x0binline_jsonR\x08messages:I\xca\xb4-%cosmwasm.wasm.v1.ContractAuthzFilterX\x8a\xe7\xb0*\x1bwasm/AcceptedMessagesFilterB,Z&github.com/CosmWasm/wasmd/x/wasm/types\xc8\xe1\x1e\x00\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmwasm.wasm.v1.authz_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z&github.com/CosmWasm/wasmd/x/wasm/types\310\341\036\000'
  _globals['_STORECODEAUTHORIZATION'].fields_by_name['grants']._loaded_options = None
  _globals['_STORECODEAUTHORIZATION'].fields_by_name['grants']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_STORECODEAUTHORIZATION']._loaded_options = None
  _globals['_STORECODEAUTHORIZATION']._serialized_options = b'\312\264-\"cosmos.authz.v1beta1.Authorization\212\347\260*\033wasm/StoreCodeAuthorization'
  _globals['_CONTRACTEXECUTIONAUTHORIZATION'].fields_by_name['grants']._loaded_options = None
  _globals['_CONTRACTEXECUTIONAUTHORIZATION'].fields_by_name['grants']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_CONTRACTEXECUTIONAUTHORIZATION']._loaded_options = None
  _globals['_CONTRACTEXECUTIONAUTHORIZATION']._serialized_options = b'\312\264-\"cosmos.authz.v1beta1.Authorization\212\347\260*#wasm/ContractExecutionAuthorization'
  _globals['_CONTRACTMIGRATIONAUTHORIZATION'].fields_by_name['grants']._loaded_options = None
  _globals['_CONTRACTMIGRATIONAUTHORIZATION'].fields_by_name['grants']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_CONTRACTMIGRATIONAUTHORIZATION']._loaded_options = None
  _globals['_CONTRACTMIGRATIONAUTHORIZATION']._serialized_options = b'\312\264-\"cosmos.authz.v1beta1.Authorization\212\347\260*#wasm/ContractMigrationAuthorization'
  _globals['_CONTRACTGRANT'].fields_by_name['contract']._loaded_options = None
  _globals['_CONTRACTGRANT'].fields_by_name['contract']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_CONTRACTGRANT'].fields_by_name['limit']._loaded_options = None
  _globals['_CONTRACTGRANT'].fields_by_name['limit']._serialized_options = b'\312\264-$cosmwasm.wasm.v1.ContractAuthzLimitX'
  _globals['_CONTRACTGRANT'].fields_by_name['filter']._loaded_options = None
  _globals['_CONTRACTGRANT'].fields_by_name['filter']._serialized_options = b'\312\264-%cosmwasm.wasm.v1.ContractAuthzFilterX'
  _globals['_MAXCALLSLIMIT']._loaded_options = None
  _globals['_MAXCALLSLIMIT']._serialized_options = b'\312\264-$cosmwasm.wasm.v1.ContractAuthzLimitX\212\347\260*\022wasm/MaxCallsLimit'
  _globals['_MAXFUNDSLIMIT'].fields_by_name['amounts']._loaded_options = None
  _globals['_MAXFUNDSLIMIT'].fields_by_name['amounts']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_MAXFUNDSLIMIT']._loaded_options = None
  _globals['_MAXFUNDSLIMIT']._serialized_options = b'\312\264-$cosmwasm.wasm.v1.ContractAuthzLimitX\212\347\260*\022wasm/MaxFundsLimit'
  _globals['_COMBINEDLIMIT'].fields_by_name['amounts']._loaded_options = None
  _globals['_COMBINEDLIMIT'].fields_by_name['amounts']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_COMBINEDLIMIT']._loaded_options = None
  _globals['_COMBINEDLIMIT']._serialized_options = b'\312\264-$cosmwasm.wasm.v1.ContractAuthzLimitX\212\347\260*\022wasm/CombinedLimit'
  _globals['_ALLOWALLMESSAGESFILTER']._loaded_options = None
  _globals['_ALLOWALLMESSAGESFILTER']._serialized_options = b'\312\264-%cosmwasm.wasm.v1.ContractAuthzFilterX\212\347\260*\033wasm/AllowAllMessagesFilter'
  _globals['_ACCEPTEDMESSAGEKEYSFILTER']._loaded_options = None
  _globals['_ACCEPTEDMESSAGEKEYSFILTER']._serialized_options = b'\312\264-%cosmwasm.wasm.v1.ContractAuthzFilterX\212\347\260*\036wasm/AcceptedMessageKeysFilter'
  _globals['_ACCEPTEDMESSAGESFILTER'].fields_by_name['messages']._loaded_options = None
  _globals['_ACCEPTEDMESSAGESFILTER'].fields_by_name['messages']._serialized_options = b'\372\336\037\022RawContractMessage\232\347\260*\013inline_json'
  _globals['_ACCEPTEDMESSAGESFILTER']._loaded_options = None
  _globals['_ACCEPTEDMESSAGESFILTER']._serialized_options = b'\312\264-%cosmwasm.wasm.v1.ContractAuthzFilterX\212\347\260*\033wasm/AcceptedMessagesFilter'
  _globals['_STORECODEAUTHORIZATION']._serialized_start=208
  _globals['_STORECODEAUTHORIZATION']._serialized_end=368
  _globals['_CONTRACTEXECUTIONAUTHORIZATION']._serialized_start=371
  _globals['_CONTRACTEXECUTIONAUTHORIZATION']._serialized_end=551
  _globals['_CONTRACTMIGRATIONAUTHORIZATION']._serialized_start=554
  _globals['_CONTRACTMIGRATIONAUTHORIZATION']._serialized_end=734
  _globals['_CODEGRANT']._serialized_start=736
  _globals['_CODEGRANT']._serialized_end=863
  _globals['_CONTRACTGRANT']._serialized_start=866
  _globals['_CONTRACTGRANT']._serialized_end=1110
  _globals['_MAXCALLSLIMIT']._serialized_start=1112
  _globals['_MAXCALLSLIMIT']._serialized_end=1222
  _globals['_MAXFUNDSLIMIT']._serialized_start=1225
  _globals['_MAXFUNDSLIMIT']._serialized_end=1430
  _globals['_COMBINEDLIMIT']._serialized_start=1433
  _globals['_COMBINEDLIMIT']._serialized_end=1679
  _globals['_ALLOWALLMESSAGESFILTER']._serialized_start=1681
  _globals['_ALLOWALLMESSAGESFILTER']._serialized_end=1780
  _globals['_ACCEPTEDMESSAGEKEYSFILTER']._serialized_start=1782
  _globals['_ACCEPTEDMESSAGEKEYSFILTER']._serialized_end=1907
  _globals['_ACCEPTEDMESSAGESFILTER']._serialized_start=1910
  _globals['_ACCEPTEDMESSAGESFILTER']._serialized_end=2077
# @@protoc_insertion_point(module_scope)
