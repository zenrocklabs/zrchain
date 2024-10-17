from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from amino import amino_pb2 as _amino_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgVerifyInvariant(_message.Message):
    __slots__ = ("sender", "invariant_module_name", "invariant_route")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    INVARIANT_MODULE_NAME_FIELD_NUMBER: _ClassVar[int]
    INVARIANT_ROUTE_FIELD_NUMBER: _ClassVar[int]
    sender: str
    invariant_module_name: str
    invariant_route: str
    def __init__(self, sender: _Optional[str] = ..., invariant_module_name: _Optional[str] = ..., invariant_route: _Optional[str] = ...) -> None: ...

class MsgVerifyInvariantResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "constant_fee")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    CONSTANT_FEE_FIELD_NUMBER: _ClassVar[int]
    authority: str
    constant_fee: _coin_pb2.Coin
    def __init__(self, authority: _Optional[str] = ..., constant_fee: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
