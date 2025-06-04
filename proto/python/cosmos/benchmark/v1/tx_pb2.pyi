from cosmos.benchmark.v1 import benchmark_pb2 as _benchmark_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgLoadTest(_message.Message):
    __slots__ = ("caller", "ops")
    CALLER_FIELD_NUMBER: _ClassVar[int]
    OPS_FIELD_NUMBER: _ClassVar[int]
    caller: bytes
    ops: _containers.RepeatedCompositeFieldContainer[_benchmark_pb2.Op]
    def __init__(self, caller: _Optional[bytes] = ..., ops: _Optional[_Iterable[_Union[_benchmark_pb2.Op, _Mapping]]] = ...) -> None: ...

class MsgLoadTestResponse(_message.Message):
    __slots__ = ("total_time", "total_errors")
    TOTAL_TIME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_ERRORS_FIELD_NUMBER: _ClassVar[int]
    total_time: int
    total_errors: int
    def __init__(self, total_time: _Optional[int] = ..., total_errors: _Optional[int] = ...) -> None: ...
