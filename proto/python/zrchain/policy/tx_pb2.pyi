from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.policy import params_pb2 as _params_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _params_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewPolicy(_message.Message):
    __slots__ = ("creator", "name", "policy", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    POLICY_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    name: str
    policy: _any_pb2.Any
    btl: int
    def __init__(self, creator: _Optional[str] = ..., name: _Optional[str] = ..., policy: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgNewPolicyResponse(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class MsgRevokeAction(_message.Message):
    __slots__ = ("creator", "action_id")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ACTION_ID_FIELD_NUMBER: _ClassVar[int]
    creator: str
    action_id: int
    def __init__(self, creator: _Optional[str] = ..., action_id: _Optional[int] = ...) -> None: ...

class MsgRevokeActionResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgApproveAction(_message.Message):
    __slots__ = ("creator", "action_type", "action_id", "additional_signatures")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ACTION_TYPE_FIELD_NUMBER: _ClassVar[int]
    ACTION_ID_FIELD_NUMBER: _ClassVar[int]
    ADDITIONAL_SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    creator: str
    action_type: str
    action_id: int
    additional_signatures: _containers.RepeatedCompositeFieldContainer[_any_pb2.Any]
    def __init__(self, creator: _Optional[str] = ..., action_type: _Optional[str] = ..., action_id: _Optional[int] = ..., additional_signatures: _Optional[_Iterable[_Union[_any_pb2.Any, _Mapping]]] = ...) -> None: ...

class MsgApproveActionResponse(_message.Message):
    __slots__ = ("status",)
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: str
    def __init__(self, status: _Optional[str] = ...) -> None: ...

class MsgAddSignMethod(_message.Message):
    __slots__ = ("creator", "config")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    creator: str
    config: _any_pb2.Any
    def __init__(self, creator: _Optional[str] = ..., config: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class MsgAddSignMethodResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRemoveSignMethod(_message.Message):
    __slots__ = ("creator", "id")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    creator: str
    id: str
    def __init__(self, creator: _Optional[str] = ..., id: _Optional[str] = ...) -> None: ...

class MsgRemoveSignMethodResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgAddMultiGrant(_message.Message):
    __slots__ = ("creator", "grantee", "msgs")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    GRANTEE_FIELD_NUMBER: _ClassVar[int]
    MSGS_FIELD_NUMBER: _ClassVar[int]
    creator: str
    grantee: str
    msgs: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, creator: _Optional[str] = ..., grantee: _Optional[str] = ..., msgs: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgAddMultiGrantResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRemoveMultiGrant(_message.Message):
    __slots__ = ("creator", "grantee", "msgs")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    GRANTEE_FIELD_NUMBER: _ClassVar[int]
    MSGS_FIELD_NUMBER: _ClassVar[int]
    creator: str
    grantee: str
    msgs: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, creator: _Optional[str] = ..., grantee: _Optional[str] = ..., msgs: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgRemoveMultiGrantResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
