from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MsgIncreaseCounter(_message.Message):
    __slots__ = ("signer", "count")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    signer: str
    count: int
    def __init__(self, signer: _Optional[str] = ..., count: _Optional[int] = ...) -> None: ...

class MsgIncreaseCountResponse(_message.Message):
    __slots__ = ("new_count",)
    NEW_COUNT_FIELD_NUMBER: _ClassVar[int]
    new_count: int
    def __init__(self, new_count: _Optional[int] = ...) -> None: ...
