from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("mpc_keyring", "zr_sign_address", "keyring_commission", "keyring_commission_destination", "min_gas_fee")
    MPC_KEYRING_FIELD_NUMBER: _ClassVar[int]
    ZR_SIGN_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    KEYRING_COMMISSION_FIELD_NUMBER: _ClassVar[int]
    KEYRING_COMMISSION_DESTINATION_FIELD_NUMBER: _ClassVar[int]
    MIN_GAS_FEE_FIELD_NUMBER: _ClassVar[int]
    mpc_keyring: str
    zr_sign_address: str
    keyring_commission: int
    keyring_commission_destination: str
    min_gas_fee: str
    def __init__(self, mpc_keyring: _Optional[str] = ..., zr_sign_address: _Optional[str] = ..., keyring_commission: _Optional[int] = ..., keyring_commission_destination: _Optional[str] = ..., min_gas_fee: _Optional[str] = ...) -> None: ...
