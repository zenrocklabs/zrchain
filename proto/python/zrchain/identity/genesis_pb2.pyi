from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.identity import params_pb2 as _params_pb2
from zrchain.identity import keyring_pb2 as _keyring_pb2
from zrchain.identity import workspace_pb2 as _workspace_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GenesisState(_message.Message):
    __slots__ = ("params", "port_id", "keyrings", "workspaces")
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    PORT_ID_FIELD_NUMBER: _ClassVar[int]
    KEYRINGS_FIELD_NUMBER: _ClassVar[int]
    WORKSPACES_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    port_id: str
    keyrings: _containers.RepeatedCompositeFieldContainer[_keyring_pb2.Keyring]
    workspaces: _containers.RepeatedCompositeFieldContainer[_workspace_pb2.Workspace]
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ..., port_id: _Optional[str] = ..., keyrings: _Optional[_Iterable[_Union[_keyring_pb2.Keyring, _Mapping]]] = ..., workspaces: _Optional[_Iterable[_Union[_workspace_pb2.Workspace, _Mapping]]] = ...) -> None: ...
