from cosmos.tx.v1beta1 import tx_pb2 as _tx_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgAuthenticate(_message.Message):
    __slots__ = ("bundler", "raw_tx", "tx", "signer_index")
    BUNDLER_FIELD_NUMBER: _ClassVar[int]
    RAW_TX_FIELD_NUMBER: _ClassVar[int]
    TX_FIELD_NUMBER: _ClassVar[int]
    SIGNER_INDEX_FIELD_NUMBER: _ClassVar[int]
    bundler: str
    raw_tx: _tx_pb2.TxRaw
    tx: _tx_pb2.Tx
    signer_index: int
    def __init__(self, bundler: _Optional[str] = ..., raw_tx: _Optional[_Union[_tx_pb2.TxRaw, _Mapping]] = ..., tx: _Optional[_Union[_tx_pb2.Tx, _Mapping]] = ..., signer_index: _Optional[int] = ...) -> None: ...

class MsgAuthenticateResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryAuthenticationMethods(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryAuthenticationMethodsResponse(_message.Message):
    __slots__ = ("authentication_methods",)
    AUTHENTICATION_METHODS_FIELD_NUMBER: _ClassVar[int]
    authentication_methods: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, authentication_methods: _Optional[_Iterable[str]] = ...) -> None: ...
