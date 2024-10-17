from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgFundCommunityPool(_message.Message):
    __slots__ = ("amount", "depositor")
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    DEPOSITOR_FIELD_NUMBER: _ClassVar[int]
    amount: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    depositor: str
    def __init__(self, amount: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., depositor: _Optional[str] = ...) -> None: ...

class MsgFundCommunityPoolResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgCommunityPoolSpend(_message.Message):
    __slots__ = ("authority", "recipient", "amount")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    authority: str
    recipient: str
    amount: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, authority: _Optional[str] = ..., recipient: _Optional[str] = ..., amount: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class MsgCommunityPoolSpendResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSubmitBudgetProposal(_message.Message):
    __slots__ = ("authority", "recipient_address", "total_budget", "start_time", "tranches", "period")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BUDGET_FIELD_NUMBER: _ClassVar[int]
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    TRANCHES_FIELD_NUMBER: _ClassVar[int]
    PERIOD_FIELD_NUMBER: _ClassVar[int]
    authority: str
    recipient_address: str
    total_budget: _coin_pb2.Coin
    start_time: _timestamp_pb2.Timestamp
    tranches: int
    period: _duration_pb2.Duration
    def __init__(self, authority: _Optional[str] = ..., recipient_address: _Optional[str] = ..., total_budget: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., start_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., tranches: _Optional[int] = ..., period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class MsgSubmitBudgetProposalResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgClaimBudget(_message.Message):
    __slots__ = ("recipient_address",)
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    recipient_address: str
    def __init__(self, recipient_address: _Optional[str] = ...) -> None: ...

class MsgClaimBudgetResponse(_message.Message):
    __slots__ = ("amount",)
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    amount: _coin_pb2.Coin
    def __init__(self, amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgCreateContinuousFund(_message.Message):
    __slots__ = ("authority", "recipient", "percentage", "expiry")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    PERCENTAGE_FIELD_NUMBER: _ClassVar[int]
    EXPIRY_FIELD_NUMBER: _ClassVar[int]
    authority: str
    recipient: str
    percentage: str
    expiry: _timestamp_pb2.Timestamp
    def __init__(self, authority: _Optional[str] = ..., recipient: _Optional[str] = ..., percentage: _Optional[str] = ..., expiry: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class MsgCreateContinuousFundResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgCancelContinuousFund(_message.Message):
    __slots__ = ("authority", "recipient_address")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    recipient_address: str
    def __init__(self, authority: _Optional[str] = ..., recipient_address: _Optional[str] = ...) -> None: ...

class MsgCancelContinuousFundResponse(_message.Message):
    __slots__ = ("canceled_time", "canceled_height", "recipient_address", "withdrawn_allocated_fund")
    CANCELED_TIME_FIELD_NUMBER: _ClassVar[int]
    CANCELED_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    WITHDRAWN_ALLOCATED_FUND_FIELD_NUMBER: _ClassVar[int]
    canceled_time: _timestamp_pb2.Timestamp
    canceled_height: int
    recipient_address: str
    withdrawn_allocated_fund: _coin_pb2.Coin
    def __init__(self, canceled_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., canceled_height: _Optional[int] = ..., recipient_address: _Optional[str] = ..., withdrawn_allocated_fund: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgWithdrawContinuousFund(_message.Message):
    __slots__ = ("recipient_address",)
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    recipient_address: str
    def __init__(self, recipient_address: _Optional[str] = ...) -> None: ...

class MsgWithdrawContinuousFundResponse(_message.Message):
    __slots__ = ("amount",)
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    amount: _coin_pb2.Coin
    def __init__(self, amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...
