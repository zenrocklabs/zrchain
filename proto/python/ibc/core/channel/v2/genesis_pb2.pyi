from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("acknowledgements", "commitments", "receipts", "async_packets", "send_sequences")
    ACKNOWLEDGEMENTS_FIELD_NUMBER: _ClassVar[int]
    COMMITMENTS_FIELD_NUMBER: _ClassVar[int]
    RECEIPTS_FIELD_NUMBER: _ClassVar[int]
    ASYNC_PACKETS_FIELD_NUMBER: _ClassVar[int]
    SEND_SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    acknowledgements: _containers.RepeatedCompositeFieldContainer[PacketState]
    commitments: _containers.RepeatedCompositeFieldContainer[PacketState]
    receipts: _containers.RepeatedCompositeFieldContainer[PacketState]
    async_packets: _containers.RepeatedCompositeFieldContainer[PacketState]
    send_sequences: _containers.RepeatedCompositeFieldContainer[PacketSequence]
    def __init__(self, acknowledgements: _Optional[_Iterable[_Union[PacketState, _Mapping]]] = ..., commitments: _Optional[_Iterable[_Union[PacketState, _Mapping]]] = ..., receipts: _Optional[_Iterable[_Union[PacketState, _Mapping]]] = ..., async_packets: _Optional[_Iterable[_Union[PacketState, _Mapping]]] = ..., send_sequences: _Optional[_Iterable[_Union[PacketSequence, _Mapping]]] = ...) -> None: ...

class PacketState(_message.Message):
    __slots__ = ("client_id", "sequence", "data")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    sequence: int
    data: bytes
    def __init__(self, client_id: _Optional[str] = ..., sequence: _Optional[int] = ..., data: _Optional[bytes] = ...) -> None: ...

class PacketSequence(_message.Message):
    __slots__ = ("client_id", "sequence")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    sequence: int
    def __init__(self, client_id: _Optional[str] = ..., sequence: _Optional[int] = ...) -> None: ...
