from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryCommunityPoolRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryCommunityPoolResponse(_message.Message):
    __slots__ = ("pool",)
    POOL_FIELD_NUMBER: _ClassVar[int]
    pool: _containers.RepeatedCompositeFieldContainer[_coin_pb2.DecCoin]
    def __init__(self, pool: _Optional[_Iterable[_Union[_coin_pb2.DecCoin, _Mapping]]] = ...) -> None: ...

class QueryUnclaimedBudgetRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class QueryUnclaimedBudgetResponse(_message.Message):
    __slots__ = ("total_budget", "claimed_amount", "unclaimed_amount", "next_claim_from", "period", "tranches_left")
    TOTAL_BUDGET_FIELD_NUMBER: _ClassVar[int]
    CLAIMED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    UNCLAIMED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    NEXT_CLAIM_FROM_FIELD_NUMBER: _ClassVar[int]
    PERIOD_FIELD_NUMBER: _ClassVar[int]
    TRANCHES_LEFT_FIELD_NUMBER: _ClassVar[int]
    total_budget: _coin_pb2.Coin
    claimed_amount: _coin_pb2.Coin
    unclaimed_amount: _coin_pb2.Coin
    next_claim_from: _timestamp_pb2.Timestamp
    period: _duration_pb2.Duration
    tranches_left: int
    def __init__(self, total_budget: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., claimed_amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., unclaimed_amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., next_claim_from: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., tranches_left: _Optional[int] = ...) -> None: ...
