from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from ibc.applications.interchain_accounts.host.v1 import host_pb2 as _host_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgUpdateParams(_message.Message):
    __slots__ = ("signer", "params")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    signer: str
    params: _host_pb2.Params
    def __init__(self, signer: _Optional[str] = ..., params: _Optional[_Union[_host_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgModuleQuerySafe(_message.Message):
    __slots__ = ("signer", "requests")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    REQUESTS_FIELD_NUMBER: _ClassVar[int]
    signer: str
    requests: _containers.RepeatedCompositeFieldContainer[_host_pb2.QueryRequest]
    def __init__(self, signer: _Optional[str] = ..., requests: _Optional[_Iterable[_Union[_host_pb2.QueryRequest, _Mapping]]] = ...) -> None: ...

class MsgModuleQuerySafeResponse(_message.Message):
    __slots__ = ("height", "responses")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    RESPONSES_FIELD_NUMBER: _ClassVar[int]
    height: int
    responses: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, height: _Optional[int] = ..., responses: _Optional[_Iterable[bytes]] = ...) -> None: ...
