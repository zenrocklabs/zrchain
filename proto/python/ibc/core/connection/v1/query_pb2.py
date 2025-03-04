# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: ibc/core/connection/v1/query.proto
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
    'ibc/core/connection/v1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as cosmos_dot_base_dot_query_dot_v1beta1_dot_pagination__pb2
from ibc.core.client.v1 import client_pb2 as ibc_dot_core_dot_client_dot_v1_dot_client__pb2
from ibc.core.connection.v1 import connection_pb2 as ibc_dot_core_dot_connection_dot_v1_dot_connection__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\"ibc/core/connection/v1/query.proto\x12\x16ibc.core.connection.v1\x1a\x14gogoproto/gogo.proto\x1a*cosmos/base/query/v1beta1/pagination.proto\x1a\x1fibc/core/client/v1/client.proto\x1a\'ibc/core/connection/v1/connection.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x19google/protobuf/any.proto\"=\n\x16QueryConnectionRequest\x12#\n\rconnection_id\x18\x01 \x01(\tR\x0c\x63onnectionId\"\xbb\x01\n\x17QueryConnectionResponse\x12\x45\n\nconnection\x18\x01 \x01(\x0b\x32%.ibc.core.connection.v1.ConnectionEndR\nconnection\x12\x14\n\x05proof\x18\x02 \x01(\x0cR\x05proof\x12\x43\n\x0cproof_height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\"a\n\x17QueryConnectionsRequest\x12\x46\n\npagination\x18\x01 \x01(\x0b\x32&.cosmos.base.query.v1beta1.PageRequestR\npagination\"\xed\x01\n\x18QueryConnectionsResponse\x12N\n\x0b\x63onnections\x18\x01 \x03(\x0b\x32,.ibc.core.connection.v1.IdentifiedConnectionR\x0b\x63onnections\x12G\n\npagination\x18\x02 \x01(\x0b\x32\'.cosmos.base.query.v1beta1.PageResponseR\npagination\x12\x38\n\x06height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x06height\"<\n\x1dQueryClientConnectionsRequest\x12\x1b\n\tclient_id\x18\x01 \x01(\tR\x08\x63lientId\"\xa6\x01\n\x1eQueryClientConnectionsResponse\x12)\n\x10\x63onnection_paths\x18\x01 \x03(\tR\x0f\x63onnectionPaths\x12\x14\n\x05proof\x18\x02 \x01(\x0cR\x05proof\x12\x43\n\x0cproof_height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\"H\n!QueryConnectionClientStateRequest\x12#\n\rconnection_id\x18\x01 \x01(\tR\x0c\x63onnectionId\"\xe2\x01\n\"QueryConnectionClientStateResponse\x12\x61\n\x17identified_client_state\x18\x01 \x01(\x0b\x32).ibc.core.client.v1.IdentifiedClientStateR\x15identifiedClientState\x12\x14\n\x05proof\x18\x02 \x01(\x0cR\x05proof\x12\x43\n\x0cproof_height\x18\x03 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\"\x9d\x01\n$QueryConnectionConsensusStateRequest\x12#\n\rconnection_id\x18\x01 \x01(\tR\x0c\x63onnectionId\x12\'\n\x0frevision_number\x18\x02 \x01(\x04R\x0erevisionNumber\x12\'\n\x0frevision_height\x18\x03 \x01(\x04R\x0erevisionHeight\"\xde\x01\n%QueryConnectionConsensusStateResponse\x12=\n\x0f\x63onsensus_state\x18\x01 \x01(\x0b\x32\x14.google.protobuf.AnyR\x0e\x63onsensusState\x12\x1b\n\tclient_id\x18\x02 \x01(\tR\x08\x63lientId\x12\x14\n\x05proof\x18\x03 \x01(\x0cR\x05proof\x12\x43\n\x0cproof_height\x18\x04 \x01(\x0b\x32\x1a.ibc.core.client.v1.HeightB\x04\xc8\xde\x1f\x00R\x0bproofHeight\"\x1e\n\x1cQueryConnectionParamsRequest\"W\n\x1dQueryConnectionParamsResponse\x12\x36\n\x06params\x18\x01 \x01(\x0b\x32\x1e.ibc.core.connection.v1.ParamsR\x06params2\xb9\t\n\x05Query\x12\xaa\x01\n\nConnection\x12..ibc.core.connection.v1.QueryConnectionRequest\x1a/.ibc.core.connection.v1.QueryConnectionResponse\";\x82\xd3\xe4\x93\x02\x35\x12\x33/ibc/core/connection/v1/connections/{connection_id}\x12\x9d\x01\n\x0b\x43onnections\x12/.ibc.core.connection.v1.QueryConnectionsRequest\x1a\x30.ibc.core.connection.v1.QueryConnectionsResponse\"+\x82\xd3\xe4\x93\x02%\x12#/ibc/core/connection/v1/connections\x12\xc2\x01\n\x11\x43lientConnections\x12\x35.ibc.core.connection.v1.QueryClientConnectionsRequest\x1a\x36.ibc.core.connection.v1.QueryClientConnectionsResponse\">\x82\xd3\xe4\x93\x02\x38\x12\x36/ibc/core/connection/v1/client_connections/{client_id}\x12\xd8\x01\n\x15\x43onnectionClientState\x12\x39.ibc.core.connection.v1.QueryConnectionClientStateRequest\x1a:.ibc.core.connection.v1.QueryConnectionClientStateResponse\"H\x82\xd3\xe4\x93\x02\x42\x12@/ibc/core/connection/v1/connections/{connection_id}/client_state\x12\x98\x02\n\x18\x43onnectionConsensusState\x12<.ibc.core.connection.v1.QueryConnectionConsensusStateRequest\x1a=.ibc.core.connection.v1.QueryConnectionConsensusStateResponse\"\x7f\x82\xd3\xe4\x93\x02y\x12w/ibc/core/connection/v1/connections/{connection_id}/consensus_state/revision/{revision_number}/height/{revision_height}\x12\xa7\x01\n\x10\x43onnectionParams\x12\x34.ibc.core.connection.v1.QueryConnectionParamsRequest\x1a\x35.ibc.core.connection.v1.QueryConnectionParamsResponse\"&\x82\xd3\xe4\x93\x02 \x12\x1e/ibc/core/connection/v1/paramsB?Z=github.com/cosmos/ibc-go/v10/modules/core/03-connection/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'ibc.core.connection.v1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z=github.com/cosmos/ibc-go/v10/modules/core/03-connection/types'
  _globals['_QUERYCONNECTIONRESPONSE'].fields_by_name['proof_height']._loaded_options = None
  _globals['_QUERYCONNECTIONRESPONSE'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_QUERYCONNECTIONSRESPONSE'].fields_by_name['height']._loaded_options = None
  _globals['_QUERYCONNECTIONSRESPONSE'].fields_by_name['height']._serialized_options = b'\310\336\037\000'
  _globals['_QUERYCLIENTCONNECTIONSRESPONSE'].fields_by_name['proof_height']._loaded_options = None
  _globals['_QUERYCLIENTCONNECTIONSRESPONSE'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_QUERYCONNECTIONCLIENTSTATERESPONSE'].fields_by_name['proof_height']._loaded_options = None
  _globals['_QUERYCONNECTIONCLIENTSTATERESPONSE'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_QUERYCONNECTIONCONSENSUSSTATERESPONSE'].fields_by_name['proof_height']._loaded_options = None
  _globals['_QUERYCONNECTIONCONSENSUSSTATERESPONSE'].fields_by_name['proof_height']._serialized_options = b'\310\336\037\000'
  _globals['_QUERY'].methods_by_name['Connection']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Connection']._serialized_options = b'\202\323\344\223\0025\0223/ibc/core/connection/v1/connections/{connection_id}'
  _globals['_QUERY'].methods_by_name['Connections']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Connections']._serialized_options = b'\202\323\344\223\002%\022#/ibc/core/connection/v1/connections'
  _globals['_QUERY'].methods_by_name['ClientConnections']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ClientConnections']._serialized_options = b'\202\323\344\223\0028\0226/ibc/core/connection/v1/client_connections/{client_id}'
  _globals['_QUERY'].methods_by_name['ConnectionClientState']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ConnectionClientState']._serialized_options = b'\202\323\344\223\002B\022@/ibc/core/connection/v1/connections/{connection_id}/client_state'
  _globals['_QUERY'].methods_by_name['ConnectionConsensusState']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ConnectionConsensusState']._serialized_options = b'\202\323\344\223\002y\022w/ibc/core/connection/v1/connections/{connection_id}/consensus_state/revision/{revision_number}/height/{revision_height}'
  _globals['_QUERY'].methods_by_name['ConnectionParams']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ConnectionParams']._serialized_options = b'\202\323\344\223\002 \022\036/ibc/core/connection/v1/params'
  _globals['_QUERYCONNECTIONREQUEST']._serialized_start=259
  _globals['_QUERYCONNECTIONREQUEST']._serialized_end=320
  _globals['_QUERYCONNECTIONRESPONSE']._serialized_start=323
  _globals['_QUERYCONNECTIONRESPONSE']._serialized_end=510
  _globals['_QUERYCONNECTIONSREQUEST']._serialized_start=512
  _globals['_QUERYCONNECTIONSREQUEST']._serialized_end=609
  _globals['_QUERYCONNECTIONSRESPONSE']._serialized_start=612
  _globals['_QUERYCONNECTIONSRESPONSE']._serialized_end=849
  _globals['_QUERYCLIENTCONNECTIONSREQUEST']._serialized_start=851
  _globals['_QUERYCLIENTCONNECTIONSREQUEST']._serialized_end=911
  _globals['_QUERYCLIENTCONNECTIONSRESPONSE']._serialized_start=914
  _globals['_QUERYCLIENTCONNECTIONSRESPONSE']._serialized_end=1080
  _globals['_QUERYCONNECTIONCLIENTSTATEREQUEST']._serialized_start=1082
  _globals['_QUERYCONNECTIONCLIENTSTATEREQUEST']._serialized_end=1154
  _globals['_QUERYCONNECTIONCLIENTSTATERESPONSE']._serialized_start=1157
  _globals['_QUERYCONNECTIONCLIENTSTATERESPONSE']._serialized_end=1383
  _globals['_QUERYCONNECTIONCONSENSUSSTATEREQUEST']._serialized_start=1386
  _globals['_QUERYCONNECTIONCONSENSUSSTATEREQUEST']._serialized_end=1543
  _globals['_QUERYCONNECTIONCONSENSUSSTATERESPONSE']._serialized_start=1546
  _globals['_QUERYCONNECTIONCONSENSUSSTATERESPONSE']._serialized_end=1768
  _globals['_QUERYCONNECTIONPARAMSREQUEST']._serialized_start=1770
  _globals['_QUERYCONNECTIONPARAMSREQUEST']._serialized_end=1800
  _globals['_QUERYCONNECTIONPARAMSRESPONSE']._serialized_start=1802
  _globals['_QUERYCONNECTIONPARAMSRESPONSE']._serialized_end=1889
  _globals['_QUERY']._serialized_start=1892
  _globals['_QUERY']._serialized_end=3101
# @@protoc_insertion_point(module_scope)
