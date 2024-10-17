from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Permissions(_message.Message):
    __slots__ = ("level", "limit_type_urls")
    class Level(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        LEVEL_NONE_UNSPECIFIED: _ClassVar[Permissions.Level]
        LEVEL_SOME_MSGS: _ClassVar[Permissions.Level]
        LEVEL_ALL_MSGS: _ClassVar[Permissions.Level]
        LEVEL_SUPER_ADMIN: _ClassVar[Permissions.Level]
    LEVEL_NONE_UNSPECIFIED: Permissions.Level
    LEVEL_SOME_MSGS: Permissions.Level
    LEVEL_ALL_MSGS: Permissions.Level
    LEVEL_SUPER_ADMIN: Permissions.Level
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    LIMIT_TYPE_URLS_FIELD_NUMBER: _ClassVar[int]
    level: Permissions.Level
    limit_type_urls: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, level: _Optional[_Union[Permissions.Level, str]] = ..., limit_type_urls: _Optional[_Iterable[str]] = ...) -> None: ...

class GenesisAccountPermissions(_message.Message):
    __slots__ = ("address", "permissions")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    PERMISSIONS_FIELD_NUMBER: _ClassVar[int]
    address: str
    permissions: Permissions
    def __init__(self, address: _Optional[str] = ..., permissions: _Optional[_Union[Permissions, _Mapping]] = ...) -> None: ...

class GenesisState(_message.Message):
    __slots__ = ("account_permissions", "disabled_type_urls")
    ACCOUNT_PERMISSIONS_FIELD_NUMBER: _ClassVar[int]
    DISABLED_TYPE_URLS_FIELD_NUMBER: _ClassVar[int]
    account_permissions: _containers.RepeatedCompositeFieldContainer[GenesisAccountPermissions]
    disabled_type_urls: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, account_permissions: _Optional[_Iterable[_Union[GenesisAccountPermissions, _Mapping]]] = ..., disabled_type_urls: _Optional[_Iterable[str]] = ...) -> None: ...
