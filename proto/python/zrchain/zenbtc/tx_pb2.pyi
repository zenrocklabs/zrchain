from amino import amino_pb2 as _amino_pb2
from cosmos.msg.v1 import msg_pb2 as _msg_pb2
from cosmos_proto import cosmos_pb2 as _cosmos_pb2
from gogoproto import gogo_pb2 as _gogo_pb2
from zrchain.zenbtc import params_pb2 as _params_pb2
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

class MsgVerifyDepositBlockInclusion(_message.Message):
    __slots__ = ("creator", "chain_name", "block_height", "raw_tx", "index", "proof", "deposit_addr", "amount")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    CHAIN_NAME_FIELD_NUMBER: _ClassVar[int]
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    RAW_TX_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    PROOF_FIELD_NUMBER: _ClassVar[int]
    DEPOSIT_ADDR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    creator: str
    chain_name: str
    block_height: int
    raw_tx: str
    index: int
    proof: _containers.RepeatedScalarFieldContainer[str]
    deposit_addr: str
    amount: int
    def __init__(self, creator: _Optional[str] = ..., chain_name: _Optional[str] = ..., block_height: _Optional[int] = ..., raw_tx: _Optional[str] = ..., index: _Optional[int] = ..., proof: _Optional[_Iterable[str]] = ..., deposit_addr: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class MsgVerifyDepositBlockInclusionResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class MsgSubmitUnlockTransaction(_message.Message):
    __slots__ = ("creator", "chain", "txID", "withdrawalAddr", "amount")
    CREATOR_FIELD_NUMBER: _ClassVar[int]
    CHAIN_FIELD_NUMBER: _ClassVar[int]
    TXID_FIELD_NUMBER: _ClassVar[int]
    WITHDRAWALADDR_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    creator: str
    chain: str
    txID: str
    withdrawalAddr: str
    amount: int
    def __init__(self, creator: _Optional[str] = ..., chain: _Optional[str] = ..., txID: _Optional[str] = ..., withdrawalAddr: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class MsgSubmitUnlockTransactionResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
