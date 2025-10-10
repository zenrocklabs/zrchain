from amino import amino_pb2 as _amino_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Asset(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ASSET_UNSPECIFIED: _ClassVar[Asset]
    ASSET_ZENBTC: _ClassVar[Asset]
    ASSET_ZENZEC: _ClassVar[Asset]
ASSET_UNSPECIFIED: Asset
ASSET_ZENBTC: Asset
ASSET_ZENZEC: Asset

class Params(_message.Message):
    __slots__ = ("assets",)
    ASSETS_FIELD_NUMBER: _ClassVar[int]
    assets: _containers.RepeatedCompositeFieldContainer[AssetParams]
    def __init__(self, assets: _Optional[_Iterable[_Union[AssetParams, _Mapping]]] = ...) -> None: ...

class AssetParams(_message.Message):
    __slots__ = ("asset", "deposit_keyring_addr", "staker_key_id", "eth_minter_key_id", "unstaker_key_id", "completer_key_id", "rewards_deposit_key_id", "change_address_key_ids", "proxy_address", "eth_token_addr", "controller_addr", "solana")
    ASSET_FIELD_NUMBER: _ClassVar[int]
    DEPOSIT_KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    STAKER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    ETH_MINTER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    UNSTAKER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    COMPLETER_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    REWARDS_DEPOSIT_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    CHANGE_ADDRESS_KEY_IDS_FIELD_NUMBER: _ClassVar[int]
    PROXY_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    ETH_TOKEN_ADDR_FIELD_NUMBER: _ClassVar[int]
    CONTROLLER_ADDR_FIELD_NUMBER: _ClassVar[int]
    SOLANA_FIELD_NUMBER: _ClassVar[int]
    asset: Asset
    deposit_keyring_addr: str
    staker_key_id: int
    eth_minter_key_id: int
    unstaker_key_id: int
    completer_key_id: int
    rewards_deposit_key_id: int
    change_address_key_ids: _containers.RepeatedScalarFieldContainer[int]
    proxy_address: str
    eth_token_addr: str
    controller_addr: str
    solana: Solana
    def __init__(self, asset: _Optional[_Union[Asset, str]] = ..., deposit_keyring_addr: _Optional[str] = ..., staker_key_id: _Optional[int] = ..., eth_minter_key_id: _Optional[int] = ..., unstaker_key_id: _Optional[int] = ..., completer_key_id: _Optional[int] = ..., rewards_deposit_key_id: _Optional[int] = ..., change_address_key_ids: _Optional[_Iterable[int]] = ..., proxy_address: _Optional[str] = ..., eth_token_addr: _Optional[str] = ..., controller_addr: _Optional[str] = ..., solana: _Optional[_Union[Solana, _Mapping]] = ...) -> None: ...

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
