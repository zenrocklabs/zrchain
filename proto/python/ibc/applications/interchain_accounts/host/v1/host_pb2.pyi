from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("host_enabled", "allow_messages")
    HOST_ENABLED_FIELD_NUMBER: _ClassVar[int]
    ALLOW_MESSAGES_FIELD_NUMBER: _ClassVar[int]
    host_enabled: bool
    allow_messages: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, host_enabled: bool = ..., allow_messages: _Optional[_Iterable[str]] = ...) -> None: ...

class QueryRequest(_message.Message):
    __slots__ = ("path", "data")
    PATH_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    path: str
    data: bytes
    def __init__(self, path: _Optional[str] = ..., data: _Optional[bytes] = ...) -> None: ...
