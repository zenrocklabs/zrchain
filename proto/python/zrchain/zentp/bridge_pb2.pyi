from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zentp import params_pb2 as _params_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class BridgeStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    BRIDGE_STATUS_UNSPECIFIED: _ClassVar[BridgeStatus]
    BRIDGE_STATUS_PENDING: _ClassVar[BridgeStatus]
    BRIDGE_STATUS_COMPLETED: _ClassVar[BridgeStatus]
    BRIDGE_STATUS_FAILED: _ClassVar[BridgeStatus]
BRIDGE_STATUS_UNSPECIFIED: BridgeStatus
BRIDGE_STATUS_PENDING: BridgeStatus
BRIDGE_STATUS_COMPLETED: BridgeStatus
BRIDGE_STATUS_FAILED: BridgeStatus

class Bridge(_message.Message):
    __slots__ = ("id", "denom", "creator", "source_address", "source_chain", "destination_chain", "amount", "recipient_address", "tx_id", "tx_hash", "state", "block_height", "awaiting_event_since")
    ID_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SOURCE_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    SOURCE_CHAIN_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CHAIN_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    TX_ID_FIELD_NUMBER: _ClassVar[int]
    TX_HASH_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    AWAITING_EVENT_SINCE_FIELD_NUMBER: _ClassVar[int]
    id: int
    denom: str
    creator: str
    source_address: str
    source_chain: str
    destination_chain: str
    amount: int
    recipient_address: str
    tx_id: int
    tx_hash: str
    state: BridgeStatus
    block_height: int
    awaiting_event_since: int
    def __init__(self, id: _Optional[int] = ..., denom: _Optional[str] = ..., creator: _Optional[str] = ..., source_address: _Optional[str] = ..., source_chain: _Optional[str] = ..., destination_chain: _Optional[str] = ..., amount: _Optional[int] = ..., recipient_address: _Optional[str] = ..., tx_id: _Optional[int] = ..., tx_hash: _Optional[str] = ..., state: _Optional[_Union[BridgeStatus, str]] = ..., block_height: _Optional[int] = ..., awaiting_event_since: _Optional[int] = ...) -> None: ...
