from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class EventEpochEnd(_message.Message):
    __slots__ = ("epoch_number",)
    EPOCH_NUMBER_FIELD_NUMBER: _ClassVar[int]
    epoch_number: int
    def __init__(self, epoch_number: _Optional[int] = ...) -> None: ...

class EventEpochStart(_message.Message):
    __slots__ = ("epoch_number", "epoch_start_time")
    EPOCH_NUMBER_FIELD_NUMBER: _ClassVar[int]
    EPOCH_START_TIME_FIELD_NUMBER: _ClassVar[int]
    epoch_number: int
    epoch_start_time: int
    def __init__(self, epoch_number: _Optional[int] = ..., epoch_start_time: _Optional[int] = ...) -> None: ...
