from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from zrchain.zentp import params_pb2 as _params_pb2
from zrchain.zentp import bridge_pb2 as _bridge_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryParamsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryParamsResponse(_message.Message):
    __slots__ = ("params",)
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class QueryMintsRequest(_message.Message):
    __slots__ = ("id", "creator", "status", "denom", "pagination")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    status: _bridge_pb2.BridgeStatus
    denom: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., status: _Optional[_Union[_bridge_pb2.BridgeStatus, str]] = ..., denom: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryMintsResponse(_message.Message):
    __slots__ = ("mints", "pagination")
    MINTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    mints: _containers.RepeatedCompositeFieldContainer[_bridge_pb2.Bridge]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, mints: _Optional[_Iterable[_Union[_bridge_pb2.Bridge, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryBurnsRequest(_message.Message):
    __slots__ = ("id", "denom", "status", "pagination")
    ID_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    id: int
    denom: str
    status: _bridge_pb2.BridgeStatus
    pagination: _pagination_pb2.PageRequest
    def __init__(self, id: _Optional[int] = ..., denom: _Optional[str] = ..., status: _Optional[_Union[_bridge_pb2.BridgeStatus, str]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryBurnsResponse(_message.Message):
    __slots__ = ("burns", "pagination")
    BURNS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    burns: _containers.RepeatedCompositeFieldContainer[_bridge_pb2.Bridge]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, burns: _Optional[_Iterable[_Union[_bridge_pb2.Bridge, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryStatsRequest(_message.Message):
    __slots__ = ("address", "denom")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    address: str
    denom: str
    def __init__(self, address: _Optional[str] = ..., denom: _Optional[str] = ...) -> None: ...

class QueryStatsResponse(_message.Message):
    __slots__ = ("total_mints", "total_burns")
    TOTAL_MINTS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BURNS_FIELD_NUMBER: _ClassVar[int]
    total_mints: int
    total_burns: int
    def __init__(self, total_mints: _Optional[int] = ..., total_burns: _Optional[int] = ...) -> None: ...
