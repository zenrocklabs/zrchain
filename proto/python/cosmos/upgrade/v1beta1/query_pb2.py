# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/upgrade/v1beta1/query.proto
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
    'cosmos/upgrade/v1beta1/query.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from cosmos.upgrade.v1beta1 import upgrade_pb2 as cosmos_dot_upgrade_dot_v1beta1_dot_upgrade__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\"cosmos/upgrade/v1beta1/query.proto\x12\x16\x63osmos.upgrade.v1beta1\x1a\x1cgoogle/api/annotations.proto\x1a$cosmos/upgrade/v1beta1/upgrade.proto\"\x19\n\x17QueryCurrentPlanRequest\"L\n\x18QueryCurrentPlanResponse\x12\x30\n\x04plan\x18\x01 \x01(\x0b\x32\x1c.cosmos.upgrade.v1beta1.PlanR\x04plan\"-\n\x17QueryAppliedPlanRequest\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\"2\n\x18QueryAppliedPlanResponse\x12\x16\n\x06height\x18\x01 \x01(\x03R\x06height\"I\n\"QueryUpgradedConsensusStateRequest\x12\x1f\n\x0blast_height\x18\x01 \x01(\x03R\nlastHeight:\x02\x18\x01\"i\n#QueryUpgradedConsensusStateResponse\x12\x38\n\x18upgraded_consensus_state\x18\x02 \x01(\x0cR\x16upgradedConsensusState:\x02\x18\x01J\x04\x08\x01\x10\x02\"=\n\x1aQueryModuleVersionsRequest\x12\x1f\n\x0bmodule_name\x18\x01 \x01(\tR\nmoduleName\"m\n\x1bQueryModuleVersionsResponse\x12N\n\x0fmodule_versions\x18\x01 \x03(\x0b\x32%.cosmos.upgrade.v1beta1.ModuleVersionR\x0emoduleVersions\"\x17\n\x15QueryAuthorityRequest\"2\n\x16QueryAuthorityResponse\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress2\xf4\x06\n\x05Query\x12\x9e\x01\n\x0b\x43urrentPlan\x12/.cosmos.upgrade.v1beta1.QueryCurrentPlanRequest\x1a\x30.cosmos.upgrade.v1beta1.QueryCurrentPlanResponse\",\x82\xd3\xe4\x93\x02&\x12$/cosmos/upgrade/v1beta1/current_plan\x12\xa5\x01\n\x0b\x41ppliedPlan\x12/.cosmos.upgrade.v1beta1.QueryAppliedPlanRequest\x1a\x30.cosmos.upgrade.v1beta1.QueryAppliedPlanResponse\"3\x82\xd3\xe4\x93\x02-\x12+/cosmos/upgrade/v1beta1/applied_plan/{name}\x12\xdc\x01\n\x16UpgradedConsensusState\x12:.cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest\x1a;.cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse\"I\x88\x02\x01\x82\xd3\xe4\x93\x02@\x12>/cosmos/upgrade/v1beta1/upgraded_consensus_state/{last_height}\x12\xaa\x01\n\x0eModuleVersions\x12\x32.cosmos.upgrade.v1beta1.QueryModuleVersionsRequest\x1a\x33.cosmos.upgrade.v1beta1.QueryModuleVersionsResponse\"/\x82\xd3\xe4\x93\x02)\x12\'/cosmos/upgrade/v1beta1/module_versions\x12\x95\x01\n\tAuthority\x12-.cosmos.upgrade.v1beta1.QueryAuthorityRequest\x1a..cosmos.upgrade.v1beta1.QueryAuthorityResponse\")\x82\xd3\xe4\x93\x02#\x12!/cosmos/upgrade/v1beta1/authorityB\x1eZ\x1c\x63osmossdk.io/x/upgrade/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.upgrade.v1beta1.query_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\034cosmossdk.io/x/upgrade/types'
  _globals['_QUERYUPGRADEDCONSENSUSSTATEREQUEST']._loaded_options = None
  _globals['_QUERYUPGRADEDCONSENSUSSTATEREQUEST']._serialized_options = b'\030\001'
  _globals['_QUERYUPGRADEDCONSENSUSSTATERESPONSE']._loaded_options = None
  _globals['_QUERYUPGRADEDCONSENSUSSTATERESPONSE']._serialized_options = b'\030\001'
  _globals['_QUERY'].methods_by_name['CurrentPlan']._loaded_options = None
  _globals['_QUERY'].methods_by_name['CurrentPlan']._serialized_options = b'\202\323\344\223\002&\022$/cosmos/upgrade/v1beta1/current_plan'
  _globals['_QUERY'].methods_by_name['AppliedPlan']._loaded_options = None
  _globals['_QUERY'].methods_by_name['AppliedPlan']._serialized_options = b'\202\323\344\223\002-\022+/cosmos/upgrade/v1beta1/applied_plan/{name}'
  _globals['_QUERY'].methods_by_name['UpgradedConsensusState']._loaded_options = None
  _globals['_QUERY'].methods_by_name['UpgradedConsensusState']._serialized_options = b'\210\002\001\202\323\344\223\002@\022>/cosmos/upgrade/v1beta1/upgraded_consensus_state/{last_height}'
  _globals['_QUERY'].methods_by_name['ModuleVersions']._loaded_options = None
  _globals['_QUERY'].methods_by_name['ModuleVersions']._serialized_options = b'\202\323\344\223\002)\022\'/cosmos/upgrade/v1beta1/module_versions'
  _globals['_QUERY'].methods_by_name['Authority']._loaded_options = None
  _globals['_QUERY'].methods_by_name['Authority']._serialized_options = b'\202\323\344\223\002#\022!/cosmos/upgrade/v1beta1/authority'
  _globals['_QUERYCURRENTPLANREQUEST']._serialized_start=130
  _globals['_QUERYCURRENTPLANREQUEST']._serialized_end=155
  _globals['_QUERYCURRENTPLANRESPONSE']._serialized_start=157
  _globals['_QUERYCURRENTPLANRESPONSE']._serialized_end=233
  _globals['_QUERYAPPLIEDPLANREQUEST']._serialized_start=235
  _globals['_QUERYAPPLIEDPLANREQUEST']._serialized_end=280
  _globals['_QUERYAPPLIEDPLANRESPONSE']._serialized_start=282
  _globals['_QUERYAPPLIEDPLANRESPONSE']._serialized_end=332
  _globals['_QUERYUPGRADEDCONSENSUSSTATEREQUEST']._serialized_start=334
  _globals['_QUERYUPGRADEDCONSENSUSSTATEREQUEST']._serialized_end=407
  _globals['_QUERYUPGRADEDCONSENSUSSTATERESPONSE']._serialized_start=409
  _globals['_QUERYUPGRADEDCONSENSUSSTATERESPONSE']._serialized_end=514
  _globals['_QUERYMODULEVERSIONSREQUEST']._serialized_start=516
  _globals['_QUERYMODULEVERSIONSREQUEST']._serialized_end=577
  _globals['_QUERYMODULEVERSIONSRESPONSE']._serialized_start=579
  _globals['_QUERYMODULEVERSIONSRESPONSE']._serialized_end=688
  _globals['_QUERYAUTHORITYREQUEST']._serialized_start=690
  _globals['_QUERYAUTHORITYREQUEST']._serialized_end=713
  _globals['_QUERYAUTHORITYRESPONSE']._serialized_start=715
  _globals['_QUERYAUTHORITYRESPONSE']._serialized_end=765
  _globals['_QUERY']._serialized_start=768
  _globals['_QUERY']._serialized_end=1652
# @@protoc_insertion_point(module_scope)
