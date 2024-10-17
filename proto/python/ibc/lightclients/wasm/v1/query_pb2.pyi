from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryChecksumsRequest(_message.Message):
    __slots__ = ("pagination",)
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryChecksumsResponse(_message.Message):
    __slots__ = ("checksums", "pagination")
    CHECKSUMS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    checksums: _containers.RepeatedScalarFieldContainer[str]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, checksums: _Optional[_Iterable[str]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryCodeRequest(_message.Message):
    __slots__ = ("checksum",)
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    checksum: str
    def __init__(self, checksum: _Optional[str] = ...) -> None: ...

class QueryCodeResponse(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...
