from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ActionStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ACTION_STATUS_UNSPECIFIED: _ClassVar[ActionStatus]
    ACTION_STATUS_PENDING: _ClassVar[ActionStatus]
    ACTION_STATUS_COMPLETED: _ClassVar[ActionStatus]
    ACTION_STATUS_REVOKED: _ClassVar[ActionStatus]
    ACTION_STATUS_TIMEOUT: _ClassVar[ActionStatus]
ACTION_STATUS_UNSPECIFIED: ActionStatus
ACTION_STATUS_PENDING: ActionStatus
ACTION_STATUS_COMPLETED: ActionStatus
ACTION_STATUS_REVOKED: ActionStatus
ACTION_STATUS_TIMEOUT: ActionStatus

class Action(_message.Message):
    __slots__ = ("id", "approvers", "status", "policy_id", "msg", "creator", "btl", "policy_data")
    ID_FIELD_NUMBER: _ClassVar[int]
    APPROVERS_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    POLICY_DATA_FIELD_NUMBER: _ClassVar[int]
    id: int
    approvers: _containers.RepeatedScalarFieldContainer[str]
    status: ActionStatus
    policy_id: int
    msg: _any_pb2.Any
    creator: str
    btl: int
    policy_data: _containers.RepeatedCompositeFieldContainer[KeyValue]
    def __init__(self, id: _Optional[int] = ..., approvers: _Optional[_Iterable[str]] = ..., status: _Optional[_Union[ActionStatus, str]] = ..., policy_id: _Optional[int] = ..., msg: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., creator: _Optional[str] = ..., btl: _Optional[int] = ..., policy_data: _Optional[_Iterable[_Union[KeyValue, _Mapping]]] = ...) -> None: ...

class KeyValue(_message.Message):
    __slots__ = ("key", "value")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    key: str
    value: bytes
    def __init__(self, key: _Optional[str] = ..., value: _Optional[bytes] = ...) -> None: ...

class ActionResponse(_message.Message):
    __slots__ = ("id", "approvers", "status", "policy_id", "msg", "creator", "btl", "policy_data")
    ID_FIELD_NUMBER: _ClassVar[int]
    APPROVERS_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    POLICY_DATA_FIELD_NUMBER: _ClassVar[int]
    id: int
    approvers: _containers.RepeatedScalarFieldContainer[str]
    status: str
    policy_id: int
    msg: _any_pb2.Any
    creator: str
    btl: int
    policy_data: _containers.RepeatedCompositeFieldContainer[KeyValue]
    def __init__(self, id: _Optional[int] = ..., approvers: _Optional[_Iterable[str]] = ..., status: _Optional[str] = ..., policy_id: _Optional[int] = ..., msg: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., creator: _Optional[str] = ..., btl: _Optional[int] = ..., policy_data: _Optional[_Iterable[_Union[KeyValue, _Mapping]]] = ...) -> None: ...
