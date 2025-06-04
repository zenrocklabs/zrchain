from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.epochs.v1beta1 import genesis_pb2 as _genesis_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryEpochInfosRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryEpochInfosResponse(_message.Message):
    __slots__ = ("epochs",)
    EPOCHS_FIELD_NUMBER: _ClassVar[int]
    epochs: _containers.RepeatedCompositeFieldContainer[_genesis_pb2.EpochInfo]
    def __init__(self, epochs: _Optional[_Iterable[_Union[_genesis_pb2.EpochInfo, _Mapping]]] = ...) -> None: ...

class QueryCurrentEpochRequest(_message.Message):
    __slots__ = ("identifier",)
    IDENTIFIER_FIELD_NUMBER: _ClassVar[int]
    identifier: str
    def __init__(self, identifier: _Optional[str] = ...) -> None: ...

class QueryCurrentEpochResponse(_message.Message):
    __slots__ = ("current_epoch",)
    CURRENT_EPOCH_FIELD_NUMBER: _ClassVar[int]
    current_epoch: int
    def __init__(self, current_epoch: _Optional[int] = ...) -> None: ...
