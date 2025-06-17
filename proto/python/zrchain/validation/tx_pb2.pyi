from amino import amino_pb2 as _amino_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from tendermint.abci import types_pb2 as _types_pb2
from zrchain.validation import hybrid_validation_pb2 as _hybrid_validation_pb2
from zrchain.validation import staking_pb2 as _staking_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class EventType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    EVENT_TYPE_UNSPECIFIED: _ClassVar[EventType]
    EVENT_TYPE_ZENBTC_MINT: _ClassVar[EventType]
    EVENT_TYPE_ZENBTC_BURN: _ClassVar[EventType]
    EVENT_TYPE_ZENTP_MINT: _ClassVar[EventType]
    EVENT_TYPE_ZENTP_BURN: _ClassVar[EventType]
EVENT_TYPE_UNSPECIFIED: EventType
EVENT_TYPE_ZENBTC_MINT: EventType
EVENT_TYPE_ZENBTC_BURN: EventType
EVENT_TYPE_ZENTP_MINT: EventType
EVENT_TYPE_ZENTP_BURN: EventType

class MsgCreateValidator(_message.Message):
    __slots__ = ("description", "commission", "min_self_delegation", "delegator_address", "validator_address", "pubkey", "value")
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    COMMISSION_FIELD_NUMBER: _ClassVar[int]
    MIN_SELF_DELEGATION_FIELD_NUMBER: _ClassVar[int]
    DELEGATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    PUBKEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    description: _staking_pb2.Description
    commission: _staking_pb2.CommissionRates
    min_self_delegation: str
    delegator_address: str
    validator_address: str
    pubkey: _any_pb2.Any
    value: _coin_pb2.Coin
    def __init__(self, description: _Optional[_Union[_staking_pb2.Description, _Mapping]] = ..., commission: _Optional[_Union[_staking_pb2.CommissionRates, _Mapping]] = ..., min_self_delegation: _Optional[str] = ..., delegator_address: _Optional[str] = ..., validator_address: _Optional[str] = ..., pubkey: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., value: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgCreateValidatorResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgEditValidator(_message.Message):
    __slots__ = ("description", "validator_address", "commission_rate", "min_self_delegation")
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    COMMISSION_RATE_FIELD_NUMBER: _ClassVar[int]
    MIN_SELF_DELEGATION_FIELD_NUMBER: _ClassVar[int]
    description: _staking_pb2.Description
    validator_address: str
    commission_rate: str
    min_self_delegation: str
    def __init__(self, description: _Optional[_Union[_staking_pb2.Description, _Mapping]] = ..., validator_address: _Optional[str] = ..., commission_rate: _Optional[str] = ..., min_self_delegation: _Optional[str] = ...) -> None: ...

class MsgEditValidatorResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgDelegate(_message.Message):
    __slots__ = ("delegator_address", "validator_address", "amount")
    DELEGATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    delegator_address: str
    validator_address: str
    amount: _coin_pb2.Coin
    def __init__(self, delegator_address: _Optional[str] = ..., validator_address: _Optional[str] = ..., amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgDelegateResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgBeginRedelegate(_message.Message):
    __slots__ = ("delegator_address", "validator_src_address", "validator_dst_address", "amount")
    DELEGATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_SRC_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_DST_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    delegator_address: str
    validator_src_address: str
    validator_dst_address: str
    amount: _coin_pb2.Coin
    def __init__(self, delegator_address: _Optional[str] = ..., validator_src_address: _Optional[str] = ..., validator_dst_address: _Optional[str] = ..., amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgBeginRedelegateResponse(_message.Message):
    __slots__ = ("completion_time",)
    COMPLETION_TIME_FIELD_NUMBER: _ClassVar[int]
    completion_time: _timestamp_pb2.Timestamp
    def __init__(self, completion_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class MsgUndelegate(_message.Message):
    __slots__ = ("delegator_address", "validator_address", "amount")
    DELEGATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    delegator_address: str
    validator_address: str
    amount: _coin_pb2.Coin
    def __init__(self, delegator_address: _Optional[str] = ..., validator_address: _Optional[str] = ..., amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgUndelegateResponse(_message.Message):
    __slots__ = ("completion_time", "amount")
    COMPLETION_TIME_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    completion_time: _timestamp_pb2.Timestamp
    amount: _coin_pb2.Coin
    def __init__(self, completion_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgCancelUnbondingDelegation(_message.Message):
    __slots__ = ("delegator_address", "validator_address", "amount", "creation_height")
    DELEGATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    CREATION_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    delegator_address: str
    validator_address: str
    amount: _coin_pb2.Coin
    creation_height: int
    def __init__(self, delegator_address: _Optional[str] = ..., validator_address: _Optional[str] = ..., amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., creation_height: _Optional[int] = ...) -> None: ...

class MsgCancelUnbondingDelegationResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "Params", "HVParams")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    HVPARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    Params: _staking_pb2.Params
    HVParams: _hybrid_validation_pb2.HVParams
    def __init__(self, authority: _Optional[str] = ..., Params: _Optional[_Union[_staking_pb2.Params, _Mapping]] = ..., HVParams: _Optional[_Union[_hybrid_validation_pb2.HVParams, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateHVParams(_message.Message):
    __slots__ = ("authority", "HVParams")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    HVPARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    HVParams: _hybrid_validation_pb2.HVParams
    def __init__(self, authority: _Optional[str] = ..., HVParams: _Optional[_Union[_hybrid_validation_pb2.HVParams, _Mapping]] = ...) -> None: ...

class MsgUpdateHVParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgTriggerEventBackfill(_message.Message):
    __slots__ = ("authority", "tx_hash", "caip2_chain_id", "event_type")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    TX_HASH_FIELD_NUMBER: _ClassVar[int]
    CAIP2_CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    EVENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    authority: str
    tx_hash: str
    caip2_chain_id: str
    event_type: EventType
    def __init__(self, authority: _Optional[str] = ..., tx_hash: _Optional[str] = ..., caip2_chain_id: _Optional[str] = ..., event_type: _Optional[_Union[EventType, str]] = ...) -> None: ...

class MsgTriggerEventBackfillResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class BackfillRequests(_message.Message):
    __slots__ = ("requests",)
    REQUESTS_FIELD_NUMBER: _ClassVar[int]
    requests: _containers.RepeatedCompositeFieldContainer[MsgTriggerEventBackfill]
    def __init__(self, requests: _Optional[_Iterable[_Union[MsgTriggerEventBackfill, _Mapping]]] = ...) -> None: ...
