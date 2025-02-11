from amino import amino_pb2 as _amino_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from zrchain.treasury import key_pb2 as _key_pb2
from zrchain.treasury import mpcsign_pb2 as _mpcsign_pb2
from zrchain.treasury import params_pb2 as _params_pb2
from zrchain.treasury import wallet_pb2 as _wallet_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryParamsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryParamsResponse(_message.Message):
    __slots__ = ("params",)
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    params: _params_pb2.Params
    def __init__(self, params: _Optional[_Union[_params_pb2.Params, _Mapping]] = ...) -> None: ...

class QueryKeyRequestsRequest(_message.Message):
    __slots__ = ("keyring_addr", "status", "workspace_addr", "pagination")
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    keyring_addr: str
    status: _key_pb2.KeyRequestStatus
    workspace_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, keyring_addr: _Optional[str] = ..., status: _Optional[_Union[_key_pb2.KeyRequestStatus, str]] = ..., workspace_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryKeyRequestsResponse(_message.Message):
    __slots__ = ("key_requests", "pagination")
    KEY_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    key_requests: _containers.RepeatedCompositeFieldContainer[_key_pb2.KeyReqResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, key_requests: _Optional[_Iterable[_Union[_key_pb2.KeyReqResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryKeyByIDRequest(_message.Message):
    __slots__ = ("id", "wallet_type", "prefixes")
    ID_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    PREFIXES_FIELD_NUMBER: _ClassVar[int]
    id: int
    wallet_type: _wallet_pb2.WalletType
    prefixes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, id: _Optional[int] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., prefixes: _Optional[_Iterable[str]] = ...) -> None: ...

class QueryKeyByIDResponse(_message.Message):
    __slots__ = ("key", "wallets")
    KEY_FIELD_NUMBER: _ClassVar[int]
    WALLETS_FIELD_NUMBER: _ClassVar[int]
    key: _key_pb2.KeyResponse
    wallets: _containers.RepeatedCompositeFieldContainer[WalletResponse]
    def __init__(self, key: _Optional[_Union[_key_pb2.KeyResponse, _Mapping]] = ..., wallets: _Optional[_Iterable[_Union[WalletResponse, _Mapping]]] = ...) -> None: ...

class QueryKeysRequest(_message.Message):
    __slots__ = ("workspace_addr", "wallet_type", "prefixes", "pagination")
    WORKSPACE_ADDR_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    PREFIXES_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    workspace_addr: str
    wallet_type: _wallet_pb2.WalletType
    prefixes: _containers.RepeatedScalarFieldContainer[str]
    pagination: _pagination_pb2.PageRequest
    def __init__(self, workspace_addr: _Optional[str] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., prefixes: _Optional[_Iterable[str]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryKeysResponse(_message.Message):
    __slots__ = ("keys", "pagination")
    KEYS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    keys: _containers.RepeatedCompositeFieldContainer[KeyAndWalletResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, keys: _Optional[_Iterable[_Union[KeyAndWalletResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class KeyAndWalletResponse(_message.Message):
    __slots__ = ("key", "wallets")
    KEY_FIELD_NUMBER: _ClassVar[int]
    WALLETS_FIELD_NUMBER: _ClassVar[int]
    key: _key_pb2.KeyResponse
    wallets: _containers.RepeatedCompositeFieldContainer[WalletResponse]
    def __init__(self, key: _Optional[_Union[_key_pb2.KeyResponse, _Mapping]] = ..., wallets: _Optional[_Iterable[_Union[WalletResponse, _Mapping]]] = ...) -> None: ...

class WalletResponse(_message.Message):
    __slots__ = ("address", "type")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    address: str
    type: str
    def __init__(self, address: _Optional[str] = ..., type: _Optional[str] = ...) -> None: ...

class QueryKeyRequestByIDRequest(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class QueryKeyRequestByIDResponse(_message.Message):
    __slots__ = ("key_request",)
    KEY_REQUEST_FIELD_NUMBER: _ClassVar[int]
    key_request: _key_pb2.KeyReqResponse
    def __init__(self, key_request: _Optional[_Union[_key_pb2.KeyReqResponse, _Mapping]] = ...) -> None: ...

class QuerySignatureRequestsRequest(_message.Message):
    __slots__ = ("keyring_addr", "status", "pagination")
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    keyring_addr: str
    status: _mpcsign_pb2.SignRequestStatus
    pagination: _pagination_pb2.PageRequest
    def __init__(self, keyring_addr: _Optional[str] = ..., status: _Optional[_Union[_mpcsign_pb2.SignRequestStatus, str]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QuerySignatureRequestsResponse(_message.Message):
    __slots__ = ("sign_requests", "pagination")
    SIGN_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    sign_requests: _containers.RepeatedCompositeFieldContainer[_mpcsign_pb2.SignReqResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, sign_requests: _Optional[_Iterable[_Union[_mpcsign_pb2.SignReqResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QuerySignatureRequestByIDRequest(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class QuerySignatureRequestByIDResponse(_message.Message):
    __slots__ = ("sign_request",)
    SIGN_REQUEST_FIELD_NUMBER: _ClassVar[int]
    sign_request: _mpcsign_pb2.SignReqResponse
    def __init__(self, sign_request: _Optional[_Union[_mpcsign_pb2.SignReqResponse, _Mapping]] = ...) -> None: ...

class QuerySignTransactionRequestsRequest(_message.Message):
    __slots__ = ("request_id", "key_id", "wallet_type", "status", "pagination")
    REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    request_id: int
    key_id: int
    wallet_type: _wallet_pb2.WalletType
    status: _mpcsign_pb2.SignRequestStatus
    pagination: _pagination_pb2.PageRequest
    def __init__(self, request_id: _Optional[int] = ..., key_id: _Optional[int] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., status: _Optional[_Union[_mpcsign_pb2.SignRequestStatus, str]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QuerySignTransactionRequestsResponse(_message.Message):
    __slots__ = ("sign_transaction_requests", "pagination")
    SIGN_TRANSACTION_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    sign_transaction_requests: _containers.RepeatedCompositeFieldContainer[SignTransactionRequestsResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, sign_transaction_requests: _Optional[_Iterable[_Union[SignTransactionRequestsResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class SignTransactionRequestsResponse(_message.Message):
    __slots__ = ("sign_transaction_requests", "sign_request")
    SIGN_TRANSACTION_REQUESTS_FIELD_NUMBER: _ClassVar[int]
    SIGN_REQUEST_FIELD_NUMBER: _ClassVar[int]
    sign_transaction_requests: _mpcsign_pb2.SignTxReqResponse
    sign_request: _mpcsign_pb2.SignReqResponse
    def __init__(self, sign_transaction_requests: _Optional[_Union[_mpcsign_pb2.SignTxReqResponse, _Mapping]] = ..., sign_request: _Optional[_Union[_mpcsign_pb2.SignReqResponse, _Mapping]] = ...) -> None: ...

class QuerySignTransactionRequestByIDRequest(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class QuerySignTransactionRequestByIDResponse(_message.Message):
    __slots__ = ("sign_transaction_request",)
    SIGN_TRANSACTION_REQUEST_FIELD_NUMBER: _ClassVar[int]
    sign_transaction_request: _mpcsign_pb2.SignTxReqResponse
    def __init__(self, sign_transaction_request: _Optional[_Union[_mpcsign_pb2.SignTxReqResponse, _Mapping]] = ...) -> None: ...

class QueryZrSignKeysRequest(_message.Message):
    __slots__ = ("address", "walletType", "pagination")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    WALLETTYPE_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    address: str
    walletType: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, address: _Optional[str] = ..., walletType: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class ZrSignKeyEntry(_message.Message):
    __slots__ = ("walletType", "address", "index", "id")
    WALLETTYPE_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    walletType: str
    address: str
    index: int
    id: int
    def __init__(self, walletType: _Optional[str] = ..., address: _Optional[str] = ..., index: _Optional[int] = ..., id: _Optional[int] = ...) -> None: ...

class QueryZrSignKeysResponse(_message.Message):
    __slots__ = ("Keys", "pagination")
    KEYS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    Keys: _containers.RepeatedCompositeFieldContainer[ZrSignKeyEntry]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, Keys: _Optional[_Iterable[_Union[ZrSignKeyEntry, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryKeyByAddressRequest(_message.Message):
    __slots__ = ("address", "keyring_addr", "key_type", "wallet_type", "prefixes")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    KEYRING_ADDR_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    PREFIXES_FIELD_NUMBER: _ClassVar[int]
    address: str
    keyring_addr: str
    key_type: _key_pb2.KeyType
    wallet_type: _wallet_pb2.WalletType
    prefixes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, address: _Optional[str] = ..., keyring_addr: _Optional[str] = ..., key_type: _Optional[_Union[_key_pb2.KeyType, str]] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., prefixes: _Optional[_Iterable[str]] = ...) -> None: ...

class QueryKeyByAddressResponse(_message.Message):
    __slots__ = ("response",)
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    response: KeyAndWalletResponse
    def __init__(self, response: _Optional[_Union[KeyAndWalletResponse, _Mapping]] = ...) -> None: ...

class QueryZenbtcWalletsRequest(_message.Message):
    __slots__ = ("recipient_addr", "chain_type", "mint_chain_id", "return_addr", "pagination")
    RECIPIENT_ADDR_FIELD_NUMBER: _ClassVar[int]
    CHAIN_TYPE_FIELD_NUMBER: _ClassVar[int]
    MINT_CHAIN_ID_FIELD_NUMBER: _ClassVar[int]
    RETURN_ADDR_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    recipient_addr: str
    chain_type: _wallet_pb2.WalletType
    mint_chain_id: str
    return_addr: str
    pagination: _pagination_pb2.PageRequest
    def __init__(self, recipient_addr: _Optional[str] = ..., chain_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., mint_chain_id: _Optional[str] = ..., return_addr: _Optional[str] = ..., pagination: _Optional[_Union[_pagination_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class QueryZenbtcWalletsResponse(_message.Message):
    __slots__ = ("zenbtc_wallets", "pagination")
    ZENBTC_WALLETS_FIELD_NUMBER: _ClassVar[int]
    PAGINATION_FIELD_NUMBER: _ClassVar[int]
    zenbtc_wallets: _containers.RepeatedCompositeFieldContainer[KeyAndWalletResponse]
    pagination: _pagination_pb2.PageResponse
    def __init__(self, zenbtc_wallets: _Optional[_Iterable[_Union[KeyAndWalletResponse, _Mapping]]] = ..., pagination: _Optional[_Union[_pagination_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class QueryFeeExcemptsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryFeeExcemptsResponse(_message.Message):
    __slots__ = ("no_fee_msgs",)
    NO_FEE_MSGS_FIELD_NUMBER: _ClassVar[int]
    no_fee_msgs: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, no_fee_msgs: _Optional[_Iterable[str]] = ...) -> None: ...
