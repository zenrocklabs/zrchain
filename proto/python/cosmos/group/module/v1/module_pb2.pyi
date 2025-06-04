from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Module(_message.Message):
    __slots__ = ("max_execution_period", "max_metadata_len")
    MAX_EXECUTION_PERIOD_FIELD_NUMBER: _ClassVar[int]
    MAX_METADATA_LEN_FIELD_NUMBER: _ClassVar[int]
    max_execution_period: _duration_pb2.Duration
    max_metadata_len: int
    def __init__(self, max_execution_period: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., max_metadata_len: _Optional[int] = ...) -> None: ...
