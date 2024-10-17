from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from ibc.core.client.v1 import client_pb2 as _client_pb2
from ibc.applications.transfer.v1 import transfer_pb2 as _transfer_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgTransfer(_message.Message):
    __slots__ = ("source_port", "source_channel", "token", "sender", "receiver", "timeout_height", "timeout_timestamp", "memo", "tokens", "forwarding")
    SOURCE_PORT_FIELD_NUMBER: _ClassVar[int]
    SOURCE_CHANNEL_FIELD_NUMBER: _ClassVar[int]
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    TOKENS_FIELD_NUMBER: _ClassVar[int]
    FORWARDING_FIELD_NUMBER: _ClassVar[int]
    source_port: str
    source_channel: str
    token: _coin_pb2.Coin
    sender: str
    receiver: str
    timeout_height: _client_pb2.Height
    timeout_timestamp: int
    memo: str
    tokens: _containers.RepeatedCompositeFieldContainer[_coin_pb2.Coin]
    forwarding: _transfer_pb2.Forwarding
    def __init__(self, source_port: _Optional[str] = ..., source_channel: _Optional[str] = ..., token: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., timeout_height: _Optional[_Union[_client_pb2.Height, _Mapping]] = ..., timeout_timestamp: _Optional[int] = ..., memo: _Optional[str] = ..., tokens: _Optional[_Iterable[_Union[_coin_pb2.Coin, _Mapping]]] = ..., forwarding: _Optional[_Union[_transfer_pb2.Forwarding, _Mapping]] = ...) -> None: ...

class MsgTransferResponse(_message.Message):
    __slots__ = ("sequence",)
    SEQUENCE_FIELD_NUMBER: _ClassVar[int]
    sequence: int
    def __init__(self, sequence: _Optional[int] = ...) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("signer", "params")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    signer: str
    params: _transfer_pb2.Params
    def __init__(self, signer: _Optional[str] = ..., params: _Optional[_Union[_transfer_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
