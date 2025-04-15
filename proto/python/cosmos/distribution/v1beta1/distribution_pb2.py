# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cosmos/distribution/v1beta1/distribution.proto
# Protobuf Python Version: 6.30.2
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
    2,
    '',
    'cosmos/distribution/v1beta1/distribution.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from amino import amino_pb2 as amino_dot_amino__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n.cosmos/distribution/v1beta1/distribution.proto\x12\x1b\x63osmos.distribution.v1beta1\x1a\x14gogoproto/gogo.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x11\x61mino/amino.proto\"\x9a\x03\n\x06Params\x12[\n\rcommunity_tax\x18\x01 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x0c\x63ommunityTax\x12j\n\x14\x62\x61se_proposer_reward\x18\x02 \x01(\tB8\x18\x01\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x12\x62\x61seProposerReward\x12l\n\x15\x62onus_proposer_reward\x18\x03 \x01(\tB8\x18\x01\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x13\x62onusProposerReward\x12\x32\n\x15withdraw_addr_enabled\x18\x04 \x01(\x08R\x13withdrawAddrEnabled:%\x8a\xe7\xb0* cosmos-sdk/x/distribution/Params\"\xd6\x01\n\x1aValidatorHistoricalRewards\x12\x8e\x01\n\x17\x63umulative_reward_ratio\x18\x01 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB8\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoins\xa8\xe7\xb0*\x01R\x15\x63umulativeRewardRatio\x12\'\n\x0freference_count\x18\x02 \x01(\rR\x0ereferenceCount\"\xa3\x01\n\x17ValidatorCurrentRewards\x12p\n\x07rewards\x18\x01 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB8\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoins\xa8\xe7\xb0*\x01R\x07rewards\x12\x16\n\x06period\x18\x02 \x01(\x04R\x06period\"\x98\x01\n\x1eValidatorAccumulatedCommission\x12v\n\ncommission\x18\x01 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB8\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoins\xa8\xe7\xb0*\x01R\ncommission\"\x8f\x01\n\x1bValidatorOutstandingRewards\x12p\n\x07rewards\x18\x01 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB8\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoins\xa8\xe7\xb0*\x01R\x07rewards\"\x8f\x01\n\x13ValidatorSlashEvent\x12)\n\x10validator_period\x18\x01 \x01(\x04R\x0fvalidatorPeriod\x12M\n\x08\x66raction\x18\x02 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\x08\x66raction\"\x89\x01\n\x14ValidatorSlashEvents\x12q\n\x16validator_slash_events\x18\x01 \x03(\x0b\x32\x30.cosmos.distribution.v1beta1.ValidatorSlashEventB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x14validatorSlashEvents\"\x80\x02\n\x07\x46\x65\x65Pool\x12\x7f\n\x0e\x63ommunity_pool\x18\x01 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB:\x18\x01\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoins\xa8\xe7\xb0*\x01R\rcommunityPool\x12t\n\x0c\x64\x65\x63imal_pool\x18\x02 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB3\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoinsR\x0b\x64\x65\x63imalPool\"\x97\x02\n\x1a\x43ommunityPoolSpendProposal\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x1c\n\trecipient\x18\x03 \x01(\tR\trecipient\x12y\n\x06\x61mount\x18\x04 \x03(\x0b\x32\x19.cosmos.base.v1beta1.CoinBF\xc8\xde\x1f\x00\xaa\xdf\x1f(github.com/cosmos/cosmos-sdk/types.Coins\x9a\xe7\xb0*\x0clegacy_coins\xa8\xe7\xb0*\x01R\x06\x61mount:(\x18\x01\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\xca\xb4-\x1a\x63osmos.gov.v1beta1.Content\"\xd4\x01\n\x15\x44\x65legatorStartingInfo\x12\'\n\x0fprevious_period\x18\x01 \x01(\x04R\x0epreviousPeriod\x12L\n\x05stake\x18\x02 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x05stake\x12\x44\n\x06height\x18\x03 \x01(\x04\x42,\xea\xde\x1f\x0f\x63reation_height\xa2\xe7\xb0*\x0f\x63reation_height\xa8\xe7\xb0*\x01R\x06height\"\xe1\x01\n\x19\x44\x65legationDelegatorReward\x12N\n\x11validator_address\x18\x01 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x10validatorAddress\x12n\n\x06reward\x18\x02 \x03(\x0b\x32\x1c.cosmos.base.v1beta1.DecCoinB8\xc8\xde\x1f\x00\xaa\xdf\x1f+github.com/cosmos/cosmos-sdk/types.DecCoins\xa8\xe7\xb0*\x01R\x06reward:\x04\x88\xa0\x1f\x00\"\xd5\x01\n%CommunityPoolSpendProposalWithDeposit\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x02 \x01(\tR\x0b\x64\x65scription\x12\x1c\n\trecipient\x18\x03 \x01(\tR\trecipient\x12\x16\n\x06\x61mount\x18\x04 \x01(\tR\x06\x61mount\x12\x18\n\x07\x64\x65posit\x18\x05 \x01(\tR\x07\x64\x65posit:$\x18\x01\x88\xa0\x1f\x00\xca\xb4-\x1a\x63osmos.gov.v1beta1.ContentB\'Z!cosmossdk.io/x/distribution/types\xa8\xe2\x1e\x01\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cosmos.distribution.v1beta1.distribution_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z!cosmossdk.io/x/distribution/types\250\342\036\001'
  _globals['_PARAMS'].fields_by_name['community_tax']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['community_tax']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['base_proposer_reward']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['base_proposer_reward']._serialized_options = b'\030\001\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['bonus_proposer_reward']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['bonus_proposer_reward']._serialized_options = b'\030\001\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\212\347\260* cosmos-sdk/x/distribution/Params'
  _globals['_VALIDATORHISTORICALREWARDS'].fields_by_name['cumulative_reward_ratio']._loaded_options = None
  _globals['_VALIDATORHISTORICALREWARDS'].fields_by_name['cumulative_reward_ratio']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins\250\347\260*\001'
  _globals['_VALIDATORCURRENTREWARDS'].fields_by_name['rewards']._loaded_options = None
  _globals['_VALIDATORCURRENTREWARDS'].fields_by_name['rewards']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins\250\347\260*\001'
  _globals['_VALIDATORACCUMULATEDCOMMISSION'].fields_by_name['commission']._loaded_options = None
  _globals['_VALIDATORACCUMULATEDCOMMISSION'].fields_by_name['commission']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins\250\347\260*\001'
  _globals['_VALIDATOROUTSTANDINGREWARDS'].fields_by_name['rewards']._loaded_options = None
  _globals['_VALIDATOROUTSTANDINGREWARDS'].fields_by_name['rewards']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins\250\347\260*\001'
  _globals['_VALIDATORSLASHEVENT'].fields_by_name['fraction']._loaded_options = None
  _globals['_VALIDATORSLASHEVENT'].fields_by_name['fraction']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_VALIDATORSLASHEVENTS'].fields_by_name['validator_slash_events']._loaded_options = None
  _globals['_VALIDATORSLASHEVENTS'].fields_by_name['validator_slash_events']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_FEEPOOL'].fields_by_name['community_pool']._loaded_options = None
  _globals['_FEEPOOL'].fields_by_name['community_pool']._serialized_options = b'\030\001\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins\250\347\260*\001'
  _globals['_FEEPOOL'].fields_by_name['decimal_pool']._loaded_options = None
  _globals['_FEEPOOL'].fields_by_name['decimal_pool']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins'
  _globals['_COMMUNITYPOOLSPENDPROPOSAL'].fields_by_name['amount']._loaded_options = None
  _globals['_COMMUNITYPOOLSPENDPROPOSAL'].fields_by_name['amount']._serialized_options = b'\310\336\037\000\252\337\037(github.com/cosmos/cosmos-sdk/types.Coins\232\347\260*\014legacy_coins\250\347\260*\001'
  _globals['_COMMUNITYPOOLSPENDPROPOSAL']._loaded_options = None
  _globals['_COMMUNITYPOOLSPENDPROPOSAL']._serialized_options = b'\030\001\210\240\037\000\350\240\037\000\312\264-\032cosmos.gov.v1beta1.Content'
  _globals['_DELEGATORSTARTINGINFO'].fields_by_name['stake']._loaded_options = None
  _globals['_DELEGATORSTARTINGINFO'].fields_by_name['stake']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_DELEGATORSTARTINGINFO'].fields_by_name['height']._loaded_options = None
  _globals['_DELEGATORSTARTINGINFO'].fields_by_name['height']._serialized_options = b'\352\336\037\017creation_height\242\347\260*\017creation_height\250\347\260*\001'
  _globals['_DELEGATIONDELEGATORREWARD'].fields_by_name['validator_address']._loaded_options = None
  _globals['_DELEGATIONDELEGATORREWARD'].fields_by_name['validator_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_DELEGATIONDELEGATORREWARD'].fields_by_name['reward']._loaded_options = None
  _globals['_DELEGATIONDELEGATORREWARD'].fields_by_name['reward']._serialized_options = b'\310\336\037\000\252\337\037+github.com/cosmos/cosmos-sdk/types.DecCoins\250\347\260*\001'
  _globals['_DELEGATIONDELEGATORREWARD']._loaded_options = None
  _globals['_DELEGATIONDELEGATORREWARD']._serialized_options = b'\210\240\037\000'
  _globals['_COMMUNITYPOOLSPENDPROPOSALWITHDEPOSIT']._loaded_options = None
  _globals['_COMMUNITYPOOLSPENDPROPOSALWITHDEPOSIT']._serialized_options = b'\030\001\210\240\037\000\312\264-\032cosmos.gov.v1beta1.Content'
  _globals['_PARAMS']._serialized_start=180
  _globals['_PARAMS']._serialized_end=590
  _globals['_VALIDATORHISTORICALREWARDS']._serialized_start=593
  _globals['_VALIDATORHISTORICALREWARDS']._serialized_end=807
  _globals['_VALIDATORCURRENTREWARDS']._serialized_start=810
  _globals['_VALIDATORCURRENTREWARDS']._serialized_end=973
  _globals['_VALIDATORACCUMULATEDCOMMISSION']._serialized_start=976
  _globals['_VALIDATORACCUMULATEDCOMMISSION']._serialized_end=1128
  _globals['_VALIDATOROUTSTANDINGREWARDS']._serialized_start=1131
  _globals['_VALIDATOROUTSTANDINGREWARDS']._serialized_end=1274
  _globals['_VALIDATORSLASHEVENT']._serialized_start=1277
  _globals['_VALIDATORSLASHEVENT']._serialized_end=1420
  _globals['_VALIDATORSLASHEVENTS']._serialized_start=1423
  _globals['_VALIDATORSLASHEVENTS']._serialized_end=1560
  _globals['_FEEPOOL']._serialized_start=1563
  _globals['_FEEPOOL']._serialized_end=1819
  _globals['_COMMUNITYPOOLSPENDPROPOSAL']._serialized_start=1822
  _globals['_COMMUNITYPOOLSPENDPROPOSAL']._serialized_end=2101
  _globals['_DELEGATORSTARTINGINFO']._serialized_start=2104
  _globals['_DELEGATORSTARTINGINFO']._serialized_end=2316
  _globals['_DELEGATIONDELEGATORREWARD']._serialized_start=2319
  _globals['_DELEGATIONDELEGATORREWARD']._serialized_end=2544
  _globals['_COMMUNITYPOOLSPENDPROPOSALWITHDEPOSIT']._serialized_start=2547
  _globals['_COMMUNITYPOOLSPENDPROPOSALWITHDEPOSIT']._serialized_end=2760
# @@protoc_insertion_point(module_scope)
