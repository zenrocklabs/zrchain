from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Module(_message.Message):
    __slots__ = ("max_metadata_len", "authority", "max_title_len", "max_summary_len")
    MAX_METADATA_LEN_FIELD_NUMBER: _ClassVar[int]
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MAX_TITLE_LEN_FIELD_NUMBER: _ClassVar[int]
    MAX_SUMMARY_LEN_FIELD_NUMBER: _ClassVar[int]
    max_metadata_len: int
    authority: str
    max_title_len: int
    max_summary_len: int
    def __init__(self, max_metadata_len: _Optional[int] = ..., authority: _Optional[str] = ..., max_title_len: _Optional[int] = ..., max_summary_len: _Optional[int] = ...) -> None: ...
