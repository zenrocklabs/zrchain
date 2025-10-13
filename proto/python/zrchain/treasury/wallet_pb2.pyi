from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from typing import ClassVar as _ClassVar

DESCRIPTOR: _descriptor.FileDescriptor

class WalletType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    WALLET_TYPE_UNSPECIFIED: _ClassVar[WalletType]
    WALLET_TYPE_NATIVE: _ClassVar[WalletType]
    WALLET_TYPE_EVM: _ClassVar[WalletType]
    WALLET_TYPE_BTC_TESTNET: _ClassVar[WalletType]
    WALLET_TYPE_BTC_MAINNET: _ClassVar[WalletType]
    WALLET_TYPE_BTC_REGNET: _ClassVar[WalletType]
    WALLET_TYPE_SOLANA: _ClassVar[WalletType]
    WALLET_TYPE_ZCASH_MAINNET: _ClassVar[WalletType]
    WALLET_TYPE_ZCASH_TESTNET: _ClassVar[WalletType]
    WALLET_TYPE_ZCASH_REGNET: _ClassVar[WalletType]
WALLET_TYPE_UNSPECIFIED: WalletType
WALLET_TYPE_NATIVE: WalletType
WALLET_TYPE_EVM: WalletType
WALLET_TYPE_BTC_TESTNET: WalletType
WALLET_TYPE_BTC_MAINNET: WalletType
WALLET_TYPE_BTC_REGNET: WalletType
WALLET_TYPE_SOLANA: WalletType
WALLET_TYPE_ZCASH_MAINNET: WalletType
WALLET_TYPE_ZCASH_TESTNET: WalletType
WALLET_TYPE_ZCASH_REGNET: WalletType
