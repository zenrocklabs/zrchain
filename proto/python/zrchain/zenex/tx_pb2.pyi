from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zenex import params_pb2 as _params_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _params_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSwap(_message.Message):
    __slots__ = ("creator", "pair", "workspace", "amount_in", "sender_key", "recipient_key")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    PAIR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_IN_FIELD_NUMBER: _ClassVar[int]
    YIELD_FIELD_NUMBER: _ClassVar[int]
    SENDER_KEY_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_KEY_FIELD_NUMBER: _ClassVar[int]
    creator: str
    pair: str
    workspace: str
    amount_in: str
    sender_key: int
    recipient_key: int
    def __init__(self, creator: _Optional[str] = ..., pair: _Optional[str] = ..., workspace: _Optional[str] = ..., amount_in: _Optional[str] = ..., sender_key: _Optional[int] = ..., recipient_key: _Optional[int] = ..., **kwargs) -> None: ...

class MsgSwapResponse(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...
