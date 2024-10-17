from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Workspace(_message.Message):
    __slots__ = ("address", "creator", "owners", "child_workspaces", "admin_policy_id", "sign_policy_id", "alias")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    OWNERS_FIELD_NUMBER: _ClassVar[int]
    CHILD_WORKSPACES_FIELD_NUMBER: _ClassVar[int]
    ADMIN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ALIAS_FIELD_NUMBER: _ClassVar[int]
    address: str
    creator: str
    owners: _containers.RepeatedScalarFieldContainer[str]
    child_workspaces: _containers.RepeatedScalarFieldContainer[str]
    admin_policy_id: int
    sign_policy_id: int
    alias: str
    def __init__(self, address: _Optional[str] = ..., creator: _Optional[str] = ..., owners: _Optional[_Iterable[str]] = ..., child_workspaces: _Optional[_Iterable[str]] = ..., admin_policy_id: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., alias: _Optional[str] = ...) -> None: ...
