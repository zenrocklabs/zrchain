from cosmos.protocolpool.v1 import types_pb2 as _types_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("continuous_funds", "params")
    CONTINUOUS_FUNDS_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    continuous_funds: _containers.RepeatedCompositeFieldContainer[_types_pb2.ContinuousFund]
    params: _types_pb2.Params
    def __init__(self, continuous_funds: _Optional[_Iterable[_Union[_types_pb2.ContinuousFund, _Mapping]]] = ..., params: _Optional[_Union[_types_pb2.Params, _Mapping]] = ...) -> None: ...
