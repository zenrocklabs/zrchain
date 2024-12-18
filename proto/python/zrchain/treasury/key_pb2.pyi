from zrchain.treasury import wallet_pb2 as _wallet_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class KeyRequestStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    KEY_REQUEST_STATUS_UNSPECIFIED: _ClassVar[KeyRequestStatus]
    KEY_REQUEST_STATUS_PENDING: _ClassVar[KeyRequestStatus]
    KEY_REQUEST_STATUS_PARTIAL: _ClassVar[KeyRequestStatus]
    KEY_REQUEST_STATUS_FULFILLED: _ClassVar[KeyRequestStatus]
    KEY_REQUEST_STATUS_REJECTED: _ClassVar[KeyRequestStatus]

class KeyType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    KEY_TYPE_UNSPECIFIED: _ClassVar[KeyType]
    KEY_TYPE_ECDSA_SECP256K1: _ClassVar[KeyType]
    KEY_TYPE_EDDSA_ED25519: _ClassVar[KeyType]
    KEY_TYPE_BITCOIN_SECP256K1: _ClassVar[KeyType]
KEY_REQUEST_STATUS_UNSPECIFIED: KeyRequestStatus
KEY_REQUEST_STATUS_PENDING: KeyRequestStatus
KEY_REQUEST_STATUS_PARTIAL: KeyRequestStatus
KEY_REQUEST_STATUS_FULFILLED: KeyRequestStatus
KEY_REQUEST_STATUS_REJECTED: KeyRequestStatus
KEY_TYPE_UNSPECIFIED: KeyType
KEY_TYPE_ECDSA_SECP256K1: KeyType
KEY_TYPE_EDDSA_ED25519: KeyType
KEY_TYPE_BITCOIN_SECP256K1: KeyType

class KeyRequest(_message.Message):
    __slots__ = ("id", "creator", "workspace_addr", "keyring_addr", "key_type", "status", "keyring_party_signatures", "reject_reason", "index", "sign_policy_id", "zenbtc_metadata")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ZENBTC_METADATA_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    workspace_addr: str
    keyring_addr: str
    key_type: KeyType
    status: KeyRequestStatus
    keyring_party_signatures: _containers.RepeatedScalarFieldContainer[bytes]
    reject_reason: str
    index: int
    sign_policy_id: int
    zenbtc_metadata: ZenBTCMetadata
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., workspace_addr: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., key_type: _Optional[_Union[KeyType, str]] = ..., status: _Optional[_Union[KeyRequestStatus, str]] = ..., keyring_party_signatures: _Optional[_Iterable[bytes]] = ..., reject_reason: _Optional[str] = ..., index: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., zenbtc_metadata: _Optional[_Union[ZenBTCMetadata, _Mapping]] = ...) -> None: ...

class KeyReqResponse(_message.Message):
    __slots__ = ("id", "creator", "workspace_addr", "keyring_addr", "key_type", "status", "keyring_party_signatures", "reject_reason", "index", "sign_policy_id", "zenbtc_metadata")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ZENBTC_METADATA_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    workspace_addr: str
    keyring_addr: str
    key_type: str
    status: str
    keyring_party_signatures: _containers.RepeatedScalarFieldContainer[bytes]
    reject_reason: str
    index: int
    sign_policy_id: int
    zenbtc_metadata: ZenBTCMetadata
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., workspace_addr: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., key_type: _Optional[str] = ..., status: _Optional[str] = ..., keyring_party_signatures: _Optional[_Iterable[bytes]] = ..., reject_reason: _Optional[str] = ..., index: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., zenbtc_metadata: _Optional[_Union[ZenBTCMetadata, _Mapping]] = ...) -> None: ...

class Key(_message.Message):
    __slots__ = ("id", "workspace_addr", "keyring_addr", "type", "public_key", "index", "sign_policy_id", "zenbtc_metadata")
    ID_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ZENBTC_METADATA_FIELD_NUMBER: _ClassVar[int]
    id: int
    workspace_addr: str
    keyring_addr: str
    type: KeyType
    public_key: bytes
    index: int
    sign_policy_id: int
    zenbtc_metadata: ZenBTCMetadata
    def __init__(self, id: _Optional[int] = ..., workspace_addr: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., type: _Optional[_Union[KeyType, str]] = ..., public_key: _Optional[bytes] = ..., index: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., zenbtc_metadata: _Optional[_Union[ZenBTCMetadata, _Mapping]] = ...) -> None: ...

class KeyResponse(_message.Message):
    __slots__ = ("id", "workspace_addr", "keyring_addr", "type", "public_key", "index", "sign_policy_id", "zenbtc_metadata")
    ID_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ZENBTC_METADATA_FIELD_NUMBER: _ClassVar[int]
    id: int
    workspace_addr: str
    keyring_addr: str
    type: str
    public_key: bytes
    index: int
    sign_policy_id: int
    zenbtc_metadata: ZenBTCMetadata
    def __init__(self, id: _Optional[int] = ..., workspace_addr: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., type: _Optional[str] = ..., public_key: _Optional[bytes] = ..., index: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., zenbtc_metadata: _Optional[_Union[ZenBTCMetadata, _Mapping]] = ...) -> None: ...

class ZenBTCMetadata(_message.Message):
    __slots__ = ("recipient_addr", "chain_type", "chain_id", "return_address")
    RECIPIENT_ADDR_FIELD_NUMBER: _ClassVar[int]
    CHAIN_TYPE_FIELD_NUMBER: _ClassVar[int]
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    RETURN_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    recipient_addr: str
    chain_type: _wallet_pb2.WalletType
    chain_id: int
    return_address: str
    def __init__(self, recipient_addr: _Optional[str] = ..., chain_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., chain_id: _Optional[int] = ..., return_address: _Optional[str] = ...) -> None: ...

class PendingMintTransaction(_message.Message):
    __slots__ = ("chain_id", "chain_type", "recipient_address", "amount", "creator", "key_id")
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    CHAIN_TYPE_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    chain_id: int
    chain_type: _wallet_pb2.WalletType
    recipient_address: str
    amount: int
    creator: str
    key_id: int
    def __init__(self, chain_id: _Optional[int] = ..., chain_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., recipient_address: _Optional[str] = ..., amount: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ...) -> None: ...

class PendingMintTransactions(_message.Message):
    __slots__ = ("txs",)
    TXS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedCompositeFieldContainer[PendingMintTransaction]
    def __init__(self, txs: _Optional[_Iterable[_Union[PendingMintTransaction, _Mapping]]] = ...) -> None: ...
