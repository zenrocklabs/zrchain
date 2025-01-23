from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Keyring(_message.Message):
    __slots__ = ("address", "creator", "description", "admins", "parties", "party_threshold", "key_req_fee", "sig_req_fee", "is_active", "delegate_fees", "mpc_default_btl", "mpc_minimum_btl")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    ADMINS_FIELD_NUMBER: _ClassVar[int]
    PARTIES_FIELD_NUMBER: _ClassVar[int]
    PARTY_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    KEY_REQ_FEE_FIELD_NUMBER: _ClassVar[int]
    SIG_REQ_FEE_FIELD_NUMBER: _ClassVar[int]
    IS_ACTIVE_FIELD_NUMBER: _ClassVar[int]
    DELEGATE_FEES_FIELD_NUMBER: _ClassVar[int]
    MPC_DEFAULT_BTL_FIELD_NUMBER: _ClassVar[int]
    MPC_MINIMUM_BTL_FIELD_NUMBER: _ClassVar[int]
    address: str
    creator: str
    description: str
    admins: _containers.RepeatedScalarFieldContainer[str]
    parties: _containers.RepeatedScalarFieldContainer[str]
    party_threshold: int
    key_req_fee: int
    sig_req_fee: int
    is_active: bool
    delegate_fees: bool
    mpc_default_btl: int
    mpc_minimum_btl: int
    def __init__(self, address: _Optional[str] = ..., creator: _Optional[str] = ..., description: _Optional[str] = ..., admins: _Optional[_Iterable[str]] = ..., parties: _Optional[_Iterable[str]] = ..., party_threshold: _Optional[int] = ..., key_req_fee: _Optional[int] = ..., sig_req_fee: _Optional[int] = ..., is_active: bool = ..., delegate_fees: bool = ..., mpc_default_btl: _Optional[int] = ..., mpc_minimum_btl: _Optional[int] = ...) -> None: ...
