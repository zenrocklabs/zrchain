from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from ibc.applications.transfer.v1 import transfer_pb2 as _transfer_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Allocation(_message.Message):
    __slots__ = ("source_port", "source_channel", "spend_limit", "allow_list", "allowed_packet_data")
    SOURCE_PORT_FIELD_NUMBER: _ClassVar[int]
    SOURCE_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    SPEND_LIMIT_FIELD_NUMBER: _ClassVar[int]
    ALLOW_LIST_FIELD_NUMBER: _ClassVar[int]
    ALLOWED_PACKET_DATA_FIELD_NUMBER: _ClassVar[int]
    source_port: str
    source_channel: str
    spend_limit: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    allow_list: _containers.RepeatedScalarFieldContainer[str]
    allowed_packet_data: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, source_port: _Optional[str] = ..., source_channel: _Optional[str] = ..., spend_limit: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., allow_list: _Optional[_Iterable[str]] = ..., allowed_packet_data: _Optional[_Iterable[str]] = ...) -> None: ...

class TransferAuthorization(_message.Message):
    __slots__ = ("allocations",)
    ALLOCATIONS_FIELD_NUMBER: _ClassVar[int]
    allocations: _containers.RepeatedCompositeFieldContainer[Allocation]
    def __init__(self, allocations: _Optional[_Iterable[_Union[Allocation, _Mapping]]] = ...) -> None: ...
