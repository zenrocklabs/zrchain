from amino import amino_pb2 as _amino_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from cosmos.query.v1 import query_pb2 as _query_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from zrchain.validation import hybrid_validation_pb2 as _hybrid_validation_pb2
from zrchain.validation import staking_pb2 as _staking_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryValidatorsRequest(_message.Message):
    __slots__ = ("status", "pagination")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    status: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, status: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryValidatorsResponse(_message.Message):
    __slots__ = ("validators", "pagination")
    VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    validators: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.ValidatorHV]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, validators: _Optional[_Iterable[_Union[_hybrid_validation_pb2.ValidatorHV, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryValidatorRequest(_message.Message):
    __slots__ = ("validator_addr",)
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    validator_addr: str
    def __init__(self, validator_addr: _Optional[str] = ...) -> None: ...

class QueryValidatorResponse(_message.Message):
    __slots__ = ("validator",)
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    validator: _hybrid_validation_pb2.ValidatorHV
    def __init__(self, validator: _Optional[_Union[_hybrid_validation_pb2.ValidatorHV, _Mapping]] = ...) -> None: ...

class QueryValidatorDelegationsRequest(_message.Message):
    __slots__ = ("validator_addr", "pagination")
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    validator_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, validator_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryValidatorDelegationsResponse(_message.Message):
    __slots__ = ("delegation_responses", "pagination")
    DELEGATION_RESPONSES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    delegation_responses: _containers.RepeatedCompositeFieldContainer[_staking_pb2.DelegationResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, delegation_responses: _Optional[_Iterable[_Union[_staking_pb2.DelegationResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryValidatorUnbondingDelegationsRequest(_message.Message):
    __slots__ = ("validator_addr", "pagination")
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    validator_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, validator_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryValidatorUnbondingDelegationsResponse(_message.Message):
    __slots__ = ("unbonding_responses", "pagination")
    UNBONDING_RESPONSES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    unbonding_responses: _containers.RepeatedCompositeFieldContainer[_staking_pb2.UnbondingDelegation]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, unbonding_responses: _Optional[_Iterable[_Union[_staking_pb2.UnbondingDelegation, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryDelegationRequest(_message.Message):
    __slots__ = ("delegator_addr", "validator_addr")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    validator_addr: str
    def __init__(self, delegator_addr: _Optional[str] = ..., validator_addr: _Optional[str] = ...) -> None: ...

class QueryDelegationResponse(_message.Message):
    __slots__ = ("delegation_response",)
    DELEGATION_RESPONSE_FIELD_NUMBER: _ClassVar[int]
    delegation_response: _staking_pb2.DelegationResponse
    def __init__(self, delegation_response: _Optional[_Union[_staking_pb2.DelegationResponse, _Mapping]] = ...) -> None: ...

class QueryUnbondingDelegationRequest(_message.Message):
    __slots__ = ("delegator_addr", "validator_addr")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    validator_addr: str
    def __init__(self, delegator_addr: _Optional[str] = ..., validator_addr: _Optional[str] = ...) -> None: ...

class QueryUnbondingDelegationResponse(_message.Message):
    __slots__ = ("unbond",)
    UNBOND_FIELD_NUMBER: _ClassVar[int]
    unbond: _staking_pb2.UnbondingDelegation
    def __init__(self, unbond: _Optional[_Union[_staking_pb2.UnbondingDelegation, _Mapping]] = ...) -> None: ...

class QueryDelegatorDelegationsRequest(_message.Message):
    __slots__ = ("delegator_addr", "pagination")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, delegator_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryDelegatorDelegationsResponse(_message.Message):
    __slots__ = ("delegation_responses", "pagination")
    DELEGATION_RESPONSES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    delegation_responses: _containers.RepeatedCompositeFieldContainer[_staking_pb2.DelegationResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, delegation_responses: _Optional[_Iterable[_Union[_staking_pb2.DelegationResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryDelegatorUnbondingDelegationsRequest(_message.Message):
    __slots__ = ("delegator_addr", "pagination")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, delegator_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryDelegatorUnbondingDelegationsResponse(_message.Message):
    __slots__ = ("unbonding_responses", "pagination")
    UNBONDING_RESPONSES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    unbonding_responses: _containers.RepeatedCompositeFieldContainer[_staking_pb2.UnbondingDelegation]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, unbonding_responses: _Optional[_Iterable[_Union[_staking_pb2.UnbondingDelegation, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryRedelegationsRequest(_message.Message):
    __slots__ = ("delegator_addr", "src_validator_addr", "dst_validator_addr", "pagination")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    SRC_VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    DST_VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    src_validator_addr: str
    dst_validator_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, delegator_addr: _Optional[str] = ..., src_validator_addr: _Optional[str] = ..., dst_validator_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryRedelegationsResponse(_message.Message):
    __slots__ = ("redelegation_responses", "pagination")
    REDELEGATION_RESPONSES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    redelegation_responses: _containers.RepeatedCompositeFieldContainer[_staking_pb2.RedelegationResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, redelegation_responses: _Optional[_Iterable[_Union[_staking_pb2.RedelegationResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryDelegatorValidatorsRequest(_message.Message):
    __slots__ = ("delegator_addr", "pagination")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, delegator_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryDelegatorValidatorsResponse(_message.Message):
    __slots__ = ("validators", "pagination")
    VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    validators: _containers.RepeatedCompositeFieldContainer[_hybrid_validation_pb2.ValidatorHV]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, validators: _Optional[_Iterable[_Union[_hybrid_validation_pb2.ValidatorHV, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryDelegatorValidatorRequest(_message.Message):
    __slots__ = ("delegator_addr", "validator_addr")
    DELEGATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    delegator_addr: str
    validator_addr: str
    def __init__(self, delegator_addr: _Optional[str] = ..., validator_addr: _Optional[str] = ...) -> None: ...

class QueryDelegatorValidatorResponse(_message.Message):
    __slots__ = ("validator",)
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    validator: _hybrid_validation_pb2.ValidatorHV
    def __init__(self, validator: _Optional[_Union[_hybrid_validation_pb2.ValidatorHV, _Mapping]] = ...) -> None: ...

class QueryHistoricalInfoRequest(_message.Message):
    __slots__ = ("height",)
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    height: int
    def __init__(self, height: _Optional[int] = ...) -> None: ...

class QueryHistoricalInfoResponse(_message.Message):
    __slots__ = ("hist",)
    HIST_FIELD_NUMBER: _ClassVar[int]
    hist: _hybrid_validation_pb2.HistoricalInfoHV
    def __init__(self, hist: _Optional[_Union[_hybrid_validation_pb2.HistoricalInfoHV, _Mapping]] = ...) -> None: ...

class QueryPoolRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryPoolResponse(_message.Message):
    __slots__ = ("pool",)
    POOL_FIELD_NUMBER: _ClassVar[int]
    pool: _staking_pb2.Pool
    def __init__(self, pool: _Optional[_Union[_staking_pb2.Pool, _Mapping]] = ...) -> None: ...

class QueryParamsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryParamsResponse(_message.Message):
    __slots__ = ("Params", "HVParams")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    HVPARAMS_FIELD_NUMBER: _ClassVar[int]
    Params: _staking_pb2.Params
    HVParams: _hybrid_validation_pb2.HVParams
    def __init__(self, Params: _Optional[_Union[_staking_pb2.Params, _Mapping]] = ..., HVParams: _Optional[_Union[_hybrid_validation_pb2.HVParams, _Mapping]] = ...) -> None: ...

class QueryPowerRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ValidatorPower(_message.Message):
    __slots__ = ("operator_addr", "cons_addr", "power", "jailed", "status")
    OPERATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    CONS_ADDR_FIELD_NUMBER: _ClassVar[int]
    POWER_FIELD_NUMBER: _ClassVar[int]
    JAILED_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    operator_addr: str
    cons_addr: str
    power: int
    jailed: bool
    status: _staking_pb2.BondStatus
    def __init__(self, operator_addr: _Optional[str] = ..., cons_addr: _Optional[str] = ..., power: _Optional[int] = ..., jailed: bool = ..., status: _Optional[_Union[_staking_pb2.BondStatus, str]] = ...) -> None: ...

class QueryPowerResponse(_message.Message):
    __slots__ = ("validator_power", "total_power", "height")
    VALIDATOR_POWER_FIELD_NUMBER: _ClassVar[int]
    TOTAL_POWER_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    validator_power: _containers.RepeatedCompositeFieldContainer[ValidatorPower]
    total_power: int
    height: int
    def __init__(self, validator_power: _Optional[_Iterable[_Union[ValidatorPower, _Mapping]]] = ..., total_power: _Optional[int] = ..., height: _Optional[int] = ...) -> None: ...

class QueryPendingMintTransactionsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryPendingMintTransactionsResponse(_message.Message):
    __slots__ = ("pending_mint_transactions",)
    PENDING_MINT_TRANSACTIONS_FIELD_NUMBER: _ClassVar[int]
    pending_mint_transactions: _containers.RepeatedCompositeFieldContainer[PendingMintTransactionResponse]
    def __init__(self, pending_mint_transactions: _Optional[_Iterable[_Union[PendingMintTransactionResponse, _Mapping]]] = ...) -> None: ...

class PendingMintTransactionResponse(_message.Message):
    __slots__ = ("chain_id", "chain_type", "recipient_address", "amount", "creator", "key_id")
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    CHAIN_TYPE_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    chain_id: int
    chain_type: str
    recipient_address: str
    amount: int
    creator: str
    key_id: int
    def __init__(self, chain_id: _Optional[int] = ..., chain_type: _Optional[str] = ..., recipient_address: _Optional[str] = ..., amount: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ...) -> None: ...
