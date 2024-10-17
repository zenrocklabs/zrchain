from tendermint.types import block_pb2 as _block_pb2
from tendermint.types import types_pb2 as _types_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class BlockRequest(_message.Message):
    __slots__ = ("height",)
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    height: int
    def __init__(self, height: _Optional[int] = ...) -> None: ...

class NoBlockResponse(_message.Message):
    __slots__ = ("height",)
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    height: int
    def __init__(self, height: _Optional[int] = ...) -> None: ...

class BlockResponse(_message.Message):
    __slots__ = ("block", "ext_commit")
    BLOCK_FIELD_NUMBER: _ClassVar[int]
    EXT_COMMIT_FIELD_NUMBER: _ClassVar[int]
    block: _block_pb2.Block
    ext_commit: _types_pb2.ExtendedCommit
    def __init__(self, block: _Optional[_Union[_block_pb2.Block, _Mapping]] = ..., ext_commit: _Optional[_Union[_types_pb2.ExtendedCommit, _Mapping]] = ...) -> None: ...

class StatusRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class StatusResponse(_message.Message):
    __slots__ = ("height", "base")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    height: int
    base: int
    def __init__(self, height: _Optional[int] = ..., base: _Optional[int] = ...) -> None: ...

class Message(_message.Message):
    __slots__ = ("block_request", "no_block_response", "block_response", "status_request", "status_response")
    BLOCK_REQUEST_FIELD_NUMBER: _ClassVar[int]
    NO_BLOCK_RESPONSE_FIELD_NUMBER: _ClassVar[int]
    BLOCK_RESPONSE_FIELD_NUMBER: _ClassVar[int]
    STATUS_REQUEST_FIELD_NUMBER: _ClassVar[int]
    STATUS_RESPONSE_FIELD_NUMBER: _ClassVar[int]
    block_request: BlockRequest
    no_block_response: NoBlockResponse
    block_response: BlockResponse
    status_request: StatusRequest
    status_response: StatusResponse
    def __init__(self, block_request: _Optional[_Union[BlockRequest, _Mapping]] = ..., no_block_response: _Optional[_Union[NoBlockResponse, _Mapping]] = ..., block_response: _Optional[_Union[BlockResponse, _Mapping]] = ..., status_request: _Optional[_Union[StatusRequest, _Mapping]] = ..., status_response: _Optional[_Union[StatusResponse, _Mapping]] = ...) -> None: ...
