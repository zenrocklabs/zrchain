from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from ibc.core.channel.v2 import genesis_pb2 as _genesis_pb2
from ibc.core.client.v1 import client_pb2 as _client_pb2
from google.api import annotations_pb2 as _annotations_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryNextSequenceSendRequest(_message.Message):
    __slots__ = ("client_id",)
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    def __init__(self, client_id: _Optional[str] = ...) -> None: ...

class QueryNextSequenceSendResponse(_message.Message):
    __slots__ = ("next_sequence_send", "proof", "proof_height")
    NEXT_SEQUENCE_SEND_FIELD_NUMBER: _ClassVar[int]
    PROOF_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    next_sequence_send: int
    proof: bytes
    proof_height: _client_pb2.Height
    def __init__(self, next_sequence_send: _Optional[int] = ..., proof: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryPacketCommitmentRequest(_message.Message):
    __slots__ = ("client_id", "sequence")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    sequence: int
    def __init__(self, client_id: _Optional[str] = ..., sequence: _Optional[int] = ...) -> None: ...

class QueryPacketCommitmentResponse(_message.Message):
    __slots__ = ("commitment", "proof", "proof_height")
    COMMITMENT_FIELD_NUMBER: _ClassVar[int]
    PROOF_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    commitment: bytes
    proof: bytes
    proof_height: _client_pb2.Height
    def __init__(self, commitment: _Optional[bytes] = ..., proof: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryPacketCommitmentsRequest(_message.Message):
    __slots__ = ("client_id", "pagination")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, client_id: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryPacketCommitmentsResponse(_message.Message):
    __slots__ = ("commitments", "pagination", "height")
    COMMITMENTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    commitments: _containers.RepeatedCompositeFieldContainer[_genesis_pb2.PacketState]
    pagination: _pagination_pb2.PageResponse
    height: _client_pb2.Height
    def __init__(self, commitments: _Optional[_Iterable[_Union[_genesis_pb2.PacketState, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ..., height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryPacketAcknowledgementRequest(_message.Message):
    __slots__ = ("client_id", "sequence")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    sequence: int
    def __init__(self, client_id: _Optional[str] = ..., sequence: _Optional[int] = ...) -> None: ...

class QueryPacketAcknowledgementResponse(_message.Message):
    __slots__ = ("acknowledgement", "proof", "proof_height")
    ACKNOWLEDGEMENT_FIELD_NUMBER: _ClassVar[int]
    PROOF_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    acknowledgement: bytes
    proof: bytes
    proof_height: _client_pb2.Height
    def __init__(self, acknowledgement: _Optional[bytes] = ..., proof: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryPacketAcknowledgementsRequest(_message.Message):
    __slots__ = ("client_id", "pagination", "packet_commitment_sequences")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    PACKET_COMMITMENT_SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    pagination: _pagination_pb2.PageRequest
    packet_commitment_sequences: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, client_id: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ..., packet_commitment_sequences: _Optional[_Iterable[int]] = ...) -> None: ...

class QueryPacketAcknowledgementsResponse(_message.Message):
    __slots__ = ("acknowledgements", "pagination", "height")
    ACKNOWLEDGEMENTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    acknowledgements: _containers.RepeatedCompositeFieldContainer[_genesis_pb2.PacketState]
    pagination: _pagination_pb2.PageResponse
    height: _client_pb2.Height
    def __init__(self, acknowledgements: _Optional[_Iterable[_Union[_genesis_pb2.PacketState, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ..., height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryPacketReceiptRequest(_message.Message):
    __slots__ = ("client_id", "sequence")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    sequence: int
    def __init__(self, client_id: _Optional[str] = ..., sequence: _Optional[int] = ...) -> None: ...

class QueryPacketReceiptResponse(_message.Message):
    __slots__ = ("received", "proof", "proof_height")
    RECEIVED_FIELD_NUMBER: _ClassVar[int]
    PROOF_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    received: bool
    proof: bytes
    proof_height: _client_pb2.Height
    def __init__(self, received: bool = ..., proof: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryUnreceivedPacketsRequest(_message.Message):
    __slots__ = ("client_id", "sequences")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    sequences: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, client_id: _Optional[str] = ..., sequences: _Optional[_Iterable[int]] = ...) -> None: ...

class QueryUnreceivedPacketsResponse(_message.Message):
    __slots__ = ("sequences", "height")
    SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    sequences: _containers.RepeatedScalarFieldContainer[int]
    height: _client_pb2.Height
    def __init__(self, sequences: _Optional[_Iterable[int]] = ..., height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class QueryUnreceivedAcksRequest(_message.Message):
    __slots__ = ("client_id", "packet_ack_sequences")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    PACKET_ACK_SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    packet_ack_sequences: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, client_id: _Optional[str] = ..., packet_ack_sequences: _Optional[_Iterable[int]] = ...) -> None: ...

class QueryUnreceivedAcksResponse(_message.Message):
    __slots__ = ("sequences", "height")
    SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    sequences: _containers.RepeatedScalarFieldContainer[int]
    height: _client_pb2.Height
    def __init__(self, sequences: _Optional[_Iterable[int]] = ..., height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...
