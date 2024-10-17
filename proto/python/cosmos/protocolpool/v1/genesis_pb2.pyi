from cosmos.protocolpool.v1 import types_pb2 as _types_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("continuous_fund", "budget", "to_distribute")
    CONTINUOUS_FUND_FIELD_NUMBER: _ClassVar[int]
    BUDGET_FIELD_NUMBER: _ClassVar[int]
    TO_DISTRIBUTE_FIELD_NUMBER: _ClassVar[int]
    continuous_fund: _containers.RepeatedCompositeFieldContainer[_types_pb2.ContinuousFund]
    budget: _containers.RepeatedCompositeFieldContainer[_types_pb2.Budget]
    to_distribute: str
    def __init__(self, continuous_fund: _Optional[_Iterable[_Union[_types_pb2.ContinuousFund, _Mapping]]] = ..., budget: _Optional[_Iterable[_Union[_types_pb2.Budget, _Mapping]]] = ..., to_distribute: _Optional[str] = ...) -> None: ...
