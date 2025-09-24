from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zenex import params_pb2 as _params_pb2
from zrchain.treasury import wallet_pb2 as _wallet_pb2
from zrchain.zenex import swap_pb2 as _swap_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

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

class MsgSwapRequest(_message.Message):
    __slots__ = ("creator", "pair", "workspace", "amount_in", "rock_key_id", "btc_key_id")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    PAIR_FIELD_NUMBER: _ClassVar[int]
    WORKSPACE_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_IN_FIELD_NUMBER: _ClassVar[int]
    ROCK_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    BTC_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    creator: str
    pair: _swap_pb2.TradePair
    workspace: str
    amount_in: int
    rock_key_id: int
    btc_key_id: int
    def __init__(self, creator: _Optional[str] = ..., pair: _Optional[_Union[_swap_pb2.TradePair, str]] = ..., workspace: _Optional[str] = ..., amount_in: _Optional[int] = ..., rock_key_id: _Optional[int] = ..., btc_key_id: _Optional[int] = ...) -> None: ...

class MsgSwapRequestResponse(_message.Message):
    __slots__ = ("swap_id",)
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    swap_id: int
    def __init__(self, swap_id: _Optional[int] = ...) -> None: ...

class MsgZenexTransferRequest(_message.Message):
    __slots__ = ("creator", "swap_id", "data_for_signing", "wallet_type", "cache_id", "unsigned_plus_tx", "reject_reason")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    DATA_FOR_SIGNING_FIELD_NUMBER: _ClassVar[int]
    WALLET_TYPE_FIELD_NUMBER: _ClassVar[int]
    CACHE_ID_FIELD_NUMBER: _ClassVar[int]
    UNSIGNED_PLUS_TX_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    creator: str
    swap_id: int
    data_for_signing: _containers.RepeatedCompositeFieldContainer[_swap_pb2.InputHashes]
    wallet_type: _wallet_pb2.WalletType
    cache_id: bytes
    unsigned_plus_tx: bytes
    reject_reason: str
    def __init__(self, creator: _Optional[str] = ..., swap_id: _Optional[int] = ..., data_for_signing: _Optional[_Iterable[_Union[_swap_pb2.InputHashes, _Mapping]]] = ..., wallet_type: _Optional[_Union[_wallet_pb2.WalletType, str]] = ..., cache_id: _Optional[bytes] = ..., unsigned_plus_tx: _Optional[bytes] = ..., reject_reason: _Optional[str] = ...) -> None: ...

class MsgZenexTransferRequestResponse(_message.Message):
    __slots__ = ("sign_req_id",)
    SIGN_REQ_ID_FIELD_NUMBER: _ClassVar[int]
    sign_req_id: int
    def __init__(self, sign_req_id: _Optional[int] = ...) -> None: ...

class MsgAcknowledgePoolTransfer(_message.Message):
    __slots__ = ("creator", "swap_id", "source_tx_hash", "status", "reject_reason")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SWAP_ID_FIELD_NUMBER: _ClassVar[int]
    SOURCE_TX_HASH_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    REJECT_REASON_FIELD_NUMBER: _ClassVar[int]
    creator: str
    swap_id: int
    source_tx_hash: str
    status: _swap_pb2.SwapStatus
    reject_reason: str
    def __init__(self, creator: _Optional[str] = ..., swap_id: _Optional[int] = ..., source_tx_hash: _Optional[str] = ..., status: _Optional[_Union[_swap_pb2.SwapStatus, str]] = ..., reject_reason: _Optional[str] = ...) -> None: ...

class MsgAcknowledgePoolTransferResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
