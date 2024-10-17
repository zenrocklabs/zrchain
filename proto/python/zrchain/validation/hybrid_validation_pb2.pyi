from amino import amino_pb2 as _amino_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from tendermint.types import types_pb2 as _types_pb2
from zrchain.validation import staking_pb2 as _staking_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ValidatorHV(_message.Message):
    __slots__ = ("operator_address", "consensus_pubkey", "jailed", "status", "tokensNative", "tokensAVS", "delegator_shares", "description", "unbonding_height", "unbonding_time", "commission", "min_self_delegation", "unbonding_on_hold_ref_count", "unbonding_ids")
    OPERATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_PUBKEY_FIELD_NUMBER: _ClassVar[int]
    JAILED_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    TOKENSNATIVE_FIELD_NUMBER: _ClassVar[int]
    TOKENSAVS_FIELD_NUMBER: _ClassVar[int]
    DELEGATOR_SHARES_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_TIME_FIELD_NUMBER: _ClassVar[int]
    COMMISSION_FIELD_NUMBER: _ClassVar[int]
    MIN_SELF_DELEGATION_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_ON_HOLD_REF_COUNT_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_IDS_FIELD_NUMBER: _ClassVar[int]
    operator_address: str
    consensus_pubkey: _any_pb2.Any
    jailed: bool
    status: _staking_pb2.BondStatus
    tokensNative: str
    tokensAVS: str
    delegator_shares: str
    description: _staking_pb2.Description
    unbonding_height: int
    unbonding_time: _timestamp_pb2.Timestamp
    commission: _staking_pb2.Commission
    min_self_delegation: str
    unbonding_on_hold_ref_count: int
    unbonding_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, operator_address: _Optional[str] = ..., consensus_pubkey: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., jailed: bool = ..., status: _Optional[_Union[_staking_pb2.BondStatus, str]] = ..., tokensNative: _Optional[str] = ..., tokensAVS: _Optional[str] = ..., delegator_shares: _Optional[str] = ..., description: _Optional[_Union[_staking_pb2.Description, _Mapping]] = ..., unbonding_height: _Optional[int] = ..., unbonding_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., commission: _Optional[_Union[_staking_pb2.Commission, _Mapping]] = ..., min_self_delegation: _Optional[str] = ..., unbonding_on_hold_ref_count: _Optional[int] = ..., unbonding_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class HistoricalInfoHV(_message.Message):
    __slots__ = ("header", "valset")
    HEADER_FIELD_NUMBER: _ClassVar[int]
    VALSET_FIELD_NUMBER: _ClassVar[int]
    header: _types_pb2.Header
    valset: _containers.RepeatedCompositeFieldContainer[ValidatorHV]
    def __init__(self, header: _Optional[_Union[_types_pb2.Header, _Mapping]] = ..., valset: _Optional[_Iterable[_Union[ValidatorHV, _Mapping]]] = ...) -> None: ...

class SlashEvent(_message.Message):
    __slots__ = ("blockHeight", "validatorAddr", "percentageSlashed", "tokensSlashedNative", "tokensSlashedAVS")
    BLOCKHEIGHT_FIELD_NUMBER: _ClassVar[int]
    VALIDATORADDR_FIELD_NUMBER: _ClassVar[int]
    PERCENTAGESLASHED_FIELD_NUMBER: _ClassVar[int]
    TOKENSSLASHEDNATIVE_FIELD_NUMBER: _ClassVar[int]
    TOKENSSLASHEDAVS_FIELD_NUMBER: _ClassVar[int]
    blockHeight: int
    validatorAddr: str
    percentageSlashed: str
    tokensSlashedNative: str
    tokensSlashedAVS: str
    def __init__(self, blockHeight: _Optional[int] = ..., validatorAddr: _Optional[str] = ..., percentageSlashed: _Optional[str] = ..., tokensSlashedNative: _Optional[str] = ..., tokensSlashedAVS: _Optional[str] = ...) -> None: ...

class HVParams(_message.Message):
    __slots__ = ("AVSRewardsRate", "BlockTime", "ZenBTCParams")
    AVSREWARDSRATE_FIELD_NUMBER: _ClassVar[int]
    BLOCKTIME_FIELD_NUMBER: _ClassVar[int]
    ZENBTCPARAMS_FIELD_NUMBER: _ClassVar[int]
    AVSRewardsRate: str
    BlockTime: int
    ZenBTCParams: ZenBTCParams
    def __init__(self, AVSRewardsRate: _Optional[str] = ..., BlockTime: _Optional[int] = ..., ZenBTCParams: _Optional[_Union[ZenBTCParams, _Mapping]] = ...) -> None: ...

class ZenBTCParams(_message.Message):
    __slots__ = ("zenBTCEthContractAddr", "zenBTCDepositKeyringAddr", "zenBTCMinterKeyID")
    ZENBTCETHCONTRACTADDR_FIELD_NUMBER: _ClassVar[int]
    ZENBTCDEPOSITKEYRINGADDR_FIELD_NUMBER: _ClassVar[int]
    ZENBTCMINTERKEYID_FIELD_NUMBER: _ClassVar[int]
    zenBTCEthContractAddr: str
    zenBTCDepositKeyringAddr: str
    zenBTCMinterKeyID: int
    def __init__(self, zenBTCEthContractAddr: _Optional[str] = ..., zenBTCDepositKeyringAddr: _Optional[str] = ..., zenBTCMinterKeyID: _Optional[int] = ...) -> None: ...

class ValidationInfo(_message.Message):
    __slots__ = ("non_voting_validators", "mismatched_vote_extensions")
    NON_VOTING_VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    MISMATCHED_VOTE_EXTENSIONS_FIELD_NUMBER: _ClassVar[int]
    non_voting_validators: _containers.RepeatedScalarFieldContainer[str]
    mismatched_vote_extensions: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, non_voting_validators: _Optional[_Iterable[str]] = ..., mismatched_vote_extensions: _Optional[_Iterable[str]] = ...) -> None: ...

class WithdrawalInfo(_message.Message):
    __slots__ = ("withdrawal_address", "amount", "retry_count")
    WITHDRAWAL_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    RETRY_COUNT_FIELD_NUMBER: _ClassVar[int]
    withdrawal_address: str
    amount: int
    retry_count: int
    def __init__(self, withdrawal_address: _Optional[str] = ..., amount: _Optional[int] = ..., retry_count: _Optional[int] = ...) -> None: ...
