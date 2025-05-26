from cosmos.protocolpool.v1 import types_pb2 as _types_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgFundCommunityPool(_message.Message):
    __slots__ = ("depositor", "amount")
    DEPOSITOR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    depositor: str
    amount: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, depositor: _Optional[str] = ..., amount: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

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
    __slots__ = ("authority", "recipient")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    authority: str
    recipient: str
    def __init__(self, authority: _Optional[str] = ..., recipient: _Optional[str] = ...) -> None: ...

class MsgCancelContinuousFundResponse(_message.Message):
    __slots__ = ("canceled_time", "canceled_height", "recipient")
    CANCELED_TIME_FIELD_NUMBER: _ClassVar[int]
    CANCELED_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    canceled_time: _timestamp_pb2.Timestamp
    canceled_height: int
    recipient: str
    def __init__(self, canceled_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., canceled_height: _Optional[int] = ..., recipient: _Optional[str] = ...) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _types_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_types_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
