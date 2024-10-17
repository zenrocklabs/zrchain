from ibc.applications.transfer.v2 import token_pb2 as _token_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from ibc.applications.transfer.v1 import transfer_pb2 as _transfer_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class FungibleTokenPacketData(_message.Message):
    __slots__ = ("denom", "amount", "sender", "receiver", "memo")
    DENOM_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    denom: str
    amount: str
    sender: str
    receiver: str
    memo: str
    def __init__(self, denom: _Optional[str] = ..., amount: _Optional[str] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., memo: _Optional[str] = ...) -> None: ...

class FungibleTokenPacketDataV2(_message.Message):
    __slots__ = ("tokens", "sender", "receiver", "memo", "forwarding")
    TOKENS_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    FORWARDING_FIELD_NUMBER: _ClassVar[int]
    tokens: _containers.RepeatedCompositeFieldContainer[_token_pb2.Token]
    sender: str
    receiver: str
    memo: str
    forwarding: ForwardingPacketData
    def __init__(self, tokens: _Optional[_Iterable[_Union[_token_pb2.Token, _Mapping]]] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., memo: _Optional[str] = ..., forwarding: _Optional[_Union[ForwardingPacketData, _Mapping]] = ...) -> None: ...

class ForwardingPacketData(_message.Message):
    __slots__ = ("destination_memo", "hops")
    DESTINATION_MEMO_FIELD_NUMBER: _ClassVar[int]
    HOPS_FIELD_NUMBER: _ClassVar[int]
    destination_memo: str
    hops: _containers.RepeatedCompositeFieldContainer[_transfer_pb2.Hop]
    def __init__(self, destination_memo: _Optional[str] = ..., hops: _Optional[_Iterable[_Union[_transfer_pb2.Hop, _Mapping]]] = ...) -> None: ...
