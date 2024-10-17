from gogoproto import gogo_pb2 as _gogo_pb2
from ibc.core.client.v1 import client_pb2 as _client_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ClientState(_message.Message):
    __slots__ = ("data", "checksum", "latest_height")
    DATA_FIELD_NUMBER: _ClassVar[int]
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    LATEST_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    checksum: bytes
    latest_height: _client_pb2.Height
    def __init__(self, data: _Optional[bytes] = ..., checksum: _Optional[bytes] = ..., latest_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ...) -> None: ...

class ConsensusState(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class ClientMessage(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class Checksums(_message.Message):
    __slots__ = ("checksums",)
    CHECKSUMS_FIELD_NUMBER: _ClassVar[int]
    checksums: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, checksums: _Optional[_Iterable[bytes]] = ...) -> None: ...
