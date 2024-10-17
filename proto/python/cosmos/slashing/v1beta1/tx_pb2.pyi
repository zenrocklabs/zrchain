from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.slashing.v1beta1 import slashing_pb2 as _slashing_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from amino import amino_pb2 as _amino_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgUnjail(_message.Message):
    __slots__ = ("validator_addr",)
    VALIDATOR_ADDR_FIELD_NUMBER: _ClassVar[int]
    validator_addr: str
    def __init__(self, validator_addr: _Optional[str] = ...) -> None: ...

class MsgUnjailResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _slashing_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_slashing_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
