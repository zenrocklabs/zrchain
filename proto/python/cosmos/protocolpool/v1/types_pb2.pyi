from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ContinuousFund(_message.Message):
    __slots__ = ("recipient", "percentage", "expiry")
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    PERCENTAGE_FIELD_NUMBER: _ClassVar[int]
    EXPIRY_FIELD_NUMBER: _ClassVar[int]
    recipient: str
    percentage: str
    expiry: _timestamp_pb2.Timestamp
    def __init__(self, recipient: _Optional[str] = ..., percentage: _Optional[str] = ..., expiry: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class Params(_message.Message):
    __slots__ = ("enabled_distribution_denoms", "distribution_frequency")
    ENABLED_DISTRIBUTION_DENOMS_FIELD_NUMBER: _ClassVar[int]
    DISTRIBUTION_FREQUENCY_FIELD_NUMBER: _ClassVar[int]
    enabled_distribution_denoms: _containers.RepeatedScalarFieldContainer[str]
    distribution_frequency: int
    def __init__(self, enabled_distribution_denoms: _Optional[_Iterable[str]] = ..., distribution_frequency: _Optional[int] = ...) -> None: ...
