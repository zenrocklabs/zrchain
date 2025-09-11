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
    SWAP_STATUS_SWAP_TRANSFER_REQUESTED: _ClassVar[SwapStatus]
    SWAP_STATUS_REJECTED: _ClassVar[SwapStatus]
    SWAP_STATUS_COMPLETED: _ClassVar[SwapStatus]

class TradePair(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    TRADE_PAIR_UNSPECIFIED: _ClassVar[TradePair]
    TRADE_PAIR_ROCK_BTC: _ClassVar[TradePair]
    TRADE_PAIR_BTC_ROCK: _ClassVar[TradePair]
SWAP_STATUS_UNSPECIFIED: SwapStatus
SWAP_STATUS_REQUESTED: SwapStatus
SWAP_STATUS_SWAP_TRANSFER_REQUESTED: SwapStatus
SWAP_STATUS_REJECTED: SwapStatus
SWAP_STATUS_COMPLETED: SwapStatus
TRADE_PAIR_UNSPECIFIED: TradePair
TRADE_PAIR_ROCK_BTC: TradePair
TRADE_PAIR_BTC_ROCK: TradePair

class Swap(_message.Message):
    __slots__ = ("creator", "swap_id", "status", "pair", "data", "rock_key_id", "btc_key_id", "zenex_pool_key_id", "workspace", "sign_tx_id", "source_tx_hash", "reject_reason")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAIR_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    ROCK_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    BTC_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    ZENEX_POOL_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    SIGN_TX_ID_FIELD_NUMBER: _ClassVar[int]
    SOURCE_TX_HASH_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    creator: str
    swap_id: int
    status: SwapStatus
    pair: TradePair
    data: SwapData
    rock_key_id: int
    btc_key_id: int
    zenex_pool_key_id: int
    workspace: str
    sign_tx_id: int
    source_tx_hash: str
    reject_reason: str
    def __init__(self, creator: _Optional[str] = ..., swap_id: _Optional[int] = ..., status: _Optional[_Union[SwapStatus, str]] = ..., pair: _Optional[_Union[TradePair, str]] = ..., data: _Optional[_Union[SwapData, _Mapping]] = ..., rock_key_id: _Optional[int] = ..., btc_key_id: _Optional[int] = ..., zenex_pool_key_id: _Optional[int] = ..., workspace: _Optional[str] = ..., sign_tx_id: _Optional[int] = ..., source_tx_hash: _Optional[str] = ..., reject_reason: _Optional[str] = ...) -> None: ...

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
    amount_in: int
    amount_out: int
    def __init__(self, base_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ..., quote_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ..., price: _Optional[str] = ..., amount_in: _Optional[int] = ..., amount_out: _Optional[int] = ...) -> None: ...

class SwapPair(_message.Message):
    __slots__ = ("base_token", "quote_token")
    BASE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    QUOTE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    base_token: _asset_data_pb2.AssetData
    quote_token: _asset_data_pb2.AssetData
    def __init__(self, base_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ..., quote_token: _Optional[_Union[_asset_data_pb2.AssetData, _Mapping]] = ...) -> None: ...
