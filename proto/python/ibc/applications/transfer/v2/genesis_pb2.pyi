from ibc.applications.transfer.v1 import transfer_pb2 as _transfer_pb2
from ibc.applications.transfer.v2 import token_pb2 as _token_pb2
from ibc.core.channel.v1 import channel_pb2 as _channel_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("port_id", "denoms", "params", "total_escrowed", "forwarded_packets")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    DENOMS_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_ESCROWED_FIELD_NUMBER: _ClassVar[int]
    FORWARDED_PACKETS_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    denoms: _containers.RepeatedCompositeFieldContainer[_token_pb2.Denom]
    params: _transfer_pb2.Params
    total_escrowed: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    forwarded_packets: _containers.RepeatedCompositeFieldContainer[ForwardedPacket]
    def __init__(self, port_id: _Optional[str] = ..., denoms: _Optional[_Iterable[_Union[_token_pb2.Denom, _Mapping]]] = ..., params: _Optional[_Union[_transfer_pb2.Params, _Mapping]] = ..., total_escrowed: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., forwarded_packets: _Optional[_Iterable[_Union[ForwardedPacket, _Mapping]]] = ...) -> None: ...

class ForwardedPacket(_message.Message):
    __slots__ = ("forward_key", "packet")
    FORWARD_KEY_FIELD_NUMBER: _ClassVar[int]
    PACKET_FIELD_NUMBER: _ClassVar[int]
    forward_key: _channel_pb2.PacketId
    packet: _channel_pb2.Packet
    def __init__(self, forward_key: _Optional[_Union[_channel_pb2.PacketId, _Mapping]] = ..., packet: _Optional[_Union[_channel_pb2.Packet, _Mapping]] = ...) -> None: ...
