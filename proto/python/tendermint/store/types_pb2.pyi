from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class BlockStoreState(_message.Message):
    __slots__ = ("base", "height")
    BASE_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    base: int
    height: int
    def __init__(self, base: _Optional[int] = ..., height: _Optional[int] = ...) -> None: ...
