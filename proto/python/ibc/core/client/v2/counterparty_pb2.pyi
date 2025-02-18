from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class CounterpartyInfo(_message.Message):
    __slots__ = ("merkle_prefix", "client_id")
    MERKLE_PREFIX_FIELD_NUMBER: _ClassVar[int]
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    merkle_prefix: _containers.RepeatedScalarFieldContainer[bytes]
    client_id: str
    def __init__(self, merkle_prefix: _Optional[_Iterable[bytes]] = ..., client_id: _Optional[str] = ...) -> None: ...
