from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class QueryGetCountRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryGetCountResponse(_message.Message):
    __slots__ = ("total_count",)
    TOTAL_COUNT_FIELD_NUMBER: _ClassVar[int]
    total_count: int
    def __init__(self, total_count: _Optional[int] = ...) -> None: ...
