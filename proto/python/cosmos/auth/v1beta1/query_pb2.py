# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/auth/v1beta1/query.proto
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
    'cosmos/auth/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.auth.v1beta1 import auth_pb2 as cosmos_dot_auth_dot_v1beta1_dot_auth__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.query.v1 import query_pb2 as cosmos_dot_query_dot_v1_dot_query__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1f\x63osmos/auth/v1beta1/query.proto\x12\x13\x63osmos.auth.v1beta1\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1e\x63osmos/auth/v1beta1/auth.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1b\x63osmos/query/v1/query.proto\"^\n\x14QueryAccountsRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xb4\x01\n\x15QueryAccountsResponse\x12R\n\x08\x61\x63\x63ounts\x18\x01 \x03(\x0b\x32\x14.google.protobuf.AnyB \xca\xb4-\x1c\x63osmos.auth.v1beta1.AccountIR\x08\x61\x63\x63ounts\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"S\n\x13QueryAccountRequest\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"h\n\x14QueryAccountResponse\x12P\n\x07\x61\x63\x63ount\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyB \xca\xb4-\x1c\x63osmos.auth.v1beta1.AccountIR\x07\x61\x63\x63ount\"\x14\n\x12QueryParamsRequest\"P\n\x13QueryParamsResponse\x12\x39\n\x06params\x18\x01 \x01(\x0b\x32\x1b.cosmos.auth.v1beta1.ParamsB\x04\xc8\xde\x1f\x00R\x06params\"\x1c\n\x1aQueryModuleAccountsRequest\"w\n\x1bQueryModuleAccountsResponse\x12X\n\x08\x61\x63\x63ounts\x18\x01 \x03(\x0b\x32\x14.google.protobuf.AnyB&\xca\xb4-\"cosmos.auth.v1beta1.ModuleAccountIR\x08\x61\x63\x63ounts\"5\n\x1fQueryModuleAccountByNameRequest\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\"z\n QueryModuleAccountByNameResponse\x12V\n\x07\x61\x63\x63ount\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyB&\xca\xb4-\"cosmos.auth.v1beta1.ModuleAccountIR\x07\x61\x63\x63ount\"\x15\n\x13\x42\x65\x63h32PrefixRequest\";\n\x14\x42\x65\x63h32PrefixResponse\x12#\n\rbech32_prefix\x18\x01 \x01(\tR\x0c\x62\x65\x63h32Prefix\"B\n\x1b\x41\x64\x64ressBytesToStringRequest\x12#\n\raddress_bytes\x18\x01 \x01(\x0cR\x0c\x61\x64\x64ressBytes\"E\n\x1c\x41\x64\x64ressBytesToStringResponse\x12%\n\x0e\x61\x64\x64ress_string\x18\x01 \x01(\tR\raddressString\"D\n\x1b\x41\x64\x64ressStringToBytesRequest\x12%\n\x0e\x61\x64\x64ress_string\x18\x01 \x01(\tR\raddressString\"C\n\x1c\x41\x64\x64ressStringToBytesResponse\x12#\n\raddress_bytes\x18\x01 \x01(\x0cR\x0c\x61\x64\x64ressBytes\"S\n\x1eQueryAccountAddressByIDRequest\x12\x12\n\x02id\x18\x01 \x01(\x03\x42\x02\x18\x01R\x02id\x12\x1d\n\naccount_id\x18\x02 \x01(\x04R\taccountId\"d\n\x1fQueryAccountAddressByIDResponse\x12\x41\n\x0f\x61\x63\x63ount_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x0e\x61\x63\x63ountAddress\"M\n\x17QueryAccountInfoRequest\x12\x32\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x07\x61\x64\x64ress\"P\n\x18QueryAccountInfoResponse\x12\x34\n\x04info\x18\x01 \x01(\x0b\x32 .cosmos.auth.v1beta1.BaseAccountR\x04info2\xef\x0c\n\x05Query\x12\x8d\x01\n\x08\x41\x63\x63ounts\x12).cosmos.auth.v1beta1.QueryAccountsRequest\x1a*.cosmos.auth.v1beta1.QueryAccountsResponse\"*\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02\x1f\x12\x1d/cosmos/auth/v1beta1/accounts\x12\x94\x01\n\x07\x41\x63\x63ount\x12(.cosmos.auth.v1beta1.QueryAccountRequest\x1a).cosmos.auth.v1beta1.QueryAccountResponse\"4\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02)\x12\'/cosmos/auth/v1beta1/accounts/{address}\x12\xb5\x01\n\x12\x41\x63\x63ountAddressByID\x12\x33.cosmos.auth.v1beta1.QueryAccountAddressByIDRequest\x1a\x34.cosmos.auth.v1beta1.QueryAccountAddressByIDResponse\"4\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02)\x12\'/cosmos/auth/v1beta1/address_by_id/{id}\x12\x85\x01\n\x06Params\x12\'.cosmos.auth.v1beta1.QueryParamsRequest\x1a(.cosmos.auth.v1beta1.QueryParamsResponse\"(\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02\x1d\x12\x1b/cosmos/auth/v1beta1/params\x12\xa6\x01\n\x0eModuleAccounts\x12/.cosmos.auth.v1beta1.QueryModuleAccountsRequest\x1a\x30.cosmos.auth.v1beta1.QueryModuleAccountsResponse\"1\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02&\x12$/cosmos/auth/v1beta1/module_accounts\x12\xbc\x01\n\x13ModuleAccountByName\x12\x34.cosmos.auth.v1beta1.QueryModuleAccountByNameRequest\x1a\x35.cosmos.auth.v1beta1.QueryModuleAccountByNameResponse\"8\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02-\x12+/cosmos/auth/v1beta1/module_accounts/{name}\x12\x88\x01\n\x0c\x42\x65\x63h32Prefix\x12(.cosmos.auth.v1beta1.Bech32PrefixRequest\x1a).cosmos.auth.v1beta1.Bech32PrefixResponse\"#\x82\xd3\xe4\x93\x02\x1d\x12\x1b/cosmos/auth/v1beta1/bech32\x12\xb0\x01\n\x14\x41\x64\x64ressBytesToString\x12\x30.cosmos.auth.v1beta1.AddressBytesToStringRequest\x1a\x31.cosmos.auth.v1beta1.AddressBytesToStringResponse\"3\x82\xd3\xe4\x93\x02-\x12+/cosmos/auth/v1beta1/bech32/{address_bytes}\x12\xb1\x01\n\x14\x41\x64\x64ressStringToBytes\x12\x30.cosmos.auth.v1beta1.AddressStringToBytesRequest\x1a\x31.cosmos.auth.v1beta1.AddressStringToBytesResponse\"4\x82\xd3\xe4\x93\x02.\x12,/cosmos/auth/v1beta1/bech32/{address_string}\x12\xa4\x01\n\x0b\x41\x63\x63ountInfo\x12,.cosmos.auth.v1beta1.QueryAccountInfoRequest\x1a-.cosmos.auth.v1beta1.QueryAccountInfoResponse\"8\x88\xe7\xb0*\x01\x82\xd3\xe4\x93\x02-\x12+/cosmos/auth/v1beta1/account_info/{address}B\x1bZ\x19\x63osmossdk.io/x/auth/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.auth.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\031cosmossdk.io/x/auth/types'
  _globals['_QUERYACCOUNTSRESPONSE'].fields_by_name['accounts']._loaded_options = None
  _globals['_QUERYACCOUNTSRESPONSE'].fields_by_name['accounts']._serialized_options = b'\312\264-\034cosmos.auth.v1beta1.AccountI'
  _globals['_QUERYACCOUNTREQUEST'].fields_by_name['address']._loaded_options = None
  _globals['_QUERYACCOUNTREQUEST'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYACCOUNTREQUEST']._loaded_options = None
  _globals['_QUERYACCOUNTREQUEST']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_QUERYACCOUNTRESPONSE'].fields_by_name['account']._loaded_options = None
  _globals['_QUERYACCOUNTRESPONSE'].fields_by_name['account']._serialized_options = b'\312\264-\034cosmos.auth.v1beta1.AccountI'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._serialized_options = b'\310\336\037\000'
  _globals['_QUERYMODULEACCOUNTSRESPONSE'].fields_by_name['accounts']._loaded_options = None
  _globals['_QUERYMODULEACCOUNTSRESPONSE'].fields_by_name['accounts']._serialized_options = b'\312\264-\"cosmos.auth.v1beta1.ModuleAccountI'
  _globals['_QUERYMODULEACCOUNTBYNAMERESPONSE'].fields_by_name['account']._loaded_options = None
  _globals['_QUERYMODULEACCOUNTBYNAMERESPONSE'].fields_by_name['account']._serialized_options = b'\312\264-\"cosmos.auth.v1beta1.ModuleAccountI'
  _globals['_QUERYACCOUNTADDRESSBYIDREQUEST'].fields_by_name['id']._loaded_options = None
  _globals['_QUERYACCOUNTADDRESSBYIDREQUEST'].fields_by_name['id']._serialized_options = b'\030\001'
  _globals['_QUERYACCOUNTADDRESSBYIDRESPONSE'].fields_by_name['account_address']._loaded_options = None
  _globals['_QUERYACCOUNTADDRESSBYIDRESPONSE'].fields_by_name['account_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYACCOUNTINFOREQUEST'].fields_by_name['address']._loaded_options = None
  _globals['_QUERYACCOUNTINFOREQUEST'].fields_by_name['address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERY'].methods_by_name['Accounts']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Accounts']._serialized_options = b'\210\347\260*\001\202\323\344\223\002\037\022\035/cosmos/auth/v1beta1/accounts'
  _globals['_QUERY'].methods_by_name['Account']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Account']._serialized_options = b'\210\347\260*\001\202\323\344\223\002)\022\'/cosmos/auth/v1beta1/accounts/{address}'
  _globals['_QUERY'].methods_by_name['AccountAddressByID']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AccountAddressByID']._serialized_options = b'\210\347\260*\001\202\323\344\223\002)\022\'/cosmos/auth/v1beta1/address_by_id/{id}'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\210\347\260*\001\202\323\344\223\002\035\022\033/cosmos/auth/v1beta1/params'
  _globals['_QUERY'].methods_by_name['ModuleAccounts']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ModuleAccounts']._serialized_options = b'\210\347\260*\001\202\323\344\223\002&\022$/cosmos/auth/v1beta1/module_accounts'
  _globals['_QUERY'].methods_by_name['ModuleAccountByName']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ModuleAccountByName']._serialized_options = b'\210\347\260*\001\202\323\344\223\002-\022+/cosmos/auth/v1beta1/module_accounts/{name}'
  _globals['_QUERY'].methods_by_name['Bech32Prefix']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Bech32Prefix']._serialized_options = b'\202\323\344\223\002\035\022\033/cosmos/auth/v1beta1/bech32'
  _globals['_QUERY'].methods_by_name['AddressBytesToString']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AddressBytesToString']._serialized_options = b'\202\323\344\223\002-\022+/cosmos/auth/v1beta1/bech32/{address_bytes}'
  _globals['_QUERY'].methods_by_name['AddressStringToBytes']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AddressStringToBytes']._serialized_options = b'\202\323\344\223\002.\022,/cosmos/auth/v1beta1/bech32/{address_string}'
  _globals['_QUERY'].methods_by_name['AccountInfo']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AccountInfo']._serialized_options = b'\210\347\260*\001\202\323\344\223\002-\022+/cosmos/auth/v1beta1/account_info/{address}'
  _globals['_QUERYACCOUNTSREQUEST']._serialized_start=267
  _globals['_QUERYACCOUNTSREQUEST']._serialized_end=361
  _globals['_QUERYACCOUNTSRESPONSE']._serialized_start=364
  _globals['_QUERYACCOUNTSRESPONSE']._serialized_end=544
  _globals['_QUERYACCOUNTREQUEST']._serialized_start=546
  _globals['_QUERYACCOUNTREQUEST']._serialized_end=629
  _globals['_QUERYACCOUNTRESPONSE']._serialized_start=631
  _globals['_QUERYACCOUNTRESPONSE']._serialized_end=735
  _globals['_QUERYPARAMSREQUEST']._serialized_start=737
  _globals['_QUERYPARAMSREQUEST']._serialized_end=757
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=759
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=839
  _globals['_QUERYMODULEACCOUNTSREQUEST']._serialized_start=841
  _globals['_QUERYMODULEACCOUNTSREQUEST']._serialized_end=869
  _globals['_QUERYMODULEACCOUNTSRESPONSE']._serialized_start=871
  _globals['_QUERYMODULEACCOUNTSRESPONSE']._serialized_end=990
  _globals['_QUERYMODULEACCOUNTBYNAMEREQUEST']._serialized_start=992
  _globals['_QUERYMODULEACCOUNTBYNAMEREQUEST']._serialized_end=1045
  _globals['_QUERYMODULEACCOUNTBYNAMERESPONSE']._serialized_start=1047
  _globals['_QUERYMODULEACCOUNTBYNAMERESPONSE']._serialized_end=1169
  _globals['_BECH32PREFIXREQUEST']._serialized_start=1171
  _globals['_BECH32PREFIXREQUEST']._serialized_end=1192
  _globals['_BECH32PREFIXRESPONSE']._serialized_start=1194
  _globals['_BECH32PREFIXRESPONSE']._serialized_end=1253
  _globals['_ADDRESSBYTESTOSTRINGREQUEST']._serialized_start=1255
  _globals['_ADDRESSBYTESTOSTRINGREQUEST']._serialized_end=1321
  _globals['_ADDRESSBYTESTOSTRINGRESPONSE']._serialized_start=1323
  _globals['_ADDRESSBYTESTOSTRINGRESPONSE']._serialized_end=1392
  _globals['_ADDRESSSTRINGTOBYTESREQUEST']._serialized_start=1394
  _globals['_ADDRESSSTRINGTOBYTESREQUEST']._serialized_end=1462
  _globals['_ADDRESSSTRINGTOBYTESRESPONSE']._serialized_start=1464
  _globals['_ADDRESSSTRINGTOBYTESRESPONSE']._serialized_end=1531
  _globals['_QUERYACCOUNTADDRESSBYIDREQUEST']._serialized_start=1533
  _globals['_QUERYACCOUNTADDRESSBYIDREQUEST']._serialized_end=1616
  _globals['_QUERYACCOUNTADDRESSBYIDRESPONSE']._serialized_start=1618
  _globals['_QUERYACCOUNTADDRESSBYIDRESPONSE']._serialized_end=1718
  _globals['_QUERYACCOUNTINFOREQUEST']._serialized_start=1720
  _globals['_QUERYACCOUNTINFOREQUEST']._serialized_end=1797
  _globals['_QUERYACCOUNTINFORESPONSE']._serialized_start=1799
  _globals['_QUERYACCOUNTINFORESPONSE']._serialized_end=1879
  _globals['_QUERY']._serialized_start=1882
  _globals['_QUERY']._serialized_end=3529
# @@protoc_insertion_point(module_scope)
