from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class EventDataRoundState(_message.Message):
    __slots__ = ("height", "round", "step")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    STEP_FIELD_NUMBER: _ClassVar[int]
    height: int
    round: int
    step: str
    def __init__(self, height: _Optional[int] = ..., round: _Optional[int] = ..., step: _Optional[str] = ...) -> None: ...
