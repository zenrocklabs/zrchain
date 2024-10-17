from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos.circuit.v1 import types_pb2 as _types_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgAuthorizeCircuitBreaker(_message.Message):
    __slots__ = ("granter", "grantee", "permissions")
    GRANTER_FIELD_NUMBER: _ClassVar[int]
    GRANTEE_FIELD_NUMBER: _ClassVar[int]
    PERMISSIONS_FIELD_NUMBER: _ClassVar[int]
    granter: str
    grantee: str
    permissions: _types_pb2.Permissions
    def __init__(self, granter: _Optional[str] = ..., grantee: _Optional[str] = ..., permissions: _Optional[_Union[_types_pb2.Permissions, _Mapping]] = ...) -> None: ...

class MsgAuthorizeCircuitBreakerResponse(_message.Message):
    __slots__ = ("success",)
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...

class MsgTripCircuitBreaker(_message.Message):
    __slots__ = ("authority", "msg_type_urls")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MSG_TYPE_URLS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    msg_type_urls: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, authority: _Optional[str] = ..., msg_type_urls: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgTripCircuitBreakerResponse(_message.Message):
    __slots__ = ("success",)
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...

class MsgResetCircuitBreaker(_message.Message):
    __slots__ = ("authority", "msg_type_urls")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MSG_TYPE_URLS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    msg_type_urls: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, authority: _Optional[str] = ..., msg_type_urls: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgResetCircuitBreakerResponse(_message.Message):
    __slots__ = ("success",)
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...
