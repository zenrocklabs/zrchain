from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Minter(_message.Message):
    __slots__ = ("inflation", "annual_provisions")
    INFLATION_FIELD_NUMBER: _ClassVar[int]
    ANNUAL_PROVISIONS_FIELD_NUMBER: _ClassVar[int]
    inflation: str
    annual_provisions: str
    def __init__(self, inflation: _Optional[str] = ..., annual_provisions: _Optional[str] = ...) -> None: ...

class Params(_message.Message):
    __slots__ = ("mint_denom", "inflation_rate_change", "inflation_max", "inflation_min", "goal_bonded", "blocks_per_year", "staking_yield", "burn_rate", "protocol_wallet_rate", "protocol_wallet_address", "additional_staking_rewards", "additional_mpc_rewards", "additional_burn_rate")
    MINT_DENOM_FIELD_NUMBER: _ClassVar[int]
    INFLATION_RATE_CHANGE_FIELD_NUMBER: _ClassVar[int]
    INFLATION_MAX_FIELD_NUMBER: _ClassVar[int]
    INFLATION_MIN_FIELD_NUMBER: _ClassVar[int]
    GOAL_BONDED_FIELD_NUMBER: _ClassVar[int]
    BLOCKS_PER_YEAR_FIELD_NUMBER: _ClassVar[int]
    STAKING_YIELD_FIELD_NUMBER: _ClassVar[int]
    BURN_RATE_FIELD_NUMBER: _ClassVar[int]
    PROTOCOL_WALLET_RATE_FIELD_NUMBER: _ClassVar[int]
    PROTOCOL_WALLET_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    ADDITIONAL_STAKING_REWARDS_FIELD_NUMBER: _ClassVar[int]
    ADDITIONAL_MPC_REWARDS_FIELD_NUMBER: _ClassVar[int]
    ADDITIONAL_BURN_RATE_FIELD_NUMBER: _ClassVar[int]
    mint_denom: str
    inflation_rate_change: str
    inflation_max: str
    inflation_min: str
    goal_bonded: str
    blocks_per_year: int
    staking_yield: str
    burn_rate: str
    protocol_wallet_rate: str
    protocol_wallet_address: str
    additional_staking_rewards: str
    additional_mpc_rewards: str
    additional_burn_rate: str
    def __init__(self, mint_denom: _Optional[str] = ..., inflation_rate_change: _Optional[str] = ..., inflation_max: _Optional[str] = ..., inflation_min: _Optional[str] = ..., goal_bonded: _Optional[str] = ..., blocks_per_year: _Optional[int] = ..., staking_yield: _Optional[str] = ..., burn_rate: _Optional[str] = ..., protocol_wallet_rate: _Optional[str] = ..., protocol_wallet_address: _Optional[str] = ..., additional_staking_rewards: _Optional[str] = ..., additional_mpc_rewards: _Optional[str] = ..., additional_burn_rate: _Optional[str] = ...) -> None: ...
