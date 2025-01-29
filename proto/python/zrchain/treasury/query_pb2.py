# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/treasury/query.proto
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
    'zrchain/treasury/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from zrchain.treasury import key_pb2 as zrchain_dot_treasury_dot_key__pb2
from zrchain.treasury import mpcsign_pb2 as zrchain_dot_treasury_dot_mpcsign__pb2
from zrchain.treasury import params_pb2 as zrchain_dot_treasury_dot_params__pb2
from zrchain.treasury import wallet_pb2 as zrchain_dot_treasury_dot_wallet__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1czrchain/treasury/query.proto\x12\x10zrchain.treasury\x1a\x11\x61mino/amino.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1azrchain/treasury/key.proto\x1a\x1ezrchain/treasury/mpcsign.proto\x1a\x1dzrchain/treasury/params.proto\x1a\x1dzrchain/treasury/wallet.proto\"\x14\n\x12QueryParamsRequest\"R\n\x13QueryParamsResponse\x12;\n\x06params\x18\x01 \x01(\x0b\x32\x18.zrchain.treasury.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\"\xe7\x01\n\x17QueryKeyRequestsRequest\x12!\n\x0ckeyring_addr\x18\x01 \x01(\tR\x0bkeyringAddr\x12:\n\x06status\x18\x02 \x01(\x0e\x32\".zrchain.treasury.KeyRequestStatusR\x06status\x12%\n\x0eworkspace_addr\x18\x03 \x01(\tR\rworkspaceAddr\x12\x46\n\npagination\x18\x04 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xa8\x01\n\x18QueryKeyRequestsResponse\x12\x43\n\x0ckey_requests\x18\x01 \x03(\x0b\x32 .zrchain.treasury.KeyReqResponseR\x0bkeyRequests\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"\x80\x01\n\x13QueryKeyByIDRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12=\n\x0bwallet_type\x18\x02 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\nwalletType\x12\x1a\n\x08prefixes\x18\x03 \x03(\tR\x08prefixes\"\x83\x01\n\x14QueryKeyByIDResponse\x12/\n\x03key\x18\x01 \x01(\x0b\x32\x1d.zrchain.treasury.KeyResponseR\x03key\x12:\n\x07wallets\x18\x02 \x03(\x0b\x32 .zrchain.treasury.WalletResponseR\x07wallets\"\xdc\x01\n\x10QueryKeysRequest\x12%\n\x0eworkspace_addr\x18\x01 \x01(\tR\rworkspaceAddr\x12=\n\x0bwallet_type\x18\x02 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\nwalletType\x12\x1a\n\x08prefixes\x18\x03 \x03(\tR\x08prefixes\x12\x46\n\npagination\x18\x04 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x98\x01\n\x11QueryKeysResponse\x12:\n\x04keys\x18\x01 \x03(\x0b\x32&.zrchain.treasury.KeyAndWalletResponseR\x04keys\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"\x83\x01\n\x14KeyAndWalletResponse\x12/\n\x03key\x18\x01 \x01(\x0b\x32\x1d.zrchain.treasury.KeyResponseR\x03key\x12:\n\x07wallets\x18\x02 \x03(\x0b\x32 .zrchain.treasury.WalletResponseR\x07wallets\">\n\x0eWalletResponse\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12\x12\n\x04type\x18\x02 \x01(\tR\x04type\",\n\x1aQueryKeyRequestByIDRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\"`\n\x1bQueryKeyRequestByIDResponse\x12\x41\n\x0bkey_request\x18\x01 \x01(\x0b\x32 .zrchain.treasury.KeyReqResponseR\nkeyRequest\"\xc7\x01\n\x1dQuerySignatureRequestsRequest\x12!\n\x0ckeyring_addr\x18\x01 \x01(\tR\x0bkeyringAddr\x12;\n\x06status\x18\x02 \x01(\x0e\x32#.zrchain.treasury.SignRequestStatusR\x06status\x12\x46\n\npagination\x18\x03 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xb1\x01\n\x1eQuerySignatureRequestsResponse\x12\x46\n\rsign_requests\x18\x01 \x03(\x0b\x32!.zrchain.treasury.SignReqResponseR\x0csignRequests\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"2\n QuerySignatureRequestByIDRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\"i\n!QuerySignatureRequestByIDResponse\x12\x44\n\x0csign_request\x18\x01 \x01(\x0b\x32!.zrchain.treasury.SignReqResponseR\x0bsignRequest\"\x9f\x02\n#QuerySignTransactionRequestsRequest\x12\x1d\n\nrequest_id\x18\x01 \x01(\x04R\trequestId\x12\x15\n\x06key_id\x18\x02 \x01(\x04R\x05keyId\x12=\n\x0bwallet_type\x18\x03 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\nwalletType\x12;\n\x06status\x18\x04 \x01(\x0e\x32#.zrchain.treasury.SignRequestStatusR\x06status\x12\x46\n\npagination\x18\x05 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xde\x01\n$QuerySignTransactionRequestsResponse\x12m\n\x19sign_transaction_requests\x18\x01 \x03(\x0b\x32\x31.zrchain.treasury.SignTransactionRequestsResponseR\x17signTransactionRequests\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"\xc8\x01\n\x1fSignTransactionRequestsResponse\x12_\n\x19sign_transaction_requests\x18\x01 \x01(\x0b\x32#.zrchain.treasury.SignTxReqResponseR\x17signTransactionRequests\x12\x44\n\x0csign_request\x18\x02 \x01(\x0b\x32!.zrchain.treasury.SignReqResponseR\x0bsignRequest\"8\n&QuerySignTransactionRequestByIDRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\"\x88\x01\n\'QuerySignTransactionRequestByIDResponse\x12]\n\x18sign_transaction_request\x18\x01 \x01(\x0b\x32#.zrchain.treasury.SignTxReqResponseR\x16signTransactionRequest\"\x9a\x01\n\x16QueryZrSignKeysRequest\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12\x1e\n\nwalletType\x18\x02 \x01(\tR\nwalletType\x12\x46\n\npagination\x18\x03 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"p\n\x0eZrSignKeyEntry\x12\x1e\n\nwalletType\x18\x01 \x01(\tR\nwalletType\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ress\x12\x14\n\x05index\x18\x03 \x01(\x04R\x05index\x12\x0e\n\x02id\x18\x04 \x01(\x04R\x02id\"\x98\x01\n\x17QueryZrSignKeysResponse\x12\x34\n\x04Keys\x18\x01 \x03(\x0b\x32 .zrchain.treasury.ZrSignKeyEntryR\x04Keys\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"\xe8\x01\n\x18QueryKeyByAddressRequest\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12!\n\x0ckeyring_addr\x18\x02 \x01(\tR\x0bkeyringAddr\x12\x34\n\x08key_type\x18\x03 \x01(\x0e\x32\x19.zrchain.treasury.KeyTypeR\x07keyType\x12=\n\x0bwallet_type\x18\x04 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\nwalletType\x12\x1a\n\x08prefixes\x18\x05 \x03(\tR\x08prefixes\"_\n\x19QueryKeyByAddressResponse\x12\x42\n\x08response\x18\x01 \x01(\x0b\x32&.zrchain.treasury.KeyAndWalletResponseR\x08response\"\x8c\x02\n\x19QueryZenbtcWalletsRequest\x12%\n\x0erecipient_addr\x18\x01 \x01(\tR\rrecipientAddr\x12;\n\nchain_type\x18\x02 \x01(\x0e\x32\x1c.zrchain.treasury.WalletTypeR\tchainType\x12\"\n\rmint_chain_id\x18\x03 \x01(\x04R\x0bmintChainId\x12\x1f\n\x0breturn_addr\x18\x04 \x01(\tR\nreturnAddr\x12\x46\n\npagination\x18\x05 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xb4\x01\n\x1aQueryZenbtcWalletsResponse\x12M\n\x0ezenbtc_wallets\x18\x01 \x03(\x0b\x32&.zrchain.treasury.KeyAndWalletResponseR\rzenbtcWallets\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination2\xcf\x10\n\x05Query\x12w\n\x06Params\x12$.zrchain.treasury.QueryParamsRequest\x1a%.zrchain.treasury.QueryParamsResponse\" \x82\xd3\xe4\x93\x02\x1a\x12\x18/zrchain/treasury/params\x12\xb5\x01\n\x0bKeyRequests\x12).zrchain.treasury.QueryKeyRequestsRequest\x1a*.zrchain.treasury.QueryKeyRequestsResponse\"O\x82\xd3\xe4\x93\x02I\x12G/zrchain/treasury/key_requests/{keyring_addr}/{status}/{workspace_addr}\x12\x9f\x01\n\x0eKeyRequestByID\x12,.zrchain.treasury.QueryKeyRequestByIDRequest\x1a-.zrchain.treasury.QueryKeyRequestByIDResponse\"0\x82\xd3\xe4\x93\x02*\x12(/zrchain/treasury/key_request_by_id/{id}\x12\x99\x01\n\x04Keys\x12\".zrchain.treasury.QueryKeysRequest\x1a#.zrchain.treasury.QueryKeysResponse\"H\x82\xd3\xe4\x93\x02\x42\x12@/zrchain/treasury/keys/{workspace_addr}/{wallet_type}/{prefixes}\x12\x9b\x01\n\x07KeyByID\x12%.zrchain.treasury.QueryKeyByIDRequest\x1a&.zrchain.treasury.QueryKeyByIDResponse\"A\x82\xd3\xe4\x93\x02;\x12\x39/zrchain/treasury/key_by_id/{id}/{wallet_type}/{prefixes}\x12\xbc\x01\n\x11SignatureRequests\x12/.zrchain.treasury.QuerySignatureRequestsRequest\x1a\x30.zrchain.treasury.QuerySignatureRequestsResponse\"D\x82\xd3\xe4\x93\x02>\x12</zrchain/treasury/signature_requests/{keyring_addr}/{status}\x12\xb7\x01\n\x14SignatureRequestByID\x12\x32.zrchain.treasury.QuerySignatureRequestByIDRequest\x1a\x33.zrchain.treasury.QuerySignatureRequestByIDResponse\"6\x82\xd3\xe4\x93\x02\x30\x12./zrchain/treasury/signature_request_by_id/{id}\x12\xdc\x01\n\x17SignTransactionRequests\x12\x35.zrchain.treasury.QuerySignTransactionRequestsRequest\x1a\x36.zrchain.treasury.QuerySignTransactionRequestsResponse\"R\x82\xd3\xe4\x93\x02L\x12J/zrchain/treasury/sign_transaction_request/{wallet_type}/{key_id}/{status}\x12\xd0\x01\n\x1aSignTransactionRequestByID\x12\x38.zrchain.treasury.QuerySignTransactionRequestByIDRequest\x1a\x39.zrchain.treasury.QuerySignTransactionRequestByIDResponse\"=\x82\xd3\xe4\x93\x02\x37\x12\x35/zrchain/treasury/sign_transaction_request_by_id/{id}\x12\xa0\x01\n\nZrSignKeys\x12(.zrchain.treasury.QueryZrSignKeysRequest\x1a).zrchain.treasury.QueryZrSignKeysResponse\"=\x82\xd3\xe4\x93\x02\x37\x12\x35/zrchain/treasury/zr_sign_keys/{address}/{walletType}\x12\x94\x01\n\x0cKeyByAddress\x12*.zrchain.treasury.QueryKeyByAddressRequest\x1a+.zrchain.treasury.QueryKeyByAddressResponse\"+\x82\xd3\xe4\x93\x02%\x12#/zrchain/v4/treasury/key_by_address\x12\xd3\x01\n\rZenbtcWallets\x12+.zrchain.treasury.QueryZenbtcWalletsRequest\x1a,.zrchain.treasury.QueryZenbtcWalletsResponse\"g\x82\xd3\xe4\x93\x02\x61\x12_/zrchain/v5/treasury/zenbtc_wallets/{recipient_addr}/{chain_type}/{mint_chain_id}/{return_addr}B;Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.treasury.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z9github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\032\022\030/zrchain/treasury/params'
  _globals['_QUERY'].methods_by_name['KeyRequests']._loaded_options = None
  _globals['_QUERY'].methods_by_name['KeyRequests']._serialized_options = b'\202\323\344\223\002I\022G/zrchain/treasury/key_requests/{keyring_addr}/{status}/{workspace_addr}'
  _globals['_QUERY'].methods_by_name['KeyRequestByID']._loaded_options = None
  _globals['_QUERY'].methods_by_name['KeyRequestByID']._serialized_options = b'\202\323\344\223\002*\022(/zrchain/treasury/key_request_by_id/{id}'
  _globals['_QUERY'].methods_by_name['Keys']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Keys']._serialized_options = b'\202\323\344\223\002B\022@/zrchain/treasury/keys/{workspace_addr}/{wallet_type}/{prefixes}'
  _globals['_QUERY'].methods_by_name['KeyByID']._loaded_options = None
  _globals['_QUERY'].methods_by_name['KeyByID']._serialized_options = b'\202\323\344\223\002;\0229/zrchain/treasury/key_by_id/{id}/{wallet_type}/{prefixes}'
  _globals['_QUERY'].methods_by_name['SignatureRequests']._loaded_options = None
  _globals['_QUERY'].methods_by_name['SignatureRequests']._serialized_options = b'\202\323\344\223\002>\022</zrchain/treasury/signature_requests/{keyring_addr}/{status}'
  _globals['_QUERY'].methods_by_name['SignatureRequestByID']._loaded_options = None
  _globals['_QUERY'].methods_by_name['SignatureRequestByID']._serialized_options = b'\202\323\344\223\0020\022./zrchain/treasury/signature_request_by_id/{id}'
  _globals['_QUERY'].methods_by_name['SignTransactionRequests']._loaded_options = None
  _globals['_QUERY'].methods_by_name['SignTransactionRequests']._serialized_options = b'\202\323\344\223\002L\022J/zrchain/treasury/sign_transaction_request/{wallet_type}/{key_id}/{status}'
  _globals['_QUERY'].methods_by_name['SignTransactionRequestByID']._loaded_options = None
  _globals['_QUERY'].methods_by_name['SignTransactionRequestByID']._serialized_options = b'\202\323\344\223\0027\0225/zrchain/treasury/sign_transaction_request_by_id/{id}'
  _globals['_QUERY'].methods_by_name['ZrSignKeys']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ZrSignKeys']._serialized_options = b'\202\323\344\223\0027\0225/zrchain/treasury/zr_sign_keys/{address}/{walletType}'
  _globals['_QUERY'].methods_by_name['KeyByAddress']._loaded_options = None
  _globals['_QUERY'].methods_by_name['KeyByAddress']._serialized_options = b'\202\323\344\223\002%\022#/zrchain/v4/treasury/key_by_address'
  _globals['_QUERY'].methods_by_name['ZenbtcWallets']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ZenbtcWallets']._serialized_options = b'\202\323\344\223\002a\022_/zrchain/v5/treasury/zenbtc_wallets/{recipient_addr}/{chain_type}/{mint_chain_id}/{return_addr}'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=287
  _globals['_QUERYPARAMSREQUEST']._serialized_end=307
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=309
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=391
  _globals['_QUERYKEYREQUESTSREQUEST']._serialized_start=394
  _globals['_QUERYKEYREQUESTSREQUEST']._serialized_end=625
  _globals['_QUERYKEYREQUESTSRESPONSE']._serialized_start=628
  _globals['_QUERYKEYREQUESTSRESPONSE']._serialized_end=796
  _globals['_QUERYKEYBYIDREQUEST']._serialized_start=799
  _globals['_QUERYKEYBYIDREQUEST']._serialized_end=927
  _globals['_QUERYKEYBYIDRESPONSE']._serialized_start=930
  _globals['_QUERYKEYBYIDRESPONSE']._serialized_end=1061
  _globals['_QUERYKEYSREQUEST']._serialized_start=1064
  _globals['_QUERYKEYSREQUEST']._serialized_end=1284
  _globals['_QUERYKEYSRESPONSE']._serialized_start=1287
  _globals['_QUERYKEYSRESPONSE']._serialized_end=1439
  _globals['_KEYANDWALLETRESPONSE']._serialized_start=1442
  _globals['_KEYANDWALLETRESPONSE']._serialized_end=1573
  _globals['_WALLETRESPONSE']._serialized_start=1575
  _globals['_WALLETRESPONSE']._serialized_end=1637
  _globals['_QUERYKEYREQUESTBYIDREQUEST']._serialized_start=1639
  _globals['_QUERYKEYREQUESTBYIDREQUEST']._serialized_end=1683
  _globals['_QUERYKEYREQUESTBYIDRESPONSE']._serialized_start=1685
  _globals['_QUERYKEYREQUESTBYIDRESPONSE']._serialized_end=1781
  _globals['_QUERYSIGNATUREREQUESTSREQUEST']._serialized_start=1784
  _globals['_QUERYSIGNATUREREQUESTSREQUEST']._serialized_end=1983
  _globals['_QUERYSIGNATUREREQUESTSRESPONSE']._serialized_start=1986
  _globals['_QUERYSIGNATUREREQUESTSRESPONSE']._serialized_end=2163
  _globals['_QUERYSIGNATUREREQUESTBYIDREQUEST']._serialized_start=2165
  _globals['_QUERYSIGNATUREREQUESTBYIDREQUEST']._serialized_end=2215
  _globals['_QUERYSIGNATUREREQUESTBYIDRESPONSE']._serialized_start=2217
  _globals['_QUERYSIGNATUREREQUESTBYIDRESPONSE']._serialized_end=2322
  _globals['_QUERYSIGNTRANSACTIONREQUESTSREQUEST']._serialized_start=2325
  _globals['_QUERYSIGNTRANSACTIONREQUESTSREQUEST']._serialized_end=2612
  _globals['_QUERYSIGNTRANSACTIONREQUESTSRESPONSE']._serialized_start=2615
  _globals['_QUERYSIGNTRANSACTIONREQUESTSRESPONSE']._serialized_end=2837
  _globals['_SIGNTRANSACTIONREQUESTSRESPONSE']._serialized_start=2840
  _globals['_SIGNTRANSACTIONREQUESTSRESPONSE']._serialized_end=3040
  _globals['_QUERYSIGNTRANSACTIONREQUESTBYIDREQUEST']._serialized_start=3042
  _globals['_QUERYSIGNTRANSACTIONREQUESTBYIDREQUEST']._serialized_end=3098
  _globals['_QUERYSIGNTRANSACTIONREQUESTBYIDRESPONSE']._serialized_start=3101
  _globals['_QUERYSIGNTRANSACTIONREQUESTBYIDRESPONSE']._serialized_end=3237
  _globals['_QUERYZRSIGNKEYSREQUEST']._serialized_start=3240
  _globals['_QUERYZRSIGNKEYSREQUEST']._serialized_end=3394
  _globals['_ZRSIGNKEYENTRY']._serialized_start=3396
  _globals['_ZRSIGNKEYENTRY']._serialized_end=3508
  _globals['_QUERYZRSIGNKEYSRESPONSE']._serialized_start=3511
  _globals['_QUERYZRSIGNKEYSRESPONSE']._serialized_end=3663
  _globals['_QUERYKEYBYADDRESSREQUEST']._serialized_start=3666
  _globals['_QUERYKEYBYADDRESSREQUEST']._serialized_end=3898
  _globals['_QUERYKEYBYADDRESSRESPONSE']._serialized_start=3900
  _globals['_QUERYKEYBYADDRESSRESPONSE']._serialized_end=3995
  _globals['_QUERYZENBTCWALLETSREQUEST']._serialized_start=3998
  _globals['_QUERYZENBTCWALLETSREQUEST']._serialized_end=4266
  _globals['_QUERYZENBTCWALLETSRESPONSE']._serialized_start=4269
  _globals['_QUERYZENBTCWALLETSRESPONSE']._serialized_end=4449
  _globals['_QUERY']._serialized_start=4452
  _globals['_QUERY']._serialized_end=6579
# @@protoc_insertion_point(module_scope)
