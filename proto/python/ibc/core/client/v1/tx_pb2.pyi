from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos.upgrade.v1beta1 import upgrade_pb2 as _upgrade_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import any_pb2 as _any_pb2
from ibc.core.client.v1 import client_pb2 as _client_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MsgCreateClient(_message.Message):
    __slots__ = ("client_state", "consensus_state", "signer")
    CLIENT_STATE_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_STATE_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    client_state: _any_pb2.Any
    consensus_state: _any_pb2.Any
    signer: str
    def __init__(self, client_state: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., consensus_state: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgCreateClientResponse(_message.Message):
    __slots__ = ("client_id",)
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    def __init__(self, client_id: _Optional[str] = ...) -> None: ...

class MsgUpdateClient(_message.Message):
    __slots__ = ("client_id", "client_message", "signer")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    CLIENT_MESSAGE_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    client_message: _any_pb2.Any
    signer: str
    def __init__(self, client_id: _Optional[str] = ..., client_message: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgUpdateClientResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpgradeClient(_message.Message):
    __slots__ = ("client_id", "client_state", "consensus_state", "proof_upgrade_client", "proof_upgrade_consensus_state", "signer")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    CLIENT_STATE_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_STATE_FIELD_NUMBER: _ClassVar[int]
    PROOF_UPGRADE_CLIENT_FIELD_NUMBER: _ClassVar[int]
    PROOF_UPGRADE_CONSENSUS_STATE_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    client_state: _any_pb2.Any
    consensus_state: _any_pb2.Any
    proof_upgrade_client: bytes
    proof_upgrade_consensus_state: bytes
    signer: str
    def __init__(self, client_id: _Optional[str] = ..., client_state: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., consensus_state: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., proof_upgrade_client: _Optional[bytes] = ..., proof_upgrade_consensus_state: _Optional[bytes] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgUpgradeClientResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSubmitMisbehaviour(_message.Message):
    __slots__ = ("client_id", "misbehaviour", "signer")
    CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    MISBEHAVIOUR_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    client_id: str
    misbehaviour: _any_pb2.Any
    signer: str
    def __init__(self, client_id: _Optional[str] = ..., misbehaviour: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgSubmitMisbehaviourResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgRecoverClient(_message.Message):
    __slots__ = ("subject_client_id", "substitute_client_id", "signer")
    SUBJECT_CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SUBSTITUTE_CLIENT_ID_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    subject_client_id: str
    substitute_client_id: str
    signer: str
    def __init__(self, subject_client_id: _Optional[str] = ..., substitute_client_id: _Optional[str] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgRecoverClientResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgIBCSoftwareUpgrade(_message.Message):
    __slots__ = ("plan", "upgraded_client_state", "signer")
    PLAN_FIELD_NUMBER: _ClassVar[int]
    UPGRADED_CLIENT_STATE_FIELD_NUMBER: _ClassVar[int]
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    plan: _upgrade_pb2.Plan
    upgraded_client_state: _any_pb2.Any
    signer: str
    def __init__(self, plan: _Optional[_Union[_upgrade_pb2.Plan, _Mapping]] = ..., upgraded_client_state: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., signer: _Optional[str] = ...) -> None: ...

class MsgIBCSoftwareUpgradeResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgUpdateParams(_message.Message):
    __slots__ = ("signer", "params")
    SIGNER_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    signer: str
    params: _client_pb2.Params
    def __init__(self, signer: _Optional[str] = ..., params: _Optional[_Union[_client_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
