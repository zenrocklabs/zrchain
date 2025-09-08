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

class MsgSwapRequest(_message.Message):
    __slots__ = ("creator", "pair", "workspace", "amount_in", "rock_key_id", "btc_key_id", "destination_caip2")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    PAIR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_IN_FIELD_NUMBER: _ClassVar[int]
    ROCK_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    BTC_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CAIP2_FIELD_NUMBER: _ClassVar[int]
    creator: str
    pair: str
    workspace: str
    amount_in: int
    rock_key_id: int
    btc_key_id: int
    destination_caip2: str
    def __init__(self, creator: _Optional[str] = ..., pair: _Optional[str] = ..., workspace: _Optional[str] = ..., amount_in: _Optional[int] = ..., rock_key_id: _Optional[int] = ..., btc_key_id: _Optional[int] = ..., destination_caip2: _Optional[str] = ...) -> None: ...

class MsgSwapRequestResponse(_message.Message):
    __slots__ = ("swap_id",)
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    swap_id: int
    def __init__(self, swap_id: _Optional[int] = ...) -> None: ...
