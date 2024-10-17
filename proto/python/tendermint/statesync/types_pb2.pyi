from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Message(_message.Message):
    __slots__ = ("snapshots_request", "snapshots_response", "chunk_request", "chunk_response")
    SNAPSHOTS_REQUEST_FIELD_NUMBER: _ClassVar[int]
    SNAPSHOTS_RESPONSE_FIELD_NUMBER: _ClassVar[int]
    CHUNK_REQUEST_FIELD_NUMBER: _ClassVar[int]
    CHUNK_RESPONSE_FIELD_NUMBER: _ClassVar[int]
    snapshots_request: SnapshotsRequest
    snapshots_response: SnapshotsResponse
    chunk_request: ChunkRequest
    chunk_response: ChunkResponse
    def __init__(self, snapshots_request: _Optional[_Union[SnapshotsRequest, _Mapping]] = ..., snapshots_response: _Optional[_Union[SnapshotsResponse, _Mapping]] = ..., chunk_request: _Optional[_Union[ChunkRequest, _Mapping]] = ..., chunk_response: _Optional[_Union[ChunkResponse, _Mapping]] = ...) -> None: ...

class SnapshotsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class SnapshotsResponse(_message.Message):
    __slots__ = ("height", "format", "chunks", "hash", "metadata")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    FORMAT_FIELD_NUMBER: _ClassVar[int]
    CHUNKS_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    height: int
    format: int
    chunks: int
    hash: bytes
    metadata: bytes
    def __init__(self, height: _Optional[int] = ..., format: _Optional[int] = ..., chunks: _Optional[int] = ..., hash: _Optional[bytes] = ..., metadata: _Optional[bytes] = ...) -> None: ...

class ChunkRequest(_message.Message):
    __slots__ = ("height", "format", "index")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    FORMAT_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    height: int
    format: int
    index: int
    def __init__(self, height: _Optional[int] = ..., format: _Optional[int] = ..., index: _Optional[int] = ...) -> None: ...

class ChunkResponse(_message.Message):
    __slots__ = ("height", "format", "index", "chunk", "missing")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    FORMAT_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    CHUNK_FIELD_NUMBER: _ClassVar[int]
    MISSING_FIELD_NUMBER: _ClassVar[int]
    height: int
    format: int
    index: int
    chunk: bytes
    missing: bool
    def __init__(self, height: _Optional[int] = ..., format: _Optional[int] = ..., index: _Optional[int] = ..., chunk: _Optional[bytes] = ..., missing: bool = ...) -> None: ...
