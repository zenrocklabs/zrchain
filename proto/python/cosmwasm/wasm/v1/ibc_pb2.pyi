from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MsgIBCSend(_message.Message):
    __slots__ = ("channel", "timeout_height", "timeout_timestamp", "data")
    CHANNEL_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    channel: str
    timeout_height: int
    timeout_timestamp: int
    data: bytes
    def __init__(self, channel: _Optional[str] = ..., timeout_height: _Optional[int] = ..., timeout_timestamp: _Optional[int] = ..., data: _Optional[bytes] = ...) -> None: ...

class MsgIBCSendResponse(_message.Message):
    __slots__ = ("sequence",)
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    sequence: int
    def __init__(self, sequence: _Optional[int] = ...) -> None: ...

class MsgIBCWriteAcknowledgementResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgIBCCloseChannel(_message.Message):
    __slots__ = ("channel",)
    CHANNEL_FIELD_NUMBER: _ClassVar[int]
    channel: str
    def __init__(self, channel: _Optional[str] = ...) -> None: ...
