from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Asset(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNSPECIFIED: _ClassVar[Asset]
    ROCK: _ClassVar[Asset]
    BTC: _ClassVar[Asset]
    ETH: _ClassVar[Asset]
    ZEC: _ClassVar[Asset]
UNSPECIFIED: Asset
ROCK: Asset
BTC: Asset
ETH: Asset
ZEC: Asset

class AssetData(_message.Message):
    __slots__ = ("asset", "priceUSD", "precision")
    ASSET_FIELD_NUMBER: _ClassVar[int]
    PRICEUSD_FIELD_NUMBER: _ClassVar[int]
    PRECISION_FIELD_NUMBER: _ClassVar[int]
    asset: Asset
    priceUSD: str
    precision: int
    def __init__(self, asset: _Optional[_Union[Asset, str]] = ..., priceUSD: _Optional[str] = ..., precision: _Optional[int] = ...) -> None: ...
