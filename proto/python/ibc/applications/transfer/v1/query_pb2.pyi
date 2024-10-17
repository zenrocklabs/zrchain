from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from ibc.applications.transfer.v1 import transfer_pb2 as _transfer_pb2
from google.api import annotations_pb2 as _annotations_pb2
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
    params: _transfer_pb2.Params
    def __init__(self, params: _Optional[_Union[_transfer_pb2.Params, _Mapping]] = ...) -> None: ...

class QueryDenomHashRequest(_message.Message):
    __slots__ = ("trace",)
    TRACE_FIELD_NUMBER: _ClassVar[int]
    trace: str
    def __init__(self, trace: _Optional[str] = ...) -> None: ...

class QueryDenomHashResponse(_message.Message):
    __slots__ = ("hash",)
    HASH_FIELD_NUMBER: _ClassVar[int]
    hash: str
    def __init__(self, hash: _Optional[str] = ...) -> None: ...

class QueryEscrowAddressRequest(_message.Message):
    __slots__ = ("port_id", "channel_id")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ...) -> None: ...

class QueryEscrowAddressResponse(_message.Message):
    __slots__ = ("escrow_address",)
    ESCROW_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    escrow_address: str
    def __init__(self, escrow_address: _Optional[str] = ...) -> None: ...

class QueryTotalEscrowForDenomRequest(_message.Message):
    __slots__ = ("denom",)
    DENOM_FIELD_NUMBER: _ClassVar[int]
    denom: str
    def __init__(self, denom: _Optional[str] = ...) -> None: ...

class QueryTotalEscrowForDenomResponse(_message.Message):
    __slots__ = ("amount",)
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    amount: _coin_pb2.Coin
    def __init__(self, amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...
