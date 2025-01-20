from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.identity import params_pb2 as _params_pb2
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

class MsgNewWorkspace(_message.Message):
    __slots__ = ("creator", "admin_policy_id", "sign_policy_id", "additional_owners")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ADMIN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ADDITIONAL_OWNERS_FIELD_NUMBER: _ClassVar[int]
    creator: str
    admin_policy_id: int
    sign_policy_id: int
    additional_owners: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, creator: _Optional[str] = ..., admin_policy_id: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., additional_owners: _Optional[_Iterable[str]] = ...) -> None: ...

class MsgNewWorkspaceResponse(_message.Message):
    __slots__ = ("addr",)
    ADDR_FIELD_NUMBER: _ClassVar[int]
    addr: str
    def __init__(self, addr: _Optional[str] = ...) -> None: ...

class MsgAddWorkspaceOwner(_message.Message):
    __slots__ = ("creator", "workspace_addr", "new_owner", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    NEW_OWNER_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    workspace_addr: str
    new_owner: str
    btl: int
    def __init__(self, creator: _Optional[str] = ..., workspace_addr: _Optional[str] = ..., new_owner: _Optional[str] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgAddWorkspaceOwnerResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgAppendChildWorkspace(_message.Message):
    __slots__ = ("creator", "parent_workspace_addr", "child_workspace_addr", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    PARENT_WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    CHILD_WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    parent_workspace_addr: str
    child_workspace_addr: str
    btl: int
    def __init__(self, creator: _Optional[str] = ..., parent_workspace_addr: _Optional[str] = ..., child_workspace_addr: _Optional[str] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgAppendChildWorkspaceResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewChildWorkspace(_message.Message):
    __slots__ = ("creator", "parent_workspace_addr", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    PARENT_WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    parent_workspace_addr: str
    btl: int
    def __init__(self, creator: _Optional[str] = ..., parent_workspace_addr: _Optional[str] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgNewChildWorkspaceResponse(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class MsgRemoveWorkspaceOwner(_message.Message):
    __slots__ = ("creator", "workspace_addr", "owner", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    OWNER_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    workspace_addr: str
    owner: str
    btl: int
    def __init__(self, creator: _Optional[str] = ..., workspace_addr: _Optional[str] = ..., owner: _Optional[str] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgRemoveWorkspaceOwnerResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewKeyring(_message.Message):
    __slots__ = ("creator", "description", "party_threshold", "key_req_fee", "sig_req_fee", "delegate_fees", "mpc_minimum_timeout", "mpc_default_timeout")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    PARTY_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    KEY_REQ_FEE_FIELD_NUMBER: _ClassVar[int]
    SIG_REQ_FEE_FIELD_NUMBER: _ClassVar[int]
    DELEGATE_FEES_FIELD_NUMBER: _ClassVar[int]
    MPC_MINIMUM_TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    MPC_DEFAULT_TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    creator: str
    description: str
    party_threshold: int
    key_req_fee: int
    sig_req_fee: int
    delegate_fees: bool
    mpc_minimum_timeout: int
    mpc_default_timeout: int
    def __init__(self, creator: _Optional[str] = ..., description: _Optional[str] = ..., party_threshold: _Optional[int] = ..., key_req_fee: _Optional[int] = ..., sig_req_fee: _Optional[int] = ..., delegate_fees: bool = ..., mpc_minimum_timeout: _Optional[int] = ..., mpc_default_timeout: _Optional[int] = ...) -> None: ...

class MsgNewKeyringResponse(_message.Message):
    __slots__ = ("addr",)
    ADDR_FIELD_NUMBER: _ClassVar[int]
    addr: str
    def __init__(self, addr: _Optional[str] = ...) -> None: ...

class MsgAddKeyringParty(_message.Message):
    __slots__ = ("creator", "keyring_addr", "party", "increase_threshold")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    PARTY_FIELD_NUMBER: _ClassVar[int]
    INCREASE_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring_addr: str
    party: str
    increase_threshold: bool
    def __init__(self, creator: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., party: _Optional[str] = ..., increase_threshold: bool = ...) -> None: ...

class MsgAddKeyringPartyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateKeyring(_message.Message):
    __slots__ = ("creator", "keyring_addr", "party_threshold", "key_req_fee", "sig_req_fee", "description", "is_active", "mpc_minimum_timeout", "mpc_default_timeout")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    PARTY_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    KEY_REQ_FEE_FIELD_NUMBER: _ClassVar[int]
    SIG_REQ_FEE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    IS_ACTIVE_FIELD_NUMBER: _ClassVar[int]
    MPC_MINIMUM_TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    MPC_DEFAULT_TIMEOUT_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring_addr: str
    party_threshold: int
    key_req_fee: int
    sig_req_fee: int
    description: str
    is_active: bool
    mpc_minimum_timeout: int
    mpc_default_timeout: int
    def __init__(self, creator: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., party_threshold: _Optional[int] = ..., key_req_fee: _Optional[int] = ..., sig_req_fee: _Optional[int] = ..., description: _Optional[str] = ..., is_active: bool = ..., mpc_minimum_timeout: _Optional[int] = ..., mpc_default_timeout: _Optional[int] = ...) -> None: ...

class MsgUpdateKeyringResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRemoveKeyringParty(_message.Message):
    __slots__ = ("creator", "keyring_addr", "party", "decrease_threshold")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    PARTY_FIELD_NUMBER: _ClassVar[int]
    DECREASE_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring_addr: str
    party: str
    decrease_threshold: bool
    def __init__(self, creator: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., party: _Optional[str] = ..., decrease_threshold: bool = ...) -> None: ...

class MsgRemoveKeyringPartyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgAddKeyringAdmin(_message.Message):
    __slots__ = ("creator", "keyring_addr", "admin")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring_addr: str
    admin: str
    def __init__(self, creator: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., admin: _Optional[str] = ...) -> None: ...

class MsgAddKeyringAdminResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRemoveKeyringAdmin(_message.Message):
    __slots__ = ("creator", "keyring_addr", "admin")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    ADMIN_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring_addr: str
    admin: str
    def __init__(self, creator: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., admin: _Optional[str] = ...) -> None: ...

class MsgRemoveKeyringAdminResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateWorkspace(_message.Message):
    __slots__ = ("creator", "workspace_addr", "admin_policy_id", "sign_policy_id", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    ADMIN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    workspace_addr: str
    admin_policy_id: int
    sign_policy_id: int
    btl: int
    def __init__(self, creator: _Optional[str] = ..., workspace_addr: _Optional[str] = ..., admin_policy_id: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgUpdateWorkspaceResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgDeactivateKeyring(_message.Message):
    __slots__ = ("creator", "keyring_addr")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring_addr: str
    def __init__(self, creator: _Optional[str] = ..., keyring_addr: _Optional[str] = ...) -> None: ...

class MsgDeactivateKeyringResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
