from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from ibc.core.client.v1 import client_pb2 as _client_pb2
from ibc.core.channel.v1 import channel_pb2 as _channel_pb2
from ibc.core.channel.v1 import upgrade_pb2 as _upgrade_pb2
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

class MsgChannelOpenInit(_message.Message):
    __slots__ = ("port_id", "channel", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel: _channel_pb2.Channel
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel: _Optional[_Union[_channel_pb2.Channel, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelOpenInitResponse(_message.Message):
    __slots__ = ("channel_id", "version")
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    channel_id: str
    version: str
    def __init__(self, channel_id: _Optional[str] = ..., version: _Optional[str] = ...) -> None: ...

class MsgChannelOpenTry(_message.Message):
    __slots__ = ("port_id", "previous_channel_id", "channel", "counterparty_version", "proof_init", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    PREVIOUS_CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_VERSION_FIELD_NUMBER: _ClassVar[int]
    PROOF_INIT_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    previous_channel_id: str
    channel: _channel_pb2.Channel
    counterparty_version: str
    proof_init: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., previous_channel_id: _Optional[str] = ..., channel: _Optional[_Union[_channel_pb2.Channel, _Mapping]] = ..., counterparty_version: _Optional[str] = ..., proof_init: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelOpenTryResponse(_message.Message):
    __slots__ = ("version", "channel_id")
    VERSION_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    version: str
    channel_id: str
    def __init__(self, version: _Optional[str] = ..., channel_id: _Optional[str] = ...) -> None: ...

class MsgChannelOpenAck(_message.Message):
    __slots__ = ("port_id", "channel_id", "counterparty_channel_id", "counterparty_version", "proof_try", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_VERSION_FIELD_NUMBER: _ClassVar[int]
    PROOF_TRY_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    counterparty_channel_id: str
    counterparty_version: str
    proof_try: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., counterparty_channel_id: _Optional[str] = ..., counterparty_version: _Optional[str] = ..., proof_try: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelOpenAckResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgChannelOpenConfirm(_message.Message):
    __slots__ = ("port_id", "channel_id", "proof_ack", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    PROOF_ACK_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    proof_ack: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., proof_ack: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelOpenConfirmResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgChannelCloseInit(_message.Message):
    __slots__ = ("port_id", "channel_id", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelCloseInitResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgChannelCloseConfirm(_message.Message):
    __slots__ = ("port_id", "channel_id", "proof_init", "proof_height", "signer", "counterparty_upgrade_sequence")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    PROOF_INIT_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    proof_init: bytes
    proof_height: _client_pb2.Height
    signer: str
    counterparty_upgrade_sequence: int
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., proof_init: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ..., counterparty_upgrade_sequence: _Optional[int] = ...) -> None: ...

class MsgChannelCloseConfirmResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRecvPacket(_message.Message):
    __slots__ = ("packet", "proof_commitment", "proof_height", "signer")
    PACKET_FIELD_NUMBER: _ClassVar[int]
    PROOF_COMMITMENT_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    packet: _channel_pb2.Packet
    proof_commitment: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, packet: _Optional[_Union[_channel_pb2.Packet, _Mapping]] = ..., proof_commitment: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgRecvPacketResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgTimeout(_message.Message):
    __slots__ = ("packet", "proof_unreceived", "proof_height", "next_sequence_recv", "signer")
    PACKET_FIELD_NUMBER: _ClassVar[int]
    PROOF_UNRECEIVED_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    NEXT_SEQUENCE_RECV_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    packet: _channel_pb2.Packet
    proof_unreceived: bytes
    proof_height: _client_pb2.Height
    next_sequence_recv: int
    signer: str
    def __init__(self, packet: _Optional[_Union[_channel_pb2.Packet, _Mapping]] = ..., proof_unreceived: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., next_sequence_recv: _Optional[int] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgTimeoutResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgTimeoutOnClose(_message.Message):
    __slots__ = ("packet", "proof_unreceived", "proof_close", "proof_height", "next_sequence_recv", "signer", "counterparty_upgrade_sequence")
    PACKET_FIELD_NUMBER: _ClassVar[int]
    PROOF_UNRECEIVED_FIELD_NUMBER: _ClassVar[int]
    PROOF_CLOSE_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    NEXT_SEQUENCE_RECV_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    packet: _channel_pb2.Packet
    proof_unreceived: bytes
    proof_close: bytes
    proof_height: _client_pb2.Height
    next_sequence_recv: int
    signer: str
    counterparty_upgrade_sequence: int
    def __init__(self, packet: _Optional[_Union[_channel_pb2.Packet, _Mapping]] = ..., proof_unreceived: _Optional[bytes] = ..., proof_close: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., next_sequence_recv: _Optional[int] = ..., signer: _Optional[str] = ..., counterparty_upgrade_sequence: _Optional[int] = ...) -> None: ...

class MsgTimeoutOnCloseResponse(_message.Message):
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
    packet: _channel_pb2.Packet
    acknowledgement: bytes
    proof_acked: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, packet: _Optional[_Union[_channel_pb2.Packet, _Mapping]] = ..., acknowledgement: _Optional[bytes] = ..., proof_acked: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgAcknowledgementResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgChannelUpgradeInit(_message.Message):
    __slots__ = ("port_id", "channel_id", "fields", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    fields: _upgrade_pb2.UpgradeFields
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., fields: _Optional[_Union[_upgrade_pb2.UpgradeFields, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeInitResponse(_message.Message):
    __slots__ = ("upgrade", "upgrade_sequence")
    UPGRADE_FIELD_NUMBER: _ClassVar[int]
    UPGRADE_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    upgrade: _upgrade_pb2.Upgrade
    upgrade_sequence: int
    def __init__(self, upgrade: _Optional[_Union[_upgrade_pb2.Upgrade, _Mapping]] = ..., upgrade_sequence: _Optional[int] = ...) -> None: ...

class MsgChannelUpgradeTry(_message.Message):
    __slots__ = ("port_id", "channel_id", "proposed_upgrade_connection_hops", "counterparty_upgrade_fields", "counterparty_upgrade_sequence", "proof_channel", "proof_upgrade", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    PROPOSED_UPGRADE_CONNECTION_HOPS_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_FIELDS_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    PROOF_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    PROOF_UPGRADE_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    proposed_upgrade_connection_hops: _containers.RepeatedScalarFieldContainer[str]
    counterparty_upgrade_fields: _upgrade_pb2.UpgradeFields
    counterparty_upgrade_sequence: int
    proof_channel: bytes
    proof_upgrade: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., proposed_upgrade_connection_hops: _Optional[_Iterable[str]] = ..., counterparty_upgrade_fields: _Optional[_Union[_upgrade_pb2.UpgradeFields, _Mapping]] = ..., counterparty_upgrade_sequence: _Optional[int] = ..., proof_channel: _Optional[bytes] = ..., proof_upgrade: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeTryResponse(_message.Message):
    __slots__ = ("upgrade", "upgrade_sequence", "result")
    UPGRADE_FIELD_NUMBER: _ClassVar[int]
    UPGRADE_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    upgrade: _upgrade_pb2.Upgrade
    upgrade_sequence: int
    result: ResponseResultType
    def __init__(self, upgrade: _Optional[_Union[_upgrade_pb2.Upgrade, _Mapping]] = ..., upgrade_sequence: _Optional[int] = ..., result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgChannelUpgradeAck(_message.Message):
    __slots__ = ("port_id", "channel_id", "counterparty_upgrade", "proof_channel", "proof_upgrade", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_FIELD_NUMBER: _ClassVar[int]
    PROOF_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    PROOF_UPGRADE_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    counterparty_upgrade: _upgrade_pb2.Upgrade
    proof_channel: bytes
    proof_upgrade: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., counterparty_upgrade: _Optional[_Union[_upgrade_pb2.Upgrade, _Mapping]] = ..., proof_channel: _Optional[bytes] = ..., proof_upgrade: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeAckResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgChannelUpgradeConfirm(_message.Message):
    __slots__ = ("port_id", "channel_id", "counterparty_channel_state", "counterparty_upgrade", "proof_channel", "proof_upgrade", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_CHANNEL_STATE_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_FIELD_NUMBER: _ClassVar[int]
    PROOF_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    PROOF_UPGRADE_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    counterparty_channel_state: _channel_pb2.State
    counterparty_upgrade: _upgrade_pb2.Upgrade
    proof_channel: bytes
    proof_upgrade: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., counterparty_channel_state: _Optional[_Union[_channel_pb2.State, str]] = ..., counterparty_upgrade: _Optional[_Union[_upgrade_pb2.Upgrade, _Mapping]] = ..., proof_channel: _Optional[bytes] = ..., proof_upgrade: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeConfirmResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: ResponseResultType
    def __init__(self, result: _Optional[_Union[ResponseResultType, str]] = ...) -> None: ...

class MsgChannelUpgradeOpen(_message.Message):
    __slots__ = ("port_id", "channel_id", "counterparty_channel_state", "counterparty_upgrade_sequence", "proof_channel", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_CHANNEL_STATE_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_UPGRADE_SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    PROOF_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    counterparty_channel_state: _channel_pb2.State
    counterparty_upgrade_sequence: int
    proof_channel: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., counterparty_channel_state: _Optional[_Union[_channel_pb2.State, str]] = ..., counterparty_upgrade_sequence: _Optional[int] = ..., proof_channel: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeOpenResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgChannelUpgradeTimeout(_message.Message):
    __slots__ = ("port_id", "channel_id", "counterparty_channel", "proof_channel", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    PROOF_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    counterparty_channel: _channel_pb2.Channel
    proof_channel: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., counterparty_channel: _Optional[_Union[_channel_pb2.Channel, _Mapping]] = ..., proof_channel: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeTimeoutResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgChannelUpgradeCancel(_message.Message):
    __slots__ = ("port_id", "channel_id", "error_receipt", "proof_error_receipt", "proof_height", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_RECEIPT_FIELD_NUMBER: _ClassVar[int]
    PROOF_ERROR_RECEIPT_FIELD_NUMBER: _ClassVar[int]
    PROOF_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    error_receipt: _upgrade_pb2.ErrorReceipt
    proof_error_receipt: bytes
    proof_height: _client_pb2.Height
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., error_receipt: _Optional[_Union[_upgrade_pb2.ErrorReceipt, _Mapping]] = ..., proof_error_receipt: _Optional[bytes] = ..., proof_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgChannelUpgradeCancelResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _channel_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_channel_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgPruneAcknowledgements(_message.Message):
    __slots__ = ("port_id", "channel_id", "limit", "signer")
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    port_id: str
    channel_id: str
    limit: int
    signer: str
    def __init__(self, port_id: _Optional[str] = ..., channel_id: _Optional[str] = ..., limit: _Optional[int] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgPruneAcknowledgementsResponse(_message.Message):
    __slots__ = ("total_pruned_sequences", "total_remaining_sequences")
    TOTAL_PRUNED_SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_REMAINING_SEQUENCES_FIELD_NUMBER: _ClassVar[int]
    total_pruned_sequences: int
    total_remaining_sequences: int
    def __init__(self, total_pruned_sequences: _Optional[int] = ..., total_remaining_sequences: _Optional[int] = ...) -> None: ...
