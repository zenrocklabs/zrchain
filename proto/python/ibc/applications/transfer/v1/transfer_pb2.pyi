from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("send_enabled", "receive_enabled")
    SEND_ENABLED_FIELD_NUMBER: _ClassVar[int]
    RECEIVE_ENABLED_FIELD_NUMBER: _ClassVar[int]
    send_enabled: bool
    receive_enabled: bool
    def __init__(self, send_enabled: bool = ..., receive_enabled: bool = ...) -> None: ...
