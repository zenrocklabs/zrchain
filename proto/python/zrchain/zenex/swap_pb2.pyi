from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.validation import asset_data_pb2 as _asset_data_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SwapStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SWAP_STATUS_UNSPECIFIED: _ClassVar[SwapStatus]
    SWAP_STATUS_REQUESTED: _ClassVar[SwapStatus]
    SWAP_STATUS_REJECTED: _ClassVar[SwapStatus]
    SWAP_STATUS_COMPLETED: _ClassVar[SwapStatus]
SWAP_STATUS_UNSPECIFIED: SwapStatus
SWAP_STATUS_REQUESTED: SwapStatus
SWAP_STATUS_REJECTED: SwapStatus
SWAP_STATUS_COMPLETED: SwapStatus

class Swap(_message.Message):
    __slots__ = ("creator", "swap_id", "status", "pair", "data", "sender_key_id", "recipient_key_id", "workspace", "zenbtc_yield")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAIR_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    SENDER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    ZENBTC_YIELD_FIELD_NUMBER: _ClassVar[int]
    creator: str
    swap_id: int
    status: SwapStatus
    pair: str
    data: SwapData
    sender_key_id: int
    recipient_key_id: int
    workspace: str
    zenbtc_yield: bool
    def __init__(self, creator: _Optional[str] = ..., swap_id: _Optional[int] = ..., status: _Optional[_Union[SwapStatus, str]] = ..., pair: _Optional[str] = ..., data: _Optional[_Union[SwapData, _Mapping]] = ..., sender_key_id: _Optional[int] = ..., recipient_key_id: _Optional[int] = ..., workspace: _Optional[str] = ..., zenbtc_yield: bool = ...) -> None: ...

class SwapData(_message.Message):
    __slots__ = ("base_token", "quote_token", "price", "amount_in", "amount_out")
    BASE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    QUOTE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_IN_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_OUT_FIELD_NUMBER: _ClassVar[int]
    base_token: _asset_data_pb2.AssetData
    quote_token: _asset_data_pb2.AssetData
    price: str
    amount_in: str
    amount_out: str
    def __init__(self, base_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ..., quote_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ..., price: _Optional[str] = ..., amount_in: _Optional[str] = ..., amount_out: _Optional[str] = ...) -> None: ...

class SwapPair(_message.Message):
    __slots__ = ("base_token", "quote_token")
    BASE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    QUOTE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    base_token: _asset_data_pb2.AssetData
    quote_token: _asset_data_pb2.AssetData
    def __init__(self, base_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ..., quote_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ...) -> None: ...
