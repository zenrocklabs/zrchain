from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MintTransactionStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    MINT_TRANSACTION_STATUS_UNSPECIFIED: _ClassVar[MintTransactionStatus]
    MINT_TRANSACTION_STATUS_DEPOSITED: _ClassVar[MintTransactionStatus]
    MINT_TRANSACTION_STATUS_STAKED: _ClassVar[MintTransactionStatus]
    MINT_TRANSACTION_STATUS_MINTED: _ClassVar[MintTransactionStatus]

class WalletType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    WALLET_TYPE_UNSPECIFIED: _ClassVar[WalletType]
    WALLET_TYPE_NATIVE: _ClassVar[WalletType]
    WALLET_TYPE_EVM: _ClassVar[WalletType]
    WALLET_TYPE_BTC_TESTNET: _ClassVar[WalletType]
    WALLET_TYPE_BTC_MAINNET: _ClassVar[WalletType]
    WALLET_TYPE_BTC_REGNET: _ClassVar[WalletType]
    WALLET_TYPE_SOLANA: _ClassVar[WalletType]
MINT_TRANSACTION_STATUS_UNSPECIFIED: MintTransactionStatus
MINT_TRANSACTION_STATUS_DEPOSITED: MintTransactionStatus
MINT_TRANSACTION_STATUS_STAKED: MintTransactionStatus
MINT_TRANSACTION_STATUS_MINTED: MintTransactionStatus
WALLET_TYPE_UNSPECIFIED: WalletType
WALLET_TYPE_NATIVE: WalletType
WALLET_TYPE_EVM: WalletType
WALLET_TYPE_BTC_TESTNET: WalletType
WALLET_TYPE_BTC_MAINNET: WalletType
WALLET_TYPE_BTC_REGNET: WalletType
WALLET_TYPE_SOLANA: WalletType

class NonceData(_message.Message):
    __slots__ = ("nonce", "counter", "skip", "prev_nonce")
    NONCE_FIELD_NUMBER: _ClassVar[int]
    COUNTER_FIELD_NUMBER: _ClassVar[int]
    SKIP_FIELD_NUMBER: _ClassVar[int]
    PREV_NONCE_FIELD_NUMBER: _ClassVar[int]
    nonce: int
    counter: int
    skip: bool
    prev_nonce: int
    def __init__(self, nonce: _Optional[int] = ..., counter: _Optional[int] = ..., skip: bool = ..., prev_nonce: _Optional[int] = ...) -> None: ...

class RequestedBitcoinHeaders(_message.Message):
    __slots__ = ("heights",)
    HEIGHTS_FIELD_NUMBER: _ClassVar[int]
    heights: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, heights: _Optional[_Iterable[int]] = ...) -> None: ...

class LockTransaction(_message.Message):
    __slots__ = ("raw_tx", "vout", "sender", "mint_recipient", "amount", "block_height")
    RAW_TX_FIELD_NUMBER: _ClassVar[int]
    VOUT_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    MINT_RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    raw_tx: str
    vout: int
    sender: str
    mint_recipient: str
    amount: int
    block_height: int
    def __init__(self, raw_tx: _Optional[str] = ..., vout: _Optional[int] = ..., sender: _Optional[str] = ..., mint_recipient: _Optional[str] = ..., amount: _Optional[int] = ..., block_height: _Optional[int] = ...) -> None: ...

class PendingMintTransaction(_message.Message):
    __slots__ = ("chain_id", "chain_type", "recipient_address", "amount", "creator", "key_id", "caip2_chain_id", "id", "status", "zrchain_tx_id", "block_height", "awaiting_event_since", "tx_hash")
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    CHAIN_TYPE_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    CAIP2_CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    ZRCHAIN_TX_ID_FIELD_NUMBER: _ClassVar[int]
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    AWAITING_EVENT_SINCE_FIELD_NUMBER: _ClassVar[int]
    TX_HASH_FIELD_NUMBER: _ClassVar[int]
    chain_id: int
    chain_type: WalletType
    recipient_address: str
    amount: int
    creator: str
    key_id: int
    caip2_chain_id: str
    id: int
    status: MintTransactionStatus
    zrchain_tx_id: int
    block_height: int
    awaiting_event_since: int
    tx_hash: str
    def __init__(self, chain_id: _Optional[int] = ..., chain_type: _Optional[_Union[WalletType, str]] = ..., recipient_address: _Optional[str] = ..., amount: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ..., caip2_chain_id: _Optional[str] = ..., id: _Optional[int] = ..., status: _Optional[_Union[MintTransactionStatus, str]] = ..., zrchain_tx_id: _Optional[int] = ..., block_height: _Optional[int] = ..., awaiting_event_since: _Optional[int] = ..., tx_hash: _Optional[str] = ...) -> None: ...

class PendingMintTransactions(_message.Message):
    __slots__ = ("txs",)
    TXS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedCompositeFieldContainer[PendingMintTransaction]
    def __init__(self, txs: _Optional[_Iterable[_Union[PendingMintTransaction, _Mapping]]] = ...) -> None: ...
