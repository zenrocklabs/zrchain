from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Params(_message.Message):
    __slots__ = ("zrchain_relayer_key_id", "solana")
    ZRCHAIN_RELAYER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    SOLANA_FIELD_NUMBER: _ClassVar[int]
    zrchain_relayer_key_id: int
    solana: Solana
    def __init__(self, zrchain_relayer_key_id: _Optional[int] = ..., solana: _Optional[_Union[Solana, _Mapping]] = ...) -> None: ...

class Solana(_message.Message):
    __slots__ = ("program_id", "nonce_account_pub_key", "nonce_authority_pub_key", "rpc_url", "mint_address", "fee_wallet", "fee")
    PROGRAM_ID_FIELD_NUMBER: _ClassVar[int]
    NONCE_ACCOUNT_PUB_KEY_FIELD_NUMBER: _ClassVar[int]
    NONCE_AUTHORITY_PUB_KEY_FIELD_NUMBER: _ClassVar[int]
    RPC_URL_FIELD_NUMBER: _ClassVar[int]
    MINT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    FEE_WALLET_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    program_id: str
    nonce_account_pub_key: str
    nonce_authority_pub_key: str
    rpc_url: str
    mint_address: str
    fee_wallet: str
    fee: int
    def __init__(self, program_id: _Optional[str] = ..., nonce_account_pub_key: _Optional[str] = ..., nonce_authority_pub_key: _Optional[str] = ..., rpc_url: _Optional[str] = ..., mint_address: _Optional[str] = ..., fee_wallet: _Optional[str] = ..., fee: _Optional[int] = ...) -> None: ...
