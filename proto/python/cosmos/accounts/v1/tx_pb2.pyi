from google.protobuf import any_pb2 as _any_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from cosmos.tx.v1beta1 import tx_pb2 as _tx_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgInit(_message.Message):
    __slots__ = ("sender", "account_type", "message", "funds")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_TYPE_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    sender: str
    account_type: str
    message: _any_pb2.Any
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, sender: _Optional[str] = ..., account_type: _Optional[str] = ..., message: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class MsgInitResponse(_message.Message):
    __slots__ = ("account_address", "response")
    ACCOUNT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    account_address: str
    response: _any_pb2.Any
    def __init__(self, account_address: _Optional[str] = ..., response: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class MsgExecute(_message.Message):
    __slots__ = ("sender", "target", "message", "funds")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    TARGET_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    FUNDS_FIELD_NUMBER: _ClassVar[int]
    sender: str
    target: str
    message: _any_pb2.Any
    funds: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    def __init__(self, sender: _Optional[str] = ..., target: _Optional[str] = ..., message: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., funds: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ...) -> None: ...

class MsgExecuteResponse(_message.Message):
    __slots__ = ("response",)
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    response: _any_pb2.Any
    def __init__(self, response: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class MsgExecuteBundle(_message.Message):
    __slots__ = ("bundler", "txs")
    BUNDLER_FIELD_NUMBER: _ClassVar[int]
    TXS_FIELD_NUMBER: _ClassVar[int]
    bundler: str
    txs: _containers.RepeatedCompositeFieldContainer[_tx_pb2.TxRaw]
    def __init__(self, bundler: _Optional[str] = ..., txs: _Optional[_Iterable[_Union[_tx_pb2.TxRaw, _Mapping]]] = ...) -> None: ...

class BundledTxResponse(_message.Message):
    __slots__ = ("exec_responses", "error")
    EXEC_RESPONSES_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    exec_responses: _any_pb2.Any
    error: str
    def __init__(self, exec_responses: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., error: _Optional[str] = ...) -> None: ...

class MsgExecuteBundleResponse(_message.Message):
    __slots__ = ("responses",)
    RESPONSES_FIELD_NUMBER: _ClassVar[int]
    responses: _containers.RepeatedCompositeFieldContainer[BundledTxResponse]
    def __init__(self, responses: _Optional[_Iterable[_Union[BundledTxResponse, _Mapping]]] = ...) -> None: ...
