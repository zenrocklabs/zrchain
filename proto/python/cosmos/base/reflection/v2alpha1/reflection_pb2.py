# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/base/reflection/v2alpha1/reflection.proto
<<<<<<< HEAD
# Protobuf Python Version: 6.30.1
=======
# Protobuf Python Version: 6.30.0
>>>>>>> main
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
<<<<<<< HEAD
    1,
=======
    0,
>>>>>>> main
    '',
    'cosmos/base/reflection/v2alpha1/reflection.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n0cosmos/base/reflection/v2alpha1/reflection.proto\x12\x1f\x63osmos.base.reflection.v2alpha1\x1a\x1cgoogle/api/annotations.proto\"\xe7\x03\n\rAppDescriptor\x12\x46\n\x05\x61uthn\x18\x01 \x01(\x0b\x32\x30.cosmos.base.reflection.v2alpha1.AuthnDescriptorR\x05\x61uthn\x12\x46\n\x05\x63hain\x18\x02 \x01(\x0b\x32\x30.cosmos.base.reflection.v2alpha1.ChainDescriptorR\x05\x63hain\x12\x46\n\x05\x63odec\x18\x03 \x01(\x0b\x32\x30.cosmos.base.reflection.v2alpha1.CodecDescriptorR\x05\x63odec\x12^\n\rconfiguration\x18\x04 \x01(\x0b\x32\x38.cosmos.base.reflection.v2alpha1.ConfigurationDescriptorR\rconfiguration\x12_\n\x0equery_services\x18\x05 \x01(\x0b\x32\x38.cosmos.base.reflection.v2alpha1.QueryServicesDescriptorR\rqueryServices\x12=\n\x02tx\x18\x06 \x01(\x0b\x32-.cosmos.base.reflection.v2alpha1.TxDescriptorR\x02tx\"n\n\x0cTxDescriptor\x12\x1a\n\x08\x66ullname\x18\x01 \x01(\tR\x08\x66ullname\x12\x42\n\x04msgs\x18\x02 \x03(\x0b\x32..cosmos.base.reflection.v2alpha1.MsgDescriptorR\x04msgs\"h\n\x0f\x41uthnDescriptor\x12U\n\nsign_modes\x18\x01 \x03(\x0b\x32\x36.cosmos.base.reflection.v2alpha1.SigningModeDescriptorR\tsignModes\"\x91\x01\n\x15SigningModeDescriptor\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x16\n\x06number\x18\x02 \x01(\x05R\x06number\x12L\n#authn_info_provider_method_fullname\x18\x03 \x01(\tR\x1f\x61uthnInfoProviderMethodFullname\"!\n\x0f\x43hainDescriptor\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\"g\n\x0f\x43odecDescriptor\x12T\n\ninterfaces\x18\x01 \x03(\x0b\x32\x34.cosmos.base.reflection.v2alpha1.InterfaceDescriptorR\ninterfaces\"\xb2\x02\n\x13InterfaceDescriptor\x12\x1a\n\x08\x66ullname\x18\x01 \x01(\tR\x08\x66ullname\x12\x86\x01\n\x1cinterface_accepting_messages\x18\x02 \x03(\x0b\x32\x44.cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptorR\x1ainterfaceAcceptingMessages\x12v\n\x16interface_implementers\x18\x03 \x03(\x0b\x32?.cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptorR\x15interfaceImplementers\"W\n\x1eInterfaceImplementerDescriptor\x12\x1a\n\x08\x66ullname\x18\x01 \x01(\tR\x08\x66ullname\x12\x19\n\x08type_url\x18\x02 \x01(\tR\x07typeUrl\"w\n#InterfaceAcceptingMessageDescriptor\x12\x1a\n\x08\x66ullname\x18\x01 \x01(\tR\x08\x66ullname\x12\x34\n\x16\x66ield_descriptor_names\x18\x02 \x03(\tR\x14\x66ieldDescriptorNames\"\\\n\x17\x43onfigurationDescriptor\x12\x41\n\x1d\x62\x65\x63h32_account_address_prefix\x18\x01 \x01(\tR\x1a\x62\x65\x63h32AccountAddressPrefix\"1\n\rMsgDescriptor\x12 \n\x0cmsg_type_url\x18\x01 \x01(\tR\nmsgTypeUrl\"\x1b\n\x19GetAuthnDescriptorRequest\"d\n\x1aGetAuthnDescriptorResponse\x12\x46\n\x05\x61uthn\x18\x01 \x01(\x0b\x32\x30.cosmos.base.reflection.v2alpha1.AuthnDescriptorR\x05\x61uthn\"\x1b\n\x19GetChainDescriptorRequest\"d\n\x1aGetChainDescriptorResponse\x12\x46\n\x05\x63hain\x18\x01 \x01(\x0b\x32\x30.cosmos.base.reflection.v2alpha1.ChainDescriptorR\x05\x63hain\"\x1b\n\x19GetCodecDescriptorRequest\"d\n\x1aGetCodecDescriptorResponse\x12\x46\n\x05\x63odec\x18\x01 \x01(\x0b\x32\x30.cosmos.base.reflection.v2alpha1.CodecDescriptorR\x05\x63odec\"#\n!GetConfigurationDescriptorRequest\"v\n\"GetConfigurationDescriptorResponse\x12P\n\x06\x63onfig\x18\x01 \x01(\x0b\x32\x38.cosmos.base.reflection.v2alpha1.ConfigurationDescriptorR\x06\x63onfig\"#\n!GetQueryServicesDescriptorRequest\"x\n\"GetQueryServicesDescriptorResponse\x12R\n\x07queries\x18\x01 \x01(\x0b\x32\x38.cosmos.base.reflection.v2alpha1.QueryServicesDescriptorR\x07queries\"\x18\n\x16GetTxDescriptorRequest\"X\n\x17GetTxDescriptorResponse\x12=\n\x02tx\x18\x01 \x01(\x0b\x32-.cosmos.base.reflection.v2alpha1.TxDescriptorR\x02tx\"y\n\x17QueryServicesDescriptor\x12^\n\x0equery_services\x18\x01 \x03(\x0b\x32\x37.cosmos.base.reflection.v2alpha1.QueryServiceDescriptorR\rqueryServices\"\xa3\x01\n\x16QueryServiceDescriptor\x12\x1a\n\x08\x66ullname\x18\x01 \x01(\tR\x08\x66ullname\x12\x1b\n\tis_module\x18\x02 \x01(\x08R\x08isModule\x12P\n\x07methods\x18\x03 \x03(\x0b\x32\x36.cosmos.base.reflection.v2alpha1.QueryMethodDescriptorR\x07methods\"S\n\x15QueryMethodDescriptor\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12&\n\x0f\x66ull_query_path\x18\x02 \x01(\tR\rfullQueryPath2\xa7\n\n\x11ReflectionService\x12\xcb\x01\n\x12GetAuthnDescriptor\x12:.cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest\x1a;.cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse\"<\x82\xd3\xe4\x93\x02\x36\x12\x34/cosmos/base/reflection/v1beta1/app_descriptor/authn\x12\xcb\x01\n\x12GetChainDescriptor\x12:.cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest\x1a;.cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse\"<\x82\xd3\xe4\x93\x02\x36\x12\x34/cosmos/base/reflection/v1beta1/app_descriptor/chain\x12\xcb\x01\n\x12GetCodecDescriptor\x12:.cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest\x1a;.cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse\"<\x82\xd3\xe4\x93\x02\x36\x12\x34/cosmos/base/reflection/v1beta1/app_descriptor/codec\x12\xeb\x01\n\x1aGetConfigurationDescriptor\x12\x42.cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest\x1a\x43.cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse\"D\x82\xd3\xe4\x93\x02>\x12</cosmos/base/reflection/v1beta1/app_descriptor/configuration\x12\xec\x01\n\x1aGetQueryServicesDescriptor\x12\x42.cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest\x1a\x43.cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse\"E\x82\xd3\xe4\x93\x02?\x12=/cosmos/base/reflection/v1beta1/app_descriptor/query_services\x12\xca\x01\n\x0fGetTxDescriptor\x12\x37.cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest\x1a\x38.cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse\"D\x82\xd3\xe4\x93\x02>\x12</cosmos/base/reflection/v1beta1/app_descriptor/tx_descriptorB>Z<github.com/cosmos/cosmos-sdk/server/grpc/reflection/v2alpha1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.base.reflection.v2alpha1.reflection_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z<github.com/cosmos/cosmos-sdk/server/grpc/reflection/v2alpha1'
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetAuthnDescriptor']._loaded_options = None
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetAuthnDescriptor']._serialized_options = b'\202\323\344\223\0026\0224/cosmos/base/reflection/v1beta1/app_descriptor/authn'
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetChainDescriptor']._loaded_options = None
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetChainDescriptor']._serialized_options = b'\202\323\344\223\0026\0224/cosmos/base/reflection/v1beta1/app_descriptor/chain'
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetCodecDescriptor']._loaded_options = None
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetCodecDescriptor']._serialized_options = b'\202\323\344\223\0026\0224/cosmos/base/reflection/v1beta1/app_descriptor/codec'
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetConfigurationDescriptor']._loaded_options = None
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetConfigurationDescriptor']._serialized_options = b'\202\323\344\223\002>\022</cosmos/base/reflection/v1beta1/app_descriptor/configuration'
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetQueryServicesDescriptor']._loaded_options = None
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetQueryServicesDescriptor']._serialized_options = b'\202\323\344\223\002?\022=/cosmos/base/reflection/v1beta1/app_descriptor/query_services'
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetTxDescriptor']._loaded_options = None
  _globals['_REFLECTIONSERVICE'].methods_by_name['GetTxDescriptor']._serialized_options = b'\202\323\344\223\002>\022</cosmos/base/reflection/v1beta1/app_descriptor/tx_descriptor'
  _globals['_APPDESCRIPTOR']._serialized_start=116
  _globals['_APPDESCRIPTOR']._serialized_end=603
  _globals['_TXDESCRIPTOR']._serialized_start=605
  _globals['_TXDESCRIPTOR']._serialized_end=715
  _globals['_AUTHNDESCRIPTOR']._serialized_start=717
  _globals['_AUTHNDESCRIPTOR']._serialized_end=821
  _globals['_SIGNINGMODEDESCRIPTOR']._serialized_start=824
  _globals['_SIGNINGMODEDESCRIPTOR']._serialized_end=969
  _globals['_CHAINDESCRIPTOR']._serialized_start=971
  _globals['_CHAINDESCRIPTOR']._serialized_end=1004
  _globals['_CODECDESCRIPTOR']._serialized_start=1006
  _globals['_CODECDESCRIPTOR']._serialized_end=1109
  _globals['_INTERFACEDESCRIPTOR']._serialized_start=1112
  _globals['_INTERFACEDESCRIPTOR']._serialized_end=1418
  _globals['_INTERFACEIMPLEMENTERDESCRIPTOR']._serialized_start=1420
  _globals['_INTERFACEIMPLEMENTERDESCRIPTOR']._serialized_end=1507
  _globals['_INTERFACEACCEPTINGMESSAGEDESCRIPTOR']._serialized_start=1509
  _globals['_INTERFACEACCEPTINGMESSAGEDESCRIPTOR']._serialized_end=1628
  _globals['_CONFIGURATIONDESCRIPTOR']._serialized_start=1630
  _globals['_CONFIGURATIONDESCRIPTOR']._serialized_end=1722
  _globals['_MSGDESCRIPTOR']._serialized_start=1724
  _globals['_MSGDESCRIPTOR']._serialized_end=1773
  _globals['_GETAUTHNDESCRIPTORREQUEST']._serialized_start=1775
  _globals['_GETAUTHNDESCRIPTORREQUEST']._serialized_end=1802
  _globals['_GETAUTHNDESCRIPTORRESPONSE']._serialized_start=1804
  _globals['_GETAUTHNDESCRIPTORRESPONSE']._serialized_end=1904
  _globals['_GETCHAINDESCRIPTORREQUEST']._serialized_start=1906
  _globals['_GETCHAINDESCRIPTORREQUEST']._serialized_end=1933
  _globals['_GETCHAINDESCRIPTORRESPONSE']._serialized_start=1935
  _globals['_GETCHAINDESCRIPTORRESPONSE']._serialized_end=2035
  _globals['_GETCODECDESCRIPTORREQUEST']._serialized_start=2037
  _globals['_GETCODECDESCRIPTORREQUEST']._serialized_end=2064
  _globals['_GETCODECDESCRIPTORRESPONSE']._serialized_start=2066
  _globals['_GETCODECDESCRIPTORRESPONSE']._serialized_end=2166
  _globals['_GETCONFIGURATIONDESCRIPTORREQUEST']._serialized_start=2168
  _globals['_GETCONFIGURATIONDESCRIPTORREQUEST']._serialized_end=2203
  _globals['_GETCONFIGURATIONDESCRIPTORRESPONSE']._serialized_start=2205
  _globals['_GETCONFIGURATIONDESCRIPTORRESPONSE']._serialized_end=2323
  _globals['_GETQUERYSERVICESDESCRIPTORREQUEST']._serialized_start=2325
  _globals['_GETQUERYSERVICESDESCRIPTORREQUEST']._serialized_end=2360
  _globals['_GETQUERYSERVICESDESCRIPTORRESPONSE']._serialized_start=2362
  _globals['_GETQUERYSERVICESDESCRIPTORRESPONSE']._serialized_end=2482
  _globals['_GETTXDESCRIPTORREQUEST']._serialized_start=2484
  _globals['_GETTXDESCRIPTORREQUEST']._serialized_end=2508
  _globals['_GETTXDESCRIPTORRESPONSE']._serialized_start=2510
  _globals['_GETTXDESCRIPTORRESPONSE']._serialized_end=2598
  _globals['_QUERYSERVICESDESCRIPTOR']._serialized_start=2600
  _globals['_QUERYSERVICESDESCRIPTOR']._serialized_end=2721
  _globals['_QUERYSERVICEDESCRIPTOR']._serialized_start=2724
  _globals['_QUERYSERVICEDESCRIPTOR']._serialized_end=2887
  _globals['_QUERYMETHODDESCRIPTOR']._serialized_start=2889
  _globals['_QUERYMETHODDESCRIPTOR']._serialized_end=2972
  _globals['_REFLECTIONSERVICE']._serialized_start=2975
  _globals['_REFLECTIONSERVICE']._serialized_end=4294
# @@protoc_insertion_point(module_scope)
