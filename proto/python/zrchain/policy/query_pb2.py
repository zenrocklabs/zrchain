# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/policy/query.proto
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
    'zrchain/policy/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from amino import amino_pb2 as amino_dot_amino__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from zrchain.policy import action_pb2 as zrchain_dot_policy_dot_action__pb2
from zrchain.policy import params_pb2 as zrchain_dot_policy_dot_params__pb2
from zrchain.policy import policy_pb2 as zrchain_dot_policy_dot_policy__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1azrchain/policy/query.proto\x12\x0ezrchain.policy\x1a\x11\x61mino/amino.proto\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x1bzrchain/policy/action.proto\x1a\x1bzrchain/policy/params.proto\x1a\x1bzrchain/policy/policy.proto\x1a\x19google/protobuf/any.proto\"\x14\n\x12QueryParamsRequest\"P\n\x13QueryParamsResponse\x12\x39\n\x06params\x18\x01 \x01(\x0b\x32\x16.zrchain.policy.ParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x06params\"\xad\x01\n\x13QueryActionsRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ress\x12\x34\n\x06status\x18\x03 \x01(\x0e\x32\x1c.zrchain.policy.ActionStatusR\x06status\"\x9f\x01\n\x14QueryActionsResponse\x12G\n\npagination\x18\x01 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\x12>\n\x07\x61\x63tions\x18\x02 \x03(\x0b\x32\x1e.zrchain.policy.ActionResponseB\x04\xc8\xde\x1f\x00R\x07\x61\x63tions\"^\n\x14QueryPoliciesRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"r\n\x0ePolicyResponse\x12.\n\x06policy\x18\x01 \x01(\x0b\x32\x16.zrchain.policy.PolicyR\x06policy\x12\x30\n\x08metadata\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyR\x08metadata\"\xa2\x01\n\x15QueryPoliciesResponse\x12G\n\npagination\x18\x01 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\x12@\n\x08policies\x18\x02 \x03(\x0b\x32\x1e.zrchain.policy.PolicyResponseB\x04\xc8\xde\x1f\x00R\x08policies\"(\n\x16QueryPolicyByIdRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\"Q\n\x17QueryPolicyByIdResponse\x12\x36\n\x06policy\x18\x01 \x01(\x0b\x32\x1e.zrchain.policy.PolicyResponseR\x06policy\"\x84\x01\n QuerySignMethodsByAddressRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ress\"\x9a\x01\n!QuerySignMethodsByAddressResponse\x12G\n\npagination\x18\x01 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\x12,\n\x06\x63onfig\x18\x02 \x03(\x0b\x32\x14.google.protobuf.AnyR\x06\x63onfig\"\x83\x01\n\x1dQueryPoliciesByCreatorRequest\x12\x1a\n\x08\x63reators\x18\x01 \x03(\tR\x08\x63reators\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x9d\x01\n\x1eQueryPoliciesByCreatorResponse\x12\x32\n\x08policies\x18\x01 \x03(\x0b\x32\x16.zrchain.policy.PolicyR\x08policies\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"/\n\x1dQueryActionDetailsByIdRequest\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\"\x82\x02\n\x1eQueryActionDetailsByIdResponse\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12.\n\x06\x61\x63tion\x18\x02 \x01(\x0b\x32\x16.zrchain.policy.ActionR\x06\x61\x63tion\x12.\n\x06policy\x18\x03 \x01(\x0b\x32\x16.zrchain.policy.PolicyR\x06policy\x12\x1c\n\tapprovers\x18\x04 \x03(\tR\tapprovers\x12+\n\x11pending_approvers\x18\x05 \x03(\tR\x10pendingApprovers\x12%\n\x0e\x63urrent_height\x18\x06 \x01(\x04R\rcurrentHeight2\x85\x08\n\x05Query\x12q\n\x06Params\x12\".zrchain.policy.QueryParamsRequest\x1a#.zrchain.policy.QueryParamsResponse\"\x1e\x82\xd3\xe4\x93\x02\x18\x12\x16/zrchain/policy/params\x12u\n\x07\x41\x63tions\x12#.zrchain.policy.QueryActionsRequest\x1a$.zrchain.policy.QueryActionsResponse\"\x1f\x82\xd3\xe4\x93\x02\x19\x12\x17/zrchain/policy/actions\x12y\n\x08Policies\x12$.zrchain.policy.QueryPoliciesRequest\x1a%.zrchain.policy.QueryPoliciesResponse\" \x82\xd3\xe4\x93\x02\x1a\x12\x18/zrchain/policy/policies\x12\x88\x01\n\nPolicyById\x12&.zrchain.policy.QueryPolicyByIdRequest\x1a\'.zrchain.policy.QueryPolicyByIdResponse\")\x82\xd3\xe4\x93\x02#\x12!/zrchain/policy/policy_by_id/{id}\x12\xb6\x01\n\x14SignMethodsByAddress\x12\x30.zrchain.policy.QuerySignMethodsByAddressRequest\x1a\x31.zrchain.policy.QuerySignMethodsByAddressResponse\"9\x82\xd3\xe4\x93\x02\x33\x12\x31/zrchain/policy/sign_methods_by_address/{address}\x12\xaa\x01\n\x11PoliciesByCreator\x12-.zrchain.policy.QueryPoliciesByCreatorRequest\x1a..zrchain.policy.QueryPoliciesByCreatorResponse\"6\x82\xd3\xe4\x93\x02\x30\x12./zrchain/policy/policies_by_creator/{creators}\x12\xa5\x01\n\x11\x41\x63tionDetailsById\x12-.zrchain.policy.QueryActionDetailsByIdRequest\x1a..zrchain.policy.QueryActionDetailsByIdResponse\"1\x82\xd3\xe4\x93\x02+\x12)/zrchain/policy/action_details_by_id/{id}B9Z7github.com/Zenrock-Foundation/zrchain/v5/x/policy/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.policy.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z7github.com/Zenrock-Foundation/zrchain/v5/x/policy/types'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYACTIONSRESPONSE'].fields_by_name['actions']._loaded_options = None
  _globals['_QUERYACTIONSRESPONSE'].fields_by_name['actions']._serialized_options = b'\310\336\037\000'
  _globals['_QUERYPOLICIESRESPONSE'].fields_by_name['policies']._loaded_options = None
  _globals['_QUERYPOLICIESRESPONSE'].fields_by_name['policies']._serialized_options = b'\310\336\037\000'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\030\022\026/zrchain/policy/params'
  _globals['_QUERY'].methods_by_name['Actions']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Actions']._serialized_options = b'\202\323\344\223\002\031\022\027/zrchain/policy/actions'
  _globals['_QUERY'].methods_by_name['Policies']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Policies']._serialized_options = b'\202\323\344\223\002\032\022\030/zrchain/policy/policies'
  _globals['_QUERY'].methods_by_name['PolicyById']._loaded_options = None
  _globals['_QUERY'].methods_by_name['PolicyById']._serialized_options = b'\202\323\344\223\002#\022!/zrchain/policy/policy_by_id/{id}'
  _globals['_QUERY'].methods_by_name['SignMethodsByAddress']._loaded_options = None
  _globals['_QUERY'].methods_by_name['SignMethodsByAddress']._serialized_options = b'\202\323\344\223\0023\0221/zrchain/policy/sign_methods_by_address/{address}'
  _globals['_QUERY'].methods_by_name['PoliciesByCreator']._loaded_options = None
  _globals['_QUERY'].methods_by_name['PoliciesByCreator']._serialized_options = b'\202\323\344\223\0020\022./zrchain/policy/policies_by_creator/{creators}'
  _globals['_QUERY'].methods_by_name['ActionDetailsById']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ActionDetailsById']._serialized_options = b'\202\323\344\223\002+\022)/zrchain/policy/action_details_by_id/{id}'
  _globals['_QUERYPARAMSREQUEST']._serialized_start=275
  _globals['_QUERYPARAMSREQUEST']._serialized_end=295
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=297
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=377
  _globals['_QUERYACTIONSREQUEST']._serialized_start=380
  _globals['_QUERYACTIONSREQUEST']._serialized_end=553
  _globals['_QUERYACTIONSRESPONSE']._serialized_start=556
  _globals['_QUERYACTIONSRESPONSE']._serialized_end=715
  _globals['_QUERYPOLICIESREQUEST']._serialized_start=717
  _globals['_QUERYPOLICIESREQUEST']._serialized_end=811
  _globals['_POLICYRESPONSE']._serialized_start=813
  _globals['_POLICYRESPONSE']._serialized_end=927
  _globals['_QUERYPOLICIESRESPONSE']._serialized_start=930
  _globals['_QUERYPOLICIESRESPONSE']._serialized_end=1092
  _globals['_QUERYPOLICYBYIDREQUEST']._serialized_start=1094
  _globals['_QUERYPOLICYBYIDREQUEST']._serialized_end=1134
  _globals['_QUERYPOLICYBYIDRESPONSE']._serialized_start=1136
  _globals['_QUERYPOLICYBYIDRESPONSE']._serialized_end=1217
  _globals['_QUERYSIGNMETHODSBYADDRESSREQUEST']._serialized_start=1220
  _globals['_QUERYSIGNMETHODSBYADDRESSREQUEST']._serialized_end=1352
  _globals['_QUERYSIGNMETHODSBYADDRESSRESPONSE']._serialized_start=1355
  _globals['_QUERYSIGNMETHODSBYADDRESSRESPONSE']._serialized_end=1509
  _globals['_QUERYPOLICIESBYCREATORREQUEST']._serialized_start=1512
  _globals['_QUERYPOLICIESBYCREATORREQUEST']._serialized_end=1643
  _globals['_QUERYPOLICIESBYCREATORRESPONSE']._serialized_start=1646
  _globals['_QUERYPOLICIESBYCREATORRESPONSE']._serialized_end=1803
  _globals['_QUERYACTIONDETAILSBYIDREQUEST']._serialized_start=1805
  _globals['_QUERYACTIONDETAILSBYIDREQUEST']._serialized_end=1852
  _globals['_QUERYACTIONDETAILSBYIDRESPONSE']._serialized_start=1855
  _globals['_QUERYACTIONDETAILSBYIDRESPONSE']._serialized_end=2113
  _globals['_QUERY']._serialized_start=2116
  _globals['_QUERY']._serialized_end=3145
# @@protoc_insertion_point(module_scope)
