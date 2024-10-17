from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.protobuf import any_pb2 as _any_pb2
from zrchain.treasury import key_pb2 as _key_pb2
from zrchain.treasury import mpcsign_pb2 as _mpcsign_pb2
from zrchain.treasury import params_pb2 as _params_pb2
from zrchain.treasury import wallet_pb2 as _wallet_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class VerificationVersion(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[VerificationVersion]
    BITCOIN_PLUS: _ClassVar[VerificationVersion]

class SolanaNetworkType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNDEFINED: _ClassVar[SolanaNetworkType]
    MAINNET: _ClassVar[SolanaNetworkType]
    DEVNET: _ClassVar[SolanaNetworkType]
    TESTNET: _ClassVar[SolanaNetworkType]
UNKNOWN: VerificationVersion
BITCOIN_PLUS: VerificationVersion
UNDEFINED: SolanaNetworkType
MAINNET: SolanaNetworkType
DEVNET: SolanaNetworkType
TESTNET: SolanaNetworkType

class MsgUpdateParams(_message.Message):
    __slots__ = ("authority", "params")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    authority: str
    params: _params_pb2.Params
    def __init__(self, authority: _Optional[str] = ..., params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class MsgUpdateParamsResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewKeyRequest(_message.Message):
    __slots__ = ("creator", "workspace_addr", "keyring_addr", "key_type", "btl", "index", "ext_requester", "ext_key_type", "sign_policy_id", "zenbtc_metadata")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    EXT_REQUESTER_FIELD_NUMBER: _ClassVar[int]
    EXT_KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    ZENBTC_METADATA_FIELD_NUMBER: _ClassVar[int]
    creator: str
    workspace_addr: str
    keyring_addr: str
    key_type: str
    btl: int
    index: int
    ext_requester: str
    ext_key_type: int
    sign_policy_id: int
    zenbtc_metadata: _key_pb2.ZenBTCMetdata
    def __init__(self, creator: _Optional[str] = ..., workspace_addr: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., key_type: _Optional[str] = ..., btl: _Optional[int] = ..., index: _Optional[int] = ..., ext_requester: _Optional[str] = ..., ext_key_type: _Optional[int] = ..., sign_policy_id: _Optional[int] = ..., zenbtc_metadata: _Optional[_Union[_key_pb2.ZenBTCMetdata, _Mapping]] = ...) -> None: ...

class MsgNewKeyRequestResponse(_message.Message):
    __slots__ = ("key_req_id",)
    KEY_REQ_ID_FIELD_NUMBER: _ClassVar[int]
    key_req_id: int
    def __init__(self, key_req_id: _Optional[int] = ...) -> None: ...

class MsgFulfilKeyRequest(_message.Message):
    __slots__ = ("creator", "request_id", "status", "key", "reject_reason", "keyring_party_signature")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    KEY_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    creator: str
    request_id: int
    status: _key_pb2.KeyRequestStatus
    key: MsgNewKey
    reject_reason: str
    keyring_party_signature: bytes
    def __init__(self, creator: _Optional[str] = ..., request_id: _Optional[int] = ..., status: _Optional[_Union[_key_pb2.KeyRequestStatus, str]] = ..., key: _Optional[_Union[MsgNewKey, _Mapping]] = ..., reject_reason: _Optional[str] = ..., keyring_party_signature: _Optional[bytes] = ...) -> None: ...

class MsgNewKey(_message.Message):
    __slots__ = ("public_key",)
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    public_key: bytes
    def __init__(self, public_key: _Optional[bytes] = ...) -> None: ...

class MsgFulfilKeyRequestResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewSignatureRequest(_message.Message):
    __slots__ = ("creator", "key_id", "data_for_signing", "btl", "cache_id", "verify_signing_data", "verify_signing_data_version")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    DATA_FOR_SIGNING_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    CACHE_ID_FIELD_NUMBER: _ClassVar[int]
    VERIFY_SIGNING_DATA_FIELD_NUMBER: _ClassVar[int]
    VERIFY_SIGNING_DATA_VERSION_FIELD_NUMBER: _ClassVar[int]
    creator: str
    key_id: int
    data_for_signing: str
    btl: int
    cache_id: bytes
    verify_signing_data: bytes
    verify_signing_data_version: VerificationVersion
    def __init__(self, creator: _Optional[str] = ..., key_id: _Optional[int] = ..., data_for_signing: _Optional[str] = ..., btl: _Optional[int] = ..., cache_id: _Optional[bytes] = ..., verify_signing_data: _Optional[bytes] = ..., verify_signing_data_version: _Optional[_Union[VerificationVersion, str]] = ...) -> None: ...

class MsgNewSignatureRequestResponse(_message.Message):
    __slots__ = ("sig_req_id",)
    SIG_REQ_ID_FIELD_NUMBER: _ClassVar[int]
    sig_req_id: int
    def __init__(self, sig_req_id: _Optional[int] = ...) -> None: ...

class MsgFulfilSignatureRequest(_message.Message):
    __slots__ = ("creator", "request_id", "status", "keyring_party_signature", "signed_data", "reject_reason")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    SIGNED_DATA_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    creator: str
    request_id: int
    status: _mpcsign_pb2.SignRequestStatus
    keyring_party_signature: bytes
    signed_data: bytes
    reject_reason: str
    def __init__(self, creator: _Optional[str] = ..., request_id: _Optional[int] = ..., status: _Optional[_Union[_mpcsign_pb2.SignRequestStatus, str]] = ..., keyring_party_signature: _Optional[bytes] = ..., signed_data: _Optional[bytes] = ..., reject_reason: _Optional[str] = ...) -> None: ...

class MsgFulfilSignatureRequestResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MetadataEthereum(_message.Message):
    __slots__ = ("chain_id",)
    CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    chain_id: int
    def __init__(self, chain_id: _Optional[int] = ...) -> None: ...

class MetadataSolana(_message.Message):
    __slots__ = ("network",)
    NETWORK_FIELD_NUMBER: _ClassVar[int]
    network: SolanaNetworkType
    def __init__(self, network: _Optional[_Union[SolanaNetworkType, str]] = ...) -> None: ...

class MsgNewSignTransactionRequest(_message.Message):
    __slots__ = ("creator", "key_id", "wallet_type", "unsigned_transaction", "metadata", "btl", "cache_id", "no_broadcast")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    UNSIGNED_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    CACHE_ID_FIELD_NUMBER: _ClassVar[int]
    NO_BROADCAST_FIELD_NUMBER: _ClassVar[int]
    creator: str
    key_id: int
    wallet_type: _wallet_pb2.WalletType
    unsigned_transaction: bytes
    metadata: _any_pb2.Any
    btl: int
    cache_id: bytes
    no_broadcast: bool
    def __init__(self, creator: _Optional[str] = ..., key_id: _Optional[int] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., unsigned_transaction: _Optional[bytes] = ..., metadata: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., btl: _Optional[int] = ..., cache_id: _Optional[bytes] = ..., no_broadcast: bool = ...) -> None: ...

class MsgNewSignTransactionRequestResponse(_message.Message):
    __slots__ = ("id", "signature_request_id")
    ID_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    signature_request_id: int
    def __init__(self, id: _Optional[int] = ..., signature_request_id: _Optional[int] = ...) -> None: ...

class MsgTransferFromKeyring(_message.Message):
    __slots__ = ("creator", "keyring", "recipient", "amount", "denom")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEYRING_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    creator: str
    keyring: str
    recipient: str
    amount: int
    denom: str
    def __init__(self, creator: _Optional[str] = ..., keyring: _Optional[str] = ..., recipient: _Optional[str] = ..., amount: _Optional[int] = ..., denom: _Optional[str] = ...) -> None: ...

class MsgTransferFromKeyringResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewICATransactionRequest(_message.Message):
    __slots__ = ("creator", "key_id", "input_payload", "connection_id", "relative_timeout_timestamp", "btl")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    INPUT_PAYLOAD_FIELD_NUMBER: _ClassVar[int]
    CONNECTION_ID_FIELD_NUMBER: _ClassVar[int]
    RELATIVE_TIMEOUT_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    creator: str
    key_id: int
    input_payload: str
    connection_id: str
    relative_timeout_timestamp: int
    btl: int
    def __init__(self, creator: _Optional[str] = ..., key_id: _Optional[int] = ..., input_payload: _Optional[str] = ..., connection_id: _Optional[str] = ..., relative_timeout_timestamp: _Optional[int] = ..., btl: _Optional[int] = ...) -> None: ...

class MsgNewICATransactionRequestResponse(_message.Message):
    __slots__ = ("id", "signature_request_id")
    ID_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    signature_request_id: int
    def __init__(self, id: _Optional[int] = ..., signature_request_id: _Optional[int] = ...) -> None: ...

class MsgFulfilICATransactionRequest(_message.Message):
    __slots__ = ("creator", "request_id", "status", "keyring_party_signature", "signed_data", "reject_reason")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    SIGNED_DATA_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    creator: str
    request_id: int
    status: _mpcsign_pb2.SignRequestStatus
    keyring_party_signature: bytes
    signed_data: bytes
    reject_reason: str
    def __init__(self, creator: _Optional[str] = ..., request_id: _Optional[int] = ..., status: _Optional[_Union[_mpcsign_pb2.SignRequestStatus, str]] = ..., keyring_party_signature: _Optional[bytes] = ..., signed_data: _Optional[bytes] = ..., reject_reason: _Optional[str] = ...) -> None: ...

class MsgFulfilICATransactionRequestResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgNewZrSignSignatureRequest(_message.Message):
    __slots__ = ("creator", "address", "key_type", "wallet_index", "cache_id", "data", "verify_signing_data", "verify_signing_data_version", "wallet_type", "metadata", "no_broadcast", "btl", "tx")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    WALLET_INDEX_FIELD_NUMBER: _ClassVar[int]
    CACHE_ID_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    VERIFY_SIGNING_DATA_FIELD_NUMBER: _ClassVar[int]
    VERIFY_SIGNING_DATA_VERSION_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    NO_BROADCAST_FIELD_NUMBER: _ClassVar[int]
    BTL_FIELD_NUMBER: _ClassVar[int]
    TX_FIELD_NUMBER: _ClassVar[int]
    creator: str
    address: str
    key_type: int
    wallet_index: int
    cache_id: bytes
    data: str
    verify_signing_data: bytes
    verify_signing_data_version: VerificationVersion
    wallet_type: _wallet_pb2.WalletType
    metadata: _any_pb2.Any
    no_broadcast: bool
    btl: int
    tx: bool
    def __init__(self, creator: _Optional[str] = ..., address: _Optional[str] = ..., key_type: _Optional[int] = ..., wallet_index: _Optional[int] = ..., cache_id: _Optional[bytes] = ..., data: _Optional[str] = ..., verify_signing_data: _Optional[bytes] = ..., verify_signing_data_version: _Optional[_Union[VerificationVersion, str]] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., metadata: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., no_broadcast: bool = ..., btl: _Optional[int] = ..., tx: bool = ...) -> None: ...

class MsgNewZrSignSignatureRequestResponse(_message.Message):
    __slots__ = ("req_id",)
    REQ_ID_FIELD_NUMBER: _ClassVar[int]
    req_id: int
    def __init__(self, req_id: _Optional[int] = ...) -> None: ...

class MsgUpdateKeyPolicy(_message.Message):
    __slots__ = ("creator", "key_id", "sign_policy_id")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    SIGN_POLICY_ID_FIELD_NUMBER: _ClassVar[int]
    creator: str
    key_id: int
    sign_policy_id: int
    def __init__(self, creator: _Optional[str] = ..., key_id: _Optional[int] = ..., sign_policy_id: _Optional[int] = ...) -> None: ...

class MsgUpdateKeyPolicyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
