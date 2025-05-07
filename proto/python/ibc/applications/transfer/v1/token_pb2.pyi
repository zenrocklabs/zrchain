from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Token(_message.Message):
    __slots__ = ("denom", "amount")
    DENOM_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    denom: Denom
    amount: str
    def __init__(self, denom: _Optional[_Union[Denom, _Mapping]] = ..., amount: _Optional[str] = ...) -> None: ...

class Denom(_message.Message):
    __slots__ = ("base", "trace")
    BASE_FIELD_NUMBER: _ClassVar[int]
    TRACE_FIELD_NUMBER: _ClassVar[int]
    base: str
    trace: _containers.RepeatedCompositeFieldContainer[Hop]
    def __init__(self, base: _Optional[str] = ..., trace: _Optional[_Iterable[_Union[Hop, _Mapping]]] = ...) -> None: ...

class Hop(_message.Message):
    __slots__ = ("port_id", "channel_id")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ...) -> None: ...
