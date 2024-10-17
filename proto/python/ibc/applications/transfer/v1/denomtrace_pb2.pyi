from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class DenomTrace(_message.Message):
    __slots__ = ("path", "base_denom")
    PATH_FIELD_NUMBER: _ClassVar[int]
    BASE_DENOM_FIELD_NUMBER: _ClassVar[int]
    path: str
    base_denom: str
    def __init__(self, path: _Optional[str] = ..., base_denom: _Optional[str] = ...) -> None: ...
