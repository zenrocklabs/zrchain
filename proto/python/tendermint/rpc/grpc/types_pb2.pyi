from tendermint.abci import types_pb2 as _types_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class RequestPing(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RequestBroadcastTx(_message.Message):
    __slots__ = ("tx",)
    TX_FIELD_NUMBER: _ClassVar[int]
    tx: bytes
    def __init__(self, tx: _Optional[bytes] = ...) -> None: ...

class ResponsePing(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ResponseBroadcastTx(_message.Message):
    __slots__ = ("check_tx", "tx_result")
    CHECK_TX_FIELD_NUMBER: _ClassVar[int]
    TX_RESULT_FIELD_NUMBER: _ClassVar[int]
    check_tx: _types_pb2.ResponseCheckTx
    tx_result: _types_pb2.ExecTxResult
    def __init__(self, check_tx: _Optional[_Union[_types_pb2.ResponseCheckTx, _Mapping]] = ..., tx_result: _Optional[_Union[_types_pb2.ExecTxResult, _Mapping]] = ...) -> None: ...
