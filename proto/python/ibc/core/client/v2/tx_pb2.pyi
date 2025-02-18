from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MsgRegisterCounterparty(_message.Message):
    __slots__ = ("client_id", "counterparty_merkle_prefix", "counterparty_client_id", "signer")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_MERKLE_PREFIX_FIELD_NUMBER: _ClassVar[int]
    COUNTERPARTY_CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    counterparty_merkle_prefix: _containers.RepeatedScalarFieldContainer[bytes]
    counterparty_client_id: str
    signer: str
    def __init__(self, client_id: _Optional[str] = ..., counterparty_merkle_prefix: _Optional[_Iterable[bytes]] = ..., counterparty_client_id: _Optional[str] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgRegisterCounterpartyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
