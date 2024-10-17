from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Txs(_message.Message):
    __slots__ = ("txs",)
    TXS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, txs: _Optional[_Iterable[bytes]] = ...) -> None: ...

class Message(_message.Message):
    __slots__ = ("txs",)
    TXS_FIELD_NUMBER: _ClassVar[int]
    txs: Txs
    def __init__(self, txs: _Optional[_Union[Txs, _Mapping]] = ...) -> None: ...
