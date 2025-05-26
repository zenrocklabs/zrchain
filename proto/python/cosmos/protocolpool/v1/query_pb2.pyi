from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.protocolpool.v1 import types_pb2 as _types_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryCommunityPoolRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryCommunityPoolResponse(_message.Message):
    __slots__ = ("pool",)
    POOL_FIELD_NUMBER: _ClassVar[int]
    pool: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, pool: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class QueryContinuousFundRequest(_message.Message):
    __slots__ = ("recipient",)
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    recipient: str
    def __init__(self, recipient: _Optional[str] = ...) -> None: ...

class QueryContinuousFundResponse(_message.Message):
    __slots__ = ("continuous_fund",)
    CONTINUOUS_FUND_FIELD_NUMBER: _ClassVar[int]
    continuous_fund: _types_pb2.ContinuousFund
    def __init__(self, continuous_fund: _Optional[_Union[_types_pb2.ContinuousFund, _Mapping]] = ...) -> None: ...

class QueryContinuousFundsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryContinuousFundsResponse(_message.Message):
    __slots__ = ("continuous_funds",)
    CONTINUOUS_FUNDS_FIELD_NUMBER: _ClassVar[int]
    continuous_funds: _containers.RepeatedCompositeFieldContainer[_types_pb2.ContinuousFund]
    def __init__(self, continuous_funds: _Optional[_Iterable[_Union[_types_pb2.ContinuousFund, _Mapping]]] = ...) -> None: ...

class QueryParamsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryParamsResponse(_message.Message):
    __slots__ = ("params",)
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    params: _types_pb2.Params
    def __init__(self, params: _Optional[_Union[_types_pb2.Params, _Mapping]] = ...) -> None: ...
