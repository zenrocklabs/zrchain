from tendermint.abci import types_pb2 as _types_pb2
from cosmos.store.v1beta1 import listening_pb2 as _listening_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListenFinalizeBlockRequest(_message.Message):
    __slots__ = ("req", "res")
    REQ_FIELD_NUMBER: _ClassVar[int]
    RES_FIELD_NUMBER: _ClassVar[int]
    req: _types_pb2.RequestFinalizeBlock
    res: _types_pb2.ResponseFinalizeBlock
    def __init__(self, req: _Optional[_Union[_types_pb2.RequestFinalizeBlock, _Mapping]] = ..., res: _Optional[_Union[_types_pb2.ResponseFinalizeBlock, _Mapping]] = ...) -> None: ...

class ListenFinalizeBlockResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ListenCommitRequest(_message.Message):
    __slots__ = ("block_height", "res", "change_set")
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    RES_FIELD_NUMBER: _ClassVar[int]
    CHANGE_SET_FIELD_NUMBER: _ClassVar[int]
    block_height: int
    res: _types_pb2.ResponseCommit
    change_set: _containers.RepeatedCompositeFieldContainer[_listening_pb2.StoreKVPair]
    def __init__(self, block_height: _Optional[int] = ..., res: _Optional[_Union[_types_pb2.ResponseCommit, _Mapping]] = ..., change_set: _Optional[_Iterable[_Union[_listening_pb2.StoreKVPair, _Mapping]]] = ...) -> None: ...

class ListenCommitResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
