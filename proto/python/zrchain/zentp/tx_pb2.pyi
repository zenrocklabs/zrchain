from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zentp import params_pb2 as _params_pb2
from zrchain.zentp import dct_pb2 as _dct_pb2
from cosmos.base.v1beta1 import coin_pb2 as _coin_pb2
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

class MsgBridge(_message.Message):
    __slots__ = ("creator", "source_address", "amount", "denom", "destination_chain", "recipient_address")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    SOURCE_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CHAIN_FIELD_NUMBER: _ClassVar[int]
    RECIPIENT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    creator: str
    source_address: str
    amount: int
    denom: str
    destination_chain: str
    recipient_address: str
    def __init__(self, creator: _Optional[str] = ..., source_address: _Optional[str] = ..., amount: _Optional[int] = ..., denom: _Optional[str] = ..., destination_chain: _Optional[str] = ..., recipient_address: _Optional[str] = ...) -> None: ...

class MsgBridgeResponse(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class MsgBurn(_message.Message):
    __slots__ = ("authority", "module_account", "denom", "amount")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    MODULE_ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    authority: str
    module_account: str
    denom: str
    amount: int
    def __init__(self, authority: _Optional[str] = ..., module_account: _Optional[str] = ..., denom: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class MsgBurnResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSetSolanaROCKSupply(_message.Message):
    __slots__ = ("authority", "amount")
    AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    authority: str
    amount: int
    def __init__(self, authority: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class MsgSetSolanaROCKSupplyResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgInitDct(_message.Message):
    __slots__ = ("creator", "amount", "destination_chain")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    DESTINATION_CHAIN_FIELD_NUMBER: _ClassVar[int]
    creator: str
    amount: _coin_pb2.Coin
    destination_chain: str
    def __init__(self, creator: _Optional[str] = ..., amount: _Optional[_Union[_coin_pb2.Coin, _Mapping]] = ..., destination_chain: _Optional[str] = ...) -> None: ...

class MsgInitDctResponse(_message.Message):
    __slots__ = ("dct",)
    DCT_FIELD_NUMBER: _ClassVar[int]
    dct: _dct_pb2.Dct
    def __init__(self, dct: _Optional[_Union[_dct_pb2.Dct, _Mapping]] = ...) -> None: ...

class MsgInitDctKeys(_message.Message):
    __slots__ = ("creator", "denom", "unsigned_tx")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    DENOM_FIELD_NUMBER: _ClassVar[int]
    UNSIGNED_TX_FIELD_NUMBER: _ClassVar[int]
    creator: str
    denom: str
    unsigned_tx: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, creator: _Optional[str] = ..., denom: _Optional[str] = ..., unsigned_tx: _Optional[_Iterable[bytes]] = ...) -> None: ...

class MsgInitDctKeysResponse(_message.Message):
    __slots__ = ("sign_req_ids",)
    SIGN_REQ_IDS_FIELD_NUMBER: _ClassVar[int]
    sign_req_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, sign_req_ids: _Optional[_Iterable[int]] = ...) -> None: ...
