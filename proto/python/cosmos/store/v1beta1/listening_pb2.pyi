from cometbft.abci.v2 import types_pb2 as _types_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class StoreKVPair(_message.Message):
    __slots__ = ("store_key", "delete", "key", "value")
    STORE_KEY_FIELD_NUMBER: _ClassVar[int]
    DELETE_FIELD_NUMBER: _ClassVar[int]
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    store_key: str
    delete: bool
    key: bytes
    value: bytes
    def __init__(self, store_key: _Optional[str] = ..., delete: bool = ..., key: _Optional[bytes] = ..., value: _Optional[bytes] = ...) -> None: ...

class BlockMetadata(_message.Message):
    __slots__ = ("response_commit", "request_finalize_block", "response_finalize_block")
    RESPONSE_COMMIT_FIELD_NUMBER: _ClassVar[int]
    REQUEST_FINALIZE_BLOCK_FIELD_NUMBER: _ClassVar[int]
    RESPONSE_FINALIZE_BLOCK_FIELD_NUMBER: _ClassVar[int]
    response_commit: _types_pb2.CommitResponse
    request_finalize_block: _types_pb2.FinalizeBlockRequest
    response_finalize_block: _types_pb2.FinalizeBlockResponse
    def __init__(self, response_commit: _Optional[_Union[_types_pb2.CommitResponse, _Mapping]] = ..., request_finalize_block: _Optional[_Union[_types_pb2.FinalizeBlockRequest, _Mapping]] = ..., response_finalize_block: _Optional[_Union[_types_pb2.FinalizeBlockResponse, _Mapping]] = ...) -> None: ...
