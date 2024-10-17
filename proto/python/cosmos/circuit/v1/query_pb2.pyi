from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from cosmos.circuit.v1 import types_pb2 as _types_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.query.v1 import query_pb2 as _query_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryAccountRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class AccountResponse(_message.Message):
    __slots__ = ("permission",)
    PERMISSION_FIELD_NUMBER: _ClassVar[int]
    permission: _types_pb2.Permissions
    def __init__(self, permission: _Optional[_Union[_types_pb2.Permissions, _Mapping]] = ...) -> None: ...

class QueryAccountsRequest(_message.Message):
    __slots__ = ("pagination",)
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class AccountsResponse(_message.Message):
    __slots__ = ("accounts", "pagination")
    ACCOUNTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    accounts: _containers.RepeatedCompositeFieldContainer[_types_pb2.GenesisAccountPermissions]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, accounts: _Optional[_Iterable[_Union[_types_pb2.GenesisAccountPermissions, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryDisabledListRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class DisabledListResponse(_message.Message):
    __slots__ = ("disabled_list",)
    DISABLED_LIST_FIELD_NUMBER: _ClassVar[int]
    disabled_list: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, disabled_list: _Optional[_Iterable[str]] = ...) -> None: ...
