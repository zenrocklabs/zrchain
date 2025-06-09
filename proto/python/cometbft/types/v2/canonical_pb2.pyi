from gogoproto import gogo_pb2 as _gogo_pb2
from cometbft.types.v2 import types_pb2 as _types_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CanonicalBlockID(_message.Message):
    __slots__ = ("hash", "part_set_header")
    HASH_FIELD_NUMBER: _ClassVar[int]
    PART_SET_HEADER_FIELD_NUMBER: _ClassVar[int]
    hash: bytes
    part_set_header: CanonicalPartSetHeader
    def __init__(self, hash: _Optional[bytes] = ..., part_set_header: _Optional[_Union[CanonicalPartSetHeader, _Mapping]] = ...) -> None: ...

class CanonicalPartSetHeader(_message.Message):
    __slots__ = ("total", "hash")
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    total: int
    hash: bytes
    def __init__(self, total: _Optional[int] = ..., hash: _Optional[bytes] = ...) -> None: ...

class CanonicalProposal(_message.Message):
    __slots__ = ("type", "height", "round", "pol_round", "block_id", "timestamp", "chain_id")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    POL_ROUND_FIELD_NUMBER: _ClassVar[int]
    BLOCK_ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    type: _types_pb2.SignedMsgType
    height: int
    round: int
    pol_round: int
    block_id: CanonicalBlockID
    timestamp: _timestamp_pb2.Timestamp
    chain_id: str
    def __init__(self, type: _Optional[_Union[_types_pb2.SignedMsgType, str]] = ..., height: _Optional[int] = ..., round: _Optional[int] = ..., pol_round: _Optional[int] = ..., block_id: _Optional[_Union[CanonicalBlockID, _Mapping]] = ..., timestamp: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., chain_id: _Optional[str] = ...) -> None: ...

class CanonicalVote(_message.Message):
    __slots__ = ("type", "height", "round", "block_id", "timestamp", "chain_id")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    BLOCK_ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    type: _types_pb2.SignedMsgType
    height: int
    round: int
    block_id: CanonicalBlockID
    timestamp: _timestamp_pb2.Timestamp
    chain_id: str
    def __init__(self, type: _Optional[_Union[_types_pb2.SignedMsgType, str]] = ..., height: _Optional[int] = ..., round: _Optional[int] = ..., block_id: _Optional[_Union[CanonicalBlockID, _Mapping]] = ..., timestamp: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., chain_id: _Optional[str] = ...) -> None: ...

class CanonicalVoteExtension(_message.Message):
    __slots__ = ("extension", "height", "round", "chain_id")
    EXTENSION_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    extension: bytes
    height: int
    round: int
    chain_id: str
    def __init__(self, extension: _Optional[bytes] = ..., height: _Optional[int] = ..., round: _Optional[int] = ..., chain_id: _Optional[str] = ...) -> None: ...
