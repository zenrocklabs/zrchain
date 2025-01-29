# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/treasury/key.proto
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
    'zrchain/treasury/key.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from zrchain.treasury import wallet_pb2 as zrchain_dot_treasury_dot_wallet__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1azrchain/treasury/key.proto\x12\x10zrchain.treasury\x1a\x1dzrchain/treasury/wallet.proto\x1a\x19google/protobuf/any.proto\"\xd8\x03\n\nKeyRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x18\n\x07\x63reator\x18\x02 \x01(\tR\x07\x63reator\x12%\n\x0eworkspace_addr\x18\x03 \x01(\tR\rworkspaceAddr\x12!\n\x0ckeyring_addr\x18\x04 \x01(\tR\x0bkeyringAddr\x12\x34\n\x08key_type\x18\x05 \x01(\x0e\x32\x19.zrchain.treasury.KeyTypeR\x07keyType\x12:\n\x06status\x18\x06 \x01(\x0e\x32\".zrchain.treasury.KeyRequestStatusR\x06status\x12\x38\n\x18keyring_party_signatures\x18\x07 \x03(\x0cR\x16keyringPartySignatures\x12#\n\rreject_reason\x18\x08 \x01(\tR\x0crejectReason\x12\x14\n\x05index\x18\t \x01(\x04R\x05index\x12$\n\x0esign_policy_id\x18\n \x01(\x04R\x0csignPolicyId\x12I\n\x0fzenbtc_metadata\x18\x0b \x01(\x0b\x32 .zrchain.treasury.ZenBTCMetadataR\x0ezenbtcMetadata\"\x9d\x03\n\x0eKeyReqResponse\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x18\n\x07\x63reator\x18\x02 \x01(\tR\x07\x63reator\x12%\n\x0eworkspace_addr\x18\x03 \x01(\tR\rworkspaceAddr\x12!\n\x0ckeyring_addr\x18\x04 \x01(\tR\x0bkeyringAddr\x12\x19\n\x08key_type\x18\x05 \x01(\tR\x07keyType\x12\x16\n\x06status\x18\x06 \x01(\tR\x06status\x12\x38\n\x18keyring_party_signatures\x18\x07 \x03(\x0cR\x16keyringPartySignatures\x12#\n\rreject_reason\x18\x08 \x01(\tR\x0crejectReason\x12\x14\n\x05index\x18\t \x01(\x04R\x05index\x12$\n\x0esign_policy_id\x18\n \x01(\x04R\x0csignPolicyId\x12I\n\x0fzenbtc_metadata\x18\x0b \x01(\x0b\x32 .zrchain.treasury.ZenBTCMetadataR\x0ezenbtcMetadata\"\xb4\x02\n\x03Key\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12%\n\x0eworkspace_addr\x18\x02 \x01(\tR\rworkspaceAddr\x12!\n\x0ckeyring_addr\x18\x03 \x01(\tR\x0bkeyringAddr\x12-\n\x04type\x18\x04 \x01(\x0e\x32\x19.zrchain.treasury.KeyTypeR\x04type\x12\x1d\n\npublic_key\x18\x05 \x01(\x0cR\tpublicKey\x12\x14\n\x05index\x18\t \x01(\x04R\x05index\x12$\n\x0esign_policy_id\x18\n \x01(\x04R\x0csignPolicyId\x12I\n\x0fzenbtc_metadata\x18\x0b \x01(\x0b\x32 .zrchain.treasury.ZenBTCMetadataR\x0ezenbtcMetadata\"\xa1\x02\n\x0bKeyResponse\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12%\n\x0eworkspace_addr\x18\x02 \x01(\tR\rworkspaceAddr\x12!\n\x0ckeyring_addr\x18\x03 \x01(\tR\x0bkeyringAddr\x12\x12\n\x04type\x18\x04 \x01(\tR\x04type\x12\x1d\n\npublic_key\x18\x05 \x01(\x0cR\tpublicKey\x12\x14\n\x05index\x18\t \x01(\x04R\x05index\x12$\n\x0esign_policy_id\x18\n \x01(\x04R\x0csignPolicyId\x12I\n\x0fzenbtc_metadata\x18\x0b \x01(\x0b\x32 .zrchain.treasury.ZenBTCMetadataR\x0ezenbtcMetadata\"\xf3\x01\n\x0eZenBTCMetadata\x12%\n\x0erecipient_addr\x18\x01 \x01(\tR\rrecipientAddr\x12;\n\nchain_type\x18\x02 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\tchainType\x12\x19\n\x08\x63hain_id\x18\x03 \x01(\x04R\x07\x63hainId\x12%\n\x0ereturn_address\x18\x04 \x01(\tR\rreturnAddress\x12;\n\x0e\x63hain_metadata\x18\x05 \x01(\x0b\x32\x14.google.protobuf.AnyR\rchainMetadata*\xb9\x01\n\x10KeyRequestStatus\x12\"\n\x1eKEY_REQUEST_STATUS_UNSPECIFIED\x10\x00\x12\x1e\n\x1aKEY_REQUEST_STATUS_PENDING\x10\x01\x12\x1e\n\x1aKEY_REQUEST_STATUS_PARTIAL\x10\x02\x12 \n\x1cKEY_REQUEST_STATUS_FULFILLED\x10\x03\x12\x1f\n\x1bKEY_REQUEST_STATUS_REJECTED\x10\x04*}\n\x07KeyType\x12\x18\n\x14KEY_TYPE_UNSPECIFIED\x10\x00\x12\x1c\n\x18KEY_TYPE_ECDSA_SECP256K1\x10\x01\x12\x1a\n\x16KEY_TYPE_EDDSA_ED25519\x10\x02\x12\x1e\n\x1aKEY_TYPE_BITCOIN_SECP256K1\x10\x03\x42;Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.treasury.key_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types'
  _globals['_KEYREQUESTSTATUS']._serialized_start=1847
  _globals['_KEYREQUESTSTATUS']._serialized_end=2032
  _globals['_KEYTYPE']._serialized_start=2034
  _globals['_KEYTYPE']._serialized_end=2159
  _globals['_KEYREQUEST']._serialized_start=107
  _globals['_KEYREQUEST']._serialized_end=579
  _globals['_KEYREQRESPONSE']._serialized_start=582
  _globals['_KEYREQRESPONSE']._serialized_end=995
  _globals['_KEY']._serialized_start=998
  _globals['_KEY']._serialized_end=1306
  _globals['_KEYRESPONSE']._serialized_start=1309
  _globals['_KEYRESPONSE']._serialized_end=1598
  _globals['_ZENBTCMETADATA']._serialized_start=1601
  _globals['_ZENBTCMETADATA']._serialized_end=1844
# @@protoc_insertion_point(module_scope)
