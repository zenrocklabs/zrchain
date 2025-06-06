# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/auth/v1beta1/auth.proto
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
    'cosmos/auth/v1beta1/auth.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1e\x63osmos/auth/v1beta1/auth.proto\x12\x13\x63osmos.auth.v1beta1\x1a\x11\x61mino/amino.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\"\xa1\x02\n\x0b\x42\x61seAccount\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\x12V\n\x07pub_key\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyB\'\xea\xde\x1f\x14public_key,omitempty\xa2\xe7\xb0*\npublic_keyR\x06pubKey\x12%\n\x0e\x61\x63\x63ount_number\x18\x03 \x01(\x04R\raccountNumber\x12\x1a\n\x08sequence\x18\x04 \x01(\x04R\x08sequence:C\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\xca\xb4-\x1c\x63osmos.auth.v1beta1.AccountI\x8a\xe7\xb0*\x16\x63osmos-sdk/BaseAccount\"\xec\x01\n\rModuleAccount\x12I\n\x0c\x62\x61se_account\x18\x01 \x01(\x0b\x32 .cosmos.auth.v1beta1.BaseAccountB\x04\xd0\xde\x1f\x01R\x0b\x62\x61seAccount\x12\x12\n\x04name\x18\x02 \x01(\tR\x04name\x12 \n\x0bpermissions\x18\x03 \x03(\tR\x0bpermissions:Z\x88\xa0\x1f\x00\xca\xb4-\"cosmos.auth.v1beta1.ModuleAccountI\x8a\xe7\xb0*\x18\x63osmos-sdk/ModuleAccount\x92\xe7\xb0*\x0emodule_account\"\x97\x01\n\x10ModuleCredential\x12\x1f\n\x0bmodule_name\x18\x01 \x01(\tR\nmoduleName\x12\'\n\x0f\x64\x65rivation_keys\x18\x02 \x03(\x0cR\x0e\x64\x65rivationKeys:9\xd2\xb4-\x0f\x63osmos-sdk 0.47\x8a\xe7\xb0*!cosmos-sdk/GroupAccountCredential\"\xd7\x02\n\x06Params\x12.\n\x13max_memo_characters\x18\x01 \x01(\x04R\x11maxMemoCharacters\x12 \n\x0ctx_sig_limit\x18\x02 \x01(\x04R\ntxSigLimit\x12\x30\n\x15tx_size_cost_per_byte\x18\x03 \x01(\x04R\x11txSizeCostPerByte\x12O\n\x17sig_verify_cost_ed25519\x18\x04 \x01(\x04\x42\x18\xe2\xde\x1f\x14SigVerifyCostED25519R\x14sigVerifyCostEd25519\x12U\n\x19sig_verify_cost_secp256k1\x18\x05 \x01(\x04\x42\x1a\xe2\xde\x1f\x16SigVerifyCostSecp256k1R\x16sigVerifyCostSecp256k1:!\xe8\xa0\x1f\x01\x8a\xe7\xb0*\x18\x63osmos-sdk/x/auth/ParamsB+Z)github.com/cosmos/cosmos-sdk/x/auth/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.auth.v1beta1.auth_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)github.com/cosmos/cosmos-sdk/x/auth/types'
  _globals['_BASEACCOUNT'].fields_by_name['address']._loaded_options = None
  _globals['_BASEACCOUNT'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_BASEACCOUNT'].fields_by_name['pub_key']._loaded_options = None
  _globals['_BASEACCOUNT'].fields_by_name['pub_key']._serialized_options = b'\352\336\037\024public_key,omitempty\242\347\260*\npublic_key'
  _globals['_BASEACCOUNT']._loaded_options = None
  _globals['_BASEACCOUNT']._serialized_options = b'\210\240\037\000\350\240\037\000\312\264-\034cosmos.auth.v1beta1.AccountI\212\347\260*\026cosmos-sdk/BaseAccount'
  _globals['_MODULEACCOUNT'].fields_by_name['base_account']._loaded_options = None
  _globals['_MODULEACCOUNT'].fields_by_name['base_account']._serialized_options = b'\320\336\037\001'
  _globals['_MODULEACCOUNT']._loaded_options = None
  _globals['_MODULEACCOUNT']._serialized_options = b'\210\240\037\000\312\264-\"cosmos.auth.v1beta1.ModuleAccountI\212\347\260*\030cosmos-sdk/ModuleAccount\222\347\260*\016module_account'
  _globals['_MODULECREDENTIAL']._loaded_options = None
  _globals['_MODULECREDENTIAL']._serialized_options = b'\322\264-\017cosmos-sdk 0.47\212\347\260*!cosmos-sdk/GroupAccountCredential'
  _globals['_PARAMS'].fields_by_name['sig_verify_cost_ed25519']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['sig_verify_cost_ed25519']._serialized_options = b'\342\336\037\024SigVerifyCostED25519'
  _globals['_PARAMS'].fields_by_name['sig_verify_cost_secp256k1']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['sig_verify_cost_secp256k1']._serialized_options = b'\342\336\037\026SigVerifyCostSecp256k1'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\350\240\037\001\212\347\260*\030cosmos-sdk/x/auth/Params'
  _globals['_BASEACCOUNT']._serialized_start=151
  _globals['_BASEACCOUNT']._serialized_end=440
  _globals['_MODULEACCOUNT']._serialized_start=443
  _globals['_MODULEACCOUNT']._serialized_end=679
  _globals['_MODULECREDENTIAL']._serialized_start=682
  _globals['_MODULECREDENTIAL']._serialized_end=833
  _globals['_PARAMS']._serialized_start=836
  _globals['_PARAMS']._serialized_end=1179
# @@protoc_insertion_point(module_scope)
