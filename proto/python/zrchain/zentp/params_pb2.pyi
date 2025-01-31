from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("solana_relayer_key_id", "zrchain_relayer_key_id")
    SOLANA_RELAYER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    ZRCHAIN_RELAYER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    solana_relayer_key_id: int
    zrchain_relayer_key_id: int
    def __init__(self, solana_relayer_key_id: _Optional[int] = ..., zrchain_relayer_key_id: _Optional[int] = ...) -> None: ...
