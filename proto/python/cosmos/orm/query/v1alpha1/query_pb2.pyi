from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import any_pb2 as _any_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetRequest(_message.Message):
    __slots__ = ("message_name", "index", "values")
    MESSAGE_NAME_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    VALUES_FIELD_NUMBER: _ClassVar[int]
    message_name: str
    index: str
    values: _containers.RepeatedCompositeFieldContainer[IndexValue]
    def __init__(self, message_name: _Optional[str] = ..., index: _Optional[str] = ..., values: _Optional[_Iterable[_Union[IndexValue, _Mapping]]] = ...) -> None: ...

class GetResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: _any_pb2.Any
    def __init__(self, result: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class ListRequest(_message.Message):
    __slots__ = ("message_name", "index", "prefix", "range", "pagination")
    class Prefix(_message.Message):
        __slots__ = ("values",)
        VALUES_FIELD_NUMBER: _ClassVar[int]
        values: _containers.RepeatedCompositeFieldContainer[IndexValue]
        def __init__(self, values: _Optional[_Iterable[_Union[IndexValue, _Mapping]]] = ...) -> None: ...
    class Range(_message.Message):
        __slots__ = ("start", "end")
        START_FIELD_NUMBER: _ClassVar[int]
        END_FIELD_NUMBER: _ClassVar[int]
        start: _containers.RepeatedCompositeFieldContainer[IndexValue]
        end: _containers.RepeatedCompositeFieldContainer[IndexValue]
        def __init__(self, start: _Optional[_Iterable[_Union[IndexValue, _Mapping]]] = ..., end: _Optional[_Iterable[_Union[IndexValue, _Mapping]]] = ...) -> None: ...
    MESSAGE_NAME_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    PREFIX_FIELD_NUMBER: _ClassVar[int]
    RANGE_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    message_name: str
    index: str
    prefix: ListRequest.Prefix
    range: ListRequest.Range
    pagination: _pagination_pb2.PageRequest
    def __init__(self, message_name: _Optional[str] = ..., index: _Optional[str] = ..., prefix: _Optional[_Union[ListRequest.Prefix, _Mapping]] = ..., range: _Optional[_Union[ListRequest.Range, _Mapping]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class ListResponse(_message.Message):
    __slots__ = ("results", "pagination")
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[_any_pb2.Any]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, results: _Optional[_Iterable[_Union[_any_pb2.Any, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class IndexValue(_message.Message):
    __slots__ = ("uint", "int", "str", "bytes", "enum", "bool", "timestamp", "duration")
    UINT_FIELD_NUMBER: _ClassVar[int]
    INT_FIELD_NUMBER: _ClassVar[int]
    STR_FIELD_NUMBER: _ClassVar[int]
    BYTES_FIELD_NUMBER: _ClassVar[int]
    ENUM_FIELD_NUMBER: _ClassVar[int]
    BOOL_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    DURATION_FIELD_NUMBER: _ClassVar[int]
    uint: int
    int: int
    str: str
    bytes: bytes
    enum: str
    bool: bool
    timestamp: _timestamp_pb2.Timestamp
    duration: _duration_pb2.Duration
    def __init__(self, uint: _Optional[int] = ..., int: _Optional[int] = ..., str: _Optional[str] = ..., bytes: _Optional[bytes] = ..., enum: _Optional[str] = ..., bool: bool = ..., timestamp: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., duration: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...
