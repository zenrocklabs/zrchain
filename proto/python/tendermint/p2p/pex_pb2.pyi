from tendermint.p2p import types_pb2 as _types_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PexRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class PexAddrs(_message.Message):
    __slots__ = ("addrs",)
    ADDRS_FIELD_NUMBER: _ClassVar[int]
    addrs: _containers.RepeatedCompositeFieldContainer[_types_pb2.NetAddress]
    def __init__(self, addrs: _Optional[_Iterable[_Union[_types_pb2.NetAddress, _Mapping]]] = ...) -> None: ...

class Message(_message.Message):
    __slots__ = ("pex_request", "pex_addrs")
    PEX_REQUEST_FIELD_NUMBER: _ClassVar[int]
    PEX_ADDRS_FIELD_NUMBER: _ClassVar[int]
    pex_request: PexRequest
    pex_addrs: PexAddrs
    def __init__(self, pex_request: _Optional[_Union[PexRequest, _Mapping]] = ..., pex_addrs: _Optional[_Union[PexAddrs, _Mapping]] = ...) -> None: ...
