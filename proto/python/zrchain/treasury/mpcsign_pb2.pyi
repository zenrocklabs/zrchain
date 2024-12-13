from google.protobuf import any_pb2 as _any_pb2
from zrchain.treasury import key_pb2 as _key_pb2
from zrchain.treasury import wallet_pb2 as _wallet_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SignRequestStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SIGN_REQUEST_STATUS_UNSPECIFIED: _ClassVar[SignRequestStatus]
    SIGN_REQUEST_STATUS_PENDING: _ClassVar[SignRequestStatus]
    SIGN_REQUEST_STATUS_PARTIAL: _ClassVar[SignRequestStatus]
    SIGN_REQUEST_STATUS_FULFILLED: _ClassVar[SignRequestStatus]
    SIGN_REQUEST_STATUS_REJECTED: _ClassVar[SignRequestStatus]
SIGN_REQUEST_STATUS_UNSPECIFIED: SignRequestStatus
SIGN_REQUEST_STATUS_PENDING: SignRequestStatus
SIGN_REQUEST_STATUS_PARTIAL: SignRequestStatus
SIGN_REQUEST_STATUS_FULFILLED: SignRequestStatus
SIGN_REQUEST_STATUS_REJECTED: SignRequestStatus

class SignRequest(_message.Message):
    __slots__ = ("id", "creator", "key_id", "key_type", "data_for_signing", "status", "signed_data", "keyring_party_signatures", "reject_reason", "metadata", "parent_req_id", "child_req_ids", "cache_id", "key_ids")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    DATA_FOR_SIGNING_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    SIGNED_DATA_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    PARENT_REQ_ID_FIELD_NUMBER: _ClassVar[int]
    CHILD_REQ_IDS_FIELD_NUMBER: _ClassVar[int]
    CACHE_ID_FIELD_NUMBER: _ClassVar[int]
    KEY_IDS_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    key_id: int
    key_type: _key_pb2.KeyType
    data_for_signing: _containers.RepeatedScalarFieldContainer[bytes]
    status: SignRequestStatus
    signed_data: _containers.RepeatedCompositeFieldContainer[SignedDataWithID]
    keyring_party_signatures: _containers.RepeatedScalarFieldContainer[bytes]
    reject_reason: str
    metadata: _any_pb2.Any
    parent_req_id: int
    child_req_ids: _containers.RepeatedScalarFieldContainer[int]
    cache_id: bytes
    key_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ..., key_type: _Optional[_Union[_key_pb2.KeyType, str]] = ..., data_for_signing: _Optional[_Iterable[bytes]] = ..., status: _Optional[_Union[SignRequestStatus, str]] = ..., signed_data: _Optional[_Iterable[_Union[SignedDataWithID, _Mapping]]] = ..., keyring_party_signatures: _Optional[_Iterable[bytes]] = ..., reject_reason: _Optional[str] = ..., metadata: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., parent_req_id: _Optional[int] = ..., child_req_ids: _Optional[_Iterable[int]] = ..., cache_id: _Optional[bytes] = ..., key_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class SignedDataWithID(_message.Message):
    __slots__ = ("sign_request_id", "signed_data")
    SIGN_REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    SIGNED_DATA_FIELD_NUMBER: _ClassVar[int]
    sign_request_id: int
    signed_data: bytes
    def __init__(self, sign_request_id: _Optional[int] = ..., signed_data: _Optional[bytes] = ...) -> None: ...

class SignTransactionRequest(_message.Message):
    __slots__ = ("id", "creator", "key_id", "wallet_type", "unsigned_transaction", "sign_request_id", "no_broadcast", "key_ids")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    UNSIGNED_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    SIGN_REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    NO_BROADCAST_FIELD_NUMBER: _ClassVar[int]
    KEY_IDS_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    key_id: int
    wallet_type: _wallet_pb2.WalletType
    unsigned_transaction: bytes
    sign_request_id: int
    no_broadcast: bool
    key_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., unsigned_transaction: _Optional[bytes] = ..., sign_request_id: _Optional[int] = ..., no_broadcast: bool = ..., key_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class SignReqResponse(_message.Message):
    __slots__ = ("id", "creator", "key_id", "key_type", "data_for_signing", "status", "signed_data", "keyring_party_signatures", "reject_reason", "metadata", "parent_req_id", "child_req_ids", "cache_id", "key_ids")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    DATA_FOR_SIGNING_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    SIGNED_DATA_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    PARENT_REQ_ID_FIELD_NUMBER: _ClassVar[int]
    CHILD_REQ_IDS_FIELD_NUMBER: _ClassVar[int]
    CACHE_ID_FIELD_NUMBER: _ClassVar[int]
    KEY_IDS_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    key_id: int
    key_type: str
    data_for_signing: _containers.RepeatedScalarFieldContainer[bytes]
    status: str
    signed_data: _containers.RepeatedCompositeFieldContainer[SignedDataWithID]
    keyring_party_signatures: _containers.RepeatedScalarFieldContainer[bytes]
    reject_reason: str
    metadata: _any_pb2.Any
    parent_req_id: int
    child_req_ids: _containers.RepeatedScalarFieldContainer[int]
    cache_id: bytes
    key_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ..., key_type: _Optional[str] = ..., data_for_signing: _Optional[_Iterable[bytes]] = ..., status: _Optional[str] = ..., signed_data: _Optional[_Iterable[_Union[SignedDataWithID, _Mapping]]] = ..., keyring_party_signatures: _Optional[_Iterable[bytes]] = ..., reject_reason: _Optional[str] = ..., metadata: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., parent_req_id: _Optional[int] = ..., child_req_ids: _Optional[_Iterable[int]] = ..., cache_id: _Optional[bytes] = ..., key_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class SignTxReqResponse(_message.Message):
    __slots__ = ("id", "creator", "key_id", "wallet_type", "unsigned_transaction", "sign_request_id", "no_broadcast", "key_ids")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    UNSIGNED_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    SIGN_REQUEST_ID_FIELD_NUMBER: _ClassVar[int]
    NO_BROADCAST_FIELD_NUMBER: _ClassVar[int]
    KEY_IDS_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    key_id: int
    wallet_type: str
    unsigned_transaction: bytes
    sign_request_id: int
    no_broadcast: bool
    key_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ..., wallet_type: _Optional[str] = ..., unsigned_transaction: _Optional[bytes] = ..., sign_request_id: _Optional[int] = ..., no_broadcast: bool = ..., key_ids: _Optional[_Iterable[int]] = ...) -> None: ...

class ICATransactionRequest(_message.Message):
    __slots__ = ("id", "creator", "key_id", "key_type", "input_msg", "status", "signed_data", "keyring_party_signatures", "reject_reason", "key_ids")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    KEY_ID_FIELD_NUMBER: _ClassVar[int]
    KEY_TYPE_FIELD_NUMBER: _ClassVar[int]
    INPUT_MSG_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    SIGNED_DATA_FIELD_NUMBER: _ClassVar[int]
    KEYRING_PARTY_SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    KEY_IDS_FIELD_NUMBER: _ClassVar[int]
    id: int
    creator: str
    key_id: int
    key_type: _key_pb2.KeyType
    input_msg: bytes
    status: SignRequestStatus
    signed_data: _containers.RepeatedScalarFieldContainer[bytes]
    keyring_party_signatures: _containers.RepeatedScalarFieldContainer[bytes]
    reject_reason: str
    key_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, id: _Optional[int] = ..., creator: _Optional[str] = ..., key_id: _Optional[int] = ..., key_type: _Optional[_Union[_key_pb2.KeyType, str]] = ..., input_msg: _Optional[bytes] = ..., status: _Optional[_Union[SignRequestStatus, str]] = ..., signed_data: _Optional[_Iterable[bytes]] = ..., keyring_party_signatures: _Optional[_Iterable[bytes]] = ..., reject_reason: _Optional[str] = ..., key_ids: _Optional[_Iterable[int]] = ...) -> None: ...
