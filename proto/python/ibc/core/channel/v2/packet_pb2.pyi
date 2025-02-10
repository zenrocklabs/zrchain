from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PacketStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PACKET_STATUS_UNSPECIFIED: _ClassVar[PacketStatus]
    PACKET_STATUS_SUCCESS: _ClassVar[PacketStatus]
    PACKET_STATUS_FAILURE: _ClassVar[PacketStatus]
    PACKET_STATUS_ASYNC: _ClassVar[PacketStatus]
PACKET_STATUS_UNSPECIFIED: PacketStatus
PACKET_STATUS_SUCCESS: PacketStatus
PACKET_STATUS_FAILURE: PacketStatus
PACKET_STATUS_ASYNC: PacketStatus

class Packet(_message.Message):
    __slots__ = ("sequence", "source_client", "destination_client", "timeout_timestamp", "payloads")
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_CLIENT_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CLIENT_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    PAYLOADS_FIELD_NUMBER: _ClassVar[int]
    sequence: int
    source_client: str
    destination_client: str
    timeout_timestamp: int
    payloads: _containers.RepeatedCompositeFieldContainer[Payload]
    def __init__(self, sequence: _Optional[int] = ..., source_client: _Optional[str] = ..., destination_client: _Optional[str] = ..., timeout_timestamp: _Optional[int] = ..., payloads: _Optional[_Iterable[_Union[Payload, _Mapping]]] = ...) -> None: ...

class Payload(_message.Message):
    __slots__ = ("source_port", "destination_port", "version", "encoding", "value")
    SOURCE_PORT_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_PORT_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    ENCODING_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    source_port: str
    destination_port: str
    version: str
    encoding: str
    value: bytes
    def __init__(self, source_port: _Optional[str] = ..., destination_port: _Optional[str] = ..., version: _Optional[str] = ..., encoding: _Optional[str] = ..., value: _Optional[bytes] = ...) -> None: ...

class Acknowledgement(_message.Message):
    __slots__ = ("app_acknowledgements",)
    APP_ACKNOWLEDGEMENTS_FIELD_NUMBER: _ClassVar[int]
    app_acknowledgements: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, app_acknowledgements: _Optional[_Iterable[bytes]] = ...) -> None: ...

class RecvPacketResult(_message.Message):
    __slots__ = ("status", "acknowledgement")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    ACKNOWLEDGEMENT_FIELD_NUMBER: _ClassVar[int]
    status: PacketStatus
    acknowledgement: bytes
    def __init__(self, status: _Optional[_Union[PacketStatus, str]] = ..., acknowledgement: _Optional[bytes] = ...) -> None: ...
