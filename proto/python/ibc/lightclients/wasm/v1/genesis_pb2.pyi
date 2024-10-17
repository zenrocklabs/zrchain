from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("contracts",)
    CONTRACTS_FIELD_NUMBER: _ClassVar[int]
    contracts: _containers.RepeatedCompositeFieldContainer[Contract]
    def __init__(self, contracts: _Optional[_Iterable[_Union[Contract, _Mapping]]] = ...) -> None: ...

class Contract(_message.Message):
    __slots__ = ("code_bytes",)
    CODE_BYTES_FIELD_NUMBER: _ClassVar[int]
    code_bytes: bytes
    def __init__(self, code_bytes: _Optional[bytes] = ...) -> None: ...
