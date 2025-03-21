from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Config(_message.Message):
    __slots__ = ("allowed_relayers",)
    ALLOWED_RELAYERS_FIELD_NUMBER: _ClassVar[int]
    allowed_relayers: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, allowed_relayers: _Optional[_Iterable[str]] = ...) -> None: ...
