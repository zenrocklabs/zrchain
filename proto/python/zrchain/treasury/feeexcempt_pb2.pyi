from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class FeeExcemptMsgs(_message.Message):
    __slots__ = ("msgs",)
    MSGS_FIELD_NUMBER: _ClassVar[int]
    msgs: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, msgs: _Optional[_Iterable[str]] = ...) -> None: ...
