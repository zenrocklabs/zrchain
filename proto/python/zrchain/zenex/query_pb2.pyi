from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from zrchain.zenex import params_pb2 as _params_pb2
from zrchain.zenex import swap_pb2 as _swap_pb2
from zrchain.validation import asset_data_pb2 as _asset_data_pb2
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

class QuerySwapsRequest(_message.Message):
    __slots__ = ("creator", "swap_id", "status", "pair", "workspace", "source_tx_hash", "pagination")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAIR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_TX_HASH_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    creator: str
    swap_id: int
    status: _swap_pb2.SwapStatus
    pair: _swap_pb2.TradePair
    workspace: str
    source_tx_hash: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, creator: _Optional[str] = ..., swap_id: _Optional[int] = ..., status: _Optional[_Union[_swap_pb2.SwapStatus, str]] = ..., pair: _Optional[_Union[_swap_pb2.TradePair, str]] = ..., workspace: _Optional[str] = ..., source_tx_hash: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QuerySwapsResponse(_message.Message):
    __slots__ = ("swaps", "pagination")
    SWAPS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    swaps: _containers.RepeatedCompositeFieldContainer[_swap_pb2.Swap]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, swaps: _Optional[_Iterable[_Union[_swap_pb2.Swap, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryRockPoolRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryRockPoolResponse(_message.Message):
    __slots__ = ("rock_balance", "redeemable_assets")
    ROCK_BALANCE_FIELD_NUMBER: _ClassVar[int]
    REDEEMABLE_ASSETS_FIELD_NUMBER: _ClassVar[int]
    rock_balance: int
    redeemable_assets: _containers.RepeatedCompositeFieldContainer[RedeemableAsset]
    def __init__(self, rock_balance: _Optional[int] = ..., redeemable_assets: _Optional[_Iterable[_Union[RedeemableAsset, _Mapping]]] = ...) -> None: ...

class RedeemableAsset(_message.Message):
    __slots__ = ("asset", "amount")
    ASSET_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    asset: _asset_data_pb2.Asset
    amount: int
    def __init__(self, asset: _Optional[_Union[_asset_data_pb2.Asset, str]] = ..., amount: _Optional[int] = ...) -> None: ...
