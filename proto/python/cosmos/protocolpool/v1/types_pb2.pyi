from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Budget(_message.Message):
    __slots__ = ("recipient_address", "total_budget", "claimed_amount", "start_time", "next_claim_from", "tranches", "tranches_left", "period")
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BUDGET_FIELD_NUMBER: _ClassVar[int]
    CLAIMED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    NEXT_CLAIM_FROM_FIELD_NUMBER: _ClassVar[int]
    TRANCHES_FIELD_NUMBER: _ClassVar[int]
    TRANCHES_LEFT_FIELD_NUMBER: _ClassVar[int]
    PERIOD_FIELD_NUMBER: _ClassVar[int]
    recipient_address: str
    total_budget: _coin_pb2.Coin
    claimed_amount: _coin_pb2.Coin
    start_time: _timestamp_pb2.Timestamp
    next_claim_from: _timestamp_pb2.Timestamp
    tranches: int
    tranches_left: int
    period: _duration_pb2.Duration
    def __init__(self, recipient_address: _Optional[str] = ..., total_budget: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., claimed_amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., start_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., next_claim_from: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., tranches: _Optional[int] = ..., tranches_left: _Optional[int] = ..., period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class ContinuousFund(_message.Message):
    __slots__ = ("recipient", "percentage", "expiry")
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    PERCENTAGE_FIELD_NUMBER: _ClassVar[int]
    EXPIRY_FIELD_NUMBER: _ClassVar[int]
    recipient: str
    percentage: str
    expiry: _timestamp_pb2.Timestamp
    def __init__(self, recipient: _Optional[str] = ..., percentage: _Optional[str] = ..., expiry: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...
