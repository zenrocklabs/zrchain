# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/gov/v1beta1/query.proto
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
    'cosmos/gov/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.gov.v1beta1 import gov_pb2 as cosmos_dot_gov_dot_v1beta1_dot_gov__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1e\x63osmos/gov/v1beta1/query.proto\x12\x12\x63osmos.gov.v1beta1\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x14gogoproto/gogo.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1c\x63osmos/gov/v1beta1/gov.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x11\x61mino/amino.proto\"7\n\x14QueryProposalRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"\\\n\x15QueryProposalResponse\x12\x43\n\x08proposal\x18\x01 \x01(\x0b\x32\x1c.cosmos.gov.v1beta1.ProposalB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x08proposal\"\x9e\x02\n\x15QueryProposalsRequest\x12K\n\x0fproposal_status\x18\x01 \x01(\x0e\x32\".cosmos.gov.v1beta1.ProposalStatusR\x0eproposalStatus\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter\x12\x36\n\tdepositor\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor\x12\x46\n\npagination\x18\x04 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"\xa8\x01\n\x16QueryProposalsResponse\x12\x45\n\tproposals\x18\x01 \x03(\x0b\x32\x1c.cosmos.gov.v1beta1.ProposalB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\tproposals\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"m\n\x10QueryVoteRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"L\n\x11QueryVoteResponse\x12\x37\n\x04vote\x18\x01 \x01(\x0b\x32\x18.cosmos.gov.v1beta1.VoteB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x04vote\"|\n\x11QueryVotesRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x98\x01\n\x12QueryVotesResponse\x12\x39\n\x05votes\x18\x01 \x03(\x0b\x32\x18.cosmos.gov.v1beta1.VoteB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x05votes\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"5\n\x12QueryParamsRequest\x12\x1f\n\x0bparams_type\x18\x01 \x01(\tR\nparamsType\"\x8b\x02\n\x13QueryParamsResponse\x12P\n\rvoting_params\x18\x01 \x01(\x0b\x32 .cosmos.gov.v1beta1.VotingParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0cvotingParams\x12S\n\x0e\x64\x65posit_params\x18\x02 \x01(\x0b\x32!.cosmos.gov.v1beta1.DepositParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\rdepositParams\x12M\n\x0ctally_params\x18\x03 \x01(\x0b\x32\x1f.cosmos.gov.v1beta1.TallyParamsB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0btallyParams\"x\n\x13QueryDepositRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x36\n\tdepositor\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"X\n\x14QueryDepositResponse\x12@\n\x07\x64\x65posit\x18\x01 \x01(\x0b\x32\x1b.cosmos.gov.v1beta1.DepositB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x07\x64\x65posit\"\x7f\n\x14QueryDepositsRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xa4\x01\n\x15QueryDepositsResponse\x12\x42\n\x08\x64\x65posits\x18\x01 \x03(\x0b\x32\x1b.cosmos.gov.v1beta1.DepositB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x08\x64\x65posits\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\":\n\x17QueryTallyResultRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"\\\n\x18QueryTallyResultResponse\x12@\n\x05tally\x18\x01 \x01(\x0b\x32\x1f.cosmos.gov.v1beta1.TallyResultB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x05tally2\xd4\t\n\x05Query\x12\x94\x01\n\x08Proposal\x12(.cosmos.gov.v1beta1.QueryProposalRequest\x1a).cosmos.gov.v1beta1.QueryProposalResponse\"3\x82\xd3\xe4\x93\x02-\x12+/cosmos/gov/v1beta1/proposals/{proposal_id}\x12\x89\x01\n\tProposals\x12).cosmos.gov.v1beta1.QueryProposalsRequest\x1a*.cosmos.gov.v1beta1.QueryProposalsResponse\"%\x82\xd3\xe4\x93\x02\x1f\x12\x1d/cosmos/gov/v1beta1/proposals\x12\x96\x01\n\x04Vote\x12$.cosmos.gov.v1beta1.QueryVoteRequest\x1a%.cosmos.gov.v1beta1.QueryVoteResponse\"A\x82\xd3\xe4\x93\x02;\x12\x39/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}\x12\x91\x01\n\x05Votes\x12%.cosmos.gov.v1beta1.QueryVotesRequest\x1a&.cosmos.gov.v1beta1.QueryVotesResponse\"9\x82\xd3\xe4\x93\x02\x33\x12\x31/cosmos/gov/v1beta1/proposals/{proposal_id}/votes\x12\x8b\x01\n\x06Params\x12&.cosmos.gov.v1beta1.QueryParamsRequest\x1a\'.cosmos.gov.v1beta1.QueryParamsResponse\"0\x82\xd3\xe4\x93\x02*\x12(/cosmos/gov/v1beta1/params/{params_type}\x12\xa6\x01\n\x07\x44\x65posit\x12\'.cosmos.gov.v1beta1.QueryDepositRequest\x1a(.cosmos.gov.v1beta1.QueryDepositResponse\"H\x82\xd3\xe4\x93\x02\x42\x12@/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}\x12\x9d\x01\n\x08\x44\x65posits\x12(.cosmos.gov.v1beta1.QueryDepositsRequest\x1a).cosmos.gov.v1beta1.QueryDepositsResponse\"<\x82\xd3\xe4\x93\x02\x36\x12\x34/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits\x12\xa3\x01\n\x0bTallyResult\x12+.cosmos.gov.v1beta1.QueryTallyResultRequest\x1a,.cosmos.gov.v1beta1.QueryTallyResultResponse\"9\x82\xd3\xe4\x93\x02\x33\x12\x31/cosmos/gov/v1beta1/proposals/{proposal_id}/tallyB\"Z cosmossdk.io/x/gov/types/v1beta1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.gov.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z cosmossdk.io/x/gov/types/v1beta1'
  _globals['_QUERYPROPOSALRESPONSE'].fields_by_name['proposal']._loaded_options = None
  _globals['_QUERYPROPOSALRESPONSE'].fields_by_name['proposal']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['voter']._loaded_options = None
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['depositor']._loaded_options = None
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYPROPOSALSREQUEST']._loaded_options = None
  _globals['_QUERYPROPOSALSREQUEST']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_QUERYPROPOSALSRESPONSE'].fields_by_name['proposals']._loaded_options = None
  _globals['_QUERYPROPOSALSRESPONSE'].fields_by_name['proposals']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYVOTEREQUEST'].fields_by_name['voter']._loaded_options = None
  _globals['_QUERYVOTEREQUEST'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYVOTEREQUEST']._loaded_options = None
  _globals['_QUERYVOTEREQUEST']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_QUERYVOTERESPONSE'].fields_by_name['vote']._loaded_options = None
  _globals['_QUERYVOTERESPONSE'].fields_by_name['vote']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYVOTESRESPONSE'].fields_by_name['votes']._loaded_options = None
  _globals['_QUERYVOTESRESPONSE'].fields_by_name['votes']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['voting_params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['voting_params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['deposit_params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['deposit_params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['tally_params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['tally_params']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYDEPOSITREQUEST'].fields_by_name['depositor']._loaded_options = None
  _globals['_QUERYDEPOSITREQUEST'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYDEPOSITREQUEST']._loaded_options = None
  _globals['_QUERYDEPOSITREQUEST']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_QUERYDEPOSITRESPONSE'].fields_by_name['deposit']._loaded_options = None
  _globals['_QUERYDEPOSITRESPONSE'].fields_by_name['deposit']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYDEPOSITSRESPONSE'].fields_by_name['deposits']._loaded_options = None
  _globals['_QUERYDEPOSITSRESPONSE'].fields_by_name['deposits']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERYTALLYRESULTRESPONSE'].fields_by_name['tally']._loaded_options = None
  _globals['_QUERYTALLYRESULTRESPONSE'].fields_by_name['tally']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_QUERY'].methods_by_name['Proposal']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Proposal']._serialized_options = b'\202\323\344\223\002-\022+/cosmos/gov/v1beta1/proposals/{proposal_id}'
  _globals['_QUERY'].methods_by_name['Proposals']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Proposals']._serialized_options = b'\202\323\344\223\002\037\022\035/cosmos/gov/v1beta1/proposals'
  _globals['_QUERY'].methods_by_name['Vote']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Vote']._serialized_options = b'\202\323\344\223\002;\0229/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}'
  _globals['_QUERY'].methods_by_name['Votes']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Votes']._serialized_options = b'\202\323\344\223\0023\0221/cosmos/gov/v1beta1/proposals/{proposal_id}/votes'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002*\022(/cosmos/gov/v1beta1/params/{params_type}'
  _globals['_QUERY'].methods_by_name['Deposit']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Deposit']._serialized_options = b'\202\323\344\223\002B\022@/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}'
  _globals['_QUERY'].methods_by_name['Deposits']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Deposits']._serialized_options = b'\202\323\344\223\0026\0224/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits'
  _globals['_QUERY'].methods_by_name['TallyResult']._loaded_options = None
  _globals['_QUERY'].methods_by_name['TallyResult']._serialized_options = b'\202\323\344\223\0023\0221/cosmos/gov/v1beta1/proposals/{proposal_id}/tally'
  _globals['_QUERYPROPOSALREQUEST']._serialized_start=226
  _globals['_QUERYPROPOSALREQUEST']._serialized_end=281
  _globals['_QUERYPROPOSALRESPONSE']._serialized_start=283
  _globals['_QUERYPROPOSALRESPONSE']._serialized_end=375
  _globals['_QUERYPROPOSALSREQUEST']._serialized_start=378
  _globals['_QUERYPROPOSALSREQUEST']._serialized_end=664
  _globals['_QUERYPROPOSALSRESPONSE']._serialized_start=667
  _globals['_QUERYPROPOSALSRESPONSE']._serialized_end=835
  _globals['_QUERYVOTEREQUEST']._serialized_start=837
  _globals['_QUERYVOTEREQUEST']._serialized_end=946
  _globals['_QUERYVOTERESPONSE']._serialized_start=948
  _globals['_QUERYVOTERESPONSE']._serialized_end=1024
  _globals['_QUERYVOTESREQUEST']._serialized_start=1026
  _globals['_QUERYVOTESREQUEST']._serialized_end=1150
  _globals['_QUERYVOTESRESPONSE']._serialized_start=1153
  _globals['_QUERYVOTESRESPONSE']._serialized_end=1305
  _globals['_QUERYPARAMSREQUEST']._serialized_start=1307
  _globals['_QUERYPARAMSREQUEST']._serialized_end=1360
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=1363
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=1630
  _globals['_QUERYDEPOSITREQUEST']._serialized_start=1632
  _globals['_QUERYDEPOSITREQUEST']._serialized_end=1752
  _globals['_QUERYDEPOSITRESPONSE']._serialized_start=1754
  _globals['_QUERYDEPOSITRESPONSE']._serialized_end=1842
  _globals['_QUERYDEPOSITSREQUEST']._serialized_start=1844
  _globals['_QUERYDEPOSITSREQUEST']._serialized_end=1971
  _globals['_QUERYDEPOSITSRESPONSE']._serialized_start=1974
  _globals['_QUERYDEPOSITSRESPONSE']._serialized_end=2138
  _globals['_QUERYTALLYRESULTREQUEST']._serialized_start=2140
  _globals['_QUERYTALLYRESULTREQUEST']._serialized_end=2198
  _globals['_QUERYTALLYRESULTRESPONSE']._serialized_start=2200
  _globals['_QUERYTALLYRESULTRESPONSE']._serialized_end=2292
  _globals['_QUERY']._serialized_start=2295
  _globals['_QUERY']._serialized_end=3531
# @@protoc_insertion_point(module_scope)
