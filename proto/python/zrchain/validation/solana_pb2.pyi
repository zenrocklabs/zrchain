from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SolanaNonce(_message.Message):
    __slots__ = ("nonce",)
    NONCE_FIELD_NUMBER: _ClassVar[int]
    nonce: bytes
    def __init__(self, nonce: _Optional[bytes] = ...) -> None: ...

class SolanaCounters(_message.Message):
    __slots__ = ("mint_counter", "redemption_counter")
    MINT_COUNTER_FIELD_NUMBER: _ClassVar[int]
    REDEMPTION_COUNTER_FIELD_NUMBER: _ClassVar[int]
    mint_counter: int
    redemption_counter: int
    def __init__(self, mint_counter: _Optional[int] = ..., redemption_counter: _Optional[int] = ...) -> None: ...
