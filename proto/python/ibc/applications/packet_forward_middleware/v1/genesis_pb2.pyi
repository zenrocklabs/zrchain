from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("in_flight_packets",)
    class InFlightPacketsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: InFlightPacket
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[InFlightPacket, _Mapping]] = ...) -> None: ...
    IN_FLIGHT_PACKETS_FIELD_NUMBER: _ClassVar[int]
    in_flight_packets: _containers.MessageMap[str, InFlightPacket]
    def __init__(self, in_flight_packets: _Optional[_Mapping[str, InFlightPacket]] = ...) -> None: ...

class InFlightPacket(_message.Message):
    __slots__ = ("original_sender_address", "refund_channel_id", "refund_port_id", "packet_src_channel_id", "packet_src_port_id", "packet_timeout_timestamp", "packet_timeout_height", "packet_data", "refund_sequence", "retries_remaining", "timeout", "nonrefundable")
    ORIGINAL_SENDER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    REFUND_CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    REFUND_PORT_ID_FIELD_NUMBER: _ClassVar[int]
    PACKET_SRC_CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    PACKET_SRC_PORT_ID_FIELD_NUMBER: _ClassVar[int]
    PACKET_TIMEOUT_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    PACKET_TIMEOUT_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    PACKET_DATA_FIELD_NUMBER: _ClassVar[int]
    REFUND_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    RETRIES_REMAINING_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    NONREFUNDABLE_FIELD_NUMBER: _ClassVar[int]
    original_sender_address: str
    refund_channel_id: str
    refund_port_id: str
    packet_src_channel_id: str
    packet_src_port_id: str
    packet_timeout_timestamp: int
    packet_timeout_height: str
    packet_data: bytes
    refund_sequence: int
    retries_remaining: int
    timeout: int
    nonrefundable: bool
    def __init__(self, original_sender_address: _Optional[str] = ..., refund_channel_id: _Optional[str] = ..., refund_port_id: _Optional[str] = ..., packet_src_channel_id: _Optional[str] = ..., packet_src_port_id: _Optional[str] = ..., packet_timeout_timestamp: _Optional[int] = ..., packet_timeout_height: _Optional[str] = ..., packet_data: _Optional[bytes] = ..., refund_sequence: _Optional[int] = ..., retries_remaining: _Optional[int] = ..., timeout: _Optional[int] = ..., nonrefundable: bool = ...) -> None: ...
