from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class EpochInfo(_message.Message):
    __slots__ = ("identifier", "start_time", "duration", "current_epoch", "current_epoch_start_time", "epoch_counting_started", "current_epoch_start_height")
    IDENTIFIER_FIELD_NUMBER: _ClassVar[int]
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    DURATION_FIELD_NUMBER: _ClassVar[int]
    CURRENT_EPOCH_FIELD_NUMBER: _ClassVar[int]
    CURRENT_EPOCH_START_TIME_FIELD_NUMBER: _ClassVar[int]
    EPOCH_COUNTING_STARTED_FIELD_NUMBER: _ClassVar[int]
    CURRENT_EPOCH_START_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    identifier: str
    start_time: _timestamp_pb2.Timestamp
    duration: _duration_pb2.Duration
    current_epoch: int
    current_epoch_start_time: _timestamp_pb2.Timestamp
    epoch_counting_started: bool
    current_epoch_start_height: int
    def __init__(self, identifier: _Optional[str] = ..., start_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., duration: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., current_epoch: _Optional[int] = ..., current_epoch_start_time: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., epoch_counting_started: bool = ..., current_epoch_start_height: _Optional[int] = ...) -> None: ...

class GenesisState(_message.Message):
    __slots__ = ("epochs",)
    EPOCHS_FIELD_NUMBER: _ClassVar[int]
    epochs: _containers.RepeatedCompositeFieldContainer[EpochInfo]
    def __init__(self, epochs: _Optional[_Iterable[_Union[EpochInfo, _Mapping]]] = ...) -> None: ...
