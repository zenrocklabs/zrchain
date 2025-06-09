from google.api import annotations_pb2 as _annotations_pb2
from cometbft.types.v2 import params_pb2 as _params_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryParamsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryParamsResponse(_message.Message):
    __slots__ = ("params",)
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.ConsensusParams
    def __init__(self, params: _Optional[_Union[_params_pb2.ConsensusParams, _Mapping]] = ...) -> None: ...
