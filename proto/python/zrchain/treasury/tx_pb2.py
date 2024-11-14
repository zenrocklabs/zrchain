# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/treasury/tx.proto
# Protobuf Python Version: 5.28.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    3,
    '',
    'zrchain/treasury/tx.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos.msg.v1 import msg_pb2 as cosmos_dot_msg_dot_v1_dot_msg__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from zrchain.treasury import key_pb2 as zrchain_dot_treasury_dot_key__pb2
from zrchain.treasury import mpcsign_pb2 as zrchain_dot_treasury_dot_mpcsign__pb2
from zrchain.treasury import params_pb2 as zrchain_dot_treasury_dot_params__pb2
from zrchain.treasury import wallet_pb2 as zrchain_dot_treasury_dot_wallet__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x19zrchain/treasury/tx.proto\x12\x10zrchain.treasury\x1a\x11\x61mino/amino.proto\x1a\x17\x63osmos/msg/v1/msg.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a\x1azrchain/treasury/key.proto\x1a\x1ezrchain/treasury/mpcsign.proto\x1a\x1dzrchain/treasury/params.proto\x1a\x1dzrchain/treasury/wallet.proto\"\xde\x01\n\x0fMsgUpdateParams\x12\x36\n\tauthority\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tauthority\x12;\n\x06params\x18\x02 \x01(\x0b\x32\x18.zrchain.treasury.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params:V\x82\xe7\xb0*\tauthority\x8a\xe7\xb0*Cgithub.com/Zenrock-Foundation/zrchain/v5/x/treasury/MsgUpdateParams\"\x19\n\x17MsgUpdateParamsResponse\"\xff\x02\n\x10MsgNewKeyRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12%\n\x0eworkspace_addr\x18\x02 \x01(\tR\rworkspaceAddr\x12!\n\x0ckeyring_addr\x18\x03 \x01(\tR\x0bkeyringAddr\x12\x19\n\x08key_type\x18\x04 \x01(\tR\x07keyType\x12\x10\n\x03\x62tl\x18\x05 \x01(\x04R\x03\x62tl\x12\x14\n\x05index\x18\x06 \x01(\x04R\x05index\x12#\n\rext_requester\x18\x07 \x01(\tR\x0c\x65xtRequester\x12 \n\x0c\x65xt_key_type\x18\x08 \x01(\x04R\nextKeyType\x12$\n\x0esign_policy_id\x18\t \x01(\x04R\x0csignPolicyId\x12I\n\x0fzenbtc_metadata\x18\n \x01(\x0b\x32 .zrchain.treasury.ZenBTCMetadataR\x0ezenbtcMetadata:\x0c\x82\xe7\xb0*\x07\x63reator\"8\n\x18MsgNewKeyRequestResponse\x12\x1c\n\nkey_req_id\x18\x01 \x01(\x04R\x08keyReqId\"\xb2\x02\n\x13MsgFulfilKeyRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x1d\n\nrequest_id\x18\x02 \x01(\x04R\trequestId\x12:\n\x06status\x18\x03 \x01(\x0e\x32\".zrchain.treasury.KeyRequestStatusR\x06status\x12/\n\x03key\x18\x04 \x01(\x0b\x32\x1b.zrchain.treasury.MsgNewKeyH\x00R\x03key\x12%\n\rreject_reason\x18\x05 \x01(\tH\x00R\x0crejectReason\x12\x36\n\x17keyring_party_signature\x18\x06 \x01(\x0cR\x15keyringPartySignature:\x0c\x82\xe7\xb0*\x07\x63reatorB\x08\n\x06result\"*\n\tMsgNewKey\x12\x1d\n\npublic_key\x18\x01 \x01(\x0cR\tpublicKey\"\x1d\n\x1bMsgFulfilKeyRequestResponse\"\xc6\x02\n\x16MsgNewSignatureRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x17\n\x07key_ids\x18\x02 \x03(\x04R\x06keyIds\x12(\n\x10\x64\x61ta_for_signing\x18\x03 \x01(\tR\x0e\x64\x61taForSigning\x12\x10\n\x03\x62tl\x18\x04 \x01(\x04R\x03\x62tl\x12\x19\n\x08\x63\x61\x63he_id\x18\x05 \x01(\x0cR\x07\x63\x61\x63heId\x12.\n\x13verify_signing_data\x18\x06 \x01(\x0cR\x11verifySigningData\x12\x64\n\x1bverify_signing_data_version\x18\x07 \x01(\x0e\x32%.zrchain.treasury.VerificationVersionR\x18verifySigningDataVersion:\x0c\x82\xe7\xb0*\x07\x63reator\">\n\x1eMsgNewSignatureRequestResponse\x12\x1c\n\nsig_req_id\x18\x01 \x01(\x04R\x08sigReqId\"\x9d\x02\n\x19MsgFulfilSignatureRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x1d\n\nrequest_id\x18\x02 \x01(\x04R\trequestId\x12;\n\x06status\x18\x03 \x01(\x0e\x32#.zrchain.treasury.SignRequestStatusR\x06status\x12\x36\n\x17keyring_party_signature\x18\x04 \x01(\x0cR\x15keyringPartySignature\x12\x1f\n\x0bsigned_data\x18\x05 \x01(\x0cR\nsignedData\x12#\n\rreject_reason\x18\x06 \x01(\tR\x0crejectReason:\x0c\x82\xe7\xb0*\x07\x63reator\"#\n!MsgFulfilSignatureRequestResponse\"-\n\x10MetadataEthereum\x12\x19\n\x08\x63hain_id\x18\x01 \x01(\x04R\x07\x63hainId\"O\n\x0eMetadataSolana\x12=\n\x07network\x18\x01 \x01(\x0e\x32#.zrchain.treasury.SolanaNetworkTypeR\x07network\"\xd1\x02\n\x1cMsgNewSignTransactionRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x15\n\x06key_id\x18\x02 \x01(\x04R\x05keyId\x12=\n\x0bwallet_type\x18\x03 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\nwalletType\x12\x31\n\x14unsigned_transaction\x18\x04 \x01(\x0cR\x13unsignedTransaction\x12\x30\n\x08metadata\x18\x05 \x01(\x0b\x32\x14.google.protobuf.AnyR\x08metadata\x12\x10\n\x03\x62tl\x18\x06 \x01(\x04R\x03\x62tl\x12\x19\n\x08\x63\x61\x63he_id\x18\x07 \x01(\x0cR\x07\x63\x61\x63heId\x12!\n\x0cno_broadcast\x18\x08 \x01(\x08R\x0bnoBroadcast:\x0c\x82\xe7\xb0*\x07\x63reator\"h\n$MsgNewSignTransactionRequestResponse\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x30\n\x14signature_request_id\x18\x02 \x01(\x04R\x12signatureRequestId\"\xa6\x01\n\x16MsgTransferFromKeyring\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x18\n\x07keyring\x18\x02 \x01(\tR\x07keyring\x12\x1c\n\trecipient\x18\x03 \x01(\tR\trecipient\x12\x16\n\x06\x61mount\x18\x04 \x01(\x04R\x06\x61mount\x12\x14\n\x05\x64\x65nom\x18\x05 \x01(\tR\x05\x64\x65nom:\x0c\x82\xe7\xb0*\x07\x63reator\" \n\x1eMsgTransferFromKeyringResponse\"\xf6\x01\n\x1bMsgNewICATransactionRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x15\n\x06key_id\x18\x02 \x01(\x04R\x05keyId\x12#\n\rinput_payload\x18\x03 \x01(\tR\x0cinputPayload\x12#\n\rconnection_id\x18\x04 \x01(\tR\x0c\x63onnectionId\x12<\n\x1arelative_timeout_timestamp\x18\x05 \x01(\x04R\x18relativeTimeoutTimestamp\x12\x10\n\x03\x62tl\x18\x06 \x01(\x04R\x03\x62tl:\x0c\x82\xe7\xb0*\x07\x63reator\"g\n#MsgNewICATransactionRequestResponse\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x30\n\x14signature_request_id\x18\x02 \x01(\x04R\x12signatureRequestId\"\xa2\x02\n\x1eMsgFulfilICATransactionRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x1d\n\nrequest_id\x18\x02 \x01(\x04R\trequestId\x12;\n\x06status\x18\x03 \x01(\x0e\x32#.zrchain.treasury.SignRequestStatusR\x06status\x12\x36\n\x17keyring_party_signature\x18\x04 \x01(\x0cR\x15keyringPartySignature\x12\x1f\n\x0bsigned_data\x18\x05 \x01(\x0cR\nsignedData\x12#\n\rreject_reason\x18\x06 \x01(\tR\x0crejectReason:\x0c\x82\xe7\xb0*\x07\x63reator\"(\n&MsgFulfilICATransactionRequestResponse\"\x99\x04\n\x1cMsgNewZrSignSignatureRequest\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ress\x12\x19\n\x08key_type\x18\x03 \x01(\x04R\x07keyType\x12!\n\x0cwallet_index\x18\x04 \x01(\x04R\x0bwalletIndex\x12\x19\n\x08\x63\x61\x63he_id\x18\x05 \x01(\x0cR\x07\x63\x61\x63heId\x12\x12\n\x04\x64\x61ta\x18\x06 \x01(\tR\x04\x64\x61ta\x12.\n\x13verify_signing_data\x18\x07 \x01(\x0cR\x11verifySigningData\x12\x64\n\x1bverify_signing_data_version\x18\x08 \x01(\x0e\x32%.zrchain.treasury.VerificationVersionR\x18verifySigningDataVersion\x12=\n\x0bwallet_type\x18\t \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\nwalletType\x12\x30\n\x08metadata\x18\n \x01(\x0b\x32\x14.google.protobuf.AnyR\x08metadata\x12!\n\x0cno_broadcast\x18\x0b \x01(\x08R\x0bnoBroadcast\x12\x10\n\x03\x62tl\x18\x0c \x01(\x04R\x03\x62tl\x12\x0e\n\x02tx\x18\r \x01(\x08R\x02tx:\x0c\x82\xe7\xb0*\x07\x63reator\"=\n$MsgNewZrSignSignatureRequestResponse\x12\x15\n\x06req_id\x18\x01 \x01(\x04R\x05reqId\"y\n\x12MsgUpdateKeyPolicy\x12\x18\n\x07\x63reator\x18\x01 \x01(\tR\x07\x63reator\x12\x15\n\x06key_id\x18\x02 \x01(\x04R\x05keyId\x12$\n\x0esign_policy_id\x18\x03 \x01(\x04R\x0csignPolicyId:\x0c\x82\xe7\xb0*\x07\x63reator\"\x1c\n\x1aMsgUpdateKeyPolicyResponse*4\n\x13VerificationVersion\x12\x0b\n\x07UNKNOWN\x10\x00\x12\x10\n\x0c\x42ITCOIN_PLUS\x10\x01*H\n\x11SolanaNetworkType\x12\r\n\tUNDEFINED\x10\x00\x12\x0b\n\x07MAINNET\x10\x01\x12\n\n\x06\x44\x45VNET\x10\x02\x12\x0b\n\x07TESTNET\x10\x03\x32\x99\n\n\x03Msg\x12\\\n\x0cUpdateParams\x12!.zrchain.treasury.MsgUpdateParams\x1a).zrchain.treasury.MsgUpdateParamsResponse\x12_\n\rNewKeyRequest\x12\".zrchain.treasury.MsgNewKeyRequest\x1a*.zrchain.treasury.MsgNewKeyRequestResponse\x12h\n\x10\x46ulfilKeyRequest\x12%.zrchain.treasury.MsgFulfilKeyRequest\x1a-.zrchain.treasury.MsgFulfilKeyRequestResponse\x12q\n\x13NewSignatureRequest\x12(.zrchain.treasury.MsgNewSignatureRequest\x1a\x30.zrchain.treasury.MsgNewSignatureRequestResponse\x12z\n\x16\x46ulfilSignatureRequest\x12+.zrchain.treasury.MsgFulfilSignatureRequest\x1a\x33.zrchain.treasury.MsgFulfilSignatureRequestResponse\x12\x83\x01\n\x19NewSignTransactionRequest\x12..zrchain.treasury.MsgNewSignTransactionRequest\x1a\x36.zrchain.treasury.MsgNewSignTransactionRequestResponse\x12q\n\x13TransferFromKeyring\x12(.zrchain.treasury.MsgTransferFromKeyring\x1a\x30.zrchain.treasury.MsgTransferFromKeyringResponse\x12\x80\x01\n\x18NewICATransactionRequest\x12-.zrchain.treasury.MsgNewICATransactionRequest\x1a\x35.zrchain.treasury.MsgNewICATransactionRequestResponse\x12\x89\x01\n\x1b\x46ulfilICATransactionRequest\x12\x30.zrchain.treasury.MsgFulfilICATransactionRequest\x1a\x38.zrchain.treasury.MsgFulfilICATransactionRequestResponse\x12\x83\x01\n\x19NewZrSignSignatureRequest\x12..zrchain.treasury.MsgNewZrSignSignatureRequest\x1a\x36.zrchain.treasury.MsgNewZrSignSignatureRequestResponse\x12\x65\n\x0fUpdateKeyPolicy\x12$.zrchain.treasury.MsgUpdateKeyPolicy\x1a,.zrchain.treasury.MsgUpdateKeyPolicyResponse\x1a\x05\x80\xe7\xb0*\x01\x42;Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.treasury.tx_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['authority']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._loaded_options = None
  _globals['_MSGUPDATEPARAMS'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_MSGUPDATEPARAMS']._loaded_options = None
  _globals['_MSGUPDATEPARAMS']._serialized_options = b'\202\347\260*\tauthority\212\347\260*Cgithub.com/Zenrock-Foundation/zrchain/v5/x/treasury/MsgUpdateParams'
  _globals['_MSGNEWKEYREQUEST']._loaded_options = None
  _globals['_MSGNEWKEYREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGFULFILKEYREQUEST']._loaded_options = None
  _globals['_MSGFULFILKEYREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGNEWSIGNATUREREQUEST']._loaded_options = None
  _globals['_MSGNEWSIGNATUREREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGFULFILSIGNATUREREQUEST']._loaded_options = None
  _globals['_MSGFULFILSIGNATUREREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGNEWSIGNTRANSACTIONREQUEST']._loaded_options = None
  _globals['_MSGNEWSIGNTRANSACTIONREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGTRANSFERFROMKEYRING']._loaded_options = None
  _globals['_MSGTRANSFERFROMKEYRING']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGNEWICATRANSACTIONREQUEST']._loaded_options = None
  _globals['_MSGNEWICATRANSACTIONREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGFULFILICATRANSACTIONREQUEST']._loaded_options = None
  _globals['_MSGFULFILICATRANSACTIONREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGNEWZRSIGNSIGNATUREREQUEST']._loaded_options = None
  _globals['_MSGNEWZRSIGNSIGNATUREREQUEST']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSGUPDATEKEYPOLICY']._loaded_options = None
  _globals['_MSGUPDATEKEYPOLICY']._serialized_options = b'\202\347\260*\007creator'
  _globals['_MSG']._loaded_options = None
  _globals['_MSG']._serialized_options = b'\200\347\260*\001'
  _globals['_VERIFICATIONVERSION']._serialized_start=4309
  _globals['_VERIFICATIONVERSION']._serialized_end=4361
  _globals['_SOLANANETWORKTYPE']._serialized_start=4363
  _globals['_SOLANANETWORKTYPE']._serialized_end=4435
  _globals['_MSGUPDATEPARAMS']._serialized_start=290
  _globals['_MSGUPDATEPARAMS']._serialized_end=512
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_start=514
  _globals['_MSGUPDATEPARAMSRESPONSE']._serialized_end=539
  _globals['_MSGNEWKEYREQUEST']._serialized_start=542
  _globals['_MSGNEWKEYREQUEST']._serialized_end=925
  _globals['_MSGNEWKEYREQUESTRESPONSE']._serialized_start=927
  _globals['_MSGNEWKEYREQUESTRESPONSE']._serialized_end=983
  _globals['_MSGFULFILKEYREQUEST']._serialized_start=986
  _globals['_MSGFULFILKEYREQUEST']._serialized_end=1292
  _globals['_MSGNEWKEY']._serialized_start=1294
  _globals['_MSGNEWKEY']._serialized_end=1336
  _globals['_MSGFULFILKEYREQUESTRESPONSE']._serialized_start=1338
  _globals['_MSGFULFILKEYREQUESTRESPONSE']._serialized_end=1367
  _globals['_MSGNEWSIGNATUREREQUEST']._serialized_start=1370
  _globals['_MSGNEWSIGNATUREREQUEST']._serialized_end=1696
  _globals['_MSGNEWSIGNATUREREQUESTRESPONSE']._serialized_start=1698
  _globals['_MSGNEWSIGNATUREREQUESTRESPONSE']._serialized_end=1760
  _globals['_MSGFULFILSIGNATUREREQUEST']._serialized_start=1763
  _globals['_MSGFULFILSIGNATUREREQUEST']._serialized_end=2048
  _globals['_MSGFULFILSIGNATUREREQUESTRESPONSE']._serialized_start=2050
  _globals['_MSGFULFILSIGNATUREREQUESTRESPONSE']._serialized_end=2085
  _globals['_METADATAETHEREUM']._serialized_start=2087
  _globals['_METADATAETHEREUM']._serialized_end=2132
  _globals['_METADATASOLANA']._serialized_start=2134
  _globals['_METADATASOLANA']._serialized_end=2213
  _globals['_MSGNEWSIGNTRANSACTIONREQUEST']._serialized_start=2216
  _globals['_MSGNEWSIGNTRANSACTIONREQUEST']._serialized_end=2553
  _globals['_MSGNEWSIGNTRANSACTIONREQUESTRESPONSE']._serialized_start=2555
  _globals['_MSGNEWSIGNTRANSACTIONREQUESTRESPONSE']._serialized_end=2659
  _globals['_MSGTRANSFERFROMKEYRING']._serialized_start=2662
  _globals['_MSGTRANSFERFROMKEYRING']._serialized_end=2828
  _globals['_MSGTRANSFERFROMKEYRINGRESPONSE']._serialized_start=2830
  _globals['_MSGTRANSFERFROMKEYRINGRESPONSE']._serialized_end=2862
  _globals['_MSGNEWICATRANSACTIONREQUEST']._serialized_start=2865
  _globals['_MSGNEWICATRANSACTIONREQUEST']._serialized_end=3111
  _globals['_MSGNEWICATRANSACTIONREQUESTRESPONSE']._serialized_start=3113
  _globals['_MSGNEWICATRANSACTIONREQUESTRESPONSE']._serialized_end=3216
  _globals['_MSGFULFILICATRANSACTIONREQUEST']._serialized_start=3219
  _globals['_MSGFULFILICATRANSACTIONREQUEST']._serialized_end=3509
  _globals['_MSGFULFILICATRANSACTIONREQUESTRESPONSE']._serialized_start=3511
  _globals['_MSGFULFILICATRANSACTIONREQUESTRESPONSE']._serialized_end=3551
  _globals['_MSGNEWZRSIGNSIGNATUREREQUEST']._serialized_start=3554
  _globals['_MSGNEWZRSIGNSIGNATUREREQUEST']._serialized_end=4091
  _globals['_MSGNEWZRSIGNSIGNATUREREQUESTRESPONSE']._serialized_start=4093
  _globals['_MSGNEWZRSIGNSIGNATUREREQUESTRESPONSE']._serialized_end=4154
  _globals['_MSGUPDATEKEYPOLICY']._serialized_start=4156
  _globals['_MSGUPDATEKEYPOLICY']._serialized_end=4277
  _globals['_MSGUPDATEKEYPOLICYRESPONSE']._serialized_start=4279
  _globals['_MSGUPDATEKEYPOLICYRESPONSE']._serialized_end=4307
  _globals['_MSG']._serialized_start=4438
  _globals['_MSG']._serialized_end=5743
# @@protoc_insertion_point(module_scope)
