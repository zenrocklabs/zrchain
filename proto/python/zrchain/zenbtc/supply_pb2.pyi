from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Supply(_message.Message):
    __slots__ = ("custodiedBTC", "mintedZenBTC", "pendingZenBTC")
    CUSTODIEDBTC_FIELD_NUMBER: _ClassVar[int]
    MINTEDZENBTC_FIELD_NUMBER: _ClassVar[int]
    PENDINGZENBTC_FIELD_NUMBER: _ClassVar[int]
    custodiedBTC: int
    mintedZenBTC: int
    pendingZenBTC: int
    def __init__(self, custodiedBTC: _Optional[int] = ..., mintedZenBTC: _Optional[int] = ..., pendingZenBTC: _Optional[int] = ...) -> None: ...
