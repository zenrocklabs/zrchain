# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: zrchain/validation/staking.proto
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
    'zrchain/validation/staking.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from gogoproto import gogo_pb2 as gogoproto_dot_gogo__pb2
from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2
from google.protobuf import duration_pb2 as google_dot_protobuf_dot_duration__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from cosmos_proto import cosmos_pb2 as cosmos__proto_dot_cosmos__pb2
from cosmos.base.v1beta1 import coin_pb2 as cosmos_dot_base_dot_v1beta1_dot_coin__pb2
from amino import amino_pb2 as amino_dot_amino__pb2
from tendermint.types import types_pb2 as tendermint_dot_types_dot_types__pb2
from tendermint.abci import types_pb2 as tendermint_dot_abci_dot_types__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n zrchain/validation/staking.proto\x12\x12zrchain.validation\x1a\x14gogoproto/gogo.proto\x1a\x19google/protobuf/any.proto\x1a\x1egoogle/protobuf/duration.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x19\x63osmos_proto/cosmos.proto\x1a\x1e\x63osmos/base/v1beta1/coin.proto\x1a\x11\x61mino/amino.proto\x1a\x1ctendermint/types/types.proto\x1a\x1btendermint/abci/types.proto\"\x96\x02\n\x0f\x43ommissionRates\x12J\n\x04rate\x18\x01 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x04rate\x12Q\n\x08max_rate\x18\x02 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x07maxRate\x12^\n\x0fmax_change_rate\x18\x03 \x01(\tB6\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\rmaxChangeRate:\x04\xe8\xa0\x1f\x01\"\xbd\x01\n\nCommission\x12]\n\x10\x63ommission_rates\x18\x01 \x01(\x0b\x32#.zrchain.validation.CommissionRatesB\r\xc8\xde\x1f\x00\xd0\xde\x1f\x01\xa8\xe7\xb0*\x01R\x0f\x63ommissionRates\x12J\n\x0bupdate_time\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\r\xc8\xde\x1f\x00\x90\xdf\x1f\x01\xa8\xe7\xb0*\x01R\nupdateTime:\x04\xe8\xa0\x1f\x01\"\xa8\x01\n\x0b\x44\x65scription\x12\x18\n\x07moniker\x18\x01 \x01(\tR\x07moniker\x12\x1a\n\x08identity\x18\x02 \x01(\tR\x08identity\x12\x18\n\x07website\x18\x03 \x01(\tR\x07website\x12)\n\x10security_contact\x18\x04 \x01(\tR\x0fsecurityContact\x12\x18\n\x07\x64\x65tails\x18\x05 \x01(\tR\x07\x64\x65tails:\x04\xe8\xa0\x1f\x01\"F\n\x0cValAddresses\x12\x36\n\taddresses\x18\x01 \x03(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\taddresses\"\xa9\x01\n\x06\x44VPair\x12\x45\n\x11\x64\x65legator_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10\x64\x65legatorAddress\x12N\n\x11validator_address\x18\x02 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x10validatorAddress:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"F\n\x07\x44VPairs\x12;\n\x05pairs\x18\x01 \x03(\x0b\x32\x1a.zrchain.validation.DVPairB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x05pairs\"\x8b\x02\n\nDVVTriplet\x12\x45\n\x11\x64\x65legator_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10\x64\x65legatorAddress\x12U\n\x15validator_src_address\x18\x02 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x13validatorSrcAddress\x12U\n\x15validator_dst_address\x18\x03 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x13validatorDstAddress:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"T\n\x0b\x44VVTriplets\x12\x45\n\x08triplets\x18\x01 \x03(\x0b\x32\x1e.zrchain.validation.DVVTripletB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x08triplets\"\xf8\x01\n\nDelegation\x12\x45\n\x11\x64\x65legator_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10\x64\x65legatorAddress\x12N\n\x11validator_address\x18\x02 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x10validatorAddress\x12I\n\x06shares\x18\x03 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\x06shares:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"\x89\x02\n\x13UnbondingDelegation\x12\x45\n\x11\x64\x65legator_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10\x64\x65legatorAddress\x12N\n\x11validator_address\x18\x02 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x10validatorAddress\x12Q\n\x07\x65ntries\x18\x03 \x03(\x0b\x32,.zrchain.validation.UnbondingDelegationEntryB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x07\x65ntries:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"\x9b\x03\n\x18UnbondingDelegationEntry\x12\'\n\x0f\x63reation_height\x18\x01 \x01(\x03R\x0e\x63reationHeight\x12R\n\x0f\x63ompletion_time\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\r\xc8\xde\x1f\x00\x90\xdf\x1f\x01\xa8\xe7\xb0*\x01R\x0e\x63ompletionTime\x12T\n\x0finitial_balance\x18\x03 \x01(\tB+\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xd2\xb4-\ncosmos.IntR\x0einitialBalance\x12\x45\n\x07\x62\x61lance\x18\x04 \x01(\tB+\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xd2\xb4-\ncosmos.IntR\x07\x62\x61lance\x12!\n\x0cunbonding_id\x18\x05 \x01(\x04R\x0bunbondingId\x12<\n\x1bunbonding_on_hold_ref_count\x18\x06 \x01(\x03R\x17unbondingOnHoldRefCount:\x04\xe8\xa0\x1f\x01\"\x9f\x03\n\x11RedelegationEntry\x12\'\n\x0f\x63reation_height\x18\x01 \x01(\x03R\x0e\x63reationHeight\x12R\n\x0f\x63ompletion_time\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.TimestampB\r\xc8\xde\x1f\x00\x90\xdf\x1f\x01\xa8\xe7\xb0*\x01R\x0e\x63ompletionTime\x12T\n\x0finitial_balance\x18\x03 \x01(\tB+\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xd2\xb4-\ncosmos.IntR\x0einitialBalance\x12P\n\nshares_dst\x18\x04 \x01(\tB1\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xd2\xb4-\ncosmos.DecR\tsharesDst\x12!\n\x0cunbonding_id\x18\x05 \x01(\x04R\x0bunbondingId\x12<\n\x1bunbonding_on_hold_ref_count\x18\x06 \x01(\x03R\x17unbondingOnHoldRefCount:\x04\xe8\xa0\x1f\x01\"\xd9\x02\n\x0cRedelegation\x12\x45\n\x11\x64\x65legator_address\x18\x01 \x01(\tB\x18\xd2\xb4-\x14\x63osmos.AddressStringR\x10\x64\x65legatorAddress\x12U\n\x15validator_src_address\x18\x02 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x13validatorSrcAddress\x12U\n\x15validator_dst_address\x18\x03 \x01(\tB!\xd2\xb4-\x1d\x63osmos.ValidatorAddressStringR\x13validatorDstAddress\x12J\n\x07\x65ntries\x18\x04 \x03(\x0b\x32%.zrchain.validation.RedelegationEntryB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x07\x65ntries:\x08\x88\xa0\x1f\x00\xe8\xa0\x1f\x00\"\x9c\x03\n\x06Params\x12O\n\x0eunbonding_time\x18\x01 \x01(\x0b\x32\x19.google.protobuf.DurationB\r\xc8\xde\x1f\x00\x98\xdf\x1f\x01\xa8\xe7\xb0*\x01R\runbondingTime\x12%\n\x0emax_validators\x18\x02 \x01(\rR\rmaxValidators\x12\x1f\n\x0bmax_entries\x18\x03 \x01(\rR\nmaxEntries\x12-\n\x12historical_entries\x18\x04 \x01(\rR\x11historicalEntries\x12\x1d\n\nbond_denom\x18\x05 \x01(\tR\tbondDenom\x12\x84\x01\n\x13min_commission_rate\x18\x06 \x01(\tBT\xc8\xde\x1f\x00\xda\xde\x1f\x1b\x63osmossdk.io/math.LegacyDec\xf2\xde\x1f\x1ayaml:\"min_commission_rate\"\xd2\xb4-\ncosmos.Dec\xa8\xe7\xb0*\x01R\x11minCommissionRate:$\xe8\xa0\x1f\x01\x8a\xe7\xb0*\x1bzrchain/x/validation/Params\"\xa5\x01\n\x12\x44\x65legationResponse\x12I\n\ndelegation\x18\x01 \x01(\x0b\x32\x1e.zrchain.validation.DelegationB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\ndelegation\x12>\n\x07\x62\x61lance\x18\x02 \x01(\x0b\x32\x19.cosmos.base.v1beta1.CoinB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x07\x62\x61lance:\x04\xe8\xa0\x1f\x00\"\xc9\x01\n\x19RedelegationEntryResponse\x12_\n\x12redelegation_entry\x18\x01 \x01(\x0b\x32%.zrchain.validation.RedelegationEntryB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x11redelegationEntry\x12\x45\n\x07\x62\x61lance\x18\x04 \x01(\tB+\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xd2\xb4-\ncosmos.IntR\x07\x62\x61lance:\x04\xe8\xa0\x1f\x01\"\xc1\x01\n\x14RedelegationResponse\x12O\n\x0credelegation\x18\x01 \x01(\x0b\x32 .zrchain.validation.RedelegationB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x0credelegation\x12R\n\x07\x65ntries\x18\x02 \x03(\x0b\x32-.zrchain.validation.RedelegationEntryResponseB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x07\x65ntries:\x04\xe8\xa0\x1f\x00\"\xeb\x01\n\x04Pool\x12q\n\x11not_bonded_tokens\x18\x01 \x01(\tBE\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xea\xde\x1f\x11not_bonded_tokens\xd2\xb4-\ncosmos.Int\xa8\xe7\xb0*\x01R\x0fnotBondedTokens\x12\x66\n\rbonded_tokens\x18\x02 \x01(\tBA\xc8\xde\x1f\x00\xda\xde\x1f\x15\x63osmossdk.io/math.Int\xea\xde\x1f\rbonded_tokens\xd2\xb4-\ncosmos.Int\xa8\xe7\xb0*\x01R\x0c\x62ondedTokens:\x08\xe8\xa0\x1f\x01\xf0\xa0\x1f\x01\"Y\n\x10ValidatorUpdates\x12\x45\n\x07updates\x18\x01 \x03(\x0b\x32 .tendermint.abci.ValidatorUpdateB\t\xc8\xde\x1f\x00\xa8\xe7\xb0*\x01R\x07updates*\xb6\x01\n\nBondStatus\x12,\n\x17\x42OND_STATUS_UNSPECIFIED\x10\x00\x1a\x0f\x8a\x9d \x0bUnspecified\x12&\n\x14\x42OND_STATUS_UNBONDED\x10\x01\x1a\x0c\x8a\x9d \x08Unbonded\x12(\n\x15\x42OND_STATUS_UNBONDING\x10\x02\x1a\r\x8a\x9d \tUnbonding\x12\"\n\x12\x42OND_STATUS_BONDED\x10\x03\x1a\n\x8a\x9d \x06\x42onded\x1a\x04\x88\xa3\x1e\x00*]\n\nInfraction\x12\x1a\n\x16INFRACTION_UNSPECIFIED\x10\x00\x12\x1a\n\x16INFRACTION_DOUBLE_SIGN\x10\x01\x12\x17\n\x13INFRACTION_DOWNTIME\x10\x02\x42=Z;github.com/Zenrock-Foundation/zrchain/v5/x/validation/typesb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'zrchain.validation.staking_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z;github.com/Zenrock-Foundation/zrchain/v5/x/validation/types'
  _globals['_BONDSTATUS']._loaded_options = None
  _globals['_BONDSTATUS']._serialized_options = b'\210\243\036\000'
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_UNSPECIFIED"]._loaded_options = None
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_UNSPECIFIED"]._serialized_options = b'\212\235 \013Unspecified'
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_UNBONDED"]._loaded_options = None
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_UNBONDED"]._serialized_options = b'\212\235 \010Unbonded'
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_UNBONDING"]._loaded_options = None
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_UNBONDING"]._serialized_options = b'\212\235 \tUnbonding'
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_BONDED"]._loaded_options = None
  _globals['_BONDSTATUS'].values_by_name["BOND_STATUS_BONDED"]._serialized_options = b'\212\235 \006Bonded'
  _globals['_COMMISSIONRATES'].fields_by_name['rate']._loaded_options = None
  _globals['_COMMISSIONRATES'].fields_by_name['rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_COMMISSIONRATES'].fields_by_name['max_rate']._loaded_options = None
  _globals['_COMMISSIONRATES'].fields_by_name['max_rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_COMMISSIONRATES'].fields_by_name['max_change_rate']._loaded_options = None
  _globals['_COMMISSIONRATES'].fields_by_name['max_change_rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_COMMISSIONRATES']._loaded_options = None
  _globals['_COMMISSIONRATES']._serialized_options = b'\350\240\037\001'
  _globals['_COMMISSION'].fields_by_name['commission_rates']._loaded_options = None
  _globals['_COMMISSION'].fields_by_name['commission_rates']._serialized_options = b'\310\336\037\000\320\336\037\001\250\347\260*\001'
  _globals['_COMMISSION'].fields_by_name['update_time']._loaded_options = None
  _globals['_COMMISSION'].fields_by_name['update_time']._serialized_options = b'\310\336\037\000\220\337\037\001\250\347\260*\001'
  _globals['_COMMISSION']._loaded_options = None
  _globals['_COMMISSION']._serialized_options = b'\350\240\037\001'
  _globals['_DESCRIPTION']._loaded_options = None
  _globals['_DESCRIPTION']._serialized_options = b'\350\240\037\001'
  _globals['_VALADDRESSES'].fields_by_name['addresses']._loaded_options = None
  _globals['_VALADDRESSES'].fields_by_name['addresses']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_DVPAIR'].fields_by_name['delegator_address']._loaded_options = None
  _globals['_DVPAIR'].fields_by_name['delegator_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_DVPAIR'].fields_by_name['validator_address']._loaded_options = None
  _globals['_DVPAIR'].fields_by_name['validator_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_DVPAIR']._loaded_options = None
  _globals['_DVPAIR']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_DVPAIRS'].fields_by_name['pairs']._loaded_options = None
  _globals['_DVPAIRS'].fields_by_name['pairs']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_DVVTRIPLET'].fields_by_name['delegator_address']._loaded_options = None
  _globals['_DVVTRIPLET'].fields_by_name['delegator_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_DVVTRIPLET'].fields_by_name['validator_src_address']._loaded_options = None
  _globals['_DVVTRIPLET'].fields_by_name['validator_src_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_DVVTRIPLET'].fields_by_name['validator_dst_address']._loaded_options = None
  _globals['_DVVTRIPLET'].fields_by_name['validator_dst_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_DVVTRIPLET']._loaded_options = None
  _globals['_DVVTRIPLET']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_DVVTRIPLETS'].fields_by_name['triplets']._loaded_options = None
  _globals['_DVVTRIPLETS'].fields_by_name['triplets']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_DELEGATION'].fields_by_name['delegator_address']._loaded_options = None
  _globals['_DELEGATION'].fields_by_name['delegator_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_DELEGATION'].fields_by_name['validator_address']._loaded_options = None
  _globals['_DELEGATION'].fields_by_name['validator_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_DELEGATION'].fields_by_name['shares']._loaded_options = None
  _globals['_DELEGATION'].fields_by_name['shares']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_DELEGATION']._loaded_options = None
  _globals['_DELEGATION']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_UNBONDINGDELEGATION'].fields_by_name['delegator_address']._loaded_options = None
  _globals['_UNBONDINGDELEGATION'].fields_by_name['delegator_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_UNBONDINGDELEGATION'].fields_by_name['validator_address']._loaded_options = None
  _globals['_UNBONDINGDELEGATION'].fields_by_name['validator_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_UNBONDINGDELEGATION'].fields_by_name['entries']._loaded_options = None
  _globals['_UNBONDINGDELEGATION'].fields_by_name['entries']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_UNBONDINGDELEGATION']._loaded_options = None
  _globals['_UNBONDINGDELEGATION']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_UNBONDINGDELEGATIONENTRY'].fields_by_name['completion_time']._loaded_options = None
  _globals['_UNBONDINGDELEGATIONENTRY'].fields_by_name['completion_time']._serialized_options = b'\310\336\037\000\220\337\037\001\250\347\260*\001'
  _globals['_UNBONDINGDELEGATIONENTRY'].fields_by_name['initial_balance']._loaded_options = None
  _globals['_UNBONDINGDELEGATIONENTRY'].fields_by_name['initial_balance']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\322\264-\ncosmos.Int'
  _globals['_UNBONDINGDELEGATIONENTRY'].fields_by_name['balance']._loaded_options = None
  _globals['_UNBONDINGDELEGATIONENTRY'].fields_by_name['balance']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\322\264-\ncosmos.Int'
  _globals['_UNBONDINGDELEGATIONENTRY']._loaded_options = None
  _globals['_UNBONDINGDELEGATIONENTRY']._serialized_options = b'\350\240\037\001'
  _globals['_REDELEGATIONENTRY'].fields_by_name['completion_time']._loaded_options = None
  _globals['_REDELEGATIONENTRY'].fields_by_name['completion_time']._serialized_options = b'\310\336\037\000\220\337\037\001\250\347\260*\001'
  _globals['_REDELEGATIONENTRY'].fields_by_name['initial_balance']._loaded_options = None
  _globals['_REDELEGATIONENTRY'].fields_by_name['initial_balance']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\322\264-\ncosmos.Int'
  _globals['_REDELEGATIONENTRY'].fields_by_name['shares_dst']._loaded_options = None
  _globals['_REDELEGATIONENTRY'].fields_by_name['shares_dst']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\322\264-\ncosmos.Dec'
  _globals['_REDELEGATIONENTRY']._loaded_options = None
  _globals['_REDELEGATIONENTRY']._serialized_options = b'\350\240\037\001'
  _globals['_REDELEGATION'].fields_by_name['delegator_address']._loaded_options = None
  _globals['_REDELEGATION'].fields_by_name['delegator_address']._serialized_options = b'\322\264-\024cosmos.AddressString'
  _globals['_REDELEGATION'].fields_by_name['validator_src_address']._loaded_options = None
  _globals['_REDELEGATION'].fields_by_name['validator_src_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_REDELEGATION'].fields_by_name['validator_dst_address']._loaded_options = None
  _globals['_REDELEGATION'].fields_by_name['validator_dst_address']._serialized_options = b'\322\264-\035cosmos.ValidatorAddressString'
  _globals['_REDELEGATION'].fields_by_name['entries']._loaded_options = None
  _globals['_REDELEGATION'].fields_by_name['entries']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_REDELEGATION']._loaded_options = None
  _globals['_REDELEGATION']._serialized_options = b'\210\240\037\000\350\240\037\000'
  _globals['_PARAMS'].fields_by_name['unbonding_time']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['unbonding_time']._serialized_options = b'\310\336\037\000\230\337\037\001\250\347\260*\001'
  _globals['_PARAMS'].fields_by_name['min_commission_rate']._loaded_options = None
  _globals['_PARAMS'].fields_by_name['min_commission_rate']._serialized_options = b'\310\336\037\000\332\336\037\033cosmossdk.io/math.LegacyDec\362\336\037\032yaml:\"min_commission_rate\"\322\264-\ncosmos.Dec\250\347\260*\001'
  _globals['_PARAMS']._loaded_options = None
  _globals['_PARAMS']._serialized_options = b'\350\240\037\001\212\347\260*\033zrchain/x/validation/Params'
  _globals['_DELEGATIONRESPONSE'].fields_by_name['delegation']._loaded_options = None
  _globals['_DELEGATIONRESPONSE'].fields_by_name['delegation']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_DELEGATIONRESPONSE'].fields_by_name['balance']._loaded_options = None
  _globals['_DELEGATIONRESPONSE'].fields_by_name['balance']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_DELEGATIONRESPONSE']._loaded_options = None
  _globals['_DELEGATIONRESPONSE']._serialized_options = b'\350\240\037\000'
  _globals['_REDELEGATIONENTRYRESPONSE'].fields_by_name['redelegation_entry']._loaded_options = None
  _globals['_REDELEGATIONENTRYRESPONSE'].fields_by_name['redelegation_entry']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_REDELEGATIONENTRYRESPONSE'].fields_by_name['balance']._loaded_options = None
  _globals['_REDELEGATIONENTRYRESPONSE'].fields_by_name['balance']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\322\264-\ncosmos.Int'
  _globals['_REDELEGATIONENTRYRESPONSE']._loaded_options = None
  _globals['_REDELEGATIONENTRYRESPONSE']._serialized_options = b'\350\240\037\001'
  _globals['_REDELEGATIONRESPONSE'].fields_by_name['redelegation']._loaded_options = None
  _globals['_REDELEGATIONRESPONSE'].fields_by_name['redelegation']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_REDELEGATIONRESPONSE'].fields_by_name['entries']._loaded_options = None
  _globals['_REDELEGATIONRESPONSE'].fields_by_name['entries']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_REDELEGATIONRESPONSE']._loaded_options = None
  _globals['_REDELEGATIONRESPONSE']._serialized_options = b'\350\240\037\000'
  _globals['_POOL'].fields_by_name['not_bonded_tokens']._loaded_options = None
  _globals['_POOL'].fields_by_name['not_bonded_tokens']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\352\336\037\021not_bonded_tokens\322\264-\ncosmos.Int\250\347\260*\001'
  _globals['_POOL'].fields_by_name['bonded_tokens']._loaded_options = None
  _globals['_POOL'].fields_by_name['bonded_tokens']._serialized_options = b'\310\336\037\000\332\336\037\025cosmossdk.io/math.Int\352\336\037\rbonded_tokens\322\264-\ncosmos.Int\250\347\260*\001'
  _globals['_POOL']._loaded_options = None
  _globals['_POOL']._serialized_options = b'\350\240\037\001\360\240\037\001'
  _globals['_VALIDATORUPDATES'].fields_by_name['updates']._loaded_options = None
  _globals['_VALIDATORUPDATES'].fields_by_name['updates']._serialized_options = b'\310\336\037\000\250\347\260*\001'
  _globals['_BONDSTATUS']._serialized_start=4635
  _globals['_BONDSTATUS']._serialized_end=4817
  _globals['_INFRACTION']._serialized_start=4819
  _globals['_INFRACTION']._serialized_end=4912
  _globals['_COMMISSIONRATES']._serialized_start=308
  _globals['_COMMISSIONRATES']._serialized_end=586
  _globals['_COMMISSION']._serialized_start=589
  _globals['_COMMISSION']._serialized_end=778
  _globals['_DESCRIPTION']._serialized_start=781
  _globals['_DESCRIPTION']._serialized_end=949
  _globals['_VALADDRESSES']._serialized_start=951
  _globals['_VALADDRESSES']._serialized_end=1021
  _globals['_DVPAIR']._serialized_start=1024
  _globals['_DVPAIR']._serialized_end=1193
  _globals['_DVPAIRS']._serialized_start=1195
  _globals['_DVPAIRS']._serialized_end=1265
  _globals['_DVVTRIPLET']._serialized_start=1268
  _globals['_DVVTRIPLET']._serialized_end=1535
  _globals['_DVVTRIPLETS']._serialized_start=1537
  _globals['_DVVTRIPLETS']._serialized_end=1621
  _globals['_DELEGATION']._serialized_start=1624
  _globals['_DELEGATION']._serialized_end=1872
  _globals['_UNBONDINGDELEGATION']._serialized_start=1875
  _globals['_UNBONDINGDELEGATION']._serialized_end=2140
  _globals['_UNBONDINGDELEGATIONENTRY']._serialized_start=2143
  _globals['_UNBONDINGDELEGATIONENTRY']._serialized_end=2554
  _globals['_REDELEGATIONENTRY']._serialized_start=2557
  _globals['_REDELEGATIONENTRY']._serialized_end=2972
  _globals['_REDELEGATION']._serialized_start=2975
  _globals['_REDELEGATION']._serialized_end=3320
  _globals['_PARAMS']._serialized_start=3323
  _globals['_PARAMS']._serialized_end=3735
  _globals['_DELEGATIONRESPONSE']._serialized_start=3738
  _globals['_DELEGATIONRESPONSE']._serialized_end=3903
  _globals['_REDELEGATIONENTRYRESPONSE']._serialized_start=3906
  _globals['_REDELEGATIONENTRYRESPONSE']._serialized_end=4107
  _globals['_REDELEGATIONRESPONSE']._serialized_start=4110
  _globals['_REDELEGATIONRESPONSE']._serialized_end=4303
  _globals['_POOL']._serialized_start=4306
  _globals['_POOL']._serialized_end=4541
  _globals['_VALIDATORUPDATES']._serialized_start=4543
  _globals['_VALIDATORUPDATES']._serialized_end=4632
# @@protoc_insertion_point(module_scope)
