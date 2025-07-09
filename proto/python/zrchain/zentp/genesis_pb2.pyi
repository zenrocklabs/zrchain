from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zentp import params_pb2 as _params_pb2
from zrchain.zentp import bridge_pb2 as _bridge_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("params", "mints", "burns", "solana_rock_supply", "zentp_fees")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    MINTS_FIELD_NUMBER: _ClassVar[int]
    BURNS_FIELD_NUMBER: _ClassVar[int]
    SOLANA_ROCK_SUPPLY_FIELD_NUMBER: _ClassVar[int]
    ZENTP_FEES_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    mints: _containers.RepeatedCompositeFieldContainer[_bridge_pb2.Bridge]
    burns: _containers.RepeatedCompositeFieldContainer[_bridge_pb2.Bridge]
    solana_rock_supply: int
    zentp_fees: _coin_pb2.Coin
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ..., mints: _Optional[_Iterable[_Union[_bridge_pb2.Bridge, _Mapping]]] = ..., burns: _Optional[_Iterable[_Union[_bridge_pb2.Bridge, _Mapping]]] = ..., solana_rock_supply: _Optional[int] = ..., zentp_fees: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...
