from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from zenrock.policy import action_pb2 as _action_pb2
from zenrock.policy import params_pb2 as _params_pb2
from zenrock.policy import policy_pb2 as _policy_pb2
from google.protobuf import any_pb2 as _any_pb2
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

class QueryActionsRequest(_message.Message):
    __slots__ = ("pagination",)
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryActionsResponse(_message.Message):
    __slots__ = ("pagination", "actions")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    ACTIONS_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageResponse
    actions: _containers.RepeatedCompositeFieldContainer[_action_pb2.Action]
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ..., actions: _Optional[_Iterable[_Union[_action_pb2.Action, _Mapping]]] = ...) -> None: ...

class QueryPoliciesRequest(_message.Message):
    __slots__ = ("pagination",)
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class PolicyResponse(_message.Message):
    __slots__ = ("policy", "metadata")
    POLICY_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    policy: _policy_pb2.Policy
    metadata: _any_pb2.Any
    def __init__(self, policy: _Optional[_Union[_policy_pb2.Policy, _Mapping]] = ..., metadata: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class QueryPoliciesResponse(_message.Message):
    __slots__ = ("pagination", "policies")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    POLICIES_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageResponse
    policies: _containers.RepeatedCompositeFieldContainer[PolicyResponse]
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ..., policies: _Optional[_Iterable[_Union[PolicyResponse, _Mapping]]] = ...) -> None: ...

class QueryPolicyByIdRequest(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class QueryPolicyByIdResponse(_message.Message):
    __slots__ = ("policy",)
    POLICY_FIELD_NUMBER: _ClassVar[int]
    policy: PolicyResponse
    def __init__(self, policy: _Optional[_Union[PolicyResponse, _Mapping]] = ...) -> None: ...

class QueryActionsByAddressRequest(_message.Message):
    __slots__ = ("pagination", "address", "status")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    address: str
    status: _action_pb2.ActionStatus
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ..., address: _Optional[str] = ..., status: _Optional[_Union[_action_pb2.ActionStatus, str]] = ...) -> None: ...

class QueryActionsByAddressResponse(_message.Message):
    __slots__ = ("pagination", "actions")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    ACTIONS_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageResponse
    actions: _containers.RepeatedCompositeFieldContainer[_action_pb2.Action]
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ..., actions: _Optional[_Iterable[_Union[_action_pb2.Action, _Mapping]]] = ...) -> None: ...

class QueryGetSignMethodsForAddressRequest(_message.Message):
    __slots__ = ("pagination", "address")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageRequest
    address: str
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ..., address: _Optional[str] = ...) -> None: ...

class QueryGetSignMethodsForAddressResponse(_message.Message):
    __slots__ = ("pagination", "config")
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    pagination: _pagination_pb2.PageResponse
    config: _containers.RepeatedCompositeFieldContainer[_any_pb2.Any]
    def __init__(self, pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ..., config: _Optional[_Iterable[_Union[_any_pb2.Any, _Mapping]]] = ...) -> None: ...
