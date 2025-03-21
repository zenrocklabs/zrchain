from ibc.applications.transfer.v2 import token_pb2 as _token_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

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
    __slots__ = ("token", "sender", "receiver", "memo")
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    token: _token_pb2.Token
    sender: str
    receiver: str
    memo: str
    def __init__(self, token: _Optional[_Union[_token_pb2.Token, _Mapping]] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., memo: _Optional[str] = ...) -> None: ...
