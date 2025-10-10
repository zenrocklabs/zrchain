from amino import amino_pb2 as _amino_pb2
from cosmos.base.query.v1beta1 import pagination_pb2 as _pagination_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from google.api import annotations_pb2 as _annotations_pb2
from zrchain.zenbtc import mint_pb2 as _mint_pb2
from zrchain.zenbtc import params_pb2 as _params_pb2
from zrchain.zenbtc import redemptions_pb2 as _redemptions_pb2
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

class QueryLockTransactionsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QueryLockTransactionsResponse(_message.Message):
    __slots__ = ("lock_transactions",)
    LOCK_TRANSACTIONS_FIELD_NUMBER: _ClassVar[int]
    lock_transactions: _containers.RepeatedCompositeFieldContainer[_mint_pb2.LockTransaction]
    def __init__(self, lock_transactions: _Optional[_Iterable[_Union[_mint_pb2.LockTransaction, _Mapping]]] = ...) -> None: ...

class QueryRedemptionsRequest(_message.Message):
    __slots__ = ("start_index", "status")
    START_INDEX_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    start_index: int
    status: _redemptions_pb2.RedemptionStatus
    def __init__(self, start_index: _Optional[int] = ..., status: _Optional[_Union[_redemptions_pb2.RedemptionStatus, str]] = ...) -> None: ...

class QueryRedemptionsResponse(_message.Message):
    __slots__ = ("redemptions",)
    REDEMPTIONS_FIELD_NUMBER: _ClassVar[int]
    redemptions: _containers.RepeatedCompositeFieldContainer[_redemptions_pb2.Redemption]
    def __init__(self, redemptions: _Optional[_Iterable[_Union[_redemptions_pb2.Redemption, _Mapping]]] = ...) -> None: ...

class QueryPendingMintTransactionsRequest(_message.Message):
    __slots__ = ("start_index", "status")
    START_INDEX_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    start_index: int
    status: _mint_pb2.MintTransactionStatus
    def __init__(self, start_index: _Optional[int] = ..., status: _Optional[_Union[_mint_pb2.MintTransactionStatus, str]] = ...) -> None: ...

class QueryPendingMintTransactionsResponse(_message.Message):
    __slots__ = ("pending_mint_transactions",)
    PENDING_MINT_TRANSACTIONS_FIELD_NUMBER: _ClassVar[int]
    pending_mint_transactions: _containers.RepeatedCompositeFieldContainer[_mint_pb2.PendingMintTransaction]
    def __init__(self, pending_mint_transactions: _Optional[_Iterable[_Union[_mint_pb2.PendingMintTransaction, _Mapping]]] = ...) -> None: ...

class QueryPendingMintTransactionRequest(_message.Message):
    __slots__ = ("tx_hash",)
    TX_HASH_FIELD_NUMBER: _ClassVar[int]
    tx_hash: str
    def __init__(self, tx_hash: _Optional[str] = ...) -> None: ...

class QueryPendingMintTransactionResponse(_message.Message):
    __slots__ = ("pending_mint_transaction",)
    PENDING_MINT_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    pending_mint_transaction: _mint_pb2.PendingMintTransaction
    def __init__(self, pending_mint_transaction: _Optional[_Union[_mint_pb2.PendingMintTransaction, _Mapping]] = ...) -> None: ...

class QuerySupplyRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class QuerySupplyResponse(_message.Message):
    __slots__ = ("custodiedBTC", "totalZenBTC", "mintedZenBTC", "pendingZenBTC", "exchangeRate")
    CUSTODIEDBTC_FIELD_NUMBER: _ClassVar[int]
    TOTALZENBTC_FIELD_NUMBER: _ClassVar[int]
    MINTEDZENBTC_FIELD_NUMBER: _ClassVar[int]
    PENDINGZENBTC_FIELD_NUMBER: _ClassVar[int]
    EXCHANGERATE_FIELD_NUMBER: _ClassVar[int]
    custodiedBTC: int
    totalZenBTC: int
    mintedZenBTC: int
    pendingZenBTC: int
    exchangeRate: str
    def __init__(self, custodiedBTC: _Optional[int] = ..., totalZenBTC: _Optional[int] = ..., mintedZenBTC: _Optional[int] = ..., pendingZenBTC: _Optional[int] = ..., exchangeRate: _Optional[str] = ...) -> None: ...

class QueryBurnEventsRequest(_message.Message):
    __slots__ = ("start_index", "txID", "logIndex", "caip2chainID", "status")
    START_INDEX_FIELD_NUMBER: _ClassVar[int]
    TXID_FIELD_NUMBER: _ClassVar[int]
    LOGINDEX_FIELD_NUMBER: _ClassVar[int]
    CAIP2CHAINID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    start_index: int
    txID: str
    logIndex: int
    caip2chainID: str
    status: _redemptions_pb2.BurnStatus
    def __init__(self, start_index: _Optional[int] = ..., txID: _Optional[str] = ..., logIndex: _Optional[int] = ..., caip2chainID: _Optional[str] = ..., status: _Optional[_Union[_redemptions_pb2.BurnStatus, str]] = ...) -> None: ...

class QueryBurnEventsResponse(_message.Message):
    __slots__ = ("burnEvents",)
    BURNEVENTS_FIELD_NUMBER: _ClassVar[int]
    burnEvents: _containers.RepeatedCompositeFieldContainer[_redemptions_pb2.BurnEvent]
    def __init__(self, burnEvents: _Optional[_Iterable[_Union[_redemptions_pb2.BurnEvent, _Mapping]]] = ...) -> None: ...
