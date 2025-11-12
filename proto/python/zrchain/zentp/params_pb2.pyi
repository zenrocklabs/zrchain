from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("solana", "bridge_fee", "flat_fee")
    SOLANA_FIELD_NUMBER: _ClassVar[int]
    BRIDGE_FEE_FIELD_NUMBER: _ClassVar[int]
    FLAT_FEE_FIELD_NUMBER: _ClassVar[int]
    solana: Solana
    bridge_fee: str
    flat_fee: int
    def __init__(self, solana: _Optional[_Union[Solana, _Mapping]] = ..., bridge_fee: _Optional[str] = ..., flat_fee: _Optional[int] = ...) -> None: ...

class Solana(_message.Message):
    __slots__ = ("signer_key_id", "program_id", "nonce_account_key", "nonce_authority_key", "mint_address", "fee_wallet", "fee", "multisig_key_address", "btl", "event_store_program_id")
    SIGNER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    PROGRAM_ID_FIELD_NUMBER: _ClassVar[int]
    NONCE_ACCOUNT_KEY_FIELD_NUMBER: _ClassVar[int]
    NONCE_AUTHORITY_KEY_FIELD_NUMBER: _ClassVar[int]
    MINT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    FEE_WALLET_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    MULTISIG_KEY_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    EVENT_STORE_PROGRAM_ID_FIELD_NUMBER: _ClassVar[int]
    signer_key_id: int
    program_id: str
    nonce_account_key: int
    nonce_authority_key: int
    mint_address: str
    fee_wallet: str
    fee: int
    multisig_key_address: str
    btl: int
    event_store_program_id: str
    def __init__(self, signer_key_id: _Optional[int] = ..., program_id: _Optional[str] = ..., nonce_account_key: _Optional[int] = ..., nonce_authority_key: _Optional[int] = ..., mint_address: _Optional[str] = ..., fee_wallet: _Optional[str] = ..., fee: _Optional[int] = ..., multisig_key_address: _Optional[str] = ..., btl: _Optional[int] = ..., event_store_program_id: _Optional[str] = ...) -> None: ...
