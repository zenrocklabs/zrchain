from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from ibc.core.channel.v2 import packet_pb2 as _packet_pb2
from ibc.core.client.v1 import client_pb2 as _client_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ResponseResultType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RESPONSE_RESULT_TYPE_UNSPECIFIED: _ClassVar[ResponseResultType]
    RESPONSE_RESULT_TYPE_NOOP: _ClassVar[ResponseResultType]
    RESPONSE_RESULT_TYPE_SUCCESS: _ClassVar[ResponseResultType]
    RESPONSE_RESULT_TYPE_FAILURE: _ClassVar[ResponseResultType]
RESPONSE_RESULT_TYPE_UNSPECIFIED: ResponseResultType
RESPONSE_RESULT_TYPE_NOOP: ResponseResultType
RESPONSE_RESULT_TYPE_SUCCESS: ResponseResultType
RESPONSE_RESULT_TYPE_FAILURE: ResponseResultType

class MsgSendPacket(_message.Message):
    __slots__ = ("source_client", "timeout_timestamp", "payloads", "signer")
    SOURCE_CLIENT_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    PAYLOADS_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    source_client: str
    timeout_timestamp: int
    payloads: _containers.RepeatedCompositeFieldContainer[_packet_pb2.Payload]
    signer: str
    def __init__(self, source_client: _Optional[str] = ..., timeout_timestamp: _Optional[int] = ..., payloads: _Optional[_Iterable[_Union[_packet_pb2.Payload, _Mapping]]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgSendPacketResponse(_message.Message):
    __slots__ = ("sequence",)
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    sequence: int
    def __init__(self, sequence: _Optional[int] = ...) -> None: ...

class MsgRecvPacket(_message.Message):
    __slots__ = ("packet", "proof_commitment", "proof_height", "signer")
    PACKET_FIELD_NUMBER: _ClassVar[int]
    PROOF_COMMITMENT_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    packet: _packet_pb2.Packet
    proof_commitment: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, packet: _Optional[_Union[_packet_pb2.Packet, _Mapping]] = ..., proof_commitment: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgRecvPacketResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgTimeout(_message.Message):
    __slots__ = ("packet", "proof_unreceived", "proof_height", "signer")
    PACKET_FIELD_NUMBER: _ClassVar[int]
    PROOF_UNRECEIVED_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    packet: _packet_pb2.Packet
    proof_unreceived: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, packet: _Optional[_Union[_packet_pb2.Packet, _Mapping]] = ..., proof_unreceived: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgTimeoutResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgAcknowledgement(_message.Message):
    __slots__ = ("packet", "acknowledgement", "proof_acked", "proof_height", "signer")
    PACKET_FIELD_NUMBER: _ClassVar[int]
    ACKNOWLEDGEMENT_FIELD_NUMBER: _ClassVar[int]
    PROOF_ACKED_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    packet: _packet_pb2.Packet
    acknowledgement: _packet_pb2.Acknowledgement
    proof_acked: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, packet: _Optional[_Union[_packet_pb2.Packet, _Mapping]] = ..., acknowledgement: _Optional[_Union[_packet_pb2.Acknowledgement, _Mapping]] = ..., proof_acked: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgAcknowledgementResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...
