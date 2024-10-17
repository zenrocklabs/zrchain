from cosmos.app.v1alpha1 import module_pb2 as _module_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Module(_message.Message):
    __slots__ = ("bech32_prefix", "module_account_permissions", "authority")
    BECH32_PREFIX_FIELD_NUMBER: _ClassVar[int]
    MODULE_ACCOUNT_PERMISSIONS_FIELD_NUMBER: _ClassVar[int]
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    bech32_prefix: str
    module_account_permissions: _containers.RepeatedCompositeFieldContainer[ModuleAccountPermission]
    authority: str
    def __init__(self, bech32_prefix: _Optional[str] = ..., module_account_permissions: _Optional[_Iterable[_Union[ModuleAccountPermission, _Mapping]]] = ..., authority: _Optional[str] = ...) -> None: ...

class ModuleAccountPermission(_message.Message):
    __slots__ = ("account", "permissions")
    ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    PERMISSIONS_FIELD_NUMBER: _ClassVar[int]
    account: str
    permissions: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, account: _Optional[str] = ..., permissions: _Optional[_Iterable[str]] = ...) -> None: ...
