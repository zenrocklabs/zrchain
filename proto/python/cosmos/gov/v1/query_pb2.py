# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/gov/v1/query.proto
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
    'cosmos/gov/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.gov.v1 import gov_pb2 as cosmos_dot_gov_dot_v1_dot_gov__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x19\x63osmos/gov/v1/query.proto\x12\rcosmos.gov.v1\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x17\x63osmos/gov/v1/gov.proto\x1a\x19\x63osmos_proto/cosmos.proto\"\x1a\n\x18QueryConstitutionRequest\"?\n\x19QueryConstitutionResponse\x12\"\n\x0c\x63onstitution\x18\x01 \x01(\tR\x0c\x63onstitution\"7\n\x14QueryProposalRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"L\n\x15QueryProposalResponse\x12\x33\n\x08proposal\x18\x01 \x01(\x0b\x32\x17.cosmos.gov.v1.ProposalR\x08proposal\"\x8f\x02\n\x15QueryProposalsRequest\x12\x46\n\x0fproposal_status\x18\x01 \x01(\x0e\x32\x1d.cosmos.gov.v1.ProposalStatusR\x0eproposalStatus\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter\x12\x36\n\tdepositor\x18\x03 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor\x12\x46\n\npagination\x18\x04 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x98\x01\n\x16QueryProposalsResponse\x12\x35\n\tproposals\x18\x01 \x03(\x0b\x32\x17.cosmos.gov.v1.ProposalR\tproposals\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"c\n\x10QueryVoteRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12.\n\x05voter\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x05voter\"<\n\x11QueryVoteResponse\x12\'\n\x04vote\x18\x01 \x01(\x0b\x32\x13.cosmos.gov.v1.VoteR\x04vote\"|\n\x11QueryVotesRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x88\x01\n\x12QueryVotesResponse\x12)\n\x05votes\x18\x01 \x03(\x0b\x32\x13.cosmos.gov.v1.VoteR\x05votes\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\"9\n\x12QueryParamsRequest\x12#\n\x0bparams_type\x18\x01 \x01(\tB\x02\x18\x01R\nparamsType\"\x96\x02\n\x13QueryParamsResponse\x12\x44\n\rvoting_params\x18\x01 \x01(\x0b\x32\x1b.cosmos.gov.v1.VotingParamsB\x02\x18\x01R\x0cvotingParams\x12G\n\x0e\x64\x65posit_params\x18\x02 \x01(\x0b\x32\x1c.cosmos.gov.v1.DepositParamsB\x02\x18\x01R\rdepositParams\x12\x41\n\x0ctally_params\x18\x03 \x01(\x0b\x32\x1a.cosmos.gov.v1.TallyParamsB\x02\x18\x01R\x0btallyParams\x12-\n\x06params\x18\x04 \x01(\x0b\x32\x15.cosmos.gov.v1.ParamsR\x06params\"n\n\x13QueryDepositRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x36\n\tdepositor\x18\x02 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\tdepositor\"H\n\x14QueryDepositResponse\x12\x30\n\x07\x64\x65posit\x18\x01 \x01(\x0b\x32\x16.cosmos.gov.v1.DepositR\x07\x64\x65posit\"\x7f\n\x14QueryDepositsRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\x12\x46\n\npagination\x18\x02 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\x94\x01\n\x15QueryDepositsResponse\x12\x32\n\x08\x64\x65posits\x18\x01 \x03(\x0b\x32\x16.cosmos.gov.v1.DepositR\x08\x64\x65posits\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\":\n\x17QueryTallyResultRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"L\n\x18QueryTallyResultResponse\x12\x30\n\x05tally\x18\x01 \x01(\x0b\x32\x1a.cosmos.gov.v1.TallyResultR\x05tally\"B\n\x1fQueryProposalVoteOptionsRequest\x12\x1f\n\x0bproposal_id\x18\x01 \x01(\x04R\nproposalId\"i\n QueryProposalVoteOptionsResponse\x12\x45\n\x0cvote_options\x18\x01 \x01(\x0b\x32\".cosmos.gov.v1.ProposalVoteOptionsR\x0bvoteOptions\"9\n\x1eQueryMessageBasedParamsRequest\x12\x17\n\x07msg_url\x18\x01 \x01(\tR\x06msgUrl\"\\\n\x1fQueryMessageBasedParamsResponse\x12\x39\n\x06params\x18\x01 \x01(\x0b\x32!.cosmos.gov.v1.MessageBasedParamsR\x06params2\xaa\x0c\n\x05Query\x12\x86\x01\n\x0c\x43onstitution\x12\'.cosmos.gov.v1.QueryConstitutionRequest\x1a(.cosmos.gov.v1.QueryConstitutionResponse\"#\x82\xd3\xe4\x93\x02\x1d\x12\x1b/cosmos/gov/v1/constitution\x12\x85\x01\n\x08Proposal\x12#.cosmos.gov.v1.QueryProposalRequest\x1a$.cosmos.gov.v1.QueryProposalResponse\".\x82\xd3\xe4\x93\x02(\x12&/cosmos/gov/v1/proposals/{proposal_id}\x12z\n\tProposals\x12$.cosmos.gov.v1.QueryProposalsRequest\x1a%.cosmos.gov.v1.QueryProposalsResponse\" \x82\xd3\xe4\x93\x02\x1a\x12\x18/cosmos/gov/v1/proposals\x12\x87\x01\n\x04Vote\x12\x1f.cosmos.gov.v1.QueryVoteRequest\x1a .cosmos.gov.v1.QueryVoteResponse\"<\x82\xd3\xe4\x93\x02\x36\x12\x34/cosmos/gov/v1/proposals/{proposal_id}/votes/{voter}\x12\x82\x01\n\x05Votes\x12 .cosmos.gov.v1.QueryVotesRequest\x1a!.cosmos.gov.v1.QueryVotesResponse\"4\x82\xd3\xe4\x93\x02.\x12,/cosmos/gov/v1/proposals/{proposal_id}/votes\x12n\n\x06Params\x12!.cosmos.gov.v1.QueryParamsRequest\x1a\".cosmos.gov.v1.QueryParamsResponse\"\x1d\x82\xd3\xe4\x93\x02\x17\x12\x15/cosmos/gov/v1/params\x12\x97\x01\n\x07\x44\x65posit\x12\".cosmos.gov.v1.QueryDepositRequest\x1a#.cosmos.gov.v1.QueryDepositResponse\"C\x82\xd3\xe4\x93\x02=\x12;/cosmos/gov/v1/proposals/{proposal_id}/deposits/{depositor}\x12\x8e\x01\n\x08\x44\x65posits\x12#.cosmos.gov.v1.QueryDepositsRequest\x1a$.cosmos.gov.v1.QueryDepositsResponse\"7\x82\xd3\xe4\x93\x02\x31\x12//cosmos/gov/v1/proposals/{proposal_id}/deposits\x12\x94\x01\n\x0bTallyResult\x12&.cosmos.gov.v1.QueryTallyResultRequest\x1a\'.cosmos.gov.v1.QueryTallyResultResponse\"4\x82\xd3\xe4\x93\x02.\x12,/cosmos/gov/v1/proposals/{proposal_id}/tally\x12\xb3\x01\n\x13ProposalVoteOptions\x12..cosmos.gov.v1.QueryProposalVoteOptionsRequest\x1a/.cosmos.gov.v1.QueryProposalVoteOptionsResponse\";\x82\xd3\xe4\x93\x02\x35\x12\x33/cosmos/gov/v1/proposals/{proposal_id}/vote_options\x12\x9c\x01\n\x12MessageBasedParams\x12-.cosmos.gov.v1.QueryMessageBasedParamsRequest\x1a..cosmos.gov.v1.QueryMessageBasedParamsResponse\"\'\x82\xd3\xe4\x93\x02!\x12\x1f/cosmos/gov/v1/params/{msg_url}B\x1dZ\x1b\x63osmossdk.io/x/gov/types/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.gov.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\033cosmossdk.io/x/gov/types/v1'
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['voter']._loaded_options = None
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['depositor']._loaded_options = None
  _globals['_QUERYPROPOSALSREQUEST'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYVOTEREQUEST'].fields_by_name['voter']._loaded_options = None
  _globals['_QUERYVOTEREQUEST'].fields_by_name['voter']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERYPARAMSREQUEST'].fields_by_name['params_type']._loaded_options = None
  _globals['_QUERYPARAMSREQUEST'].fields_by_name['params_type']._serialized_options = b'\030\001'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['voting_params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['voting_params']._serialized_options = b'\030\001'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['deposit_params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['deposit_params']._serialized_options = b'\030\001'
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['tally_params']._loaded_options = None
  _globals['_QUERYPARAMSRESPONSE'].fields_by_name['tally_params']._serialized_options = b'\030\001'
  _globals['_QUERYDEPOSITREQUEST'].fields_by_name['depositor']._loaded_options = None
  _globals['_QUERYDEPOSITREQUEST'].fields_by_name['depositor']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_QUERY'].methods_by_name['Constitution']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Constitution']._serialized_options = b'\202\323\344\223\002\035\022\033/cosmos/gov/v1/constitution'
  _globals['_QUERY'].methods_by_name['Proposal']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Proposal']._serialized_options = b'\202\323\344\223\002(\022&/cosmos/gov/v1/proposals/{proposal_id}'
  _globals['_QUERY'].methods_by_name['Proposals']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Proposals']._serialized_options = b'\202\323\344\223\002\032\022\030/cosmos/gov/v1/proposals'
  _globals['_QUERY'].methods_by_name['Vote']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Vote']._serialized_options = b'\202\323\344\223\0026\0224/cosmos/gov/v1/proposals/{proposal_id}/votes/{voter}'
  _globals['_QUERY'].methods_by_name['Votes']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Votes']._serialized_options = b'\202\323\344\223\002.\022,/cosmos/gov/v1/proposals/{proposal_id}/votes'
  _globals['_QUERY'].methods_by_name['Params']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Params']._serialized_options = b'\202\323\344\223\002\027\022\025/cosmos/gov/v1/params'
  _globals['_QUERY'].methods_by_name['Deposit']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Deposit']._serialized_options = b'\202\323\344\223\002=\022;/cosmos/gov/v1/proposals/{proposal_id}/deposits/{depositor}'
  _globals['_QUERY'].methods_by_name['Deposits']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Deposits']._serialized_options = b'\202\323\344\223\0021\022//cosmos/gov/v1/proposals/{proposal_id}/deposits'
  _globals['_QUERY'].methods_by_name['TallyResult']._loaded_options = None
  _globals['_QUERY'].methods_by_name['TallyResult']._serialized_options = b'\202\323\344\223\002.\022,/cosmos/gov/v1/proposals/{proposal_id}/tally'
  _globals['_QUERY'].methods_by_name['ProposalVoteOptions']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ProposalVoteOptions']._serialized_options = b'\202\323\344\223\0025\0223/cosmos/gov/v1/proposals/{proposal_id}/vote_options'
  _globals['_QUERY'].methods_by_name['MessageBasedParams']._loaded_options = None
  _globals['_QUERY'].methods_by_name['MessageBasedParams']._serialized_options = b'\202\323\344\223\002!\022\037/cosmos/gov/v1/params/{msg_url}'
  _globals['_QUERYCONSTITUTIONREQUEST']._serialized_start=170
  _globals['_QUERYCONSTITUTIONREQUEST']._serialized_end=196
  _globals['_QUERYCONSTITUTIONRESPONSE']._serialized_start=198
  _globals['_QUERYCONSTITUTIONRESPONSE']._serialized_end=261
  _globals['_QUERYPROPOSALREQUEST']._serialized_start=263
  _globals['_QUERYPROPOSALREQUEST']._serialized_end=318
  _globals['_QUERYPROPOSALRESPONSE']._serialized_start=320
  _globals['_QUERYPROPOSALRESPONSE']._serialized_end=396
  _globals['_QUERYPROPOSALSREQUEST']._serialized_start=399
  _globals['_QUERYPROPOSALSREQUEST']._serialized_end=670
  _globals['_QUERYPROPOSALSRESPONSE']._serialized_start=673
  _globals['_QUERYPROPOSALSRESPONSE']._serialized_end=825
  _globals['_QUERYVOTEREQUEST']._serialized_start=827
  _globals['_QUERYVOTEREQUEST']._serialized_end=926
  _globals['_QUERYVOTERESPONSE']._serialized_start=928
  _globals['_QUERYVOTERESPONSE']._serialized_end=988
  _globals['_QUERYVOTESREQUEST']._serialized_start=990
  _globals['_QUERYVOTESREQUEST']._serialized_end=1114
  _globals['_QUERYVOTESRESPONSE']._serialized_start=1117
  _globals['_QUERYVOTESRESPONSE']._serialized_end=1253
  _globals['_QUERYPARAMSREQUEST']._serialized_start=1255
  _globals['_QUERYPARAMSREQUEST']._serialized_end=1312
  _globals['_QUERYPARAMSRESPONSE']._serialized_start=1315
  _globals['_QUERYPARAMSRESPONSE']._serialized_end=1593
  _globals['_QUERYDEPOSITREQUEST']._serialized_start=1595
  _globals['_QUERYDEPOSITREQUEST']._serialized_end=1705
  _globals['_QUERYDEPOSITRESPONSE']._serialized_start=1707
  _globals['_QUERYDEPOSITRESPONSE']._serialized_end=1779
  _globals['_QUERYDEPOSITSREQUEST']._serialized_start=1781
  _globals['_QUERYDEPOSITSREQUEST']._serialized_end=1908
  _globals['_QUERYDEPOSITSRESPONSE']._serialized_start=1911
  _globals['_QUERYDEPOSITSRESPONSE']._serialized_end=2059
  _globals['_QUERYTALLYRESULTREQUEST']._serialized_start=2061
  _globals['_QUERYTALLYRESULTREQUEST']._serialized_end=2119
  _globals['_QUERYTALLYRESULTRESPONSE']._serialized_start=2121
  _globals['_QUERYTALLYRESULTRESPONSE']._serialized_end=2197
  _globals['_QUERYPROPOSALVOTEOPTIONSREQUEST']._serialized_start=2199
  _globals['_QUERYPROPOSALVOTEOPTIONSREQUEST']._serialized_end=2265
  _globals['_QUERYPROPOSALVOTEOPTIONSRESPONSE']._serialized_start=2267
  _globals['_QUERYPROPOSALVOTEOPTIONSRESPONSE']._serialized_end=2372
  _globals['_QUERYMESSAGEBASEDPARAMSREQUEST']._serialized_start=2374
  _globals['_QUERYMESSAGEBASEDPARAMSREQUEST']._serialized_end=2431
  _globals['_QUERYMESSAGEBASEDPARAMSRESPONSE']._serialized_start=2433
  _globals['_QUERYMESSAGEBASEDPARAMSRESPONSE']._serialized_end=2525
  _globals['_QUERY']._serialized_start=2528
  _globals['_QUERY']._serialized_end=4106
# @@protoc_insertion_point(module_scope)
