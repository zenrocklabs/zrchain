from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Coin(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNSPECIFIED: _ClassVar[Coin]
    ZENBTC: _ClassVar[Coin]
    ROCK: _ClassVar[Coin]
UNSPECIFIED: Coin
ZENBTC: Coin
ROCK: Coin

class LatestBitcoinBlockHeaderRequest(_message.Message):
    __slots__ = ("ChainName",)
    CHAINNAME_FIELD_NUMBER: _ClassVar[int]
    ChainName: str
    def __init__(self, ChainName: _Optional[str] = ...) -> None: ...

class BitcoinBlockHeaderByHeightRequest(_message.Message):
    __slots__ = ("BlockHeight", "ChainName")
    BLOCKHEIGHT_FIELD_NUMBER: _ClassVar[int]
    CHAINNAME_FIELD_NUMBER: _ClassVar[int]
    BlockHeight: int
    ChainName: str
    def __init__(self, BlockHeight: _Optional[int] = ..., ChainName: _Optional[str] = ...) -> None: ...

class BitcoinBlockHeaderResponse(_message.Message):
    __slots__ = ("blockHeader", "BlockHeight", "TipHeight")
    BLOCKHEADER_FIELD_NUMBER: _ClassVar[int]
    BLOCKHEIGHT_FIELD_NUMBER: _ClassVar[int]
    TIPHEIGHT_FIELD_NUMBER: _ClassVar[int]
    blockHeader: BTCBlockHeader
    BlockHeight: int
    TipHeight: int
    def __init__(self, blockHeader: _Optional[_Union[BTCBlockHeader, _Mapping]] = ..., BlockHeight: _Optional[int] = ..., TipHeight: _Optional[int] = ...) -> None: ...

class BTCBlockHeader(_message.Message):
    __slots__ = ("Version", "PrevBlock", "MerkleRoot", "TimeStamp", "Bits", "Nonce", "BlockHash")
    VERSION_FIELD_NUMBER: _ClassVar[int]
    PREVBLOCK_FIELD_NUMBER: _ClassVar[int]
    MERKLEROOT_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    BITS_FIELD_NUMBER: _ClassVar[int]
    NONCE_FIELD_NUMBER: _ClassVar[int]
    BLOCKHASH_FIELD_NUMBER: _ClassVar[int]
    Version: int
    PrevBlock: str
    MerkleRoot: str
    TimeStamp: int
    Bits: int
    Nonce: int
    BlockHash: str
    def __init__(self, Version: _Optional[int] = ..., PrevBlock: _Optional[str] = ..., MerkleRoot: _Optional[str] = ..., TimeStamp: _Optional[int] = ..., Bits: _Optional[int] = ..., Nonce: _Optional[int] = ..., BlockHash: _Optional[str] = ...) -> None: ...

class Redemption(_message.Message):
    __slots__ = ("id", "destination_address", "amount")
    ID_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    id: int
    destination_address: bytes
    amount: int
    def __init__(self, id: _Optional[int] = ..., destination_address: _Optional[bytes] = ..., amount: _Optional[int] = ...) -> None: ...

class BurnEvent(_message.Message):
    __slots__ = ("txID", "logIndex", "chainID", "destinationAddr", "amount", "IsZenBTC")
    TXID_FIELD_NUMBER: _ClassVar[int]
    LOGINDEX_FIELD_NUMBER: _ClassVar[int]
    CHAINID_FIELD_NUMBER: _ClassVar[int]
    DESTINATIONADDR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    ISZENBTC_FIELD_NUMBER: _ClassVar[int]
    txID: str
    logIndex: int
    chainID: str
    destinationAddr: bytes
    amount: int
    IsZenBTC: bool
    def __init__(self, txID: _Optional[str] = ..., logIndex: _Optional[int] = ..., chainID: _Optional[str] = ..., destinationAddr: _Optional[bytes] = ..., amount: _Optional[int] = ..., IsZenBTC: bool = ...) -> None: ...

class SidecarStateRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class SidecarStateResponse(_message.Message):
    __slots__ = ("EigenDelegations", "EthBlockHeight", "EthGasLimit", "EthBaseFee", "EthTipCap", "SolanaLamportsPerSignature", "EthBurnEvents", "Redemptions", "ROCKUSDPrice", "BTCUSDPrice", "ETHUSDPrice", "SolanaMintEvents", "SolanaBurnEvents")
    EIGENDELEGATIONS_FIELD_NUMBER: _ClassVar[int]
    ETHBLOCKHEIGHT_FIELD_NUMBER: _ClassVar[int]
    ETHGASLIMIT_FIELD_NUMBER: _ClassVar[int]
    ETHBASEFEE_FIELD_NUMBER: _ClassVar[int]
    ETHTIPCAP_FIELD_NUMBER: _ClassVar[int]
    SOLANALAMPORTSPERSIGNATURE_FIELD_NUMBER: _ClassVar[int]
    ETHBURNEVENTS_FIELD_NUMBER: _ClassVar[int]
    REDEMPTIONS_FIELD_NUMBER: _ClassVar[int]
    ROCKUSDPRICE_FIELD_NUMBER: _ClassVar[int]
    BTCUSDPRICE_FIELD_NUMBER: _ClassVar[int]
    ETHUSDPRICE_FIELD_NUMBER: _ClassVar[int]
    SOLANAMINTEVENTS_FIELD_NUMBER: _ClassVar[int]
    SOLANABURNEVENTS_FIELD_NUMBER: _ClassVar[int]
    EigenDelegations: bytes
    EthBlockHeight: int
    EthGasLimit: int
    EthBaseFee: int
    EthTipCap: int
    SolanaLamportsPerSignature: int
    EthBurnEvents: _containers.RepeatedCompositeFieldContainer[BurnEvent]
    Redemptions: _containers.RepeatedCompositeFieldContainer[Redemption]
    ROCKUSDPrice: str
    BTCUSDPrice: str
    ETHUSDPrice: str
    SolanaMintEvents: _containers.RepeatedCompositeFieldContainer[SolanaMintEvent]
    SolanaBurnEvents: _containers.RepeatedCompositeFieldContainer[BurnEvent]
    def __init__(self, EigenDelegations: _Optional[bytes] = ..., EthBlockHeight: _Optional[int] = ..., EthGasLimit: _Optional[int] = ..., EthBaseFee: _Optional[int] = ..., EthTipCap: _Optional[int] = ..., SolanaLamportsPerSignature: _Optional[int] = ..., EthBurnEvents: _Optional[_Iterable[_Union[BurnEvent, _Mapping]]] = ..., Redemptions: _Optional[_Iterable[_Union[Redemption, _Mapping]]] = ..., ROCKUSDPrice: _Optional[str] = ..., BTCUSDPrice: _Optional[str] = ..., ETHUSDPrice: _Optional[str] = ..., SolanaMintEvents: _Optional[_Iterable[_Union[SolanaMintEvent, _Mapping]]] = ..., SolanaBurnEvents: _Optional[_Iterable[_Union[BurnEvent, _Mapping]]] = ...) -> None: ...

class SidecarStateByEthHeightRequest(_message.Message):
    __slots__ = ("EthBlockHeight",)
    ETHBLOCKHEIGHT_FIELD_NUMBER: _ClassVar[int]
    EthBlockHeight: int
    def __init__(self, EthBlockHeight: _Optional[int] = ...) -> None: ...

class LatestEthereumNonceForAccountRequest(_message.Message):
    __slots__ = ("Address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    Address: str
    def __init__(self, Address: _Optional[str] = ...) -> None: ...

class LatestEthereumNonceForAccountResponse(_message.Message):
    __slots__ = ("Nonce",)
    NONCE_FIELD_NUMBER: _ClassVar[int]
    Nonce: int
    def __init__(self, Nonce: _Optional[int] = ...) -> None: ...

class SolanaAccountInfoRequest(_message.Message):
    __slots__ = ("PubKey",)
    PUBKEY_FIELD_NUMBER: _ClassVar[int]
    PubKey: str
    def __init__(self, PubKey: _Optional[str] = ...) -> None: ...

class SolanaAccountInfoResponse(_message.Message):
    __slots__ = ("Account",)
    ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    Account: bytes
    def __init__(self, Account: _Optional[bytes] = ...) -> None: ...

class SolanaMintEvent(_message.Message):
    __slots__ = ("Coint", "SigHash", "Recipient", "Date", "Value", "Fee", "Mint", "TxSig")
    COINT_FIELD_NUMBER: _ClassVar[int]
    SIGHASH_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    DATE_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    MINT_FIELD_NUMBER: _ClassVar[int]
    TXSIG_FIELD_NUMBER: _ClassVar[int]
    Coint: Coin
    SigHash: bytes
    Recipient: bytes
    Date: int
    Value: int
    Fee: int
    Mint: bytes
    TxSig: str
    def __init__(self, Coint: _Optional[_Union[Coin, str]] = ..., SigHash: _Optional[bytes] = ..., Recipient: _Optional[bytes] = ..., Date: _Optional[int] = ..., Value: _Optional[int] = ..., Fee: _Optional[int] = ..., Mint: _Optional[bytes] = ..., TxSig: _Optional[str] = ...) -> None: ...
