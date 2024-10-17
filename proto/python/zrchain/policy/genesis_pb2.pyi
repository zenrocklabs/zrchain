from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.policy import params_pb2 as _params_pb2
from zrchain.policy import action_pb2 as _action_pb2
from zrchain.policy import policy_pb2 as _policy_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("params", "port_id", "policies", "actions", "sign_methods")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    POLICIES_FIELD_NUMBER: _ClassVar[int]
    ACTIONS_FIELD_NUMBER: _ClassVar[int]
    SIGN_METHODS_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    port_id: str
    policies: _containers.RepeatedCompositeFieldContainer[_policy_pb2.Policy]
    actions: _containers.RepeatedCompositeFieldContainer[_action_pb2.Action]
    sign_methods: _containers.RepeatedCompositeFieldContainer[GenesisSignMethod]
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ..., port_id: _Optional[str] = ..., policies: _Optional[_Iterable[_Union[_policy_pb2.Policy, _Mapping]]] = ..., actions: _Optional[_Iterable[_Union[_action_pb2.Action, _Mapping]]] = ..., sign_methods: _Optional[_Iterable[_Union[GenesisSignMethod, _Mapping]]] = ...) -> None: ...

class GenesisSignMethod(_message.Message):
    __slots__ = ("owner", "id", "config")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    owner: str
    id: str
    config: _any_pb2.Any
    def __init__(self, owner: _Optional[str] = ..., id: _Optional[str] = ..., config: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...
