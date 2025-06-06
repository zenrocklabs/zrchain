from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Module(_message.Message):
    __slots__ = ("max_metadata_len", "authority")
    MAX_METADATA_LEN_FIELD_NUMBER: _ClassVar[int]
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    max_metadata_len: int
    authority: str
    def __init__(self, max_metadata_len: _Optional[int] = ..., authority: _Optional[str] = ...) -> None: ...
