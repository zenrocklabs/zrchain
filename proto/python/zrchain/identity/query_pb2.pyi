from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from zrchain.identity import params_pb2 as _params_pb2
from zrchain.identity import workspace_pb2 as _workspace_pb2
from zrchain.identity import keyring_pb2 as _keyring_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryParamsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryParamsResponse(_message.Message):
    __slots__ = ("params",)
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class QueryWorkspacesRequest(_message.Message):
    __slots__ = ("pagination", "owner", "creator")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    OWNER_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    owner: str
    creator: str
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ..., owner: _Optional[str] = ..., creator: _Optional[str] = ...) -> None: ...

class QueryWorkspacesResponse(_message.Message):
    __slots__ = ("workspaces", "pagination")
    WORKSPACES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    workspaces: _containers.RepeatedCompositeFieldContainer[_workspace_pb2.Workspace]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, workspaces: _Optional[_Iterable[_Union[_workspace_pb2.Workspace, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryWorkspaceByAddressRequest(_message.Message):
    __slots__ = ("workspace_addr",)
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    workspace_addr: str
    def __init__(self, workspace_addr: _Optional[str] = ...) -> None: ...

class QueryWorkspaceByAddressResponse(_message.Message):
    __slots__ = ("workspace",)
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    workspace: _workspace_pb2.Workspace
    def __init__(self, workspace: _Optional[_Union[_workspace_pb2.Workspace, _Mapping]] = ...) -> None: ...

class QueryKeyringsRequest(_message.Message):
    __slots__ = ("pagination",)
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryKeyringsResponse(_message.Message):
    __slots__ = ("keyrings", "pagination")
    KEYRINGS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    keyrings: _containers.RepeatedCompositeFieldContainer[_keyring_pb2.Keyring]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, keyrings: _Optional[_Iterable[_Union[_keyring_pb2.Keyring, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryKeyringByAddressRequest(_message.Message):
    __slots__ = ("keyring_addr",)
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    keyring_addr: str
    def __init__(self, keyring_addr: _Optional[str] = ...) -> None: ...

class QueryKeyringByAddressResponse(_message.Message):
    __slots__ = ("keyring",)
    KEYRING_FIELD_NUMBER: _ClassVar[int]
    keyring: _keyring_pb2.Keyring
    def __init__(self, keyring: _Optional[_Union[_keyring_pb2.Keyring, _Mapping]] = ...) -> None: ...
