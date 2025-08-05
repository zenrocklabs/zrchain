from zrchain.zentp import params_pb2 as _params_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Dct(_message.Message):
    __slots__ = ("denom", "solana")
    DENOM_FIELD_NUMBER: _ClassVar[int]
    SOLANA_FIELD_NUMBER: _ClassVar[int]
    denom: str
    solana: _params_pb2.Solana
    def __init__(self, denom: _Optional[str] = ..., solana: _Optional[_Union[_params_pb2.Solana, _Mapping]] = ...) -> None: ...
